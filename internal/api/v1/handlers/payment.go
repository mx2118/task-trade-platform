package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-platform-api/internal/api/v1/middleware"
	"task-platform-api/internal/services"
	"task-platform-api/pkg/utils"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler 创建支付处理器
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// PrePayRequest 预支付请求
type PrePayRequest struct {
	TaskID      uint64  `json:"task_id" binding:"required"`
	OrderType   string  `json:"order_type" binding:"required,oneof=task_publish task_take"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	ReturnURL   string  `json:"return_url"`
	ClientIP    string  `json:"client_ip"`
	Remark      string  `json:"remark"`
}

// PrePayResponse 预支付响应
type PrePayResponse struct {
	OrderNo string  `json:"order_no"`
	TradeNo string  `json:"trade_no"`
	PayURL  string  `json:"pay_url"`
	QRCode  string  `json:"qrcode"`
	Amount  float64 `json:"amount"`
	ExpireTime int64 `json:"expire_time"`
}

// PrePay 预支付
// @Summary 预支付
// @Description 创建预支付订单
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body PrePayRequest true "预支付请求"
// @Success 200 {object} utils.Response{data=PrePayResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/prepay [post]
func (h *PaymentHandler) PrePay(c *gin.Context) {
	var req PrePayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "")
		return
	}

	// 构建预支付请求
	prePayReq := &services.CreatePrePayOrderRequest{
		UserID:    userID,
		TaskID:    req.TaskID,
		OrderType: services.OrderType(req.OrderType),
		Amount:    req.Amount,
		ReturnURL: req.ReturnURL,
		ClientIP:  utils.GetClientIP(c),
		Remark:    req.Remark,
	}

	// 根据订单类型设置主题和描述
	switch req.OrderType {
	case "task_publish":
		prePayReq.Subject = "任务发布费用"
		prePayReq.Description = "任务发布预缴费用，包含任务本金和服务费"
	case "task_take":
		prePayReq.Subject = "任务接取费用"
		prePayReq.Description = "任务接取预缴费用，包含保证金和服务费"
	}

	// 调用支付服务
	trade, prePayResp, err := h.paymentService.CreatePrePayOrder(c.Request.Context(), prePayReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建预支付订单失败", err.Error())
		return
	}

	// 构建响应
	resp := &PrePayResponse{
		OrderNo:    trade.OrderNo,
		TradeNo:    prePayResp.TradeNo,
		PayURL:     prePayResp.PayURL,
		QRCode:     prePayResp.QRCode,
		Amount:     trade.Amount,
		ExpireTime: prePayResp.ExpireTime,
	}

	utils.SuccessResponse(c, resp)
}

// CallbackRequest 回调请求
type CallbackRequest map[string]string

// Callback 支付回调
// @Summary 支付回调
// @Description 接收收钱吧支付回调
// @Tags 支付
// @Accept x-www-form-urlencoded
// @Produce plain
// @Param signature formData string true "签名"
// @Param order_no formData string true "订单号"
// @Param trade_no formData string true "交易号"
// @Param status formData string true "状态"
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /api/v1/pay/callback [post]
func (h *PaymentHandler) Callback(c *gin.Context) {
	// 解析表单数据
	data := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	// 处理回调
	err := h.paymentService.ProcessPaymentCallback(c.Request.Context(), data)
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// RefundCallbackRequest 退款回调
type RefundCallbackRequest map[string]string

// RefundCallback 退款回调
// @Summary 退款回调
// @Description 接收收钱吧退款回调
// @Tags 支付
// @Accept x-www-form-urlencoded
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /api/v1/pay/refund_callback [post]
func (h *PaymentHandler) RefundCallback(c *gin.Context) {
	// 解析表单数据
	data := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	// 验证签名 TODO: 实现退款回调处理逻辑

	c.String(http.StatusOK, "success")
}

// TransferCallbackRequest 转账回调
type TransferCallbackRequest map[string]string

// TransferCallback 转账回调
// @Summary 转账回调
// @Description 接收收钱吧转账回调
// @Tags 支付
// @Accept x-www-form-urlencoded
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /api/v1/pay/transfer_callback [post]
func (h *PaymentHandler) TransferCallback(c *gin.Context) {
	// 解析表单数据
	data := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	// 验证签名 TODO: 实现转账回调处理逻辑

	c.String(http.StatusOK, "success")
}

// RefundRequest 退款请求
type RefundRequest struct {
	OrderNo string  `json:"order_no" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Reason  string  `json:"reason" binding:"required"`
}

// Refund 退款
// @Summary 退款
// @Description 申请退款
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body RefundRequest true "退款请求"
// @Success 200 {object} utils.Response{data=models.Refund}
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/refund [post]
func (h *PaymentHandler) Refund(c *gin.Context) {
	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "")
		return
	}

	// 构建退款请求
	refundReq := &services.RefundOrderRequest{
		OrderNo: req.OrderNo,
		Amount:  req.Amount,
		Reason:  req.Reason,
	}

	// 调用退款服务
	refund, err := h.paymentService.RefundOrder(c.Request.Context(), refundReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "退款失败", err.Error())
		return
	}

	utils.SuccessResponse(c, refund)
}

// QueryStatusRequest 查询状态请求
type QueryStatusRequest struct {
	OrderNo string `form:"order_no" binding:"required"`
}

// QueryStatus 查询交易状态
// @Summary 查询交易状态
// @Description 查询订单交易状态
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param order_no query string true "订单号"
// @Success 200 {object} utils.Response{data=models.Trade}
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/query [get]
func (h *PaymentHandler) QueryStatus(c *gin.Context) {
	var req QueryStatusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 查询交易状态
	trade, err := h.paymentService.QueryTradeStatus(c.Request.Context(), req.OrderNo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	utils.SuccessResponse(c, trade)
}

// SettlementRequest 结算请求
type SettlementRequest struct {
	TaskID         uint64  `json:"task_id" binding:"required"`
	PublisherAmount float64 `json:"publisher_amount" binding:"required,gt=0"`
	TakerAmount    float64 `json:"taker_amount" binding:"required,gt=0"`
}

// Settlement 任务结算
// @Summary 任务结算
// @Description 完成任务结算
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body SettlementRequest true "结算请求"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/settlement [post]
func (h *PaymentHandler) Settlement(c *gin.Context) {
	var req SettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取当前用户ID（需要验证是否为任务发布方）
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "")
		return
	}

	// 调用结算服务
	err := h.paymentService.SettlementTask(c.Request.Context(), req.TaskID, req.PublisherAmount, req.TakerAmount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "结算失败", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "结算成功"})
}

// UserBalanceRequest 用户余额查询请求
type UserBalanceRequest struct{}

// UserBalanceResponse 用户余额响应
type UserBalanceResponse struct {
	UserID       uint64  `json:"user_id"`
	TotalIncome  float64 `json:"total_income"`  // 总收入
	TotalExpense float64 `json:"total_expense"` // 总支出
	AvailableBalance float64 `json:"available_balance"` // 可用余额
	FrozenBalance   float64 `json:"frozen_balance"`    // 冻结余额
}

// GetUserBalance 获取用户余额
// @Summary 获取用户余额
// @Description 获取用户账户余额信息
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response{data=UserBalanceResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/balance [get]
func (h *PaymentHandler) GetUserBalance(c *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "")
		return
	}

	// TODO: 实现用户余额查询逻辑
	// 这里返回模拟数据
	resp := &UserBalanceResponse{
		UserID:          userID,
		TotalIncome:     1000.00,
		TotalExpense:    500.00,
		AvailableBalance: 500.00,
		FrozenBalance:   0.00,
	}

	utils.SuccessResponse(c, resp)
}

// TransactionListRequest 交易列表请求
type TransactionListRequest struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
	TradeType string `form:"trade_type"`
	Status    string `form:"status"`
}

// TransactionListResponse 交易列表响应
type TransactionListResponse struct {
	List     []models.Trade `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// GetTransactionList 获取交易列表
// @Summary 获取交易列表
// @Description 获取用户交易记录列表
// @Tags 支付
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码"
// @Param page_size query int false "页大小"
// @Param trade_type query string false "交易类型"
// @Param status query string false "交易状态"
// @Success 200 {object} utils.Response{data=TransactionListResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/pay/transactions [get]
func (h *PaymentHandler) GetTransactionList(c *gin.Context) {
	var req TransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "")
		return
	}

	// TODO: 实现交易列表查询逻辑
	// 这里返回空列表
	resp := &TransactionListResponse{
		List:     []models.Trade{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.SuccessResponse(c, resp)
}