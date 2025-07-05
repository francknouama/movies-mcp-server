package movie

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

// Movie represents a movie aggregate in the domain
type Movie struct {
	shared.AggregateRoot
	id        shared.MovieID
	title     string
	director  string
	year      shared.Year
	rating    shared.Rating
	genres    []string
	posterURL string
	createdAt time.Time
	updatedAt time.Time
}

// NewMovie creates a new Movie with validation
func NewMovie(title, director string, year int) (*Movie, error) {
	// Generate a zero ID - this will be replaced when saved to repository
	id, err := shared.NewMovieID(0) // Zero ID indicates new movie
	if err != nil {
		return nil, err
	}

	return NewMovieWithID(id, title, director, year)
}

// NewMovieWithID creates a new Movie with a specific ID (for repository reconstruction)
func NewMovieWithID(id shared.MovieID, title, director string, year int) (*Movie, error) {
	// Validate inputs
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if strings.TrimSpace(director) == "" {
		return nil, errors.New("director cannot be empty")
	}

	movieYear, err := shared.NewYear(year)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	movie := &Movie{
		AggregateRoot: shared.NewAggregateRoot(),
		id:            id,
		title:         strings.TrimSpace(title),
		director:      strings.TrimSpace(director),
		year:          movieYear,
		genres:        make([]string, 0),
		createdAt:     now,
		updatedAt:     now,
	}

	// Emit domain event for new movie creation (only for non-zero IDs)
	if !id.IsZero() { // Skip zero ID
		event := NewMovieCreatedEvent(movie, movie.Version()+1)
		movie.AddEvent(event)
	}

	return movie, nil
}

// ID returns the movie's unique identifier
func (m *Movie) ID() shared.MovieID {
	return m.id
}

// Title returns the movie's title
func (m *Movie) Title() string {
	return m.title
}

// Director returns the movie's director
func (m *Movie) Director() string {
	return m.director
}

// Year returns the movie's release year
func (m *Movie) Year() shared.Year {
	return m.year
}

// Rating returns the movie's rating
func (m *Movie) Rating() shared.Rating {
	return m.rating
}

// Genres returns a copy of the movie's genres
func (m *Movie) Genres() []string {
	genres := make([]string, len(m.genres))
	copy(genres, m.genres)
	return genres
}

// PosterURL returns the movie's poster URL
func (m *Movie) PosterURL() string {
	return m.posterURL
}

// CreatedAt returns when the movie was created
func (m *Movie) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt returns when the movie was last updated
func (m *Movie) UpdatedAt() time.Time {
	return m.updatedAt
}

// SetRating sets the movie's rating with validation
func (m *Movie) SetRating(rating float64) error {
	newRating, err := shared.NewRating(rating)
	if err != nil {
		return err
	}

	// Emit domain event if rating actually changed
	if m.rating.Value() != newRating.Value() {
		event := NewMovieRatingChangedEvent(m.id, m.rating, newRating, m.Version()+1)
		m.AddEvent(event)
	}

	m.rating = newRating
	m.touch()
	return nil
}

// AddGenre adds a genre to the movie with validation
func (m *Movie) AddGenre(genre string) error {
	genre = strings.TrimSpace(genre)
	if genre == "" {
		return errors.New("genre cannot be empty")
	}

	// Check for duplicates
	for _, g := range m.genres {
		if g == genre {
			return errors.New("genre already exists")
		}
	}

	m.genres = append(m.genres, genre)

	// Emit domain event for genre addition
	event := NewMovieGenreAddedEvent(m.id, genre, m.Version()+1)
	m.AddEvent(event)

	m.touch()
	return nil
}

// HasGenre checks if the movie has a specific genre
func (m *Movie) HasGenre(genre string) bool {
	for _, g := range m.genres {
		if g == genre {
			return true
		}
	}
	return false
}

// SetPosterURL sets the movie's poster URL with validation
func (m *Movie) SetPosterURL(posterURL string) error {
	oldPosterURL := m.posterURL

	if posterURL == "" {
		if oldPosterURL != "" {
			// Emit domain event for poster URL change
			event := NewMoviePosterChangedEvent(m.id, oldPosterURL, "", m.Version()+1)
			m.AddEvent(event)
		}
		m.posterURL = ""
		m.touch()
		return nil
	}

	// Validate URL format
	parsedURL, err := url.Parse(posterURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Only allow HTTP and HTTPS schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("poster URL must use HTTP or HTTPS scheme")
	}

	// Emit domain event if poster URL actually changed
	if oldPosterURL != posterURL {
		event := NewMoviePosterChangedEvent(m.id, oldPosterURL, posterURL, m.Version()+1)
		m.AddEvent(event)
	}

	m.posterURL = posterURL
	m.touch()
	return nil
}

// Validate performs comprehensive validation of the movie
func (m *Movie) Validate() error {
	if strings.TrimSpace(m.title) == "" {
		return errors.New("title cannot be empty")
	}
	if strings.TrimSpace(m.director) == "" {
		return errors.New("director cannot be empty")
	}
	if m.year.IsZero() {
		return errors.New("year must be set")
	}
	// Rating is optional, but if set, must be valid (already validated in SetRating)
	// Genres are optional
	// PosterURL is optional, but if set, must be valid (already validated in SetPosterURL)
	return nil
}

// touch updates the updatedAt timestamp
func (m *Movie) touch() {
	m.updatedAt = time.Now()
}

// SetID sets the movie's ID (used by repository when saving)
func (m *Movie) SetID(id shared.MovieID) {
	m.id = id
	m.touch()
}
