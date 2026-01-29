package utils

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
    Code      int         `json:"code"`      // 状态码
    Message   string      `json:"message"`   // 消息
    Data      interface{} `json:"data"`      // 数据
    Timestamp int64       `json:"timestamp"` // 时间戳
}

// PageResponse 分页响应结构
type PageResponse struct {
    List       interface{} `json:"list"`       // 数据列表
    Pagination Pagination  `json:"pagination"` // 分页信息
}

// Pagination 分页信息
type Pagination struct {
    Page       int `json:"page"`        // 当前页码
    PageSize   int `json:"page_size"`   // 每页数量
    Total      int64 `json:"total"`      // 总记录数
    TotalPages int  `json:"total_pages"` // 总页数
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
    response := Response{
        Code:      http.StatusOK,
        Message:   "success",
        Data:      data,
        Timestamp: getCurrentTimestamp(),
    }
    c.JSON(http.StatusOK, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
    response := Response{
        Code:      code,
        Message:   message,
        Data:      nil,
        Timestamp: getCurrentTimestamp(),
    }
    c.JSON(code, response)
}

// SuccessPageResponse 分页成功响应
func SuccessPageResponse(c *gin.Context, list interface{}, pagination Pagination) {
    data := PageResponse{
        List:       list,
        Pagination: pagination,
    }
    
    response := Response{
        Code:      http.StatusOK,
        Message:   "success",
        Data:      data,
        Timestamp: getCurrentTimestamp(),
    }
    c.JSON(http.StatusOK, response)
}

// CreatedResponse 创建成功响应
func CreatedResponse(c *gin.Context, data interface{}) {
    response := Response{
        Code:      http.StatusCreated,
        Message:   "created",
        Data:      data,
        Timestamp: getCurrentTimestamp(),
    }
    c.JSON(http.StatusCreated, response)
}

// BadRequestResponse 错误请求响应
func BadRequestResponse(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusBadRequest, message)
}

// UnauthorizedResponse 未授权响应
func UnauthorizedResponse(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse 禁止访问响应
func ForbiddenResponse(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusForbidden, message)
}

// NotFoundResponse 未找到响应
func NotFoundResponse(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse 内部服务器错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusInternalServerError, message)
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
    return 0 // TODO: 实现时间戳生成
}