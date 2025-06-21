package movie

import (
	"testing"
	"time"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/domain/shared"
)

func TestNewMovie(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		director string
		year     int
		wantErr  bool
	}{
		{
			name:     "valid movie",
			title:    "Inception",
			director: "Christopher Nolan",
			year:     2010,
			wantErr:  false,
		},
		{
			name:     "empty title",
			title:    "",
			director: "Christopher Nolan",
			year:     2010,
			wantErr:  true,
		},
		{
			name:     "whitespace only title",
			title:    "   ",
			director: "Christopher Nolan",
			year:     2010,
			wantErr:  true,
		},
		{
			name:     "empty director",
			title:    "Inception",
			director: "",
			year:     2010,
			wantErr:  true,
		},
		{
			name:     "whitespace only director",
			title:    "Inception",
			director: "   ",
			year:     2010,
			wantErr:  true,
		},
		{
			name:     "invalid year",
			title:    "Inception",
			director: "Christopher Nolan",
			year:     1800,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie, err := NewMovie(tt.title, tt.director, tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if movie.Title() != tt.title {
					t.Errorf("NewMovie() title = %v, want %v", movie.Title(), tt.title)
				}
				if movie.Director() != tt.director {
					t.Errorf("NewMovie() director = %v, want %v", movie.Director(), tt.director)
				}
				if movie.Year().Value() != tt.year {
					t.Errorf("NewMovie() year = %v, want %v", movie.Year().Value(), tt.year)
				}
				if movie.ID().IsZero() {
					t.Error("NewMovie() should not have zero ID after creation")
				}
			}
		})
	}
}

func TestMovie_SetRating(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	tests := []struct {
		name    string
		rating  float64
		wantErr bool
	}{
		{
			name:    "valid rating",
			rating:  7.5,
			wantErr: false,
		},
		{
			name:    "minimum rating",
			rating:  0.0,
			wantErr: false,
		},
		{
			name:    "maximum rating",
			rating:  10.0,
			wantErr: false,
		},
		{
			name:    "invalid negative rating",
			rating:  -1.0,
			wantErr: true,
		},
		{
			name:    "invalid too high rating",
			rating:  11.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := movie.SetRating(tt.rating)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && movie.Rating().Value() != tt.rating {
				t.Errorf("SetRating() rating = %v, want %v", movie.Rating().Value(), tt.rating)
			}
		})
	}
}

func TestMovie_AddGenre(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	tests := []struct {
		name    string
		genre   string
		wantErr bool
	}{
		{
			name:    "valid genre",
			genre:   "Action",
			wantErr: false,
		},
		{
			name:    "empty genre",
			genre:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only genre",
			genre:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := movie.AddGenre(tt.genre)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				genres := movie.Genres()
				found := false
				for _, g := range genres {
					if g == tt.genre {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("AddGenre() genre %v not found in genres %v", tt.genre, genres)
				}
			}
		})
	}
}

func TestMovie_AddGenre_Duplicate(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	// Add genre first time
	err = movie.AddGenre("Action")
	if err != nil {
		t.Fatalf("Failed to add genre: %v", err)
	}

	// Try to add same genre again
	err = movie.AddGenre("Action")
	if err == nil {
		t.Error("Expected error when adding duplicate genre")
	}

	// Verify only one instance exists
	genres := movie.Genres()
	count := 0
	for _, g := range genres {
		if g == "Action" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected 1 instance of 'Action' genre, got %d", count)
	}
}

func TestMovie_SetPosterURL(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	tests := []struct {
		name      string
		posterURL string
		wantErr   bool
	}{
		{
			name:      "valid HTTP URL",
			posterURL: "http://example.com/poster.jpg",
			wantErr:   false,
		},
		{
			name:      "valid HTTPS URL",
			posterURL: "https://example.com/poster.jpg",
			wantErr:   false,
		},
		{
			name:      "empty URL (removing poster)",
			posterURL: "",
			wantErr:   false,
		},
		{
			name:      "invalid URL format",
			posterURL: "not-a-url",
			wantErr:   true,
		},
		{
			name:      "invalid scheme",
			posterURL: "ftp://example.com/poster.jpg",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := movie.SetPosterURL(tt.posterURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPosterURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && movie.PosterURL() != tt.posterURL {
				t.Errorf("SetPosterURL() posterURL = %v, want %v", movie.PosterURL(), tt.posterURL)
			}
		})
	}
}

func TestMovie_UpdateTimestamp(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	originalUpdatedAt := movie.UpdatedAt()
	
	// Sleep to ensure timestamp changes
	time.Sleep(1 * time.Millisecond)
	
	movie.touch() // This should be called internally by business operations
	
	if !movie.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated after modification")
	}
}

func TestMovie_Validation(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	// Valid movie should pass validation
	if err := movie.Validate(); err != nil {
		t.Errorf("Valid movie should pass validation, got: %v", err)
	}

	// Test validation with genres
	movie.AddGenre("Action")
	movie.AddGenre("Thriller")
	if err := movie.Validate(); err != nil {
		t.Errorf("Movie with genres should pass validation, got: %v", err)
	}
}

func TestNewMovieWithID(t *testing.T) {
	id, err := shared.NewMovieID(123)
	if err != nil {
		t.Fatalf("Failed to create movie ID: %v", err)
	}

	movie, err := NewMovieWithID(id, "Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create movie with ID: %v", err)
	}

	if movie.ID().Value() != 123 {
		t.Errorf("Expected movie ID 123, got %d", movie.ID().Value())
	}
}

func TestMovie_HasGenre(t *testing.T) {
	movie, err := NewMovie("Test Movie", "Test Director", 2020)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	movie.AddGenre("Action")
	movie.AddGenre("Thriller")

	if !movie.HasGenre("Action") {
		t.Error("Expected movie to have Action genre")
	}

	if !movie.HasGenre("Thriller") {
		t.Error("Expected movie to have Thriller genre")
	}

	if movie.HasGenre("Comedy") {
		t.Error("Expected movie to not have Comedy genre")
	}
}