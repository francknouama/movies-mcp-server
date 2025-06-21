package composition

import (
	"database/sql"

	actorApp "movies-mcp-server/internal/application/actor"
	movieApp "movies-mcp-server/internal/application/movie"
	"movies-mcp-server/internal/infrastructure/postgres"
	"movies-mcp-server/internal/interfaces/mcp"
	"movies-mcp-server/internal/schemas"
)

// Container holds all application dependencies
type Container struct {
	// Repositories
	MovieRepository *postgres.MovieRepository
	ActorRepository *postgres.ActorRepository
	
	// Application Services
	MovieService *movieApp.Service
	ActorService *actorApp.Service
	
	// Interface Handlers
	MovieHandlers         *mcp.MovieHandlers
	ActorHandlers         *mcp.ActorHandlers
	PromptHandlers        *mcp.PromptHandlers
	CompoundToolHandlers  *mcp.CompoundToolHandlers
	ContextManager        *mcp.ContextManager
	ToolValidator         *mcp.ToolValidator
}

// NewContainer creates and wires up all application dependencies
func NewContainer(db *sql.DB) *Container {
	// Infrastructure Layer - Repositories
	movieRepo := postgres.NewMovieRepository(db)
	actorRepo := postgres.NewActorRepository(db)
	
	// Application Layer - Services
	movieService := movieApp.NewService(movieRepo)
	actorService := actorApp.NewService(actorRepo)
	
	// Interface Layer - Handlers
	movieHandlers := mcp.NewMovieHandlers(movieService)
	actorHandlers := mcp.NewActorHandlers(actorService)
	promptHandlers := mcp.NewPromptHandlers()
	compoundToolHandlers := mcp.NewCompoundToolHandlers(movieService)
	contextManager := mcp.NewContextManager(movieService)
	toolValidator := mcp.NewToolValidator(schemas.GetToolSchemas())
	
	return &Container{
		MovieRepository:       movieRepo,
		ActorRepository:       actorRepo,
		MovieService:          movieService,
		ActorService:          actorService,
		MovieHandlers:         movieHandlers,
		ActorHandlers:         actorHandlers,
		PromptHandlers:        promptHandlers,
		CompoundToolHandlers:  compoundToolHandlers,
		ContextManager:        contextManager,
		ToolValidator:         toolValidator,
	}
}

// NewTestContainer creates minimal dependencies for testing MCP protocol
func NewTestContainer() *Container {
	// For MCP protocol testing, we mainly need the core protocol handlers
	// We'll create minimal handlers that can respond to protocol requests
	promptHandlers := mcp.NewPromptHandlers()
	toolValidator := mcp.NewToolValidator(schemas.GetToolSchemas())
	
	return &Container{
		MovieRepository:      nil, // Not needed for protocol testing
		ActorRepository:      nil, // Not needed for protocol testing
		MovieService:         nil, // Not needed for basic protocol testing
		ActorService:         nil, // Not needed for basic protocol testing
		MovieHandlers:        nil, // Will be handled by protocol-level testing
		ActorHandlers:        nil, // Will be handled by protocol-level testing
		PromptHandlers:       promptHandlers,
		CompoundToolHandlers: nil, // Will be handled by protocol-level testing
		ContextManager:       nil, // Will be handled by protocol-level testing
		ToolValidator:        toolValidator,
	}
}