// Package config provides configuration management for the movies MCP server.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Image    ImageConfig
}

// DatabaseConfig holds database-specific configuration.
type DatabaseConfig struct {
	Name            string        // SQLite database file path
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationsPath  string
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	LogLevel string
	Timeout  time.Duration
}

// ImageConfig holds image-related configuration.
type ImageConfig struct {
	MaxSize          int64
	AllowedTypes     []string
	EnableThumbnails bool
	ThumbnailSize    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Name:            getEnv("DB_NAME", "movies.db"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 1),  // SQLite works best with 1
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 1),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "0"),
			MigrationsPath:  getEnv("MIGRATIONS_PATH", "file://migrations"),
		},
		Server: ServerConfig{
			LogLevel: getEnv("LOG_LEVEL", "info"),
			Timeout:  getEnvAsDuration("SERVER_TIMEOUT", "30s"),
		},
		Image: ImageConfig{
			MaxSize:          getEnvAsInt64("MAX_IMAGE_SIZE", 5*1024*1024), // 5MB default
			AllowedTypes:     getEnvAsStringSlice("ALLOWED_IMAGE_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
			EnableThumbnails: getEnvAsBool("ENABLE_THUMBNAILS", true),
			ThumbnailSize:    getEnv("THUMBNAIL_SIZE", "200x200"),
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
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Image.MaxSize <= 0 {
		return fmt.Errorf("MAX_IMAGE_SIZE must be positive")
	}
	if len(c.Image.AllowedTypes) == 0 {
		return fmt.Errorf("ALLOWED_IMAGE_TYPES cannot be empty")
	}
	return nil
}

// ConnectionString returns the SQLite database file path
func (c *DatabaseConfig) ConnectionString() string {
	return c.Name
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
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
	duration, err := time.ParseDuration(defaultValue)
	if err != nil {
		// If default value is also invalid, return 0
		return 0
	}
	return duration
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	value, exists := os.LookupEnv(key)
	if exists {
		if value == "" {
			return []string{}
		}
		return strings.Split(value, ",")
	}
	return defaultValue
}
