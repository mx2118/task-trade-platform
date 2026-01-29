package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var startTime = time.Now()

// 数据库连接配置
const (
	dbUser     = "root"
	dbPassword = "root123456"
	dbHost     = "localhost"
	dbPort     = "3306"
	dbName     = "task_platform"
)

// API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 任务结构
type Task struct {
	TaskID          int64   `json:"task_id"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	Amount          float64 `json:"amount"`
	Status          int     `json:"status"`
	ViewCount       int     `json:"view_count"`
	ApplyCount      int     `json:"apply_count"`
	CategoryID      int64   `json:"category_id"`
	CategoryName    string  `json:"category_name,omitempty"`
	PublisherID     int64   `json:"publisher_id"`
	PublisherName   string  `json:"publisher_name,omitempty"`
	PublisherAvatar string  `json:"publisher_avatar,omitempty"`
	Deadline        string  `json:"deadline"`
	CreateTime      string  `json:"create_time"`
}

// 任务分类结构
type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	TaskCount   int    `json:"task_count"`
}

// 用户结构
type User struct {
	UserID      int64   `json:"user_id"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	CreditScore float64 `json:"credit_score"`
	Level       int     `json:"level"`
}

func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置字符编码
	_, err = db.Exec("SET NAMES utf8mb4")
	if err != nil {
		return fmt.Errorf("设置字符编码失败: %v", err)
	}

	// 设置连接池
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("✅ 数据库连接成功!")
	return nil
}

func main() {
	// 初始化数据库
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 设置CORS中间件
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Content-Type", "application/json")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	// API路由
	http.HandleFunc("/api/", corsMiddleware(apiHandler))
	http.HandleFunc("/api/tasks", corsMiddleware(tasksHandler))
	http.HandleFunc("/api/tasks/list", corsMiddleware(tasksHandler))
	http.HandleFunc("/api/tasks/stats", corsMiddleware(taskStatsHandler))
	http.HandleFunc("/api/categories", corsMiddleware(categoriesHandler))
	http.HandleFunc("/api/users/stats", corsMiddleware(userStatsHandler))
	http.HandleFunc("/health", corsMiddleware(healthHandler))

	fmt.Println("🚀 任务交易平台API启动中...")
	fmt.Printf("📡 服务器地址: http://49.234.39.189:8080\n")
	fmt.Printf("🔧 API接口: http://49.234.39.189:8080/api/\n")
	fmt.Printf("💾 数据库: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	fmt.Println("=====================================")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Code:    200,
		Message: "任务交易平台API运行正常",
		Data: map[string]interface{}{
			"service":   "Task Trade Platform",
			"version":   "2.0.0",
			"time":      time.Now(),
			"database":  "Connected to MySQL",
			"endpoints": []string{"/api/tasks", "/api/categories", "/api/users/stats", "/health"},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	status := r.URL.Query().Get("status")
	categoryID := r.URL.Query().Get("category_id")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// 默认分页参数
	pageNum := 1
	pageSize := 20
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			pageSize = l
		}
	}
	offset := (pageNum - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE t.deleted_at IS NULL"
	args := []interface{}{}

	if status != "" {
		whereClause += " AND t.status = ?"
		args = append(args, status)
	}
	if categoryID != "" {
		whereClause += " AND t.category_id = ?"
		args = append(args, categoryID)
	}

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM tasks t ` + whereClause
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("查询任务总数失败: %v", err)
		response := APIResponse{Code: 500, Message: "查询任务总数失败"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建SQL查询 - JOIN用户表获取发布者信息
	query := `
		SELECT t.task_id, t.title, t.content, t.amount, t.status, 
		       t.view_count, t.apply_count, t.category_id, t.publisher_id,
		       t.deadline, t.create_time, 
		       c.name as category_name,
		       u.nickname as publisher_name, u.avatar as publisher_avatar
		FROM tasks t
		LEFT JOIN task_categories c ON t.category_id = c.id
		LEFT JOIN users u ON t.publisher_id = u.user_id
		` + whereClause + `
		ORDER BY t.create_time DESC 
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("查询任务失败: %v", err)
		response := APIResponse{Code: 500, Message: "查询任务失败"}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		var categoryName, publisherName, publisherAvatar sql.NullString
		var deadline, createTime time.Time

		err := rows.Scan(
			&task.TaskID, &task.Title, &task.Content, &task.Amount, &task.Status,
			&task.ViewCount, &task.ApplyCount, &task.CategoryID, &task.PublisherID,
			&deadline, &createTime, &categoryName, &publisherName, &publisherAvatar,
		)
		if err != nil {
			log.Printf("扫描任务数据失败: %v", err)
			continue
		}

		if categoryName.Valid {
			task.CategoryName = categoryName.String
		}
		if publisherName.Valid {
			task.PublisherName = publisherName.String
		}
		if publisherAvatar.Valid {
			task.PublisherAvatar = publisherAvatar.String
		}
		task.Deadline = deadline.Format("2006-01-02 15:04:05")
		task.CreateTime = createTime.Format("2006-01-02 15:04:05")

		tasks = append(tasks, task)
	}

	response := APIResponse{
		Code:    200,
		Message: "获取任务列表成功",
		Data: map[string]interface{}{
			"total": total,
			"page":  pageNum,
			"limit": pageSize,
			"tasks": tasks,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT c.id, c.name, c.icon, c.description, 
		       COUNT(t.task_id) as task_count
		FROM task_categories c
		LEFT JOIN tasks t ON c.id = t.category_id AND t.deleted_at IS NULL AND t.status = 1
		WHERE c.status = 1
		GROUP BY c.id, c.name, c.icon, c.description
		ORDER BY c.sort_order ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询分类失败: %v", err)
		response := APIResponse{Code: 500, Message: "查询分类失败"}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var cat Category
		var description sql.NullString

		err := rows.Scan(&cat.ID, &cat.Name, &cat.Icon, &description, &cat.TaskCount)
		if err != nil {
			log.Printf("扫描分类数据失败: %v", err)
			continue
		}

		if description.Valid {
			cat.Description = description.String
		}

		categories = append(categories, cat)
	}

	response := APIResponse{
		Code:    200,
		Message: "获取分类列表成功",
		Data: map[string]interface{}{
			"total":      len(categories),
			"categories": categories,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func userStatsHandler(w http.ResponseWriter, r *http.Request) {
	var totalUsers, activeUsers int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE status = 1").Scan(&totalUsers)
	db.QueryRow("SELECT COUNT(DISTINCT publisher_id) FROM tasks WHERE status IN (1,2,3)").Scan(&activeUsers)

	response := APIResponse{
		Code:    200,
		Message: "获取用户统计成功",
		Data: map[string]interface{}{
			"total_users":  totalUsers,
			"active_users": activeUsers,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func taskStatsHandler(w http.ResponseWriter, r *http.Request) {
	var total, draft, pending, inProgress, reviewing, completed, cancelled int
	
	// 获取分类ID参数
	categoryID := r.URL.Query().Get("category_id")
	
	// 构建基础查询条件
	baseWhere := "WHERE deleted_at IS NULL"
	var args []interface{}
	
	// 如果有分类筛选，添加到条件中
	if categoryID != "" {
		baseWhere += " AND category_id = ?"
		args = append(args, categoryID)
	}
	
	// 查询各状态的任务数量
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere, args...).Scan(&total)
	
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 0", args...).Scan(&draft)
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 1", args...).Scan(&pending)
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 2", args...).Scan(&inProgress)
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 3", args...).Scan(&reviewing)
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 4", args...).Scan(&completed)
	db.QueryRow("SELECT COUNT(*) FROM tasks "+baseWhere+" AND status = 5", args...).Scan(&cancelled)

	response := APIResponse{
		Code:    200,
		Message: "获取任务统计成功",
		Data: map[string]interface{}{
			"total":       total,
			"draft":       draft,
			"pending":     pending,
			"in_progress": inProgress,
			"reviewing":   reviewing,
			"completed":   completed,
			"cancelled":   cancelled,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var dbStatus string
	if err := db.Ping(); err != nil {
		dbStatus = "disconnected"
	} else {
		dbStatus = "connected"
	}

	response := map[string]interface{}{
		"status":        "ok",
		"timestamp":     time.Now(),
		"version":       "2.0.0",
		"uptime":        time.Since(startTime).String(),
		"database":      dbStatus,
		"database_type": "MySQL",
	}

	json.NewEncoder(w).Encode(response)
}
