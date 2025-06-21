package shared

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event that occurs within the system
type DomainEvent interface {
	// EventID returns the unique identifier for this event
	EventID() string

	// EventType returns the type of event (e.g., "MovieCreated")
	EventType() string

	// AggregateID returns the ID of the aggregate that generated this event
	AggregateID() string

	// AggregateType returns the type of aggregate (e.g., "Movie", "Actor")
	AggregateType() string

	// OccurredAt returns when the event occurred
	OccurredAt() time.Time

	// Version returns the aggregate version when this event was created
	Version() int
}

// BaseDomainEvent provides common functionality for all domain events
type BaseDomainEvent struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	version       int
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string, version int) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:       uuid.New().String(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		version:       version,
	}
}

// EventID returns the unique identifier for this event
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType returns the type of event
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the ID of the aggregate that generated this event
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the type of aggregate
func (e BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Version returns the aggregate version when this event was created
func (e BaseDomainEvent) Version() int {
	return e.version
}

// AggregateRoot provides base functionality for aggregate roots that generate domain events
type AggregateRoot struct {
	uncommittedEvents []DomainEvent
	version           int
}

// NewAggregateRoot creates a new aggregate root
func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{
		uncommittedEvents: make([]DomainEvent, 0),
		version:           0,
	}
}

// Version returns the current version of the aggregate
func (ar *AggregateRoot) Version() int {
	return ar.version
}

// AddEvent adds a domain event to the uncommitted events list
func (ar *AggregateRoot) AddEvent(event DomainEvent) {
	ar.uncommittedEvents = append(ar.uncommittedEvents, event)
	ar.version++
}

// UncommittedEvents returns all uncommitted domain events
func (ar *AggregateRoot) UncommittedEvents() []DomainEvent {
	// Return a copy to prevent external modification
	events := make([]DomainEvent, len(ar.uncommittedEvents))
	copy(events, ar.uncommittedEvents)
	return events
}

// MarkEventsAsCommitted clears the uncommitted events (called after successful persistence)
func (ar *AggregateRoot) MarkEventsAsCommitted() {
	ar.uncommittedEvents = make([]DomainEvent, 0)
}

// HasUncommittedEvents returns true if there are uncommitted events
func (ar *AggregateRoot) HasUncommittedEvents() bool {
	return len(ar.uncommittedEvents) > 0
}
