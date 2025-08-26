package response

// 系统级错误码 (10000-19999)
const (
	// 成功
	CodeSuccess = 200

	// 系统错误 (10000-10999)
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

	// 客户端错误 (20000-29999)
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

// 业务错误码 (30000-99999)
const (
	// 用户相关错误 (30000-30999)
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

	// 应用相关错误 (31000-31999)
	CodeAppNotFound           = 31000
	CodeAppExists             = 31001
	CodeAppDisabled           = 31002
	CodeAppDeleted            = 31003
	CodeAppNameInvalid        = 31004
	CodeAppDescriptionTooLong = 31005
	CodeAppConfigInvalid      = 31006
	CodeAppVersionInvalid     = 31007
	CodeAppStatusInvalid      = 31008
	CodeAppPermissionDenied   = 31009
	CodeAppQuotaExceeded      = 31010
	CodeAppDependencyError    = 31011
	CodeAppDeploymentFailed   = 31012
	CodeAppBackupFailed       = 31013

	// 订单相关错误 (32000-32999)
	CodeOrderNotFound        = 32000
	CodeOrderExists          = 32001
	CodeOrderCancelled       = 32002
	CodeOrderCompleted       = 32003
	CodeOrderExpired         = 32004
	CodeOrderPaid            = 32005
	CodeOrderRefunded        = 32006
	CodeOrderProcessing      = 32007
	CodeOrderShipped         = 32008
	CodeOrderDelivered       = 32009
	CodeOrderReturned        = 32010
	CodeOrderAmountInvalid   = 32011
	CodeOrderQuantityInvalid = 32012
	CodeOrderStatusInvalid   = 32013

	// 支付相关错误 (33000-33999)
	CodePaymentFailed           = 33000
	CodePaymentTimeout          = 33001
	CodePaymentCancelled        = 33002
	CodePaymentRefunded         = 33003
	CodePaymentProcessing       = 33004
	CodePaymentCompleted        = 33005
	CodeInsufficientBalance     = 33006
	CodePaymentMethodInvalid    = 33007
	CodePaymentAmountInvalid    = 33008
	CodePaymentCurrencyInvalid  = 33009
	CodePaymentGatewayError     = 33010
	CodePaymentSignatureInvalid = 33011
	CodePaymentOrderNotFound    = 33012
	CodePaymentDuplicate        = 33013

	// 文件相关错误 (34000-34999)
	CodeFileNotFound             = 34000
	CodeFileTooBig               = 34001
	CodeFileTypeNotSupported     = 34002
	CodeFileUploadFailed         = 34003
	CodeFileDeleteFailed         = 34004
	CodeFileCorrupted            = 34005
	CodeFilePermissionDenied     = 34006
	CodeFileStorageQuotaExceeded = 34007
	CodeFileVirusDetected        = 34008
	CodeFileProcessingFailed     = 34009

	// 权限相关错误 (35000-35999)
	CodePermissionDenied    = 35000
	CodeRoleNotFound        = 35001
	CodeRoleExists          = 35002
	CodeInvalidRole         = 35003
	CodePermissionNotFound  = 35004
	CodePermissionExists    = 35005
	CodeInvalidPermission   = 35006
	CodeAccessTokenExpired  = 35007
	CodeRefreshTokenExpired = 35008
	CodeInvalidToken        = 35009
	CodeTokenNotFound       = 35010
	CodeSessionExpired      = 35011
	CodeLoginRequired       = 35012
	CodeAccountLocked       = 35013
)

// 错误码消息映射表
var errorCodeMap = map[int]string{
	// 系统错误
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

	// 客户端错误
	CodeValidationError:      "参数验证失败",
	CodeMissingParameter:     "缺少必需参数",
	CodeInvalidParameter:     "无效参数",
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

	// 用户相关错误
	CodeUserNotFound:             "用户不存在",
	CodeUserExists:               "用户已存在",
	CodeUserDisabled:             "用户已禁用",
	CodeUserLocked:               "用户已锁定",
	CodeUserDeleted:              "用户已删除",
	CodePasswordError:            "密码错误",
	CodePasswordExpired:          "密码已过期",
	CodePasswordTooWeak:          "密码强度不足",
	CodeEmailExists:              "邮箱已存在",
	CodePhoneExists:              "手机号已存在",
	CodeUsernameExists:           "用户名已存在",
	CodeUserProfileIncomplete:    "用户资料不完整",
	CodeUserVerificationRequired: "需要验证用户",
	CodeUserAccountSuspended:     "用户账户已暂停",

	// 应用相关错误
	CodeAppNotFound:           "应用不存在",
	CodeAppExists:             "应用已存在",
	CodeAppDisabled:           "应用已禁用",
	CodeAppDeleted:            "应用已删除",
	CodeAppNameInvalid:        "应用名称无效",
	CodeAppDescriptionTooLong: "应用描述过长",
	CodeAppConfigInvalid:      "应用配置无效",
	CodeAppVersionInvalid:     "应用版本无效",
	CodeAppStatusInvalid:      "应用状态无效",
	CodeAppPermissionDenied:   "应用权限被拒绝",
	CodeAppQuotaExceeded:      "应用配额已超出",
	CodeAppDependencyError:    "应用依赖错误",
	CodeAppDeploymentFailed:   "应用部署失败",
	CodeAppBackupFailed:       "应用备份失败",

	// 订单相关错误
	CodeOrderNotFound:        "订单不存在",
	CodeOrderExists:          "订单已存在",
	CodeOrderCancelled:       "订单已取消",
	CodeOrderCompleted:       "订单已完成",
	CodeOrderExpired:         "订单已过期",
	CodeOrderPaid:            "订单已支付",
	CodeOrderRefunded:        "订单已退款",
	CodeOrderProcessing:      "订单处理中",
	CodeOrderShipped:         "订单已发货",
	CodeOrderDelivered:       "订单已送达",
	CodeOrderReturned:        "订单已退货",
	CodeOrderAmountInvalid:   "订单金额无效",
	CodeOrderQuantityInvalid: "订单数量无效",
	CodeOrderStatusInvalid:   "订单状态无效",

	// 支付相关错误
	CodePaymentFailed:           "支付失败",
	CodePaymentTimeout:          "支付超时",
	CodePaymentCancelled:        "支付已取消",
	CodePaymentRefunded:         "支付已退款",
	CodePaymentProcessing:       "支付处理中",
	CodePaymentCompleted:        "支付已完成",
	CodeInsufficientBalance:     "余额不足",
	CodePaymentMethodInvalid:    "支付方式无效",
	CodePaymentAmountInvalid:    "支付金额无效",
	CodePaymentCurrencyInvalid:  "支付货币无效",
	CodePaymentGatewayError:     "支付网关错误",
	CodePaymentSignatureInvalid: "支付签名无效",
	CodePaymentOrderNotFound:    "支付订单不存在",
	CodePaymentDuplicate:        "重复支付",

	// 文件相关错误
	CodeFileNotFound:             "文件不存在",
	CodeFileTooBig:               "文件过大",
	CodeFileTypeNotSupported:     "不支持的文件类型",
	CodeFileUploadFailed:         "文件上传失败",
	CodeFileDeleteFailed:         "文件删除失败",
	CodeFileCorrupted:            "文件已损坏",
	CodeFilePermissionDenied:     "文件权限被拒绝",
	CodeFileStorageQuotaExceeded: "文件存储配额已超出",
	CodeFileVirusDetected:        "检测到病毒文件",
	CodeFileProcessingFailed:     "文件处理失败",

	// 权限相关错误
	CodePermissionDenied:    "权限被拒绝",
	CodeRoleNotFound:        "角色不存在",
	CodeRoleExists:          "角色已存在",
	CodeInvalidRole:         "无效角色",
	CodePermissionNotFound:  "权限不存在",
	CodePermissionExists:    "权限已存在",
	CodeInvalidPermission:   "无效权限",
	CodeAccessTokenExpired:  "访问令牌已过期",
	CodeRefreshTokenExpired: "刷新令牌已过期",
	CodeInvalidToken:        "无效令牌",
	CodeTokenNotFound:       "令牌不存在",
	CodeSessionExpired:      "会话已过期",
	CodeLoginRequired:       "需要登录",
	CodeAccountLocked:       "账户已锁定",
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code int) string {
	if message, exists := errorCodeMap[code]; exists {
		return message
	}
	return "未知错误"
}

// IsSystemError 判断是否为系统错误
func IsSystemError(code int) bool {
	return code >= 10000 && code < 20000
}

// IsClientError 判断是否为客户端错误
func IsClientError(code int) bool {
	return code >= 20000 && code < 30000
}

// IsBusinessError 判断是否为业务错误
func IsBusinessError(code int) bool {
	return code >= 30000 && code < 100000
}

// GetErrorType 获取错误类型
func GetErrorType(code int) string {
	switch {
	case IsSystemError(code):
		return "system"
	case IsClientError(code):
		return "client"
	case IsBusinessError(code):
		return "business"
	default:
		return "unknown"
	}
}
