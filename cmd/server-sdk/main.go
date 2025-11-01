package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	actorApp "github.com/francknouama/movies-mcp-server/internal/application/actor"
	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
	"github.com/francknouama/movies-mcp-server/internal/config"
	"github.com/francknouama/movies-mcp-server/internal/domain/actor"
	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/infrastructure/postgres"
	"github.com/francknouama/movies-mcp-server/internal/infrastructure/sqlite"
	"github.com/francknouama/movies-mcp-server/internal/mcp/resources"
	"github.com/francknouama/movies-mcp-server/internal/mcp/tools"
)

var (
	// Build-time variables (set by goreleaser)
	version = "dev-sdk"
	commit  = "none"
	date    = "unknown"
)

const name = "movies-mcp-server-sdk"

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
		fmt.Printf("%s version %s (SDK-based)\n", name, version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
		fmt.Printf("SDK: github.com/modelcontextprotocol/go-sdk v1.1.0\n")
		os.Exit(0)
	}

	if *showHelp {
		fmt.Printf("Movies MCP Server (SDK Edition) - Official Golang MCP SDK Implementation\n")
		fmt.Printf("Built with Clean Architecture and the official Model Context Protocol SDK\n\n")
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nThe server communicates via stdin/stdout using the MCP protocol.\n")
		fmt.Printf("\nFeatures:\n")
		fmt.Printf("  - Official MCP SDK integration\n")
		fmt.Printf("  - Type-safe tool handlers with automatic schema generation\n")
		fmt.Printf("  - 23 tools across movie/actor management, search, and analysis\n")
		fmt.Printf("  - 3 database resources for movie data and statistics\n")
		fmt.Printf("  - Clean Architecture with Domain-Driven Design\n")
		fmt.Printf("  - SQLite or PostgreSQL with automatic migrations\n")
		os.Exit(0)
	}

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

	fmt.Fprintf(os.Stderr, "Connected to database: %s (driver: %s)\n", cfg.Database.Name, cfg.Database.Driver)
	fmt.Fprintf(os.Stderr, "Starting Movies MCP Server with Official SDK...\n")

	// Initialize repositories based on driver
	var movieRepo movie.Repository
	var actorRepo actor.Repository

	if cfg.Database.Driver == "sqlite" {
		movieRepo = sqlite.NewMovieRepository(db)
		actorRepo = sqlite.NewActorRepository(db)
	} else {
		movieRepo = postgres.NewMovieRepository(db)
		actorRepo = postgres.NewActorRepository(db)
	}

	// Initialize services
	movieService := movieApp.NewService(movieRepo)
	actorService := actorApp.NewService(actorRepo)

	// Initialize SDK-based tool handlers
	movieTools := tools.NewMovieTools(movieService)
	actorTools := tools.NewActorTools(actorService)
	compoundTools := tools.NewCompoundTools(movieService)
	contextTools := tools.NewContextTools(movieService)

	// Initialize resource handlers
	dbResources := resources.NewDatabaseResources(movieService)

	// Create MCP server with SDK
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    name,
			Version: version,
		},
		nil, // Options
	)

	fmt.Fprintf(os.Stderr, "Registering tools with SDK...\n")

	// Register Movie Tools (8 tools)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_movie",
		Description: "Get a movie by ID",
	}, movieTools.GetMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_movie",
		Description: "Add a new movie to the database",
	}, movieTools.AddMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_movie",
		Description: "Update an existing movie",
	}, movieTools.UpdateMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_movie",
		Description: "Delete a movie by ID",
	}, movieTools.DeleteMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_top_movies",
		Description: "Get top-rated movies",
	}, movieTools.ListTopMovies)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_movies",
		Description: "Search for movies with various filters",
	}, movieTools.SearchMovies)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_by_decade",
		Description: "Search movies by decade (e.g., 1990s, 90s)",
	}, movieTools.SearchByDecade)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_by_rating_range",
		Description: "Search movies by rating range",
	}, movieTools.SearchByRatingRange)

	// Register Actor Tools (9 tools)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_actor",
		Description: "Get an actor by ID",
	}, actorTools.GetActor)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_actor",
		Description: "Add a new actor to the database",
	}, actorTools.AddActor)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_actor",
		Description: "Update an existing actor",
	}, actorTools.UpdateActor)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_actor",
		Description: "Delete an actor by ID",
	}, actorTools.DeleteActor)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "link_actor_to_movie",
		Description: "Link an actor to a movie",
	}, actorTools.LinkActorToMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "unlink_actor_from_movie",
		Description: "Unlink an actor from a movie",
	}, actorTools.UnlinkActorFromMovie)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_movie_cast",
		Description: "Get all actors in a movie",
	}, actorTools.GetMovieCast)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_actor_movies",
		Description: "Get all movies for an actor",
	}, actorTools.GetActorMovies)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_actors",
		Description: "Search for actors with various filters",
	}, actorTools.SearchActors)

	// Register Compound Tools (3 tools)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "bulk_movie_import",
		Description: "Import multiple movies at once",
	}, compoundTools.BulkMovieImport)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "movie_recommendation_engine",
		Description: "Get personalized movie recommendations based on preferences",
	}, compoundTools.MovieRecommendationEngine)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "director_career_analysis",
		Description: "Analyze a director's career trajectory and filmography",
	}, compoundTools.DirectorCareerAnalysis)

	// Register Context Management Tools (3 tools)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_search_context",
		Description: "Create a paginated context for large search results",
	}, contextTools.CreateSearchContext)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_context_page",
		Description: "Get a specific page from a search context",
	}, contextTools.GetContextPage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_context_info",
		Description: "Get metadata about a search context",
	}, contextTools.GetContextInfo)

	fmt.Fprintf(os.Stderr, "✓ Registered 23 tools successfully\n")
	fmt.Fprintf(os.Stderr, "  - Movie tools: 8\n")
	fmt.Fprintf(os.Stderr, "  - Actor tools: 9\n")
	fmt.Fprintf(os.Stderr, "  - Compound tools: 3\n")
	fmt.Fprintf(os.Stderr, "  - Context tools: 3\n")

	fmt.Fprintf(os.Stderr, "Registering resources with SDK...\n")

	// Register Database Resources (3 resources)
	server.AddResource(dbResources.AllMoviesResource(), dbResources.HandleAllMovies)
	server.AddResource(dbResources.DatabaseStatsResource(), dbResources.HandleDatabaseStats)
	server.AddResource(dbResources.PosterCollectionResource(), dbResources.HandlePosterCollection)

	fmt.Fprintf(os.Stderr, "✓ Registered 3 resources successfully\n")
	fmt.Fprintf(os.Stderr, "  - movies://database/all\n")
	fmt.Fprintf(os.Stderr, "  - movies://database/stats\n")
	fmt.Fprintf(os.Stderr, "  - movies://posters/collection\n")

	fmt.Fprintf(os.Stderr, "\nServer ready - listening on stdin/stdout\n")
	fmt.Fprintf(os.Stderr, "Using official MCP SDK v1.1.0\n\n")

	// Run server with stdio transport
	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// connectToDatabase establishes a connection to the database with retries
func connectToDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	var driverName string
	var dsn string

	if cfg.Driver == "sqlite" {
		driverName = "sqlite"
		dsn = cfg.ConnectionString()
	} else {
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	}

	var db *sql.DB
	var err error

	// Retry connection with exponential backoff
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open(driverName, dsn)
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

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Construct from individual components based on driver
		driver := os.Getenv("DB_DRIVER")
		if driver == "" {
			driver = "sqlite"
		}

		if driver == "sqlite" {
			dbname := os.Getenv("DB_NAME")
			if dbname == "" {
				dbname = "movies.db"
			}
			dbURL = fmt.Sprintf("sqlite://%s", dbname)
		} else {
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
	}

	// Run the migration tool
	fmt.Fprintf(os.Stderr, "Running migrations...\n")
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
