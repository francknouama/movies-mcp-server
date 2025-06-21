package movie

import (
	"strconv"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// Movie Domain Events

// MovieCreatedEvent is raised when a new movie is created
type MovieCreatedEvent struct {
	shared.BaseDomainEvent
	MovieID   shared.MovieID
	Title     string
	Director  string
	Year      shared.Year
	Genres    []string
	Rating    shared.Rating
	PosterURL string
}

// NewMovieCreatedEvent creates a new MovieCreatedEvent
func NewMovieCreatedEvent(movie *Movie, version int) *MovieCreatedEvent {
	return &MovieCreatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MovieCreated",
			strconv.Itoa(movie.ID().Value()),
			"Movie",
			version,
		),
		MovieID:   movie.ID(),
		Title:     movie.Title(),
		Director:  movie.Director(),
		Year:      movie.Year(),
		Genres:    movie.Genres(),
		Rating:    movie.Rating(),
		PosterURL: movie.PosterURL(),
	}
}

// MovieUpdatedEvent is raised when a movie is updated
type MovieUpdatedEvent struct {
	shared.BaseDomainEvent
	MovieID   shared.MovieID
	Title     string
	Director  string
	Year      shared.Year
	Genres    []string
	Rating    shared.Rating
	PosterURL string
}

// NewMovieUpdatedEvent creates a new MovieUpdatedEvent
func NewMovieUpdatedEvent(movie *Movie, version int) *MovieUpdatedEvent {
	return &MovieUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MovieUpdated",
			strconv.Itoa(movie.ID().Value()),
			"Movie",
			version,
		),
		MovieID:   movie.ID(),
		Title:     movie.Title(),
		Director:  movie.Director(),
		Year:      movie.Year(),
		Genres:    movie.Genres(),
		Rating:    movie.Rating(),
		PosterURL: movie.PosterURL(),
	}
}

// MovieDeletedEvent is raised when a movie is deleted
type MovieDeletedEvent struct {
	shared.BaseDomainEvent
	MovieID shared.MovieID
	Title   string
}

// NewMovieDeletedEvent creates a new MovieDeletedEvent
func NewMovieDeletedEvent(movieID shared.MovieID, title string, version int) *MovieDeletedEvent {
	return &MovieDeletedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MovieDeleted",
			strconv.Itoa(movieID.Value()),
			"Movie",
			version,
		),
		MovieID: movieID,
		Title:   title,
	}
}

// MovieRatingChangedEvent is raised when a movie's rating is changed
type MovieRatingChangedEvent struct {
	shared.BaseDomainEvent
	MovieID   shared.MovieID
	OldRating shared.Rating
	NewRating shared.Rating
}

// NewMovieRatingChangedEvent creates a new MovieRatingChangedEvent
func NewMovieRatingChangedEvent(movieID shared.MovieID, oldRating, newRating shared.Rating, version int) *MovieRatingChangedEvent {
	return &MovieRatingChangedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MovieRatingChanged",
			strconv.Itoa(movieID.Value()),
			"Movie",
			version,
		),
		MovieID:   movieID,
		OldRating: oldRating,
		NewRating: newRating,
	}
}

// MovieGenreAddedEvent is raised when a genre is added to a movie
type MovieGenreAddedEvent struct {
	shared.BaseDomainEvent
	MovieID shared.MovieID
	Genre   string
}

// NewMovieGenreAddedEvent creates a new MovieGenreAddedEvent
func NewMovieGenreAddedEvent(movieID shared.MovieID, genre string, version int) *MovieGenreAddedEvent {
	return &MovieGenreAddedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MovieGenreAdded",
			strconv.Itoa(movieID.Value()),
			"Movie",
			version,
		),
		MovieID: movieID,
		Genre:   genre,
	}
}

// MoviePosterChangedEvent is raised when a movie's poster URL is changed
type MoviePosterChangedEvent struct {
	shared.BaseDomainEvent
	MovieID      shared.MovieID
	OldPosterURL string
	NewPosterURL string
}

// NewMoviePosterChangedEvent creates a new MoviePosterChangedEvent
func NewMoviePosterChangedEvent(movieID shared.MovieID, oldPosterURL, newPosterURL string, version int) *MoviePosterChangedEvent {
	return &MoviePosterChangedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(
			"MoviePosterChanged",
			strconv.Itoa(movieID.Value()),
			"Movie",
			version,
		),
		MovieID:      movieID,
		OldPosterURL: oldPosterURL,
		NewPosterURL: newPosterURL,
	}
}
