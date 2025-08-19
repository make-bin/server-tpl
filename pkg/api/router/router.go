package router

import (
	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/api"
	"github.com/make-bin/server-tpl/pkg/api/middleware"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()

	// 国际化解析
	r.Use(middleware.I18n())

	// 使用自定义错误处理中间件
	r.Use(middleware.ErrorHandler())

	// 添加其他中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 注册所有API接口
		for _, apiInterface := range api.GetRegisterAPIInterfaces() {
			apiInterface.InitAPIServiceRoute(v1)
		}
	}

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		middleware.SuccessResponse(c, gin.H{
			"status": "ok",
		})
	})

	return r
}
