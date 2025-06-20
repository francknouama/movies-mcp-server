package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"movies-mcp-server/pkg/logging"
	"movies-mcp-server/internal/models"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// JSON-RPC standard error codes
	ParseError     ErrorCode = "PARSE_ERROR"
	InvalidRequest ErrorCode = "INVALID_REQUEST"
	MethodNotFound ErrorCode = "METHOD_NOT_FOUND"
	InvalidParams  ErrorCode = "INVALID_PARAMS"
	InternalError  ErrorCode = "INTERNAL_ERROR"
	
	// Application-specific error codes
	ValidationError    ErrorCode = "VALIDATION_ERROR"
	NotFoundError      ErrorCode = "NOT_FOUND_ERROR"
	ConflictError      ErrorCode = "CONFLICT_ERROR"
	AuthError          ErrorCode = "AUTH_ERROR"
	RateLimitError     ErrorCode = "RATE_LIMIT_ERROR"
	DatabaseError      ErrorCode = "DATABASE_ERROR"
	ImageProcessingError ErrorCode = "IMAGE_PROCESSING_ERROR"
	ResourceError      ErrorCode = "RESOURCE_ERROR"
	TimeoutError       ErrorCode = "TIMEOUT_ERROR"
	ServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// Severity represents error severity levels
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// ApplicationError represents a structured application error
type ApplicationError struct {
	Code        ErrorCode              `json:"code"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Severity    Severity               `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	RequestID   string                 `json:"request_id,omitempty"`
	StackTrace  []string               `json:"stack_trace,omitempty"`
	Cause       error                  `json:"-"` // Original error, not serialized
	Retryable   bool                   `json:"retryable"`
	Component   string                 `json:"component,omitempty"`
}

// Error implements the error interface
func (e *ApplicationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ApplicationError) Unwrap() error {
	return e.Cause
}

// NewApplicationError creates a new application error
func NewApplicationError(code ErrorCode, message string) *ApplicationError {
	return &ApplicationError{
		Code:      code,
		Message:   message,
		Severity:  SeverityMedium,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Retryable: false,
	}
}

// WithCause adds a cause to the error
func (e *ApplicationError) WithCause(cause error) *ApplicationError {
	e.Cause = cause
	return e
}

// WithDetails adds details to the error
func (e *ApplicationError) WithDetails(details map[string]interface{}) *ApplicationError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// WithSeverity sets the error severity
func (e *ApplicationError) WithSeverity(severity Severity) *ApplicationError {
	e.Severity = severity
	return e
}

// WithRequestID sets the request ID
func (e *ApplicationError) WithRequestID(requestID string) *ApplicationError {
	e.RequestID = requestID
	return e
}

// WithComponent sets the component where the error occurred
func (e *ApplicationError) WithComponent(component string) *ApplicationError {
	e.Component = component
	return e
}

// WithStackTrace adds stack trace information
func (e *ApplicationError) WithStackTrace() *ApplicationError {
	e.StackTrace = getStackTrace()
	return e
}

// AsRetryable marks the error as retryable
func (e *ApplicationError) AsRetryable() *ApplicationError {
	e.Retryable = true
	return e
}

// ToJSONRPCError converts to JSON-RPC error format
func (e *ApplicationError) ToJSONRPCError() *models.JSONRPCError {
	// Map application error codes to JSON-RPC codes
	var rpcCode int
	switch e.Code {
	case ParseError:
		rpcCode = models.ParseError
	case InvalidRequest:
		rpcCode = models.InvalidRequest
	case MethodNotFound:
		rpcCode = models.MethodNotFound
	case InvalidParams, ValidationError:
		rpcCode = models.InvalidParams
	case InternalError, DatabaseError, ImageProcessingError, ServiceUnavailable:
		rpcCode = models.InternalError
	default:
		rpcCode = models.InternalError
	}

	return &models.JSONRPCError{
		Code:    rpcCode,
		Message: e.Message,
		Data:    e.Details,
	}
}

// getStackTrace captures the current stack trace
func getStackTrace() []string {
	var traces []string
	for i := 2; i < 10; i++ { // Skip first 2 frames (this function and caller)
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		// Simplify file path to relative path
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			file = file[idx+1:]
		}
		
		traces = append(traces, fmt.Sprintf("%s:%d", file, line))
	}
	return traces
}

// ErrorHandler provides centralized error handling
type ErrorHandler struct {
	logger *logging.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *logging.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// Handle processes an error and returns appropriate response
func (eh *ErrorHandler) Handle(err error, requestID interface{}) *models.JSONRPCError {
	if err == nil {
		return nil
	}

	var appError *ApplicationError
	
	// Check if it's already an ApplicationError
	if ae, ok := err.(*ApplicationError); ok {
		appError = ae
	} else {
		// Wrap unknown errors
		appError = eh.WrapError(err, InternalError, "An unexpected error occurred")
	}

	// Set request ID if not already set
	if appError.RequestID == "" && requestID != nil {
		appError.RequestID = fmt.Sprintf("%v", requestID)
	}

	// Log the error
	eh.LogError(appError)

	return appError.ToJSONRPCError()
}

// WrapError wraps a generic error into an ApplicationError
func (eh *ErrorHandler) WrapError(err error, code ErrorCode, message string) *ApplicationError {
	return NewApplicationError(code, message).
		WithCause(err).
		WithStackTrace().
		WithSeverity(eh.determineSeverity(code))
}

// LogError logs an application error
func (eh *ErrorHandler) LogError(err *ApplicationError) {
	fields := []interface{}{
		"error_code", string(err.Code),
		"error_message", err.Message,
		"severity", string(err.Severity),
		"retryable", err.Retryable,
		"timestamp", err.Timestamp.Format(time.RFC3339),
	}

	if err.RequestID != "" {
		fields = append(fields, "request_id", err.RequestID)
	}

	if err.Component != "" {
		fields = append(fields, "component", err.Component)
	}

	if err.Cause != nil {
		fields = append(fields, "cause", err.Cause.Error())
	}

	if len(err.Details) > 0 {
		detailsJSON, _ := json.Marshal(err.Details)
		fields = append(fields, "details", string(detailsJSON))
	}

	if len(err.StackTrace) > 0 {
		fields = append(fields, "stack_trace", strings.Join(err.StackTrace, " -> "))
	}

	// Log based on severity
	switch err.Severity {
	case SeverityLow, SeverityMedium:
		eh.logger.Warn("application_error", fields...)
	case SeverityHigh, SeverityCritical:
		eh.logger.Error("application_error", fields...)
	default:
		eh.logger.Error("application_error", fields...)
	}
}

// determineSeverity determines error severity based on error code
func (eh *ErrorHandler) determineSeverity(code ErrorCode) Severity {
	switch code {
	case ParseError, InvalidRequest, InvalidParams, ValidationError:
		return SeverityLow
	case MethodNotFound, NotFoundError, ConflictError:
		return SeverityLow
	case AuthError, RateLimitError:
		return SeverityMedium
	case DatabaseError, TimeoutError:
		return SeverityHigh
	case InternalError, ServiceUnavailable, ImageProcessingError:
		return SeverityHigh
	default:
		return SeverityMedium
	}
}

// Common error constructors

// NewValidationError creates a validation error
func NewValidationError(message string, field string, value interface{}) *ApplicationError {
	return NewApplicationError(ValidationError, message).
		WithDetails(map[string]interface{}{
			"field": field,
			"value": value,
		}).
		WithSeverity(SeverityLow)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, id interface{}) *ApplicationError {
	return NewApplicationError(NotFoundError, fmt.Sprintf("%s not found", resource)).
		WithDetails(map[string]interface{}{
			"resource": resource,
			"id":       id,
		}).
		WithSeverity(SeverityLow)
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, cause error) *ApplicationError {
	return NewApplicationError(DatabaseError, "Database operation failed").
		WithCause(cause).
		WithDetails(map[string]interface{}{
			"operation": operation,
		}).
		WithSeverity(SeverityHigh).
		WithComponent("database").
		AsRetryable()
}

// NewImageProcessingError creates an image processing error
func NewImageProcessingError(operation string, cause error) *ApplicationError {
	return NewApplicationError(ImageProcessingError, "Image processing failed").
		WithCause(cause).
		WithDetails(map[string]interface{}{
			"operation": operation,
		}).
		WithSeverity(SeverityMedium).
		WithComponent("image_processor")
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(operation string, timeout time.Duration) *ApplicationError {
	return NewApplicationError(TimeoutError, "Operation timed out").
		WithDetails(map[string]interface{}{
			"operation":     operation,
			"timeout_ms":    timeout.Milliseconds(),
		}).
		WithSeverity(SeverityHigh).
		AsRetryable()
}

// NewResourceError creates a resource error
func NewResourceError(uri string, operation string, cause error) *ApplicationError {
	return NewApplicationError(ResourceError, "Resource operation failed").
		WithCause(cause).
		WithDetails(map[string]interface{}{
			"uri":       uri,
			"operation": operation,
		}).
		WithSeverity(SeverityMedium).
		WithComponent("resources")
}

// ErrorRecovery provides error recovery utilities
type ErrorRecovery struct {
	logger *logging.Logger
}

// NewErrorRecovery creates a new error recovery instance
func NewErrorRecovery(logger *logging.Logger) *ErrorRecovery {
	return &ErrorRecovery{logger: logger}
}

// Recover recovers from panics and converts them to errors
func (er *ErrorRecovery) Recover() {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered: %v", r)
		appError := NewApplicationError(InternalError, "Internal server error").
			WithCause(err).
			WithSeverity(SeverityCritical).
			WithStackTrace()
		
		er.logger.Error("panic_recovered",
			"panic_value", r,
			"stack_trace", strings.Join(appError.StackTrace, " -> "),
		)
	}
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if appErr, ok := err.(*ApplicationError); ok {
		return appErr.Retryable
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := err.(*ApplicationError); ok {
		return appErr.Code
	}
	return InternalError
}

// GetErrorSeverity extracts the error severity from an error
func GetErrorSeverity(err error) Severity {
	if appErr, ok := err.(*ApplicationError); ok {
		return appErr.Severity
	}
	return SeverityMedium
}