package main

import (
	"fmt"
	"log"
	
	"task-platform-api/internal/config"
	"task-platform-api/pkg/database"
)

func main() {
	fmt.Println("=== 测试数据库连接 ===")
	
	// 加载配置
	cfg, err := config.Load("./configs/config-optimized.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	
	fmt.Printf("数据库配置: Host=%s, Port=%d, User=%s, DB=%s\n",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Database)
	
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
		ConnMaxLifetime: 3600,
	}
	
	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	
	fmt.Println("✓ 数据库连接成功!")
	
	// 测试查询
	var count int64
	db.Raw("SELECT COUNT(*) FROM users").Scan(&count)
	fmt.Printf("用户表记录数: %d\n", count)
	
	db.Raw("SELECT COUNT(*) FROM tasks").Scan(&count)
	fmt.Printf("任务表记录数: %d\n", count)
	
	db.Raw("SELECT COUNT(*) FROM task_categories").Scan(&count)
	fmt.Printf("任务分类记录数: %d\n", count)
}
