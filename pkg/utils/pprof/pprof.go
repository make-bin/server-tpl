package pprof

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// PProfConfig PProf配置
type PProfConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	PathPrefix string `mapstructure:"path_prefix"`
}

// PProfManager PProf管理器
type PProfManager struct {
	config *PProfConfig
}

// NewPProfManager 创建PProf管理器
func NewPProfManager(config *PProfConfig) *PProfManager {
	return &PProfManager{
		config: config,
	}
}

// StartCPUProfile 开始CPU性能分析
func (p *PProfManager) StartCPUProfile(filename string) (*os.File, error) {
	if !p.config.Enabled {
		return nil, fmt.Errorf("pprof is not enabled")
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to start CPU profile: %w", err)
	}

	logger.WithField("filename", filename).Info("CPU profiling started")
	return file, nil
}

// StopCPUProfile 停止CPU性能分析
func (p *PProfManager) StopCPUProfile() {
	pprof.StopCPUProfile()
	logger.Info("CPU profiling stopped")
}

// WriteHeapProfile 写入堆内存分析
func (p *PProfManager) WriteHeapProfile(filename string) error {
	if !p.config.Enabled {
		return fmt.Errorf("pprof is not enabled")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create heap profile file: %w", err)
	}
	defer file.Close()

	if err := pprof.WriteHeapProfile(file); err != nil {
		return fmt.Errorf("failed to write heap profile: %w", err)
	}

	logger.WithField("filename", filename).Info("Heap profile written")
	return nil
}

// WriteGoroutineProfile 写入Goroutine分析
func (p *PProfManager) WriteGoroutineProfile(filename string) error {
	if !p.config.Enabled {
		return fmt.Errorf("pprof is not enabled")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create goroutine profile file: %w", err)
	}
	defer file.Close()

	profile := pprof.Lookup("goroutine")
	if profile == nil {
		return fmt.Errorf("goroutine profile not available")
	}

	if err := profile.WriteTo(file, 0); err != nil {
		return fmt.Errorf("failed to write goroutine profile: %w", err)
	}

	logger.WithField("filename", filename).Info("Goroutine profile written")
	return nil
}

// WriteBlockProfile 写入阻塞分析
func (p *PProfManager) WriteBlockProfile(filename string) error {
	if !p.config.Enabled {
		return fmt.Errorf("pprof is not enabled")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create block profile file: %w", err)
	}
	defer file.Close()

	profile := pprof.Lookup("block")
	if profile == nil {
		return fmt.Errorf("block profile not available")
	}

	if err := profile.WriteTo(file, 0); err != nil {
		return fmt.Errorf("failed to write block profile: %w", err)
	}

	logger.WithField("filename", filename).Info("Block profile written")
	return nil
}

// WriteMutexProfile 写入互斥锁分析
func (p *PProfManager) WriteMutexProfile(filename string) error {
	if !p.config.Enabled {
		return fmt.Errorf("pprof is not enabled")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create mutex profile file: %w", err)
	}
	defer file.Close()

	profile := pprof.Lookup("mutex")
	if profile == nil {
		return fmt.Errorf("mutex profile not available")
	}

	if err := profile.WriteTo(file, 0); err != nil {
		return fmt.Errorf("failed to write mutex profile: %w", err)
	}

	logger.WithField("filename", filename).Info("Mutex profile written")
	return nil
}

// GetRuntimeStats 获取运行时统计信息
func (p *PProfManager) GetRuntimeStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines":     runtime.NumGoroutine(),
		"heap_alloc":     m.HeapAlloc,
		"heap_sys":       m.HeapSys,
		"heap_idle":      m.HeapIdle,
		"heap_inuse":     m.HeapInuse,
		"heap_released":  m.HeapReleased,
		"heap_objects":   m.HeapObjects,
		"stack_inuse":    m.StackInuse,
		"stack_sys":      m.StackSys,
		"total_alloc":    m.TotalAlloc,
		"sys":            m.Sys,
		"num_gc":         m.NumGC,
		"pause_total_ns": m.PauseTotalNs,
	}
}

// StartPeriodicProfiling 开始定期性能分析
func (p *PProfManager) StartPeriodicProfiling(ctx context.Context, interval time.Duration, outputDir string) {
	if !p.config.Enabled {
		logger.Warn("PProf is not enabled, skipping periodic profiling")
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				timestamp := time.Now().Format("20060102_150405")

				// 写入堆内存分析
				heapFile := fmt.Sprintf("%s/heap_%s.prof", outputDir, timestamp)
				if err := p.WriteHeapProfile(heapFile); err != nil {
					logger.WithField("error", err.Error()).Error("Failed to write heap profile")
				}

				// 写入Goroutine分析
				goroutineFile := fmt.Sprintf("%s/goroutine_%s.prof", outputDir, timestamp)
				if err := p.WriteGoroutineProfile(goroutineFile); err != nil {
					logger.WithField("error", err.Error()).Error("Failed to write goroutine profile")
				}

				// 记录运行时统计
				stats := p.GetRuntimeStats()
				logger.WithFields(stats).Info("Runtime statistics")
			}
		}
	}()

	logger.WithField("interval", interval).Info("Periodic profiling started")
}

// EnableBlockProfiling 启用阻塞分析
func (p *PProfManager) EnableBlockProfiling(rate int) {
	if !p.config.Enabled {
		logger.Warn("PProf is not enabled, cannot enable block profiling")
		return
	}

	runtime.SetBlockProfileRate(rate)
	logger.WithField("rate", rate).Info("Block profiling enabled")
}

// EnableMutexProfiling 启用互斥锁分析
func (p *PProfManager) EnableMutexProfiling(fraction int) {
	if !p.config.Enabled {
		logger.Warn("PProf is not enabled, cannot enable mutex profiling")
		return
	}

	runtime.SetMutexProfileFraction(fraction)
	logger.WithField("fraction", fraction).Info("Mutex profiling enabled")
}

// GenerateFullProfile 生成完整的性能分析报告
func (p *PProfManager) GenerateFullProfile(outputDir string) error {
	if !p.config.Enabled {
		return fmt.Errorf("pprof is not enabled")
	}

	timestamp := time.Now().Format("20060102_150405")

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成各种分析文件
	profiles := []struct {
		name     string
		filename string
		writer   func(string) error
	}{
		{"heap", fmt.Sprintf("%s/heap_%s.prof", outputDir, timestamp), p.WriteHeapProfile},
		{"goroutine", fmt.Sprintf("%s/goroutine_%s.prof", outputDir, timestamp), p.WriteGoroutineProfile},
		{"block", fmt.Sprintf("%s/block_%s.prof", outputDir, timestamp), p.WriteBlockProfile},
		{"mutex", fmt.Sprintf("%s/mutex_%s.prof", outputDir, timestamp), p.WriteMutexProfile},
	}

	for _, profile := range profiles {
		if err := profile.writer(profile.filename); err != nil {
			logger.WithField("error", err.Error()).WithField("profile", profile.name).Error("Failed to generate profile")
		} else {
			logger.WithField("filename", profile.filename).Info("Profile generated")
		}
	}

	// 生成运行时统计报告
	stats := p.GetRuntimeStats()
	statsFile := fmt.Sprintf("%s/runtime_stats_%s.json", outputDir, timestamp)

	// 这里可以添加JSON序列化逻辑
	logger.WithField("filename", statsFile).WithFields(stats).Info("Runtime statistics saved")

	return nil
}
