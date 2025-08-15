# 开发指南

本文档提供了使用Go HTTP Server Template进行开发的详细指南。

## 项目架构

项目采用分层架构设计，主要包含以下几个层次：

### 1. API层 (pkg/api/)
- 处理HTTP请求和响应
- 参数验证和错误处理
- 路由定义和管理
- 使用DTO进行数据传输

### 2. 领域层 (pkg/domain/)
- 定义业务模型和实体
- 实现核心业务逻辑
- 定义服务接口
- 业务规则和约束

### 3. 基础设施层 (pkg/infrastructure/)
- 数据持久化实现
- 外部服务集成
- 中间件实现
- 配置管理

### 4. 工具层 (pkg/utils/)
- 通用工具函数
- 依赖注入容器
- 辅助功能

## 开发流程

### 1. 添加新的业务功能

#### 步骤1: 定义数据模型
在 `pkg/domain/model/` 中定义数据模型：

```go
// pkg/domain/model/user.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    Email     string         `gorm:"size:100;unique;not null" json:"email"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### 步骤2: 定义服务接口
在 `pkg/domain/service/` 中定义服务接口：

```go
// pkg/domain/service/user.go
package service

import (
    "context"
    "github.com/make-bin/server-tpl/pkg/domain/model"
)

type UserService interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uint) (*model.User, error)
    List(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
}
```

#### 步骤3: 实现服务
在 `pkg/domain/service/` 中实现服务：

```go
// pkg/domain/service/user_impl.go
package service

import (
    "context"
    "github.com/make-bin/server-tpl/pkg/domain/model"
    "github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

type userService struct {
    store datastore.DataStore
}

func NewUserService(store datastore.DataStore) UserService {
    return &userService{store: store}
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
    return s.store.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := s.store.GetByID(ctx, id, &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ... 其他方法实现
```

#### 步骤4: 定义DTO
在 `pkg/api/dto/v1/` 中定义数据传输对象：

```go
// pkg/api/dto/v1/user.go
package v1

import (
    "time"
)

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required" validate:"required"`
    Email string `json:"email" binding:"required,email" validate:"required,email"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
    Users []UserResponse `json:"users"`
    Total int64         `json:"total"`
}
```

#### 步骤5: 实现转换器
在 `pkg/api/assembler/v1/` 中实现DTO与模型的转换：

```go
// pkg/api/assembler/v1/user.go
package v1

import (
    "github.com/make-bin/server-tpl/pkg/api/dto/v1"
    "github.com/make-bin/server-tpl/pkg/domain/model"
)

func ToUserModel(req *v1.CreateUserRequest) *model.User {
    return &model.User{
        Name:  req.Name,
        Email: req.Email,
    }
}

func ToUserResponse(user *model.User) *v1.UserResponse {
    return &v1.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}

func ToUserListResponse(users []*model.User, total int64) *v1.ListUsersResponse {
    responses := make([]v1.UserResponse, len(users))
    for i, user := range users {
        responses[i] = *ToUserResponse(user)
    }
    
    return &v1.ListUsersResponse{
        Users: responses,
        Total: total,
    }
}
```

#### 步骤6: 实现API处理器
在 `pkg/api/` 中实现API处理器：

```go
// pkg/api/user.go
package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/make-bin/server-tpl/pkg/api/dto/v1"
    "github.com/make-bin/server-tpl/pkg/api/assembler/v1"
    "github.com/make-bin/server-tpl/pkg/domain/service"
)

type userAPI struct {
    userService service.UserService
}

func NewUserAPI(userService service.UserService) APIInterface {
    return &userAPI{userService: userService}
}

func (api *userAPI) InitAPIServiceRoute(rg *gin.RouterGroup) {
    users := rg.Group("/users")
    {
        users.POST("", api.Create)
        users.GET("", api.List)
        users.GET("/:id", api.GetByID)
        users.PUT("/:id", api.Update)
        users.DELETE("/:id", api.Delete)
    }
}

func (api *userAPI) Create(c *gin.Context) {
    var req v1.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := assembler.ToUserModel(&req)
    if err := api.userService.Create(c.Request.Context(), user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, assembler.ToUserResponse(user))
}

func (api *userAPI) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    user, err := api.userService.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, assembler.ToUserResponse(user))
}

// ... 其他方法实现
```

#### 步骤7: 注册API
在 `pkg/api/interface.go` 中注册新的API：

```go
func init() {
    RegisterAPIInterface(newApplication())
    RegisterAPIInterface(newUser()) // 添加新的API
}
```

### 2. 配置管理

#### 环境变量配置
```bash
export SERVER_PORT=8080
export SERVER_HOST=0.0.0.0
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=password
export DATABASE_DBNAME=server_tpl
```

#### 配置文件
在 `configs/config.yaml` 中配置：

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  mode: "debug"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "server_tpl"
  sslmode: "disable"

log:
  level: "debug"
  format: "text"
  output: "stdout"
```

### 3. 依赖注入

#### 注册服务
在 `pkg/domain/service/interface.go` 中注册服务：

```go
func InitServiceBean() []interface{} {
    return []interface{}{
        NewApplicationService(),
        NewUserService(), // 添加新的服务
    }
}
```

#### 注册API
在 `pkg/api/interface.go` 中注册API：

```go
func init() {
    RegisterAPIInterface(newApplication())
    RegisterAPIInterface(newUser()) // 添加新的API
}
```

### 4. 测试

#### 单元测试
```go
// pkg/domain/service/user_test.go
package service

import (
    "testing"
    "context"
    "github.com/make-bin/server-tpl/pkg/domain/model"
)

func TestUserService_Create(t *testing.T) {
    // 测试实现
}
```

#### 集成测试
```go
// pkg/e2e/user_test.go
package e2e

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
)

func TestUserAPI_Create(t *testing.T) {
    // 集成测试实现
}
```

### 5. 日志记录

#### 使用日志
```go
import "github.com/make-bin/server-tpl/pkg/utils/logger"

// 基本日志
logger.Info("Application started")
logger.Error("An error occurred")

// 带字段的日志
logger.WithField("user_id", 123).Info("User logged in")
logger.WithFields(logrus.Fields{
    "user_id": 123,
    "action":  "login",
}).Info("User action")
```

### 6. 错误处理

#### 定义错误类型
```go
// pkg/domain/errors/errors.go
package errors

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

#### 使用错误
```go
func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    if err := s.store.GetByID(ctx, id, &user); err != nil {
        return nil, errors.Wrap(ErrUserNotFound, err.Error())
    }
    return &user, nil
}
```

## 最佳实践

### 1. 代码组织
- 保持包的结构清晰
- 使用有意义的包名
- 避免循环依赖

### 2. 错误处理
- 始终检查错误
- 提供有意义的错误信息
- 使用错误包装

### 3. 日志记录
- 记录关键操作
- 使用适当的日志级别
- 包含上下文信息

### 4. 测试
- 编写单元测试
- 测试边界条件
- 保持测试覆盖率

### 5. 性能
- 使用连接池
- 避免N+1查询
- 使用适当的索引

### 6. 安全
- 验证输入
- 使用HTTPS
- 实现认证和授权

## 常见问题

### Q: 如何处理数据库迁移？
A: 使用GORM的AutoMigrate功能，在应用启动时自动迁移表结构。

### Q: 如何实现分页？
A: 在服务层实现分页逻辑，使用offset和limit参数。

### Q: 如何处理并发？
A: 使用适当的锁机制，避免竞态条件。

### Q: 如何实现缓存？
A: 在服务层实现缓存逻辑，使用Redis或其他缓存系统。

## 更多资源

- [Go官方文档](https://golang.org/doc/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [GORM文档](https://gorm.io/docs/)
- [Viper文档](https://github.com/spf13/viper)
- [Logrus文档](https://github.com/sirupsen/logrus)
