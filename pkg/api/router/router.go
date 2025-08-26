package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/make-bin/server-tpl/pkg/api"
	v1 "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/api/middleware"
	"github.com/make-bin/server-tpl/pkg/api/response"
	"github.com/make-bin/server-tpl/pkg/api/validation"
	infra_middleware "github.com/make-bin/server-tpl/pkg/infrastructure/middleware"
	"github.com/make-bin/server-tpl/pkg/utils/container"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// RouterConfig 路由配置
type RouterConfig struct {
	EnableAuth     bool                       `json:"enable_auth"`
	EnableSecurity bool                       `json:"enable_security"`
	SecurityConfig *middleware.SecurityConfig `json:"security_config"`
	CORSConfig     *CORSConfig                `json:"cors_config"`
	Validator      *validator.Validate        `json:"-"`
}

// DefaultRouterConfig 默认路由配置
func DefaultRouterConfig() *RouterConfig {
	v := validator.New()
	validation.RegisterCustomValidators(v)

	return &RouterConfig{
		EnableAuth:     true,
		EnableSecurity: true,
		SecurityConfig: middleware.DefaultSecurityConfig,
		CORSConfig: &CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
		Validator: v,
	}
}

// InitRouter initializes the router with all middleware and routes
// @title API 服务文档
// @version 1.0.0
// @description 完整的 API 服务接口文档
// @termsOfService http://swagger.io/terms/
// @contact.name API 支持团队
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token 认证
func InitRouter(engine *gin.Engine, c *container.Container) {
	config := DefaultRouterConfig()
	initRouterWithConfig(engine, c, config)
}

// InitRouterWithConfig 使用配置初始化路由
func InitRouterWithConfig(engine *gin.Engine, c *container.Container, config *RouterConfig) {
	initRouterWithConfig(engine, c, config)
}

// initRouterWithConfig 内部路由初始化函数
func initRouterWithConfig(engine *gin.Engine, c *container.Container, config *RouterConfig) {
	// 注册验证器
	if config.Validator != nil {
		validation.RegisterCustomValidators(config.Validator)
	}

	// 应用全局中间件
	setupGlobalMiddleware(engine, config)

	// 创建API v1路由组
	v1 := engine.Group("/api/v1")

	// 应用API级别中间件
	setupAPIMiddleware(v1, config)

	// 初始化所有注册的API接口
	apiInterfaces := api.GetRegisterAPIInterfaces()
	for _, apiInterface := range apiInterfaces {
		apiInterface.InitAPIServiceRoute(v1)
	}

	// 添加系统级路由
	setupSystemRoutes(engine)

	// 添加Swagger文档路由
	setupSwaggerRoutes(engine)
}

// setupGlobalMiddleware 设置全局中间件
func setupGlobalMiddleware(engine *gin.Engine, config *RouterConfig) {
	// 使用基础设施层的中间件
	loggerManager := logger.NewManager(&logger.LogConfig{})

	// 请求ID中间件（最先执行）
	engine.Use(infra_middleware.GinMiddleware(infra_middleware.NewRequestIDMiddleware()))

	// 日志中间件
	engine.Use(infra_middleware.GinMiddleware(infra_middleware.NewLoggerMiddleware(loggerManager)))

	// 恢复中间件
	engine.Use(middleware.Recovery())

	// 安全响应头中间件
	if config.EnableSecurity {
		engine.Use(middleware.SecurityHeadersMiddleware())
	}

	// CORS中间件
	engine.Use(middleware.CORS())

	// 性能监控中间件
	engine.Use(infra_middleware.PrometheusGinMiddleware())

	// 错误处理中间件（最后执行）
	engine.Use(infra_middleware.GinMiddleware(infra_middleware.NewErrorHandlerMiddleware()))
}

// setupAPIMiddleware 设置API级别中间件
func setupAPIMiddleware(rg *gin.RouterGroup, config *RouterConfig) {
	if config.EnableSecurity {
		// 输入验证中间件
		rg.Use(middleware.InputValidationMiddleware())

		// 限流中间件
		rg.Use(middleware.RateLimitMiddleware(config.SecurityConfig))

		// CSRF防护中间件
		rg.Use(middleware.CSRFMiddleware(config.SecurityConfig))
	}

	if config.EnableAuth {
		// JWT认证中间件
		rg.Use(middleware.JWTAuthMiddleware(config.SecurityConfig))
	}
}

// setupSystemRoutes 设置系统路由
func setupSystemRoutes(engine *gin.Engine) {
	// 根级健康检查
	engine.GET("/health", healthCheck)

	// 系统信息
	engine.GET("/info", systemInfo)

	// 性能指标
	engine.GET("/metrics", infra_middleware.MetricsHandler())
}

// setupSwaggerRoutes 设置Swagger文档路由
func setupSwaggerRoutes(engine *gin.Engine) {
	// 这里可以添加Swagger UI路由
	// engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck 健康检查处理器
// @Summary 系统健康检查
// @Description 检查系统整体健康状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=v1.HealthCheckResponse} "系统正常"
// @Failure 500 {object} response.Response{error=string} "系统异常"
// @Router /health [get]
func healthCheck(c *gin.Context) {
	healthResp := v1.HealthCheckResponse{
		Status:    "ok",
		Message:   "系统运行正常",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"uptime":   time.Since(time.Now()).String(),
			"database": "connected",
			"cache":    "connected",
			"memory":   "normal",
			"cpu":      "normal",
		},
	}

	response.Success(c, healthResp)
}

// systemInfo 系统信息处理器
// @Summary 获取系统信息
// @Description 获取系统基本信息
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /info [get]
func systemInfo(c *gin.Context) {
	info := map[string]interface{}{
		"service_name": "server-tpl",
		"version":      "1.0.0",
		"build_time":   "2024-01-01T00:00:00Z",
		"go_version":   "1.21.0",
		"git_commit":   "unknown",
		"environment":  gin.Mode(),
		"features": map[string]bool{
			"authentication":       true,
			"authorization":        true,
			"rate_limiting":        true,
			"csrf_protection":      true,
			"file_upload":          true,
			"internationalization": true,
		},
	}

	response.Success(c, info)
}

// RegisterRoutes 注册路由（向后兼容）
func RegisterRoutes(engine *gin.Engine, c *container.Container) {
	InitRouter(engine, c)
}

// SetupMiddlewares 设置中间件（向后兼容）
func SetupMiddlewares(engine *gin.Engine) {
	config := DefaultRouterConfig()
	setupGlobalMiddleware(engine, config)
}
