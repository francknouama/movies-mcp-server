package database

import (
	"fmt"
	"time"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	movies      map[int]*Movie
	nextID      int
	shouldError bool
	errorMsg    string
}

// NewMockDatabase creates a new mock database
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		movies: make(map[int]*Movie),
		nextID: 1,
	}
}

// SetError configures the mock to return an error
func (m *MockDatabase) SetError(shouldError bool, errorMsg string) {
	m.shouldError = shouldError
	m.errorMsg = errorMsg
}

// AddTestMovie adds a movie directly to the mock for testing
func (m *MockDatabase) AddTestMovie(movie *Movie) {
	if movie.ID == 0 {
		movie.ID = m.nextID
		m.nextID++
	}
	if movie.CreatedAt.IsZero() {
		movie.CreatedAt = time.Now()
	}
	if movie.UpdatedAt.IsZero() {
		movie.UpdatedAt = time.Now()
	}
	m.movies[movie.ID] = movie
}

// CreateMovie implements Database.CreateMovie
func (m *MockDatabase) CreateMovie(movie *Movie) error {
	if m.shouldError {
		return fmt.Errorf(m.errorMsg)
	}
	
	movie.ID = m.nextID
	m.nextID++
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	
	m.movies[movie.ID] = movie
	return nil
}

// GetMovie implements Database.GetMovie
func (m *MockDatabase) GetMovie(id int) (*Movie, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	movie, exists := m.movies[id]
	if !exists {
		return nil, fmt.Errorf("movie not found")
	}
	
	// Return a copy to prevent test interference
	copy := *movie
	return &copy, nil
}

// UpdateMovie implements Database.UpdateMovie
func (m *MockDatabase) UpdateMovie(movie *Movie) error {
	if m.shouldError {
		return fmt.Errorf(m.errorMsg)
	}
	
	if _, exists := m.movies[movie.ID]; !exists {
		return fmt.Errorf("movie not found")
	}
	
	movie.UpdatedAt = time.Now()
	m.movies[movie.ID] = movie
	return nil
}

// DeleteMovie implements Database.DeleteMovie
func (m *MockDatabase) DeleteMovie(id int) error {
	if m.shouldError {
		return fmt.Errorf(m.errorMsg)
	}
	
	if _, exists := m.movies[id]; !exists {
		return fmt.Errorf("movie not found")
	}
	
	delete(m.movies, id)
	return nil
}

// SearchMovies implements Database.SearchMovies
func (m *MockDatabase) SearchMovies(sq SearchQuery) ([]*Movie, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	var results []*Movie
	
	for _, movie := range m.movies {
		match := false
		
		switch sq.Type {
		case "title":
			if contains(movie.Title, sq.Query) {
				match = true
			}
		case "director":
			if contains(movie.Director, sq.Query) {
				match = true
			}
		case "genre":
			for _, genre := range movie.Genre {
				if genre == sq.Query {
					match = true
					break
				}
			}
		case "year":
			if fmt.Sprintf("%d", movie.Year) == sq.Query {
				match = true
			}
		case "fulltext":
			// Simple fulltext simulation
			description := ""
			if movie.Description.Valid {
				description = movie.Description.String
			}
			if contains(movie.Title, sq.Query) || 
			   contains(movie.Director, sq.Query) || 
			   contains(description, sq.Query) {
				match = true
			}
		default:
			// Default to title search
			if contains(movie.Title, sq.Query) {
				match = true
			}
		}
		
		if match {
			copy := *movie
			results = append(results, &copy)
		}
	}
	
	// Apply sorting (simplified - just by rating desc)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			ratingJ := results[j].Rating.Float64
			if !results[j].Rating.Valid {
				ratingJ = 0.0
			}
			ratingI := results[i].Rating.Float64
			if !results[i].Rating.Valid {
				ratingI = 0.0
			}
			if ratingJ > ratingI {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	
	// Apply pagination
	start := sq.Offset
	end := sq.Offset + sq.Limit
	if start > len(results) {
		return []*Movie{}, nil
	}
	if end > len(results) {
		end = len(results)
	}
	
	return results[start:end], nil
}

// ListTopMovies implements Database.ListTopMovies
func (m *MockDatabase) ListTopMovies(limit int, genreFilter string) ([]*Movie, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	var results []*Movie
	
	for _, movie := range m.movies {
		if genreFilter != "" {
			hasGenre := false
			for _, genre := range movie.Genre {
				if genre == genreFilter {
					hasGenre = true
					break
				}
			}
			if !hasGenre {
				continue
			}
		}
		
		copy := *movie
		results = append(results, &copy)
	}
	
	// Sort by rating descending
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			ratingJ := results[j].Rating.Float64
			if !results[j].Rating.Valid {
				ratingJ = 0.0
			}
			ratingI := results[i].Rating.Float64
			if !results[i].Rating.Valid {
				ratingI = 0.0
			}
			if ratingJ > ratingI {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	
	// Apply limit
	if limit < len(results) {
		results = results[:limit]
	}
	
	return results, nil
}

// GetMoviePoster implements Database.GetMoviePoster
func (m *MockDatabase) GetMoviePoster(id int) ([]byte, string, error) {
	if m.shouldError {
		return nil, "", fmt.Errorf(m.errorMsg)
	}
	
	movie, exists := m.movies[id]
	if !exists {
		return nil, "", fmt.Errorf("movie not found")
	}
	
	return movie.PosterData, movie.PosterType, nil
}

// UpdateMoviePoster implements Database.UpdateMoviePoster
func (m *MockDatabase) UpdateMoviePoster(id int, data []byte, mimeType string) error {
	if m.shouldError {
		return fmt.Errorf(m.errorMsg)
	}
	
	movie, exists := m.movies[id]
	if !exists {
		return fmt.Errorf("movie not found")
	}
	
	movie.PosterData = data
	movie.PosterType = mimeType
	movie.UpdatedAt = time.Now()
	return nil
}

// GetStats implements Database.GetStats
func (m *MockDatabase) GetStats() (*DatabaseStats, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	stats := &DatabaseStats{
		TotalMovies:   len(m.movies),
		MoviesPerYear: make(map[int]int),
		GenreCount:    make(map[string]int),
	}
	
	var totalRating float64
	var posterCount int
	
	for _, movie := range m.movies {
		if movie.Rating.Valid {
			totalRating += movie.Rating.Float64
		}
		stats.MoviesPerYear[movie.Year]++
		
		for _, genre := range movie.Genre {
			stats.GenreCount[genre]++
		}
		
		if len(movie.PosterData) > 0 {
			posterCount++
		}
	}
	
	if len(m.movies) > 0 {
		stats.AverageRating = totalRating / float64(len(m.movies))
	}
	
	stats.TotalPosters = posterCount
	stats.DatabaseSize = "1.2 MB" // Mock value
	
	return stats, nil
}

// GetAllMovies returns all movies (for database/all resource)
func (m *MockDatabase) GetAllMovies() ([]*Movie, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	var movies []*Movie
	for _, movie := range m.movies {
		copy := *movie
		movies = append(movies, &copy)
	}
	
	// Sort by ID for consistent results
	for i := 0; i < len(movies)-1; i++ {
		for j := i + 1; j < len(movies); j++ {
			if movies[j].ID < movies[i].ID {
				movies[i], movies[j] = movies[j], movies[i]
			}
		}
	}
	
	return movies, nil
}

// GetMoviesWithPosters returns all movies that have poster data
func (m *MockDatabase) GetMoviesWithPosters() ([]*Movie, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	var movies []*Movie
	for _, movie := range m.movies {
		if len(movie.PosterData) > 0 {
			copy := *movie
			movies = append(movies, &copy)
		}
	}
	
	return movies, nil
}

// GetGenres implements Database.GetGenres
func (m *MockDatabase) GetGenres() ([]GenreCount, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	genreMap := make(map[string]int)
	for _, movie := range m.movies {
		for _, genre := range movie.Genre {
			genreMap[genre]++
		}
	}
	
	var results []GenreCount
	for genre, count := range genreMap {
		results = append(results, GenreCount{
			Genre: genre,
			Count: count,
		})
	}
	
	return results, nil
}

// GetDirectors implements Database.GetDirectors
func (m *MockDatabase) GetDirectors() ([]DirectorCount, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	
	directorMap := make(map[string]struct {
		count  int
		rating float64
	})
	
	for _, movie := range m.movies {
		director := directorMap[movie.Director]
		director.count++
		if movie.Rating.Valid {
			director.rating += movie.Rating.Float64
		}
		directorMap[movie.Director] = director
	}
	
	var results []DirectorCount
	for name, data := range directorMap {
		results = append(results, DirectorCount{
			Director:      name,
			MovieCount:    data.count,
			AverageRating: data.rating / float64(data.count),
		})
	}
	
	return results, nil
}

// Ping implements Database.Ping
func (m *MockDatabase) Ping() error {
	if m.shouldError {
		return fmt.Errorf(m.errorMsg)
	}
	return nil
}

// Close implements Database.Close
func (m *MockDatabase) Close() error {
	return nil
}

// Helper function for case-insensitive contains
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		   (s == substr || 
		    len(s) >= len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr ||
		     findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}