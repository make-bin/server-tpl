package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
		want    *Error
	}{
		{
			name:    "创建基本错误",
			code:    400,
			message: "Bad request",
			want: &Error{
				Code:    400,
				Message: "Bad request",
			},
		},
		{
			name:    "创建内部服务器错误",
			code:    500,
			message: "Internal server error",
			want: &Error{
				Code:    500,
				Message: "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)

			assert.Equal(t, tt.want.Code, err.Code)
			assert.Equal(t, tt.want.Message, err.Message)
			assert.NotEmpty(t, err.Stack)
			assert.Nil(t, err.cause)
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "基本错误消息",
			err: &Error{
				Code:    400,
				Message: "Bad request",
			},
			expected: "Bad request",
		},
		{
			name: "带原始错误的错误消息",
			err: &Error{
				Code:    400,
				Message: "Bad request",
				cause:   fmt.Errorf("original error"),
			},
			expected: "Bad request: original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	err := &Error{
		Code:    400,
		Message: "Bad request",
		cause:   originalErr,
	}

	result := err.Unwrap()
	assert.Equal(t, originalErr, result)
}

func TestError_WithDetails(t *testing.T) {
	originalErr := &Error{
		Code:    400,
		Message: "Bad request",
		Stack:   "stack trace",
	}

	details := "Additional details about the error"
	result := originalErr.WithDetails(details)

	assert.Equal(t, originalErr.Code, result.Code)
	assert.Equal(t, originalErr.Message, result.Message)
	assert.Equal(t, details, result.Details)
	assert.Equal(t, originalErr.Stack, result.Stack)
	assert.Equal(t, originalErr.cause, result.cause)
}

func TestError_WithCause(t *testing.T) {
	originalErr := &Error{
		Code:    400,
		Message: "Bad request",
		Details: "Some details",
		Stack:   "stack trace",
	}

	cause := fmt.Errorf("caused by this")
	result := originalErr.WithCause(cause)

	assert.Equal(t, originalErr.Code, result.Code)
	assert.Equal(t, originalErr.Message, result.Message)
	assert.Equal(t, originalErr.Details, result.Details)
	assert.Equal(t, originalErr.Stack, result.Stack)
	assert.Equal(t, cause, result.cause)
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		message  string
		expected *Error
	}{
		{
			name:     "包装现有错误",
			err:      fmt.Errorf("original error"),
			code:     400,
			message:  "Bad request",
			expected: &Error{Code: 400, Message: "Bad request"},
		},
		{
			name:     "包装nil错误",
			err:      nil,
			code:     400,
			message:  "Bad request",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.err, tt.code, tt.message)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Code, result.Code)
				assert.Equal(t, tt.expected.Message, result.Message)
				assert.NotEmpty(t, result.Stack)
				if tt.err != nil {
					assert.Equal(t, tt.err, result.cause)
				}
			}
		})
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "格式化包装错误",
			err:      fmt.Errorf("original error"),
			code:     400,
			format:   "Bad request: %s",
			args:     []interface{}{"invalid input"},
			expected: "Bad request: invalid input",
		},
		{
			name:     "包装nil错误",
			err:      nil,
			code:     400,
			format:   "Bad request: %s",
			args:     []interface{}{"invalid input"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrapf(tt.err, tt.code, tt.format, tt.args...)

			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected, result.Message)
				assert.Equal(t, tt.code, result.Code)
				assert.NotEmpty(t, result.Stack)
				assert.Equal(t, tt.err, result.cause)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "相同错误代码",
			err:      New(400, "Bad request"),
			target:   New(400, "Different message"),
			expected: true,
		},
		{
			name:     "不同错误代码",
			err:      New(400, "Bad request"),
			target:   New(500, "Internal server error"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			target:   New(400, "Bad request"),
			expected: false,
		},
		{
			name:     "nil目标",
			err:      New(400, "Bad request"),
			target:   nil,
			expected: false,
		},
		{
			name:     "都是nil",
			err:      nil,
			target:   nil,
			expected: true,
		},
		{
			name:     "非Error类型错误",
			err:      fmt.Errorf("standard error"),
			target:   New(400, "Bad request"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Is(tt.err, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestError_ToResponse(t *testing.T) {
	err := &Error{
		Code:    400,
		Message: "Bad request",
		Details: "Invalid input parameters",
	}

	response := err.ToResponse()

	assert.Equal(t, err.Code, response.Code)
	assert.Equal(t, err.Message, response.Message)
	assert.Equal(t, err.Details, response.Details)
}

func TestFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected *ErrorResponse
	}{
		{
			name: "自定义错误",
			err: &Error{
				Code:    400,
				Message: "Bad request",
				Details: "Invalid input",
			},
			expected: &ErrorResponse{
				Code:    400,
				Message: "Bad request",
				Details: "Invalid input",
			},
		},
		{
			name: "标准错误",
			err:  fmt.Errorf("standard error"),
			expected: &ErrorResponse{
				Code:    500,
				Message: "Internal server error",
				Details: "standard error",
			},
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromError(tt.err)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Code, result.Code)
				assert.Equal(t, tt.expected.Message, result.Message)
				assert.Equal(t, tt.expected.Details, result.Details)
			}
		})
	}
}

func TestCommonErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected int
	}{
		{name: "ErrInternalServer", err: ErrInternalServer, expected: 500},
		{name: "ErrBadRequest", err: ErrBadRequest, expected: 400},
		{name: "ErrUnauthorized", err: ErrUnauthorized, expected: 401},
		{name: "ErrForbidden", err: ErrForbidden, expected: 403},
		{name: "ErrNotFound", err: ErrNotFound, expected: 404},
		{name: "ErrConflict", err: ErrConflict, expected: 409},
		{name: "ErrValidation", err: ErrValidation, expected: 422},
		{name: "ErrDatabaseConnection", err: ErrDatabaseConnection, expected: 500},
		{name: "ErrDatabaseQuery", err: ErrDatabaseQuery, expected: 500},
		{name: "ErrDatabaseTransaction", err: ErrDatabaseTransaction, expected: 500},
		{name: "ErrInvalidInput", err: ErrInvalidInput, expected: 400},
		{name: "ErrResourceNotFound", err: ErrResourceNotFound, expected: 404},
		{name: "ErrResourceExists", err: ErrResourceExists, expected: 409},
		{name: "ErrPermissionDenied", err: ErrPermissionDenied, expected: 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
			assert.NotEmpty(t, tt.err.Stack)
		})
	}
}

func TestGetStackTrace(t *testing.T) {
	// 测试堆栈跟踪功能
	err := New(500, "Test error")

	// 验证堆栈跟踪不为空
	assert.NotEmpty(t, err.Stack)

	// 验证堆栈跟踪包含文件路径和行号
	assert.Contains(t, err.Stack, ":")
}

// 基准测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(400, "Bad request")
	}
}

func BenchmarkWrap(b *testing.B) {
	originalErr := fmt.Errorf("original error")
	for i := 0; i < b.N; i++ {
		_ = Wrap(originalErr, 400, "Bad request")
	}
}

func BenchmarkWrapf(b *testing.B) {
	originalErr := fmt.Errorf("original error")
	for i := 0; i < b.N; i++ {
		_ = Wrapf(originalErr, 400, "Bad request: %s", "invalid input")
	}
}

func BenchmarkIs(b *testing.B) {
	err1 := New(400, "Bad request")
	err2 := New(400, "Different message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Is(err1, err2)
	}
}
