package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/movie"
)

// Integration tests for MovieRepository
// These tests require a PostgreSQL database connection

func setupTestDB(t *testing.T) *sql.DB {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration tests")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	// Clean up test data
	_, err := db.Exec("DELETE FROM movie_actors")
	if err != nil {
		t.Logf("Warning: failed to clean up movie_actors: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM movies")
	if err != nil {
		t.Logf("Warning: failed to clean up movies: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM actors")
	if err != nil {
		t.Logf("Warning: failed to clean up actors: %v", err)
	}

	db.Close()
}

func createTestMovie(t *testing.T) *movie.Movie {
	movie, err := movie.NewMovie("Test Movie", "Test Director", 2023)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}
	
	movie.SetRating(8.5)
	movie.AddGenre("Action")
	movie.AddGenre("Thriller")
	movie.SetPosterURL("https://example.com/poster.jpg")
	
	return movie
}

func TestMovieRepository_Integration_Save_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	testMovie := createTestMovie(t)

	ctx := context.Background()
	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save movie: %v", err)
	}

	// Verify the movie was assigned an ID
	if testMovie.ID().IsZero() {
		t.Error("Expected movie to be assigned an ID after save")
	}

	// Verify we can retrieve the movie
	retrieved, err := repo.FindByID(ctx, testMovie.ID())
	if err != nil {
		t.Fatalf("Failed to retrieve saved movie: %v", err)
	}

	// Verify movie data
	if retrieved.Title() != testMovie.Title() {
		t.Errorf("Expected title %s, got %s", testMovie.Title(), retrieved.Title())
	}
	if retrieved.Director() != testMovie.Director() {
		t.Errorf("Expected director %s, got %s", testMovie.Director(), retrieved.Director())
	}
	if retrieved.Year().Value() != testMovie.Year().Value() {
		t.Errorf("Expected year %d, got %d", testMovie.Year().Value(), retrieved.Year().Value())
	}
	if retrieved.Rating().Value() != testMovie.Rating().Value() {
		t.Errorf("Expected rating %f, got %f", testMovie.Rating().Value(), retrieved.Rating().Value())
	}
	if retrieved.PosterURL() != testMovie.PosterURL() {
		t.Errorf("Expected poster URL %s, got %s", testMovie.PosterURL(), retrieved.PosterURL())
	}

	// Verify genres
	expectedGenres := map[string]bool{"Action": true, "Thriller": true}
	retrievedGenres := retrieved.Genres()
	if len(retrievedGenres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(retrievedGenres))
	}
	for _, genre := range retrievedGenres {
		if !expectedGenres[genre] {
			t.Errorf("Unexpected genre: %s", genre)
		}
	}
}

func TestMovieRepository_Integration_Save_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	testMovie := createTestMovie(t)

	ctx := context.Background()
	
	// Save initially
	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save movie: %v", err)
	}

	originalID := testMovie.ID()

	// Update the movie
	updatedMovie, err := movie.NewMovieWithID(originalID, "Updated Title", "Updated Director", 2024)
	if err != nil {
		t.Fatalf("Failed to create updated movie: %v", err)
	}
	updatedMovie.SetRating(9.0)
	updatedMovie.AddGenre("Drama")

	// Save the update
	err = repo.Save(ctx, updatedMovie)
	if err != nil {
		t.Fatalf("Failed to update movie: %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.FindByID(ctx, originalID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated movie: %v", err)
	}

	if retrieved.Title() != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got %s", retrieved.Title())
	}
	if retrieved.Director() != "Updated Director" {
		t.Errorf("Expected updated director 'Updated Director', got %s", retrieved.Director())
	}
	if retrieved.Rating().Value() != 9.0 {
		t.Errorf("Expected updated rating 9.0, got %f", retrieved.Rating().Value())
	}
}

func TestMovieRepository_Integration_FindByCriteria(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create test movies
	movies := []*movie.Movie{
		func() *movie.Movie {
			m, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
			m.SetRating(8.8)
			m.AddGenre("Sci-Fi")
			return m
		}(),
		func() *movie.Movie {
			m, _ := movie.NewMovie("Interstellar", "Christopher Nolan", 2014)
			m.SetRating(8.6)
			m.AddGenre("Sci-Fi")
			return m
		}(),
		func() *movie.Movie {
			m, _ := movie.NewMovie("The Matrix", "The Wachowskis", 1999)
			m.SetRating(8.7)
			m.AddGenre("Action")
			return m
		}(),
	}

	// Save all movies
	for _, movie := range movies {
		if err := repo.Save(ctx, movie); err != nil {
			t.Fatalf("Failed to save movie %s: %v", movie.Title(), err)
		}
	}

	// Test search by director
	criteria := movie.SearchCriteria{
		Director: "Christopher Nolan",
		Limit:    10,
	}
	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search by director: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 movies by Christopher Nolan, got %d", len(results))
	}

	// Test search by genre
	criteria = movie.SearchCriteria{
		Genre: "Sci-Fi",
		Limit: 10,
	}
	results, err = repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search by genre: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 Sci-Fi movies, got %d", len(results))
	}

	// Test search by year range
	criteria = movie.SearchCriteria{
		MinYear: 2000,
		MaxYear: 2015,
		Limit:   10,
	}
	results, err = repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search by year range: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 movies from 2000-2015, got %d", len(results))
	}
}

func TestMovieRepository_Integration_FindTopRated(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies with different ratings
	movies := []*movie.Movie{
		func() *movie.Movie {
			m, _ := movie.NewMovie("Great Movie", "Director 1", 2020)
			m.SetRating(9.5)
			return m
		}(),
		func() *movie.Movie {
			m, _ := movie.NewMovie("Good Movie", "Director 2", 2021)
			m.SetRating(8.0)
			return m
		}(),
		func() *movie.Movie {
			m, _ := movie.NewMovie("Average Movie", "Director 3", 2022)
			m.SetRating(7.0)
			return m
		}(),
		func() *movie.Movie {
			m, _ := movie.NewMovie("Unrated Movie", "Director 4", 2023)
			// No rating set
			return m
		}(),
	}

	// Save all movies
	for _, movie := range movies {
		if err := repo.Save(ctx, movie); err != nil {
			t.Fatalf("Failed to save movie: %v", err)
		}
	}

	// Get top 2 rated movies
	results, err := repo.FindTopRated(ctx, 2)
	if err != nil {
		t.Fatalf("Failed to get top rated movies: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 top rated movies, got %d", len(results))
	}

	// Verify they're in descending order by rating
	if results[0].Rating().Value() < results[1].Rating().Value() {
		t.Error("Top rated movies not in descending order")
	}

	// Verify the highest rated movie is first
	if results[0].Title() != "Great Movie" {
		t.Errorf("Expected 'Great Movie' to be top rated, got %s", results[0].Title())
	}
}

func TestMovieRepository_Integration_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	testMovie := createTestMovie(t)

	ctx := context.Background()
	
	// Save movie
	err := repo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save movie: %v", err)
	}

	movieID := testMovie.ID()

	// Delete movie
	err = repo.Delete(ctx, movieID)
	if err != nil {
		t.Fatalf("Failed to delete movie: %v", err)
	}

	// Verify movie is gone
	_, err = repo.FindByID(ctx, movieID)
	if err == nil {
		t.Error("Expected error when finding deleted movie")
	}
}

func TestMovieRepository_Integration_CountAll(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewMovieRepository(db)
	ctx := context.Background()

	// Check initial count
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("Failed to count movies: %v", err)
	}
	initialCount := count

	// Add some movies
	for i := 0; i < 3; i++ {
		movie, _ := movie.NewMovie("Movie", "Director", 2020+i)
		if err := repo.Save(ctx, movie); err != nil {
			t.Fatalf("Failed to save movie: %v", err)
		}
	}

	// Check count again
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("Failed to count movies: %v", err)
	}

	expectedCount := initialCount + 3
	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}
}