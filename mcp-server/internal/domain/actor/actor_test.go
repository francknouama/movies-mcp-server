package actor

import (
	"testing"
	"time"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

func TestNewActor(t *testing.T) {
	tests := []struct {
		name      string
		actorName string
		birthYear int
		wantErr   bool
	}{
		{
			name:      "valid actor",
			actorName: "Leonardo DiCaprio",
			birthYear: 1974,
			wantErr:   false,
		},
		{
			name:      "empty name",
			actorName: "",
			birthYear: 1974,
			wantErr:   true,
		},
		{
			name:      "whitespace only name",
			actorName: "   ",
			birthYear: 1974,
			wantErr:   true,
		},
		{
			name:      "invalid birth year - too old",
			actorName: "Test Actor",
			birthYear: 1800,
			wantErr:   true,
		},
		{
			name:      "invalid birth year - future",
			actorName: "Test Actor",
			birthYear: 2030,
			wantErr:   true,
		},
		{
			name:      "valid recent birth year",
			actorName: "Young Actor",
			birthYear: 2005,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor, err := NewActor(tt.actorName, tt.birthYear)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if actor.Name() != tt.actorName {
					t.Errorf("NewActor() name = %v, want %v", actor.Name(), tt.actorName)
				}
				if actor.BirthYear().Value() != tt.birthYear {
					t.Errorf("NewActor() birthYear = %v, want %v", actor.BirthYear().Value(), tt.birthYear)
				}
				if actor.ID().IsZero() {
					t.Error("NewActor() should not have zero ID after creation")
				}
			}
		})
	}
}

func TestActor_AddMovie(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	movieID, err := shared.NewMovieID(123)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	// Add movie
	err = actor.AddMovie(movieID)
	if err != nil {
		t.Errorf("AddMovie() error = %v", err)
	}

	// Check if movie was added
	movies := actor.MovieIDs()
	found := false
	for _, id := range movies {
		if id.Value() == movieID.Value() {
			found = true
			break
		}
	}
	if !found {
		t.Error("Movie ID not found in actor's movie list")
	}
}

func TestActor_AddMovie_Duplicate(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	movieID, err := shared.NewMovieID(123)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	// Add movie first time
	err = actor.AddMovie(movieID)
	if err != nil {
		t.Fatalf("Failed to add movie first time: %v", err)
	}

	// Try to add same movie again
	err = actor.AddMovie(movieID)
	if err == nil {
		t.Error("Expected error when adding duplicate movie")
	}

	// Verify only one instance exists
	movies := actor.MovieIDs()
	count := 0
	for _, id := range movies {
		if id.Value() == movieID.Value() {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected 1 instance of movie ID, got %d", count)
	}
}

func TestActor_RemoveMovie(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	movieID, err := shared.NewMovieID(123)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	// Add movie first
	err = actor.AddMovie(movieID)
	if err != nil {
		t.Fatalf("Failed to add movie: %v", err)
	}

	// Remove movie
	err = actor.RemoveMovie(movieID)
	if err != nil {
		t.Errorf("RemoveMovie() error = %v", err)
	}

	// Check if movie was removed
	movies := actor.MovieIDs()
	for _, id := range movies {
		if id.Value() == movieID.Value() {
			t.Error("Movie ID should have been removed")
		}
	}
}

func TestActor_RemoveMovie_NotFound(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	movieID, err := shared.NewMovieID(123)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	// Try to remove movie that was never added
	err = actor.RemoveMovie(movieID)
	if err == nil {
		t.Error("Expected error when removing non-existent movie")
	}
}

func TestActor_HasMovie(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	movieID1, _ := shared.NewMovieID(123)
	movieID2, _ := shared.NewMovieID(456)

	// Add one movie
	actor.AddMovie(movieID1)

	// Test HasMovie
	if !actor.HasMovie(movieID1) {
		t.Error("Expected actor to have movie ID 123")
	}

	if actor.HasMovie(movieID2) {
		t.Error("Expected actor to not have movie ID 456")
	}
}

func TestActor_SetBio(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	bio := "This is a test actor biography."
	actor.SetBio(bio)

	if actor.Bio() != bio {
		t.Errorf("Expected bio %v, got %v", bio, actor.Bio())
	}

	// Test setting empty bio
	actor.SetBio("")
	if actor.Bio() != "" {
		t.Error("Expected empty bio")
	}
}

func TestActor_UpdateTimestamp(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	originalUpdatedAt := actor.UpdatedAt()

	// Sleep to ensure timestamp changes
	time.Sleep(1 * time.Millisecond)

	// Modify actor to trigger timestamp update
	actor.SetBio("New bio")

	if !actor.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated after modification")
	}
}

func TestActor_Validation(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	// Valid actor should pass validation
	if err := actor.Validate(); err != nil {
		t.Errorf("Valid actor should pass validation, got: %v", err)
	}

	// Test validation with movies and bio
	movieID, _ := shared.NewMovieID(123)
	actor.AddMovie(movieID)
	actor.SetBio("Test biography")

	if err := actor.Validate(); err != nil {
		t.Errorf("Actor with movies and bio should pass validation, got: %v", err)
	}
}

func TestNewActorWithID(t *testing.T) {
	id, err := shared.NewActorID(456)
	if err != nil {
		t.Fatalf("Failed to create actor ID: %v", err)
	}

	actor, err := NewActorWithID(id, "Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create actor with ID: %v", err)
	}

	if actor.ID().Value() != 456 {
		t.Errorf("Expected actor ID 456, got %d", actor.ID().Value())
	}
}

func TestActor_MovieCount(t *testing.T) {
	actor, err := NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	if actor.MovieCount() != 0 {
		t.Errorf("Expected 0 movies, got %d", actor.MovieCount())
	}

	// Add movies
	movieID1, _ := shared.NewMovieID(123)
	movieID2, _ := shared.NewMovieID(456)

	actor.AddMovie(movieID1)
	if actor.MovieCount() != 1 {
		t.Errorf("Expected 1 movie, got %d", actor.MovieCount())
	}

	actor.AddMovie(movieID2)
	if actor.MovieCount() != 2 {
		t.Errorf("Expected 2 movies, got %d", actor.MovieCount())
	}
}
