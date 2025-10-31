package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/config"
	"github.com/francknouama/movies-mcp-server/internal/infrastructure/postgres"
	"github.com/francknouama/movies-mcp-server/internal/mcp/tools"
)

// This is a proof-of-concept demonstrating the SDK-based approach
// It shows how to register one tool (get_movie) with the official MCP SDK

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Starting SDK-based MCP server (POC)")

	// Connect to database
	db, err := connectToDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Initialize repositories
	movieRepo := postgres.NewMovieRepository(db)

	// Initialize services
	movieService := movie.NewService(movieRepo)

	// Initialize SDK-based tools
	movieTools := tools.NewMovieTools(movieService)

	// Create MCP server with SDK
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "movies-mcp-server-sdk-poc",
			Version: "0.1.0",
		},
		nil, // Options
	)

	// Register the get_movie tool using SDK
	// The SDK automatically:
	// 1. Generates JSON schema from GetMovieInput struct tags
	// 2. Handles JSON-RPC protocol
	// 3. Validates input against schema
	// 4. Marshals output to proper MCP format
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_movie",
			Description: "Get a movie by ID",
		},
		movieTools.GetMovie,
	)

	log.Println("Registered tools: get_movie")
	log.Println("Server ready - listening on stdin/stdout")

	// Run server with stdio transport
	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// connectToDatabase establishes a connection to PostgreSQL with retries
func connectToDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	var db *sql.DB
	var err error

	// Retry connection with exponential backoff
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		// Test connection
		err = db.Ping()
		if err == nil {
			break
		}

		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v", closeErr)
		}

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
		}

		waitTime := time.Duration(1<<i) * time.Second
		log.Printf("Database connection failed, retrying in %v... (attempt %d/%d)\n",
			waitTime, i+1, maxRetries)
		time.Sleep(waitTime)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
