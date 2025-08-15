package factory

import (
	"fmt"

	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/memory"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/opengauss"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/postgresql"
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMySQL      DatabaseType = "mysql"
	DatabaseTypeOpenGauss  DatabaseType = "opengauss"
	DatabaseTypeMemory     DatabaseType = "memory"
)

// NewDataStore 创建数据存储实例
func NewDataStore(config *datastore.Config) (datastore.DataStore, error) {
	switch DatabaseType(config.Type) {
	case DatabaseTypePostgreSQL:
		return postgresql.NewPostgreSQL(config), nil
	case DatabaseTypeOpenGauss:
		return opengauss.NewOpenGauss(config), nil
	case DatabaseTypeMemory:
		return memory.NewMemory(config), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// NewDataStoreWithType 根据类型创建数据存储实例
func NewDataStoreWithType(dbType DatabaseType, config *datastore.Config) (datastore.DataStore, error) {
	config.Type = string(dbType)
	return NewDataStore(config)
}
