# 测试规范文档

## 概述

本文档描述了项目的测试策略、规范和最佳实践。项目采用分层测试架构，包括单元测试、集成测试和端到端测试。

## 测试架构

### 测试分层

1. **单元测试 (Unit Tests)**
   - 测试单个函数或方法
   - 使用Mock隔离依赖
   - 快速执行，高覆盖率

2. **集成测试 (Integration Tests)**
   - 测试组件间的交互
   - 使用测试数据库
   - 验证API接口

3. **端到端测试 (E2E Tests)**
   - 测试完整业务流程
   - 使用真实环境配置
   - 验证系统整体功能

## 测试工具和框架

### 核心工具

- **Go Testing**: 标准测试框架
- **Testify**: 断言和Mock库
- **httptest**: HTTP测试工具
- **golangci-lint**: 代码质量检查

### 依赖管理

```bash
# 安装测试依赖
go get -u github.com/stretchr/testify

# 使用vendor目录
go mod vendor
```

## 测试规范

### 文件命名

- 测试文件以 `_test.go` 结尾
- 测试文件与源文件在同一目录
- 测试包名与源文件包名相同

### 函数命名

```go
// 基本功能测试
func TestFunctionName(t *testing.T)

// 特定场景测试
func TestFunctionName_WithSpecificCase(t *testing.T)

// 错误场景测试
func TestFunctionName_ErrorCase(t *testing.T)

// 边界条件测试
func TestFunctionName_EdgeCase(t *testing.T)
```

### 测试结构

```go
func TestUserService_CreateUser(t *testing.T) {
    // 1. 准备测试数据
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "正常创建用户",
            input: CreateUserRequest{
                Name:  "testuser",
                Email: "test@example.com",
            },
            want: &User{
                Name:  "testuser",
                Email: "test@example.com",
            },
            wantErr: false,
        },
        // 更多测试用例...
    }

    // 2. 执行测试用例
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置测试环境
            mockStore := &MockDataStore{}
            service := NewUserService(mockStore)
            
            // 执行被测试的方法
            got, err := service.CreateUser(context.Background(), tt.input)
            
            // 验证结果
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Mock使用规范

### Mock数据存储

```go
// MockDataStore 模拟数据存储
type MockDataStore struct {
    mock.Mock
}

func (m *MockDataStore) Get(ctx context.Context, entity Entity) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

// 在测试中使用
func TestUserService_GetUser(t *testing.T) {
    mockStore := &MockDataStore{}
    user := &User{ID: "1", Name: "test"}
    
    mockStore.On("Get", mock.Anything, mock.AnythingOfType("*model.User")).
        Return(nil).
        Run(func(args mock.Arguments) {
            arg := args.Get(1).(*User)
            *arg = *user
        })
    
    service := NewUserService(mockStore)
    result, err := service.GetUser(context.Background(), "1")
    
    assert.NoError(t, err)
    assert.Equal(t, user, result)
    mockStore.AssertExpectations(t)
}
```

## 测试覆盖率要求

### 覆盖率目标

- **核心业务逻辑**: ≥ 80%
- **工具函数**: ≥ 90%
- **API层**: ≥ 70%
- **整体项目**: ≥ 75%

### 覆盖率检查

```bash
# 生成覆盖率报告
go test -mod=vendor -coverprofile=coverage.out ./...

# 查看HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率摘要
go tool cover -func=coverage.out
```

## 基准测试

### 基准测试规范

```go
func BenchmarkUserService_CreateUser(b *testing.B) {
    mockStore := &MockDataStore{}
    mockStore.On("Add", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

    service := NewUserService(mockStore)
    user := createTestUser()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user.ID = uint(i + 1)
        _ = service.CreateUser(context.Background(), user)
    }
}
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -mod=vendor -bench=. ./...

# 运行特定包的基准测试
go test -mod=vendor -bench=. ./pkg/domain/service

# 生成基准测试报告
go test -mod=vendor -bench=. -benchmem ./...
```

## 测试命令

### 基本测试命令

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行基准测试
make test-bench

# 运行测试并生成覆盖率
make test-coverage

# 运行完整测试套件
make test-full
```

### 测试脚本

```bash
# 运行测试脚本
./scripts/test.sh

# 带详细输出
./scripts/test.sh -v

# 生成覆盖率报告
./scripts/test.sh -c

# 运行基准测试
./scripts/test.sh -b

# 运行所有测试
./scripts/test.sh -a
```

## 集成测试

### API集成测试

```go
func TestUserAPI_CreateUser(t *testing.T) {
    // 设置测试服务器
    router := gin.New()
    api := NewUserAPI()
    api.RegisterRoutes(router.Group("/api/v1"))
    
    // 创建测试请求
    reqBody := `{"name":"testuser","email":"test@example.com"}`
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    // 记录响应
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "testuser", response["name"])
}
```

### 数据库集成测试

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // 自动迁移
    err = db.AutoMigrate(&User{}, &Application{})
    require.NoError(t, err)
    
    return db
}

func TestUserService_Integration(t *testing.T) {
    db := setupTestDB(t)
    store := NewPostgreSQLStore(db)
    service := NewUserService(store)
    
    user := &User{
        Name:  "testuser",
        Email: "test@example.com",
    }
    
    // 测试创建用户
    err := service.CreateUser(context.Background(), user)
    assert.NoError(t, err)
    
    // 测试查询用户
    found, err := service.GetUserByEmail(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.Name)
}
```

## 测试最佳实践

### 测试隔离

- 每个测试用例独立运行
- 使用 `t.Parallel()` 并行执行测试
- 避免测试间的状态共享

### 测试可读性

- 使用描述性的测试名称
- 遵循 AAA 模式（Arrange, Act, Assert）
- 添加必要的注释说明

### 测试维护性

- 提取公共的测试辅助函数
- 使用测试数据工厂
- 定期重构测试代码

### 测试性能

- 避免在测试中执行耗时操作
- 使用内存数据库进行测试
- 合理使用测试并行化

## 持续集成

### CI/CD配置

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.23'
      - run: go mod download
      - run: go mod vendor
      - run: make test-full
      - run: make lint
```

### 测试报告

- 测试结果报告
- 覆盖率报告
- 性能基准报告
- 代码质量报告

## 故障排查

### 常见问题

1. **测试失败**
   - 检查Mock设置
   - 验证测试数据
   - 查看错误日志

2. **覆盖率不足**
   - 识别未覆盖的代码路径
   - 添加边界条件测试
   - 检查错误处理分支

3. **测试性能问题**
   - 优化测试数据准备
   - 使用并行测试
   - 减少外部依赖

### 调试技巧

```bash
# 运行单个测试
go test -mod=vendor -v -run TestUserService_CreateUser

# 运行测试并显示详细输出
go test -mod=vendor -v -count=1 ./...

# 生成测试覆盖率报告
go test -mod=vendor -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 总结

遵循本测试规范可以确保：

1. **代码质量**: 通过高覆盖率测试保证代码质量
2. **可维护性**: 通过良好的测试结构提高可维护性
3. **可靠性**: 通过自动化测试提高系统可靠性
4. **开发效率**: 通过快速反馈提高开发效率

定期审查和更新测试规范，确保测试策略与项目发展保持一致。
