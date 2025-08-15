package service

import (
	"context"
	"strconv"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// VariablesService 变量服务接口
type VariablesService interface {
	CreateVariable(ctx context.Context, variable *model.Variable) error
	GetVariableByID(ctx context.Context, id uint) (*model.Variable, error)
	GetVariablesByAppID(ctx context.Context, appID uint) ([]*model.Variable, error)
	UpdateVariable(ctx context.Context, variable *model.Variable) error
	DeleteVariable(ctx context.Context, id uint) error
	ListVariables(ctx context.Context, offset, limit int) ([]*model.Variable, error)
	CountVariables(ctx context.Context) (int64, error)
}

type variablesService struct {
	Store datastore.DataStore `inject:"datastore"`
}

func NewVariablesService() VariablesService {
	return &variablesService{}
}

func (s *variablesService) CreateVariable(ctx context.Context, variable *model.Variable) error {
	return s.Store.Add(ctx, variable)
}

func (s *variablesService) GetVariableByID(ctx context.Context, id uint) (*model.Variable, error) {
	variable := &model.Variable{}
	variable.BaseEntity.ID = id
	if err := s.Store.Get(ctx, variable); err != nil {
		return nil, err
	}
	return variable, nil
}

func (s *variablesService) GetVariablesByAppID(ctx context.Context, appID uint) ([]*model.Variable, error) {
	variable := &model.Variable{ApplicationID: appID}

	options := &datastore.ListOptions{
		FilterOptions: datastore.FilterOptions{
			In: []datastore.InQueryOption{
				{Key: "application_id", Values: []string{strconv.FormatUint(uint64(appID), 10)}},
			},
		},
		SortBy: []datastore.SortOption{
			{Key: "id", Order: datastore.SortOrderAscending},
		},
	}

	entities, err := s.Store.List(ctx, variable, options)
	if err != nil {
		return nil, err
	}

	variables := make([]*model.Variable, len(entities))
	for i, entity := range entities {
		variables[i] = entity.(*model.Variable)
	}

	return variables, nil
}

func (s *variablesService) UpdateVariable(ctx context.Context, variable *model.Variable) error {
	return s.Store.Put(ctx, variable)
}

func (s *variablesService) DeleteVariable(ctx context.Context, id uint) error {
	variable := &model.Variable{}
	variable.BaseEntity.ID = id
	return s.Store.Delete(ctx, variable)
}

func (s *variablesService) ListVariables(ctx context.Context, offset, limit int) ([]*model.Variable, error) {
	variable := &model.Variable{}

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

	entities, err := s.Store.List(ctx, variable, options)
	if err != nil {
		return nil, err
	}

	variables := make([]*model.Variable, len(entities))
	for i, entity := range entities {
		variables[i] = entity.(*model.Variable)
	}

	return variables, nil
}

func (s *variablesService) CountVariables(ctx context.Context) (int64, error) {
	variable := &model.Variable{}
	return s.Store.Count(ctx, variable, nil)
}
