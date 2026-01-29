package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    go_redis "github.com/go-redis/redis/v8"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "task-platform-api/internal/config"
    "task-platform-api/internal/api/v1/handlers"
    "task-platform-api/internal/api/v1/middleware"
    "task-platform-api/internal/api/v1/routes"
    "task-platform-api/internal/models"
)

var (
    configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
)

func main() {
    flag.Parse()

    // 加载配置
    cfg, err := config.Load(*configPath)
    if err != nil {
        log.Fatalf("加载配置文件失败: %v", err)
    }

    // 初始化日志
    zapLogger, err := initLogger(cfg.Log)
    if err != nil {
        log.Fatalf("初始化日志失败: %v", err)
    }
    defer zapLogger.Sync()

    // 初始化数据库
    db, err := initDatabase(cfg.Database)
    if err != nil {
        zapLogger.Fatal("初始化数据库失败", zap.Error(err))
    }

    // 初始化Redis
    rdb, err := initRedis(cfg.Redis)
    if err != nil {
        zapLogger.Fatal("初始化Redis失败", zap.Error(err))
    }

    // 设置Gin模式
    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // 创建路由
    router := routes.SetupRouter(db, rdb, cfg, zapLogger)

    // 创建HTTP服务器
    srv := &http.Server{
        Addr:         ":" + cfg.Server.Port,
        Handler:      router,
        ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
    }

    // 启动服务器
    go func() {
        zapLogger.Info("服务器启动", zap.String("port", cfg.Server.Port))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            zapLogger.Fatal("服务器启动失败", zap.Error(err))
        }
    }()

    // 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    zapLogger.Info("服务器正在关闭...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        zapLogger.Fatal("服务器强制关闭", zap.Error(err))
    }

    zapLogger.Info("服务器已关闭")
}

// initDatabase 初始化数据库连接
func initDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
    )

    var gormLogger logger.Interface
    switch cfg.LogLevel {
    case "debug":
        gormLogger = logger.Default.LogMode(logger.Info)
    case "warn":
        gormLogger = logger.Default.LogMode(logger.Warn)
    case "error":
        gormLogger = logger.Default.LogMode(logger.Error)
    default:
        gormLogger = logger.Default.LogMode(logger.Info)
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
    if err != nil {
        return nil, fmt.Errorf("连接数据库失败: %w", err)
    }

    // 设置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
    }

    sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

    // 自动迁移数据表
    err = db.AutoMigrate(
        &models.User{},
        &models.UserSession{},
        &models.UserCredit{},
        &models.Task{},
        &models.TaskStage{},
        &models.TaskDelivery{},
        &models.Trade{},
        &models.Settlement{},
        &models.Refund{},
        &models.Violation{},
        &models.Complaint{},
        &models.Notification{},
    )
    if err != nil {
        return nil, fmt.Errorf("数据库迁移失败: %w", err)
    }

    return db, nil
}

// initRedis 初始化Redis连接
func initRedis(cfg config.RedisConfig) (*go_redis.Client, error) {
    rdb := go_redis.NewClient(&go_redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
        PoolSize: cfg.PoolSize,
    })

    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("连接Redis失败: %w", err)
    }

    return rdb, nil
}

// initLogger 初始化日志
func initLogger(cfg config.LogConfig) (*zap.Logger, error) {
    var zapConfig zap.Config
    if cfg.Format == "json" {
        zapConfig = zap.NewProductionConfig()
    } else {
        zapConfig = zap.NewDevelopmentConfig()
    }

    // 设置日志级别
    switch cfg.Level {
    case "debug":
        zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
    case "info":
        zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    case "warn":
        zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
    case "error":
        zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
    default:
        zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    }

    // 设置输出路径
    if cfg.Output == "file" && cfg.FilePath != "" {
        zapConfig.OutputPaths = []string{cfg.FilePath}
        zapConfig.ErrorOutputPaths = []string{cfg.FilePath}
    }

    return zapConfig.Build()
}