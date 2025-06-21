package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	logger := logrus.New()
	
	// Set output to stderr (following MCP pattern)
	logger.SetOutput(os.Stderr)
	
	// Set log level from environment or default to info
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	
	// Set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	
	return &Logger{Logger: logger}
}

// WithField adds a field to the log entry
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the log entry
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}