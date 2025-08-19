# 8. 开发最佳实践

## 8.1 命名规范
- 包名：使用小写字母，避免下划线
- 函数名：使用驼峰命名法
- 常量：使用大写字母和下划线
- 变量：使用驼峰命名法
- 接口名：以Service结尾
- 结构体：使用驼峰命名法

## 8.2 错误处理
- 使用统一的错误处理范式
- 定义业务错误码系统
- 统一的错误响应格式
- 错误分类和追踪
- 记录错误日志

## 8.3 日志规范
- 使用Logrus进行结构化日志
- 定义日志级别和格式
- 记录关键业务操作
- 错误日志包含上下文信息

## 8.4 配置管理
- 使用Viper管理配置
- 支持多环境配置
- 配置热重载
- 敏感信息使用环境变量

## 8.5 测试规范

### 8.5.1 单元测试规范

#### 8.5.1.1 测试文件命名
- 测试文件以 `_test.go` 结尾
- 测试文件与源文件在同一目录下
- 测试包名与源文件包名相同

#### 8.5.1.2 测试函数命名
```go
// 测试函数命名规范
func TestFunctionName(t *testing.T)           // 基本功能测试
func TestFunctionName_WithSpecificCase(t *testing.T)  // 特定场景测试
func TestFunctionName_ErrorCase(t *testing.T) // 错误场景测试
func TestFunctionName_EdgeCase(t *testing.T)  // 边界条件测试
```

#### 8.5.1.3 测试用例结构
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
        {
            name: "邮箱格式错误",
            input: CreateUserRequest{
                Name:  "testuser",
                Email: "invalid-email",
            },
            want:    nil,
            wantErr: true,
        },
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

#### 8.5.1.4 测试覆盖率要求
- **核心业务逻辑**: 覆盖率 ≥ 80%
- **工具函数**: 覆盖率 ≥ 90%
- **API层**: 覆盖率 ≥ 70%
- **整体项目**: 覆盖率 ≥ 75%

#### 8.5.1.5 Mock和Stub使用
```go
// 1. 使用接口进行Mock
type MockDataStore struct {
    mock.Mock
}

func (m *MockDataStore) Get(ctx context.Context, entity Entity) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

// 2. 在测试中使用Mock
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

#### 8.5.1.6 测试数据管理
```go
// 1. 使用测试辅助函数
func createTestUser() *User {
    return &User{
        ID:        "test-id",
        Name:      "test-user",
        Email:     "test@example.com",
        CreateTime: time.Now(),
        UpdateTime: time.Now(),
    }
}

// 2. 使用测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // 自动迁移
    err = db.AutoMigrate(&User{}, &Application{})
    require.NoError(t, err)
    
    return db
}
```

#### 8.5.1.7 性能测试
```go
func BenchmarkUserService_CreateUser(b *testing.B) {
    mockStore := &MockDataStore{}
    service := NewUserService(mockStore)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.CreateUser(context.Background(), CreateUserRequest{
            Name:  fmt.Sprintf("user%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        })
    }
}
```

### 8.5.2 集成测试规范

#### 8.5.2.1 测试环境配置
```go
// 使用测试配置文件
func TestMain(m *testing.M) {
    // 设置测试环境
    os.Setenv("ENV", "test")
    os.Setenv("DB_TYPE", "memory")
    
    // 运行测试
    code := m.Run()
    
    // 清理资源
    os.Exit(code)
}
```

#### 8.5.2.2 API集成测试
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

### 8.5.3 端到端测试规范

#### 8.5.3.1 完整流程测试
```go
func TestE2E_UserWorkflow(t *testing.T) {
    // 启动测试服务器
    server := startTestServer(t)
    defer server.Close()
    
    // 1. 创建用户
    userID := createUser(t, server.URL)
    
    // 2. 查询用户
    user := getUser(t, server.URL, userID)
    assert.Equal(t, "testuser", user.Name)
    
    // 3. 更新用户
    updateUser(t, server.URL, userID, "updateduser")
    
    // 4. 删除用户
    deleteUser(t, server.URL, userID)
    
    // 5. 验证删除
    assertUserNotExists(t, server.URL, userID)
}
```

### 8.5.4 测试工具和框架

#### 8.5.4.1 推荐测试框架
- **标准库**: `testing` 包
- **断言库**: `testify/assert` 和 `testify/require`
- **Mock库**: `testify/mock`
- **HTTP测试**: `httptest` 包
- **数据库测试**: 内存数据库或测试容器

#### 8.5.4.2 测试命令
```bash
# 运行所有测试
go test -mod=vendor ./...

# 运行特定包的测试
go test -mod=vendor ./pkg/domain/service

# 运行测试并显示覆盖率
go test -mod=vendor -cover ./...

# 生成覆盖率报告
go test -mod=vendor -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
go test -mod=vendor -bench=. ./...

# 运行测试并显示详细信息
go test -mod=vendor -v ./...
```

### 8.5.5 测试最佳实践

#### 8.5.5.1 测试隔离
- 每个测试用例独立运行
- 使用 `t.Parallel()` 并行执行测试
- 避免测试间的状态共享

#### 8.5.5.2 测试可读性
- 使用描述性的测试名称
- 遵循 AAA 模式（Arrange, Act, Assert）
- 添加必要的注释说明

#### 8.5.5.3 测试维护性
- 提取公共的测试辅助函数
- 使用测试数据工厂
- 定期重构测试代码

#### 8.5.5.4 测试性能
- 避免在测试中执行耗时操作
- 使用内存数据库进行测试
- 合理使用测试并行化

## 8.6 安全规范
- 输入验证和清理
- SQL注入防护
- XSS防护
- CORS配置
- 认证和授权

## 8.7 构建规范
- 使用 vendor 目录确保构建一致性
- 支持离线构建
- 多阶段 Docker 构建
- 构建产物优化
