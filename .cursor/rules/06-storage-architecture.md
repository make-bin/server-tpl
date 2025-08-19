# 6. 存储架构设计

## 6.1 抽象接口层
项目采用统一的存储抽象接口，支持多种存储后端：

```go
// Entity 数据库实体接口
type Entity interface {
    SetCreateTime(time.Time)
    SetUpdateTime(time.Time)
    PrimaryKey() string
    TableName() string
    ShortTableName() string
    Index() map[string]interface{}
}

// DataStore 数据存储接口
type DataStore interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    HealthCheck(ctx context.Context) error
    BeginTx(ctx context.Context) (Transaction, error)
    
    Add(ctx context.Context, entity Entity) error
    BatchAdd(ctx context.Context, entities []Entity) error
    Put(ctx context.Context, entity Entity) error
    Delete(ctx context.Context, entity Entity) error
    Get(ctx context.Context, entity Entity) error
    List(ctx context.Context, query Entity, options *ListOptions) ([]Entity, error)
    Count(ctx context.Context, entity Entity, options *FilterOptions) (int64, error)
    IsExist(ctx context.Context, entity Entity) (bool, error)
}
```

## 6.2 支持的存储后端
1. **PostgreSQL**: 生产环境主数据库
2. **OpenGauss**: 兼容 PostgreSQL 的国产数据库
3. **Memory**: 内存存储，用于测试和开发

## 6.3 工厂模式
使用工厂模式创建存储实例：

```go
// 创建存储实例
config := &datastore.Config{
    Type:     "memory", // 或 "postgresql", "opengauss"
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "password",
    Database: "server_tpl",
}

store, err := factory.NewDataStore(config)
```
