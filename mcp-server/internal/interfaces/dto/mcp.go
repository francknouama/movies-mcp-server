package dto

import (
	"github.com/francknouama/movies-mcp-server/shared-mcp/pkg/protocol"
)

// Type aliases for shared protocol types - Phase 2 BDD Remediation
type JSONRPCRequest = protocol.JSONRPCRequest
type JSONRPCResponse = protocol.JSONRPCResponse
type JSONRPCError = protocol.JSONRPCError

// MCP Protocol Types - using shared library
type InitializeRequest = protocol.InitializeRequest
type ClientCapabilities = protocol.ClientCapabilities
type ClientInfo = protocol.ClientInfo
type InitializeResponse = protocol.InitializeResponse
type ServerCapabilities = protocol.ServerCapabilities
type ToolsCapability = protocol.ToolsCapability
type ResourcesCapability = protocol.ResourcesCapability
type PromptsCapability = protocol.PromptsCapability
type ServerInfo = protocol.ServerInfo

// Tool and related types - using shared library
type Tool = protocol.Tool
type InputSchema = protocol.InputSchema
type SchemaProperty = protocol.SchemaProperty
type ToolsListResponse = protocol.ToolsListResponse
type ToolCallRequest = protocol.ToolCallRequest
type ToolCallResponse = protocol.ToolCallResponse
type ContentBlock = protocol.ContentBlock
type ImageSource = protocol.ImageSource

// Resource and related types - using shared library
type Resource = protocol.Resource
type ResourcesListResponse = protocol.ResourcesListResponse
type ResourceReadRequest = protocol.ResourceReadRequest
type ResourceReadResponse = protocol.ResourceReadResponse
type ResourceContent = protocol.ResourceContent

// Prompt and related types - using shared library
type Prompt = protocol.Prompt
type PromptArgument = protocol.PromptArgument
type PromptsListResponse = protocol.PromptsListResponse
type PromptGetRequest = protocol.PromptGetRequest
type PromptGetResponse = protocol.PromptGetResponse
type PromptMessage = protocol.PromptMessage
type PromptMessageContent = protocol.PromptMessageContent

// Note: ResourceTemplate not yet in shared library - keep local definition
// ResourceTemplate represents an MCP resource template
type ResourceTemplate struct {
	URITemplate string `json:"uriTemplate"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourceTemplatesListResponse represents the response for resources/templates/list
type ResourceTemplatesListResponse struct {
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
}

// Error codes - using shared library constants
const (
	ParseError     = protocol.ParseError
	InvalidRequest = protocol.InvalidRequest
	MethodNotFound = protocol.MethodNotFound
	InvalidParams  = protocol.InvalidParams
	InternalError  = protocol.InternalError
)

// NewJSONRPCError creates a new JSON-RPC error
func NewJSONRPCError(code int, message string, data interface{}) *JSONRPCError {
	return &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewJSONRPCResponse creates a new JSON-RPC response
func NewJSONRPCResponse(id interface{}, result interface{}, err *JSONRPCError) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: protocol.JSONRPC2Version,
		ID:      id,
		Result:  result,
		Error:   err,
	}
}
