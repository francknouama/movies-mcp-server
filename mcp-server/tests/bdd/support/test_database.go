package support

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gopkg.in/yaml.v2"
)

// TestDatabase provides database operations for BDD tests with fixture management
type TestDatabase struct {
	db       *sql.DB
	fixtures map[string][]interface{}
}

// Movie represents a movie fixture
type Movie struct {
	ID       int      `yaml:"id" json:"id"`
	Title    string   `yaml:"title" json:"title"`
	Director string   `yaml:"director" json:"director"`
	Year     int      `yaml:"year" json:"year"`
	Genre    string   `yaml:"genre" json:"genre"`
	Rating   float64  `yaml:"rating" json:"rating"`
	Genres   []string `yaml:"genres" json:"genres"`
}

// Actor represents an actor fixture
type Actor struct {
	ID        int    `yaml:"id" json:"id"`
	Name      string `yaml:"name" json:"name"`
	BirthYear int    `yaml:"birth_year" json:"birth_year"`
	Bio       string `yaml:"bio" json:"bio"`
	MovieIDs  []int  `yaml:"movie_ids" json:"movie_ids"`
}

// Fixtures represents the structure of a fixture file
type Fixtures struct {
	Movies []Movie `yaml:"movies"`
	Actors []Actor `yaml:"actors"`
}

// NewTestDatabase creates a new test database instance
func NewTestDatabase() (*TestDatabase, error) {
	// Get database configuration from environment variables with defaults
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbName := getEnvOrDefault("DB_NAME", "movies_mcp_test")
	dbUser := getEnvOrDefault("DB_USER", "movies_user")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "movies_password")

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection with retry logic
	if err := testDatabaseConnection(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Verify that required tables exist
	if err := verifyDatabaseSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database schema verification failed: %w", err)
	}

	return &TestDatabase{
		db:       db,
		fixtures: make(map[string][]interface{}),
	}, nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// testDatabaseConnection tests the database connection with retries
func testDatabaseConnection(db *sql.DB) error {
	// Simple ping test
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test basic query
	var count int
	err := db.QueryRow("SELECT 1").Scan(&count)
	if err != nil {
		return fmt.Errorf("basic query test failed: %w", err)
	}

	return nil
}

// verifyDatabaseSchema verifies that required tables exist
func verifyDatabaseSchema(db *sql.DB) error {
	requiredTables := []string{"movies", "actors", "movie_actors"}

	for _, table := range requiredTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' AND table_name = $1
			)
		`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("required table %s does not exist - please run database migrations", table)
		}
	}

	return nil
}

// LoadFixtures loads test data from a YAML fixture file
func (tdb *TestDatabase) LoadFixtures(fixtureName string) error {
	// Validate fixture name to prevent path traversal
	if !isValidFixtureName(fixtureName) {
		return fmt.Errorf("invalid fixture name: %s", fixtureName)
	}

	fixturesDir := "fixtures"
	fixturePath := filepath.Join(fixturesDir, fixtureName+".yaml")

	data, err := os.ReadFile(filepath.Clean(fixturePath))
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

// insertMovie inserts a movie fixture into the database
func (tdb *TestDatabase) insertMovie(movie Movie) error {
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

// insertActor inserts an actor fixture into the database
func (tdb *TestDatabase) insertActor(actor Actor) error {
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
func (tdb *TestDatabase) CleanupAfterScenario() error {
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

// CleanupAfterSuite cleans up after the entire test suite
func (tdb *TestDatabase) CleanupAfterSuite() error {
	if tdb.db != nil {
		return tdb.db.Close()
	}
	return nil
}

// GetDB returns the underlying database connection for custom queries
func (tdb *TestDatabase) GetDB() *sql.DB {
	return tdb.db
}

// ExecuteQuery executes a custom SQL query (for test verification)
func (tdb *TestDatabase) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return tdb.db.Query(query, args...)
}

// CountRows counts rows in a table with optional WHERE condition
func (tdb *TestDatabase) CountRows(table string, whereClause string, args ...interface{}) (int, error) {
	// Validate table name to prevent SQL injection (test context only)
	validTables := map[string]bool{
		"movies":       true,
		"actors":       true,
		"movie_actors": true,
	}
	if !validTables[table] {
		return 0, fmt.Errorf("invalid table name: %s", table)
	}
	
	// #nosec G201 - table name is validated against whitelist above
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	var count int
	err := tdb.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// isValidFixtureName validates fixture names to prevent path traversal attacks
func isValidFixtureName(name string) bool {
	// Only allow alphanumeric characters, hyphens, and underscores
	// No path separators or relative path indicators
	if strings.Contains(name, "/") || strings.Contains(name, "\\") ||
		strings.Contains(name, "..") || strings.Contains(name, ".") {
		return false
	}
	// Additional check for empty names
	return len(strings.TrimSpace(name)) > 0
}

// VerifyMovieExists checks if a movie with given attributes exists
func (tdb *TestDatabase) VerifyMovieExists(title, director string, year int) (bool, error) {
	query := `SELECT COUNT(*) FROM movies WHERE title = $1 AND director = $2 AND year = $3`
	var count int
	err := tdb.db.QueryRow(query, title, director, year).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// VerifyActorExists checks if an actor with given attributes exists
func (tdb *TestDatabase) VerifyActorExists(name string, birthYear int) (bool, error) {
	query := `SELECT COUNT(*) FROM actors WHERE name = $1 AND birth_year = $2`
	var count int
	err := tdb.db.QueryRow(query, name, birthYear).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
