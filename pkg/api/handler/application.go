package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	v1 "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/api/response"
	"github.com/make-bin/server-tpl/pkg/api/validation"
	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/domain/service"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// ApplicationHandler 应用处理器
type ApplicationHandler struct {
	applicationService service.ApplicationServiceInterface
	validator          *validator.Validate
}

// NewApplicationHandler 创建应用处理器
func NewApplicationHandler(applicationService service.ApplicationServiceInterface) *ApplicationHandler {
	validator := validator.New()
	validation.RegisterCustomValidators(validator)

	return &ApplicationHandler{
		applicationService: applicationService,
		validator:          validator,
	}
}

// CreateApplication godoc
// @Summary 创建应用
// @Description 创建新的应用
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param request body v1.CreateApplicationRequest true "应用创建请求"
// @Success 201 {object} response.Response{data=v1.ApplicationResponse} "应用创建成功"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 409 {object} response.Response{error=string} "应用已存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications [post]
// @Security BearerAuth
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req v1.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			details := response.ParseValidationErrors(validationErrors)
			response.ValidationError(c, details)
		} else {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "validation_error", err)
		}
		return
	}

	// 转换为领域模型
	app := &model.Application{
		Name:        req.Name,
		Description: req.Description,
	}

	// 创建应用
	createdApp, err := h.applicationService.CreateApplication(c.Request.Context(), app)
	if err != nil {
		logger.Error("Failed to create application: %v", err)
		if errors.Is(err, model.ErrApplicationNotFound) {
			response.BusinessError(c, response.CodeAppNotFound, "app_exists", err)
		} else {
			response.InternalServerError(c, "internal_error", err)
		}
		return
	}

	// 转换响应
	resp := h.convertToApplicationResponse(createdApp)
	response.Created(c, resp, "app_created")
}

// GetApplication godoc
// @Summary 获取应用详情
// @Description 根据应用ID获取应用详细信息
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param id path int true "应用ID" minimum(1)
// @Success 200 {object} response.Response{data=v1.ApplicationResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 404 {object} response.Response{error=string} "应用不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications/{id} [get]
// @Security BearerAuth
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParameter, "invalid_parameter", err)
		return
	}

	app, err := h.applicationService.GetApplicationByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("Failed to get application: %v", err)
		if errors.Is(err, model.ErrApplicationNotFound) {
			response.NotFound(c, "app_not_found", err)
		} else {
			response.InternalServerError(c, "internal_error", err)
		}
		return
	}

	resp := h.convertToApplicationResponse(app)
	response.Success(c, resp)
}

// ListApplications godoc
// @Summary 获取应用列表
// @Description 分页获取应用列表
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) minimum(1)
// @Param size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param keyword query string false "搜索关键词" maxlength(100)
// @Param sort_by query string false "排序字段" example("created_at")
// @Param sort_desc query bool false "排序方向" default(true)
// @Param status query string false "应用状态" Enums(active, inactive, deleted)
// @Success 200 {object} response.Response{data=response.PaginationResponse{items=[]v1.ApplicationResponse}} "获取成功"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications [get]
// @Security BearerAuth
func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	var req v1.ListApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			details := response.ParseValidationErrors(validationErrors)
			response.ValidationError(c, details)
		} else {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "validation_error", err)
		}
		return
	}

	// 设置默认值
	req.PageRequest.Validate()

	// 调用服务
	apps, total, err := h.applicationService.ListApplications(c.Request.Context(), req.Page, req.Size)
	if err != nil {
		logger.Error("Failed to list applications: %v", err)
		response.InternalServerError(c, "internal_error", err)
		return
	}

	// 转换响应
	items := make([]v1.ApplicationResponse, len(apps))
	for i, app := range apps {
		items[i] = h.convertToApplicationResponse(app)
	}

	response.Page(c, items, req.Page, req.Size, int(total))
}

// UpdateApplication godoc
// @Summary 更新应用
// @Description 更新应用信息
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param id path int true "应用ID" minimum(1)
// @Param request body v1.UpdateApplicationRequest true "应用更新请求"
// @Success 200 {object} response.Response{data=v1.ApplicationResponse} "更新成功"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 404 {object} response.Response{error=string} "应用不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications/{id} [put]
// @Security BearerAuth
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParameter, "invalid_parameter", err)
		return
	}

	var req v1.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			details := response.ParseValidationErrors(validationErrors)
			response.ValidationError(c, details)
		} else {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "validation_error", err)
		}
		return
	}

	// 获取现有应用
	app, err := h.applicationService.GetApplicationByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("Failed to get application: %v", err)
		if errors.Is(err, model.ErrApplicationNotFound) {
			response.NotFound(c, "app_not_found", err)
		} else {
			response.InternalServerError(c, "internal_error", err)
		}
		return
	}

	// 更新字段
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Description != "" {
		app.Description = req.Description
	}

	// 更新应用
	updatedApp, err := h.applicationService.UpdateApplication(c.Request.Context(), app)
	if err != nil {
		logger.Error("Failed to update application: %v", err)
		response.InternalServerError(c, "internal_error", err)
		return
	}

	resp := h.convertToApplicationResponse(updatedApp)
	response.WithMessage(c, resp, "app_updated")
}

// DeleteApplication godoc
// @Summary 删除应用
// @Description 删除指定的应用
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param id path int true "应用ID" minimum(1)
// @Success 204 "删除成功"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 404 {object} response.Response{error=string} "应用不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications/{id} [delete]
// @Security BearerAuth
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParameter, "invalid_parameter", err)
		return
	}

	err = h.applicationService.DeleteApplication(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("Failed to delete application: %v", err)
		if errors.Is(err, model.ErrApplicationNotFound) {
			response.NotFound(c, "app_not_found", err)
		} else {
			response.InternalServerError(c, "internal_error", err)
		}
		return
	}

	response.NoContent(c)
}

// GetApplicationStats godoc
// @Summary 获取应用统计
// @Description 获取应用统计信息
// @Tags 应用管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=v1.ApplicationStatsResponse} "获取成功"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications/stats [get]
// @Security BearerAuth
func (h *ApplicationHandler) GetApplicationStats(c *gin.Context) {
	// 这里应该调用统计服务，暂时返回模拟数据
	stats := v1.ApplicationStatsResponse{
		TotalApps:    150,
		ActiveApps:   120,
		InactiveApps: 20,
		DeletedApps:  10,
		TodayNewApps: 5,
		MonthNewApps: 25,
	}

	response.Success(c, stats)
}

// BatchDeleteApplications godoc
// @Summary 批量删除应用
// @Description 批量删除多个应用
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param request body v1.BatchDeleteApplicationsRequest true "批量删除请求"
// @Success 200 {object} response.Response{data=v1.BulkOperationResponse} "操作完成"
// @Failure 400 {object} response.Response{error=string} "参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /applications/batch-delete [post]
// @Security BearerAuth
func (h *ApplicationHandler) BatchDeleteApplications(c *gin.Context) {
	var req v1.BatchDeleteApplicationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			details := response.ParseValidationErrors(validationErrors)
			response.ValidationError(c, details)
		} else {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "validation_error", err)
		}
		return
	}

	var failures []v1.BulkFailureItem
	successCount := 0

	for _, id := range req.IDs {
		err := h.applicationService.DeleteApplication(c.Request.Context(), id)
		if err != nil {
			failures = append(failures, v1.BulkFailureItem{
				ID:     strconv.FormatUint(uint64(id), 10),
				Reason: err.Error(),
			})
		} else {
			successCount++
		}
	}

	result := v1.BulkOperationResponse{
		SuccessCount: successCount,
		FailureCount: len(failures),
		TotalCount:   len(req.IDs),
		Failures:     failures,
	}

	response.Success(c, result)
}

// HealthCheck godoc
// @Summary 健康检查
// @Description 检查应用服务健康状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=v1.HealthCheckResponse} "服务正常"
// @Failure 500 {object} response.Response{error=string} "服务异常"
// @Router /applications/health [get]
func (h *ApplicationHandler) HealthCheck(c *gin.Context) {
	healthResp := v1.HealthCheckResponse{
		Status:  "ok",
		Message: "应用服务运行正常",
		Version: "1.0.0",
	}

	response.Success(c, healthResp)
}

// convertToApplicationResponse 转换为应用响应
func (h *ApplicationHandler) convertToApplicationResponse(app *model.Application) v1.ApplicationResponse {
	return v1.ApplicationResponse{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		Status:      "active", // 这里应该从模型中获取状态
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}
