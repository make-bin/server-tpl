# Go HTTP Server Makefile
# 支持 vendor 目录缓存的构建系统

# 变量定义
BINARY_NAME=server-tpl
MAIN_PATH=./cmd/simple_server
TEST_STORAGE_PATH=./cmd/test_storage
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -s -w"

# 默认目标
.PHONY: all
all: clean vendor build

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@rm -f test_storage
	@rm -f simple_server
	@rm -f *_vendor
	@go clean -cache -testcache

# 初始化 vendor 目录
.PHONY: vendor
vendor:
	@echo "Initializing vendor directory..."
	@go mod tidy
	@go mod vendor
	@echo "Vendor directory initialized successfully"

# 更新依赖
.PHONY: deps
deps:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@go mod vendor
	@echo "Dependencies updated successfully"

# 使用 vendor 目录构建
.PHONY: build
build: vendor
	@echo "Building with vendor directory..."
	@go build -mod=vendor ${LDFLAGS} -o ${BINARY_NAME} ${MAIN_PATH}
	@go build -mod=vendor -o test_storage ${TEST_STORAGE_PATH}
	@echo "Build completed successfully"

# 构建测试程序
.PHONY: build-test
build-test: vendor
	@echo "Building test programs..."
	@go build -mod=vendor -o test_storage_vendor ${TEST_STORAGE_PATH}
	@go build -mod=vendor -o simple_server_vendor ${MAIN_PATH}
	@echo "Test programs built successfully"

# 运行测试
.PHONY: test
test: vendor
	@echo "Running tests..."
	@go test -mod=vendor -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage: vendor
	@echo "Running tests with coverage..."
	@go test -mod=vendor -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 运行单元测试
.PHONY: test-unit
test-unit: vendor
	@echo "Running unit tests..."
	@go test -mod=vendor -v ./pkg/domain/service/... ./pkg/utils/...

# 运行基准测试
.PHONY: test-bench
test-bench: vendor
	@echo "Running benchmark tests..."
	@go test -mod=vendor -bench=. ./pkg/domain/service/... ./pkg/utils/...

# 运行完整测试套件
.PHONY: test-full
test-full: vendor
	@echo "Running full test suite..."
	@./scripts/test.sh -v -c -b

# 运行测试脚本
.PHONY: test-script
test-script:
	@echo "Running test script..."
	@./scripts/test.sh "$(ARGS)"

# 运行存储层测试
.PHONY: test-storage
test-storage: build-test
	@echo "Running storage layer tests..."
	@./test_storage_vendor

# 运行服务器测试
.PHONY: test-server
test-server: build-test
	@echo "Starting server for testing..."
	@./simple_server_vendor &
	@sleep 3
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8080/health
	@echo ""
	@echo "Testing storage endpoint..."
	@curl -s http://localhost:8080/test
	@echo ""
	@pkill -f simple_server_vendor || true
	@echo "Server test completed"

# 运行 Prometheus 测试
.PHONY: test-prometheus
test-prometheus: build-test
	@echo "Building Prometheus test program..."
	@go build -mod=vendor -o test_prometheus cmd/test_prometheus/main.go
	@echo "Starting Prometheus test server..."
	@./test_prometheus &
	@sleep 3
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8080/health
	@echo ""
	@echo "Testing business operations..."
	@curl -s http://localhost:8080/api/users
	@echo ""
	@echo "Testing cache operations..."
	@curl -s http://localhost:8080/api/cache-test?key=cached
	@echo ""
	@echo "Testing Prometheus metrics endpoint..."
	@curl -s http://localhost:9090/metrics | head -10
	@echo ""
	@pkill -f test_prometheus || true
	@echo "Prometheus test completed"

# 代码检查
.PHONY: lint
lint: vendor
	@echo "Running code linting..."
	@golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# 生成 API 文档
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@swag init -g ${MAIN_PATH}/main.go

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${BINARY_NAME}:${VERSION} .
	@docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest

# Docker 运行
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 ${BINARY_NAME}:latest

# 开发模式运行
.PHONY: dev
dev: build
	@echo "Starting development server..."
	@./${BINARY_NAME}

# 生产模式构建
.PHONY: prod
prod: clean vendor
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -mod=vendor ${LDFLAGS} -a -installsuffix cgo -o ${BINARY_NAME} ${MAIN_PATH}
	@echo "Production build completed"

# 安装依赖工具
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Development tools installed successfully"

# 验证构建
.PHONY: verify
verify: build test-storage test-server
	@echo "Build verification completed successfully"

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, vendor, and build"
	@echo "  clean         - Clean build artifacts"
	@echo "  vendor        - Initialize vendor directory"
	@echo "  deps          - Update dependencies"
	@echo "  build         - Build with vendor directory"
	@echo "  build-test    - Build test programs"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-bench    - Run benchmark tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-full     - Run full test suite with script"
	@echo "  test-script   - Run test script with custom args"
	@echo "  test-storage  - Run storage layer tests"
	@echo "  test-server   - Run server tests"
	@echo "  lint          - Run code linting"
	@echo "  fmt           - Format code"
	@echo "  docs          - Generate API documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  dev           - Run in development mode"
	@echo "  prod          - Build for production"
	@echo "  install-tools - Install development tools"
	@echo "  verify        - Verify build and tests"
	@echo "  help          - Show this help message"
