# 1. 项目结构

这是一个基于Go语言的HTTP服务器项目，采用分层架构设计：

## 1.1 目录结构
```
.
├── cmd/                    # 应用程序入口点
│   ├── main.go            # 主程序入口
│   ├── test_storage/      # 存储层测试程序
│   │   └── main.go        # 存储功能测试
│   └── simple_server/     # 简化服务器程序
│       └── main.go        # 基础服务器实现
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
│   │   │   ├── base.go    # 基础实体接口和结构
│   │   │   ├── user.go    # 用户模型
│   │   │   ├── application.go # 应用模型
│   │   │   └── variable.go # 变量模型
│   │   └── service/       # 核心业务接口实现
│   │       ├── interface.go # 服务接口定义
│   │       ├── user.go    # 用户服务实现
│   │       ├── application.go # 应用服务实现
│   │       └── variables.go # 变量服务实现
│   ├── infrastructure/    # 基础设施层 - 数据库，外部服务等Client实现
│   │   ├── datastore/     # ORM数据实现
│   │   │   ├── interface.go # ORM抽象接口
│   │   │   ├── factory/   # 存储工厂
│   │   │   │   └── simple_factory.go # 简单工厂实现
│   │   │   ├── postgresql/ # PostgreSQL实现
│   │   │   │   └── postgresql.go
│   │   │   ├── opengauss/ # OpenGauss实现
│   │   │   │   └── opengauss.go
│   │   │   ├── memory/    # 内存存储实现
│   │   │   │   └── memory.go
│   │   │   └── README.md  # 存储层文档
│   │   └── middleware/    # 外部服务中间件（Redis、Prometheus等）
│   │       ├── redis.go
│   │       ├── prometheus.go
│   │       └── prometheus_example.go
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
├── docs/                  # 文档（包含swagger api文件）
├── deploy/                # 部署文件，通常基于kubernetes部署应用
├── vendor/                # 依赖包缓存目录（本地缓存）
├── .gitignore
├── Makefile               # 编译相关文件
├── go.mod                 # Go模块文件
├── go.sum                 # Go依赖校验文件
├── Dockerfile             # Docker容器化文件
├── .golangci.yml          # Go代码检查配置
├── REFACTOR_SUMMARY.md    # 重构总结文档
├── RUNNING_SUMMARY.md     # 运行总结文档
└── README.md              # 项目说明文档
```
