package factory

import (
	"fmt"

	"github.com/make-bin/server-tpl/pkg/infrastructure/cache"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/memory"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/opengauss"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/postgresql"
	"github.com/make-bin/server-tpl/pkg/utils/config"
)

// DatastoreType represents the type of datastore
type DatastoreType string

const (
	PostgreSQL DatastoreType = "postgresql"
	OpenGauss  DatastoreType = "opengauss"
	Memory     DatastoreType = "memory"
)

// SimpleFactory is a factory for creating datastore instances
type SimpleFactory struct{}

// NewSimpleFactory creates a new SimpleFactory instance
func NewSimpleFactory() *SimpleFactory {
	return &SimpleFactory{}
}

// CreateDatastore creates a datastore instance based on the configuration (backward compatibility)
func (f *SimpleFactory) CreateDatastore(cfg *config.Config) (datastore.DatastoreInterface, error) {
	switch DatastoreType(cfg.Database.Type) {
	case PostgreSQL:
		return postgresql.New(cfg)
	case OpenGauss:
		return opengauss.New(cfg)
	case Memory:
		return memory.New()
	default:
		return nil, fmt.Errorf("unsupported datastore type: %s", cfg.Database.Type)
	}
}

// CreateDataStore creates a new DataStore instance with monitoring
func (f *SimpleFactory) CreateDataStore(cfg *config.Config) (datastore.DatastoreInterface, error) {
	var store datastore.DatastoreInterface
	var err error

	switch DatastoreType(cfg.Database.Type) {
	case PostgreSQL:
		store, err = postgresql.New(cfg)
	case OpenGauss:
		store, err = opengauss.New(cfg)
	case Memory:
		store, err = memory.New()
	default:
		return nil, fmt.Errorf("unsupported datastore type: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, err
	}

	// Return the store directly (monitoring integration to be implemented later)
	return store, nil
}

// CreateMonitoredDataStore creates a new DataStore instance with monitoring wrapper
func (f *SimpleFactory) CreateMonitoredDataStore(cfg *config.Config) (datastore.DatastoreInterface, error) {
	// First create the basic datastore
	store, err := f.CreateDataStore(cfg)
	if err != nil {
		return nil, err
	}

	// For now, just return the store as-is
	// TODO: Implement monitoring wrapper
	return store, nil
}

// CreateCache creates a cache instance based on configuration
func (f *SimpleFactory) CreateCache(cfg *config.Config) (datastore.Cache, error) {
	cacheConfig := &datastore.CacheConfig{
		Type:     "memory", // Default to memory cache
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Database,
		TTL:      cfg.Redis.DialTimeout, // Use dial timeout as default TTL
	}

	// For now, always create memory cache
	// In future, this could create Redis cache based on config
	return cache.NewMemoryCache(cacheConfig), nil
}

// DataStoreFactory provides factory methods for data store creation
type DataStoreFactory interface {
	CreateDataStore(cfg *config.Config) (datastore.DatastoreInterface, error)
	CreateCache(cfg *config.Config) (datastore.Cache, error)
	CreateMonitoredDataStore(cfg *config.Config) (datastore.DatastoreInterface, error)
}

// NewDataStoreFactory creates a new data store factory
func NewDataStoreFactory() DataStoreFactory {
	return &SimpleFactory{}
}

// Note: CreateMonitoredDataStore is already defined above

// NewDataStore creates a new DataStore instance (convenience function)
// TODO: Fix interface compatibility issues
// func NewDataStore(cfg *config.Config) (datastore.DataStore, error) {
// 	factory := NewDataStoreFactory()
// 	return factory.CreateDataStore(cfg)
// }

// ConfigFromAppConfig converts app config to datastore config
func ConfigFromAppConfig(appConfig *config.Config) *datastore.Config {
	return &datastore.Config{
		Type:            appConfig.Database.Type,
		Host:            appConfig.Database.Host,
		Port:            appConfig.Database.Port,
		User:            appConfig.Database.User,
		Password:        appConfig.Database.Password,
		Database:        appConfig.Database.Database,
		SSLMode:         appConfig.Database.SSLMode,
		TimeZone:        "UTC",
		MaxOpenConns:    appConfig.Database.MaxOpenConns,
		MaxIdleConns:    appConfig.Database.MaxIdleConns,
		ConnMaxLifetime: appConfig.Database.ConnMaxLifetime,
	}
}

// CacheConfigFromAppConfig converts app config to cache config
func CacheConfigFromAppConfig(appConfig *config.Config) *datastore.CacheConfig {
	return &datastore.CacheConfig{
		Type:     "memory", // Default
		Host:     appConfig.Redis.Host,
		Port:     appConfig.Redis.Port,
		Password: appConfig.Redis.Password,
		Database: appConfig.Redis.Database,
		TTL:      appConfig.Redis.DialTimeout,
	}
}
