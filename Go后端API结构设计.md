# Go后端API项目结构设计

## 1. 项目目录结构

```
task-platform-api/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口
├── internal/                     # 私有应用程序代码
│   ├── api/                      # API层
│   │   ├── v1/                  # API版本1
│   │   │   ├── handlers/        # HTTP处理器
│   │   │   │   ├── auth.go
│   │   │   │   ├── user.go
│   │   │   │   ├── task.go
│   │   │   │   ├── order.go
│   │   │   │   ├── payment.go
│   │   │   │   ├── review.go
│   │   │   │   └── message.go
│   │   │   ├── middleware/      # 中间件
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logging.go
│   │   │   │   ├── ratelimit.go
│   │   │   │   └── validation.go
│   │   │   ├── routes/          # 路由定义
│   │   │   │   └── routes.go
│   │   │   └── validators/      # 请求验证器
│   │   │       ├── auth.go
│   │   │       ├── user.go
│   │   │       ├── task.go
│   │   │       └── order.go
│   │   └── docs/                # API文档
│   │       └── swagger.yaml
│   ├── config/                   # 配置管理
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── redis.go
│   │   └── jwt.go
│   ├── models/                   # 数据模型
│   │   ├── user.go
│   │   ├── task.go
│   │   ├── order.go
│   │   ├── payment.go
│   │   ├── review.go
│   │   └── message.go
│   ├── repository/               # 数据访问层
│   │   ├── interfaces/          # 仓库接口
│   │   │   ├── user.go
│   │   │   ├── task.go
│   │   │   ├── order.go
│   │   │   └── payment.go
│   │   ├── mysql/               # MySQL实现
│   │   │   ├── user.go
│   │   │   ├── task.go
│   │   │   ├── order.go
│   │   │   └── payment.go
│   │   └── redis/               # Redis实现
│   │       ├── session.go
│   │       └── cache.go
│   ├── service/                  # 业务逻辑层
│   │   ├── interfaces/          # 服务接口
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   ├── task.go
│   │   │   ├── order.go
│   │   │   └── payment.go
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── task.go
│   │   ├── order.go
│   │   └── payment.go
│   ├── pkg/                      # 内部共享包
│   │   ├── utils/               # 工具函数
│   │   │   ├── crypto.go
│   │   │   ├── validator.go
│   │   │   ├── response.go
│   │   │   └── file.go
│   │   ├── logger/              # 日志包
│   │   │   └── logger.go
│   │   ├── errors/              # 错误处理
│   │   │   └── errors.go
│   │   └── queue/               # 消息队列
│   │       ├── rabbitmq.go
│   │       └── producer.go
│   └── domain/                   # 领域模型
│       ├── user/
│       │   ├── entity.go
│       │   ├── repository.go
│       │   └── service.go
│       ├── task/
│       │   ├── entity.go
│       │   ├── repository.go
│       │   └── service.go
│       └── order/
│           ├── entity.go
│           ├── repository.go
│           └── service.go
├── pkg/                          # 公共包
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── mysql.go
│   │   └── redis.go
│   ├── middleware/
│   │   ├── recovery.go
│   │   └── request_id.go
│   └── validator/
│       └── validator.go
├── scripts/                      # 脚本文件
│   ├── migrate.go               # 数据库迁移
│   ├── seed.go                  # 数据填充
│   └── build.sh                 # 构建脚本
├── configs/                      # 配置文件
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── docs/                         # 文档
│   ├── api.md
│   └── deployment.md
├── deployments/                  # 部署文件
│   ├── docker/
│   │   └── Dockerfile
│   ├── kubernetes/
│   │   └── deployment.yaml
│   └── nginx/
│       └── nginx.conf
├── tests/                        # 测试文件
│   ├── integration/
│   ├── unit/
│   └── mock/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 2. 核心代码结构设计

### 2.1 主程序入口

#### cmd/server/main.go
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    "task-platform-api/internal/api/v1/routes"
    "task-platform-api/internal/config"
    "task-platform-api/pkg/database"
    "task-platform-api/pkg/logger"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 初始化日志
    logger.Init(cfg.Log)
    
    // 初始化数据库
    db, err := database.InitMySQL(cfg.Database)
    if err != nil {
        logger.Fatal("Failed to initialize database", err)
    }
    
    // 初始化Redis
    redis, err := database.InitRedis(cfg.Redis)
    if err != nil {
        logger.Fatal("Failed to initialize redis", err)
    }
    
    // 设置Gin模式
    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // 创建路由
    router := routes.SetupRouter(db, redis)
    
    // 创建HTTP服务器
    srv := &http.Server{
        Addr:         ":" + cfg.Server.Port,
        Handler:      router,
        ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
    }
    
    // 启动服务器
    go func() {
        logger.Info("Server starting on port " + cfg.Server.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed to start", err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    logger.Info("Server shutting down...")
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", err)
    }
    
    logger.Info("Server exited")
}
```

### 2.2 配置管理

#### internal/config/config.go
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
    OSS      OSSConfig      `mapstructure:"oss"`
    SMS      SMSConfig      `mapstructure:"sms"`
    Email    EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
    Port         string `mapstructure:"port"`
    Mode         string `mapstructure:"mode"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host        string `mapstructure:"host"`
    Port        int    `mapstructure:"port"`
    Username    string `mapstructure:"username"`
    Password    string `mapstructure:"password"`
    Database    string `mapstructure:"database"`
    MaxIdleConn int    `mapstructure:"max_idle_conn"`
    MaxOpenConn int    `mapstructure:"max_open_conn"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime int    `mapstructure:"expire_time"`
}

type LogConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
    Output string `mapstructure:"output"`
}

type OSSConfig struct {
    Endpoint        string `mapstructure:"endpoint"`
    AccessKeyID     string `mapstructure:"access_key_id"`
    AccessKeySecret string `mapstructure:"access_key_secret"`
    Bucket          string `mapstructure:"bucket"`
}

type SMSConfig struct {
    Provider string `mapstructure:"provider"`
    APIKey   string `mapstructure:"api_key"`
    Secret   string `mapstructure:"secret"`
}

type EmailConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    From     string `mapstructure:"from"`
}

func Load(env ...string) (*Config, error) {
    v := viper.New()
    
    // 设置配置文件路径
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath("./configs")
    v.AddConfigPath(".")
    
    // 设置环境变量前缀
    v.SetEnvPrefix("TASK_PLATFORM")
    v.AutomaticEnv()
    
    // 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 2.3 路由设置

#### internal/api/v1/routes/routes.go
```go
package routes

import (
    "task-platform-api/internal/api/v1/handlers"
    "task-platform-api/internal/api/v1/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Handler struct {
    AuthHandler    *handlers.AuthHandler
    UserHandler    *handlers.UserHandler
    TaskHandler    *handlers.TaskHandler
    OrderHandler   *handlers.OrderHandler
    PaymentHandler *handlers.PaymentHandler
    ReviewHandler  *handlers.ReviewHandler
    MessageHandler *handlers.MessageHandler
}

func SetupRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
    router := gin.New()
    
    // 全局中间件
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery())
    router.Use(middleware.CORS())
    router.Use(middleware.RequestID())
    
    // 初始化处理器
    h := &Handler{
        AuthHandler:    handlers.NewAuthHandler(db, redis),
        UserHandler:    handlers.NewUserHandler(db, redis),
        TaskHandler:    handlers.NewTaskHandler(db, redis),
        OrderHandler:   handlers.NewOrderHandler(db, redis),
        PaymentHandler: handlers.NewPaymentHandler(db, redis),
        ReviewHandler:  handlers.NewReviewHandler(db, redis),
        MessageHandler: handlers.NewMessageHandler(db, redis),
    }
    
    // API路由组
    v1 := router.Group("/api/v1")
    {
        // 无需认证的路由
        public := v1.Group("/")
        {
            public.POST("/auth/register", h.AuthHandler.Register)
            public.POST("/auth/login", h.AuthHandler.Login)
            public.POST("/auth/refresh", h.AuthHandler.RefreshToken)
            public.POST("/auth/send-sms", h.AuthHandler.SendSMS)
            public.GET("/tasks", h.TaskHandler.ListTasks)
            public.GET("/tasks/:id", h.TaskHandler.GetTask)
            public.GET("/categories", h.TaskHandler.ListCategories)
        }
        
        // 需要认证的路由
        protected := v1.Group("/")
        protected.Use(middleware.JWTAuth())
        {
            // 用户相关
            user := protected.Group("/user")
            {
                user.GET("/profile", h.UserHandler.GetProfile)
                user.PUT("/profile", h.UserHandler.UpdateProfile)
                user.POST("/upload-avatar", h.UserHandler.UploadAvatar)
                user.GET("/skills", h.UserHandler.GetUserSkills)
                user.POST("/skills", h.UserHandler.AddUserSkill)
                user.PUT("/skills/:id", h.UserHandler.UpdateUserSkill)
                user.DELETE("/skills/:id", h.UserHandler.DeleteUserSkill)
            }
            
            // 任务相关
            task := protected.Group("/tasks")
            {
                task.POST("/", h.TaskHandler.CreateTask)
                task.PUT("/:id", h.TaskHandler.UpdateTask)
                task.DELETE("/:id", h.TaskHandler.DeleteTask)
                task.POST("/:id/apply", h.TaskHandler.ApplyTask)
                task.GET("/my-published", h.TaskHandler.GetMyPublishedTasks)
                task.GET("/my-applied", h.TaskHandler.GetMyAppliedTasks)
            }
            
            // 订单相关
            order := protected.Group("/orders")
            {
                order.GET("/", h.OrderHandler.ListOrders)
                order.GET("/:id", h.OrderHandler.GetOrder)
                order.POST("/:id/accept", h.OrderHandler.AcceptOrder)
                order.POST("/:id/complete", h.OrderHandler.CompleteOrder)
                order.POST("/:id/confirm", h.OrderHandler.ConfirmOrder)
                order.POST("/:id/cancel", h.OrderHandler.CancelOrder)
            }
            
            // 支付相关
            payment := protected.Group("/payments")
            {
                payment.POST("/create", h.PaymentHandler.CreatePayment)
                payment.POST("/callback", h.PaymentHandler.PaymentCallback)
                payment.GET("/wallet", h.PaymentHandler.GetWallet)
                payment.POST("/withdraw", h.PaymentHandler.Withdraw)
            }
            
            // 评价相关
            review := protected.Group("/reviews")
            {
                review.POST("/", h.ReviewHandler.CreateReview)
                review.GET("/received", h.ReviewHandler.GetReceivedReviews)
                review.GET("/given", h.ReviewHandler.GetGivenReviews)
            }
            
            // 消息相关
            message := protected.Group("/messages")
            {
                message.GET("/", h.MessageHandler.ListMessages)
                message.POST("/", h.MessageHandler.SendMessage)
                message.PUT("/:id/read", h.MessageHandler.MarkAsRead)
                message.GET("/notifications", h.MessageHandler.GetNotifications)
                message.PUT("/notifications/:id/read", h.MessageHandler.MarkNotificationAsRead)
            }
        }
        
        // 管理员路由
        admin := v1.Group("/admin")
        admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
        {
            admin.GET("/users", h.UserHandler.ListUsers)
            admin.PUT("/users/:id/status", h.UserHandler.UpdateUserStatus)
            admin.GET("/tasks/pending", h.TaskHandler.GetPendingTasks)
            admin.PUT("/tasks/:id/approve", h.TaskHandler.ApproveTask)
            admin.GET("/complaints", h.ReviewHandler.ListComplaints)
            admin.PUT("/complaints/:id/handle", h.ReviewHandler.HandleComplaint)
        }
    }
    
    return router
}
```

### 2.4 数据模型

#### internal/models/user.go
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID          uint64         `json:"id" gorm:"primaryKey"`
    Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
    Phone       string         `json:"phone" gorm:"uniqueIndex;size:20"`
    Email       string         `json:"email" gorm:"uniqueIndex;size:100"`
    PasswordHash string       `json:"-" gorm:"size:255;not null"`
    Nickname    string         `json:"nickname" gorm:"size:50;not null"`
    Avatar      string         `json:"avatar" gorm:"size:255"`
    Gender      int8           `json:"gender" gorm:"default:0;comment:0-未知,1-男,2-女"`
    Bio         string         `json:"bio" gorm:"type:text"`
    Status      int8           `json:"status" gorm:"default:1;comment:0-禁用,1-正常,2-待审核"`
    CreditScore float32        `json:"credit_score" gorm:"type:decimal(3,1);default:5.0;comment:信用评分(0-10)"`
    Level       int            `json:"level" gorm:"default:1"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserProfile struct {
    ID         uint64    `json:"id" gorm:"primaryKey"`
    UserID     uint64    `json:"user_id" gorm:"uniqueIndex;not null"`
    RealName   string    `json:"real_name" gorm:"size:50"`
    IDCard     string    `json:"id_card" gorm:"size:18"`
    Skills     string    `json:"skills" gorm:"type:json"`
    Experience string    `json:"experience" gorm:"type:text"`
    Location   string    `json:"location" gorm:"size:200"`
    WorkType   string    `json:"work_type" gorm:"size:50"`
    Balance    float64   `json:"balance" gorm:"type:decimal(10,2);default:0.00"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    User       User      `json:"user" gorm:"foreignKey:UserID"`
}

type UserSkill struct {
    ID           uint64    `json:"id" gorm:"primaryKey"`
    UserID       uint64    `json:"user_id" gorm:"index;not null"`
    CategoryID   uint64    `json:"category_id" gorm:"index;not null"`
    SkillName    string    `json:"skill_name" gorm:"size:100;not null"`
    Proficiency int       `json:"proficiency" gorm:"comment:熟练度1-5"`
    CreatedAt    time.Time `json:"created_at"`
    
    User         User      `json:"user" gorm:"foreignKey:UserID"`
    Category     TaskCategory `json:"category" gorm:"foreignKey:CategoryID"`
}

func (User) TableName() string {
    return "users"
}

func (UserProfile) TableName() string {
    return "user_profiles"
}

func (UserSkill) TableName() string {
    return "user_skills"
}
```

## 3. 依赖管理

### go.mod
```go
module task-platform-api

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/spf13/viper v1.16.0
    github.com/go-redis/redis/v8 v8.11.5
    gorm.io/gorm v1.25.4
    gorm.io/driver/mysql v1.5.1
    github.com/streadway/amqp v1.1.0
    github.com/aliyun/aliyun-oss-go-sdk v2.2.7+incompatible
    github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.661
    github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.661
    github.com/sirupsen/logrus v1.9.3
    github.com/go-playground/validator/v10 v10.15.3
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/swag v1.16.1
    github.com/golang/mock v1.6.0
    github.com/stretchr/testify v1.8.4
    go.uber.org/zap v1.25.0
)
```

## 4. 构建和部署

### Makefile
```makefile
.PHONY: build run test clean docker-build docker-run

# 变量定义
APP_NAME=task-platform-api
VERSION=1.0.0
BUILD_DIR=build

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/server/main.go
	@echo "Build completed!"

# 运行
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/server/main.go

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 清理
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean

# 生成Swagger文档
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o internal/api/docs

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	@go run scripts/migrate.go

# 数据填充
seed:
	@echo "Seeding database..."
	@go run scripts/seed.go

# Docker构建
docker-build:
	@echo "Building Docker image..."
	@docker build -f deployments/docker/Dockerfile -t $(APP_NAME):$(VERSION) .

# Docker运行
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run

# 生成mock文件
mock:
	@echo "Generating mock files..."
	@mockgen -source=internal/repository/interfaces/user.go -destination=tests/mock/user_mock.go
	@mockgen -source=internal/repository/interfaces/task.go -destination=tests/mock/task_mock.go
```

## 5. 中间件设计

### 认证中间件
### 日志中间件
### 限流中间件
### CORS中间件