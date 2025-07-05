package actor

import (
	"context"
	"fmt"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/actor"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// Service provides application-level actor operations
type Service struct {
	actorRepo actor.Repository
}

// NewService creates a new actor application service
func NewService(actorRepo actor.Repository) *Service {
	return &Service{
		actorRepo: actorRepo,
	}
}

// CreateActorCommand represents the command to create a new actor
type CreateActorCommand struct {
	Name      string
	BirthYear int
	Bio       string
}

// UpdateActorCommand represents the command to update an existing actor
type UpdateActorCommand struct {
	ID        int
	Name      string
	BirthYear int
	Bio       string
}

// SearchActorsQuery represents the query to search for actors
type SearchActorsQuery struct {
	Name         string
	MinBirthYear int
	MaxBirthYear int
	MovieID      int // Find actors who appeared in this movie
	Limit        int
	Offset       int
	OrderBy      string
	OrderDir     string
}

// ActorDTO represents an actor data transfer object
type ActorDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Bio       string `json:"bio,omitempty"`
	MovieIDs  []int  `json:"movie_ids"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateActor creates a new actor
func (s *Service) CreateActor(ctx context.Context, cmd CreateActorCommand) (*ActorDTO, error) {
	// Create domain actor
	domainActor, err := actor.NewActor(cmd.Name, cmd.BirthYear)
	if err != nil {
		return nil, fmt.Errorf("failed to create actor: %w", err)
	}

	// Set bio if provided
	if cmd.Bio != "" {
		domainActor.SetBio(cmd.Bio)
	}

	// Validate the actor
	if err := domainActor.Validate(); err != nil {
		return nil, fmt.Errorf("actor validation failed: %w", err)
	}

	// Save to repository
	if err := s.actorRepo.Save(ctx, domainActor); err != nil {
		return nil, fmt.Errorf("failed to save actor: %w", err)
	}

	return s.toDTO(domainActor), nil
}

// GetActor retrieves an actor by ID
func (s *Service) GetActor(ctx context.Context, id int) (*ActorDTO, error) {
	actorID, err := shared.NewActorID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid actor ID: %w", err)
	}

	domainActor, err := s.actorRepo.FindByID(ctx, actorID)
	if err != nil {
		return nil, fmt.Errorf("actor not found: %w", err)
	}

	return s.toDTO(domainActor), nil
}

// UpdateActor updates an existing actor
func (s *Service) UpdateActor(ctx context.Context, cmd UpdateActorCommand) (*ActorDTO, error) {
	actorID, err := shared.NewActorID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid actor ID: %w", err)
	}

	// Get existing actor to preserve movie links
	existingActor, err := s.actorRepo.FindByID(ctx, actorID)
	if err != nil {
		return nil, fmt.Errorf("actor not found: %w", err)
	}

	// Create new actor with updated values
	updatedActor, err := actor.NewActorWithID(actorID, cmd.Name, cmd.BirthYear)
	if err != nil {
		return nil, fmt.Errorf("failed to create updated actor: %w", err)
	}

	// Set bio
	updatedActor.SetBio(cmd.Bio)

	// Preserve existing movie links
	for _, movieID := range existingActor.MovieIDs() {
		if err := updatedActor.AddMovie(movieID); err != nil {
			return nil, fmt.Errorf("failed to preserve movie link: %w", err)
		}
	}

	// Validate the updated actor
	if err := updatedActor.Validate(); err != nil {
		return nil, fmt.Errorf("actor validation failed: %w", err)
	}

	// Save updated actor
	if err := s.actorRepo.Save(ctx, updatedActor); err != nil {
		return nil, fmt.Errorf("failed to save updated actor: %w", err)
	}

	return s.toDTO(updatedActor), nil
}

// DeleteActor deletes an actor by ID
func (s *Service) DeleteActor(ctx context.Context, id int) error {
	actorID, err := shared.NewActorID(id)
	if err != nil {
		return fmt.Errorf("invalid actor ID: %w", err)
	}

	if err := s.actorRepo.Delete(ctx, actorID); err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}

	return nil
}

// LinkActorToMovie links an actor to a movie
func (s *Service) LinkActorToMovie(ctx context.Context, actorID, movieID int) error {
	actorDomainID, err := shared.NewActorID(actorID)
	if err != nil {
		return fmt.Errorf("invalid actor ID: %w", err)
	}

	movieDomainID, err := shared.NewMovieID(movieID)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	// Get existing actor
	domainActor, err := s.actorRepo.FindByID(ctx, actorDomainID)
	if err != nil {
		return fmt.Errorf("actor not found: %w", err)
	}

	// Add movie to actor's filmography
	if err := domainActor.AddMovie(movieDomainID); err != nil {
		return fmt.Errorf("failed to link actor to movie: %w", err)
	}

	// Save updated actor
	if err := s.actorRepo.Save(ctx, domainActor); err != nil {
		return fmt.Errorf("failed to save actor: %w", err)
	}

	return nil
}

// UnlinkActorFromMovie removes the link between an actor and a movie
func (s *Service) UnlinkActorFromMovie(ctx context.Context, actorID, movieID int) error {
	actorDomainID, err := shared.NewActorID(actorID)
	if err != nil {
		return fmt.Errorf("invalid actor ID: %w", err)
	}

	movieDomainID, err := shared.NewMovieID(movieID)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	// Get existing actor
	domainActor, err := s.actorRepo.FindByID(ctx, actorDomainID)
	if err != nil {
		return fmt.Errorf("actor not found: %w", err)
	}

	// Remove movie from actor's filmography
	if err := domainActor.RemoveMovie(movieDomainID); err != nil {
		return fmt.Errorf("failed to unlink actor from movie: %w", err)
	}

	// Save updated actor
	if err := s.actorRepo.Save(ctx, domainActor); err != nil {
		return fmt.Errorf("failed to save actor: %w", err)
	}

	return nil
}

// SearchActors searches for actors based on criteria
func (s *Service) SearchActors(ctx context.Context, query SearchActorsQuery) ([]*ActorDTO, error) {
	criteria := actor.SearchCriteria{
		Name:         query.Name,
		MinBirthYear: query.MinBirthYear,
		MaxBirthYear: query.MaxBirthYear,
		Limit:        query.Limit,
		Offset:       query.Offset,
	}

	// Set MovieID if provided
	if query.MovieID != 0 {
		movieID, err := shared.NewMovieID(query.MovieID)
		if err != nil {
			return nil, fmt.Errorf("invalid movie ID: %w", err)
		}
		criteria.MovieID = movieID
	}

	// Set default limit if not provided
	if criteria.Limit == 0 {
		criteria.Limit = 50
	}

	// Set order by
	switch query.OrderBy {
	case "name":
		criteria.OrderBy = actor.OrderByName
	case "birth_year":
		criteria.OrderBy = actor.OrderByBirthYear
	case "created_at":
		criteria.OrderBy = actor.OrderByCreatedAt
	case "updated_at":
		criteria.OrderBy = actor.OrderByUpdatedAt
	default:
		criteria.OrderBy = actor.OrderByName
	}

	// Set order direction
	if query.OrderDir == "desc" {
		criteria.OrderDir = actor.OrderDesc
	} else {
		criteria.OrderDir = actor.OrderAsc
	}

	domainActors, err := s.actorRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search actors: %w", err)
	}

	var dtos []*ActorDTO
	for _, domainActor := range domainActors {
		dtos = append(dtos, s.toDTO(domainActor))
	}

	return dtos, nil
}

// GetActorsByMovie retrieves all actors who appeared in a specific movie
func (s *Service) GetActorsByMovie(ctx context.Context, movieID int) ([]*ActorDTO, error) {
	movieDomainID, err := shared.NewMovieID(movieID)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	domainActors, err := s.actorRepo.FindByMovieID(ctx, movieDomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get actors by movie: %w", err)
	}

	var dtos []*ActorDTO
	for _, domainActor := range domainActors {
		dtos = append(dtos, s.toDTO(domainActor))
	}

	return dtos, nil
}

// toDTO converts a domain actor to a DTO
func (s *Service) toDTO(domainActor *actor.Actor) *ActorDTO {
	movieIDs := make([]int, len(domainActor.MovieIDs()))
	for i, movieID := range domainActor.MovieIDs() {
		movieIDs[i] = movieID.Value()
	}

	dto := &ActorDTO{
		ID:        domainActor.ID().Value(),
		Name:      domainActor.Name(),
		BirthYear: domainActor.BirthYear().Value(),
		Bio:       domainActor.Bio(),
		MovieIDs:  movieIDs,
		CreatedAt: domainActor.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: domainActor.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	return dto
}
