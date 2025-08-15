# 中间件结构说明

## 概述

项目中的中间件分为两个层次：

1. **HTTP中间件** (`pkg/api/middleware/`) - 处理HTTP请求的Gin中间件
2. **基础设施中间件** (`pkg/infrastructure/middleware/`) - 外部服务的基础组件

## HTTP中间件 (pkg/api/middleware/)

HTTP中间件位于API层，专门处理HTTP请求相关的功能。

### 文件结构

```
pkg/api/middleware/
├── error_handler.go  # 错误处理中间件
├── cors.go          # CORS跨域中间件
├── logger.go        # 日志中间件
└── recovery.go      # 恢复中间件
```

### 功能说明

#### 1. error_handler.go
- **功能**: 统一错误处理和响应
- **特性**:
  - 自动转换错误为HTTP响应
  - Panic自动恢复
  - 错误日志记录
  - HTTP状态码自动映射

#### 2. cors.go
- **功能**: 跨域资源共享处理
- **特性**:
  - 支持多种HTTP方法
  - 可配置允许的域名
  - 支持自定义头部

#### 3. logger.go
- **功能**: HTTP请求日志记录
- **特性**:
  - 结构化日志输出
  - 请求响应时间统计
  - 错误信息记录

#### 4. recovery.go
- **功能**: Panic恢复处理
- **特性**:
  - 自动捕获Panic
  - 返回友好的错误响应
  - 记录错误堆栈

### 使用方式

```go
import "github.com/make-bin/server-tpl/pkg/api/middleware"

// 在路由中使用
r := gin.New()
r.Use(middleware.ErrorHandler())
r.Use(middleware.Logger())
r.Use(middleware.CORS())
```

## 基础设施中间件 (pkg/infrastructure/middleware/)

基础设施中间件位于基础设施层，提供外部服务的基础组件实现。

### 文件结构

```
pkg/infrastructure/middleware/
├── redis.go         # Redis中间件
├── kafka.go         # Kafka中间件
└── elasticsearch.go # Elasticsearch中间件
```

### 功能说明

#### 1. redis.go
- **功能**: Redis缓存和存储中间件
- **特性**:
  - 连接池管理
  - 键值对操作
  - 哈希、列表、集合、有序集合操作
  - 过期时间管理
  - 健康检查

#### 2. kafka.go
- **功能**: Kafka消息队列中间件
- **特性**:
  - 生产者和消费者
  - 消息发送和接收
  - 分区管理
  - 偏移量控制
  - 健康检查

#### 3. elasticsearch.go
- **功能**: Elasticsearch搜索引擎中间件
- **特性**:
  - 文档索引和搜索
  - 批量操作
  - 索引管理
  - 统计信息
  - 健康检查

### 使用方式

```go
import "github.com/make-bin/server-tpl/pkg/infrastructure/middleware"

// Redis中间件
redisConfig := &middleware.RedisConfig{
    Host:     "localhost",
    Port:     6379,
    Password: "",
    DB:       0,
    PoolSize: 10,
}
redisMiddleware, err := middleware.NewRedisMiddleware(redisConfig)

// Kafka中间件
kafkaConfig := &middleware.KafkaConfig{
    Brokers: []string{"localhost:9092"},
    Version: "2.8.0",
    GroupID: "my-group",
}
kafkaMiddleware, err := middleware.NewKafkaMiddleware(kafkaConfig)

// Elasticsearch中间件
esConfig := &middleware.ElasticsearchConfig{
    URLs:     []string{"http://localhost:9200"},
    Username: "",
    Password: "",
    Index:    "my-index",
}
esMiddleware, err := middleware.NewElasticsearchMiddleware(esConfig)
```

## 中间件分层设计

### 设计原则

1. **职责分离**: HTTP中间件专注于HTTP请求处理，基础设施中间件专注于外部服务集成
2. **层次清晰**: API层处理HTTP相关，基础设施层处理外部服务
3. **可扩展性**: 易于添加新的中间件组件
4. **可维护性**: 清晰的代码结构和文档

### 架构优势

1. **模块化**: 每个中间件都是独立的模块
2. **可配置**: 支持灵活的配置选项
3. **可测试**: 每个中间件都可以独立测试
4. **可复用**: 中间件可以在不同项目中复用

## 配置管理

### HTTP中间件配置

HTTP中间件的配置通常通过代码或环境变量进行：

```go
// 错误处理中间件配置
r.Use(middleware.ErrorHandler())

// CORS中间件配置
r.Use(middleware.CORS())

// 日志中间件配置
r.Use(middleware.Logger())
```

### 基础设施中间件配置

基础设施中间件通过配置文件进行配置：

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

HTTP中间件通过中间件链自动进行健康检查：

```go
// 健康检查端点
r.GET("/health", func(c *gin.Context) {
    middleware.SuccessResponse(c, gin.H{
        "status":  "ok",
        "message": "service is running",
    })
})
```

### 基础设施中间件健康检查

基础设施中间件提供独立的健康检查方法：

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

HTTP中间件通过统一的错误处理机制：

```go
// 错误响应
middleware.ErrorResponse(c, err)

// 成功响应
middleware.SuccessResponse(c, data)
```

### 基础设施中间件错误处理

基础设施中间件通过返回错误进行处理：

```go
// 检查错误
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
// 其他中间件...
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

通过合理的中间件分层设计，项目实现了：

1. **清晰的职责分离**: HTTP中间件和基础设施中间件各司其职
2. **良好的可扩展性**: 易于添加新的中间件组件
3. **统一的错误处理**: 一致的错误处理机制
4. **完善的配置管理**: 灵活的配置选项
5. **可靠的健康检查**: 全面的健康检查机制

这种设计为项目的稳定运行和后续扩展提供了可靠的技术基础。
