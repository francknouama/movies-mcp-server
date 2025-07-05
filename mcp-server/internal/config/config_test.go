package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save current environment
	oldEnv := os.Environ()
	defer func() {
		// Restore environment
		os.Clearenv()
		for _, e := range oldEnv {
			pair := splitEnvVar(e)
			os.Setenv(pair[0], pair[1])
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			want: &Config{
				Database: DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Name:            "movies_mcp",
					User:            "movies_user",
					Password:        "movies_password",
					SSLMode:         "disable",
					MaxOpenConns:    25,
					MaxIdleConns:    5,
					ConnMaxLifetime: time.Hour,
					MigrationsPath:  "file://migrations",
				},
				Server: ServerConfig{
					LogLevel: "info",
					Timeout:  30 * time.Second,
				},
				Image: ImageConfig{
					MaxSize:          5 * 1024 * 1024,
					AllowedTypes:     []string{"image/jpeg", "image/png", "image/webp"},
					EnableThumbnails: true,
					ThumbnailSize:    "200x200",
				},
			},
			wantErr: false,
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"DB_HOST":              "db.example.com",
				"DB_PORT":              "5433",
				"DB_NAME":              "custom_db",
				"DB_USER":              "custom_user",
				"DB_PASSWORD":          "custom_pass",
				"DB_SSLMODE":           "require",
				"DB_MAX_OPEN_CONNS":    "50",
				"DB_MAX_IDLE_CONNS":    "10",
				"DB_CONN_MAX_LIFETIME": "2h",
				"MIGRATIONS_PATH":      "file://custom/migrations",
				"LOG_LEVEL":            "debug",
				"SERVER_TIMEOUT":       "1m",
				"MAX_IMAGE_SIZE":       "10485760",
				"ALLOWED_IMAGE_TYPES":  "image/jpeg,image/png",
				"ENABLE_THUMBNAILS":    "false",
				"THUMBNAIL_SIZE":       "300x300",
			},
			want: &Config{
				Database: DatabaseConfig{
					Host:            "db.example.com",
					Port:            5433,
					Name:            "custom_db",
					User:            "custom_user",
					Password:        "custom_pass",
					SSLMode:         "require",
					MaxOpenConns:    50,
					MaxIdleConns:    10,
					ConnMaxLifetime: 2 * time.Hour,
					MigrationsPath:  "file://custom/migrations",
				},
				Server: ServerConfig{
					LogLevel: "debug",
					Timeout:  time.Minute,
				},
				Image: ImageConfig{
					MaxSize:          10485760,
					AllowedTypes:     []string{"image/jpeg", "image/png"},
					EnableThumbnails: false,
					ThumbnailSize:    "300x300",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"DB_PORT": "99999",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty host",
			envVars: map[string]string{
				"DB_HOST": "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid image size",
			envVars: map[string]string{
				"MAX_IMAGE_SIZE": "-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty allowed types",
			envVars: map[string]string{
				"ALLOWED_IMAGE_TYPES": "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 5432,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: &Config{
				Database: DatabaseConfig{
					Port: 5432,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "DB_HOST is required",
		},
		{
			name: "invalid port - zero",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 0,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "DB_PORT must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 99999,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "DB_PORT must be between 1 and 65535",
		},
		{
			name: "missing db name",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 5432,
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "DB_NAME is required",
		},
		{
			name: "missing db user",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 5432,
					Name: "test_db",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "DB_USER is required",
		},
		{
			name: "invalid image size",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 5432,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      0,
					AllowedTypes: []string{"image/jpeg"},
				},
			},
			wantErr: true,
			errMsg:  "MAX_IMAGE_SIZE must be positive",
		},
		{
			name: "empty allowed types",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 5432,
					Name: "test_db",
					User: "test_user",
				},
				Image: ImageConfig{
					MaxSize:      1024,
					AllowedTypes: []string{},
				},
			},
			wantErr: true,
			errMsg:  "ALLOWED_IMAGE_TYPES cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Config.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDatabaseConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name   string
		config DatabaseConfig
		want   string
	}{
		{
			name: "basic connection string",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "test_user",
				Password: "test_pass",
				Name:     "test_db",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=test_user password=test_pass dbname=test_db sslmode=disable",
		},
		{
			name: "with special characters",
			config: DatabaseConfig{
				Host:     "db.example.com",
				Port:     5433,
				User:     "user@example",
				Password: "p@ss!word",
				Name:     "my-db",
				SSLMode:  "require",
			},
			want: "host=db.example.com port=5433 user=user@example password=p@ss!word dbname=my-db sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.ConnectionString(); got != tt.want {
				t.Errorf("DatabaseConfig.ConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvHelpers(t *testing.T) {
	// Save current environment
	oldEnv := os.Environ()
	defer func() {
		// Restore environment
		os.Clearenv()
		for _, e := range oldEnv {
			pair := splitEnvVar(e)
			os.Setenv(pair[0], pair[1])
		}
	}()

	t.Run("getEnv", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		if got := getEnv("MISSING_VAR", "default"); got != "default" {
			t.Errorf("getEnv() = %v, want %v", got, "default")
		}

		// Test existing value
		os.Setenv("EXISTING_VAR", "value")
		if got := getEnv("EXISTING_VAR", "default"); got != "value" {
			t.Errorf("getEnv() = %v, want %v", got, "value")
		}
	})

	t.Run("getEnvAsInt", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		if got := getEnvAsInt("MISSING_VAR", 42); got != 42 {
			t.Errorf("getEnvAsInt() = %v, want %v", got, 42)
		}

		// Test valid int
		os.Setenv("INT_VAR", "123")
		if got := getEnvAsInt("INT_VAR", 42); got != 123 {
			t.Errorf("getEnvAsInt() = %v, want %v", got, 123)
		}

		// Test invalid int
		os.Setenv("INVALID_INT", "not-a-number")
		if got := getEnvAsInt("INVALID_INT", 42); got != 42 {
			t.Errorf("getEnvAsInt() = %v, want %v", got, 42)
		}
	})

	t.Run("getEnvAsInt64", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		if got := getEnvAsInt64("MISSING_VAR", 42); got != 42 {
			t.Errorf("getEnvAsInt64() = %v, want %v", got, 42)
		}

		// Test valid int64
		os.Setenv("INT64_VAR", "9223372036854775807")
		if got := getEnvAsInt64("INT64_VAR", 42); got != 9223372036854775807 {
			t.Errorf("getEnvAsInt64() = %v, want %v", got, 9223372036854775807)
		}
	})

	t.Run("getEnvAsBool", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		if got := getEnvAsBool("MISSING_VAR", true); got != true {
			t.Errorf("getEnvAsBool() = %v, want %v", got, true)
		}

		// Test true values
		for _, v := range []string{"true", "True", "TRUE", "1"} {
			os.Setenv("BOOL_VAR", v)
			if got := getEnvAsBool("BOOL_VAR", false); got != true {
				t.Errorf("getEnvAsBool() with %s = %v, want %v", v, got, true)
			}
		}

		// Test false values
		for _, v := range []string{"false", "False", "FALSE", "0"} {
			os.Setenv("BOOL_VAR", v)
			if got := getEnvAsBool("BOOL_VAR", true); got != false {
				t.Errorf("getEnvAsBool() with %s = %v, want %v", v, got, false)
			}
		}
	})

	t.Run("getEnvAsDuration", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		if got := getEnvAsDuration("MISSING_VAR", "1h"); got != time.Hour {
			t.Errorf("getEnvAsDuration() = %v, want %v", got, time.Hour)
		}

		// Test valid duration
		os.Setenv("DURATION_VAR", "30m")
		if got := getEnvAsDuration("DURATION_VAR", "1h"); got != 30*time.Minute {
			t.Errorf("getEnvAsDuration() = %v, want %v", got, 30*time.Minute)
		}

		// Test invalid duration (should use default)
		os.Setenv("INVALID_DURATION", "not-a-duration")
		if got := getEnvAsDuration("INVALID_DURATION", "1h"); got != time.Hour {
			t.Errorf("getEnvAsDuration() = %v, want %v", got, time.Hour)
		}
	})

	t.Run("getEnvAsStringSlice", func(t *testing.T) {
		os.Clearenv()

		// Test default value
		defaultSlice := []string{"a", "b", "c"}
		if got := getEnvAsStringSlice("MISSING_VAR", defaultSlice); !reflect.DeepEqual(got, defaultSlice) {
			t.Errorf("getEnvAsStringSlice() = %v, want %v", got, defaultSlice)
		}

		// Test comma-separated values
		os.Setenv("SLICE_VAR", "one,two,three")
		want := []string{"one", "two", "three"}
		if got := getEnvAsStringSlice("SLICE_VAR", defaultSlice); !reflect.DeepEqual(got, want) {
			t.Errorf("getEnvAsStringSlice() = %v, want %v", got, want)
		}

		// Test single value
		os.Setenv("SINGLE_VAR", "single")
		want = []string{"single"}
		if got := getEnvAsStringSlice("SINGLE_VAR", defaultSlice); !reflect.DeepEqual(got, want) {
			t.Errorf("getEnvAsStringSlice() = %v, want %v", got, want)
		}
	})
}

// Helper function to split environment variable
func splitEnvVar(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env, ""}
}

