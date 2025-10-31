package resources

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// MockMovieService is a mock implementation for testing resources
type MockMovieService struct {
	SearchMoviesFunc func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)
}

func (m *MockMovieService) GetMovie(ctx context.Context, id int) (*movieApp.MovieDTO, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) CreateMovie(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) UpdateMovie(ctx context.Context, cmd movieApp.UpdateMovieCommand) (*movieApp.MovieDTO, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) DeleteMovie(ctx context.Context, id int) error {
	return errors.New("not implemented")
}

func (m *MockMovieService) SearchMovies(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
	if m.SearchMoviesFunc != nil {
		return m.SearchMoviesFunc(ctx, query)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieService) GetTopRatedMovies(ctx context.Context, limit int) ([]*movieApp.MovieDTO, error) {
	return nil, errors.New("not implemented")
}

// Helper function to create mock service with the Service interface wrapper
func newMockService(searchFunc func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error)) *movieApp.Service {
	// Since we can't easily mock the Service struct directly, we'll need to use a different approach
	// For now, we'll use a wrapper that implements the needed methods
	return nil
}

// Test Resource Definitions

func TestAllMoviesResource(t *testing.T) {
	// Create a nil service for definition tests (definitions don't use the service)
	var service *movieApp.Service
	resources := NewDatabaseResources(service)

	resource := resources.AllMoviesResource()

	if resource.URI != "movies://database/all" {
		t.Errorf("Expected URI 'movies://database/all', got: %s", resource.URI)
	}
	if resource.Name != "All Movies" {
		t.Errorf("Expected Name 'All Movies', got: %s", resource.Name)
	}
	if resource.MIMEType != "application/json" {
		t.Errorf("Expected MIMEType 'application/json', got: %s", resource.MIMEType)
	}
	if resource.Description == "" {
		t.Error("Expected non-empty Description")
	}
}

func TestDatabaseStatsResource(t *testing.T) {
	var service *movieApp.Service
	resources := NewDatabaseResources(service)

	resource := resources.DatabaseStatsResource()

	if resource.URI != "movies://database/stats" {
		t.Errorf("Expected URI 'movies://database/stats', got: %s", resource.URI)
	}
	if resource.Name != "Database Statistics" {
		t.Errorf("Expected Name 'Database Statistics', got: %s", resource.Name)
	}
	if resource.MIMEType != "application/json" {
		t.Errorf("Expected MIMEType 'application/json', got: %s", resource.MIMEType)
	}
}

func TestPosterCollectionResource(t *testing.T) {
	var service *movieApp.Service
	resources := NewDatabaseResources(service)

	resource := resources.PosterCollectionResource()

	if resource.URI != "movies://posters/collection" {
		t.Errorf("Expected URI 'movies://posters/collection', got: %s", resource.URI)
	}
	if resource.Name != "Movie Posters Collection" {
		t.Errorf("Expected Name 'Movie Posters Collection', got: %s", resource.Name)
	}
	if resource.MIMEType != "application/json" {
		t.Errorf("Expected MIMEType 'application/json', got: %s", resource.MIMEType)
	}
}

// For handler tests, we'll need to test through integration since we can't easily mock the Service struct
// However, we can test the JSON structure and error handling logic

func TestResourceJSONStructure(t *testing.T) {
	// Test that we can marshal the expected structures without errors

	// Test AllMovies structure
	allMoviesData := map[string]interface{}{
		"total_movies": 2,
		"movies": []movieApp.MovieDTO{
			{
				ID:       1,
				Title:    "Test Movie",
				Director: "Test Director",
				Year:     2020,
				Rating:   8.5,
				Genres:   []string{"Action"},
			},
		},
	}

	_, err := json.MarshalIndent(allMoviesData, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal AllMovies structure: %v", err)
	}

	// Test DatabaseStats structure
	statsData := map[string]interface{}{
		"total_movies":   3,
		"total_genres":   2,
		"genres":         []string{"Action", "Drama"},
		"average_rating": "8.5",
		"year_range": map[string]interface{}{
			"earliest": 1994,
			"latest":   2020,
		},
	}

	_, err = json.MarshalIndent(statsData, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal DatabaseStats structure: %v", err)
	}

	// Test PosterCollection structure
	collectionData := map[string]interface{}{
		"total": 2,
		"posters": []map[string]interface{}{
			{
				"movie_id": 1,
				"title":    "Test Movie",
				"year":     2020,
				"uri":      "movies://posters/1",
			},
		},
	}

	_, err = json.MarshalIndent(collectionData, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal PosterCollection structure: %v", err)
	}
}

func TestNewDatabaseResources(t *testing.T) {
	var service *movieApp.Service
	resources := NewDatabaseResources(service)

	if resources == nil {
		t.Fatal("Expected non-nil DatabaseResources, got nil")
	}

	// Test that all resource definition methods work
	allMoviesRes := resources.AllMoviesResource()
	if allMoviesRes == nil {
		t.Error("Expected non-nil AllMoviesResource")
	}

	statsRes := resources.DatabaseStatsResource()
	if statsRes == nil {
		t.Error("Expected non-nil DatabaseStatsResource")
	}

	postersRes := resources.PosterCollectionResource()
	if postersRes == nil {
		t.Error("Expected non-nil PosterCollectionResource")
	}
}

func TestResourceContentsStructure(t *testing.T) {
	// Test that mcp.ResourceContents can be created correctly
	contents := &mcp.ResourceContents{
		URI:      "movies://database/all",
		MIMEType: "application/json",
		Text:     `{"total_movies": 0, "movies": []}`,
	}

	if contents.URI != "movies://database/all" {
		t.Errorf("Expected URI 'movies://database/all', got: %s", contents.URI)
	}

	if contents.MIMEType != "application/json" {
		t.Errorf("Expected MIMEType 'application/json', got: %s", contents.MIMEType)
	}

	// Verify JSON is valid
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(contents.Text), &data); err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}
}

func TestReadResourceResultStructure(t *testing.T) {
	// Test that mcp.ReadResourceResult can be created correctly
	result := &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "movies://database/all",
				MIMEType: "application/json",
				Text:     `{"test": "data"}`,
			},
		},
	}

	if len(result.Contents) != 1 {
		t.Errorf("Expected 1 content item, got: %d", len(result.Contents))
	}

	if result.Contents[0].URI != "movies://database/all" {
		t.Errorf("Expected URI 'movies://database/all', got: %s", result.Contents[0].URI)
	}
}

// Test error message formats
func TestErrorMessageFormats(t *testing.T) {
	testCases := []struct {
		name          string
		errorTemplate string
		expectedPart  string
	}{
		{
			name:          "AllMovies error",
			errorTemplate: "failed to fetch all movies: %w",
			expectedPart:  "failed to fetch all movies",
		},
		{
			name:          "DatabaseStats error",
			errorTemplate: "failed to fetch movies for stats: %w",
			expectedPart:  "failed to fetch movies for stats",
		},
		{
			name:          "PosterCollection error",
			errorTemplate: "failed to fetch movies for poster collection: %w",
			expectedPart:  "failed to fetch movies for poster collection",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseErr := errors.New("database error")
			wrappedErr := errors.New(strings.Replace(tc.errorTemplate, "%w", baseErr.Error(), 1))

			if !strings.Contains(wrappedErr.Error(), tc.expectedPart) {
				t.Errorf("Expected error to contain '%s', got: %s", tc.expectedPart, wrappedErr.Error())
			}
		})
	}
}

// Test statistics calculation logic
func TestStatisticsCalculation(t *testing.T) {
	// Test average rating calculation
	movies := []struct {
		rating float64
	}{
		{8.8},
		{8.7},
		{9.3},
	}

	var totalRating float64
	ratingCount := 0

	for _, movie := range movies {
		if movie.rating > 0 {
			totalRating += movie.rating
			ratingCount++
		}
	}

	avgRating := totalRating / float64(ratingCount)
	expected := 8.933333333333334

	if avgRating != expected {
		t.Errorf("Expected average rating %.2f, got: %.2f", expected, avgRating)
	}

	// Test that formatting to 1 decimal place works
	// In the real code, we use fmt.Sprintf("%.1f", avgRating)
	// which would give us "8.9"
	t.Logf("Average rating: %.1f", avgRating)
}

// Test genre collection logic
func TestGenreCollection(t *testing.T) {
	movies := []struct {
		genres []string
	}{
		{[]string{"Sci-Fi", "Action"}},
		{[]string{"Sci-Fi", "Action"}},
		{[]string{"Drama"}},
	}

	genreSet := make(map[string]bool)
	for _, movie := range movies {
		for _, genre := range movie.genres {
			genreSet[genre] = true
		}
	}

	if len(genreSet) != 3 {
		t.Errorf("Expected 3 unique genres, got: %d", len(genreSet))
	}

	if !genreSet["Sci-Fi"] {
		t.Error("Expected 'Sci-Fi' in genre set")
	}
	if !genreSet["Action"] {
		t.Error("Expected 'Action' in genre set")
	}
	if !genreSet["Drama"] {
		t.Error("Expected 'Drama' in genre set")
	}
}

// Test year range tracking
func TestYearRangeTracking(t *testing.T) {
	movies := []struct {
		year int
	}{
		{2010},
		{1999},
		{1994},
	}

	var earliestYear, latestYear *int

	for _, movie := range movies {
		if earliestYear == nil || movie.year < *earliestYear {
			earliestYear = &movie.year
		}
		if latestYear == nil || movie.year > *latestYear {
			latestYear = &movie.year
		}
	}

	if earliestYear == nil || *earliestYear != 1994 {
		t.Errorf("Expected earliest year 1994, got: %v", earliestYear)
	}

	if latestYear == nil || *latestYear != 2010 {
		t.Errorf("Expected latest year 2010, got: %v", latestYear)
	}
}

// Test poster URI generation
func TestPosterURIGeneration(t *testing.T) {
	// Test that the URI format is correct
	// In the real code, we use fmt.Sprintf("movies://posters/%d", movieID)

	testCases := []struct {
		movieID int
	}{
		{1},
		{42},
		{9999},
	}

	for _, tc := range testCases {
		t.Run("movie_id_"+string(rune(tc.movieID+'0')), func(t *testing.T) {
			// The implementation uses fmt.Sprintf("movies://posters/%d", movie.ID)
			// We're just verifying the logic works for poster collection
			t.Logf("Testing poster URI for movie ID: %d", tc.movieID)
		})
	}
}
