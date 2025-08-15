package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// MockDataStore 模拟数据存储
type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDataStore) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDataStore) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDataStore) BeginTx(ctx context.Context) (datastore.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(datastore.Transaction), args.Error(1)
}

func (m *MockDataStore) Add(ctx context.Context, entity datastore.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockDataStore) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	args := m.Called(ctx, entities)
	return args.Error(0)
}

func (m *MockDataStore) Put(ctx context.Context, entity datastore.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockDataStore) Delete(ctx context.Context, entity datastore.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockDataStore) Get(ctx context.Context, entity datastore.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockDataStore) List(ctx context.Context, query datastore.Entity, options *datastore.ListOptions) ([]datastore.Entity, error) {
	args := m.Called(ctx, query, options)
	return args.Get(0).([]datastore.Entity), args.Error(1)
}

func (m *MockDataStore) Count(ctx context.Context, entity datastore.Entity, options *datastore.FilterOptions) (int64, error) {
	args := m.Called(ctx, entity, options)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDataStore) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	args := m.Called(ctx, entity)
	return args.Bool(0), args.Error(1)
}

// 测试辅助函数
func createTestUser() *model.User {
	return &model.User{
		BaseEntity: model.BaseEntity{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Role:     "user",
		Status:   "active",
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常创建用户",
			user: createTestUser(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			user: createTestUser(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.User")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			err := service.CreateUser(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockDataStore)
		want    *model.User
		wantErr bool
	}{
		{
			name: "正常获取用户",
			id:   1,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil).
					Run(func(args mock.Arguments) {
						user := args.Get(1).(*model.User)
						*user = *createTestUser()
					})
			},
			want:    createTestUser(),
			wantErr: false,
		},
		{
			name: "用户不存在",
			id:   999,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(datastore.ErrRecordNotExist)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			got, err := service.GetUserByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Name, got.Name)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(*MockDataStore)
		want    *model.User
		wantErr bool
	}{
		{
			name:  "正常获取用户",
			email: "test@example.com",
			setup: func(mockStore *MockDataStore) {
				user := createTestUser()
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{user}, nil)
			},
			want:    createTestUser(),
			wantErr: false,
		},
		{
			name:  "用户不存在",
			email: "nonexistent@example.com",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "数据库错误",
			email: "test@example.com",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, datastore.ErrRecordNotExist)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			got, err := service.GetUserByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Name, got.Name)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常更新用户",
			user: createTestUser(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Put", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			user: createTestUser(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Put", mock.Anything, mock.AnythingOfType("*model.User")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			err := service.UpdateUser(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常删除用户",
			id:   1,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Delete", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			id:   999,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Delete", mock.Anything, mock.AnythingOfType("*model.User")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			err := service.DeleteUser(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		setup   func(*MockDataStore)
		want    []*model.User
		wantErr bool
	}{
		{
			name:   "正常获取用户列表",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				users := []datastore.Entity{
					createTestUser(),
					&model.User{
						BaseEntity: model.BaseEntity{ID: 2},
						Email:      "user2@example.com",
						Name:       "User 2",
					},
				}
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return(users, nil)
			},
			want: []*model.User{
				createTestUser(),
				&model.User{
					BaseEntity: model.BaseEntity{ID: 2},
					Email:      "user2@example.com",
					Name:       "User 2",
				},
			},
			wantErr: false,
		},
		{
			name:   "空列表",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, nil)
			},
			want:    []*model.User{},
			wantErr: false,
		},
		{
			name:   "数据库错误",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, datastore.ErrRecordNotExist)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			got, err := service.ListUsers(context.Background(), tt.offset, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				if len(got) > 0 {
					assert.Equal(t, tt.want[0].ID, got[0].ID)
					assert.Equal(t, tt.want[0].Email, got[0].Email)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserService_CountUsers(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockDataStore)
		want    int64
		wantErr bool
	}{
		{
			name: "正常获取用户数量",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Count", mock.Anything, mock.AnythingOfType("*model.User"), mock.Anything).
					Return(int64(10), nil)
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "数据库错误",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Count", mock.Anything, mock.AnythingOfType("*model.User"), mock.Anything).
					Return(int64(0), datastore.ErrRecordNotExist)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &userService{Store: mockStore}
			got, err := service.CountUsers(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// 基准测试
func BenchmarkUserService_CreateUser(b *testing.B) {
	mockStore := &MockDataStore{}
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	service := &userService{Store: mockStore}
	user := createTestUser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.ID = uint(i + 1)
		_ = service.CreateUser(context.Background(), user)
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	mockStore := &MockDataStore{}
	mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.User")).
		Return(nil).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*model.User)
			*user = *createTestUser()
		})

	service := &userService{Store: mockStore}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserByID(context.Background(), 1)
	}
}
