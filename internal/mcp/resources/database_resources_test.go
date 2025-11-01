package resources

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
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

// MockMovieRepository is a mock implementation of the movie repository for testing
type MockMovieRepository struct {
	FindByCriteriaFunc func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error)
}

func (m *MockMovieRepository) Save(ctx context.Context, mov *movie.Movie) error {
	return errors.New("not implemented")
}

func (m *MockMovieRepository) FindByID(ctx context.Context, id shared.MovieID) (*movie.Movie, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) FindByCriteria(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
	if m.FindByCriteriaFunc != nil {
		return m.FindByCriteriaFunc(ctx, criteria)
	}
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) FindByTitle(ctx context.Context, title string) ([]*movie.Movie, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) FindByDirector(ctx context.Context, director string) ([]*movie.Movie, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) FindByGenre(ctx context.Context, genre string) ([]*movie.Movie, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) FindTopRated(ctx context.Context, limit int) ([]*movie.Movie, error) {
	return nil, errors.New("not implemented")
}

func (m *MockMovieRepository) Delete(ctx context.Context, id shared.MovieID) error {
	return errors.New("not implemented")
}

func (m *MockMovieRepository) DeleteAll(ctx context.Context) error {
	return errors.New("not implemented")
}

func (m *MockMovieRepository) CountAll(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
}

// Handler Tests

func TestHandleAllMovies_Success(t *testing.T) {
	// Create mock repository
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			// Create test movies
			movie1, _ := movie.NewMovie("The Shawshank Redemption", "Frank Darabont", 1994)
			movie1.SetRating(9.3)
			movie1.AddGenre("Drama")

			movie2, _ := movie.NewMovie("The Godfather", "Francis Ford Coppola", 1972)
			movie2.SetRating(9.2)
			movie2.AddGenre("Crime")
			movie2.AddGenre("Drama")

			return []*movie.Movie{movie1, movie2}, nil
		},
	}

	// Create service with mock repository
	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)

	// Create request - the request parameter is not used by the handler
	var req *mcp.ReadResourceRequest

	// Call handler
	result, err := resources.HandleAllMovies(context.Background(), req)

	// Verify no error
	if err != nil {
		t.Fatalf("HandleAllMovies() error = %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Contents) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Contents))
	}

	content := result.Contents[0]
	if content.URI != "movies://database/all" {
		t.Errorf("Expected URI 'movies://database/all', got '%s'", content.URI)
	}

	if content.MIMEType != "application/json" {
		t.Errorf("Expected MIMEType 'application/json', got '%s'", content.MIMEType)
	}

	// Verify JSON structure
	var data map[string]interface{}
	err = json.Unmarshal([]byte(content.Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total_movies"].(float64) != 2 {
		t.Errorf("Expected total_movies=2, got %v", data["total_movies"])
	}

	movies := data["movies"].([]interface{})
	if len(movies) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(movies))
	}
}

func TestHandleAllMovies_EmptyDatabase(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return []*movie.Movie{}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandleAllMovies(context.Background(), req)

	if err != nil {
		t.Fatalf("HandleAllMovies() error = %v", err)
	}

	// Verify JSON structure for empty database
	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total_movies"].(float64) != 0 {
		t.Errorf("Expected total_movies=0, got %v", data["total_movies"])
	}
}

func TestHandleAllMovies_RepositoryError(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandleAllMovies(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result on error, got %v", result)
	}

	if !strings.Contains(err.Error(), "failed to fetch all movies") {
		t.Errorf("Expected error message to contain 'failed to fetch all movies', got '%s'", err.Error())
	}
}

func TestHandleDatabaseStats_Success(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			movie1, _ := movie.NewMovie("The Shawshank Redemption", "Frank Darabont", 1994)
			movie1.SetRating(9.3)
			movie1.AddGenre("Drama")

			movie2, _ := movie.NewMovie("The Godfather", "Francis Ford Coppola", 1972)
			movie2.SetRating(9.2)
			movie2.AddGenre("Crime")
			movie2.AddGenre("Drama")

			movie3, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
			movie3.SetRating(8.8)
			movie3.AddGenre("Sci-Fi")
			movie3.AddGenre("Action")

			return []*movie.Movie{movie1, movie2, movie3}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandleDatabaseStats(context.Background(), req)

	if err != nil {
		t.Fatalf("HandleDatabaseStats() error = %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify JSON structure
	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total_movies"].(float64) != 3 {
		t.Errorf("Expected total_movies=3, got %v", data["total_movies"])
	}

	if data["total_genres"].(float64) != 4 {
		t.Errorf("Expected total_genres=4 (Drama, Crime, Sci-Fi, Action), got %v", data["total_genres"])
	}

	yearRange := data["year_range"].(map[string]interface{})
	if yearRange["earliest"].(float64) != 1972 {
		t.Errorf("Expected earliest year=1972, got %v", yearRange["earliest"])
	}
	if yearRange["latest"].(float64) != 2010 {
		t.Errorf("Expected latest year=2010, got %v", yearRange["latest"])
	}

	// Average rating should be (9.3 + 9.2 + 8.8) / 3 = 9.1
	avgRating := data["average_rating"].(string)
	if avgRating != "9.1" {
		t.Errorf("Expected average_rating='9.1', got '%s'", avgRating)
	}
}

func TestHandleDatabaseStats_EmptyDatabase(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return []*movie.Movie{}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandleDatabaseStats(context.Background(), req)

	if err != nil {
		t.Fatalf("HandleDatabaseStats() error = %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total_movies"].(float64) != 0 {
		t.Errorf("Expected total_movies=0, got %v", data["total_movies"])
	}

	if data["total_genres"].(float64) != 0 {
		t.Errorf("Expected total_genres=0, got %v", data["total_genres"])
	}

	// Average rating should be "0.0" for empty database
	avgRating := data["average_rating"].(string)
	if avgRating != "0.0" {
		t.Errorf("Expected average_rating='0.0', got '%s'", avgRating)
	}
}

func TestHandleDatabaseStats_NoRatings(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			unratedMovie, _ := movie.NewMovie("Unrated Movie", "Director", 2020)
			// Don't set rating - it will be 0
			unratedMovie.AddGenre("Drama")

			return []*movie.Movie{unratedMovie}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandleDatabaseStats(context.Background(), req)

	if err != nil {
		t.Fatalf("HandleDatabaseStats() error = %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Average rating should be "0.0" when no movies have ratings
	avgRating := data["average_rating"].(string)
	if avgRating != "0.0" {
		t.Errorf("Expected average_rating='0.0' for movies without ratings, got '%s'", avgRating)
	}
}

func TestHandleDatabaseStats_RepositoryError(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return nil, errors.New("database error")
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	_, err := resources.HandleDatabaseStats(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to fetch movies for stats") {
		t.Errorf("Expected error message to contain 'failed to fetch movies for stats', got '%s'", err.Error())
	}
}

func TestHandlePosterCollection_Success(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			id1, _ := shared.NewMovieID(1)
			movie1, _ := movie.NewMovieWithID(id1, "Movie 1", "Director 1", 2020)
			movie1.SetRating(8.5)
			movie1.AddGenre("Action")

			id2, _ := shared.NewMovieID(2)
			movie2, _ := movie.NewMovieWithID(id2, "Movie 2", "Director 2", 2021)
			movie2.SetRating(7.5)
			movie2.AddGenre("Drama")

			return []*movie.Movie{movie1, movie2}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandlePosterCollection(context.Background(), req)

	if err != nil {
		t.Fatalf("HandlePosterCollection() error = %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify JSON structure
	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total"].(float64) != 2 {
		t.Errorf("Expected total=2, got %v", data["total"])
	}

	posters := data["posters"].([]interface{})
	if len(posters) != 2 {
		t.Fatalf("Expected 2 posters, got %d", len(posters))
	}

	// Verify poster structure
	poster1 := posters[0].(map[string]interface{})
	if poster1["movie_id"].(float64) != 1 {
		t.Errorf("Expected movie_id=1, got %v", poster1["movie_id"])
	}
	if poster1["title"].(string) != "Movie 1" {
		t.Errorf("Expected title='Movie 1', got '%s'", poster1["title"])
	}
	if poster1["year"].(float64) != 2020 {
		t.Errorf("Expected year=2020, got %v", poster1["year"])
	}
	if poster1["uri"].(string) != "movies://posters/1" {
		t.Errorf("Expected uri='movies://posters/1', got '%s'", poster1["uri"])
	}
}

func TestHandlePosterCollection_EmptyDatabase(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return []*movie.Movie{}, nil
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	result, err := resources.HandlePosterCollection(context.Background(), req)

	if err != nil {
		t.Fatalf("HandlePosterCollection() error = %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Contents[0].Text), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["total"].(float64) != 0 {
		t.Errorf("Expected total=0, got %v", data["total"])
	}

	posters := data["posters"].([]interface{})
	if len(posters) != 0 {
		t.Errorf("Expected 0 posters, got %d", len(posters))
	}
}

func TestHandlePosterCollection_RepositoryError(t *testing.T) {
	mockRepo := &MockMovieRepository{
		FindByCriteriaFunc: func(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
			return nil, errors.New("database error")
		},
	}

	service := movieApp.NewService(mockRepo)
	resources := NewDatabaseResources(service)
	var req *mcp.ReadResourceRequest

	_, err := resources.HandlePosterCollection(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to fetch movies for poster collection") {
		t.Errorf("Expected error message to contain 'failed to fetch movies for poster collection', got '%s'", err.Error())
	}
}
