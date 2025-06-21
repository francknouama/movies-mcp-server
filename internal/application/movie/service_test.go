package movie

import (
	"context"
	"errors"
	"testing"

	"movies-mcp-server/internal/domain/movie"
	"movies-mcp-server/internal/domain/shared"
)

// MockMovieRepository implements movie.Repository for testing
type MockMovieRepository struct {
	movies       map[int]*movie.Movie
	nextID       int
	findByIDFunc func(ctx context.Context, id shared.MovieID) (*movie.Movie, error)
	saveFunc     func(ctx context.Context, m *movie.Movie) error
	deleteFunc   func(ctx context.Context, id shared.MovieID) error
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
	
	// Always assign a new ID for new movies (ID 1 is the temporary ID from NewMovie)
	if movie.ID().Value() == 1 || movie.ID().IsZero() {
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
		
		if match {
			result = append(result, movieItem)
		}
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
	var result []*movie.Movie
	count := 0
	for _, movie := range m.movies {
		if count >= limit {
			break
		}
		if !movie.Rating().IsZero() {
			result = append(result, movie)
			count++
		}
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