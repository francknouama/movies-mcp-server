package database

import (
	"database/sql"
	"time"
)

// Movie represents a movie in the database
type Movie struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Director    string    `json:"director" db:"director"`
	Year        int       `json:"year" db:"year"`
	Genre       []string  `json:"genre" db:"genre"`
	Rating      sql.NullFloat64 `json:"rating" db:"rating"`
	Description sql.NullString  `json:"description" db:"description"`
	Duration    sql.NullInt32   `json:"duration" db:"duration"` // minutes
	Language    sql.NullString  `json:"language" db:"language"`
	Country     sql.NullString  `json:"country" db:"country"`
	PosterData  []byte    `json:"-" db:"poster_data"`     // Binary image data
	PosterType  string    `json:"poster_type,omitempty" db:"poster_type"` // MIME type
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Database interface defines all database operations
type Database interface {
	// CRUD operations
	CreateMovie(movie *Movie) error
	GetMovie(id int) (*Movie, error)
	UpdateMovie(movie *Movie) error
	DeleteMovie(id int) error
	
	// Search operations
	SearchMovies(query SearchQuery) ([]*Movie, error)
	ListTopMovies(limit int, genreFilter string) ([]*Movie, error)
	
	// Image operations
	GetMoviePoster(id int) ([]byte, string, error) // returns data, mime type, error
	UpdateMoviePoster(id int, data []byte, mimeType string) error
	
	// Utility operations
	GetStats() (*DatabaseStats, error)
	GetGenres() ([]GenreCount, error)
	GetDirectors() ([]DirectorCount, error)
	
	// Health and cleanup
	Ping() error
	Close() error
}

// SearchQuery represents search parameters
type SearchQuery struct {
	Query      string
	Type       string // title, director, genre, year, fulltext
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}

// DatabaseStats represents database statistics
type DatabaseStats struct {
	TotalMovies    int            `json:"total_movies"`
	AverageRating  float64        `json:"average_rating"`
	MoviesPerYear  map[int]int    `json:"movies_per_year"`
	GenreCount     map[string]int `json:"genre_count"`
	TotalPosters   int            `json:"total_posters"`
	DatabaseSize   string         `json:"database_size"`
}

// GenreCount represents genre statistics
type GenreCount struct {
	Genre string `json:"genre" db:"genre"`
	Count int    `json:"count" db:"count"`
}

// DirectorCount represents director statistics
type DirectorCount struct {
	Director      string  `json:"director" db:"director"`
	MovieCount    int     `json:"movie_count" db:"movie_count"`
	AverageRating float64 `json:"average_rating" db:"average_rating"`
}