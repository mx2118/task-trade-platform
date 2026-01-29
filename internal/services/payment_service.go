package services

import (
	"context"
	"fmt"
	"time"

	"task-platform-api/internal/models"
	"task-platform-api/pkg/payment"
	"task-platform-api/pkg/utils"
)

// CreatePrePayOrderRequest 预支付订单请求
type CreatePrePayOrderRequest struct {
	UserID      uint64 `json:"user_id"`
	TaskID      uint64 `json:"task_id"`
	OrderType   string `json:"order_type"`
	Amount      float64 `json:"amount"`
	Remark      string `json:"remark"`
	ClientIP    string `json:"client_ip"`
}

// PaymentService 支付服务
type PaymentService struct {
	db         interface{} // 临时使用interface{}避免GORM导入问题
	sqbClient *payment.ShouqianbaClient
}

// NewPaymentService 创建支付服务
func NewPaymentService(db interface{}, sqbClient *payment.ShouqianbaClient) *PaymentService {
	return &PaymentService{
		db:         db,
		sqbClient: sqbClient,
	}
}

// CreatePrePayOrder 创建预支付订单
func (s *PaymentService) CreatePrePayOrder(ctx context.Context, req *CreatePrePayOrderRequest) (*models.Trade, *payment.PrePayResponseData, error) {
	// 生成订单号
	orderNo := utils.GenerateOrderNo()

	// 创建交易记录
	taskIDPtr := &req.TaskID
	expireTime := time.Now().Add(15 * time.Minute)
	trade := &models.Trade{
		InternalNo:   orderNo,
		UserID:       req.UserID,
		TaskID:       taskIDPtr,
		TradeType:    string(req.OrderType),
		Amount:       req.Amount,
		Status:       0, // 待支付
		PayTime:      nil,
		Description:  req.Remark,
		ExpireTime:   &expireTime,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 保存交易记录 - 暂时跳过实际数据库操作
	// if err := s.db.WithContext(ctx).Create(trade).Error; err != nil {
	//     return nil, nil, fmt.Errorf("创建交易记录失败: %w", err)
	// }

// 调用收钱吧预支付接口 - 暂时返回模拟数据
	prePayResp := &payment.PrePayResponseData{
		TradeNo: "mock_trade_no_" + orderNo,
		PayURL:  "https://mock.pay.url/" + orderNo,
		QRCode:  "mock_qr_code_" + orderNo,
	}

	// 更新交易记录
	trade.ThirdPartyNo = prePayResp.TradeNo
	// s.db.WithContext(ctx).Save(trade)

	return trade, prePayResp, nil
}

// ProcessPaymentCallback 处理支付回调
func (s *PaymentService) ProcessPaymentCallback(ctx context.Context, data map[string]string) error {
	// 验证签名
	if !s.sqbClient.VerifyNotification(data) {
		return fmt.Errorf("签名验证失败")
	}

	tradeNo := data["trade_no"]
	amountStr := data["amount"]

	// 查询交易记录 - 暂时跳过
	/*
		var trade models.Trade
		if err := s.db.WithContext(ctx).Where("internal_no = ?", orderNo).First(&trade).Error; err != nil {
			return fmt.Errorf("交易记录不存在: %w", err)
		}

		// 检查交易状态
		if trade.Status == 1 { // 1=已支付
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
	*/

	_ = tradeNo
	_ = amountStr

	return fmt.Errorf("支付回调处理暂未实现")
}

// QueryPaymentStatus 查询支付状态
func (s *PaymentService) QueryPaymentStatus(ctx context.Context, orderNo string) (*models.Trade, error) {
	// 暂时返回模拟数据
	trade := &models.Trade{
		ID:          1,
		InternalNo:  orderNo,
		Status:      0, // 待支付
		Amount:      100.00,
		TradeType:   "prepay",
		Description: "模拟交易",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return trade, nil
}

// CreateRefund 创建退款
func (s *PaymentService) CreateRefund(ctx context.Context, tradeID uint64, reason string) error {
	// 暂时不实现
	return fmt.Errorf("退款功能暂未实现")
}

// ProcessRefundCallback 处理退款回调
func (s *PaymentService) ProcessRefundCallback(ctx context.Context, data map[string]string) error {
	// 暂时不实现
	return fmt.Errorf("退款回调处理暂未实现")
}