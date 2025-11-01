package errors

import (
	"testing"
)

func TestGodogMCPError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *GodogMCPError
		wantMsg string
	}{
		{
			name: "basic error",
			err: &GodogMCPError{
				Code:    -32600,
				Message: "Invalid request",
			},
			wantMsg: "godog-mcp error -32600: Invalid request",
		},
		{
			name: "error with data",
			err: &GodogMCPError{
				Code:    -32602,
				Message: "Invalid parameters",
				Data:    map[string]string{"param": "value"},
			},
			wantMsg: "godog-mcp error -32602: Invalid parameters",
		},
		{
			name: "custom error code",
			err: &GodogMCPError{
				Code:    -40001,
				Message: "Godog not found",
				Data:    nil,
			},
			wantMsg: "godog-mcp error -40001: Godog not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("GodogMCPError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestNewGodogMCPError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		data     any
		wantCode int
		wantMsg  string
	}{
		{
			name:     "create error without data",
			code:     -32700,
			message:  "Parse error",
			data:     nil,
			wantCode: -32700,
			wantMsg:  "Parse error",
		},
		{
			name:     "create error with string data",
			code:     -32603,
			message:  "Internal error",
			data:     "additional info",
			wantCode: -32603,
			wantMsg:  "Internal error",
		},
		{
			name:     "create error with map data",
			code:     -40001,
			message:  "Custom error",
			data:     map[string]interface{}{"key": "value"},
			wantCode: -40001,
			wantMsg:  "Custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewGodogMCPError(tt.code, tt.message, tt.data)

			if err.Code != tt.wantCode {
				t.Errorf("NewGodogMCPError() Code = %v, want %v", err.Code, tt.wantCode)
			}

			if err.Message != tt.wantMsg {
				t.Errorf("NewGodogMCPError() Message = %v, want %v", err.Message, tt.wantMsg)
			}

			if tt.data != nil && err.Data == nil {
				t.Error("NewGodogMCPError() Data should not be nil")
			}
		})
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"ParseError", ParseError, -32700},
		{"InvalidRequest", InvalidRequest, -32600},
		{"MethodNotFound", MethodNotFound, -32601},
		{"InvalidParams", InvalidParams, -32602},
		{"InternalError", InternalError, -32603},
		{"GodogNotFound", GodogNotFound, -40001},
		{"FeatureParseError", FeatureParseError, -40002},
		{"TestExecutionError", TestExecutionError, -40003},
		{"StepDefinitionError", StepDefinitionError, -40004},
		{"ReportGenerationError", ReportGenerationError, -40005},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("Error code %s = %d, want %d", tt.name, tt.code, tt.want)
			}
		})
	}
}

func TestNewParseError(t *testing.T) {
	message := "Failed to parse JSON"
	err := NewParseError(message)

	if err.Code != ParseError {
		t.Errorf("NewParseError() Code = %d, want %d", err.Code, ParseError)
	}

	if err.Message != message {
		t.Errorf("NewParseError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewParseError() Data should be nil")
	}
}

func TestNewInvalidRequest(t *testing.T) {
	message := "Request is invalid"
	err := NewInvalidRequest(message)

	if err.Code != InvalidRequest {
		t.Errorf("NewInvalidRequest() Code = %d, want %d", err.Code, InvalidRequest)
	}

	if err.Message != message {
		t.Errorf("NewInvalidRequest() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewInvalidRequest() Data should be nil")
	}
}

func TestNewMethodNotFound(t *testing.T) {
	method := "unknown_method"
	err := NewMethodNotFound(method)

	if err.Code != MethodNotFound {
		t.Errorf("NewMethodNotFound() Code = %d, want %d", err.Code, MethodNotFound)
	}

	expectedMsg := "Method not found: unknown_method"
	if err.Message != expectedMsg {
		t.Errorf("NewMethodNotFound() Message = %v, want %v", err.Message, expectedMsg)
	}

	if err.Data != nil {
		t.Error("NewMethodNotFound() Data should be nil")
	}
}

func TestNewInvalidParams(t *testing.T) {
	message := "Invalid parameters provided"
	err := NewInvalidParams(message)

	if err.Code != InvalidParams {
		t.Errorf("NewInvalidParams() Code = %d, want %d", err.Code, InvalidParams)
	}

	if err.Message != message {
		t.Errorf("NewInvalidParams() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewInvalidParams() Data should be nil")
	}
}

func TestNewInternalError(t *testing.T) {
	message := "Internal server error"
	err := NewInternalError(message)

	if err.Code != InternalError {
		t.Errorf("NewInternalError() Code = %d, want %d", err.Code, InternalError)
	}

	if err.Message != message {
		t.Errorf("NewInternalError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewInternalError() Data should be nil")
	}
}

func TestNewGodogNotFound(t *testing.T) {
	message := "Godog feature not found"
	err := NewGodogNotFound(message)

	if err.Code != GodogNotFound {
		t.Errorf("NewGodogNotFound() Code = %d, want %d", err.Code, GodogNotFound)
	}

	if err.Message != message {
		t.Errorf("NewGodogNotFound() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewGodogNotFound() Data should be nil")
	}
}

func TestNewFeatureParseError(t *testing.T) {
	message := "Failed to parse feature file"
	err := NewFeatureParseError(message)

	if err.Code != FeatureParseError {
		t.Errorf("NewFeatureParseError() Code = %d, want %d", err.Code, FeatureParseError)
	}

	if err.Message != message {
		t.Errorf("NewFeatureParseError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewFeatureParseError() Data should be nil")
	}
}

func TestNewTestExecutionError(t *testing.T) {
	message := "Test execution failed"
	err := NewTestExecutionError(message)

	if err.Code != TestExecutionError {
		t.Errorf("NewTestExecutionError() Code = %d, want %d", err.Code, TestExecutionError)
	}

	if err.Message != message {
		t.Errorf("NewTestExecutionError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewTestExecutionError() Data should be nil")
	}
}

func TestNewStepDefinitionError(t *testing.T) {
	message := "Step definition error"
	err := NewStepDefinitionError(message)

	if err.Code != StepDefinitionError {
		t.Errorf("NewStepDefinitionError() Code = %d, want %d", err.Code, StepDefinitionError)
	}

	if err.Message != message {
		t.Errorf("NewStepDefinitionError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewStepDefinitionError() Data should be nil")
	}
}

func TestNewReportGenerationError(t *testing.T) {
	message := "Failed to generate report"
	err := NewReportGenerationError(message)

	if err.Code != ReportGenerationError {
		t.Errorf("NewReportGenerationError() Code = %d, want %d", err.Code, ReportGenerationError)
	}

	if err.Message != message {
		t.Errorf("NewReportGenerationError() Message = %v, want %v", err.Message, message)
	}

	if err.Data != nil {
		t.Error("NewReportGenerationError() Data should be nil")
	}
}

func TestGodogMCPError_AsStandardError(t *testing.T) {
	err := NewInternalError("test error")

	// Test that it satisfies the error interface
	var _ error = err

	// Test that we can call Error() method
	if errMsg := err.Error(); errMsg == "" {
		t.Error("Error() should return non-empty string")
	}
}

func TestGodogMCPError_WithVariousDataTypes(t *testing.T) {
	tests := []struct {
		name     string
		dataType string
		data     any
	}{
		{
			name:     "nil data",
			dataType: "nil",
			data:     nil,
		},
		{
			name:     "string data",
			dataType: "string",
			data:     "error details",
		},
		{
			name:     "int data",
			dataType: "int",
			data:     42,
		},
		{
			name:     "map data",
			dataType: "map",
			data:     map[string]interface{}{"key": "value", "count": 10},
		},
		{
			name:     "slice data",
			dataType: "slice",
			data:     []string{"error1", "error2"},
		},
		{
			name:     "struct data",
			dataType: "struct",
			data: struct {
				Field1 string
				Field2 int
			}{"test", 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewGodogMCPError(InternalError, "test message", tt.data)

			if tt.data == nil && err.Data != nil {
				t.Error("Data should be nil when nil is provided")
			}

			if tt.data != nil && err.Data == nil {
				t.Error("Data should not be nil when data is provided")
			}
		})
	}
}
