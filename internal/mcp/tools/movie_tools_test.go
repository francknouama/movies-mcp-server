package tools

import (
	"context"
	"errors"
	"testing"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// MockMovieService is a mock implementation for testing
type MockMovieService struct {
	GetMovieFunc          func(ctx context.Context, id int) (*movieApp.MovieDTO, error)
	CreateMovieFunc       func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error)
	UpdateMovieFunc       func(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error)
	DeleteMovieFunc       func(ctx context.Context, id int) error
	SearchMoviesFunc      func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
	GetTopRatedMoviesFunc func(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error)
}

func (m *MockMovieService) GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
	if m.GetMovieFunc != nil {
		return m.GetMovieFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
	if m.CreateMovieFunc != nil {
		return m.CreateMovieFunc(ctx, cmd)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
	if m.UpdateMovieFunc != nil {
		return m.UpdateMovieFunc(ctx, cmd)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) DeleteMovie(ctx context.Context, id int) error {
	if m.DeleteMovieFunc != nil {
		return m.DeleteMovieFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *MockMovieService) SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
	if m.SearchMoviesFunc != nil {
		return m.SearchMoviesFunc(ctx, query)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) GetTopRatedMovies(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error) {
	if m.GetTopRatedMoviesFunc != nil {
		return m.GetTopRatedMoviesFunc(ctx, limit)
	}
	return nil, errors.New("not implemented")
}

func TestGetMovie_Success(t *testing.T) {
	// Arrange
	mockService := &MockMovieService{
		GetMovieFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
			return &movieApp.MovieDTO{
				ID:        42,
				Title:     "The Matrix",
				Director:  "The Wachowskis",
				Year:      1999,
				Rating:    8.7,
				Genres:    []string{"Action", "Sci-Fi"},
				PosterURL: "https://example.com/matrix.jpg",
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-01T00:00:00Z",
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := GetMovieInput{
		MovieID: 42,
	}

	// Act
	result, output, err := tools.GetMovie(ctx, nil, input)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	// Verify output
	if output.ID != 42 {
		t.Errorf("Expected ID to be 42, got: %d", output.ID)
	}

	if output.Title != "The Matrix" {
		t.Errorf("Expected title to be 'The Matrix', got: %s", output.Title)
	}

	if output.Director != "The Wachowskis" {
		t.Errorf("Expected director to be 'The Wachowskis', got: %s", output.Director)
	}

	if output.Year != 1999 {
		t.Errorf("Expected year to be 1999, got: %d", output.Year)
	}

	if output.Rating != 8.7 {
		t.Errorf("Expected rating to be 8.7, got: %f", output.Rating)
	}

	if len(output.Genres) != 2 {
		t.Errorf("Expected 2 genres, got: %d", len(output.Genres))
	}

	if output.PosterURL != "https://example.com/matrix.jpg" {
		t.Errorf("Expected poster URL to be 'https://example.com/matrix.jpg', got: %s", output.PosterURL)
	}

	if output.CreatedAt != "2025-01-01T00:00:00Z" {
		t.Errorf("Expected created_at to be '2025-01-01T00:00:00Z', got: %s", output.CreatedAt)
	}

	if output.UpdatedAt != "2025-01-01T00:00:00Z" {
		t.Errorf("Expected updated_at to be '2025-01-01T00:00:00Z', got: %s", output.UpdatedAt)
	}
}

func TestGetMovie_NotFound(t *testing.T) {
	// Arrange
	mockService := &MockMovieService{
		GetMovieFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
			return nil, errors.New("movie not found")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := GetMovieInput{
		MovieID: 999,
	}

	// Act
	result, output, err := tools.GetMovie(ctx, nil, input)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "movie not found" {
		t.Errorf("Expected error 'movie not found', got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	// Output should be empty on error
	if output.ID != 0 {
		t.Errorf("Expected empty output on error, got ID: %d", output.ID)
	}
}

func TestGetMovie_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockMovieService{
		GetMovieFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := GetMovieInput{
		MovieID: 42,
	}

	// Act
	result, output, err := tools.GetMovie(ctx, nil, input)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should wrap the error
	expectedErrorSubstring := "failed to get movie"
	if err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("Expected error to start with '%s', got: %v", expectedErrorSubstring, err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	// Output should be empty on error
	if output.ID != 0 {
		t.Errorf("Expected empty output on error, got ID: %d", output.ID)
	}
}

func TestGetMovie_ContextCancellation(t *testing.T) {
	// Arrange
	mockService := &MockMovieService{
		GetMovieFunc: func(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
			// Check if context is cancelled
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return &movieApp.MovieDTO{
				ID:    42,
				Title: "The Matrix",
			}, nil
		},
	}

	tools := NewMovieTools(mockService)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	input := GetMovieInput{
		MovieID: 42,
	}

	// Act
	result, output, err := tools.GetMovie(ctx, nil, input)

	// Assert
	if err == nil {
		t.Fatal("Expected error due to cancelled context, got nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	// Output should be empty on error
	if output.ID != 0 {
		t.Errorf("Expected empty output on error, got ID: %d", output.ID)
	}
}
