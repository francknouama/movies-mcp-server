package main

import (
	"database/sql"
	"fmt"
	"log"

	"movies-mcp-server/internal/config"
	"movies-mcp-server/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Test direct connection
	fmt.Println("Testing database connection...")
	fmt.Printf("Connecting to: %s\n", cfg.Database.ConnectionString())

	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test ping
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection successful!")

	fmt.Println("✓ Note: Run migrations separately using 'make db-migrate'")

	// Test database interface
	fmt.Println("\nTesting database interface...")
	pgDB, err := database.NewPostgresDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database instance: %v", err)
	}
	defer pgDB.Close()

	// Test creating a movie
	movie := &database.Movie{
		Title:       "Test Movie",
		Director:    "Test Director",
		Year:        2024,
		Genre:       []string{"Test"},
		Rating:      sql.NullFloat64{Float64: 8.5, Valid: true},
		Description: sql.NullString{String: "A test movie", Valid: true},
		Duration:    sql.NullInt32{Int32: 120, Valid: true},
		Language:    sql.NullString{String: "English", Valid: true},
		Country:     sql.NullString{String: "USA", Valid: true},
	}

	if err := pgDB.CreateMovie(movie); err != nil {
		log.Fatalf("Failed to create movie: %v", err)
	}
	fmt.Printf("✓ Created movie with ID: %d\n", movie.ID)

	// Test getting the movie
	retrieved, err := pgDB.GetMovie(movie.ID)
	if err != nil {
		log.Fatalf("Failed to get movie: %v", err)
	}
	fmt.Printf("✓ Retrieved movie: %s (%d)\n", retrieved.Title, retrieved.Year)

	// Test database stats
	stats, err := pgDB.GetStats()
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	fmt.Printf("\nDatabase Stats:\n")
	fmt.Printf("- Total movies: %d\n", stats.TotalMovies)
	fmt.Printf("- Database size: %s\n", stats.DatabaseSize)

	fmt.Println("\n✅ All tests passed!")
}