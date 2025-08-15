# 错误处理指南

本文档介绍了项目中错误处理的使用方法和最佳实践。

## 错误处理架构

项目采用分层的错误处理架构：

1. **错误定义层** (`pkg/utils/errors/`) - 定义错误类型和基础错误
2. **业务错误码层** (`pkg/utils/bcode/`) - 定义业务错误码
3. **中间件层** (`pkg/infrastructure/middleware/`) - 统一错误处理和响应
4. **API层** - 使用错误处理工具

## 错误类型

### 1. 基础错误 (pkg/utils/errors/)

```go
import "github.com/make-bin/server-tpl/pkg/utils/errors"

// 创建新错误
err := errors.New(400, "Bad request")

// 包装现有错误
err = errors.Wrap(originalErr, 500, "Internal server error")

// 添加详细信息
err = err.WithDetails("Additional context")

// 添加原始错误
err = err.WithCause(originalErr)
```

### 2. 业务错误码 (pkg/utils/bcode/)

```go
import "github.com/make-bin/server-tpl/pkg/utils/bcode"

// 使用预定义错误码
err := bcode.ErrUserNotFound.WithDetails("User ID: 123")
err := bcode.ErrAppAlreadyExists.WithDetails("App name: my-app")
err := bcode.ErrDBConnectionFailed.WithDetails("Connection timeout")

// 创建自定义业务错误码
customErr := bcode.New(4001, "Custom business error", "CUSTOM")
```

## 错误码范围

| 范围 | 模块 | 说明 |
|------|------|------|
| 1000-1999 | 系统错误 | 系统级错误，如配置、网络等 |
| 2000-2999 | 用户错误 | 用户相关错误，如认证、权限等 |
| 3000-3999 | 应用错误 | 应用业务逻辑错误 |
| 4000-4999 | 配置错误 | 配置相关错误 |
| 5000-5999 | 数据库错误 | 数据库操作错误 |
| 6000-6999 | 认证授权错误 | 认证和授权相关错误 |
| 7000-7999 | 文件操作错误 | 文件读写等操作错误 |
| 8000-8999 | 外部服务错误 | 外部服务调用错误 |

## 在API中使用错误处理

### 1. 基本用法

```go
func (api *userAPI) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := api.userService.GetByID(c.Request.Context(), userID)
    if err != nil {
        // 使用中间件统一处理错误
        middleware.ErrorResponse(c, err)
        return
    }
    
    // 成功响应
    middleware.SuccessResponse(c, user)
}
```

### 2. 业务错误处理

```go
func (api *appAPI) CreateApp(c *gin.Context) {
    var req CreateAppRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 参数验证错误
        validationErr := errors.ErrValidation.WithDetails(err.Error())
        middleware.ErrorResponse(c, validationErr)
        return
    }
    
    // 检查应用名是否已存在
    if exists, _ := api.appService.ExistsByName(req.Name); exists {
        err := bcode.ErrAppAlreadyExists.WithDetails("App name: " + req.Name)
        middleware.ErrorResponse(c, err)
        return
    }
    
    app, err := api.appService.Create(c.Request.Context(), &req)
    if err != nil {
        middleware.ErrorResponse(c, err)
        return
    }
    
    middleware.CreatedResponse(c, app)
}
```

### 3. 数据库错误处理

```go
func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
    var user model.User
    err := s.store.GetByID(ctx, id, &user)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, bcode.ErrUserNotFound.WithDetails("User ID: " + id)
        }
        return nil, bcode.ErrDBQueryFailed.WithCause(err)
    }
    return &user, nil
}
```

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "example",
    "status": "active"
  }
}
```

### 错误响应

```json
{
  "code": 3000,
  "message": "Application not found",
  "details": "Application: my-app"
}
```

### 列表响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {"id": 1, "name": "item1"},
      {"id": 2, "name": "item2"}
    ],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

## 中间件响应函数

### 1. 成功响应

```go
// 普通成功响应
middleware.SuccessResponse(c, data)

// 创建成功响应
middleware.CreatedResponse(c, data)

// 无内容响应
middleware.NoContentResponse(c)

// 列表响应
middleware.ListResponse(c, list, total, page, size)
```

### 2. 错误响应

```go
// 统一错误响应
middleware.ErrorResponse(c, err)
```

## 错误处理最佳实践

### 1. 错误分类

- **系统错误** (1000-1999): 系统级问题，需要运维处理
- **用户错误** (2000-2999): 用户操作问题，可以提示用户
- **业务错误** (3000-3999): 业务逻辑问题，需要业务处理
- **技术错误** (5000-8999): 技术实现问题，需要开发处理

### 2. 错误信息

- 提供清晰的错误描述
- 包含必要的上下文信息
- 避免暴露敏感信息
- 使用统一的错误格式

### 3. 错误日志

- 记录错误的详细信息
- 包含请求上下文
- 记录堆栈跟踪
- 区分错误级别

### 4. 错误恢复

- 提供降级方案
- 实现重试机制
- 设置超时控制
- 监控错误率

## 示例API

项目提供了示例API来演示错误处理的使用：

### 错误演示

```bash
# 系统错误
GET /api/v1/examples/error/system

# 用户错误
GET /api/v1/examples/error/user

# 应用错误
GET /api/v1/examples/error/app

# 数据库错误
GET /api/v1/examples/error/db

# 认证错误
GET /api/v1/examples/error/auth

# 文件错误
GET /api/v1/examples/error/file

# 外部服务错误
GET /api/v1/examples/error/external

# 验证错误
GET /api/v1/examples/error/validation

# 自定义错误
GET /api/v1/examples/error/custom
```

### 成功响应演示

```bash
# 成功响应
GET /api/v1/examples/success

# 列表响应
GET /api/v1/examples/list?page=1&size=10

# 创建响应
POST /api/v1/examples/create
Content-Type: application/json

{
  "name": "test-app",
  "status": "active"
}
```

### Panic处理演示

```bash
# 字符串panic
GET /api/v1/examples/panic?type=string

# 错误panic
GET /api/v1/examples/panic?type=error

# 空值panic
GET /api/v1/examples/panic?type=nil
```

## 自定义错误码

### 1. 添加新的错误码

```go
// 在 pkg/utils/bcode/bcode.go 中添加
var (
    // 自定义模块错误码 (9000-9999)
    ErrCustomBusinessError = New(9000, "Custom business error", "CUSTOM")
    ErrCustomValidation   = New(9001, "Custom validation error", "CUSTOM")
)
```

### 2. 在业务中使用

```go
func (s *customService) Process() error {
    // 业务逻辑
    if someCondition {
        return bcode.ErrCustomBusinessError.WithDetails("Specific error details")
    }
    
    if validationFailed {
        return bcode.ErrCustomValidation.WithDetails("Validation failed")
    }
    
    return nil
}
```

## 错误监控

### 1. 错误统计

- 按错误码统计错误率
- 按模块统计错误分布
- 按时间统计错误趋势
- 设置错误告警阈值

### 2. 错误追踪

- 记录错误发生的完整链路
- 关联用户请求和错误
- 分析错误根因
- 优化错误处理逻辑

## 总结

通过统一的错误处理机制，项目实现了：

1. **标准化**: 统一的错误格式和响应结构
2. **分类化**: 按模块和类型分类错误
3. **可追踪**: 完整的错误信息和堆栈跟踪
4. **可监控**: 错误统计和告警机制
5. **易维护**: 清晰的错误处理代码结构

这种错误处理方式提高了系统的可维护性和用户体验，为业务发展提供了可靠的技术基础。
