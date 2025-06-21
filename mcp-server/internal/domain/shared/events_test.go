package shared

import (
	"testing"
	"time"
)

func TestDomainEvent(t *testing.T) {
	// Test creating a base domain event
	event := NewBaseDomainEvent("TestEvent", "123", "TestAggregate", 1)
	
	// Verify event properties
	if event.EventType() != "TestEvent" {
		t.Errorf("Expected event type 'TestEvent', got '%s'", event.EventType())
	}
	
	if event.AggregateID() != "123" {
		t.Errorf("Expected aggregate ID '123', got '%s'", event.AggregateID())
	}
	
	if event.AggregateType() != "TestAggregate" {
		t.Errorf("Expected aggregate type 'TestAggregate', got '%s'", event.AggregateType())
	}
	
	if event.Version() != 1 {
		t.Errorf("Expected version 1, got %d", event.Version())
	}
	
	// Verify event ID is not empty
	if event.EventID() == "" {
		t.Error("Event ID should not be empty")
	}
	
	// Verify event occurred recently
	if time.Since(event.OccurredAt()) > time.Second {
		t.Error("Event should have occurred within the last second")
	}
}

func TestAggregateRoot(t *testing.T) {
	// Create a new aggregate root
	ar := NewAggregateRoot()
	
	// Verify initial state
	if ar.Version() != 0 {
		t.Errorf("Expected initial version 0, got %d", ar.Version())
	}
	
	if ar.HasUncommittedEvents() {
		t.Error("New aggregate root should not have uncommitted events")
	}
	
	if len(ar.UncommittedEvents()) != 0 {
		t.Errorf("Expected 0 uncommitted events, got %d", len(ar.UncommittedEvents()))
	}
	
	// Add an event
	event := NewBaseDomainEvent("TestEvent", "123", "TestAggregate", 1)
	ar.AddEvent(event)
	
	// Verify state after adding event
	if ar.Version() != 1 {
		t.Errorf("Expected version 1 after adding event, got %d", ar.Version())
	}
	
	if !ar.HasUncommittedEvents() {
		t.Error("Aggregate root should have uncommitted events after adding event")
	}
	
	if len(ar.UncommittedEvents()) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(ar.UncommittedEvents()))
	}
	
	// Add another event
	event2 := NewBaseDomainEvent("TestEvent2", "123", "TestAggregate", 2)
	ar.AddEvent(event2)
	
	// Verify state after adding second event
	if ar.Version() != 2 {
		t.Errorf("Expected version 2 after adding second event, got %d", ar.Version())
	}
	
	if len(ar.UncommittedEvents()) != 2 {
		t.Errorf("Expected 2 uncommitted events, got %d", len(ar.UncommittedEvents()))
	}
	
	// Mark events as committed
	ar.MarkEventsAsCommitted()
	
	// Verify state after committing events
	if ar.HasUncommittedEvents() {
		t.Error("Aggregate root should not have uncommitted events after committing")
	}
	
	if len(ar.UncommittedEvents()) != 0 {
		t.Errorf("Expected 0 uncommitted events after committing, got %d", len(ar.UncommittedEvents()))
	}
	
	// Version should remain the same after committing
	if ar.Version() != 2 {
		t.Errorf("Expected version 2 after committing, got %d", ar.Version())
	}
}