// Package errors provides error handling and custom error types for the MCP server.
package errors

import (
	"fmt"
)

// GodogMCPError represents a custom error for the Godog MCP server.
type GodogMCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *GodogMCPError) Error() string {
	return fmt.Sprintf("godog-mcp error %d: %s", e.Code, e.Message)
}

// NewGodogMCPError creates a new GodogMCPError.
func NewGodogMCPError(code int, message string, data any) *GodogMCPError {
	return &GodogMCPError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Common error codes following MCP specification.
const (
	// Standard JSON-RPC error codes.
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603

	// Custom Godog MCP error codes.
	GodogNotFound         = -40001
	FeatureParseError     = -40002
	TestExecutionError    = -40003
	StepDefinitionError   = -40004
	ReportGenerationError = -40005
)

// Pre-defined error constructors.
// NewParseError creates a new parse error.
func NewParseError(message string) *GodogMCPError {
	return NewGodogMCPError(ParseError, message, nil)
}

// NewInvalidRequest creates a new invalid request error.
func NewInvalidRequest(message string) *GodogMCPError {
	return NewGodogMCPError(InvalidRequest, message, nil)
}

// NewMethodNotFound creates a new method not found error.
func NewMethodNotFound(method string) *GodogMCPError {
	return NewGodogMCPError(MethodNotFound, fmt.Sprintf("Method not found: %s", method), nil)
}

// NewInvalidParams creates a new invalid params error.
func NewInvalidParams(message string) *GodogMCPError {
	return NewGodogMCPError(InvalidParams, message, nil)
}

// NewInternalError creates a new internal error.
func NewInternalError(message string) *GodogMCPError {
	return NewGodogMCPError(InternalError, message, nil)
}

// NewGodogNotFound creates a new godog not found error.
func NewGodogNotFound(message string) *GodogMCPError {
	return NewGodogMCPError(GodogNotFound, message, nil)
}

// NewFeatureParseError creates a new feature parse error.
func NewFeatureParseError(message string) *GodogMCPError {
	return NewGodogMCPError(FeatureParseError, message, nil)
}

// NewTestExecutionError creates a new test execution error.
func NewTestExecutionError(message string) *GodogMCPError {
	return NewGodogMCPError(TestExecutionError, message, nil)
}

// NewStepDefinitionError creates a new step definition error.
func NewStepDefinitionError(message string) *GodogMCPError {
	return NewGodogMCPError(StepDefinitionError, message, nil)
}

// NewReportGenerationError creates a new report generation error.
func NewReportGenerationError(message string) *GodogMCPError {
	return NewGodogMCPError(ReportGenerationError, message, nil)
}
