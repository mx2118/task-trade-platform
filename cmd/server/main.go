package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"task-platform-api/internal/api/v1/handlers"
	"task-platform-api/internal/api/v1/routes"
	"task-platform-api/internal/config"
	"task-platform-api/internal/services"
	"task-platform-api/pkg/logger"
	"task-platform-api/pkg/database"
	"task-platform-api/pkg/redis"
)

var (
	configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
)

func initLogger(config config.LogConfig) (*zap.Logger, error) {
	loggerConfig := logger.Config{
		Level:      config.Level,
		Format:     config.Format,
		OutputPath: config.FilePath,
	}
	return logger.New(loggerConfig)
}

func main() {
	flag.Parse()

	// 加载配置
	if *configPath == "./configs/config.yaml" {
		*configPath = "./configs/config-optimized.yaml"
	}
	
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
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		Charset:         "utf8mb4",
		MaxIdleConns:    cfg.Database.MaxIdleConn,
		MaxOpenConns:    cfg.Database.MaxOpenConn,
		ConnMaxLifetime: 3600, // 1 hour
	}
	db, err := database.New(dbConfig)
	if err != nil {
		zapLogger.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// 初始化Redis
	redisConfig := redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	}
	rdb, err := redis.New(redisConfig)
	if err != nil {
		zapLogger.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer rdb.Close()

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建处理器
	authHandler := handlers.NewAuthHandler(db, rdb, cfg, zapLogger)
	paymentService := services.NewPaymentService(db, nil) // 暂时不传入支付客户端
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	
	// 创建路由
	router := gin.New()
	routes.SetupRoutes(router, authHandler, paymentHandler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		zapLogger.Info("启动服务器", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("正在关闭服务器...")

	// 设置5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("服务器强制关闭", zap.Error(err))
	}

	zapLogger.Info("服务器已关闭")
}