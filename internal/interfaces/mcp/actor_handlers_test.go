package mcp

import (
	"context"
	"errors"
	"testing"

	actorApp "github.com/francknouama/movies-mcp-server/internal/application/actor"
	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

// ActorServiceInterface defines the interface needed by actor handlers
type ActorServiceInterface interface {
	CreateActor(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error)
	GetActor(ctx context.Context, id int) (*actorApp.ActorDTO, error)
	UpdateActor(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error)
	DeleteActor(ctx context.Context, id int) error
	LinkActorToMovie(ctx context.Context, actorID, movieID int) error
	UnlinkActorFromMovie(ctx context.Context, actorID, movieID int) error
	GetActorsByMovie(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error)
	SearchActors(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error)
}

// MockActorService implements the ActorServiceInterface for testing
type MockActorService struct {
	CreateFunc           func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error)
	GetByIDFunc          func(ctx context.Context, id int) (*actorApp.ActorDTO, error)
	UpdateFunc           func(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error)
	DeleteFunc           func(ctx context.Context, id int) error
	LinkToMovieFunc      func(ctx context.Context, actorID, movieID int) error
	UnlinkFromMovieFunc  func(ctx context.Context, actorID, movieID int) error
	GetActorsByMovieFunc func(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error)
	SearchActorsFunc     func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error)
}

func (m *MockActorService) CreateActor(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *MockActorService) GetActor(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockActorService) UpdateActor(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *MockActorService) DeleteActor(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockActorService) LinkActorToMovie(ctx context.Context, actorID, movieID int) error {
	if m.LinkToMovieFunc != nil {
		return m.LinkToMovieFunc(ctx, actorID, movieID)
	}
	return nil
}

func (m *MockActorService) UnlinkActorFromMovie(ctx context.Context, actorID, movieID int) error {
	if m.UnlinkFromMovieFunc != nil {
		return m.UnlinkFromMovieFunc(ctx, actorID, movieID)
	}
	return nil
}

func (m *MockActorService) GetActorsByMovie(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error) {
	if m.GetActorsByMovieFunc != nil {
		return m.GetActorsByMovieFunc(ctx, movieID)
	}
	return nil, nil
}

func (m *MockActorService) SearchActors(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
	if m.SearchActorsFunc != nil {
		return m.SearchActorsFunc(ctx, query)
	}
	return nil, nil
}

// ActorHandlersTestable is a version of ActorHandlers that accepts an interface
type ActorHandlersTestable struct {
	actorService ActorServiceInterface
}

func NewActorHandlersTestable(actorService ActorServiceInterface) *ActorHandlersTestable {
	return &ActorHandlersTestable{
		actorService: actorService,
	}
}

// HandleAddActor handles the add_actor tool call
func (h *ActorHandlersTestable) HandleAddActor(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	// Parse required parameters
	name, ok := arguments["name"].(string)
	if !ok || name == "" {
		sendError(id, dto.InvalidParams, "Invalid name parameter", nil)
		return
	}

	birthYear, ok := arguments["birth_year"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid birth_year parameter", nil)
		return
	}

	bio, _ := arguments["bio"].(string)

	// Convert to application command
	cmd := actorApp.CreateActorCommand{
		Name:      name,
		BirthYear: int(birthYear),
		Bio:       bio,
	}

	// Create actor
	ctx := context.Background()
	actorDTO, err := h.actorService.CreateActor(ctx, cmd)
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to create actor", err.Error())
		return
	}

	// Convert to response format
	response := map[string]interface{}{
		"id":         actorDTO.ID,
		"name":       actorDTO.Name,
		"birth_year": actorDTO.BirthYear,
		"bio":        actorDTO.Bio,
	}

	sendResult(id, response)
}

// HandleGetActor handles the get_actor tool call
func (h *ActorHandlersTestable) HandleGetActor(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	actorID, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid actor_id parameter", nil)
		return
	}

	ctx := context.Background()
	actorDTO, err := h.actorService.GetActor(ctx, int(actorID))
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to get actor", err.Error())
		return
	}

	response := map[string]interface{}{
		"id":         actorDTO.ID,
		"name":       actorDTO.Name,
		"birth_year": actorDTO.BirthYear,
		"bio":        actorDTO.Bio,
	}

	sendResult(id, response)
}

// HandleUpdateActor handles the update_actor tool call
func (h *ActorHandlersTestable) HandleUpdateActor(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	actorID, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid actor_id parameter", nil)
		return
	}

	name, _ := arguments["name"].(string)
	birthYear, _ := arguments["birth_year"].(float64)
	bio, _ := arguments["bio"].(string)

	cmd := actorApp.UpdateActorCommand{
		ID:        int(actorID),
		Name:      name,
		BirthYear: int(birthYear),
		Bio:       bio,
	}

	ctx := context.Background()
	actorDTO, err := h.actorService.UpdateActor(ctx, cmd)
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to update actor", err.Error())
		return
	}

	response := map[string]interface{}{
		"id":         actorDTO.ID,
		"name":       actorDTO.Name,
		"birth_year": actorDTO.BirthYear,
		"bio":        actorDTO.Bio,
	}

	sendResult(id, response)
}

// HandleDeleteActor handles the delete_actor tool call
func (h *ActorHandlersTestable) HandleDeleteActor(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	actorID, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid actor_id parameter", nil)
		return
	}

	ctx := context.Background()
	err := h.actorService.DeleteActor(ctx, int(actorID))
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to delete actor", err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "Actor deleted successfully",
	}

	sendResult(id, response)
}

// HandleLinkActorToMovie handles the link_actor_to_movie tool call
func (h *ActorHandlersTestable) HandleLinkActorToMovie(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	actorID, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid actor_id parameter", nil)
		return
	}

	movieID, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "Invalid movie_id parameter", nil)
		return
	}

	ctx := context.Background()
	err := h.actorService.LinkActorToMovie(ctx, int(actorID), int(movieID))
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to link actor to movie", err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "Actor linked to movie successfully",
	}

	sendResult(id, response)
}

// HandleSearchActors handles the search_actors tool call
func (h *ActorHandlersTestable) HandleSearchActors(id interface{}, arguments map[string]interface{}, sendResult func(interface{}, interface{}), sendError func(interface{}, int, string, interface{})) {
	name, _ := arguments["name"].(string)
	limit := 10
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	query := actorApp.SearchActorsQuery{
		Name:  name,
		Limit: limit,
	}

	ctx := context.Background()
	actorDTOs, err := h.actorService.SearchActors(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search actors", err.Error())
		return
	}

	var actors []interface{}
	for _, actorDTO := range actorDTOs {
		actors = append(actors, map[string]interface{}{
			"id":         actorDTO.ID,
			"name":       actorDTO.Name,
			"birth_year": actorDTO.BirthYear,
			"bio":        actorDTO.Bio,
		})
	}

	response := map[string]interface{}{
		"actors": actors,
		"total":  len(actors),
	}

	sendResult(id, response)
}

func TestActorHandlers_HandleAddActor(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
	}{
		{
			name: "successful actor creation",
			arguments: map[string]interface{}{
				"name":       "Test Actor",
				"birth_year": float64(1990),
				"bio":        "Test bio",
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					CreateFunc: func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
						return &actorApp.ActorDTO{
							ID:        1,
							Name:      cmd.Name,
							BirthYear: cmd.BirthYear,
							Bio:       cmd.Bio,
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				if resultMap["id"] != 1 {
					t.Errorf("Expected id 1, got %v", resultMap["id"])
				}
				if resultMap["name"] != "Test Actor" {
					t.Errorf("Expected name 'Test Actor', got %v", resultMap["name"])
				}
			},
		},
		{
			name:        "missing name parameter",
			arguments:   map[string]interface{}{"birth_year": float64(1990)},
			mockService: func() *MockActorService { return &MockActorService{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:        "invalid birth_year type",
			arguments:   map[string]interface{}{"name": "Test Actor", "birth_year": "not_a_number"},
			mockService: func() *MockActorService { return &MockActorService{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name: "service error",
			arguments: map[string]interface{}{
				"name":       "Test Actor",
				"birth_year": float64(1990),
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					CreateFunc: func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
						return nil, errors.New("service error")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name: "domain validation error",
			arguments: map[string]interface{}{
				"name":       "",
				"birth_year": float64(1990),
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					CreateFunc: func(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error) {
						return nil, errors.New("name cannot be empty")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleAddActor(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestActorHandlers_HandleGetActor(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
	}{
		{
			name:      "successful get actor",
			arguments: map[string]interface{}{"actor_id": float64(1)},
			mockService: func() *MockActorService {
				return &MockActorService{
					GetByIDFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
						return &actorApp.ActorDTO{
							ID:        id,
							Name:      "Test Actor",
							BirthYear: 1990,
							Bio:       "Test bio",
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				if resultMap["id"] != 1 {
					t.Errorf("Expected id 1, got %v", resultMap["id"])
				}
			},
		},
		{
			name:        "missing actor_id",
			arguments:   map[string]interface{}{},
			mockService: func() *MockActorService { return &MockActorService{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name:      "actor not found",
			arguments: map[string]interface{}{"actor_id": float64(999)},
			mockService: func() *MockActorService {
				return &MockActorService{
					GetByIDFunc: func(ctx context.Context, id int) (*actorApp.ActorDTO, error) {
						return nil, errors.New("actor not found")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleGetActor(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestActorHandlers_HandleDeleteActor(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
	}{
		{
			name:      "successful delete",
			arguments: map[string]interface{}{"actor_id": float64(1)},
			mockService: func() *MockActorService {
				return &MockActorService{
					DeleteFunc: func(ctx context.Context, id int) error {
						return nil
					},
				}
			},
			expectError: false,
		},
		{
			name:      "actor not found",
			arguments: map[string]interface{}{"actor_id": float64(999)},
			mockService: func() *MockActorService {
				return &MockActorService{
					DeleteFunc: func(ctx context.Context, id int) error {
						return errors.New("actor not found")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotError *dto.JSONRPCError
			var gotResult interface{}

			handlers.HandleDeleteActor(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				// Check for success message
				if resultMap, ok := gotResult.(map[string]interface{}); ok {
					if resultMap["message"] != "Actor deleted successfully" {
						t.Errorf("Unexpected success message: %v", resultMap["message"])
					}
				}
			}
		})
	}
}

func TestActorHandlers_HandleSearchActors(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
	}{
		{
			name: "successful search",
			arguments: map[string]interface{}{
				"name":  "Test",
				"limit": float64(10),
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					SearchActorsFunc: func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
						return []*actorApp.ActorDTO{
							{ID: 1, Name: "Test Actor 1", BirthYear: 1990},
							{ID: 2, Name: "Test Actor 2", BirthYear: 1991},
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				actors, ok := resultMap["actors"].([]interface{})
				if !ok {
					t.Fatalf("Expected actors to be an array")
				}
				if len(actors) != 2 {
					t.Errorf("Expected 2 actors, got %d", len(actors))
				}
			},
		},
		{
			name:      "empty search",
			arguments: map[string]interface{}{},
			mockService: func() *MockActorService {
				return &MockActorService{
					SearchActorsFunc: func(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error) {
						return []*actorApp.ActorDTO{}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				actors, _ := resultMap["actors"].([]interface{})
				if len(actors) != 0 {
					t.Errorf("Expected 0 actors, got %d", len(actors))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleSearchActors(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestActorHandlers_HandleUpdateActor(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
		checkResult func(t *testing.T, result interface{})
	}{
		{
			name: "successful actor update",
			arguments: map[string]interface{}{
				"actor_id":   float64(1),
				"name":       "Updated Actor",
				"birth_year": float64(1985),
				"bio":        "Updated bio",
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					UpdateFunc: func(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error) {
						return &actorApp.ActorDTO{
							ID:        cmd.ID,
							Name:      cmd.Name,
							BirthYear: cmd.BirthYear,
							Bio:       cmd.Bio,
						}, nil
					},
				}
			},
			expectError: false,
			checkResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected result to be a map")
				}
				if resultMap["name"] != "Updated Actor" {
					t.Errorf("Expected name 'Updated Actor', got %v", resultMap["name"])
				}
			},
		},
		{
			name:        "missing actor_id parameter",
			arguments:   map[string]interface{}{"name": "Updated Actor"},
			mockService: func() *MockActorService { return &MockActorService{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleUpdateActor(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, gotResult)
				}
			}
		})
	}
}

func TestActorHandlers_HandleLinkActorToMovie(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		mockService func() *MockActorService
		expectError bool
		errorCode   int
	}{
		{
			name: "successful link",
			arguments: map[string]interface{}{
				"actor_id": float64(1),
				"movie_id": float64(2),
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					LinkToMovieFunc: func(ctx context.Context, actorID, movieID int) error {
						return nil
					},
				}
			},
			expectError: false,
		},
		{
			name:        "missing actor_id",
			arguments:   map[string]interface{}{"movie_id": float64(2)},
			mockService: func() *MockActorService { return &MockActorService{} },
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
		{
			name: "service error",
			arguments: map[string]interface{}{
				"actor_id": float64(1),
				"movie_id": float64(2),
			},
			mockService: func() *MockActorService {
				return &MockActorService{
					LinkToMovieFunc: func(ctx context.Context, actorID, movieID int) error {
						return errors.New("link failed")
					},
				}
			},
			expectError: true,
			errorCode:   dto.InvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewActorHandlersTestable(tt.mockService())

			var gotResult interface{}
			var gotError *dto.JSONRPCError

			handlers.HandleLinkActorToMovie(
				1,
				tt.arguments,
				func(id interface{}, result interface{}) {
					gotResult = result
				},
				func(id interface{}, code int, message string, data interface{}) {
					gotError = &dto.JSONRPCError{
						Code:    code,
						Message: message,
						Data:    data,
					}
				},
			)

			if tt.expectError {
				if gotError == nil {
					t.Fatal("Expected error but got none")
				}
				if gotError.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, gotError.Code)
				}
			} else {
				if gotError != nil {
					t.Fatalf("Unexpected error: %v", gotError)
				}
				// Check for success message
				if resultMap, ok := gotResult.(map[string]interface{}); ok {
					if resultMap["message"] != "Actor linked to movie successfully" {
						t.Errorf("Unexpected success message: %v", resultMap["message"])
					}
				}
			}
		})
	}
}
