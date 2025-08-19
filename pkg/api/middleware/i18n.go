package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/i18n"
)

const localeKey = "locale"

// I18n resolves the request locale and stores it in the Gin context.
// Priority: query "lang" -> header "X-Lang" -> header "Accept-Language" -> default.
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		bundle := i18n.Default()
		locale := bundle.ResolveLocaleFromRequest(c.Request)
		c.Set(localeKey, locale)
		c.Next()
	}
}

// getLocale fetches the locale stored by I18n middleware, falling back to default.
func getLocale(c *gin.Context) string {
	if v, ok := c.Get(localeKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return i18n.Default().SupportedLocales()[0]
}
