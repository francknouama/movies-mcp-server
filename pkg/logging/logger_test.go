package logging

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	if logger.Logger == nil {
		t.Fatal("NewLogger().Logger is nil")
	}
}

func TestNewLogger_DefaultLevel(t *testing.T) {
	// Clear any existing LOG_LEVEL env var
	originalLevel := os.Getenv("LOG_LEVEL")
	os.Unsetenv("LOG_LEVEL")
	defer func() {
		if originalLevel != "" {
			os.Setenv("LOG_LEVEL", originalLevel)
		}
	}()

	logger := NewLogger()

	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected default log level to be Info, got: %v", logger.GetLevel())
	}
}

func TestNewLogger_WithEnvLevel(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectedLevel logrus.Level
	}{
		{
			name:          "debug level",
			envValue:      "debug",
			expectedLevel: logrus.DebugLevel,
		},
		{
			name:          "info level",
			envValue:      "info",
			expectedLevel: logrus.InfoLevel,
		},
		{
			name:          "warn level",
			envValue:      "warn",
			expectedLevel: logrus.WarnLevel,
		},
		{
			name:          "error level",
			envValue:      "error",
			expectedLevel: logrus.ErrorLevel,
		},
		{
			name:          "invalid level defaults to info",
			envValue:      "invalid",
			expectedLevel: logrus.InfoLevel,
		},
	}

	originalLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		if originalLevel != "" {
			os.Setenv("LOG_LEVEL", originalLevel)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.envValue)
			logger := NewLogger()

			if logger.GetLevel() != tt.expectedLevel {
				t.Errorf("Expected log level %v, got: %v", tt.expectedLevel, logger.GetLevel())
			}
		})
	}
}

func TestNewLogger_OutputIsStderr(t *testing.T) {
	logger := NewLogger()

	// The logger should output to stderr
	if logger.Out != os.Stderr {
		t.Error("Expected logger output to be os.Stderr")
	}
}

func TestNewLogger_FormatterIsTextFormatter(t *testing.T) {
	logger := NewLogger()

	if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
		t.Errorf("Expected formatter to be *logrus.TextFormatter, got: %T", logger.Formatter)
	}

	formatter := logger.Formatter.(*logrus.TextFormatter)
	if !formatter.FullTimestamp {
		t.Error("Expected FullTimestamp to be true")
	}
}

func TestLogger_WithField(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	entry := logger.WithField("key", "value")
	if entry == nil {
		t.Fatal("WithField() returned nil")
	}

	entry.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "\"key\":\"value\"") {
		t.Errorf("Expected log output to contain field, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain message, got: %s", output)
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	fields := logrus.Fields{
		"field1": "value1",
		"field2": "value2",
		"field3": 123,
	}

	entry := logger.WithFields(fields)
	if entry == nil {
		t.Fatal("WithFields() returned nil")
	}

	entry.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "\"field1\":\"value1\"") {
		t.Errorf("Expected log output to contain field1, got: %s", output)
	}
	if !strings.Contains(output, "\"field2\":\"value2\"") {
		t.Errorf("Expected log output to contain field2, got: %s", output)
	}
	if !strings.Contains(output, "\"field3\":123") {
		t.Errorf("Expected log output to contain field3, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain message, got: %s", output)
	}
}

func TestLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel logrus.Level
		logFunc  func(*Logger, string)
		message  string
		expected bool
	}{
		{
			name:     "debug level logs debug",
			logLevel: logrus.DebugLevel,
			logFunc:  func(l *Logger, msg string) { l.Debug(msg) },
			message:  "debug message",
			expected: true,
		},
		{
			name:     "info level skips debug",
			logLevel: logrus.InfoLevel,
			logFunc:  func(l *Logger, msg string) { l.Debug(msg) },
			message:  "debug message",
			expected: false,
		},
		{
			name:     "info level logs info",
			logLevel: logrus.InfoLevel,
			logFunc:  func(l *Logger, msg string) { l.Info(msg) },
			message:  "info message",
			expected: true,
		},
		{
			name:     "warn level logs warn",
			logLevel: logrus.WarnLevel,
			logFunc:  func(l *Logger, msg string) { l.Warn(msg) },
			message:  "warn message",
			expected: true,
		},
		{
			name:     "warn level skips info",
			logLevel: logrus.WarnLevel,
			logFunc:  func(l *Logger, msg string) { l.Info(msg) },
			message:  "info message",
			expected: false,
		},
		{
			name:     "error level logs error",
			logLevel: logrus.ErrorLevel,
			logFunc:  func(l *Logger, msg string) { l.Error(msg) },
			message:  "error message",
			expected: true,
		},
		{
			name:     "error level skips warn",
			logLevel: logrus.ErrorLevel,
			logFunc:  func(l *Logger, msg string) { l.Warn(msg) },
			message:  "warn message",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.SetOutput(&buf)
			logger.SetLevel(tt.logLevel)

			tt.logFunc(logger, tt.message)

			output := buf.String()
			contains := strings.Contains(output, tt.message)

			if contains != tt.expected {
				if tt.expected {
					t.Errorf("Expected log output to contain message %q, got: %s", tt.message, output)
				} else {
					t.Errorf("Expected log output to not contain message %q, got: %s", tt.message, output)
				}
			}
		})
	}
}

func TestLogger_WithFieldChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Test chaining WithField calls
	entry := logger.WithField("key1", "value1")
	entry = entry.WithField("key2", "value2")
	entry.Info("chained message")

	output := buf.String()
	if !strings.Contains(output, "\"key1\":\"value1\"") {
		t.Errorf("Expected log output to contain key1, got: %s", output)
	}
	if !strings.Contains(output, "\"key2\":\"value2\"") {
		t.Errorf("Expected log output to contain key2, got: %s", output)
	}
	if !strings.Contains(output, "chained message") {
		t.Errorf("Expected log output to contain message, got: %s", output)
	}
}

func TestLogger_DifferentDataTypes(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	fields := logrus.Fields{
		"string":  "text",
		"int":     42,
		"float":   3.14,
		"bool":    true,
		"nil":     nil,
		"slice":   []int{1, 2, 3},
		"map":     map[string]string{"nested": "value"},
		"struct":  struct{ Name string }{Name: "test"},
	}

	logger.WithFields(fields).Info("test with various types")

	output := buf.String()
	if !strings.Contains(output, "\"string\":\"text\"") {
		t.Error("Expected output to contain string field")
	}
	if !strings.Contains(output, "\"int\":42") {
		t.Error("Expected output to contain int field")
	}
	if !strings.Contains(output, "\"bool\":true") {
		t.Error("Expected output to contain bool field")
	}
}

func TestLogger_EmptyFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)

	// Test with empty fields
	logger.WithFields(logrus.Fields{}).Info("message with no fields")

	output := buf.String()
	if !strings.Contains(output, "message with no fields") {
		t.Errorf("Expected log output to contain message, got: %s", output)
	}
}

func TestLogger_LongMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)

	longMessage := strings.Repeat("This is a very long log message. ", 100)
	logger.Info(longMessage)

	output := buf.String()
	if !strings.Contains(output, longMessage) {
		t.Error("Expected log output to contain full long message")
	}
}

func TestLogger_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	message := "Message with special chars: \n\t\"quotes\" 'single' \\backslash"
	logger.Info(message)

	output := buf.String()
	// JSON formatter should escape special characters
	if len(output) == 0 {
		t.Error("Expected non-empty log output")
	}
}

func TestLogger_ConcurrentLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger()
	logger.SetOutput(&buf)

	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			logger.WithField("goroutine", id).Info("concurrent log")
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	output := buf.String()
	if !strings.Contains(output, "concurrent log") {
		t.Error("Expected log output from concurrent goroutines")
	}
}
