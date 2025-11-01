package sqlite

import (
	"context"
	"database/sql"
	"testing"

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
