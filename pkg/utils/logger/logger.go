package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Standard log fields
const (
	// Request related fields
	FieldRequestID    = "request_id"
	FieldUserID       = "user_id"
	FieldIP           = "ip"
	FieldUserAgent    = "user_agent"
	FieldMethod       = "method"
	FieldPath         = "path"
	FieldStatusCode   = "status_code"
	FieldResponseTime = "response_time"

	// Business related fields
	FieldOperation  = "operation"
	FieldResource   = "resource"
	FieldResourceID = "resource_id"
	FieldAction     = "action"

	// Error related fields
	FieldError      = "error"
	FieldErrorCode  = "error_code"
	FieldErrorType  = "error_type"
	FieldStackTrace = "stack_trace"

	// System related fields
	FieldService     = "service"
	FieldVersion     = "version"
	FieldEnvironment = "environment"
	FieldHost        = "host"
	FieldPID         = "pid"
)

// Manager interface for log management
type Manager interface {
	GetLogger() *logrus.Logger
	WithContext(ctx context.Context) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	SetLevel(level string) error
	SetFormat(format string) error
	SetOutput(output string) error
}

// LogManager implements the Manager interface
type LogManager struct {
	logger     *logrus.Logger
	config     *LogConfig
	lumberjack *lumberjack.Logger
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string            `mapstructure:"level"`
	Format     string            `mapstructure:"format"` // json, text
	Output     string            `mapstructure:"output"` // stdout, file, both
	FilePath   string            `mapstructure:"file_path"`
	MaxSize    int               `mapstructure:"max_size"` // MB
	MaxBackups int               `mapstructure:"max_backups"`
	MaxAge     int               `mapstructure:"max_age"` // days
	Compress   bool              `mapstructure:"compress"`
	Fields     map[string]string `mapstructure:"fields"` // Default fields
	BufferSize int               `mapstructure:"buffer_size"`
	Async      bool              `mapstructure:"async"`
}

var (
	defaultManager *LogManager
)

// NewManager creates a new log manager
func NewManager(config *LogConfig) Manager {
	logger := logrus.New()

	manager := &LogManager{
		logger: logger,
		config: config,
	}

	// Set log level
	manager.SetLevel(config.Level)

	// Set log format
	manager.SetFormat(config.Format)

	// Set log output
	manager.SetOutput(config.Output)

	// Add default fields
	if config.Fields != nil {
		for key, value := range config.Fields {
			logger.WithField(key, value)
		}
	}

	return manager
}

// GetLogger returns the logrus logger instance
func (m *LogManager) GetLogger() *logrus.Logger {
	return m.logger
}

// WithContext returns a logger entry with context
func (m *LogManager) WithContext(ctx context.Context) *logrus.Entry {
	entry := m.logger.WithContext(ctx)

	// Add request ID if available
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		entry = entry.WithField(FieldRequestID, requestID)
	}

	// Add user ID if available
	if userID := ctx.Value(FieldUserID); userID != nil {
		entry = entry.WithField(FieldUserID, userID)
	}

	return entry
}

// WithFields returns a logger entry with fields
func (m *LogManager) WithFields(fields logrus.Fields) *logrus.Entry {
	return m.logger.WithFields(fields)
}

// SetLevel sets the log level
func (m *LogManager) SetLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}

	m.logger.SetLevel(logLevel)
	return nil
}

// SetFormat sets the log format
func (m *LogManager) SetFormat(format string) error {
	switch strings.ToLower(format) {
	case "json":
		m.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		m.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		return fmt.Errorf("unsupported log format: %s", format)
	}

	return nil
}

// SetOutput sets the log output
func (m *LogManager) SetOutput(output string) error {
	switch strings.ToLower(output) {
	case "stdout":
		m.logger.SetOutput(os.Stdout)
	case "file":
		if err := m.setupFileOutput(); err != nil {
			return err
		}
	case "both":
		if err := m.setupBothOutput(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported log output: %s", output)
	}

	return nil
}

// setupFileOutput sets up file output with lumberjack rotation
func (m *LogManager) setupFileOutput() error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(m.config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	m.lumberjack = &lumberjack.Logger{
		Filename:   m.config.FilePath,
		MaxSize:    m.config.MaxSize,
		MaxBackups: m.config.MaxBackups,
		MaxAge:     m.config.MaxAge,
		Compress:   m.config.Compress,
	}

	m.logger.SetOutput(m.lumberjack)
	return nil
}

// setupBothOutput sets up both stdout and file output
func (m *LogManager) setupBothOutput() error {
	if err := m.setupFileOutput(); err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, m.lumberjack)
	m.logger.SetOutput(multiWriter)
	return nil
}

// Init initializes the default logger (backward compatibility)
func Init(level string) {
	config := &LogConfig{
		Level:      level,
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	defaultManager = NewManager(config).(*LogManager)
}

// InitWithConfig initializes the default logger with config
func InitWithConfig(config *LogConfig) {
	defaultManager = NewManager(config).(*LogManager)
}

// GetDefaultLogger returns the default logger
func GetDefaultLogger() *logrus.Logger {
	if defaultManager == nil {
		Init("info")
	}
	return defaultManager.GetLogger()
}

// Convenience methods for backward compatibility
func Debug(format string, args ...interface{}) {
	if defaultManager == nil {
		Init("info")
	}
	defaultManager.logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	if defaultManager == nil {
		Init("info")
	}
	defaultManager.logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	if defaultManager == nil {
		Init("info")
	}
	defaultManager.logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	if defaultManager == nil {
		Init("info")
	}
	defaultManager.logger.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	if defaultManager == nil {
		Init("info")
	}
	defaultManager.logger.Fatalf(format, args...)
}

// WithFields creates a logger entry with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	if defaultManager == nil {
		Init("info")
	}
	return defaultManager.logger.WithFields(fields)
}

// WithContext creates a logger entry with context
func WithContext(ctx context.Context) *logrus.Entry {
	if defaultManager == nil {
		Init("info")
	}
	return defaultManager.WithContext(ctx)
}
