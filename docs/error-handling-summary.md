# 错误处理功能总结

## 功能概述

本次为项目添加了完整的错误处理范式，包括：

1. **统一错误处理工具** (`pkg/utils/errors/`)
2. **业务错误码系统** (`pkg/utils/bcode/`)
3. **错误处理中间件** (`pkg/infrastructure/middleware/`)
4. **示例API** (`pkg/api/example.go`)
5. **完整文档** (`docs/error-handling.md`)

## 新增文件结构

```
pkg/
├── utils/
│   ├── errors/
│   │   ├── errors.go      # 错误处理工具
│   │   └── errors_test.go # 错误处理测试
│   └── bcode/
│       ├── bcode.go       # 业务错误码定义
│       └── bcode_test.go  # 业务错误码测试
├── infrastructure/
│   └── middleware/
│       └── error_handler.go # 错误处理中间件
└── api/
    └── example.go         # 示例API
```

## 核心功能

### 1. 错误处理工具 (pkg/utils/errors/)

#### 主要特性
- **自定义错误类型**: 支持错误码、消息、详情、堆栈跟踪
- **错误包装**: 支持包装现有错误，保留原始错误信息
- **错误链**: 支持错误原因链，便于调试
- **统一响应**: 转换为统一的错误响应格式

#### 核心方法
```go
// 创建新错误
err := errors.New(400, "Bad request")

// 包装错误
err = errors.Wrap(originalErr, 500, "Internal error")

// 添加详情
err = err.WithDetails("Additional context")

// 添加原因
err = err.WithCause(originalErr)

// 转换为响应
resp := err.ToResponse()
```

### 2. 业务错误码系统 (pkg/utils/bcode/)

#### 错误码范围
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

#### 预定义错误码
```go
// 系统错误
ErrSystemInternal     // 1000: 系统内部错误
ErrSystemUnavailable  // 1001: 系统暂时不可用
ErrSystemMaintenance  // 1002: 系统维护中

// 用户错误
ErrUserNotFound       // 2000: 用户不存在
ErrUserAlreadyExists  // 2001: 用户已存在
ErrUserPasswordWrong  // 2002: 密码错误

// 应用错误
ErrAppNotFound        // 3000: 应用不存在
ErrAppAlreadyExists   // 3001: 应用已存在
ErrAppCreateFailed    // 3100: 应用创建失败

// 数据库错误
ErrDBConnectionFailed // 5000: 数据库连接失败
ErrDBQueryFailed      // 5100: 数据库查询失败
ErrDBTransactionFailed // 5200: 数据库事务失败

// 认证错误
ErrAuthFailed         // 6000: 认证失败
ErrAuthExpired        // 6001: 认证过期
ErrAuthPermissionDenied // 6100: 权限被拒绝
```

#### 使用方法
```go
// 使用预定义错误码
err := bcode.ErrUserNotFound.WithDetails("User ID: 123")
err := bcode.ErrAppAlreadyExists.WithDetails("App name: my-app")

// 创建自定义错误码
customErr := bcode.New(9000, "Custom business error", "CUSTOM")
```

### 3. 错误处理中间件 (pkg/infrastructure/middleware/)

#### 主要功能
- **统一错误响应**: 自动转换错误为HTTP响应
- **Panic恢复**: 自动恢复panic并返回错误响应
- **错误日志**: 自动记录错误日志
- **状态码映射**: 根据错误码自动设置HTTP状态码

#### 响应函数
```go
// 错误响应
middleware.ErrorResponse(c, err)

// 成功响应
middleware.SuccessResponse(c, data)
middleware.CreatedResponse(c, data)
middleware.NoContentResponse(c)
middleware.ListResponse(c, list, total, page, size)
```

#### 响应格式
```json
// 成功响应
{
  "code": 200,
  "message": "success",
  "data": {...}
}

// 错误响应
{
  "code": 3000,
  "message": "Application not found",
  "details": "Application: my-app"
}

// 列表响应
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

### 4. 示例API (pkg/api/example.go)

#### 演示功能
- **错误类型演示**: 展示不同类型的错误处理
- **成功响应演示**: 展示成功响应的格式
- **列表响应演示**: 展示分页列表响应
- **创建响应演示**: 展示创建操作的响应
- **Panic处理演示**: 展示panic的自动恢复

#### API端点
```bash
# 错误演示
GET /api/v1/examples/error/system
GET /api/v1/examples/error/user
GET /api/v1/examples/error/app
GET /api/v1/examples/error/db
GET /api/v1/examples/error/auth
GET /api/v1/examples/error/file
GET /api/v1/examples/error/external
GET /api/v1/examples/error/validation
GET /api/v1/examples/error/custom

# 成功响应演示
GET /api/v1/examples/success
GET /api/v1/examples/list?page=1&size=10
POST /api/v1/examples/create

# Panic处理演示
GET /api/v1/examples/panic?type=string
GET /api/v1/examples/panic?type=error
GET /api/v1/examples/panic?type=nil
```

## 测试覆盖

### 1. 错误处理工具测试
- ✅ 错误创建和包装
- ✅ 错误详情和原因
- ✅ 错误比较和转换
- ✅ 响应格式转换

### 2. 业务错误码测试
- ✅ 错误码创建
- ✅ 预定义错误码验证
- ✅ 错误信息获取
- ✅ 错误转换

### 3. 测试结果
```bash
$ go test ./pkg/utils/errors/ -v
PASS
ok      github.com/make-bin/server-tpl/pkg/utils/errors  0.002s

$ go test ./pkg/utils/bcode/ -v
PASS
ok      github.com/make-bin/server-tpl/pkg/utils/bcode   0.004s
```

## 使用示例

### 1. 在API中使用

```go
func (api *userAPI) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := api.userService.GetByID(c.Request.Context(), userID)
    if err != nil {
        middleware.ErrorResponse(c, err)
        return
    }
    
    middleware.SuccessResponse(c, user)
}
```

### 2. 在服务中使用

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

### 3. 参数验证

```go
func (api *appAPI) CreateApp(c *gin.Context) {
    var req CreateAppRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErr := errors.ErrValidation.WithDetails(err.Error())
        middleware.ErrorResponse(c, validationErr)
        return
    }
    
    // 业务逻辑...
}
```

## 最佳实践

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

## 优势特点

### 1. 标准化
- 统一的错误格式和响应结构
- 一致的错误处理流程
- 标准化的错误码定义

### 2. 分类化
- 按模块和类型分类错误
- 清晰的错误码范围
- 便于错误统计和分析

### 3. 可追踪
- 完整的错误信息和堆栈跟踪
- 错误原因链支持
- 便于调试和问题定位

### 4. 可监控
- 错误统计和告警机制
- 错误率监控
- 性能影响分析

### 5. 易维护
- 清晰的错误处理代码结构
- 模块化的错误定义
- 易于扩展和修改

## 总结

通过本次错误处理功能的添加，项目实现了：

1. **完整的错误处理体系**: 从错误定义到响应处理的完整链路
2. **业务错误码系统**: 标准化的业务错误码定义和管理
3. **统一的错误响应**: 一致的API错误响应格式
4. **自动错误恢复**: 自动处理panic和异常情况
5. **完善的测试覆盖**: 确保错误处理功能的可靠性
6. **详细的使用文档**: 便于开发人员理解和使用

这个错误处理系统为项目提供了可靠的技术基础，提高了系统的可维护性和用户体验，为业务发展提供了强有力的支撑。
