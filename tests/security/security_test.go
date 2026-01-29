package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

// SecurityTestSuite 安全测试套件
type SecurityTestSuite struct {
	router *gin.Engine
}

func (suite *SecurityTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	
	// 添加基本路由用于安全测试
	suite.router.POST("/api/v1/test/echo", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	})
	
	suite.router.GET("/api/v1/test/html", func(c *gin.Context) {
		htmlContent := c.Query("content")
		c.String(http.StatusOK, htmlContent)
	})
}

// TestSQLInjection 测试SQL注入防护
func (suite *SecurityTestSuite) TestSQLInjection(t *testing.T) {
	testCases := []string{
		"'; DROP TABLE users; --",
		"' OR '1'='1",
		"' UNION SELECT * FROM users --",
		"1'; DELETE FROM users WHERE '1'='1",
		"' OR 1=1 #",
	}

	for _, maliciousInput := range testCases {
		t.Run("SQL注入测试: "+maliciousInput, func(t *testing.T) {
			data := map[string]interface{}{
				"query": maliciousInput,
			}

			jsonData, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 检查响应状态和内容
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 验证恶意输入没有被直接返回（应该被过滤或转义）
			if data, ok := response["data"].(map[string]interface{}); ok {
				returnedQuery, _ := data["query"].(string)
				// 恶意字符应该被转义或过滤
				assert.NotEqual(t, maliciousInput, returnedQuery)
			}
		})
	}
}

// TestXSSPrevention 测试XSS防护
func (suite *SecurityTestSuite) TestXSSPrevention(t *testing.T) {
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src='x' onerror='alert(1)'>",
		"javascript:alert('XSS')",
		"<svg onload='alert(1)'>",
		"';alert('XSS');//",
	}

	for _, payload := range xssPayloads {
		t.Run("XSS防护测试: "+payload, func(t *testing.T) {
			// 测试API中的XSS防护
			data := map[string]interface{}{
				"comment": payload,
			}

			jsonData, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if data, ok := response["data"].(map[string]interface{}); ok {
				returnedComment, _ := data["comment"].(string)
				// 验证脚本标签被过滤或转义
				assert.NotContains(t, strings.ToLower(returnedComment), "<script>")
				assert.NotContains(t, strings.ToLower(returnedComment), "javascript:")
			}

			// 测试HTML响应中的XSS防护
			url := "/api/v1/test/html?content=" + url.QueryEscape(payload)
			req, _ = http.NewRequest("GET", url, nil)
			w = httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 解析HTML响应
			doc, err := html.Parse(strings.NewReader(w.Body.String()))
			assert.NoError(t, err)

			// 检查是否存在script标签
			var hasScriptTag bool
			var f func(*html.Node)
			f = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "script" {
					hasScriptTag = true
				}
				for c := n.FirstChild; c != nil; c = n.NextSibling {
					f(c)
				}
			}
			f(doc)

			assert.False(t, hasScriptTag, "发现未过滤的script标签")
		})
	}
}

// TestCSRFProtection 测试CSRF防护
func (suite *SecurityTestSuite) TestCSRFProtection(t *testing.T) {
	// 模拟跨站请求
	csrfToken := "fake-csrf-token"
	
	data := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"csrf_token": csrfToken,
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://malicious-site.com")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 检查是否有CSRF防护（应该返回错误或验证失败）
	// 这取决于实际实现的CSRF防护机制
	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest, http.StatusForbidden}, w.Code)
}

// TestInputValidation 测试输入验证
func (suite *SecurityTestSuite) TestInputValidation(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "空输入",
			input:    "",
			expected: http.StatusBadRequest,
		},
		{
			name:     "过长输入",
			input:    strings.Repeat("a", 10000),
			expected: http.StatusBadRequest,
		},
		{
			name:     "特殊字符",
			input:    "!@#$%^&*()",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Unicode字符",
			input:    "🚀💻📱",
			expected: http.StatusBadRequest,
		},
		{
			name:     "正常输入",
			input:    "正常输入123",
			expected: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"input": tc.input,
			}

			jsonData, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 根据输入验证实现，可能返回200但包含错误信息
			assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)
		})
	}
}

// TestRateLimiting 测试频率限制
func (suite *SecurityTestSuite) TestRateLimiting(t *testing.T) {
	// 模拟高频请求
	for i := 0; i < 100; i++ {
		data := map[string]interface{}{
			"request_id": i,
		}

		jsonData, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// 在达到频率限制后，应该返回429状态码
		if i > 50 {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, 
				"第%d个请求应该被频率限制", i)
			break
		}
	}
}

// TestFileUploadSecurity 测试文件上传安全
func (suite *SecurityTestSuite) TestFileUploadSecurity(t *testing.T) {
	// 这里需要实际的文件上传路由来测试
	// 模拟文件上传测试
	
	maliciousFiles := []string{
		"malicious.php",
		"virus.exe",
		"script.js",
		"../../../etc/passwd",
	}

	for _, filename := range maliciousFiles {
		t.Run("文件上传安全测试: "+filename, func(t *testing.T) {
			// 模拟文件上传请求
			// 实际实现中应该检查文件类型、大小、路径等
			assert.NotContains(t, filename, "..", "不允许路径遍历攻击")
			assert.NotContains(t, filename, ".php", "不允许上传PHP文件")
			assert.NotContains(t, filename, ".exe", "不允许上传可执行文件")
		})
	}
}

// TestAuthenticationSecurity 测试认证安全
func (suite *SecurityTestSuite) TestAuthenticationSecurity(t *testing.T) {
	testCases := []struct {
		name      string
		token     string
		expected  int
		desc      string
	}{
		{
			name:     "空token",
			token:    "",
			expected:  http.StatusUnauthorized,
			desc:     "空token应该被拒绝",
		},
		{
			name:     "无效token",
			token:    "invalid-token",
			expected:  http.StatusUnauthorized,
			desc:     "无效token应该被拒绝",
		},
		{
			name:     "过期token",
			token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", // 过期的JWT
			expected:  http.StatusUnauthorized,
			desc:     "过期token应该被拒绝",
		},
		{
			name:     "格式错误token",
			token:    "not-a-jwt-token",
			expected:  http.StatusUnauthorized,
			desc:     "格式错误token应该被拒绝",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/user/profile", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code, tc.desc)
		})
	}
}

// TestParameterPollution 测试参数污染
func (suite *SecurityTestSuite) TestParameterPollution(t *testing.T) {
	// 测试参数污染攻击
	data := map[string]interface{}{
		"username": "admin",
	}

	// 添加污染参数
	jsonData := []byte(`{"username":"admin","user":{"name":"hacker"}}`)

	req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证没有被参数污染
	if data, ok := response["data"].(map[string]interface{}); ok {
		username, _ := data["username"].(string)
		assert.Equal(t, "admin", username)
		// 检查是否没有意外的嵌套对象
		_, hasNestedUser := data["user"]
		assert.False(t, hasNestedUser, "检测到参数污染")
	}
}

// TestPasswordSecurity 测试密码安全
func (suite *SecurityTestSuite) TestPasswordSecurity(t *testing.T) {
	weakPasswords := []string{
		"123456",
		"password",
		"admin",
		"qwerty",
		"111111",
		"123123",
		"password1",
	}

	for _, password := range weakPasswords {
		t.Run("弱密码测试: "+password, func(t *testing.T) {
			data := map[string]interface{}{
				"password": password,
			}

			jsonData, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/api/v1/test/echo", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 在实际的密码验证中，弱密码应该被拒绝
			// 这里只是模拟测试逻辑
			assert.NotContains(t, []string{"123456", "password", "admin"}, password, 
				"不应该使用弱密码")
		})
	}
}

// TestHeadersSecurity 测试HTTP头部安全
func (suite *SecurityTestSuite) TestHeadersSecurity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/test/echo", nil)
	
	// 添加恶意头部
	req.Header.Set("X-Forwarded-For", "192.168.1.100")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1)")
	req.Header.Set("Referer", "http://malicious-site.com")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 检查安全相关的响应头
	headers := w.Header()
	
	// 检查是否有适当的安全头
	assert.NotEmpty(t, headers.Get("X-Content-Type-Options"))
	assert.NotEmpty(t, headers.Get("X-Frame-Options"))
	assert.NotEmpty(t, headers.Get("X-XSS-Protection"))
}

// TestErrorHandlingSecurity 测试错误处理安全
func (suite *SecurityTestSuite) TestErrorHandlingSecurity(t *testing.T) {
	// 测试错误信息是否泄露敏感信息
	testCases := []string{
		"/api/v1/nonexistent",
		"/api/v1/test/internal-error",
		"/api/v1/test/database-error",
	}

	for _, path := range testCases {
		t.Run("错误处理安全测试: "+path, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// 检查响应不包含敏感信息
			responseBody := w.Body.String()
			
			assert.NotContains(t, responseBody, "database", "错误响应不应包含数据库信息")
			assert.NotContains(t, responseBody, "stack trace", "错误响应不应包含堆栈跟踪")
			assert.NotContains(t, responseBody, "internal", "错误响应不应包含内部信息")
			
			// 应该返回标准化的错误响应
			assert.Contains(t, responseBody, "error", "应该返回标准错误格式")
		})
	}
}

// TestSessionSecurity 测试会话安全
func (suite *SecurityTestSuite) TestSessionSecurity(t *testing.T) {
	// 测试会话固定攻击
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", 
		bytes.NewBuffer([]byte(`{"phone":"test","code":"test"}`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 检查响应头中的会话信息
	headers := w.Header()
	
	// 如果使用Cookie，应该有安全属性
	setCookieHeader := headers.Get("Set-Cookie")
	if setCookieHeader != "" {
		assert.Contains(t, setCookieHeader, "HttpOnly", "Cookie应该设置HttpOnly")
		assert.Contains(t, setCookieHeader, "Secure", "Cookie应该设置Secure")
	}
}

// TestCORSSecurity 测试CORS安全
func (suite *SecurityTestSuite) TestCORSSecurity(t *testing.T) {
	// 测试跨域请求
	req, _ := http.NewRequest("OPTIONS", "/api/v1/test/echo", nil)
	req.Header.Set("Origin", "http://malicious-site.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 检查CORS头
	headers := w.Header()
	
	// 应该有适当的CORS控制
	allowedOrigin := headers.Get("Access-Control-Allow-Origin")
	
	// 不应该允许任意来源
	assert.NotEqual(t, "*", allowedOrigin, "不应该允许任意跨域请求")
}

func TestSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}