package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"task-platform-api/internal/config"
	"task-platform-api/internal/models"
	"task-platform-api/pkg/payment"
	"task-platform-api/pkg/utils"
)

// PaymentService 支付服务
type PaymentService struct {
	db        *gorm.DB
	sqbClient *payment.ShouqianbaClient
	config    *config.Config
}

// NewPaymentService 创建支付服务
func NewPaymentService(db *gorm.DB, cfg *config.Config) *PaymentService {
	return &PaymentService{
		db:        db,
		sqbClient: payment.NewShouqianbaClient(&cfg.Shouqianba),
		config:    cfg,
	}
}

// OrderType 订单类型
type OrderType string

const (
	OrderTypeTaskPublish OrderType = "task_publish" // 任务发布预缴
	OrderTypeTaskTake    OrderType = "task_take"    // 任务接取预缴
	OrderTypeDeposit     OrderType = "deposit"      // 保证金
	OrderTypeServiceFee  OrderType = "service_fee"  // 服务费
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"    // 待支付
	PaymentStatusPaid      PaymentStatus = "paid"       // 已支付
	PaymentStatusFailed    PaymentStatus = "failed"     // 支付失败
	PaymentStatusRefunded  PaymentStatus = "refunded"   // 已退款
	PaymentStatusSettled   PaymentStatus = "settled"    // 已结算
)

// CreatePrePayOrder 创建预支付订单
func (s *PaymentService) CreatePrePayOrder(ctx context.Context, req *CreatePrePayOrderRequest) (*models.Trade, *payment.PrePayResponseData, error) {
	// 生成订单号
	orderNo := utils.GenerateOrderNo()

	// 创建交易记录
	trade := &models.Trade{
		OrderNo:     orderNo,
		UserID:      req.UserID,
		TaskID:      req.TaskID,
		TradeType:   string(req.OrderType),
		Amount:      req.Amount,
		Status:      string(PaymentStatusPending),
		PayTime:     nil,
		Remark:      req.Remark,
		ExpiredAt:   time.Now().Add(15 * time.Minute), // 15分钟过期
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存交易记录
	if err := s.db.WithContext(ctx).Create(trade).Error; err != nil {
		return nil, nil, fmt.Errorf("创建交易记录失败: %w", err)
	}

	// 构建预支付请求
	prePayReq := &payment.PrePayRequest{
		OrderNo:     orderNo,
		Amount:      req.Amount,
		Subject:     req.Subject,
		Description: req.Description,
		NotifyURL:   s.getServerURL() + "/api/v1/pay/callback",
		ReturnURL:   req.ReturnURL,
		ExpireTime:  900, // 15分钟
		ClientIP:    req.ClientIP,
		Extra:       fmt.Sprintf("user_id=%d&task_id=%d", req.UserID, req.TaskID),
	}

	// 调用收钱吧预支付接口
	prePayResp, err := s.sqbClient.PrePay(prePayReq)
	if err != nil {
		// 更新交易状态为失败
		trade.Status = string(PaymentStatusFailed)
		s.db.WithContext(ctx).Save(trade)
		return nil, nil, fmt.Errorf("预支付失败: %w", err)
	}

	// 更新交易记录
	trade.TradeNo = prePayResp.TradeNo
	s.db.WithContext(ctx).Save(trade)

	return trade, prePayResp, nil
}

// ProcessPaymentCallback 处理支付回调
func (s *PaymentService) ProcessPaymentCallback(ctx context.Context, data map[string]string) error {
	// 验证签名
	if !s.sqbClient.VerifyNotification(data) {
		return errors.New("签名验证失败")
	}

	orderNo := data["order_no"]
	tradeNo := data["trade_no"]
	status := data["status"]
	amountStr := data["amount"]
	payTimeStr := data["pay_time"]
	payMethod := data["pay_method"]
	transactionID := data["transaction_id"]

	// 查询交易记录
	var trade models.Trade
	if err := s.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&trade).Error; err != nil {
		return fmt.Errorf("交易记录不存在: %w", err)
	}

	// 检查交易状态
	if trade.Status == string(PaymentStatusPaid) {
		return nil // 已处理，直接返回
	}

	// 解析金额
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("金额格式错误: %w", err)
	}

	// 验证金额
	if amount != trade.Amount {
		return fmt.Errorf("金额不匹配: 预期%.2f, 实际%.2f", trade.Amount, amount)
	}

	// 解析支付时间
	payTime, err := time.Parse("2006-01-02 15:04:05", payTimeStr)
	if err != nil {
		return fmt.Errorf("支付时间格式错误: %w", err)
	}

	// 更新交易状态
	updates := map[string]interface{}{
		"trade_no":      tradeNo,
		"status":        string(PaymentStatusPaid),
		"pay_time":      &payTime,
		"pay_method":    payMethod,
		"transaction_id": transactionID,
		"updated_at":    time.Now(),
	}

	if err := s.db.WithContext(ctx).Model(&trade).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新交易记录失败: %w", err)
	}

	// 根据订单类型处理后续逻辑
	if err := s.handlePaidOrder(ctx, &trade); err != nil {
		return fmt.Errorf("处理已支付订单失败: %w", err)
	}

	return nil
}

// handlePaidOrder 处理已支付订单的后续逻辑
func (s *PaymentService) handlePaidOrder(ctx context.Context, trade *models.Trade) error {
	switch trade.TradeType {
	case string(OrderTypeTaskPublish):
		return s.handleTaskPublishPaid(ctx, trade)
	case string(OrderTypeTaskTake):
		return s.handleTaskTakePaid(ctx, trade)
	default:
		return nil
	}
}

// handleTaskPublishPaid 处理任务发布支付完成
func (s *PaymentService) handleTaskPublishPaid(ctx context.Context, trade *models.Trade) error {
	// 查询任务
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ?", trade.TaskID).First(&task).Error; err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 更新任务状态为待审核
	if err := s.db.WithContext(ctx).Model(&task).Update("status", "pending_audit").Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	// 记录操作日志
	log := &models.TaskLog{
		TaskID:     task.ID,
		UserID:     task.PublisherID,
		Action:     "pay_publish",
		Content:    "支付任务发布费用",
		OldStatus:  task.Status,
		NewStatus:  "pending_audit",
		CreatedAt:  time.Now(),
	}

	return s.db.WithContext(ctx).Create(log).Error
}

// handleTaskTakePaid 处理任务接取支付完成
func (s *PaymentService) handleTaskTakePaid(ctx context.Context, trade *models.Trade) error {
	// 查询任务
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ?", trade.TaskID).First(&task).Error; err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 开始事务
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新任务接取信息
	updates := map[string]interface{}{
		"taker_id":   trade.UserID,
		"status":     "in_progress",
		"take_time":  time.Now(),
		"updated_at": time.Now(),
	}

	if err := tx.Model(&task).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新任务信息失败: %w", err)
	}

	// 记录操作日志
	log := &models.TaskLog{
		TaskID:     task.ID,
		UserID:     trade.UserID,
		Action:     "take_task",
		Content:    "接取任务",
		OldStatus:  "open",
		NewStatus:  "in_progress",
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(log).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录操作日志失败: %w", err)
	}

	// 提交事务
	return tx.Commit().Error
}

// RefundOrder 退款订单
func (s *PaymentService) RefundOrder(ctx context.Context, req *RefundOrderRequest) (*models.Refund, error) {
	// 查询原交易记录
	var trade models.Trade
	if err := s.db.WithContext(ctx).Where("order_no = ?", req.OrderNo).First(&trade).Error; err != nil {
		return nil, fmt.Errorf("查询原交易记录失败: %w", err)
	}

	// 检查交易状态
	if trade.Status != string(PaymentStatusPaid) {
		return nil, errors.New("只有已支付的订单才能退款")
	}

	// 生成退款单号
	refundNo := utils.GenerateOrderNo()

	// 创建退款记录
	refund := &models.Refund{
		RefundNo:   refundNo,
		TradeID:    trade.ID,
		OrderNo:    trade.OrderNo,
		UserID:     trade.UserID,
		TaskID:     trade.TaskID,
		Amount:     req.Amount,
		Reason:     req.Reason,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存退款记录
	if err := s.db.WithContext(ctx).Create(refund).Error; err != nil {
		return nil, fmt.Errorf("创建退款记录失败: %w", err)
	}

	// 调用收钱吧退款接口
	refundReq := &payment.RefundRequest{
		OrderNo:   trade.OrderNo,
		RefundNo:  refundNo,
		Amount:    req.Amount,
		Reason:    req.Reason,
		NotifyURL: s.getServerURL() + "/api/v1/pay/refund_callback",
	}

	refundResp, err := s.sqbClient.Refund(refundReq)
	if err != nil {
		// 更新退款状态为失败
		refund.Status = "failed"
		s.db.WithContext(ctx).Save(refund)
		return nil, fmt.Errorf("退款请求失败: %w", err)
	}

	// 更新退款记录
	refund.Status = refundResp.Status
	refund.RefundTime = &refundResp.RefundTime
	s.db.WithContext(ctx).Save(refund)

	return refund, nil
}

// SettlementTask 任务结算
func (s *PaymentService) SettlementTask(ctx context.Context, taskID uint64, publisherAmount, takerAmount float64) error {
	// 查询任务
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error; err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 查询接取方信息
	var taker models.User
	if err := s.db.WithContext(ctx).Where("id = ?", task.TakerID).First(&taker).Error; err != nil {
		return fmt.Errorf("查询接取方信息失败: %w", err)
	}

	// 开始事务
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建结算记录
	settlement := &models.Settlement{
		TaskID:         taskID,
		PublisherID:    task.PublisherID,
		TakerID:        task.TakerID,
		PublisherAmount: publisherAmount,
		TakerAmount:    takerAmount,
		PlatformFee:    task.Amount - publisherAmount - takerAmount,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := tx.Create(settlement).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建结算记录失败: %w", err)
	}

	// 转账给接取方
	transferNo := utils.GenerateOrderNo()
	transferReq := &payment.TransferRequest{
		OrderNo:   transferNo,
		AccountNo: taker.AlipayAccount, // 假设有支付宝账户字段
		Amount:    takerAmount,
		RealName:  taker.RealName,
		Memo:      fmt.Sprintf("任务结算: %s", task.Title),
		NotifyURL: s.getServerURL() + "/api/v1/pay/transfer_callback",
	}

	transferResp, err := s.sqbClient.Transfer(transferReq)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("转账失败: %w", err)
	}

	// 更新结算记录
	settlement.Status = "completed"
	settlement.SettleTime = &transferResp.TransferTime
	settlement.TransferNo = transferResp.TransferNo
	tx.Save(settlement)

	// 更新任务状态
	if err := tx.Model(&task).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	// 记录操作日志
	log := &models.TaskLog{
		TaskID:     taskID,
		UserID:     task.PublisherID,
		Action:     "settlement",
		Content:    fmt.Sprintf("任务结算完成: 发布方%.2f, 接取方%.2f", publisherAmount, takerAmount),
		OldStatus:  "completed",
		NewStatus:  "completed",
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(log).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录操作日志失败: %w", err)
	}

	// 提交事务
	return tx.Commit().Error
}

// QueryTradeStatus 查询交易状态
func (s *PaymentService) QueryTradeStatus(ctx context.Context, orderNo string) (*models.Trade, error) {
	var trade models.Trade
	err := s.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&trade).Error
	if err != nil {
		return nil, fmt.Errorf("查询交易记录失败: %w", err)
	}

	// 如果是待支付状态，向收钱吧查询最新状态
	if trade.Status == string(PaymentStatusPending) {
		req := &payment.PayStatusRequest{
			OrderNo: orderNo,
			TradeNo: trade.TradeNo,
		}

		payStatus, err := s.sqbClient.QueryPayStatus(req)
		if err == nil && payStatus.Status == "SUCCESS" {
			// 更新本地状态
			trade.Status = string(PaymentStatusPaid)
			trade.PayTime = &payStatus.PayTime
			trade.PayMethod = payStatus.PayMethod
			trade.TransactionID = payStatus.TransactionID
			s.db.WithContext(ctx).Save(&trade)

			// 处理后续逻辑
			s.handlePaidOrder(ctx, &trade)
		}
	}

	return &trade, nil
}

// getServerURL 获取服务器URL
func (s *PaymentService) getServerURL() string {
	if s.config.Server.IsHTTPS {
		return "https://" + s.config.Server.Domain
	}
	return "http://" + s.config.Server.Domain
}

// Request types
type CreatePrePayOrderRequest struct {
	UserID      uint64   `json:"user_id"`
	TaskID      uint64   `json:"task_id"`
	OrderType   OrderType `json:"order_type"`
	Amount      float64  `json:"amount"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	ReturnURL   string   `json:"return_url"`
	ClientIP    string   `json:"client_ip"`
	Remark      string   `json:"remark"`
}

type RefundOrderRequest struct {
	OrderNo string  `json:"order_no"`
	Amount  float64 `json:"amount"`
	Reason  string  `json:"reason"`
}