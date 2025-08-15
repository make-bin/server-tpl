package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/api/middleware"
	"github.com/make-bin/server-tpl/pkg/utils/bcode"
	"github.com/make-bin/server-tpl/pkg/utils/errors"
)

type exampleAPI struct{}

func NewExampleAPI() APIInterface {
	return &exampleAPI{}
}

func (api *exampleAPI) InitAPIServiceRoute(rg *gin.RouterGroup) {
	examples := rg.Group("/examples")
	{
		examples.GET("/error/:type", api.DemonstrateError)
		examples.GET("/success", api.DemonstrateSuccess)
		examples.GET("/list", api.DemonstrateList)
		examples.POST("/create", api.DemonstrateCreate)
		examples.GET("/panic", api.DemonstratePanic)
	}
}

// DemonstrateError 演示不同类型的错误处理
func (api *exampleAPI) DemonstrateError(c *gin.Context) {
	errorType := c.Param("type")

	switch errorType {
	case "system":
		// 系统错误
		err := bcode.ErrSystemInternal.WithDetails("Database connection failed")
		middleware.ErrorResponse(c, err)

	case "user":
		// 用户错误
		err := bcode.ErrUserNotFound.WithDetails("User ID: 123")
		middleware.ErrorResponse(c, err)

	case "app":
		// 应用错误
		err := bcode.ErrAppNotFound.WithDetails("Application: my-app")
		middleware.ErrorResponse(c, err)

	case "db":
		// 数据库错误
		err := bcode.ErrDBConnectionFailed.WithDetails("PostgreSQL connection timeout")
		middleware.ErrorResponse(c, err)

	case "auth":
		// 认证错误
		err := bcode.ErrAuthFailed.WithDetails("Invalid credentials")
		middleware.ErrorResponse(c, err)

	case "file":
		// 文件错误
		err := bcode.ErrFileNotFound.WithDetails("File: config.yaml")
		middleware.ErrorResponse(c, err)

	case "external":
		// 外部服务错误
		err := bcode.ErrExternalServiceUnavailable.WithDetails("Payment service down")
		middleware.ErrorResponse(c, err)

	case "validation":
		// 验证错误
		err := errors.ErrValidation.WithDetails("Email format is invalid")
		middleware.ErrorResponse(c, err)

	case "custom":
		// 自定义错误
		customErr := errors.New(4001, "Custom business error")
		middleware.ErrorResponse(c, customErr)

	default:
		// 默认错误
		err := errors.ErrBadRequest.WithDetails("Unknown error type: " + errorType)
		middleware.ErrorResponse(c, err)
	}
}

// DemonstrateSuccess 演示成功响应
func (api *exampleAPI) DemonstrateSuccess(c *gin.Context) {
	data := gin.H{
		"id":      1,
		"name":    "Example Item",
		"status":  "active",
		"created": "2023-01-01T00:00:00Z",
	}

	middleware.SuccessResponse(c, data)
}

// DemonstrateList 演示列表响应
func (api *exampleAPI) DemonstrateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 模拟数据
	list := []gin.H{
		{"id": 1, "name": "Item 1", "status": "active"},
		{"id": 2, "name": "Item 2", "status": "inactive"},
		{"id": 3, "name": "Item 3", "status": "active"},
	}

	total := int64(100)

	middleware.ListResponse(c, list, total, page, size)
}

// DemonstrateCreate 演示创建响应
func (api *exampleAPI) DemonstrateCreate(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证错误
		validationErr := errors.ErrValidation.WithDetails(err.Error())
		middleware.ErrorResponse(c, validationErr)
		return
	}

	// 模拟业务逻辑
	if req.Name == "error" {
		// 模拟业务错误
		err := bcode.ErrAppAlreadyExists.WithDetails("Application name already exists")
		middleware.ErrorResponse(c, err)
		return
	}

	// 创建成功
	createdData := gin.H{
		"id":      123,
		"name":    req.Name,
		"status":  req.Status,
		"created": "2023-01-01T00:00:00Z",
	}

	middleware.CreatedResponse(c, createdData)
}

// DemonstratePanic 演示panic处理
func (api *exampleAPI) DemonstratePanic(c *gin.Context) {
	panicType := c.Query("type")

	switch panicType {
	case "string":
		panic("This is a string panic")
	case "error":
		panic(errors.New(5000, "This is an error panic"))
	case "nil":
		panic(nil)
	default:
		panic("Unknown panic type")
	}
}
