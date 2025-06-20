package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"movies-mcp-server/internal/config"
	"movies-mcp-server/pkg/image"
)

// Movie represents a movie with its poster URL
type Movie struct {
	ID        int
	Title     string
	PosterURL string
}

// posterURLs maps movie titles to their poster URLs
var posterURLs = map[string]string{
	"The Shawshank Redemption":                        "https://image.tmdb.org/t/p/w500/q6y0Go1tsGEsmtFryDOJo3dEmqu.jpg",
	"The Godfather":                                   "https://image.tmdb.org/t/p/w500/3bhkrj58Vtu7enYsRolD1fZdja1.jpg",
	"The Dark Knight":                                 "https://image.tmdb.org/t/p/w500/qJ2tW6WMUDux911r6m7haRef0WH.jpg",
	"The Godfather Part II":                           "https://image.tmdb.org/t/p/w500/hek3koDUyRQk7FIhPXsa6mT2Zc3.jpg",
	"12 Angry Men":                                    "https://image.tmdb.org/t/p/w500/ow3wq89wM8qd5X7hWKxiRfsFf9C.jpg",
	"Schindlers List":                                 "https://image.tmdb.org/t/p/w500/sF1U4EUQS8YHUYjNl3pMGNIQyr0.jpg",
	"The Lord of the Rings: The Return of the King":  "https://image.tmdb.org/t/p/w500/rCzpDGLbOoPwLjy3OAm5NUPOTrC.jpg",
	"Pulp Fiction":                                    "https://image.tmdb.org/t/p/w500/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg",
	"The Good, the Bad and the Ugly":                 "https://image.tmdb.org/t/p/w500/bX2xnavhMYjWDoZp1VM6VnU1xwe.jpg",
	"Fight Club":                                      "https://image.tmdb.org/t/p/w500/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg",
	"Forrest Gump":                                    "https://image.tmdb.org/t/p/w500/arw2vcBveWOVZr6pxd9XTd1TdQa.jpg",
	"Inception":                                       "https://image.tmdb.org/t/p/w500/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
	"The Lord of the Rings: The Two Towers":          "https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9L.jpg",
	"Star Wars: Episode V - The Empire Strikes Back": "https://image.tmdb.org/t/p/w500/nNAeTmF4CtdSgMDplXTDPOpYzsX.jpg",
	"The Lord of the Rings: The Fellowship of the Ring": "https://image.tmdb.org/t/p/w500/6oom5QYQ2yQTMJIbnvbkBL9cHo6.jpg",
	"Goodfellas":                                      "https://image.tmdb.org/t/p/w500/aKuFiU82s5ISJpGZp7YkIr3kCUd.jpg",
	"One Flew Over the Cuckoos Nest":                 "https://image.tmdb.org/t/p/w500/kjWsMh72V6d8KRLV4EOoSJLT1H7.jpg",
	"The Matrix":                                      "https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg",
	"Seven Samurai":                                   "https://image.tmdb.org/t/p/w500/8OKmBV5BUFzmozIC3pPWKHy17kx.jpg",
	"City of God":                                     "https://image.tmdb.org/t/p/w500/k7eYdWvhYQyRQoU2TB2A2Xu2TfD.jpg",
}

func main() {
	fmt.Println("Movies MCP Server - Image Population Tool")
	fmt.Println("=========================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create image processor
	imageProcessor := image.NewImageProcessor(&cfg.Image)

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✅ Connected to database successfully")

	// Get all movies without poster data
	movies, err := getMoviesWithoutPosters(db)
	if err != nil {
		log.Fatalf("Failed to get movies: %v", err)
	}

	fmt.Printf("📋 Found %d movies without poster data\n", len(movies))

	// Process each movie
	successCount := 0
	errorCount := 0

	for _, movie := range movies {
		fmt.Printf("\n🎬 Processing: %s (ID: %d)\n", movie.Title, movie.ID)

		posterURL, exists := posterURLs[movie.Title]
		if !exists {
			fmt.Printf("   ⚠️  No poster URL found for: %s\n", movie.Title)
			errorCount++
			continue
		}

		// Download poster image
		fmt.Printf("   📥 Downloading from: %s\n", posterURL)
		posterData, mimeType, err := imageProcessor.DownloadImageFromURL(posterURL)
		if err != nil {
			fmt.Printf("   ❌ Failed to download poster: %v\n", err)
			errorCount++
			continue
		}

		fmt.Printf("   ✅ Downloaded %d bytes (%s)\n", len(posterData), mimeType)

		// Update database with poster data
		err = updateMoviePoster(db, movie.ID, posterData, mimeType)
		if err != nil {
			fmt.Printf("   ❌ Failed to update database: %v\n", err)
			errorCount++
			continue
		}

		fmt.Printf("   ✅ Updated database successfully\n")
		successCount++
	}

	// Print summary
	fmt.Printf("\n📊 Processing Summary:\n")
	fmt.Printf("   ✅ Successful: %d\n", successCount)
	fmt.Printf("   ❌ Failed: %d\n", errorCount)
	fmt.Printf("   📁 Total: %d\n", len(movies))

	if successCount > 0 {
		fmt.Println("\n🎉 Image population completed successfully!")
	} else {
		fmt.Println("\n⚠️  No images were successfully populated.")
		os.Exit(1)
	}
}

// getMoviesWithoutPosters retrieves all movies that don't have poster data
func getMoviesWithoutPosters(db *sql.DB) ([]Movie, error) {
	query := `
		SELECT id, title 
		FROM movies 
		WHERE poster_data IS NULL OR poster_type IS NULL
		ORDER BY id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie row: %w", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// updateMoviePoster updates a movie's poster data in the database
func updateMoviePoster(db *sql.DB, movieID int, posterData []byte, mimeType string) error {
	query := `
		UPDATE movies 
		SET poster_data = $1, poster_type = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3
	`

	result, err := db.Exec(query, posterData, mimeType, movieID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated for movie ID %d", movieID)
	}

	return nil
}