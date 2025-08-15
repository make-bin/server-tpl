package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/domain/service"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/factory"
	"github.com/make-bin/server-tpl/pkg/utils/container"
)

func main() {
	// 初始化配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.type", "memory")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "server_tpl")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_idle", 10)
	viper.SetDefault("database.max_open", 100)
	viper.SetDefault("database.timeout", 30)

	// 创建数据存储
	dbConfig := &datastore.Config{
		Type:     viper.GetString("database.type"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.database"),
		SSLMode:  viper.GetString("database.sslmode"),
		MaxIdle:  viper.GetInt("database.max_idle"),
		MaxOpen:  viper.GetInt("database.max_open"),
		Timeout:  viper.GetInt("database.timeout"),
	}

	store, err := factory.NewDataStore(dbConfig)
	if err != nil {
		log.Fatalf("Failed to create datastore: %v", err)
	}

	// 连接数据库
	ctx := context.Background()
	if err := store.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Disconnect(ctx)

	// 健康检查
	if err := store.HealthCheck(ctx); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// 初始化依赖注入容器
	container := container.NewContainer()
	if err := container.ProvideWithName("datastore", store); err != nil {
		log.Fatalf("Failed to provide datastore: %v", err)
	}

	// 初始化服务
	if err := container.Provides(service.InitServiceBean()...); err != nil {
		log.Fatalf("Failed to provide services: %v", err)
	}

	// 填充依赖
	if err := container.Populate(); err != nil {
		log.Fatalf("Failed to populate container: %v", err)
	}

	// 创建 Gin 引擎
	engine := gin.Default()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 健康检查路由
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// 简单的测试路由
	engine.GET("/test", func(c *gin.Context) {
		// 创建测试用户
		user := &model.User{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "password123",
			Role:     "user",
			Status:   "active",
		}

		// 直接使用存储层
		if err := store.Add(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Test user created successfully",
			"user":    user,
		})
	})

	// 启动服务器
	port := viper.GetString("server.port")
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Database type: %s\n", viper.GetString("database.type"))
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("Test endpoint: http://localhost:%s/test\n", port)

	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
