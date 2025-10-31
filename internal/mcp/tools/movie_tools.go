package tools

import (
	"context"
	"fmt"
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
