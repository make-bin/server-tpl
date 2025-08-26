# 依赖管理规范文档
# Dependency Management Documentation

## 概述 (Overview)

本项目采用 Go Modules 和 vendor 目录相结合的依赖管理策略，确保构建的一致性和离线构建能力。

This project uses a combination of Go Modules and vendor directory for dependency management to ensure build consistency and offline build capability.

## 目录结构 (Directory Structure)

```
.
├── go.mod              # Go模块定义和直接依赖
├── go.sum              # 依赖校验和文件
├── vendor/             # 依赖缓存目录（包含在版本控制中）
├── scripts/deps.sh     # 依赖管理脚本
└── Makefile           # 构建和依赖管理命令
```

## 核心组件 (Core Components)

### 1. Go Modules (`go.mod`)

定义项目的模块名称和直接依赖：

```go
module github.com/make-bin/server-tpl

go 1.21

require (
    github.com/gin-contrib/cors v1.5.0
    github.com/gin-gonic/gin v1.9.1
    github.com/go-playground/validator/v10 v10.16.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    // ... 其他依赖
)
```

### 2. Vendor 目录

- **目的**: 本地缓存所有依赖，确保构建一致性
- **版本控制**: 包含在 Git 版本控制中
- **大小管理**: 定期清理和重新生成

### 3. 依赖管理脚本 (`scripts/deps.sh`)

提供全面的依赖管理功能：

```bash
./scripts/deps.sh [命令]
```

**可用命令：**

- `init`: 初始化依赖管理
- `update`: 更新依赖版本
- `vendor`: 生成 vendor 目录
- `clean`: 清理 vendor 目录
- `verify`: 验证依赖完整性
- `check`: 检查过期依赖
- `build`: 使用 vendor 构建
- `test`: 使用 vendor 测试
- `security`: 安全性检查

## 使用指南 (Usage Guide)

### 初始化项目依赖

```bash
# 方法1: 使用 Makefile
make deps-init

# 方法2: 使用脚本
./scripts/deps.sh init

# 手动步骤
go mod tidy
go mod download
go mod vendor
go mod verify
```

### 添加新依赖

```bash
# 1. 添加依赖
go get github.com/example/package@v1.2.3

# 2. 更新 vendor
make vendor

# 3. 验证依赖
make deps-verify
```

### 更新依赖

```bash
# 更新所有依赖（谨慎使用）
make deps-update

# 更新特定依赖
go get -u github.com/example/package

# 重新生成 vendor
make vendor
```

### 构建项目

```bash
# 使用 vendor 构建（推荐）
make build

# 或者
make build-vendor

# Linux 构建
make build-linux
```

### 运行测试

```bash
# 使用 vendor 运行测试
make test

# 或者
make test-vendor

# 测试覆盖率
make test-coverage
```

## 最佳实践 (Best Practices)

### 1. 版本锁定

- 使用 `go.mod` 锁定依赖版本
- 使用 `go.sum` 验证依赖完整性
- 定期更新依赖，确保安全性

### 2. Vendor 管理

```bash
# 定期重新生成 vendor
make clean-vendor
make vendor

# 验证 vendor 完整性
make deps-verify
```

### 3. 构建一致性

- 始终使用 `-mod=vendor` 构建
- CI/CD 环境使用相同的依赖管理策略
- 团队成员使用相同的 Go 版本

### 4. 安全性

```bash
# 定期检查安全漏洞
make deps-security

# 检查过期依赖
make deps-check
```

## CI/CD 集成 (CI/CD Integration)

### GitHub Actions 示例

```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Verify dependencies
      run: |
        go mod verify
        make deps-verify
    
    - name: Run tests
      run: make test-vendor
    
    - name: Build
      run: make build-vendor
```

### Docker 构建示例

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY vendor/ vendor/

# 使用 vendor 构建
RUN go build -mod=vendor -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server /usr/local/bin/server
CMD ["server"]
```

## 故障排除 (Troubleshooting)

### 常见问题

1. **vendor 目录过大**
   ```bash
   # 清理并重新生成
   make clean-vendor
   make vendor
   ```

2. **依赖版本冲突**
   ```bash
   # 清理模块缓存
   go clean -modcache
   go mod tidy
   make vendor
   ```

3. **构建失败**
   ```bash
   # 验证依赖
   make deps-verify
   
   # 重新初始化
   make deps-init
   ```

### 调试命令

```bash
# 显示依赖信息
make deps-info

# 检查依赖状态
go list -m all

# 查看 vendor 大小
du -sh vendor/
```

## 性能优化 (Performance Optimization)

### 1. 构建缓存

```bash
# 设置构建缓存
export GOCACHE=/path/to/build/cache
export GOMODCACHE=/path/to/mod/cache
```

### 2. 并行构建

```bash
# 设置并行度
export GOMAXPROCS=4
```

### 3. Vendor 优化

- 定期清理不用的依赖
- 使用 `go mod tidy` 清理 `go.mod`
- 避免包含测试文件和文档

## 监控和维护 (Monitoring and Maintenance)

### 定期任务

1. **每周检查**
   ```bash
   make deps-check     # 检查过期依赖
   make deps-security  # 安全检查
   ```

2. **每月更新**
   ```bash
   make deps-update    # 更新依赖
   make test          # 运行测试
   ```

3. **季度审查**
   - 评估依赖必要性
   - 检查许可证合规性
   - 性能影响分析

### 指标监控

- vendor 目录大小
- 依赖数量
- 构建时间
- 安全漏洞数量

## 参考资料 (References)

- [Go Modules Reference](https://golang.org/ref/mod)
- [Vendor Directory](https://golang.org/cmd/go/#hdr-Vendor_Directories)
- [Go Security](https://golang.org/security/)
- [Best Practices for Go Modules](https://blog.golang.org/using-go-modules)
