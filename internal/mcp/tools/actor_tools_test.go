package tools

import (
	"context"
	"errors"
	"testing"

	actorApp "github.com/francknouama/movies-mcp-server/internal/application/actor"
)

// MockActorService implements ActorService for testing
type MockActorService struct {
	CreateActorFunc        func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error)
	GetActorFunc           func(ctx context.Context, id int) (*actorApp.ActorDTO, error)
	UpdateActorFunc        func(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error)
	DeleteActorFunc        func(ctx context.Context, id int) error
	LinkActorToMovieFunc   func(ctx context.Context, actorID, movieID int) error
	UnlinkActorFromMovieFunc func(ctx context.Context, actorID, movieID int) error
	GetActorsByMovieFunc   func(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error)
	SearchActorsFunc       func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error)
}

func (m *MockActorService) CreateActor(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
	if m.CreateActorFunc != nil {
		return m.CreateActorFunc(ctx, cmd)
	}
	return nil, errors.New("CreateActorFunc not implemented")
}

func (m *MockActorService) GetActor(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
	if m.GetActorFunc != nil {
		return m.GetActorFunc(ctx, id)
	}
	return nil, errors.New("GetActorFunc not implemented")
}

func (m *MockActorService) UpdateActor(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error) {
	if m.UpdateActorFunc != nil {
		return m.UpdateActorFunc(ctx, cmd)
	}
	return nil, errors.New("UpdateActorFunc not implemented")
}

func (m *MockActorService) DeleteActor(ctx context.Context, id int) error {
	if m.DeleteActorFunc != nil {
		return m.DeleteActorFunc(ctx, id)
	}
	return errors.New("DeleteActorFunc not implemented")
}

func (m *MockActorService) LinkActorToMovie(ctx context.Context, actorID, movieID int) error {
	if m.LinkActorToMovieFunc != nil {
		return m.LinkActorToMovieFunc(ctx, actorID, movieID)
	}
	return errors.New("LinkActorToMovieFunc not implemented")
}

func (m *MockActorService) UnlinkActorFromMovie(ctx context.Context, actorID, movieID int) error {
	if m.UnlinkActorFromMovieFunc != nil {
		return m.UnlinkActorFromMovieFunc(ctx, actorID, movieID)
	}
	return errors.New("UnlinkActorFromMovieFunc not implemented")
}

func (m *MockActorService) GetActorsByMovie(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error) {
	if m.GetActorsByMovieFunc != nil {
		return m.GetActorsByMovieFunc(ctx, movieID)
	}
	return nil, errors.New("GetActorsByMovieFunc not implemented")
}

func (m *MockActorService) SearchActors(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
	if m.SearchActorsFunc != nil {
		return m.SearchActorsFunc(ctx, query)
	}
	return nil, errors.New("SearchActorsFunc not implemented")
}

// ===== GetActor Tests =====

func TestGetActor_Success(t *testing.T) {
	mockService := &MockActorService{
		GetActorFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
			return &actorApp.ActorDTO{
				ID:        1,
				Name:      "Leonardo DiCaprio",
				BirthYear: 1974,
				Bio:       "Academy Award-winning actor",
				MovieIDs:  []int{1, 2, 3},
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.GetActor(context.Background(), nil, GetActorInput{ActorID: 1})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Name != "Leonardo DiCaprio" {
		t.Errorf("Expected name 'Leonardo DiCaprio', got: %s", output.Name)
	}

	if output.BirthYear != 1974 {
		t.Errorf("Expected birth year 1974, got: %d", output.BirthYear)
	}

	if len(output.MovieIDs) != 3 {
		t.Errorf("Expected 3 movie IDs, got: %d", len(output.MovieIDs))
	}
}

func TestGetActor_NotFound(t *testing.T) {
	mockService := &MockActorService{
		GetActorFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
			return nil, errors.New("actor not found")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.GetActor(context.Background(), nil, GetActorInput{ActorID: 999})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "actor not found" {
		t.Errorf("Expected 'actor not found' error, got: %v", err)
	}
}

// ===== AddActor Tests =====

func TestAddActor_Success(t *testing.T) {
	mockService := &MockActorService{
		CreateActorFunc: func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
			return &actorApp.ActorDTO{
				ID:        1,
				Name:      cmd.Name,
				BirthYear: cmd.BirthYear,
				Bio:       cmd.Bio,
				MovieIDs:  []int{},
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.AddActor(context.Background(), nil, AddActorInput{
		Name:      "Tom Hanks",
		BirthYear: 1956,
		Bio:       "Two-time Academy Award winner",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Name != "Tom Hanks" {
		t.Errorf("Expected name 'Tom Hanks', got: %s", output.Name)
	}

	if output.BirthYear != 1956 {
		t.Errorf("Expected birth year 1956, got: %d", output.BirthYear)
	}
}

func TestAddActor_ServiceError(t *testing.T) {
	mockService := &MockActorService{
		CreateActorFunc: func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.AddActor(context.Background(), nil, AddActorInput{
		Name: "Test Actor",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== UpdateActor Tests =====

func TestUpdateActor_Success(t *testing.T) {
	mockService := &MockActorService{
		UpdateActorFunc: func(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error) {
			return &actorApp.ActorDTO{
				ID:        cmd.ID,
				Name:      cmd.Name,
				BirthYear: cmd.BirthYear,
				Bio:       cmd.Bio,
				MovieIDs:  []int{1, 2},
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-02T00:00:00Z",
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.UpdateActor(context.Background(), nil, UpdateActorInput{
		ID:        1,
		Name:      "Leonardo Wilhelm DiCaprio",
		BirthYear: 1974,
		Bio:       "Updated biography",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Name != "Leonardo Wilhelm DiCaprio" {
		t.Errorf("Expected updated name, got: %s", output.Name)
	}
}

func TestUpdateActor_NotFound(t *testing.T) {
	mockService := &MockActorService{
		UpdateActorFunc: func(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error) {
			return nil, errors.New("actor not found")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.UpdateActor(context.Background(), nil, UpdateActorInput{
		ID:   999,
		Name: "Test",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== DeleteActor Tests =====

func TestDeleteActor_Success(t *testing.T) {
	mockService := &MockActorService{
		DeleteActorFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.DeleteActor(context.Background(), nil, DeleteActorInput{ActorID: 1})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Message != "Actor deleted successfully" {
		t.Errorf("Expected success message, got: %s", output.Message)
	}
}

func TestDeleteActor_NotFound(t *testing.T) {
	mockService := &MockActorService{
		DeleteActorFunc: func(ctx context.Context, id int) error {
			return errors.New("actor not found")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.DeleteActor(context.Background(), nil, DeleteActorInput{ActorID: 999})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== LinkActorToMovie Tests =====

func TestLinkActorToMovie_Success(t *testing.T) {
	mockService := &MockActorService{
		LinkActorToMovieFunc: func(ctx context.Context, actorID, movieID int) error {
			return nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.LinkActorToMovie(context.Background(), nil, LinkActorToMovieInput{
		ActorID: 1,
		MovieID: 42,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Message != "Actor linked to movie successfully" {
		t.Errorf("Expected success message, got: %s", output.Message)
	}
}

func TestLinkActorToMovie_ServiceError(t *testing.T) {
	mockService := &MockActorService{
		LinkActorToMovieFunc: func(ctx context.Context, actorID, movieID int) error {
			return errors.New("link failed")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.LinkActorToMovie(context.Background(), nil, LinkActorToMovieInput{
		ActorID: 1,
		MovieID: 999,
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== UnlinkActorFromMovie Tests =====

func TestUnlinkActorFromMovie_Success(t *testing.T) {
	mockService := &MockActorService{
		UnlinkActorFromMovieFunc: func(ctx context.Context, actorID, movieID int) error {
			return nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.UnlinkActorFromMovie(context.Background(), nil, UnlinkActorFromMovieInput{
		ActorID: 1,
		MovieID: 42,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Message != "Actor unlinked from movie successfully" {
		t.Errorf("Expected success message, got: %s", output.Message)
	}
}

func TestUnlinkActorFromMovie_ServiceError(t *testing.T) {
	mockService := &MockActorService{
		UnlinkActorFromMovieFunc: func(ctx context.Context, actorID, movieID int) error {
			return errors.New("unlink failed")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.UnlinkActorFromMovie(context.Background(), nil, UnlinkActorFromMovieInput{
		ActorID: 1,
		MovieID: 999,
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== GetMovieCast Tests =====

func TestGetMovieCast_Success(t *testing.T) {
	mockService := &MockActorService{
		GetActorsByMovieFunc: func(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error) {
			return []*actorApp.ActorDTO{
				{
					ID:        1,
					Name:      "Keanu Reeves",
					BirthYear: 1964,
					Bio:       "Canadian actor",
					MovieIDs:  []int{42},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
				{
					ID:        2,
					Name:      "Laurence Fishburne",
					BirthYear: 1961,
					Bio:       "American actor",
					MovieIDs:  []int{42},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.GetMovieCast(context.Background(), nil, GetMovieCastInput{MovieID: 42})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Actors) != 2 {
		t.Errorf("Expected 2 actors, got: %d", len(output.Actors))
	}

	if output.Actors[0].Name != "Keanu Reeves" {
		t.Errorf("Expected first actor 'Keanu Reeves', got: %s", output.Actors[0].Name)
	}
}

func TestGetMovieCast_EmptyResult(t *testing.T) {
	mockService := &MockActorService{
		GetActorsByMovieFunc: func(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error) {
			return []*actorApp.ActorDTO{}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.GetMovieCast(context.Background(), nil, GetMovieCastInput{MovieID: 999})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Actors) != 0 {
		t.Errorf("Expected 0 actors, got: %d", len(output.Actors))
	}
}

// ===== GetActorMovies Tests =====

func TestGetActorMovies_Success(t *testing.T) {
	mockService := &MockActorService{
		GetActorFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
			return &actorApp.ActorDTO{
				ID:        1,
				Name:      "Tom Hanks",
				BirthYear: 1956,
				Bio:       "Actor",
				MovieIDs:  []int{1, 2, 3, 4, 5},
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.GetActorMovies(context.Background(), nil, GetActorMoviesInput{ActorID: 1})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.ActorName != "Tom Hanks" {
		t.Errorf("Expected actor name 'Tom Hanks', got: %s", output.ActorName)
	}

	if len(output.MovieIDs) != 5 {
		t.Errorf("Expected 5 movie IDs, got: %d", len(output.MovieIDs))
	}
}

func TestGetActorMovies_NotFound(t *testing.T) {
	mockService := &MockActorService{
		GetActorFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
			return nil, errors.New("actor not found")
		},
	}

	tools := NewActorTools(mockService)
	_, _, err := tools.GetActorMovies(context.Background(), nil, GetActorMoviesInput{ActorID: 999})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== SearchActors Tests =====

func TestSearchActors_Success(t *testing.T) {
	mockService := &MockActorService{
		SearchActorsFunc: func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
			return []*actorApp.ActorDTO{
				{
					ID:        1,
					Name:      "Leonardo DiCaprio",
					BirthYear: 1974,
					Bio:       "Actor",
					MovieIDs:  []int{1, 2},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
				{
					ID:        2,
					Name:      "Leonardo Wilhelm",
					BirthYear: 1980,
					Bio:       "Another actor",
					MovieIDs:  []int{3},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.SearchActors(context.Background(), nil, SearchActorsInput{
		Name: "Leonardo",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Actors) != 2 {
		t.Errorf("Expected 2 actors, got: %d", len(output.Actors))
	}

	if output.Actors[0].Name != "Leonardo DiCaprio" {
		t.Errorf("Expected first actor 'Leonardo DiCaprio', got: %s", output.Actors[0].Name)
	}
}

func TestSearchActors_WithBirthYearFilter(t *testing.T) {
	var capturedQuery actorApp.SearchActorsQuery

	mockService := &MockActorService{
		SearchActorsFunc: func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
			capturedQuery = query
			return []*actorApp.ActorDTO{
				{
					ID:        1,
					Name:      "Tom Hanks",
					BirthYear: 1956,
					Bio:       "Actor",
					MovieIDs:  []int{1},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
			}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.SearchActors(context.Background(), nil, SearchActorsInput{
		Name:         "Tom",
		MinBirthYear: 1950,
		MaxBirthYear: 1960,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if capturedQuery.MinBirthYear != 1950 {
		t.Errorf("Expected min birth year filter 1950, got: %d", capturedQuery.MinBirthYear)
	}

	if capturedQuery.MaxBirthYear != 1960 {
		t.Errorf("Expected max birth year filter 1960, got: %d", capturedQuery.MaxBirthYear)
	}

	if len(output.Actors) != 1 {
		t.Errorf("Expected 1 actor, got: %d", len(output.Actors))
	}
}

func TestSearchActors_EmptyResults(t *testing.T) {
	mockService := &MockActorService{
		SearchActorsFunc: func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
			return []*actorApp.ActorDTO{}, nil
		},
	}

	tools := NewActorTools(mockService)
	_, output, err := tools.SearchActors(context.Background(), nil, SearchActorsInput{
		Name: "NonexistentActor",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Actors) != 0 {
		t.Errorf("Expected 0 actors, got: %d", len(output.Actors))
	}
}
