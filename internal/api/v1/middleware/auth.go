package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.uber.org/zap"

    "task-platform-api/internal/config"
    "task-platform-api/internal/models"
    "task-platform-api/pkg/utils"
)

type Claims struct {
    UserID   uint64 `json:"user_id"`
    OpenID   string `json:"openid"`
    AuthType string `json:"auth_type"`
    jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func JWTAuth(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getTokenFromRequest(c)
        if token == "" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "未提供认证令牌")
            c.Abort()
            return
        }

        claims, err := parseToken(token, cfg.Secret)
        if err != nil {
            logger.Warn("JWT令牌解析失败", zap.Error(err), zap.String("token", token))
            utils.ErrorResponse(c, http.StatusUnauthorized, "无效的认证令牌")
            c.Abort()
            return
        }

        // 检查令牌是否过期
        if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
            utils.ErrorResponse(c, http.StatusUnauthorized, "认证令牌已过期")
            c.Abort()
            return
        }

        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("openid", claims.OpenID)
        c.Set("auth_type", claims.AuthType)

        c.Next()
    }
}

// OptionalAuth 可选认证中间件
func OptionalAuth(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getTokenFromRequest(c)
        if token != "" {
            claims, err := parseToken(token, cfg.Secret)
            if err == nil && (claims.ExpiresAt == nil || claims.ExpiresAt.Time.After(time.Now())) {
                c.Set("user_id", claims.UserID)
                c.Set("openid", claims.OpenID)
                c.Set("auth_type", claims.AuthType)
            }
        }
        c.Next()
    }
}

// RequireRole 角色权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            utils.ErrorResponse(c, http.StatusUnauthorized, "未认证用户")
            c.Abort()
            return
        }

        // 这里应该从数据库获取用户信息，简化处理
        // 在实际应用中，建议使用缓存来优化
        userRole := "user" // 默认角色
        
        hasPermission := false
        for _, role := range roles {
            if userRole == role {
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
            c.Abort()
            return
        }

        c.Set("role", userRole)
        c.Next()
    }
}

// RequireNormalUser 正常用户状态中间件
func RequireNormalUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            utils.ErrorResponse(c, http.StatusUnauthorized, "未认证用户")
            c.Abort()
            return
        }

        var user models.User
        if err := db.First(&user, userID).Error; err != nil {
            utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
            c.Abort()
            return
        }

        if !user.IsNormal() {
            utils.ErrorResponse(c, http.StatusForbidden, "用户状态异常")
            c.Abort()
            return
        }

        c.Set("user", &user)
        c.Next()
    }
}

// RateLimit 限流中间件
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
    // 这里应该使用Redis等实现分布式限流
    // 简化实现，实际应用中建议使用成熟的限流库
    return func(c *gin.Context) {
        // TODO: 实现限流逻辑
        c.Next()
    }
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *models.User, cfg *config.JWTConfig) (string, error) {
    now := time.Now()
    claims := Claims{
        UserID:   user.ID,
        OpenID:   user.OpenID,
        AuthType: user.AuthType,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.ExpireTime) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    "task-platform",
            Subject:   string(user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(user *models.User, cfg *config.JWTConfig) (string, error) {
    now := time.Now()
    claims := Claims{
        UserID:   user.ID,
        OpenID:   user.OpenID,
        AuthType: user.AuthType,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.RefreshExpireTime) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    "task-platform",
            Subject:   string(user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}

// getTokenFromRequest 从请求中获取令牌
func getTokenFromRequest(c *gin.Context) string {
    // 从Authorization header获取
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" {
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) == 2 && parts[0] == "Bearer" {
            return parts[1]
        }
    }

    // 从查询参数获取
    return c.Query("token")
}

// parseToken 解析令牌
func parseToken(tokenString, secret string) (*Claims, error) {
    return ParseToken(tokenString, secret)
}