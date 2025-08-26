package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/errors"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// ErrorHandler handles errors in a centralized way
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  errors.ErrInternalServer,
			})
		} else if err, ok := recovered.(*errors.AppError); ok {
			logger.Error("Application error: %v", err)
			c.JSON(err.HTTPStatus, gin.H{
				"error": err.Message,
				"code":  err.Code,
			})
		} else {
			logger.Error("Unknown panic recovered: %v", recovered)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  errors.ErrInternalServer,
			})
		}
		c.Abort()
	})
}
