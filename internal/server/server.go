package server

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"movies-mcp-server/internal/config"
	"movies-mcp-server/internal/database"
	"movies-mcp-server/internal/models"
	"movies-mcp-server/pkg/errors"
	"movies-mcp-server/pkg/image"
	"movies-mcp-server/pkg/validation"
)

// MoviesServer represents the MCP server
type MoviesServer struct {
	input          io.Reader
	output         io.Writer
	logger         *log.Logger
	db             database.Database
	imageProcessor *image.ImageProcessor
	validator      *validation.RequestValidator
}

// New creates a new MoviesServer instance
func New(db database.Database) *MoviesServer {
	// Create default image configuration
	imageCfg := &config.ImageConfig{
		MaxSize:          5 * 1024 * 1024, // 5MB
		AllowedTypes:     []string{"image/jpeg", "image/png", "image/webp"},
		EnableThumbnails: false, // Disabled for now
		ThumbnailSize:    "200x200",
	}
	
	return &MoviesServer{
		input:          os.Stdin,
		output:         os.Stdout,
		logger:         log.New(os.Stderr, "[movies-mcp] ", log.LstdFlags),
		db:             db,
		imageProcessor: image.NewImageProcessor(imageCfg),
		validator:      validation.NewRequestValidator(),
	}
}

// NewWithConfig creates a new MoviesServer instance with custom config
func NewWithConfig(db database.Database, cfg *config.Config) *MoviesServer {
	return &MoviesServer{
		input:          os.Stdin,
		output:         os.Stdout,
		logger:         log.New(os.Stderr, "[movies-mcp] ", log.LstdFlags),
		db:             db,
		imageProcessor: image.NewImageProcessor(&cfg.Image),
		validator:      validation.NewRequestValidator(),
	}
}

// Run starts the server and handles incoming requests
func (s *MoviesServer) Run() error {
	if s.logger != nil {
		s.logger.Println("Starting Movies MCP Server...")
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
		
		var request models.JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.sendError(nil, models.ParseError, "Parse error", err.Error())
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
func (s *MoviesServer) handleRequest(req *models.JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "notifications/initialized":
		s.handleNotificationInitialized(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolCall(req)
	case "resources/list":
		s.handleResourcesList(req)
	case "resources/read":
		s.handleResourceRead(req)
	default:
		// Only send "Method not found" for actual requests (with ID), not notifications
		if req.ID != nil {
			s.sendError(req.ID, models.MethodNotFound, "Method not found", nil)
		}
		// For notifications without ID, don't respond
	}
}

// handleInitialize handles the initialize request
func (s *MoviesServer) handleInitialize(req *models.JSONRPCRequest) {
	var params models.InitializeRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, models.InvalidParams, "Invalid params", err.Error())
		return
	}

	// Validate initialize request parameters
	paramsMap := make(map[string]interface{})
	if err := json.Unmarshal(req.Params, &paramsMap); err == nil {
		if err := s.validator.ValidateInitializeRequest(paramsMap); err != nil {
			appErr := err.(*errors.ApplicationError)
			s.sendError(req.ID, models.InvalidParams, appErr.Message, appErr.Details)
			return
		}
	}
	
	response := models.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: models.ServerCapabilities{
			Tools: &models.ToolsCapability{},
			Resources: &models.ResourcesCapability{
				Subscribe: false,
			},
		},
		ServerInfo: models.ServerInfo{
			Name:    "movies-mcp-server",
			Version: "0.1.0",
		},
	}
	
	s.sendResult(req.ID, response)
}

// handleNotificationInitialized handles the notifications/initialized notification
func (s *MoviesServer) handleNotificationInitialized(_ *models.JSONRPCRequest) {
	// This is a notification, so we don't send any response
	// The client is just informing us that initialization is complete
	if s.logger != nil {
		s.logger.Println("Client initialization completed")
	}
}

// handleToolsList handles the tools/list request
func (s *MoviesServer) handleToolsList(req *models.JSONRPCRequest) {
	tools := []models.Tool{
		{
			Name:        "get_movie",
			Description: "Get a movie by ID",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The movie ID",
					},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "add_movie",
			Description: "Add a new movie to the database",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Movie title",
					},
					"director": map[string]interface{}{
						"type":        "string",
						"description": "Movie director",
					},
					"year": map[string]interface{}{
						"type":        "integer",
						"description": "Release year",
						"minimum":     1888,
						"maximum":     2100,
					},
					"genre": map[string]interface{}{
						"type":        "array",
						"description": "Movie genres",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"rating": map[string]interface{}{
						"type":        "number",
						"description": "Movie rating (0-10)",
						"minimum":     0,
						"maximum":     10,
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Movie description/plot",
					},
					"duration": map[string]interface{}{
						"type":        "integer",
						"description": "Duration in minutes",
						"minimum":     1,
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Primary language",
					},
					"country": map[string]interface{}{
						"type":        "string",
						"description": "Country of origin",
					},
					"poster_url": map[string]interface{}{
						"type":        "string",
						"description": "URL to movie poster image",
					},
				},
				Required: []string{"title", "director", "year"},
			},
		},
		{
			Name:        "update_movie",
			Description: "Update an existing movie",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "Movie ID to update",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Movie title",
					},
					"director": map[string]interface{}{
						"type":        "string",
						"description": "Movie director",
					},
					"year": map[string]interface{}{
						"type":        "integer",
						"description": "Release year",
						"minimum":     1888,
						"maximum":     2100,
					},
					"genre": map[string]interface{}{
						"type":        "array",
						"description": "Movie genres",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"rating": map[string]interface{}{
						"type":        "number",
						"description": "Movie rating (0-10)",
						"minimum":     0,
						"maximum":     10,
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Movie description/plot",
					},
					"duration": map[string]interface{}{
						"type":        "integer",
						"description": "Duration in minutes",
						"minimum":     1,
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Primary language",
					},
					"country": map[string]interface{}{
						"type":        "string",
						"description": "Country of origin",
					},
					"poster_url": map[string]interface{}{
						"type":        "string",
						"description": "URL to movie poster image",
					},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "delete_movie",
			Description: "Delete a movie by ID",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The movie ID to delete",
					},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "search_movies",
			Description: "Search for movies by title, director, genre, year, or full-text search",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"type": map[string]interface{}{
						"type":        "string",
						"description": "Search type: title, director, genre, year, fulltext",
						"enum":        []string{"title", "director", "genre", "year", "fulltext"},
						"default":     "title",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (1-100)",
						"default":     10,
						"minimum":     1,
						"maximum":     100,
					},
					"offset": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to skip for pagination",
						"default":     0,
						"minimum":     0,
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "list_top_movies",
			Description: "List top-rated movies with optional genre filtering",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (1-100)",
						"default":     10,
						"minimum":     1,
						"maximum":     100,
					},
					"genre": map[string]interface{}{
						"type":        "string",
						"description": "Filter by genre (optional)",
					},
				},
			},
		},
		{
			Name:        "get_movie_poster",
			Description: "Get the poster image for a movie by ID",
			InputSchema: models.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The movie ID",
					},
				},
				Required: []string{"id"},
			},
		},
	}
	
	response := models.ToolsListResponse{
		Tools: tools,
	}
	
	s.sendResult(req.ID, response)
}

// handleToolCall handles tool execution requests
func (s *MoviesServer) handleToolCall(req *models.JSONRPCRequest) {
	var params models.ToolCallRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, models.InvalidParams, "Invalid params", err.Error())
		return
	}

	// Validate tool call request parameters
	paramsMap := make(map[string]interface{})
	if err := json.Unmarshal(req.Params, &paramsMap); err == nil {
		if err := s.validator.ValidateToolsCallRequest(paramsMap); err != nil {
			appErr := err.(*errors.ApplicationError)
			s.sendError(req.ID, models.InvalidParams, appErr.Message, appErr.Details)
			return
		}
	}
	
	switch params.Name {
	case "get_movie":
		s.handleGetMovie(req.ID, params.Arguments)
	case "add_movie":
		s.handleAddMovie(req.ID, params.Arguments)
	case "update_movie":
		s.handleUpdateMovie(req.ID, params.Arguments)
	case "delete_movie":
		s.handleDeleteMovie(req.ID, params.Arguments)
	case "search_movies":
		s.handleSearchMovies(req.ID, params.Arguments)
	case "list_top_movies":
		s.handleListTopMovies(req.ID, params.Arguments)
	case "get_movie_poster":
		s.handleGetMoviePoster(req.ID, params.Arguments)
	default:
		s.sendError(req.ID, models.MethodNotFound, "Unknown tool", fmt.Sprintf("Tool '%s' not found", params.Name))
	}
}

// handleGetMovie handles the get_movie tool
func (s *MoviesServer) handleGetMovie(id interface{}, args map[string]interface{}) {
	// Extract movie ID from arguments
	movieIDArg, ok := args["id"]
	if !ok {
		s.sendError(id, models.InvalidParams, "Missing required parameter", "id is required")
		return
	}
	
	var movieID int
	switch v := movieIDArg.(type) {
	case float64:
		movieID = int(v)
	case int:
		movieID = v
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be a valid integer")
			return
		}
		movieID = parsed
	default:
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be an integer")
		return
	}
	
	if movieID <= 0 {
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be positive")
		return
	}
	
	// Get movie from database
	dbMovie, err := s.db.GetMovie(movieID)
	if err != nil {
		if err.Error() == "movie not found" {
			s.sendError(id, models.InvalidParams, "Movie not found", fmt.Sprintf("Movie with ID %d not found", movieID))
			return
		}
		s.sendError(id, models.InternalError, "Database error", err.Error())
		return
	}
	
	// Convert to models.Movie
	movie := &models.Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		Director:    dbMovie.Director,
		Year:        dbMovie.Year,
		Genre:       dbMovie.Genre,
		Rating:      dbMovie.Rating.Float64,
		Description: dbMovie.Description.String,
		Duration:    int(dbMovie.Duration.Int32),
		Language:    dbMovie.Language.String,
		Country:     dbMovie.Country.String,
		PosterType:  dbMovie.PosterType,
		CreatedAt:   dbMovie.CreatedAt,
		UpdatedAt:   dbMovie.UpdatedAt,
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: fmt.Sprintf("Found movie: %s\n\n%s", movie.ToSummary(), movie.ToJSON()),
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleAddMovie handles the add_movie tool
func (s *MoviesServer) handleAddMovie(id interface{}, args map[string]interface{}) {
	// Validate movie data using request validator
	if err := s.validator.ValidateMovieData(args); err != nil {
		appErr := err.(*errors.ApplicationError)
		s.sendError(id, models.InvalidParams, appErr.Message, appErr.Details)
		return
	}

	// Parse arguments into MovieCreateRequest
	req, err := models.ParseMovieArguments(args)
	if err != nil {
		s.sendError(id, models.InvalidParams, "Invalid arguments", err.Error())
		return
	}
	
	// Validate the request (additional model-level validation)
	if err := req.Validate(); err != nil {
		s.sendError(id, models.InvalidParams, "Validation failed", err.Error())
		return
	}
	
	// Convert to database.Movie
	dbMovie := &database.Movie{
		Title:       req.Title,
		Director:    req.Director,
		Year:        req.Year,
		Genre:       req.Genre,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Language:    sql.NullString{String: req.Language, Valid: req.Language != ""},
		Country:     sql.NullString{String: req.Country, Valid: req.Country != ""},
	}
	
	if req.Rating != nil {
		dbMovie.Rating = sql.NullFloat64{Float64: *req.Rating, Valid: true}
	}
	if req.Duration != nil {
		dbMovie.Duration = sql.NullInt32{Int32: int32(*req.Duration), Valid: true}
	}
	
	// Handle poster URL if provided
	if req.PosterURL != "" {
		if s.imageProcessor != nil {
			posterData, posterType, err := s.imageProcessor.DownloadImageFromURL(req.PosterURL)
			if err != nil {
				// Log the error but don't fail the movie creation
				if s.logger != nil {
					s.logger.Printf("Failed to download poster from %s: %v", req.PosterURL, err)
				}
			} else {
				dbMovie.PosterData = posterData
				dbMovie.PosterType = posterType
			}
		}
	}
	
	// Create movie in database
	if err := s.db.CreateMovie(dbMovie); err != nil {
		s.sendError(id, models.InternalError, "Failed to create movie", err.Error())
		return
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully created movie '%s' with ID %d", dbMovie.Title, dbMovie.ID),
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleUpdateMovie handles the update_movie tool
func (s *MoviesServer) handleUpdateMovie(id interface{}, args map[string]interface{}) {
	// Validate movie data using request validator (skip required validation for updates)
	updateArgs := make(map[string]interface{})
	for k, v := range args {
		if k != "id" { // Don't validate ID field with movie data validator
			updateArgs[k] = v
		}
	}
	if len(updateArgs) > 0 {
		if err := s.validator.ValidateMovieUpdateData(updateArgs); err != nil {
			appErr := err.(*errors.ApplicationError)
			s.sendError(id, models.InvalidParams, appErr.Message, appErr.Details)
			return
		}
	}

	// Parse arguments into MovieUpdateRequest
	req, err := models.ParseUpdateArguments(args)
	if err != nil {
		s.sendError(id, models.InvalidParams, "Invalid arguments", err.Error())
		return
	}
	
	// Validate the request (additional model-level validation)
	if err := req.Validate(); err != nil {
		s.sendError(id, models.InvalidParams, "Validation failed", err.Error())
		return
	}
	
	// Get existing movie to update
	dbMovie, err := s.db.GetMovie(req.ID)
	if err != nil {
		if err.Error() == "movie not found" {
			s.sendError(id, models.InvalidParams, "Movie not found", fmt.Sprintf("Movie with ID %d not found", req.ID))
			return
		}
		s.sendError(id, models.InternalError, "Database error", err.Error())
		return
	}
	
	// Update fields that were provided
	if req.Title != nil {
		dbMovie.Title = *req.Title
	}
	if req.Director != nil {
		dbMovie.Director = *req.Director
	}
	if req.Year != nil {
		dbMovie.Year = *req.Year
	}
	if req.Genre != nil {
		dbMovie.Genre = req.Genre
	}
	if req.Rating != nil {
		dbMovie.Rating = sql.NullFloat64{Float64: *req.Rating, Valid: true}
	}
	if req.Description != nil {
		dbMovie.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Duration != nil {
		dbMovie.Duration = sql.NullInt32{Int32: int32(*req.Duration), Valid: true}
	}
	if req.Language != nil {
		dbMovie.Language = sql.NullString{String: *req.Language, Valid: true}
	}
	if req.Country != nil {
		dbMovie.Country = sql.NullString{String: *req.Country, Valid: true}
	}
	
	// Handle poster URL updates
	if req.PosterURL != nil {
		if *req.PosterURL == "" {
			// Empty URL means remove poster
			dbMovie.PosterData = nil
			dbMovie.PosterType = ""
		} else if s.imageProcessor != nil {
			// Download new poster
			posterData, posterType, err := s.imageProcessor.DownloadImageFromURL(*req.PosterURL)
			if err != nil {
				// Log the error but don't fail the update
				if s.logger != nil {
					s.logger.Printf("Failed to download poster from %s: %v", *req.PosterURL, err)
				}
			} else {
				dbMovie.PosterData = posterData
				dbMovie.PosterType = posterType
			}
		}
	}
	
	// Update movie in database
	if err := s.db.UpdateMovie(dbMovie); err != nil {
		s.sendError(id, models.InternalError, "Failed to update movie", err.Error())
		return
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully updated movie '%s' (ID: %d)", dbMovie.Title, dbMovie.ID),
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleDeleteMovie handles the delete_movie tool
func (s *MoviesServer) handleDeleteMovie(id interface{}, args map[string]interface{}) {
	// Extract movie ID from arguments
	movieIDArg, ok := args["id"]
	if !ok {
		s.sendError(id, models.InvalidParams, "Missing required parameter", "id is required")
		return
	}
	
	var movieID int
	switch v := movieIDArg.(type) {
	case float64:
		movieID = int(v)
	case int:
		movieID = v
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be a valid integer")
			return
		}
		movieID = parsed
	default:
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be an integer")
		return
	}
	
	if movieID <= 0 {
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be positive")
		return
	}
	
	// Get movie info before deletion for confirmation message
	movie, err := s.db.GetMovie(movieID)
	if err != nil {
		if err.Error() == "movie not found" {
			s.sendError(id, models.InvalidParams, "Movie not found", fmt.Sprintf("Movie with ID %d not found", movieID))
			return
		}
		s.sendError(id, models.InternalError, "Database error", err.Error())
		return
	}
	
	// Delete movie from database
	if err := s.db.DeleteMovie(movieID); err != nil {
		s.sendError(id, models.InternalError, "Failed to delete movie", err.Error())
		return
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully deleted movie '%s' (ID: %d)", movie.Title, movieID),
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleSearchMovies handles the search_movies tool
func (s *MoviesServer) handleSearchMovies(id interface{}, args map[string]interface{}) {
	// Validate search query parameters
	if err := s.validator.ValidateSearchQuery(args); err != nil {
		appErr := err.(*errors.ApplicationError)
		s.sendError(id, models.InvalidParams, appErr.Message, appErr.Details)
		return
	}

	// Extract search parameters
	queryArg, ok := args["query"]
	if !ok {
		s.sendError(id, models.InvalidParams, "Missing required parameter", "query is required")
		return
	}
	
	query, ok := queryArg.(string)
	if !ok {
		s.sendError(id, models.InvalidParams, "Invalid parameter", "query must be a string")
		return
	}
	
	if query == "" {
		s.sendError(id, models.InvalidParams, "Invalid parameter", "query cannot be empty")
		return
	}
	
	// Extract optional parameters
	searchType := "title" // default
	if typeArg, ok := args["type"]; ok {
		if t, ok := typeArg.(string); ok {
			searchType = t
		}
	}
	
	limit := 10 // default
	if limitArg, ok := args["limit"]; ok {
		switch v := limitArg.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		}
	}
	
	offset := 0 // default
	if offsetArg, ok := args["offset"]; ok {
		switch v := offsetArg.(type) {
		case float64:
			offset = int(v)
		case int:
			offset = v
		}
	}
	
	// Create search query
	searchQuery := database.SearchQuery{
		Query:  query,
		Type:   searchType,
		Limit:  limit,
		Offset: offset,
	}
	
	// Perform search
	movies, err := s.db.SearchMovies(searchQuery)
	if err != nil {
		s.sendError(id, models.InternalError, "Search failed", err.Error())
		return
	}
	
	// Format response
	var resultText string
	if len(movies) == 0 {
		resultText = fmt.Sprintf("No movies found for query '%s' (type: %s)", query, searchType)
	} else {
		resultText = fmt.Sprintf("Found %d movies for query '%s' (type: %s):\n\n", len(movies), query, searchType)
		for i, movie := range movies {
			rating := movie.Rating.Float64
			if !movie.Rating.Valid {
				rating = 0.0
			}
			resultText += fmt.Sprintf("%d. %s (%d) - Directed by %s - Rating: %.1f\n", 
				i+1+offset, movie.Title, movie.Year, movie.Director, rating)
			if len(movie.Genre) > 0 {
				resultText += fmt.Sprintf("   Genres: %v\n", movie.Genre)
			}
			if movie.Description.Valid && movie.Description.String != "" {
				desc := movie.Description.String
				if len(desc) > 100 {
					desc = desc[:100] + "..."
				}
				resultText += fmt.Sprintf("   %s\n", desc)
			}
			resultText += "\n"
		}
		
		if len(movies) == limit {
			resultText += fmt.Sprintf("Note: Results limited to %d. Use 'offset' parameter for pagination.", limit)
		}
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleListTopMovies handles the list_top_movies tool
func (s *MoviesServer) handleListTopMovies(id interface{}, args map[string]interface{}) {
	// Extract parameters
	limit := 10 // default
	if limitArg, ok := args["limit"]; ok {
		switch v := limitArg.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		}
	}
	
	genre := "" // default (no filter)
	if genreArg, ok := args["genre"]; ok {
		if g, ok := genreArg.(string); ok {
			genre = g
		}
	}
	
	// Get top movies
	movies, err := s.db.ListTopMovies(limit, genre)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to list top movies", err.Error())
		return
	}
	
	// Format response
	var resultText string
	if len(movies) == 0 {
		if genre != "" {
			resultText = fmt.Sprintf("No movies found for genre '%s'", genre)
		} else {
			resultText = "No movies found in database"
		}
	} else {
		if genre != "" {
			resultText = fmt.Sprintf("Top %d %s movies:\n\n", len(movies), genre)
		} else {
			resultText = fmt.Sprintf("Top %d movies:\n\n", len(movies))
		}
		
		for i, movie := range movies {
			rating := movie.Rating.Float64
			if !movie.Rating.Valid {
				rating = 0.0
			}
			resultText += fmt.Sprintf("%d. %s (%d) - Rating: %.1f/10\n", 
				i+1, movie.Title, movie.Year, rating)
			resultText += fmt.Sprintf("   Directed by: %s\n", movie.Director)
			if len(movie.Genre) > 0 {
				resultText += fmt.Sprintf("   Genres: %v\n", movie.Genre)
			}
			if movie.Duration.Valid && movie.Duration.Int32 > 0 {
				resultText += fmt.Sprintf("   Duration: %d minutes\n", movie.Duration.Int32)
			}
			if movie.Language.Valid && movie.Language.String != "" {
				resultText += fmt.Sprintf("   Language: %s\n", movie.Language.String)
			}
			resultText += "\n"
		}
	}
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleGetMoviePoster handles the get_movie_poster tool
func (s *MoviesServer) handleGetMoviePoster(id interface{}, args map[string]interface{}) {
	// Extract movie ID from arguments
	movieIDArg, ok := args["id"]
	if !ok {
		s.sendError(id, models.InvalidParams, "Missing required parameter", "id is required")
		return
	}
	
	var movieID int
	switch v := movieIDArg.(type) {
	case float64:
		movieID = int(v)
	case int:
		movieID = v
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be a valid integer")
			return
		}
		movieID = parsed
	default:
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be an integer")
		return
	}
	
	if movieID <= 0 {
		s.sendError(id, models.InvalidParams, "Invalid parameter", "id must be positive")
		return
	}
	
	// Get poster data from database
	posterData, posterType, err := s.db.GetMoviePoster(movieID)
	if err != nil {
		if err.Error() == "movie not found" {
			s.sendError(id, models.InvalidParams, "Movie not found", fmt.Sprintf("Movie with ID %d not found", movieID))
			return
		}
		s.sendError(id, models.InternalError, "Database error", err.Error())
		return
	}
	
	if len(posterData) == 0 {
		s.sendError(id, models.InvalidParams, "No poster available", fmt.Sprintf("Movie %d has no poster image", movieID))
		return
	}
	
	// Encode poster data as base64
	encodedData := base64.StdEncoding.EncodeToString(posterData)
	
	response := models.ToolCallResponse{
		Content: []models.ContentBlock{
			{
				Type: "image",
				Source: &models.ImageSource{
					Type:      "base64",
					MediaType: posterType,
					Data:      encodedData,
				},
			},
		},
	}
	
	s.sendResult(id, response)
}

// handleResourcesList handles the resources/list request
func (s *MoviesServer) handleResourcesList(req *models.JSONRPCRequest) {
	resources := []models.Resource{
		{
			URI:         "movies://database/all",
			Name:        "All Movies",
			Description: "Complete list of all movies in the database",
			MimeType:    "application/json",
		},
		{
			URI:         "movies://database/stats",
			Name:        "Database Statistics",
			Description: "Statistics about the movie database",
			MimeType:    "application/json",
		},
		{
			URI:         "movies://database/genres",
			Name:        "Movie Genres",
			Description: "List of all genres with movie counts",
			MimeType:    "application/json",
		},
		{
			URI:         "movies://database/directors",
			Name:        "Movie Directors",
			Description: "List of all directors with movie counts and ratings",
			MimeType:    "application/json",
		},
		{
			URI:         "movies://posters/collection",
			Name:        "Poster Collection",
			Description: "Collection of all movie posters with metadata",
			MimeType:    "application/json",
		},
	}
	
	response := models.ResourcesListResponse{
		Resources: resources,
	}
	
	s.sendResult(req.ID, response)
}

// handleResourceRead handles resource read requests
func (s *MoviesServer) handleResourceRead(req *models.JSONRPCRequest) {
	var params models.ResourceReadRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, models.InvalidParams, "Invalid params", err.Error())
		return
	}

	// Validate resource read request parameters
	paramsMap := make(map[string]interface{})
	if err := json.Unmarshal(req.Params, &paramsMap); err == nil {
		if err := s.validator.ValidateResourcesReadRequest(paramsMap); err != nil {
			appErr := err.(*errors.ApplicationError)
			s.sendError(req.ID, models.InvalidParams, appErr.Message, appErr.Details)
			return
		}
	}
	
	// Parse and route based on URI
	switch {
	case params.URI == "movies://database/all":
		s.handleResourceDatabaseAll(req.ID, params)
	case params.URI == "movies://database/stats":
		s.handleResourceDatabaseStats(req.ID, params)
	case params.URI == "movies://database/genres":
		s.handleResourceDatabaseGenres(req.ID, params)
	case params.URI == "movies://database/directors":
		s.handleResourceDatabaseDirectors(req.ID, params)
	case params.URI == "movies://posters/collection":
		s.handleResourcePostersCollection(req.ID, params)
	case strings.HasPrefix(params.URI, "movies://posters/"):
		s.handleResourceIndividualPoster(req.ID, params)
	default:
		s.sendError(req.ID, models.InvalidParams, "Unknown resource URI", fmt.Sprintf("Resource '%s' not found", params.URI))
	}
}

// sendResult sends a successful response
func (s *MoviesServer) sendResult(id interface{}, result interface{}) {
	response := models.NewJSONRPCResponse(id, result, nil)
	s.sendResponse(response)
}

// sendError sends an error response
func (s *MoviesServer) sendError(id interface{}, code int, message string, data interface{}) {
	err := models.NewJSONRPCError(code, message, data)
	response := models.NewJSONRPCResponse(id, nil, err)
	s.sendResponse(response)
}

// sendResponse marshals and sends a response
func (s *MoviesServer) sendResponse(response models.JSONRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		if s.logger != nil {
			s.logger.Printf("Error marshaling response: %v", err)
		}
		return
	}
	
	if s.logger != nil {
		s.logger.Printf("Sending: %s", string(data))
	}
	fmt.Fprintf(s.output, "%s\n", data)
}

// Resource handler methods

// handleResourceDatabaseAll handles the movies://database/all resource
func (s *MoviesServer) handleResourceDatabaseAll(id interface{}, params models.ResourceReadRequest) {
	// Check if database has GetAllMovies method (for mock), otherwise use search
	type AllMoviesGetter interface {
		GetAllMovies() ([]*database.Movie, error)
	}
	
	var movies []*database.Movie
	var err error
	
	if allGetter, ok := s.db.(AllMoviesGetter); ok {
		// Use GetAllMovies for mock database
		movies, err = allGetter.GetAllMovies()
	} else {
		// Use search for real database - search with very broad criteria
		movies, err = s.db.SearchMovies(database.SearchQuery{
			Query: "%", // Wildcard to match all titles
			Type:  "title",
			Limit: 10000,
		})
	}
	
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to retrieve movies", err.Error())
		return
	}

	// Convert to JSON
	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to serialize movies", err.Error())
		return
	}

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: "application/json",
		Text:     string(moviesJSON),
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}

// handleResourceDatabaseStats handles the movies://database/stats resource
func (s *MoviesServer) handleResourceDatabaseStats(id interface{}, params models.ResourceReadRequest) {
	stats, err := s.db.GetStats()
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to retrieve database statistics", err.Error())
		return
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to serialize statistics", err.Error())
		return
	}

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: "application/json",
		Text:     string(statsJSON),
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}

// handleResourceDatabaseGenres handles the movies://database/genres resource
func (s *MoviesServer) handleResourceDatabaseGenres(id interface{}, params models.ResourceReadRequest) {
	genres, err := s.db.GetGenres()
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to retrieve genres", err.Error())
		return
	}

	genresJSON, err := json.Marshal(genres)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to serialize genres", err.Error())
		return
	}

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: "application/json",
		Text:     string(genresJSON),
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}

// handleResourceDatabaseDirectors handles the movies://database/directors resource
func (s *MoviesServer) handleResourceDatabaseDirectors(id interface{}, params models.ResourceReadRequest) {
	directors, err := s.db.GetDirectors()
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to retrieve directors", err.Error())
		return
	}

	directorsJSON, err := json.Marshal(directors)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to serialize directors", err.Error())
		return
	}

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: "application/json",
		Text:     string(directorsJSON),
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}

// handleResourcePostersCollection handles the movies://posters/collection resource
func (s *MoviesServer) handleResourcePostersCollection(id interface{}, params models.ResourceReadRequest) {
	// Check if database has GetMoviesWithPosters method (for mock), otherwise use search
	type PostersGetter interface {
		GetMoviesWithPosters() ([]*database.Movie, error)
	}
	
	var movies []*database.Movie
	var err error
	
	if postersGetter, ok := s.db.(PostersGetter); ok {
		// Use GetMoviesWithPosters for mock database
		movies, err = postersGetter.GetMoviesWithPosters()
	} else {
		// Get all movies first, then filter
		movies, err = s.db.SearchMovies(database.SearchQuery{
			Query: "%", // Wildcard to match all titles
			Type:  "title",
			Limit: 10000,
		})
		if err != nil {
			s.sendError(id, models.InternalError, "Failed to retrieve movies", err.Error())
			return
		}
		
		// Filter to only movies with posters
		var moviesWithPosters []*database.Movie
		for _, movie := range movies {
			if movie.PosterType != "" {
				moviesWithPosters = append(moviesWithPosters, movie)
			}
		}
		movies = moviesWithPosters
	}
	
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to retrieve movies", err.Error())
		return
	}

	// Filter movies with posters and create collection
	var posterCollection []map[string]interface{}
	for _, movie := range movies {
		if movie.PosterType != "" { // Has poster
			posterItem := map[string]interface{}{
				"movie_id":   movie.ID,
				"title":      movie.Title,
				"director":   movie.Director,
				"year":       movie.Year,
				"poster_uri": fmt.Sprintf("movies://posters/%d", movie.ID),
				"mime_type":  movie.PosterType,
			}
			posterCollection = append(posterCollection, posterItem)
		}
	}

	collectionJSON, err := json.Marshal(posterCollection)
	if err != nil {
		s.sendError(id, models.InternalError, "Failed to serialize poster collection", err.Error())
		return
	}

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: "application/json",
		Text:     string(collectionJSON),
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}

// handleResourceIndividualPoster handles movies://posters/{id} resources
func (s *MoviesServer) handleResourceIndividualPoster(id interface{}, params models.ResourceReadRequest) {
	// Extract movie ID from URI (movies://posters/123)
	parts := strings.Split(params.URI, "/")
	if len(parts) < 4 {
		s.sendError(id, models.InvalidParams, "Invalid poster URI", "URI format should be movies://posters/{id}")
		return
	}

	movieIDStr := parts[3] // movies://posters/{id} -> parts[3] is the ID
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		s.sendError(id, models.InvalidParams, "Invalid movie ID", "Movie ID must be a valid integer")
		return
	}

	// Get poster data from database
	posterData, posterType, err := s.db.GetMoviePoster(movieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.sendError(id, models.InvalidParams, "Movie not found", fmt.Sprintf("Movie with ID %d not found", movieID))
		} else {
			s.sendError(id, models.InternalError, "Failed to retrieve poster", err.Error())
		}
		return
	}

	if len(posterData) == 0 {
		s.sendError(id, models.InvalidParams, "No poster available", fmt.Sprintf("Movie %d has no poster", movieID))
		return
	}

	// Encode poster data as base64
	encodedData := base64.StdEncoding.EncodeToString(posterData)

	content := models.ResourceContent{
		URI:      params.URI,
		MimeType: posterType,
		Blob:     encodedData,
	}

	response := models.ResourceReadResponse{
		Contents: []models.ResourceContent{content},
	}

	s.sendResult(id, response)
}