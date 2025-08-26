package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// Recovery returns a middleware that recovers from any panics
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered: %v", recovered)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		c.Abort()
	})
}
