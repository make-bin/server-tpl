package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusConfig Prometheus配置
type PrometheusConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	MetricsPath string `mapstructure:"metrics_path"`
}

// PrometheusMiddleware Prometheus中间件
type PrometheusMiddleware struct {
	config *PrometheusConfig

	// HTTP 请求指标
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	// 业务指标
	businessOperationsTotal   *prometheus.CounterVec
	businessOperationDuration *prometheus.HistogramVec
	businessErrorsTotal       *prometheus.CounterVec

	// 系统指标
	systemMemoryUsage *prometheus.GaugeVec
	systemCPUUsage    *prometheus.GaugeVec
	systemGoroutines  prometheus.Gauge
	systemHeapAlloc   prometheus.Gauge
	systemHeapSys     prometheus.Gauge

	// 数据库指标
	databaseConnections   *prometheus.GaugeVec
	databaseQueries       *prometheus.CounterVec
	databaseQueryDuration *prometheus.HistogramVec
	databaseErrors        *prometheus.CounterVec

	// 缓存指标
	cacheHits   *prometheus.CounterVec
	cacheMisses *prometheus.CounterVec
	cacheSize   *prometheus.GaugeVec
}

// NewPrometheusMiddleware 创建Prometheus中间件
func NewPrometheusMiddleware(config *PrometheusConfig) (*PrometheusMiddleware, error) {
	if !config.Enabled {
		return &PrometheusMiddleware{config: config}, nil
	}

	// 初始化HTTP指标
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	// 初始化业务指标
	businessOperationsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_operations_total",
			Help: "Total number of business operations",
		},
		[]string{"operation", "status"},
	)

	businessOperationDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "business_operation_duration_seconds",
			Help:    "Business operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	businessErrorsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_errors_total",
			Help: "Total number of business errors",
		},
		[]string{"operation", "error_type"},
	)

	// 初始化系统指标
	systemMemoryUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_bytes",
			Help: "System memory usage in bytes",
		},
		[]string{"type"},
	)

	systemCPUUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "System CPU usage percentage",
		},
		[]string{"core"},
	)

	systemGoroutines := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_goroutines",
			Help: "Number of goroutines",
		},
	)

	systemHeapAlloc := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_heap_alloc_bytes",
			Help: "Heap memory allocated in bytes",
		},
	)

	systemHeapSys := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_heap_sys_bytes",
			Help: "Heap memory obtained from system in bytes",
		},
	)

	// 初始化数据库指标
	databaseConnections := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"database", "status"},
	)

	databaseQueries := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"database", "operation"},
	)

	databaseQueryDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"database", "operation"},
	)

	databaseErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"database", "error_type"},
	)

	// 初始化缓存指标
	cacheHits := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	cacheMisses := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	cacheSize := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Cache size",
		},
		[]string{"cache_type"},
	)

	// 注册指标
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestSize,
		httpResponseSize,
		businessOperationsTotal,
		businessOperationDuration,
		businessErrorsTotal,
		systemMemoryUsage,
		systemCPUUsage,
		systemGoroutines,
		systemHeapAlloc,
		systemHeapSys,
		databaseConnections,
		databaseQueries,
		databaseQueryDuration,
		databaseErrors,
		cacheHits,
		cacheMisses,
		cacheSize,
	)

	logger.Info("Prometheus middleware initialized successfully")

	return &PrometheusMiddleware{
		config:                    config,
		httpRequestsTotal:         httpRequestsTotal,
		httpRequestDuration:       httpRequestDuration,
		httpRequestSize:           httpRequestSize,
		httpResponseSize:          httpResponseSize,
		businessOperationsTotal:   businessOperationsTotal,
		businessOperationDuration: businessOperationDuration,
		businessErrorsTotal:       businessErrorsTotal,
		systemMemoryUsage:         systemMemoryUsage,
		systemCPUUsage:            systemCPUUsage,
		systemGoroutines:          systemGoroutines,
		systemHeapAlloc:           systemHeapAlloc,
		systemHeapSys:             systemHeapSys,
		databaseConnections:       databaseConnections,
		databaseQueries:           databaseQueries,
		databaseQueryDuration:     databaseQueryDuration,
		databaseErrors:            databaseErrors,
		cacheHits:                 cacheHits,
		cacheMisses:               cacheMisses,
		cacheSize:                 cacheSize,
	}, nil
}

// HTTPMiddleware HTTP中间件
func (p *PrometheusMiddleware) HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !p.config.Enabled {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 记录请求大小
		requestSize := c.Request.ContentLength
		if requestSize > 0 {
			p.httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
		}

		// 处理请求
		c.Next()

		// 记录响应大小
		responseSize := c.Writer.Size()
		if responseSize > 0 {
			p.httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
		}

		// 记录请求总数和持续时间
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		p.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		p.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// MetricsHandler 指标处理器
func (p *PrometheusMiddleware) MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// StartMetricsServer 启动指标服务器（已废弃，现在复用Gin HTTP服务器）
// 此方法保留是为了向后兼容，实际不再启动独立的HTTP服务器
func (p *PrometheusMiddleware) StartMetricsServer(ctx context.Context) error {
	if !p.config.Enabled {
		return nil
	}

	// 启动系统指标收集
	go p.collectSystemMetrics(ctx)

	logger.WithField("path", p.config.MetricsPath).
		Info("Prometheus metrics will be served through Gin HTTP server")

	return nil
}

// collectSystemMetrics 收集系统指标
func (p *PrometheusMiddleware) collectSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			p.systemGoroutines.Set(float64(runtime.NumGoroutine()))
			p.systemHeapAlloc.Set(float64(m.HeapAlloc))
			p.systemHeapSys.Set(float64(m.HeapSys))

			p.systemMemoryUsage.WithLabelValues("heap_alloc").Set(float64(m.HeapAlloc))
			p.systemMemoryUsage.WithLabelValues("heap_sys").Set(float64(m.HeapSys))
			p.systemMemoryUsage.WithLabelValues("heap_idle").Set(float64(m.HeapIdle))
			p.systemMemoryUsage.WithLabelValues("heap_inuse").Set(float64(m.HeapInuse))
		}
	}
}

// RecordBusinessOperation 记录业务操作
func (p *PrometheusMiddleware) RecordBusinessOperation(operation, status string, duration time.Duration) {
	if !p.config.Enabled {
		return
	}

	p.businessOperationsTotal.WithLabelValues(operation, status).Inc()
	p.businessOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordBusinessError 记录业务错误
func (p *PrometheusMiddleware) RecordBusinessError(operation, errorType string) {
	if !p.config.Enabled {
		return
	}

	p.businessErrorsTotal.WithLabelValues(operation, errorType).Inc()
}

// RecordDatabaseQuery 记录数据库查询
func (p *PrometheusMiddleware) RecordDatabaseQuery(database, operation string, duration time.Duration) {
	if !p.config.Enabled {
		return
	}

	p.databaseQueries.WithLabelValues(database, operation).Inc()
	p.databaseQueryDuration.WithLabelValues(database, operation).Observe(duration.Seconds())
}

// RecordDatabaseError 记录数据库错误
func (p *PrometheusMiddleware) RecordDatabaseError(database, errorType string) {
	if !p.config.Enabled {
		return
	}

	p.databaseErrors.WithLabelValues(database, errorType).Inc()
}

// SetDatabaseConnections 设置数据库连接数
func (p *PrometheusMiddleware) SetDatabaseConnections(database, status string, count int) {
	if !p.config.Enabled {
		return
	}

	p.databaseConnections.WithLabelValues(database, status).Set(float64(count))
}

// RecordCacheHit 记录缓存命中
func (p *PrometheusMiddleware) RecordCacheHit(cacheType string) {
	if !p.config.Enabled {
		return
	}

	p.cacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (p *PrometheusMiddleware) RecordCacheMiss(cacheType string) {
	if !p.config.Enabled {
		return
	}

	p.cacheMisses.WithLabelValues(cacheType).Inc()
}

// SetCacheSize 设置缓存大小
func (p *PrometheusMiddleware) SetCacheSize(cacheType string, size int) {
	if !p.config.Enabled {
		return
	}

	p.cacheSize.WithLabelValues(cacheType).Set(float64(size))
}

// GetMetrics 获取指标数据
func (p *PrometheusMiddleware) GetMetrics() (map[string]interface{}, error) {
	if !p.config.Enabled {
		return nil, fmt.Errorf("prometheus metrics not enabled")
	}

	// 这里可以实现自定义指标聚合逻辑
	metrics := map[string]interface{}{
		"http_requests_total":       "available",
		"http_request_duration":     "available",
		"business_operations_total": "available",
		"database_queries_total":    "available",
		"cache_hits_total":          "available",
		"system_goroutines":         "available",
		"system_heap_alloc":         "available",
	}

	return metrics, nil
}

// Close 关闭中间件
func (p *PrometheusMiddleware) Close(ctx context.Context) error {
	if !p.config.Enabled {
		return nil
	}

	logger.Info("Prometheus middleware closed")
	return nil
}
