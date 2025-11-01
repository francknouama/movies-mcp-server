package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Add _time_format parameter to parse timestamps
	db, err := sql.Open("sqlite", ":memory:?_time_format=sqlite")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create schema
	schema := `
	CREATE TABLE movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		director TEXT NOT NULL,
		year INTEGER NOT NULL,
		rating REAL,
		genre TEXT NOT NULL DEFAULT '[]',
		description TEXT,
		duration INTEGER,
		language TEXT,
		country TEXT,
		poster_data BLOB,
		poster_type TEXT,
		poster_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return db
}

func TestMovieRepository_Save_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a new movie (no ID)
	domainMovie, err := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	if err != nil {
		t.Fatalf("failed to create domain movie: %v", err)
	}

	// Save the movie
	err = repo.Save(ctx, domainMovie)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify ID was assigned
	if domainMovie.ID().IsZero() {
		t.Error("Expected movie to have ID assigned after save")
	}

	// Verify movie can be retrieved
	retrieved, err := repo.FindByID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if retrieved.Title() != "Inception" {
		t.Errorf("Expected title 'Inception', got %s", retrieved.Title())
	}
	if retrieved.Director() != "Christopher Nolan" {
		t.Errorf("Expected director 'Christopher Nolan', got %s", retrieved.Director())
	}
	if retrieved.Year().Value() != 2010 {
		t.Errorf("Expected year 2010, got %d", retrieved.Year().Value())
	}
}

func TestMovieRepository_Save_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie
	domainMovie, _ := movie.NewMovie("Original Title", "Original Director", 2010)
	_ = repo.Save(ctx, domainMovie)

	// Update the movie's rating (only mutable field we can update directly)
	_ = domainMovie.SetRating(8.5)

	// Save the update
	err := repo.Save(ctx, domainMovie)
	if err != nil {
		t.Fatalf("Save() update error = %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.FindByID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if retrieved.Rating().Value() != 8.5 {
		t.Errorf("Expected rating 8.5, got %f", retrieved.Rating().Value())
	}
}

func TestMovieRepository_Save_WithGenres(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie with genres
	domainMovie, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	_ = domainMovie.AddGenre("Sci-Fi")
	_ = domainMovie.AddGenre("Thriller")
	_ = domainMovie.AddGenre("Action")

	// Save the movie
	err := repo.Save(ctx, domainMovie)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Retrieve and verify genres
	retrieved, err := repo.FindByID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	genres := retrieved.Genres()
	if len(genres) != 3 {
		t.Errorf("Expected 3 genres, got %d", len(genres))
	}

	expectedGenres := map[string]bool{"Sci-Fi": true, "Thriller": true, "Action": true}
	for _, genre := range genres {
		if !expectedGenres[genre] {
			t.Errorf("Unexpected genre: %s", genre)
		}
	}
}

func TestMovieRepository_Save_WithPosterURL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie with poster URL
	domainMovie, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	_ = domainMovie.SetPosterURL("https://example.com/poster.jpg")

	// Save the movie
	err := repo.Save(ctx, domainMovie)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Retrieve and verify poster URL
	retrieved, err := repo.FindByID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if retrieved.PosterURL() != "https://example.com/poster.jpg" {
		t.Errorf("Expected poster URL 'https://example.com/poster.jpg', got %s", retrieved.PosterURL())
	}
}

func TestMovieRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	movieID, _ := shared.NewMovieID(999)
	_, err := repo.FindByID(ctx, movieID)
	if err == nil {
		t.Error("Expected error for non-existent movie")
	}
}

func TestMovieRepository_FindByCriteria_ByTitle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []struct {
		title    string
		director string
		year     int
	}{
		{"Inception", "Christopher Nolan", 2010},
		{"Interstellar", "Christopher Nolan", 2014},
		{"The Matrix", "The Wachowskis", 1999},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by partial title (case-insensitive)
	criteria := movie.SearchCriteria{
		Title: "incep",
		Limit: 10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 movie, got %d", len(results))
	}

	if len(results) > 0 && results[0].Title() != "Inception" {
		t.Errorf("Expected 'Inception', got %s", results[0].Title())
	}
}

func TestMovieRepository_FindByCriteria_ByDirector(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []struct {
		title    string
		director string
		year     int
	}{
		{"Inception", "Christopher Nolan", 2010},
		{"Interstellar", "Christopher Nolan", 2014},
		{"The Dark Knight", "Christopher Nolan", 2008},
		{"The Matrix", "The Wachowskis", 1999},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by director (case-insensitive)
	criteria := movie.SearchCriteria{
		Director: "nolan",
		Limit:    10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 movies by Nolan, got %d", len(results))
	}
}

func TestMovieRepository_FindByCriteria_ByGenre(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies with genres
	movies := []struct {
		title    string
		director string
		year     int
		genres   []string
	}{
		{"Inception", "Christopher Nolan", 2010, []string{"Sci-Fi", "Thriller"}},
		{"Interstellar", "Christopher Nolan", 2014, []string{"Sci-Fi", "Drama"}},
		{"The Matrix", "The Wachowskis", 1999, []string{"Sci-Fi", "Action"}},
		{"Pulp Fiction", "Quentin Tarantino", 1994, []string{"Crime", "Drama"}},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		for _, genre := range m.genres {
			_ = domainMovie.AddGenre(genre)
		}
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by genre
	criteria := movie.SearchCriteria{
		Genre: "Sci-Fi",
		Limit: 10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 Sci-Fi movies, got %d", len(results))
	}
}

func TestMovieRepository_FindByCriteria_ByYearRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []struct {
		title    string
		director string
		year     int
	}{
		{"Movie 1", "Director A", 2005},
		{"Movie 2", "Director B", 2010},
		{"Movie 3", "Director C", 2015},
		{"Movie 4", "Director D", 2020},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by year range
	criteria := movie.SearchCriteria{
		MinYear: 2010,
		MaxYear: 2015,
		Limit:   10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 movies (2010-2015), got %d", len(results))
	}
}

func TestMovieRepository_FindByCriteria_ByRatingRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies with ratings
	movies := []struct {
		title    string
		director string
		year     int
		rating   float64
	}{
		{"Movie A", "Director A", 2020, 7.5},
		{"Movie B", "Director B", 2020, 8.5},
		{"Movie C", "Director C", 2020, 9.0},
		{"Movie D", "Director D", 2020, 9.5},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = domainMovie.SetRating(m.rating)
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by rating range
	criteria := movie.SearchCriteria{
		MinRating: 8.5,
		MaxRating: 9.5,
		Limit:     10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 movies (rating 8.5-9.5), got %d", len(results))
	}
}

func TestMovieRepository_FindByCriteria_WithOrderBy(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []struct {
		title    string
		director string
		year     int
		rating   float64
	}{
		{"Movie C", "Director C", 2015, 8.0},
		{"Movie A", "Director A", 2010, 9.0},
		{"Movie B", "Director B", 2020, 7.5},
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = domainMovie.SetRating(m.rating)
		_ = repo.Save(ctx, domainMovie)
	}

	tests := []struct {
		name        string
		orderBy     movie.OrderBy
		orderDir    movie.OrderDirection
		firstTitle  string
		description string
	}{
		{
			name:        "order by title ASC",
			orderBy:     movie.OrderByTitle,
			orderDir:    movie.OrderAsc,
			firstTitle:  "Movie A",
			description: "Should sort alphabetically",
		},
		{
			name:        "order by title DESC",
			orderBy:     movie.OrderByTitle,
			orderDir:    movie.OrderDesc,
			firstTitle:  "Movie C",
			description: "Should sort reverse alphabetically",
		},
		{
			name:        "order by year ASC",
			orderBy:     movie.OrderByYear,
			orderDir:    movie.OrderAsc,
			firstTitle:  "Movie A",
			description: "Should sort by oldest year",
		},
		{
			name:        "order by rating DESC",
			orderBy:     movie.OrderByRating,
			orderDir:    movie.OrderDesc,
			firstTitle:  "Movie A",
			description: "Should sort by highest rating",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := movie.SearchCriteria{
				OrderBy:  tt.orderBy,
				OrderDir: tt.orderDir,
				Limit:    10,
			}

			results, err := repo.FindByCriteria(ctx, criteria)
			if err != nil {
				t.Fatalf("FindByCriteria() error = %v", err)
			}

			if len(results) == 0 {
				t.Fatal("Expected results, got none")
			}

			if results[0].Title() != tt.firstTitle {
				t.Errorf("%s: Expected first movie '%s', got '%s'", tt.description, tt.firstTitle, results[0].Title())
			}
		})
	}
}

func TestMovieRepository_FindByCriteria_WithLimitAndOffset(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create 10 test movies
	for i := 1; i <= 10; i++ {
		domainMovie, _ := movie.NewMovie("Movie "+string(rune('A'+i-1)), "Director", 2020)
		_ = repo.Save(ctx, domainMovie)
	}

	// Test limit
	criteria := movie.SearchCriteria{
		Limit: 3,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 movies with limit, got %d", len(results))
	}

	// Test offset
	criteria = movie.SearchCriteria{
		Limit:  3,
		Offset: 5,
	}

	results, err = repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 movies with offset, got %d", len(results))
	}
}

func TestMovieRepository_FindByTitle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	domainMovie, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	_ = repo.Save(ctx, domainMovie)

	// Search by title
	results, err := repo.FindByTitle(ctx, "Inception")
	if err != nil {
		t.Fatalf("FindByTitle() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 movie, got %d", len(results))
	}
}

func TestMovieRepository_FindByDirector(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []string{"Inception", "Interstellar", "The Dark Knight"}
	for _, title := range movies {
		domainMovie, _ := movie.NewMovie(title, "Christopher Nolan", 2010)
		_ = repo.Save(ctx, domainMovie)
	}

	// Search by director
	results, err := repo.FindByDirector(ctx, "Christopher Nolan")
	if err != nil {
		t.Fatalf("FindByDirector() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 movies, got %d", len(results))
	}
}

func TestMovieRepository_FindByGenre(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies with Sci-Fi genre
	movies := []string{"Inception", "Interstellar", "The Matrix"}
	for _, title := range movies {
		domainMovie, _ := movie.NewMovie(title, "Director", 2010)
		_ = domainMovie.AddGenre("Sci-Fi")
		_ = repo.Save(ctx, domainMovie)
	}

	// Create a movie without Sci-Fi genre
	otherMovie, _ := movie.NewMovie("Pulp Fiction", "Tarantino", 1994)
	_ = otherMovie.AddGenre("Crime")
	_ = repo.Save(ctx, otherMovie)

	// Search by genre
	results, err := repo.FindByGenre(ctx, "Sci-Fi")
	if err != nil {
		t.Fatalf("FindByGenre() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 Sci-Fi movies, got %d", len(results))
	}
}

func TestMovieRepository_FindTopRated(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies with ratings
	movies := []struct {
		title  string
		rating float64
	}{
		{"Movie A", 9.5},
		{"Movie B", 8.7},
		{"Movie C", 9.8},
		{"Movie D", 8.0},
		{"Movie E", 0.0}, // No rating
	}

	for _, m := range movies {
		domainMovie, _ := movie.NewMovie(m.title, "Director", 2020)
		if m.rating > 0 {
			_ = domainMovie.SetRating(m.rating)
		}
		_ = repo.Save(ctx, domainMovie)
	}

	// Get top 2 rated
	results, err := repo.FindTopRated(ctx, 2)
	if err != nil {
		t.Fatalf("FindTopRated() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 top rated movies, got %d", len(results))
	}

	// First should be highest rated (9.8)
	if len(results) > 0 && results[0].Rating().Value() != 9.8 {
		t.Errorf("Expected first movie rating 9.8, got %f", results[0].Rating().Value())
	}

	// Second should be second highest (9.5)
	if len(results) > 1 && results[1].Rating().Value() != 9.5 {
		t.Errorf("Expected second movie rating 9.5, got %f", results[1].Rating().Value())
	}
}

func TestMovieRepository_CountAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Initially should be 0
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Add 5 movies
	for i := 1; i <= 5; i++ {
		domainMovie, _ := movie.NewMovie("Movie", "Director", 2020)
		_ = repo.Save(ctx, domainMovie)
	}

	// Should now be 5
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestMovieRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie
	domainMovie, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	_ = repo.Save(ctx, domainMovie)

	// Verify it exists
	_, err := repo.FindByID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("Movie should exist before delete")
	}

	// Delete the movie
	err = repo.Delete(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's deleted
	_, err = repo.FindByID(ctx, domainMovie.ID())
	if err == nil {
		t.Error("Expected error when finding deleted movie")
	}
}

func TestMovieRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	movieID, _ := shared.NewMovieID(999)
	err := repo.Delete(ctx, movieID)
	if err == nil {
		t.Error("Expected error when deleting non-existent movie")
	}
}

func TestMovieRepository_DeleteAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Add multiple movies
	for i := 1; i <= 5; i++ {
		domainMovie, _ := movie.NewMovie("Movie", "Director", 2020)
		_ = repo.Save(ctx, domainMovie)
	}

	// Verify count
	count, _ := repo.CountAll(ctx)
	if count != 5 {
		t.Errorf("Expected 5 movies before DeleteAll, got %d", count)
	}

	// Delete all
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify count is 0
	count, _ = repo.CountAll(ctx)
	if count != 0 {
		t.Errorf("Expected 0 movies after DeleteAll, got %d", count)
	}
}

func TestMovieRepository_ComplexSearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create diverse test data
	testMovies := []struct {
		title    string
		director string
		year     int
		rating   float64
		genres   []string
	}{
		{"Inception", "Christopher Nolan", 2010, 8.8, []string{"Sci-Fi", "Thriller"}},
		{"Interstellar", "Christopher Nolan", 2014, 8.6, []string{"Sci-Fi", "Drama"}},
		{"The Dark Knight", "Christopher Nolan", 2008, 9.0, []string{"Action", "Crime"}},
		{"The Matrix", "The Wachowskis", 1999, 8.7, []string{"Sci-Fi", "Action"}},
		{"Pulp Fiction", "Quentin Tarantino", 1994, 8.9, []string{"Crime", "Drama"}},
	}

	for _, m := range testMovies {
		domainMovie, _ := movie.NewMovie(m.title, m.director, m.year)
		_ = domainMovie.SetRating(m.rating)
		for _, genre := range m.genres {
			_ = domainMovie.AddGenre(genre)
		}
		_ = repo.Save(ctx, domainMovie)
	}

	// Complex search: Sci-Fi movies by Nolan after 2009 with rating > 8.5
	criteria := movie.SearchCriteria{
		Director:  "Nolan",
		Genre:     "Sci-Fi",
		MinYear:   2009,
		MinRating: 8.5,
		OrderBy:   movie.OrderByRating,
		OrderDir:  movie.OrderDesc,
		Limit:     10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should find Inception and Interstellar
	if len(results) != 2 {
		t.Errorf("Expected 2 movies matching complex criteria, got %d", len(results))
	}

	// First should be Inception (higher rating)
	if len(results) > 0 && results[0].Title() != "Inception" {
		t.Errorf("Expected 'Inception' first (highest rating), got '%s'", results[0].Title())
	}
}

// Error scenario tests for better coverage

func TestMovieRepository_Save_Insert_InvalidData(t *testing.T) {
	// Create movie with invalid year (should fail domain validation)
	_, err := movie.NewMovie("", "Director", -1)
	if err == nil {
		t.Error("Expected error for invalid movie data")
	}
}

func TestMovieRepository_toDomainModel_InvalidGenresJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)

	// Insert a movie with invalid JSON genres directly into DB
	_, err := db.Exec(`
		INSERT INTO movies (title, director, year, genre, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "Test Movie", "Director", 2020, "[invalid json", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Try to retrieve it - should fail when parsing genres
	ctx := context.Background()
	_, err = repo.FindByTitle(ctx, "Test Movie")
	if err == nil {
		t.Error("Expected error when parsing invalid genres JSON")
	}
}

func TestMovieRepository_toDomainModel_LegacyGenres(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Insert a movie with legacy non-JSON genre format
	_, err := db.Exec(`
		INSERT INTO movies (title, director, year, genre, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "Legacy Movie", "Director", 2020, "Action", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Retrieve it - should handle legacy format gracefully
	result, err := repo.FindByTitle(ctx, "Legacy Movie")
	if err != nil {
		t.Fatalf("FindByTitle() error = %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected to find legacy movie")
	}

	// Should have one genre
	genres := result[0].Genres()
	if len(genres) != 1 || genres[0] != "Action" {
		t.Errorf("Expected genres [Action], got %v", genres)
	}
}

func TestMovieRepository_toDomainModel_InvalidRating(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Insert a movie with invalid rating directly into DB (bypassing domain validation)
	_, err := db.Exec(`
		INSERT INTO movies (title, director, year, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "Bad Rating Movie", "Director", 2020, 15.0, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Try to retrieve it - should fail domain validation for rating
	_, err = repo.FindByTitle(ctx, "Bad Rating Movie")
	if err == nil {
		t.Error("Expected error when domain validation fails for invalid rating")
	}
}

func TestMovieRepository_toDomainModel_InvalidPosterURL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Insert a movie with invalid poster URL format (not a valid URL)
	_, err := db.Exec(`
		INSERT INTO movies (title, director, year, poster_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "Invalid Poster Movie", "Director", 2020, "not-a-valid-url", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Try to retrieve it - should fail domain validation for invalid URL format
	_, err = repo.FindByTitle(ctx, "Invalid Poster Movie")
	if err == nil {
		t.Error("Expected error when domain validation fails for invalid poster URL format")
	}
}

func TestMovieRepository_FindByCriteria_WithInvalidCriteria(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Insert test movie
	testMovie, _ := movie.NewMovie("Test Movie", "Director", 2020)
	testMovie.AddGenre("Action")
	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save test movie: %v", err)
	}

	// Test with genre filter that has no matches
	criteria := movie.SearchCriteria{
		Genre: "NonExistent Genre",
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 movies for non-existent genre, got %d", len(results))
	}
}

func TestMovieRepository_DeleteAll_WithData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Insert multiple test movies
	for i := 1; i <= 3; i++ {
		testMovie, _ := movie.NewMovie(fmt.Sprintf("Movie %d", i), "Director", 2020)
		err := repo.Save(ctx, testMovie)
		if err != nil {
			t.Fatalf("Failed to save test movie: %v", err)
		}
	}

	// Verify movies exist
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 3 {
		t.Fatalf("Expected 3 movies before DeleteAll, got %d", count)
	}

	// Delete all
	err = repo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify all deleted
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 movies after DeleteAll, got %d", count)
	}
}

func TestMovieRepository_DeleteAll_EmptyTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Delete from empty table (should succeed)
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Errorf("DeleteAll() on empty table error = %v", err)
	}

	// Verify count is 0
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0 after DeleteAll on empty table, got %d", count)
	}
}

func TestMovieRepository_FindByCriteria_WithRatingRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with different ratings
	movie1, _ := movie.NewMovie("Low Rated", "Director", 2020)
	movie1.SetRating(5.0)
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("High Rated", "Director", 2020)
	movie2.SetRating(9.0)
	repo.Save(ctx, movie2)

	// Search for high-rated movies
	criteria := movie.SearchCriteria{
		MinRating: 8.0,
		MaxRating: 10.0,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should only find high-rated movie
	if len(results) != 1 {
		t.Errorf("Expected 1 high-rated movie, got %d", len(results))
	}

	if len(results) > 0 && results[0].Title() != "High Rated" {
		t.Errorf("Expected 'High Rated', got '%s'", results[0].Title())
	}
}

func TestMovieRepository_FindByCriteria_WithYearRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies from different years
	movie1, _ := movie.NewMovie("Old Movie", "Director", 1990)
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("New Movie", "Director", 2020)
	repo.Save(ctx, movie2)

	// Search for movies from 2015-2025
	criteria := movie.SearchCriteria{
		MinYear: 2015,
		MaxYear: 2025,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should only find new movie
	if len(results) != 1 {
		t.Errorf("Expected 1 movie in year range, got %d", len(results))
	}

	if len(results) > 0 && results[0].Title() != "New Movie" {
		t.Errorf("Expected 'New Movie', got '%s'", results[0].Title())
	}
}

func TestMovieRepository_FindByCriteria_WithGenre(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with different genres
	movie1, _ := movie.NewMovie("Action Movie", "Director", 2020)
	movie1.AddGenre("Action")
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Drama Movie", "Director", 2020)
	movie2.AddGenre("Drama")
	repo.Save(ctx, movie2)

	// Search for action movies
	criteria := movie.SearchCriteria{
		Genre: "Action",
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should only find action movie
	if len(results) != 1 {
		t.Errorf("Expected 1 action movie, got %d", len(results))
	}

	if len(results) > 0 && results[0].Title() != "Action Movie" {
		t.Errorf("Expected 'Action Movie', got '%s'", results[0].Title())
	}
}

func TestMovieRepository_FindByCriteria_OrderByYear(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies from different years
	movie1, _ := movie.NewMovie("Movie 2020", "Director", 2020)
	movie2, _ := movie.NewMovie("Movie 2015", "Director", 2015)
	movie3, _ := movie.NewMovie("Movie 2018", "Director", 2018)
	repo.Save(ctx, movie1)
	repo.Save(ctx, movie2)
	repo.Save(ctx, movie3)

	// Search ordered by year ascending
	criteria := movie.SearchCriteria{
		OrderBy: movie.OrderByYear,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 movies, got %d", len(results))
	}

	// Verify ordering
	if results[0].Year().Value() != 2015 || results[1].Year().Value() != 2018 || results[2].Year().Value() != 2020 {
		t.Errorf("Movies not ordered by year ascending")
	}
}

func TestMovieRepository_FindByCriteria_OrderDescending(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies
	movie1, _ := movie.NewMovie("A Movie", "Director", 2020)
	movie2, _ := movie.NewMovie("Z Movie", "Director", 2020)
	repo.Save(ctx, movie1)
	repo.Save(ctx, movie2)

	// Search ordered by title descending
	criteria := movie.SearchCriteria{
		OrderBy:  movie.OrderByTitle,
		OrderDir: movie.OrderDesc,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 movies, got %d", len(results))
	}

	// Verify descending order (Z should come before A)
	if results[0].Title() != "Z Movie" {
		t.Errorf("Expected 'Z Movie' first in descending order, got '%s'", results[0].Title())
	}
}

func TestMovieRepository_FindByCriteria_WithDirector(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with different directors
	movie1, _ := movie.NewMovie("Movie 1", "Christopher Nolan", 2020)
	movie2, _ := movie.NewMovie("Movie 2", "Steven Spielberg", 2020)
	repo.Save(ctx, movie1)
	repo.Save(ctx, movie2)

	// Search for Nolan's movies
	criteria := movie.SearchCriteria{
		Director: "Nolan",
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should find Nolan's movie
	if len(results) != 1 {
		t.Errorf("Expected 1 movie by Nolan, got %d", len(results))
	}

	if len(results) > 0 && results[0].Director() != "Christopher Nolan" {
		t.Errorf("Expected Christopher Nolan, got '%s'", results[0].Director())
	}
}

func TestMovieRepository_Save_WithMultipleGenres(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie with multiple genres
	testMovie, _ := movie.NewMovie("Multi-Genre Movie", "Director", 2020)
	testMovie.AddGenre("Action")
	testMovie.AddGenre("Sci-Fi")
	testMovie.AddGenre("Thriller")

	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Retrieve and verify
	results, err := repo.FindByTitle(ctx, "Multi-Genre Movie")
	if err != nil {
		t.Fatalf("FindByTitle() error = %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected to find the movie")
	}

	genres := results[0].Genres()
	if len(genres) != 3 {
		t.Errorf("Expected 3 genres, got %d", len(genres))
	}
}

func TestMovieRepository_Update_WithChangedFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie
	testMovie, _ := movie.NewMovie("Original Title", "Original Director", 2020)
	testMovie.SetRating(7.0)
	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Update the movie
	testMovie.SetRating(9.0)
	testMovie.AddGenre("Updated Genre")
	err = repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Update Save() error = %v", err)
	}

	// Retrieve and verify changes
	results, err := repo.FindByID(ctx, testMovie.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if results.Rating().Value() != 9.0 {
		t.Errorf("Expected rating 9.0, got %v", results.Rating().Value())
	}

	genres := results.Genres()
	if len(genres) == 0 {
		t.Error("Expected genres to be saved")
	}
}

func TestMovieRepository_FindByCriteria_OrderByDirector(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with different directors
	movie1, _ := movie.NewMovie("Movie A", "Zemeckis", 1990)
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Movie B", "Anderson", 1995)
	repo.Save(ctx, movie2)

	movie3, _ := movie.NewMovie("Movie C", "Nolan", 2000)
	repo.Save(ctx, movie3)

	// Search ordered by director ascending
	criteria := movie.SearchCriteria{
		OrderBy:  movie.OrderByDirector,
		OrderDir: movie.OrderAsc,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 movies, got %d", len(results))
	}

	// Verify alphabetical ordering by director: Anderson, Nolan, Zemeckis
	if results[0].Director() != "Anderson" {
		t.Errorf("Expected 'Anderson' first, got '%s'", results[0].Director())
	}
	if results[1].Director() != "Nolan" {
		t.Errorf("Expected 'Nolan' second, got '%s'", results[1].Director())
	}
	if results[2].Director() != "Zemeckis" {
		t.Errorf("Expected 'Zemeckis' third, got '%s'", results[2].Director())
	}
}

func TestMovieRepository_FindByCriteria_OrderByCreatedAt(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with time delays
	movie1, _ := movie.NewMovie("First Movie", "Director A", 1990)
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Second Movie", "Director B", 1995)
	repo.Save(ctx, movie2)

	movie3, _ := movie.NewMovie("Third Movie", "Director C", 2000)
	repo.Save(ctx, movie3)

	// Search ordered by created_at ascending
	criteria := movie.SearchCriteria{
		OrderBy:  movie.OrderByCreatedAt,
		OrderDir: movie.OrderAsc,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 movies, got %d", len(results))
	}

	// Verify ordering by creation time
	if results[0].Title() != "First Movie" {
		t.Errorf("Expected 'First Movie' first, got '%s'", results[0].Title())
	}
	if results[2].Title() != "Third Movie" {
		t.Errorf("Expected 'Third Movie' last, got '%s'", results[2].Title())
	}
}

func TestMovieRepository_FindByCriteria_OrderByUpdatedAt(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies
	movie1, _ := movie.NewMovie("Movie One", "Director A", 1990)
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Movie Two", "Director B", 1995)
	repo.Save(ctx, movie2)

	// Update movie1 to change its updated_at timestamp
	movie1.SetRating(9.5)
	repo.Save(ctx, movie1)

	// Search ordered by updated_at descending (most recently updated first)
	criteria := movie.SearchCriteria{
		OrderBy:  movie.OrderByUpdatedAt,
		OrderDir: movie.OrderDesc,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 movies, got %d", len(results))
	}

	// Movie One should be first because it was updated most recently
	if results[0].Title() != "Movie One" {
		t.Errorf("Expected 'Movie One' first (most recently updated), got '%s'", results[0].Title())
	}
}

func TestMovieRepository_DeleteAll_Comprehensive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create multiple movies with various data
	movie1, _ := movie.NewMovie("Movie 1", "Director 1", 1990)
	movie1.SetRating(8.5)
	movie1.AddGenre("Action")
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Movie 2", "Director 2", 1995)
	movie2.SetRating(7.0)
	movie2.AddGenre("Drama")
	repo.Save(ctx, movie2)

	movie3, _ := movie.NewMovie("Movie 3", "Director 3", 2000)
	repo.Save(ctx, movie3)

	// Verify movies exist
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 movies before DeleteAll, got %d", count)
	}

	// Delete all movies
	err = repo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify all movies are deleted
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() after DeleteAll error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 movies after DeleteAll, got %d", count)
	}

	// Verify we can still query (empty results)
	results, err := repo.FindByCriteria(ctx, movie.SearchCriteria{})
	if err != nil {
		t.Fatalf("FindByCriteria() after DeleteAll error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results after DeleteAll, got %d", len(results))
	}
}

func TestMovieRepository_ComplexSearchScenarios(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with various attributes
	movie1, _ := movie.NewMovie("The Shawshank Redemption", "Frank Darabont", 1994)
	movie1.SetRating(9.3)
	movie1.AddGenre("Drama")
	movie1.SetPosterURL("http://example.com/poster1.jpg")
	repo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("The Godfather", "Francis Ford Coppola", 1972)
	movie2.SetRating(9.2)
	movie2.AddGenre("Crime")
	movie2.AddGenre("Drama")
	repo.Save(ctx, movie2)

	movie3, _ := movie.NewMovie("The Dark Knight", "Christopher Nolan", 2008)
	movie3.SetRating(9.0)
	movie3.AddGenre("Action")
	movie3.AddGenre("Crime")
	repo.Save(ctx, movie3)

	movie4, _ := movie.NewMovie("Pulp Fiction", "Quentin Tarantino", 1994)
	movie4.SetRating(8.9)
	movie4.AddGenre("Crime")
	repo.Save(ctx, movie4)

	// Test 1: Find movies by year and rating range
	criteria1 := movie.SearchCriteria{
		MinYear:   1990,
		MaxYear:   2000,
		MinRating: 9.0,
		OrderBy:   movie.OrderByRating,
		OrderDir:  movie.OrderDesc,
	}

	results1, err := repo.FindByCriteria(ctx, criteria1)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results1) != 1 {
		t.Errorf("Expected 1 movie (Shawshank), got %d", len(results1))
	}

	// Test 2: Find movies with Crime genre, ordered by year
	criteria2 := movie.SearchCriteria{
		Genre:    "Crime",
		OrderBy:  movie.OrderByYear,
		OrderDir: movie.OrderAsc,
	}

	results2, err := repo.FindByCriteria(ctx, criteria2)
	if err != nil {
		t.Fatalf("FindByCriteria() by genre error = %v", err)
	}

	if len(results2) != 3 {
		t.Errorf("Expected 3 movies with Crime genre, got %d", len(results2))
	}

	// Verify ordering by year: Godfather (1972), Pulp Fiction (1994), Dark Knight (2008)
	if results2[0].Year().Value() != 1972 {
		t.Errorf("Expected first movie year 1972, got %d", results2[0].Year().Value())
	}

	// Test 3: Find movies by director
	criteria3 := movie.SearchCriteria{
		Director: "Christopher Nolan",
	}

	results3, err := repo.FindByCriteria(ctx, criteria3)
	if err != nil {
		t.Fatalf("FindByCriteria() by director error = %v", err)
	}

	if len(results3) != 1 {
		t.Errorf("Expected 1 Nolan movie, got %d", len(results3))
	}

	// Test 4: Find movies with limit and offset
	criteria4 := movie.SearchCriteria{
		Limit:    2,
		Offset:   1,
		OrderBy:  movie.OrderByTitle,
		OrderDir: movie.OrderAsc,
	}

	results4, err := repo.FindByCriteria(ctx, criteria4)
	if err != nil {
		t.Fatalf("FindByCriteria() with limit/offset error = %v", err)
	}

	if len(results4) != 2 {
		t.Errorf("Expected 2 movies with limit=2, got %d", len(results4))
	}

	// Test 5: Update movie and verify changes persist
	movie1.SetRating(9.5)
	movie1.AddGenre("Prison")
	err = repo.Save(ctx, movie1)
	if err != nil {
		t.Fatalf("Save() after update error = %v", err)
	}

	updated, err := repo.FindByID(ctx, movie1.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if updated.Rating().Value() != 9.5 {
		t.Errorf("Expected updated rating 9.5, got %v", updated.Rating().Value())
	}

	genres := updated.Genres()
	hasPrison := false
	for _, g := range genres {
		if g == "Prison" {
			hasPrison = true
			break
		}
	}
	if !hasPrison {
		t.Error("Expected 'Prison' genre to be added")
	}

	// Test 6: Find top rated movies
	topRated, err := repo.FindTopRated(ctx, 2)
	if err != nil {
		t.Fatalf("FindTopRated() error = %v", err)
	}

	if len(topRated) != 2 {
		t.Errorf("Expected 2 top rated movies, got %d", len(topRated))
	}

	// First should be the updated Shawshank with 9.5 rating
	if topRated[0].Title() != "The Shawshank Redemption" {
		t.Errorf("Expected 'The Shawshank Redemption' as top rated, got '%s'", topRated[0].Title())
	}
}
