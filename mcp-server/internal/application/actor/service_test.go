package actor

import (
	"context"
	"errors"
	"testing"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/actor"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// MockActorRepository implements actor.Repository for testing
type MockActorRepository struct {
	actors       map[int]*actor.Actor
	nextID       int
	findByIDFunc func(ctx context.Context, id shared.ActorID) (*actor.Actor, error)
	saveFunc     func(ctx context.Context, a *actor.Actor) error
	deleteFunc   func(ctx context.Context, id shared.ActorID) error
}

func NewMockActorRepository() *MockActorRepository {
	return &MockActorRepository{
		actors: make(map[int]*actor.Actor),
		nextID: 1,
	}
}

func (m *MockActorRepository) FindByID(ctx context.Context, id shared.ActorID) (*actor.Actor, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	if actor, exists := m.actors[id.Value()]; exists {
		return actor, nil
	}
	return nil, errors.New("actor not found")
}

func (m *MockActorRepository) Save(ctx context.Context, actor *actor.Actor) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, actor)
	}
	
	// Only assign a new ID for new actors (ID 1 is the temporary ID from NewActor)
	if actor.ID().Value() == 1 {
		// Assign new ID
		id, _ := shared.NewActorID(m.nextID)
		actor.SetID(id)
		m.nextID++
	}
	// For existing actors (ID > 1), don't change the ID
	
	m.actors[actor.ID().Value()] = actor
	return nil
}

func (m *MockActorRepository) Delete(ctx context.Context, id shared.ActorID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	
	if _, exists := m.actors[id.Value()]; !exists {
		return errors.New("actor not found")
	}
	delete(m.actors, id.Value())
	return nil
}

func (m *MockActorRepository) FindByCriteria(ctx context.Context, criteria actor.SearchCriteria) ([]*actor.Actor, error) {
	var result []*actor.Actor
	for _, actorItem := range m.actors {
		match := true
		
		// Filter by name
		if criteria.Name != "" && actorItem.Name() != criteria.Name {
			match = false
		}
		
		// Filter by birth year range
		if criteria.MinBirthYear > 0 && actorItem.BirthYear().Value() < criteria.MinBirthYear {
			match = false
		}
		if criteria.MaxBirthYear > 0 && actorItem.BirthYear().Value() > criteria.MaxBirthYear {
			match = false
		}
		
		// Filter by movie ID
		if !criteria.MovieID.IsZero() && !actorItem.HasMovie(criteria.MovieID) {
			match = false
		}
		
		if match {
			result = append(result, actorItem)
		}
	}
	return result, nil
}

func (m *MockActorRepository) FindByName(ctx context.Context, name string) ([]*actor.Actor, error) {
	var result []*actor.Actor
	for _, actor := range m.actors {
		if actor.Name() == name {
			result = append(result, actor)
		}
	}
	return result, nil
}

func (m *MockActorRepository) FindByMovieID(ctx context.Context, movieID shared.MovieID) ([]*actor.Actor, error) {
	var result []*actor.Actor
	seen := make(map[int]bool)
	
	for _, actor := range m.actors {
		if actor.HasMovie(movieID) && !seen[actor.ID().Value()] {
			result = append(result, actor)
			seen[actor.ID().Value()] = true
		}
	}
	return result, nil
}

func (m *MockActorRepository) CountAll(ctx context.Context) (int, error) {
	return len(m.actors), nil
}

func (m *MockActorRepository) DeleteAll(ctx context.Context) error {
	m.actors = make(map[int]*actor.Actor)
	return nil
}

func TestService_CreateActor(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	tests := []struct {
		name    string
		cmd     CreateActorCommand
		wantErr bool
	}{
		{
			name: "valid actor",
			cmd: CreateActorCommand{
				Name:      "Leonardo DiCaprio",
				BirthYear: 1974,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			cmd: CreateActorCommand{
				Name:      "",
				BirthYear: 1974,
			},
			wantErr: true,
		},
		{
			name: "invalid birth year",
			cmd: CreateActorCommand{
				Name:      "Test Actor",
				BirthYear: 1800,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateActor(context.Background(), tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Name != tt.cmd.Name {
					t.Errorf("CreateActor() name = %v, want %v", result.Name, tt.cmd.Name)
				}
				if result.BirthYear != tt.cmd.BirthYear {
					t.Errorf("CreateActor() birthYear = %v, want %v", result.BirthYear, tt.cmd.BirthYear)
				}
			}
		})
	}
}

func TestService_CreateActor_WithBio(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	cmd := CreateActorCommand{
		Name:      "Leonardo DiCaprio",
		BirthYear: 1974,
		Bio:       "Academy Award-winning actor known for his versatile performances.",
	}

	result, err := service.CreateActor(context.Background(), cmd)
	if err != nil {
		t.Fatalf("CreateActor() error = %v", err)
	}

	if result.Bio != cmd.Bio {
		t.Errorf("CreateActor() bio = %v, want %v", result.Bio, cmd.Bio)
	}
}

func TestService_GetActor(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create an actor first
	createCmd := CreateActorCommand{
		Name:      "Test Actor",
		BirthYear: 1980,
	}
	created, err := service.CreateActor(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Get the actor
	result, err := service.GetActor(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetActor() error = %v", err)
	}

	if result.Name != created.Name {
		t.Errorf("GetActor() name = %v, want %v", result.Name, created.Name)
	}
}

func TestService_GetActor_NotFound(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	_, err := service.GetActor(context.Background(), 999)
	if err == nil {
		t.Error("Expected error for non-existent actor")
	}
}

func TestService_UpdateActor(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create an actor first
	createCmd := CreateActorCommand{
		Name:      "Original Name",
		BirthYear: 1980,
	}
	created, err := service.CreateActor(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Update the actor
	updateCmd := UpdateActorCommand{
		ID:        created.ID,
		Name:      "Updated Name",
		BirthYear: 1985,
		Bio:       "Updated biography",
	}

	result, err := service.UpdateActor(context.Background(), updateCmd)
	if err != nil {
		t.Fatalf("UpdateActor() error = %v", err)
	}

	if result.Name != updateCmd.Name {
		t.Errorf("UpdateActor() name = %v, want %v", result.Name, updateCmd.Name)
	}
	if result.BirthYear != updateCmd.BirthYear {
		t.Errorf("UpdateActor() birthYear = %v, want %v", result.BirthYear, updateCmd.BirthYear)
	}
	if result.Bio != updateCmd.Bio {
		t.Errorf("UpdateActor() bio = %v, want %v", result.Bio, updateCmd.Bio)
	}
}

func TestService_DeleteActor(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create an actor first
	createCmd := CreateActorCommand{
		Name:      "Test Actor",
		BirthYear: 1980,
	}
	created, err := service.CreateActor(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Delete the actor
	err = service.DeleteActor(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("DeleteActor() error = %v", err)
	}

	// Verify it's deleted
	_, err = service.GetActor(context.Background(), created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted actor")
	}
}

func TestService_LinkActorToMovie(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create an actor first
	createCmd := CreateActorCommand{
		Name:      "Test Actor",
		BirthYear: 1980,
	}
	created, err := service.CreateActor(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	// Link to a movie
	movieID := 123
	err = service.LinkActorToMovie(context.Background(), created.ID, movieID)
	if err != nil {
		t.Fatalf("LinkActorToMovie() error = %v", err)
	}

	// Verify the link
	result, err := service.GetActor(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Failed to get actor after linking: %v", err)
	}

	if len(result.MovieIDs) != 1 {
		t.Errorf("Expected 1 movie ID, got %d", len(result.MovieIDs))
	}

	if result.MovieIDs[0] != movieID {
		t.Errorf("Expected movie ID %d, got %d", movieID, result.MovieIDs[0])
	}
}

func TestService_SearchActors(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create test actors
	actors := []CreateActorCommand{
		{Name: "Leonardo DiCaprio", BirthYear: 1974},
		{Name: "Matt Damon", BirthYear: 1970},
		{Name: "Tom Hanks", BirthYear: 1956},
	}

	for _, cmd := range actors {
		_, err := service.CreateActor(context.Background(), cmd)
		if err != nil {
			t.Fatalf("Failed to create test actor: %v", err)
		}
	}

	// Search by birth year range
	query := SearchActorsQuery{
		MinBirthYear: 1970,
		MaxBirthYear: 1980,
		Limit:        10,
	}

	results, err := service.SearchActors(context.Background(), query)
	if err != nil {
		t.Fatalf("SearchActors() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 actors born between 1970-1980, got %d", len(results))
	}
}

func TestService_GetActorsByMovie(t *testing.T) {
	// Use a fresh repository to avoid interference from other tests
	freshRepo := NewMockActorRepository()
	freshService := NewService(freshRepo)

	// Create test actors
	actor1Cmd := CreateActorCommand{Name: "Actor 1", BirthYear: 1980}
	actor2Cmd := CreateActorCommand{Name: "Actor 2", BirthYear: 1985}

	actor1, err := freshService.CreateActor(context.Background(), actor1Cmd)
	if err != nil {
		t.Fatalf("Failed to create actor 1: %v", err)
	}

	actor2, err := freshService.CreateActor(context.Background(), actor2Cmd)
	if err != nil {
		t.Fatalf("Failed to create actor 2: %v", err)
	}

	// Link both actors to the same movie
	movieID := 123
	freshService.LinkActorToMovie(context.Background(), actor1.ID, movieID)
	freshService.LinkActorToMovie(context.Background(), actor2.ID, movieID)

	// Get actors by movie
	results, err := freshService.GetActorsByMovie(context.Background(), movieID)
	if err != nil {
		t.Fatalf("GetActorsByMovie() error = %v", err)
	}

	if len(results) != 2 {
		// Debug on failure
		count, _ := freshRepo.CountAll(context.Background())
		t.Logf("Total actors in repo: %d", count)
		for _, result := range results {
			t.Logf("Actor: %s (ID: %d)", result.Name, result.ID)
		}
		t.Errorf("Expected 2 actors for movie %d, got %d", movieID, len(results))
	}
}

func TestService_RepositoryError(t *testing.T) {
	repo := NewMockActorRepository()
	repo.saveFunc = func(ctx context.Context, a *actor.Actor) error {
		return errors.New("database error")
	}
	
	service := NewService(repo)

	cmd := CreateActorCommand{
		Name:      "Test Actor",
		BirthYear: 1980,
	}

	_, err := service.CreateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error from repository")
	}
}