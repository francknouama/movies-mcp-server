package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"movies-mcp-server/internal/database"
	"movies-mcp-server/internal/models"
)

// mockWriter captures output for testing
type mockWriter struct {
	buffer bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

// Helper function to create a test server
func createTestServer() (*MoviesServer, *database.MockDatabase, *mockWriter) {
	db := database.NewMockDatabase()
	output := &mockWriter{}
	server := &MoviesServer{
		input:          strings.NewReader(""),
		output:         output,
		logger:         nil, // Disable logging in tests
		db:             db,
		imageProcessor: nil, // Not needed for most tests
	}
	return server, db, output
}

// Helper function to parse JSON response
func parseResponse(t *testing.T, output string) models.JSONRPCResponse {
	var response models.JSONRPCResponse
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return response
}

// Test search_movies handler
func TestHandleSearchMovies(t *testing.T) {
	server, db, output := createTestServer()

	// Add test data
	db.AddTestMovie(&database.Movie{
		ID:          1,
		Title:       "Inception",
		Director:    "Christopher Nolan",
		Year:        2010,
		Genre:       []string{"Action", "Sci-Fi", "Thriller"},
		Rating:      sql.NullFloat64{Float64: 8.8, Valid: true},
		Description: sql.NullString{String: "A thief who steals corporate secrets through dream-sharing technology", Valid: true},
	})
	db.AddTestMovie(&database.Movie{
		ID:          2,
		Title:       "The Dark Knight",
		Director:    "Christopher Nolan",
		Year:        2008,
		Genre:       []string{"Action", "Crime", "Drama"},
		Rating:      sql.NullFloat64{Float64: 9.0, Valid: true},
		Description: sql.NullString{String: "Batman faces the Joker", Valid: true},
	})
	db.AddTestMovie(&database.Movie{
		ID:       3,
		Title:    "Interstellar",
		Director: "Christopher Nolan",
		Year:     2014,
		Genre:    []string{"Adventure", "Drama", "Sci-Fi"},
		Rating:   sql.NullFloat64{Float64: 8.6, Valid: true},
	})

	tests := []struct {
		name         string
		args         map[string]interface{}
		expectError  bool
		errorMessage string
		expectCount  int
		checkContent func(string) bool
	}{
		{
			name: "Search by title",
			args: map[string]interface{}{
				"query": "Inception",
				"type":  "title",
			},
			expectError: false,
			expectCount: 1,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Inception") &&
					strings.Contains(content, "Found 1 movies")
			},
		},
		{
			name: "Search by director",
			args: map[string]interface{}{
				"query": "Nolan",
				"type":  "director",
			},
			expectError: false,
			expectCount: 3,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Found 3 movies") &&
					strings.Contains(content, "Christopher Nolan")
			},
		},
		{
			name: "Search by genre",
			args: map[string]interface{}{
				"query": "Sci-Fi",
				"type":  "genre",
			},
			expectError: false,
			expectCount: 2,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Found 2 movies")
			},
		},
		{
			name: "Search by year",
			args: map[string]interface{}{
				"query": "2010",
				"type":  "year",
			},
			expectError: false,
			expectCount: 1,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Inception")
			},
		},
		{
			name: "Full-text search",
			args: map[string]interface{}{
				"query": "Batman",
				"type":  "fulltext",
			},
			expectError: false,
			expectCount: 1,
			checkContent: func(content string) bool {
				return strings.Contains(content, "The Dark Knight")
			},
		},
		{
			name: "Search with pagination",
			args: map[string]interface{}{
				"query":  "Nolan",
				"type":   "director",
				"limit":  2,
				"offset": 1,
			},
			expectError: false,
			expectCount: 2,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Found 2 movies")
			},
		},
		{
			name:         "Missing query parameter",
			args:         map[string]interface{}{},
			expectError:  true,
			errorMessage: "At least one search parameter must be provided",
		},
		{
			name: "Empty query",
			args: map[string]interface{}{
				"query": "",
			},
			expectError:  true,
			errorMessage: "At least one search parameter must be provided",
		},
		{
			name: "No results found",
			args: map[string]interface{}{
				"query": "NonexistentMovie",
			},
			expectError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "No movies found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			server.handleSearchMovies(1, tt.args)

			response := parseResponse(t, output.buffer.String())

			if tt.expectError {
				if response.Error == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(response.Error.Message, tt.errorMessage) {
					t.Errorf("Expected error message '%s' but got '%s'", 
						tt.errorMessage, response.Error.Message)
				}
			} else {
				if response.Error != nil {
					t.Errorf("Unexpected error: %v", response.Error)
				} else if response.Result != nil {
					var result models.ToolCallResponse
					resultBytes, _ := json.Marshal(response.Result)
					json.Unmarshal(resultBytes, &result)
					
					if len(result.Content) > 0 && tt.checkContent != nil {
						if !tt.checkContent(result.Content[0].Text) {
							t.Errorf("Content check failed: %s", result.Content[0].Text)
						}
					}
				}
			}
		})
	}
}

// Test list_top_movies handler
func TestHandleListTopMovies(t *testing.T) {
	server, db, output := createTestServer()

	// Add test data
	db.AddTestMovie(&database.Movie{
		ID:       1,
		Title:    "The Shawshank Redemption",
		Director: "Frank Darabont",
		Year:     1994,
		Genre:    []string{"Drama"},
		Rating:   sql.NullFloat64{Float64: 9.3, Valid: true},
		Duration: sql.NullInt32{Int32: 142, Valid: true},
		Language: sql.NullString{String: "English", Valid: true},
	})
	db.AddTestMovie(&database.Movie{
		ID:       2,
		Title:    "The Godfather",
		Director: "Francis Ford Coppola",
		Year:     1972,
		Genre:    []string{"Crime", "Drama"},
		Rating:   sql.NullFloat64{Float64: 9.2, Valid: true},
		Duration: sql.NullInt32{Int32: 175, Valid: true},
	})
	db.AddTestMovie(&database.Movie{
		ID:       3,
		Title:    "The Dark Knight",
		Director: "Christopher Nolan",
		Year:     2008,
		Genre:    []string{"Action", "Crime", "Drama"},
		Rating:   sql.NullFloat64{Float64: 9.0, Valid: true},
	})
	db.AddTestMovie(&database.Movie{
		ID:       4,
		Title:    "Pulp Fiction",
		Director: "Quentin Tarantino",
		Year:     1994,
		Genre:    []string{"Crime", "Drama"},
		Rating:   sql.NullFloat64{Float64: 8.9, Valid: true},
	})

	tests := []struct {
		name         string
		args         map[string]interface{}
		expectError  bool
		errorMessage string
		checkContent func(string) bool
	}{
		{
			name: "List top movies default",
			args: map[string]interface{}{},
			checkContent: func(content string) bool {
				return strings.Contains(content, "Top") &&
					strings.Contains(content, "The Shawshank Redemption") &&
					strings.Contains(content, "9.3/10")
			},
		},
		{
			name: "List top movies with limit",
			args: map[string]interface{}{
				"limit": 2,
			},
			checkContent: func(content string) bool {
				return strings.Contains(content, "Top 2 movies") &&
					strings.Contains(content, "The Shawshank Redemption") &&
					strings.Contains(content, "The Godfather") &&
					!strings.Contains(content, "The Dark Knight")
			},
		},
		{
			name: "List top movies by genre",
			args: map[string]interface{}{
				"genre": "Drama",
			},
			checkContent: func(content string) bool {
				return strings.Contains(content, "Top") &&
					strings.Contains(content, "Drama movies") &&
					strings.Contains(content, "The Shawshank Redemption")
			},
		},
		{
			name: "List top movies with genre and limit",
			args: map[string]interface{}{
				"genre": "Crime",
				"limit": 2,
			},
			checkContent: func(content string) bool {
				return strings.Contains(content, "Top 2 Crime movies") &&
					strings.Contains(content, "The Godfather") &&
					!strings.Contains(content, "The Shawshank Redemption")
			},
		},
		{
			name: "No movies for genre",
			args: map[string]interface{}{
				"genre": "Horror",
			},
			checkContent: func(content string) bool {
				return strings.Contains(content, "No movies found for genre 'Horror'")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			server.handleListTopMovies(1, tt.args)

			response := parseResponse(t, output.buffer.String())

			if tt.expectError {
				if response.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if response.Error != nil {
					t.Errorf("Unexpected error: %v", response.Error)
				} else if response.Result != nil {
					var result models.ToolCallResponse
					resultBytes, _ := json.Marshal(response.Result)
					json.Unmarshal(resultBytes, &result)
					
					if len(result.Content) > 0 && tt.checkContent != nil {
						if !tt.checkContent(result.Content[0].Text) {
							t.Errorf("Content check failed: %s", result.Content[0].Text)
						}
					}
				}
			}
		})
	}
}

// Test error handling with database errors
func TestHandleSearchMoviesWithDBError(t *testing.T) {
	server, db, output := createTestServer()
	
	// Configure mock to return an error
	db.SetError(true, "database connection failed")
	
	server.handleSearchMovies(1, map[string]interface{}{
		"query": "test",
	})
	
	response := parseResponse(t, output.buffer.String())
	
	if response.Error == nil {
		t.Errorf("Expected error but got none")
	} else if !strings.Contains(response.Error.Message, "Search failed") {
		t.Errorf("Expected 'Search failed' error but got: %s", response.Error.Message)
	}
}

// Test edge cases
func TestSearchMoviesEdgeCases(t *testing.T) {
	server, db, output := createTestServer()
	
	// Add a movie with special characters
	db.AddTestMovie(&database.Movie{
		ID:          1,
		Title:       "Test & Movie: Part 2",
		Director:    "Director O'Neill",
		Year:        2020,
		Genre:       []string{"Action"},
		Rating:      sql.NullFloat64{Float64: 7.5, Valid: true},
		Description: sql.NullString{String: "A movie with special characters & symbols", Valid: true},
	})
	
	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Search with special characters",
			args: map[string]interface{}{
				"query": "Test & Movie",
				"type":  "title",
			},
		},
		{
			name: "Search with apostrophe",
			args: map[string]interface{}{
				"query": "O'Neill",
				"type":  "director",
			},
		},
		{
			name: "Very large limit",
			args: map[string]interface{}{
				"query": "Test",
				"limit": 1000,
			},
		},
		{
			name: "Invalid search type defaults to title",
			args: map[string]interface{}{
				"query": "Test",
				"type":  "invalid_type",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			server.handleSearchMovies(1, tt.args)
			
			response := parseResponse(t, output.buffer.String())
			if response.Error != nil {
				t.Errorf("Unexpected error: %v", response.Error)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSearchMovies(b *testing.B) {
	server, db, _ := createTestServer()
	
	// Add 100 test movies
	for i := 1; i <= 100; i++ {
		db.AddTestMovie(&database.Movie{
			ID:       i,
			Title:    fmt.Sprintf("Movie %d", i),
			Director: fmt.Sprintf("Director %d", i%10),
			Year:     2000 + i%20,
			Genre:    []string{"Drama", "Action"},
			Rating:   sql.NullFloat64{Float64: float64(i%10) / 2, Valid: true},
		})
	}
	
	args := map[string]interface{}{
		"query": "Movie",
		"type":  "title",
		"limit": 10,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.handleSearchMovies(1, args)
	}
}

func BenchmarkListTopMovies(b *testing.B) {
	server, db, _ := createTestServer()
	
	// Add 100 test movies
	for i := 1; i <= 100; i++ {
		db.AddTestMovie(&database.Movie{
			ID:       i,
			Title:    fmt.Sprintf("Movie %d", i),
			Director: fmt.Sprintf("Director %d", i%10),
			Year:     2000 + i%20,
			Genre:    []string{"Drama", "Action"},
			Rating:   sql.NullFloat64{Float64: float64(i%10) / 2, Valid: true},
		})
	}
	
	args := map[string]interface{}{
		"limit": 10,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.handleListTopMovies(1, args)
	}
}