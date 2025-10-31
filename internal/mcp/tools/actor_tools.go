package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	actorApp "github.com/francknouama/movies-mcp-server/internal/application/actor"
)

// ActorService defines the interface for actor operations
type ActorService interface {
	CreateActor(ctx context.Context, cmd actorApp.CreateActorCommand) (*actorApp.ActorDTO, error)
	GetActor(ctx context.Context, id int) (*actorApp.ActorDTO, error)
	UpdateActor(ctx context.Context, cmd actorApp.UpdateActorCommand) (*actorApp.ActorDTO, error)
	DeleteActor(ctx context.Context, id int) error
	LinkActorToMovie(ctx context.Context, actorID, movieID int) error
	UnlinkActorFromMovie(ctx context.Context, actorID, movieID int) error
	GetActorsByMovie(ctx context.Context, movieID int) ([]*actorApp.ActorDTO, error)
	SearchActors(ctx context.Context, query actorApp.SearchActorsQuery) ([]*actorApp.ActorDTO, error)
}

// ActorTools provides SDK-based MCP handlers for actor operations
type ActorTools struct {
	actorService ActorService
}

// NewActorTools creates a new actor tools instance
func NewActorTools(actorService ActorService) *ActorTools {
	return &ActorTools{
		actorService: actorService,
	}
}

// ===== Actor Output Type (shared) =====

// ActorOutput defines the common output schema for actor data
type ActorOutput struct {
	ID        int    `json:"id" jsonschema:"description=Actor ID"`
	Name      string `json:"name" jsonschema:"description=Actor name"`
	BirthYear int    `json:"birth_year,omitempty" jsonschema:"description=Birth year"`
	Bio       string `json:"bio,omitempty" jsonschema:"description=Biography"`
	MovieIDs  []int  `json:"movie_ids" jsonschema:"description=List of movie IDs the actor appears in"`
	CreatedAt string `json:"created_at" jsonschema:"description=Creation timestamp"`
	UpdatedAt string `json:"updated_at" jsonschema:"description=Last update timestamp"`
}

// ===== get_actor Tool =====

// GetActorInput defines the input schema for get_actor tool
type GetActorInput struct {
	ActorID int `json:"actor_id" jsonschema:"required,description=The actor ID to retrieve"`
}

// GetActor handles the get_actor tool call
func (t *ActorTools) GetActor(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetActorInput,
) (*mcp.CallToolResult, ActorOutput, error) {
	actorDTO, err := t.actorService.GetActor(ctx, input.ActorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ActorOutput{}, fmt.Errorf("actor not found")
		}
		return nil, ActorOutput{}, fmt.Errorf("failed to get actor: %w", err)
	}

	output := ActorOutput{
		ID:        actorDTO.ID,
		Name:      actorDTO.Name,
		BirthYear: actorDTO.BirthYear,
		Bio:       actorDTO.Bio,
		MovieIDs:  actorDTO.MovieIDs,
		CreatedAt: actorDTO.CreatedAt,
		UpdatedAt: actorDTO.UpdatedAt,
	}

	return nil, output, nil
}

// ===== add_actor Tool =====

// AddActorInput defines the input schema for add_actor tool
type AddActorInput struct {
	Name      string `json:"name" jsonschema:"required,description=Actor name"`
	BirthYear int    `json:"birth_year,omitempty" jsonschema:"description=Birth year"`
	Bio       string `json:"bio,omitempty" jsonschema:"description=Biography"`
}

// AddActor handles the add_actor tool call
func (t *ActorTools) AddActor(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddActorInput,
) (*mcp.CallToolResult, ActorOutput, error) {
	cmd := actorApp.CreateActorCommand{
		Name:      input.Name,
		BirthYear: input.BirthYear,
		Bio:       input.Bio,
	}

	actorDTO, err := t.actorService.CreateActor(ctx, cmd)
	if err != nil {
		return nil, ActorOutput{}, fmt.Errorf("failed to create actor: %w", err)
	}

	output := ActorOutput{
		ID:        actorDTO.ID,
		Name:      actorDTO.Name,
		BirthYear: actorDTO.BirthYear,
		Bio:       actorDTO.Bio,
		MovieIDs:  actorDTO.MovieIDs,
		CreatedAt: actorDTO.CreatedAt,
		UpdatedAt: actorDTO.UpdatedAt,
	}

	return nil, output, nil
}

// ===== update_actor Tool =====

// UpdateActorInput defines the input schema for update_actor tool
type UpdateActorInput struct {
	ID        int    `json:"id" jsonschema:"required,description=Actor ID"`
	Name      string `json:"name" jsonschema:"required,description=Actor name"`
	BirthYear int    `json:"birth_year,omitempty" jsonschema:"description=Birth year"`
	Bio       string `json:"bio,omitempty" jsonschema:"description=Biography"`
}

// UpdateActor handles the update_actor tool call
func (t *ActorTools) UpdateActor(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UpdateActorInput,
) (*mcp.CallToolResult, ActorOutput, error) {
	cmd := actorApp.UpdateActorCommand{
		ID:        input.ID,
		Name:      input.Name,
		BirthYear: input.BirthYear,
		Bio:       input.Bio,
	}

	actorDTO, err := t.actorService.UpdateActor(ctx, cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ActorOutput{}, fmt.Errorf("actor not found")
		}
		return nil, ActorOutput{}, fmt.Errorf("failed to update actor: %w", err)
	}

	output := ActorOutput{
		ID:        actorDTO.ID,
		Name:      actorDTO.Name,
		BirthYear: actorDTO.BirthYear,
		Bio:       actorDTO.Bio,
		MovieIDs:  actorDTO.MovieIDs,
		CreatedAt: actorDTO.CreatedAt,
		UpdatedAt: actorDTO.UpdatedAt,
	}

	return nil, output, nil
}

// ===== delete_actor Tool =====

// DeleteActorInput defines the input schema for delete_actor tool
type DeleteActorInput struct {
	ActorID int `json:"actor_id" jsonschema:"required,description=The actor ID to delete"`
}

// DeleteActorOutput defines the output schema for delete_actor tool
type DeleteActorOutput struct {
	Message string `json:"message" jsonschema:"description=Success message"`
}

// DeleteActor handles the delete_actor tool call
func (t *ActorTools) DeleteActor(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DeleteActorInput,
) (*mcp.CallToolResult, DeleteActorOutput, error) {
	err := t.actorService.DeleteActor(ctx, input.ActorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, DeleteActorOutput{}, fmt.Errorf("actor not found")
		}
		return nil, DeleteActorOutput{}, fmt.Errorf("failed to delete actor: %w", err)
	}

	output := DeleteActorOutput{
		Message: "Actor deleted successfully",
	}

	return nil, output, nil
}

// ===== link_actor_to_movie Tool =====

// LinkActorToMovieInput defines the input schema for link_actor_to_movie tool
type LinkActorToMovieInput struct {
	ActorID int `json:"actor_id" jsonschema:"required,description=Actor ID"`
	MovieID int `json:"movie_id" jsonschema:"required,description=Movie ID"`
}

// LinkActorToMovieOutput defines the output schema for link_actor_to_movie tool
type LinkActorToMovieOutput struct {
	Message string `json:"message" jsonschema:"description=Success message"`
}

// LinkActorToMovie handles the link_actor_to_movie tool call
func (t *ActorTools) LinkActorToMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LinkActorToMovieInput,
) (*mcp.CallToolResult, LinkActorToMovieOutput, error) {
	err := t.actorService.LinkActorToMovie(ctx, input.ActorID, input.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, LinkActorToMovieOutput{}, fmt.Errorf("actor or movie not found")
		}
		if strings.Contains(err.Error(), "already exists") {
			return nil, LinkActorToMovieOutput{}, fmt.Errorf("actor is already linked to this movie")
		}
		return nil, LinkActorToMovieOutput{}, fmt.Errorf("failed to link actor to movie: %w", err)
	}

	output := LinkActorToMovieOutput{
		Message: "Actor linked to movie successfully",
	}

	return nil, output, nil
}

// ===== unlink_actor_from_movie Tool =====

// UnlinkActorFromMovieInput defines the input schema for unlink_actor_from_movie tool
type UnlinkActorFromMovieInput struct {
	ActorID int `json:"actor_id" jsonschema:"required,description=Actor ID"`
	MovieID int `json:"movie_id" jsonschema:"required,description=Movie ID"`
}

// UnlinkActorFromMovieOutput defines the output schema for unlink_actor_from_movie tool
type UnlinkActorFromMovieOutput struct {
	Message string `json:"message" jsonschema:"description=Success message"`
}

// UnlinkActorFromMovie handles the unlink_actor_from_movie tool call
func (t *ActorTools) UnlinkActorFromMovie(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UnlinkActorFromMovieInput,
) (*mcp.CallToolResult, UnlinkActorFromMovieOutput, error) {
	err := t.actorService.UnlinkActorFromMovie(ctx, input.ActorID, input.MovieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, UnlinkActorFromMovieOutput{}, fmt.Errorf("actor, movie, or link not found")
		}
		return nil, UnlinkActorFromMovieOutput{}, fmt.Errorf("failed to unlink actor from movie: %w", err)
	}

	output := UnlinkActorFromMovieOutput{
		Message: "Actor unlinked from movie successfully",
	}

	return nil, output, nil
}

// ===== get_movie_cast Tool =====

// GetMovieCastInput defines the input schema for get_movie_cast tool
type GetMovieCastInput struct {
	MovieID int `json:"movie_id" jsonschema:"required,description=Movie ID to get cast for"`
}

// GetMovieCastOutput defines the output schema for get_movie_cast tool
type GetMovieCastOutput struct {
	Actors      []ActorOutput `json:"actors" jsonschema:"description=List of actors in the movie"`
	Total       int           `json:"total" jsonschema:"description=Total number of actors"`
	Description string        `json:"description" jsonschema:"description=Description of results"`
}

// GetMovieCast handles the get_movie_cast tool call
func (t *ActorTools) GetMovieCast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetMovieCastInput,
) (*mcp.CallToolResult, GetMovieCastOutput, error) {
	actorDTOs, err := t.actorService.GetActorsByMovie(ctx, input.MovieID)
	if err != nil {
		return nil, GetMovieCastOutput{}, fmt.Errorf("failed to get movie cast: %w", err)
	}

	actors := make([]ActorOutput, len(actorDTOs))
	for i, actorDTO := range actorDTOs {
		actors[i] = ActorOutput{
			ID:        actorDTO.ID,
			Name:      actorDTO.Name,
			BirthYear: actorDTO.BirthYear,
			Bio:       actorDTO.Bio,
			MovieIDs:  actorDTO.MovieIDs,
			CreatedAt: actorDTO.CreatedAt,
			UpdatedAt: actorDTO.UpdatedAt,
		}
	}

	output := GetMovieCastOutput{
		Actors:      actors,
		Total:       len(actors),
		Description: fmt.Sprintf("Cast of movie %d", input.MovieID),
	}

	return nil, output, nil
}

// ===== get_actor_movies Tool =====

// GetActorMoviesInput defines the input schema for get_actor_movies tool
type GetActorMoviesInput struct {
	ActorID int `json:"actor_id" jsonschema:"required,description=Actor ID to get movies for"`
}

// GetActorMoviesOutput defines the output schema for get_actor_movies tool
type GetActorMoviesOutput struct {
	ActorID     int    `json:"actor_id" jsonschema:"description=Actor ID"`
	ActorName   string `json:"actor_name" jsonschema:"description=Actor name"`
	MovieIDs    []int  `json:"movie_ids" jsonschema:"description=List of movie IDs"`
	TotalMovies int    `json:"total_movies" jsonschema:"description=Total number of movies"`
}

// GetActorMovies handles the get_actor_movies tool call
func (t *ActorTools) GetActorMovies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetActorMoviesInput,
) (*mcp.CallToolResult, GetActorMoviesOutput, error) {
	actorDTO, err := t.actorService.GetActor(ctx, input.ActorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, GetActorMoviesOutput{}, fmt.Errorf("actor not found")
		}
		return nil, GetActorMoviesOutput{}, fmt.Errorf("failed to get actor: %w", err)
	}

	output := GetActorMoviesOutput{
		ActorID:     actorDTO.ID,
		ActorName:   actorDTO.Name,
		MovieIDs:    actorDTO.MovieIDs,
		TotalMovies: len(actorDTO.MovieIDs),
	}

	return nil, output, nil
}

// ===== search_actors Tool =====

// SearchActorsInput defines the input schema for search_actors tool
type SearchActorsInput struct {
	Name         string `json:"name,omitempty" jsonschema:"description=Search by actor name"`
	MinBirthYear int    `json:"min_birth_year,omitempty" jsonschema:"description=Minimum birth year"`
	MaxBirthYear int    `json:"max_birth_year,omitempty" jsonschema:"description=Maximum birth year"`
	MovieID      int    `json:"movie_id,omitempty" jsonschema:"description=Filter actors by movie ID"`
	Limit        int    `json:"limit,omitempty" jsonschema:"description=Maximum number of results,default=20"`
	Offset       int    `json:"offset,omitempty" jsonschema:"description=Number of results to skip for pagination,default=0"`
	OrderBy      string `json:"order_by,omitempty" jsonschema:"description=Field to order by (name/birth_year),default=name"`
	OrderDir     string `json:"order_dir,omitempty" jsonschema:"description=Order direction (asc/desc),default=asc"`
}

// SearchActorsOutput defines the output schema for search_actors tool
type SearchActorsOutput struct {
	Actors      []ActorOutput `json:"actors" jsonschema:"description=List of matching actors"`
	Total       int           `json:"total" jsonschema:"description=Total number of actors found"`
	Description string        `json:"description" jsonschema:"description=Description of search results"`
}

// SearchActors handles the search_actors tool call
func (t *ActorTools) SearchActors(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchActorsInput,
) (*mcp.CallToolResult, SearchActorsOutput, error) {
	query := actorApp.SearchActorsQuery{
		Name:         input.Name,
		MinBirthYear: input.MinBirthYear,
		MaxBirthYear: input.MaxBirthYear,
		MovieID:      input.MovieID,
		Limit:        input.Limit,
		Offset:       input.Offset,
		OrderBy:      input.OrderBy,
		OrderDir:     input.OrderDir,
	}

	// Set default limit
	if query.Limit == 0 {
		query.Limit = 20
	}

	actorDTOs, err := t.actorService.SearchActors(ctx, query)
	if err != nil {
		return nil, SearchActorsOutput{}, fmt.Errorf("failed to search actors: %w", err)
	}

	actors := make([]ActorOutput, len(actorDTOs))
	for i, actorDTO := range actorDTOs {
		actors[i] = ActorOutput{
			ID:        actorDTO.ID,
			Name:      actorDTO.Name,
			BirthYear: actorDTO.BirthYear,
			Bio:       actorDTO.Bio,
			MovieIDs:  actorDTO.MovieIDs,
			CreatedAt: actorDTO.CreatedAt,
			UpdatedAt: actorDTO.UpdatedAt,
		}
	}

	output := SearchActorsOutput{
		Actors:      actors,
		Total:       len(actors),
		Description: "Search results",
	}

	return nil, output, nil
}
