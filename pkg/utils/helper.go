package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	
	// 生成6位随机数
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	
	return fmt.Sprintf("SQ%s%s", timestamp, randomStr[:6])
}

// GenerateNonce 生成随机字符串
func GenerateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, rand.NewSource(time.Now().Unix()).Int63())
		b[i] = charset[n%int64(len(charset))]
	}
	return string(b)
}

// GetClientIP 获取客户端真实IP
func GetClientIP(c *gin.Context) string {
	// 尝试从X-Forwarded-For获取
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// 尝试从X-Real-IP获取
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 尝试从X-Forwarded获取
	xForwarded := c.GetHeader("X-Forwarded")
	if xForwarded != "" {
		return xForwarded
	}

	// 使用RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	local, domain := parts[0], parts[1]
	
	// 检查local部分
	if len(local) == 0 || len(local) > 64 {
		return false
	}
	
	// 检查domain部分
	if len(domain) == 0 || len(domain) > 255 {
		return false
	}
	
	// 更严格的邮箱验证可以用正则表达式，这里简化处理
	return strings.Contains(domain, ".")
}

// IsValidPhone 验证手机号格式（中国大陆）
func IsValidPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}
	
	// 移除非数字字符
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
	
	// 中国大陆手机号：11位，1开头
	return len(digits) == 11 && digits[0] == '1'
}

// FormatAmount 格式化金额显示
func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// ParseAmount 解析金额字符串
func ParseAmount(amountStr string) (float64, error) {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的金额格式: %w", err)
	}
	
	if amount < 0 {
		return 0, fmt.Errorf("金额不能为负数")
	}
	
	return amount, nil
}

// CalculateServiceFee 计算服务费
func CalculateServiceFee(amount float64, rate float64) float64 {
	return amount * rate
}

// CalculateDeposit 计算保证金
func CalculateDeposit(amount float64, rate float64) float64 {
	return amount * rate
}

// RoundToMoney 四舍五入到分
func RoundToMoney(amount float64) float64 {
	return float64(int64(amount*100+0.5)) / 100
}

// IsExpired 检查是否过期
func IsExpired(expiredAt time.Time) bool {
	return time.Now().After(expiredAt)
}

// GetExpireDuration 获取过期时间间隔
func GetExpireDuration(expireTime int) time.Duration {
	return time.Duration(expireTime) * time.Second
}

// SanitizeString 清理字符串（移除前后空格，特殊字符等）
func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	
	// 将多个连续空格替换为单个空格
	for strings.Contains(input, "  ") {
		input = strings.ReplaceAll(input, "  ", " ")
	}
	
	return input
}

// ValidateRequired 验证必填字段
func ValidateRequired(fields map[string]string) []string {
	var errors []string
	
	for fieldName, fieldValue := range fields {
		if strings.TrimSpace(fieldValue) == "" {
			errors = append(errors, fmt.Sprintf("%s 不能为空", fieldName))
		}
	}
	
	return errors
}

// BuildURL 构建URL
func BuildURL(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}
	
	var builder strings.Builder
	builder.WriteString(baseURL)
	builder.WriteString("?")
	
	first := true
	for key, value := range params {
		if !first {
			builder.WriteString("&")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
		first = false
	}
	
	return builder.String()
}

// GetRequestHeader 获取请求头
func GetRequestHeader(c *gin.Context, key string) string {
	return c.GetHeader(key)
}

// IsHTTPS 判断是否为HTTPS请求
func IsHTTPS(c *gin.Context) bool {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = c.Request.URL.Scheme
	}
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme == "https"
}

// GetBaseURL 获取基础URL
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if IsHTTPS(c) {
		scheme = "https"
	}
	
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	
	return fmt.Sprintf("%s://%s", scheme, host)
}

// Pagination 计算分页参数
func Pagination(page, pageSize int) (offset int, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	offset = (page - 1) * pageSize
	limit = pageSize
	
	return offset, limit
}

// Contains 检查字符串数组是否包含某个字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// InStringArray 检查字符串是否在数组中（忽略大小写）
func InStringArray(slice []string, item string) bool {
	item = strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == item {
			return true
		}
	}
	return false
}

// FilterEmptyStrings 过滤空字符串
func FilterEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

// UniqueStrings 去重字符串数组
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	
	return result
}

// FileSizeToString 文件大小转字符串
func FileSizeToString(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// GetContentTypeByFileExtension 根据文件扩展名获取Content-Type
func GetContentTypeByFileExtension(filename string) string {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	
	mimeTypes := map[string]string{
		"txt":  "text/plain",
		"html": "text/html",
		"css":  "text/css",
		"js":   "application/javascript",
		"json": "application/json",
		"xml":  "application/xml",
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"ppt":  "application/vnd.ms-powerpoint",
		"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"svg":  "image/svg+xml",
		"mp3":  "audio/mpeg",
		"wav":  "audio/wav",
		"mp4":  "video/mp4",
		"avi":  "video/x-msvideo",
		"zip":  "application/zip",
		"rar":  "application/x-rar-compressed",
	}
	
	if contentType, ok := mimeTypes[ext]; ok {
		return contentType
	}
	
	return "application/octet-stream"
}