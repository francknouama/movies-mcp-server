package server

import (
	"encoding/json"
	"io"
	"strings"

	"shared-mcp/pkg/errors"
	"shared-mcp/pkg/logging"

	"github.com/francknouama/movies-mcp-server/godog-server/internal/godog"
)

// MCPServer handles the Model Context Protocol communication
type MCPServer struct {
	input       io.Reader
	output      io.Writer
	logger      *logging.Logger
	godogRunner *godog.Runner
	tools       map[string]ToolHandler
	resources   map[string]ResourceHandler
}

// ToolHandler represents a function that handles tool calls
type ToolHandler func(arguments map[string]any) (any, error)

// ResourceHandler represents a function that handles resource requests
type ResourceHandler func(uri string) (any, error)

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response
type MCPResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// MCPError represents an MCP error response
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(input io.Reader, output io.Writer, logger *logging.Logger, godogRunner *godog.Runner) *MCPServer {
	server := &MCPServer{
		input:       input,
		output:      output,
		logger:      logger,
		godogRunner: godogRunner,
		tools:       make(map[string]ToolHandler),
		resources:   make(map[string]ResourceHandler),
	}

	// Register basic MCP tools
	server.registerTools()
	server.registerResources()

	return server
}

// Run starts the MCP server and handles requests
func (s *MCPServer) Run() error {
	decoder := json.NewDecoder(s.input)
	encoder := json.NewEncoder(s.output)

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				s.logger.Info("MCP server shutting down")
				return nil
			}
			s.logger.WithField("error", err).Error("Failed to decode request")
			continue
		}

		s.logger.WithField("method", request.Method).Debug("Received MCP request")

		response := s.handleRequest(request)

		if err := encoder.Encode(response); err != nil {
			s.logger.WithField("error", err).Error("Failed to encode response")
			continue
		}
	}
}

// handleRequest processes an MCP request and returns a response
func (s *MCPServer) handleRequest(request MCPRequest) MCPResponse {
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "tools/list":
		return s.handleListTools(request)
	case "tools/call":
		return s.handleToolCall(request)
	case "resources/list":
		return s.handleListResources(request)
	case "resources/read":
		return s.handleResourceRead(request)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.MethodNotFound,
				Message: "Method not found: " + request.Method,
			},
		}
	}
}

// handleInitialize handles the MCP initialize request
func (s *MCPServer) handleInitialize(request MCPRequest) MCPResponse {
	s.logger.Info("Initializing MCP server")

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "godog-mcp-server",
				"version": "0.1.0",
			},
		},
	}
}

// handleListTools handles the tools/list request
func (s *MCPServer) handleListTools(request MCPRequest) MCPResponse {
	tools := make([]map[string]any, 0, len(s.tools))

	// Add basic tools
	tools = append(tools, map[string]any{
		"name":        "validate_feature",
		"description": "Parse and validate a Gherkin feature file",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"file_path": map[string]any{
					"type":        "string",
					"description": "Path to the .feature file",
				},
			},
			"required": []string{"file_path"},
		},
	})

	tools = append(tools, map[string]any{
		"name":        "list_features",
		"description": "List all available feature files",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"directory": map[string]any{
					"type":        "string",
					"description": "Directory to search for feature files (optional)",
				},
				"include_content": map[string]any{
					"type":        "boolean",
					"description": "Include feature content in response (default: false)",
				},
			},
		},
	})

	tools = append(tools, map[string]any{
		"name":        "get_feature_content",
		"description": "Get the content of a specific feature file",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"file_path": map[string]any{
					"type":        "string",
					"description": "Path to the feature file",
				},
				"include_parsed": map[string]any{
					"type":        "boolean",
					"description": "Include parsed Gherkin structure (default: true)",
				},
			},
			"required": []string{"file_path"},
		},
	})

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]any{
			"tools": tools,
		},
	}
}

// handleToolCall handles the tools/call request
func (s *MCPServer) handleToolCall(request MCPRequest) MCPResponse {
	params, ok := request.Params.(map[string]any)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InvalidParams,
				Message: "Tool name is required",
			},
		}
	}

	arguments, ok := params["arguments"].(map[string]any)
	if !ok {
		arguments = make(map[string]any)
	}

	handler, exists := s.tools[toolName]
	if !exists {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.MethodNotFound,
				Message: "Tool not found: " + toolName,
			},
		}
	}

	result, err := handler(arguments)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InternalError,
				Message: err.Error(),
			},
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

// handleListResources handles the resources/list request
func (s *MCPServer) handleListResources(request MCPRequest) MCPResponse {
	resources := []map[string]any{
		{
			"uri":         "godog://features/all",
			"name":        "All Features",
			"description": "List of all available feature files",
			"mimeType":    "application/json",
		},
		{
			"uri":         "godog://features/list",
			"name":        "Feature File List",
			"description": "Lightweight list of feature files without content",
			"mimeType":    "application/json",
		},
		{
			"uri":         "godog://features/{name}",
			"name":        "Individual Feature",
			"description": "Content of a specific feature file (use actual filename)",
			"mimeType":    "text/plain",
		},
		{
			"uri":         "godog://reports/latest",
			"name":        "Latest Test Results",
			"description": "Most recent test execution results",
			"mimeType":    "application/json",
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]any{
			"resources": resources,
		},
	}
}

// handleResourceRead handles the resources/read request
func (s *MCPServer) handleResourceRead(request MCPRequest) MCPResponse {
	params, ok := request.Params.(map[string]any)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InvalidParams,
				Message: "URI is required",
			},
		}
	}

	// Handle individual feature file URIs (pattern: godog://features/{filename})
	if strings.HasPrefix(uri, "godog://features/") && !strings.HasSuffix(uri, "/all") && !strings.HasSuffix(uri, "/list") {
		filename := strings.TrimPrefix(uri, "godog://features/")
		content, err := s.handleIndividualFeatureResource(filename)
		if err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: MCPError{
					Code:    errors.InternalError,
					Message: err.Error(),
				},
			}
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]any{
				"contents": []map[string]any{
					{
						"uri":      uri,
						"mimeType": "text/plain",
						"text":     content,
					},
				},
			},
		}
	}

	handler, exists := s.resources[uri]
	if !exists {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.MethodNotFound,
				Message: "Resource not found: " + uri,
			},
		}
	}

	content, err := handler(uri)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: MCPError{
				Code:    errors.InternalError,
				Message: err.Error(),
			},
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]any{
			"contents": []map[string]any{
				{
					"uri":      uri,
					"mimeType": "application/json",
					"text":     content,
				},
			},
		},
	}
}

// registerTools registers available tool handlers
func (s *MCPServer) registerTools() {
	s.tools["validate_feature"] = s.handleValidateFeature
	s.tools["list_features"] = s.handleListFeatures
	s.tools["get_feature_content"] = s.handleGetFeatureContent
}

// registerResources registers available resource handlers
func (s *MCPServer) registerResources() {
	s.resources["godog://features/all"] = s.handleFeaturesResource
	s.resources["godog://features/list"] = s.handleFeatureListResource
	s.resources["godog://reports/latest"] = s.handleLatestReportResource
}

// Tool handlers
func (s *MCPServer) handleValidateFeature(arguments map[string]any) (any, error) {
	filePath, ok := arguments["file_path"].(string)
	if !ok {
		return nil, errors.NewInvalidParams("file_path is required")
	}

	return s.godogRunner.ValidateFeature(filePath)
}

func (s *MCPServer) handleListFeatures(arguments map[string]any) (any, error) {
	directory, _ := arguments["directory"].(string)
	includeContent, _ := arguments["include_content"].(bool)
	return s.godogRunner.ListFeatures(directory, includeContent)
}

func (s *MCPServer) handleGetFeatureContent(arguments map[string]any) (any, error) {
	filePath, ok := arguments["file_path"].(string)
	if !ok {
		return nil, errors.NewInvalidParams("file_path is required")
	}

	includeParsed, _ := arguments["include_parsed"].(bool)
	if !includeParsed {
		includeParsed = true // default to true
	}

	return s.godogRunner.GetFeatureContent(filePath, includeParsed)
}

// Resource handlers
func (s *MCPServer) handleFeaturesResource(uri string) (any, error) {
	return s.godogRunner.ListFeatures("", true) // include content for resource access
}

func (s *MCPServer) handleFeatureListResource(uri string) (any, error) {
	return s.godogRunner.ListFeatures("", false) // lightweight list without content
}

func (s *MCPServer) handleIndividualFeatureResource(filename string) (any, error) {
	// Find the feature file by name
	featureListResult, err := s.godogRunner.ListFeatures("", false)
	if err != nil {
		return nil, err
	}

	featureListMap, ok := featureListResult.(map[string]any)
	if !ok {
		return nil, errors.NewInternalError("Invalid feature list format")
	}

	featureFiles, ok := featureListMap["feature_files"].([]map[string]any)
	if !ok {
		return nil, errors.NewInternalError("Invalid feature files format")
	}

	// Find the feature file by name
	var targetPath string
	for _, featureFile := range featureFiles {
		if name, ok := featureFile["name"].(string); ok && name == filename {
			if path, ok := featureFile["file_path"].(string); ok {
				targetPath = path
				break
			}
		}
	}

	if targetPath == "" {
		return nil, errors.NewFeatureParseError("Feature file not found: " + filename)
	}

	// Get the content
	contentResult, err := s.godogRunner.GetFeatureContent(targetPath, false)
	if err != nil {
		return nil, err
	}

	if contentMap, ok := contentResult.(map[string]any); ok {
		if rawContent, ok := contentMap["raw_content"].(string); ok {
			return rawContent, nil
		}
	}

	return nil, errors.NewInternalError("Failed to retrieve feature content")
}

func (s *MCPServer) handleLatestReportResource(uri string) (any, error) {
	return s.godogRunner.GetLatestReport()
}
