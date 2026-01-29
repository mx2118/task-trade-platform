package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"task-platform-api/internal/config"
	"task-platform-api/pkg/utils"
)

// ShouqianbaClient 收钱吧客户端
type ShouqianbaClient struct {
	config *config.ShouqianbaConfig
	client *http.Client
}

// PrePayRequest 预支付请求
type PrePayRequest struct {
	OrderNo     string  `json:"order_no"`     // 商户订单号
	Amount      float64 `json:"amount"`       // 支付金额，单位：元
	Subject     string  `json:"subject"`      // 支付主题
	Description string  `json:"description"`  // 支付描述
	NotifyURL   string  `json:"notify_url"`   // 异步通知地址
	ReturnURL   string  `json:"return_url"`   // 同步回调地址
	ExpireTime  int     `json:"expire_time"`  // 订单过期时间，秒
	ClientIP    string  `json:"client_ip"`    // 客户端IP
	Extra       string  `json:"extra"`        // 附加参数
}

// PrePayResponse 预支付响应
type PrePayResponse struct {
	Code    string          `json:"code"`    // 状态码
	Message string          `json:"message"` // 消息
	Data    *PrePayResponseData `json:"data"` // 数据
}

// PrePayResponseData 预支付响应数据
type PrePayResponseData struct {
	OrderNo    string `json:"order_no"`    // 商户订单号
	TradeNo    string `json:"trade_no"`    // 平台交易号
	PayURL     string `json:"pay_url"`     // 支付链接
	QRCode     string `json:"qrcode"`      // 二维码内容
	ExpireTime int64  `json:"expire_time"` // 过期时间
}

// PayStatusRequest 支付状态查询请求
type PayStatusRequest struct {
	OrderNo string `json:"order_no"` // 商户订单号
	TradeNo string `json:"trade_no"` // 平台交易号
}

// PayStatusResponse 支付状态查询响应
type PayStatusResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    *PayStatusData    `json:"data"`
}

// PayStatusData 支付状态数据
type PayStatusData struct {
	OrderNo     string    `json:"order_no"`      // 商户订单号
	TradeNo     string    `json:"trade_no"`      // 平台交易号
	Status      string    `json:"status"`        // 支付状态
	Amount      float64   `json:"amount"`        // 支付金额
	PayTime     time.Time `json:"pay_time"`      // 支付时间
	PayMethod   string    `json:"pay_method"`    // 支付方式
	TransactionID string  `json:"transaction_id"` // 第三方交易号
}

// RefundRequest 退款请求
type RefundRequest struct {
	OrderNo   string  `json:"order_no"`   // 原订单号
	RefundNo  string  `json:"refund_no"`  // 退款订单号
	Amount    float64 `json:"amount"`     // 退款金额
	Reason    string  `json:"reason"`     // 退款原因
	NotifyURL string  `json:"notify_url"` // 退款通知地址
}

// RefundResponse 退款响应
type RefundResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Data    *RefundData    `json:"data"`
}

// RefundData 退款数据
type RefundData struct {
	RefundNo   string    `json:"refund_no"`   // 退款订单号
	OrderNo    string    `json:"order_no"`    // 原订单号
	Amount     float64   `json:"amount"`      // 退款金额
	Status     string    `json:"status"`      // 退款状态
	RefundTime time.Time `json:"refund_time"` // 退款时间
}

// TransferRequest 转账请求（用于任务结算）
type TransferRequest struct {
	OrderNo   string  `json:"order_no"`   // 商户订单号
	AccountNo string  `json:"account_no"`  // 收款账户
	Amount    float64 `json:"amount"`     // 转账金额
	RealName  string  `json:"real_name"`  // 真实姓名
	BankCode  string  `json:"bank_code"`  // 银行代码
	Memo      string  `json:"memo"`       // 转账备注
	NotifyURL string  `json:"notify_url"` // 异步通知地址
}

// TransferResponse 转账响应
type TransferResponse struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    *TransferData   `json:"data"`
}

// TransferData 转账数据
type TransferData struct {
	OrderNo     string    `json:"order_no"`     // 商户订单号
	TransferNo  string    `json:"transfer_no"`  // 平台转账单号
	Status      string    `json:"status"`       // 转账状态
	Amount      float64   `json:"amount"`       // 转账金额
	TransferTime time.Time `json:"transfer_time"` // 转账时间
}

// NotificationData 回调通知数据
type NotificationData struct {
	OrderNo       string    `json:"order_no"`       // 商户订单号
	TradeNo       string    `json:"trade_no"`       // 平台交易号
	Status        string    `json:"status"`         // 交易状态
	Amount        float64   `json:"amount"`         // 交易金额
	PayTime       time.Time `json:"pay_time"`       // 支付时间
	PayMethod     string    `json:"pay_method"`     // 支付方式
	TransactionID string    `json:"transaction_id"`  // 第三方交易号
	Signature     string    `json:"signature"`      // 签名
	Extra         string    `json:"extra"`          // 附加参数
}

// NewShouqianbaClient 创建收钱吧客户端
func NewShouqianbaClient(cfg *config.ShouqianbaConfig) *ShouqianbaClient {
	return &ShouqianbaClient{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PrePay 预支付
func (c *ShouqianbaClient) PrePay(req *PrePayRequest) (*PrePayResponseData, error) {
	// 设置默认值
	if req.ExpireTime == 0 {
		req.ExpireTime = 900 // 15分钟
	}

	// 构建请求参数
	params := map[string]string{
		"appid":       c.config.AppID,
		"mch_no":      c.config.MerchantNo,
		"order_no":    req.OrderNo,
		"amount":      fmt.Sprintf("%.2f", req.Amount),
		"subject":     req.Subject,
		"description": req.Description,
		"notify_url":  req.NotifyURL,
		"return_url":  req.ReturnURL,
		"expire_time": strconv.Itoa(req.ExpireTime),
		"client_ip":   req.ClientIP,
		"extra":       req.Extra,
		"timestamp":   strconv.FormatInt(time.Now().Unix(), 10),
		"nonce_str":   utils.GenerateNonce(32),
	}

	// 生成签名
	signature := c.generateSignature(params)
	params["sign"] = signature

	// 发送请求
	resp, err := c.postRequest("/api/pay/prepay", params)
	if err != nil {
		return nil, fmt.Errorf("预支付请求失败: %w", err)
	}

	var result PrePayResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "200" {
		return nil, fmt.Errorf("预支付失败: %s", result.Message)
	}

	return result.Data, nil
}

// QueryPayStatus 查询支付状态
func (c *ShouqianbaClient) QueryPayStatus(req *PayStatusRequest) (*PayStatusData, error) {
	params := map[string]string{
		"appid":     c.config.AppID,
		"mch_no":    c.config.MerchantNo,
		"order_no":  req.OrderNo,
		"trade_no":  req.TradeNo,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"nonce_str": utils.GenerateNonce(32),
	}

	// 生成签名
	signature := c.generateSignature(params)
	params["sign"] = signature

	// 发送请求
	resp, err := c.postRequest("/api/pay/query", params)
	if err != nil {
		return nil, fmt.Errorf("查询支付状态失败: %w", err)
	}

	var result PayStatusResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "200" {
		return nil, fmt.Errorf("查询失败: %s", result.Message)
	}

	return result.Data, nil
}

// Refund 退款
func (c *ShouqianbaClient) Refund(req *RefundRequest) (*RefundData, error) {
	params := map[string]string{
		"appid":       c.config.AppID,
		"mch_no":      c.config.MerchantNo,
		"order_no":    req.OrderNo,
		"refund_no":   req.RefundNo,
		"amount":      fmt.Sprintf("%.2f", req.Amount),
		"reason":      req.Reason,
		"notify_url":  req.NotifyURL,
		"timestamp":   strconv.FormatInt(time.Now().Unix(), 10),
		"nonce_str":   utils.GenerateNonce(32),
	}

	// 生成签名
	signature := c.generateSignature(params)
	params["sign"] = signature

	// 发送请求
	resp, err := c.postRequest("/api/pay/refund", params)
	if err != nil {
		return nil, fmt.Errorf("退款请求失败: %w", err)
	}

	var result RefundResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "200" {
		return nil, fmt.Errorf("退款失败: %s", result.Message)
	}

	return result.Data, nil
}

// Transfer 转账（用于任务结算）
func (c *ShouqianbaClient) Transfer(req *TransferRequest) (*TransferData, error) {
	params := map[string]string{
		"appid":       c.config.AppID,
		"mch_no":      c.config.MerchantNo,
		"order_no":    req.OrderNo,
		"account_no":  req.AccountNo,
		"amount":      fmt.Sprintf("%.2f", req.Amount),
		"real_name":   req.RealName,
		"bank_code":   req.BankCode,
		"memo":        req.Memo,
		"notify_url":  req.NotifyURL,
		"timestamp":   strconv.FormatInt(time.Now().Unix(), 10),
		"nonce_str":   utils.GenerateNonce(32),
	}

	// 生成签名
	signature := c.generateSignature(params)
	params["sign"] = signature

	// 发送请求
	resp, err := c.postRequest("/api/pay/transfer", params)
	if err != nil {
		return nil, fmt.Errorf("转账请求失败: %w", err)
	}

	var result TransferResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "200" {
		return nil, fmt.Errorf("转账失败: %s", result.Message)
	}

	return result.Data, nil
}

// VerifyNotification 验证回调通知签名
func (c *ShouqianbaClient) VerifyNotification(data map[string]string) bool {
	signature := data["signature"]
	if signature == "" {
		return false
	}

	// 删除签名参数
	params := make(map[string]string)
	for k, v := range data {
		if k != "signature" {
			params[k] = v
		}
	}

	// 生成签名并对比
	return c.generateSignature(params) == signature
}

// generateSignature 生成签名
func (c *ShouqianbaClient) generateSignature(params map[string]string) string {
	// 按键名排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var builder strings.Builder
	for _, key := range keys {
		if params[key] != "" {
			builder.WriteString(key)
			builder.WriteString("=")
			builder.WriteString(params[key])
			builder.WriteString("&")
		}
	}
	builder.WriteString("key=")
	builder.WriteString(c.config.SecretKey)

	// MD5加密并转为大写
	hash := md5.Sum([]byte(builder.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// postRequest 发送POST请求
func (c *ShouqianbaClient) postRequest(path string, params map[string]string) ([]byte, error) {
	// 构建请求URL
	baseURL := c.config.APIURL
	if c.config.Sandbox {
		baseURL = c.config.SandboxURL
	}

	requestURL := baseURL + path

	// 构建表单数据
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	// 发送请求
	resp, err := c.client.PostForm(requestURL, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}