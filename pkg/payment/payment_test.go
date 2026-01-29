package payment

import (
	"testing"
	"time"

	"task-platform-api/internal/config"
)

func TestGenerateSignature(t *testing.T) {
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key",
	}
	
	client := NewShouqianbaClient(cfg)
	
	params := map[string]string{
		"appid":     "test-appid",
		"mch_no":    "test-mch-no",
		"order_no":  "SQ20240101120000123456",
		"amount":    "100.00",
		"timestamp": "1704067200",
		"nonce_str": "test-nonce",
	}
	
	signature := client.generateSignature(params)
	
	// 测试签名生成
	if signature == "" {
		t.Error("签名不能为空")
	}
	
	// 测试参数顺序一致性
	params2 := map[string]string{
		"amount":    "100.00",
		"appid":     "test-appid",
		"mch_no":    "test-mch-no",
		"nonce_str": "test-nonce",
		"order_no":  "SQ20240101120000123456",
		"timestamp": "1704067200",
	}
	
	signature2 := client.generateSignature(params2)
	
	if signature != signature2 {
		t.Error("相同参数应生成相同签名，顺序不影响结果")
	}
}

func TestVerifyNotification(t *testing.T) {
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key",
	}
	
	client := NewShouqianbaClient(cfg)
	
	// 构造测试数据
	params := map[string]string{
		"appid":     "test-appid",
		"mch_no":    "test-mch-no",
		"order_no":  "SQ20240101120000123456",
		"trade_no":  "SQ_TRADE_123456",
		"status":    "SUCCESS",
		"amount":    "100.00",
		"pay_time":  "2024-01-01 12:00:00",
		"pay_method": "ALIPAY",
		"timestamp": "1704067200",
		"nonce_str": "test-nonce",
	}
	
	// 生成签名
	signature := client.generateSignature(params)
	params["signature"] = signature
	
	// 验证签名
	if !client.VerifyNotification(params) {
		t.Error("签名验证失败")
	}
	
	// 修改参数，签名验证应该失败
	params["amount"] = "200.00"
	if client.VerifyNotification(params) {
		t.Error("参数修改后签名验证应该失败")
	}
}

func TestPrePayRequest(t *testing.T) {
	// 这是一个集成测试，需要真实的API环境
	// 在CI/CD环境中应该跳过或使用模拟服务
	t.Skip("需要真实的API环境")
	
	cfg := &config.ShouqianbaConfig{
		AppID:      "test-appid",
		MerchantNo: "test-mch-no",
		SecretKey:  "test-secret-key",
		APIURL:     "https://api.shouqianba.com",
		Sandbox:    true,
	}
	
	client := NewShouqianbaClient(cfg)
	
	req := &PrePayRequest{
		OrderNo:     "SQ20240101120000123456",
		Amount:      100.00,
		Subject:     "测试任务",
		Description: "测试任务描述",
		NotifyURL:   "http://localhost:8080/api/v1/pay/callback",
		ClientIP:    "127.0.0.1",
	}
	
	// 这个测试需要实际的API调用
	_, err := client.PrePay(req)
	if err != nil {
		t.Errorf("预支付请求失败: %v", err)
	}
}

func TestPaymentStatus(t *testing.T) {
	t.Skip("需要真实的API环境")
	
	cfg := &config.ShouqianbaConfig{
		AppID:      "test-appid",
		MerchantNo: "test-mch-no",
		SecretKey:  "test-secret-key",
		APIURL:     "https://api.shouqianba.com",
		Sandbox:    true,
	}
	
	client := NewShouqianbaClient(cfg)
	
	req := &PayStatusRequest{
		OrderNo: "SQ20240101120000123456",
	}
	
	_, err := client.QueryPayStatus(req)
	if err != nil {
		t.Errorf("查询支付状态失败: %v", err)
	}
}

func TestRefundRequest(t *testing.T) {
	t.Skip("需要真实的API环境")
	
	cfg := &config.ShouqianbaConfig{
		AppID:      "test-appid",
		MerchantNo: "test-mch-no",
		SecretKey:  "test-secret-key",
		APIURL:     "https://api.shouqianba.com",
		Sandbox:    true,
	}
	
	client := NewShouqianbaClient(cfg)
	
	req := &RefundRequest{
		OrderNo:   "SQ20240101120000123456",
		RefundNo:  "RF20240101120000123456",
		Amount:    100.00,
		Reason:    "测试退款",
		NotifyURL: "http://localhost:8080/api/v1/pay/refund_callback",
	}
	
	_, err := client.Refund(req)
	if err != nil {
		t.Errorf("退款请求失败: %v", err)
	}
}

func TestTransferRequest(t *testing.T) {
	t.Skip("需要真实的API环境")
	
	cfg := &config.ShouqianbaConfig{
		AppID:      "test-appid",
		MerchantNo: "test-mch-no",
		SecretKey:  "test-secret-key",
		APIURL:     "https://api.shouqianba.com",
		Sandbox:    true,
	}
	
	client := NewShouqianbaClient(cfg)
	
	req := &TransferRequest{
		OrderNo:   "TR20240101120000123456",
		AccountNo: "test@example.com",
		Amount:    100.00,
		RealName:  "测试用户",
		Memo:      "任务结算测试",
		NotifyURL: "http://localhost:8080/api/v1/pay/transfer_callback",
	}
	
	_, err := client.Transfer(req)
	if err != nil {
		t.Errorf("转账请求失败: %v", err)
	}
}

// 基准测试
func BenchmarkGenerateSignature(b *testing.B) {
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key-for-benchmark-testing",
	}
	
	client := NewShouqianbaClient(cfg)
	
	params := map[string]string{
		"appid":     "test-appid",
		"mch_no":    "test-mch-no",
		"order_no":  "SQ20240101120000123456",
		"amount":    "100.00",
		"subject":   "测试任务标题",
		"timestamp": "1704067200",
		"nonce_str": "test-nonce-string-for-benchmark",
		"extra":     "additional-parameters-for-testing",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.generateSignature(params)
	}
}

// 并发安全测试
func TestConcurrentSignature(t *testing.T) {
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key",
	}
	
	client := NewShouqianbaClient(cfg)
	
	params := map[string]string{
		"appid":     "test-appid",
		"mch_no":    "test-mch-no",
		"order_no":  "SQ20240101120000123456",
		"amount":    "100.00",
		"timestamp": "1704067200",
		"nonce_str": "test-nonce",
	}
	
	// 并发生成签名
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			signature := client.generateSignature(params)
			if signature == "" {
				t.Error("并发生成签名失败")
			}
			done <- true
		}()
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// 边界条件测试
func TestEdgeCases(t *testing.T) {
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key",
	}
	
	client := NewShouqianbaClient(cfg)
	
	// 测试空参数
	emptyParams := map[string]string{}
	signature := client.generateSignature(emptyParams)
	expectedSignature := client.generateSignature(map[string]string{"key": "test-secret-key"})
	
	if signature != expectedSignature {
		t.Error("空参数签名计算错误")
	}
	
	// 测试特殊字符
	specialParams := map[string]string{
		"appid":     "test@app-id",
		"mch_no":    "test&mch-no",
		"order_no":  "SQ#2024-01-01#1200#00123456",
		"amount":    "100.00",
		"timestamp": "1704067200",
		"nonce_str": "test-nonce!@#$%^&*()",
	}
	
	signature2 := client.generateSignature(specialParams)
	if signature2 == "" {
		t.Error("特殊字符参数签名生成失败")
	}
}

// 性能测试
func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}
	
	cfg := &config.ShouqianbaConfig{
		SecretKey: "test-secret-key-for-performance-testing",
	}
	
	client := NewShouqianbaClient(cfg)
	
	start := time.Now()
	
	// 执行大量签名生成
	for i := 0; i < 10000; i++ {
		params := map[string]string{
			"appid":     "test-appid",
			"mch_no":    "test-mch-no",
			"order_no":  "SQ20240101120000123456",
			"amount":    "100.00",
			"timestamp": "1704067200",
			"nonce_str": "test-nonce",
		}
		client.generateSignature(params)
	}
	
	duration := time.Since(start)
	t.Logf("10000次签名生成耗时: %v, 平均每次: %v", duration, duration/10000)
	
	// 性能要求：每次签名生成应该在1ms以内
	avgTime := duration / 10000
	if avgTime > time.Millisecond {
		t.Errorf("性能不达标，平均签名生成时间: %v > 1ms", avgTime)
	}
}