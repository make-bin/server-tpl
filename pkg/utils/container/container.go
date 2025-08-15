package container

import (
	"fmt"
	"sync"
	"time"
)

type Container struct {
	beans map[string]interface{}
	mu    sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		beans: make(map[string]interface{}),
	}
}

// Provides 提供多个bean
func (c *Container) Provides(beans ...interface{}) error {
	for _, bean := range beans {
		if err := c.ProvideWithName("", bean); err != nil {
			return err
		}
	}
	return nil
}

// ProvideWithName 提供带名称的bean
func (c *Container) ProvideWithName(name string, bean interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "" {
		name = fmt.Sprintf("%T", bean)
	}

	c.beans[name] = bean
	return nil
}

// Get 获取bean
func (c *Container) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bean, exists := c.beans[name]
	return bean, exists
}

// Populate 填充依赖字段
func (c *Container) Populate() error {
	start := time.Now()
	defer func() {
		fmt.Printf("[INFO]populate the bean container take time %s\n", time.Since(start))
	}()

	// 这里可以实现更复杂的依赖注入逻辑
	// 目前只是简单的注册
	return nil
}
