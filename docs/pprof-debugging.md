# Go PProf 调试规范

## 概述

本项目集成了Go语言的PProf性能分析工具，提供全面的性能分析和调试能力。PProf通过HTTP端点提供实时性能数据，支持CPU、内存、Goroutine、阻塞等多种分析类型。

## 架构设计

### 设计原则

1. **安全性优先**：生产环境默认关闭，仅在需要时启用
2. **易于使用**：提供简洁的API接口和配置
3. **全面覆盖**：支持所有PProf分析类型
4. **集成友好**：与现有HTTP服务器无缝集成

### 技术实现

- **HTTP端点**：通过Gin路由提供PProf端点
- **配置管理**：使用Viper管理PProf配置
- **工具封装**：提供PProfManager简化使用
- **安全控制**：支持启用/禁用控制

## 配置说明

### 配置文件

```yaml
# configs/config.yaml
pprof:
  enabled: false  # 生产环境建议关闭
  path_prefix: /debug/pprof
```

### 环境变量

```bash
PPROF_ENABLED=false
PPROF_PATH_PREFIX=/debug/pprof
```

## 功能特性

### 1. HTTP端点

启用PProf后，以下端点将可用：

- `/debug/pprof/` - PProf主页
- `/debug/pprof/cmdline` - 命令行参数
- `/debug/pprof/profile` - CPU分析（30秒采样）
- `/debug/pprof/symbol` - 符号表
- `/debug/pprof/trace` - 执行跟踪
- `/debug/pprof/allocs` - 内存分配分析
- `/debug/pprof/block` - 阻塞分析
- `/debug/pprof/goroutine` - Goroutine分析
- `/debug/pprof/heap` - 堆内存分析
- `/debug/pprof/mutex` - 互斥锁分析
- `/debug/pprof/threadcreate` - 线程创建分析

### 2. 分析类型

#### CPU分析
- 分析CPU使用情况
- 识别热点函数
- 优化计算密集型操作

#### 内存分析
- 分析内存分配模式
- 检测内存泄漏
- 优化内存使用

#### Goroutine分析
- 分析Goroutine数量
- 检测Goroutine泄漏
- 优化并发性能

#### 阻塞分析
- 分析程序阻塞情况
- 识别性能瓶颈
- 优化I/O操作

#### 互斥锁分析
- 分析锁竞争情况
- 检测死锁风险
- 优化并发安全

## 使用指南

### 1. 启用PProf

在配置文件中启用PProf：

```yaml
pprof:
  enabled: true
  path_prefix: /debug/pprof
```

### 2. 访问PProf界面

启动应用后，访问PProf主页：

```
http://localhost:8080/debug/pprof/
```

### 3. 使用PProf工具

#### 命令行工具

```bash
# CPU分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 堆内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine分析
go tool pprof http://localhost:8080/debug/pprof/goroutine

# 阻塞分析
go tool pprof http://localhost:8080/debug/pprof/block

# 互斥锁分析
go tool pprof http://localhost:8080/debug/pprof/mutex
```

#### 生成火焰图

```bash
# 安装火焰图工具
go install github.com/google/pprof@latest

# 生成火焰图
pprof -http=:8081 http://localhost:8080/debug/pprof/profile
```

### 4. 编程接口

#### 基本使用

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

// 获取运行时统计
stats := pprofManager.GetRuntimeStats()
```

#### 高级功能

```go
// 启用阻塞分析
pprofManager.EnableBlockProfiling(1)

// 启用互斥锁分析
pprofManager.EnableMutexProfiling(1)

// 定期性能分析
ctx := context.Background()
pprofManager.StartPeriodicProfiling(ctx, 5*time.Minute, "./profiles")

// 生成完整分析报告
if err := pprofManager.GenerateFullProfile("./profiles"); err != nil {
    return err
}
```

## 性能分析场景

### 1. CPU性能分析

#### 场景
- 应用响应缓慢
- CPU使用率过高
- 计算密集型操作优化

#### 步骤
1. 访问 `/debug/pprof/profile`
2. 使用 `go tool pprof` 分析
3. 查看热点函数
4. 优化算法或数据结构

#### 示例
```bash
# 获取CPU分析数据
curl -o cpu.prof http://localhost:8080/debug/pprof/profile

# 分析CPU使用情况
go tool pprof cpu.prof

# 在pprof交互界面中
(pprof) top
(pprof) list functionName
(pprof) web
```

### 2. 内存泄漏分析

#### 场景
- 内存使用持续增长
- 应用频繁GC
- 内存不足错误

#### 步骤
1. 访问 `/debug/pprof/heap`
2. 分析内存分配模式
3. 识别泄漏对象
4. 修复内存管理问题

#### 示例
```bash
# 获取堆内存分析
curl -o heap.prof http://localhost:8080/debug/pprof/heap

# 分析内存使用
go tool pprof heap.prof

# 查看内存分配
(pprof) top
(pprof) list allocFunction
```

### 3. Goroutine泄漏分析

#### 场景
- Goroutine数量异常增长
- 应用响应变慢
- 资源耗尽

#### 步骤
1. 访问 `/debug/pprof/goroutine`
2. 分析Goroutine状态
3. 识别泄漏原因
4. 修复并发问题

#### 示例
```bash
# 获取Goroutine分析
curl -o goroutine.prof http://localhost:8080/debug/pprof/goroutine

# 分析Goroutine
go tool pprof goroutine.prof

# 查看Goroutine状态
(pprof) top
(pprof) traces
```

### 4. 阻塞分析

#### 场景
- 应用响应延迟
- 并发性能差
- 锁竞争严重

#### 步骤
1. 启用阻塞分析
2. 访问 `/debug/pprof/block`
3. 分析阻塞原因
4. 优化并发逻辑

#### 示例
```go
// 启用阻塞分析
pprofManager.EnableBlockProfiling(1)

// 获取阻塞分析
curl -o block.prof http://localhost:8080/debug/pprof/block

// 分析阻塞情况
go tool pprof block.prof
```

## 最佳实践

### 1. 安全考虑

- **生产环境**：默认关闭PProf，仅在需要时临时启用
- **访问控制**：考虑添加认证和授权
- **网络隔离**：限制PProf端点的网络访问
- **资源限制**：避免长时间运行CPU分析

### 2. 性能考虑

- **采样频率**：合理设置阻塞和互斥锁采样频率
- **文件大小**：定期清理分析文件
- **存储空间**：监控分析文件占用的磁盘空间
- **网络带宽**：避免频繁的远程分析

### 3. 分析策略

- **问题导向**：根据具体问题选择分析类型
- **渐进分析**：从概览到细节逐步分析
- **对比分析**：对比不同时间点的分析结果
- **持续监控**：建立性能基准和监控

### 4. 工具集成

- **CI/CD集成**：在测试环境中启用PProf
- **监控集成**：与Prometheus监控系统集成
- **日志集成**：记录分析操作和结果
- **告警集成**：设置性能告警阈值

## 故障排查

### 1. PProf端点不可访问

**问题**：无法访问PProf端点

**排查步骤**：
1. 检查配置中的 `enabled` 设置
2. 确认应用正常启动
3. 检查防火墙设置
4. 验证路径前缀配置

**解决方案**：
```yaml
pprof:
  enabled: true
  path_prefix: /debug/pprof
```

### 2. 分析文件过大

**问题**：分析文件占用过多磁盘空间

**排查步骤**：
1. 检查分析频率设置
2. 确认文件清理策略
3. 监控磁盘使用情况

**解决方案**：
```go
// 降低采样频率
pprofManager.EnableBlockProfiling(1000) // 每1000纳秒采样一次

// 定期清理文件
// 实现文件清理逻辑
```

### 3. 性能影响

**问题**：PProf影响应用性能

**排查步骤**：
1. 检查采样频率设置
2. 确认分析持续时间
3. 监控系统资源使用

**解决方案**：
```go
// 仅在需要时启用
if debugMode {
    pprofManager.EnableBlockProfiling(1)
}

// 限制分析时间
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## 监控和告警

### 1. 性能指标

- **分析文件大小**：监控分析文件占用的磁盘空间
- **分析频率**：监控分析操作的频率
- **响应时间**：监控PProf端点的响应时间
- **错误率**：监控分析操作的错误率

### 2. 告警规则

```yaml
# Prometheus告警规则
groups:
  - name: pprof
    rules:
      - alert: PProfFileSizeTooLarge
        expr: pprof_file_size_bytes > 1000000000  # 1GB
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "PProf file size too large"
          description: "PProf analysis file size is {{ $value }} bytes"

      - alert: PProfEndpointDown
        expr: up{job="go-http-server"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "PProf endpoint is down"
          description: "Cannot access PProf endpoints"
```

### 3. Grafana仪表板

创建PProf监控仪表板，包括：
- 分析文件大小趋势
- 分析操作频率
- 端点响应时间
- 错误率统计

## 总结

Go PProf调试系统为项目提供了强大的性能分析能力：

1. **全面覆盖**：支持CPU、内存、Goroutine、阻塞等多种分析
2. **易于使用**：提供简洁的API和配置接口
3. **安全可控**：支持启用/禁用控制，适合生产环境
4. **集成友好**：与现有HTTP服务器无缝集成
5. **工具丰富**：支持命令行工具和可视化分析

通过合理使用PProf，可以快速识别和解决性能问题，提升应用的整体性能和质量。
