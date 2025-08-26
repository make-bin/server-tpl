package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container interface defines dependency injection container
type Container interface {
	Provide(name string, bean interface{}) error
	Get(name string) (interface{}, bool)
	Populate() error
	Close() error
}

// Lifecycle represents the lifecycle of a dependency
type Lifecycle int

const (
	Singleton Lifecycle = iota
	Prototype
	Session
	Request
)

// Bean represents a registered dependency
type Bean struct {
	Name      string
	Value     interface{}
	Type      reflect.Type
	Lifecycle Lifecycle
	Factory   func() interface{}
	Tags      map[string]string
}

// DIContainer implements a comprehensive dependency injection container
type DIContainer struct {
	beans    map[string]*Bean
	services map[reflect.Type]interface{}
	mutex    sync.RWMutex
}

// New creates a new container instance
func New() Container {
	return &DIContainer{
		beans:    make(map[string]*Bean),
		services: make(map[reflect.Type]interface{}),
	}
}

// Provide registers a bean in the container
func (c *DIContainer) Provide(name string, bean interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.beans[name]; exists {
		return fmt.Errorf("bean with name '%s' already exists", name)
	}

	beanType := reflect.TypeOf(bean)
	c.beans[name] = &Bean{
		Name:      name,
		Value:     bean,
		Type:      beanType,
		Lifecycle: Singleton,
		Tags:      make(map[string]string),
	}

	// Also register by type for backward compatibility
	c.services[beanType] = bean

	return nil
}

// Get retrieves a bean from the container
func (c *DIContainer) Get(name string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	bean, exists := c.beans[name]
	if !exists {
		return nil, false
	}

	switch bean.Lifecycle {
	case Singleton:
		return bean.Value, true
	case Prototype:
		if bean.Factory != nil {
			return bean.Factory(), true
		}
		return bean.Value, true
	default:
		return bean.Value, true
	}
}

// Populate performs dependency injection on registered beans
func (c *DIContainer) Populate() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, bean := range c.beans {
		if err := c.injectDependencies(bean.Value); err != nil {
			return fmt.Errorf("failed to inject dependencies for bean '%s': %w", bean.Name, err)
		}
	}

	return nil
}

// Close cleans up the container
func (c *DIContainer) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Cleanup logic if needed
	c.beans = make(map[string]*Bean)
	c.services = make(map[reflect.Type]interface{})

	return nil
}

// injectDependencies performs dependency injection using struct tags
func (c *DIContainer) injectDependencies(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	if targetValue.Kind() != reflect.Struct {
		return nil
	}

	targetType := targetValue.Type()

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)

		// Check for inject tag
		injectTag := fieldType.Tag.Get("inject")
		if injectTag == "" {
			continue
		}

		if !field.CanSet() {
			continue
		}

		// Find dependency by name
		if dependency, exists := c.Get(injectTag); exists {
			dependencyValue := reflect.ValueOf(dependency)
			if dependencyValue.Type().AssignableTo(field.Type()) {
				field.Set(dependencyValue)
			}
		}
	}

	return nil
}

// RegisterWithLifecycle registers a bean with specific lifecycle
func (c *DIContainer) RegisterWithLifecycle(name string, bean interface{}, lifecycle Lifecycle) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.beans[name]; exists {
		return fmt.Errorf("bean with name '%s' already exists", name)
	}

	beanType := reflect.TypeOf(bean)
	c.beans[name] = &Bean{
		Name:      name,
		Value:     bean,
		Type:      beanType,
		Lifecycle: lifecycle,
		Tags:      make(map[string]string),
	}

	return nil
}

// RegisterFactory registers a bean factory
func (c *DIContainer) RegisterFactory(name string, factory func() interface{}, lifecycle Lifecycle) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.beans[name]; exists {
		return fmt.Errorf("bean with name '%s' already exists", name)
	}

	// Create initial instance to get type
	instance := factory()
	beanType := reflect.TypeOf(instance)

	c.beans[name] = &Bean{
		Name:      name,
		Value:     instance,
		Type:      beanType,
		Lifecycle: lifecycle,
		Factory:   factory,
		Tags:      make(map[string]string),
	}

	return nil
}

// GetBeansByType returns all beans of a specific type
func (c *DIContainer) GetBeansByType(beanType reflect.Type) []interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []interface{}
	for _, bean := range c.beans {
		if bean.Type == beanType || bean.Type.Implements(beanType) {
			result = append(result, bean.Value)
		}
	}

	return result
}

// GetBeanNames returns all registered bean names
func (c *DIContainer) GetBeanNames() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var names []string
	for name := range c.beans {
		names = append(names, name)
	}

	return names
}

// HasBean checks if a bean exists
func (c *DIContainer) HasBean(name string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, exists := c.beans[name]
	return exists
}

// Legacy methods for backward compatibility

// Register registers a service in the container (backward compatibility)
func (c *DIContainer) Register(service interface{}) {
	serviceType := reflect.TypeOf(service)
	name := serviceType.String()
	c.Provide(name, service)
}

// RegisterAs registers a service with a specific interface type (backward compatibility)
func (c *DIContainer) RegisterAs(service interface{}, as interface{}) {
	interfaceType := reflect.TypeOf(as).Elem()
	name := interfaceType.String()
	c.Provide(name, service)
}

// Resolve resolves a service from the container (backward compatibility)
func (c *DIContainer) Resolve(serviceType interface{}) (interface{}, error) {
	t := reflect.TypeOf(serviceType).Elem()
	name := t.String()

	if service, exists := c.Get(name); exists {
		return service, nil
	}

	return nil, fmt.Errorf("service of type %s not found", t.String())
}

// MustResolve resolves a service from the container and panics if not found (backward compatibility)
func (c *DIContainer) MustResolve(serviceType interface{}) interface{} {
	service, err := c.Resolve(serviceType)
	if err != nil {
		panic(err)
	}
	return service
}

// Has checks if a service is registered in the container (backward compatibility)
func (c *DIContainer) Has(serviceType interface{}) bool {
	t := reflect.TypeOf(serviceType).Elem()
	name := t.String()
	return c.HasBean(name)
}

// Clear removes all services from the container (backward compatibility)
func (c *DIContainer) Clear() {
	c.Close()
}

// Count returns the number of registered services (backward compatibility)
func (c *DIContainer) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.beans)
}
