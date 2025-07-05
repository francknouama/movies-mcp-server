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
	actors              map[int]*actor.Actor
	nextID              int
	findByIDFunc        func(ctx context.Context, id shared.ActorID) (*actor.Actor, error)
	saveFunc            func(ctx context.Context, a *actor.Actor) error
	deleteFunc          func(ctx context.Context, id shared.ActorID) error
	findByCriteriaFunc  func(ctx context.Context, criteria actor.SearchCriteria) ([]*actor.Actor, error)
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

	// Only assign a new ID for new actors (ID 0 is the temporary ID from NewActor)
	if actor.ID().Value() == 0 {
		// Assign new ID
		id, _ := shared.NewActorID(m.nextID)
		actor.SetID(id)
		m.nextID++
	}
	// For existing actors (ID > 0), don't change the ID

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
	if m.findByCriteriaFunc != nil {
		return m.findByCriteriaFunc(ctx, criteria)
	}

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

	// Apply offset
	if criteria.Offset > 0 && criteria.Offset < len(result) {
		result = result[criteria.Offset:]
	} else if criteria.Offset >= len(result) {
		result = []*actor.Actor{}
	}

	// Apply limit
	if criteria.Limit > 0 && criteria.Limit < len(result) {
		result = result[:criteria.Limit]
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

// Test UnlinkActorFromMovie - currently has 0% coverage
func TestService_UnlinkActorFromMovie(t *testing.T) {
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

	// Verify the link exists
	result, err := service.GetActor(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Failed to get actor: %v", err)
	}
	if len(result.MovieIDs) != 1 {
		t.Fatalf("Expected 1 movie ID, got %d", len(result.MovieIDs))
	}

	// Unlink from the movie
	err = service.UnlinkActorFromMovie(context.Background(), created.ID, movieID)
	if err != nil {
		t.Fatalf("UnlinkActorFromMovie() error = %v", err)
	}

	// Verify the link is removed
	result, err = service.GetActor(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Failed to get actor after unlinking: %v", err)
	}
	if len(result.MovieIDs) != 0 {
		t.Errorf("Expected 0 movie IDs after unlinking, got %d", len(result.MovieIDs))
	}
}

func TestService_UnlinkActorFromMovie_InvalidActorID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.UnlinkActorFromMovie(context.Background(), -1, 123)
	if err == nil {
		t.Error("Expected error for invalid actor ID")
	}
}

func TestService_UnlinkActorFromMovie_InvalidMovieID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.UnlinkActorFromMovie(context.Background(), 1, -1)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

func TestService_UnlinkActorFromMovie_ActorNotFound(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.UnlinkActorFromMovie(context.Background(), 999, 123)
	if err == nil {
		t.Error("Expected error for non-existent actor")
	}
}

// Additional edge cases for existing functions
func TestService_CreateActor_RepositoryError(t *testing.T) {
	repo := NewMockActorRepository()
	repo.saveFunc = func(ctx context.Context, a *actor.Actor) error {
		return errors.New("save failed")
	}
	service := NewService(repo)

	cmd := CreateActorCommand{
		Name:      "Test Actor",
		BirthYear: 1980,
	}

	_, err := service.CreateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error from repository save")
	}
}

func TestService_UpdateActor_InvalidID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	cmd := UpdateActorCommand{
		ID:        -1,
		Name:      "Test",
		BirthYear: 1980,
	}

	_, err := service.UpdateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for invalid actor ID")
	}
}

func TestService_UpdateActor_ActorNotFound(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	cmd := UpdateActorCommand{
		ID:        999,
		Name:      "Test",
		BirthYear: 1980,
	}

	_, err := service.UpdateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for non-existent actor")
	}
}

func TestService_UpdateActor_ValidationError(t *testing.T) {
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

	// Try to update with invalid data
	cmd := UpdateActorCommand{
		ID:        created.ID,
		Name:      "", // Empty name should fail
		BirthYear: 1980,
	}

	_, err = service.UpdateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error for empty name")
	}
}

func TestService_UpdateActor_RepositoryError(t *testing.T) {
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

	// Set up repository to fail on save
	repo.saveFunc = func(ctx context.Context, a *actor.Actor) error {
		return errors.New("save failed")
	}

	cmd := UpdateActorCommand{
		ID:        created.ID,
		Name:      "Updated Name",
		BirthYear: 1985,
	}

	_, err = service.UpdateActor(context.Background(), cmd)
	if err == nil {
		t.Error("Expected error from repository save")
	}
}

func TestService_DeleteActor_InvalidID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.DeleteActor(context.Background(), -1)
	if err == nil {
		t.Error("Expected error for invalid actor ID")
	}
}

func TestService_DeleteActor_RepositoryError(t *testing.T) {
	repo := NewMockActorRepository()
	repo.deleteFunc = func(ctx context.Context, id shared.ActorID) error {
		return errors.New("delete failed")
	}
	service := NewService(repo)

	err := service.DeleteActor(context.Background(), 1)
	if err == nil {
		t.Error("Expected error from repository delete")
	}
}

func TestService_LinkActorToMovie_InvalidActorID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.LinkActorToMovie(context.Background(), -1, 123)
	if err == nil {
		t.Error("Expected error for invalid actor ID")
	}
}

func TestService_LinkActorToMovie_InvalidMovieID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.LinkActorToMovie(context.Background(), 1, -1)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

func TestService_LinkActorToMovie_ActorNotFound(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	err := service.LinkActorToMovie(context.Background(), 999, 123)
	if err == nil {
		t.Error("Expected error for non-existent actor")
	}
}

func TestService_LinkActorToMovie_RepositoryError(t *testing.T) {
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

	// Set up repository to fail on save
	repo.saveFunc = func(ctx context.Context, a *actor.Actor) error {
		return errors.New("save failed")
	}

	err = service.LinkActorToMovie(context.Background(), created.ID, 123)
	if err == nil {
		t.Error("Expected error from repository save")
	}
}

func TestService_SearchActors_ExtendedCriteria(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create test actors
	actors := []CreateActorCommand{
		{Name: "Leonardo DiCaprio", BirthYear: 1974},
		{Name: "Matt Damon", BirthYear: 1970},
		{Name: "Tom Hanks", BirthYear: 1956},
		{Name: "Leonardo da Vinci", BirthYear: 1452}, // This should fail validation but let's test search
	}

	var createdActors []*ActorDTO
	for _, cmd := range actors {
		if cmd.BirthYear >= 1888 { // Only create valid actors
			created, err := service.CreateActor(context.Background(), cmd)
			if err != nil {
				t.Fatalf("Failed to create test actor: %v", err)
			}
			createdActors = append(createdActors, created)
		}
	}

	tests := []struct {
		name     string
		query    SearchActorsQuery
		expected int
	}{
		{
			name: "search by name",
			query: SearchActorsQuery{
				Name:  "Leonardo DiCaprio",
				Limit: 10,
			},
			expected: 1,
		},
		{
			name: "search with min birth year only",
			query: SearchActorsQuery{
				MinBirthYear: 1970,
				Limit:        10,
			},
			expected: 2, // Leonardo and Matt
		},
		{
			name: "search with max birth year only",
			query: SearchActorsQuery{
				MaxBirthYear: 1960,
				Limit:        10,
			},
			expected: 1, // Tom Hanks
		},
		{
			name: "search with limit",
			query: SearchActorsQuery{
				Limit: 2,
			},
			expected: 2,
		},
		{
			name: "search with offset",
			query: SearchActorsQuery{
				Limit:  10,
				Offset: 1,
			},
			expected: 2, // Should skip first result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchActors(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("SearchActors() error = %v", err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d actors, got %d", tt.expected, len(results))
			}
		})
	}
}

func TestService_GetActor_InvalidID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	_, err := service.GetActor(context.Background(), -1)
	if err == nil {
		t.Error("Expected error for invalid actor ID")
	}
}

func TestService_GetActorsByMovie_InvalidID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	_, err := service.GetActorsByMovie(context.Background(), -1)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

// Additional comprehensive tests for SearchActors to improve coverage
func TestService_SearchActors_ComprehensiveCoverage(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	// Create test actors with movies
	actors := []struct {
		cmd CreateActorCommand
		movieIDs []int
	}{
		{CreateActorCommand{Name: "Leonardo DiCaprio", BirthYear: 1974}, []int{1, 2}},
		{CreateActorCommand{Name: "Matt Damon", BirthYear: 1970}, []int{2, 3}},
		{CreateActorCommand{Name: "Tom Hanks", BirthYear: 1956}, []int{3}},
		{CreateActorCommand{Name: "Robert De Niro", BirthYear: 1943}, []int{4}},
	}

	for _, actorData := range actors {
		created, err := service.CreateActor(context.Background(), actorData.cmd)
		if err != nil {
			t.Fatalf("Failed to create test actor: %v", err)
		}

		// Link to movies
		for _, movieID := range actorData.movieIDs {
			err = service.LinkActorToMovie(context.Background(), created.ID, movieID)
			if err != nil {
				t.Fatalf("Failed to link actor to movie: %v", err)
			}
		}
	}

	tests := []struct {
		name     string
		query    SearchActorsQuery
		expected int
		desc     string
	}{
		{
			name: "empty query with default limit",
			query: SearchActorsQuery{},
			expected: 4,
			desc: "Should return all actors with default limit",
		},
		{
			name: "zero limit defaults to 50",
			query: SearchActorsQuery{Limit: 0},
			expected: 4,
			desc: "Should apply default limit of 50",
		},
		{
			name: "search by specific movie",
			query: SearchActorsQuery{MovieID: 2},
			expected: 2,
			desc: "Should find actors in movie 2 (Leonardo and Matt)",
		},
		{
			name: "search by movie with limit",
			query: SearchActorsQuery{MovieID: 2, Limit: 1},
			expected: 1,
			desc: "Should respect limit when searching by movie",
		},
		{
			name: "search by birth year range",
			query: SearchActorsQuery{MinBirthYear: 1970, MaxBirthYear: 1980, Limit: 10},
			expected: 2,
			desc: "Should find Leonardo and Matt born between 1970-1980",
		},
		{
			name: "search with min birth year only",
			query: SearchActorsQuery{MinBirthYear: 1960, Limit: 10},
			expected: 2,
			desc: "Should find actors born after 1960",
		},
		{
			name: "search with max birth year only",
			query: SearchActorsQuery{MaxBirthYear: 1950, Limit: 10},
			expected: 1,
			desc: "Should find actors born before 1950",
		},
		{
			name: "search with offset",
			query: SearchActorsQuery{Limit: 2, Offset: 2},
			expected: 2,
			desc: "Should skip first 2 actors and return next 2",
		},
		{
			name: "search with high offset",
			query: SearchActorsQuery{Limit: 10, Offset: 10},
			expected: 0,
			desc: "Should return empty when offset exceeds total",
		},
		{
			name: "search by name exact match",
			query: SearchActorsQuery{Name: "Tom Hanks", Limit: 10},
			expected: 1,
			desc: "Should find exact name match",
		},
		{
			name: "search by non-existent name",
			query: SearchActorsQuery{Name: "Non Existent", Limit: 10},
			expected: 0,
			desc: "Should return empty for non-existent name",
		},
		{
			name: "search by non-existent movie",
			query: SearchActorsQuery{MovieID: 999, Limit: 10},
			expected: 0,
			desc: "Should return empty for non-existent movie",
		},
		{
			name: "complex search with multiple criteria",
			query: SearchActorsQuery{
				MovieID: 2,
				MinBirthYear: 1970,
				MaxBirthYear: 1980,
				Limit: 10,
			},
			expected: 2,
			desc: "Should find actors matching all criteria",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchActors(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("SearchActors() error = %v", err)
			}

			if len(results) != tt.expected {
				t.Errorf("%s: Expected %d actors, got %d", tt.desc, tt.expected, len(results))
			}
		})
	}
}

func TestService_SearchActors_RepositoryError(t *testing.T) {
	repo := NewMockActorRepository()
	repo.findByCriteriaFunc = func(ctx context.Context, criteria actor.SearchCriteria) ([]*actor.Actor, error) {
		return nil, errors.New("repository error")
	}
	service := NewService(repo)

	query := SearchActorsQuery{Name: "Test", Limit: 10}
	_, err := service.SearchActors(context.Background(), query)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

func TestService_SearchActors_InvalidMovieID(t *testing.T) {
	repo := NewMockActorRepository()
	service := NewService(repo)

	query := SearchActorsQuery{MovieID: -1, Limit: 10}
	_, err := service.SearchActors(context.Background(), query)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}
