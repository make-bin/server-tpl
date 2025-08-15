# 中间件重构总结

## 重构概述

本次重构将项目中的中间件进行了重新组织，实现了更清晰的分层架构：

1. **HTTP中间件** 移动到 `pkg/api/middleware/`
2. **基础设施中间件** 保留在 `pkg/infrastructure/middleware/` 并新增外部服务组件

## 重构前后对比

### 重构前
```
pkg/infrastructure/middleware/
├── error_handler.go  # HTTP错误处理
├── cors.go          # HTTP跨域处理
├── logger.go        # HTTP日志记录
└── recovery.go      # HTTP恢复处理
```

### 重构后
```
pkg/api/middleware/                    # HTTP中间件
├── error_handler.go  # HTTP错误处理
├── cors.go          # HTTP跨域处理
├── logger.go        # HTTP日志记录
└── recovery.go      # HTTP恢复处理

pkg/infrastructure/middleware/         # 基础设施中间件
├── redis.go         # Redis中间件
├── kafka.go         # Kafka中间件
└── elasticsearch.go # Elasticsearch中间件
```

## 重构内容

### 1. 文件移动

#### 移动的文件
- `pkg/infrastructure/middleware/error_handler.go` → `pkg/api/middleware/error_handler.go`
- `pkg/infrastructure/middleware/cors.go` → `pkg/api/middleware/cors.go`
- `pkg/infrastructure/middleware/logger.go` → `pkg/api/middleware/logger.go`
- `pkg/infrastructure/middleware/recovery.go` → `pkg/api/middleware/recovery.go`

#### 新增的文件
- `pkg/infrastructure/middleware/redis.go` - Redis中间件
- `pkg/infrastructure/middleware/kafka.go` - Kafka中间件
- `pkg/infrastructure/middleware/elasticsearch.go` - Elasticsearch中间件

### 2. 导入路径更新

#### 更新的文件
- `pkg/api/router/router.go` - 更新中间件导入路径
- `pkg/api/example.go` - 更新中间件导入路径

#### 更新内容
```go
// 更新前
import "github.com/make-bin/server-tpl/pkg/infrastructure/middleware"

// 更新后
import "github.com/make-bin/server-tpl/pkg/api/middleware"
```

### 3. 文档更新

#### 更新的文档
- `.cursor/rules/project.mdc` - 项目结构规则
- `README.md` - 项目说明文档
- `docs/middleware-structure.md` - 新增中间件结构说明

## 新增基础设施中间件

### 1. Redis中间件 (redis.go)

#### 功能特性
- 连接池管理
- 键值对操作
- 哈希、列表、集合、有序集合操作
- 过期时间管理
- 健康检查

#### 配置结构
```go
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"pool_size"`
}
```

#### 使用示例
```go
redisConfig := &middleware.RedisConfig{
    Host:     "localhost",
    Port:     6379,
    Password: "",
    DB:       0,
    PoolSize: 10,
}
redisMiddleware, err := middleware.NewRedisMiddleware(redisConfig)

// 使用Redis
err = redisMiddleware.Set(ctx, "key", "value", time.Hour)
value, err := redisMiddleware.Get(ctx, "key")
```

### 2. Kafka中间件 (kafka.go)

#### 功能特性
- 生产者和消费者
- 消息发送和接收
- 分区管理
- 偏移量控制
- 健康检查

#### 配置结构
```go
type KafkaConfig struct {
    Brokers []string `mapstructure:"brokers"`
    Version string   `mapstructure:"version"`
    GroupID string   `mapstructure:"group_id"`
}
```

#### 使用示例
```go
kafkaConfig := &middleware.KafkaConfig{
    Brokers: []string{"localhost:9092"},
    Version: "2.8.0",
    GroupID: "my-group",
}
kafkaMiddleware, err := middleware.NewKafkaMiddleware(kafkaConfig)

// 发送消息
err = kafkaMiddleware.SendMessage("topic", []byte("key"), []byte("value"))

// 消费消息
err = kafkaMiddleware.ConsumeMessages("topic", func(msg *sarama.ConsumerMessage) error {
    // 处理消息
    return nil
})
```

### 3. Elasticsearch中间件 (elasticsearch.go)

#### 功能特性
- 文档索引和搜索
- 批量操作
- 索引管理
- 统计信息
- 健康检查

#### 配置结构
```go
type ElasticsearchConfig struct {
    URLs     []string `mapstructure:"urls"`
    Username string   `mapstructure:"username"`
    Password string   `mapstructure:"password"`
    Index    string   `mapstructure:"index"`
}
```

#### 使用示例
```go
esConfig := &middleware.ElasticsearchConfig{
    URLs:     []string{"http://localhost:9200"},
    Username: "",
    Password: "",
    Index:    "my-index",
}
esMiddleware, err := middleware.NewElasticsearchMiddleware(esConfig)

// 索引文档
err = esMiddleware.IndexDocument("index", "id", document)

// 搜索文档
err = esMiddleware.SearchDocuments("index", query, result)
```

## 架构优势

### 1. 职责分离
- **HTTP中间件**: 专注于HTTP请求处理
- **基础设施中间件**: 专注于外部服务集成

### 2. 层次清晰
- **API层**: 处理HTTP相关功能
- **基础设施层**: 处理外部服务集成

### 3. 可扩展性
- 易于添加新的HTTP中间件
- 易于添加新的外部服务中间件

### 4. 可维护性
- 清晰的代码结构
- 完善的文档说明

## 配置管理

### HTTP中间件配置
HTTP中间件通过代码配置，支持环境变量覆盖：

```go
r := gin.New()
r.Use(middleware.ErrorHandler())
r.Use(middleware.Logger())
r.Use(middleware.CORS())
```

### 基础设施中间件配置
基础设施中间件通过配置文件配置：

```yaml
# configs/config.yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

kafka:
  brokers:
    - "localhost:9092"
  version: "2.8.0"
  group_id: "my-group"

elasticsearch:
  urls:
    - "http://localhost:9200"
  username: ""
  password: ""
  index: "my-index"
```

## 健康检查

### HTTP中间件健康检查
通过中间件链自动进行健康检查：

```go
r.GET("/health", func(c *gin.Context) {
    middleware.SuccessResponse(c, gin.H{
        "status":  "ok",
        "message": "service is running",
    })
})
```

### 基础设施中间件健康检查
提供独立的健康检查方法：

```go
// Redis健康检查
err := redisMiddleware.HealthCheck(ctx)

// Kafka健康检查
err := kafkaMiddleware.HealthCheck(ctx)

// Elasticsearch健康检查
err := esMiddleware.HealthCheck(ctx)
```

## 错误处理

### HTTP中间件错误处理
统一的错误处理机制：

```go
// 错误响应
middleware.ErrorResponse(c, err)

// 成功响应
middleware.SuccessResponse(c, data)
```

### 基础设施中间件错误处理
通过返回错误进行处理：

```go
if err != nil {
    logger.WithError(err).Error("Middleware operation failed")
    return err
}
```

## 最佳实践

### 1. 中间件顺序
HTTP中间件的使用顺序很重要：

```go
r := gin.New()
r.Use(middleware.ErrorHandler())  // 错误处理应该在最前面
r.Use(middleware.Logger())        // 日志记录
r.Use(middleware.CORS())          // 跨域处理
```

### 2. 配置管理
- 使用配置文件管理中间件配置
- 支持环境变量覆盖
- 提供默认配置值

### 3. 错误处理
- 统一的错误处理机制
- 详细的错误日志记录
- 友好的错误响应

### 4. 性能优化
- 连接池管理
- 超时控制
- 重试机制

### 5. 监控和日志
- 健康检查
- 性能指标
- 结构化日志

## 总结

通过本次中间件重构，项目实现了：

1. **清晰的职责分离**: HTTP中间件和基础设施中间件各司其职
2. **良好的可扩展性**: 易于添加新的中间件组件
3. **统一的错误处理**: 一致的错误处理机制
4. **完善的配置管理**: 灵活的配置选项
5. **可靠的健康检查**: 全面的健康检查机制

这种设计为项目的稳定运行和后续扩展提供了可靠的技术基础，使项目架构更加清晰和易于维护。
