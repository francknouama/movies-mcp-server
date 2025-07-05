package movie

import (
	"testing"

	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
)

func TestMovieDomainEvents(t *testing.T) {
	// Create a movie
	movie, err := NewMovie("Test Movie", "Test Director", 2023)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Verify no events for new movie with temporary ID
	if movie.HasUncommittedEvents() {
		t.Error("New movie with temporary ID should not have domain events")
	}

	// Set a real ID to trigger events
	realID, err := shared.NewMovieID(42)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}
	movie.SetID(realID)

	// Now test rating change event
	initialEventCount := len(movie.UncommittedEvents())
	err = movie.SetRating(8.5)
	if err != nil {
		t.Fatalf("Failed to set rating: %v", err)
	}

	// Should have one more event
	if len(movie.UncommittedEvents()) != initialEventCount+1 {
		t.Errorf("Expected %d events after rating change, got %d", initialEventCount+1, len(movie.UncommittedEvents()))
	}

	// Test genre addition event
	eventCountBeforeGenre := len(movie.UncommittedEvents())
	err = movie.AddGenre("Action")
	if err != nil {
		t.Fatalf("Failed to add genre: %v", err)
	}

	// Should have one more event
	if len(movie.UncommittedEvents()) != eventCountBeforeGenre+1 {
		t.Errorf("Expected %d events after genre addition, got %d", eventCountBeforeGenre+1, len(movie.UncommittedEvents()))
	}

	// Test poster URL change event
	eventCountBeforePoster := len(movie.UncommittedEvents())
	err = movie.SetPosterURL("https://example.com/poster.jpg")
	if err != nil {
		t.Fatalf("Failed to set poster URL: %v", err)
	}

	// Should have one more event
	if len(movie.UncommittedEvents()) != eventCountBeforePoster+1 {
		t.Errorf("Expected %d events after poster URL change, got %d", eventCountBeforePoster+1, len(movie.UncommittedEvents()))
	}

	// Verify we can get all events
	events := movie.UncommittedEvents()
	if len(events) == 0 {
		t.Error("Movie should have domain events")
	}

	// Test that setting the same rating doesn't create duplicate events
	eventCountBeforeDuplicate := len(movie.UncommittedEvents())
	err = movie.SetRating(8.5) // Same rating
	if err != nil {
		t.Fatalf("Failed to set rating: %v", err)
	}

	// Should have the same number of events (no duplicate)
	if len(movie.UncommittedEvents()) != eventCountBeforeDuplicate {
		t.Errorf("Expected %d events after setting same rating, got %d", eventCountBeforeDuplicate, len(movie.UncommittedEvents()))
	}
}

func TestMovieCreatedEvent(t *testing.T) {
	// Create a movie with a real ID to trigger creation event
	id, err := shared.NewMovieID(42)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	movie, err := NewMovieWithID(id, "Test Movie", "Test Director", 2023)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Should have a creation event
	events := movie.UncommittedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 creation event, got %d", len(events))
	}

	if len(events) > 0 {
		event := events[0]
		if event.EventType() != "MovieCreated" {
			t.Errorf("Expected MovieCreated event, got %s", event.EventType())
		}

		if event.AggregateType() != "Movie" {
			t.Errorf("Expected Movie aggregate type, got %s", event.AggregateType())
		}

		if event.AggregateID() != "42" {
			t.Errorf("Expected aggregate ID '42', got %s", event.AggregateID())
		}
	}
}

func TestMovieEventTypes(t *testing.T) {
	id, _ := shared.NewMovieID(42)
	rating, _ := shared.NewRating(8.0)
	newRating, _ := shared.NewRating(9.0)

	// Test MovieCreatedEvent
	movie, _ := NewMovieWithID(id, "Test Movie", "Test Director", 2023)
	createdEvent := NewMovieCreatedEvent(movie, 1)
	if createdEvent.EventType() != "MovieCreated" {
		t.Errorf("Expected MovieCreated event type, got %s", createdEvent.EventType())
	}

	// Test MovieRatingChangedEvent
	ratingEvent := NewMovieRatingChangedEvent(id, rating, newRating, 2)
	if ratingEvent.EventType() != "MovieRatingChanged" {
		t.Errorf("Expected MovieRatingChanged event type, got %s", ratingEvent.EventType())
	}
	if ratingEvent.OldRating.Value() != 8.0 {
		t.Errorf("Expected old rating 8.0, got %f", ratingEvent.OldRating.Value())
	}
	if ratingEvent.NewRating.Value() != 9.0 {
		t.Errorf("Expected new rating 9.0, got %f", ratingEvent.NewRating.Value())
	}

	// Test MovieGenreAddedEvent
	genreEvent := NewMovieGenreAddedEvent(id, "Action", 3)
	if genreEvent.EventType() != "MovieGenreAdded" {
		t.Errorf("Expected MovieGenreAdded event type, got %s", genreEvent.EventType())
	}
	if genreEvent.Genre != "Action" {
		t.Errorf("Expected genre 'Action', got %s", genreEvent.Genre)
	}

	// Test MoviePosterChangedEvent
	posterEvent := NewMoviePosterChangedEvent(id, "old.jpg", "new.jpg", 4)
	if posterEvent.EventType() != "MoviePosterChanged" {
		t.Errorf("Expected MoviePosterChanged event type, got %s", posterEvent.EventType())
	}
	if posterEvent.OldPosterURL != "old.jpg" {
		t.Errorf("Expected old poster URL 'old.jpg', got %s", posterEvent.OldPosterURL)
	}
	if posterEvent.NewPosterURL != "new.jpg" {
		t.Errorf("Expected new poster URL 'new.jpg', got %s", posterEvent.NewPosterURL)
	}

	// Test MovieDeletedEvent
	deletedEvent := NewMovieDeletedEvent(id, "Test Movie", 5)
	if deletedEvent.EventType() != "MovieDeleted" {
		t.Errorf("Expected MovieDeleted event type, got %s", deletedEvent.EventType())
	}
	if deletedEvent.Title != "Test Movie" {
		t.Errorf("Expected title 'Test Movie', got %s", deletedEvent.Title)
	}
}
