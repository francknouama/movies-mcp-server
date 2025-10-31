package tools

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// ContextTools provides SDK-based MCP handlers for context management
type ContextTools struct {
	movieService MovieService
	contexts     map[string]*DataContext
	mutex        sync.RWMutex
	ttl          time.Duration
}

// DataContext represents a paginated data context
type DataContext struct {
	ID        string
	Query     movieApp.SearchMoviesQuery
	Total     int
	PageSize  int
	CreatedAt time.Time
	ExpiresAt time.Time
	Data      []*movieApp.MovieDTO
}

// NewContextTools creates a new context tools instance
func NewContextTools(movieService MovieService) *ContextTools {
	return &ContextTools{
		movieService: movieService,
		contexts:     make(map[string]*DataContext),
		ttl:          time.Hour,
	}
}

// ===== create_search_context Tool =====

// CreateSearchContextInput defines the input schema for create_search_context tool
type CreateSearchContextInput struct {
	Query    SearchMoviesInput `json:"query" jsonschema:"required,description=Search query for movies"`
	PageSize int               `json:"page_size,omitempty" jsonschema:"description=Number of items per page,default=50"`
}

// CreateSearchContextOutput defines the output schema for create_search_context tool
type CreateSearchContextOutput struct {
	ContextID  string `json:"context_id" jsonschema:"description=Unique context identifier"`
	Total      int    `json:"total" jsonschema:"description=Total number of results"`
	PageSize   int    `json:"page_size" jsonschema:"description=Items per page"`
	TotalPages int    `json:"total_pages" jsonschema:"description=Total number of pages"`
	CreatedAt  string `json:"created_at" jsonschema:"description=Context creation time"`
	ExpiresAt  string `json:"expires_at" jsonschema:"description=Context expiration time"`
}

// CreateSearchContext handles the create_search_context tool call
func (t *ContextTools) CreateSearchContext(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CreateSearchContextInput,
) (*mcp.CallToolResult, CreateSearchContextOutput, error) {
	// Set defaults
	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	// Build search query
	query := movieApp.SearchMoviesQuery{
		Title:     input.Query.Title,
		Director:  input.Query.Director,
		Genre:     input.Query.Genre,
		MinYear:   input.Query.MinYear,
		MaxYear:   input.Query.MaxYear,
		MinRating: input.Query.MinRating,
		MaxRating: input.Query.MaxRating,
		Limit:     10000, // Large limit to get all results
	}

	// Get all movies
	movies, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, CreateSearchContextOutput{}, fmt.Errorf("failed to search movies: %w", err)
	}

	total := len(movies)
	totalPages := (total + pageSize - 1) / pageSize

	// Generate unique context ID
	contextID := uuid.New().String()

	// Create context
	dataContext := &DataContext{
		ID:        contextID,
		Query:     query,
		Total:     total,
		PageSize:  pageSize,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(t.ttl),
		Data:      movies,
	}

	// Store context
	t.mutex.Lock()
	t.contexts[contextID] = dataContext
	t.mutex.Unlock()

	// Clean up expired contexts in background
	go t.cleanupExpiredContexts()

	output := CreateSearchContextOutput{
		ContextID:  contextID,
		Total:      total,
		PageSize:   pageSize,
		TotalPages: totalPages,
		CreatedAt:  dataContext.CreatedAt.Format(time.RFC3339),
		ExpiresAt:  dataContext.ExpiresAt.Format(time.RFC3339),
	}

	return nil, output, nil
}

// ===== get_context_page Tool =====

// GetContextPageInput defines the input schema for get_context_page tool
type GetContextPageInput struct {
	ContextID string `json:"context_id" jsonschema:"required,description=Context identifier"`
	Page      int    `json:"page" jsonschema:"required,description=Page number (1-based)"`
	PageSize  int    `json:"page_size,omitempty" jsonschema:"description=Override page size"`
}

// GetContextPageOutput defines the output schema for get_context_page tool
type GetContextPageOutput struct {
	ContextID   string           `json:"context_id" jsonschema:"description=Context identifier"`
	Page        int              `json:"page" jsonschema:"description=Current page number"`
	PageSize    int              `json:"page_size" jsonschema:"description=Items per page"`
	Total       int              `json:"total" jsonschema:"description=Total items"`
	TotalPages  int              `json:"total_pages" jsonschema:"description=Total pages"`
	HasNext     bool             `json:"has_next" jsonschema:"description=Has next page"`
	HasPrevious bool             `json:"has_previous" jsonschema:"description=Has previous page"`
	Data        []GetMovieOutput `json:"data" jsonschema:"description=Page data"`
}

// GetContextPage handles the get_context_page tool call
func (t *ContextTools) GetContextPage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetContextPageInput,
) (*mcp.CallToolResult, GetContextPageOutput, error) {
	t.mutex.RLock()
	dataContext, exists := t.contexts[input.ContextID]
	t.mutex.RUnlock()

	if !exists {
		return nil, GetContextPageOutput{}, fmt.Errorf("context not found: %s", input.ContextID)
	}

	if time.Now().After(dataContext.ExpiresAt) {
		t.mutex.Lock()
		delete(t.contexts, input.ContextID)
		t.mutex.Unlock()
		return nil, GetContextPageOutput{}, fmt.Errorf("context expired: %s", input.ContextID)
	}

	// Override page size if provided
	pageSize := dataContext.PageSize
	if input.PageSize > 0 && input.PageSize <= 1000 {
		pageSize = input.PageSize
	}

	totalPages := (dataContext.Total + pageSize - 1) / pageSize

	// Validate page number
	page := input.Page
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	// Calculate offset and limit
	offset := (page - 1) * pageSize
	end := offset + pageSize
	if end > len(dataContext.Data) {
		end = len(dataContext.Data)
	}

	// Get page data
	var pageData []GetMovieOutput
	if offset < len(dataContext.Data) {
		for _, movie := range dataContext.Data[offset:end] {
			pageData = append(pageData, GetMovieOutput{
				ID:        movie.ID,
				Title:     movie.Title,
				Director:  movie.Director,
				Year:      movie.Year,
				Rating:    movie.Rating,
				Genres:    movie.Genres,
				PosterURL: movie.PosterURL,
				CreatedAt: movie.CreatedAt,
				UpdatedAt: movie.UpdatedAt,
			})
		}
	}

	output := GetContextPageOutput{
		ContextID:   input.ContextID,
		Page:        page,
		PageSize:    pageSize,
		Total:       dataContext.Total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
		Data:        pageData,
	}

	return nil, output, nil
}

// ===== get_context_info Tool =====

// GetContextInfoInput defines the input schema for get_context_info tool
type GetContextInfoInput struct {
	ContextID string `json:"context_id" jsonschema:"required,description=Context identifier"`
}

// GetContextInfoOutput defines the output schema for get_context_info tool
type GetContextInfoOutput struct {
	ContextID  string `json:"context_id" jsonschema:"description=Context identifier"`
	Total      int    `json:"total" jsonschema:"description=Total items"`
	PageSize   int    `json:"page_size" jsonschema:"description=Items per page"`
	TotalPages int    `json:"total_pages" jsonschema:"description=Total pages"`
	CreatedAt  string `json:"created_at" jsonschema:"description=Creation time"`
	ExpiresAt  string `json:"expires_at" jsonschema:"description=Expiration time"`
}

// GetContextInfo handles the get_context_info tool call
func (t *ContextTools) GetContextInfo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetContextInfoInput,
) (*mcp.CallToolResult, GetContextInfoOutput, error) {
	t.mutex.RLock()
	dataContext, exists := t.contexts[input.ContextID]
	t.mutex.RUnlock()

	if !exists {
		return nil, GetContextInfoOutput{}, fmt.Errorf("context not found: %s", input.ContextID)
	}

	if time.Now().After(dataContext.ExpiresAt) {
		t.mutex.Lock()
		delete(t.contexts, input.ContextID)
		t.mutex.Unlock()
		return nil, GetContextInfoOutput{}, fmt.Errorf("context expired: %s", input.ContextID)
	}

	totalPages := (dataContext.Total + dataContext.PageSize - 1) / dataContext.PageSize

	output := GetContextInfoOutput{
		ContextID:  dataContext.ID,
		Total:      dataContext.Total,
		PageSize:   dataContext.PageSize,
		TotalPages: totalPages,
		CreatedAt:  dataContext.CreatedAt.Format(time.RFC3339),
		ExpiresAt:  dataContext.ExpiresAt.Format(time.RFC3339),
	}

	return nil, output, nil
}

// Helper method to clean up expired contexts
func (t *ContextTools) cleanupExpiredContexts() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()
	for id, ctx := range t.contexts {
		if now.After(ctx.ExpiresAt) {
			delete(t.contexts, id)
		}
	}
}
