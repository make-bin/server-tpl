package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Error 自定义错误类型
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Stack   string `json:"stack,omitempty"`
	cause   error
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	return e.cause
}

// WithDetails 添加详细信息
func (e *Error) WithDetails(details string) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
		Stack:   e.Stack,
		cause:   e.cause,
	}
}

// WithCause 添加原始错误
func (e *Error) WithCause(cause error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
		Stack:   e.Stack,
		cause:   cause,
	}
}

// New 创建新的错误
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Stack:   getStackTrace(),
	}
}

// Wrap 包装现有错误
func Wrap(err error, code int, message string) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Code:    code,
		Message: message,
		Stack:   getStackTrace(),
		cause:   err,
	}
}

// Wrapf 格式化包装错误
func Wrapf(err error, code int, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Stack:   getStackTrace(),
		cause:   err,
	}
}

// Is 检查错误类型
func Is(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}

	if e, ok := err.(*Error); ok {
		if t, ok := target.(*Error); ok {
			return e.Code == t.Code
		}
	}

	return false
}

// As 类型断言
func As(err error, target interface{}) bool {
	return false
}

// getStackTrace 获取堆栈跟踪
func getStackTrace() string {
	var stack []string
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		// 跳过内部函数
		if strings.Contains(fn.Name(), "pkg/utils/errors") {
			continue
		}

		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}

	return strings.Join(stack, "\n")
}

// Common error constructors
var (
	// 通用错误
	ErrInternalServer = New(500, "Internal server error")
	ErrBadRequest     = New(400, "Bad request")
	ErrUnauthorized   = New(401, "Unauthorized")
	ErrForbidden      = New(403, "Forbidden")
	ErrNotFound       = New(404, "Not found")
	ErrConflict       = New(409, "Conflict")
	ErrValidation     = New(422, "Validation failed")

	// 数据库错误
	ErrDatabaseConnection  = New(500, "Database connection failed")
	ErrDatabaseQuery       = New(500, "Database query failed")
	ErrDatabaseTransaction = New(500, "Database transaction failed")

	// 业务错误
	ErrInvalidInput     = New(400, "Invalid input")
	ErrResourceNotFound = New(404, "Resource not found")
	ErrResourceExists   = New(409, "Resource already exists")
	ErrPermissionDenied = New(403, "Permission denied")
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details string      `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ToResponse 转换为错误响应
func (e *Error) ToResponse() *ErrorResponse {
	return &ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}

// FromError 从error创建ErrorResponse
func FromError(err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e.ToResponse()
	}

	// 默认内部服务器错误
	return &ErrorResponse{
		Code:    500,
		Message: "Internal server error",
		Details: err.Error(),
	}
}
