package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Logger wraps slog.Logger with additional context and functionality
type Logger struct {
	logger *slog.Logger
	ctx    context.Context
}

// New creates a new structured logger
func New(level LogLevel) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	// Create handler with JSON output for production
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().UTC().Format(time.RFC3339))
			}
			return a
		},
	})

	logger := slog.New(handler)
	return &Logger{
		logger: logger,
		ctx:    context.Background(),
	}
}

// WithContext returns a logger with additional context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger,
		ctx:    ctx,
	}
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}
	return &Logger{
		logger: l.logger.With(args...),
		ctx:    l.ctx,
	}
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.DebugContext(l.ctx, msg, fields...)
}

// Info logs an info message with optional fields
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.InfoContext(l.ctx, msg, fields...)
}

// Warn logs a warning message with optional fields
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.WarnContext(l.ctx, msg, fields...)
}

// Error logs an error message with optional fields
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logger.ErrorContext(l.ctx, msg, fields...)
}

// LogRequest logs an incoming MCP request
func (l *Logger) LogRequest(requestID interface{}, method string, duration time.Duration) {
	l.Info("mcp_request",
		"request_id", requestID,
		"method", method,
		"duration_ms", duration.Milliseconds(),
		"component", "mcp_server",
	)
}

// LogError logs an error with additional context
func (l *Logger) LogError(err error, context string, fields ...interface{}) {
	allFields := append([]interface{}{
		"error", err.Error(),
		"context", context,
	}, fields...)
	l.Error("error_occurred", allFields...)
}

// LogDatabaseOperation logs database operations
func (l *Logger) LogDatabaseOperation(operation string, table string, duration time.Duration, err error) {
	fields := []interface{}{
		"operation", operation,
		"table", table,
		"duration_ms", duration.Milliseconds(),
		"component", "database",
	}
	
	if err != nil {
		fields = append(fields, "error", err.Error(), "success", false)
		l.Error("database_operation", fields...)
	} else {
		fields = append(fields, "success", true)
		l.Debug("database_operation", fields...)
	}
}

// LogImageOperation logs image processing operations
func (l *Logger) LogImageOperation(operation string, size int64, mimeType string, duration time.Duration, err error) {
	fields := []interface{}{
		"operation", operation,
		"image_size_bytes", size,
		"mime_type", mimeType,
		"duration_ms", duration.Milliseconds(),
		"component", "image_processor",
	}
	
	if err != nil {
		fields = append(fields, "error", err.Error(), "success", false)
		l.Error("image_operation", fields...)
	} else {
		fields = append(fields, "success", true)
		l.Info("image_operation", fields...)
	}
}

// LogServerStart logs server startup
func (l *Logger) LogServerStart(version string, config map[string]interface{}) {
	l.Info("server_start",
		"version", version,
		"config", config,
		"component", "mcp_server",
	)
}

// LogServerShutdown logs server shutdown
func (l *Logger) LogServerShutdown(reason string) {
	l.Info("server_shutdown",
		"reason", reason,
		"component", "mcp_server",
	)
}

// LogPerformanceMetric logs performance metrics
func (l *Logger) LogPerformanceMetric(metric string, value float64, unit string, tags map[string]string) {
	fields := []interface{}{
		"metric", metric,
		"value", value,
		"unit", unit,
		"component", "metrics",
	}
	
	for key, val := range tags {
		fields = append(fields, key, val)
	}
	
	l.Info("performance_metric", fields...)
}

// LogHealthCheck logs health check results
func (l *Logger) LogHealthCheck(component string, status string, duration time.Duration, details map[string]interface{}) {
	fields := []interface{}{
		"component", component,
		"status", status,
		"duration_ms", duration.Milliseconds(),
		"check_type", "health_check",
	}
	
	if details != nil {
		detailsJSON, _ := json.Marshal(details)
		fields = append(fields, "details", string(detailsJSON))
	}
	
	if status == "healthy" {
		l.Info("health_check", fields...)
	} else {
		l.Warn("health_check", fields...)
	}
}

// LogSecurity logs security-related events
func (l *Logger) LogSecurity(event string, severity string, details map[string]interface{}) {
	fields := []interface{}{
		"security_event", event,
		"severity", severity,
		"component", "security",
	}
	
	for key, value := range details {
		fields = append(fields, key, value)
	}
	
	switch severity {
	case "high", "critical":
		l.Error("security_event", fields...)
	case "medium":
		l.Warn("security_event", fields...)
	default:
		l.Info("security_event", fields...)
	}
}

// RequestLogger creates a middleware-style logger for tracking requests
type RequestLogger struct {
	logger *Logger
}

// NewRequestLogger creates a new request logger
func NewRequestLogger(logger *Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// LogRequest logs request details with timing
func (rl *RequestLogger) LogRequest(requestID interface{}, method string, startTime time.Time, err error) {
	duration := time.Since(startTime)
	
	if err != nil {
		rl.logger.Error("mcp_request_failed",
			"request_id", requestID,
			"method", method,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
			"component", "mcp_server",
		)
	} else {
		rl.logger.Info("mcp_request_completed",
			"request_id", requestID,
			"method", method,
			"duration_ms", duration.Milliseconds(),
			"component", "mcp_server",
		)
	}
}

// ContextKey is used for logger context keys
type ContextKey string

const (
	LoggerContextKey ContextKey = "logger"
	RequestIDKey     ContextKey = "request_id"
)

// FromContext extracts logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*Logger); ok {
		return logger
	}
	// Return default logger if none in context
	return New(LevelInfo)
}

// ToContext adds logger to context
func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) interface{} {
	if id := ctx.Value(RequestIDKey); id != nil {
		return id
	}
	return "unknown"
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID interface{}) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// Helper function to format duration for logging
func FormatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.2fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}