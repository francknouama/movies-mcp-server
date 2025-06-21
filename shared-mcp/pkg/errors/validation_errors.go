package errors

import "fmt"

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Message string                 `json:"message"`
	Field   string                 `json:"field"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message, field string, data map[string]interface{}) *ValidationError {
	return &ValidationError{
		Message: message,
		Field:   field,
		Data:    data,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
