package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

// MovieHandlers provides MCP handlers for movie operations
type MovieHandlers struct {
	movieService *movieApp.Service
}

// NewMovieHandlers creates a new movie handlers instance
func NewMovieHandlers(movieService *movieApp.Service) *MovieHandlers {
	return &MovieHandlers{
		movieService: movieService,
	}
}

// HandleGetMovie handles the get_movie tool call
func (h *MovieHandlers) HandleGetMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse movie ID
	movieIDFloat, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "movie_id is required and must be a number", nil)
		return
	}
	movieID := int(movieIDFloat)

	// Get movie from service
	ctx := context.Background()
	movieDTO, err := h.movieService.GetMovie(ctx, movieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Movie not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to get movie", err.Error())
		}
		return
	}

	// Convert to response format
	response := h.toMovieResponse(movieDTO)
	sendResult(id, response)
}

// HandleAddMovie handles the add_movie tool call
func (h *MovieHandlers) HandleAddMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseCreateMovieRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid movie data", err.Error())
		return
	}

	// Convert to application command
	cmd := movieApp.CreateMovieCommand{
		Title:     req.Title,
		Director:  req.Director,
		Year:      req.Year,
		Rating:    req.Rating,
		Genres:    req.Genres,
		PosterURL: req.PosterURL,
	}

	// Create movie
	ctx := context.Background()
	movieDTO, err := h.movieService.CreateMovie(ctx, cmd)
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to create movie", err.Error())
		return
	}

	// Convert to response format
	response := h.toMovieResponse(movieDTO)
	sendResult(id, response)
}

// HandleUpdateMovie handles the update_movie tool call
func (h *MovieHandlers) HandleUpdateMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseUpdateMovieRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid movie data", err.Error())
		return
	}

	// Convert to application command
	cmd := movieApp.UpdateMovieCommand{
		ID:        req.ID,
		Title:     req.Title,
		Director:  req.Director,
		Year:      req.Year,
		Rating:    req.Rating,
		Genres:    req.Genres,
		PosterURL: req.PosterURL,
	}

	// Update movie
	ctx := context.Background()
	movieDTO, err := h.movieService.UpdateMovie(ctx, cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Movie not found", nil)
		} else {
			sendError(id, dto.InvalidParams, "Failed to update movie", err.Error())
		}
		return
	}

	// Convert to response format
	response := h.toMovieResponse(movieDTO)
	sendResult(id, response)
}

// HandleDeleteMovie handles the delete_movie tool call
func (h *MovieHandlers) HandleDeleteMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse movie ID
	movieIDFloat, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "movie_id is required and must be a number", nil)
		return
	}
	movieID := int(movieIDFloat)

	// Delete movie
	ctx := context.Background()
	err := h.movieService.DeleteMovie(ctx, movieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Movie not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to delete movie", err.Error())
		}
		return
	}

	sendResult(id, map[string]string{"message": "Movie deleted successfully"})
}

// HandleSearchMovies handles the search_movies tool call
func (h *MovieHandlers) HandleSearchMovies(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseSearchMoviesRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid search parameters", err.Error())
		return
	}

	// Convert to application query
	query := movieApp.SearchMoviesQuery{
		Title:     req.Title,
		Director:  req.Director,
		Genre:     req.Genre,
		MinYear:   req.MinYear,
		MaxYear:   req.MaxYear,
		MinRating: req.MinRating,
		MaxRating: req.MaxRating,
		Limit:     req.Limit,
		Offset:    req.Offset,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
	}

	// Set default limit
	if query.Limit == 0 {
		query.Limit = 20
	}

	// Search movies
	ctx := context.Background()
	movieDTOs, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies", err.Error())
		return
	}

	// Convert to response format
	response := &dto.MoviesListResponse{
		Movies:      make([]*dto.MovieResponse, len(movieDTOs)),
		Total:       len(movieDTOs),
		Description: "Search results",
	}

	for i, movieDTO := range movieDTOs {
		response.Movies[i] = h.toMovieResponse(movieDTO)
	}

	sendResult(id, response)
}

// HandleListTopMovies handles the list_top_movies tool call
func (h *MovieHandlers) HandleListTopMovies(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse limit
	limit := 10 // Default
	if limitFloat, ok := arguments["limit"].(float64); ok {
		limit = int(limitFloat)
	}

	// Get top movies
	ctx := context.Background()
	movieDTOs, err := h.movieService.GetTopRatedMovies(ctx, limit)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to get top movies", err.Error())
		return
	}

	// Convert to response format
	response := &dto.MoviesListResponse{
		Movies:      make([]*dto.MovieResponse, len(movieDTOs)),
		Total:       len(movieDTOs),
		Description: fmt.Sprintf("Top %d rated movies", limit),
	}

	for i, movieDTO := range movieDTOs {
		response.Movies[i] = h.toMovieResponse(movieDTO)
	}

	sendResult(id, response)
}

// HandleSearchByDecade handles the search_by_decade tool call
func (h *MovieHandlers) HandleSearchByDecade(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse decade
	decade, ok := arguments["decade"].(string)
	if !ok || decade == "" {
		sendError(id, dto.InvalidParams, "decade is required", nil)
		return
	}

	// Parse decade to year range
	minYear, maxYear, err := h.parseDecade(decade)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid decade format", err.Error())
		return
	}

	// Search movies in decade
	query := movieApp.SearchMoviesQuery{
		MinYear:  minYear,
		MaxYear:  maxYear,
		Limit:    50,
		OrderBy:  "year",
		OrderDir: "asc",
	}

	ctx := context.Background()
	movieDTOs, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies by decade", err.Error())
		return
	}

	// Convert to response format
	response := &dto.MoviesListResponse{
		Movies:      make([]*dto.MovieResponse, len(movieDTOs)),
		Total:       len(movieDTOs),
		Description: fmt.Sprintf("Movies from the %s", decade),
	}

	for i, movieDTO := range movieDTOs {
		response.Movies[i] = h.toMovieResponse(movieDTO)
	}

	sendResult(id, response)
}

// HandleSearchByRatingRange handles the search_by_rating_range tool call
func (h *MovieHandlers) HandleSearchByRatingRange(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse rating range
	var minRating, maxRating float64
	var hasMin, hasMax bool

	if val, ok := arguments["min_rating"].(float64); ok {
		minRating = val
		hasMin = true
	}
	if val, ok := arguments["max_rating"].(float64); ok {
		maxRating = val
		hasMax = true
	}

	if !hasMin && !hasMax {
		sendError(id, dto.InvalidParams, "At least one of min_rating or max_rating is required", nil)
		return
	}

	// Validate ratings
	if hasMin && (minRating < 0 || minRating > 10) {
		sendError(id, dto.InvalidParams, "Rating must be between 0 and 10", nil)
		return
	}
	if hasMax && (maxRating < 0 || maxRating > 10) {
		sendError(id, dto.InvalidParams, "Rating must be between 0 and 10", nil)
		return
	}
	if hasMin && hasMax && minRating > maxRating {
		sendError(id, dto.InvalidParams, "min_rating cannot be greater than max_rating", nil)
		return
	}

	// Search movies by rating range
	query := movieApp.SearchMoviesQuery{
		MinRating: minRating,
		MaxRating: maxRating,
		Limit:     50,
		OrderBy:   "rating",
		OrderDir:  "desc",
	}

	ctx := context.Background()
	movieDTOs, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies by rating", err.Error())
		return
	}

	// Create description
	var description string
	if hasMin && hasMax {
		description = fmt.Sprintf("Movies with rating between %.1f and %.1f", minRating, maxRating)
	} else if hasMin {
		description = fmt.Sprintf("Movies with rating >= %.1f", minRating)
	} else {
		description = fmt.Sprintf("Movies with rating <= %.1f", maxRating)
	}

	// Convert to response format
	response := &dto.MoviesListResponse{
		Movies:      make([]*dto.MovieResponse, len(movieDTOs)),
		Total:       len(movieDTOs),
		Description: description,
	}

	for i, movieDTO := range movieDTOs {
		response.Movies[i] = h.toMovieResponse(movieDTO)
	}

	sendResult(id, response)
}

// Utility methods

func (h *MovieHandlers) parseCreateMovieRequest(arguments map[string]any) (*dto.CreateMovieRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.CreateMovieRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *MovieHandlers) parseUpdateMovieRequest(arguments map[string]any) (*dto.UpdateMovieRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.UpdateMovieRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *MovieHandlers) parseSearchMoviesRequest(arguments map[string]any) (*dto.SearchMoviesRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.SearchMoviesRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *MovieHandlers) parseDecade(decade string) (int, int, error) {
	decade = strings.TrimSpace(decade)

	// Handle formats like "1990s", "90s", "1990"
	decade = strings.TrimSuffix(decade, "s")

	var baseYear int
	if len(decade) == 2 {
		// Handle "90" -> 1990
		year, err := strconv.Atoi(decade)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid decade format")
		}
		if year >= 0 && year <= 30 {
			baseYear = 2000 + year
		} else {
			baseYear = 1900 + year
		}
	} else if len(decade) == 4 {
		// Handle "1990" -> 1990
		year, err := strconv.Atoi(decade)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid decade format")
		}
		baseYear = year
	} else {
		return 0, 0, fmt.Errorf("invalid decade format")
	}

	// Convert to decade boundaries
	decadeStart := (baseYear / 10) * 10
	decadeEnd := decadeStart + 9

	return decadeStart, decadeEnd, nil
}

func (h *MovieHandlers) toMovieResponse(movieDTO *movieApp.MovieDTO) *dto.MovieResponse {
	return &dto.MovieResponse{
		ID:        movieDTO.ID,
		Title:     movieDTO.Title,
		Director:  movieDTO.Director,
		Year:      movieDTO.Year,
		Rating:    movieDTO.Rating,
		Genres:    movieDTO.Genres,
		PosterURL: movieDTO.PosterURL,
		CreatedAt: movieDTO.CreatedAt,
		UpdatedAt: movieDTO.UpdatedAt,
	}
}
