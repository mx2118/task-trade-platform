package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "task-platform-api/internal/config"
)

// CORS 跨域中间件
func CORS(cfg *config.SecurityConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // 检查是否在允许的源列表中
        allowed := false
        for _, allowedOrigin := range cfg.AllowOrigins {
            if allowedOrigin == "*" || allowedOrigin == origin {
                allowed = true
                break
            }
        }

        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
        }

        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        // 处理预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

// CustomHeaders 自定义响应头中间件
func CustomHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 安全相关头
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
        
        // 缓存相关头
        if c.Request.Method == "GET" {
            c.Header("Cache-Control", "public, max-age=300")
        } else {
            c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
        }

        c.Next()
    }
}

// Compression 压缩中间件
func Compression() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Header("Vary", "Accept-Encoding")
        c.Next()
    })
}