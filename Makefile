# Makefile for Go HTTP Server Template

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build parameters
BINARY_NAME=server-tpl
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/main.go

# Docker parameters
DOCKER_IMAGE=server-tpl
DOCKER_TAG=latest

# Default target
all: test build

# Build the binary
build:
	$(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -mod=vendor -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
#	rm -rf vendor/

# Run tests
test:
	$(GOTEST) -mod=vendor -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -mod=vendor -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Update dependencies
deps-update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Initialize dependency management with vendor
deps-init:
	./scripts/deps.sh init

# Generate vendor directory
vendor:
	$(GOMOD) vendor

# Build with vendor
build-vendor:
	$(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(MAIN_PATH)

# Test with vendor
test-vendor:
	$(GOTEST) -mod=vendor -v ./...

# Clean vendor directory
clean-vendor:
	rm -rf vendor/

# Verify dependencies
deps-verify:
	$(GOMOD) verify
	./scripts/deps.sh verify

# Check for outdated dependencies
deps-check:
	./scripts/deps.sh check

# Dependency security check
deps-security:
	./scripts/deps.sh security

# Show dependency information
deps-info:
	./scripts/deps.sh info

# Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	$(GOLINT) run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Security check
security:
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	gosec ./...

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
docker-compose-up:
	docker-compose up -d

# Docker compose down
docker-compose-down:
	docker-compose down

# Generate swagger documentation
swagger:
	swag init -g cmd/main.go -o docs/

# Initialize project
init:
	$(GOMOD) init github.com/make-bin/server-tpl
	$(GOMOD) tidy

# Development setup
dev-setup: install-tools deps
	$(GOLINT) --version
	$(GOCMD) version

# CI/CD targets
ci: deps lint vet security test

# Production build
prod-build: clean ci build-linux

# Help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests and build"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  run           - Build and run the application"
	@echo "  deps          - Download dependencies"
	@echo "  deps-update   - Update dependencies"
	@echo "  deps-init     - Initialize dependency management with vendor"
	@echo "  vendor        - Generate vendor directory"
	@echo "  build-vendor  - Build with vendor"
	@echo "  test-vendor   - Test with vendor"
	@echo "  clean-vendor  - Clean vendor directory"
	@echo "  deps-verify   - Verify dependencies"
	@echo "  deps-check    - Check for outdated dependencies"
	@echo "  deps-security - Dependency security check"
	@echo "  deps-info     - Show dependency information"
	@echo "  install-tools - Install development tools"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  security      - Run security check"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-compose-up   - Start with docker-compose"
	@echo "  docker-compose-down - Stop docker-compose"
	@echo "  swagger       - Generate swagger documentation"
	@echo "  init          - Initialize Go module"
	@echo "  dev-setup     - Setup development environment"
	@echo "  ci            - Run CI pipeline"
	@echo "  prod-build    - Production build"
	@echo "  help          - Show this help"

.PHONY: all build build-linux clean test test-coverage run deps deps-update deps-init vendor build-vendor test-vendor clean-vendor deps-verify deps-check deps-security deps-info install-tools lint fmt vet security docker-build docker-run docker-compose-up docker-compose-down swagger init dev-setup ci prod-build help
