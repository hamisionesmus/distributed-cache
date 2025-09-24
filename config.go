package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server" toml:"server" yaml:"server"`
	Cache    CacheConfig    `json:"cache" toml:"cache" yaml:"cache"`
	Cluster  ClusterConfig  `json:"cluster" toml:"cluster" yaml:"cluster"`
	Storage  StorageConfig  `json:"storage" toml:"storage" yaml:"storage"`
	Metrics  MetricsConfig  `json:"metrics" toml:"metrics" yaml:"metrics"`
	Security SecurityConfig `json:"security" toml:"security" yaml:"security"`
	Logging  LoggingConfig  `json:"logging" toml:"logging" yaml:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host            string        `json:"host" toml:"host" yaml:"host"`
	Port            int           `json:"port" toml:"port" yaml:"port"`
	HTTPPort        int           `json:"http_port" toml:"http_port" yaml:"http_port"`
	ReadTimeout     time.Duration `json:"read_timeout" toml:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" toml:"write_timeout" yaml:"write_timeout"`
	MaxConnections  int           `json:"max_connections" toml:"max_connections" yaml:"max_connections"`
	EnableHTTP      bool          `json:"enable_http" toml:"enable_http" yaml:"enable_http"`
	EnableTLS       bool          `json:"enable_tls" toml:"enable_tls" yaml:"enable_tls"`
	TLSCertFile     string        `json:"tls_cert_file" toml:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile      string        `json:"tls_key_file" toml:"tls_key_file" yaml:"tls_key_file"`
	EnableCORS      bool          `json:"enable_cors" toml:"enable_cors" yaml:"enable_cors"`
	CORSOrigins     []string      `json:"cors_origins" toml:"cors_origins" yaml:"cors_origins"`
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	MaxMemory         int64         `json:"max_memory" toml:"max_memory" yaml:"max_memory"`
	DefaultTTL        time.Duration `json:"default_ttl" toml:"default_ttl" yaml:"default_ttl"`
	CleanupInterval   time.Duration `json:"cleanup_interval" toml:"cleanup_interval" yaml:"cleanup_interval"`
	EvictionPolicy    string        `json:"eviction_policy" toml:"eviction_policy" yaml:"eviction_policy"`
	EnableCompression bool          `json:"enable_compression" toml:"enable_compression" yaml:"enable_compression"`
	CompressionLevel  int           `json:"compression_level" toml:"compression_level" yaml:"compression_level"`
	ShardCount        int           `json:"shard_count" toml:"shard_count" yaml:"shard_count"`
	EnableMetrics     bool          `json:"enable_metrics" toml:"enable_metrics" yaml:"enable_metrics"`
}

// ClusterConfig holds clustering configuration
type ClusterConfig struct {
	Enabled         bool     `json:"enabled" toml:"enabled" yaml:"enabled"`
	NodeID          string   `json:"node_id" toml:"node_id" yaml:"node_id"`
	Seeds           []string `json:"seeds" toml:"seeds" yaml:"seeds"`
	Port            int      `json:"port" toml:"port" yaml:"port"`
	GossipInterval  time.Duration `json:"gossip_interval" toml:"gossip_interval" yaml:"gossip_interval"`
	ProbeInterval   time.Duration `json:"probe_interval" toml:"probe_interval" yaml:"probe_interval"`
	ProbeTimeout    time.Duration `json:"probe_timeout" toml:"probe_timeout" yaml:"probe_timeout"`
	SuspicionMult   int      `json:"suspicion_mult" toml:"suspicion_mult" yaml:"suspicion_mult"`
	ReconnectIntvl  time.Duration `json:"reconnect_interval" toml:"reconnect_interval" yaml:"reconnect_interval"`
	ReconnectTimeout time.Duration `json:"reconnect_timeout" toml:"reconnect_timeout" yaml:"reconnect_timeout"`
}

// StorageConfig holds persistence configuration
type StorageConfig struct {
	Enabled           bool          `json:"enabled" toml:"enabled" yaml:"enabled"`
	Type              string        `json:"type" toml:"type" yaml:"type"`
	Path              string        `json:"path" toml:"path" yaml:"path"`
	SyncInterval      time.Duration `json:"sync_interval" toml:"sync_interval" yaml:"sync_interval"`
	MaxFileSize       int64         `json:"max_file_size" toml:"max_file_size" yaml:"max_file_size"`
	Compression       bool          `json:"compression" toml:"compression" yaml:"compression"`
	Encryption        bool          `json:"encryption" toml:"encryption" yaml:"encryption"`
	EncryptionKey     string        `json:"encryption_key" toml:"encryption_key" yaml:"encryption_key"`
	BackupEnabled     bool          `json:"backup_enabled" toml:"backup_enabled" yaml:"backup_enabled"`
	BackupInterval    time.Duration `json:"backup_interval" toml:"backup_interval" yaml:"backup_interval"`
	BackupRetention   int           `json:"backup_retention" toml:"backup_retention" yaml:"backup_retention"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled         bool          `json:"enabled" toml:"enabled" yaml:"enabled"`
	Interval        time.Duration `json:"interval" toml:"interval" yaml:"interval"`
	RetentionPeriod time.Duration `json:"retention_period" toml:"retention_period" yaml:"retention_period"`
	PrometheusPort  int           `json:"prometheus_port" toml:"prometheus_port" yaml:"prometheus_port"`
	EnableHistogram bool          `json:"enable_histogram" toml:"enable_histogram" yaml:"enable_histogram"`
	Buckets         []float64     `json:"buckets" toml:"buckets" yaml:"buckets"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EnableAuth       bool     `json:"enable_auth" toml:"enable_auth" yaml:"enable_auth"`
	AuthType         string   `json:"auth_type" toml:"auth_type" yaml:"auth_type"`
	JWTSecret        string   `json:"jwt_secret" toml:"jwt_secret" yaml:"jwt_secret"`
	JWTExpiry        time.Duration `json:"jwt_expiry" toml:"jwt_expiry" yaml:"jwt_expiry"`
	EnableACL        bool     `json:"enable_acl" toml:"enable_acl" yaml:"enable_acl"`
	ACLFile          string   `json:"acl_file" toml:"acl_file" yaml:"acl_file"`
	EnableTLS        bool     `json:"enable_tls" toml:"enable_tls" yaml:"enable_tls"`
	TLSCertFile      string   `json:"tls_cert_file" toml:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile       string   `json:"tls_key_file" toml:"tls_key_file" yaml:"tls_key_file"`
	EnableRateLimit  bool     `json:"enable_rate_limit" toml:"enable_rate_limit" yaml:"enable_rate_limit"`
	RateLimitRPM     int      `json:"rate_limit_rpm" toml:"rate_limit_rpm" yaml:"rate_limit_rpm"`
	EnableIPFilter   bool     `json:"enable_ip_filter" toml:"enable_ip_filter" yaml:"enable_ip_filter"`
	AllowedIPs       []string `json:"allowed_ips" toml:"allowed_ips" yaml:"allowed_ips"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `json:"level" toml:"level" yaml:"level"`
	Format     string `json:"format" toml:"format" yaml:"format"`
	Output     string `json:"output" toml:"output" yaml:"output"`
	File       string `json:"file" toml:"file" yaml:"file"`
	MaxSize    int64  `json:"max_size" toml:"max_size" yaml:"max_size"`
	MaxFiles   int    `json:"max_files" toml:"max_files" yaml:"max_files"`
	Compress   bool   `json:"compress" toml:"compress" yaml:"compress"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:           "0.0.0.0",
			Port:           6379,
			HTTPPort:       8080,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxConnections: 10000,
			EnableHTTP:     true,
			EnableTLS:      false,
			EnableCORS:     true,
			CORSOrigins:    []string{"*"},
		},
		Cache: CacheConfig{
			MaxMemory:         512 * 1024 * 1024, // 512MB
			DefaultTTL:        24 * time.Hour,
			CleanupInterval:   10 * time.Minute,
			EvictionPolicy:    "lru",
			EnableCompression: true,
			CompressionLevel:  6,
			ShardCount:        16,
			EnableMetrics:     true,
		},
		Cluster: ClusterConfig{
			Enabled:         false,
			GossipInterval:  1 * time.Second,
			ProbeInterval:   5 * time.Second,
			ProbeTimeout:    3 * time.Second,
			SuspicionMult:   5,
			ReconnectIntvl:  10 * time.Second,
			ReconnectTimeout: 6 * time.Second,
		},
		Storage: StorageConfig{
			Enabled:         false,
			Type:            "aof",
			Path:            "./data",
			SyncInterval:    1 * time.Second,
			MaxFileSize:     1024 * 1024 * 1024, // 1GB
			Compression:     true,
			BackupEnabled:   false,
			BackupInterval:  24 * time.Hour,
			BackupRetention: 7,
		},
		Metrics: MetricsConfig{
			Enabled:         true,
			Interval:        10 * time.Second,
			RetentionPeriod: 7 * 24 * time.Hour,
			PrometheusPort:  9090,
			EnableHistogram: true,
			Buckets:         []float64{.005, .01, .025, .05, .1, .25, .5, 1.0, 2.5, 5.0, 10.0},
		},
		Security: SecurityConfig{
			EnableAuth:      false,
			AuthType:        "jwt",
			JWTExpiry:       24 * time.Hour,
			EnableACL:       false,
			EnableRateLimit: true,
			RateLimitRPM:    1000,
		},
		Logging: LoggingConfig{
			Level:    "info",
			Format:   "json",
			Output:   "stdout",
			MaxSize:  100 * 1024 * 1024, // 100MB
			MaxFiles: 5,
			Compress: true,
		},
	}
}

// LoadConfig loads configuration from file and command line flags
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Parse command line flags
	var configFile string
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.StringVar(&config.Server.Host, "host", config.Server.Host, "Server host")
	flag.IntVar(&config.Server.Port, "port", config.Server.Port, "Server port")
	flag.IntVar(&config.Server.HTTPPort, "http-port", config.Server.HTTPPort, "HTTP server port")
	flag.Int64Var(&config.Cache.MaxMemory, "max-memory", config.Cache.MaxMemory, "Maximum memory usage")
	flag.BoolVar(&config.Cluster.Enabled, "cluster", config.Cluster.Enabled, "Enable clustering")
	flag.Parse()

	// Load from file if specified
	if configFile != "" {
		if err := loadFromFile(config, configFile); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from a file
func loadFromFile(config *Config, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(filename, ".json"):
		return json.Unmarshal(data, config)
	case strings.HasSuffix(filename, ".toml"):
		return toml.Unmarshal(data, config)
	case strings.HasSuffix(filename, ".yaml"), strings.HasSuffix(filename, ".yml"):
		return yaml.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported config file format")
	}
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) {
	// Server config
	if v := os.Getenv("CACHE_HOST"); v != "" {
		config.Server.Host = v
	}
	if v := os.Getenv("CACHE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			config.Server.Port = port
		}
	}
	if v := os.Getenv("CACHE_HTTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			config.Server.HTTPPort = port
		}
	}

	// Cache config
	if v := os.Getenv("CACHE_MAX_MEMORY"); v != "" {
		if mem, err := strconv.ParseInt(v, 10, 64); err == nil {
			config.Cache.MaxMemory = mem
		}
	}

	// Cluster config
	if v := os.Getenv("CACHE_CLUSTER_ENABLED"); v != "" {
		if enabled, err := strconv.ParseBool(v); err == nil {
			config.Cluster.Enabled = enabled
		}
	}
	if v := os.Getenv("CACHE_CLUSTER_SEEDS"); v != "" {
		config.Cluster.Seeds = strings.Split(v, ",")
	}

	// Security config
	if v := os.Getenv("CACHE_AUTH_ENABLED"); v != "" {
		if enabled, err := strconv.ParseBool(v); err == nil {
			config.Security.EnableAuth = enabled
		}
	}
	if v := os.Getenv("CACHE_JWT_SECRET"); v != "" {
		config.Security.JWTSecret = v
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Server.HTTPPort < 1 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTPPort)
	}

	// Validate cache config
	if c.Cache.MaxMemory < 1024*1024 { // 1MB minimum
		return fmt.Errorf("max memory too small: %d", c.Cache.MaxMemory)
	}
	if c.Cache.ShardCount < 1 {
		return fmt.Errorf("shard count must be at least 1")
	}

	// Validate cluster config
	if c.Cluster.Enabled {
		if len(c.Cluster.Seeds) == 0 {
			return fmt.Errorf("cluster seeds required when clustering is enabled")
		}
	}

	// Validate security config
	if c.Security.EnableAuth {
		if c.Security.JWTSecret == "" {
			return fmt.Errorf("JWT secret required when auth is enabled")
		}
		if c.Security.JWTExpiry < time.Minute {
			return fmt.Errorf("JWT expiry too short")
		}
	}

	return nil
}

// Save saves the configuration to a file
func (c *Config) Save(filename string) error {
	var data []byte
	var err error

	switch {
	case strings.HasSuffix(filename, ".json"):
		data, err = json.MarshalIndent(c, "", "  ")
	case strings.HasSuffix(filename, ".toml"):
		data, err = toml.Marshal(*c)
	case strings.HasSuffix(filename, ".yaml"), strings.HasSuffix(filename, ".yml"):
		data, err = yaml.Marshal(c)
	default:
		return fmt.Errorf("unsupported config file format")
	}

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}