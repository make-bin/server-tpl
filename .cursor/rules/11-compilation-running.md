# 11. 编译和运行

## 11.1 开发环境构建
```bash
# 使用 vendor 目录构建
go build -mod=vendor ./cmd/simple_server
go build -mod=vendor ./cmd/test_storage

# 运行测试
go test -mod=vendor ./...
```

## 11.2 生产环境构建
```bash
# 更新依赖
go mod tidy
go mod vendor

# 构建生产版本
go build -mod=vendor -ldflags="-s -w" ./cmd/simple_server
```

## 11.3 Docker 构建
```dockerfile
# 使用 vendor 目录进行多阶段构建
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o server ./cmd/simple_server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```
