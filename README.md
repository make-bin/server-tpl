# Go HTTP Server Template

一个基于Go语言的HTTP服务器模板项目，采用分层架构设计，支持依赖注入、配置管理、日志记录等功能。

## 技术栈

- **HTTP框架**: Gin
- **配置管理**: Viper
- **日志框架**: Logrus
- **数据库**: PostgreSQL + GORM
- **依赖注入**: Google Wire
- **验证**: go-playground/validator
- **API文档**: Swagger
- **代码检查**: golangci-lint
- **监控**: Prometheus (复用Gin HTTP服务器)

## 项目结构

```
.
├── cmd/                    # 应用程序入口点
│   └── main.go            # main函数文件入口
├── pkg/                   # 存放核心代码
│   ├── server/            # 服务器实现
│   │   └── server.go      # 服务server实现文件
│   ├── api/               # HTTP API层
│   │   ├── dto/v1/        # HTTP API body struct 实现
│   │   ├── assembler/v1/  # HTTP API body structs 和 models struct 转换实现
│   │   ├── router/        # HTTP路由实现
│   │   │   └── router.go  # HTTP router 实现
│   │   ├── middleware/    # HTTP中间件（Gin中间件）
│   │   │   ├── error_handler.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   ├── application.go # API应用层实现
│   │   └── interface.go   # API接口定义
│   ├── domain/            # 领域层 - 核心业务逻辑实现
│   │   ├── model/         # 数据模型定义
│   │   └── service/       # 核心业务接口实现
│   ├── infrastructure/    # 基础设施层 - 数据库，外部服务等Client实现
│   │   ├── datastore/     # ORM数据实现
│   │   │   ├── datastore.go # ORM抽象接口
│   │   │   └── pgsql/     # PostgreSQL实现
│   │   │       └── pgsql.go
│   │   └── middleware/    # 外部服务中间件（Redis、Kafka、Elasticsearch等）
│   │       ├── redis.go
│   │       ├── kafka.go
│   │       └── elasticsearch.go
│   ├── utils/             # 通用库工具实现
│   │   ├── container/     # 依赖注入容器
│   │   │   └── container.go
│   │   ├── config/        # 配置管理
│   │   │   └── config.go
│   │   ├── logger/        # 日志管理
│   │   │   └── logger.go
│   │   ├── errors/        # 错误处理
│   │   │   └── errors.go
│   │   └── bcode/         # 业务错误码
│   │       └── bcode.go
│   └── e2e/               # 集成测试代码
├── configs/               # 配置文件
│   └── config.yaml        # 应用配置
├── docs/                  # 文档
│   └── swagger.yaml       # Swagger API文档
├── deploy/                # 部署文件
│   ├── k8s/               # Kubernetes部署文件
│   │   └── deployment.yaml
│   └── docker-compose.yml # Docker Compose开发环境
├── .gitignore
├── Makefile               # 编译相关文件
├── go.mod                 # Go模块文件
├── go.sum                 # Go依赖校验文件
├── Dockerfile             # Docker容器化文件
├── .golangci.yml          # Go代码检查配置
└── README.md              # 项目说明文档
```

## 快速开始

### 环境要求

- Go 1.19+
- PostgreSQL 12+
- Docker (可选)

### 本地开发

1. 克隆项目
```bash
git clone <repository-url>
cd server-tpl
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
```bash
# 创建数据库
createdb server_tpl
```

4. 运行应用
```bash
# 开发模式
make dev

# 或者直接运行
go run cmd/main.go
```

### 使用Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

### 构建和部署

```bash
# 构建应用
make build

# 构建Docker镜像
make docker-build

# 运行测试
make test

# 代码检查
make lint
```

## 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| SERVER_PORT | 8080 | 服务器端口 |
| SERVER_HOST | 0.0.0.0 | 服务器地址 |
| SERVER_MODE | release | 运行模式 (debug/release) |
| DATABASE_HOST | localhost | 数据库地址 |
| DATABASE_PORT | 5432 | 数据库端口 |
| DATABASE_USER | postgres | 数据库用户名 |
| DATABASE_PASSWORD | password | 数据库密码 |
| DATABASE_DBNAME | server_tpl | 数据库名称 |
| LOG_LEVEL | info | 日志级别 |
| LOG_FORMAT | json | 日志格式 |

### 配置文件

配置文件位于 `configs/config.yaml`，支持YAML格式：

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  mode: "release"
  timeout: 30

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "server_tpl"
  sslmode: "disable"

log:
  level: "info"
  format: "json"
  output: "stdout"
```

## API文档

启动应用后，可以通过以下方式访问API文档：

- Swagger UI: http://localhost:8080/swagger/index.html
- API文档文件: `docs/swagger.yaml`

### 主要API端点

- `GET /health` - 健康检查
- `GET /api/v1/applications` - 获取应用列表
- `POST /api/v1/applications` - 创建应用
- `GET /api/v1/examples/error/:type` - 错误处理演示
- `GET /api/v1/examples/success` - 成功响应演示
- `GET /api/v1/examples/list` - 列表响应演示

## 开发指南

### 添加新的API

1. 在 `pkg/api/dto/v1/` 中定义DTO结构
2. 在 `pkg/api/assembler/v1/` 中实现转换逻辑
3. 在 `pkg/domain/model/` 中定义数据模型
4. 在 `pkg/domain/service/` 中实现业务逻辑
5. 在 `pkg/api/` 中实现API处理器
6. 在 `pkg/api/router/` 中添加路由

### 依赖注入

项目使用自定义容器进行依赖注入：

```go
type applicationService struct {
    Store datastore.DataStore `inject:"datastore"`
}
```

### 错误处理

项目提供统一的错误处理机制：

```go
import (
    "github.com/make-bin/server-tpl/pkg/utils/errors"
"github.com/make-bin/server-tpl/pkg/utils/bcode"
"github.com/make-bin/server-tpl/pkg/infrastructure/middleware"
)

// 使用业务错误码
err := bcode.ErrUserNotFound.WithDetails("User ID: 123")
middleware.ErrorResponse(c, err)

// 使用基础错误
err = errors.New(400, "Bad request")
middleware.ErrorResponse(c, err)
```

### 日志记录

使用Logrus进行结构化日志记录：

```go
import "github.com/make-bin/server-tpl/pkg/utils/logger"

logger.Info("Application started")
logger.WithField("user_id", 123).Info("User logged in")
```

### 监控

项目集成了Prometheus监控，通过复用Gin HTTP服务器提供指标端点：

#### 配置

```yaml
prometheus:
  enabled: true
  metrics_path: /metrics
```

#### 访问指标

启动应用后，通过以下URL访问Prometheus指标：

```
http://localhost:8080/metrics
```

#### 指标类型

- **HTTP指标**: 请求总数、响应时间、请求/响应大小
- **业务指标**: 业务操作计数、操作耗时、错误统计
- **系统指标**: 内存使用、CPU使用、Goroutine数量
- **数据库指标**: 连接数、查询次数、查询耗时、错误统计
- **缓存指标**: 命中率、缓存大小

#### 使用示例

```go
// 获取Prometheus中间件实例
prometheus := server.GetPrometheus()

// 记录业务操作
prometheus.RecordBusinessOperation("user_login", "success", duration)

// 记录数据库查询
prometheus.RecordDatabaseQuery("postgres", "select", duration)

// 记录缓存操作
prometheus.RecordCacheHit("redis")
```

详细文档请参考：[Prometheus监控规范](docs/prometheus-monitoring.md)

### 调试

项目集成了Go PProf性能分析工具，提供全面的性能分析和调试能力：

#### 配置

```yaml
pprof:
  enabled: true  # 开发环境启用，生产环境建议关闭
  path_prefix: /debug/pprof
```

#### 访问PProf界面

启动应用后，访问PProf主页：

```
http://localhost:8080/debug/pprof/
```

#### 使用PProf工具

```bash
# CPU分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 堆内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine分析
go tool pprof http://localhost:8080/debug/pprof/goroutine

# 阻塞分析
go tool pprof http://localhost:8080/debug/pprof/block
```

#### 编程接口

```go
import "github.com/make-bin/server-tpl/pkg/utils/pprof"

// 创建PProf管理器
config := &pprof.PProfConfig{
    Enabled:    true,
    PathPrefix: "/debug/pprof",
}
pprofManager := pprof.NewPProfManager(config)

// CPU分析
cpuFile, err := pprofManager.StartCPUProfile("cpu.prof")
if err != nil {
    return err
}
defer func() {
    pprofManager.StopCPUProfile()
    cpuFile.Close()
}()

// 堆内存分析
if err := pprofManager.WriteHeapProfile("heap.prof"); err != nil {
    return err
}
```

详细文档请参考：[PProf调试规范](docs/pprof-debugging.md)

## 测试

```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test ./pkg/e2e/...

# 生成测试覆盖率报告
make test-coverage
```

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t server-tpl .

# 运行容器
docker run -p 8080:8080 server-tpl
```

### Kubernetes部署

```bash
# 部署到Kubernetes
kubectl apply -f deploy/k8s/
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目维护者: [Your Name]
- 邮箱: [your.email@example.com]
- 项目链接: [https://github.com/your-username/server-tpl]
