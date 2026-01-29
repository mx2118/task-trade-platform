package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	go_redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"task-platform-api/internal/api/v1/handlers"
	"task-platform-api/internal/api/v1/middleware"
	"task-platform-api/internal/config"
	"task-platform-api/internal/services"
	"task-platform-api/pkg/utils"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB, rdb *go_redis.Client, cfg *config.Config, logger *zap.Logger) *gin.Engine {
    router := gin.New()

    // 全局中间件
    router.Use(middleware.Logging(logger))
    router.Use(gin.Recovery())
    router.Use(middleware.RequestID())
    router.Use(middleware.CustomHeaders())
    
    if cfg.Security.EnableCORS {
        router.Use(middleware.CORS(&cfg.Security))
    }
    
    // 限流中间件
    if cfg.Security.RateLimit > 0 {
        router.Use(middleware.RateLimit(cfg.Security.RateLimit))
    }

    // 初始化处理器
    initHandlers(router, db, rdb, cfg, logger)

    // 健康检查
    router.GET("/health", func(c *gin.Context) {
        utils.SuccessResponse(c, gin.H{
            "status":    "ok",
            "timestamp": time.Now().Unix(),
            "version":   "1.0.0",
        })
    })

    return router
}

// initHandlers 初始化处理器
func initHandlers(router *gin.Engine, db *gorm.DB, rdb *go_redis.Client, cfg *config.Config, logger *zap.Logger) {
    // 创建服务实例
    paymentService := services.NewPaymentService(db, cfg)
    
    // 创建处理器实例
    authHandler := handlers.NewAuthHandler(db, rdb, cfg, logger)
    userHandler := handlers.NewUserHandler(db, rdb, cfg, logger)
    taskHandler := handlers.NewTaskHandler(db, rdb, cfg, logger)
    paymentHandler := handlers.NewPaymentHandler(paymentService)

    // API路由组
    v1 := router.Group("/api/v1")
    {
        // 无需认证的公共接口
        public := v1.Group("/")
        {
            // 认证相关
            public.POST("/auth/wechat/login", authHandler.WechatLogin)
            public.POST("/auth/alipay/login", authHandler.AlipayLogin)
            public.POST("/auth/refresh", authHandler.RefreshToken)
            
            // 任务相关（公开浏览）
            public.GET("/tasks", taskHandler.ListTasks)
            public.GET("/tasks/:id", taskHandler.GetTask)
            public.GET("/categories", taskHandler.ListCategories)
            
            // 系统信息
            public.GET("/system/info", handlers.GetSystemInfo)
            public.GET("/system/announce", handlers.GetAnnouncements)
        }

        // 需要认证的接口
        protected := v1.Group("/")
        protected.Use(middleware.JWTAuth(&cfg.JWT, logger))
        protected.Use(middleware.RequireNormalUser(db))
        {
            // 用户相关
            userGroup := protected.Group("/user")
            {
                userGroup.GET("/profile", userHandler.GetProfile)
                userGroup.PUT("/profile", userHandler.UpdateProfile)
                userGroup.POST("/upload-avatar", userHandler.UploadAvatar)
                userGroup.GET("/wallet", userHandler.GetWallet)
                userGroup.GET("/transactions", userHandler.GetTransactions)
                userGroup.POST("/logout", authHandler.Logout)
            }

            // 任务相关
            taskGroup := protected.Group("/tasks")
            {
                taskGroup.POST("/", taskHandler.CreateTask)
                taskGroup.PUT("/:id", taskHandler.UpdateTask)
                taskGroup.DELETE("/:id", taskHandler.DeleteTask)
                taskGroup.POST("/:id/apply", taskHandler.ApplyTask)
                taskGroup.POST("/:id/take", taskHandler.TakeTask)
                taskGroup.POST("/:id/deliver", taskHandler.DeliverTask)
                taskGroup.POST("/:id/accept", taskHandler.AcceptTask)
                taskGroup.GET("/my-published", taskHandler.GetMyPublishedTasks)
                taskGroup.GET("/my-applied", taskHandler.GetMyAppliedTasks)
                taskGroup.GET("/my-taken", taskHandler.GetMyTakenTasks)
            }

            // 支付相关
            paymentGroup := protected.Group("/pay")
            {
                paymentGroup.POST("/prepay", paymentHandler.PrePay)
                paymentGroup.GET("/query", paymentHandler.QueryStatus)
                paymentGroup.POST("/refund", paymentHandler.Refund)
                paymentGroup.POST("/settlement", paymentHandler.Settlement)
                paymentGroup.GET("/balance", paymentHandler.GetUserBalance)
                paymentGroup.GET("/transactions", paymentHandler.GetTransactionList)
            }

            // 评价和申诉
            reviewGroup := protected.Group("/reviews")
            {
                reviewGroup.POST("/", handlers.CreateReview)
                reviewGroup.GET("/received", handlers.GetReceivedReviews)
                reviewGroup.GET("/given", handlers.GetGivenReviews)
            }

            complaintGroup := protected.Group("/complaints")
            {
                complaintGroup.POST("/", handlers.CreateComplaint)
                complaintGroup.GET("/", handlers.GetMyComplaints)
                complaintGroup.GET("/:id", handlers.GetComplaintDetail)
            }

            // 通知
            notificationGroup := protected.Group("/notifications")
            {
                notificationGroup.GET("/", handlers.GetNotifications)
                notificationGroup.PUT("/:id/read", handlers.MarkNotificationRead)
                notificationGroup.PUT("/read-all", handlers.MarkAllNotificationsRead)
            }

            // 文件上传
            protected.POST("/upload", handlers.UploadFile)
        }

        // 管理员接口
        admin := v1.Group("/admin")
        admin.Use(middleware.JWTAuth(&cfg.JWT, logger))
        admin.Use(middleware.RequireRole("admin"))
        {
            // 用户管理
            admin.GET("/users", handlers.AdminListUsers)
            admin.GET("/users/:id", handlers.AdminGetUser)
            admin.PUT("/users/:id/status", handlers.AdminUpdateUserStatus)
            admin.PUT("/users/:id/credit", handlers.AdminUpdateUserCredit)
            
            // 任务管理
            admin.GET("/tasks", handlers.AdminListTasks)
            admin.GET("/tasks/pending", handlers.AdminGetPendingTasks)
            admin.PUT("/tasks/:id/approve", handlers.AdminApproveTask)
            admin.PUT("/tasks/:id/reject", handlers.AdminRejectTask)
            
            // 交易管理
            admin.GET("/trades", handlers.AdminListTrades)
            admin.GET("/trades/:id", handlers.AdminGetTrade)
            admin.POST("/trades/:id/manual-settle", handlers.AdminManualSettle)
            
            // 申诉处理
            admin.GET("/complaints", handlers.AdminListComplaints)
            admin.PUT("/complaints/:id/handle", handlers.AdminHandleComplaint)
            
            // 系统设置
            admin.GET("/system/stats", handlers.AdminGetSystemStats)
            admin.PUT("/system/config", handlers.AdminUpdateSystemConfig)
        }
    }

    // 收钱吧回调接口（特殊处理，不需要JWT但需要签名验证）
    callback := v1.Group("/pay")
    {
        callback.POST("/callback", paymentHandler.Callback)
        callback.POST("/refund_callback", paymentHandler.RefundCallback)
        callback.POST("/transfer_callback", paymentHandler.TransferCallback)
    }
}