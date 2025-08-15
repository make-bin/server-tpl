package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/make-bin/server-tpl/pkg/api"
	api_middleware "github.com/make-bin/server-tpl/pkg/api/middleware"
	"github.com/make-bin/server-tpl/pkg/domain/service"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/factory"
	infra_middleware "github.com/make-bin/server-tpl/pkg/infrastructure/middleware"
	"github.com/make-bin/server-tpl/pkg/utils/container"
)

type Server struct {
	engine        *gin.Engine
	beanContainer *container.Container
	dataStore     datastore.DataStore
	prometheus    *infra_middleware.PrometheusMiddleware
}

func NewServer() *Server {
	return &Server{
		engine:        gin.Default(),
		beanContainer: container.NewContainer(),
	}
}

func (s *Server) Start() error {
	// 初始化配置
	if err := s.initConfig(); err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}

	// 初始化数据存储
	if err := s.initDataStore(); err != nil {
		return fmt.Errorf("failed to init datastore: %w", err)
	}

	// 初始化 Prometheus
	if err := s.initPrometheus(); err != nil {
		return fmt.Errorf("failed to init prometheus: %w", err)
	}

	// 初始化依赖注入容器
	if err := s.initContainer(); err != nil {
		return fmt.Errorf("failed to init container: %w", err)
	}

	// 初始化中间件
	s.initMiddleware()

	// 初始化路由
	s.initRoutes()

	// 启动服务器
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	return s.engine.Run(":" + port)
}

// GetPrometheus 获取 Prometheus 中间件实例
func (s *Server) GetPrometheus() *infra_middleware.PrometheusMiddleware {
	return s.prometheus
}

func (s *Server) initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.type", "memory") // 默认使用内存存储
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "server_tpl")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_idle", 10)
	viper.SetDefault("database.max_open", 100)
	viper.SetDefault("database.timeout", 30)

	// Prometheus 配置默认值
	viper.SetDefault("prometheus.enabled", true)
	viper.SetDefault("prometheus.metrics_path", "/metrics")

	// PProf 调试配置默认值
	viper.SetDefault("pprof.enabled", false)
	viper.SetDefault("pprof.path_prefix", "/debug/pprof")

	return viper.ReadInConfig()
}

func (s *Server) initDataStore() error {
	// 从配置中读取数据库配置
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

	// 使用工厂创建数据存储实例
	var err error
	s.dataStore, err = factory.NewDataStore(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create datastore: %w", err)
	}

	// 连接数据库
	ctx := context.Background()
	if err := s.dataStore.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 健康检查
	if err := s.dataStore.HealthCheck(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (s *Server) initPrometheus() error {
	// 创建 Prometheus 配置
	config := &infra_middleware.PrometheusConfig{
		Enabled:     viper.GetBool("prometheus.enabled"),
		MetricsPath: viper.GetString("prometheus.metrics_path"),
	}

	// 创建 Prometheus 中间件
	var err error
	s.prometheus, err = infra_middleware.NewPrometheusMiddleware(config)
	if err != nil {
		return fmt.Errorf("failed to create prometheus middleware: %w", err)
	}

	// 启动系统指标收集（复用Gin HTTP服务器，无需独立启动）
	ctx := context.Background()
	if err := s.prometheus.StartMetricsServer(ctx); err != nil {
		return fmt.Errorf("failed to start prometheus metrics collection: %w", err)
	}

	return nil
}

func (s *Server) initContainer() error {
	// 注册数据存储到容器
	if err := s.beanContainer.ProvideWithName("datastore", s.dataStore); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}

	// 初始化领域服务
	if err := s.beanContainer.Provides(service.InitServiceBean()...); err != nil {
		return err
	}

	// 初始化API接口
	if err := s.beanContainer.Provides(api.InitAPI()...); err != nil {
		return err
	}

	// 填充依赖
	if err := s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}

	return nil
}

func (s *Server) initMiddleware() {
	// 添加 Prometheus HTTP 中间件（如果启用）
	if s.prometheus != nil && viper.GetBool("prometheus.enabled") {
		s.engine.Use(s.prometheus.HTTPMiddleware())
	}

	// 添加全局中间件
	s.engine.Use(api_middleware.CORS())
	s.engine.Use(api_middleware.Logger())
	s.engine.Use(api_middleware.Recovery())
}

func (s *Server) initRoutes() {
	// API路由组
	apiGroup := s.engine.Group("/api/v1")

	// 获取注册的API接口并初始化路由
	apiInterfaces := api.GetRegisterAPIInterfaces()
	for _, apiInterface := range apiInterfaces {
		apiInterface.InitAPIServiceRoute(apiGroup)
	}

	// 健康检查
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Prometheus 指标端点（如果启用）
	if s.prometheus != nil && viper.GetBool("prometheus.enabled") {
		s.engine.GET(viper.GetString("prometheus.metrics_path"), s.prometheus.MetricsHandler())
	}

	// PProf 调试端点（如果启用）
	if viper.GetBool("pprof.enabled") {
		s.initPProfRoutes()
	}
}

// initPProfRoutes 初始化PProf调试路由
func (s *Server) initPProfRoutes() {
	pprofPrefix := viper.GetString("pprof.path_prefix")

	// 注册所有pprof路由
	s.engine.GET(pprofPrefix+"/", gin.WrapF(pprof.Index))
	s.engine.GET(pprofPrefix+"/cmdline", gin.WrapF(pprof.Cmdline))
	s.engine.GET(pprofPrefix+"/profile", gin.WrapF(pprof.Profile))
	s.engine.GET(pprofPrefix+"/symbol", gin.WrapF(pprof.Symbol))
	s.engine.GET(pprofPrefix+"/trace", gin.WrapF(pprof.Trace))
	s.engine.GET(pprofPrefix+"/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	s.engine.GET(pprofPrefix+"/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	s.engine.GET(pprofPrefix+"/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	s.engine.GET(pprofPrefix+"/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	s.engine.GET(pprofPrefix+"/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	s.engine.GET(pprofPrefix+"/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
}
