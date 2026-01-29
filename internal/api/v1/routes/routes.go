package routes

import (
	"github.com/gin-gonic/gin"

	"task-platform-api/internal/api/v1/handlers"
)

// SetupRoutes 设置路由
func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	paymentHandler *handlers.PaymentHandler,
) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "服务运行正常",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "注册接口"})
			})
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "登录接口"})
			})
			auth.POST("/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "登出接口"})
			})
		}

		// 支付相关路由
		pay := v1.Group("/pay")
		{
			pay.POST("/prepay", paymentHandler.PrePay)
			pay.GET("/status/:order_no", paymentHandler.QueryStatus)
			pay.POST("/callback", paymentHandler.PaymentCallback)
		}

		// 用户相关路由
		user := v1.Group("/user")
		{
			user.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "用户信息接口"})
			})
		}

		// 任务相关路由
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "任务列表接口"})
			})
			tasks.GET("/:id", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "任务详情接口"})
			})
		}

		// 系统信息
		v1.GET("/system/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"name": "任务交易平台",
				"version": "1.0.0",
				"description": "一个现代化的任务发布和接取平台",
			})
		})

		v1.GET("/announcements", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": []gin.H{
					{
						"id": 1,
						"title": "欢迎使用任务交易平台",
						"content": "平台现已正式上线，欢迎各位用户使用",
						"type": "system",
						"created_at": "2024-01-01T00:00:00Z",
					},
				},
			})
		})
	}
}