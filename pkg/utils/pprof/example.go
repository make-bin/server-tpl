package pprof

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// PProfExample PProf使用示例
func PProfExample() {
	// 1. 创建配置
	config := &PProfConfig{
		Enabled:    true,
		PathPrefix: "/debug/pprof",
	}

	// 2. 创建PProf管理器
	pprofManager := NewPProfManager(config)

	// 3. 启用阻塞和互斥锁分析
	pprofManager.EnableBlockProfiling(1) // 每纳秒采样一次
	pprofManager.EnableMutexProfiling(1) // 采样所有互斥锁

	// 4. 开始定期性能分析
	ctx := context.Background()
	pprofManager.StartPeriodicProfiling(ctx, 5*time.Minute, "./profiles")

	// 5. 创建Gin引擎并添加PProf路由
	engine := gin.Default()

	// 6. 示例：CPU密集型操作
	engine.GET("/cpu-intensive", func(c *gin.Context) {
		// 开始CPU分析
		cpuFile, err := pprofManager.StartCPUProfile("cpu_profile.prof")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer func() {
			pprofManager.StopCPUProfile()
			cpuFile.Close()
		}()

		// 执行CPU密集型操作
		performCPUIntensiveTask()

		c.JSON(200, gin.H{"message": "CPU intensive task completed"})
	})

	// 7. 示例：内存密集型操作
	engine.GET("/memory-intensive", func(c *gin.Context) {
		// 执行内存密集型操作
		performMemoryIntensiveTask()

		// 生成堆内存分析
		if err := pprofManager.WriteHeapProfile("heap_profile.prof"); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Memory intensive task completed"})
	})

	// 8. 示例：Goroutine泄漏
	engine.GET("/goroutine-leak", func(c *gin.Context) {
		// 创建一些Goroutine（模拟泄漏）
		for i := 0; i < 100; i++ {
			go func() {
				time.Sleep(time.Hour) // 长时间运行的Goroutine
			}()
		}

		// 生成Goroutine分析
		if err := pprofManager.WriteGoroutineProfile("goroutine_profile.prof"); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Goroutine leak simulation completed"})
	})

	// 9. 示例：阻塞操作
	engine.GET("/blocking", func(c *gin.Context) {
		// 执行阻塞操作
		performBlockingTask()

		// 生成阻塞分析
		if err := pprofManager.WriteBlockProfile("block_profile.prof"); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Blocking task completed"})
	})

	// 10. 示例：获取运行时统计
	engine.GET("/runtime-stats", func(c *gin.Context) {
		stats := pprofManager.GetRuntimeStats()
		c.JSON(200, stats)
	})

	// 11. 示例：生成完整分析报告
	engine.GET("/generate-profile", func(c *gin.Context) {
		if err := pprofManager.GenerateFullProfile("./profiles"); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Full profile generated"})
	})

	// 启动服务器
	engine.Run(":8080")
}

// performCPUIntensiveTask 执行CPU密集型任务
func performCPUIntensiveTask() {
	// 模拟CPU密集型计算
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

// performMemoryIntensiveTask 执行内存密集型任务
func performMemoryIntensiveTask() {
	// 模拟内存分配
	var data []int
	for i := 0; i < 1000000; i++ {
		data = append(data, i)
	}
	// 故意不释放内存，模拟内存泄漏
}

// performBlockingTask 执行阻塞任务
func performBlockingTask() {
	// 模拟阻塞操作
	time.Sleep(100 * time.Millisecond)
}

// BusinessServiceWithPProf 业务服务中使用PProf的示例
type BusinessServiceWithPProf struct {
	pprofManager *PProfManager
}

// NewBusinessServiceWithPProf 创建业务服务示例
func NewBusinessServiceWithPProf(pprofManager *PProfManager) *BusinessServiceWithPProf {
	return &BusinessServiceWithPProf{
		pprofManager: pprofManager,
	}
}

// ProcessData 处理数据（带性能分析）
func (s *BusinessServiceWithPProf) ProcessData(data []int) error {
	// 开始CPU分析
	cpuFile, err := s.pprofManager.StartCPUProfile("process_data_cpu.prof")
	if err != nil {
		return err
	}
	defer func() {
		s.pprofManager.StopCPUProfile()
		cpuFile.Close()
	}()

	// 处理数据
	for _, item := range data {
		// 模拟数据处理
		_ = item * 2
	}

	// 生成堆内存分析
	if err := s.pprofManager.WriteHeapProfile("process_data_heap.prof"); err != nil {
		return err
	}

	return nil
}

// MonitorPerformance 监控性能
func (s *BusinessServiceWithPProf) MonitorPerformance() map[string]interface{} {
	// 获取运行时统计
	stats := s.pprofManager.GetRuntimeStats()

	// 添加业务指标
	stats["business_metric"] = "custom_value"

	return stats
}

// GeneratePerformanceReport 生成性能报告
func (s *BusinessServiceWithPProf) GeneratePerformanceReport() error {
	// 生成完整分析报告
	return s.pprofManager.GenerateFullProfile("./business_profiles")
}
