package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP request duration histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Active connections gauge
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)
)

// PrometheusMiddleware implements the Middleware interface for Prometheus metrics
type PrometheusMiddleware struct {
	*BaseMiddleware
}

// NewPrometheusMiddleware creates a new Prometheus middleware
func NewPrometheusMiddleware() *PrometheusMiddleware {
	return &PrometheusMiddleware{
		BaseMiddleware: NewBaseMiddleware("prometheus", 15),
	}
}

// Handle processes the request with Prometheus metrics collection
func (m *PrometheusMiddleware) Handle(ctx context.Context, req *http.Request, next Handler) (*http.Response, error) {
	start := time.Now()

	// Increment active connections
	activeConnections.Inc()
	defer activeConnections.Dec()

	// Execute next handler
	resp, err := next.Handle(ctx, req)

	// Calculate duration
	duration := time.Since(start).Seconds()

	// Extract labels
	method := req.Method
	endpoint := req.URL.Path
	if endpoint == "" {
		endpoint = "unknown"
	}

	status := "200"
	if resp != nil {
		status = strconv.Itoa(resp.StatusCode)
	}
	if err != nil {
		status = "500"
	}

	// Record metrics
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)

	return resp, err
}

// PrometheusGinMiddleware returns a Gin middleware for Prometheus metrics collection (legacy)
func PrometheusGinMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Increment active connections
		activeConnections.Inc()
		defer activeConnections.Dec()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Extract labels
		method := c.Request.Method
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	})
}

// MetricsHandler returns a Gin handler for the /metrics endpoint
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
