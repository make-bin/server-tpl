package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// Logger returns a gin.HandlerFunc for request logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("[%s] %s %s %d %s %s",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
		return ""
	})
}
