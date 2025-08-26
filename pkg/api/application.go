package api

import (
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/make-bin/server-tpl/pkg/api/handler"
	"github.com/make-bin/server-tpl/pkg/domain/service"
)

// ApplicationAPI 应用API结构
type ApplicationAPI struct {
	handler *handler.ApplicationHandler
}

// application 支持依赖注入的应用API结构
type application struct {
	ApplicationService service.ApplicationServiceInterface `inject:""`
	handler            *handler.ApplicationHandler
}

// init 注册API接口和验证器
func init() {
	RegisterAPIInterface(newApplication())
	RegisterValidationInterface("namevalidator", func(fl validator.FieldLevel) bool {
		name, ok := fl.Field().Interface().(string)
		if ok {
			if len(name) == 1 {
				return unicode.IsLower([]rune(name)[0])
			}
			matched, _ := regexp.MatchString("^[a-z][a-z0-9-]*[a-z0-9]$", name)
			return matched
		}
		return true
	})
	RegisterValidationInterface("versionvalidator", func(fl validator.FieldLevel) bool {
		version, ok := fl.Field().Interface().(string)
		if ok {
			matched, _ := regexp.MatchString("^[^\u4e00-\u9fa5]+$", version)
			return matched
		}
		return true
	})
}

// newApplication 创建依赖注入版本的应用API
func newApplication() APIInterface {
	return &application{}
}

// NewApplicationAPI 创建应用API实例
func NewApplicationAPI(applicationService service.ApplicationServiceInterface) *ApplicationAPI {
	return &ApplicationAPI{
		handler: handler.NewApplicationHandler(applicationService),
	}
}

// InitAPIServiceRoute 初始化应用API路由
// @title 应用管理API
// @version 1.0
// @description 应用管理相关接口
// @BasePath /api/v1
func (a *ApplicationAPI) InitAPIServiceRoute(rg *gin.RouterGroup) {
	applicationGroup := rg.Group("/applications")
	{
		// 应用CRUD操作
		applicationGroup.POST("", a.handler.CreateApplication)
		applicationGroup.GET("", a.handler.ListApplications)
		applicationGroup.GET("/:id", a.handler.GetApplication)
		applicationGroup.PUT("/:id", a.handler.UpdateApplication)
		applicationGroup.DELETE("/:id", a.handler.DeleteApplication)

		// 统计和批量操作
		applicationGroup.GET("/stats", a.handler.GetApplicationStats)
		applicationGroup.POST("/batch-delete", a.handler.BatchDeleteApplications)

		// 健康检查
		applicationGroup.GET("/health", a.handler.HealthCheck)
	}
}

// InitAPIServiceRoute 依赖注入版本的路由初始化
func (a *application) InitAPIServiceRoute(rg *gin.RouterGroup) {
	// 创建handler（注入后才能使用）
	if a.ApplicationService != nil {
		a.handler = handler.NewApplicationHandler(a.ApplicationService)
	}

	applicationGroup := rg.Group("/applications")
	{
		if a.handler != nil {
			// 应用CRUD操作
			applicationGroup.POST("", a.handler.CreateApplication)
			applicationGroup.GET("", a.handler.ListApplications)
			applicationGroup.GET("/:id", a.handler.GetApplication)
			applicationGroup.PUT("/:id", a.handler.UpdateApplication)
			applicationGroup.DELETE("/:id", a.handler.DeleteApplication)

			// 统计和批量操作
			applicationGroup.GET("/stats", a.handler.GetApplicationStats)
			applicationGroup.POST("/batch-delete", a.handler.BatchDeleteApplications)

			// 健康检查
			applicationGroup.GET("/health", a.handler.HealthCheck)
		}
	}
}
