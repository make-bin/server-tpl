package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/make-bin/server-tpl/pkg/api/response"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
	"golang.org/x/time/rate"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret        string   `json:"jwt_secret"`
	RateLimitRPS     int      `json:"rate_limit_rps"`
	RateLimitBurst   int      `json:"rate_limit_burst"`
	MaxFileSize      int64    `json:"max_file_size"`
	AllowedFileTypes []string `json:"allowed_file_types"`
	CSRFEnabled      bool     `json:"csrf_enabled"`
	EncryptionKey    string   `json:"encryption_key"`
}

// DefaultSecurityConfig 默认安全配置
var DefaultSecurityConfig = &SecurityConfig{
	JWTSecret:        "your-secret-key",
	RateLimitRPS:     100,
	RateLimitBurst:   200,
	MaxFileSize:      10 * 1024 * 1024, // 10MB
	AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
	CSRFEnabled:      true,
	EncryptionKey:    "your-encryption-key-32-characters",
}

// SecurityHeadersMiddleware 安全响应头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// XSS保护
		c.Header("X-XSS-Protection", "1; mode=block")

		// 严格传输安全
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过某些路径
		if isSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 从请求头获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "unauthorized", fmt.Errorf("未提供认证令牌"))
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		// 验证JWT token
		claims, err := validateJWTToken(token, config.JWTSecret)
		if err != nil {
			response.Unauthorized(c, "invalid_token", err)
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_permissions", claims.Permissions)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// RequireRole 角色授权中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Forbidden(c, "permission_denied", fmt.Errorf("权限不足"))
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "permission_denied", fmt.Errorf("权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 权限授权中间件
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			response.Forbidden(c, "permission_denied", fmt.Errorf("权限不足"))
			c.Abort()
			return
		}

		// 检查用户权限
		userPerms, ok := userPermissions.([]string)
		if !ok {
			response.Forbidden(c, "permission_denied", fmt.Errorf("权限不足"))
			c.Abort()
			return
		}

		// 检查是否有所需权限
		hasPermission := false
		for _, requiredPerm := range permissions {
			for _, userPerm := range userPerms {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "permission_denied", fmt.Errorf("权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config *SecurityConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(config.RateLimitRPS), config.RateLimitBurst)

	return func(c *gin.Context) {
		// 获取客户端标识
		clientID := getClientID(c)

		// 检查限流
		if !limiter.Allow() {
			logger.Warn("Rate limit exceeded for client: %s", clientID)
			response.TooManyRequests(c, "rate_limit_exceeded", fmt.Errorf("请求频率超限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CSRFMiddleware CSRF防护中间件
func CSRFMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.CSRFEnabled {
			c.Next()
			return
		}

		// 只对非GET请求进行CSRF检查
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 检查CSRF token
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			response.Forbidden(c, "csrf_token_missing", fmt.Errorf("缺少CSRF令牌"))
			c.Abort()
			return
		}

		// 验证CSRF token
		if !validateCSRFToken(c, csrfToken) {
			response.Forbidden(c, "csrf_token_invalid", fmt.Errorf("CSRF令牌无效"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// FileUploadSecurityMiddleware 文件上传安全中间件
func FileUploadSecurityMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			if err != http.ErrMissingFile {
				response.Error(c, http.StatusBadRequest, response.CodeFileUploadFailed, "file_upload_failed", err)
				c.Abort()
			}
			return
		}
		defer file.Close()

		// 1. 检查文件大小
		if header.Size > config.MaxFileSize {
			response.Error(c, http.StatusBadRequest, response.CodeFileTooBig, "file_too_big", fmt.Errorf("文件大小超过限制"))
			c.Abort()
			return
		}

		// 2. 检查文件类型
		contentType := header.Header.Get("Content-Type")
		if !isAllowedFileType(contentType, config.AllowedFileTypes) {
			response.Error(c, http.StatusBadRequest, response.CodeFileTypeNotSupported, "file_type_not_supported", fmt.Errorf("不支持的文件类型"))
			c.Abort()
			return
		}

		// 3. 检查文件扩展名
		if !isAllowedFileExtension(header.Filename) {
			response.Error(c, http.StatusBadRequest, response.CodeFileTypeNotSupported, "file_extension_not_allowed", fmt.Errorf("不支持的文件扩展名"))
			c.Abort()
			return
		}

		// 4. 生成安全的文件名
		safeFileName := generateSafeFileName(header.Filename)

		c.Set("safe_file_name", safeFileName)
		c.Set("original_file_name", header.Filename)
		c.Set("file_size", header.Size)
		c.Set("content_type", contentType)

		c.Next()
	}
}

// InputValidationMiddleware 输入验证中间件
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查SQL注入
		if containsSQLInjection(c.Request.URL.RawQuery) {
			logger.Warn("SQL injection attempt detected from %s", getClientID(c))
			response.Error(c, http.StatusBadRequest, response.CodeInvalidParameter, "invalid_parameter", fmt.Errorf("检测到非法参数"))
			c.Abort()
			return
		}

		// 检查XSS攻击
		if containsXSS(c.Request.URL.RawQuery) {
			logger.Warn("XSS attempt detected from %s", getClientID(c))
			response.Error(c, http.StatusBadRequest, response.CodeInvalidParameter, "invalid_parameter", fmt.Errorf("检测到非法参数"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// 辅助函数

// isSkipPath 检查是否跳过认证的路径
func isSkipPath(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/v1/applications/health",
		"/swagger",
		"/metrics",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// validateJWTToken 验证JWT token
func validateJWTToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// getClientID 获取客户端标识
func getClientID(c *gin.Context) string {
	// 优先使用X-Forwarded-For
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	// 使用X-Real-IP
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 使用RemoteAddr
	return c.ClientIP()
}

// validateCSRFToken 验证CSRF token
func validateCSRFToken(c *gin.Context, token string) bool {
	// 从session中获取存储的token
	sessionToken, exists := c.Get("csrf_token")
	if !exists {
		return false
	}

	return token == sessionToken
}

// isAllowedFileType 检查允许的文件类型
func isAllowedFileType(contentType string, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		if strings.TrimSpace(t) == contentType {
			return true
		}
	}
	return false
}

// isAllowedFileExtension 检查允许的文件扩展名
func isAllowedFileExtension(filename string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt"}

	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// generateSafeFileName 生成安全的文件名
func generateSafeFileName(originalName string) string {
	ext := originalName[strings.LastIndex(originalName, "."):]
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%s%s", timestamp, ext)
}

// containsSQLInjection 检查SQL注入
func containsSQLInjection(input string) bool {
	sqlInjectionPatterns := []string{
		`(?i)(union\s+select)`,
		`(?i)(drop\s+table)`,
		`(?i)(delete\s+from)`,
		`(?i)(insert\s+into)`,
		`(?i)(update\s+.+set)`,
		`(?i)(or\s+1=1)`,
		`(?i)(and\s+1=1)`,
		`(?i)('|"|;|--|#)`,
	}

	for _, pattern := range sqlInjectionPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}
	return false
}

// containsXSS 检查XSS攻击
func containsXSS(input string) bool {
	xssPatterns := []string{
		`(?i)<script.*?>`,
		`(?i)<\/script>`,
		`(?i)javascript:`,
		`(?i)vbscript:`,
		`(?i)onload\s*=`,
		`(?i)onerror\s*=`,
		`(?i)onclick\s*=`,
	}

	for _, pattern := range xssPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}
	return false
}

// EncryptSensitiveData 加密敏感数据
func EncryptSensitiveData(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key)[:32]) // 确保密钥长度为32
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSensitiveData 解密敏感数据
func DecryptSensitiveData(encryptedData, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key)[:32]) // 确保密钥长度为32
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// MaskSensitiveData 敏感数据脱敏
func MaskSensitiveData(dataType, data string) string {
	switch dataType {
	case "phone":
		return maskPhone(data)
	case "email":
		return maskEmail(data)
	case "idcard":
		return maskIDCard(data)
	default:
		return data
	}
}

// maskPhone 手机号脱敏
func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// maskEmail 邮箱脱敏
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	maskedUsername := username[:1] + "***" + username[len(username)-1:]
	return maskedUsername + "@" + domain
}

// maskIDCard 身份证号脱敏
func maskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:4] + "********" + idCard[len(idCard)-4:]
}
