package support

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite" // SQLite driver
)

// SQLiteTestDatabase provides database operations for BDD tests with fixture management
type SQLiteTestDatabase struct {
	db          *sql.DB
	dbPath      string
	fixtures    map[string][]interface{}
	isTemporary bool
}

// NewSQLiteTestDatabase creates a new SQLite test database instance
// If dbPath is empty, creates a temporary database
func NewSQLiteTestDatabase(dbPath string) (*SQLiteTestDatabase, error) {
	isTemporary := false

	// If no path provided, create a temporary database
	if dbPath == "" {
		tempDir := os.TempDir()
		tempFile, err := os.CreateTemp(tempDir, "movies_mcp_test_*.db")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary database file: %w", err)
		}
		dbPath = tempFile.Name()
		tempFile.Close()
		isTemporary = true
	}

	// Open SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		if isTemporary {
			_ = os.Remove(dbPath)
		}
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		if isTemporary {
			_ = os.Remove(dbPath)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		if isTemporary {
			_ = os.Remove(dbPath)
		}
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	testDB := &SQLiteTestDatabase{
		db:          db,
		dbPath:      dbPath,
		fixtures:    make(map[string][]interface{}),
		isTemporary: isTemporary,
	}

	// Run database migrations
	if err := testDB.runMigrations(); err != nil {
		_ = testDB.Cleanup()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return testDB, nil
}

// runMigrations applies database migrations to the test database
func (tdb *SQLiteTestDatabase) runMigrations() error {
	// Find migrations directory
	migrationsDir := findMigrationsDir()
	if migrationsDir == "" {
		return fmt.Errorf("migrations directory not found")
	}

	// Check for tern migration tool
	if _, err := exec.LookPath("tern"); err == nil {
		// Use tern if available
		return tdb.runTernMigrations(migrationsDir)
	}

	// Otherwise, run migrations manually
	return tdb.runManualMigrations(migrationsDir)
}

// findMigrationsDir locates the migrations directory
func findMigrationsDir() string {
	// Try common locations
	locations := []string{
		"./migrations",
		"../migrations",
		"../../migrations",
		"../../../migrations",
		"../../../../migrations",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			absPath, err := filepath.Abs(loc)
			if err == nil {
				return absPath
			}
		}
	}

	// Try finding from working directory upwards
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		migrationsPath := filepath.Join(wd, "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return ""
}

// runTernMigrations runs migrations using tern
func (tdb *SQLiteTestDatabase) runTernMigrations(migrationsDir string) error {
	cmd := exec.Command("tern", "migrate",
		"--migrations", migrationsDir,
		"--conn-string", tdb.dbPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tern migration failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// runManualMigrations runs migrations by executing SQL files directly
func (tdb *SQLiteTestDatabase) runManualMigrations(migrationsDir string) error {
	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort and execute .sql files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsDir, file.Name()))
		}
	}

	for _, sqlFile := range sqlFiles {
		if err := tdb.executeMigrationFile(sqlFile); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filepath.Base(sqlFile), err)
		}
	}

	return nil
}

// executeMigrationFile executes a single migration SQL file
func (tdb *SQLiteTestDatabase) executeMigrationFile(filePath string) error {
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Split by semicolons and execute each statement
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := tdb.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %w\nStatement: %s", err, stmt)
		}
	}

	return nil
}

// LoadFixtures loads test data from a YAML fixture file
func (tdb *SQLiteTestDatabase) LoadFixtures(fixtureName string) error {
	inserter := NewDatabaseFixtureInserter(tdb.db)
	return LoadFixturesFromFile(fixtureName, inserter)
}

// CleanupAfterScenario cleans up test data after each scenario
func (tdb *SQLiteTestDatabase) CleanupAfterScenario() error {
	// Delete all data from test tables in reverse dependency order
	tables := []string{
		"movie_actors", // Junction table first
		"actors",
		"movies",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := tdb.db.Exec(query); err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	// Reset auto-increment counters (SQLite specific)
	if _, err := tdb.db.Exec("DELETE FROM sqlite_sequence"); err != nil {
		// Ignore error if table doesn't exist
	}

	return nil
}

// Cleanup closes the database and removes temporary files
func (tdb *SQLiteTestDatabase) Cleanup() error {
	var errors []error

	// Close database connection
	if tdb.db != nil {
		if err := tdb.db.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close database: %w", err))
		}
	}

	// Remove temporary database file
	if tdb.isTemporary && tdb.dbPath != "" {
		if err := os.Remove(tdb.dbPath); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("failed to remove temporary database: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

// GetDB returns the underlying database connection for custom queries
func (tdb *SQLiteTestDatabase) GetDB() *sql.DB {
	return tdb.db
}

// GetDBPath returns the path to the database file
func (tdb *SQLiteTestDatabase) GetDBPath() string {
	return tdb.dbPath
}

// ExecuteQuery executes a custom SQL query (for test verification)
func (tdb *SQLiteTestDatabase) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return tdb.db.Query(query, args...)
}

// CountRows counts rows in a table with optional WHERE condition
func (tdb *SQLiteTestDatabase) CountRows(table string, whereClause string, args ...interface{}) (int, error) {
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

// VerifyMovieExists checks if a movie with given attributes exists
func (tdb *SQLiteTestDatabase) VerifyMovieExists(title, director string, year int) (bool, error) {
	query := `SELECT COUNT(*) FROM movies WHERE title = ? AND director = ? AND year = ?`
	var count int
	err := tdb.db.QueryRow(query, title, director, year).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// VerifyActorExists checks if an actor with given attributes exists
func (tdb *SQLiteTestDatabase) VerifyActorExists(name string, birthYear int) (bool, error) {
	query := `SELECT COUNT(*) FROM actors WHERE name = ? AND birth_year = ?`
	var count int
	err := tdb.db.QueryRow(query, name, birthYear).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// VerifySchemaExists verifies that the database schema is properly initialized
func (tdb *SQLiteTestDatabase) VerifySchemaExists() error {
	requiredTables := []string{"movies", "actors", "movie_actors"}

	for _, table := range requiredTables {
		var exists bool
		query := `SELECT COUNT(*) > 0 FROM sqlite_master WHERE type='table' AND name=?`
		err := tdb.db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("required table %s does not exist - please run database migrations", table)
		}
	}

	return nil
}
