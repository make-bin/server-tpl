package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/make-bin/server-tpl/pkg/utils/i18n"
)

// Response 标准响应结构
type Response struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	requestID := getRequestID(c)
	response := Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   getMessage(c, "success"),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	c.JSON(http.StatusOK, response)
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, code int, message string, err error) {
	requestID := getRequestID(c)
	response := Response{
		Success:   false,
		Code:      code,
		Message:   getMessage(c, message),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// Page 分页响应
func Page(c *gin.Context, items interface{}, page, size, total int) {
	pages := (total + size - 1) / size
	pagination := Pagination{
		Page:  page,
		Size:  size,
		Total: total,
		Pages: pages,
	}

	data := PaginationResponse{
		Items:      items,
		Pagination: pagination,
	}

	Success(c, data)
}

// ValidationError 参数验证错误
func ValidationError(c *gin.Context, details []ErrorDetail) {
	requestID := getRequestID(c)
	response := Response{
		Success:   false,
		Code:      CodeValidationError,
		Message:   getMessage(c, "validation_error"),
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	c.JSON(http.StatusBadRequest, response)
}

// BusinessError 业务错误响应
func BusinessError(c *gin.Context, code int, messageKey string, err error) {
	statusCode := getHTTPStatusFromCode(code)
	Error(c, statusCode, code, messageKey, err)
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return "unknown"
}

// getMessage 获取本地化消息
func getMessage(c *gin.Context, key string) string {
	// 使用国际化工具获取消息
	if translator, exists := c.Get("translator"); exists {
		if t, ok := translator.(i18n.Translator); ok {
			return t.Translate(key)
		}
	}

	// 如果没有翻译器，返回预定义消息
	if message, exists := getDefaultMessage(key); exists {
		return message
	}

	return key
}

// getDefaultMessage 获取默认消息
func getDefaultMessage(key string) (string, bool) {
	messages := map[string]string{
		"success":          "操作成功",
		"validation_error": "参数验证失败",
		"user_not_found":   "用户不存在",
		"user_created":     "用户创建成功",
		"user_updated":     "用户更新成功",
		"user_deleted":     "用户删除成功",
		"app_not_found":    "应用不存在",
		"app_created":      "应用创建成功",
		"app_updated":      "应用更新成功",
		"app_deleted":      "应用删除成功",
		"internal_error":   "服务器内部错误",
		"unauthorized":     "未授权访问",
		"forbidden":        "权限不足",
		"not_found":        "资源不存在",
	}

	message, exists := messages[key]
	return message, exists
}

// getHTTPStatusFromCode 根据业务错误码获取HTTP状态码
func getHTTPStatusFromCode(code int) int {
	switch {
	case code >= 20000 && code < 21000: // 客户端错误
		return http.StatusBadRequest
	case code >= 21000 && code < 22000: // 认证错误
		return http.StatusUnauthorized
	case code >= 22000 && code < 23000: // 权限错误
		return http.StatusForbidden
	case code >= 23000 && code < 24000: // 资源不存在
		return http.StatusNotFound
	case code >= 24000 && code < 25000: // 冲突错误
		return http.StatusConflict
	case code >= 25000 && code < 26000: // 限流错误
		return http.StatusTooManyRequests
	case code >= 30000 && code < 40000: // 业务错误
		return http.StatusBadRequest
	case code >= 10000 && code < 20000: // 系统错误
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ParseValidationErrors 解析验证错误
func ParseValidationErrors(err error) []ErrorDetail {
	var details []ErrorDetail

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			detail := ErrorDetail{
				Field:  validationError.Field(),
				Reason: getValidationErrorMessage(validationError),
			}
			details = append(details, detail)
		}
	}

	return details
}

// getValidationErrorMessage 获取验证错误消息
func getValidationErrorMessage(ve validator.FieldError) string {
	switch ve.Tag() {
	case "required":
		return "此字段是必需的"
	case "email":
		return "请提供有效的邮箱地址"
	case "min":
		return "值太小，最小值为 " + ve.Param()
	case "max":
		return "值太大，最大值为 " + ve.Param()
	case "len":
		return "长度必须为 " + ve.Param()
	case "alphanum":
		return "只能包含字母和数字"
	case "phone":
		return "请提供有效的手机号"
	case "username":
		return "用户名格式不正确"
	case "password":
		return "密码强度不足"
	case "gte":
		return "值必须大于或等于 " + ve.Param()
	case "lte":
		return "值必须小于或等于 " + ve.Param()
	case "oneof":
		return "值必须是以下之一：" + ve.Param()
	default:
		return "字段验证失败"
	}
}

// WithMessage 自定义消息响应
func WithMessage(c *gin.Context, data interface{}, messageKey string) {
	requestID := getRequestID(c)
	response := Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   getMessage(c, messageKey),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	c.JSON(http.StatusOK, response)
}

// NoContent 无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}, messageKey string) {
	requestID := getRequestID(c)
	response := Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   getMessage(c, messageKey),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	c.JSON(http.StatusCreated, response)
}

// Accepted 已接受响应
func Accepted(c *gin.Context, data interface{}, messageKey string) {
	requestID := getRequestID(c)
	response := Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   getMessage(c, messageKey),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}

	c.JSON(http.StatusAccepted, response)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, messageKey, err)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusForbidden, CodeForbidden, messageKey, err)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusNotFound, CodeNotFound, messageKey, err)
}

// Conflict 冲突响应
func Conflict(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusConflict, CodeConflict, messageKey, err)
}

// TooManyRequests 请求过于频繁响应
func TooManyRequests(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusTooManyRequests, CodeTooManyRequests, messageKey, err)
}

// InternalServerError 服务器内部错误响应
func InternalServerError(c *gin.Context, messageKey string, err error) {
	Error(c, http.StatusInternalServerError, CodeInternalServerError, messageKey, err)
}
