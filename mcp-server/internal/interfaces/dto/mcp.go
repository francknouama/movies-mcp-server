package dto

import (
	"encoding/json"
)

// JSON-RPC 2.0 Base Types

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol Types

// InitializeRequest represents the initialize method parameters
type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

// ClientCapabilities represents client capabilities
type ClientCapabilities struct {
	// Add specific capabilities as needed
}

// ClientInfo represents client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResponse represents the initialize method response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

// ToolsCapability indicates tool support
type ToolsCapability struct{}

// ResourcesCapability indicates resource support
type ResourcesCapability struct {
	Subscribe bool `json:"subscribe,omitempty"`
}

// PromptsCapability indicates prompt support
type PromptsCapability struct{}

// ServerInfo represents server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the JSON schema for tool input
type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// ToolsListResponse represents the response for tools/list
type ToolsListResponse struct {
	Tools []Tool `json:"tools"`
}

// ToolCallRequest represents a tool call request
type ToolCallRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolCallResponse represents a tool call response
type ToolCallResponse struct {
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents a content block in tool response
type ContentBlock struct {
	Type   string       `json:"type"`
	Text   string       `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
}

// ImageSource represents an image source for content blocks
type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	Data      string `json:"data"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string                 `json:"uri"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// ResourcesListResponse represents the response for resources/list
type ResourcesListResponse struct {
	Resources []Resource `json:"resources"`
}

// ResourceReadRequest represents a resource read request
type ResourceReadRequest struct {
	URI string `json:"uri"`
}

// ResourceReadResponse represents a resource read response
type ResourceReadResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents resource content
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"` // base64 encoded binary data
}

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

// Prompt represents an MCP prompt template
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument represents an argument for a prompt
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptsListResponse represents the response for prompts/list
type PromptsListResponse struct {
	Prompts []Prompt `json:"prompts"`
}

// PromptGetRequest represents a prompt get request
type PromptGetRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// PromptGetResponse represents a prompt get response
type PromptGetResponse struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// PromptMessage represents a message in a prompt
type PromptMessage struct {
	Role    string               `json:"role"`
	Content PromptMessageContent `json:"content"`
}

// PromptMessageContent can be either a string or structured content
type PromptMessageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
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
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
		Error:   err,
	}
}
