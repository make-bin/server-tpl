# 数据存储架构

本项目实现了基于抽象接口的数据存储架构，支持多种数据库后端。

## 架构设计

### 核心接口

#### Entity 接口
所有数据模型都必须实现 `Entity` 接口：

```go
type Entity interface {
    SetCreateTime(time.Time)
    SetUpdateTime(time.Time)
    PrimaryKey() string
    TableName() string
    ShortTableName() string
    Index() map[string]interface{}
}
```

#### DataStore 接口
数据存储的核心接口，提供统一的 CRUD 操作：

```go
type DataStore interface {
    Add(ctx context.Context, entity Entity) error
    BatchAdd(ctx context.Context, entities []Entity) error
    Put(ctx context.Context, entity Entity) error
    Delete(ctx context.Context, entity Entity) error
    Get(ctx context.Context, entity Entity) error
    List(ctx context.Context, query Entity, options *ListOptions) ([]Entity, error)
    Count(ctx context.Context, entity Entity, options *FilterOptions) (int64, error)
    IsExist(ctx context.Context, entity Entity) (bool, error)
    
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    HealthCheck(ctx context.Context) error
    BeginTx(ctx context.Context) (Transaction, error)
}
```

### 数据模型

所有数据模型都继承自 `BaseEntity`：

```go
type BaseEntity struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 示例模型

```go
type User struct {
    BaseEntity
    Email    string `json:"email" gorm:"uniqueIndex;not null"`
    Name     string `json:"name" gorm:"not null"`
    Password string `json:"-" gorm:"not null"`
    Role     string `json:"role" gorm:"default:'user'"`
    Status   string `json:"status" gorm:"default:'active'"`
}

func (User) TableName() string {
    return "users"
}

func (User) ShortTableName() string {
    return "user"
}

func (u *User) Index() map[string]interface{} {
    index := u.BaseEntity.Index()
    index["email"] = u.Email
    index["name"] = u.Name
    index["role"] = u.Role
    index["status"] = u.Status
    return index
}
```

## 支持的存储后端

### 1. PostgreSQL
- 位置：`pkg/infrastructure/datastore/postgresql/`
- 使用 GORM + PostgreSQL 驱动
- 支持自动迁移表结构
- 支持连接池配置

### 2. OpenGauss
- 位置：`pkg/infrastructure/datastore/opengauss/`
- 使用 GORM + PostgreSQL 驱动（兼容 OpenGauss）
- 支持自动迁移表结构
- 支持连接池配置

### 3. Memory
- 位置：`pkg/infrastructure/datastore/memory/`
- 内存存储，用于测试和开发
- 支持所有 CRUD 操作
- 支持过滤、排序和分页

## 使用方法

### 1. 创建数据存储实例

```go
import "github.com/make-bin/server-tpl/pkg/infrastructure/datastore/factory"

config := &datastore.Config{
    Type:     "postgresql", // 或 "opengauss", "memory"
    Host:     "localhost",
    Port:     5432,
    User:     "username",
    Password: "password",
    Database: "dbname",
    SSLMode:  "disable",
    MaxIdle:  10,
    MaxOpen:  100,
    Timeout:  30,
}

store, err := factory.NewDataStore(config)
if err != nil {
    log.Fatal(err)
}
```

### 2. 连接数据库

```go
ctx := context.Background()
if err := store.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer store.Disconnect(ctx)
```

### 3. 基本 CRUD 操作

```go
// 创建
user := &model.User{
    Email:    "test@example.com",
    Name:     "Test User",
    Password: "password123",
}
err := store.Add(ctx, user)

// 查询
queryUser := &model.User{}
queryUser.BaseEntity.ID = 1
err = store.Get(ctx, queryUser)

// 更新
queryUser.Name = "Updated Name"
err = store.Put(ctx, queryUser)

// 删除
err = store.Delete(ctx, queryUser)
```

### 4. 列表查询

```go
// 基本列表
options := &datastore.ListOptions{
    Page:     1,
    PageSize: 10,
    SortBy: []datastore.SortOption{
        {Key: "id", Order: datastore.SortOrderAscending},
    },
}
users, err := store.List(ctx, &model.User{}, options)

// 带过滤的列表
options := &datastore.ListOptions{
    FilterOptions: datastore.FilterOptions{
        Queries: []datastore.FuzzyQueryOption{
            {Key: "name", Query: "test"},
        },
        In: []datastore.InQueryOption{
            {Key: "status", Values: []string{"active"}},
        },
    },
    Page:     1,
    PageSize: 10,
}
users, err := store.List(ctx, &model.User{}, options)
```

### 5. 统计查询

```go
count, err := store.Count(ctx, &model.User{}, nil)
```

## 查询选项

### FilterOptions
- `Queries`: 模糊查询
- `In`: IN 查询
- `IsNotExist`: 不存在查询

### ListOptions
- `FilterOptions`: 过滤条件
- `Page`: 页码（从1开始）
- `PageSize`: 每页大小
- `SortBy`: 排序条件

### SortOption
- `Key`: 排序字段
- `Order`: 排序方向（`SortOrderAscending` 或 `SortOrderDescending`）

## 错误处理

所有存储操作都返回统一的错误类型：

- `ErrNilEntity`: 实体为空
- `ErrTableNameEmpty`: 表名为空
- `ErrPrimaryEmpty`: 主键为空
- `ErrRecordExist`: 记录已存在
- `ErrRecordNotExist`: 记录不存在
- `ErrIndexInvalid`: 索引无效
- `ErrEntityInvalid`: 实体无效

## 事务支持

```go
tx, err := store.BeginTx(ctx)
if err != nil {
    return err
}
defer tx.Rollback()

// 在事务中执行操作
err = tx.GetDataStore().Add(ctx, entity)
if err != nil {
    return err
}

return tx.Commit()
```

## 扩展新的存储后端

要实现新的存储后端，需要：

1. 实现 `DataStore` 接口
2. 在 `factory` 包中注册新的类型
3. 实现 `Transaction` 接口（如果需要事务支持）

示例：

```go
type MyDataStore struct {
    // 实现细节
}

func (m *MyDataStore) Add(ctx context.Context, entity datastore.Entity) error {
    // 实现添加逻辑
}

// 实现其他接口方法...

func NewMyDataStore(config *datastore.Config) *MyDataStore {
    return &MyDataStore{}
}
```

然后在 factory 中注册：

```go
case DatabaseTypeMyDB:
    return mydb.NewMyDataStore(config), nil
```
