package middleware

import (
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/errors"
	"github.com/make-bin/server-tpl/pkg/utils/i18n"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Errorf("Panic recovered: %s", err)
			logger.Errorf("Stack trace: %s", debug.Stack())

			bundle := i18n.Default()
			locale := getLocale(c)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": bundle.Translate(locale, "error.500"),
				"details": err,
				"locale":  locale,
			})
			return
		}

		if err, ok := recovered.(error); ok {
			logger.Errorf("Panic recovered: %v", err)
			logger.Errorf("Stack trace: %s", debug.Stack())

			bundle := i18n.Default()
			locale := getLocale(c)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": bundle.Translate(locale, "error.500"),
				"details": err.Error(),
				"locale":  locale,
			})
			return
		}

		logger.Errorf("Panic recovered: %v", recovered)
		logger.Errorf("Stack trace: %s", debug.Stack())

		bundle := i18n.Default()
		locale := getLocale(c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": bundle.Translate(locale, "error.500"),
			"details": "Unknown error occurred",
			"locale":  locale,
		})
	})
}

// ErrorResponse 统一错误响应
func ErrorResponse(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 记录错误日志
	logger.WithField("path", c.Request.URL.Path).
		WithField("method", c.Request.Method).
		WithField("error", err.Error()).
		Error("API error occurred")

	// 转换为错误响应
	errorResp := errors.FromError(err)

	// 根据错误码设置HTTP状态码
	statusCode := getStatusCode(errorResp.Code)

	bundle := i18n.Default()
	locale := getLocale(c)
	msgKey := translateErrorKey(errorResp.Code)
	c.JSON(statusCode, gin.H{
		"code":    errorResp.Code,
		"message": bundle.Translate(locale, msgKey),
		"details": errorResp.Details,
		"data":    errorResp.Data,
		"locale":  locale,
	})
}

// getStatusCode 根据错误码获取HTTP状态码
func getStatusCode(errorCode int) int {
	switch {
	case errorCode >= 1000 && errorCode < 2000:
		// 系统错误
		return http.StatusInternalServerError
	case errorCode >= 2000 && errorCode < 3000:
		// 用户错误
		switch {
		case errorCode >= 2000 && errorCode < 2100:
			return http.StatusNotFound
		case errorCode >= 2100 && errorCode < 2200:
			return http.StatusForbidden
		case errorCode >= 2200 && errorCode < 2300:
			return http.StatusBadRequest
		default:
			return http.StatusBadRequest
		}
	case errorCode >= 3000 && errorCode < 4000:
		// 应用错误
		switch {
		case errorCode >= 3000 && errorCode < 3100:
			return http.StatusNotFound
		case errorCode >= 3100 && errorCode < 3200:
			return http.StatusInternalServerError
		case errorCode >= 3200 && errorCode < 3300:
			return http.StatusBadRequest
		default:
			return http.StatusBadRequest
		}
	case errorCode >= 4000 && errorCode < 5000:
		// 配置错误
		return http.StatusInternalServerError
	case errorCode >= 5000 && errorCode < 6000:
		// 数据库错误
		return http.StatusInternalServerError
	case errorCode >= 6000 && errorCode < 7000:
		// 认证授权错误
		switch {
		case errorCode >= 6000 && errorCode < 6100:
			return http.StatusUnauthorized
		case errorCode >= 6100 && errorCode < 6200:
			return http.StatusForbidden
		default:
			return http.StatusUnauthorized
		}
	case errorCode >= 7000 && errorCode < 8000:
		// 文件操作错误
		switch {
		case errorCode >= 7000 && errorCode < 7100:
			return http.StatusNotFound
		case errorCode >= 7100 && errorCode < 7200:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	case errorCode >= 8000 && errorCode < 9000:
		// 外部服务错误
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// SuccessResponse 统一成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	bundle := i18n.Default()
	locale := getLocale(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": bundle.Translate(locale, "common.success"),
		"data":    data,
		"locale":  locale,
	})
}

// CreatedResponse 创建成功响应
func CreatedResponse(c *gin.Context, data interface{}) {
	bundle := i18n.Default()
	locale := getLocale(c)
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": bundle.Translate(locale, "common.created"),
		"data":    data,
		"locale":  locale,
	})
}

// NoContentResponse 无内容响应
func NoContentResponse(c *gin.Context) {
	bundle := i18n.Default()
	locale := getLocale(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": bundle.Translate(locale, "common.success"),
		"data":    nil,
		"locale":  locale,
	})
}

// ListResponse 列表响应
func ListResponse(c *gin.Context, data interface{}, total int64, page, size int) {
	bundle := i18n.Default()
	locale := getLocale(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": bundle.Translate(locale, "common.success"),
		"data": gin.H{
			"list":  data,
			"total": total,
			"page":  page,
			"size":  size,
		},
		"locale": locale,
	})
}

// translateErrorKey converts numeric error code into a translation key
// like "error.404" or defaults to "error.500".
func translateErrorKey(code int) string {
	switch code {
	case 400, 401, 403, 404, 409, 422, 500:
		return "error." + strconv.Itoa(code)
	default:
		if code >= 2000 && code < 2100 {
			return "error.404"
		}
		if code >= 2100 && code < 2200 {
			return "error.403"
		}
		if code >= 2200 && code < 2300 {
			return "error.400"
		}
		return "error.500"
	}
}
