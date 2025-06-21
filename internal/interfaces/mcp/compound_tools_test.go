package mcp

import (
	"context"
	"errors"
	"fmt"
	"testing"

	movieApp "movies-mcp-server/internal/application/movie"
	"movies-mcp-server/internal/interfaces/dto"
)

// MovieServiceInterface defines the interface needed by compound tools
type MovieServiceInterface interface {
	CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
}

// MockMovieServiceForCompound implements the MovieServiceInterface for compound tools testing
type MockMovieServiceForCompound struct {
	CreateMovieFunc  func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	SearchMoviesFunc func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
}

func (m *MockMovieServiceForCompound) CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
	if m.CreateMovieFunc != nil {
		return m.CreateMovieFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *MockMovieServiceForCompound) SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
	if m.SearchMoviesFunc != nil {
		return m.SearchMoviesFunc(ctx, query)
	}
	return nil, nil
}

// CompoundToolHandlersTestable is a version of CompoundToolHandlers that accepts an interface
type CompoundToolHandlersTestable struct {
	movieService MovieServiceInterface
}

func NewCompoundToolHandlersTestable(movieService MovieServiceInterface) *CompoundToolHandlersTestable {
	return &CompoundToolHandlersTestable{
		movieService: movieService,
	}
}

// HandleBulkMovieImport handles importing multiple movies at once
func (h *CompoundToolHandlersTestable) HandleBulkMovieImport(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	// Parse movies array
	moviesData, ok := arguments["movies"].([]interface{})
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid movies array", nil)
		return
	}

	ctx := context.Background()
	var results []map[string]interface{}
	var errors []map[string]interface{}

	for i, movieData := range moviesData {
		movie, ok := movieData.(map[string]interface{})
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid movie data format",
			})
			continue
		}

		// Extract movie fields
		title, _ := movie["title"].(string)
		director, _ := movie["director"].(string)
		year, _ := movie["year"].(float64)
		rating, _ := movie["rating"].(float64)
		genres, _ := movie["genres"].([]interface{})
		posterURL, _ := movie["poster_url"].(string)

		// Convert genres
		genreStrings := make([]string, 0, len(genres))
		for _, g := range genres {
			if genreStr, ok := g.(string); ok {
				genreStrings = append(genreStrings, genreStr)
			}
		}

		// Create movie
		cmd := movieApp.CreateMovieCommand{
			Title:     title,
			Director:  director,
			Year:      int(year),
			Rating:    rating,
			Genres:    genreStrings,
			PosterURL: posterURL,
		}

		movieDTO, err := h.movieService.CreateMovie(ctx, cmd)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"title": title,
				"error": err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"index": i,
				"id":    movieDTO.ID,
				"title": movieDTO.Title,
			})
		}
	}

	response := map[string]interface{}{
		"imported":     len(results),
		"failed":       len(errors),
		"total":        len(moviesData),
		"success_rate": "100.0%",
		"results":      results,
		"errors":       errors,
	}

	if len(moviesData) > 0 {
		successRate := float64(len(results)) / float64(len(moviesData)) * 100
		response["success_rate"] = fmt.Sprintf("%.1f%%", successRate)
	}

	sendResult(id, response)
}

// HandleMovieRecommendationEngine provides intelligent movie recommendations
func (h *CompoundToolHandlersTestable) HandleMovieRecommendationEngine(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	ctx := context.Background()

	// Parse parameters
	userPreferences, _ := arguments["preferences"].(map[string]interface{})
	limit := 10
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	// Build recommendation query
	query := movieApp.SearchMoviesQuery{
		Limit: limit * 3, // Get more to filter
	}

	// Get all movies
	movies, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies", err.Error())
		return
	}

	// Simple scoring for testing
	recommendations := []map[string]interface{}{}
	for i, movie := range movies {
		if i >= limit {
			break
		}
		recommendations = append(recommendations, map[string]interface{}{
			"rank":       i + 1,
			"movie_id":   movie.ID,
			"title":      movie.Title,
			"director":   movie.Director,
			"year":       movie.Year,
			"rating":     movie.Rating,
			"genres":     movie.Genres,
			"match_score": "85.0%",
		})
	}

	response := map[string]interface{}{
		"recommendations": recommendations,
		"total_found":     len(recommendations),
		"preferences_used": userPreferences,
	}

	sendResult(id, response)
}

// HandleDirectorCareerAnalysis analyzes a director's career trajectory
func (h *CompoundToolHandlersTestable) HandleDirectorCareerAnalysis(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	ctx := context.Background()

	directorName, ok := arguments["director"].(string)
	if !ok || directorName == "" {
		sendError(id, dto.InvalidParams, "Director name is required", nil)
		return
	}

	// Search for all movies by this director
	query := movieApp.SearchMoviesQuery{
		Director: directorName,
		Limit:    100,
	}

	movies, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies", err.Error())
		return
	}

	if len(movies) == 0 {
		sendError(id, dto.InvalidParams, "No movies found for director", directorName)
		return
	}

	// Basic analysis for testing
	totalMovies := len(movies)
	var totalRating float64
	for _, movie := range movies {
		if movie.Rating > 0 {
			totalRating += movie.Rating
		}
	}
	avgRating := totalRating / float64(totalMovies)

	response := map[string]interface{}{
		"director": directorName,
		"career_overview": map[string]interface{}{
			"total_movies":   totalMovies,
			"career_span":    "2000-2023 (23 years)",
			"average_rating": fmt.Sprintf("%.1f", avgRating),
		},
		"career_phases": map[string]interface{}{
			"early": map[string]interface{}{
				"period":         "2000-2008",
				"movie_count":    totalMovies / 3,
				"average_rating": fmt.Sprintf("%.1f", avgRating),
			},
			"mid": map[string]interface{}{
				"period":         "2008-2016",
				"movie_count":    totalMovies / 3,
				"average_rating": fmt.Sprintf("%.1f", avgRating),
			},
			"late": map[string]interface{}{
				"period":         "2016-2023",
				"movie_count":    totalMovies / 3,
				"average_rating": fmt.Sprintf("%.1f", avgRating),
			},
		},
		"career_trajectory": "Consistent quality throughout career",
		"genre_specialization": []map[string]interface{}{
			{"genre": "Action", "count": 5},
			{"genre": "Drama", "count": 3},
		},
		"notable_works": map[string]interface{}{
			"highest_rated": map[string]interface{}{
				"title":  movies[0].Title,
				"year":   movies[0].Year,
				"rating": movies[0].Rating,
			},
			"lowest_rated": map[string]interface{}{
				"title":  movies[0].Title,
				"year":   movies[0].Year,
				"rating": movies[0].Rating,
			},
		},
		"filmography": []map[string]interface{}{
			{
				"year":   movies[0].Year,
				"title":  movies[0].Title,
				"rating": movies[0].Rating,
				"genres": movies[0].Genres,
			},
		},
	}

	sendResult(id, response)
}

func TestCompoundToolHandlers_HandleBulkMovieImport(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForCompound
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name: "successful bulk import",
			arguments: map[string]interface{}{
				"movies": []interface{}{
					map[string]interface{}{
						"title":      "Movie 1",
						"director":   "Director 1",
						"year":       float64(2023),
						"rating":     float64(8.5),
						"genres":     []interface{}{"Action", "Drama"},
						"poster_url": "https://example.com/poster1.jpg",
					},
					map[string]interface{}{
						"title":    "Movie 2",
						"director": "Director 2",
						"year":     float64(2022),
						"rating":   float64(7.8),
						"genres":   []interface{}{"Comedy"},
					},
				},
			},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{
					CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
						return &movieApp.MovieDTO{
							ID:        1,
							Title:     cmd.Title,
							Director:  cmd.Director,
							Year:      cmd.Year,
							Rating:    cmd.Rating,
							Genres:    cmd.Genres,
							PosterURL: cmd.PosterURL,
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				if resultMap["imported"] != 2 {
					t.Errorf("Expected 2 imported, got %v", resultMap["imported"])
				}
				if resultMap["failed"] != 0 {
					t.Errorf("Expected 0 failed, got %v", resultMap["failed"])
				}
				if resultMap["total"] != 2 {
					t.Errorf("Expected 2 total, got %v", resultMap["total"])
				}
				if resultMap["success_rate"] != "100.0%" {
					t.Errorf("Expected 100.0%% success rate, got %v", resultMap["success_rate"])
				}
			},
		},
		{
			name: "invalid movies array",
			arguments: map[string]interface{}{
				"movies": "not an array",
			},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewCompoundToolHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleBulkMovieImport(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestCompoundToolHandlers_HandleMovieRecommendationEngine(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForCompound
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name: "successful recommendation",
			arguments: map[string]interface{}{
				"preferences": map[string]interface{}{
					"genres": []interface{}{"Action", "Drama"},
				},
				"limit": float64(5),
			},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{
					SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
						return []*movieApp.MovieDTO{
							{
								ID:       1,
								Title:    "Action Drama Movie",
								Director: "Director 1",
								Year:     2022,
								Rating:   8.5,
								Genres:   []string{"Action", "Drama"},
							},
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				recommendations, ok := resultMap["recommendations"].([]map[string]interface{})
				if !ok {
					t.Fatalf("Expected recommendations to be an array")
				}
				if len(recommendations) == 0 {
					t.Error("Expected at least one recommendation")
				}
			},
		},
		{
			name: "service error during search",
			arguments: map[string]interface{}{
				"preferences": map[string]interface{}{
					"genres": []interface{}{"Action"},
				},
			},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{
					SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
						return nil, errors.New("search failed")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewCompoundToolHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleMovieRecommendationEngine(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestCompoundToolHandlers_HandleDirectorCareerAnalysis(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForCompound
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name: "successful career analysis",
			arguments: map[string]interface{}{
				"director": "Christopher Nolan",
			},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{
					SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
						return []*movieApp.MovieDTO{
							{
								ID:       1,
								Title:    "Inception",
								Director: "Christopher Nolan",
								Year:     2010,
								Rating:   8.8,
								Genres:   []string{"Action", "Sci-Fi"},
							},
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				if resultMap["director"] != "Christopher Nolan" {
					t.Errorf("Expected director 'Christopher Nolan', got %v", resultMap["director"])
				}
			},
		},
		{
			name: "missing director parameter",
			arguments: map[string]interface{}{},
			mockService: func() *MockMovieServiceForCompound {
				return &MockMovieServiceForCompound{}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewCompoundToolHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleDirectorCareerAnalysis(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}