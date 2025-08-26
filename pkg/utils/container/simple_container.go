package container

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// SimpleContainer 简单的依赖注入容器，按照规范实现
type SimpleContainer struct {
	beans map[string]interface{}
	mu    sync.RWMutex
}

// NewContainer 创建新的容器实例
func NewContainer() *SimpleContainer {
	return &SimpleContainer{
		beans: make(map[string]interface{}),
	}
}

// Provides 提供多个bean
func (c *SimpleContainer) Provides(beans ...interface{}) error {
	for _, bean := range beans {
		if err := c.ProvideWithName("", bean); err != nil {
			return err
		}
	}
	return nil
}

// ProvideWithName 提供带名称的bean
func (c *SimpleContainer) ProvideWithName(name string, bean interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "" {
		name = fmt.Sprintf("%T", bean)
	}

	// 检查是否已存在
	if _, exists := c.beans[name]; exists {
		return fmt.Errorf("bean with name '%s' already exists", name)
	}

	c.beans[name] = bean
	logger.Debug("Registered bean: %s", name)
	return nil
}

// Get 获取bean
func (c *SimpleContainer) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bean, exists := c.beans[name]
	return bean, exists
}

// GetByType 根据类型获取bean
func (c *SimpleContainer) GetByType(beanType reflect.Type) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, bean := range c.beans {
		if reflect.TypeOf(bean) == beanType {
			return bean, true
		}

		// 检查是否实现了接口
		if beanType.Kind() == reflect.Interface && reflect.TypeOf(bean).Implements(beanType) {
			return bean, true
		}
	}

	return nil, false
}

// Populate 填充依赖字段
func (c *SimpleContainer) Populate() error {
	start := time.Now()
	defer func() {
		logger.Info("Populate the bean container take time %s", time.Since(start))
	}()

	c.mu.Lock()
	defer c.mu.Unlock()

	// 遍历所有bean，进行依赖注入
	for name, bean := range c.beans {
		if err := c.injectDependencies(bean); err != nil {
			return fmt.Errorf("failed to inject dependencies for bean '%s': %w", name, err)
		}
	}

	return nil
}

// injectDependencies 注入依赖
func (c *SimpleContainer) injectDependencies(target interface{}) error {
	targetValue := reflect.ValueOf(target)

	// 如果是指针，获取元素
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	// 只处理结构体
	if targetValue.Kind() != reflect.Struct {
		return nil
	}

	targetType := targetValue.Type()

	// 遍历所有字段
	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)

		// 检查inject标签
		injectTag := fieldType.Tag.Get("inject")
		if injectTag == "" {
			continue
		}

		// 字段必须可设置
		if !field.CanSet() {
			logger.Warn("Field %s cannot be set", fieldType.Name)
			continue
		}

		// 根据标签值查找依赖
		var dependency interface{}
		var found bool

		if injectTag == "" {
			// 如果标签为空，按类型查找
			dependency, found = c.GetByType(field.Type())
		} else {
			// 按名称查找
			dependency, found = c.Get(injectTag)
		}

		if !found {
			// 尝试按类型名查找
			typeName := field.Type().String()
			dependency, found = c.Get(typeName)
		}

		if found {
			dependencyValue := reflect.ValueOf(dependency)
			if dependencyValue.Type().AssignableTo(field.Type()) {
				field.Set(dependencyValue)
				logger.Debug("Injected dependency for field: %s", fieldType.Name)
			} else {
				logger.Warn("Dependency type mismatch for field %s: expected %s, got %s",
					fieldType.Name, field.Type(), dependencyValue.Type())
			}
		} else {
			logger.Warn("Dependency not found for field: %s with inject tag: %s", fieldType.Name, injectTag)
		}
	}

	return nil
}

// ListBeans 列出所有注册的bean
func (c *SimpleContainer) ListBeans() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var names []string
	for name := range c.beans {
		names = append(names, name)
	}
	return names
}

// HasBean 检查bean是否存在
func (c *SimpleContainer) HasBean(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.beans[name]
	return exists
}

// Clear 清空容器
func (c *SimpleContainer) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.beans = make(map[string]interface{})
	logger.Info("Container cleared")
}

// Count 返回bean数量
func (c *SimpleContainer) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.beans)
}

// 为了向后兼容，保留原有的接口方法

// Register 注册服务（向后兼容）
func (c *SimpleContainer) Register(service interface{}) {
	serviceType := reflect.TypeOf(service)
	name := serviceType.String()
	c.ProvideWithName(name, service)
}

// RegisterAs 以特定接口类型注册服务（向后兼容）
func (c *SimpleContainer) RegisterAs(service interface{}, as interface{}) {
	interfaceType := reflect.TypeOf(as).Elem()
	name := interfaceType.String()
	c.ProvideWithName(name, service)
}

// Resolve 解析服务（向后兼容）
func (c *SimpleContainer) Resolve(serviceType interface{}) (interface{}, error) {
	t := reflect.TypeOf(serviceType).Elem()
	name := t.String()

	if service, exists := c.Get(name); exists {
		return service, nil
	}

	return nil, fmt.Errorf("service of type %s not found", t.String())
}

// MustResolve 解析服务，如果未找到则panic（向后兼容）
func (c *SimpleContainer) MustResolve(serviceType interface{}) interface{} {
	service, err := c.Resolve(serviceType)
	if err != nil {
		panic(err)
	}
	return service
}

// Has 检查服务是否已注册（向后兼容）
func (c *SimpleContainer) Has(serviceType interface{}) bool {
	t := reflect.TypeOf(serviceType).Elem()
	name := t.String()
	return c.HasBean(name)
}
