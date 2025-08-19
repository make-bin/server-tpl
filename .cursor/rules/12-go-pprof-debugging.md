# 12. Go PProf 调试规范

## 12.1 概述
项目集成了Go语言的PProf性能分析工具，提供全面的性能分析和调试能力。PProf通过HTTP端点提供实时性能数据，支持CPU、内存、Goroutine、阻塞等多种分析类型。

## 12.2 配置规范
```yaml
# configs/config.yaml
pprof:
  enabled: false  # 生产环境建议关闭
  path_prefix: /debug/pprof
```

## 12.3 功能特性

### 12.3.1 HTTP端点
启用PProf后，以下端点将可用：
- `/debug/pprof/` - PProf主页
- `/debug/pprof/profile` - CPU分析（30秒采样）
- `/debug/pprof/heap` - 堆内存分析
- `/debug/pprof/goroutine` - Goroutine分析
- `/debug/pprof/block` - 阻塞分析
- `/debug/pprof/mutex` - 互斥锁分析
- `/debug/pprof/allocs` - 内存分配分析
- `/debug/pprof/trace` - 执行跟踪

### 12.3.2 分析类型
- **CPU分析**: 分析CPU使用情况，识别热点函数
- **内存分析**: 分析内存分配模式，检测内存泄漏
- **Goroutine分析**: 分析Goroutine数量，检测泄漏
- **阻塞分析**: 分析程序阻塞情况，识别性能瓶颈
- **互斥锁分析**: 分析锁竞争情况，检测死锁风险

## 12.4 使用规范

### 12.4.1 基本使用
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

### 12.4.2 高级功能
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

## 12.5 命令行工具使用

### 12.5.1 基本命令
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

### 12.5.2 火焰图生成
```bash
# 安装火焰图工具
go install github.com/google/pprof@latest

# 生成火焰图
pprof -http=:8081 http://localhost:8080/debug/pprof/profile
```

## 12.6 性能分析场景

### 12.6.1 CPU性能分析
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

### 12.6.2 内存泄漏分析
```bash
# 获取堆内存分析
curl -o heap.prof http://localhost:8080/debug/pprof/heap

# 分析内存使用
go tool pprof heap.prof

# 查看内存分配
(pprof) top
(pprof) list allocFunction
```

### 12.6.3 Goroutine泄漏分析
```bash
# 获取Goroutine分析
curl -o goroutine.prof http://localhost:8080/debug/pprof/goroutine

# 分析Goroutine
go tool pprof goroutine.prof

# 查看Goroutine状态
(pprof) top
(pprof) traces
```

## 12.7 最佳实践

### 12.7.1 安全考虑
- **生产环境**: 默认关闭PProf，仅在需要时临时启用
- **访问控制**: 考虑添加认证和授权
- **网络隔离**: 限制PProf端点的网络访问
- **资源限制**: 避免长时间运行CPU分析

### 12.7.2 性能考虑
- **采样频率**: 合理设置阻塞和互斥锁采样频率
- **文件大小**: 定期清理分析文件
- **存储空间**: 监控分析文件占用的磁盘空间
- **网络带宽**: 避免频繁的远程分析

### 12.7.3 分析策略
- **问题导向**: 根据具体问题选择分析类型
- **渐进分析**: 从概览到细节逐步分析
- **对比分析**: 对比不同时间点的分析结果
- **持续监控**: 建立性能基准和监控

## 12.8 监控和告警

### 12.8.1 性能指标
- **分析文件大小**: 监控分析文件占用的磁盘空间
- **分析频率**: 监控分析操作的频率
- **响应时间**: 监控PProf端点的响应时间
- **错误率**: 监控分析操作的错误率

### 12.8.2 告警规则
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

## 12.9 故障排查

### 12.9.1 PProf端点不可访问
**问题**: 无法访问PProf端点
**解决方案**:
```yaml
pprof:
  enabled: true
  path_prefix: /debug/pprof
```

### 12.9.2 分析文件过大
**问题**: 分析文件占用过多磁盘空间
**解决方案**:
```go
// 降低采样频率
pprofManager.EnableBlockProfiling(1000) // 每1000纳秒采样一次

// 定期清理文件
// 实现文件清理逻辑
```

### 12.9.3 性能影响
**问题**: PProf影响应用性能
**解决方案**:
```go
// 仅在需要时启用
if debugMode {
    pprofManager.EnableBlockProfiling(1)
}

// 限制分析时间
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## 12.10 工具集成

### 12.10.1 项目结构
```
pkg/utils/pprof/
├── pprof.go      # PProf管理器实现
└── example.go    # 使用示例
```

### 12.10.2 文档
- `docs/pprof-debugging.md` - 详细的PProf调试规范
- 包含使用指南、最佳实践、故障排查等

## 12.11 总结

Go PProf调试系统为项目提供了强大的性能分析能力：

1. **全面覆盖**: 支持CPU、内存、Goroutine、阻塞等多种分析
2. **易于使用**: 提供简洁的API和配置接口
3. **安全可控**: 支持启用/禁用控制，适合生产环境
4. **集成友好**: 与现有HTTP服务器无缝集成
5. **工具丰富**: 支持命令行工具和可视化分析

通过合理使用PProf，可以快速识别和解决性能问题，提升应用的整体性能和质量。
