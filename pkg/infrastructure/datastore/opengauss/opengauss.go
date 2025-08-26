package opengauss

import (
	"context"
	"fmt"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/utils/config"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenGauss implements DatastoreInterface using OpenGauss
type OpenGauss struct {
	db *gorm.DB
}

// New creates a new OpenGauss datastore instance
func New(cfg *config.Config) (datastore.DatastoreInterface, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		"UTC", // Default timezone
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenGauss: %w", err)
	}

	logger.Info("Connected to OpenGauss database")

	return &OpenGauss{db: db}, nil
}

// CreateApplication creates a new application
func (o *OpenGauss) CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	if err := o.db.WithContext(ctx).Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// GetApplicationByID retrieves an application by ID
func (o *OpenGauss) GetApplicationByID(ctx context.Context, id uint) (*model.Application, error) {
	var app model.Application
	if err := o.db.WithContext(ctx).First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, datastore.ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

// GetApplicationByName retrieves an application by name
func (o *OpenGauss) GetApplicationByName(ctx context.Context, name string) (*model.Application, error) {
	var app model.Application
	if err := o.db.WithContext(ctx).Where("name = ?", name).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, datastore.ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

// ListApplications retrieves a paginated list of applications
func (o *OpenGauss) ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error) {
	var apps []*model.Application
	var total int64

	// Count total records
	if err := o.db.WithContext(ctx).Model(&model.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	if err := o.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// UpdateApplication updates an existing application
func (o *OpenGauss) UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	if err := o.db.WithContext(ctx).Save(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// DeleteApplication deletes an application by ID
func (o *OpenGauss) DeleteApplication(ctx context.Context, id uint) error {
	result := o.db.WithContext(ctx).Delete(&model.Application{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return datastore.ErrNotFound
	}
	return nil
}

// Migrate runs database migrations
func (o *OpenGauss) Migrate() error {
	return o.db.AutoMigrate(&model.Application{})
}

// Close closes the database connection
func (o *OpenGauss) Close() error {
	sqlDB, err := o.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks the database connection
func (o *OpenGauss) HealthCheck() error {
	sqlDB, err := o.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
