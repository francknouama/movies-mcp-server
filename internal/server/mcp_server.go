package server

import (
	"io"
	"log"

	"movies-mcp-server/internal/composition"
	"movies-mcp-server/internal/interfaces/dto"
	"movies-mcp-server/pkg/validation"
)

// MCPServer represents the new clean MCP server implementation
type MCPServer struct {
	protocol        *Protocol
	router          *Router
	registry        *Registry
	resourceManager *ResourceManager
	container       *composition.Container
}

// NewMCPServer creates a new clean MCP server instance
func NewMCPServer(input io.Reader, output io.Writer, logger *log.Logger, container *composition.Container) *MCPServer {
	// Create the registry for auto-registration
	registry := NewRegistry()
	
	// Create the resource manager
	resourceManager := NewResourceManager(registry)
	
	// Register default resources
	resourceManager.RegisterDefaultResources()
	
	// Create validator
	validator := validation.NewRequestValidator()
	
	// Create protocol handler
	protocol := NewProtocol(input, output, logger)
	
	// Create router
	router := NewRouter(registry, validator)
	
	server := &MCPServer{
		protocol:        protocol,
		router:          router,
		registry:        registry,
		resourceManager: resourceManager,
		container:       container,
	}
	
	// Auto-register tools from container
	server.registerHandlers()
	
	// Validate all registrations
	if err := registry.ValidateRegistrations(); err != nil {
		if logger != nil {
			logger.Printf("Registration validation failed: %v", err)
		}
	}
	
	return server
}

// Run starts the server and handles incoming requests
func (s *MCPServer) Run() error {
	return s.protocol.Listen(s.router)
}

// registerHandlers registers all handlers from the container
func (s *MCPServer) registerHandlers() {
	// Get all tool schemas from the existing validator
	allSchemas := s.container.ToolValidator.GetSchemas()
	
	// Organize schemas by category
	movieSchemas := filterSchemasByPrefix(allSchemas, []string{"get_movie", "add_movie", "update_movie", "delete_movie", "list_top_movies"})
	actorSchemas := filterSchemasByPrefix(allSchemas, []string{"add_actor", "get_actor", "update_actor", "delete_actor", "link_actor_to_movie", "unlink_actor_from_movie", "get_movie_cast", "get_actor_movies", "search_actors"})
	searchSchemas := filterSchemasByPrefix(allSchemas, []string{"search_movies", "search_by_decade", "search_by_rating_range", "search_similar_movies"})
	compoundSchemas := filterSchemasByPrefix(allSchemas, []string{"bulk_movie_import", "movie_recommendation_engine", "director_career_analysis"})
	contextSchemas := filterSchemasByPrefix(allSchemas, []string{"create_search_context", "get_context_page", "get_context_info"})
	validationSchemas := filterSchemasByPrefix(allSchemas, []string{"validate_tool_call"})
	
	// Register movie tools if available
	if s.container.MovieHandlers != nil {
		s.registerMovieHandlers(movieSchemas)
	}
	
	// Register actor tools if available
	if s.container.ActorHandlers != nil {
		s.registerActorHandlers(actorSchemas)
	}
	
	// Register search tools if available
	if s.container.MovieHandlers != nil {
		s.registerSearchHandlers(searchSchemas)
	}
	
	// Register compound tools if available
	if s.container.CompoundToolHandlers != nil {
		s.registerCompoundHandlers(compoundSchemas)
	}
	
	// Register context tools if available
	if s.container.ContextManager != nil {
		s.registerContextHandlers(contextSchemas)
	}
	
	// Register validation tools if available
	if s.container.ToolValidator != nil {
		s.registerValidationHandlers(validationSchemas)
	}
}

// registerMovieHandlers registers movie-related tool handlers
func (s *MCPServer) registerMovieHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"get_movie":        s.container.MovieHandlers.HandleGetMovie,
		"add_movie":        s.container.MovieHandlers.HandleAddMovie,
		"update_movie":     s.container.MovieHandlers.HandleUpdateMovie,
		"delete_movie":     s.container.MovieHandlers.HandleDeleteMovie,
		"list_top_movies":  s.container.MovieHandlers.HandleListTopMovies,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerActorHandlers registers actor-related tool handlers
func (s *MCPServer) registerActorHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"add_actor":                s.container.ActorHandlers.HandleAddActor,
		"get_actor":                s.container.ActorHandlers.HandleGetActor,
		"update_actor":             s.container.ActorHandlers.HandleUpdateActor,
		"delete_actor":             s.container.ActorHandlers.HandleDeleteActor,
		"link_actor_to_movie":      s.container.ActorHandlers.HandleLinkActorToMovie,
		"unlink_actor_from_movie":  s.container.ActorHandlers.HandleUnlinkActorFromMovie,
		"get_movie_cast":           s.container.ActorHandlers.HandleGetMovieCast,
		"get_actor_movies":         s.container.ActorHandlers.HandleGetActorMovies,
		"search_actors":            s.container.ActorHandlers.HandleSearchActors,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerSearchHandlers registers search-related tool handlers
func (s *MCPServer) registerSearchHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"search_movies":            s.container.MovieHandlers.HandleSearchMovies,
		"search_by_decade":         s.container.MovieHandlers.HandleSearchByDecade,
		"search_by_rating_range":   s.container.MovieHandlers.HandleSearchByRatingRange,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerCompoundHandlers registers compound tool handlers
func (s *MCPServer) registerCompoundHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"bulk_movie_import":           s.container.CompoundToolHandlers.HandleBulkMovieImport,
		"movie_recommendation_engine": s.container.CompoundToolHandlers.HandleMovieRecommendationEngine,
		"director_career_analysis":    s.container.CompoundToolHandlers.HandleDirectorCareerAnalysis,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerContextHandlers registers context management tool handlers
func (s *MCPServer) registerContextHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"create_search_context": s.container.ContextManager.HandleCreateContext,
		"get_context_page":      s.container.ContextManager.HandleGetPage,
		"get_context_info":      s.container.ContextManager.HandleContextInfo,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerValidationHandlers registers validation tool handlers
func (s *MCPServer) registerValidationHandlers(schemas []dto.Tool) {
	handlers := map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)){
		"validate_tool_call": s.container.ToolValidator.HandleValidateToolCall,
	}
	
	s.registerToolsWithSchemas(handlers, schemas)
}

// registerToolsWithSchemas registers tools with their schemas
func (s *MCPServer) registerToolsWithSchemas(handlers map[string]func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)), schemas []dto.Tool) {
	// Register each handler with its corresponding schema
	for name, handler := range handlers {
		wrappedHandler := s.wrapHandler(handler)
		
		// Find the matching schema
		for _, schema := range schemas {
			if schema.Name == name {
				s.registry.RegisterTool(name, wrappedHandler, schema)
				break
			}
		}
	}
}

// filterSchemasByPrefix filters schemas by matching tool names
func filterSchemasByPrefix(allSchemas []dto.Tool, toolNames []string) []dto.Tool {
	var filtered []dto.Tool
	nameMap := make(map[string]bool)
	for _, name := range toolNames {
		nameMap[name] = true
	}
	
	for _, schema := range allSchemas {
		if nameMap[schema.Name] {
			filtered = append(filtered, schema)
		}
	}
	return filtered
}

// wrapHandler adapts the old handler signature to the new one
func (s *MCPServer) wrapHandler(oldHandler func(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any))) ToolHandlerFunc {
	return func(id any, arguments map[string]any, sender ResponseSender) {
		oldHandler(id, arguments, sender.SendResult, sender.SendError)
	}
}