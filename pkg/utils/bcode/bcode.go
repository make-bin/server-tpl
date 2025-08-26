package bcode

import (
	"fmt"
	"sync"
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

// Business error codes
const (
	// User related errors (30000-30999)
	CodeUserNotFound             = 30000
	CodeUserExists               = 30001
	CodeUserDisabled             = 30002
	CodeUserLocked               = 30003
	CodeUserDeleted              = 30004
	CodePasswordError            = 30005
	CodePasswordExpired          = 30006
	CodePasswordTooWeak          = 30007
	CodeEmailExists              = 30008
	CodePhoneExists              = 30009
	CodeUsernameExists           = 30010
	CodeUserProfileIncomplete    = 30011
	CodeUserVerificationRequired = 30012
	CodeUserAccountSuspended     = 30013

	// Order related errors (31000-31999)
	CodeOrderNotFound        = 31000
	CodeOrderExists          = 31001
	CodeOrderCancelled       = 31002
	CodeOrderCompleted       = 31003
	CodeOrderExpired         = 31004
	CodeOrderPaid            = 31005
	CodeOrderRefunded        = 31006
	CodeOrderProcessing      = 31007
	CodeOrderShipped         = 31008
	CodeOrderDelivered       = 31009
	CodeOrderReturned        = 31010
	CodeOrderAmountInvalid   = 31011
	CodeOrderQuantityInvalid = 31012
	CodeOrderStatusInvalid   = 31013

	// Payment related errors (32000-32999)
	CodePaymentFailed           = 32000
	CodePaymentTimeout          = 32001
	CodePaymentCancelled        = 32002
	CodePaymentRefunded         = 32003
	CodePaymentProcessing       = 32004
	CodePaymentCompleted        = 32005
	CodeInsufficientBalance     = 32006
	CodePaymentMethodInvalid    = 32007
	CodePaymentAmountInvalid    = 32008
	CodePaymentCurrencyInvalid  = 32009
	CodePaymentGatewayError     = 32010
	CodePaymentSignatureInvalid = 32011
	CodePaymentOrderNotFound    = 32012
	CodePaymentDuplicate        = 32013

	// Product related errors (33000-33999)
	CodeProductNotFound              = 33000
	CodeProductExists                = 33001
	CodeProductDisabled              = 33002
	CodeProductDeleted               = 33003
	CodeProductOutOfStock            = 33004
	CodeProductPriceInvalid          = 33005
	CodeProductCategoryInvalid       = 33006
	CodeProductImageInvalid          = 33007
	CodeProductDescriptionInvalid    = 33008
	CodeProductSkuExists             = 33009
	CodeProductBarcodeExists         = 33010
	CodeProductInventoryInsufficient = 33011
	CodeProductExpired               = 33012
	CodeProductRecalled              = 33013

	// Application related errors (34000-34999)
	CodeApplicationNotFound           = 34000
	CodeApplicationExists             = 34001
	CodeApplicationDisabled           = 34002
	CodeApplicationVersionMismatch    = 34003
	CodeApplicationNameRequired       = 34004
	CodeApplicationNameTooLong        = 34005
	CodeApplicationDescriptionTooLong = 34006
)

// ErrorCode represents a business error code definition
type ErrorCode struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Module      string `json:"module"`
	HTTPStatus  int    `json:"http_status"`
	Retryable   bool   `json:"retryable"`
	LogLevel    string `json:"log_level"`
	Solution    string `json:"solution"`
}

// ErrorCodeRegistry manages error code registration and retrieval
type ErrorCodeRegistry struct {
	codes map[int]*ErrorCode
	mutex sync.RWMutex
}

// NewErrorCodeRegistry creates a new error code registry
func NewErrorCodeRegistry() *ErrorCodeRegistry {
	registry := &ErrorCodeRegistry{
		codes: make(map[int]*ErrorCode),
	}

	// Register default error codes
	registry.registerDefaultCodes()

	return registry
}

// RegisterErrorCode registers a new error code
func (r *ErrorCodeRegistry) RegisterErrorCode(code *ErrorCode) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.codes[code.Code]; exists {
		return fmt.Errorf("error code %d already exists", code.Code)
	}

	r.codes[code.Code] = code
	return nil
}

// GetErrorCode retrieves an error code by code number
func (r *ErrorCodeRegistry) GetErrorCode(code int) (*ErrorCode, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	errorCode, exists := r.codes[code]
	return errorCode, exists
}

// GetAllErrorCodes returns all registered error codes
func (r *ErrorCodeRegistry) GetAllErrorCodes() map[int]*ErrorCode {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[int]*ErrorCode)
	for k, v := range r.codes {
		result[k] = v
	}
	return result
}

// GetErrorsByCategory returns error codes by category
func (r *ErrorCodeRegistry) GetErrorsByCategory(category string) []*ErrorCode {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*ErrorCode
	for _, code := range r.codes {
		if code.Category == category {
			result = append(result, code)
		}
	}
	return result
}

// GetErrorsByModule returns error codes by module
func (r *ErrorCodeRegistry) GetErrorsByModule(module string) []*ErrorCode {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*ErrorCode
	for _, code := range r.codes {
		if code.Module == module {
			result = append(result, code)
		}
	}
	return result
}

// ErrorCategory represents an error category
type ErrorCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CodeRange   string `json:"code_range"`
	Module      string `json:"module"`
}

// Error categories
var ErrorCategories = []ErrorCategory{
	{
		Name:        "系统错误",
		Description: "系统内部错误，包括数据库、缓存、网络等",
		CodeRange:   "10000-19999",
		Module:      "system",
	},
	{
		Name:        "客户端错误",
		Description: "客户端请求错误，包括参数验证、权限等",
		CodeRange:   "20000-29999",
		Module:      "client",
	},
	{
		Name:        "业务错误",
		Description: "业务逻辑错误，包括用户、订单、支付等",
		CodeRange:   "30000-99999",
		Module:      "business",
	},
	{
		Name:        "第三方服务错误",
		Description: "第三方服务调用错误",
		CodeRange:   "100000-199999",
		Module:      "third_party",
	},
}

// registerDefaultCodes registers default error codes
func (r *ErrorCodeRegistry) registerDefaultCodes() {
	defaultCodes := []*ErrorCode{
		// User related errors
		{Code: CodeUserNotFound, Message: "用户不存在", Category: "用户管理", Module: "user", HTTPStatus: 404, Retryable: false, LogLevel: "warn"},
		{Code: CodeUserExists, Message: "用户已存在", Category: "用户管理", Module: "user", HTTPStatus: 409, Retryable: false, LogLevel: "warn"},
		{Code: CodeUserDisabled, Message: "用户已禁用", Category: "用户管理", Module: "user", HTTPStatus: 403, Retryable: false, LogLevel: "warn"},
		{Code: CodeUserLocked, Message: "用户已锁定", Category: "用户管理", Module: "user", HTTPStatus: 423, Retryable: false, LogLevel: "warn"},
		{Code: CodePasswordError, Message: "密码错误", Category: "用户管理", Module: "user", HTTPStatus: 401, Retryable: false, LogLevel: "warn"},

		// Order related errors
		{Code: CodeOrderNotFound, Message: "订单不存在", Category: "订单管理", Module: "order", HTTPStatus: 404, Retryable: false, LogLevel: "warn"},
		{Code: CodeOrderCancelled, Message: "订单已取消", Category: "订单管理", Module: "order", HTTPStatus: 400, Retryable: false, LogLevel: "info"},
		{Code: CodeOrderCompleted, Message: "订单已完成", Category: "订单管理", Module: "order", HTTPStatus: 400, Retryable: false, LogLevel: "info"},

		// Payment related errors
		{Code: CodePaymentFailed, Message: "支付失败", Category: "支付管理", Module: "payment", HTTPStatus: 400, Retryable: true, LogLevel: "error"},
		{Code: CodePaymentTimeout, Message: "支付超时", Category: "支付管理", Module: "payment", HTTPStatus: 408, Retryable: true, LogLevel: "warn"},
		{Code: CodeInsufficientBalance, Message: "余额不足", Category: "支付管理", Module: "payment", HTTPStatus: 400, Retryable: false, LogLevel: "warn"},

		// Product related errors
		{Code: CodeProductNotFound, Message: "商品不存在", Category: "商品管理", Module: "product", HTTPStatus: 404, Retryable: false, LogLevel: "warn"},
		{Code: CodeProductOutOfStock, Message: "商品缺货", Category: "商品管理", Module: "product", HTTPStatus: 409, Retryable: false, LogLevel: "warn"},

		// Application related errors
		{Code: CodeApplicationNotFound, Message: "应用不存在", Category: "应用管理", Module: "application", HTTPStatus: 404, Retryable: false, LogLevel: "warn"},
		{Code: CodeApplicationExists, Message: "应用已存在", Category: "应用管理", Module: "application", HTTPStatus: 409, Retryable: false, LogLevel: "warn"},
		{Code: CodeApplicationNameRequired, Message: "应用名称必填", Category: "应用管理", Module: "application", HTTPStatus: 400, Retryable: false, LogLevel: "warn"},
		{Code: CodeApplicationNameTooLong, Message: "应用名称过长", Category: "应用管理", Module: "application", HTTPStatus: 400, Retryable: false, LogLevel: "warn"},
	}

	for _, code := range defaultCodes {
		r.codes[code.Code] = code
	}
}

// Global registry instance
var defaultRegistry *ErrorCodeRegistry

func init() {
	defaultRegistry = NewErrorCodeRegistry()
}

// GetDefaultRegistry returns the default error code registry
func GetDefaultRegistry() *ErrorCodeRegistry {
	return defaultRegistry
}

// GetErrorMessage returns the error message for a given code
func GetErrorMessage(code int) string {
	if errorCode, exists := defaultRegistry.GetErrorCode(code); exists {
		return errorCode.Message
	}
	return "未知错误"
}

// GetErrorCode returns the error code definition
func GetErrorCode(code int) (*ErrorCode, bool) {
	return defaultRegistry.GetErrorCode(code)
}

// RegisterErrorCode registers a new error code in the default registry
func RegisterErrorCode(code *ErrorCode) error {
	return defaultRegistry.RegisterErrorCode(code)
}

// BCode represents a business code with message (backward compatibility)
type BCode struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (b *BCode) Error() string {
	return b.Message
}

// NewBCode creates a new business code (backward compatibility)
func NewBCode(code, message string) *BCode {
	return &BCode{
		Code:    code,
		Message: message,
	}
}

// Legacy constants for backward compatibility
const (
	Success                       = "SUCCESS"
	InvalidParameter              = "INVALID_PARAMETER"
	MissingParameter              = "MISSING_PARAMETER"
	InvalidFormat                 = "INVALID_FORMAT"
	AuthenticationFailed          = "AUTHENTICATION_FAILED"
	TokenExpired                  = "TOKEN_EXPIRED"
	TokenInvalid                  = "TOKEN_INVALID"
	InsufficientPrivileges        = "INSUFFICIENT_PRIVILEGES"
	ApplicationNotFound           = "APPLICATION_NOT_FOUND"
	ApplicationAlreadyExists      = "APPLICATION_ALREADY_EXISTS"
	ApplicationNameRequired       = "APPLICATION_NAME_REQUIRED"
	ApplicationNameTooLong        = "APPLICATION_NAME_TOO_LONG"
	ApplicationDescriptionTooLong = "APPLICATION_DESCRIPTION_TOO_LONG"
	DatabaseConnectionFailed      = "DATABASE_CONNECTION_FAILED"
	DatabaseOperationFailed       = "DATABASE_OPERATION_FAILED"
	RecordNotFound                = "RECORD_NOT_FOUND"
	DuplicateRecord               = "DUPLICATE_RECORD"
	ExternalServiceUnavailable    = "EXTERNAL_SERVICE_UNAVAILABLE"
	ExternalServiceTimeout        = "EXTERNAL_SERVICE_TIMEOUT"
	ExternalServiceError          = "EXTERNAL_SERVICE_ERROR"
	FileNotFound                  = "FILE_NOT_FOUND"
	FileUploadFailed              = "FILE_UPLOAD_FAILED"
	FileSizeTooLarge              = "FILE_SIZE_TOO_LARGE"
	InvalidFileType               = "INVALID_FILE_TYPE"
	RateLimitExceeded             = "RATE_LIMIT_EXCEEDED"
	TooManyRequests               = "TOO_MANY_REQUESTS"
	ServiceUnavailable            = "SERVICE_UNAVAILABLE"
	MaintenanceMode               = "MAINTENANCE_MODE"
	InternalError                 = "INTERNAL_ERROR"
)

// Legacy predefined business codes (backward compatibility)
var (
	SuccessCode                       = NewBCode(Success, "操作成功")
	InvalidParameterCode              = NewBCode(InvalidParameter, "参数无效")
	MissingParameterCode              = NewBCode(MissingParameter, "缺少必需参数")
	InvalidFormatCode                 = NewBCode(InvalidFormat, "格式无效")
	AuthenticationFailedCode          = NewBCode(AuthenticationFailed, "认证失败")
	TokenExpiredCode                  = NewBCode(TokenExpired, "令牌已过期")
	TokenInvalidCode                  = NewBCode(TokenInvalid, "令牌无效")
	InsufficientPrivilegesCode        = NewBCode(InsufficientPrivileges, "权限不足")
	ApplicationNotFoundCode           = NewBCode(ApplicationNotFound, "应用不存在")
	ApplicationAlreadyExistsCode      = NewBCode(ApplicationAlreadyExists, "应用已存在")
	ApplicationNameRequiredCode       = NewBCode(ApplicationNameRequired, "应用名称必填")
	ApplicationNameTooLongCode        = NewBCode(ApplicationNameTooLong, "应用名称过长")
	ApplicationDescriptionTooLongCode = NewBCode(ApplicationDescriptionTooLong, "应用描述过长")
	DatabaseConnectionFailedCode      = NewBCode(DatabaseConnectionFailed, "数据库连接失败")
	DatabaseOperationFailedCode       = NewBCode(DatabaseOperationFailed, "数据库操作失败")
	RecordNotFoundCode                = NewBCode(RecordNotFound, "记录不存在")
	DuplicateRecordCode               = NewBCode(DuplicateRecord, "记录重复")
	ExternalServiceUnavailableCode    = NewBCode(ExternalServiceUnavailable, "外部服务不可用")
	ExternalServiceTimeoutCode        = NewBCode(ExternalServiceTimeout, "外部服务超时")
	ExternalServiceErrorCode          = NewBCode(ExternalServiceError, "外部服务错误")
	FileNotFoundCode                  = NewBCode(FileNotFound, "文件不存在")
	FileUploadFailedCode              = NewBCode(FileUploadFailed, "文件上传失败")
	FileSizeTooLargeCode              = NewBCode(FileSizeTooLarge, "文件大小超限")
	InvalidFileTypeCode               = NewBCode(InvalidFileType, "文件类型无效")
	RateLimitExceededCode             = NewBCode(RateLimitExceeded, "请求频率超限")
	TooManyRequestsCode               = NewBCode(TooManyRequests, "请求过多")
	ServiceUnavailableCode            = NewBCode(ServiceUnavailable, "服务不可用")
	MaintenanceModeCode               = NewBCode(MaintenanceMode, "系统维护中")
	InternalErrorCode                 = NewBCode(InternalError, "系统内部错误")
)
