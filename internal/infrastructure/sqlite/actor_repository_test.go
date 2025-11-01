package sqlite

import (
	"context"
	"database/sql"
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
