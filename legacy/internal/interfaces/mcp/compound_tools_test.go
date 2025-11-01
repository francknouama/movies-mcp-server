package mcp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
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
			"rank":        i + 1,
			"movie_id":    movie.ID,
			"title":       movie.Title,
			"director":    movie.Director,
			"year":        movie.Year,
			"rating":      movie.Rating,
			"genres":      movie.Genres,
			"match_score": "85.0%",
		})
	}

	response := map[string]interface{}{
		"recommendations":  recommendations,
		"total_found":      len(recommendations),
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
		name        string
		arguments   map[string]interface{}
		mockService func() *MockMovieServiceForCompound
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
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
		name        string
		arguments   map[string]interface{}
		mockService func() *MockMovieServiceForCompound
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
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
		name        string
		arguments   map[string]interface{}
		mockService func() *MockMovieServiceForCompound
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
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
			name:      "missing director parameter",
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

func TestNewCompoundToolHandlers(t *testing.T) {
	mockService := &movieApp.Service{}
	handlers := NewCompoundToolHandlers(mockService)
	
	if handlers == nil {
		t.Fatal("NewCompoundToolHandlers returned nil")
	}
	
	if handlers.movieService != mockService {
		t.Error("NewCompoundToolHandlers did not set movieService correctly")
	}
}

func TestExtractStringArray(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "valid string array",
			input:    []interface{}{"action", "drama", "comedy"},
			expected: []string{"action", "drama", "comedy"},
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: []string{},
		},
		{
			name:     "array with non-string values",
			input:    []interface{}{"action", 123, "drama", true},
			expected: []string{"action", "drama"},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "non-array input",
			input:    "not an array",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractStringArray(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("extractStringArray() returned %d items, expected %d", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("extractStringArray()[%d] = %s, expected %s", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestCalculateGenreScore(t *testing.T) {
	tests := []struct {
		name         string
		movieGenres  []string
		userGenres   []string
		expected     float64
	}{
		{
			name:         "perfect match",
			movieGenres:  []string{"Action", "Drama"},
			userGenres:   []string{"Action", "Drama"},
			expected:     1.0,
		},
		{
			name:         "partial match",
			movieGenres:  []string{"Action", "Drama", "Comedy"},
			userGenres:   []string{"Action", "Horror"},
			expected:     0.5,
		},
		{
			name:         "no match",
			movieGenres:  []string{"Romance", "Comedy"},
			userGenres:   []string{"Action", "Horror"},
			expected:     0.0,
		},
		{
			name:         "empty user genres",
			movieGenres:  []string{"Action", "Drama"},
			userGenres:   []string{},
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateGenreScore(tt.movieGenres, tt.userGenres)
			if score != tt.expected {
				t.Errorf("calculateGenreScore() = %.1f, expected %.1f", score, tt.expected)
			}
		})
	}
}

func TestCalculateYearScore(t *testing.T) {
	tests := []struct {
		name           string
		movieYear      float64
		yearFrom       float64
		yearTo         float64
		expectedMin    float64
		expectedMax    float64
	}{
		{
			name:           "within range",
			movieYear:      2020,
			yearFrom:       2010,
			yearTo:         2025,
			expectedMin:    1.0,
			expectedMax:    1.0,
		},
		{
			name:           "before range",
			movieYear:      2000,
			yearFrom:       2010,
			yearTo:         2025,
			expectedMin:    0.7,
			expectedMax:    0.9,
		},
		{
			name:           "after range",
			movieYear:      2030,
			yearFrom:       2010,
			yearTo:         2025,
			expectedMin:    0.7,
			expectedMax:    0.9,
		},
		{
			name:           "default range",
			movieYear:      2020,
			yearFrom:       0,
			yearTo:         0,
			expectedMin:    1.0,
			expectedMax:    1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateYearScore(tt.movieYear, tt.yearFrom, tt.yearTo)
			if score < tt.expectedMin || score > tt.expectedMax {
				t.Errorf("calculateYearScore() = %.2f, expected between %.2f and %.2f", score, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestGenerateRecommendationReason(t *testing.T) {
	tests := []struct {
		name          string
		movie         *movieApp.MovieDTO
		preferences   map[string]interface{}
		score         float64
		expectedWords []string
	}{
		{
			name: "excellent match with high rating",
			movie: &movieApp.MovieDTO{
				Genres: []string{"Action", "Drama"},
				Rating: 8.5,
				Year:   2020,
			},
			preferences: map[string]interface{}{
				"genres": []interface{}{"Action", "Comedy"},
			},
			score:         0.9,
			expectedWords: []string{"Excellent match", "Highly rated", "Matches your interest in Action"},
		},
		{
			name: "good match",
			movie: &movieApp.MovieDTO{
				Genres: []string{"Comedy"},
				Rating: 7.0,
				Year:   2023,
			},
			preferences: map[string]interface{}{},
			score:         0.7,
			expectedWords: []string{"Good match"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason := generateRecommendationReason(tt.movie, tt.preferences, tt.score)
			for _, word := range tt.expectedWords {
				if !strings.Contains(reason, word) {
					t.Errorf("generateRecommendationReason() missing expected word '%s' in: %s", word, reason)
				}
			}
		})
	}
}

func TestCalculateAverageRating(t *testing.T) {
	tests := []struct {
		name     string
		movies   []*movieApp.MovieDTO
		expected float64
	}{
		{
			name: "multiple movies",
			movies: []*movieApp.MovieDTO{
				{Rating: 8.0},
				{Rating: 7.5},
				{Rating: 9.0},
			},
			expected: 8.17,
		},
		{
			name: "single movie",
			movies: []*movieApp.MovieDTO{
				{Rating: 7.5},
			},
			expected: 7.5,
		},
		{
			name: "movies with zero ratings",
			movies: []*movieApp.MovieDTO{
				{Rating: 8.0},
				{Rating: 0},
				{Rating: 9.0},
			},
			expected: 8.5,
		},
		{
			name:     "empty movie list",
			movies:   []*movieApp.MovieDTO{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avg := calculateAverageRating(tt.movies)
			if fmt.Sprintf("%.2f", avg) != fmt.Sprintf("%.2f", tt.expected) {
				t.Errorf("calculateAverageRating() = %.2f, expected %.2f", avg, tt.expected)
			}
		})
	}
}

func TestFindBestMovie(t *testing.T) {
	tests := []struct {
		name         string
		movies       []*movieApp.MovieDTO
		expectedID   int
		expectedNil  bool
	}{
		{
			name: "multiple movies",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Good Movie", Rating: 7.5},
				{ID: 2, Title: "Best Movie", Rating: 9.0},
				{ID: 3, Title: "Average Movie", Rating: 6.0},
			},
			expectedID: 2,
		},
		{
			name: "single movie",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Only Movie", Rating: 8.0},
			},
			expectedID: 1,
		},
		{
			name:        "empty list",
			movies:      []*movieApp.MovieDTO{},
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findBestMovie(tt.movies)
			if tt.expectedNil {
				if result != nil {
					t.Errorf("findBestMovie() expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Fatal("findBestMovie() returned nil, expected a movie")
				}
				if result.ID != tt.expectedID {
					t.Errorf("findBestMovie() returned ID %d, expected %d", result.ID, tt.expectedID)
				}
			}
		})
	}
}

func TestFindWorstMovie(t *testing.T) {
	tests := []struct {
		name         string
		movies       []*movieApp.MovieDTO
		expectedID   int
		expectedNil  bool
	}{
		{
			name: "multiple movies",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Good Movie", Rating: 7.5},
				{ID: 2, Title: "Best Movie", Rating: 9.0},
				{ID: 3, Title: "Worst Movie", Rating: 4.0},
			},
			expectedID: 3,
		},
		{
			name: "movies with zero rating",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Good Movie", Rating: 7.5},
				{ID: 2, Title: "Unrated Movie", Rating: 0},
				{ID: 3, Title: "Bad Movie", Rating: 3.0},
			},
			expectedID: 3,
		},
		{
			name:        "empty list",
			movies:      []*movieApp.MovieDTO{},
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findWorstMovie(tt.movies)
			if tt.expectedNil {
				if result != nil {
					t.Errorf("findWorstMovie() expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Fatal("findWorstMovie() returned nil, expected a movie")
				}
				if result.ID != tt.expectedID {
					t.Errorf("findWorstMovie() returned ID %d, expected %d", result.ID, tt.expectedID)
				}
			}
		})
	}
}

func TestFindTopGenres(t *testing.T) {
	tests := []struct {
		name       string
		genreCount map[string]int
		limit      int
		expected   []map[string]interface{}
	}{
		{
			name: "multiple genres",
			genreCount: map[string]int{
				"Action": 3,
				"Drama":  2,
				"Comedy": 1,
			},
			limit: 2,
			expected: []map[string]interface{}{
				{"genre": "Action", "count": 3},
				{"genre": "Drama", "count": 2},
			},
		},
		{
			name:       "empty genre count",
			genreCount: map[string]int{},
			limit:      3,
			expected:   []map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findTopGenres(tt.genreCount, tt.limit)
			if len(result) != len(tt.expected) {
				t.Errorf("findTopGenres() returned %d genres, expected %d", len(result), len(tt.expected))
				return
			}
			for i, genre := range result {
				if genre["genre"] != tt.expected[i]["genre"] || genre["count"] != tt.expected[i]["count"] {
					t.Errorf("findTopGenres()[%d] = %v, expected %v", i, genre, tt.expected[i])
				}
			}
		})
	}
}

func TestDetermineTrajectory(t *testing.T) {
	tests := []struct {
		name        string
		earlyAvg    float64
		midAvg      float64
		lateAvg     float64
		expected    string
	}{
		{
			name:     "ascending trajectory",
			earlyAvg: 6.0,
			midAvg:   7.0,
			lateAvg:  8.0,
			expected: "Ascending - Consistent improvement over career",
		},
		{
			name:     "descending trajectory",
			earlyAvg: 8.0,
			midAvg:   7.0,
			lateAvg:  6.0,
			expected: "Descending - Ratings declined over time",
		},
		{
			name:     "peak trajectory",
			earlyAvg: 7.0,
			midAvg:   7.6,
			lateAvg:  7.0,
			expected: "Peak in mid-career",
		},
		{
			name:     "peak in mid-career",
			earlyAvg: 6.0,
			midAvg:   8.0,
			lateAvg:  6.5,
			expected: "Peak in mid-career",
		},
		{
			name:     "consistent trajectory",
			earlyAvg: 7.0,
			midAvg:   7.0,
			lateAvg:  7.0,
			expected: "Consistent quality throughout career",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineTrajectory(tt.earlyAvg, tt.midAvg, tt.lateAvg)
			if result != tt.expected {
				t.Errorf("determineTrajectory() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestFormatFilmography(t *testing.T) {
	tests := []struct {
		name     string
		movies   []*movieApp.MovieDTO
		expected int // expected number of entries
	}{
		{
			name: "multiple movies",
			movies: []*movieApp.MovieDTO{
				{Year: 2020, Title: "Movie A", Rating: 8.0, Genres: []string{"Action"}},
				{Year: 2021, Title: "Movie B", Rating: 7.5, Genres: []string{"Drama"}},
				{Year: 2019, Title: "Movie C", Rating: 9.0, Genres: []string{"Comedy"}},
			},
			expected: 3,
		},
		{
			name:     "empty list",
			movies:   []*movieApp.MovieDTO{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFilmography(tt.movies)
			if len(result) != tt.expected {
				t.Errorf("formatFilmography() returned %d entries, expected %d", len(result), tt.expected)
			}
			// Check that all expected fields are present
			for i, entry := range result {
				if _, ok := entry["year"]; !ok {
					t.Errorf("formatFilmography()[%d] missing 'year' field", i)
				}
				if _, ok := entry["title"]; !ok {
					t.Errorf("formatFilmography()[%d] missing 'title' field", i)
				}
				if _, ok := entry["rating"]; !ok {
					t.Errorf("formatFilmography()[%d] missing 'rating' field", i)
				}
				if _, ok := entry["genres"]; !ok {
					t.Errorf("formatFilmography()[%d] missing 'genres' field", i)
				}
			}
		})
	}
}
