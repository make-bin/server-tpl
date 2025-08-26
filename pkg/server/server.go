package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/make-bin/server-tpl/pkg/api"
	"github.com/make-bin/server-tpl/pkg/api/router"
	"github.com/make-bin/server-tpl/pkg/api/validation"
	"github.com/make-bin/server-tpl/pkg/domain/service"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/factory"
	"github.com/make-bin/server-tpl/pkg/utils/config"
	"github.com/make-bin/server-tpl/pkg/utils/container"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// Server HTTP服务器结构
type Server struct {
	config        *config.Config
	httpServer    *http.Server
	beanContainer *container.SimpleContainer
	dataStore     datastore.DatastoreInterface
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	return &Server{
		config:        cfg,
		beanContainer: container.NewContainer(),
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	logger.Info("Starting server initialization...")

	// 1. 初始化依赖注入容器
	if err := s.initContainer(); err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}

	// 2. 设置Gin模式
	if s.config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 3. 创建Gin引擎
	engine := gin.New()

	// 4. 初始化路由
	router.InitRouter(engine, nil) // 传入nil，因为我们使用依赖注入

	// 5. 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Server starting on port %d", s.config.Server.Port)
	return s.httpServer.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
	}

	// 清理容器
	if s.beanContainer != nil {
		s.beanContainer.Clear()
	}

	logger.Info("Server shutdown completed")
	return nil
}

// initContainer 初始化依赖注入容器
func (s *Server) initContainer() error {
	logger.Info("Initializing dependency injection container...")

	// 1. 注册工具类实现（config，logger，validator等）
	if err := s.registerUtilities(); err != nil {
		return fmt.Errorf("failed to register utilities: %w", err)
	}

	// 2. 注册基础设施组件（数据库、缓存等）
	if err := s.registerInfrastructure(); err != nil {
		return fmt.Errorf("failed to register infrastructure: %w", err)
	}

	// 3. 注册领域服务
	if err := s.registerDomainServices(); err != nil {
		return fmt.Errorf("failed to register domain services: %w", err)
	}

	// 4. 注册API层组件
	if err := s.registerAPIComponents(); err != nil {
		return fmt.Errorf("failed to register API components: %w", err)
	}

	// 5. 调用Populate()完成依赖注入
	if err := s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("failed to populate the bean container: %w", err)
	}

	logger.Info("Container initialization completed successfully")
	return nil
}

// registerUtilities 注册工具类
func (s *Server) registerUtilities() error {
	// 注册配置
	if err := s.beanContainer.ProvideWithName("config", s.config); err != nil {
		return fmt.Errorf("failed to register config: %w", err)
	}

	// 注册验证器
	v := validator.New()
	validation.RegisterCustomValidators(v)
	if err := s.beanContainer.ProvideWithName("validator", v); err != nil {
		return fmt.Errorf("failed to register validator: %w", err)
	}

	// 注册验证接口
	validationInterfaces := api.GetRegisterValidationInterfaces()
	for name, fn := range validationInterfaces {
		v.RegisterValidation(name, fn)
	}

	logger.Debug("Utilities registered successfully")
	return nil
}

// registerInfrastructure 注册基础设施组件
func (s *Server) registerInfrastructure() error {
	// 创建数据存储
	datastoreFactory := factory.NewSimpleFactory()
	datastore, err := datastoreFactory.CreateDatastore(s.config)
	if err != nil {
		return fmt.Errorf("failed to create datastore: %w", err)
	}

	// 执行数据库迁移
	if err := datastore.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 注册数据存储
	s.dataStore = datastore
	if err := s.beanContainer.ProvideWithName("datastore", datastore); err != nil {
		return fmt.Errorf("failed to register datastore: %w", err)
	}

	logger.Debug("Infrastructure components registered successfully")
	return nil
}

// registerDomainServices 注册领域服务
func (s *Server) registerDomainServices() error {
	// 注册服务beans
	serviceBeans := service.InitServiceBean()
	if err := s.beanContainer.Provides(serviceBeans...); err != nil {
		return fmt.Errorf("failed to register service beans: %w", err)
	}

	logger.Debug("Domain services registered successfully")
	return nil
}

// registerAPIComponents 注册API层组件
func (s *Server) registerAPIComponents() error {
	// 注册API接口
	apiBeans := api.InitAPI()
	if err := s.beanContainer.Provides(apiBeans...); err != nil {
		return fmt.Errorf("failed to register API beans: %w", err)
	}

	logger.Debug("API components registered successfully")
	return nil
}

// GetContainer 获取容器实例（用于测试或其他需要）
func (s *Server) GetContainer() *container.SimpleContainer {
	return s.beanContainer
}

// GetDataStore 获取数据存储实例（用于测试或其他需要）
func (s *Server) GetDataStore() datastore.DatastoreInterface {
	return s.dataStore
}

// HealthCheck 检查服务器健康状态
func (s *Server) HealthCheck() error {
	// 检查数据库连接
	if s.dataStore != nil {
		if err := s.dataStore.HealthCheck(); err != nil {
			return fmt.Errorf("datastore health check failed: %w", err)
		}
	}

	// 检查容器状态
	if s.beanContainer == nil {
		return fmt.Errorf("bean container not initialized")
	}

	return nil
}
