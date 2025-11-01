package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/francknouama/movies-mcp-server/internal/domain/actor"
	"github.com/francknouama/movies-mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/internal/domain/shared"
	_ "modernc.org/sqlite"
)

// setupActorTestDB creates an in-memory SQLite database for actor testing
func setupActorTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Add _time_format parameter to parse timestamps
	db, err := sql.Open("sqlite", ":memory:?_time_format=sqlite")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create schema for actors and movies (matching repository expectations)
	schema := `
	CREATE TABLE actors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		birth_year INTEGER,
		bio TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		director TEXT NOT NULL,
		year INTEGER NOT NULL,
		rating REAL,
		genre TEXT NOT NULL DEFAULT '[]',
		poster_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE movie_actors (
		movie_id INTEGER NOT NULL,
		actor_id INTEGER NOT NULL,
		role TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (movie_id, actor_id),
		FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
		FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	// Verify tables were created
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('actors', 'movies', 'movie_actors')").Scan(&tableCount)
	if err != nil {
		t.Fatalf("failed to verify tables: %v", err)
	}
	if tableCount != 3 {
		t.Fatalf("expected 3 tables to be created, got %d", tableCount)
	}

	return db
}

func TestActorRepository_Save_Insert(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create a new actor (no ID)
	domainActor, err := actor.NewActor("Tom Hanks", 1956)
	if err != nil {
		t.Fatalf("failed to create domain actor: %v", err)
	}

	// Save the actor
	err = repo.Save(ctx, domainActor)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify ID was assigned
	if domainActor.ID().IsZero() {
		t.Error("Expected actor to have ID assigned after save")
	}

	// Verify actor can be retrieved
	retrieved, err := repo.FindByID(ctx, domainActor.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if retrieved.Name() != "Tom Hanks" {
		t.Errorf("Expected name 'Tom Hanks', got %s", retrieved.Name())
	}
	if retrieved.BirthYear().Value() != 1956 {
		t.Errorf("Expected birth year 1956, got %d", retrieved.BirthYear().Value())
	}
}

func TestActorRepository_Save_Update(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create and save an actor
	domainActor, _ := actor.NewActor("Tom Hanks", 1956)
	_ = repo.Save(ctx, domainActor)

	// Update the actor's bio
	domainActor.SetBio("Academy Award winner known for Forrest Gump")

	// Save the update
	err := repo.Save(ctx, domainActor)
	if err != nil {
		t.Fatalf("Save() update error = %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.FindByID(ctx, domainActor.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if retrieved.Bio() != "Academy Award winner known for Forrest Gump" {
		t.Errorf("Expected bio to be updated, got %s", retrieved.Bio())
	}
}

func TestActorRepository_Save_WithMovies(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save movies first
	movie1, _ := movie.NewMovie("Forrest Gump", "Robert Zemeckis", 1994)
	_ = movieRepo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Cast Away", "Robert Zemeckis", 2000)
	_ = movieRepo.Save(ctx, movie2)

	// Create actor with movies
	domainActor, _ := actor.NewActor("Tom Hanks", 1956)
	_ = domainActor.AddMovie(movie1.ID())
	_ = domainActor.AddMovie(movie2.ID())

	// Save the actor
	err := actorRepo.Save(ctx, domainActor)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Retrieve and verify movies
	retrieved, err := actorRepo.FindByID(ctx, domainActor.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	movieIDs := retrieved.MovieIDs()
	if len(movieIDs) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(movieIDs))
	}
}

func TestActorRepository_FindByID_NotFound(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	actorID, _ := shared.NewActorID(999)
	_, err := repo.FindByID(ctx, actorID)
	if err == nil {
		t.Error("Expected error for non-existent actor")
	}
}

func TestActorRepository_FindByName(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create test actors
	actors := []struct {
		name      string
		birthYear int
	}{
		{"Tom Hanks", 1956},
		{"Tom Cruise", 1962},
		{"Brad Pitt", 1963},
	}

	for _, a := range actors {
		domainActor, _ := actor.NewActor(a.name, a.birthYear)
		_ = repo.Save(ctx, domainActor)
	}

	// Search by partial name (case-insensitive)
	results, err := repo.FindByName(ctx, "tom")
	if err != nil {
		t.Fatalf("FindByName() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors with 'tom' in name, got %d", len(results))
	}
}

func TestActorRepository_FindByMovieID(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie
	domainMovie, _ := movie.NewMovie("Inception", "Christopher Nolan", 2010)
	_ = movieRepo.Save(ctx, domainMovie)

	// Create actors for the movie
	actor1, _ := actor.NewActor("Leonardo DiCaprio", 1974)
	_ = actor1.AddMovie(domainMovie.ID())
	_ = actorRepo.Save(ctx, actor1)

	actor2, _ := actor.NewActor("Joseph Gordon-Levitt", 1981)
	_ = actor2.AddMovie(domainMovie.ID())
	_ = actorRepo.Save(ctx, actor2)

	// Create an actor not in the movie
	actor3, _ := actor.NewActor("Tom Hanks", 1956)
	_ = actorRepo.Save(ctx, actor3)

	// Search by movie
	results, err := actorRepo.FindByMovieID(ctx, domainMovie.ID())
	if err != nil {
		t.Fatalf("FindByMovieID() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors in the movie, got %d", len(results))
	}
}

func TestActorRepository_FindByCriteria_ByBirthYearRange(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create test actors with different birth years
	actors := []struct {
		name      string
		birthYear int
	}{
		{"Actor 1", 1950},
		{"Actor 2", 1960},
		{"Actor 3", 1970},
		{"Actor 4", 1980},
	}

	for _, a := range actors {
		domainActor, _ := actor.NewActor(a.name, a.birthYear)
		_ = repo.Save(ctx, domainActor)
	}

	// Search by birth year range using criteria
	criteria := actor.SearchCriteria{
		MinBirthYear: 1960,
		MaxBirthYear: 1975,
		Limit:        10,
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors born between 1960-1975, got %d", len(results))
	}
}

func TestActorRepository_CountAll(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Initially should be 0
	count, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Add 5 actors
	for i := 1; i <= 5; i++ {
		domainActor, _ := actor.NewActor("Actor", 1980)
		_ = repo.Save(ctx, domainActor)
	}

	// Should now be 5
	count, err = repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestActorRepository_Delete(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Create an actor
	domainActor, _ := actor.NewActor("Tom Hanks", 1956)
	_ = repo.Save(ctx, domainActor)

	// Verify it exists
	_, err := repo.FindByID(ctx, domainActor.ID())
	if err != nil {
		t.Fatalf("Actor should exist before delete")
	}

	// Delete the actor
	err = repo.Delete(ctx, domainActor.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's deleted
	_, err = repo.FindByID(ctx, domainActor.ID())
	if err == nil {
		t.Error("Expected error when finding deleted actor")
	}
}

func TestActorRepository_Delete_NotFound(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	actorID, _ := shared.NewActorID(999)
	err := repo.Delete(ctx, actorID)
	if err == nil {
		t.Error("Expected error when deleting non-existent actor")
	}
}

func TestActorRepository_DeleteAll(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Add multiple actors
	for i := 1; i <= 5; i++ {
		domainActor, _ := actor.NewActor("Actor", 1980)
		_ = repo.Save(ctx, domainActor)
	}

	// Verify count
	count, _ := repo.CountAll(ctx)
	if count != 5 {
		t.Errorf("Expected 5 actors before DeleteAll, got %d", count)
	}

	// Delete all
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify count is 0
	count, _ = repo.CountAll(ctx)
	if count != 0 {
		t.Errorf("Expected 0 actors after DeleteAll, got %d", count)
	}
}

func TestActorRepository_UpdateMovieRelationships(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies
	movie1, _ := movie.NewMovie("Movie 1", "Director 1", 2010)
	_ = movieRepo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Movie 2", "Director 2", 2015)
	_ = movieRepo.Save(ctx, movie2)

	movie3, _ := movie.NewMovie("Movie 3", "Director 3", 2020)
	_ = movieRepo.Save(ctx, movie3)

	// Create actor with initial movies
	domainActor, _ := actor.NewActor("Actor", 1980)
	_ = domainActor.AddMovie(movie1.ID())
	_ = domainActor.AddMovie(movie2.ID())
	_ = actorRepo.Save(ctx, domainActor)

	// Verify initial movies
	retrieved, _ := actorRepo.FindByID(ctx, domainActor.ID())
	if len(retrieved.MovieIDs()) != 2 {
		t.Errorf("Expected 2 initial movies, got %d", len(retrieved.MovieIDs()))
	}

	// Update relationships: remove movie1, keep movie2, add movie3
	_ = domainActor.RemoveMovie(movie1.ID())
	_ = domainActor.AddMovie(movie3.ID())
	_ = actorRepo.Save(ctx, domainActor)

	// Verify updated movies
	retrieved, _ = actorRepo.FindByID(ctx, domainActor.ID())
	movieIDs := retrieved.MovieIDs()
	if len(movieIDs) != 2 {
		t.Errorf("Expected 2 movies after update, got %d", len(movieIDs))
	}

	// Check that movie1 is gone and movie3 is added
	hasMovie1 := false
	hasMovie3 := false
	for _, movieID := range movieIDs {
		if movieID.Value() == movie1.ID().Value() {
			hasMovie1 = true
		}
		if movieID.Value() == movie3.ID().Value() {
			hasMovie3 = true
		}
	}

	if hasMovie1 {
		t.Error("Movie 1 should have been removed")
	}
	if !hasMovie3 {
		t.Error("Movie 3 should have been added")
	}
}

func TestActorRepository_DeleteCascade(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie and actor
	domainMovie, _ := movie.NewMovie("Test Movie", "Director", 2020)
	_ = movieRepo.Save(ctx, domainMovie)

	domainActor, _ := actor.NewActor("Test Actor", 1980)
	_ = domainActor.AddMovie(domainMovie.ID())
	_ = actorRepo.Save(ctx, domainActor)

	// Verify relationship exists
	actorsInMovie, _ := actorRepo.FindByMovieID(ctx, domainMovie.ID())
	if len(actorsInMovie) != 1 {
		t.Errorf("Expected 1 actor in movie before delete, got %d", len(actorsInMovie))
	}

	// Delete the actor
	_ = actorRepo.Delete(ctx, domainActor.ID())

	// Verify relationship is gone (cascade delete)
	actorsInMovie, _ = actorRepo.FindByMovieID(ctx, domainMovie.ID())
	if len(actorsInMovie) != 0 {
		t.Errorf("Expected 0 actors in movie after cascade delete, got %d", len(actorsInMovie))
	}
}

// Error scenario tests for better coverage

func TestActorRepository_Save_Insert_InvalidData(t *testing.T) {
	// Create actor with invalid data (empty name should fail domain validation)
	_, err := actor.NewActor("", 1990)
	if err == nil {
		t.Error("Expected error for invalid actor data")
	}
}

func TestActorRepository_FindByCriteria_NoResults(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Search with criteria that won't match any actors
	criteria := actor.SearchCriteria{
		Name: "NonExistent Actor Name That Should Not Match",
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 actors for non-existent name, got %d", len(results))
	}
}

func TestActorRepository_FindByCriteria_WithBirthYearOnly(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Insert test actors
	actor1, _ := actor.NewActor("Actor 1990", 1990)
	actor2, _ := actor.NewActor("Actor 1985", 1985)
	actorRepo.Save(ctx, actor1)
	actorRepo.Save(ctx, actor2)

	// Search by birth year range
	criteria := actor.SearchCriteria{
		MinBirthYear: 1989,
		MaxBirthYear: 1991,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 actor in birth year range, got %d", len(results))
	}

	if len(results) > 0 && results[0].Name() != "Actor 1990" {
		t.Errorf("Expected 'Actor 1990', got '%s'", results[0].Name())
	}
}

func TestActorRepository_DeleteAll_WithData(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Insert multiple test actors
	for i := 1; i <= 3; i++ {
		testActor, _ := actor.NewActor(fmt.Sprintf("Actor %d", i), 1980+i)
		err := actorRepo.Save(ctx, testActor)
		if err != nil {
			t.Fatalf("Failed to save test actor: %v", err)
		}
	}

	// Verify actors exist
	count, err := actorRepo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 3 {
		t.Fatalf("Expected 3 actors before DeleteAll, got %d", count)
	}

	// Delete all
	err = actorRepo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify all deleted
	count, err = actorRepo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 actors after DeleteAll, got %d", count)
	}
}

func TestActorRepository_Update_WithMovieRelationships(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create and save a movie
	testMovie, _ := movie.NewMovie("Test Movie", "Director", 2020)
	err := movieRepo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save test movie: %v", err)
	}

	// Create and save an actor with the movie
	testActor, _ := actor.NewActor("Test Actor", 1990)
	testActor.AddMovie(testMovie.ID())
	err = actorRepo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to save test actor: %v", err)
	}

	// Update the actor's bio
	testActor.SetBio("Updated bio")
	err = actorRepo.Save(ctx, testActor)
	if err != nil {
		t.Fatalf("Failed to update test actor: %v", err)
	}

	// Verify the update persisted
	retrieved, err := actorRepo.FindByID(ctx, testActor.ID())
	if err != nil {
		t.Fatalf("Failed to retrieve updated actor: %v", err)
	}

	if retrieved.Bio() != "Updated bio" {
		t.Errorf("Expected bio 'Updated bio', got '%s'", retrieved.Bio())
	}

	// Verify movie relationship is still there
	if len(retrieved.MovieIDs()) != 1 {
		t.Errorf("Expected 1 movie relationship after update, got %d", len(retrieved.MovieIDs()))
	}
}

func TestActorRepository_FindByCriteria_WithMovieID(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie
	testMovie, _ := movie.NewMovie("Test Movie", "Director", 2020)
	err := movieRepo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save movie: %v", err)
	}

	// Create actors with and without the movie
	actor1, _ := actor.NewActor("Actor in Movie", 1990)
	actor1.AddMovie(testMovie.ID())
	err = actorRepo.Save(ctx, actor1)
	if err != nil {
		t.Fatalf("Failed to save actor1: %v", err)
	}

	actor2, _ := actor.NewActor("Actor not in Movie", 1985)
	err = actorRepo.Save(ctx, actor2)
	if err != nil {
		t.Fatalf("Failed to save actor2: %v", err)
	}

	// Search by movie ID
	criteria := actor.SearchCriteria{
		MovieID: testMovie.ID(),
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should only find actor1
	if len(results) != 1 {
		t.Errorf("Expected 1 actor for movie, got %d", len(results))
	}

	if len(results) > 0 && results[0].Name() != "Actor in Movie" {
		t.Errorf("Expected 'Actor in Movie', got '%s'", results[0].Name())
	}
}

func TestActorRepository_FindByCriteria_WithName(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Create actors
	actor1, _ := actor.NewActor("John Smith", 1990)
	actor2, _ := actor.NewActor("Jane Doe", 1985)
	actorRepo.Save(ctx, actor1)
	actorRepo.Save(ctx, actor2)

	// Search by partial name
	criteria := actor.SearchCriteria{
		Name: "John",
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	// Should find John Smith
	if len(results) != 1 {
		t.Errorf("Expected 1 actor matching 'John', got %d", len(results))
	}

	if len(results) > 0 && results[0].Name() != "John Smith" {
		t.Errorf("Expected 'John Smith', got '%s'", results[0].Name())
	}
}

func TestActorRepository_DeleteAll_EmptyTable(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Delete from empty table (should succeed)
	err := actorRepo.DeleteAll(ctx)
	if err != nil {
		t.Errorf("DeleteAll() on empty table error = %v", err)
	}

	// Verify count is 0
	count, err := actorRepo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0 after DeleteAll on empty table, got %d", count)
	}
}

func TestActorRepository_FindByCriteria_OrderByBirthYear(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Create actors with different birth years
	actor1, _ := actor.NewActor("Actor 1990", 1990)
	actor2, _ := actor.NewActor("Actor 1985", 1985)
	actor3, _ := actor.NewActor("Actor 1995", 1995)
	actorRepo.Save(ctx, actor1)
	actorRepo.Save(ctx, actor2)
	actorRepo.Save(ctx, actor3)

	// Search ordered by birth year ascending
	criteria := actor.SearchCriteria{
		OrderBy:  actor.OrderByBirthYear,
		OrderDir: actor.OrderAsc,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 actors, got %d", len(results))
	}

	// Verify ordering
	if results[0].BirthYear().Value() != 1985 || results[1].BirthYear().Value() != 1990 || results[2].BirthYear().Value() != 1995 {
		t.Errorf("Actors not ordered by birth year ascending")
	}
}

func TestActorRepository_FindByCriteria_OrderDescending(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Create actors
	actor1, _ := actor.NewActor("Alice", 1990)
	actor2, _ := actor.NewActor("Bob", 1985)
	actorRepo.Save(ctx, actor1)
	actorRepo.Save(ctx, actor2)

	// Search ordered by name descending
	criteria := actor.SearchCriteria{
		OrderBy:  actor.OrderByName,
		OrderDir: actor.OrderDesc,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 actors, got %d", len(results))
	}

	// Verify descending order (Bob should come before Alice)
	if results[0].Name() != "Bob" {
		t.Errorf("Expected 'Bob' first in descending order, got '%s'", results[0].Name())
	}
}

func TestActorRepository_FindByCriteria_OrderByCreatedAt(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Create actors with slight time delays
	actor1, _ := actor.NewActor("First Actor", 1990)
	actorRepo.Save(ctx, actor1)

	actor2, _ := actor.NewActor("Second Actor", 1985)
	actorRepo.Save(ctx, actor2)

	actor3, _ := actor.NewActor("Third Actor", 1995)
	actorRepo.Save(ctx, actor3)

	// Search ordered by created_at ascending
	criteria := actor.SearchCriteria{
		OrderBy:  actor.OrderByCreatedAt,
		OrderDir: actor.OrderAsc,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 actors, got %d", len(results))
	}

	// Verify ordering by creation time
	if results[0].Name() != "First Actor" {
		t.Errorf("Expected 'First Actor' first, got '%s'", results[0].Name())
	}
	if results[2].Name() != "Third Actor" {
		t.Errorf("Expected 'Third Actor' last, got '%s'", results[2].Name())
	}
}

func TestActorRepository_FindByCriteria_OrderByUpdatedAt(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	ctx := context.Background()

	// Create actors
	actor1, _ := actor.NewActor("Actor One", 1990)
	actorRepo.Save(ctx, actor1)

	actor2, _ := actor.NewActor("Actor Two", 1985)
	actorRepo.Save(ctx, actor2)

	// Update actor1 to change its updated_at timestamp
	actor1.SetBio("Updated bio")
	actorRepo.Save(ctx, actor1)

	// Search ordered by updated_at descending (most recently updated first)
	criteria := actor.SearchCriteria{
		OrderBy:  actor.OrderByUpdatedAt,
		OrderDir: actor.OrderDesc,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 actors, got %d", len(results))
	}

	// Actor One should be first because it was updated most recently
	if results[0].Name() != "Actor One" {
		t.Errorf("Expected 'Actor One' first (most recently updated), got '%s'", results[0].Name())
	}
}

func TestActorRepository_Save_WithNullBirthYear(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	repo := NewActorRepository(db)
	ctx := context.Background()

	// Manually insert an actor with NULL birth_year
	_, err := db.Exec(`
		INSERT INTO actors (name, birth_year, bio, created_at, updated_at)
		VALUES (?, NULL, ?, datetime('now'), datetime('now'))
	`, "Actor Without Birth Year", "Test bio")
	if err != nil {
		t.Fatalf("Failed to insert actor with NULL birth year: %v", err)
	}

	// Find the actor
	criteria := actor.SearchCriteria{
		Name: "Actor Without Birth Year",
	}

	results, err := repo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 actor, got %d", len(results))
	}

	// Verify default birth year (1900) is used for NULL
	foundActor := results[0]
	if foundActor.BirthYear().Value() != 1900 {
		t.Errorf("Expected default birth year 1900 for NULL, got %d", foundActor.BirthYear().Value())
	}
}

func TestActorRepository_DeleteAll_WithMovieRelationships(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create a movie
	testMovie, _ := movie.NewMovie("Test Movie", "Test Director", 2020)
	err := movieRepo.Save(ctx, testMovie)
	if err != nil {
		t.Fatalf("Failed to save movie: %v", err)
	}

	// Create actors linked to the movie
	actor1, _ := actor.NewActor("Actor One", 1990)
	err = actorRepo.Save(ctx, actor1)
	if err != nil {
		t.Fatalf("Failed to save actor1: %v", err)
	}

	actor2, _ := actor.NewActor("Actor Two", 1985)
	err = actorRepo.Save(ctx, actor2)
	if err != nil {
		t.Fatalf("Failed to save actor2: %v", err)
	}

	// Link actors to movie
	actor1.AddMovie(testMovie.ID())
	actorRepo.Save(ctx, actor1)

	actor2.AddMovie(testMovie.ID())
	actorRepo.Save(ctx, actor2)

	// Verify movie relationships exist
	var relationshipCount int
	err = db.QueryRow("SELECT COUNT(*) FROM movie_actors").Scan(&relationshipCount)
	if err != nil {
		t.Fatalf("Failed to count movie relationships: %v", err)
	}
	if relationshipCount != 2 {
		t.Errorf("Expected 2 movie relationships, got %d", relationshipCount)
	}

	// Delete all actors
	err = actorRepo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}

	// Verify actors are deleted
	count, err := actorRepo.CountAll(ctx)
	if err != nil {
		t.Fatalf("CountAll() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after DeleteAll, got %d", count)
	}

	// Verify movie relationships are also deleted (CASCADE)
	err = db.QueryRow("SELECT COUNT(*) FROM movie_actors").Scan(&relationshipCount)
	if err != nil {
		t.Fatalf("Failed to count movie relationships after DeleteAll: %v", err)
	}
	if relationshipCount != 0 {
		t.Errorf("Expected 0 movie relationships after DeleteAll (CASCADE), got %d", relationshipCount)
	}
}

func TestActorRepository_ComplexQueryScenarios(t *testing.T) {
	db := setupActorTestDB(t)
	defer db.Close()

	actorRepo := NewActorRepository(db)
	movieRepo := NewMovieRepository(db)
	ctx := context.Background()

	// Create movies
	movie1, _ := movie.NewMovie("Action Movie", "Director A", 2020)
	movieRepo.Save(ctx, movie1)

	movie2, _ := movie.NewMovie("Drama Movie", "Director B", 2021)
	movieRepo.Save(ctx, movie2)

	// Create actors with different profiles
	actor1, _ := actor.NewActor("Young Actor", 2000)
	actor1.SetBio("Rising star")
	actor1.AddMovie(movie1.ID())
	actor1.AddMovie(movie2.ID())
	actorRepo.Save(ctx, actor1)

	actor2, _ := actor.NewActor("Veteran Actor", 1960)
	actor2.AddMovie(movie1.ID())
	actorRepo.Save(ctx, actor2)

	actor3, _ := actor.NewActor("Character Actor", 1975)
	actorRepo.Save(ctx, actor3)

	// Test 1: Find actors by birth year range
	criteria := actor.SearchCriteria{
		MinBirthYear: 1950,
		MaxBirthYear: 1980,
		OrderBy:      actor.OrderByBirthYear,
		OrderDir:     actor.OrderAsc,
	}

	results, err := actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		t.Fatalf("FindByCriteria() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors in birth year range, got %d", len(results))
	}

	// Test 2: Find actors in a specific movie
	results2, err := actorRepo.FindByMovieID(ctx, movie1.ID())
	if err != nil {
		t.Fatalf("FindByMovieID() error = %v", err)
	}

	if len(results2) != 2 {
		t.Errorf("Expected 2 actors in movie 1, got %d", len(results2))
	}

	// Test 3: Search with limit and offset
	criteria3 := actor.SearchCriteria{
		Limit:    1,
		Offset:   1,
		OrderBy:  actor.OrderByName,
		OrderDir: actor.OrderAsc,
	}

	results3, err := actorRepo.FindByCriteria(ctx, criteria3)
	if err != nil {
		t.Fatalf("FindByCriteria() with limit/offset error = %v", err)
	}

	if len(results3) != 1 {
		t.Errorf("Expected 1 actor with limit=1, got %d", len(results3))
	}

	// Test 4: Update actor to remove and add movie relationships
	actor1.RemoveMovie(movie2.ID())
	err = actorRepo.Save(ctx, actor1)
	if err != nil {
		t.Fatalf("Save() after removing movie error = %v", err)
	}

	// Verify the update
	updated, err := actorRepo.FindByID(ctx, actor1.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if len(updated.MovieIDs()) != 1 {
		t.Errorf("Expected 1 movie after removal, got %d", len(updated.MovieIDs()))
	}

	// Test 5: Search by exact name match
	criteria5 := actor.SearchCriteria{
		Name: "Young Actor",
	}

	results5, err := actorRepo.FindByCriteria(ctx, criteria5)
	if err != nil {
		t.Fatalf("FindByCriteria() by name error = %v", err)
	}

	if len(results5) != 1 {
		t.Errorf("Expected 1 actor with exact name match, got %d", len(results5))
	}

	if results5[0].Bio() != "Rising star" {
		t.Errorf("Expected bio 'Rising star', got '%s'", results5[0].Bio())
	}
}
