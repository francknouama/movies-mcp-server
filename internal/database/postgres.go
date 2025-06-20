package database

import (
	"database/sql"
	"fmt"

	"movies-mcp-server/internal/config"

	_ "github.com/lib/pq"
	"github.com/lib/pq"
)

// PostgresDatabase implements the Database interface for PostgreSQL
type PostgresDatabase struct {
	db *sql.DB
}

// NewPostgresDatabase creates a new PostgreSQL database connection
func NewPostgresDatabase(cfg *config.DatabaseConfig) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDatabase{db: db}, nil
}

// CreateMovie inserts a new movie into the database
func (p *PostgresDatabase) CreateMovie(movie *Movie) error {
	query := `
		INSERT INTO movies (title, director, year, genre, rating, description, 
			duration, language, country, poster_data, poster_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	err := p.db.QueryRow(
		query,
		movie.Title,
		movie.Director,
		movie.Year,
		pq.Array(movie.Genre),
		movie.Rating,
		movie.Description,
		movie.Duration,
		movie.Language,
		movie.Country,
		movie.PosterData,
		movie.PosterType,
	).Scan(&movie.ID, &movie.CreatedAt, &movie.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	return nil
}

// GetMovie retrieves a movie by ID
func (p *PostgresDatabase) GetMovie(id int) (*Movie, error) {
	movie := &Movie{}
	var posterType sql.NullString
	
	query := `
		SELECT id, title, director, year, genre, rating, description, 
			duration, language, country, poster_type, created_at, updated_at
		FROM movies
		WHERE id = $1`

	err := p.db.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Director,
		&movie.Year,
		pq.Array(&movie.Genre),
		&movie.Rating,
		&movie.Description,
		&movie.Duration,
		&movie.Language,
		&movie.Country,
		&posterType,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("movie not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	
	if posterType.Valid {
		movie.PosterType = posterType.String
	}

	return movie, nil
}

// UpdateMovie updates an existing movie
func (p *PostgresDatabase) UpdateMovie(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $2, director = $3, year = $4, genre = $5, rating = $6,
			description = $7, duration = $8, language = $9, country = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := p.db.QueryRow(
		query,
		movie.ID,
		movie.Title,
		movie.Director,
		movie.Year,
		pq.Array(movie.Genre),
		movie.Rating,
		movie.Description,
		movie.Duration,
		movie.Language,
		movie.Country,
	).Scan(&movie.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("movie not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}

// DeleteMovie removes a movie from the database
func (p *PostgresDatabase) DeleteMovie(id int) error {
	result, err := p.db.Exec("DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("movie not found")
	}

	return nil
}

// SearchMovies searches for movies based on the query
func (p *PostgresDatabase) SearchMovies(sq SearchQuery) ([]*Movie, error) {
	var query string
	var args []interface{}

	// Default limit and offset
	if sq.Limit <= 0 {
		sq.Limit = 10
	}
	if sq.Limit > 100 {
		sq.Limit = 100
	}

	baseQuery := `
		SELECT id, title, director, year, genre, rating, description, 
			duration, language, country, poster_type, created_at, updated_at
		FROM movies`

	switch sq.Type {
	case "title":
		query = baseQuery + " WHERE title ILIKE $1"
		args = append(args, "%"+sq.Query+"%")
	case "director":
		query = baseQuery + " WHERE director ILIKE $1"
		args = append(args, "%"+sq.Query+"%")
	case "genre":
		query = baseQuery + " WHERE $1 = ANY(genre)"
		args = append(args, sq.Query)
	case "year":
		query = baseQuery + " WHERE year = $1"
		args = append(args, sq.Query)
	case "fulltext":
		query = baseQuery + ` WHERE 
			to_tsvector('english', title || ' ' || director || ' ' || COALESCE(description, '')) 
			@@ plainto_tsquery('english', $1)`
		args = append(args, sq.Query)
	default:
		// Default to title search
		query = baseQuery + " WHERE title ILIKE $1"
		args = append(args, "%"+sq.Query+"%")
	}

	// Add sorting
	if sq.SortBy == "" {
		sq.SortBy = "rating"
	}
	if sq.SortOrder == "" {
		sq.SortOrder = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sq.SortBy, sq.SortOrder)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, sq.Limit, sq.Offset)

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	return p.scanMovies(rows)
}

// ListTopMovies returns the top-rated movies
func (p *PostgresDatabase) ListTopMovies(limit int, genreFilter string) ([]*Movie, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var query string
	var args []interface{}

	if genreFilter != "" {
		query = `
			SELECT id, title, director, year, genre, rating, description, 
				duration, language, country, poster_type, created_at, updated_at
			FROM movies
			WHERE $1 = ANY(genre)
			ORDER BY rating DESC
			LIMIT $2`
		args = append(args, genreFilter, limit)
	} else {
		query = `
			SELECT id, title, director, year, genre, rating, description, 
				duration, language, country, poster_type, created_at, updated_at
			FROM movies
			ORDER BY rating DESC
			LIMIT $1`
		args = append(args, limit)
	}

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list top movies: %w", err)
	}
	defer rows.Close()

	return p.scanMovies(rows)
}

// GetMoviePoster retrieves just the poster data for a movie
func (p *PostgresDatabase) GetMoviePoster(id int) ([]byte, string, error) {
	var posterData []byte
	var posterType string

	query := `SELECT poster_data, poster_type FROM movies WHERE id = $1`
	err := p.db.QueryRow(query, id).Scan(&posterData, &posterType)

	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("movie not found")
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to get poster: %w", err)
	}

	return posterData, posterType, nil
}

// UpdateMoviePoster updates just the poster for a movie
func (p *PostgresDatabase) UpdateMoviePoster(id int, data []byte, mimeType string) error {
	query := `
		UPDATE movies
		SET poster_data = $2, poster_type = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	result, err := p.db.Exec(query, id, data, mimeType)
	if err != nil {
		return fmt.Errorf("failed to update poster: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("movie not found")
	}

	return nil
}

// GetStats returns database statistics
func (p *PostgresDatabase) GetStats() (*DatabaseStats, error) {
	stats := &DatabaseStats{
		MoviesPerYear: make(map[int]int),
		GenreCount:    make(map[string]int),
	}

	// Get total movies and average rating
	err := p.db.QueryRow(`
		SELECT COUNT(*), COALESCE(AVG(rating), 0)
		FROM movies
	`).Scan(&stats.TotalMovies, &stats.AverageRating)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// Get movies per year
	rows, err := p.db.Query(`
		SELECT year, COUNT(*)
		FROM movies
		GROUP BY year
		ORDER BY year
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies per year: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var year, count int
		if err := rows.Scan(&year, &count); err != nil {
			return nil, fmt.Errorf("failed to scan year stats: %w", err)
		}
		stats.MoviesPerYear[year] = count
	}

	// Get genre counts
	rows2, err := p.db.Query(`
		SELECT unnest(genre) as g, COUNT(*)
		FROM movies
		GROUP BY g
		ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get genre stats: %w", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var genre string
		var count int
		if err := rows2.Scan(&genre, &count); err != nil {
			return nil, fmt.Errorf("failed to scan genre stats: %w", err)
		}
		stats.GenreCount[genre] = count
	}

	// Get poster count
	err = p.db.QueryRow(`
		SELECT COUNT(*)
		FROM movies
		WHERE poster_data IS NOT NULL
	`).Scan(&stats.TotalPosters)
	if err != nil {
		return nil, fmt.Errorf("failed to get poster count: %w", err)
	}

	// Get database size
	err = p.db.QueryRow(`
		SELECT pg_size_pretty(pg_database_size(current_database()))
	`).Scan(&stats.DatabaseSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get database size: %w", err)
	}

	return stats, nil
}

// GetGenres returns all unique genres with counts
func (p *PostgresDatabase) GetGenres() ([]GenreCount, error) {
	query := `
		SELECT unnest(genre) as genre, COUNT(*) as count
		FROM movies
		GROUP BY genre
		ORDER BY count DESC`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}
	defer rows.Close()

	var genres []GenreCount
	for rows.Next() {
		var gc GenreCount
		if err := rows.Scan(&gc.Genre, &gc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan genre: %w", err)
		}
		genres = append(genres, gc)
	}

	return genres, nil
}

// GetDirectors returns all directors with their movie counts and average ratings
func (p *PostgresDatabase) GetDirectors() ([]DirectorCount, error) {
	query := `
		SELECT director, COUNT(*) as movie_count, AVG(rating) as average_rating
		FROM movies
		GROUP BY director
		ORDER BY movie_count DESC`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get directors: %w", err)
	}
	defer rows.Close()

	var directors []DirectorCount
	for rows.Next() {
		var dc DirectorCount
		if err := rows.Scan(&dc.Director, &dc.MovieCount, &dc.AverageRating); err != nil {
			return nil, fmt.Errorf("failed to scan director: %w", err)
		}
		directors = append(directors, dc)
	}

	return directors, nil
}

// Ping tests the database connection
func (p *PostgresDatabase) Ping() error {
	return p.db.Ping()
}

// Close closes the database connection
func (p *PostgresDatabase) Close() error {
	return p.db.Close()
}

// scanMovies is a helper function to scan movie rows
func (p *PostgresDatabase) scanMovies(rows *sql.Rows) ([]*Movie, error) {
	var movies []*Movie

	for rows.Next() {
		movie := &Movie{}
		var posterType sql.NullString
		
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Director,
			&movie.Year,
			pq.Array(&movie.Genre),
			&movie.Rating,
			&movie.Description,
			&movie.Duration,
			&movie.Language,
			&movie.Country,
			&posterType,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		
		if posterType.Valid {
			movie.PosterType = posterType.String
		}
		
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return movies, nil
}