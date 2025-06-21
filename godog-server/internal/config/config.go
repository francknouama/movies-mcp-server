package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the Godog MCP server
type Config struct {
	// Server configuration
	LogLevel      string        `json:"log_level"`
	ServerTimeout time.Duration `json:"server_timeout"`

	// Godog configuration
	GodogBinary string `json:"godog_binary"`
	FeaturesDir string `json:"features_dir"`
	StepDefsDir string `json:"step_defs_dir"`
	ReportsDir  string `json:"reports_dir"`

	// Test execution configuration
	MaxParallel    int           `json:"max_parallel"`
	DefaultTimeout time.Duration `json:"default_timeout"`
	RetryCount     int           `json:"retry_count"`

	// Report configuration
	ReportFormats []string `json:"report_formats"`
	KeepReports   int      `json:"keep_reports"` // Number of reports to keep
}

// Load creates a new Config from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		// Server defaults
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "info"),
		ServerTimeout: getDurationOrDefault("SERVER_TIMEOUT", "30s"),

		// Godog defaults
		GodogBinary: getEnvOrDefault("GODOG_BINARY", "godog"),
		FeaturesDir: getEnvOrDefault("FEATURES_DIR", "./features"),
		StepDefsDir: getEnvOrDefault("STEP_DEFS_DIR", "./step_definitions"),
		ReportsDir:  getEnvOrDefault("REPORTS_DIR", "./reports"),

		// Test execution defaults
		MaxParallel:    getIntOrDefault("MAX_PARALLEL", 4),
		DefaultTimeout: getDurationOrDefault("DEFAULT_TIMEOUT", "5m"),
		RetryCount:     getIntOrDefault("RETRY_COUNT", 0),

		// Report defaults
		ReportFormats: getStringSliceOrDefault("REPORT_FORMATS", []string{"cucumber", "pretty"}),
		KeepReports:   getIntOrDefault("KEEP_REPORTS", 10),
	}

	return cfg, nil
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationOrDefault(key string, defaultValue string) time.Duration {
	value := getEnvOrDefault(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// Fallback to parsing the default
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getStringSliceOrDefault(key string, defaultValue []string) []string {
	// For now, return default. In a real implementation, you'd parse a
	// comma-separated string from the environment
	return defaultValue
}
