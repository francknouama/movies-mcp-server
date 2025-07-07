package support

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FixtureInserter defines the interface for inserting fixture data
type FixtureInserter interface {
	insertMovie(movie Movie) error
	insertActor(actor Actor) error
}

// LoadFixturesFromFile loads test data from a YAML fixture file for any database type
func LoadFixturesFromFile(fixtureName string, inserter FixtureInserter) error {
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
		err = inserter.insertMovie(movie)
		if err != nil {
			return fmt.Errorf("failed to insert movie fixture: %w", err)
		}
	}

	// Load actors
	for _, actor := range fixtures.Actors {
		err = inserter.insertActor(actor)
		if err != nil {
			return fmt.Errorf("failed to insert actor fixture: %w", err)
		}
	}

	return nil
}

// DatabaseFixtureInserter provides common fixture insertion logic
type DatabaseFixtureInserter struct {
	db *sql.DB
}

// NewDatabaseFixtureInserter creates a new database fixture inserter
func NewDatabaseFixtureInserter(db *sql.DB) *DatabaseFixtureInserter {
	return &DatabaseFixtureInserter{db: db}
}

// insertMovie inserts a movie fixture into the database
func (d *DatabaseFixtureInserter) insertMovie(movie Movie) error {
	query := `
		INSERT INTO movies (id, title, director, year, genre, rating, genres)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`

	var genres []byte
	if len(movie.Genres) > 0 {
		genresJSON := "["
		for i, genre := range movie.Genres {
			if i > 0 {
				genresJSON += ", "
			}
			genresJSON += fmt.Sprintf(`"%s"`, genre)
		}
		genresJSON += "]"
		genres = []byte(genresJSON)
	}

	_, err := d.db.Exec(query, movie.ID, movie.Title, movie.Director, movie.Year, movie.Genre, movie.Rating, genres)
	return err
}

// insertActor inserts an actor fixture into the database
func (d *DatabaseFixtureInserter) insertActor(actor Actor) error {
	query := `
		INSERT INTO actors (id, name, birth_year, bio)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := d.db.Exec(query, actor.ID, actor.Name, actor.BirthYear, actor.Bio)
	return err
}