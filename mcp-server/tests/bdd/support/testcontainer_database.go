package support

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v2"
)

// TestContainerDatabase provides database operations using Testcontainers for isolation
type TestContainerDatabase struct {
	container *postgres.PostgresContainer
	db        *sql.DB
	fixtures  map[string][]interface{}
	ctx       context.Context
}

// NewTestContainerDatabase creates a new test database instance using Testcontainers
func NewTestContainerDatabase(ctx context.Context) (*TestContainerDatabase, error) {
	// Check if Docker is available
	if !isDockerAvailable() {
		return nil, fmt.Errorf("Docker is not available - please install Docker or use fallback TestDatabase")
	}

	// Create PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("movies_mcp_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = postgresContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		_ = postgresContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		_ = postgresContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	testDB := &TestContainerDatabase{
		container: postgresContainer,
		db:        db,
		fixtures:  make(map[string][]interface{}),
		ctx:       ctx,
	}

	// Run database migrations
	if err := testDB.runMigrations(); err != nil {
		_ = testDB.Cleanup()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return testDB, nil
}

// isDockerAvailable checks if Docker is available and running
func isDockerAvailable() bool {
	// Try to create a simple container request to test Docker availability
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "hello-world",
		Cmd:   []string{"echo", "test"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})

	if err != nil {
		return false
	}

	if container != nil {
		_ = container.Terminate(ctx)
	}

	return true
}

// runMigrations applies database migrations to the test database
func (tdb *TestContainerDatabase) runMigrations() error {
	// Look for migration files relative to the BDD test directory
	migrationsDir := "../../migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found at %s", migrationsDir)
	}

	// Read and execute migration files in order
	migrationFiles := []string{
		"001_create_movies_table.up.sql",
		"002_add_indexes.up.sql",
		"003_add_search_indexes.up.sql",
		"004_create_actors_tables.up.sql",
		"005_align_schema_with_domain.up.sql",
	}

	for _, filename := range migrationFiles {
		migrationPath := filepath.Join(migrationsDir, filename)

		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if _, err := tdb.db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	return nil
}

// LoadFixtures loads test data from a YAML fixture file (same as original implementation)
func (tdb *TestContainerDatabase) LoadFixtures(fixtureName string) error {
	fixturesDir := "fixtures"
	fixturePath := filepath.Join(fixturesDir, fixtureName+".yaml")

	data, err := os.ReadFile(fixturePath)
	if err != nil {
		return fmt.Errorf("failed to read fixture file %s: %w", fixturePath, err)
	}

	var fixtures Fixtures
	err = yaml.Unmarshal(data, &fixtures)
	if err != nil {
		return fmt.Errorf("failed to parse fixture file %s: %w", fixturePath, err)
	}

	// Load movies
	for _, movie := range fixtures.Movies {
		err = tdb.insertMovie(movie)
		if err != nil {
			return fmt.Errorf("failed to insert movie fixture: %w", err)
		}
	}

	// Load actors
	for _, actor := range fixtures.Actors {
		err = tdb.insertActor(actor)
		if err != nil {
			return fmt.Errorf("failed to insert actor fixture: %w", err)
		}
	}

	return nil
}

// insertMovie inserts a movie fixture into the database (same as original)
func (tdb *TestContainerDatabase) insertMovie(movie Movie) error {
	query := `
		INSERT INTO movies (id, title, director, year, genre, rating, genres, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			director = EXCLUDED.director,
			year = EXCLUDED.year,
			genre = EXCLUDED.genre,
			rating = EXCLUDED.rating,
			genres = EXCLUDED.genres,
			updated_at = NOW()
	`

	genres := "{}"
	if len(movie.Genres) > 0 {
		genres = fmt.Sprintf("{\"%s\"}", movie.Genres[0])
		for i := 1; i < len(movie.Genres); i++ {
			genres = genres[:len(genres)-1] + fmt.Sprintf(",\"%s\"}", movie.Genres[i])
		}
	}

	_, err := tdb.db.Exec(query, movie.ID, movie.Title, movie.Director, movie.Year, movie.Genre, movie.Rating, genres)
	return err
}

// insertActor inserts an actor fixture into the database (same as original)
func (tdb *TestContainerDatabase) insertActor(actor Actor) error {
	query := `
		INSERT INTO actors (id, name, birth_year, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			birth_year = EXCLUDED.birth_year,
			bio = EXCLUDED.bio,
			updated_at = NOW()
	`

	_, err := tdb.db.Exec(query, actor.ID, actor.Name, actor.BirthYear, actor.Bio)
	if err != nil {
		return err
	}

	// Insert movie-actor relationships
	for _, movieID := range actor.MovieIDs {
		relationQuery := `
			INSERT INTO movie_actors (movie_id, actor_id)
			VALUES ($1, $2)
			ON CONFLICT (movie_id, actor_id) DO NOTHING
		`
		_, err = tdb.db.Exec(relationQuery, movieID, actor.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// CleanupAfterScenario cleans up test data after each scenario
func (tdb *TestContainerDatabase) CleanupAfterScenario() error {
	// Truncate all test tables to ensure clean state
	tables := []string{
		"movie_actors",
		"actors",
		"movies",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		_, err := tdb.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}

// Cleanup terminates the container and cleans up resources
func (tdb *TestContainerDatabase) Cleanup() error {
	var errors []error

	// Close database connection
	if tdb.db != nil {
		if err := tdb.db.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close database: %w", err))
		}
	}

	// Terminate container
	if tdb.container != nil {
		if err := tdb.container.Terminate(tdb.ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to terminate container: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

// GetDB returns the underlying database connection for custom queries
func (tdb *TestContainerDatabase) GetDB() *sql.DB {
	return tdb.db
}

// ExecuteQuery executes a custom SQL query (for test verification)
func (tdb *TestContainerDatabase) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return tdb.db.Query(query, args...)
}

// CountRows counts rows in a table with optional WHERE condition
func (tdb *TestContainerDatabase) CountRows(table string, whereClause string, args ...interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	var count int
	err := tdb.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// VerifyMovieExists checks if a movie with given attributes exists
func (tdb *TestContainerDatabase) VerifyMovieExists(title, director string, year int) (bool, error) {
	query := `SELECT COUNT(*) FROM movies WHERE title = $1 AND director = $2 AND year = $3`
	var count int
	err := tdb.db.QueryRow(query, title, director, year).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// VerifyActorExists checks if an actor with given attributes exists
func (tdb *TestContainerDatabase) VerifyActorExists(name string, birthYear int) (bool, error) {
	query := `SELECT COUNT(*) FROM actors WHERE name = $1 AND birth_year = $2`
	var count int
	err := tdb.db.QueryRow(query, name, birthYear).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetConnectionString returns the database connection string for the containerized database
func (tdb *TestContainerDatabase) GetConnectionString() (string, error) {
	return tdb.container.ConnectionString(tdb.ctx, "sslmode=disable")
}
