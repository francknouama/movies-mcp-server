package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// MovieService defines the interface for movie operations
type MovieService interface {
	GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error)
	CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error)
	DeleteMovie(ctx context.Context, id int) error
	SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
	GetTopRatedMovies(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error)
}

// MovieTools provides SDK-based MCP handlers for movie operations
type MovieTools struct {
	movieService MovieService
}

// NewMovieTools creates a new movie tools instance
func NewMovieTools(movieService MovieService) *MovieTools {
	return &MovieTools{
		movieService: movieService,
	}
}

// ===== get_movie Tool =====

// GetMovieInput defines the input schema for get_movie tool
type GetMovieInput struct {
	MovieID int `json:"movie_id" jsonschema:"required,description=The movie ID to retrieve"`
}

// GetMovieOutput defines the output schema for get_movie tool
type GetMovieOutput struct {
	ID        int      `json:"id" jsonschema:"description=Movie ID"`
	Title     string   `json:"title" jsonschema:"description=Movie title"`
	Director  string   `json:"director" jsonschema:"description=Movie director"`
	Year      int      `json:"year" jsonschema:"description=Release year"`
	Rating    float64  `json:"rating,omitempty" jsonschema:"description=Movie rating (0-10)"`
	Genres    []string `json:"genres" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
	CreatedAt string   `json:"created_at" jsonschema:"description=Creation timestamp"`
	UpdatedAt string   `json:"updated_at" jsonschema:"description=Last update timestamp"`
}

// GetMovie handles the get_movie tool call with SDK-compatible signature
func (t *MovieTools) GetMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetMovieInput,
) (*mcp.CallToolResult, GetMovieOutput, error) {
	// Get movie from service
	movieDTO, err := t.movieService.GetMovie(ctx, input.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, GetMovieOutput{}, fmt.Errorf("movie not found")
		}
		return nil, GetMovieOutput{}, fmt.Errorf("failed to get movie: %w", err)
	}

	// Convert to output format
	output := GetMovieOutput{
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

	return nil, output, nil
}

// ===== add_movie Tool =====

// AddMovieInput defines the input schema for add_movie tool
type AddMovieInput struct {
	Title     string   `json:"title" jsonschema:"required,description=Movie title"`
	Director  string   `json:"director" jsonschema:"required,description=Movie director"`
	Year      int      `json:"year" jsonschema:"required,description=Release year"`
	Rating    float64  `json:"rating,omitempty" jsonschema:"description=Movie rating (0-10),minimum=0,maximum=10"`
	Genres    []string `json:"genres,omitempty" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
}

// AddMovieOutput defines the output schema for add_movie tool
type AddMovieOutput struct {
	ID        int      `json:"id" jsonschema:"description=Created movie ID"`
	Title     string   `json:"title" jsonschema:"description=Movie title"`
	Director  string   `json:"director" jsonschema:"description=Movie director"`
	Year      int      `json:"year" jsonschema:"description=Release year"`
	Rating    float64  `json:"rating,omitempty" jsonschema:"description=Movie rating"`
	Genres    []string `json:"genres" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
	CreatedAt string   `json:"created_at" jsonschema:"description=Creation timestamp"`
	UpdatedAt string   `json:"updated_at" jsonschema:"description=Last update timestamp"`
}

// AddMovie handles the add_movie tool call
func (t *MovieTools) AddMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddMovieInput,
) (*mcp.CallToolResult, AddMovieOutput, error) {
	// Create movie command
	cmd := movieApp.CreateMovieCommand{
		Title:     input.Title,
		Director:  input.Director,
		Year:      input.Year,
		Rating:    input.Rating,
		Genres:    input.Genres,
		PosterURL: input.PosterURL,
	}

	// Create movie
	movieDTO, err := t.movieService.CreateMovie(ctx, cmd)
	if err != nil {
		return nil, AddMovieOutput{}, fmt.Errorf("failed to create movie: %w", err)
	}

	// Convert to output format
	output := AddMovieOutput{
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

	return nil, output, nil
}

// ===== update_movie Tool =====

// UpdateMovieInput defines the input schema for update_movie tool
type UpdateMovieInput struct {
	ID        int      `json:"id" jsonschema:"required,description=Movie ID"`
	Title     string   `json:"title" jsonschema:"required,description=Movie title"`
	Director  string   `json:"director" jsonschema:"required,description=Movie director"`
	Year      int      `json:"year" jsonschema:"required,description=Release year"`
	Rating    float64  `json:"rating,omitempty" jsonschema:"description=Movie rating (0-10),minimum=0,maximum=10"`
	Genres    []string `json:"genres,omitempty" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
}

// UpdateMovieOutput defines the output schema for update_movie tool
type UpdateMovieOutput struct {
	ID        int      `json:"id" jsonschema:"description=Updated movie ID"`
	Title     string   `json:"title" jsonschema:"description=Movie title"`
	Director  string   `json:"director" jsonschema:"description=Movie director"`
	Year      int      `json:"year" jsonschema:"description=Release year"`
	Rating    float64  `json:"rating,omitempty" jsonschema:"description=Movie rating"`
	Genres    []string `json:"genres" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
	CreatedAt string   `json:"created_at" jsonschema:"description=Creation timestamp"`
	UpdatedAt string   `json:"updated_at" jsonschema:"description=Last update timestamp"`
}

// UpdateMovie handles the update_movie tool call
func (t *MovieTools) UpdateMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UpdateMovieInput,
) (*mcp.CallToolResult, UpdateMovieOutput, error) {
	// Create update command
	cmd := movieApp.UpdateMovieCommand{
		ID:        input.ID,
		Title:     input.Title,
		Director:  input.Director,
		Year:      input.Year,
		Rating:    input.Rating,
		Genres:    input.Genres,
		PosterURL: input.PosterURL,
	}

	// Update movie
	movieDTO, err := t.movieService.UpdateMovie(ctx, cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, UpdateMovieOutput{}, fmt.Errorf("movie not found")
		}
		return nil, UpdateMovieOutput{}, fmt.Errorf("failed to update movie: %w", err)
	}

	// Convert to output format
	output := UpdateMovieOutput{
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

	return nil, output, nil
}

// ===== delete_movie Tool =====

// DeleteMovieInput defines the input schema for delete_movie tool
type DeleteMovieInput struct {
	MovieID int `json:"movie_id" jsonschema:"required,description=The movie ID to delete"`
}

// DeleteMovieOutput defines the output schema for delete_movie tool
type DeleteMovieOutput struct {
	Message string `json:"message" jsonschema:"description=Success message"`
}

// DeleteMovie handles the delete_movie tool call
func (t *MovieTools) DeleteMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DeleteMovieInput,
) (*mcp.CallToolResult, DeleteMovieOutput, error) {
	// Delete movie
	err := t.movieService.DeleteMovie(ctx, input.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, DeleteMovieOutput{}, fmt.Errorf("movie not found")
		}
		return nil, DeleteMovieOutput{}, fmt.Errorf("failed to delete movie: %w", err)
	}

	output := DeleteMovieOutput{
		Message: "Movie deleted successfully",
	}

	return nil, output, nil
}

// ===== list_top_movies Tool =====

// ListTopMoviesInput defines the input schema for list_top_movies tool
type ListTopMoviesInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"description=Number of movies to return,default=10"`
}

// ListTopMoviesOutput defines the output schema for list_top_movies tool
type ListTopMoviesOutput struct {
	Movies      []GetMovieOutput `json:"movies" jsonschema:"description=List of top-rated movies"`
	Total       int              `json:"total" jsonschema:"description=Total number of movies returned"`
	Description string           `json:"description" jsonschema:"description=Description of results"`
}

// ListTopMovies handles the list_top_movies tool call
func (t *MovieTools) ListTopMovies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ListTopMoviesInput,
) (*mcp.CallToolResult, ListTopMoviesOutput, error) {
	// Set default limit
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	// Get top movies
	movieDTOs, err := t.movieService.GetTopRatedMovies(ctx, limit)
	if err != nil {
		return nil, ListTopMoviesOutput{}, fmt.Errorf("failed to get top movies: %w", err)
	}

	// Convert to output format
	movies := make([]GetMovieOutput, len(movieDTOs))
	for i, movieDTO := range movieDTOs {
		movies[i] = GetMovieOutput{
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

	output := ListTopMoviesOutput{
		Movies:      movies,
		Total:       len(movies),
		Description: fmt.Sprintf("Top %d rated movies", limit),
	}

	return nil, output, nil
}

// ===== search_movies Tool =====

// SearchMoviesInput defines the input schema for search_movies tool
type SearchMoviesInput struct {
	Title     string  `json:"title,omitempty" jsonschema:"description=Search by movie title"`
	Director  string  `json:"director,omitempty" jsonschema:"description=Search by director name"`
	Genre     string  `json:"genre,omitempty" jsonschema:"description=Search by genre"`
	MinYear   int     `json:"min_year,omitempty" jsonschema:"description=Minimum release year"`
	MaxYear   int     `json:"max_year,omitempty" jsonschema:"description=Maximum release year"`
	MinRating float64 `json:"min_rating,omitempty" jsonschema:"description=Minimum rating (0-10)"`
	MaxRating float64 `json:"max_rating,omitempty" jsonschema:"description=Maximum rating (0-10)"`
	Limit     int     `json:"limit,omitempty" jsonschema:"description=Maximum number of results,default=20"`
	Offset    int     `json:"offset,omitempty" jsonschema:"description=Number of results to skip for pagination,default=0"`
	OrderBy   string  `json:"order_by,omitempty" jsonschema:"description=Field to order by (title/year/rating),default=title"`
	OrderDir  string  `json:"order_dir,omitempty" jsonschema:"description=Order direction (asc/desc),default=asc"`
}

// SearchMoviesOutput defines the output schema for search_movies tool
type SearchMoviesOutput struct {
	Movies      []GetMovieOutput `json:"movies" jsonschema:"description=List of matching movies"`
	Total       int              `json:"total" jsonschema:"description=Total number of movies found"`
	Description string           `json:"description" jsonschema:"description=Description of search results"`
}

// SearchMovies handles the search_movies tool call
func (t *MovieTools) SearchMovies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchMoviesInput,
) (*mcp.CallToolResult, SearchMoviesOutput, error) {
	// Create search query
	query := movieApp.SearchMoviesQuery{
		Title:     input.Title,
		Director:  input.Director,
		Genre:     input.Genre,
		MinYear:   input.MinYear,
		MaxYear:   input.MaxYear,
		MinRating: input.MinRating,
		MaxRating: input.MaxRating,
		Limit:     input.Limit,
		Offset:    input.Offset,
		OrderBy:   input.OrderBy,
		OrderDir:  input.OrderDir,
	}

	// Set default limit
	if query.Limit == 0 {
		query.Limit = 20
	}

	// Search movies
	movieDTOs, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, SearchMoviesOutput{}, fmt.Errorf("failed to search movies: %w", err)
	}

	// Convert to output format
	movies := make([]GetMovieOutput, len(movieDTOs))
	for i, movieDTO := range movieDTOs {
		movies[i] = GetMovieOutput{
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

	output := SearchMoviesOutput{
		Movies:      movies,
		Total:       len(movies),
		Description: "Search results",
	}

	return nil, output, nil
}

// ===== search_by_decade Tool =====

// SearchByDecadeInput defines the input schema for search_by_decade tool
type SearchByDecadeInput struct {
	Decade string `json:"decade" jsonschema:"required,description=Decade to search (e.g. '1990s' or '90s' or '1990')"`
}

// SearchByDecade handles the search_by_decade tool call
func (t *MovieTools) SearchByDecade(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchByDecadeInput,
) (*mcp.CallToolResult, SearchMoviesOutput, error) {
	// Parse decade to year range
	minYear, maxYear, err := parseDecade(input.Decade)
	if err != nil {
		return nil, SearchMoviesOutput{}, fmt.Errorf("invalid decade format: %w", err)
	}

	// Create search query
	query := movieApp.SearchMoviesQuery{
		MinYear:  minYear,
		MaxYear:  maxYear,
		Limit:    50,
		OrderBy:  "year",
		OrderDir: "asc",
	}

	// Search movies
	movieDTOs, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, SearchMoviesOutput{}, fmt.Errorf("failed to search movies by decade: %w", err)
	}

	// Convert to output format
	movies := make([]GetMovieOutput, len(movieDTOs))
	for i, movieDTO := range movieDTOs {
		movies[i] = GetMovieOutput{
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

	output := SearchMoviesOutput{
		Movies:      movies,
		Total:       len(movies),
		Description: fmt.Sprintf("Movies from the %s", input.Decade),
	}

	return nil, output, nil
}

// ===== search_by_rating_range Tool =====

// SearchByRatingRangeInput defines the input schema for search_by_rating_range tool
type SearchByRatingRangeInput struct {
	MinRating float64 `json:"min_rating,omitempty" jsonschema:"description=Minimum rating (0-10)"`
	MaxRating float64 `json:"max_rating,omitempty" jsonschema:"description=Maximum rating (0-10)"`
}

// SearchByRatingRange handles the search_by_rating_range tool call
func (t *MovieTools) SearchByRatingRange(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchByRatingRangeInput,
) (*mcp.CallToolResult, SearchMoviesOutput, error) {
	// Validate that at least one rating is provided
	if input.MinRating == 0 && input.MaxRating == 0 {
		return nil, SearchMoviesOutput{}, fmt.Errorf("at least one of min_rating or max_rating is required")
	}

	// Validate rating ranges
	if input.MinRating < 0 || input.MinRating > 10 {
		return nil, SearchMoviesOutput{}, fmt.Errorf("min_rating must be between 0 and 10")
	}
	if input.MaxRating < 0 || input.MaxRating > 10 {
		return nil, SearchMoviesOutput{}, fmt.Errorf("max_rating must be between 0 and 10")
	}
	if input.MinRating > 0 && input.MaxRating > 0 && input.MinRating > input.MaxRating {
		return nil, SearchMoviesOutput{}, fmt.Errorf("min_rating cannot be greater than max_rating")
	}

	// Create search query
	query := movieApp.SearchMoviesQuery{
		MinRating: input.MinRating,
		MaxRating: input.MaxRating,
		Limit:     50,
		OrderBy:   "rating",
		OrderDir:  "desc",
	}

	// Search movies
	movieDTOs, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, SearchMoviesOutput{}, fmt.Errorf("failed to search movies by rating: %w", err)
	}

	// Convert to output format
	movies := make([]GetMovieOutput, len(movieDTOs))
	for i, movieDTO := range movieDTOs {
		movies[i] = GetMovieOutput{
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

	// Create description
	var description string
	if input.MinRating > 0 && input.MaxRating > 0 {
		description = fmt.Sprintf("Movies with rating between %.1f and %.1f", input.MinRating, input.MaxRating)
	} else if input.MinRating > 0 {
		description = fmt.Sprintf("Movies with rating >= %.1f", input.MinRating)
	} else {
		description = fmt.Sprintf("Movies with rating <= %.1f", input.MaxRating)
	}

	output := SearchMoviesOutput{
		Movies:      movies,
		Total:       len(movies),
		Description: description,
	}

	return nil, output, nil
}

// Helper function to parse decade string
func parseDecade(decade string) (int, int, error) {
	decade = strings.TrimSpace(decade)
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
