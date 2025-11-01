package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
	"github.com/francknouama/movies-mcp-server/pkg/database"
)

// MovieRepository implements the movie.Repository interface for SQLite
type MovieRepository struct {
	*database.BaseRepository
}

// NewMovieRepository creates a new SQLite movie repository
func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// dbMovie represents the database model for movies
type dbMovie struct {
	ID          int             `db:"id"`
	Title       string          `db:"title"`
	Director    string          `db:"director"`
	Year        int             `db:"year"`
	Rating      sql.NullFloat64 `db:"rating"`
	Genres      string          `db:"genre"` // JSON-encoded array
	Description sql.NullString  `db:"description"`
	Duration    sql.NullInt64   `db:"duration"`
	Language    sql.NullString  `db:"language"`
	Country     sql.NullString  `db:"country"`
	PosterData  []byte          `db:"poster_data"`
	PosterType  sql.NullString  `db:"poster_type"`
	PosterURL   sql.NullString  `db:"poster_url"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}

// Save persists a movie (insert or update)
func (r *MovieRepository) Save(ctx context.Context, domainMovie *movie.Movie) error {
	dbMovie, err := r.toDBModel(domainMovie)
	if err != nil {
		return fmt.Errorf("failed to convert to DB model: %w", err)
	}

	if domainMovie.ID().IsZero() {
		return r.insert(ctx, dbMovie, domainMovie)
	}
	return r.update(ctx, dbMovie, domainMovie)
}

func (r *MovieRepository) insert(ctx context.Context, dbMovie *dbMovie, domainMovie *movie.Movie) error {
	query := `
		INSERT INTO movies (title, director, year, rating, genre, poster_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	id, err := r.InsertWithID(ctx, query,
		dbMovie.Title,
		dbMovie.Director,
		dbMovie.Year,
		dbMovie.Rating,
		dbMovie.Genres,
		dbMovie.PosterURL,
		dbMovie.CreatedAt.Time,
		dbMovie.UpdatedAt.Time,
	)

	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	// Update domain movie with the new ID
	movieID, err := shared.NewMovieID(id)
	if err != nil {
		return fmt.Errorf("failed to create movie ID: %w", err)
	}
	domainMovie.SetID(movieID)

	return nil
}

func (r *MovieRepository) update(ctx context.Context, dbMovie *dbMovie, domainMovie *movie.Movie) error {
	query := `
		UPDATE movies
		SET title = ?, director = ?, year = ?, rating = ?, genre = ?,
		    poster_url = ?, updated_at = ?
		WHERE id = ?`

	return r.Update(ctx, query, "movie",
		dbMovie.Title,
		dbMovie.Director,
		dbMovie.Year,
		dbMovie.Rating,
		dbMovie.Genres,
		dbMovie.PosterURL,
		dbMovie.UpdatedAt.Time,
		domainMovie.ID().Value(),
	)
}

// FindByID retrieves a movie by its ID
func (r *MovieRepository) FindByID(ctx context.Context, id shared.MovieID) (*movie.Movie, error) {
	query := `
		SELECT id, title, director, year, rating, genre, poster_url, created_at, updated_at
		FROM movies
		WHERE id = ?`

	var dbMovie dbMovie
	err := r.QueryRowContext(ctx, query, id.Value()).Scan(
		&dbMovie.ID,
		&dbMovie.Title,
		&dbMovie.Director,
		&dbMovie.Year,
		&dbMovie.Rating,
		&dbMovie.Genres,
		&dbMovie.PosterURL,
		&dbMovie.CreatedAt,
		&dbMovie.UpdatedAt,
	)

	if err != nil {
		return nil, r.WrapNotFound(err, "movie")
	}

	return r.toDomainModel(&dbMovie)
}

// FindByCriteria retrieves movies based on search criteria
func (r *MovieRepository) FindByCriteria(ctx context.Context, criteria movie.SearchCriteria) ([]*movie.Movie, error) {
	query, args := r.buildSearchQuery(criteria)

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	var movies []*movie.Movie
	for rows.Next() {
		var dbMovie dbMovie
		err := rows.Scan(
			&dbMovie.ID,
			&dbMovie.Title,
			&dbMovie.Director,
			&dbMovie.Year,
			&dbMovie.Rating,
			&dbMovie.Genres,
			&dbMovie.PosterURL,
			&dbMovie.CreatedAt,
			&dbMovie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		domainMovie, err := r.toDomainModel(&dbMovie)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain model: %w", err)
		}

		movies = append(movies, domainMovie)
	}

	return movies, nil
}

func (r *MovieRepository) buildSearchQuery(criteria movie.SearchCriteria) (string, []interface{}) {
	query := `
		SELECT id, title, director, year, rating, genre, poster_url, created_at, updated_at
		FROM movies WHERE 1=1`

	var args []interface{}

	// Add WHERE conditions using ? placeholders
	if criteria.Title != "" {
		query += " AND title LIKE ? COLLATE NOCASE"
		args = append(args, "%"+criteria.Title+"%")
	}

	if criteria.Director != "" {
		query += " AND director LIKE ? COLLATE NOCASE"
		args = append(args, "%"+criteria.Director+"%")
	}

	if criteria.Genre != "" {
		// SQLite JSON search: check if genre array contains the value
		query += " AND EXISTS (SELECT 1 FROM json_each(genre) WHERE value = ?)"
		args = append(args, criteria.Genre)
	}

	if criteria.MinYear > 0 {
		query += " AND year >= ?"
		args = append(args, criteria.MinYear)
	}

	if criteria.MaxYear > 0 {
		query += " AND year <= ?"
		args = append(args, criteria.MaxYear)
	}

	if criteria.MinRating > 0 {
		query += " AND rating >= ?"
		args = append(args, criteria.MinRating)
	}

	if criteria.MaxRating > 0 {
		query += " AND rating <= ?"
		args = append(args, criteria.MaxRating)
	}

	// Add ORDER BY
	orderField := "title"
	switch criteria.OrderBy {
	case movie.OrderByDirector:
		orderField = "director"
	case movie.OrderByYear:
		orderField = "year"
	case movie.OrderByRating:
		orderField = "rating"
	case movie.OrderByCreatedAt:
		orderField = "created_at"
	case movie.OrderByUpdatedAt:
		orderField = "updated_at"
	}

	orderDir := "ASC"
	if criteria.OrderDir == movie.OrderDesc {
		orderDir = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderField, orderDir)

	// Add LIMIT and OFFSET
	if criteria.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, criteria.Limit)
	}

	if criteria.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, criteria.Offset)
	}

	return query, args
}

// FindByTitle searches movies by title (partial match)
func (r *MovieRepository) FindByTitle(ctx context.Context, title string) ([]*movie.Movie, error) {
	criteria := movie.SearchCriteria{
		Title: title,
		Limit: 100, // Default limit
	}
	return r.FindByCriteria(ctx, criteria)
}

// FindByDirector retrieves movies by director
func (r *MovieRepository) FindByDirector(ctx context.Context, director string) ([]*movie.Movie, error) {
	criteria := movie.SearchCriteria{
		Director: director,
		Limit:    100, // Default limit
	}
	return r.FindByCriteria(ctx, criteria)
}

// FindByGenre retrieves movies that have a specific genre
func (r *MovieRepository) FindByGenre(ctx context.Context, genre string) ([]*movie.Movie, error) {
	criteria := movie.SearchCriteria{
		Genre: genre,
		Limit: 100, // Default limit
	}
	return r.FindByCriteria(ctx, criteria)
}

// FindTopRated retrieves top-rated movies
func (r *MovieRepository) FindTopRated(ctx context.Context, limit int) ([]*movie.Movie, error) {
	criteria := movie.SearchCriteria{
		OrderBy:   movie.OrderByRating,
		OrderDir:  movie.OrderDesc,
		Limit:     limit,
		MinRating: 0.1, // Only movies with ratings
	}
	return r.FindByCriteria(ctx, criteria)
}

// CountAll returns the total number of movies
func (r *MovieRepository) CountAll(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM movies"
	return r.Count(ctx, query)
}

// Delete removes a movie by ID
func (r *MovieRepository) Delete(ctx context.Context, id shared.MovieID) error {
	query := "DELETE FROM movies WHERE id = ?"
	return r.BaseRepository.Delete(ctx, query, "movie", id.Value())
}

// DeleteAll removes all movies (for testing)
func (r *MovieRepository) DeleteAll(ctx context.Context) error {
	query := "DELETE FROM movies"
	_, err := r.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all movies: %w", err)
	}
	return nil
}

// toDBModel converts a domain movie to a database model
func (r *MovieRepository) toDBModel(domainMovie *movie.Movie) (*dbMovie, error) {
	// Encode genres as JSON
	genresJSON, err := json.Marshal(domainMovie.Genres())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal genres: %w", err)
	}

	dbMovie := &dbMovie{
		ID:       domainMovie.ID().Value(),
		Title:    domainMovie.Title(),
		Director: domainMovie.Director(),
		Year:     domainMovie.Year().Value(),
		Genres:   string(genresJSON),
	}

	// Handle optional rating
	if !domainMovie.Rating().IsZero() {
		dbMovie.Rating = sql.NullFloat64{
			Float64: domainMovie.Rating().Value(),
			Valid:   true,
		}
	}

	// Handle optional poster URL
	if domainMovie.PosterURL() != "" {
		dbMovie.PosterURL = sql.NullString{
			String: domainMovie.PosterURL(),
			Valid:  true,
		}
	}

	// Handle timestamps
	dbMovie.CreatedAt = sql.NullTime{
		Time:  domainMovie.CreatedAt(),
		Valid: true,
	}
	dbMovie.UpdatedAt = sql.NullTime{
		Time:  domainMovie.UpdatedAt(),
		Valid: true,
	}

	return dbMovie, nil
}

// toDomainModel converts a database model to a domain movie
func (r *MovieRepository) toDomainModel(dbMovie *dbMovie) (*movie.Movie, error) {
	movieID, err := shared.NewMovieID(dbMovie.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	domainMovie, err := movie.NewMovieWithID(movieID, dbMovie.Title, dbMovie.Director, dbMovie.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain movie: %w", err)
	}

	// Set rating if present
	if dbMovie.Rating.Valid {
		if err := domainMovie.SetRating(dbMovie.Rating.Float64); err != nil {
			return nil, fmt.Errorf("failed to set rating: %w", err)
		}
	}

	// Decode and add genres
	if dbMovie.Genres != "" && dbMovie.Genres != "null" {
		var genres []string
		if err := json.Unmarshal([]byte(dbMovie.Genres), &genres); err != nil {
			// Handle legacy non-JSON genre data gracefully
			if strings.HasPrefix(dbMovie.Genres, "[") {
				return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
			}
			// Treat as single genre if not JSON
			genres = []string{dbMovie.Genres}
		}

		for _, genre := range genres {
			if err := domainMovie.AddGenre(genre); err != nil {
				return nil, fmt.Errorf("failed to add genre: %w", err)
			}
		}
	}

	// Set poster URL if present
	if dbMovie.PosterURL.Valid {
		if err := domainMovie.SetPosterURL(dbMovie.PosterURL.String); err != nil {
			return nil, fmt.Errorf("failed to set poster URL: %w", err)
		}
	}

	return domainMovie, nil
}
