package movie

import (
	"context"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// Repository defines the interface for movie data access
type Repository interface {
	MovieReader
	MovieWriter
}

// MovieReader defines read operations for movies
type MovieReader interface {
	// FindByID retrieves a movie by its ID
	FindByID(ctx context.Context, id shared.MovieID) (*Movie, error)

	// FindByCriteria retrieves movies based on search criteria
	FindByCriteria(ctx context.Context, criteria SearchCriteria) ([]*Movie, error)

	// FindByTitle searches movies by title (partial match)
	FindByTitle(ctx context.Context, title string) ([]*Movie, error)

	// FindByDirector retrieves movies by director
	FindByDirector(ctx context.Context, director string) ([]*Movie, error)

	// FindByGenre retrieves movies that have a specific genre
	FindByGenre(ctx context.Context, genre string) ([]*Movie, error)

	// FindTopRated retrieves top-rated movies
	FindTopRated(ctx context.Context, limit int) ([]*Movie, error)

	// CountAll returns the total number of movies
	CountAll(ctx context.Context) (int, error)
}

// MovieWriter defines write operations for movies
type MovieWriter interface {
	// Save persists a movie (insert or update)
	Save(ctx context.Context, movie *Movie) error

	// Delete removes a movie by ID
	Delete(ctx context.Context, id shared.MovieID) error

	// DeleteAll removes all movies (for testing)
	DeleteAll(ctx context.Context) error
}

// SearchCriteria represents search parameters for movies
type SearchCriteria struct {
	Title     string
	Director  string
	Genre     string
	MinYear   int
	MaxYear   int
	MinRating float64
	MaxRating float64
	Limit     int
	Offset    int
	OrderBy   OrderBy
	OrderDir  OrderDirection
}

// OrderBy represents fields that can be used for ordering
type OrderBy string

const (
	OrderByTitle     OrderBy = "title"
	OrderByDirector  OrderBy = "director"
	OrderByYear      OrderBy = "year"
	OrderByRating    OrderBy = "rating"
	OrderByCreatedAt OrderBy = "created_at"
	OrderByUpdatedAt OrderBy = "updated_at"
)

// OrderDirection represents sort direction
type OrderDirection string

const (
	OrderAsc  OrderDirection = "asc"
	OrderDesc OrderDirection = "desc"
)

// NewSearchCriteria creates a new SearchCriteria with default values
func NewSearchCriteria() SearchCriteria {
	return SearchCriteria{
		Limit:    50,
		Offset:   0,
		OrderBy:  OrderByTitle,
		OrderDir: OrderAsc,
	}
}
