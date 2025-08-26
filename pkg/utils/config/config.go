package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Manager interface for configuration management
type Manager interface {
	Load(configPath string) error
	GetConfig() *Config
	WatchConfig(callback func(*Config))
	Validate() error
}

// ConfigManager implements the Manager interface
type ConfigManager struct {
	viper  *viper.Viper
	config *Config
}

// Config holds the application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	Server   ServerConfig   `mapstructure:"server"`
	Monitor  MonitorConfig  `mapstructure:"monitor"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Debug   bool   `mapstructure:"debug"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string            `mapstructure:"level"`
	Format     string            `mapstructure:"format"`
	Output     string            `mapstructure:"output"`
	FilePath   string            `mapstructure:"file_path"`
	MaxSize    int               `mapstructure:"max_size"`
	MaxBackups int               `mapstructure:"max_backups"`
	MaxAge     int               `mapstructure:"max_age"`
	Compress   bool              `mapstructure:"compress"`
	Fields     map[string]string `mapstructure:"fields"`
	BufferSize int               `mapstructure:"buffer_size"`
	Async      bool              `mapstructure:"async"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	CORS         CORSConfig    `mapstructure:"cors"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// MonitorConfig holds monitoring configuration
type MonitorConfig struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	PProf      PProfConfig      `mapstructure:"pprof"`
}

// PrometheusConfig holds Prometheus configuration
type PrometheusConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

// PProfConfig holds PProf configuration
type PProfConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	PathPrefix string `mapstructure:"path_prefix"`
	Port       int    `mapstructure:"port"`
}

// NewManager creates a new configuration manager
func NewManager() Manager {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set configuration file settings
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("./")

	// Set environment variable settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &ConfigManager{
		viper: v,
	}
}

// Load loads configuration from file and environment variables
func (m *ConfigManager) Load(configPath string) error {
	if configPath != "" {
		m.viper.SetConfigFile(configPath)
	}

	// Read configuration file
	if err := m.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal configuration
	m.config = &Config{}
	if err := m.viper.Unmarshal(m.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (m *ConfigManager) GetConfig() *Config {
	return m.config
}

// WatchConfig watches for configuration changes
func (m *ConfigManager) WatchConfig(callback func(*Config)) {
	m.viper.WatchConfig()
	m.viper.OnConfigChange(func(e fsnotify.Event) {
		newConfig := &Config{}
		if err := m.viper.Unmarshal(newConfig); err != nil {
			return
		}
		m.config = newConfig
		if callback != nil {
			callback(newConfig)
		}
	})
}

// Validate validates the configuration
func (m *ConfigManager) Validate() error {
	if m.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	// Validate app configuration
	if m.config.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	// Validate database configuration
	if m.config.Database.Type == "" {
		return fmt.Errorf("database type is required")
	}

	// Validate server configuration
	if m.config.Server.Port <= 0 || m.config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", m.config.Server.Port)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "go-http-server")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", true)

	// Database defaults
	v.SetDefault("database.type", "postgresql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "server_tpl")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.dial_timeout", "5s")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.file_path", "logs/app.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 28)
	v.SetDefault("log.compress", true)
	v.SetDefault("log.buffer_size", 1024)
	v.SetDefault("log.async", true)

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.cors.allowed_origins", []string{"http://localhost:3000"})
	v.SetDefault("server.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("server.cors.allowed_headers", []string{"Content-Type", "Authorization"})
	v.SetDefault("server.cors.allow_credentials", true)
	v.SetDefault("server.cors.max_age", 86400)

	// Monitor defaults
	v.SetDefault("monitor.prometheus.enabled", true)
	v.SetDefault("monitor.prometheus.path", "/metrics")
	v.SetDefault("monitor.prometheus.port", 9090)
	v.SetDefault("monitor.pprof.enabled", false)
	v.SetDefault("monitor.pprof.path_prefix", "/debug/pprof")
	v.SetDefault("monitor.pprof.port", 6060)
}

// Convenience methods for backward compatibility
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Env) == "development"
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Env) == "production"
}

func (c *Config) IsTest() bool {
	return strings.ToLower(c.App.Env) == "test"
}

// New creates a new configuration (backward compatibility)
func New() *Config {
	manager := NewManager()
	if err := manager.Load(""); err != nil {
		// Fallback to default configuration
		return &Config{
			App: AppConfig{
				Name:    "go-http-server",
				Version: "1.0.0",
				Env:     "development",
				Debug:   true,
			},
			Server: ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
		}
	}
	return manager.GetConfig()
}
