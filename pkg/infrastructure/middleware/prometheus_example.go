package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusExample Prometheus中间件使用示例
func PrometheusExample() {
	// 1. 创建配置
	config := &PrometheusConfig{
		Enabled:     true,
		MetricsPath: "/metrics",
	}

	// 2. 创建中间件
	prometheus, err := NewPrometheusMiddleware(config)
	if err != nil {
		panic(err)
	}

	// 3. 启动指标服务器
	ctx := context.Background()
	if err := prometheus.StartMetricsServer(ctx); err != nil {
		panic(err)
	}

	// 4. 创建 Gin 引擎
	engine := gin.Default()

	// 5. 添加 Prometheus 中间件
	engine.Use(prometheus.HTTPMiddleware())

	// 6. 添加指标端点
	engine.GET("/metrics", prometheus.MetricsHandler())

	// 7. 示例业务路由
	engine.GET("/api/users", func(c *gin.Context) {
		start := time.Now()

		// 模拟业务操作
		time.Sleep(100 * time.Millisecond)

		// 记录业务操作
		prometheus.RecordBusinessOperation("get_users", "success", time.Since(start))

		c.JSON(200, gin.H{"message": "users"})
	})

	engine.POST("/api/users", func(c *gin.Context) {
		start := time.Now()

		// 模拟业务操作
		time.Sleep(200 * time.Millisecond)

		// 记录业务操作
		prometheus.RecordBusinessOperation("create_user", "success", time.Since(start))

		c.JSON(201, gin.H{"message": "user created"})
	})

	// 8. 示例数据库操作
	engine.GET("/api/db-test", func(c *gin.Context) {
		start := time.Now()

		// 模拟数据库查询
		time.Sleep(50 * time.Millisecond)

		// 记录数据库查询
		prometheus.RecordDatabaseQuery("postgresql", "select", time.Since(start))

		c.JSON(200, gin.H{"message": "db query"})
	})

	// 9. 示例缓存操作
	engine.GET("/api/cache-test", func(c *gin.Context) {
		// 模拟缓存命中
		prometheus.RecordCacheHit("redis")

		c.JSON(200, gin.H{"message": "cache hit"})
	})

	// 10. 启动服务器
	engine.Run(":8080")
}

// BusinessServiceExample 业务服务中使用 Prometheus 的示例
type BusinessServiceExample struct {
	prometheus *PrometheusMiddleware
}

// NewBusinessServiceExample 创建业务服务示例
func NewBusinessServiceExample(prometheus *PrometheusMiddleware) *BusinessServiceExample {
	return &BusinessServiceExample{
		prometheus: prometheus,
	}
}

// CreateUser 创建用户示例
func (s *BusinessServiceExample) CreateUser(userData map[string]interface{}) error {
	start := time.Now()

	// 模拟业务逻辑
	time.Sleep(100 * time.Millisecond)

	// 模拟数据库操作
	dbStart := time.Now()
	time.Sleep(50 * time.Millisecond)
	s.prometheus.RecordDatabaseQuery("postgresql", "insert", time.Since(dbStart))

	// 模拟缓存操作
	s.prometheus.RecordCacheMiss("redis")

	// 记录业务操作
	s.prometheus.RecordBusinessOperation("create_user", "success", time.Since(start))

	return nil
}

// GetUser 获取用户示例
func (s *BusinessServiceExample) GetUser(userID string) (map[string]interface{}, error) {
	start := time.Now()

	// 模拟业务逻辑
	time.Sleep(50 * time.Millisecond)

	// 模拟数据库操作
	dbStart := time.Now()
	time.Sleep(30 * time.Millisecond)
	s.prometheus.RecordDatabaseQuery("postgresql", "select", time.Since(dbStart))

	// 模拟缓存操作
	s.prometheus.RecordCacheHit("redis")

	// 记录业务操作
	s.prometheus.RecordBusinessOperation("get_user", "success", time.Since(start))

	return map[string]interface{}{
		"id":   userID,
		"name": "Test User",
	}, nil
}

// UpdateUser 更新用户示例（包含错误处理）
func (s *BusinessServiceExample) UpdateUser(userID string, userData map[string]interface{}) error {
	start := time.Now()

	// 模拟业务逻辑
	time.Sleep(80 * time.Millisecond)

	// 模拟数据库操作
	dbStart := time.Now()
	time.Sleep(40 * time.Millisecond)
	s.prometheus.RecordDatabaseQuery("postgresql", "update", time.Since(dbStart))

	// 模拟错误情况
	if userID == "invalid" {
		s.prometheus.RecordBusinessError("update_user", "validation_error")
		s.prometheus.RecordBusinessOperation("update_user", "error", time.Since(start))
		return fmt.Errorf("invalid user ID")
	}

	// 记录业务操作
	s.prometheus.RecordBusinessOperation("update_user", "success", time.Since(start))

	return nil
}

// DatabaseServiceExample 数据库服务中使用 Prometheus 的示例
type DatabaseServiceExample struct {
	prometheus *PrometheusMiddleware
}

// NewDatabaseServiceExample 创建数据库服务示例
func NewDatabaseServiceExample(prometheus *PrometheusMiddleware) *DatabaseServiceExample {
	return &DatabaseServiceExample{
		prometheus: prometheus,
	}
}

// Query 数据库查询示例
func (s *DatabaseServiceExample) Query(query string) ([]map[string]interface{}, error) {
	start := time.Now()

	// 模拟数据库查询
	time.Sleep(100 * time.Millisecond)

	// 记录查询
	s.prometheus.RecordDatabaseQuery("postgresql", "query", time.Since(start))

	return []map[string]interface{}{
		{"id": 1, "name": "User 1"},
		{"id": 2, "name": "User 2"},
	}, nil
}

// Insert 数据库插入示例
func (s *DatabaseServiceExample) Insert(data map[string]interface{}) error {
	start := time.Now()

	// 模拟数据库插入
	time.Sleep(150 * time.Millisecond)

	// 记录插入操作
	s.prometheus.RecordDatabaseQuery("postgresql", "insert", time.Since(start))

	return nil
}

// UpdateDatabaseConnections 更新数据库连接数示例
func (s *DatabaseServiceExample) UpdateDatabaseConnections(active, idle int) {
	s.prometheus.SetDatabaseConnections("postgresql", "active", active)
	s.prometheus.SetDatabaseConnections("postgresql", "idle", idle)
}

// CacheServiceExample 缓存服务中使用 Prometheus 的示例
type CacheServiceExample struct {
	prometheus *PrometheusMiddleware
}

// NewCacheServiceExample 创建缓存服务示例
func (s *CacheServiceExample) NewCacheServiceExample(prometheus *PrometheusMiddleware) *CacheServiceExample {
	return &CacheServiceExample{
		prometheus: prometheus,
	}
}

// Get 缓存获取示例
func (s *CacheServiceExample) Get(key string) (interface{}, error) {
	// 模拟缓存查询
	time.Sleep(10 * time.Millisecond)

	// 模拟缓存命中
	if key == "cached_key" {
		s.prometheus.RecordCacheHit("redis")
		return "cached_value", nil
	}

	// 模拟缓存未命中
	s.prometheus.RecordCacheMiss("redis")
	return nil, fmt.Errorf("key not found")
}

// Set 缓存设置示例
func (s *CacheServiceExample) Set(key string, value interface{}) error {
	// 模拟缓存设置
	time.Sleep(20 * time.Millisecond)

	// 更新缓存大小
	s.prometheus.SetCacheSize("redis", 1000)

	return nil
}

// UpdateCacheStats 更新缓存统计信息示例
func (s *CacheServiceExample) UpdateCacheStats(hits, misses, size int) {
	s.prometheus.SetCacheSize("redis", size)
	// 注意：hits 和 misses 通过 Get 方法自动记录
}
