package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PerformanceMonitor implements the Monitor interface
type PerformanceMonitor struct {
	queryDuration    *prometheus.HistogramVec
	connectionGauge  *prometheus.GaugeVec
	errorCounter     *prometheus.CounterVec
	operationCounter *prometheus.CounterVec
	mutex            sync.RWMutex
	stats            map[string]*operationStats
}

// operationStats holds statistics for database operations
type operationStats struct {
	Count        int64         `json:"count"`
	TotalTime    time.Duration `json:"total_time"`
	AverageTime  time.Duration `json:"average_time"`
	LastExecuted time.Time     `json:"last_executed"`
	Errors       int64         `json:"errors"`
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() datastore.Monitor {
	monitor := &PerformanceMonitor{
		stats: make(map[string]*operationStats),
	}

	// Initialize Prometheus metrics
	monitor.queryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "datastore_query_duration_seconds",
			Help:    "Time spent executing datastore queries",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
		},
		[]string{"operation", "table"},
	)

	monitor.connectionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "datastore_connections",
			Help: "Number of active database connections",
		},
		[]string{"database"},
	)

	monitor.errorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "datastore_errors_total",
			Help: "Total number of datastore errors",
		},
		[]string{"operation", "table", "error_type"},
	)

	monitor.operationCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "datastore_operations_total",
			Help: "Total number of datastore operations",
		},
		[]string{"operation", "table"},
	)

	return monitor
}

// RecordQuery records a database query execution
func (m *PerformanceMonitor) RecordQuery(operation, table string, duration time.Duration) {
	// Update Prometheus metrics
	m.queryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
	m.operationCounter.WithLabelValues(operation, table).Inc()

	// Update internal stats
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := operation + ":" + table
	stats, exists := m.stats[key]
	if !exists {
		stats = &operationStats{}
		m.stats[key] = stats
	}

	stats.Count++
	stats.TotalTime += duration
	stats.AverageTime = stats.TotalTime / time.Duration(stats.Count)
	stats.LastExecuted = time.Now()

	// Log slow queries
	if duration > time.Second {
		logger.Warn("Slow query detected: operation=%s, table=%s, duration=%v",
			operation, table, duration)
	}
}

// RecordConnection records database connection count
func (m *PerformanceMonitor) RecordConnection(database string, connections int) {
	m.connectionGauge.WithLabelValues(database).Set(float64(connections))
}

// RecordError records a database error
func (m *PerformanceMonitor) RecordError(operation, table string, err error) {
	errorType := "unknown"
	if err != nil {
		errorType = err.Error()
		// Classify common errors
		switch err {
		case datastore.ErrNotFound:
			errorType = "not_found"
		case datastore.ErrDuplicateKey:
			errorType = "duplicate_key"
		case datastore.ErrConnectionFailed:
			errorType = "connection_failed"
		case datastore.ErrTransactionFailed:
			errorType = "transaction_failed"
		}
	}

	m.errorCounter.WithLabelValues(operation, table, errorType).Inc()

	// Update internal stats
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := operation + ":" + table
	stats, exists := m.stats[key]
	if exists {
		stats.Errors++
	}

	logger.Error("Database error: operation=%s, table=%s, error=%v",
		operation, table, err)
}

// GetStats returns performance statistics
func (m *PerformanceMonitor) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]interface{})
	for key, stats := range m.stats {
		result[key] = map[string]interface{}{
			"count":         stats.Count,
			"total_time":    stats.TotalTime.String(),
			"average_time":  stats.AverageTime.String(),
			"last_executed": stats.LastExecuted.Format(time.RFC3339),
			"errors":        stats.Errors,
		}
	}

	return result
}

// StatsCollector implements the Stats interface
type StatsCollector struct {
	monitor     datastore.Monitor
	connections map[string]int
	queries     map[string]int64
	mutex       sync.RWMutex
}

// NewStatsCollector creates a new stats collector
func NewStatsCollector(monitor datastore.Monitor) datastore.Stats {
	return &StatsCollector{
		monitor:     monitor,
		connections: make(map[string]int),
		queries:     make(map[string]int64),
	}
}

// GetStats returns general statistics
func (s *StatsCollector) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"connections": s.connections,
		"queries":     s.queries,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
}

// GetConnectionStats returns connection statistics
func (s *StatsCollector) GetConnectionStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := make(map[string]interface{})
	for db, count := range s.connections {
		stats[db] = map[string]interface{}{
			"active_connections": count,
			"last_updated":       time.Now().Format(time.RFC3339),
		}
	}

	return stats
}

// GetQueryStats returns query statistics
func (s *StatsCollector) GetQueryStats() map[string]interface{} {
	// Note: monitor doesn't have GetStats method, use local stats

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"total_queries": s.queries,
		"timestamp":     time.Now().Format(time.RFC3339),
	}
}

// UpdateConnectionCount updates connection count for a database
func (s *StatsCollector) UpdateConnectionCount(database string, count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.connections[database] = count

	if s.monitor != nil {
		s.monitor.RecordConnection(database, count)
	}
}

// IncrementQueryCount increments query count for an operation
func (s *StatsCollector) IncrementQueryCount(operation string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.queries[operation]++
}

// MonitoredDataStore wraps a DataStore with monitoring
type MonitoredDataStore struct {
	store   datastore.DataStore
	monitor datastore.Monitor
	stats   datastore.Stats
	dbName  string
}

// NewMonitoredDataStore creates a new monitored data store
func NewMonitoredDataStore(store datastore.DataStore, dbName string) *MonitoredDataStore {
	monitor := NewPerformanceMonitor()
	stats := NewStatsCollector(monitor)

	return &MonitoredDataStore{
		store:   store,
		monitor: monitor,
		stats:   stats,
		dbName:  dbName,
	}
}

// Connect connects to the database with monitoring
func (m *MonitoredDataStore) Connect(ctx context.Context) error {
	start := time.Now()
	err := m.store.Connect(ctx)
	duration := time.Since(start)

	m.monitor.RecordQuery("connect", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("connect", m.dbName, err)
	}

	return err
}

// Disconnect disconnects from the database with monitoring
func (m *MonitoredDataStore) Disconnect(ctx context.Context) error {
	start := time.Now()
	err := m.store.Disconnect(ctx)
	duration := time.Since(start)

	m.monitor.RecordQuery("disconnect", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("disconnect", m.dbName, err)
	}

	return err
}

// HealthCheck performs health check with monitoring
func (m *MonitoredDataStore) HealthCheck(ctx context.Context) error {
	start := time.Now()
	err := m.store.HealthCheck(ctx)
	duration := time.Since(start)

	m.monitor.RecordQuery("health_check", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("health_check", m.dbName, err)
	}

	return err
}

// BeginTx begins a transaction with monitoring
func (m *MonitoredDataStore) BeginTx(ctx context.Context) (datastore.Transaction, error) {
	start := time.Now()
	tx, err := m.store.BeginTx(ctx)
	duration := time.Since(start)

	m.monitor.RecordQuery("begin_tx", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("begin_tx", m.dbName, err)
	}

	// Wrap transaction with monitoring
	if tx != nil {
		return &MonitoredTransaction{
			tx:      tx,
			monitor: m.monitor,
			dbName:  m.dbName,
		}, nil
	}

	return nil, err
}

// Add adds an entity with monitoring
func (m *MonitoredDataStore) Add(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.store.Add(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("add", tableName, duration)
	if err != nil {
		m.monitor.RecordError("add", tableName, err)
	}

	return err
}

// BatchAdd adds multiple entities with monitoring
func (m *MonitoredDataStore) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	start := time.Now()
	err := m.store.BatchAdd(ctx, entities)
	duration := time.Since(start)

	tableName := "batch"
	if len(entities) > 0 {
		tableName = entities[0].TableName()
	}

	m.monitor.RecordQuery("batch_add", tableName, duration)
	if err != nil {
		m.monitor.RecordError("batch_add", tableName, err)
	}

	return err
}

// Put updates an entity with monitoring
func (m *MonitoredDataStore) Put(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.store.Put(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("put", tableName, duration)
	if err != nil {
		m.monitor.RecordError("put", tableName, err)
	}

	return err
}

// Delete removes an entity with monitoring
func (m *MonitoredDataStore) Delete(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.store.Delete(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("delete", tableName, duration)
	if err != nil {
		m.monitor.RecordError("delete", tableName, err)
	}

	return err
}

// Get retrieves an entity with monitoring
func (m *MonitoredDataStore) Get(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.store.Get(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("get", tableName, duration)
	if err != nil {
		m.monitor.RecordError("get", tableName, err)
	}

	return err
}

// List retrieves entities with monitoring
func (m *MonitoredDataStore) List(ctx context.Context, query datastore.Entity, options *datastore.ListOptions) ([]datastore.Entity, error) {
	start := time.Now()
	entities, err := m.store.List(ctx, query, options)
	duration := time.Since(start)

	tableName := query.TableName()
	m.monitor.RecordQuery("list", tableName, duration)
	if err != nil {
		m.monitor.RecordError("list", tableName, err)
	}

	return entities, err
}

// Count counts entities with monitoring
func (m *MonitoredDataStore) Count(ctx context.Context, entity datastore.Entity, options *datastore.FilterOptions) (int64, error) {
	start := time.Now()
	count, err := m.store.Count(ctx, entity, options)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("count", tableName, duration)
	if err != nil {
		m.monitor.RecordError("count", tableName, err)
	}

	return count, err
}

// IsExist checks if entity exists with monitoring
func (m *MonitoredDataStore) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	start := time.Now()
	exists, err := m.store.IsExist(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("exists", tableName, duration)
	if err != nil {
		m.monitor.RecordError("exists", tableName, err)
	}

	return exists, err
}

// Migrate performs migration with monitoring
func (m *MonitoredDataStore) Migrate(ctx context.Context, entities ...datastore.Entity) error {
	start := time.Now()
	err := m.store.Migrate(ctx, entities...)
	duration := time.Since(start)

	m.monitor.RecordQuery("migrate", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("migrate", m.dbName, err)
	}

	return err
}

// ExecuteSQL executes SQL with monitoring
func (m *MonitoredDataStore) ExecuteSQL(ctx context.Context, sql string, args ...interface{}) error {
	start := time.Now()
	err := m.store.ExecuteSQL(ctx, sql, args...)
	duration := time.Since(start)

	m.monitor.RecordQuery("execute_sql", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("execute_sql", m.dbName, err)
	}

	return err
}

// GetMonitor returns the monitor instance
func (m *MonitoredDataStore) GetMonitor() datastore.Monitor {
	return m.monitor
}

// GetStats returns the stats instance
func (m *MonitoredDataStore) GetStats() datastore.Stats {
	return m.stats
}

// MonitoredTransaction wraps a Transaction with monitoring
type MonitoredTransaction struct {
	tx      datastore.Transaction
	monitor datastore.Monitor
	dbName  string
}

// Commit commits the transaction with monitoring
func (m *MonitoredTransaction) Commit() error {
	start := time.Now()
	err := m.tx.Commit()
	duration := time.Since(start)

	m.monitor.RecordQuery("commit", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("commit", m.dbName, err)
	}

	return err
}

// Rollback rolls back the transaction with monitoring
func (m *MonitoredTransaction) Rollback() error {
	start := time.Now()
	err := m.tx.Rollback()
	duration := time.Since(start)

	m.monitor.RecordQuery("rollback", m.dbName, duration)
	if err != nil {
		m.monitor.RecordError("rollback", m.dbName, err)
	}

	return err
}

// Add adds an entity in transaction with monitoring
func (m *MonitoredTransaction) Add(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.tx.Add(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("tx_add", tableName, duration)
	if err != nil {
		m.monitor.RecordError("tx_add", tableName, err)
	}

	return err
}

// Put updates an entity in transaction with monitoring
func (m *MonitoredTransaction) Put(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.tx.Put(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("tx_put", tableName, duration)
	if err != nil {
		m.monitor.RecordError("tx_put", tableName, err)
	}

	return err
}

// Delete removes an entity in transaction with monitoring
func (m *MonitoredTransaction) Delete(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.tx.Delete(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("tx_delete", tableName, duration)
	if err != nil {
		m.monitor.RecordError("tx_delete", tableName, err)
	}

	return err
}

// Get retrieves an entity in transaction with monitoring
func (m *MonitoredTransaction) Get(ctx context.Context, entity datastore.Entity) error {
	start := time.Now()
	err := m.tx.Get(ctx, entity)
	duration := time.Since(start)

	tableName := entity.TableName()
	m.monitor.RecordQuery("tx_get", tableName, duration)
	if err != nil {
		m.monitor.RecordError("tx_get", tableName, err)
	}

	return err
}
