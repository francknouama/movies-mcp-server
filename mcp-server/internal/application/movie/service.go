package movie

import (
	"context"
	"fmt"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/movie"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// Service provides application-level movie operations
type Service struct {
	movieRepo movie.Repository
}

// NewService creates a new movie application service
func NewService(movieRepo movie.Repository) *Service {
	return &Service{
		movieRepo: movieRepo,
	}
}

// CreateMovieCommand represents the command to create a new movie
type CreateMovieCommand struct {
	Title     string
	Director  string
	Year      int
	Rating    float64
	Genres    []string
	PosterURL string
}

// UpdateMovieCommand represents the command to update an existing movie
type UpdateMovieCommand struct {
	ID        int
	Title     string
	Director  string
	Year      int
	Rating    float64
	Genres    []string
	PosterURL string
}

// SearchMoviesQuery represents the query to search for movies
type SearchMoviesQuery struct {
	Title       string
	Director    string
	Genre       string
	MinYear     int
	MaxYear     int
	MinRating   float64
	MaxRating   float64
	Limit       int
	Offset      int
	OrderBy     string
	OrderDir    string
}

// MovieDTO represents a movie data transfer object
type MovieDTO struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Director  string   `json:"director"`
	Year      int      `json:"year"`
	Rating    float64  `json:"rating"`
	Genres    []string `json:"genres"`
	PosterURL string   `json:"poster_url,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// CreateMovie creates a new movie
func (s *Service) CreateMovie(ctx context.Context, cmd CreateMovieCommand) (*MovieDTO, error) {
	// Create domain movie
	domainMovie, err := movie.NewMovie(cmd.Title, cmd.Director, cmd.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	// Set rating if provided
	if cmd.Rating > 0 {
		if err := domainMovie.SetRating(cmd.Rating); err != nil {
			return nil, fmt.Errorf("failed to set rating: %w", err)
		}
	}

	// Add genres if provided
	for _, genre := range cmd.Genres {
		if err := domainMovie.AddGenre(genre); err != nil {
			return nil, fmt.Errorf("failed to add genre %s: %w", genre, err)
		}
	}

	// Set poster URL if provided
	if cmd.PosterURL != "" {
		if err := domainMovie.SetPosterURL(cmd.PosterURL); err != nil {
			return nil, fmt.Errorf("failed to set poster URL: %w", err)
		}
	}

	// Validate the movie
	if err := domainMovie.Validate(); err != nil {
		return nil, fmt.Errorf("movie validation failed: %w", err)
	}

	// Save to repository
	if err := s.movieRepo.Save(ctx, domainMovie); err != nil {
		return nil, fmt.Errorf("failed to save movie: %w", err)
	}

	return s.toDTO(domainMovie), nil
}

// GetMovie retrieves a movie by ID
func (s *Service) GetMovie(ctx context.Context, id int) (*MovieDTO, error) {
	movieID, err := shared.NewMovieID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	domainMovie, err := s.movieRepo.FindByID(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	return s.toDTO(domainMovie), nil
}

// UpdateMovie updates an existing movie
func (s *Service) UpdateMovie(ctx context.Context, cmd UpdateMovieCommand) (*MovieDTO, error) {
	movieID, err := shared.NewMovieID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	// Verify movie exists
	_, err = s.movieRepo.FindByID(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Create new movie with updated values (immutable approach)
	updatedMovie, err := movie.NewMovieWithID(movieID, cmd.Title, cmd.Director, cmd.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to create updated movie: %w", err)
	}

	// Set rating if provided
	if cmd.Rating > 0 {
		if err := updatedMovie.SetRating(cmd.Rating); err != nil {
			return nil, fmt.Errorf("failed to set rating: %w", err)
		}
	}

	// Add genres if provided
	for _, genre := range cmd.Genres {
		if err := updatedMovie.AddGenre(genre); err != nil {
			return nil, fmt.Errorf("failed to add genre %s: %w", genre, err)
		}
	}

	// Set poster URL if provided
	if cmd.PosterURL != "" {
		if err := updatedMovie.SetPosterURL(cmd.PosterURL); err != nil {
			return nil, fmt.Errorf("failed to set poster URL: %w", err)
		}
	}

	// Validate the updated movie
	if err := updatedMovie.Validate(); err != nil {
		return nil, fmt.Errorf("movie validation failed: %w", err)
	}

	// Save updated movie
	if err := s.movieRepo.Save(ctx, updatedMovie); err != nil {
		return nil, fmt.Errorf("failed to save updated movie: %w", err)
	}

	return s.toDTO(updatedMovie), nil
}

// DeleteMovie deletes a movie by ID
func (s *Service) DeleteMovie(ctx context.Context, id int) error {
	movieID, err := shared.NewMovieID(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	if err := s.movieRepo.Delete(ctx, movieID); err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

// SearchMovies searches for movies based on criteria
func (s *Service) SearchMovies(ctx context.Context, query SearchMoviesQuery) ([]*MovieDTO, error) {
	criteria := movie.SearchCriteria{
		Title:     query.Title,
		Director:  query.Director,
		Genre:     query.Genre,
		MinYear:   query.MinYear,
		MaxYear:   query.MaxYear,
		MinRating: query.MinRating,
		MaxRating: query.MaxRating,
		Limit:     query.Limit,
		Offset:    query.Offset,
	}

	// Set default limit if not provided
	if criteria.Limit == 0 {
		criteria.Limit = 50
	}

	// Set order by
	switch query.OrderBy {
	case "title":
		criteria.OrderBy = movie.OrderByTitle
	case "director":
		criteria.OrderBy = movie.OrderByDirector
	case "year":
		criteria.OrderBy = movie.OrderByYear
	case "rating":
		criteria.OrderBy = movie.OrderByRating
	case "created_at":
		criteria.OrderBy = movie.OrderByCreatedAt
	case "updated_at":
		criteria.OrderBy = movie.OrderByUpdatedAt
	default:
		criteria.OrderBy = movie.OrderByTitle
	}

	// Set order direction
	if query.OrderDir == "desc" {
		criteria.OrderDir = movie.OrderDesc
	} else {
		criteria.OrderDir = movie.OrderAsc
	}

	domainMovies, err := s.movieRepo.FindByCriteria(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	var dtos []*MovieDTO
	for _, domainMovie := range domainMovies {
		dtos = append(dtos, s.toDTO(domainMovie))
	}

	return dtos, nil
}

// GetTopRatedMovies retrieves top-rated movies
func (s *Service) GetTopRatedMovies(ctx context.Context, limit int) ([]*MovieDTO, error) {
	if limit <= 0 {
		limit = 10
	}

	domainMovies, err := s.movieRepo.FindTopRated(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated movies: %w", err)
	}

	var dtos []*MovieDTO
	for _, domainMovie := range domainMovies {
		dtos = append(dtos, s.toDTO(domainMovie))
	}

	return dtos, nil
}

// toDTO converts a domain movie to a DTO
func (s *Service) toDTO(domainMovie *movie.Movie) *MovieDTO {
	dto := &MovieDTO{
		ID:        domainMovie.ID().Value(),
		Title:     domainMovie.Title(),
		Director:  domainMovie.Director(),
		Year:      domainMovie.Year().Value(),
		Rating:    domainMovie.Rating().Value(),
		Genres:    domainMovie.Genres(),
		PosterURL: domainMovie.PosterURL(),
		CreatedAt: domainMovie.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: domainMovie.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	return dto
}