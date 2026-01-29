package middleware

import (
    "bytes"
    "io"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// Logging 日志中间件
func Logging(logger *zap.Logger) gin.HandlerFunc {
    return gin.LoggerWithConfig(gin.LoggerConfig{
        Formatter: func(param gin.LogFormatterParams) string {
            // 自定义日志格式
            logger.Info("HTTP请求",
                zap.String("method", param.Method),
                zap.String("path", param.Path),
                zap.Int("status", param.StatusCode),
                zap.Duration("latency", param.Latency),
                zap.String("client_ip", param.ClientIP),
                zap.String("user_agent", param.Request.UserAgent()),
                zap.String("error", param.ErrorMessage),
                zap.Int("body_size", param.BodySize),
            )
            return ""
        },
        Output: io.Discard, // 我们使用zap输出，所以这里是空的
    })
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}

// RequestBodyLog 请求体日志中间件
func RequestBodyLog(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 只对特定方法的请求记录body
        if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
            // 读取body
            bodyBytes, _ := io.ReadAll(c.Request.Body)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
            
            // 记录请求体（注意：敏感信息需要过滤）
            if len(bodyBytes) > 0 {
                bodyString := string(bodyBytes)
                // 这里可以添加敏感信息过滤逻辑
                logger.Debug("请求体",
                    zap.String("path", c.Request.URL.Path),
                    zap.String("method", c.Request.Method),
                    zap.String("body", bodyString),
                )
            }
        }
        
        c.Next()
    }
}

// ResponseLog 响应日志中间件
func ResponseLog(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        
        // 拦截响应
        writer := &responseBodyWriter{
            ResponseWriter: c.Writer,
            body:          &bytes.Buffer{},
        }
        c.Writer = writer
        
        c.Next()
        
        duration := time.Since(startTime)
        
        // 记录响应信息
        logger.Info("HTTP响应",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
            zap.Int("response_size", writer.body.Len()),
            zap.String("client_ip", c.ClientIP()),
            zap.Any("headers", c.Writer.Header()),
        )
        
        // 对于错误响应，记录响应体
        if c.Writer.Status() >= 400 {
            logger.Error("错误响应",
                zap.String("path", c.Request.URL.Path),
                zap.Int("status", c.Writer.Status()),
                zap.String("response_body", writer.body.String()),
            )
        }
    }
}

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
    r.body.Write(b)
    return r.ResponseWriter.Write(b)
}

// generateRequestID 生成请求ID
func generateRequestID() string {
    // 简单实现，实际可以使用UUID或其他方式
    return "req_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}