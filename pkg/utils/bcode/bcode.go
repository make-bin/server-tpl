package bcode

import (
	"fmt"

	"github.com/make-bin/server-tpl/pkg/utils/errors"
)

// BCode 业务错误码
type BCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Module  string `json:"module"`
}

// Error 实现error接口
func (b *BCode) Error() string {
	return fmt.Sprintf("[%s] %d: %s", b.Module, b.Code, b.Message)
}

// ToError 转换为Error类型
func (b *BCode) ToError() *errors.Error {
	return errors.New(b.Code, b.Message)
}

// WithDetails 添加详细信息
func (b *BCode) WithDetails(details string) *errors.Error {
	return errors.New(b.Code, b.Message).WithDetails(details)
}

// WithCause 添加原始错误
func (b *BCode) WithCause(cause error) *errors.Error {
	return errors.New(b.Code, b.Message).WithCause(cause)
}

// New 创建新的业务错误码
func New(code int, message, module string) *BCode {
	return &BCode{
		Code:    code,
		Message: message,
		Module:  module,
	}
}

// 错误码范围定义
const (
	// 系统级错误码 (1000-1999)
	SystemErrorStart = 1000
	SystemErrorEnd   = 1999

	// 用户模块错误码 (2000-2999)
	UserErrorStart = 2000
	UserErrorEnd   = 2999

	// 应用模块错误码 (3000-3999)
	ApplicationErrorStart = 3000
	ApplicationErrorEnd   = 3999

	// 配置模块错误码 (4000-4999)
	ConfigErrorStart = 4000
	ConfigErrorEnd   = 4999

	// 数据库模块错误码 (5000-5999)
	DatabaseErrorStart = 5000
	DatabaseErrorEnd   = 5999

	// 认证授权错误码 (6000-6999)
	AuthErrorStart = 6000
	AuthErrorEnd   = 6999

	// 文件操作错误码 (7000-7999)
	FileErrorStart = 7000
	FileErrorEnd   = 7999

	// 外部服务错误码 (8000-8999)
	ExternalErrorStart = 8000
	ExternalErrorEnd   = 8999
)

// 系统级错误码
var (
	// 系统通用错误 (1000-1099)
	ErrSystemInternal    = New(1000, "System internal error", "SYSTEM")
	ErrSystemUnavailable = New(1001, "System temporarily unavailable", "SYSTEM")
	ErrSystemMaintenance = New(1002, "System under maintenance", "SYSTEM")
	ErrSystemTimeout     = New(1003, "System request timeout", "SYSTEM")
	ErrSystemOverload    = New(1004, "System overload", "SYSTEM")

	// 配置错误 (1100-1199)
	ErrConfigNotFound = New(1100, "Configuration not found", "CONFIG")
	ErrConfigInvalid  = New(1101, "Invalid configuration", "CONFIG")
	ErrConfigParse    = New(1102, "Configuration parse error", "CONFIG")
	ErrConfigRequired = New(1103, "Required configuration missing", "CONFIG")

	// 网络错误 (1200-1299)
	ErrNetworkUnreachable = New(1200, "Network unreachable", "NETWORK")
	ErrNetworkTimeout     = New(1201, "Network timeout", "NETWORK")
	ErrNetworkRefused     = New(1202, "Network connection refused", "NETWORK")
)

// 用户模块错误码
var (
	// 用户认证错误 (2000-2099)
	ErrUserNotFound       = New(2000, "User not found", "USER")
	ErrUserAlreadyExists  = New(2001, "User already exists", "USER")
	ErrUserPasswordWrong  = New(2002, "Wrong password", "USER")
	ErrUserAccountLocked  = New(2003, "User account locked", "USER")
	ErrUserAccountExpired = New(2004, "User account expired", "USER")
	ErrUserDisabled       = New(2005, "User account disabled", "USER")

	// 用户权限错误 (2100-2199)
	ErrUserPermissionDenied = New(2100, "Permission denied", "USER")
	ErrUserInsufficientRole = New(2101, "Insufficient role", "USER")
	ErrUserTokenExpired     = New(2102, "Token expired", "USER")
	ErrUserTokenInvalid     = New(2103, "Invalid token", "USER")
	ErrUserTokenMissing     = New(2104, "Token missing", "USER")

	// 用户数据错误 (2200-2299)
	ErrUserDataInvalid  = New(2200, "Invalid user data", "USER")
	ErrUserDataConflict = New(2201, "User data conflict", "USER")
	ErrUserDataRequired = New(2202, "Required user data missing", "USER")
)

// 应用模块错误码
var (
	// 应用通用错误 (3000-3099)
	ErrAppNotFound       = New(3000, "Application not found", "APP")
	ErrAppAlreadyExists  = New(3001, "Application already exists", "APP")
	ErrAppNameInvalid    = New(3002, "Invalid application name", "APP")
	ErrAppVersionInvalid = New(3003, "Invalid application version", "APP")
	ErrAppStatusInvalid  = New(3004, "Invalid application status", "APP")

	// 应用操作错误 (3100-3199)
	ErrAppCreateFailed = New(3100, "Failed to create application", "APP")
	ErrAppUpdateFailed = New(3101, "Failed to update application", "APP")
	ErrAppDeleteFailed = New(3102, "Failed to delete application", "APP")
	ErrAppDeployFailed = New(3103, "Failed to deploy application", "APP")
	ErrAppStartFailed  = New(3104, "Failed to start application", "APP")
	ErrAppStopFailed   = New(3105, "Failed to stop application", "APP")

	// 应用配置错误 (3200-3299)
	ErrAppConfigInvalid  = New(3200, "Invalid application configuration", "APP")
	ErrAppConfigMissing  = New(3201, "Application configuration missing", "APP")
	ErrAppConfigConflict = New(3202, "Application configuration conflict", "APP")
)

// 数据库模块错误码
var (
	// 数据库连接错误 (5000-5099)
	ErrDBConnectionFailed  = New(5000, "Database connection failed", "DB")
	ErrDBConnectionTimeout = New(5001, "Database connection timeout", "DB")
	ErrDBConnectionRefused = New(5002, "Database connection refused", "DB")

	// 数据库查询错误 (5100-5199)
	ErrDBQueryFailed  = New(5100, "Database query failed", "DB")
	ErrDBQueryTimeout = New(5101, "Database query timeout", "DB")
	ErrDBQuerySyntax  = New(5102, "Database query syntax error", "DB")

	// 数据库事务错误 (5200-5299)
	ErrDBTransactionFailed   = New(5200, "Database transaction failed", "DB")
	ErrDBTransactionRollback = New(5201, "Database transaction rollback", "DB")
	ErrDBTransactionCommit   = New(5202, "Database transaction commit failed", "DB")

	// 数据库约束错误 (5300-5399)
	ErrDBConstraintViolation = New(5300, "Database constraint violation", "DB")
	ErrDBUniqueViolation     = New(5301, "Database unique constraint violation", "DB")
	ErrDBForeignKeyViolation = New(5302, "Database foreign key violation", "DB")
	ErrDBNotNullViolation    = New(5303, "Database not null constraint violation", "DB")
)

// 认证授权错误码
var (
	// 认证错误 (6000-6099)
	ErrAuthFailed  = New(6000, "Authentication failed", "AUTH")
	ErrAuthExpired = New(6001, "Authentication expired", "AUTH")
	ErrAuthInvalid = New(6002, "Invalid authentication", "AUTH")
	ErrAuthMissing = New(6003, "Authentication missing", "AUTH")

	// 授权错误 (6100-6199)
	ErrAuthPermissionDenied = New(6100, "Permission denied", "AUTH")
	ErrAuthInsufficientRole = New(6101, "Insufficient role", "AUTH")
	ErrAuthResourceAccess   = New(6102, "Resource access denied", "AUTH")
	ErrAuthOperationDenied  = New(6103, "Operation denied", "AUTH")
)

// 文件操作错误码
var (
	// 文件通用错误 (7000-7099)
	ErrFileNotFound      = New(7000, "File not found", "FILE")
	ErrFileAlreadyExists = New(7001, "File already exists", "FILE")
	ErrFileAccessDenied  = New(7002, "File access denied", "FILE")
	ErrFileReadFailed    = New(7003, "File read failed", "FILE")
	ErrFileWriteFailed   = New(7004, "File write failed", "FILE")
	ErrFileDeleteFailed  = New(7005, "File delete failed", "FILE")

	// 文件格式错误 (7100-7199)
	ErrFileFormatInvalid  = New(7100, "Invalid file format", "FILE")
	ErrFileSizeExceeded   = New(7101, "File size exceeded", "FILE")
	ErrFileTypeNotAllowed = New(7102, "File type not allowed", "FILE")
)

// 外部服务错误码
var (
	// 外部服务通用错误 (8000-8099)
	ErrExternalServiceUnavailable = New(8000, "External service unavailable", "EXTERNAL")
	ErrExternalServiceTimeout     = New(8001, "External service timeout", "EXTERNAL")
	ErrExternalServiceError       = New(8002, "External service error", "EXTERNAL")
	ErrExternalServiceInvalid     = New(8003, "External service invalid response", "EXTERNAL")
)

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	switch {
	case code >= SystemErrorStart && code <= SystemErrorEnd:
		return getSystemErrorMessage(code)
	case code >= UserErrorStart && code <= UserErrorEnd:
		return getUserErrorMessage(code)
	case code >= ApplicationErrorStart && code <= ApplicationErrorEnd:
		return getApplicationErrorMessage(code)
	case code >= DatabaseErrorStart && code <= DatabaseErrorEnd:
		return getDatabaseErrorMessage(code)
	case code >= AuthErrorStart && code <= AuthErrorEnd:
		return getAuthErrorMessage(code)
	case code >= FileErrorStart && code <= FileErrorEnd:
		return getFileErrorMessage(code)
	case code >= ExternalErrorStart && code <= ExternalErrorEnd:
		return getExternalErrorMessage(code)
	default:
		return "Unknown error"
	}
}

// getSystemErrorMessage 获取系统错误信息
func getSystemErrorMessage(code int) string {
	switch code {
	case 1000:
		return "System internal error"
	case 1001:
		return "System temporarily unavailable"
	case 1002:
		return "System under maintenance"
	case 1003:
		return "System request timeout"
	case 1004:
		return "System overload"
	case 1100:
		return "Configuration not found"
	case 1101:
		return "Invalid configuration"
	case 1102:
		return "Configuration parse error"
	case 1103:
		return "Required configuration missing"
	case 1200:
		return "Network unreachable"
	case 1201:
		return "Network timeout"
	case 1202:
		return "Network connection refused"
	default:
		return "System error"
	}
}

// getUserErrorMessage 获取用户错误信息
func getUserErrorMessage(code int) string {
	switch code {
	case 2000:
		return "User not found"
	case 2001:
		return "User already exists"
	case 2002:
		return "Wrong password"
	case 2003:
		return "User account locked"
	case 2004:
		return "User account expired"
	case 2005:
		return "User account disabled"
	case 2100:
		return "Permission denied"
	case 2101:
		return "Insufficient role"
	case 2102:
		return "Token expired"
	case 2103:
		return "Invalid token"
	case 2104:
		return "Token missing"
	case 2200:
		return "Invalid user data"
	case 2201:
		return "User data conflict"
	case 2202:
		return "Required user data missing"
	default:
		return "User error"
	}
}

// getApplicationErrorMessage 获取应用错误信息
func getApplicationErrorMessage(code int) string {
	switch code {
	case 3000:
		return "Application not found"
	case 3001:
		return "Application already exists"
	case 3002:
		return "Invalid application name"
	case 3003:
		return "Invalid application version"
	case 3004:
		return "Invalid application status"
	case 3100:
		return "Failed to create application"
	case 3101:
		return "Failed to update application"
	case 3102:
		return "Failed to delete application"
	case 3103:
		return "Failed to deploy application"
	case 3104:
		return "Failed to start application"
	case 3105:
		return "Failed to stop application"
	case 3200:
		return "Invalid application configuration"
	case 3201:
		return "Application configuration missing"
	case 3202:
		return "Application configuration conflict"
	default:
		return "Application error"
	}
}

// getDatabaseErrorMessage 获取数据库错误信息
func getDatabaseErrorMessage(code int) string {
	switch code {
	case 5000:
		return "Database connection failed"
	case 5001:
		return "Database connection timeout"
	case 5002:
		return "Database connection refused"
	case 5100:
		return "Database query failed"
	case 5101:
		return "Database query timeout"
	case 5102:
		return "Database query syntax error"
	case 5200:
		return "Database transaction failed"
	case 5201:
		return "Database transaction rollback"
	case 5202:
		return "Database transaction commit failed"
	case 5300:
		return "Database constraint violation"
	case 5301:
		return "Database unique constraint violation"
	case 5302:
		return "Database foreign key violation"
	case 5303:
		return "Database not null constraint violation"
	default:
		return "Database error"
	}
}

// getAuthErrorMessage 获取认证授权错误信息
func getAuthErrorMessage(code int) string {
	switch code {
	case 6000:
		return "Authentication failed"
	case 6001:
		return "Authentication expired"
	case 6002:
		return "Invalid authentication"
	case 6003:
		return "Authentication missing"
	case 6100:
		return "Permission denied"
	case 6101:
		return "Insufficient role"
	case 6102:
		return "Resource access denied"
	case 6103:
		return "Operation denied"
	default:
		return "Authentication error"
	}
}

// getFileErrorMessage 获取文件错误信息
func getFileErrorMessage(code int) string {
	switch code {
	case 7000:
		return "File not found"
	case 7001:
		return "File already exists"
	case 7002:
		return "File access denied"
	case 7003:
		return "File read failed"
	case 7004:
		return "File write failed"
	case 7005:
		return "File delete failed"
	case 7100:
		return "Invalid file format"
	case 7101:
		return "File size exceeded"
	case 7102:
		return "File type not allowed"
	default:
		return "File error"
	}
}

// getExternalErrorMessage 获取外部服务错误信息
func getExternalErrorMessage(code int) string {
	switch code {
	case 8000:
		return "External service unavailable"
	case 8001:
		return "External service timeout"
	case 8002:
		return "External service error"
	case 8003:
		return "External service invalid response"
	default:
		return "External service error"
	}
}
