package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, "Internal Server Error: %s", err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
