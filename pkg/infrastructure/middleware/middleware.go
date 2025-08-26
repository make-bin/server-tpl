package middleware

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// Middleware interface defines a middleware component
type Middleware interface {
	// Handle processes the request and calls the next handler
	Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error)

	// Name returns the middleware name
	Name() string

	// Priority returns the middleware priority (lower numbers = higher priority)
	Priority() int
}

// Handler interface defines a request handler
type Handler interface {
	Handle(ctx context.Context, req *http.Request) (*http.Response, error)
}

// Chain represents a middleware chain
type Chain struct {
	middlewares []Middleware
	sorted      bool
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
	return &Chain{
		middlewares: make([]Middleware, 0),
		sorted:      true,
	}
}

// Use adds a middleware to the chain
func (c *Chain) Use(middleware Middleware) *Chain {
	c.middlewares = append(c.middlewares, middleware)
	c.sorted = false
	return c
}

// Handle executes the middleware chain
func (c *Chain) Handle(ctx context.Context, req *http.Request) (*http.Response, error) {
	if !c.sorted {
		c.sortMiddlewares()
	}

	if len(c.middlewares) == 0 {
		return &http.Response{StatusCode: http.StatusNotFound}, nil
	}

	// Create handler chain
	handler := &chainHandler{
		chain: c,
		index: 0,
	}

	// Execute first middleware
	return c.middlewares[0].Handle(ctx, req, handler)
}

// sortMiddlewares sorts middlewares by priority
func (c *Chain) sortMiddlewares() {
	sort.Slice(c.middlewares, func(i, j int) bool {
		return c.middlewares[i].Priority() < c.middlewares[j].Priority()
	})
	c.sorted = true
}

// chainHandler implements Handler for middleware chain execution
type chainHandler struct {
	chain *Chain
	index int
}

// Handle executes the next middleware in the chain
func (h *chainHandler) Handle(ctx context.Context, req *http.Request) (*http.Response, error) {
	h.index++

	if h.index >= len(h.chain.middlewares) {
		// End of chain, return empty response
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	return h.chain.middlewares[h.index].Handle(ctx, req, h)
}

// GinMiddleware adapts a Middleware to gin.HandlerFunc
func GinMiddleware(middleware Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context
		ctx := c.Request.Context()

		// Create gin handler wrapper
		handler := &ginHandler{c: c}

		// Execute middleware
		resp, err := middleware.Handle(ctx, c.Request, handler)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// Handle response if provided
		if resp != nil && resp.StatusCode >= 400 {
			c.Status(resp.StatusCode)
			c.Abort()
		}
	}
}

// ginHandler implements Handler for Gin framework
type ginHandler struct {
	c *gin.Context
}

// Handle continues with the next Gin handler
func (h *ginHandler) Handle(ctx context.Context, req *http.Request) (*http.Response, error) {
	h.c.Next()

	// Check if there were any errors
	if len(h.c.Errors) > 0 {
		return &http.Response{
			StatusCode: h.c.Writer.Status(),
		}, h.c.Errors.Last().Err
	}

	return &http.Response{
		StatusCode: h.c.Writer.Status(),
	}, nil
}

// BaseMiddleware provides common middleware functionality
type BaseMiddleware struct {
	name     string
	priority int
}

// NewBaseMiddleware creates a new base middleware
func NewBaseMiddleware(name string, priority int) *BaseMiddleware {
	return &BaseMiddleware{
		name:     name,
		priority: priority,
	}
}

// Name returns the middleware name
func (m *BaseMiddleware) Name() string {
	return m.name
}

// Priority returns the middleware priority
func (m *BaseMiddleware) Priority() int {
	return m.priority
}

// LoggerMiddleware implements request logging
type LoggerMiddleware struct {
	*BaseMiddleware
	logger logger.Manager
}

// NewLoggerMiddleware creates a new logger middleware
func NewLoggerMiddleware(loggerManager logger.Manager) *LoggerMiddleware {
	return &LoggerMiddleware{
		BaseMiddleware: NewBaseMiddleware("logger", 10),
		logger:         loggerManager,
	}
}

// Handle processes the request with logging
func (m *LoggerMiddleware) Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error) {
	start := time.Now()

	// Log request start
	m.logger.WithContext(ctx).WithFields(map[string]interface{}{
		"method":      req.Method,
		"path":        req.URL.Path,
		"remote_addr": req.RemoteAddr,
		"user_agent":  req.UserAgent(),
	}).Info("HTTP request started")

	// Execute next handler
	resp, err := next.Handle(ctx, req)

	// Log request completion
	duration := time.Since(start)
	statusCode := http.StatusOK
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	m.logger.WithContext(ctx).WithFields(map[string]interface{}{
		"duration":    duration.Milliseconds(),
		"status_code": statusCode,
	}).Info("HTTP request completed")

	return resp, err
}

// ErrorHandlerMiddleware implements error handling
type ErrorHandlerMiddleware struct {
	*BaseMiddleware
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		BaseMiddleware: NewBaseMiddleware("error_handler", 1000), // Low priority (runs last)
	}
}

// Handle processes the request with error handling
func (m *ErrorHandlerMiddleware) Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered in middleware: %v", r)
		}
	}()

	resp, err := next.Handle(ctx, req)
	if err != nil {
		logger.Error("Error in middleware chain: %v", err)
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CORSMiddleware implements CORS handling
type CORSMiddleware struct {
	*BaseMiddleware
	config CORSConfig
}

// CORSConfig defines CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		BaseMiddleware: NewBaseMiddleware("cors", 20),
		config:         config,
	}
}

// Handle processes the request with CORS
func (m *CORSMiddleware) Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error) {
	// This is a simplified CORS implementation
	// In real implementation, you would check origins, methods, etc.

	return next.Handle(ctx, req)
}

// RequestIDMiddleware implements request ID generation
type RequestIDMiddleware struct {
	*BaseMiddleware
}

// NewRequestIDMiddleware creates a new request ID middleware
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{
		BaseMiddleware: NewBaseMiddleware("request_id", 5),
	}
}

// Handle processes the request with request ID
func (m *RequestIDMiddleware) Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error) {
	// Generate request ID (simplified)
	requestID := generateRequestID()

	// Add to context
	ctx = context.WithValue(ctx, "request_id", requestID)

	return next.Handle(ctx, req)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simplified implementation - in real usage, use UUID or similar
	return "req_" + time.Now().Format("20060102150405")
}

// MiddlewareConfig defines middleware configuration
type MiddlewareConfig struct {
	Logger struct {
		Enabled bool   `json:"enabled"`
		Level   string `json:"level"`
		Format  string `json:"format"`
	} `json:"logger"`

	Auth struct {
		Enabled bool   `json:"enabled"`
		Type    string `json:"type"`
		Secret  string `json:"secret"`
	} `json:"auth"`

	RateLimit struct {
		Enabled bool `json:"enabled"`
		Limit   int  `json:"limit"`
		Window  int  `json:"window"`
	} `json:"rate_limit"`

	CORS CORSConfig `json:"cors"`
}

// RegisterMiddlewares registers middlewares with Gin engine
func RegisterMiddlewares(engine *gin.Engine, config *MiddlewareConfig, loggerManager logger.Manager) {
	// Register request ID middleware first
	engine.Use(GinMiddleware(NewRequestIDMiddleware()))

	// Register logger middleware
	if config.Logger.Enabled {
		engine.Use(GinMiddleware(NewLoggerMiddleware(loggerManager)))
	}

	// Register CORS middleware
	engine.Use(GinMiddleware(NewCORSMiddleware(config.CORS)))

	// Register error handler middleware last
	engine.Use(GinMiddleware(NewErrorHandlerMiddleware()))
}

// MiddlewareManager manages middleware registration and configuration
type MiddlewareManager struct {
	middlewares map[string]Middleware
	config      *MiddlewareConfig
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager(config *MiddlewareConfig) *MiddlewareManager {
	return &MiddlewareManager{
		middlewares: make(map[string]Middleware),
		config:      config,
	}
}

// Register registers a middleware
func (m *MiddlewareManager) Register(name string, middleware Middleware) {
	m.middlewares[name] = middleware
}

// Get retrieves a middleware by name
func (m *MiddlewareManager) Get(name string) (Middleware, bool) {
	middleware, exists := m.middlewares[name]
	return middleware, exists
}

// GetAll returns all registered middlewares
func (m *MiddlewareManager) GetAll() map[string]Middleware {
	return m.middlewares
}

// CreateChain creates a middleware chain with specified middlewares
func (m *MiddlewareManager) CreateChain(names ...string) *Chain {
	chain := NewChain()

	for _, name := range names {
		if middleware, exists := m.middlewares[name]; exists {
			chain.Use(middleware)
		}
	}

	return chain
}
