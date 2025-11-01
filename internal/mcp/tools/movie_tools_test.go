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

// ===== AddMovie Tests =====

func TestAddMovie_Success(t *testing.T) {
	mockService := &MockMovieService{
		CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
			return &movieApp.MovieDTO{
				ID:        1,
				Title:     cmd.Title,
				Director:  cmd.Director,
				Year:      cmd.Year,
				Rating:    cmd.Rating,
				Genres:    cmd.Genres,
				PosterURL: cmd.PosterURL,
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-01T00:00:00Z",
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := AddMovieInput{
		Title:     "Inception",
		Director:  "Christopher Nolan",
		Year:      2010,
		Rating:    8.8,
		Genres:    []string{"Sci-Fi", "Action"},
		PosterURL: "https://example.com/inception.jpg",
	}

	result, output, err := tools.AddMovie(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if output.ID != 1 {
		t.Errorf("Expected ID to be 1, got: %d", output.ID)
	}

	if output.Title != "Inception" {
		t.Errorf("Expected title 'Inception', got: %s", output.Title)
	}

	if output.Director != "Christopher Nolan" {
		t.Errorf("Expected director 'Christopher Nolan', got: %s", output.Director)
	}

	if output.Year != 2010 {
		t.Errorf("Expected year 2010, got: %d", output.Year)
	}

	if output.Rating != 8.8 {
		t.Errorf("Expected rating 8.8, got: %f", output.Rating)
	}

	if len(output.Genres) != 2 {
		t.Errorf("Expected 2 genres, got: %d", len(output.Genres))
	}
}

func TestAddMovie_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := AddMovieInput{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}

	_, _, err := tools.AddMovie(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to create movie")] != "failed to create movie" {
		t.Errorf("Expected error starting with 'failed to create movie', got: %v", err)
	}
}

// ===== UpdateMovie Tests =====

func TestUpdateMovie_Success(t *testing.T) {
	mockService := &MockMovieService{
		UpdateMovieFunc: func(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
			return &movieApp.MovieDTO{
				ID:        cmd.ID,
				Title:     cmd.Title,
				Director:  cmd.Director,
				Year:      cmd.Year,
				Rating:    cmd.Rating,
				Genres:    cmd.Genres,
				PosterURL: cmd.PosterURL,
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-02T00:00:00Z",
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := UpdateMovieInput{
		ID:        1,
		Title:     "Updated Title",
		Director:  "Updated Director",
		Year:      2021,
		Rating:    9.0,
		Genres:    []string{"Drama"},
		PosterURL: "https://example.com/updated.jpg",
	}

	result, output, err := tools.UpdateMovie(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if output.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got: %s", output.Title)
	}

	if output.Rating != 9.0 {
		t.Errorf("Expected rating 9.0, got: %f", output.Rating)
	}

	if output.UpdatedAt != "2025-01-02T00:00:00Z" {
		t.Errorf("Expected updated timestamp, got: %s", output.UpdatedAt)
	}
}

func TestUpdateMovie_NotFound(t *testing.T) {
	mockService := &MockMovieService{
		UpdateMovieFunc: func(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
			return nil, errors.New("movie not found")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := UpdateMovieInput{
		ID:       999,
		Title:    "Test",
		Director: "Test",
		Year:     2020,
	}

	_, _, err := tools.UpdateMovie(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "movie not found" {
		t.Errorf("Expected 'movie not found' error, got: %v", err)
	}
}

func TestUpdateMovie_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		UpdateMovieFunc: func(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := UpdateMovieInput{
		ID:       1,
		Title:    "Test",
		Director: "Test",
		Year:     2020,
	}

	_, _, err := tools.UpdateMovie(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to update movie")] != "failed to update movie" {
		t.Errorf("Expected error starting with 'failed to update movie', got: %v", err)
	}
}

// ===== DeleteMovie Tests =====

func TestDeleteMovie_Success(t *testing.T) {
	mockService := &MockMovieService{
		DeleteMovieFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := DeleteMovieInput{
		MovieID: 1,
	}

	result, output, err := tools.DeleteMovie(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if output.Message != "Movie deleted successfully" {
		t.Errorf("Expected success message, got: %s", output.Message)
	}
}

func TestDeleteMovie_NotFound(t *testing.T) {
	mockService := &MockMovieService{
		DeleteMovieFunc: func(ctx context.Context, id int) error {
			return errors.New("movie not found")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := DeleteMovieInput{
		MovieID: 999,
	}

	_, _, err := tools.DeleteMovie(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "movie not found" {
		t.Errorf("Expected 'movie not found' error, got: %v", err)
	}
}

func TestDeleteMovie_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		DeleteMovieFunc: func(ctx context.Context, id int) error {
			return errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := DeleteMovieInput{
		MovieID: 1,
	}

	_, _, err := tools.DeleteMovie(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to delete movie")] != "failed to delete movie" {
		t.Errorf("Expected error starting with 'failed to delete movie', got: %v", err)
	}
}

// ===== ListTopMovies Tests =====

func TestListTopMovies_Success(t *testing.T) {
	mockService := &MockMovieService{
		GetTopRatedMoviesFunc: func(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{
					ID:        1,
					Title:     "The Shawshank Redemption",
					Director:  "Frank Darabont",
					Year:      1994,
					Rating:    9.3,
					Genres:    []string{"Drama"},
					CreatedAt: "2025-01-01T00:00:00Z",
					UpdatedAt: "2025-01-01T00:00:00Z",
				},
				{
					ID:        2,
					Title:     "The Godfather",
					Director:  "Francis Ford Coppola",
					Year:      1972,
					Rating:    9.2,
					Genres:    []string{"Crime", "Drama"},
					CreatedAt: "2025-01-01T00:00:00Z",
					UpdatedAt: "2025-01-01T00:00:00Z",
				},
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := ListTopMoviesInput{
		Limit: 10,
	}

	result, output, err := tools.ListTopMovies(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if len(output.Movies) != 2 {
		t.Errorf("Expected 2 movies, got: %d", len(output.Movies))
	}

	if output.Total != 2 {
		t.Errorf("Expected total to be 2, got: %d", output.Total)
	}

	if output.Description != "Top 10 rated movies" {
		t.Errorf("Expected description 'Top 10 rated movies', got: %s", output.Description)
	}

	if output.Movies[0].Title != "The Shawshank Redemption" {
		t.Errorf("Expected first movie 'The Shawshank Redemption', got: %s", output.Movies[0].Title)
	}

	if output.Movies[0].Rating != 9.3 {
		t.Errorf("Expected first movie rating 9.3, got: %f", output.Movies[0].Rating)
	}
}

func TestListTopMovies_DefaultLimit(t *testing.T) {
	mockService := &MockMovieService{
		GetTopRatedMoviesFunc: func(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error) {
			if limit != 10 {
				t.Errorf("Expected default limit 10, got: %d", limit)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := ListTopMoviesInput{
		Limit: 0, // Should default to 10
	}

	_, _, err := tools.ListTopMovies(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestListTopMovies_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		GetTopRatedMoviesFunc: func(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := ListTopMoviesInput{
		Limit: 10,
	}

	_, _, err := tools.ListTopMovies(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to get top movies")] != "failed to get top movies" {
		t.Errorf("Expected error starting with 'failed to get top movies', got: %v", err)
	}
}

// ===== SearchMovies Tests =====

func TestSearchMovies_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{
					ID:        1,
					Title:     "Inception",
					Director:  "Christopher Nolan",
					Year:      2010,
					Rating:    8.8,
					Genres:    []string{"Sci-Fi", "Action"},
					CreatedAt: "2025-01-01T00:00:00Z",
					UpdatedAt: "2025-01-01T00:00:00Z",
				},
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchMoviesInput{
		Director: "Christopher Nolan",
		Genre:    "Sci-Fi",
		MinYear:  2010,
		MaxYear:  2020,
		Limit:    20,
	}

	result, output, err := tools.SearchMovies(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if len(output.Movies) != 1 {
		t.Errorf("Expected 1 movie, got: %d", len(output.Movies))
	}

	if output.Total != 1 {
		t.Errorf("Expected total to be 1, got: %d", output.Total)
	}

	if output.Movies[0].Title != "Inception" {
		t.Errorf("Expected movie 'Inception', got: %s", output.Movies[0].Title)
	}
}

func TestSearchMovies_DefaultLimit(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.Limit != 20 {
				t.Errorf("Expected default limit 20, got: %d", query.Limit)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchMoviesInput{
		Title: "Test",
	}

	_, _, err := tools.SearchMovies(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestSearchMovies_WithAllFilters(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Verify all filters are passed through
			if query.Title != "Inception" {
				t.Errorf("Expected title 'Inception', got: %s", query.Title)
			}
			if query.Director != "Christopher Nolan" {
				t.Errorf("Expected director 'Christopher Nolan', got: %s", query.Director)
			}
			if query.Genre != "Sci-Fi" {
				t.Errorf("Expected genre 'Sci-Fi', got: %s", query.Genre)
			}
			if query.MinYear != 2010 {
				t.Errorf("Expected min year 2010, got: %d", query.MinYear)
			}
			if query.MaxYear != 2020 {
				t.Errorf("Expected max year 2020, got: %d", query.MaxYear)
			}
			if query.MinRating != 8.0 {
				t.Errorf("Expected min rating 8.0, got: %f", query.MinRating)
			}
			if query.MaxRating != 9.0 {
				t.Errorf("Expected max rating 9.0, got: %f", query.MaxRating)
			}
			if query.OrderBy != "rating" {
				t.Errorf("Expected order by 'rating', got: %s", query.OrderBy)
			}
			if query.OrderDir != "desc" {
				t.Errorf("Expected order dir 'desc', got: %s", query.OrderDir)
			}
			if query.Limit != 30 {
				t.Errorf("Expected limit 30, got: %d", query.Limit)
			}
			if query.Offset != 10 {
				t.Errorf("Expected offset 10, got: %d", query.Offset)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchMoviesInput{
		Title:     "Inception",
		Director:  "Christopher Nolan",
		Genre:     "Sci-Fi",
		MinYear:   2010,
		MaxYear:   2020,
		MinRating: 8.0,
		MaxRating: 9.0,
		Limit:     30,
		Offset:    10,
		OrderBy:   "rating",
		OrderDir:  "desc",
	}

	_, _, err := tools.SearchMovies(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestSearchMovies_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchMoviesInput{
		Title: "Test",
	}

	_, _, err := tools.SearchMovies(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to search movies")] != "failed to search movies" {
		t.Errorf("Expected error starting with 'failed to search movies', got: %v", err)
	}
}

// ===== SearchByDecade Tests =====

func TestSearchByDecade_Success_1990s(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MinYear != 1990 || query.MaxYear != 1999 {
				t.Errorf("Expected year range 1990-1999, got: %d-%d", query.MinYear, query.MaxYear)
			}
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "The Matrix",
					Director: "The Wachowskis",
					Year:     1999,
					Rating:   8.7,
				},
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByDecadeInput{
		Decade: "1990s",
	}

	result, output, err := tools.SearchByDecade(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if len(output.Movies) != 1 {
		t.Errorf("Expected 1 movie, got: %d", len(output.Movies))
	}

	if output.Description != "Movies from the 1990s" {
		t.Errorf("Expected description 'Movies from the 1990s', got: %s", output.Description)
	}
}

func TestSearchByDecade_Success_90(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MinYear != 1990 || query.MaxYear != 1999 {
				t.Errorf("Expected year range 1990-1999, got: %d-%d", query.MinYear, query.MaxYear)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByDecadeInput{
		Decade: "90",
	}

	_, _, err := tools.SearchByDecade(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestSearchByDecade_Success_2000s(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MinYear != 2000 || query.MaxYear != 2009 {
				t.Errorf("Expected year range 2000-2009, got: %d-%d", query.MinYear, query.MaxYear)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByDecadeInput{
		Decade: "2000",
	}

	_, _, err := tools.SearchByDecade(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestSearchByDecade_InvalidFormat(t *testing.T) {
	mockService := &MockMovieService{}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByDecadeInput{
		Decade: "invalid",
	}

	_, _, err := tools.SearchByDecade(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error for invalid decade format, got nil")
	}

	if err.Error()[:len("invalid decade format")] != "invalid decade format" {
		t.Errorf("Expected error starting with 'invalid decade format', got: %v", err)
	}
}

func TestSearchByDecade_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByDecadeInput{
		Decade: "1990s",
	}

	_, _, err := tools.SearchByDecade(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to search movies by decade")] != "failed to search movies by decade" {
		t.Errorf("Expected error starting with 'failed to search movies by decade', got: %v", err)
	}
}

// ===== SearchByRatingRange Tests =====

func TestSearchByRatingRange_Success_BothLimits(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MinRating != 8.0 || query.MaxRating != 9.0 {
				t.Errorf("Expected rating range 8.0-9.0, got: %f-%f", query.MinRating, query.MaxRating)
			}
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "Inception",
					Director: "Christopher Nolan",
					Year:     2010,
					Rating:   8.8,
				},
			}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByRatingRangeInput{
		MinRating: 8.0,
		MaxRating: 9.0,
	}

	result, output, err := tools.SearchByRatingRange(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got: %v", result)
	}

	if len(output.Movies) != 1 {
		t.Errorf("Expected 1 movie, got: %d", len(output.Movies))
	}

	if output.Description != "Movies with rating between 8.0 and 9.0" {
		t.Errorf("Expected specific description, got: %s", output.Description)
	}
}

func TestSearchByRatingRange_MinOnly(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MinRating != 8.5 {
				t.Errorf("Expected min rating 8.5, got: %f", query.MinRating)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByRatingRangeInput{
		MinRating: 8.5,
	}

	_, output, err := tools.SearchByRatingRange(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Description != "Movies with rating >= 8.5" {
		t.Errorf("Expected min rating description, got: %s", output.Description)
	}
}

func TestSearchByRatingRange_MaxOnly(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			if query.MaxRating != 7.0 {
				t.Errorf("Expected max rating 7.0, got: %f", query.MaxRating)
			}
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByRatingRangeInput{
		MaxRating: 7.0,
	}

	_, output, err := tools.SearchByRatingRange(ctx, nil, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Description != "Movies with rating <= 7.0" {
		t.Errorf("Expected max rating description, got: %s", output.Description)
	}
}

func TestSearchByRatingRange_NoRating(t *testing.T) {
	mockService := &MockMovieService{}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByRatingRangeInput{}

	_, _, err := tools.SearchByRatingRange(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error for no rating provided, got nil")
	}

	if err.Error() != "at least one of min_rating or max_rating is required" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestSearchByRatingRange_InvalidMinRating(t *testing.T) {
	mockService := &MockMovieService{}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	tests := []struct {
		name      string
		minRating float64
		maxRating float64
		wantErr   string
	}{
		{
			name:      "min rating too low",
			minRating: -1.0,
			maxRating: 5.0,
			wantErr:   "min_rating must be between 0 and 10",
		},
		{
			name:      "min rating too high",
			minRating: 11.0,
			maxRating: 0,
			wantErr:   "min_rating must be between 0 and 10",
		},
		{
			name:      "max rating too low",
			minRating: 5.0,
			maxRating: -1.0,
			wantErr:   "max_rating must be between 0 and 10",
		},
		{
			name:      "max rating too high",
			minRating: 0,
			maxRating: 11.0,
			wantErr:   "max_rating must be between 0 and 10",
		},
		{
			name:      "min greater than max",
			minRating: 8.0,
			maxRating: 7.0,
			wantErr:   "min_rating cannot be greater than max_rating",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := SearchByRatingRangeInput{
				MinRating: tt.minRating,
				MaxRating: tt.maxRating,
			}

			_, _, err := tools.SearchByRatingRange(ctx, nil, input)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err.Error() != tt.wantErr {
				t.Errorf("Expected error '%s', got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestSearchByRatingRange_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewMovieTools(mockService)
	ctx := context.Background()

	input := SearchByRatingRangeInput{
		MinRating: 8.0,
		MaxRating: 9.0,
	}

	_, _, err := tools.SearchByRatingRange(ctx, nil, input)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error()[:len("failed to search movies by rating")] != "failed to search movies by rating" {
		t.Errorf("Expected error starting with 'failed to search movies by rating', got: %v", err)
	}
}
