package protocol

// InitializeRequest represents an MCP initialize request
type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

// ClientCapabilities represents the capabilities of an MCP client
type ClientCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

// ToolsCapability represents tool-related capabilities
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability represents resource-related capabilities
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// PromptsCapability represents prompt-related capabilities
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ClientInfo represents information about the MCP client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolCallRequest represents a tool call request
type ToolCallRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ResourceReadRequest represents a resource read request
type ResourceReadRequest struct {
	URI string `json:"uri"`
}

// ResourceSubscribeRequest represents a resource subscribe request
type ResourceSubscribeRequest struct {
	URI string `json:"uri"`
}

// ResourceUnsubscribeRequest represents a resource unsubscribe request
type ResourceUnsubscribeRequest struct {
	URI string `json:"uri"`
}

// PromptGetRequest represents a prompt get request
type PromptGetRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ListRequest represents a generic list request (for tools, resources, prompts)
type ListRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// PaginatedRequest represents a request with pagination
type PaginatedRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// NotificationRequest represents a notification request
type NotificationRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// ResourceUpdatedNotification represents a resource updated notification
type ResourceUpdatedNotification struct {
	URI string `json:"uri"`
}

// ToolsListChangedNotification represents a tools list changed notification
type ToolsListChangedNotification struct{}

// PromptsListChangedNotification represents a prompts list changed notification
type PromptsListChangedNotification struct{}

// CommonRequestParams represents common parameters that can be present in any request
type CommonRequestParams struct {
	ID      interface{} `json:"id,omitempty"`
	Meta    interface{} `json:"_meta,omitempty"`
	Timeout int         `json:"timeout,omitempty"` // Timeout in seconds
}