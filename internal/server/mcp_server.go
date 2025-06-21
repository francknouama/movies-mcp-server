package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"movies-mcp-server/internal/composition"
	"movies-mcp-server/internal/interfaces/dto"
	"movies-mcp-server/pkg/validation"
)

// ToolHandler defines the signature for tool handler functions
type ToolHandler func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any))

// MCPServer represents a Model Context Protocol server
type MCPServer struct {
	input        io.Reader
	output       io.Writer
	logger       *log.Logger
	validator    *validation.RequestValidator
	container    *composition.Container
	toolHandlers map[string]ToolHandler
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(input io.Reader, output io.Writer, logger *log.Logger, container *composition.Container) *MCPServer {
	validator := validation.NewRequestValidator()
	
	server := &MCPServer{
		input:     input,
		output:    output,
		logger:    logger,
		validator: validator,
		container: container,
	}
	
	server.initToolHandlers()
	return server
}

// Run starts the server and handles incoming requests
func (s *MCPServer) Run() error {
	if s.logger != nil {
		s.logger.Println("Starting MCP Server...")
	}
	
	scanner := bufio.NewScanner(s.input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		if s.logger != nil {
			s.logger.Printf("Received: %s", line)
		}
		
		var request dto.JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.sendError(nil, dto.ParseError, "Parse error", err.Error())
			continue
		}
		
		s.handleRequest(&request)
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	
	return nil
}

// handleRequest processes a single JSON-RPC request
func (s *MCPServer) handleRequest(req *dto.JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "notifications/initialized":
		s.handleInitialized(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	case "resources/list":
		s.handleResourcesList(req)
	case "resources/read":
		s.handleResourcesRead(req)
	case "resources/templates/list":
		s.handleResourceTemplatesList(req)
	case "prompts/list":
		s.handlePromptsList(req)
	case "prompts/get":
		s.handlePromptsGet(req)
	case "completion/complete":
		s.handleCompletionComplete(req)
	case "logging/setLevel":
		s.handleLoggingSetLevel(req)
	default:
		s.sendError(req.ID, dto.MethodNotFound, "Method not found", nil)
	}
}

// handleInitialize handles the MCP initialize request
func (s *MCPServer) handleInitialize(req *dto.JSONRPCRequest) {
	response := dto.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.ServerCapabilities{
			Tools:     &dto.ToolsCapability{},
			Resources: &dto.ResourcesCapability{},
			Prompts:   &dto.PromptsCapability{},
		},
		ServerInfo: dto.ServerInfo{
			Name:    "movies-mcp-server",
			Version: "1.0.0",
		},
	}
	
	s.sendResult(req.ID, response)
}

// handleInitialized handles the initialized notification
func (s *MCPServer) handleInitialized(_ *dto.JSONRPCRequest) {
	// No response needed for notifications
	if s.logger != nil {
		s.logger.Println("Client initialized")
	}
}

// handleToolsList handles the tools/list request
func (s *MCPServer) handleToolsList(req *dto.JSONRPCRequest) {
	if s.container.ToolValidator != nil {
		schemas := s.container.ToolValidator.GetSchemas()
		response := dto.ToolsListResponse{
			Tools: schemas,
		}
		s.sendResult(req.ID, response)
	} else {
		response := dto.ToolsListResponse{
			Tools: []dto.Tool{},
		}
		s.sendResult(req.ID, response)
	}
}

// handleToolsCall handles the tools/call request
func (s *MCPServer) handleToolsCall(req *dto.JSONRPCRequest) {
	var params dto.ToolCallRequest
	if err := s.unmarshalParams(req.Params, &params); err != nil {
		s.sendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}
	
	// Validate tool call if validator is available
	if s.container.ToolValidator != nil {
		result := s.container.ToolValidator.ValidateToolCall(params.Name, params.Arguments)
		if !result.Valid {
			s.sendError(req.ID, dto.InvalidParams, "Tool validation failed", result.Errors)
			return
		}
	}
	
	// Execute tool
	handler, exists := s.toolHandlers[params.Name]
	if !exists {
		s.sendError(req.ID, dto.MethodNotFound, "Tool not found", params.Name)
		return
	}
	
	handler(req.ID, params.Arguments, s.sendResult, s.sendError)
}

// handleResourcesList handles the resources/list request
func (s *MCPServer) handleResourcesList(req *dto.JSONRPCRequest) {
	resources := []dto.Resource{
		{
			URI:         "movies://database/all",
			Name:        "All Movies",
			Description: "Complete movie database",
			MimeType:    "application/json",
		},
		{
			URI:         "movies://database/stats",
			Name:        "Database Statistics",
			Description: "Movie database statistics and analytics",
			MimeType:    "application/json",
		},
	}
	
	response := dto.ResourcesListResponse{
		Resources: resources,
	}
	s.sendResult(req.ID, response)
}

// handleResourcesRead handles the resources/read request
func (s *MCPServer) handleResourcesRead(req *dto.JSONRPCRequest) {
	var params dto.ResourceReadRequest
	if err := s.unmarshalParams(req.Params, &params); err != nil {
		s.sendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}
	
	// This would typically use a resource handler from the container
	// For now, return a placeholder response
	response := dto.ResourceReadResponse{
		Contents: []dto.ResourceContent{
			{
				URI:      params.URI,
				MimeType: "application/json",
				Text:     "Resource content placeholder",
			},
		},
	}
	s.sendResult(req.ID, response)
}

// handleResourceTemplatesList handles the resources/templates/list request
func (s *MCPServer) handleResourceTemplatesList(req *dto.JSONRPCRequest) {
	response := dto.ResourceTemplatesListResponse{
		ResourceTemplates: []dto.ResourceTemplate{},
	}
	s.sendResult(req.ID, response)
}

// handlePromptsList handles the prompts/list request
func (s *MCPServer) handlePromptsList(req *dto.JSONRPCRequest) {
	if s.container.PromptHandlers != nil {
		prompts := s.container.PromptHandlers.GetPrompts()
		response := dto.PromptsListResponse{
			Prompts: prompts,
		}
		s.sendResult(req.ID, response)
	} else {
		response := dto.PromptsListResponse{
			Prompts: []dto.Prompt{},
		}
		s.sendResult(req.ID, response)
	}
}

// handlePromptsGet handles the prompts/get request
func (s *MCPServer) handlePromptsGet(req *dto.JSONRPCRequest) {
	var params dto.PromptGetRequest
	if err := s.unmarshalParams(req.Params, &params); err != nil {
		s.sendError(req.ID, dto.InvalidParams, "Invalid parameters", err.Error())
		return
	}
	
	if s.container.PromptHandlers != nil {
		// Add the name to the arguments map for the handler
		handlerArgs := make(map[string]interface{})
		handlerArgs["name"] = params.Name
		if params.Arguments != nil {
			handlerArgs["arguments"] = params.Arguments
		}
		s.container.PromptHandlers.HandlePromptGet(req.ID, handlerArgs, s.sendResult, s.sendError)
	} else {
		s.sendError(req.ID, dto.MethodNotFound, "Prompt not found", params.Name)
	}
}

// handleCompletionComplete handles the completion/complete request
func (s *MCPServer) handleCompletionComplete(req *dto.JSONRPCRequest) {
	// Completion is not implemented in this server
	s.sendError(req.ID, dto.MethodNotFound, "Completion not supported", nil)
}

// handleLoggingSetLevel handles the logging/setLevel request
func (s *MCPServer) handleLoggingSetLevel(req *dto.JSONRPCRequest) {
	// Logging level setting is not implemented
	s.sendResult(req.ID, nil)
}

// initToolHandlers initializes the tool handler map
func (s *MCPServer) initToolHandlers() {
	s.toolHandlers = map[string]ToolHandler{}
	
	// Only register handlers if the container has them
	if s.container.MovieHandlers != nil {
		s.toolHandlers["get_movie"] = s.container.MovieHandlers.HandleGetMovie
		s.toolHandlers["add_movie"] = s.container.MovieHandlers.HandleAddMovie
		s.toolHandlers["update_movie"] = s.container.MovieHandlers.HandleUpdateMovie
		s.toolHandlers["delete_movie"] = s.container.MovieHandlers.HandleDeleteMovie
		s.toolHandlers["search_movies"] = s.container.MovieHandlers.HandleSearchMovies
		s.toolHandlers["list_top_movies"] = s.container.MovieHandlers.HandleListTopMovies
		s.toolHandlers["search_by_decade"] = s.container.MovieHandlers.HandleSearchByDecade
		s.toolHandlers["search_by_rating_range"] = s.container.MovieHandlers.HandleSearchByRatingRange
		// Note: search_similar_movies not implemented yet
	}
	
	if s.container.ActorHandlers != nil {
		s.toolHandlers["add_actor"] = s.container.ActorHandlers.HandleAddActor
		s.toolHandlers["get_actor"] = s.container.ActorHandlers.HandleGetActor
		s.toolHandlers["update_actor"] = s.container.ActorHandlers.HandleUpdateActor
		s.toolHandlers["delete_actor"] = s.container.ActorHandlers.HandleDeleteActor
		s.toolHandlers["link_actor_to_movie"] = s.container.ActorHandlers.HandleLinkActorToMovie
		s.toolHandlers["unlink_actor_from_movie"] = s.container.ActorHandlers.HandleUnlinkActorFromMovie
		s.toolHandlers["get_movie_cast"] = s.container.ActorHandlers.HandleGetMovieCast
		s.toolHandlers["get_actor_movies"] = s.container.ActorHandlers.HandleGetActorMovies
		s.toolHandlers["search_actors"] = s.container.ActorHandlers.HandleSearchActors
	}
	
	if s.container.CompoundToolHandlers != nil {
		s.toolHandlers["bulk_movie_import"] = s.container.CompoundToolHandlers.HandleBulkMovieImport
		s.toolHandlers["movie_recommendation_engine"] = s.container.CompoundToolHandlers.HandleMovieRecommendationEngine
		s.toolHandlers["director_career_analysis"] = s.container.CompoundToolHandlers.HandleDirectorCareerAnalysis
	}
	
	if s.container.ContextManager != nil {
		s.toolHandlers["create_search_context"] = s.container.ContextManager.HandleCreateContext
		s.toolHandlers["get_context_page"] = s.container.ContextManager.HandleGetPage
		s.toolHandlers["get_context_info"] = s.container.ContextManager.HandleContextInfo
	}
	
	if s.container.ToolValidator != nil {
		s.toolHandlers["validate_tool_call"] = s.container.ToolValidator.HandleValidateToolCall
	}
}

// sendResult sends a successful JSON-RPC response
func (s *MCPServer) sendResult(id any, result any) {
	response := dto.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	s.sendResponse(response)
}

// sendError sends an error JSON-RPC response
func (s *MCPServer) sendError(id any, code int, message string, data interface{}) {
	response := dto.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &dto.JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	s.sendResponse(response)
}

// sendResponse sends a JSON-RPC response
func (s *MCPServer) sendResponse(response dto.JSONRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		if s.logger != nil {
			s.logger.Printf("Failed to marshal response: %v", err)
		}
		return
	}
	
	if s.logger != nil {
		s.logger.Printf("Sending: %s", string(data))
	}
	
	s.output.Write(data)
	s.output.Write([]byte("\n"))
}

// unmarshalParams unmarshals request parameters
func (s *MCPServer) unmarshalParams(params interface{}, target interface{}) error {
	if params == nil {
		return fmt.Errorf("missing parameters")
	}
	
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, target)
}