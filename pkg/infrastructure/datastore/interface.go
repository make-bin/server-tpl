package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/make-bin/server-tpl/pkg/domain/model"
)

// Common datastore errors
var (
	ErrNotFound          = errors.New("record not found")
	ErrDuplicateKey      = errors.New("duplicate key violation")
	ErrInvalidInput      = errors.New("invalid input")
	ErrConnectionFailed  = errors.New("database connection failed")
	ErrTransactionFailed = errors.New("transaction failed")
)

// Entity interface defines common methods for all entities
type Entity interface {
	SetCreateTime(time.Time)
	SetUpdateTime(time.Time)
	PrimaryKey() string
	TableName() string
	ShortTableName() string
	Index() map[string]interface{}
}

// Transaction interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	Add(ctx context.Context, entity Entity) error
	Put(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, entity Entity) error
	Get(ctx context.Context, entity Entity) error
}

// ListOptions defines options for list queries
type ListOptions struct {
	Page     int                    `json:"page"`
	Size     int                    `json:"size"`
	SortBy   string                 `json:"sort_by"`
	SortDesc bool                   `json:"sort_desc"`
	Filters  map[string]interface{} `json:"filters"`
}

// FilterOptions defines options for filter queries
type FilterOptions struct {
	Filters map[string]interface{} `json:"filters"`
}

// Config defines database configuration
type Config struct {
	Type            string        `json:"type"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	Database        string        `json:"database"`
	SSLMode         string        `json:"ssl_mode"`
	TimeZone        string        `json:"timezone"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// DataStore defines the unified storage interface
type DataStore interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	HealthCheck(ctx context.Context) error

	// Transaction management
	BeginTx(ctx context.Context) (Transaction, error)

	// CRUD operations
	Add(ctx context.Context, entity Entity) error
	BatchAdd(ctx context.Context, entities []Entity) error
	Put(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, entity Entity) error
	Get(ctx context.Context, entity Entity) error
	List(ctx context.Context, query Entity, options *ListOptions) ([]Entity, error)
	Count(ctx context.Context, entity Entity, options *FilterOptions) (int64, error)
	IsExist(ctx context.Context, entity Entity) (bool, error)

	// Migration operations
	Migrate(ctx context.Context, entities ...Entity) error
	ExecuteSQL(ctx context.Context, sql string, args ...interface{}) error
}

// DatastoreInterface defines the interface for data persistence (backward compatibility)
type DatastoreInterface interface {
	// Application operations
	CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	GetApplicationByID(ctx context.Context, id uint) (*model.Application, error)
	GetApplicationByName(ctx context.Context, name string) (*model.Application, error)
	ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error)
	UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	DeleteApplication(ctx context.Context, id uint) error

	// Database operations
	Migrate() error
	Close() error
	HealthCheck() error
}

// Cache interface for caching layer
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

// CacheConfig defines cache configuration
type CacheConfig struct {
	Type     string        `json:"type"` // redis, memory
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Password string        `json:"password"`
	Database int           `json:"database"`
	TTL      time.Duration `json:"ttl"`
}

// Performance monitoring interface
type Monitor interface {
	RecordQuery(operation, table string, duration time.Duration)
	RecordConnection(database string, connections int)
	RecordError(operation, table string, err error)
}

// Stats interface for getting storage statistics
type Stats interface {
	GetStats() map[string]interface{}
	GetConnectionStats() map[string]interface{}
	GetQueryStats() map[string]interface{}
}
