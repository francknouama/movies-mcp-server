package actor

import (
	"errors"
	"strings"
	"time"

	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
)

// Actor represents an actor aggregate in the domain
type Actor struct {
	shared.AggregateRoot
	id        shared.ActorID
	name      string
	birthYear shared.Year
	bio       string
	movieIDs  []shared.MovieID
	createdAt time.Time
	updatedAt time.Time
}

// NewActor creates a new Actor with validation
func NewActor(name string, birthYear int) (*Actor, error) {
	// Use zero ID for new actors - will be assigned by repository
	id, err := shared.NewActorID(0) // Zero ID indicates new actor
	if err != nil {
		return nil, err
	}

	return NewActorWithID(id, name, birthYear)
}

// NewActorWithID creates a new Actor with a specific ID (for repository reconstruction)
func NewActorWithID(id shared.ActorID, name string, birthYear int) (*Actor, error) {
	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	// For birth year, we're more restrictive than movie years
	// Actors should be born within reasonable human lifespans
	currentYear := time.Now().Year()
	if birthYear < 1850 || birthYear > currentYear {
		return nil, errors.New("invalid birth year")
	}

	actorYear, err := shared.NewYear(birthYear)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	actor := &Actor{
		AggregateRoot: shared.NewAggregateRoot(),
		id:            id,
		name:          strings.TrimSpace(name),
		birthYear:     actorYear,
		movieIDs:      make([]shared.MovieID, 0),
		createdAt:     now,
		updatedAt:     now,
	}

	// Emit domain event for new actor creation (only for existing actors being reconstructed)
	if !id.IsZero() { // Only emit for actors with real IDs (not new actors)
		event := NewActorCreatedEvent(actor, actor.Version()+1)
		actor.AddEvent(event)
	}

	return actor, nil
}

// ID returns the actor's unique identifier
func (a *Actor) ID() shared.ActorID {
	return a.id
}

// Name returns the actor's name
func (a *Actor) Name() string {
	return a.name
}

// BirthYear returns the actor's birth year
func (a *Actor) BirthYear() shared.Year {
	return a.birthYear
}

// Bio returns the actor's biography
func (a *Actor) Bio() string {
	return a.bio
}

// MovieIDs returns a copy of the actor's movie IDs
func (a *Actor) MovieIDs() []shared.MovieID {
	movieIDs := make([]shared.MovieID, len(a.movieIDs))
	copy(movieIDs, a.movieIDs)
	return movieIDs
}

// CreatedAt returns when the actor was created
func (a *Actor) CreatedAt() time.Time {
	return a.createdAt
}

// UpdatedAt returns when the actor was last updated
func (a *Actor) UpdatedAt() time.Time {
	return a.updatedAt
}

// SetBio sets the actor's biography
func (a *Actor) SetBio(bio string) {
	oldBio := a.bio

	// Emit domain event if bio actually changed
	if oldBio != bio {
		event := NewActorBioChangedEvent(a.id, oldBio, bio, a.Version()+1)
		a.AddEvent(event)
	}

	a.bio = bio
	a.touch()
}

// AddMovie adds a movie ID to the actor's filmography
func (a *Actor) AddMovie(movieID shared.MovieID) error {
	// Check for duplicates
	for _, id := range a.movieIDs {
		if id.Value() == movieID.Value() {
			return errors.New("movie already exists in actor's filmography")
		}
	}

	a.movieIDs = append(a.movieIDs, movieID)

	// Emit domain event for actor linked to movie
	event := NewActorLinkedToMovieEvent(a.id, movieID, a.Version()+1)
	a.AddEvent(event)

	a.touch()
	return nil
}

// RemoveMovie removes a movie ID from the actor's filmography
func (a *Actor) RemoveMovie(movieID shared.MovieID) error {
	for i, id := range a.movieIDs {
		if id.Value() == movieID.Value() {
			// Remove by slicing
			a.movieIDs = append(a.movieIDs[:i], a.movieIDs[i+1:]...)

			// Emit domain event for actor unlinked from movie
			event := NewActorUnlinkedFromMovieEvent(a.id, movieID, a.Version()+1)
			a.AddEvent(event)

			a.touch()
			return nil
		}
	}
	return errors.New("movie not found in actor's filmography")
}

// HasMovie checks if the actor has a specific movie in their filmography
func (a *Actor) HasMovie(movieID shared.MovieID) bool {
	for _, id := range a.movieIDs {
		if id.Value() == movieID.Value() {
			return true
		}
	}
	return false
}

// MovieCount returns the number of movies the actor has appeared in
func (a *Actor) MovieCount() int {
	return len(a.movieIDs)
}

// Validate performs comprehensive validation of the actor
func (a *Actor) Validate() error {
	if strings.TrimSpace(a.name) == "" {
		return errors.New("name cannot be empty")
	}
	if a.birthYear.IsZero() {
		return errors.New("birth year must be set")
	}
	// Bio is optional
	// MovieIDs are optional
	return nil
}

// touch updates the updatedAt timestamp
func (a *Actor) touch() {
	a.updatedAt = time.Now()
}

// SetID sets the actor's ID (used by repository when saving)
func (a *Actor) SetID(id shared.ActorID) {
	a.id = id
	a.touch()
}
