package bcode

import (
	"testing"
)

func TestBCode_Error(t *testing.T) {
	tests := []struct {
		name     string
		bcode    *BCode
		expected string
	}{
		{
			name: "user error",
			bcode: &BCode{
				Code:    2000,
				Message: "User not found",
				Module:  "USER",
			},
			expected: "[USER] 2000: User not found",
		},
		{
			name: "app error",
			bcode: &BCode{
				Code:    3000,
				Message: "Application not found",
				Module:  "APP",
			},
			expected: "[APP] 3000: Application not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcode.Error(); got != tt.expected {
				t.Errorf("BCode.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	bcode := New(2000, "User not found", "USER")

	if bcode.Code != 2000 {
		t.Errorf("Expected code 2000, got %d", bcode.Code)
	}

	if bcode.Message != "User not found" {
		t.Errorf("Expected message 'User not found', got %s", bcode.Message)
	}

	if bcode.Module != "USER" {
		t.Errorf("Expected module 'USER', got %s", bcode.Module)
	}
}

func TestToError(t *testing.T) {
	bcode := New(2000, "User not found", "USER")
	err := bcode.ToError()

	if err.Code != 2000 {
		t.Errorf("Expected code 2000, got %d", err.Code)
	}

	if err.Message != "User not found" {
		t.Errorf("Expected message 'User not found', got %s", err.Message)
	}
}

func TestWithDetails(t *testing.T) {
	bcode := New(2000, "User not found", "USER")
	err := bcode.WithDetails("User ID: 123")

	if err.Details != "User ID: 123" {
		t.Errorf("Expected details 'User ID: 123', got %s", err.Details)
	}

	if err.Code != 2000 {
		t.Errorf("Expected code to remain 2000, got %d", err.Code)
	}
}

func TestWithCause(t *testing.T) {
	originalErr := &BCode{Code: 1000, Message: "Original error", Module: "SYSTEM"}
	bcode := New(2000, "User not found", "USER")
	err := bcode.WithCause(originalErr)

	if err.Code != 2000 {
		t.Errorf("Expected code 2000, got %d", err.Code)
	}

	if err.Message != "User not found" {
		t.Errorf("Expected message 'User not found', got %s", err.Message)
	}
}

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "system error",
			code:     1000,
			expected: "System internal error",
		},
		{
			name:     "user not found",
			code:     2000,
			expected: "User not found",
		},
		{
			name:     "app not found",
			code:     3000,
			expected: "Application not found",
		},
		{
			name:     "db connection failed",
			code:     5000,
			expected: "Database connection failed",
		},
		{
			name:     "auth failed",
			code:     6000,
			expected: "Authentication failed",
		},
		{
			name:     "file not found",
			code:     7000,
			expected: "File not found",
		},
		{
			name:     "external service unavailable",
			code:     8000,
			expected: "External service unavailable",
		},
		{
			name:     "unknown error",
			code:     9999,
			expected: "Unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorMessage(tt.code); got != tt.expected {
				t.Errorf("GetErrorMessage(%d) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	// 测试系统错误
	if ErrSystemInternal.Code != 1000 {
		t.Errorf("Expected ErrSystemInternal code 1000, got %d", ErrSystemInternal.Code)
	}

	if ErrSystemInternal.Module != "SYSTEM" {
		t.Errorf("Expected ErrSystemInternal module 'SYSTEM', got %s", ErrSystemInternal.Module)
	}

	// 测试用户错误
	if ErrUserNotFound.Code != 2000 {
		t.Errorf("Expected ErrUserNotFound code 2000, got %d", ErrUserNotFound.Code)
	}

	if ErrUserNotFound.Module != "USER" {
		t.Errorf("Expected ErrUserNotFound module 'USER', got %s", ErrUserNotFound.Module)
	}

	// 测试应用错误
	if ErrAppNotFound.Code != 3000 {
		t.Errorf("Expected ErrAppNotFound code 3000, got %d", ErrAppNotFound.Code)
	}

	if ErrAppNotFound.Module != "APP" {
		t.Errorf("Expected ErrAppNotFound module 'APP', got %s", ErrAppNotFound.Module)
	}

	// 测试数据库错误
	if ErrDBConnectionFailed.Code != 5000 {
		t.Errorf("Expected ErrDBConnectionFailed code 5000, got %d", ErrDBConnectionFailed.Code)
	}

	if ErrDBConnectionFailed.Module != "DB" {
		t.Errorf("Expected ErrDBConnectionFailed module 'DB', got %s", ErrDBConnectionFailed.Module)
	}

	// 测试认证错误
	if ErrAuthFailed.Code != 6000 {
		t.Errorf("Expected ErrAuthFailed code 6000, got %d", ErrAuthFailed.Code)
	}

	if ErrAuthFailed.Module != "AUTH" {
		t.Errorf("Expected ErrAuthFailed module 'AUTH', got %s", ErrAuthFailed.Module)
	}

	// 测试文件错误
	if ErrFileNotFound.Code != 7000 {
		t.Errorf("Expected ErrFileNotFound code 7000, got %d", ErrFileNotFound.Code)
	}

	if ErrFileNotFound.Module != "FILE" {
		t.Errorf("Expected ErrFileNotFound module 'FILE', got %s", ErrFileNotFound.Module)
	}

	// 测试外部服务错误
	if ErrExternalServiceUnavailable.Code != 8000 {
		t.Errorf("Expected ErrExternalServiceUnavailable code 8000, got %d", ErrExternalServiceUnavailable.Code)
	}

	if ErrExternalServiceUnavailable.Module != "EXTERNAL" {
		t.Errorf("Expected ErrExternalServiceUnavailable module 'EXTERNAL', got %s", ErrExternalServiceUnavailable.Module)
	}
}
