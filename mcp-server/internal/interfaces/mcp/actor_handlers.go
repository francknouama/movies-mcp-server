package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	actorApp "github.com/francknouama/movies-mcp-server/mcp-server/internal/application/actor"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

// ActorHandlers provides MCP handlers for actor operations
type ActorHandlers struct {
	actorService *actorApp.Service
}

// NewActorHandlers creates a new actor handlers instance
func NewActorHandlers(actorService *actorApp.Service) *ActorHandlers {
	return &ActorHandlers{
		actorService: actorService,
	}
}

// HandleAddActor handles the add_actor tool call
func (h *ActorHandlers) HandleAddActor(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseCreateActorRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid actor data", err.Error())
		return
	}

	// Convert to application command
	cmd := actorApp.CreateActorCommand{
		Name:      req.Name,
		BirthYear: req.BirthYear,
		Bio:       req.Bio,
	}

	// Create actor
	ctx := context.Background()
	actorDTO, err := h.actorService.CreateActor(ctx, cmd)
	if err != nil {
		sendError(id, dto.InvalidParams, "Failed to create actor", err.Error())
		return
	}

	// Convert to response format
	response := h.toActorResponse(actorDTO)
	sendResult(id, response)
}

// HandleGetActor handles getting an actor by ID
func (h *ActorHandlers) HandleGetActor(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse actor ID
	actorIDFloat, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "actor_id is required and must be a number", nil)
		return
	}
	actorID := int(actorIDFloat)

	// Get actor from service
	ctx := context.Background()
	actorDTO, err := h.actorService.GetActor(ctx, actorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to get actor", err.Error())
		}
		return
	}

	// Convert to response format
	response := h.toActorResponse(actorDTO)
	sendResult(id, response)
}

// HandleUpdateActor handles updating an actor
func (h *ActorHandlers) HandleUpdateActor(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseUpdateActorRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid actor data", err.Error())
		return
	}

	// Convert to application command
	cmd := actorApp.UpdateActorCommand{
		ID:        req.ID,
		Name:      req.Name,
		BirthYear: req.BirthYear,
		Bio:       req.Bio,
	}

	// Update actor
	ctx := context.Background()
	actorDTO, err := h.actorService.UpdateActor(ctx, cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor not found", nil)
		} else {
			sendError(id, dto.InvalidParams, "Failed to update actor", err.Error())
		}
		return
	}

	// Convert to response format
	response := h.toActorResponse(actorDTO)
	sendResult(id, response)
}

// HandleDeleteActor handles deleting an actor
func (h *ActorHandlers) HandleDeleteActor(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse actor ID
	actorIDFloat, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "actor_id is required and must be a number", nil)
		return
	}
	actorID := int(actorIDFloat)

	// Delete actor
	ctx := context.Background()
	err := h.actorService.DeleteActor(ctx, actorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to delete actor", err.Error())
		}
		return
	}

	sendResult(id, map[string]string{"message": "Actor deleted successfully"})
}

// HandleLinkActorToMovie handles linking an actor to a movie
func (h *ActorHandlers) HandleLinkActorToMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseLinkActorToMovieRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid link data", err.Error())
		return
	}

	// Link actor to movie
	ctx := context.Background()
	err = h.actorService.LinkActorToMovie(ctx, req.ActorID, req.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor or movie not found", nil)
		} else if strings.Contains(err.Error(), "already exists") {
			sendError(id, dto.InvalidParams, "Actor is already linked to this movie", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to link actor to movie", err.Error())
		}
		return
	}

	sendResult(id, map[string]string{"message": "Actor linked to movie successfully"})
}

// HandleUnlinkActorFromMovie handles unlinking an actor from a movie
func (h *ActorHandlers) HandleUnlinkActorFromMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseLinkActorToMovieRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid unlink data", err.Error())
		return
	}

	// Unlink actor from movie
	ctx := context.Background()
	err = h.actorService.UnlinkActorFromMovie(ctx, req.ActorID, req.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor, movie, or link not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to unlink actor from movie", err.Error())
		}
		return
	}

	sendResult(id, map[string]string{"message": "Actor unlinked from movie successfully"})
}

// HandleGetMovieCast handles getting the cast of a movie
func (h *ActorHandlers) HandleGetMovieCast(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse movie ID
	movieIDFloat, ok := arguments["movie_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "movie_id is required and must be a number", nil)
		return
	}
	movieID := int(movieIDFloat)

	// Get actors by movie
	ctx := context.Background()
	actorDTOs, err := h.actorService.GetActorsByMovie(ctx, movieID)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to get movie cast", err.Error())
		return
	}

	// Convert to response format
	response := &dto.ActorsListResponse{
		Actors:      make([]*dto.ActorResponse, len(actorDTOs)),
		Total:       len(actorDTOs),
		Description: fmt.Sprintf("Cast of movie %d", movieID),
	}

	for i, actorDTO := range actorDTOs {
		response.Actors[i] = h.toActorResponse(actorDTO)
	}

	sendResult(id, response)
}

// HandleGetActorMovies handles getting movies for an actor
func (h *ActorHandlers) HandleGetActorMovies(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse actor ID
	actorIDFloat, ok := arguments["actor_id"].(float64)
	if !ok {
		sendError(id, dto.InvalidParams, "actor_id is required and must be a number", nil)
		return
	}
	actorID := int(actorIDFloat)

	// Get actor to verify it exists and get movie IDs
	ctx := context.Background()
	actorDTO, err := h.actorService.GetActor(ctx, actorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, "Actor not found", nil)
		} else {
			sendError(id, dto.InternalError, "Failed to get actor", err.Error())
		}
		return
	}

	// Return the actor's movie IDs
	response := map[string]interface{}{
		"actor_id":     actorID,
		"actor_name":   actorDTO.Name,
		"movie_ids":    actorDTO.MovieIDs,
		"total_movies": len(actorDTO.MovieIDs),
	}

	sendResult(id, response)
}

// HandleSearchActors handles searching for actors
func (h *ActorHandlers) HandleSearchActors(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
	// Parse request
	req, err := h.parseSearchActorsRequest(arguments)
	if err != nil {
		sendError(id, dto.InvalidParams, "Invalid search parameters", err.Error())
		return
	}

	// Convert to application query
	query := actorApp.SearchActorsQuery{
		Name:         req.Name,
		MinBirthYear: req.MinBirthYear,
		MaxBirthYear: req.MaxBirthYear,
		MovieID:      req.MovieID,
		Limit:        req.Limit,
		Offset:       req.Offset,
		OrderBy:      req.OrderBy,
		OrderDir:     req.OrderDir,
	}

	// Set default limit
	if query.Limit == 0 {
		query.Limit = 20
	}

	// Search actors
	ctx := context.Background()
	actorDTOs, err := h.actorService.SearchActors(ctx, query)
	if err != nil {
		sendError(id, dto.InternalError, "Failed to search actors", err.Error())
		return
	}

	// Convert to response format
	response := &dto.ActorsListResponse{
		Actors:      make([]*dto.ActorResponse, len(actorDTOs)),
		Total:       len(actorDTOs),
		Description: "Search results",
	}

	for i, actorDTO := range actorDTOs {
		response.Actors[i] = h.toActorResponse(actorDTO)
	}

	sendResult(id, response)
}

// Utility methods

func (h *ActorHandlers) parseCreateActorRequest(arguments map[string]any) (*dto.CreateActorRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.CreateActorRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *ActorHandlers) parseUpdateActorRequest(arguments map[string]any) (*dto.UpdateActorRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.UpdateActorRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *ActorHandlers) parseLinkActorToMovieRequest(arguments map[string]any) (*dto.LinkActorToMovieRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.LinkActorToMovieRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *ActorHandlers) parseSearchActorsRequest(arguments map[string]any) (*dto.SearchActorsRequest, error) {
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	var req dto.SearchActorsRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *ActorHandlers) toActorResponse(actorDTO *actorApp.ActorDTO) *dto.ActorResponse {
	return &dto.ActorResponse{
		ID:        actorDTO.ID,
		Name:      actorDTO.Name,
		BirthYear: actorDTO.BirthYear,
		Bio:       actorDTO.Bio,
		MovieIDs:  actorDTO.MovieIDs,
		CreatedAt: actorDTO.CreatedAt,
		UpdatedAt: actorDTO.UpdatedAt,
	}
}
