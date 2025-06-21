package actor

import (
	"testing"
	
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

func TestActorDomainEvents(t *testing.T) {
	// Create an actor
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Verify no events for new actor with temporary ID
	if actor.HasUncommittedEvents() {
		t.Error("New actor with temporary ID should not have domain events")
	}

	// Set a real ID to trigger events
	realID, err := shared.NewActorID(42)
	if err != nil {
		t.Fatalf("Failed to create actor ID: %v", err)
	}
	actor.SetID(realID)

	// Test bio change event
	initialEventCount := len(actor.UncommittedEvents())
	actor.SetBio("Test biography")

	// Should have one more event
	if len(actor.UncommittedEvents()) != initialEventCount+1 {
		t.Errorf("Expected %d events after bio change, got %d", initialEventCount+1, len(actor.UncommittedEvents()))
	}

	// Test movie linking event
	movieID, err := shared.NewMovieID(100)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	eventCountBeforeLink := len(actor.UncommittedEvents())
	err = actor.AddMovie(movieID)
	if err != nil {
		t.Fatalf("Failed to add movie: %v", err)
	}

	// Should have one more event
	if len(actor.UncommittedEvents()) != eventCountBeforeLink+1 {
		t.Errorf("Expected %d events after movie linking, got %d", eventCountBeforeLink+1, len(actor.UncommittedEvents()))
	}

	// Test movie unlinking event
	eventCountBeforeUnlink := len(actor.UncommittedEvents())
	err = actor.RemoveMovie(movieID)
	if err != nil {
		t.Fatalf("Failed to remove movie: %v", err)
	}

	// Should have one more event
	if len(actor.UncommittedEvents()) != eventCountBeforeUnlink+1 {
		t.Errorf("Expected %d events after movie unlinking, got %d", eventCountBeforeUnlink+1, len(actor.UncommittedEvents()))
	}

	// Test that setting the same bio doesn't create duplicate events
	eventCountBeforeDuplicate := len(actor.UncommittedEvents())
	actor.SetBio("Test biography") // Same bio

	// Should have the same number of events (no duplicate)
	if len(actor.UncommittedEvents()) != eventCountBeforeDuplicate {
		t.Errorf("Expected %d events after setting same bio, got %d", eventCountBeforeDuplicate, len(actor.UncommittedEvents()))
	}
}

func TestActorCreatedEvent(t *testing.T) {
	// Create an actor with a real ID to trigger creation event
	id, err := shared.NewActorID(42)
	if err != nil {
		t.Fatalf("Failed to create actor ID: %v", err)
	}

	actor, err := NewActorWithID(id, "Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Should have a creation event
	events := actor.UncommittedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 creation event, got %d", len(events))
	}

	if len(events) > 0 {
		event := events[0]
		if event.EventType() != "ActorCreated" {
			t.Errorf("Expected ActorCreated event, got %s", event.EventType())
		}
		
		if event.AggregateType() != "Actor" {
			t.Errorf("Expected Actor aggregate type, got %s", event.AggregateType())
		}
		
		if event.AggregateID() != "42" {
			t.Errorf("Expected aggregate ID '42', got %s", event.AggregateID())
		}
	}
}

func TestActorEventTypes(t *testing.T) {
	id, _ := shared.NewActorID(42)
	movieID, _ := shared.NewMovieID(100)

	// Test ActorCreatedEvent
	actor, _ := NewActorWithID(id, "Test Actor", 1980)
	createdEvent := NewActorCreatedEvent(actor, 1)
	if createdEvent.EventType() != "ActorCreated" {
		t.Errorf("Expected ActorCreated event type, got %s", createdEvent.EventType())
	}

	// Test ActorBioChangedEvent
	bioEvent := NewActorBioChangedEvent(id, "Old bio", "New bio", 2)
	if bioEvent.EventType() != "ActorBioChanged" {
		t.Errorf("Expected ActorBioChanged event type, got %s", bioEvent.EventType())
	}
	if bioEvent.OldBio != "Old bio" {
		t.Errorf("Expected old bio 'Old bio', got %s", bioEvent.OldBio)
	}
	if bioEvent.NewBio != "New bio" {
		t.Errorf("Expected new bio 'New bio', got %s", bioEvent.NewBio)
	}

	// Test ActorLinkedToMovieEvent
	linkedEvent := NewActorLinkedToMovieEvent(id, movieID, 3)
	if linkedEvent.EventType() != "ActorLinkedToMovie" {
		t.Errorf("Expected ActorLinkedToMovie event type, got %s", linkedEvent.EventType())
	}
	if linkedEvent.ActorID.Value() != 42 {
		t.Errorf("Expected actor ID 42, got %d", linkedEvent.ActorID.Value())
	}
	if linkedEvent.MovieID.Value() != 100 {
		t.Errorf("Expected movie ID 100, got %d", linkedEvent.MovieID.Value())
	}

	// Test ActorUnlinkedFromMovieEvent
	unlinkedEvent := NewActorUnlinkedFromMovieEvent(id, movieID, 4)
	if unlinkedEvent.EventType() != "ActorUnlinkedFromMovie" {
		t.Errorf("Expected ActorUnlinkedFromMovie event type, got %s", unlinkedEvent.EventType())
	}
	if unlinkedEvent.ActorID.Value() != 42 {
		t.Errorf("Expected actor ID 42, got %d", unlinkedEvent.ActorID.Value())
	}
	if unlinkedEvent.MovieID.Value() != 100 {
		t.Errorf("Expected movie ID 100, got %d", unlinkedEvent.MovieID.Value())
	}

	// Test ActorDeletedEvent
	deletedEvent := NewActorDeletedEvent(id, "Test Actor", 5)
	if deletedEvent.EventType() != "ActorDeleted" {
		t.Errorf("Expected ActorDeleted event type, got %s", deletedEvent.EventType())
	}
	if deletedEvent.Name != "Test Actor" {
		t.Errorf("Expected name 'Test Actor', got %s", deletedEvent.Name)
	}
}