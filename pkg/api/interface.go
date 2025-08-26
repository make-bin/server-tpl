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

// RegisterValidationInterface register validation function
func RegisterValidationInterface(name string, fn validator.Func) {
	if registerValidationInterfaces == nil {
		registerValidationInterfaces = make(map[string]validator.Func)
	}
	registerValidationInterfaces[name] = fn
}

func GetRegisterValidationInterfaces() map[string]validator.Func {
	return registerValidationInterfaces
}

func init() {
	// Todo RegisterAPIInterface(newApplication())
}
