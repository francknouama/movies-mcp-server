package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// DatabaseResources handles movie database resource operations
type DatabaseResources struct {
	movieService *movieApp.Service
}

// NewDatabaseResources creates a new database resources handler
func NewDatabaseResources(movieService *movieApp.Service) *DatabaseResources {
	return &DatabaseResources{
		movieService: movieService,
	}
}

// AllMoviesResource returns the complete movie database resource definition
func (dr *DatabaseResources) AllMoviesResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "movies://database/all",
		Name:        "All Movies",
		Description: "Complete movie database in JSON format",
		MIMEType:    "application/json",
	}
}

// DatabaseStatsResource returns the database statistics resource definition
func (dr *DatabaseResources) DatabaseStatsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "movies://database/stats",
		Name:        "Database Statistics",
		Description: "Movie database statistics and analytics",
		MIMEType:    "application/json",
	}
}

// PosterCollectionResource returns the movie posters collection resource definition
func (dr *DatabaseResources) PosterCollectionResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "movies://posters/collection",
		Name:        "Movie Posters Collection",
		Description: "Collection of all movie posters",
		MIMEType:    "application/json",
	}
}

// HandleAllMovies handles the movies://database/all resource request
func (dr *DatabaseResources) HandleAllMovies(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Fetch all movies from the database
	movies, err := dr.movieService.SearchMovies(ctx, movieApp.SearchMoviesQuery{
		// Empty query returns all movies
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all movies: %w", err)
	}

	// Convert movies to JSON
	moviesJSON, err := json.MarshalIndent(map[string]interface{}{
		"total_movies": len(movies),
		"movies":       movies,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal movies to JSON: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "movies://database/all",
				MIMEType: "application/json",
				Text:     string(moviesJSON),
			},
		},
	}, nil
}

// HandleDatabaseStats handles the movies://database/stats resource request
func (dr *DatabaseResources) HandleDatabaseStats(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Fetch all movies to compute statistics
	movies, err := dr.movieService.SearchMovies(ctx, movieApp.SearchMoviesQuery{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies for stats: %w", err)
	}

	// Compute statistics
	totalMovies := len(movies)
	genreSet := make(map[string]bool)
	var earliestYear, latestYear *int
	var totalRating float64
	ratingCount := 0

	for _, movie := range movies {
		// Collect genres
		for _, genre := range movie.Genres {
			genreSet[genre] = true
		}

		// Track year range
		if earliestYear == nil || movie.Year < *earliestYear {
			earliestYear = &movie.Year
		}
		if latestYear == nil || movie.Year > *latestYear {
			latestYear = &movie.Year
		}

		// Track ratings
		if movie.Rating > 0 {
			totalRating += movie.Rating
			ratingCount++
		}
	}

	// Convert genre set to slice
	genres := make([]string, 0, len(genreSet))
	for genre := range genreSet {
		genres = append(genres, genre)
	}

	// Calculate average rating
	var avgRating float64
	if ratingCount > 0 {
		avgRating = totalRating / float64(ratingCount)
	}

	stats := map[string]interface{}{
		"total_movies":   totalMovies,
		"total_genres":   len(genres),
		"genres":         genres,
		"average_rating": fmt.Sprintf("%.1f", avgRating),
		"year_range": map[string]interface{}{
			"earliest": earliestYear,
			"latest":   latestYear,
		},
	}

	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "movies://database/stats",
				MIMEType: "application/json",
				Text:     string(statsJSON),
			},
		},
	}, nil
}

// HandlePosterCollection handles the movies://posters/collection resource request
func (dr *DatabaseResources) HandlePosterCollection(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Fetch all movies to get poster information
	movies, err := dr.movieService.SearchMovies(ctx, movieApp.SearchMoviesQuery{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies for poster collection: %w", err)
	}

	// Build poster collection (placeholder for now)
	posters := make([]map[string]interface{}, 0, len(movies))
	for _, movie := range movies {
		posters = append(posters, map[string]interface{}{
			"movie_id": movie.ID,
			"title":    movie.Title,
			"year":     movie.Year,
			"uri":      fmt.Sprintf("movies://posters/%d", movie.ID),
		})
	}

	collection := map[string]interface{}{
		"total":   len(posters),
		"posters": posters,
	}

	collectionJSON, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal poster collection to JSON: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "movies://posters/collection",
				MIMEType: "application/json",
				Text:     string(collectionJSON),
			},
		},
	}, nil
}
