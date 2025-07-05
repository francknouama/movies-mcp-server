package actor

import (
	"context"

	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
)

// Repository defines the interface for actor data access
type Repository interface {
	ActorReader
	ActorWriter
}

// ActorReader defines read operations for actors
type ActorReader interface {
	// FindByID retrieves an actor by their ID
	FindByID(ctx context.Context, id shared.ActorID) (*Actor, error)

	// FindByCriteria retrieves actors based on search criteria
	FindByCriteria(ctx context.Context, criteria SearchCriteria) ([]*Actor, error)

	// FindByName searches actors by name (partial match)
	FindByName(ctx context.Context, name string) ([]*Actor, error)

	// FindByMovieID retrieves actors who appeared in a specific movie
	FindByMovieID(ctx context.Context, movieID shared.MovieID) ([]*Actor, error)

	// CountAll returns the total number of actors
	CountAll(ctx context.Context) (int, error)
}

// ActorWriter defines write operations for actors
type ActorWriter interface {
	// Save persists an actor (insert or update)
	Save(ctx context.Context, actor *Actor) error

	// Delete removes an actor by ID
	Delete(ctx context.Context, id shared.ActorID) error

	// DeleteAll removes all actors (for testing)
	DeleteAll(ctx context.Context) error
}

// SearchCriteria represents search parameters for actors
type SearchCriteria struct {
	Name         string
	MinBirthYear int
	MaxBirthYear int
	MovieID      shared.MovieID // Find actors who appeared in this movie
	Limit        int
	Offset       int
	OrderBy      OrderBy
	OrderDir     OrderDirection
}

// OrderBy represents fields that can be used for ordering
type OrderBy string

const (
	OrderByName      OrderBy = "name"
	OrderByBirthYear OrderBy = "birth_year"
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
		OrderBy:  OrderByName,
		OrderDir: OrderAsc,
	}
}
