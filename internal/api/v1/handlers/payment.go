package handlers

import (
	"net/http"

	"task-platform-api/internal/services"
	"task-platform-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PrePayRequest 预支付请求
type PrePayRequest struct {
	TaskID    uint64  `json:"task_id" binding:"required"`
	OrderType string  `json:"order_type" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0.01"`
	Remark    string  `json:"remark"`
}

// PrePayResponse 预支付响应
type PrePayResponse struct {
	TradeNo string `json:"trade_no"`
	PayURL  string `json:"pay_url"`
	QRCode  string `json:"qr_code"`
}

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

// PrePay 预支付
// @Summary 预支付
// @Description 创建预支付订单
// @Tags 支付
// @Accept json
// @Produce json
// @Param request body PrePayRequest true "预支付请求"
// @Success 200 {object} utils.ApiResponse{data=PrePayResponse}
// @Failure 400 {object} utils.ApiResponse
// @Failure 401 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /api/v1/pay/prepay [post]
func (h *PaymentHandler) PrePay(c *gin.Context) {
	var req PrePayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 构造预支付请求
	prePayReq := &services.CreatePrePayOrderRequest{
		UserID:    1, // 暂时使用固定用户ID
		TaskID:    req.TaskID,
		OrderType: req.OrderType,
		Amount:    req.Amount,
		Remark:    req.Remark,
		ClientIP:  c.ClientIP(),
	}

	// 调用支付服务
	trade, prePayResp, err := h.paymentService.CreatePrePayOrder(c.Request.Context(), prePayReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建预支付订单失败")
		return
	}

	// 构造响应
	resp := &PrePayResponse{
		TradeNo: prePayResp.TradeNo,
		PayURL:  prePayResp.PayURL,
		QRCode:  prePayResp.QRCode,
	}

	utils.SuccessResponse(c, gin.H{
		"trade": trade,
		"pay":   resp,
	})
}

// QueryStatus 查询支付状态
// @Summary 查询支付状态
// @Description 查询订单支付状态
// @Tags 支付
// @Accept json
// @Produce json
// @Param order_no path string true "订单号"
// @Success 200 {object} utils.ApiResponse{data=models.Trade}
// @Failure 400 {object} utils.ApiResponse
// @Failure 404 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /api/v1/pay/status/{order_no} [get]
func (h *PaymentHandler) QueryStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "订单号不能为空")
		return
	}

	trade, err := h.paymentService.QueryPaymentStatus(c.Request.Context(), orderNo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "查询支付状态失败")
		return
	}

	utils.SuccessResponse(c, trade)
}

// PaymentCallback 支付回调
// @Summary 支付回调
// @Description 处理支付回调通知
// @Tags 支付
// @Accept json
// @Produce json
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /api/v1/pay/callback [post]
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	// 解析回调数据
	data := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	err := h.paymentService.ProcessPaymentCallback(c.Request.Context(), data)
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}