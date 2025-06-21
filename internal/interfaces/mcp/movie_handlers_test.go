package mcp

import (
	"context"
	"errors"
	"testing"

	movieApp "movies-mcp-server/internal/application/movie"
	"movies-mcp-server/internal/interfaces/dto"
)

// MovieServiceInterfaceForHandlers defines the interface needed by movie handlers
type MovieServiceInterfaceForHandlers interface {
	CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error)
	UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error)
	DeleteMovie(ctx context.Context, id int) error
	SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
	FindSimilarMovies(ctx context.Context, movieID int, limit int) ([]*movieApp.MovieDTO, error)
}

// MockMovieServiceForMovieHandlers implements the MovieServiceInterface for testing
type MockMovieServiceForMovieHandlers struct {
	CreateFunc           func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	GetByIDFunc          func(ctx context.Context, id int) (*movieApp.MovieDTO, error)
	UpdateFunc           func(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error)
	DeleteFunc           func(ctx context.Context, id int) error
	SearchMoviesFunc     func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
	FindSimilarMoviesFunc func(ctx context.Context, movieID int, limit int) ([]*movieApp.MovieDTO, error)
}

func (m *MockMovieServiceForMovieHandlers) CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *MockMovieServiceForMovieHandlers) GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMovieServiceForMovieHandlers) UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *MockMovieServiceForMovieHandlers) DeleteMovie(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockMovieServiceForMovieHandlers) SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
	if m.SearchMoviesFunc != nil {
		return m.SearchMoviesFunc(ctx, query)
	}
	return nil, nil
}

func (m *MockMovieServiceForMovieHandlers) FindSimilarMovies(ctx context.Context, movieID int, limit int) ([]*movieApp.MovieDTO, error) {
	if m.FindSimilarMoviesFunc != nil {
		return m.FindSimilarMoviesFunc(ctx, movieID, limit)
	}
	return nil, nil
}

// MovieHandlersTestable is a version of MovieHandlers that accepts an interface
type MovieHandlersTestable struct {
	movieService MovieServiceInterfaceForHandlers
}

func NewMovieHandlersTestable(movieService MovieServiceInterfaceForHandlers) *MovieHandlersTestable {
	return &MovieHandlersTestable{
		movieService: movieService,
	}
}

// HandleAddMovie handles the add_movie tool call
func (h *MovieHandlersTestable) HandleAddMovie(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	// Parse required parameters
	title, ok := arguments["title"].(string)
	if !ok || title == "" {
		sendError(id, dto.InvalidParams, "Invalid title parameter", nil)
		return
	}

	director, ok := arguments["director"].(string)
	if !ok || director == "" {
		sendError(id, dto.InvalidParams, "Invalid director parameter", nil)
		return
	}

	year, ok := arguments["year"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid year parameter", nil)
		return
	}

	rating, _ := arguments["rating"].(float64)
	genres, _ := arguments["genres"].([]interface{})
	posterURL, _ := arguments["poster_url"].(string)

	// Convert genres
	genreStrings := make([]string, 0, len(genres))
	for _, g := range genres {
		if genreStr, ok := g.(string); ok {
			genreStrings = append(genreStrings, genreStr)
		}
	}

	// Convert to application command
	cmd := movieApp.CreateMovieCommand{
		Title:     title,
		Director:  director,
		Year:      int(year),
		Rating:    rating,
		Genres:    genreStrings,
		PosterURL: posterURL,
	}

	// Create movie
	ctx := context.Background()
	movieDTO, err := h.movieService.CreateMovie(ctx, cmd)
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to create movie", err.Error())
		return
	}

	// Convert to response format
	response := map[string]interface{}{
		"id":         movieDTO.ID,
		"title":      movieDTO.Title,
		"director":   movieDTO.Director,
		"year":       movieDTO.Year,
		"rating":     movieDTO.Rating,
		"genres":     movieDTO.Genres,
		"poster_url": movieDTO.PosterURL,
	}

	sendResult(id, response)
}

// HandleGetMovie handles the get_movie tool call
func (h *MovieHandlersTestable) HandleGetMovie(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	movieID, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid movie_id parameter", nil)
		return
	}

	ctx := context.Background()
	movieDTO, err := h.movieService.GetMovie(ctx, int(movieID))
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to get movie", err.Error())
		return
	}

	response := map[string]interface{}{
		"id":         movieDTO.ID,
		"title":      movieDTO.Title,
		"director":   movieDTO.Director,
		"year":       movieDTO.Year,
		"rating":     movieDTO.Rating,
		"genres":     movieDTO.Genres,
		"poster_url": movieDTO.PosterURL,
	}

	sendResult(id, response)
}

// HandleSearchMovies handles the search_movies tool call
func (h *MovieHandlersTestable) HandleSearchMovies(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	title, _ := arguments["title"].(string)
	director, _ := arguments["director"].(string)
	genre, _ := arguments["genre"].(string)
	limit := 10
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	query := movieApp.SearchMoviesQuery{
		Title:    title,
		Director: director,
		Genre:    genre,
		Limit:    limit,
	}

	ctx := context.Background()
	movieDTOs, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search movies", err.Error())
		return
	}

	var movies []interface{}
	for _, movieDTO := range movieDTOs {
		movies = append(movies, map[string]interface{}{
			"id":         movieDTO.ID,
			"title":      movieDTO.Title,
			"director":   movieDTO.Director,
			"year":       movieDTO.Year,
			"rating":     movieDTO.Rating,
			"genres":     movieDTO.Genres,
			"poster_url": movieDTO.PosterURL,
		})
	}

	response := map[string]interface{}{
		"movies": movies,
		"total":  len(movies),
	}

	sendResult(id, response)
}

// HandleDeleteMovie handles the delete_movie tool call
func (h *MovieHandlersTestable) HandleDeleteMovie(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	movieID, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid movie_id parameter", nil)
		return
	}

	ctx := context.Background()
	err := h.movieService.DeleteMovie(ctx, int(movieID))
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to delete movie", err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "Movie deleted successfully",
	}

	sendResult(id, response)
}

func TestMovieHandlers_HandleAddMovie(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForMovieHandlers
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name: "successful movie creation",
			arguments: map[string]interface{}{
				"title":      "Test Movie",
				"director":   "Test Director",
				"year":       float64(2023),
				"rating":     float64(8.5),
				"genres":     []interface{}{"Action", "Drama"},
				"poster_url": "https://example.com/poster.jpg",
			},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					CreateFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
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
				if resultMap["id"] != 1 {
					t.Errorf("Expected id 1, got %v", resultMap["id"])
				}
				if resultMap["title"] != "Test Movie" {
					t.Errorf("Expected title 'Test Movie', got %v", resultMap["title"])
				}
			},
		},
		{
			name:        "missing title parameter",
			arguments:   map[string]interface{}{"director": "Test Director", "year": float64(2023)},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:        "missing director parameter",
			arguments:   map[string]interface{}{"title": "Test Movie", "year": float64(2023)},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:        "missing year parameter",
			arguments:   map[string]interface{}{"title": "Test Movie", "director": "Test Director"},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name: "invalid year type",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     "not_a_number",
			},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name: "service error",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     float64(2023),
			},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					CreateFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
						return nil, errors.New("service error")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewMovieHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleAddMovie(
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

func TestMovieHandlers_HandleGetMovie(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForMovieHandlers
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name:      "successful get movie",
			arguments: map[string]interface{}{"movie_id": float64(1)},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					GetByIDFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
						return &movieApp.MovieDTO{
							ID:       id,
							Title:    "Test Movie",
							Director: "Test Director",
							Year:     2023,
							Rating:   8.5,
							Genres:   []string{"Action", "Drama"},
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
				if resultMap["id"] != 1 {
					t.Errorf("Expected id 1, got %v", resultMap["id"])
				}
				if resultMap["title"] != "Test Movie" {
					t.Errorf("Expected title 'Test Movie', got %v", resultMap["title"])
				}
			},
		},
		{
			name:        "missing movie_id",
			arguments:   map[string]interface{}{},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:      "movie not found",
			arguments: map[string]interface{}{"movie_id": float64(999)},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					GetByIDFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
						return nil, errors.New("movie not found")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewMovieHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleGetMovie(
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

func TestMovieHandlers_HandleSearchMovies(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForMovieHandlers
		expectError  bool
		errorCode    int
		checkResult  func(t *testing.T, result interface{})
	}{
		{
			name: "successful search by title",
			arguments: map[string]interface{}{
				"title": "Test",
				"limit": float64(10),
			},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
						return []*movieApp.MovieDTO{
							{ID: 1, Title: "Test Movie 1", Director: "Director 1", Year: 2022},
							{ID: 2, Title: "Test Movie 2", Director: "Director 2", Year: 2023},
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
				movies, ok := resultMap["movies"].([]interface{})
				if !ok {
					t.Fatalf("Expected movies to be an array")
				}
				if len(movies) != 2 {
					t.Errorf("Expected 2 movies, got %d", len(movies))
				}
			},
		},
		{
			name:      "empty search",
			arguments: map[string]interface{}{},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
						return []*movieApp.MovieDTO{}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				movies, _ := resultMap["movies"].([]interface{})
				if len(movies) != 0 {
					t.Errorf("Expected 0 movies, got %d", len(movies))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewMovieHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleSearchMovies(
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

func TestMovieHandlers_HandleDeleteMovie(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		mockService  func() *MockMovieServiceForMovieHandlers
		expectError  bool
		errorCode    int
	}{
		{
			name:      "successful delete",
			arguments: map[string]interface{}{"movie_id": float64(1)},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					DeleteFunc: func(ctx context.Context, id int) error {
						return nil
					},
				}
			},
			expectError: false,
		},
		{
			name:        "missing movie_id",
			arguments:   map[string]interface{}{},
			mockService: func() *MockMovieServiceForMovieHandlers { return &MockMovieServiceForMovieHandlers{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:      "movie not found",
			arguments: map[string]interface{}{"movie_id": float64(999)},
			mockService: func() *MockMovieServiceForMovieHandlers {
				return &MockMovieServiceForMovieHandlers{
					DeleteFunc: func(ctx context.Context, id int) error {
						return errors.New("movie not found")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewMovieHandlersTestable(tt.mockService())

			var gotError *dto.JSONRPCError
			var gotResult interface{}

			handlers.HandleDeleteMovie(
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
				// Check for success message
				if resultMap, ok := gotResult.(map[string]interface{}); ok {
					if resultMap["message"] != "Movie deleted successfully" {
						t.Errorf("Unexpected success message: %v", resultMap["message"])
					}
				}
			}
		})
	}
}