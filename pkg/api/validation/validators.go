package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("idcard", validateIDCard)
	v.RegisterValidation("chinese", validateChinese)
	v.RegisterValidation("url_path", validateURLPath)
	v.RegisterValidation("app_name", validateAppName)
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 允许空值，由required标签控制
	}

	// 中国手机号验证
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true // 允许空值，由required标签控制
	}

	// 用户名规则：3-50位，字母、数字、下划线、中划线
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	pattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true // 允许空值，由required标签控制
	}

	// 密码长度至少8位
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// 检查密码复杂度
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 至少包含大写字母、小写字母、数字和特殊字符中的三种
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if idCard == "" {
		return true // 允许空值，由required标签控制
	}

	// 18位身份证号验证
	if len(idCard) != 18 {
		return false
	}

	pattern := `^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	if !matched {
		return false
	}

	// 校验位验证
	return validateIDCardChecksum(idCard)
}

// validateIDCardChecksum 验证身份证校验位
func validateIDCardChecksum(idCard string) bool {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		digit := int(idCard[i] - '0')
		sum += digit * weights[i]
	}

	checkIndex := sum % 11
	expectedCheck := checkCodes[checkIndex]
	actualCheck := string(idCard[17])

	return expectedCheck == actualCheck || (expectedCheck == "X" && actualCheck == "x")
}

// validateChinese 验证中文字符
func validateChinese(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	if text == "" {
		return true // 允许空值，由required标签控制
	}

	pattern := `^[\x{4e00}-\x{9fa5}]+$`
	matched, _ := regexp.MatchString(pattern, text)
	return matched
}

// validateURLPath 验证URL路径
func validateURLPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // 允许空值，由required标签控制
	}

	// URL路径验证
	pattern := `^\/[a-zA-Z0-9\/_-]*$`
	matched, _ := regexp.MatchString(pattern, path)
	return matched
}

// validateAppName 验证应用名称
func validateAppName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return true // 允许空值，由required标签控制
	}

	// 应用名称规则：1-100位，字母、数字、中文、下划线、中划线、空格
	if len(name) < 1 || len(name) > 100 {
		return false
	}

	pattern := `^[a-zA-Z0-9\x{4e00}-\x{9fa5}_\- ]+$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}

// ValidationRule 验证规则结构
type ValidationRule struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// CommonValidationRules 常用验证规则
var CommonValidationRules = map[string][]ValidationRule{
	"user": {
		{Field: "username", Rule: "required,min=3,max=50,username", Message: "用户名必须为3-50位字母、数字、下划线、中划线组合"},
		{Field: "email", Rule: "required,email", Message: "请提供有效的邮箱地址"},
		{Field: "password", Rule: "required,min=8,password", Message: "密码至少8位，包含大小写字母、数字、特殊字符中的至少3种"},
		{Field: "phone", Rule: "omitempty,phone", Message: "请提供有效的手机号"},
		{Field: "age", Rule: "omitempty,gte=0,lte=150", Message: "年龄必须在0-150之间"},
	},
	"application": {
		{Field: "name", Rule: "required,min=1,max=100,app_name", Message: "应用名称必须为1-100位有效字符"},
		{Field: "description", Rule: "omitempty,max=500", Message: "应用描述不能超过500字符"},
	},
}

// GetValidationRules 获取验证规则
func GetValidationRules(category string) []ValidationRule {
	if rules, exists := CommonValidationRules[category]; exists {
		return rules
	}
	return nil
}

// ValidateStruct 验证结构体
func ValidateStruct(v *validator.Validate, data interface{}) error {
	return v.Struct(data)
}

// ValidateVar 验证单个变量
func ValidateVar(v *validator.Validate, field interface{}, tag string) error {
	return v.Var(field, tag)
}

// CustomError 自定义验证错误
type CustomError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error 实现error接口
func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError 创建自定义验证错误
func NewCustomError(field, message, value string) *CustomError {
	return &CustomError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// ValidationErrors 验证错误集合
type ValidationErrors []*CustomError

// Error 实现error接口
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return e[0].Error()
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, message, value string) {
	*e = append(*e, NewCustomError(field, message, value))
}

// HasErrors 是否有错误
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Len 错误数量
func (e ValidationErrors) Len() int {
	return len(e)
}

// First 获取第一个错误
func (e ValidationErrors) First() *CustomError {
	if len(e) > 0 {
		return e[0]
	}
	return nil
}

// All 获取所有错误
func (e ValidationErrors) All() []*CustomError {
	return e
}
