package api

import (
	"net/http"
	"regexp"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	assemblerv1 "github.com/make-bin/server-tpl/pkg/api/assembler/v1"
	dto "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	v1 "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/domain/service"
)

func init() {
	// 移除重复注册，已在 interface.go 中注册
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

type application struct {
	ApplicationService service.ApplicationService `inject:""`
	VariablesService   service.VariablesService   `inject:""`
}

func newApplication() APIInterface {
	return &application{}
}

func (a *application) InitAPIServiceRoute(rg *gin.RouterGroup) {
	// 应用相关路由
	appGroup := rg.Group("/applications")
	{
		appGroup.POST("", a.createApplication)
		appGroup.GET("", a.listApplications)
		appGroup.GET("/:id", a.getApplication)
		appGroup.PUT("/:id", a.updateApplication)
		appGroup.DELETE("/:id", a.deleteApplication)
	}

	// 变量相关路由
	varGroup := rg.Group("/variables")
	{
		varGroup.POST("", a.createVariable)
		varGroup.GET("", a.listVariables)
		varGroup.GET("/:id", a.getVariable)
		varGroup.PUT("/:id", a.updateVariable)
		varGroup.DELETE("/:id", a.deleteVariable)
		varGroup.GET("/app/:appId", a.getVariablesByAppID)
	}
}

// 应用相关API处理函数
func (a *application) createApplication(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := assemblerv1.ToApplicationModel(&req)
	if err := a.ApplicationService.CreateApplication(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, assemblerv1.ToApplicationResponse(app))
}

func (a *application) listApplications(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	apps, err := a.ApplicationService.ListApplications(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToApplicationListResponse(apps))
}

func (a *application) getApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	app, err := a.ApplicationService.GetApplicationByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToApplicationResponse(app))
}

func (a *application) updateApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := assemblerv1.ToApplicationModel(&req)
	app.BaseEntity.ID = uint(id)

	if err := a.ApplicationService.UpdateApplication(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToApplicationResponse(app))
}

func (a *application) deleteApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := a.ApplicationService.DeleteApplication(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

// 变量相关API处理函数
func (a *application) createVariable(c *gin.Context) {
	var req v1.CreateVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variable := assemblerv1.ToVariableModel(&req)
	if err := a.VariablesService.CreateVariable(c.Request.Context(), variable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, assemblerv1.ToVariableResponse(variable))
}

func (a *application) listVariables(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	variables, err := a.VariablesService.ListVariables(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToVariableListResponse(variables))
}

func (a *application) getVariable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	variable, err := a.VariablesService.GetVariableByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variable not found"})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToVariableResponse(variable))
}

func (a *application) updateVariable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req v1.UpdateVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variable := assemblerv1.ToVariableModel(&req)
	variable.ID = uint(id)

	if err := a.VariablesService.UpdateVariable(c.Request.Context(), variable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToVariableResponse(variable))
}

func (a *application) deleteVariable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := a.VariablesService.DeleteVariable(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variable deleted successfully"})
}

func (a *application) getVariablesByAppID(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("appId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid App ID"})
		return
	}

	variables, err := a.VariablesService.GetVariablesByAppID(c.Request.Context(), uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblerv1.ToVariableListResponse(variables))
}
