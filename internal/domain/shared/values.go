package shared

import (
	"errors"
	"time"
)

// MovieID represents a unique identifier for a movie
type MovieID struct {
	value int
}

// NewMovieID creates a new MovieID with validation
func NewMovieID(id int) (MovieID, error) {
	if id <= 0 {
		return MovieID{}, errors.New("movie ID must be positive")
	}
	return MovieID{value: id}, nil
}

// Value returns the underlying integer value
func (id MovieID) Value() int {
	return id.value
}

// IsZero returns true if this is a zero value
func (id MovieID) IsZero() bool {
	return id.value == 0
}

// String returns string representation
func (id MovieID) String() string {
	return "MovieID(" + string(rune(id.value)) + ")"
}

// ActorID represents a unique identifier for an actor
type ActorID struct {
	value int
}

// NewActorID creates a new ActorID with validation
func NewActorID(id int) (ActorID, error) {
	if id <= 0 {
		return ActorID{}, errors.New("actor ID must be positive")
	}
	return ActorID{value: id}, nil
}

// Value returns the underlying integer value
func (id ActorID) Value() int {
	return id.value
}

// IsZero returns true if this is a zero value
func (id ActorID) IsZero() bool {
	return id.value == 0
}

// Rating represents a movie rating between 0 and 10
type Rating struct {
	value float64
}

// NewRating creates a new Rating with validation
func NewRating(rating float64) (Rating, error) {
	if rating < 0 || rating > 10 {
		return Rating{}, errors.New("rating must be between 0 and 10")
	}
	return Rating{value: rating}, nil
}

// Value returns the underlying float64 value
func (r Rating) Value() float64 {
	return r.value
}

// IsZero returns true if this is a zero value
func (r Rating) IsZero() bool {
	return r.value == 0
}

// Year represents a movie release year
type Year struct {
	value int
}

// NewYear creates a new Year with validation
func NewYear(year int) (Year, error) {
	currentYear := time.Now().Year()
	// Allow movies from 1888 (first motion picture) to 15 years in the future
	if year < 1888 || year > currentYear+15 {
		return Year{}, errors.New("invalid year")
	}
	return Year{value: year}, nil
}

// Value returns the underlying integer value
func (y Year) Value() int {
	return y.value
}

// IsZero returns true if this is a zero value
func (y Year) IsZero() bool {
	return y.value == 0
}