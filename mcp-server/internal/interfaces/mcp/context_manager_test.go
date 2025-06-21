package mcp

import (
	"context"
	"testing"
	"time"

	movieApp "github.com/francknouama/movies-mcp-server/mcp-server/internal/application/movie"
)

// MockMovieService provides a mock implementation for testing
type MockMovieService struct {
	movies []*movieApp.MovieDTO
}

// Implement other required methods (not used in context manager tests)
func (m *MockMovieService) CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
	return nil, nil
}

func (m *MockMovieService) GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
	return nil, nil
}

func (m *MockMovieService) UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
	return nil, nil
}

func (m *MockMovieService) DeleteMovie(ctx context.Context, id int) error {
	return nil
}

func (m *MockMovieService) SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
	// Simple filtering logic for testing
	result := []*movieApp.MovieDTO{}

	for _, movie := range m.movies {
		matches := true

		if query.Title != "" && movie.Title != query.Title {
			matches = false
		}
		if query.Director != "" && movie.Director != query.Director {
			matches = false
		}
		if query.Genre != "" {
			genreMatch := false
			for _, genre := range movie.Genres {
				if genre == query.Genre {
					genreMatch = true
					break
				}
			}
			if !genreMatch {
				matches = false
			}
		}
		if query.MinYear > 0 && movie.Year < query.MinYear {
			matches = false
		}
		if query.MaxYear > 0 && movie.Year > query.MaxYear {
			matches = false
		}

		if matches {
			result = append(result, movie)
		}
	}

	// Apply limit
	if query.Limit > 0 && len(result) > query.Limit {
		result = result[:query.Limit]
	}

	return result, nil
}

func createTestMovies() []*movieApp.MovieDTO {
	return []*movieApp.MovieDTO{
		{ID: 1, Title: "The Godfather", Director: "Francis Ford Coppola", Year: 1972, Rating: 9.2, Genres: []string{"Crime", "Drama"}},
		{ID: 2, Title: "The Shawshank Redemption", Director: "Frank Darabont", Year: 1994, Rating: 9.3, Genres: []string{"Drama"}},
		{ID: 3, Title: "The Dark Knight", Director: "Christopher Nolan", Year: 2008, Rating: 9.0, Genres: []string{"Action", "Crime", "Drama"}},
		{ID: 4, Title: "Pulp Fiction", Director: "Quentin Tarantino", Year: 1994, Rating: 8.9, Genres: []string{"Crime", "Drama"}},
		{ID: 5, Title: "Schindler's List", Director: "Steven Spielberg", Year: 1993, Rating: 8.9, Genres: []string{"Biography", "Drama", "History"}},
		{ID: 6, Title: "The Lord of the Rings: The Return of the King", Director: "Peter Jackson", Year: 2003, Rating: 8.9, Genres: []string{"Adventure", "Drama", "Fantasy"}},
		{ID: 7, Title: "12 Angry Men", Director: "Sidney Lumet", Year: 1957, Rating: 8.9, Genres: []string{"Crime", "Drama"}},
		{ID: 8, Title: "The Good, the Bad and the Ugly", Director: "Sergio Leone", Year: 1966, Rating: 8.8, Genres: []string{"Western"}},
		{ID: 9, Title: "Fight Club", Director: "David Fincher", Year: 1999, Rating: 8.8, Genres: []string{"Drama"}},
		{ID: 10, Title: "Forrest Gump", Director: "Robert Zemeckis", Year: 1994, Rating: 8.8, Genres: []string{"Drama", "Romance"}},
	}
}

func TestContextManager_CreateAndRetrieveContext(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)

	// Create a context
	query := movieApp.SearchMoviesQuery{
		Title: "The Godfather",
	}

	contextInfo, err := cm.CreateMovieSearchContext(context.Background(), query, 5)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	if contextInfo.Total != 1 {
		t.Errorf("Expected total 1, got %d", contextInfo.Total)
	}

	if contextInfo.PageSize != 5 {
		t.Errorf("Expected page size 5, got %d", contextInfo.PageSize)
	}

	if contextInfo.TotalPages != 1 {
		t.Errorf("Expected total pages 1, got %d", contextInfo.TotalPages)
	}

	// Retrieve the first page
	pageReq := PageRequest{
		ContextID: contextInfo.ID,
		Page:      1,
	}

	pageResp, err := cm.GetPage(pageReq)
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}

	if len(pageResp.Data) != 1 {
		t.Errorf("Expected 1 movie in page, got %d", len(pageResp.Data))
	}

	if pageResp.HasNext {
		t.Error("Expected HasNext to be false")
	}

	if pageResp.HasPrevious {
		t.Error("Expected HasPrevious to be false")
	}
}

func TestContextManager_Pagination(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)

	// Create a context with all movies (pageSize = 3)
	query := movieApp.SearchMoviesQuery{} // Empty query returns all movies

	contextInfo, err := cm.CreateMovieSearchContext(context.Background(), query, 3)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	if contextInfo.Total != 10 {
		t.Errorf("Expected total 10, got %d", contextInfo.Total)
	}

	expectedPages := 4 // 10 movies / 3 per page = 4 pages (rounded up)
	if contextInfo.TotalPages != expectedPages {
		t.Errorf("Expected total pages %d, got %d", expectedPages, contextInfo.TotalPages)
	}

	// Test first page
	pageResp, err := cm.GetPage(PageRequest{ContextID: contextInfo.ID, Page: 1})
	if err != nil {
		t.Fatalf("Failed to get page 1: %v", err)
	}

	if len(pageResp.Data) != 3 {
		t.Errorf("Expected 3 movies in page 1, got %d", len(pageResp.Data))
	}

	if !pageResp.HasNext {
		t.Error("Expected HasNext to be true for page 1")
	}

	if pageResp.HasPrevious {
		t.Error("Expected HasPrevious to be false for page 1")
	}

	// Test middle page
	pageResp, err = cm.GetPage(PageRequest{ContextID: contextInfo.ID, Page: 2})
	if err != nil {
		t.Fatalf("Failed to get page 2: %v", err)
	}

	if len(pageResp.Data) != 3 {
		t.Errorf("Expected 3 movies in page 2, got %d", len(pageResp.Data))
	}

	if !pageResp.HasNext {
		t.Error("Expected HasNext to be true for page 2")
	}

	if !pageResp.HasPrevious {
		t.Error("Expected HasPrevious to be true for page 2")
	}

	// Test last page
	pageResp, err = cm.GetPage(PageRequest{ContextID: contextInfo.ID, Page: 4})
	if err != nil {
		t.Fatalf("Failed to get page 4: %v", err)
	}

	if len(pageResp.Data) != 1 { // Only 1 movie on last page (10 % 3 = 1)
		t.Errorf("Expected 1 movie in page 4, got %d", len(pageResp.Data))
	}

	if pageResp.HasNext {
		t.Error("Expected HasNext to be false for last page")
	}

	if !pageResp.HasPrevious {
		t.Error("Expected HasPrevious to be true for last page")
	}
}

func TestContextManager_ContextExpiration(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)
	cm.ttl = 1 * time.Millisecond // Very short TTL for testing

	// Create a context
	query := movieApp.SearchMoviesQuery{}
	contextInfo, err := cm.CreateMovieSearchContext(context.Background(), query, 5)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Try to get page - should fail due to expiration
	_, err = cm.GetPage(PageRequest{ContextID: contextInfo.ID, Page: 1})
	if err == nil {
		t.Error("Expected error due to context expiration, got nil")
	}

	if err.Error() != "context expired: "+contextInfo.ID {
		t.Errorf("Expected context expired error, got: %v", err)
	}
}

func TestContextManager_FilteredSearch(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)

	// Search for movies with "Drama" genre
	query := movieApp.SearchMoviesQuery{
		Genre: "Drama",
	}

	contextInfo, err := cm.CreateMovieSearchContext(context.Background(), query, 5)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// Should find 9 drama movies (count from createTestMovies function)
	// The Godfather, The Shawshank Redemption, The Dark Knight, Pulp Fiction,
	// Schindler's List, The Lord of the Rings, 12 Angry Men, Fight Club, Forrest Gump
	expectedDramaCount := 9
	if contextInfo.Total != expectedDramaCount {
		t.Errorf("Expected %d drama movies, got %d", expectedDramaCount, contextInfo.Total)
	}

	// Get first page
	pageResp, err := cm.GetPage(PageRequest{ContextID: contextInfo.ID, Page: 1})
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}

	if len(pageResp.Data) != 5 {
		t.Errorf("Expected 5 movies in page, got %d", len(pageResp.Data))
	}

	// Verify all returned movies have Drama genre
	for _, item := range pageResp.Data {
		movie := item.(*movieApp.MovieDTO)
		hasDrama := false
		for _, genre := range movie.Genres {
			if genre == "Drama" {
				hasDrama = true
				break
			}
		}
		if !hasDrama {
			t.Errorf("Movie %s should have Drama genre", movie.Title)
		}
	}
}

func TestContextManager_GetContextInfo(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)

	// Create a context
	query := movieApp.SearchMoviesQuery{}
	contextInfo, err := cm.CreateMovieSearchContext(context.Background(), query, 3)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// Get context info
	retrievedInfo, err := cm.GetContextInfo(contextInfo.ID)
	if err != nil {
		t.Fatalf("Failed to get context info: %v", err)
	}

	if retrievedInfo.ID != contextInfo.ID {
		t.Errorf("Expected ID %s, got %s", contextInfo.ID, retrievedInfo.ID)
	}

	if retrievedInfo.Total != contextInfo.Total {
		t.Errorf("Expected total %d, got %d", contextInfo.Total, retrievedInfo.Total)
	}

	if retrievedInfo.PageSize != contextInfo.PageSize {
		t.Errorf("Expected page size %d, got %d", contextInfo.PageSize, retrievedInfo.PageSize)
	}
}

func TestContextManager_NonExistentContext(t *testing.T) {
	mockService := &MockMovieService{movies: createTestMovies()}
	cm := NewContextManager(mockService)

	// Try to get page from non-existent context
	_, err := cm.GetPage(PageRequest{ContextID: "nonexistent", Page: 1})
	if err == nil {
		t.Error("Expected error for non-existent context, got nil")
	}

	// Try to get info from non-existent context
	_, err = cm.GetContextInfo("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent context, got nil")
	}
}
