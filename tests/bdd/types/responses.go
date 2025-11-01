package types

import "time"

// ActorResponse represents an actor in API responses
type ActorResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	BirthYear int       `json:"birth_year"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// MovieResponse represents a movie in API responses
type MovieResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Director  string    `json:"director"`
	Year      int       `json:"year"`
	Rating    float64   `json:"rating"`
	Genres    []string  `json:"genres,omitempty"`
	PosterURL string    `json:"poster_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ActorsListResponse represents a list of actors
type ActorsListResponse struct {
	Actors []*ActorResponse `json:"actors"`
	Total  int              `json:"total,omitempty"`
}

// MoviesListResponse represents a list of movies
type MoviesListResponse struct {
	Movies []*MovieResponse `json:"movies"`
	Total  int              `json:"total,omitempty"`
}

// MovieCastResponse represents the cast of a movie
type MovieCastResponse struct {
	MovieID int               `json:"movie_id"`
	Cast    []*ActorResponse `json:"cast"`
}

// ActorMoviesResponse represents movies associated with an actor
type ActorMoviesResponse struct {
	ActorID  int   `json:"actor_id"`
	MovieIDs []int `json:"movie_ids"`
}

// SearchContextResponse represents a search context for pagination
type SearchContextResponse struct {
	ContextID  string `json:"context_id"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
	PageSize   int    `json:"page_size"`
}

// ContextPageResponse represents a page of results from a search context
type ContextPageResponse struct {
	ContextID  string           `json:"context_id"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
	Items      []*MovieResponse `json:"items"`
}

// BulkImportResponse represents the result of a bulk import operation
type BulkImportResponse struct {
	Imported    int      `json:"imported"`
	Failed      int      `json:"failed"`
	FailedItems []string `json:"failed_items,omitempty"`
}

// RecommendationResponse represents a movie recommendation
type RecommendationResponse struct {
	Movies []*RecommendedMovie `json:"movies"`
}

// RecommendedMovie represents a single recommended movie with score
type RecommendedMovie struct {
	Movie *MovieResponse `json:"movie"`
	Score float64        `json:"score"`
}

// DirectorAnalysisResponse represents director career analysis
type DirectorAnalysisResponse struct {
	Director    string           `json:"director"`
	TotalMovies int              `json:"total_movies"`
	AvgRating   float64          `json:"avg_rating"`
	EarlyPhase  []*MovieResponse `json:"early_phase,omitempty"`
	MidPhase    []*MovieResponse `json:"mid_phase,omitempty"`
	LatePhase   []*MovieResponse `json:"late_phase,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// DatabaseStatsResponse represents database statistics
type DatabaseStatsResponse struct {
	MovieCount int `json:"movie_count"`
	ActorCount int `json:"actor_count"`
	TotalSize  int `json:"total_size,omitempty"`
}

// MCPCapabilitiesResponse represents MCP server capabilities
type MCPCapabilitiesResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ServerInfo      MCPServerInfo          `json:"serverInfo"`
	Capabilities    map[string]interface{} `json:"capabilities"`
}

// MCPServerInfo represents server information
type MCPServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPToolsListResponse represents the list of available tools
type MCPToolsListResponse struct {
	Tools []MCPTool `json:"tools"`
}

// MCPTool represents a single MCP tool
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPResourcesListResponse represents the list of available resources
type MCPResourcesListResponse struct {
	Resources []MCPResource `json:"resources"`
}

// MCPResource represents a single MCP resource
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType,omitempty"`
}
