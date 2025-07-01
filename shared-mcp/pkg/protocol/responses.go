package protocol

// InitializeResponse represents an MCP initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities represents the capabilities of an MCP server
type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

// ServerInfo represents information about the MCP server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolsListResponse represents a tools list response
type ToolsListResponse struct {
	Tools      []Tool  `json:"tools"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the JSON schema for tool input
type InputSchema struct {
	Type       string                            `json:"type"`
	Properties map[string]SchemaProperty        `json:"properties,omitempty"`
	Required   []string                         `json:"required,omitempty"`
	Additional map[string]interface{}           `json:"additionalProperties,omitempty"`
}

// SchemaProperty represents a property in a JSON schema
type SchemaProperty struct {
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Pattern     string                 `json:"pattern,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
	Items       *SchemaProperty        `json:"items,omitempty"`
	Properties  map[string]SchemaProperty `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

// ToolCallResponse represents a tool call response
type ToolCallResponse struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock represents a block of content in a response
type ContentBlock struct {
	Type   string       `json:"type"`
	Text   string       `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
	Data   interface{}  `json:"data,omitempty"`
}

// ImageSource represents an image source in content
type ImageSource struct {
	Type     string `json:"type"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// ResourcesListResponse represents a resources list response
type ResourcesListResponse struct {
	Resources  []Resource `json:"resources"`
	NextCursor *string    `json:"nextCursor,omitempty"`
}

// Resource represents an MCP resource definition
type Resource struct {
	URI         string                 `json:"uri"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// ResourceReadResponse represents a resource read response
type ResourceReadResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents the content of a resource
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     []byte `json:"blob,omitempty"`
}

// PromptsListResponse represents a prompts list response
type PromptsListResponse struct {
	Prompts    []Prompt `json:"prompts"`
	NextCursor *string  `json:"nextCursor,omitempty"`
}

// Prompt represents an MCP prompt definition
type Prompt struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Arguments   []PromptArgument     `json:"arguments,omitempty"`
}

// PromptArgument represents an argument for a prompt
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptGetResponse represents a prompt get response
type PromptGetResponse struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// PromptMessage represents a message in a prompt
type PromptMessage struct {
	Role    string               `json:"role"`
	Content []ContentBlock       `json:"content"`
}

// PromptMessageContent represents content in a prompt message (backward compatibility)
type PromptMessageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PaginatedResponse represents a response with pagination
type PaginatedResponse struct {
	NextCursor *string `json:"nextCursor,omitempty"`
	Total      *int    `json:"total,omitempty"`
}

// EmptyResponse represents an empty response (for notifications, etc.)
type EmptyResponse struct{}

// CommonResponseMeta represents common metadata that can be present in any response
type CommonResponseMeta struct {
	RequestID string      `json:"requestId,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Duration  int64       `json:"duration,omitempty"` // Duration in milliseconds
	Meta      interface{} `json:"_meta,omitempty"`
}