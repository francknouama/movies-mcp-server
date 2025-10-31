package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// MovieTools provides SDK-based MCP handlers for movie operations
type MovieTools struct {
	movieService *movieApp.Service
}

// NewMovieTools creates a new movie tools instance
func NewMovieTools(movieService *movieApp.Service) *MovieTools {
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
