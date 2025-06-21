package server

import (
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
	"github.com/francknouama/movies-mcp-server/shared-mcp/pkg/validation"
)

// Router handles MCP request routing and method dispatch
type Router struct {
	registry  *Registry
	validator *validation.RequestValidator
}

// NewRouter creates a new router with the given registry and validator
func NewRouter(registry *Registry, validator *validation.RequestValidator) *Router {
	return &Router{
		registry:  registry,
		validator: validator,
	}
}

// HandleRequest implements RequestHandler interface
func (r *Router) HandleRequest(req *dto.JSONRPCRequest, protocol *Protocol) {
	switch req.Method {
	case "initialize":
		r.handleInitialize(req, protocol)
	case "notifications/initialized":
		r.handleInitialized(req, protocol)
	case "tools/list":
		r.handleToolsList(req, protocol)
	case "tools/call":
		r.handleToolsCall(req, protocol)
	case "resources/list":
		r.handleResourcesList(req, protocol)
	case "resources/read":
		r.handleResourcesRead(req, protocol)
	case "resources/templates/list":
		r.handleResourceTemplatesList(req, protocol)
	case "prompts/list":
		r.handlePromptsList(req, protocol)
	case "prompts/get":
		r.handlePromptsGet(req, protocol)
	case "completion/complete":
		r.handleCompletionComplete(req, protocol)
	case "logging/setLevel":
		r.handleLoggingSetLevel(req, protocol)
	default:
		protocol.SendError(req.ID, dto.MethodNotFound, "Method not found", req.Method)
	}
}

// handleInitialize handles the MCP initialize request
func (r *Router) handleInitialize(req *dto.JSONRPCRequest, protocol *Protocol) {
	response := dto.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.ServerCapabilities{
			Tools:     &dto.ToolsCapability{},
			Resources: &dto.ResourcesCapability{},
			Prompts:   &dto.PromptsCapability{},
		},
		ServerInfo: dto.ServerInfo{
			Name:    "movies-mcp-server",
			Version: "2.0.0", // Updated for clean architecture
		},
	}

	protocol.SendResult(req.ID, response)
}

// handleInitialized handles the initialized notification
func (r *Router) handleInitialized(_ *dto.JSONRPCRequest, _ *Protocol) {
	// No response needed for notifications
}

// handleToolsList handles the tools/list request
func (r *Router) handleToolsList(req *dto.JSONRPCRequest, protocol *Protocol) {
	schemas := r.registry.GetToolSchemas()
	response := dto.ToolsListResponse{
		Tools: schemas,
	}
	protocol.SendResult(req.ID, response)
}

// handleToolsCall handles the tools/call request
func (r *Router) handleToolsCall(req *dto.JSONRPCRequest, protocol *Protocol) {
	var params dto.ToolCallRequest
	if err := protocol.UnmarshalParams(req.Params, &params); err != nil {
		protocol.SendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}

	// Validate tool call if validator is available
	if r.validator != nil {
		// Note: For now, we skip validation here since it's handled
		// by the container's ToolValidator. This can be enhanced later.
	}

	// Get and execute tool handler
	handler, exists := r.registry.GetToolHandler(params.Name)
	if !exists {
		protocol.SendError(req.ID, dto.MethodNotFound, "Tool not found", params.Name)
		return
	}

	handler(req.ID, params.Arguments, protocol)
}

// handleResourcesList handles the resources/list request
func (r *Router) handleResourcesList(req *dto.JSONRPCRequest, protocol *Protocol) {
	resources := r.registry.GetResources()
	response := dto.ResourcesListResponse{
		Resources: resources,
	}
	protocol.SendResult(req.ID, response)
}

// handleResourcesRead handles the resources/read request
func (r *Router) handleResourcesRead(req *dto.JSONRPCRequest, protocol *Protocol) {
	var params dto.ResourceReadRequest
	if err := protocol.UnmarshalParams(req.Params, &params); err != nil {
		protocol.SendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}

	// Get and execute resource handler
	handler, exists := r.registry.GetResourceHandler(params.URI)
	if !exists {
		protocol.SendError(req.ID, dto.MethodNotFound, "Resource not found", params.URI)
		return
	}

	handler(params.URI, protocol)
}

// handleResourceTemplatesList handles the resources/templates/list request
func (r *Router) handleResourceTemplatesList(req *dto.JSONRPCRequest, protocol *Protocol) {
	response := dto.ResourceTemplatesListResponse{
		ResourceTemplates: []dto.ResourceTemplate{},
	}
	protocol.SendResult(req.ID, response)
}

// handlePromptsList handles the prompts/list request
func (r *Router) handlePromptsList(req *dto.JSONRPCRequest, protocol *Protocol) {
	prompts := r.registry.GetPrompts()
	response := dto.PromptsListResponse{
		Prompts: prompts,
	}
	protocol.SendResult(req.ID, response)
}

// handlePromptsGet handles the prompts/get request
func (r *Router) handlePromptsGet(req *dto.JSONRPCRequest, protocol *Protocol) {
	var params dto.PromptGetRequest
	if err := protocol.UnmarshalParams(req.Params, &params); err != nil {
		protocol.SendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}

	// Get and execute prompt handler
	handler, exists := r.registry.GetPromptHandler(params.Name)
	if !exists {
		protocol.SendError(req.ID, dto.MethodNotFound, "Prompt not found", params.Name)
		return
	}

	handler(req.ID, params.Name, params.Arguments, protocol)
}

// handleCompletionComplete handles the completion/complete request
func (r *Router) handleCompletionComplete(req *dto.JSONRPCRequest, protocol *Protocol) {
	// Completion is not implemented in this server
	protocol.SendError(req.ID, dto.MethodNotFound, "Completion not supported", nil)
}

// handleLoggingSetLevel handles the logging/setLevel request
func (r *Router) handleLoggingSetLevel(req *dto.JSONRPCRequest, protocol *Protocol) {
	// Logging level setting is not implemented
	protocol.SendResult(req.ID, nil)
}
