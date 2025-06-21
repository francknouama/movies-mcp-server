package errors

import (
	"fmt"
)

// GodogMCPError represents a custom error for the Godog MCP server
type GodogMCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface
func (e *GodogMCPError) Error() string {
	return fmt.Sprintf("godog-mcp error %d: %s", e.Code, e.Message)
}

// NewGodogMCPError creates a new GodogMCPError
func NewGodogMCPError(code int, message string, data any) *GodogMCPError {
	return &GodogMCPError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Common error codes following MCP specification
const (
	// Standard JSON-RPC error codes
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603

	// Custom Godog MCP error codes
	GodogNotFound      = -40001
	FeatureParseError  = -40002
	TestExecutionError = -40003
	StepDefinitionError = -40004
	ReportGenerationError = -40005
)

// Pre-defined error constructors
func NewParseError(message string) *GodogMCPError {
	return NewGodogMCPError(ParseError, message, nil)
}

func NewInvalidRequest(message string) *GodogMCPError {
	return NewGodogMCPError(InvalidRequest, message, nil)
}

func NewMethodNotFound(method string) *GodogMCPError {
	return NewGodogMCPError(MethodNotFound, fmt.Sprintf("Method not found: %s", method), nil)
}

func NewInvalidParams(message string) *GodogMCPError {
	return NewGodogMCPError(InvalidParams, message, nil)
}

func NewInternalError(message string) *GodogMCPError {
	return NewGodogMCPError(InternalError, message, nil)
}

func NewGodogNotFound(message string) *GodogMCPError {
	return NewGodogMCPError(GodogNotFound, message, nil)
}

func NewFeatureParseError(message string) *GodogMCPError {
	return NewGodogMCPError(FeatureParseError, message, nil)
}

func NewTestExecutionError(message string) *GodogMCPError {
	return NewGodogMCPError(TestExecutionError, message, nil)
}

func NewStepDefinitionError(message string) *GodogMCPError {
	return NewGodogMCPError(StepDefinitionError, message, nil)
}

func NewReportGenerationError(message string) *GodogMCPError {
	return NewGodogMCPError(ReportGenerationError, message, nil)
}