package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    go_redis "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "task-platform-api/internal/config"
    "task-platform-api/internal/api/v1/middleware"
    "task-platform-api/internal/models"
    "task-platform-api/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
    db     *gorm.DB
    rdb    *go_redis.Client
    cfg    *config.Config
    logger *zap.Logger
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB, rdb *go_redis.Client, cfg *config.Config, logger *zap.Logger) *AuthHandler {
    return &AuthHandler{
        db:     db,
        rdb:    rdb,
        cfg:    cfg,
        logger: logger,
    }
}

// LoginRequest 登录请求
type LoginRequest struct {
    AuthType string `json:"auth_type" binding:"required,oneof=wechat alipay"`
    Code     string `json:"code" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    User         *UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
    ID         uint64  `json:"id"`
    OpenID     string  `json:"openid"`
    AuthType   string  `json:"auth_type"`
    Nickname   string  `json:"nickname"`
    Avatar     string  `json:"avatar"`
    Phone      string  `json:"phone"`
    Email      string  `json:"email"`
    CreditScore float32 `json:"credit_score"`
    Level      int     `json:"level"`
    CreatedAt  int64   `json:"created_at"`
}

// WechatLogin 微信登录
func (h *AuthHandler) WechatLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestResponse(c, err.Error())
        return
    }

    // 获取设备信息
    deviceInfo := c.GetHeader("User-Agent")
    ipAddress := c.ClientIP()

    // 调用微信API获取用户信息
    wechatUser, err := h.getWechatUserInfo(req.Code)
    if err != nil {
        h.logger.Error("获取微信用户信息失败", zap.Error(err))
        utils.ErrorResponse(c, http.StatusUnauthorized, "微信授权失败")
        return
    }

    // 查找或创建用户
    user, err := h.findOrCreateUser(wechatUser.OpenID, req.AuthType, wechatUser)
    if err != nil {
        h.logger.Error("用户处理失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "用户处理失败")
        return
    }

    // 生成JWT令牌
    accessToken, err := middleware.GenerateToken(user, &h.cfg.JWT)
    if err != nil {
        h.logger.Error("生成访问令牌失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "令牌生成失败")
        return
    }

    refreshToken, err := middleware.GenerateRefreshToken(user, &h.cfg.JWT)
    if err != nil {
        h.logger.Error("生成刷新令牌失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "令牌生成失败")
        return
    }

    // 保存会话
    sessionID := uuid.New().String()
    err = h.saveUserSession(sessionID, user.ID, accessToken, deviceInfo, ipAddress)
    if err != nil {
        h.logger.Error("保存用户会话失败", zap.Error(err))
        // 不影响登录流程
    }

    // 返回登录结果
    response := LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(h.cfg.JWT.ExpireTime),
        User:         h.convertUserInfo(user),
    }

    utils.SuccessResponse(c, response)
}

// AlipayLogin 支付宝登录
func (h *AuthHandler) AlipayLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestResponse(c, err.Error())
        return
    }

    // 获取设备信息
    deviceInfo := c.GetHeader("User-Agent")
    ipAddress := c.ClientIP()

    // 调用支付宝API获取用户信息
    alipayUser, err := h.getAlipayUserInfo(req.Code)
    if err != nil {
        h.logger.Error("获取支付宝用户信息失败", zap.Error(err))
        utils.ErrorResponse(c, http.StatusUnauthorized, "支付宝授权失败")
        return
    }

    // 查找或创建用户
    user, err := h.findOrCreateUser(alipayUser.UserID, req.AuthType, alipayUser)
    if err != nil {
        h.logger.Error("用户处理失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "用户处理失败")
        return
    }

    // 生成JWT令牌
    accessToken, err := middleware.GenerateToken(user, &h.cfg.JWT)
    if err != nil {
        h.logger.Error("生成访问令牌失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "令牌生成失败")
        return
    }

    refreshToken, err := middleware.GenerateRefreshToken(user, &h.cfg.JWT)
    if err != nil {
        h.logger.Error("生成刷新令牌失败", zap.Error(err))
        utils.InternalServerErrorResponse(c, "令牌生成失败")
        return
    }

    // 保存会话
    sessionID := uuid.New().String()
    err = h.saveUserSession(sessionID, user.ID, accessToken, deviceInfo, ipAddress)
    if err != nil {
        h.logger.Error("保存用户会话失败", zap.Error(err))
    }

    // 返回登录结果
    response := LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(h.cfg.JWT.ExpireTime),
        User:         h.convertUserInfo(user),
    }

    utils.SuccessResponse(c, response)
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestResponse(c, err.Error())
        return
    }

    // 解析刷新令牌
    claims, err := middleware.ParseToken(req.RefreshToken, h.cfg.JWT.Secret)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "无效的刷新令牌")
        return
    }

    // 获取用户信息
    var user models.User
    if err := h.db.First(&user, claims.UserID).Error; err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
        return
    }

    // 生成新的访问令牌
    newAccessToken, err := middleware.GenerateToken(&user, &h.cfg.JWT)
    if err != nil {
        utils.InternalServerErrorResponse(c, "令牌生成失败")
        return
    }

    // 返回新的令牌
    response := map[string]interface{}{
        "access_token": newAccessToken,
        "expires_in":   h.cfg.JWT.ExpireTime,
    }

    utils.SuccessResponse(c, response)
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    // 删除用户会话
    h.deleteUserSessions(userID.(uint64))

    utils.SuccessResponse(c, gin.H{
        "message": "退出登录成功",
    })
}

// WechatUser 微信用户信息
type WechatUser struct {
    OpenID     string `json:"openid"`
    Nickname   string `json:"nickname"`
    Avatar     string `json:"headimgurl"`
    UnionID    string `json:"unionid"`
}

// AlipayUser 支付宝用户信息
type AlipayUser struct {
    UserID   string `json:"user_id"`
    Nickname string `json:"nick_name"`
    Avatar   string `json:"avatar"`
}

// getWechatUserInfo 获取微信用户信息
func (h *AuthHandler) getWechatUserInfo(code string) (*WechatUser, error) {
    // TODO: 调用微信API
    // 这里应该实现微信小程序或公众号的授权登录逻辑
    return &WechatUser{
        OpenID:   "test_openid_" + code,
        Nickname: "测试用户",
        Avatar:   "https://example.com/avatar.jpg",
        UnionID:  "test_unionid_" + code,
    }, nil
}

// getAlipayUserInfo 获取支付宝用户信息
func (h *AuthHandler) getAlipayUserInfo(code string) (*AlipayUser, error) {
    // TODO: 调用支付宝API
    // 这里应该实现支付宝的授权登录逻辑
    return &AlipayUser{
        UserID:   "test_userid_" + code,
        Nickname: "支付宝测试用户",
        Avatar:   "https://example.com/alipay_avatar.jpg",
    }, nil
}

// findOrCreateUser 查找或创建用户
func (h *AuthHandler) findOrCreateUser(openID, authType string, userInfo interface{}) (*models.User, error) {
    var user models.User
    
    // 查找用户
    err := h.db.Where("openid = ? AND auth_type = ?", openID, authType).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            // 创建新用户
            user = models.User{
                OpenID:   openID,
                AuthType: authType,
                Status:   1, // 正常状态
            }

            // 设置用户信息
            switch v := userInfo.(type) {
            case *WechatUser:
                user.Nickname = v.Nickname
                user.Avatar = v.Avatar
                user.UnionID = v.UnionID
            case *AlipayUser:
                user.Nickname = v.Nickname
                user.Avatar = v.Avatar
            }

            // 创建用户记录
            if err := h.db.Create(&user).Error; err != nil {
                return nil, err
            }

            // 创建用户信誉记录
            credit := models.UserCredit{
                UserID: user.ID,
                Score:  5.0,
                Level:  1,
            }
            if err := h.db.Create(&credit).Error; err != nil {
                h.logger.Error("创建用户信誉记录失败", zap.Error(err))
            }

            // 创建钱包记录
            wallet := models.Wallet{
                UserID:  user.ID,
                Balance: 0,
            }
            if err := h.db.Create(&wallet).Error; err != nil {
                h.logger.Error("创建用户钱包失败", zap.Error(err))
            }

        } else {
            return nil, err
        }
    }

    return &user, nil
}

// saveUserSession 保存用户会话
func (h *AuthHandler) saveUserSession(sessionID string, userID uint64, token, deviceInfo, ipAddress string) error {
    expireTime := time.Now().Add(time.Duration(h.cfg.JWT.ExpireTime) * time.Second)
    
    session := models.UserSession{
        SessionID:  sessionID,
        UserID:     userID,
        Token:      token,
        ExpireTime: expireTime,
        DeviceInfo: deviceInfo,
        IPAddress:  ipAddress,
        UserAgent:  deviceInfo,
    }

    return h.db.Create(&session).Error
}

// deleteUserSessions 删除用户所有会话
func (h *AuthHandler) deleteUserSessions(userID uint64) error {
    return h.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

// convertUserInfo 转换用户信息
func (h *AuthHandler) convertUserInfo(user *models.User) *UserInfo {
    return &UserInfo{
        ID:          user.ID,
        OpenID:      user.OpenID,
        AuthType:    user.AuthType,
        Nickname:    user.Nickname,
        Avatar:      user.Avatar,
        Phone:       user.Phone,
        Email:       user.Email,
        CreditScore: user.CreditScore,
        Level:       user.Level,
        CreatedAt:   user.CreatedAt.Unix(),
    }
}