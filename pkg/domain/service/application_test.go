package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// 测试辅助函数
func createTestApplication() *model.Application {
	return &model.Application{
		BaseEntity: model.BaseEntity{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "test-app",
		Description: "Test application",
		Version:     "1.0.0",
		Status:      "active",
		CreatedBy:   1,
	}
}

func TestApplicationService_CreateApplication(t *testing.T) {
	tests := []struct {
		name    string
		app     *model.Application
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常创建应用",
			app:  createTestApplication(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.Application")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			app:  createTestApplication(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.Application")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &applicationService{Store: mockStore}
			err := service.CreateApplication(context.Background(), tt.app)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_GetApplicationByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockDataStore)
		want    *model.Application
		wantErr bool
	}{
		{
			name: "正常获取应用",
			id:   1,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.Application")).
					Return(nil).
					Run(func(args mock.Arguments) {
						app := args.Get(1).(*model.Application)
						*app = *createTestApplication()
					})
			},
			want:    createTestApplication(),
			wantErr: false,
		},
		{
			name: "应用不存在",
			id:   999,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.Application")).
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

			service := &applicationService{Store: mockStore}
			got, err := service.GetApplicationByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Version, got.Version)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_GetApplicationByName(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		setup   func(*MockDataStore)
		want    *model.Application
		wantErr bool
	}{
		{
			name:    "正常获取应用",
			appName: "test-app",
			setup: func(mockStore *MockDataStore) {
				app := createTestApplication()
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{app}, nil)
			},
			want:    createTestApplication(),
			wantErr: false,
		},
		{
			name:    "应用不存在",
			appName: "nonexistent-app",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "数据库错误",
			appName: "test-app",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
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

			service := &applicationService{Store: mockStore}
			got, err := service.GetApplicationByName(context.Background(), tt.appName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Version, got.Version)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_UpdateApplication(t *testing.T) {
	tests := []struct {
		name    string
		app     *model.Application
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常更新应用",
			app:  createTestApplication(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Put", mock.Anything, mock.AnythingOfType("*model.Application")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			app:  createTestApplication(),
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Put", mock.Anything, mock.AnythingOfType("*model.Application")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &applicationService{Store: mockStore}
			err := service.UpdateApplication(context.Background(), tt.app)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_DeleteApplication(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockDataStore)
		wantErr bool
	}{
		{
			name: "正常删除应用",
			id:   1,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Delete", mock.Anything, mock.AnythingOfType("*model.Application")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			id:   999,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Delete", mock.Anything, mock.AnythingOfType("*model.Application")).Return(datastore.ErrRecordNotExist)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDataStore{}
			tt.setup(mockStore)

			service := &applicationService{Store: mockStore}
			err := service.DeleteApplication(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_ListApplications(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		setup   func(*MockDataStore)
		want    []*model.Application
		wantErr bool
	}{
		{
			name:   "正常获取应用列表",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				apps := []datastore.Entity{
					createTestApplication(),
					&model.Application{
						BaseEntity: model.BaseEntity{ID: 2},
						Name:       "app-2",
						Version:    "2.0.0",
						Status:     "active",
						CreatedBy:  1,
					},
				}
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
					Return(apps, nil)
			},
			want: []*model.Application{
				createTestApplication(),
				&model.Application{
					BaseEntity: model.BaseEntity{ID: 2},
					Name:       "app-2",
					Version:    "2.0.0",
					Status:     "active",
					CreatedBy:  1,
				},
			},
			wantErr: false,
		},
		{
			name:   "空列表",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
					Return([]datastore.Entity{}, nil)
			},
			want:    []*model.Application{},
			wantErr: false,
		},
		{
			name:   "数据库错误",
			offset: 0,
			limit:  10,
			setup: func(mockStore *MockDataStore) {
				mockStore.On("List", mock.Anything, mock.AnythingOfType("*model.Application"), mock.AnythingOfType("*datastore.ListOptions")).
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

			service := &applicationService{Store: mockStore}
			got, err := service.ListApplications(context.Background(), tt.offset, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				if len(got) > 0 {
					assert.Equal(t, tt.want[0].ID, got[0].ID)
					assert.Equal(t, tt.want[0].Name, got[0].Name)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestApplicationService_CountApplications(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockDataStore)
		want    int64
		wantErr bool
	}{
		{
			name: "正常获取应用数量",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Count", mock.Anything, mock.AnythingOfType("*model.Application"), mock.Anything).
					Return(int64(5), nil)
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "数据库错误",
			setup: func(mockStore *MockDataStore) {
				mockStore.On("Count", mock.Anything, mock.AnythingOfType("*model.Application"), mock.Anything).
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

			service := &applicationService{Store: mockStore}
			got, err := service.CountApplications(context.Background())

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
func BenchmarkApplicationService_CreateApplication(b *testing.B) {
	mockStore := &MockDataStore{}
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.Application")).Return(nil)

	service := &applicationService{Store: mockStore}
	app := createTestApplication()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ID = uint(i + 1)
		app.Name = fmt.Sprintf("app-%d", i)
		_ = service.CreateApplication(context.Background(), app)
	}
}

func BenchmarkApplicationService_GetApplicationByID(b *testing.B) {
	mockStore := &MockDataStore{}
	mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.Application")).
		Return(nil).
		Run(func(args mock.Arguments) {
			app := args.Get(1).(*model.Application)
			*app = *createTestApplication()
		})

	service := &applicationService{Store: mockStore}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetApplicationByID(context.Background(), 1)
	}
}
