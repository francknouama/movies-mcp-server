package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Image    ImageConfig
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	SSLMode       string
	MaxOpenConns  int
	MaxIdleConns  int
	ConnMaxLifetime time.Duration
	MigrationsPath string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	LogLevel string
	Timeout  time.Duration
}

// ImageConfig holds image-related configuration
type ImageConfig struct {
	MaxSize         int64
	AllowedTypes    []string
	EnableThumbnails bool
	ThumbnailSize   string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:          getEnv("DB_HOST", "localhost"),
			Port:          getEnvAsInt("DB_PORT", 5432),
			Name:          getEnv("DB_NAME", "movies_mcp"),
			User:          getEnv("DB_USER", "movies_user"),
			Password:      getEnv("DB_PASSWORD", "movies_password"),
			SSLMode:       getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:  getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:  getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "1h"),
			MigrationsPath: getEnv("MIGRATIONS_PATH", "file://migrations"),
		},
		Server: ServerConfig{
			LogLevel: getEnv("LOG_LEVEL", "info"),
			Timeout:  getEnvAsDuration("SERVER_TIMEOUT", "30s"),
		},
		Image: ImageConfig{
			MaxSize:         getEnvAsInt64("MAX_IMAGE_SIZE", 5*1024*1024), // 5MB default
			AllowedTypes:    getEnvAsStringSlice("ALLOWED_IMAGE_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
			EnableThumbnails: getEnvAsBool("ENABLE_THUMBNAILS", true),
			ThumbnailSize:   getEnv("THUMBNAIL_SIZE", "200x200"),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configuration is present and valid
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Image.MaxSize <= 0 {
		return fmt.Errorf("MAX_IMAGE_SIZE must be positive")
	}
	if len(c.Image.AllowedTypes) == 0 {
		return fmt.Errorf("ALLOWED_IMAGE_TYPES cannot be empty")
	}
	return nil
}

// ConnectionString returns a PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// Return default if parsing fails
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}