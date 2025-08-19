# 10. Prometheus 监控规范

## 10.1 指标类型
项目使用 Prometheus 进行监控，支持以下指标类型：

### 10.1.1 HTTP 请求指标
- `http_requests_total`: HTTP 请求总数（Counter）
- `http_request_duration_seconds`: HTTP 请求持续时间（Histogram）
- `http_request_size_bytes`: HTTP 请求大小（Histogram）
- `http_response_size_bytes`: HTTP 响应大小（Histogram）

### 10.1.2 业务指标
- `business_operations_total`: 业务操作总数（Counter）
- `business_operation_duration_seconds`: 业务操作持续时间（Histogram）
- `business_errors_total`: 业务错误总数（Counter）

### 10.1.3 系统指标
- `system_memory_usage_bytes`: 系统内存使用量（Gauge）
- `system_cpu_usage_percent`: 系统 CPU 使用率（Gauge）
- `system_goroutines`: Goroutine 数量（Gauge）
- `system_heap_alloc_bytes`: 堆内存分配量（Gauge）
- `system_heap_sys_bytes`: 堆内存系统分配量（Gauge）

### 10.1.4 数据库指标
- `database_connections`: 数据库连接数（Gauge）
- `database_queries_total`: 数据库查询总数（Counter）
- `database_query_duration_seconds`: 数据库查询持续时间（Histogram）
- `database_errors_total`: 数据库错误总数（Counter）

### 10.1.5 缓存指标
- `cache_hits_total`: 缓存命中总数（Counter）
- `cache_misses_total`: 缓存未命中总数（Counter）
- `cache_size`: 缓存大小（Gauge）

## 10.2 配置规范
```yaml
# configs/config.yaml
prometheus:
  enabled: true
  metrics_path: "/metrics"
  port: 9090
  host: "localhost"
```

## 10.3 使用规范
```go
// 1. 创建配置
config := &middleware.PrometheusConfig{
    Enabled:     true,
    MetricsPath: "/metrics",
    Port:        9090,
    Host:        "localhost",
}

// 2. 创建中间件
prometheus, err := middleware.NewPrometheusMiddleware(config)
if err != nil {
    log.Fatal(err)
}

// 3. 添加到 Gin 引擎
engine.Use(prometheus.HTTPMiddleware())
engine.GET("/metrics", prometheus.MetricsHandler())

// 4. 记录业务指标
prometheus.RecordBusinessOperation("create_user", "success", duration)
prometheus.RecordDatabaseQuery("postgresql", "insert", duration)
prometheus.RecordCacheHit("redis")
```

## 10.4 监控最佳实践
- **指标命名**: 使用下划线分隔的小写字母
- **标签设计**: 合理使用标签，避免高基数问题
- **指标聚合**: 在应用层进行指标聚合
- **错误处理**: 记录所有错误类型和频率
- **性能监控**: 监控关键业务操作的性能
- **资源监控**: 监控系统资源使用情况

## 10.5 Grafana 仪表板
项目提供标准的 Grafana 仪表板配置，包括：
- HTTP 请求监控面板
- 业务指标监控面板
- 系统资源监控面板
- 数据库性能监控面板
- 缓存性能监控面板
