package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/gin-gonic/gin"
)

// PProfConfig holds PProf configuration
type PProfConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	PathPrefix string `mapstructure:"path_prefix"`
	Port       int    `mapstructure:"port"`
}

// PProfManager manages PProf profiling
type PProfManager struct {
	config     *PProfConfig
	cpuFile    *os.File
	traceFile  *os.File
	httpServer *http.Server
}

// RuntimeStats holds runtime statistics
type RuntimeStats struct {
	NumGoroutine  int         `json:"num_goroutine"`
	NumCPU        int         `json:"num_cpu"`
	MemAlloc      uint64      `json:"mem_alloc"`
	MemTotalAlloc uint64      `json:"mem_total_alloc"`
	MemSys        uint64      `json:"mem_sys"`
	MemLookups    uint64      `json:"mem_lookups"`
	MemMallocs    uint64      `json:"mem_mallocs"`
	MemFrees      uint64      `json:"mem_frees"`
	HeapAlloc     uint64      `json:"heap_alloc"`
	HeapSys       uint64      `json:"heap_sys"`
	HeapIdle      uint64      `json:"heap_idle"`
	HeapInuse     uint64      `json:"heap_inuse"`
	HeapReleased  uint64      `json:"heap_released"`
	HeapObjects   uint64      `json:"heap_objects"`
	StackInuse    uint64      `json:"stack_inuse"`
	StackSys      uint64      `json:"stack_sys"`
	MSpanInuse    uint64      `json:"mspan_inuse"`
	MSpanSys      uint64      `json:"mspan_sys"`
	MCacheInuse   uint64      `json:"mcache_inuse"`
	MCacheSys     uint64      `json:"mcache_sys"`
	BuckHashSys   uint64      `json:"buck_hash_sys"`
	GCSys         uint64      `json:"gc_sys"`
	OtherSys      uint64      `json:"other_sys"`
	NextGC        uint64      `json:"next_gc"`
	LastGC        uint64      `json:"last_gc"`
	PauseTotalNs  uint64      `json:"pause_total_ns"`
	PauseNs       [256]uint64 `json:"pause_ns"`
	PauseEnd      [256]uint64 `json:"pause_end"`
	NumGC         uint32      `json:"num_gc"`
	NumForcedGC   uint32      `json:"num_forced_gc"`
	GCCPUFraction float64     `json:"gc_cpu_fraction"`
	EnableGC      bool        `json:"enable_gc"`
	DebugGC       bool        `json:"debug_gc"`
	Timestamp     time.Time   `json:"timestamp"`
}

// NewPProfManager creates a new PProf manager
func NewPProfManager(config *PProfConfig) *PProfManager {
	return &PProfManager{
		config: config,
	}
}

// StartHTTPServer starts the PProf HTTP server
func (p *PProfManager) StartHTTPServer() error {
	if !p.config.Enabled {
		return nil
	}

	mux := http.NewServeMux()

	// Register pprof handlers
	mux.HandleFunc(p.config.PathPrefix+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	p.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.config.Port),
		Handler: mux,
	}

	go func() {
		if err := p.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("PProf HTTP server error: %v\n", err)
		}
	}()

	return nil
}

// StopHTTPServer stops the PProf HTTP server
func (p *PProfManager) StopHTTPServer() error {
	if p.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.httpServer.Shutdown(ctx)
	}
	return nil
}

// StartCPUProfile starts CPU profiling
func (p *PProfManager) StartCPUProfile(filename string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to start CPU profile: %w", err)
	}

	p.cpuFile = file
	return file, nil
}

// StopCPUProfile stops CPU profiling
func (p *PProfManager) StopCPUProfile() {
	if p.cpuFile != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.cpuFile = nil
	}
}

// WriteHeapProfile writes heap profile to file
func (p *PProfManager) WriteHeapProfile(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create heap profile file: %w", err)
	}
	defer file.Close()

	runtime.GC() // Get fresh heap statistics
	if err := pprof.WriteHeapProfile(file); err != nil {
		return fmt.Errorf("failed to write heap profile: %w", err)
	}

	return nil
}

// WriteGoroutineProfile writes goroutine profile to file
func (p *PProfManager) WriteGoroutineProfile(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create goroutine profile file: %w", err)
	}
	defer file.Close()

	profile := pprof.Lookup("goroutine")
	if profile == nil {
		return fmt.Errorf("goroutine profile not found")
	}

	if err := profile.WriteTo(file, 0); err != nil {
		return fmt.Errorf("failed to write goroutine profile: %w", err)
	}

	return nil
}

// StartTrace starts execution tracing
func (p *PProfManager) StartTrace(filename string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return nil, fmt.Errorf("failed to create trace directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace file: %w", err)
	}

	if err := trace.Start(file); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to start trace: %w", err)
	}

	p.traceFile = file
	return file, nil
}

// StopTrace stops execution tracing
func (p *PProfManager) StopTrace() {
	if p.traceFile != nil {
		trace.Stop()
		p.traceFile.Close()
		p.traceFile = nil
	}
}

// EnableBlockProfiling enables block profiling
func (p *PProfManager) EnableBlockProfiling(rate int) {
	runtime.SetBlockProfileRate(rate)
}

// EnableMutexProfiling enables mutex profiling
func (p *PProfManager) EnableMutexProfiling(rate int) {
	runtime.SetMutexProfileFraction(rate)
}

// GetRuntimeStats returns current runtime statistics
func (p *PProfManager) GetRuntimeStats() *RuntimeStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &RuntimeStats{
		NumGoroutine:  runtime.NumGoroutine(),
		NumCPU:        runtime.NumCPU(),
		MemAlloc:      m.Alloc,
		MemTotalAlloc: m.TotalAlloc,
		MemSys:        m.Sys,
		MemLookups:    m.Lookups,
		MemMallocs:    m.Mallocs,
		MemFrees:      m.Frees,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapReleased:  m.HeapReleased,
		HeapObjects:   m.HeapObjects,
		StackInuse:    m.StackInuse,
		StackSys:      m.StackSys,
		MSpanInuse:    m.MSpanInuse,
		MSpanSys:      m.MSpanSys,
		MCacheInuse:   m.MCacheInuse,
		MCacheSys:     m.MCacheSys,
		BuckHashSys:   m.BuckHashSys,
		GCSys:         m.GCSys,
		OtherSys:      m.OtherSys,
		NextGC:        m.NextGC,
		LastGC:        m.LastGC,
		PauseTotalNs:  m.PauseTotalNs,
		PauseNs:       m.PauseNs,
		PauseEnd:      m.PauseEnd,
		NumGC:         m.NumGC,
		NumForcedGC:   m.NumForcedGC,
		GCCPUFraction: m.GCCPUFraction,
		EnableGC:      m.EnableGC,
		DebugGC:       m.DebugGC,
		Timestamp:     time.Now(),
	}
}

// StartPeriodicProfiling starts periodic profiling
func (p *PProfManager) StartPeriodicProfiling(ctx context.Context, interval time.Duration, outputDir string) {
	if !p.config.Enabled {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			timestamp := time.Now().Format("20060102_150405")

			// Write heap profile
			heapFile := filepath.Join(outputDir, fmt.Sprintf("heap_%s.prof", timestamp))
			if err := p.WriteHeapProfile(heapFile); err != nil {
				fmt.Printf("Failed to write heap profile: %v\n", err)
			}

			// Write goroutine profile
			goroutineFile := filepath.Join(outputDir, fmt.Sprintf("goroutine_%s.prof", timestamp))
			if err := p.WriteGoroutineProfile(goroutineFile); err != nil {
				fmt.Printf("Failed to write goroutine profile: %v\n", err)
			}
		}
	}
}

// GenerateFullProfile generates a complete set of profiles
func (p *PProfManager) GenerateFullProfile(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	// Generate heap profile
	heapFile := filepath.Join(outputDir, fmt.Sprintf("heap_%s.prof", timestamp))
	if err := p.WriteHeapProfile(heapFile); err != nil {
		return fmt.Errorf("failed to write heap profile: %w", err)
	}

	// Generate goroutine profile
	goroutineFile := filepath.Join(outputDir, fmt.Sprintf("goroutine_%s.prof", timestamp))
	if err := p.WriteGoroutineProfile(goroutineFile); err != nil {
		return fmt.Errorf("failed to write goroutine profile: %w", err)
	}

	// Generate CPU profile (5 seconds)
	cpuFile := filepath.Join(outputDir, fmt.Sprintf("cpu_%s.prof", timestamp))
	file, err := p.StartCPUProfile(cpuFile)
	if err != nil {
		return fmt.Errorf("failed to start CPU profile: %w", err)
	}

	time.Sleep(5 * time.Second)
	p.StopCPUProfile()
	file.Close()

	return nil
}

// RegisterRoutes registers PProf routes with Gin router
func (p *PProfManager) RegisterRoutes(router *gin.Engine) {
	if !p.config.Enabled {
		return
	}

	pprofGroup := router.Group(p.config.PathPrefix)
	{
		pprofGroup.GET("/", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, p.config.PathPrefix+"/debug/pprof/", http.StatusMovedPermanently)
		})))
		pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.POST("/symbol", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/allocs", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/block", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/goroutine", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/heap", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/mutex", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		pprofGroup.GET("/threadcreate", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
	}

	// Add runtime stats endpoint
	pprofGroup.GET("/stats", func(c *gin.Context) {
		stats := p.GetRuntimeStats()
		c.JSON(http.StatusOK, stats)
	})
}
