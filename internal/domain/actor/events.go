package actor

import (
	"strconv"

	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
)

// Actor Domain Events

// ActorCreatedEvent is raised when a new actor is created
type ActorCreatedEvent struct {
	shared.BaseDomainEvent
	ActorID   shared.ActorID
	Name      string
	BirthYear shared.Year
	Bio       string
}

// NewActorCreatedEvent creates a new ActorCreatedEvent
func NewActorCreatedEvent(actor *Actor, version int) *ActorCreatedEvent {
	return &ActorCreatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorCreated",
			strconv.Itoa(actor.ID().Value()),
			"Actor",
			version,
		),
		ActorID:   actor.ID(),
		Name:      actor.Name(),
		BirthYear: actor.BirthYear(),
		Bio:       actor.Bio(),
	}
}

// ActorUpdatedEvent is raised when an actor is updated
type ActorUpdatedEvent struct {
	shared.BaseDomainEvent
	ActorID   shared.ActorID
	Name      string
	BirthYear shared.Year
	Bio       string
}

// NewActorUpdatedEvent creates a new ActorUpdatedEvent
func NewActorUpdatedEvent(actor *Actor, version int) *ActorUpdatedEvent {
	return &ActorUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorUpdated",
			strconv.Itoa(actor.ID().Value()),
			"Actor",
			version,
		),
		ActorID:   actor.ID(),
		Name:      actor.Name(),
		BirthYear: actor.BirthYear(),
		Bio:       actor.Bio(),
	}
}

// ActorDeletedEvent is raised when an actor is deleted
type ActorDeletedEvent struct {
	shared.BaseDomainEvent
	ActorID shared.ActorID
	Name    string
}

// NewActorDeletedEvent creates a new ActorDeletedEvent
func NewActorDeletedEvent(actorID shared.ActorID, name string, version int) *ActorDeletedEvent {
	return &ActorDeletedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorDeleted",
			strconv.Itoa(actorID.Value()),
			"Actor",
			version,
		),
		ActorID: actorID,
		Name:    name,
	}
}

// ActorLinkedToMovieEvent is raised when an actor is linked to a movie
type ActorLinkedToMovieEvent struct {
	shared.BaseDomainEvent
	ActorID shared.ActorID
	MovieID shared.MovieID
}

// NewActorLinkedToMovieEvent creates a new ActorLinkedToMovieEvent
func NewActorLinkedToMovieEvent(actorID shared.ActorID, movieID shared.MovieID, version int) *ActorLinkedToMovieEvent {
	return &ActorLinkedToMovieEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorLinkedToMovie",
			strconv.Itoa(actorID.Value()),
			"Actor",
			version,
		),
		ActorID: actorID,
		MovieID: movieID,
	}
}

// ActorUnlinkedFromMovieEvent is raised when an actor is unlinked from a movie
type ActorUnlinkedFromMovieEvent struct {
	shared.BaseDomainEvent
	ActorID shared.ActorID
	MovieID shared.MovieID
}

// NewActorUnlinkedFromMovieEvent creates a new ActorUnlinkedFromMovieEvent
func NewActorUnlinkedFromMovieEvent(actorID shared.ActorID, movieID shared.MovieID, version int) *ActorUnlinkedFromMovieEvent {
	return &ActorUnlinkedFromMovieEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorUnlinkedFromMovie",
			strconv.Itoa(actorID.Value()),
			"Actor",
			version,
		),
		ActorID: actorID,
		MovieID: movieID,
	}
}

// ActorBioChangedEvent is raised when an actor's biography is changed
type ActorBioChangedEvent struct {
	shared.BaseDomainEvent
	ActorID shared.ActorID
	OldBio  string
	NewBio  string
}

// NewActorBioChangedEvent creates a new ActorBioChangedEvent
func NewActorBioChangedEvent(actorID shared.ActorID, oldBio, newBio string, version int) *ActorBioChangedEvent {
	return &ActorBioChangedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"ActorBioChanged",
			strconv.Itoa(actorID.Value()),
			"Actor",
			version,
		),
		ActorID: actorID,
		OldBio:  oldBio,
		NewBio:  newBio,
	}
}
