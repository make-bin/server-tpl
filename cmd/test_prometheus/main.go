package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/infrastructure/middleware"
)

func main() {
	// 创建 Prometheus 配置
	config := &middleware.PrometheusConfig{
		Enabled:     true,
		MetricsPath: "/metrics",
		Port:        9090,
		Host:        "localhost",
	}

	// 创建 Prometheus 中间件
	prometheus, err := middleware.NewPrometheusMiddleware(config)
	if err != nil {
		log.Fatalf("Failed to create Prometheus middleware: %v", err)
	}

	// 启动指标服务器
	ctx := context.Background()
	if err := prometheus.StartMetricsServer(ctx); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}

	// 创建 Gin 引擎
	engine := gin.Default()

	// 添加 Prometheus 中间件
	engine.Use(prometheus.HTTPMiddleware())

	// 添加指标端点
	engine.GET("/metrics", prometheus.MetricsHandler())

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// 测试业务操作
	engine.GET("/api/users", func(c *gin.Context) {
		start := time.Now()

		// 模拟业务操作
		time.Sleep(100 * time.Millisecond)

		// 记录业务操作
		prometheus.RecordBusinessOperation("get_users", "success", time.Since(start))

		c.JSON(http.StatusOK, gin.H{
			"message": "users retrieved successfully",
			"count":   10,
		})
	})

	engine.POST("/api/users", func(c *gin.Context) {
		start := time.Now()

		// 模拟业务操作
		time.Sleep(200 * time.Millisecond)

		// 记录业务操作
		prometheus.RecordBusinessOperation("create_user", "success", time.Since(start))

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
			"id":      123,
		})
	})

	// 测试数据库操作
	engine.GET("/api/db-test", func(c *gin.Context) {
		start := time.Now()

		// 模拟数据库查询
		time.Sleep(50 * time.Millisecond)

		// 记录数据库查询
		prometheus.RecordDatabaseQuery("postgresql", "select", time.Since(start))

		c.JSON(http.StatusOK, gin.H{
			"message": "database query executed",
			"result":  "success",
		})
	})

	engine.POST("/api/db-test", func(c *gin.Context) {
		start := time.Now()

		// 模拟数据库插入
		time.Sleep(150 * time.Millisecond)

		// 记录数据库操作
		prometheus.RecordDatabaseQuery("postgresql", "insert", time.Since(start))

		c.JSON(http.StatusOK, gin.H{
			"message": "database insert executed",
			"result":  "success",
		})
	})

	// 测试缓存操作
	engine.GET("/api/cache-test", func(c *gin.Context) {
		key := c.Query("key")

		if key == "cached" {
			// 模拟缓存命中
			prometheus.RecordCacheHit("redis")
			c.JSON(http.StatusOK, gin.H{
				"message": "cache hit",
				"value":   "cached_value",
			})
		} else {
			// 模拟缓存未命中
			prometheus.RecordCacheMiss("redis")
			c.JSON(http.StatusOK, gin.H{
				"message": "cache miss",
				"value":   "fetched_from_db",
			})
		}
	})

	// 测试错误情况
	engine.GET("/api/error-test", func(c *gin.Context) {
		start := time.Now()

		// 模拟错误
		prometheus.RecordBusinessError("get_data", "validation_error")
		prometheus.RecordBusinessOperation("get_data", "error", time.Since(start))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed",
		})
	})

	// 模拟系统指标更新
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 模拟数据库连接数更新
				prometheus.SetDatabaseConnections("postgresql", "active", 5)
				prometheus.SetDatabaseConnections("postgresql", "idle", 3)

				// 模拟缓存大小更新
				prometheus.SetCacheSize("redis", 1000)

				fmt.Println("System metrics updated")
			}
		}
	}()

	// 启动服务器
	fmt.Println("Prometheus test server starting on port 8080...")
	fmt.Println("Metrics endpoint: http://localhost:9090/metrics")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("Test endpoints:")
	fmt.Println("  GET  /api/users")
	fmt.Println("  POST /api/users")
	fmt.Println("  GET  /api/db-test")
	fmt.Println("  POST /api/db-test")
	fmt.Println("  GET  /api/cache-test?key=cached")
	fmt.Println("  GET  /api/cache-test?key=miss")
	fmt.Println("  GET  /api/error-test")

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
