package database

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	"movies-mcp-server/internal/config"
)

// Integration tests for PostgreSQL database
// These tests require a running PostgreSQL instance

func TestPostgresSearchMovies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Check if database is available
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	port := 5432
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Name == "" {
		cfg.Name = "movies_mcp"
	}

	db, err := NewPostgresDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Clean up test data
	defer func() {
		db.db.Exec("DELETE FROM movies WHERE title LIKE 'Test Movie%'")
	}()

	// Create test movies
	testMovies := []*Movie{
		{
			Title:       "Test Movie 1",
			Director:    "Test Director",
			Year:        2020,
			Genre:       []string{"Action", "Drama"},
			Rating:      sql.NullFloat64{Float64: 8.5, Valid: true},
			Description: sql.NullString{String: "A test movie about testing", Valid: true},
		},
		{
			Title:       "Test Movie 2",
			Director:    "Another Director",
			Year:        2021,
			Genre:       []string{"Comedy"},
			Rating:      sql.NullFloat64{Float64: 7.0, Valid: true},
			Description: sql.NullString{String: "Another test movie", Valid: true},
		},
		{
			Title:    "Different Title",
			Director: "Test Director",
			Year:     2020,
			Genre:    []string{"Drama"},
			Rating:   sql.NullFloat64{Float64: 9.0, Valid: true},
		},
	}

	for _, movie := range testMovies {
		if err := db.CreateMovie(movie); err != nil {
			t.Fatalf("Failed to create test movie: %v", err)
		}
	}

	// Test search by title
	t.Run("SearchByTitle", func(t *testing.T) {
		results, err := db.SearchMovies(SearchQuery{
			Query: "Test Movie",
			Type:  "title",
			Limit: 10,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})

	// Test search by director
	t.Run("SearchByDirector", func(t *testing.T) {
		results, err := db.SearchMovies(SearchQuery{
			Query: "Test Director",
			Type:  "director",
			Limit: 10,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})

	// Test search by genre
	t.Run("SearchByGenre", func(t *testing.T) {
		results, err := db.SearchMovies(SearchQuery{
			Query: "Drama",
			Type:  "genre",
			Limit: 10,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(results) < 2 {
			t.Errorf("Expected at least 2 results, got %d", len(results))
		}
	})

	// Test search by year
	t.Run("SearchByYear", func(t *testing.T) {
		results, err := db.SearchMovies(SearchQuery{
			Query: "2020",
			Type:  "year",
			Limit: 10,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(results) < 2 {
			t.Errorf("Expected at least 2 results, got %d", len(results))
		}
	})

	// Test pagination
	t.Run("SearchWithPagination", func(t *testing.T) {
		// First page
		page1, err := db.SearchMovies(SearchQuery{
			Query:  "Test",
			Type:   "title",
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(page1) != 1 {
			t.Errorf("Expected 1 result, got %d", len(page1))
		}

		// Second page
		page2, err := db.SearchMovies(SearchQuery{
			Query:  "Test",
			Type:   "title",
			Limit:  1,
			Offset: 1,
		})
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
		if len(page2) > 1 {
			t.Errorf("Expected at most 1 result, got %d", len(page2))
		}

		// Pages should have different movies
		if len(page1) > 0 && len(page2) > 0 && page1[0].ID == page2[0].ID {
			t.Errorf("Pagination not working: same movie on different pages")
		}
	})
}

func TestPostgresListTopMovies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	port := 5432
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Name == "" {
		cfg.Name = "movies_mcp"
	}

	db, err := NewPostgresDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test list top movies
	t.Run("ListTopMovies", func(t *testing.T) {
		results, err := db.ListTopMovies(5, "")
		if err != nil {
			t.Errorf("ListTopMovies failed: %v", err)
		}

		// Verify sorting
		for i := 1; i < len(results); i++ {
			ratingCurr := results[i].Rating.Float64
			if !results[i].Rating.Valid {
				ratingCurr = 0.0
			}
			ratingPrev := results[i-1].Rating.Float64
			if !results[i-1].Rating.Valid {
				ratingPrev = 0.0
			}
			if ratingCurr > ratingPrev {
				t.Errorf("Movies not sorted by rating: %f > %f", 
					ratingCurr, ratingPrev)
			}
		}
	})

	// Test list top movies with genre filter
	t.Run("ListTopMoviesByGenre", func(t *testing.T) {
		results, err := db.ListTopMovies(5, "Drama")
		if err != nil {
			t.Errorf("ListTopMovies failed: %v", err)
		}

		// Verify all movies have the genre
		for _, movie := range results {
			hasGenre := false
			for _, g := range movie.Genre {
				if g == "Drama" {
					hasGenre = true
					break
				}
			}
			if !hasGenre {
				t.Errorf("Movie %s doesn't have Drama genre", movie.Title)
			}
		}
	})
}

// Test full-text search functionality
func TestPostgresFullTextSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	port := 5432
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Name == "" {
		cfg.Name = "movies_mcp"
	}

	db, err := NewPostgresDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Clean up test data
	defer func() {
		db.db.Exec("DELETE FROM movies WHERE title = 'Full Text Test Movie'")
	}()

	// Create a test movie with specific description
	testMovie := &Movie{
		Title:       "Full Text Test Movie",
		Director:    "FTS Director",
		Year:        2023,
		Genre:       []string{"Test"},
		Rating:      sql.NullFloat64{Float64: 5.0, Valid: true},
		Description: sql.NullString{String: "A movie about dreams within dreams and inception of ideas", Valid: true},
	}

	if err := db.CreateMovie(testMovie); err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	// Test full-text search
	t.Run("FullTextSearch", func(t *testing.T) {
		results, err := db.SearchMovies(SearchQuery{
			Query: "dreams inception",
			Type:  "fulltext",
			Limit: 10,
		})
		if err != nil {
			t.Errorf("Full-text search failed: %v", err)
		}

		found := false
		for _, movie := range results {
			if movie.Title == "Full Text Test Movie" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Full-text search didn't find the test movie")
		}
	})
}

// Test that indexes are being used effectively
func TestPostgresIndexUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	port := 5432
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Name == "" {
		cfg.Name = "movies_mcp"
	}

	db, err := NewPostgresDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check if indexes exist
	indexTests := []struct {
		indexName string
		query     string
	}{
		{"idx_movies_title_lower", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_title_lower'"},
		{"idx_movies_director_lower", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_director_lower'"},
		{"idx_movies_year", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_year'"},
		{"idx_movies_rating", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_rating'"},
		{"idx_movies_genre", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_genre'"},
		{"idx_movies_fulltext", "SELECT 1 FROM pg_indexes WHERE indexname = 'idx_movies_fulltext'"},
	}

	for _, test := range indexTests {
		t.Run("Index_"+test.indexName, func(t *testing.T) {
			var exists int
			err := db.db.QueryRow(test.query).Scan(&exists)
			if err == sql.ErrNoRows {
				t.Errorf("Index %s does not exist", test.indexName)
			} else if err != nil {
				t.Errorf("Error checking index %s: %v", test.indexName, err)
			}
		})
	}
}