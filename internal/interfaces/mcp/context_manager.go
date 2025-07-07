package mcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

// MovieService interface for dependency injection and testing
type MovieService interface {
	SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
}

// ContextManager handles large dataset pagination and caching for MCP resources
type ContextManager struct {
	movieService MovieService
	contexts     map[string]*DataContext
	mutex        sync.RWMutex
	ttl          time.Duration
}

// DataContext represents a paginated data context
type DataContext struct {
	ID        string        `json:"id"`
	Query     interface{}   `json:"query"`
	Total     int           `json:"total"`
	PageSize  int           `json:"page_size"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
	Data      []interface{} `json:"data,omitempty"`
}

// ContextInfo provides metadata about a data context
type ContextInfo struct {
	ID         string    `json:"id"`
	Total      int       `json:"total"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// PageRequest represents a request for a page of data
type PageRequest struct {
	ContextID string `json:"context_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size,omitempty"`
}

// PageResponse represents a page of data with metadata
type PageResponse struct {
	ContextID   string        `json:"context_id"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
	Total       int           `json:"total"`
	TotalPages  int           `json:"total_pages"`
	HasNext     bool          `json:"has_next"`
	HasPrevious bool          `json:"has_previous"`
	Data        []interface{} `json:"data"`
}

// NewContextManager creates a new context manager
func NewContextManager(movieService MovieService) *ContextManager {
	return &ContextManager{
		movieService: movieService,
		contexts:     make(map[string]*DataContext),
		ttl:          time.Hour, // Default TTL of 1 hour
	}
}

// CreateMovieSearchContext creates a paginated context for movie search results
func (cm *ContextManager) CreateMovieSearchContext(ctx context.Context, query movieApp.SearchMoviesQuery, pageSize int) (*ContextInfo, error) {
	if pageSize == 0 {
		pageSize = 50 // Default page size
	}
	if pageSize > 1000 {
		pageSize = 1000 // Maximum page size
	}

	// Get total count first (would need to implement count query in real scenario)
	// For now, we'll get all results and count them
	allMovies, err := cm.movieService.SearchMovies(ctx, movieApp.SearchMoviesQuery{
		Title:    query.Title,
		Director: query.Director,
		Genre:    query.Genre,
		MinYear:  query.MinYear,
		MaxYear:  query.MaxYear,
		Limit:    10000, // Large limit to get all results for counting
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get movie count: %w", err)
	}

	total := len(allMovies)
	totalPages := (total + pageSize - 1) / pageSize

	// Generate unique context ID
	contextID := generateContextID()

	// Convert movies to interface{} slice
	data := make([]interface{}, len(allMovies))
	for i, movie := range allMovies {
		data[i] = movie
	}

	// Create context
	dataContext := &DataContext{
		ID:        contextID,
		Query:     query,
		Total:     total,
		PageSize:  pageSize,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(cm.ttl),
		Data:      data,
	}

	// Store context
	cm.mutex.Lock()
	cm.contexts[contextID] = dataContext
	cm.mutex.Unlock()

	// Clean up expired contexts
	go cm.cleanupExpiredContexts()

	return &ContextInfo{
		ID:         contextID,
		Total:      total,
		PageSize:   pageSize,
		TotalPages: totalPages,
		CreatedAt:  dataContext.CreatedAt,
		ExpiresAt:  dataContext.ExpiresAt,
	}, nil
}

// GetPage retrieves a specific page from a context
func (cm *ContextManager) GetPage(req PageRequest) (*PageResponse, error) {
	cm.mutex.RLock()
	dataContext, exists := cm.contexts[req.ContextID]
	cm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("context not found: %s", req.ContextID)
	}

	if time.Now().After(dataContext.ExpiresAt) {
		cm.mutex.Lock()
		delete(cm.contexts, req.ContextID)
		cm.mutex.Unlock()
		return nil, fmt.Errorf("context expired: %s", req.ContextID)
	}

	// Override page size if provided
	pageSize := dataContext.PageSize
	if req.PageSize > 0 && req.PageSize <= 1000 {
		pageSize = req.PageSize
	}

	totalPages := (dataContext.Total + pageSize - 1) / pageSize

	// Validate page number
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Page > totalPages {
		req.Page = totalPages
	}

	// Calculate offset and limit
	offset := (req.Page - 1) * pageSize
	limit := pageSize

	// Ensure we don't exceed data bounds
	if offset >= len(dataContext.Data) {
		return &PageResponse{
			ContextID:   req.ContextID,
			Page:        req.Page,
			PageSize:    pageSize,
			Total:       dataContext.Total,
			TotalPages:  totalPages,
			HasNext:     false,
			HasPrevious: req.Page > 1,
			Data:        []interface{}{},
		}, nil
	}

	end := offset + limit
	if end > len(dataContext.Data) {
		end = len(dataContext.Data)
	}

	pageData := dataContext.Data[offset:end]

	return &PageResponse{
		ContextID:   req.ContextID,
		Page:        req.Page,
		PageSize:    pageSize,
		Total:       dataContext.Total,
		TotalPages:  totalPages,
		HasNext:     req.Page < totalPages,
		HasPrevious: req.Page > 1,
		Data:        pageData,
	}, nil
}

// GetContextInfo retrieves information about a context
func (cm *ContextManager) GetContextInfo(contextID string) (*ContextInfo, error) {
	cm.mutex.RLock()
	dataContext, exists := cm.contexts[contextID]
	cm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("context not found: %s", contextID)
	}

	if time.Now().After(dataContext.ExpiresAt) {
		cm.mutex.Lock()
		delete(cm.contexts, contextID)
		cm.mutex.Unlock()
		return nil, fmt.Errorf("context expired: %s", contextID)
	}

	totalPages := (dataContext.Total + dataContext.PageSize - 1) / dataContext.PageSize

	return &ContextInfo{
		ID:         dataContext.ID,
		Total:      dataContext.Total,
		PageSize:   dataContext.PageSize,
		TotalPages: totalPages,
		CreatedAt:  dataContext.CreatedAt,
		ExpiresAt:  dataContext.ExpiresAt,
	}, nil
}

// DeleteContext removes a context from memory
func (cm *ContextManager) DeleteContext(contextID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.contexts[contextID]; !exists {
		return fmt.Errorf("context not found: %s", contextID)
	}

	delete(cm.contexts, contextID)
	return nil
}

// ListActiveContexts returns information about all active contexts
func (cm *ContextManager) ListActiveContexts() []ContextInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var contexts []ContextInfo
	now := time.Now()

	for _, dataContext := range cm.contexts {
		if now.Before(dataContext.ExpiresAt) {
			totalPages := (dataContext.Total + dataContext.PageSize - 1) / dataContext.PageSize
			contexts = append(contexts, ContextInfo{
				ID:         dataContext.ID,
				Total:      dataContext.Total,
				PageSize:   dataContext.PageSize,
				TotalPages: totalPages,
				CreatedAt:  dataContext.CreatedAt,
				ExpiresAt:  dataContext.ExpiresAt,
			})
		}
	}

	return contexts
}

// cleanupExpiredContexts removes expired contexts from memory
func (cm *ContextManager) cleanupExpiredContexts() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	for id, dataContext := range cm.contexts {
		if now.After(dataContext.ExpiresAt) {
			delete(cm.contexts, id)
		}
	}
}

// HandleCreateContext handles MCP tool call for creating a paginated context
func (cm *ContextManager) HandleCreateContext(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	// Parse query parameters
	queryType, ok := arguments["type"].(string)
	if !ok || queryType == "" {
		queryType = "title"
	}

	query, ok := arguments["query"].(string)
	if !ok {
		query = ""
	}
	if query == "" {
		sendError(id, dto.InvalidParams, "Query is required", nil)
		return
	}

	pageSize := 50
	if ps, ok := arguments["page_size"].(float64); ok {
		pageSize = int(ps)
	}

	// Build search query based on type
	searchQuery := movieApp.SearchMoviesQuery{}

	switch queryType {
	case "title":
		searchQuery.Title = query
	case "director":
		searchQuery.Director = query
	case "genre":
		searchQuery.Genre = query
	}

	// Add optional filters
	if director, ok := arguments["director"].(string); ok && director != "" {
		searchQuery.Director = director
	}
	if genre, ok := arguments["genre"].(string); ok && genre != "" {
		searchQuery.Genre = genre
	}
	if year, ok := arguments["year"].(float64); ok && year > 0 {
		searchQuery.MinYear = int(year)
		searchQuery.MaxYear = int(year)
	}

	// Create context
	ctx := context.Background()
	contextInfo, err := cm.CreateMovieSearchContext(ctx, searchQuery, pageSize)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to create context", err.Error())
		return
	}

	sendResult(id, contextInfo)
}

// HandleGetPage handles MCP tool call for getting a page of data
func (cm *ContextManager) HandleGetPage(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	contextID, ok := arguments["context_id"].(string)
	if !ok || contextID == "" {
		sendError(id, dto.InvalidParams, "Context ID is required", nil)
		return
	}

	page := 1
	if p, ok := arguments["page"].(float64); ok {
		page = int(p)
	}

	pageSize := 0
	if ps, ok := arguments["page_size"].(float64); ok {
		pageSize = int(ps)
	}

	req := PageRequest{
		ContextID: contextID,
		Page:      page,
		PageSize:  pageSize,
	}

	response, err := cm.GetPage(req)
	if err != nil {
		sendError(id, dto.InvalidParams, err.Error(), nil)
		return
	}

	sendResult(id, response)
}

// HandleContextInfo handles MCP tool call for getting context information
func (cm *ContextManager) HandleContextInfo(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	contextID, ok := arguments["context_id"].(string)
	if !ok || contextID == "" {
		sendError(id, dto.InvalidParams, "Context ID is required", nil)
		return
	}

	contextInfo, err := cm.GetContextInfo(contextID)
	if err != nil {
		sendError(id, dto.InvalidParams, err.Error(), nil)
		return
	}

	sendResult(id, contextInfo)
}

// generateContextID generates a unique context ID
func generateContextID() string {
	return fmt.Sprintf("ctx_%d", time.Now().UnixNano())
}
