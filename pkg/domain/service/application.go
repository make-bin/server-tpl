package service

import (
	"context"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// ApplicationService 应用服务接口
type ApplicationService interface {
	CreateApplication(ctx context.Context, app *model.Application) error
	GetApplicationByID(ctx context.Context, id uint) (*model.Application, error)
	GetApplicationByName(ctx context.Context, name string) (*model.Application, error)
	UpdateApplication(ctx context.Context, app *model.Application) error
	DeleteApplication(ctx context.Context, id uint) error
	ListApplications(ctx context.Context, offset, limit int) ([]*model.Application, error)
	CountApplications(ctx context.Context) (int64, error)
}

type applicationService struct {
	Store datastore.DataStore `inject:"datastore"`
}

func NewApplicationService() ApplicationService {
	return &applicationService{}
}

func (s *applicationService) CreateApplication(ctx context.Context, app *model.Application) error {
	return s.Store.Add(ctx, app)
}

func (s *applicationService) GetApplicationByID(ctx context.Context, id uint) (*model.Application, error) {
	app := &model.Application{}
	app.BaseEntity.ID = id
	if err := s.Store.Get(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *applicationService) GetApplicationByName(ctx context.Context, name string) (*model.Application, error) {
	app := &model.Application{Name: name}

	// 使用 List 方法查询
	options := &datastore.ListOptions{
		FilterOptions: datastore.FilterOptions{
			In: []datastore.InQueryOption{
				{Key: "name", Values: []string{name}},
			},
		},
		Page:     1,
		PageSize: 1,
	}

	entities, err := s.Store.List(ctx, app, options)
	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, datastore.ErrRecordNotExist
	}

	return entities[0].(*model.Application), nil
}

func (s *applicationService) UpdateApplication(ctx context.Context, app *model.Application) error {
	return s.Store.Put(ctx, app)
}

func (s *applicationService) DeleteApplication(ctx context.Context, id uint) error {
	app := &model.Application{}
	app.BaseEntity.ID = id
	return s.Store.Delete(ctx, app)
}

func (s *applicationService) ListApplications(ctx context.Context, offset, limit int) ([]*model.Application, error) {
	app := &model.Application{}

	page := (offset / limit) + 1
	if offset%limit != 0 {
		page++
	}

	options := &datastore.ListOptions{
		Page:     page,
		PageSize: limit,
		SortBy: []datastore.SortOption{
			{Key: "id", Order: datastore.SortOrderAscending},
		},
	}

	entities, err := s.Store.List(ctx, app, options)
	if err != nil {
		return nil, err
	}

	apps := make([]*model.Application, len(entities))
	for i, entity := range entities {
		apps[i] = entity.(*model.Application)
	}

	return apps, nil
}

func (s *applicationService) CountApplications(ctx context.Context) (int64, error) {
	app := &model.Application{}
	return s.Store.Count(ctx, app, nil)
}
