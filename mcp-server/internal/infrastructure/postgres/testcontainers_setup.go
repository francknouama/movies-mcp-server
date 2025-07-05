//go:build integration
// +build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	container testcontainers.Container
	DB        *sql.DB
	URL       string
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("movies_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	// Clean up container when test is done
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate postgres container: %v", err)
		}
	})

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up database connection when test is done
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	})

	// Wait for database to be ready
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		if i == maxRetries-1 {
			t.Fatalf("Database not ready after %d attempts", maxRetries)
		}
		time.Sleep(time.Second)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	// Delete in correct order to respect foreign key constraints
	tables := []string{
		"movie_actors",
		"actors",
		"movies",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			// Table might not exist yet, log but don't fail
			t.Logf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func runMigrations(db *sql.DB) error {
	// Create tables if they don't exist
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			director VARCHAR(255) NOT NULL,
			year INTEGER NOT NULL,
			genre TEXT[] NOT NULL DEFAULT '{}',
			rating DECIMAL(3,1),
			description TEXT,
			duration INTEGER,
			language VARCHAR(50),
			country VARCHAR(100),
			poster_data BYTEA,
			poster_type VARCHAR(50),
			poster_url TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS actors (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			birth_year INTEGER NOT NULL,
			bio TEXT,
			image_url TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS movie_actors (
			movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
			actor_id INTEGER NOT NULL REFERENCES actors(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (movie_id, actor_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title)`,
		`CREATE INDEX IF NOT EXISTS idx_movies_year ON movies(year)`,
		`CREATE INDEX IF NOT EXISTS idx_actors_name ON actors(name)`,
		`CREATE INDEX IF NOT EXISTS idx_actors_birth_year ON actors(birth_year)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// SetupTestContainer creates a PostgreSQL container for integration tests
// This is useful when you need access to both the container and the database
func SetupTestContainer(t *testing.T) (*postgres.PostgresContainer, *sql.DB) {
	t.Helper()

	ctx := context.Background()

	// Get the project root to find migration files
	projectRoot := findProjectRoot()
	migrationsPath := filepath.Join(projectRoot, "tools", "migrate", "migrations")

	// Start PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("movies_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(filepath.Join(migrationsPath, "*.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	// Clean up container when test is done
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate postgres container: %v", err)
		}
	})

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up database connection when test is done
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	})

	// Wait for database to be ready
	if err := waitForDatabase(db); err != nil {
		t.Fatalf("Database not ready: %v", err)
	}

	return postgresContainer, db
}

func waitForDatabase(db *sql.DB) error {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("database not ready after %d seconds", maxRetries)
}

func findProjectRoot() string {
	// Simple heuristic: look for go.mod
	dir := "."
	for i := 0; i < 10; i++ {
		if _, err := filepath.Abs(filepath.Join(dir, "go.mod")); err == nil {
			abs, _ := filepath.Abs(dir)
			return abs
		}
		dir = filepath.Join(dir, "..")
	}
	return "."
}

