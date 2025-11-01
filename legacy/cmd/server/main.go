package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/francknouama/movies-mcp-server/internal/composition"
	"github.com/francknouama/movies-mcp-server/internal/config"
	"github.com/francknouama/movies-mcp-server/internal/server"
)

var (
	// Build-time variables (set by goreleaser)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const name = "movies-mcp-server"

func main() {
	var (
		showVersion    = flag.Bool("version", false, "Show version information")
		showHelp       = flag.Bool("help", false, "Show help information")
		skipMigrations = flag.Bool("skip-migrations", false, "Skip database migrations")
		migrateOnly    = flag.Bool("migrate-only", false, "Run migrations and exit")
		migrationsPath = flag.String("migrations", "./migrations", "Path to database migrations")
	)

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", name, version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
		os.Exit(0)
	}

	if *showHelp {
		fmt.Printf("Movies MCP Server - A Model Context Protocol server for movie database\n")
		fmt.Printf("Built with Clean Architecture and Domain-Driven Design\n\n")
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nThe server communicates via stdin/stdout using the MCP protocol.\n")
		fmt.Printf("\nFeatures:\n")
		fmt.Printf("  - Clean Architecture with Domain-Driven Design\n")
		fmt.Printf("  - Type-safe domain models with value objects\n")
		fmt.Printf("  - Comprehensive test coverage\n")
		fmt.Printf("  - PostgreSQL with automatic migrations\n")
		fmt.Printf("  - Actor and movie management\n")
		fmt.Printf("  - Advanced search capabilities\n")
		os.Exit(0)
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived shutdown signal, gracefully shutting down...\n")
		cancel()
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := connectToDatabase(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run database migrations
	if !*skipMigrations {
		if err := runMigrations(*migrationsPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Database migrations completed successfully\n")
	}

	// Exit if only running migrations
	if *migrateOnly {
		fmt.Fprintf(os.Stderr, "Migrations completed, exiting as requested\n")
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Connected to database: %s\n", cfg.Database.Name)
	fmt.Fprintf(os.Stderr, "Starting Movies MCP Server with Clean Architecture...\n")

	// Create dependency container
	container := composition.NewContainer(db)

	// Create logger for the server
	logger := log.New(os.Stderr, "[MCP] ", log.LstdFlags)

	// Create and run the MCP server with new clean architecture
	srv := server.NewMCPServer(os.Stdin, os.Stdout, logger, container)

	// Run server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Run()
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		fmt.Fprintf(os.Stderr, "Shutting down gracefully...\n")
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
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
		fmt.Fprintf(os.Stderr, "Database connection failed, retrying in %v... (attempt %d/%d)\n",
			waitTime, i+1, maxRetries)
		time.Sleep(waitTime)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// runMigrations applies database migrations using our custom tool
func runMigrations(migrationsPath string) error {
	// First, build the migration tool
	fmt.Fprintf(os.Stderr, "Building migration tool...\n")
	buildCmd := exec.Command("go", "build", "-o", "./migrate", "./tools/migrate")
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build migration tool: %w", err)
	}

	// Get database URL from the connection
	// We need to reconstruct the connection string from the database instance
	// For now, we'll use environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Construct from individual components
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "movies_user"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "movies_password"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "movies_mcp"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode)
	}

	// Run the migration tool
	fmt.Fprintf(os.Stderr, "Running migrations...\n")
	// #nosec G204 - dbURL and migrationsPath are validated configuration parameters
	migrateCmd := exec.Command("./migrate", dbURL, migrationsPath, "up")
	migrateCmd.Stdout = os.Stderr
	migrateCmd.Stderr = os.Stderr

	if err := migrateCmd.Run(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Clean up the built binary
	if err := os.Remove("./migrate"); err != nil {
		log.Printf("Warning: failed to remove migrate binary: %v", err)
	}

	return nil
}
