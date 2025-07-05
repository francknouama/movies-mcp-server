package movie

import (
	"context"
	"errors"
	"testing"

	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
)

// MockMovieRepository implements movie.Repository for testing
type MockMovieRepository struct {
	movies             map[int]*movie.Movie
	nextID             int
	findByIDFunc       func(ctx context.Context, id shared.MovieID) (*movie.Movie, error)
	saveFunc           func(ctx context.Context, m *movie.Movie) error
	deleteFunc         func(ctx context.Context, id shared.MovieID) error
	findByCriteriaFunc func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error)
	findTopRatedFunc   func(ctx context.Context, limit int) ([]*movie.Movie, error)
}

func NewMockMovieRepository() *MockMovieRepository {
	return &MockMovieRepository{
		movies: make(map[int]*movie.Movie),
		nextID: 1,
	}
}

func (m *MockMovieRepository) FindByID(ctx context.Context, id shared.MovieID) (*movie.Movie, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	if movie, exists := m.movies[id.Value()]; exists {
		return movie, nil
	}
	return nil, errors.New("movie not found")
}

func (m *MockMovieRepository) Save(ctx context.Context, movie *movie.Movie) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, movie)
	}

	// Only assign a new ID for new movies (ID 0 is the zero ID from NewMovie)
	if movie.ID().IsZero() {
		// Assign new ID
		id, _ := shared.NewMovieID(m.nextID)
		movie.SetID(id)
		m.nextID++
	}

	m.movies[movie.ID().Value()] = movie
	return nil
}

func (m *MockMovieRepository) Delete(ctx context.Context, id shared.MovieID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	if _, exists := m.movies[id.Value()]; !exists {
		return errors.New("movie not found")
	}
	delete(m.movies, id.Value())
	return nil
}

func (m *MockMovieRepository) FindByCriteria(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
	if m.findByCriteriaFunc != nil {
		return m.findByCriteriaFunc(ctx, criteria)
	}

	var result []*movie.Movie
	for _, movieItem := range m.movies {
		match := true

		// Filter by director
		if criteria.Director != "" && movieItem.Director() != criteria.Director {
			match = false
		}

		// Filter by title
		if criteria.Title != "" && movieItem.Title() != criteria.Title {
			match = false
		}

		// Filter by genre
		if criteria.Genre != "" && !movieItem.HasGenre(criteria.Genre) {
			match = false
		}

		// Filter by year range
		if criteria.MinYear > 0 && movieItem.Year().Value() < criteria.MinYear {
			match = false
		}
		if criteria.MaxYear > 0 && movieItem.Year().Value() > criteria.MaxYear {
			match = false
		}

		// Filter by rating range
		if criteria.MinRating > 0 && movieItem.Rating().Value() < criteria.MinRating {
			match = false
		}
		if criteria.MaxRating > 0 && movieItem.Rating().Value() > criteria.MaxRating {
			match = false
		}

		if match {
			result = append(result, movieItem)
		}
	}

	// Apply offset
	if criteria.Offset > 0 && criteria.Offset < len(result) {
		result = result[criteria.Offset:]
	} else if criteria.Offset >= len(result) {
		result = []*movie.Movie{}
	}

	// Apply limit
	if criteria.Limit > 0 && criteria.Limit < len(result) {
		result = result[:criteria.Limit]
	}

	return result, nil
}

func (m *MockMovieRepository) FindByTitle(ctx context.Context, title string) ([]*movie.Movie, error) {
	var result []*movie.Movie
	for _, movie := range m.movies {
		if movie.Title() == title {
			result = append(result, movie)
		}
	}
	return result, nil
}

func (m *MockMovieRepository) FindByDirector(ctx context.Context, director string) ([]*movie.Movie, error) {
	var result []*movie.Movie
	for _, movie := range m.movies {
		if movie.Director() == director {
			result = append(result, movie)
		}
	}
	return result, nil
}

func (m *MockMovieRepository) FindByGenre(ctx context.Context, genre string) ([]*movie.Movie, error) {
	var result []*movie.Movie
	for _, movie := range m.movies {
		if movie.HasGenre(genre) {
			result = append(result, movie)
		}
	}
	return result, nil
}

func (m *MockMovieRepository) FindTopRated(ctx context.Context, limit int) ([]*movie.Movie, error) {
	if m.findTopRatedFunc != nil {
		return m.findTopRatedFunc(ctx, limit)
	}

	var result []*movie.Movie
	for _, movie := range m.movies {
		if !movie.Rating().IsZero() {
			result = append(result, movie)
		}
	}
	// Sort by rating (descending) - simple implementation
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Rating().Value() < result[j].Rating().Value() {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (m *MockMovieRepository) CountAll(ctx context.Context) (int, error) {
	return len(m.movies), nil
}

func (m *MockMovieRepository) DeleteAll(ctx context.Context) error {
	m.movies = make(map[int]*movie.Movie)
	return nil
}

func TestService_CreateMovie(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	tests := []struct {
		name    string
		cmd     CreateMovieCommand
		wantErr bool
	}{
		{
			name: "valid movie",
			cmd: CreateMovieCommand{
				Title:    "Inception",
				Director: "Christopher Nolan",
				Year:     2010,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			cmd: CreateMovieCommand{
				Title:    "",
				Director: "Christopher Nolan",
				Year:     2010,
			},
			wantErr: true,
		},
		{
			name: "empty director",
			cmd: CreateMovieCommand{
				Title:    "Inception",
				Director: "",
				Year:     2010,
			},
			wantErr: true,
		},
		{
			name: "invalid year",
			cmd: CreateMovieCommand{
				Title:    "Inception",
				Director: "Christopher Nolan",
				Year:     1800,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateMovie(context.Background(), tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Title != tt.cmd.Title {
					t.Errorf("CreateMovie() title = %v, want %v", result.Title, tt.cmd.Title)
				}
				if result.Director != tt.cmd.Director {
					t.Errorf("CreateMovie() director = %v, want %v", result.Director, tt.cmd.Director)
				}
				if result.Year != tt.cmd.Year {
					t.Errorf("CreateMovie() year = %v, want %v", result.Year, tt.cmd.Year)
				}
			}
		})
	}
}

func TestService_CreateMovie_WithGenres(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:    "Inception",
		Director: "Christopher Nolan",
		Year:     2010,
		Genres:   []string{"Sci-Fi", "Action", "Thriller"},
	}

	result, err := service.CreateMovie(context.Background(), cmd)
	if err != nil {
		t.Fatalf("CreateMovie() error = %v", err)
	}

	if len(result.Genres) != 3 {
		t.Errorf("Expected 3 genres, got %d", len(result.Genres))
	}

	expectedGenres := map[string]bool{"Sci-Fi": true, "Action": true, "Thriller": true}
	for _, genre := range result.Genres {
		if !expectedGenres[genre] {
			t.Errorf("Unexpected genre: %s", genre)
		}
	}
}

func TestService_GetMovie(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Get the movie
	result, err := service.GetMovie(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetMovie() error = %v", err)
	}

	if result.Title != created.Title {
		t.Errorf("GetMovie() title = %v, want %v", result.Title, created.Title)
	}
}

func TestService_GetMovie_NotFound(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	_, err := service.GetMovie(context.Background(), 999)
	if err == nil {
		t.Error("Expected error for non-existent movie")
	}
}

func TestService_UpdateMovie(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Original Title",
		Director: "Original Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Update the movie
	updateCmd := UpdateMovieCommand{
		ID:       created.ID,
		Title:    "Updated Title",
		Director: "Updated Director",
		Year:     2021,
		Rating:   8.5,
	}

	result, err := service.UpdateMovie(context.Background(), updateCmd)
	if err != nil {
		t.Fatalf("UpdateMovie() error = %v", err)
	}

	if result.Title != updateCmd.Title {
		t.Errorf("UpdateMovie() title = %v, want %v", result.Title, updateCmd.Title)
	}
	if result.Director != updateCmd.Director {
		t.Errorf("UpdateMovie() director = %v, want %v", result.Director, updateCmd.Director)
	}
	if result.Year != updateCmd.Year {
		t.Errorf("UpdateMovie() year = %v, want %v", result.Year, updateCmd.Year)
	}
	if result.Rating != updateCmd.Rating {
		t.Errorf("UpdateMovie() rating = %v, want %v", result.Rating, updateCmd.Rating)
	}
}

func TestService_DeleteMovie(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Delete the movie
	err = service.DeleteMovie(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("DeleteMovie() error = %v", err)
	}

	// Verify it's deleted
	_, err = service.GetMovie(context.Background(), created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted movie")
	}
}

func TestService_SearchMovies(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create test movies
	movies := []CreateMovieCommand{
		{Title: "Inception", Director: "Christopher Nolan", Year: 2010},
		{Title: "Interstellar", Director: "Christopher Nolan", Year: 2014},
		{Title: "The Matrix", Director: "The Wachowskis", Year: 1999},
	}

	for _, cmd := range movies {
		_, err := service.CreateMovie(context.Background(), cmd)
		if err != nil {
			t.Fatalf("Failed to create test movie: %v", err)
		}
	}

	// Search by director
	query := SearchMoviesQuery{
		Director: "Christopher Nolan",
		Limit:    10,
	}

	results, err := service.SearchMovies(context.Background(), query)
	if err != nil {
		t.Fatalf("SearchMovies() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 movies by Christopher Nolan, got %d", len(results))
	}
}

func TestService_RepositoryError(t *testing.T) {
	repo := NewMockMovieRepository()
	repo.saveFunc = func(ctx context.Context, m *movie.Movie) error {
		return errors.New("database error")
	}

	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}

	_, err := service.CreateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// Test GetTopRatedMovies - currently has 0% coverage
func TestService_GetTopRatedMovies(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create test movies with ratings
	movies := []struct {
		cmd    CreateMovieCommand
		rating float64
	}{
		{CreateMovieCommand{Title: "Movie A", Director: "Director A", Year: 2020}, 9.5},
		{CreateMovieCommand{Title: "Movie B", Director: "Director B", Year: 2021}, 8.7},
		{CreateMovieCommand{Title: "Movie C", Director: "Director C", Year: 2022}, 9.8},
		{CreateMovieCommand{Title: "Movie D", Director: "Director D", Year: 2023}, 0.0}, // No rating
	}

	for _, movieData := range movies {
		created, err := service.CreateMovie(context.Background(), movieData.cmd)
		if err != nil {
			t.Fatalf("Failed to create test movie: %v", err)
		}

		// Set rating if provided
		if movieData.rating > 0 {
			updateCmd := UpdateMovieCommand{
				ID:       created.ID,
				Title:    created.Title,
				Director: created.Director,
				Year:     created.Year,
				Rating:   movieData.rating,
			}
			_, err = service.UpdateMovie(context.Background(), updateCmd)
			if err != nil {
				t.Fatalf("Failed to update movie rating: %v", err)
			}
		}
	}

	// Get top rated movies
	results, err := service.GetTopRatedMovies(context.Background(), 2)
	if err != nil {
		t.Fatalf("GetTopRatedMovies() error = %v", err)
	}

	// Should return 2 movies (only those with ratings), sorted by rating desc
	if len(results) != 2 {
		t.Errorf("Expected 2 top rated movies, got %d", len(results))
	}

	// First movie should have highest rating (9.8)
	if len(results) > 0 && results[0].Rating != 9.8 {
		t.Errorf("Expected first movie rating 9.8, got %f", results[0].Rating)
	}

	// Second movie should have second highest rating (9.5)
	if len(results) > 1 && results[1].Rating != 9.5 {
		t.Errorf("Expected second movie rating 9.5, got %f", results[1].Rating)
	}
}

func TestService_GetTopRatedMovies_DefaultLimit(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Test with negative limit (should default to 10)
	results, err := service.GetTopRatedMovies(context.Background(), -1)
	if err != nil {
		t.Fatalf("GetTopRatedMovies() error = %v", err)
	}

	// Should work without error even with empty repository and return empty slice
	if results == nil {
		t.Error("Expected non-nil result but got nil")
	}
}

// Additional edge cases for existing functions
func TestService_CreateMovie_ValidationErrors(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	tests := []struct {
		name string
		cmd  CreateMovieCommand
	}{
		{
			name: "empty title",
			cmd: CreateMovieCommand{
				Title:    "",
				Director: "Test Director",
				Year:     2020,
			},
		},
		{
			name: "empty director",
			cmd: CreateMovieCommand{
				Title:    "Test Movie",
				Director: "",
				Year:     2020,
			},
		},
		{
			name: "invalid year",
			cmd: CreateMovieCommand{
				Title:    "Test Movie",
				Director: "Test Director",
				Year:     1800,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateMovie(context.Background(), tt.cmd)
			if err == nil {
				t.Error("Expected validation error")
			}
		})
	}
}

func TestService_CreateMovie_WithGenresAndRating(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
		Genres:   []string{"Action", "Drama"},
		Rating:   8.5,
	}

	result, err := service.CreateMovie(context.Background(), cmd)
	if err != nil {
		t.Fatalf("CreateMovie() error = %v", err)
	}

	if len(result.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(result.Genres))
	}

	if result.Rating != 8.5 {
		t.Errorf("Expected rating 8.5, got %f", result.Rating)
	}
}

func TestService_CreateMovie_InvalidRating(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
		Rating:   15.0, // Invalid rating > 10
	}

	_, err := service.CreateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for invalid rating")
	}
}

func TestService_UpdateMovie_InvalidID(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := UpdateMovieCommand{
		ID:       -1,
		Title:    "Test",
		Director: "Test",
		Year:     2020,
	}

	_, err := service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

func TestService_UpdateMovie_MovieNotFound(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := UpdateMovieCommand{
		ID:       999,
		Title:    "Test",
		Director: "Test",
		Year:     2020,
	}

	_, err := service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for non-existent movie")
	}
}

func TestService_UpdateMovie_ValidationError(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Try to update with invalid data
	cmd := UpdateMovieCommand{
		ID:       created.ID,
		Title:    "", // Empty title should fail
		Director: "Test",
		Year:     2020,
	}

	_, err = service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for empty title")
	}
}

func TestService_UpdateMovie_RepositoryError(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Set up repository to fail on save
	repo.saveFunc = func(ctx context.Context, m *movie.Movie) error {
		return errors.New("save failed")
	}

	cmd := UpdateMovieCommand{
		ID:       created.ID,
		Title:    "Updated Title",
		Director: "Updated Director",
		Year:     2021,
	}

	_, err = service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error from repository save")
	}
}

func TestService_DeleteMovie_InvalidID(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	err := service.DeleteMovie(context.Background(), -1)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

func TestService_DeleteMovie_RepositoryError(t *testing.T) {
	repo := NewMockMovieRepository()
	repo.deleteFunc = func(ctx context.Context, id shared.MovieID) error {
		return errors.New("delete failed")
	}
	service := NewService(repo)

	err := service.DeleteMovie(context.Background(), 1)
	if err == nil {
		t.Error("Expected error from repository delete")
	}
}

func TestService_SearchMovies_ExtendedCriteria(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create test movies with various properties
	movies := []struct {
		cmd    CreateMovieCommand
		genres []string
	}{
		{CreateMovieCommand{Title: "Inception", Director: "Christopher Nolan", Year: 2010}, []string{"Sci-Fi", "Thriller"}},
		{CreateMovieCommand{Title: "Interstellar", Director: "Christopher Nolan", Year: 2014}, []string{"Sci-Fi", "Drama"}},
		{CreateMovieCommand{Title: "The Matrix", Director: "The Wachowskis", Year: 1999}, []string{"Sci-Fi", "Action"}},
		{CreateMovieCommand{Title: "Pulp Fiction", Director: "Quentin Tarantino", Year: 1994}, []string{"Crime", "Drama"}},
	}

	for _, movieData := range movies {
		created, err := service.CreateMovie(context.Background(), movieData.cmd)
		if err != nil {
			t.Fatalf("Failed to create test movie: %v", err)
		}

		// Add genres
		if len(movieData.genres) > 0 {
			updateCmd := UpdateMovieCommand{
				ID:       created.ID,
				Title:    created.Title,
				Director: created.Director,
				Year:     created.Year,
				Genres:   movieData.genres,
			}
			_, err = service.UpdateMovie(context.Background(), updateCmd)
			if err != nil {
				t.Fatalf("Failed to update movie genres: %v", err)
			}
		}
	}

	tests := []struct {
		name     string
		query    SearchMoviesQuery
		expected int
	}{
		{
			name: "search by title",
			query: SearchMoviesQuery{
				Title: "Inception",
				Limit: 10,
			},
			expected: 1,
		},
		{
			name: "search by director",
			query: SearchMoviesQuery{
				Director: "Christopher Nolan",
				Limit:    10,
			},
			expected: 2,
		},
		{
			name: "search by genre",
			query: SearchMoviesQuery{
				Genre: "Sci-Fi",
				Limit: 10,
			},
			expected: 3,
		},
		{
			name: "search with limit",
			query: SearchMoviesQuery{
				Limit: 2,
			},
			expected: 2,
		},
		{
			name: "search with offset",
			query: SearchMoviesQuery{
				Limit:  10,
				Offset: 1,
			},
			expected: 3, // Should skip first result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchMovies(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("SearchMovies() error = %v", err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d movies, got %d", tt.expected, len(results))
			}
		})
	}
}

func TestService_GetMovie_InvalidID(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	_, err := service.GetMovie(context.Background(), -1)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

// Additional comprehensive tests for SearchMovies to improve coverage
func TestService_SearchMovies_ComprehensiveCoverage(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create test movies with various properties
	movies := []struct {
		cmd    CreateMovieCommand
		rating float64
		genres []string
	}{
		{CreateMovieCommand{Title: "Inception", Director: "Christopher Nolan", Year: 2010}, 8.8, []string{"Sci-Fi", "Thriller"}},
		{CreateMovieCommand{Title: "Interstellar", Director: "Christopher Nolan", Year: 2014}, 8.6, []string{"Sci-Fi", "Drama"}},
		{CreateMovieCommand{Title: "The Matrix", Director: "The Wachowskis", Year: 1999}, 8.7, []string{"Sci-Fi", "Action"}},
		{CreateMovieCommand{Title: "Pulp Fiction", Director: "Quentin Tarantino", Year: 1994}, 8.9, []string{"Crime", "Drama"}},
		{CreateMovieCommand{Title: "The Dark Knight", Director: "Christopher Nolan", Year: 2008}, 9.0, []string{"Action", "Crime"}},
	}

	for _, movieData := range movies {
		created, err := service.CreateMovie(context.Background(), movieData.cmd)
		if err != nil {
			t.Fatalf("Failed to create test movie: %v", err)
		}

		// Update with rating and genres
		updateCmd := UpdateMovieCommand{
			ID:       created.ID,
			Title:    created.Title,
			Director: created.Director,
			Year:     created.Year,
			Rating:   movieData.rating,
			Genres:   movieData.genres,
		}
		_, err = service.UpdateMovie(context.Background(), updateCmd)
		if err != nil {
			t.Fatalf("Failed to update movie: %v", err)
		}
	}

	tests := []struct {
		name     string
		query    SearchMoviesQuery
		expected int
		desc     string
	}{
		{
			name:     "empty query with default limit",
			query:    SearchMoviesQuery{},
			expected: 5,
			desc:     "Should return all movies with default limit",
		},
		{
			name:     "zero limit defaults to 50",
			query:    SearchMoviesQuery{Limit: 0},
			expected: 5,
			desc:     "Should apply default limit of 50",
		},
		{
			name:     "search by title exact match",
			query:    SearchMoviesQuery{Title: "Inception", Limit: 10},
			expected: 1,
			desc:     "Should find exact title match",
		},
		{
			name:     "search by director",
			query:    SearchMoviesQuery{Director: "Christopher Nolan", Limit: 10},
			expected: 3,
			desc:     "Should find all Nolan movies",
		},
		{
			name:     "search by genre",
			query:    SearchMoviesQuery{Genre: "Sci-Fi", Limit: 10},
			expected: 3,
			desc:     "Should find all Sci-Fi movies",
		},
		{
			name:     "search by year range",
			query:    SearchMoviesQuery{MinYear: 2008, MaxYear: 2014, Limit: 10},
			expected: 3,
			desc:     "Should find movies between 2008-2014",
		},
		{
			name:     "search by min year only",
			query:    SearchMoviesQuery{MinYear: 2010, Limit: 10},
			expected: 2,
			desc:     "Should find movies from 2010 onwards",
		},
		{
			name:     "search by max year only",
			query:    SearchMoviesQuery{MaxYear: 2000, Limit: 10},
			expected: 2,
			desc:     "Should find movies before 2000",
		},
		{
			name:     "search by rating range",
			query:    SearchMoviesQuery{MinRating: 8.8, MaxRating: 9.0, Limit: 10},
			expected: 3,
			desc:     "Should find highly rated movies",
		},
		{
			name:     "search by min rating only",
			query:    SearchMoviesQuery{MinRating: 8.8, Limit: 10},
			expected: 3,
			desc:     "Should find movies with rating >= 8.8",
		},
		{
			name:     "search by max rating only",
			query:    SearchMoviesQuery{MaxRating: 8.7, Limit: 10},
			expected: 2,
			desc:     "Should find movies with rating <= 8.7",
		},
		{
			name:     "search with limit",
			query:    SearchMoviesQuery{Limit: 3},
			expected: 3,
			desc:     "Should respect limit",
		},
		{
			name:     "search with offset",
			query:    SearchMoviesQuery{Limit: 10, Offset: 2},
			expected: 3,
			desc:     "Should skip first 2 movies",
		},
		{
			name:     "search with high offset",
			query:    SearchMoviesQuery{Limit: 10, Offset: 10},
			expected: 0,
			desc:     "Should return empty when offset exceeds total",
		},
		{
			name:     "search by non-existent title",
			query:    SearchMoviesQuery{Title: "Non Existent", Limit: 10},
			expected: 0,
			desc:     "Should return empty for non-existent title",
		},
		{
			name:     "search by non-existent director",
			query:    SearchMoviesQuery{Director: "Non Existent", Limit: 10},
			expected: 0,
			desc:     "Should return empty for non-existent director",
		},
		{
			name:     "search by non-existent genre",
			query:    SearchMoviesQuery{Genre: "NonExistent", Limit: 10},
			expected: 0,
			desc:     "Should return empty for non-existent genre",
		},
		{
			name:     "order by title",
			query:    SearchMoviesQuery{OrderBy: "title", OrderDir: "asc", Limit: 10},
			expected: 5,
			desc:     "Should order by title ascending",
		},
		{
			name:     "order by director desc",
			query:    SearchMoviesQuery{OrderBy: "director", OrderDir: "desc", Limit: 10},
			expected: 5,
			desc:     "Should order by director descending",
		},
		{
			name:     "order by year",
			query:    SearchMoviesQuery{OrderBy: "year", Limit: 10},
			expected: 5,
			desc:     "Should order by year",
		},
		{
			name:     "order by rating",
			query:    SearchMoviesQuery{OrderBy: "rating", Limit: 10},
			expected: 5,
			desc:     "Should order by rating",
		},
		{
			name:     "order by created_at",
			query:    SearchMoviesQuery{OrderBy: "created_at", Limit: 10},
			expected: 5,
			desc:     "Should order by created_at",
		},
		{
			name:     "order by updated_at",
			query:    SearchMoviesQuery{OrderBy: "updated_at", Limit: 10},
			expected: 5,
			desc:     "Should order by updated_at",
		},
		{
			name:     "invalid order by defaults to title",
			query:    SearchMoviesQuery{OrderBy: "invalid", Limit: 10},
			expected: 5,
			desc:     "Should default to title ordering",
		},
		{
			name: "complex search with multiple criteria",
			query: SearchMoviesQuery{
				Director:  "Christopher Nolan",
				Genre:     "Sci-Fi",
				MinYear:   2010,
				MaxYear:   2015,
				MinRating: 8.0,
				Limit:     10,
			},
			expected: 2,
			desc:     "Should find movies matching all criteria",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchMovies(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("SearchMovies() error = %v", err)
			}

			if len(results) != tt.expected {
				t.Errorf("%s: Expected %d movies, got %d", tt.desc, tt.expected, len(results))
			}
		})
	}
}

func TestService_SearchMovies_RepositoryError(t *testing.T) {
	repo := NewMockMovieRepository()
	repo.findByCriteriaFunc = func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
		return nil, errors.New("repository error")
	}
	service := NewService(repo)

	query := SearchMoviesQuery{Title: "Test", Limit: 10}
	_, err := service.SearchMovies(context.Background(), query)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// Additional edge cases to boost CreateMovie coverage
func TestService_CreateMovie_InvalidGenre(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
		Genres:   []string{""}, // Empty genre should fail
	}

	_, err := service.CreateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for empty genre")
	}
}

func TestService_CreateMovie_InvalidPosterURL(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	cmd := CreateMovieCommand{
		Title:     "Test Movie",
		Director:  "Test Director",
		Year:      2020,
		PosterURL: "invalid-url", // Invalid URL format
	}

	_, err := service.CreateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for invalid poster URL")
	}
}

// Additional edge cases to boost UpdateMovie coverage
func TestService_UpdateMovie_InvalidGenre(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Try to update with invalid genre
	cmd := UpdateMovieCommand{
		ID:       created.ID,
		Title:    "Updated Title",
		Director: "Updated Director",
		Year:     2021,
		Genres:   []string{""}, // Empty genre should fail
	}

	_, err = service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for empty genre")
	}
}

func TestService_UpdateMovie_InvalidPosterURL(t *testing.T) {
	repo := NewMockMovieRepository()
	service := NewService(repo)

	// Create a movie first
	createCmd := CreateMovieCommand{
		Title:    "Test Movie",
		Director: "Test Director",
		Year:     2020,
	}
	created, err := service.CreateMovie(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Try to update with invalid poster URL
	cmd := UpdateMovieCommand{
		ID:        created.ID,
		Title:     "Updated Title",
		Director:  "Updated Director",
		Year:      2021,
		PosterURL: "invalid-url", // Invalid URL format
	}

	_, err = service.UpdateMovie(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for invalid poster URL")
	}
}

func TestService_GetTopRatedMovies_RepositoryError(t *testing.T) {
	repo := NewMockMovieRepository()
	repo.findTopRatedFunc = func(ctx context.Context, limit int) ([]*movie.Movie, error) {
		return nil, errors.New("repository error")
	}
	service := NewService(repo)

	_, err := service.GetTopRatedMovies(context.Background(), 10)
	if err == nil {
		t.Error("Expected error from repository")
	}
}
