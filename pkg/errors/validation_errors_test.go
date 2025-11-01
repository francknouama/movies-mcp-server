package errors

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *ValidationError
		wantMsg string
	}{
		{
			name: "error with field",
			err: &ValidationError{
				Message: "must be a valid email",
				Field:   "email",
			},
			wantMsg: "validation error for field 'email': must be a valid email",
		},
		{
			name: "error without field",
			err: &ValidationError{
				Message: "invalid input data",
				Field:   "",
			},
			wantMsg: "validation error: invalid input data",
		},
		{
			name: "error with field and data",
			err: &ValidationError{
				Message: "must be between 1 and 100",
				Field:   "age",
				Data:    map[string]interface{}{"min": 1, "max": 100},
			},
			wantMsg: "validation error for field 'age': must be between 1 and 100",
		},
		{
			name: "error with empty message",
			err: &ValidationError{
				Message: "",
				Field:   "name",
			},
			wantMsg: "validation error for field 'name': ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		field       string
		data        map[string]interface{}
		wantMessage string
		wantField   string
	}{
		{
			name:        "create with all fields",
			message:     "invalid value",
			field:       "username",
			data:        map[string]interface{}{"min_length": 3},
			wantMessage: "invalid value",
			wantField:   "username",
		},
		{
			name:        "create without data",
			message:     "required field",
			field:       "password",
			data:        nil,
			wantMessage: "required field",
			wantField:   "password",
		},
		{
			name:        "create with empty field",
			message:     "general validation error",
			field:       "",
			data:        nil,
			wantMessage: "general validation error",
			wantField:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.field, tt.data)

			if err.Message != tt.wantMessage {
				t.Errorf("NewValidationError() Message = %v, want %v", err.Message, tt.wantMessage)
			}

			if err.Field != tt.wantField {
				t.Errorf("NewValidationError() Field = %v, want %v", err.Field, tt.wantField)
			}

			if tt.data != nil && err.Data == nil {
				t.Error("NewValidationError() Data should not be nil when data is provided")
			}

			if tt.data == nil && err.Data != nil {
				t.Error("NewValidationError() Data should be nil when nil is provided")
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is validation error",
			err:  NewValidationError("test error", "field", nil),
			want: true,
		},
		{
			name: "is validation error pointer",
			err:  &ValidationError{Message: "test", Field: "field"},
			want: true,
		},
		{
			name: "not validation error - standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "not validation error - nil",
			err:  nil,
			want: false,
		},
		{
			name: "not validation error - custom error",
			err:  &GodogMCPError{Code: 100, Message: "test"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_AsStandardError(t *testing.T) {
	err := NewValidationError("test error", "field", nil)

	// Test that it satisfies the error interface
	var _ error = err

	// Test that we can call Error() method
	if errMsg := err.Error(); errMsg == "" {
		t.Error("Error() should return non-empty string")
	}
}

func TestValidationError_WithVariousData(t *testing.T) {
	tests := []struct {
		name     string
		dataType string
		data     map[string]interface{}
	}{
		{
			name:     "nil data",
			dataType: "nil",
			data:     nil,
		},
		{
			name:     "empty map",
			dataType: "empty_map",
			data:     map[string]interface{}{},
		},
		{
			name:     "string values",
			dataType: "strings",
			data:     map[string]interface{}{"key1": "value1", "key2": "value2"},
		},
		{
			name:     "mixed types",
			dataType: "mixed",
			data: map[string]interface{}{
				"string":  "value",
				"int":     42,
				"float":   3.14,
				"bool":    true,
				"nil":     nil,
				"slice":   []int{1, 2, 3},
				"map":     map[string]string{"nested": "value"},
			},
		},
		{
			name:     "numeric values",
			dataType: "numeric",
			data: map[string]interface{}{
				"min":   0,
				"max":   100,
				"value": 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError("test message", "test_field", tt.data)

			if tt.data == nil && err.Data != nil {
				t.Error("Data should be nil when nil is provided")
			}

			if tt.data != nil && err.Data == nil {
				t.Error("Data should not be nil when data is provided")
			}

			if tt.data != nil && len(tt.data) != len(err.Data) {
				t.Errorf("Data length = %d, want %d", len(err.Data), len(tt.data))
			}
		})
	}
}

func TestValidationError_MultipleFields(t *testing.T) {
	// Test creating multiple validation errors for different fields
	errors := []*ValidationError{
		NewValidationError("required", "username", nil),
		NewValidationError("invalid format", "email", nil),
		NewValidationError("too short", "password", map[string]interface{}{"min_length": 8}),
	}

	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errors))
	}

	// Verify each error is distinct
	fields := make(map[string]bool)
	for _, err := range errors {
		if fields[err.Field] {
			t.Errorf("Duplicate field: %s", err.Field)
		}
		fields[err.Field] = true
	}
}

func TestValidationError_ErrorMessageFormat(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		field     string
		expectIn  string // substring that should be in error message
		notExpect string // substring that should not be in error message
	}{
		{
			name:      "with field contains field name",
			message:   "is required",
			field:     "email",
			expectIn:  "email",
			notExpect: "",
		},
		{
			name:      "without field does not mention field",
			message:   "validation failed",
			field:     "",
			expectIn:  "validation error",
			notExpect: "field",
		},
		{
			name:      "contains message text",
			message:   "must be positive",
			field:     "amount",
			expectIn:  "must be positive",
			notExpect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.field, nil)
			errMsg := err.Error()

			if tt.expectIn != "" {
				if !contains(errMsg, tt.expectIn) {
					t.Errorf("Error message should contain '%s', got: %s", tt.expectIn, errMsg)
				}
			}

			if tt.notExpect != "" {
				if contains(errMsg, tt.notExpect) {
					t.Errorf("Error message should not contain '%s', got: %s", tt.notExpect, errMsg)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
