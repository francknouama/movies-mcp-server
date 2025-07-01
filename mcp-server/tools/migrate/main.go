package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationTool handles database migrations
type MigrationTool struct {
	db             *sql.DB
	migrationsPath string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <database_url> <migrations_path> [up|down]\n", os.Args[0])
		os.Exit(1)
	}

	dbURL := os.Args[1]
	migrationsPath := os.Args[2]
	command := "up"
	if len(os.Args) > 3 {
		command = os.Args[3]
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	tool := &MigrationTool{
		db:             db,
		migrationsPath: migrationsPath,
	}

	// Ensure migrations table exists
	if err := tool.ensureMigrationsTable(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ensure migrations table: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "up":
		if err := tool.up(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run up migrations: %v\n", err)
			os.Exit(1)
		}
	case "down":
		if err := tool.down(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run down migrations: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// ensureMigrationsTable creates the migrations tracking table if it doesn't exist
func (m *MigrationTool) ensureMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// getCurrentVersion returns the current migration version
func (m *MigrationTool) getCurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	return version, err
}

// loadMigrations loads all migration files from the migrations directory
func (m *MigrationTool) loadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(m.migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Parse migration files (format: 001_create_movies.up.sql, 001_create_movies.down.sql)
		filename := d.Name()
		if !strings.HasSuffix(filename, ".sql") {
			return nil
		}

		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return nil
		}

		versionStr := parts[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return fmt.Errorf("invalid version in filename %s: %v", filename, err)
		}

		// Extract migration name and direction
		nameAndDirection := strings.Join(parts[1:], "_")
		nameAndDirection = strings.TrimSuffix(nameAndDirection, ".sql")

		var direction string
		var name string
		if strings.HasSuffix(nameAndDirection, ".up") {
			direction = "up"
			name = strings.TrimSuffix(nameAndDirection, ".up")
		} else if strings.HasSuffix(nameAndDirection, ".down") {
			direction = "down"
			name = strings.TrimSuffix(nameAndDirection, ".down")
		} else {
			return nil // Skip files that don't match our pattern
		}

		// Read file content (validate path for security)
		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", path, err)
		}

		// Find or create migration
		var migration *Migration
		for i := range migrations {
			if migrations[i].Version == version && migrations[i].Name == name {
				migration = &migrations[i]
				break
			}
		}

		if migration == nil {
			migrations = append(migrations, Migration{
				Version: version,
				Name:    name,
			})
			migration = &migrations[len(migrations)-1]
		}

		// Set SQL content
		if direction == "up" {
			migration.UpSQL = string(content)
		} else {
			migration.DownSQL = string(content)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// up runs all pending migrations
func (m *MigrationTool) up() error {
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %v", err)
	}

	applied := 0
	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue
		}

		if migration.UpSQL == "" {
			fmt.Printf("Warning: No up migration found for version %d (%s)\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		// Execute migration in a transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %v", migration.Version, err)
		}

		// Execute the migration SQL
		if _, err := tx.Exec(migration.UpSQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %v", migration.Version, err)
		}

		// Record the migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %v", migration.Version, err)
		}

		applied++
	}

	if applied == 0 {
		fmt.Println("No migrations to apply (database is up to date)")
	} else {
		fmt.Printf("Successfully applied %d migrations\n", applied)
	}

	return nil
}

// down rolls back the last migration
func (m *MigrationTool) down() error {
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	if currentVersion == 0 {
		fmt.Println("No migrations to roll back")
		return nil
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %v", err)
	}

	// Find the migration to roll back
	var targetMigration *Migration
	for i := range migrations {
		if migrations[i].Version == currentVersion {
			targetMigration = &migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %d not found", currentVersion)
	}

	if targetMigration.DownSQL == "" {
		return fmt.Errorf("no down migration found for version %d (%s)", targetMigration.Version, targetMigration.Name)
	}

	fmt.Printf("Rolling back migration %d: %s\n", targetMigration.Version, targetMigration.Name)

	// Execute rollback in a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for rollback %d: %v", targetMigration.Version, err)
	}

	// Execute the down migration SQL
	if _, err := tx.Exec(targetMigration.DownSQL); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute rollback %d: %v", targetMigration.Version, err)
	}

	// Remove the migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", targetMigration.Version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %v", targetMigration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %v", targetMigration.Version, err)
	}

	fmt.Printf("Successfully rolled back migration %d\n", targetMigration.Version)
	return nil
}
