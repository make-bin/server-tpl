package service

import (
	"context"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// ApplicationService implements ApplicationServiceInterface
type ApplicationService struct {
	datastore datastore.DatastoreInterface
}

// applicationService 内部实现，支持依赖注入
type applicationService struct {
	Store datastore.DatastoreInterface `inject:"datastore"`
}

// NewApplicationService creates a new ApplicationService instance
func NewApplicationService(ds datastore.DatastoreInterface) ApplicationServiceInterface {
	return &ApplicationService{
		datastore: ds,
	}
}

// NewApplicationServiceForDI 创建支持依赖注入的应用服务实例
func NewApplicationServiceForDI() ApplicationServiceInterface {
	return &applicationService{}
}

// CreateApplication creates a new application
func (s *ApplicationService) CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	logger.Info("Creating application: %s", app.Name)

	// Validate domain rules
	if err := app.Validate(); err != nil {
		return nil, err
	}

	// Check if application with same name exists
	existing, err := s.datastore.GetApplicationByName(ctx, app.Name)
	if err != nil && err != datastore.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, model.NewDomainError("application with this name already exists")
	}

	// Create application
	result, err := s.datastore.CreateApplication(ctx, app)
	if err != nil {
		logger.Error("Failed to create application: %v", err)
		return nil, err
	}

	logger.Info("Application created successfully: %d", result.ID)
	return result, nil
}

// GetApplicationByID retrieves an application by ID
func (s *ApplicationService) GetApplicationByID(ctx context.Context, id uint) (*model.Application, error) {
	logger.Info("Getting application by ID: %d", id)

	app, err := s.datastore.GetApplicationByID(ctx, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		logger.Error("Failed to get application by ID: %v", err)
		return nil, err
	}

	return app, nil
}

// GetApplicationByName retrieves an application by name
func (s *ApplicationService) GetApplicationByName(ctx context.Context, name string) (*model.Application, error) {
	logger.Info("Getting application by name: %s", name)

	app, err := s.datastore.GetApplicationByName(ctx, name)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		logger.Error("Failed to get application by name: %v", err)
		return nil, err
	}

	return app, nil
}

// ListApplications retrieves a paginated list of applications
func (s *ApplicationService) ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error) {
	logger.Info("Listing applications: page=%d, pageSize=%d", page, pageSize)

	apps, total, err := s.datastore.ListApplications(ctx, page, pageSize)
	if err != nil {
		logger.Error("Failed to list applications: %v", err)
		return nil, 0, err
	}

	return apps, total, nil
}

// UpdateApplication updates an existing application
func (s *ApplicationService) UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	logger.Info("Updating application: %d", app.ID)

	// Validate domain rules
	if err := app.Validate(); err != nil {
		return nil, err
	}

	// Check if application exists
	existing, err := s.datastore.GetApplicationByID(ctx, app.ID)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		return nil, err
	}

	// Check if another application with same name exists
	if existing.Name != app.Name {
		nameExists, err := s.datastore.GetApplicationByName(ctx, app.Name)
		if err != nil && err != datastore.ErrNotFound {
			return nil, err
		}
		if nameExists != nil {
			return nil, model.NewDomainError("application with this name already exists")
		}
	}

	// Update application
	result, err := s.datastore.UpdateApplication(ctx, app)
	if err != nil {
		logger.Error("Failed to update application: %v", err)
		return nil, err
	}

	logger.Info("Application updated successfully: %d", result.ID)
	return result, nil
}

// DeleteApplication deletes an application by ID
func (s *ApplicationService) DeleteApplication(ctx context.Context, id uint) error {
	logger.Info("Deleting application: %d", id)

	// Check if application exists
	_, err := s.datastore.GetApplicationByID(ctx, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return model.ErrApplicationNotFound
		}
		return err
	}

	// Delete application
	err = s.datastore.DeleteApplication(ctx, id)
	if err != nil {
		logger.Error("Failed to delete application: %v", err)
		return err
	}

	logger.Info("Application deleted successfully: %d", id)
	return nil
}

// 为依赖注入版本实现相同的方法

// CreateApplication creates a new application (DI version)
func (s *applicationService) CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	logger.Info("Creating application: %s", app.Name)

	// Validate domain rules
	if err := app.Validate(); err != nil {
		return nil, err
	}

	// Check if application with same name exists
	existing, err := s.Store.GetApplicationByName(ctx, app.Name)
	if err != nil && err != datastore.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, model.NewDomainError("application with this name already exists")
	}

	// Create application
	result, err := s.Store.CreateApplication(ctx, app)
	if err != nil {
		logger.Error("Failed to create application: %v", err)
		return nil, err
	}

	logger.Info("Application created successfully: %d", result.ID)
	return result, nil
}

// GetApplicationByID retrieves an application by ID (DI version)
func (s *applicationService) GetApplicationByID(ctx context.Context, id uint) (*model.Application, error) {
	logger.Info("Getting application by ID: %d", id)

	app, err := s.Store.GetApplicationByID(ctx, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		logger.Error("Failed to get application by ID: %v", err)
		return nil, err
	}

	return app, nil
}

// GetApplicationByName retrieves an application by name (DI version)
func (s *applicationService) GetApplicationByName(ctx context.Context, name string) (*model.Application, error) {
	logger.Info("Getting application by name: %s", name)

	app, err := s.Store.GetApplicationByName(ctx, name)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		logger.Error("Failed to get application by name: %v", err)
		return nil, err
	}

	return app, nil
}

// ListApplications retrieves a paginated list of applications (DI version)
func (s *applicationService) ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error) {
	logger.Info("Listing applications: page=%d, pageSize=%d", page, pageSize)

	apps, total, err := s.Store.ListApplications(ctx, page, pageSize)
	if err != nil {
		logger.Error("Failed to list applications: %v", err)
		return nil, 0, err
	}

	return apps, total, nil
}

// UpdateApplication updates an existing application (DI version)
func (s *applicationService) UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	logger.Info("Updating application: %d", app.ID)

	// Validate domain rules
	if err := app.Validate(); err != nil {
		return nil, err
	}

	// Check if application exists
	existing, err := s.Store.GetApplicationByID(ctx, app.ID)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrApplicationNotFound
		}
		return nil, err
	}

	// Check if another application with same name exists
	if existing.Name != app.Name {
		nameExists, err := s.Store.GetApplicationByName(ctx, app.Name)
		if err != nil && err != datastore.ErrNotFound {
			return nil, err
		}
		if nameExists != nil {
			return nil, model.NewDomainError("application with this name already exists")
		}
	}

	// Update application
	result, err := s.Store.UpdateApplication(ctx, app)
	if err != nil {
		logger.Error("Failed to update application: %v", err)
		return nil, err
	}

	logger.Info("Application updated successfully: %d", result.ID)
	return result, nil
}

// DeleteApplication deletes an application by ID (DI version)
func (s *applicationService) DeleteApplication(ctx context.Context, id uint) error {
	logger.Info("Deleting application: %d", id)

	// Check if application exists
	_, err := s.Store.GetApplicationByID(ctx, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return model.ErrApplicationNotFound
		}
		return err
	}

	// Delete application
	err = s.Store.DeleteApplication(ctx, id)
	if err != nil {
		logger.Error("Failed to delete application: %v", err)
		return err
	}

	logger.Info("Application deleted successfully: %d", id)
	return nil
}
