package protocol

// Protocol version constants
const (
	JSONRPC2Version = "2.0"
	MCPVersion      = "2024-11-05"
)

// Standard JSON-RPC 2.0 error codes
const (
	ParseError     = -32700 // Invalid JSON was received by the server
	InvalidRequest = -32600 // The JSON sent is not a valid Request object
	MethodNotFound = -32601 // The method does not exist / is not available
	InvalidParams  = -32602 // Invalid method parameter(s)
	InternalError  = -32603 // Internal JSON-RPC error
)

// MCP-specific error codes
const (
	// Tool-related errors
	ToolNotFound     = -32000 // Tool not found or not available
	ToolExecuteError = -32001 // Error executing tool
	ToolTimeoutError = -32002 // Tool execution timeout

	// Resource-related errors
	ResourceNotFound     = -32010 // Resource not found
	ResourceReadError    = -32011 // Error reading resource
	ResourceAccessDenied = -32012 // Access denied to resource

	// Prompt-related errors
	PromptNotFound    = -32020 // Prompt not found
	PromptRenderError = -32021 // Error rendering prompt

	// General MCP errors
	CapabilityNotSupported  = -32030 // Requested capability not supported
	ProtocolVersionMismatch = -32031 // Protocol version mismatch
	InitializationError     = -32032 // Initialization error
)

// Custom application error codes
const (
	// Godog-specific errors
	GodogNotFound         = -40001 // Godog server not found
	FeatureParseError     = -40002 // Error parsing feature file
	TestExecutionError    = -40003 // Error executing test
	StepDefinitionError   = -40004 // Error in step definition
	ReportGenerationError = -40005 // Error generating report

	// Database-related errors
	DatabaseConnectionError  = -50001 // Database connection error
	DatabaseQueryError       = -50002 // Database query error
	DatabaseTransactionError = -50003 // Database transaction error

	// Validation errors
	ValidationError = -60001 // Validation error
	SchemaError     = -60002 // Schema validation error
)

// Standard MCP methods
const (
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodShutdown    = "shutdown"
	MethodExit        = "exit"

	// Tool methods
	MethodToolsList = "tools/list"
	MethodToolsCall = "tools/call"

	// Resource methods
	MethodResourcesList        = "resources/list"
	MethodResourcesRead        = "resources/read"
	MethodResourcesSubscribe   = "resources/subscribe"
	MethodResourcesUnsubscribe = "resources/unsubscribe"

	// Prompt methods
	MethodPromptsList = "prompts/list"
	MethodPromptsGet  = "prompts/get"

	// Notification methods
	MethodResourcesUpdated   = "resources/updated"
	MethodPromptsListChanged = "prompts/list_changed"
	MethodToolsListChanged   = "tools/list_changed"
)

// Error messages
const (
	ErrMsgParseError     = "Parse error"
	ErrMsgInvalidRequest = "Invalid Request"
	ErrMsgMethodNotFound = "Method not found"
	ErrMsgInvalidParams  = "Invalid params"
	ErrMsgInternalError  = "Internal error"

	ErrMsgToolNotFound     = "Tool not found"
	ErrMsgToolExecuteError = "Tool execution error"
	ErrMsgToolTimeoutError = "Tool execution timeout"

	ErrMsgResourceNotFound     = "Resource not found"
	ErrMsgResourceReadError    = "Resource read error"
	ErrMsgResourceAccessDenied = "Resource access denied"

	ErrMsgPromptNotFound    = "Prompt not found"
	ErrMsgPromptRenderError = "Prompt render error"

	ErrMsgCapabilityNotSupported  = "Capability not supported"
	ErrMsgProtocolVersionMismatch = "Protocol version mismatch"
	ErrMsgInitializationError     = "Initialization error"
)
