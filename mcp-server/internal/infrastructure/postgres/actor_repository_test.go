package postgres

import (
	"context"
	"testing"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/actor"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/movie"
)

func createTestActor(t *testing.T) *actor.Actor {
	actor, err := actor.NewActor("Test Actor", 1980)
	if err != nil {
		t.Fatalf("Failed to create test actor: %v", err)
	}

	actor.SetBio("This is a test actor biography.")

	return actor
}

func TestActorRepository_Integration_Save_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewActorRepository(db)
	testActor := createTestActor(t)

	ctx := context.Background()
	err := repo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save actor: %v", err)
	}

	// Verify the actor was assigned an ID
	if testActor.ID().IsZero() {
		t.Error("Expected actor to be assigned an ID after save")
	}

	// Verify we can retrieve the actor
	retrieved, err := repo.FindByID(ctx, testActor.ID())
	if err != nil {
		t.Fatalf("Failed to retrieve saved actor: %v", err)
	}

	// Verify actor data
	if retrieved.Name() != testActor.Name() {
		t.Errorf("Expected name %s, got %s", testActor.Name(), retrieved.Name())
	}
	if retrieved.BirthYear().Value() != testActor.BirthYear().Value() {
		t.Errorf("Expected birth year %d, got %d", testActor.BirthYear().Value(), retrieved.BirthYear().Value())
	}
	if retrieved.Bio() != testActor.Bio() {
		t.Errorf("Expected bio %s, got %s", testActor.Bio(), retrieved.Bio())
	}
}

func TestActorRepository_Integration_Save_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewActorRepository(db)
	testActor := createTestActor(t)

	ctx := context.Background()

	// Save initially
	err := repo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save actor: %v", err)
	}

	originalID := testActor.ID()

	// Update the actor
	updatedActor, err := actor.NewActorWithID(originalID, "Updated Name", 1985)
	if err != nil {
		t.Fatalf("Failed to create updated actor: %v", err)
	}
	updatedActor.SetBio("Updated biography")

	// Save the update
	err = repo.Save(ctx, updatedActor)
	if err != nil {
		t.Fatalf("Failed to update actor: %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.FindByID(ctx, originalID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated actor: %v", err)
	}

	if retrieved.Name() != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got %s", retrieved.Name())
	}
	if retrieved.BirthYear().Value() != 1985 {
		t.Errorf("Expected updated birth year 1985, got %d", retrieved.BirthYear().Value())
	}
	if retrieved.Bio() != "Updated biography" {
		t.Errorf("Expected updated bio 'Updated biography', got %s", retrieved.Bio())
	}
}

func TestActorRepository_Integration_MovieRelationships(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie first
	testMovie, err := movie.NewMovie("Test Movie", "Test Director", 2023)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	err = movieRepo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save test movie: %v", err)
	}

	// Create actor with movie relationship
	testActor := createTestActor(t)
	err = testActor.AddMovie(testMovie.ID())
	if err != nil {
		t.Fatalf("Failed to add movie to actor: %v", err)
	}

	// Save actor
	err = actorRepo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save actor: %v", err)
	}

	// Retrieve and verify relationships
	retrieved, err := actorRepo.FindByID(ctx, testActor.ID())
	if err != nil {
		t.Fatalf("Failed to retrieve actor: %v", err)
	}

	movieIDs := retrieved.MovieIDs()
	if len(movieIDs) != 1 {
		t.Errorf("Expected 1 movie relationship, got %d", len(movieIDs))
	}

	if movieIDs[0].Value() != testMovie.ID().Value() {
		t.Errorf("Expected movie ID %d, got %d", testMovie.ID().Value(), movieIDs[0].Value())
	}

	// Test FindByMovieID
	actors, err := actorRepo.FindByMovieID(ctx, testMovie.ID())
	if err != nil {
		t.Fatalf("Failed to find actors by movie ID: %v", err)
	}

	if len(actors) != 1 {
		t.Errorf("Expected 1 actor for movie, got %d", len(actors))
	}

	if actors[0].ID().Value() != testActor.ID().Value() {
		t.Errorf("Expected actor ID %d, got %d", testActor.ID().Value(), actors[0].ID().Value())
	}
}

func TestActorRepository_Integration_FindByCriteria(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create test actors
	actors := []*actor.Actor{
		func() *actor.Actor {
			a, _ := actor.NewActor("Leonardo DiCaprio", 1974)
			a.SetBio("Academy Award winner")
			return a
		}(),
		func() *actor.Actor {
			a, _ := actor.NewActor("Matt Damon", 1970)
			return a
		}(),
		func() *actor.Actor {
			a, _ := actor.NewActor("Tom Hanks", 1956)
			return a
		}(),
	}

	// Save all actors
	for _, actor := range actors {
		if err := repo.Save(ctx, actor); err != nil {
			t.Fatalf("Failed to save actor %s: %v", actor.Name(), err)
		}
	}

	// Test search by name
	criteria := actor.SearchCriteria{
		Name:  "DiCaprio",
		Limit: 10,
	}
	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search by name: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 actor with name containing 'DiCaprio', got %d", len(results))
	}

	// Test search by birth year range
	criteria = actor.SearchCriteria{
		MinBirthYear: 1970,
		MaxBirthYear: 1980,
		Limit:        10,
	}
	results, err = repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search by birth year range: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors born between 1970-1980, got %d", len(results))
	}

	// Test ordering
	criteria = actor.SearchCriteria{
		OrderBy:  actor.OrderByBirthYear,
		OrderDir: actor.OrderDesc,
		Limit:    10,
	}
	results, err = repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("Failed to search with ordering: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 actors, got %d", len(results))
	}

	// Verify descending order by birth year
	if results[0].BirthYear().Value() < results[1].BirthYear().Value() {
		t.Error("Actors not in descending birth year order")
	}
}

func TestActorRepository_Integration_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewActorRepository(db)
	testActor := createTestActor(t)

	ctx := context.Background()

	// Save actor
	err := repo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save actor: %v", err)
	}

	actorID := testActor.ID()

	// Delete actor
	err = repo.Delete(ctx, actorID)
	if err != nil {
		t.Fatalf("Failed to delete actor: %v", err)
	}

	// Verify actor is gone
	_, err = repo.FindByID(ctx, actorID)
	if err == nil {
		t.Error("Expected error when finding deleted actor")
	}
}

func TestActorRepository_Integration_DeleteWithMovieRelationships(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie
	testMovie, err := movie.NewMovie("Test Movie", "Test Director", 2023)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	err = movieRepo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save test movie: %v", err)
	}

	// Create actor with movie relationship
	testActor := createTestActor(t)
	err = testActor.AddMovie(testMovie.ID())
	if err != nil {
		t.Fatalf("Failed to add movie to actor: %v", err)
	}

	// Save actor
	err = actorRepo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save actor: %v", err)
	}

	actorID := testActor.ID()

	// Delete actor (should also delete movie relationships)
	err = actorRepo.Delete(ctx, actorID)
	if err != nil {
		t.Fatalf("Failed to delete actor: %v", err)
	}

	// Verify actor is gone
	_, err = actorRepo.FindByID(ctx, actorID)
	if err == nil {
		t.Error("Expected error when finding deleted actor")
	}

	// Verify movie still exists but relationship is gone
	retrievedMovie, err := movieRepo.FindByID(ctx, testMovie.ID())
	if err != nil {
		t.Fatalf("Movie should still exist after actor deletion: %v", err)
	}

	if retrievedMovie.ID().Value() != testMovie.ID().Value() {
		t.Error("Movie should remain after actor deletion")
	}

	// Verify no actors found for the movie
	actors, err := actorRepo.FindByMovieID(ctx, testMovie.ID())
	if err != nil {
		t.Fatalf("Failed to search actors by movie: %v", err)
	}

	if len(actors) != 0 {
		t.Errorf("Expected 0 actors for movie after actor deletion, got %d", len(actors))
	}
}

func TestActorRepository_Integration_CountAll(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Check initial count
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("Failed to count actors: %v", err)
	}
	initialCount := count

	// Add some actors
	for i := 0; i < 3; i++ {
		actor, _ := actor.NewActor("Actor", 1980+i)
		if err := repo.Save(ctx, actor); err != nil {
			t.Fatalf("Failed to save actor: %v", err)
		}
	}

	// Check count again
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("Failed to count actors: %v", err)
	}

	expectedCount := initialCount + 3
	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}
}
