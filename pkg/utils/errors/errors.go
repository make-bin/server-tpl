package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// Error code ranges
const (
	// Success status codes (200-299)
	SuccessCodeMin = 200
	SuccessCodeMax = 299

	// System error codes (10000-19999)
	SystemErrorCodeMin = 10000
	SystemErrorCodeMax = 19999

	// Client error codes (20000-29999)
	ClientErrorCodeMin = 20000
	ClientErrorCodeMax = 29999

	// Business error codes (30000-99999)
	BusinessErrorCodeMin = 30000
	BusinessErrorCodeMax = 99999

	// Third party service error codes (100000-199999)
	ThirdPartyErrorCodeMin = 100000
	ThirdPartyErrorCodeMax = 199999
)

// Error codes
const (
	// System error codes (10000-19999)
	CodeSystemError         = 10000
	CodeDatabaseError       = 10001
	CodeCacheError          = 10002
	CodeNetworkError        = 10003
	CodeServiceUnavailable  = 10004
	CodeInternalServerError = 10005
	CodeConfigError         = 10006
	CodeFileSystemError     = 10007
	CodeMemoryError         = 10008
	CodeTimeoutError        = 10009

	// Client error codes (20000-29999)
	CodeValidationError      = 20000
	CodeMissingParameter     = 20001
	CodeInvalidParameter     = 20002
	CodeParameterTypeError   = 20003
	CodeUnauthorized         = 20004
	CodeForbidden            = 20005
	CodeNotFound             = 20006
	CodeMethodNotAllowed     = 20007
	CodeConflict             = 20008
	CodeTooManyRequests      = 20009
	CodeRequestTimeout       = 20010
	CodePayloadTooLarge      = 20011
	CodeUnsupportedMediaType = 20012
)

// Error represents a business error
type Error struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	StackTrace string                 `json:"stack_trace,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorWrapper wraps errors with additional context
type ErrorWrapper struct {
	err        error
	code       int
	message    string
	details    string
	fields     map[string]interface{}
	stackTrace string
}

func (w *ErrorWrapper) Error() string {
	return w.message
}

func (w *ErrorWrapper) Unwrap() error {
	return w.err
}

func (w *ErrorWrapper) Code() int {
	return w.code
}

func (w *ErrorWrapper) Details() string {
	return w.details
}

func (w *ErrorWrapper) Fields() map[string]interface{} {
	return w.fields
}

func (w *ErrorWrapper) StackTrace() string {
	return w.stackTrace
}

// Error code mapping
var errorCodeMap = map[int]string{
	// System errors
	CodeSystemError:         "系统错误",
	CodeDatabaseError:       "数据库错误",
	CodeCacheError:          "缓存错误",
	CodeNetworkError:        "网络错误",
	CodeServiceUnavailable:  "服务不可用",
	CodeInternalServerError: "服务器内部错误",
	CodeConfigError:         "配置错误",
	CodeFileSystemError:     "文件系统错误",
	CodeMemoryError:         "内存错误",
	CodeTimeoutError:        "超时错误",

	// Client errors
	CodeValidationError:      "参数验证失败",
	CodeMissingParameter:     "缺少必需参数",
	CodeInvalidParameter:     "参数无效",
	CodeParameterTypeError:   "参数类型错误",
	CodeUnauthorized:         "未授权访问",
	CodeForbidden:            "禁止访问",
	CodeNotFound:             "资源不存在",
	CodeMethodNotAllowed:     "方法不允许",
	CodeConflict:             "资源冲突",
	CodeTooManyRequests:      "请求过于频繁",
	CodeRequestTimeout:       "请求超时",
	CodePayloadTooLarge:      "请求体过大",
	CodeUnsupportedMediaType: "不支持的媒体类型",
}

// GetErrorMessage returns the error message for a given code
func GetErrorMessage(code int) string {
	if message, exists := errorCodeMap[code]; exists {
		return message
	}
	return "未知错误"
}

// NewError creates a new business error
func NewError(code int, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// NewErrorWithDetails creates a new error with details
func NewErrorWithDetails(code int, message, details string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// NewErrorWithFields creates a new error with fields
func NewErrorWithFields(code int, message string, fields map[string]interface{}) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, code int, message string) *ErrorWrapper {
	return &ErrorWrapper{
		err:        err,
		code:       code,
		message:    message,
		stackTrace: getStackTrace(),
	}
}

// WrapErrorWithDetails wraps an error with details
func WrapErrorWithDetails(err error, code int, message, details string) *ErrorWrapper {
	return &ErrorWrapper{
		err:        err,
		code:       code,
		message:    message,
		details:    details,
		stackTrace: getStackTrace(),
	}
}

// WrapErrorWithFields wraps an error with fields
func WrapErrorWithFields(err error, code int, message string, fields map[string]interface{}) *ErrorWrapper {
	return &ErrorWrapper{
		err:        err,
		code:       code,
		message:    message,
		fields:     fields,
		stackTrace: getStackTrace(),
	}
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return string(buf)
}

// IsErrorType checks if error is of specific type
func IsErrorType(err error, errorType string) bool {
	if appErr, ok := err.(*Error); ok {
		return strings.Contains(appErr.Message, errorType)
	}
	if wrapErr, ok := err.(*ErrorWrapper); ok {
		return strings.Contains(wrapErr.message, errorType)
	}
	return false
}

// GetHTTPStatusCode maps error code to HTTP status code
func GetHTTPStatusCode(code int) int {
	switch {
	case code >= SuccessCodeMin && code <= SuccessCodeMax:
		return code
	case code >= SystemErrorCodeMin && code <= SystemErrorCodeMax:
		return http.StatusInternalServerError
	case code >= ClientErrorCodeMin && code <= ClientErrorCodeMax:
		switch code {
		case CodeUnauthorized:
			return http.StatusUnauthorized
		case CodeForbidden:
			return http.StatusForbidden
		case CodeNotFound:
			return http.StatusNotFound
		case CodeMethodNotAllowed:
			return http.StatusMethodNotAllowed
		case CodeConflict:
			return http.StatusConflict
		case CodeTooManyRequests:
			return http.StatusTooManyRequests
		case CodeRequestTimeout:
			return http.StatusRequestTimeout
		case CodePayloadTooLarge:
			return http.StatusRequestEntityTooLarge
		case CodeUnsupportedMediaType:
			return http.StatusUnsupportedMediaType
		default:
			return http.StatusBadRequest
		}
	case code >= BusinessErrorCodeMin && code <= BusinessErrorCodeMax:
		return http.StatusBadRequest
	case code >= ThirdPartyErrorCodeMin && code <= ThirdPartyErrorCodeMax:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// Predefined errors for backward compatibility
var (
	ErrInternalServerError = NewError(CodeInternalServerError, GetErrorMessage(CodeInternalServerError))
	ErrBadRequestError     = NewError(CodeValidationError, GetErrorMessage(CodeValidationError))
	ErrUnauthorizedError   = NewError(CodeUnauthorized, GetErrorMessage(CodeUnauthorized))
	ErrForbiddenError      = NewError(CodeForbidden, GetErrorMessage(CodeForbidden))
	ErrNotFoundError       = NewError(CodeNotFound, GetErrorMessage(CodeNotFound))
	ErrConflictError       = NewError(CodeConflict, GetErrorMessage(CodeConflict))
	ErrValidationError     = NewError(CodeValidationError, GetErrorMessage(CodeValidationError))
)

// Legacy constants for backward compatibility
const (
	ErrInternalServer = "INTERNAL_SERVER_ERROR"
	ErrBadRequest     = "BAD_REQUEST"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrNotFound       = "NOT_FOUND"
	ErrConflict       = "CONFLICT"
	ErrValidation     = "VALIDATION_ERROR"
)

// AppError for backward compatibility
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status"`
	Cause      error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error (backward compatibility)
func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewAppErrorWithCause creates a new application error with cause (backward compatibility)
func NewAppErrorWithCause(code, message string, httpStatus int, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Cause:      cause,
	}
}

// IsAppError checks if an error is an AppError (backward compatibility)
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error, or creates a generic one (backward compatibility)
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppErrorWithCause(ErrInternalServer, "Internal server error", http.StatusInternalServerError, err)
}
