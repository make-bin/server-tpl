package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Example metrics for demonstration purposes
var (
	// Business metrics example
	applicationCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "application_created_total",
			Help: "Total number of applications created",
		},
		[]string{"status"},
	)

	applicationProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "application_processing_duration_seconds",
			Help:    "Time spent processing applications",
			Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"operation"},
	)

	// System metrics example
	databaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
		[]string{"database"},
	)

	cacheHitRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "Cache hit ratio (0-1)",
		},
		[]string{"cache_type"},
	)
)

// RecordApplicationCreated records an application creation event
func RecordApplicationCreated(status string) {
	applicationCreated.WithLabelValues(status).Inc()
}

// RecordApplicationProcessingTime records the time spent processing an application
func RecordApplicationProcessingTime(operation string, duration float64) {
	applicationProcessingTime.WithLabelValues(operation).Observe(duration)
}

// SetDatabaseConnections sets the number of active database connections
func SetDatabaseConnections(database string, count float64) {
	databaseConnections.WithLabelValues(database).Set(count)
}

// SetCacheHitRatio sets the cache hit ratio
func SetCacheHitRatio(cacheType string, ratio float64) {
	cacheHitRatio.WithLabelValues(cacheType).Set(ratio)
}
