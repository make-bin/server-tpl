# 7. 依赖注入容器实现

## 7.1 container.go 
使用简单的容器来管理依赖注入

```go
package container
import (
	"fmt"
	"time"
	"sync"
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
```

## 7.2 Service 实现示例
初始化 application service实例

```go
// 在server初始化中
if err := s.beanContainer.ProvideWithName("datastore", s.dataStore); err != nil {
	return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
}
// init domain service
if err := s.beanContainer.Provides(service.InitServiceBean()...); nil != err {
	return err
}
```

**pkg/domain/service/application.go**
```go
type applicationService struct {
	Store          datastore.DataStore          `inject:"datastore"`
}

func NewApplicationService() ApplicationService {
	return &applicationService{}
}
```

**pkg/domain/service/interface.go**
```go
package service

// InitServiceBean convert service interface to bean type
func InitServiceBean() []interface{} {
	return []interface{}{
		NewApplicationService(),
	}
}
```

## 7.3 API 初始化示例

**pkg/api/interface.go**
```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var registeredAPIInterfaces []APIInterface
var registerValidationInterfaces map[string]validator.Func

type APIInterface interface {
	InitAPIServiceRoute(rg *gin.RouterGroup)
}

// RegisterAPIInterface register APIInterface
func RegisterAPIInterface(api APIInterface) {
	registeredAPIInterfaces = append(registeredAPIInterfaces, api)
}

func GetRegisterAPIInterfaces() []APIInterface {
	return registeredAPIInterfaces
}

// InitAPI convert APIinterface to beans type
func InitAPI() []interface{} {
	var beans []interface{}
	for i := range registeredAPIInterfaces {
		beans = append(beans, registeredAPIInterfaces[i])
	}
	return beans
}
```

**pkg/api/application.go**
```go
func init() {
	RegisterAPIInterface(newApplication())
	RegisterValidationInterface("namevalidator", func(fl validator.FieldLevel) bool {
		name, ok := fl.Field().Interface().(string)
		if ok {
			if len(name) == 1 {
				return unicode.IsLower([]rune(name)[0])
			}
			matched, _ := regexp.MatchString("^[a-z][a-z0-9-]*[a-z0-9]$", name)
			return matched
		}
		return true
	})
	RegisterValidationInterface("versionvalidator", func(fl validator.FieldLevel) bool {
		version, ok := fl.Field().Interface().(string)
		if ok {
			matched, _ := regexp.MatchString("^[^\u4e00-\u9fa5]+$", version)
			return matched
		}
		return true
	})
}

type application struct {
	ApplicationService service.ApplicationService `inject:""`
	VariablesService   service.VariablesService   `inject:""`
}

func newApplication() APIInterface {
	return &application{}
}
```

**API初始化**
```go
// init route api
if err := s.beanContainer.Provides(api.InitAPI()...); nil != err {
	return err
}
if err := s.beanContainer.Populate(); err != nil {
	return fmt.Errorf("fail to populate the bean container: %w", err)
}
```
