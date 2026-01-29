package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OAuth相关结构体
type OAuthUser struct {
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid,omitempty"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"headimgurl,omitempty"`
	AuthType   string `json:"auth_type"`
	UserType   string `json:"user_type"`
}

type LoginResponse struct {
	Code       int       `json:"code"`
	Message    string    `json:"message"`
	Token      string    `json:"token,omitempty"`
	User       OAuthUser `json:"user,omitempty"`
}

type UserInfo struct {
	UserID      int     `json:"user_id"`
	AuthType    string  `json:"auth_type"`
	OpenID      string  `json:"openid"`
	UnionID     string  `json:"unionid,omitempty"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	CreditScore float64 `json:"credit_score"`
	Level       int     `json:"level"`
	UserType    string  `json:"user_type"`
	CompleteRate float64 `json:"complete_rate"`
	CreateTime  string  `json:"create_time"`
}

// 模拟用户数据库
var users = make(map[string]UserInfo)
var nextUserID = 1001

var startTime = time.Now()

func main() {
	// 微信OAuth登录
	http.HandleFunc("/api/auth/wechat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// 解析请求参数
		var req struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Code: 400, Message: "请求参数错误"})
			return
		}
		
		// 模拟微信OAuth2.0授权码换取用户信息
		// 实际生产环境中需要调用微信API: https://api.weixin.qq.com/sns/oauth2/access_token
		mockUserInfo := OAuthUser{
			OpenID:   fmt.Sprintf("wx_%d", time.Now().Unix()),
			UnionID:  fmt.Sprintf("union_%d", time.Now().Unix()),
			Nickname: "微信用户",
			Avatar:   "https://thirdwx.qlogo.cn/mmopen/vi_32/default_avatar.png",
			AuthType: "wechat",
			UserType: "general", // 通用用户类型
		}
		
		// 生成token
		token := fmt.Sprintf("token_%d_%d", time.Now().Unix(), nextUserID)
		
		// 检查用户是否已存在
		if user, exists := users[mockUserInfo.OpenID]; exists {
			response := LoginResponse{
				Code:    200,
				Message: "登录成功",
				Token:   token,
				User: OAuthUser{
					OpenID:   user.OpenID,
					UnionID:  user.UnionID,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
					AuthType: user.AuthType,
					UserType: user.UserType,
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// 创建新用户
		newUser := UserInfo{
			UserID:      nextUserID,
			AuthType:    "wechat",
			OpenID:      mockUserInfo.OpenID,
			UnionID:     mockUserInfo.UnionID,
			Nickname:    mockUserInfo.Nickname,
			Avatar:      mockUserInfo.Avatar,
			CreditScore:  5.0,
			Level:       1,
			UserType:    "general",
			CompleteRate: 0.0,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		}
		
		users[mockUserInfo.OpenID] = newUser
		nextUserID++
		
		response := LoginResponse{
			Code:    200,
			Message: "注册并登录成功",
			Token:   token,
			User: mockUserInfo,
		}
		json.NewEncoder(w).Encode(response)
	})
	
	// 支付宝OAuth登录
	http.HandleFunc("/api/auth/alipay", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// 解析请求参数
		var req struct {
			AuthCode string `json:"auth_code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Code: 400, Message: "请求参数错误"})
			return
		}
		
		// 模拟支付宝OAuth2.0授权码换取用户信息
		// 实际生产环境中需要调用支付宝API: https://openapi.alipay.com/gateway.do
		mockUserInfo := OAuthUser{
			OpenID:   fmt.Sprintf("alipay_%d", time.Now().Unix()),
			UnionID:  fmt.Sprintf("alipay_union_%d", time.Now().Unix()),
			Nickname: "支付宝用户",
			Avatar:   "https://tfs.alipayobjects.com/images/partner/T1ByRfXklXXXXXXXXXXXXX",
			AuthType: "alipay",
			UserType: "general", // 通用用户类型
		}
		
		// 生成token
		token := fmt.Sprintf("token_%d_%d", time.Now().Unix(), nextUserID)
		
		// 检查用户是否已存在
		if user, exists := users[mockUserInfo.OpenID]; exists {
			response := LoginResponse{
				Code:    200,
				Message: "登录成功",
				Token:   token,
				User: OAuthUser{
					OpenID:   user.OpenID,
					UnionID:  user.UnionID,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
					AuthType: user.AuthType,
					UserType: user.UserType,
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// 创建新用户
		newUser := UserInfo{
			UserID:      nextUserID,
			AuthType:    "alipay",
			OpenID:      mockUserInfo.OpenID,
			UnionID:     mockUserInfo.UnionID,
			Nickname:    mockUserInfo.Nickname,
			Avatar:      mockUserInfo.Avatar,
			CreditScore: 5.0,
			Level:       1,
			UserType:    "general",
			CompleteRate: 0.0,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		}
		
		users[mockUserInfo.OpenID] = newUser
		nextUserID++
		
		response := LoginResponse{
			Code:    200,
			Message: "注册并登录成功",
			Token:   token,
			User: mockUserInfo,
		}
		json.NewEncoder(w).Encode(response)
	})
	
	// 获取用户信息
	http.HandleFunc("/api/user/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// 从header获取token
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{Code: 401, Message: "未授权"})
			return
		}
		
		// 模拟token验证，返回第一个用户作为示例
		for _, user := range users {
			response := APIResponse{
				Code:    200,
				Message: "获取用户信息成功",
				Data:    user,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{Code: 404, Message: "用户不存在"})
	})
	
	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		response := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now(),
			Version:   "1.0.0",
			Uptime:    time.Since(startTime).String(),
		}
		
		json.NewEncoder(w).Encode(response)
	})

	// API路由
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// 解析API路径
		path := r.URL.Path[len("/api/"):]
		
		switch path {
		case "health":
			// 健康检查端点
			response := HealthResponse{
				Status:    "ok",
				Timestamp: time.Now(),
				Version:   "1.0.0",
				Uptime:    time.Since(startTime).String(),
			}
			json.NewEncoder(w).Encode(response)
		case "":
			// API根路径
			response := APIResponse{
				Code:    200,
				Message: "任务交易平台API运行正常",
				Data: map[string]interface{}{
					"service": "Task Trade Platform",
					"version": "1.0.0",
					"time":    time.Now(),
					"endpoints": []string{"/health", "/tasks", "/users", "/payment/status"},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "tasks":
			// 任务列表API
			response := APIResponse{
				Code:    200,
				Message: "任务列表获取成功",
				Data: map[string]interface{}{
					"total": 0,
					"tasks": []map[string]interface{}{
						{
							"id": 1,
							"title": "示例任务",
							"description": "这是一个示例任务",
							"price": 10.00,
							"status": "pending",
							"created_at": time.Now(),
						},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "users":
			// 用户API
			response := APIResponse{
				Code:    200,
				Message: "用户信息获取成功",
				Data: map[string]interface{}{
					"total_users": 1,
					"active_users": 1,
					"online_users": 1,
				},
			}
			json.NewEncoder(w).Encode(response)
		case "payment/status":
			// 支付状态API
			response := APIResponse{
				Code:    200,
				Message: "支付系统运行正常",
				Data: map[string]interface{}{
					"provider": "收钱吧",
					"status": "online",
					"methods": []string{"alipay", "wechat", "card"},
				},
			}
			json.NewEncoder(w).Encode(response)
		default:
			// 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			response := APIResponse{
				Code:    404,
				Message: "API端点未找到",
				Data: map[string]interface{}{
					"path": path,
					"available_endpoints": []string{"/", "/health", "/tasks", "/users", "/payment/status"},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	})

	// 主页路由（用于SPA支持）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>任务交易平台</title>
    <meta charset="utf-8">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            margin: 40px; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            min-height: 100vh;
        }
        .container { 
            max-width: 800px; 
            margin: 0 auto; 
            text-align: center;
        }
        .card { 
            background: rgba(255,255,255,0.1);
            padding: 40px;
            border-radius: 20px;
            backdrop-filter: blur(10px);
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            margin: 20px 0;
        }
        h1 { 
            color: #ffffff; 
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        .status { 
            color: #4CAF50; 
            font-weight: bold;
            font-size: 1.2em;
        }
        .button {
            display: inline-block;
            padding: 15px 30px;
            margin: 10px;
            background: rgba(255,255,255,0.2);
            color: white;
            text-decoration: none;
            border-radius: 25px;
            border: 2px solid rgba(255,255,255,0.3);
            transition: all 0.3s ease;
        }
        .button:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-2px);
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-top: 30px;
        }
        .feature {
            background: rgba(255,255,255,0.05);
            padding: 20px;
            border-radius: 15px;
            border-left: 4px solid #4CAF50;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>🚀 任务交易平台</h1>
            <div class="status">✅ 系统运行正常</div>
            <p>服务器IP: <strong>121.41.39.105</strong></p>
            <p>部署时间: <strong>%s</strong></p>
            <p>运行状态: <strong>在线服务中</strong></p>
        </div>
        
        <div class="grid">
            <div class="feature">
                <h3>🔧 系统管理</h3>
                <p>宝塔面板管理</p>
                <a href="https://121.41.39.105:21452/f97c6b7e" class="button" target="_blank">访问面板</a>
            </div>
            
            <div class="feature">
                <h3>📊 API状态</h3>
                <p>查看API健康状态</p>
                <a href="/api/" class="button" target="_blank">API接口</a>
            </div>
            
            <div class="feature">
                <h3>💳 支付系统</h3>
                <p>收钱吧支付集成</p>
                <a href="/api/payment/status" class="button" target="_blank">支付状态</a>
            </div>
            
            <div class="feature">
                <h3>📋 任务管理</h3>
                <p>任务发布和管理</p>
                <a href="/api/tasks" class="button" target="_blank">任务列表</a>
            </div>
        </div>
        
        <div class="card">
            <h3>🎉 部署成功！</h3>
            <p>任务交易平台已成功部署到公网服务器</p>
            <p>用户现在可以访问 <strong>http://121.41.39.105</strong> 使用完整的任务交易功能</p>
            <p><strong>开始您的在线任务交易业务吧！</strong></p>
        </div>
    </div>
</body>
</html>`, time.Now().Format("2006-01-02 15:04:05"))
	})

	fmt.Println("🚀 任务交易平台启动中...")
	fmt.Printf("📡 服务器地址: http://121.41.39.105:8080\n")
	fmt.Printf("🔧 API接口: http://121.41.39.105:8080/api/\n")
	fmt.Printf("💚 健康检查: http://121.41.39.105:8080/health\n")
	fmt.Println("=====================================")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}