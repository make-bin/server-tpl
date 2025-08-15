# Prometheus 监控规范

## 概述

本项目采用Prometheus作为监控系统，通过复用Gin HTTP服务器提供指标端点，无需启动独立的HTTP服务器。这种设计简化了部署架构，减少了资源消耗，并提供了统一的端口管理。

## 架构设计

### 设计原则

1. **统一端口管理**：Prometheus指标端点复用Gin HTTP服务器，避免多端口管理
2. **简化部署**：无需额外的HTTP服务器进程，减少部署复杂度
3. **资源优化**：减少系统资源消耗，提高整体性能
4. **易于集成**：与现有监控系统无缝集成

### 技术实现

- **指标收集**：使用Prometheus Go客户端库
- **HTTP中间件**：集成到Gin框架中
- **指标端点**：通过Gin路由提供`/metrics`端点
- **系统指标**：后台goroutine定期收集系统资源指标

## 配置说明

### 配置文件

```yaml
prometheus:
  enabled: true                    # 是否启用Prometheus监控
  metrics_path: "/metrics"         # 指标端点路径
```

### 环境变量

```bash
PROMETHEUS_ENABLED=true
PROMETHEUS_METRICS_PATH=/metrics
```

## 指标类型

### 1. HTTP指标

#### http_requests_total
- **类型**：Counter
- **标签**：method, endpoint, status
- **说明**：HTTP请求总数统计

#### http_request_duration_seconds
- **类型**：Histogram
- **标签**：method, endpoint
- **说明**：HTTP请求响应时间分布

#### http_request_size_bytes
- **类型**：Histogram
- **标签**：method, endpoint
- **说明**：HTTP请求大小分布

#### http_response_size_bytes
- **类型**：Histogram
- **标签**：method, endpoint
- **说明**：HTTP响应大小分布

### 2. 业务指标

#### business_operations_total
- **类型**：Counter
- **标签**：operation, status
- **说明**：业务操作总数统计

#### business_operation_duration_seconds
- **类型**：Histogram
- **标签**：operation
- **说明**：业务操作耗时分布

#### business_errors_total
- **类型**：Counter
- **标签**：operation, error_type
- **说明**：业务错误统计

### 3. 系统指标

#### system_goroutines
- **类型**：Gauge
- **说明**：当前Goroutine数量

#### system_heap_alloc_bytes
- **类型**：Gauge
- **说明**：堆内存分配字节数

#### system_heap_sys_bytes
- **类型**：Gauge
- **说明**：从系统获取的堆内存字节数

#### system_memory_usage_bytes
- **类型**：Gauge
- **标签**：type (heap_alloc, heap_sys, heap_idle, heap_inuse)
- **说明**：系统内存使用情况

### 4. 数据库指标

#### database_connections
- **类型**：Gauge
- **标签**：database, status
- **说明**：数据库连接数

#### database_queries_total
- **类型**：Counter
- **标签**：database, operation
- **说明**：数据库查询总数

#### database_query_duration_seconds
- **类型**：Histogram
- **标签**：database, operation
- **说明**：数据库查询耗时分布

#### database_errors_total
- **类型**：Counter
- **标签**：database, error_type
- **说明**：数据库错误统计

### 5. 缓存指标

#### cache_hits_total
- **类型**：Counter
- **标签**：cache_type
- **说明**：缓存命中次数

#### cache_misses_total
- **类型**：Counter
- **标签**：cache_type
- **说明**：缓存未命中次数

#### cache_size
- **类型**：Gauge
- **标签**：cache_type
- **说明**：缓存大小

## 使用指南

### 1. 启用监控

在配置文件中启用Prometheus监控：

```yaml
prometheus:
  enabled: true
  metrics_path: "/metrics"
```

### 2. 访问指标

启动应用后，通过以下URL访问指标：

```
http://localhost:8080/metrics
```

### 3. 记录业务指标

在业务代码中使用Prometheus中间件记录指标：

```go
// 获取Prometheus中间件实例
prometheus := server.GetPrometheus()

// 记录业务操作
prometheus.RecordBusinessOperation("user_login", "success", duration)

// 记录业务错误
prometheus.RecordBusinessError("user_login", "invalid_credentials")

// 记录数据库查询
prometheus.RecordDatabaseQuery("postgres", "select", duration)

// 记录缓存操作
prometheus.RecordCacheHit("redis")
prometheus.RecordCacheMiss("redis")
```

### 4. 自定义指标

如需添加自定义指标，可以在Prometheus中间件中扩展：

```go
// 在PrometheusMiddleware结构体中添加新指标
customMetric := prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "custom_metric_total",
        Help: "Custom metric description",
    },
    []string{"label1", "label2"},
)

// 注册指标
prometheus.MustRegister(customMetric)

// 使用指标
customMetric.WithLabelValues("value1", "value2").Inc()
```

## Prometheus配置

### 1. 基础配置

在Prometheus配置文件中添加目标：

```yaml
scrape_configs:
  - job_name: 'go-http-server'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### 2. 告警规则

创建告警规则文件：

```yaml
groups:
  - name: go-http-server
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }} seconds"
```

### 3. Grafana仪表板

创建Grafana仪表板配置：

```json
{
  "dashboard": {
    "title": "Go HTTP Server Metrics",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

## 最佳实践

### 1. 指标命名

- 使用下划线分隔的小写字母
- 添加适当的单位后缀（如`_seconds`、`_bytes`）
- 使用描述性的名称

### 2. 标签使用

- 避免高基数标签（如用户ID）
- 使用有意义的标签值
- 保持标签数量合理

### 3. 性能考虑

- 避免在热点路径中记录过多指标
- 使用适当的指标类型（Counter、Gauge、Histogram）
- 定期清理不需要的指标

### 4. 监控策略

- 设置合理的告警阈值
- 监控关键业务指标
- 建立监控仪表板
- 定期审查监控配置

## 故障排查

### 1. 指标端点不可访问

检查项：
- 应用是否正常启动
- 配置文件中的`enabled`设置
- 防火墙设置
- 网络连接

### 2. 指标数据异常

检查项：
- 指标标签是否正确
- 业务代码中的指标记录逻辑
- Prometheus抓取配置
- 时间同步

### 3. 性能问题

检查项：
- 指标数量是否过多
- 标签基数是否过高
- 系统资源使用情况
- 网络延迟

## 总结

通过复用Gin HTTP服务器提供Prometheus指标端点，我们实现了：

1. **简化的架构**：无需额外的HTTP服务器
2. **统一的端口管理**：所有HTTP服务使用同一端口
3. **资源优化**：减少系统资源消耗
4. **易于部署**：降低部署复杂度
5. **全面的监控**：覆盖HTTP、业务、系统、数据库、缓存等各个方面

这种设计为生产环境提供了高效、可靠的监控解决方案。
