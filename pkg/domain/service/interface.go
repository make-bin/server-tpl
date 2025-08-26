package service

import (
	"context"

	"github.com/make-bin/server-tpl/pkg/domain/model"
)

// ApplicationServiceInterface defines the interface for application service
type ApplicationServiceInterface interface {
	CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	GetApplicationByID(ctx context.Context, id uint) (*model.Application, error)
	GetApplicationByName(ctx context.Context, name string) (*model.Application, error)
	ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error)
	UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	DeleteApplication(ctx context.Context, id uint) error
}

// InitServiceBean convert service interface to bean type
func InitServiceBean() []interface{} {
	return []interface{}{
		NewApplicationServiceForDI(),
	}
}
