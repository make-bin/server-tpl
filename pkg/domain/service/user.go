package service

import (
	"context"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, offset, limit int) ([]*model.User, error)
	CountUsers(ctx context.Context) (int64, error)
}

type userService struct {
	Store datastore.DataStore `inject:"datastore"`
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.Store.Add(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user := &model.User{}
	user.BaseEntity.ID = id
	if err := s.Store.Get(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{Email: email}

	// 使用 List 方法查询
	options := &datastore.ListOptions{
		FilterOptions: datastore.FilterOptions{
			In: []datastore.InQueryOption{
				{Key: "email", Values: []string{email}},
			},
		},
		Page:     1,
		PageSize: 1,
	}

	entities, err := s.Store.List(ctx, user, options)
	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, datastore.ErrRecordNotExist
	}

	return entities[0].(*model.User), nil
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.Store.Put(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	user := &model.User{}
	user.BaseEntity.ID = id
	return s.Store.Delete(ctx, user)
}

func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]*model.User, error) {
	user := &model.User{}

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

	entities, err := s.Store.List(ctx, user, options)
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, len(entities))
	for i, entity := range entities {
		users[i] = entity.(*model.User)
	}

	return users, nil
}

func (s *userService) CountUsers(ctx context.Context) (int64, error) {
	user := &model.User{}
	return s.Store.Count(ctx, user, nil)
}
