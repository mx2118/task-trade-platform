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
	"gorm.io/gorm"

	"task-platform/internal/api/v1/handlers"
	"task-platform/internal/api/v1/middleware"
	"task-platform/internal/config"
	"task-platform/internal/models"
	"task-platform/pkg/payment/shouqianba"
	"task-platform/pkg/utils"
)

// UserFlowTestSuite 用户完整流程测试套件
type UserFlowTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
	cfg    *config.Config
	userID uint
	taskID  uint
	token  string
}

// SetupSuite 设置测试套件
func (suite *UserFlowTestSuite) SetupSuite() {
	// 设置测试配置
	suite.cfg = &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Mode: gin.TestMode,
		},
		Database: config.DatabaseConfig{
			DSN: "test_user:test_pass@tcp(localhost:3306)/task_platform_test?charset=utf8mb4&parseTime=True&loc=Local",
		},
		JWT: config.JWTConfig{
			Secret: "test_secret_key",
			Expire: 7200,
		},
		Shouqianba: config.ShouqianbaConfig{
			AppID:        "test_app_id",
			MerchantNo:   "test_merchant",
			SecretKey:    "test_secret",
			APIURL:       "https://test-api.shouqianba.com",
			Sandbox:      true,
		},
	}

	// 初始化数据库
	suite.initDatabase()

	// 初始化路由
	suite.initRouter()
}

// TearDownSuite 清理测试套件
func (suite *UserFlowTestSuite) TearDownSuite() {
	if suite.db != nil {
		// 清理测试数据
		suite.db.Exec("DELETE FROM users")
		suite.db.Exec("DELETE FROM tasks")
		suite.db.Exec("DELETE FROM payments")
		suite.db.Exec("DELETE FROM transactions")
	}
}

// initDatabase 初始化测试数据库
func (suite *UserFlowTestSuite) initDatabase() {
	// 这里应该连接到测试数据库
	// 为了演示，我们使用模拟的数据库连接
	// 实际项目中应该使用真实的数据库连接
}

// initRouter 初始化路由
func (suite *UserFlowTestSuite) initRouter() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// 添加中间件
	suite.router.Use(gin.Recovery())
	suite.router.Use(middleware.CORS())

	// 初始化处理器
	userHandler := handlers.NewUserHandler(suite.db, nil, suite.cfg, nil)
	taskHandler := handlers.NewTaskHandler(suite.db, nil, suite.cfg, nil)
	paymentHandler := handlers.NewPaymentHandler(shouqianba.NewClient(suite.cfg.Shouqianba))

	// 注册路由
	api := suite.router.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
		}

		// 用户路由
		users := api.Group("/user")
		users.Use(middleware.JWTAuth(&suite.cfg.JWT, nil))
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/stats", userHandler.GetStats)
		}

		// 任务路由
		tasks := api.Group("/tasks")
		tasks.Use(middleware.JWTAuth(&suite.cfg.JWT, nil))
		{
			tasks.POST("/", taskHandler.CreateTask)
			tasks.GET("/", taskHandler.ListTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.POST("/:id/apply", taskHandler.ApplyTask)
			tasks.POST("/:id/deliver", taskHandler.DeliverTask)
			tasks.POST("/:id/confirm", taskHandler.ConfirmTask)
		}

		// 支付路由
		payments := api.Group("/pay")
		payments.Use(middleware.JWTAuth(&suite.cfg.JWT, nil))
		{
			payments.POST("/prepay", paymentHandler.PrePay)
			payments.POST("/callback", paymentHandler.Callback)
			payments.GET("/query", paymentHandler.QueryStatus)
			payments.POST("/refund", paymentHandler.Refund)
		}
	}
}

// TestUserRegistration 测试用户注册流程
func (suite *UserFlowTestSuite) TestUserRegistration() {
	// 准备注册数据
	registerData := map[string]interface{}{
		"phone":    "13800138000",
		"nickname":  "测试用户",
		"password":  "test123456",
		"code":     "123456",
	}

	// 转换为JSON
	jsonData, _ := json.Marshal(registerData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])
	assert.NotNil(suite.T(), response["data"])

	// 提取用户ID和token
	if data, ok := response["data"].(map[string]interface{}); ok {
		suite.userID = uint(data["id"].(float64))
		suite.token = data["token"].(string)
	}
}

// TestUserLogin 测试用户登录流程
func (suite *UserFlowTestSuite) TestUserLogin() {
	// 准备登录数据
	loginData := map[string]interface{}{
		"phone":   "13800138000",
		"code":    "123456",
	}

	jsonData, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 更新token
	if data, ok := response["data"].(map[string]interface{}); ok {
		suite.token = data["token"].(string)
	}
}

// TestUserProfile 测试用户资料管理
func (suite *UserFlowTestSuite) TestUserProfile() {
	// 测试获取用户资料
	req, _ := http.NewRequest("GET", "/api/v1/user/profile", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	// 测试更新用户资料
	updateData := map[string]interface{}{
		"nickname": "更新用户名",
		"bio":      "这是我的个人简介",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", "/api/v1/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestTaskCreation 测试任务创建流程
func (suite *UserFlowTestSuite) TestTaskCreation() {
	taskData := map[string]interface{}{
		"title":       "测试任务标题",
		"description": "这是一个测试任务的描述，包含了详细的任务要求和工作内容。",
		"category_id": 1,
		"price":       100.00,
		"deadline":    time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		"location":    "北京市朝阳区",
		"is_urgent":   true,
		"is_remote":   false,
	}

	jsonData, _ := json.Marshal(taskData)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])

	// 提取任务ID
	if data, ok := response["data"].(map[string]interface{}); ok {
		suite.taskID = uint(data["id"].(float64))
	}
}

// TestTaskRetrieval 测试任务获取流程
func (suite *UserFlowTestSuite) TestTaskRetrieval() {
	// 测试获取任务列表
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 测试获取任务详情
	taskURL := fmt.Sprintf("/api/v1/tasks/%d", suite.taskID)
	req, _ = http.NewRequest("GET", taskURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["data"])
}

// TestTaskApplication 测试任务申请流程
func (suite *UserFlowTestSuite) TestTaskApplication() {
	applicationData := map[string]interface{}{
		"message":        "我对这个任务很感兴趣，有相关经验可以胜任。",
		"estimated_time": "3天",
		"contact":        "13800138001",
	}

	jsonData, _ := json.Marshal(applicationData)
	applyURL := fmt.Sprintf("/api/v1/tasks/%d/apply", suite.taskID)
	req, _ := http.NewRequest("POST", applyURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 注意：不能申请自己发布的任务，所以这里应该返回错误
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestPaymentFlow 测试支付流程
func (suite *UserFlowTestSuite) TestPaymentFlow() {
	paymentData := map[string]interface{}{
		"order_type":    "task",
		"order_id":      suite.taskID,
		"amount":        100.00,
		"payment_method": "alipay",
		"return_url":    "https://example.com/return",
		"notify_url":    "https://example.com/notify",
	}

	jsonData, _ := json.Marshal(paymentData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/prepay", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["data"])
}

// TestPaymentCallback 测试支付回调
func (suite *UserFlowTestSuite) TestPaymentCallback() {
	// 模拟收钱吧支付回调数据
	callbackData := map[string]interface{}{
		"merchant_no":   "test_merchant",
		"terminal_no":   "test_terminal",
		"order_no":      fmt.Sprintf("TASK_%d_%d", suite.taskID, time.Now().Unix()),
		"trade_no":      "test_trade_no_" + fmt.Sprint(time.Now().Unix()),
		"pay_status":    "success",
		"pay_amount":    "100.00",
		"pay_time":      time.Now().Format("20060102150405"),
		"subject":       "任务支付",
		"sign":          "test_signature", // 这里应该是真实的签名
	}

	jsonData, _ := json.Marshal(callbackData)
	req, _ := http.NewRequest("POST", "/api/v1/pay/callback", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestPaymentQuery 测试支付查询
func (suite *UserFlowTestSuite) TestPaymentQuery() {
	orderNo := fmt.Sprintf("TASK_%d_%d", suite.taskID, time.Now().Unix())
	req, _ := http.NewRequest("GET", "/api/v1/pay/query?order_no="+orderNo, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["data"])
}

// TestTaskDelivery 测试任务交付流程
func (suite *UserFlowTestSuite) TestTaskDelivery() {
	deliveryData := map[string]interface{}{
		"delivery_content": "任务已完成，请查看附件。",
		"delivery_images": []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
	}

	jsonData, _ := json.Marshal(deliveryData)
	deliveryURL := fmt.Sprintf("/api/v1/tasks/%d/deliver", suite.taskID)
	req, _ := http.NewRequest("POST", deliveryURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 这里应该检查任务状态，如果任务还没有被接受，则无法交付
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestTaskConfirmation 测试任务确认流程
func (suite *UserFlowTestSuite) TestTaskConfirmation() {
	confirmURL := fmt.Sprintf("/api/v1/tasks/%d/confirm", suite.taskID)
	req, _ := http.NewRequest("POST", confirmURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 这里应该检查任务状态，如果任务还没有被交付，则无法确认
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestUserStats 测试用户统计
func (suite *UserFlowTestSuite) TestUserStats() {
	req, _ := http.NewRequest("GET", "/api/v1/user/stats", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["data"])

	// 验证统计数据的结构
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(suite.T(), data, "completed_tasks")
		assert.Contains(suite.T(), data, "total_earnings")
		assert.Contains(suite.T(), data, "credit")
	}
}

// TestErrorHandling 测试错误处理
func (suite *UserFlowTestSuite) TestErrorHandling() {
	// 测试无效的token
	req, _ := http.NewRequest("GET", "/api/v1/user/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// 测试无效的任务ID
	req, _ = http.NewRequest("GET", "/api/v1/tasks/999999", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// 测试无效的请求数据
	invalidData := map[string]interface{}{
		"invalid_field": "invalid_value",
	}

	jsonData, _ := json.Marshal(invalidData)
	req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestInputValidation 测试输入验证
func (suite *UserFlowTestSuite) TestInputValidation() {
	// 测试空的用户名
	invalidData := map[string]interface{}{
		"nickname": "",
		"password": "test123456",
	}

	jsonData, _ := json.Marshal(invalidData)
	req, _ := http.NewRequest("PUT", "/api/v1/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// 测试无效的价格数据
	invalidTaskData := map[string]interface{}{
		"title":       "测试任务",
		"description": "测试描述",
		"category_id": 1,
		"price":       -100, // 无效价格
		"deadline":    time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	jsonData, _ = json.Marshal(invalidTaskData)
	req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestConcurrentRequests 测试并发请求
func (suite *UserFlowTestSuite) TestConcurrentRequests() {
	// 模拟多个用户同时获取任务列表
	const numRequests = 10
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
			req.Header.Set("Authorization", "Bearer "+suite.token)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			results <- w.Code
		}()
	}

	// 收集所有结果
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		assert.Equal(suite.T(), http.StatusOK, statusCode)
	}
}

// TestPerformance 性能测试
func (suite *UserFlowTestSuite) TestPerformance() {
	// 测试任务列表查询性能
	start := time.Now()

	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+suite.token)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	avgDuration := duration / 100

	// 平均每个请求应该在10ms以内
	suite.T().Logf("平均响应时间: %v", avgDuration)
	assert.Less(suite.T(), avgDuration, 10*time.Millisecond)
}

// RunTests 运行所有测试
func TestUserFlowTestSuite(t *testing.T) {
	suite.Run(t, new(UserFlowTestSuite))
}