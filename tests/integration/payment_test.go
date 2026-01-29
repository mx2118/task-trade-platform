package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"task-platform/internal/api/v1/handlers"
	"task-platform/internal/config"
	"task-platform/pkg/payment/shouqianba"
)

// PaymentTestSuite 支付集成测试套件
type PaymentTestSuite struct {
	suite.Suite
	router         *gin.Engine
	cfg            *config.Config
	paymentHandler *handlers.PaymentHandler
	client         *shouqianba.Client
	token          string
	orderNo        string
	tradeNo        string
}

// SetupSuite 设置支付测试套件
func (suite *PaymentTestSuite) SetupSuite() {
	// 初始化测试配置
	suite.cfg = &config.Config{
		Server: config.ServerConfig{
			Mode: gin.TestMode,
		},
		Shouqianba: config.ShouqianbaConfig{
			AppID:        "test_app_id",
			MerchantNo:   "test_merchant_123",
			SecretKey:    "test_secret_key_123456",
			APIURL:       "https://sandbox.shouqianba.com", // 沙箱环境
			Sandbox:      true,
			NotifyURL:    "http://localhost:8080/api/v1/pay/callback",
		},
		JWT: config.JWTConfig{
			Secret: "test_secret_key",
			Expire: 7200,
		},
	}

	// 初始化支付客户端
	suite.client = shouqianba.NewClient(suite.cfg.Shouqianba)
	suite.paymentHandler = handlers.NewPaymentHandler(suite.client)

	// 初始化路由
	suite.initRouter()

	// 初始化测试数据
	suite.initTestData()
}

// TearDownSuite 清理支付测试套件
func (suite *PaymentTestSuite) TearDownSuite() {
	// 清理测试数据
}

// initRouter 初始化支付测试路由
func (suite *PaymentTestSuite) initRouter() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// 添加中间件
	suite.router.Use(gin.Recovery())

	// 注册支付路由
	api := suite.router.Group("/api/v1")
	{
		// 支付接口
		pay := api.Group("/pay")
		{
			pay.POST("/prepay", suite.paymentHandler.PrePay)
			pay.POST("/callback", suite.paymentHandler.Callback)
			pay.POST("/refund_callback", suite.paymentHandler.RefundCallback)
			pay.POST("/transfer_callback", suite.paymentHandler.TransferCallback)
			pay.GET("/query", suite.paymentHandler.QueryStatus)
			pay.POST("/refund", suite.paymentHandler.Refund)
			pay.POST("/settlement", suite.paymentHandler.Settlement)
			pay.GET("/balance", suite.paymentHandler.GetUserBalance)
			pay.GET("/transactions", suite.paymentHandler.GetTransactionList)
		}
	}
}

// initTestData 初始化测试数据
func (suite *PaymentTestSuite) initTestData() {
	// 生成测试token（实际项目中应该通过真实的登录流程获取）
	suite.token = "test_jwt_token"
	suite.orderNo = fmt.Sprintf("TEST_%d", time.Now().Unix())
	suite.tradeNo = fmt.Sprintf("TRADE_%d", time.Now().Unix())
}

// TestPrePay 测试预支付接口
func (suite *PaymentTestSuite) TestPrePay() {
	// 准备预支付数据
	prePayData := map[string]interface{}{
		"order_type":    "task",
		"order_id":      1,
		"amount":        100.00,
		"payment_method": "alipay",
		"subject":       "测试任务支付",
		"return_url":    "http://localhost:3000/payment/success",
		"notify_url":    "http://localhost:8080/api/v1/pay/callback",
	}

	jsonData, _ := json.Marshal(prePayData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	// 执行请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "order_no")
		assert.Contains(suite.T(), data, "payment_url")
		assert.Contains(suite.T(), data, "qr_code")
		assert.Contains(suite.T(), data, "expires_at")

		// 保存订单号用于后续测试
		suite.orderNo = data["order_no"].(string)
	}
}

// TestPrePayInvalidAmount 测试无效金额的预支付
func (suite *PaymentTestSuite) TestPrePayInvalidAmount() {
	prePayData := map[string]interface{}{
		"order_type":    "task",
		"order_id":      1,
		"amount":        -100.00, // 无效金额
		"payment_method": "alipay",
		"subject":       "测试任务支付",
	}

	jsonData, _ := json.Marshal(prePayData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error", response["status"])
}

// TestPrePayMissingFields 测试缺少必要字段的预支付
func (suite *PaymentTestSuite) TestPrePayMissingFields() {
	prePayData := map[string]interface{}{
		"amount": 100.00,
		// 缺少必要字段：order_type, order_id, payment_method
	}

	jsonData, _ := json.Marshal(prePayData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestPaymentCallback 测试支付回调
func (suite *PaymentTestSuite) TestPaymentCallback() {
	// 准备支付回调数据（模拟收钱吧回调格式）
	callbackData := map[string]interface{}{
		"merchant_no":   suite.cfg.Shouqianba.MerchantNo,
		"terminal_no":   "test_terminal",
		"order_no":      suite.orderNo,
		"trade_no":      suite.tradeNo,
		"pay_status":    "success",
		"pay_amount":    "100.00",
		"pay_time":      time.Now().Format("20060102150405"),
		"subject":       "测试任务支付",
		"buyer_account": "test_buyer@example.com",
		"seller_account": suite.cfg.Shouqianba.MerchantNo,
		"device_info":   "web",
		"goods_tag":     "task_payment",
		"extend":        "",
		"sign_type":     "MD5",
	}

	// 生成签名（实际项目中应该使用真实的签名算法）
	sign := suite.generateSignature(callbackData)
	callbackData["sign"] = sign

	jsonData, _ := json.Marshal(callbackData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/callback", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 执行回调请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 验证回调响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["status"])
}

// TestPaymentCallbackInvalidSignature 测试无效签名的支付回调
func (suite *PaymentTestSuite) TestPaymentCallbackInvalidSignature() {
	callbackData := map[string]interface{}{
		"merchant_no":   suite.cfg.Shouqianba.MerchantNo,
		"order_no":      suite.orderNo,
		"trade_no":      suite.tradeNo,
		"pay_status":    "success",
		"pay_amount":    "100.00",
		"sign":          "invalid_signature", // 无效签名
	}

	jsonData, _ := json.Marshal(callbackData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/callback", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error", response["status"])
}

// TestPaymentCallbackDuplicate 测试重复的支付回调
func (suite *PaymentTestSuite) TestPaymentCallbackDuplicate() {
	// 先发送一次成功的回调
	suite.TestPaymentCallback()

	// 再次发送相同的回调
	callbackData := map[string]interface{}{
		"merchant_no":   suite.cfg.Shouqianba.MerchantNo,
		"order_no":      suite.orderNo,
		"trade_no":      suite.tradeNo,
		"pay_status":    "success",
		"pay_amount":    "100.00",
		"sign":          suite.generateSignature(callbackData),
	}

	jsonData, _ := json.Marshal(callbackData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/callback", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 重复回调应该返回成功，但不应该重复处理
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestPaymentQuery 测试支付查询
func (suite *PaymentTestSuite) TestPaymentQuery() {
	// 查询支付状态
	req, _ := http.NewRequest("GET", "/api/v1/pay/query?order_no="+suite.orderNo, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "order_no")
		assert.Contains(suite.T(), data, "status")
		assert.Contains(suite.T(), data, "amount")
		assert.Contains(suite.T(), data, "created_at")
	}
}

// TestPaymentQueryInvalidOrder 测试查询无效订单号
func (suite *PaymentTestSuite) TestPaymentQueryInvalidOrder() {
	req, _ := http.NewRequest("GET", "/api/v1/pay/query?order_no=INVALID_ORDER", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// TestRefund 测试退款接口
func (suite *PaymentTestSuite) TestRefund() {
	// 先确保支付成功
	suite.TestPaymentCallback()

	// 准备退款数据
	refundData := map[string]interface{}{
		"order_no":      suite.orderNo,
		"refund_amount": 50.00, // 部分退款
		"refund_reason":  "用户申请退款",
	}

	jsonData, _ := json.Marshal(refundData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/refund", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "refund_order_no")
		assert.Contains(suite.T(), data, "refund_amount")
		assert.Contains(suite.T(), data, "status")
	}
}

// TestRefundExceedAmount 测试超额退款
func (suite *PaymentTestSuite) TestRefundExceedAmount() {
	refundData := map[string]interface{}{
		"order_no":      suite.orderNo,
		"refund_amount": 200.00, // 超出原订单金额
		"refund_reason":  "测试超额退款",
	}

	jsonData, _ := json.Marshal(refundData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/refund", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestSettlement 测试结算接口
func (suite *PaymentTestSuite) TestSettlement() {
	// 准备结算数据
	settlementData := map[string]interface{}{
		"task_id":        1,
		"taker_user_id":  2,
		"settle_amount":  90.00, // 扣除10%平台费用
		"settle_reason":  "任务完成结算",
	}

	jsonData, _ := json.Marshal(settlementData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/settlement", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "settlement_no")
		assert.Contains(suite.T(), data, "settle_amount")
		assert.Contains(suite.T(), data, "status")
	}
}

// TestGetUserBalance 测试获取用户余额
func (suite *PaymentTestSuite) TestGetUserBalance() {
	req, _ := http.NewRequest("GET", "/api/v1/pay/balance", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "balance")
		assert.Contains(suite.T(), data, "frozen_amount")
		assert.Contains(suite.T(), data, "total_income")
		assert.Contains(suite.T(), data, "total_expense")
	}
}

// TestGetTransactionList 测试获取交易记录
func (suite *PaymentTestSuite) TestGetTransactionList() {
	req, _ := http.NewRequest("GET", "/api/v1/pay/transactions?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 验证返回的数据结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "list")
		assert.Contains(suite.T(), data, "total")
		assert.Contains(suite.T(), data, "page")
		assert.Contains(suite.T(), data, "limit")
	}
}

// TestPaymentConcurrency 测试支付并发处理
func (suite *PaymentTestSuite) TestPaymentConcurrency() {
	const numRequests = 10
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			prePayData := map[string]interface{}{
				"order_type":    "task",
				"order_id":      id + 1000, // 使用不同的订单ID
				"amount":        100.00,
				"payment_method": "alipay",
				"subject":       fmt.Sprintf("并发测试任务%d", id),
			}

			jsonData, _ := json.Marshal(prePayData)
			req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+suite.token)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			results <- w.Code
		}(i)
	}

	// 收集所有结果
	successCount := 0
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		if statusCode == http.StatusOK {
			successCount++
		}
	}

	// 验证至少有一半的请求成功
	assert.GreaterOrEqual(suite.T(), successCount, numRequests/2)
}

// TestPaymentSecurity 测试支付安全性
func (suite *PaymentTestSuite) TestPaymentSecurity() {
	// 测试SQL注入攻击
	maliciousData := map[string]interface{}{
		"order_type":    "task; DROP TABLE payments; --",
		"order_id":      1,
		"amount":        100.00,
		"payment_method": "alipay",
		"subject":       "恶意输入测试",
	}

	jsonData, _ := json.Marshal(maliciousData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 应该被输入验证拦截
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestPaymentPerformance 测试支付性能
func (suite *PaymentTestSuite) TestPaymentPerformance() {
	start := time.Now()

	// 批量创建预支付订单
	for i := 0; i < 50; i++ {
		prePayData := map[string]interface{}{
			"order_type":    "task",
			"order_id":      i + 2000,
			"amount":        100.00,
			"payment_method": "alipay",
			"subject":       fmt.Sprintf("性能测试任务%d", i),
		}

		jsonData, _ := json.Marshal(prePayData)
		req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.token)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	avgDuration := duration / 50

	// 平均每个请求应该在20ms以内
	suite.T().Logf("平均支付创建时间: %v", avgDuration)
	assert.Less(suite.T(), avgDuration, 20*time.Millisecond)
}

// generateSignature 生成测试签名（模拟收钱吧签名算法）
func (suite *PaymentTestSuite) generateSignature(params map[string]interface{}) string {
	// 这里应该实现真实的收钱吧签名算法
	// 为了测试目的，返回一个简单的模拟签名
	return "test_signature"
}

// RunPaymentTests 运行所有支付测试
func TestPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentTestSuite))
}