package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Upstream UpstreamConfig
	Cache    CacheConfig
	Log      LogConfig
}

type ServerConfig struct {
	Host   string
	Port   int
	Prefix string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type UpstreamConfig struct {
	BaseURL   string
	Timeout   time.Duration
	RateLimit int // requests per second
}

type CacheConfig struct {
	TTL time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

// Global config instance
var cfg *Config

// Init initializes the configuration from file and/or environment
func Init(cfgFile string) error {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in common locations
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.mchess")
		viper.AddConfigPath("/etc/mchess")
	}

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("No config file found, using defaults")
		} else {
			return fmt.Errorf("error reading config: %w", err)
		}
	} else {
		slog.Info("Using config file", "path", viper.ConfigFileUsed())
	}

	// Parse into struct
	cfg = &Config{}
	if err := parseConfig(); err != nil {
		return err
	}

	return nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.prefix", "/api")

	// Database defaults
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "mchess")
	viper.SetDefault("db.password", "mchess")
	viper.SetDefault("db.name", "mchess")

	// Upstream defaults
	viper.SetDefault("upstream.baseUrl", "https://member.schack.se/public/api/v1")
	viper.SetDefault("upstream.timeout", "30s")
	viper.SetDefault("upstream.rateLimit", 10)

	// Cache defaults
	viper.SetDefault("cache.ttl", "24h")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
}

func parseConfig() error {
	cfg.Server = ServerConfig{
		Host:   viper.GetString("server.host"),
		Port:   viper.GetInt("server.port"),
		Prefix: viper.GetString("server.prefix"),
	}

	cfg.DB = DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Name:     viper.GetString("db.name"),
	}

	timeout, err := time.ParseDuration(viper.GetString("upstream.timeout"))
	if err != nil {
		timeout = 30 * time.Second
	}
	cfg.Upstream = UpstreamConfig{
		BaseURL:   viper.GetString("upstream.baseUrl"),
		Timeout:   timeout,
		RateLimit: viper.GetInt("upstream.rateLimit"),
	}

	ttl, err := time.ParseDuration(viper.GetString("cache.ttl"))
	if err != nil {
		ttl = 24 * time.Hour
	}
	cfg.Cache = CacheConfig{
		TTL: ttl,
	}

	cfg.Log = LogConfig{
		Level:  viper.GetString("log.level"),
		Format: viper.GetString("log.format"),
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		panic("config not initialized - call Init first")
	}
	return cfg
}

// ServerAddr returns the server address in host:port format
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// DBConnectionString returns the PostgreSQL connection string
func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name)
}

// SetupLogger configures the global slog logger based on config
func (c *Config) SetupLogger() {
	var level slog.Level
	switch c.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if c.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
