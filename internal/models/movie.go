package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Movie represents a movie in the system
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Year        int       `json:"year"`
	Genre       []string  `json:"genre"`
	Rating      float64   `json:"rating"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration,omitempty"` // minutes
	Language    string    `json:"language,omitempty"`
	Country     string    `json:"country,omitempty"`
	PosterType  string    `json:"poster_type,omitempty"` // MIME type, only when poster exists
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MovieCreateRequest represents the input for creating a movie
type MovieCreateRequest struct {
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Year        int       `json:"year"`
	Genre       []string  `json:"genre,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	Description string    `json:"description,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	Language    string    `json:"language,omitempty"`
	Country     string    `json:"country,omitempty"`
	PosterURL   string    `json:"poster_url,omitempty"` // URL to download poster from
}

// MovieUpdateRequest represents the input for updating a movie
type MovieUpdateRequest struct {
	ID          int       `json:"id"`
	Title       *string   `json:"title,omitempty"`
	Director    *string   `json:"director,omitempty"`
	Year        *int      `json:"year,omitempty"`
	Genre       []string  `json:"genre,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	Description *string   `json:"description,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	Language    *string   `json:"language,omitempty"`
	Country     *string   `json:"country,omitempty"`
	PosterURL   *string   `json:"poster_url,omitempty"`
}

// Validate checks if the movie create request is valid
func (req *MovieCreateRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(req.Director) == "" {
		return fmt.Errorf("director is required")
	}
	if req.Year < 1888 || req.Year > 2100 {
		return fmt.Errorf("year must be between 1888 and 2100")
	}
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 10) {
		return fmt.Errorf("rating must be between 0 and 10")
	}
	if req.Duration != nil && *req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// Validate checks if the movie update request is valid
func (req *MovieUpdateRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("id is required and must be positive")
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.Director != nil && strings.TrimSpace(*req.Director) == "" {
		return fmt.Errorf("director cannot be empty")
	}
	if req.Year != nil && (*req.Year < 1888 || *req.Year > 2100) {
		return fmt.Errorf("year must be between 1888 and 2100")
	}
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 10) {
		return fmt.Errorf("rating must be between 0 and 10")
	}
	if req.Duration != nil && *req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// ToJSON converts a movie to a formatted JSON string for MCP responses
func (m *Movie) ToJSON() string {
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data)
}

// ToSummary creates a brief text summary of the movie
func (m *Movie) ToSummary() string {
	summary := fmt.Sprintf("%s (%d) - Directed by %s", m.Title, m.Year, m.Director)
	if m.Rating > 0 {
		summary += fmt.Sprintf(" - Rating: %.1f/10", m.Rating)
	}
	if len(m.Genre) > 0 {
		summary += fmt.Sprintf(" - Genre: %s", strings.Join(m.Genre, ", "))
	}
	if m.Duration > 0 {
		hours := m.Duration / 60
		minutes := m.Duration % 60
		if hours > 0 {
			summary += fmt.Sprintf(" - Duration: %dh %dm", hours, minutes)
		} else {
			summary += fmt.Sprintf(" - Duration: %dm", minutes)
		}
	}
	return summary
}

// ParseMovieArguments extracts movie creation arguments from MCP tool call arguments
func ParseMovieArguments(args map[string]interface{}) (*MovieCreateRequest, error) {
	req := &MovieCreateRequest{}

	// Required fields
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title is required and must be a string")
	}
	req.Title = title

	director, ok := args["director"].(string)
	if !ok || director == "" {
		return nil, fmt.Errorf("director is required and must be a string")
	}
	req.Director = director

	year, ok := args["year"]
	if !ok {
		return nil, fmt.Errorf("year is required")
	}
	switch v := year.(type) {
	case float64:
		req.Year = int(v)
	case int:
		req.Year = v
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("year must be a valid integer")
		}
		req.Year = parsed
	default:
		return nil, fmt.Errorf("year must be an integer")
	}

	// Optional fields
	if genre, ok := args["genre"]; ok {
		switch v := genre.(type) {
		case []interface{}:
			for _, g := range v {
				if str, ok := g.(string); ok {
					req.Genre = append(req.Genre, str)
				}
			}
		case []string:
			req.Genre = v
		case string:
			// Single genre as string
			req.Genre = []string{v}
		}
	}

	if rating, ok := args["rating"]; ok {
		switch v := rating.(type) {
		case float64:
			req.Rating = &v
		case int:
			f := float64(v)
			req.Rating = &f
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("rating must be a valid number")
			}
			req.Rating = &parsed
		}
	}

	if description, ok := args["description"].(string); ok {
		req.Description = description
	}

	if duration, ok := args["duration"]; ok {
		switch v := duration.(type) {
		case float64:
			d := int(v)
			req.Duration = &d
		case int:
			req.Duration = &v
		case string:
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("duration must be a valid integer")
			}
			req.Duration = &parsed
		}
	}

	if language, ok := args["language"].(string); ok {
		req.Language = language
	}

	if country, ok := args["country"].(string); ok {
		req.Country = country
	}

	if posterURL, ok := args["poster_url"].(string); ok {
		req.PosterURL = posterURL
	}

	return req, nil
}

// ParseUpdateArguments extracts movie update arguments from MCP tool call arguments
func ParseUpdateArguments(args map[string]interface{}) (*MovieUpdateRequest, error) {
	req := &MovieUpdateRequest{}

	// Required ID
	id, ok := args["id"]
	if !ok {
		return nil, fmt.Errorf("id is required")
	}
	switch v := id.(type) {
	case float64:
		req.ID = int(v)
	case int:
		req.ID = v
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("id must be a valid integer")
		}
		req.ID = parsed
	default:
		return nil, fmt.Errorf("id must be an integer")
	}

	// Optional fields - only set if provided
	if title, ok := args["title"].(string); ok {
		req.Title = &title
	}

	if director, ok := args["director"].(string); ok {
		req.Director = &director
	}

	if year, ok := args["year"]; ok {
		switch v := year.(type) {
		case float64:
			y := int(v)
			req.Year = &y
		case int:
			req.Year = &v
		case string:
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("year must be a valid integer")
			}
			req.Year = &parsed
		}
	}

	if genre, ok := args["genre"]; ok {
		switch v := genre.(type) {
		case []interface{}:
			for _, g := range v {
				if str, ok := g.(string); ok {
					req.Genre = append(req.Genre, str)
				}
			}
		case []string:
			req.Genre = v
		case string:
			req.Genre = []string{v}
		}
	}

	if rating, ok := args["rating"]; ok {
		switch v := rating.(type) {
		case float64:
			req.Rating = &v
		case int:
			f := float64(v)
			req.Rating = &f
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("rating must be a valid number")
			}
			req.Rating = &parsed
		}
	}

	if description, ok := args["description"].(string); ok {
		req.Description = &description
	}

	if duration, ok := args["duration"]; ok {
		switch v := duration.(type) {
		case float64:
			d := int(v)
			req.Duration = &d
		case int:
			req.Duration = &v
		case string:
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("duration must be a valid integer")
			}
			req.Duration = &parsed
		}
	}

	if language, ok := args["language"].(string); ok {
		req.Language = &language
	}

	if country, ok := args["country"].(string); ok {
		req.Country = &country
	}

	if posterURL, ok := args["poster_url"].(string); ok {
		req.PosterURL = &posterURL
	}

	return req, nil
}