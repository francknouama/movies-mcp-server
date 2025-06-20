package server

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"movies-mcp-server/internal/database"
	"movies-mcp-server/internal/models"
)

// Test resources/list handler
func TestHandleResourcesList(t *testing.T) {
	server, db, output := createTestServer()
	_ = db // db not needed for this test

	// Create a mock request
	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/list",
		Params: json.RawMessage(`{}`),
	}

	server.handleResourcesList(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	// Parse the result
	var result models.ResourcesListResponse
	resultBytes, _ := json.Marshal(response.Result)
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		t.Fatalf("Failed to parse resources list result: %v", err)
	}

	// Expected resources based on Phase 5 requirements
	expectedResources := []string{
		"movies://database/all",
		"movies://database/stats", 
		"movies://database/genres",
		"movies://database/directors",
		"movies://posters/collection",
	}

	if len(result.Resources) < len(expectedResources) {
		t.Errorf("Expected at least %d resources, got %d", len(expectedResources), len(result.Resources))
	}

	// Check that expected resources are present
	resourceURIs := make(map[string]bool)
	for _, resource := range result.Resources {
		resourceURIs[resource.URI] = true
		
		// Validate resource structure
		if resource.URI == "" {
			t.Errorf("Resource URI cannot be empty")
		}
		if resource.Name == "" {
			t.Errorf("Resource name cannot be empty for URI: %s", resource.URI)
		}
		if resource.MimeType == "" {
			t.Errorf("Resource mime type cannot be empty for URI: %s", resource.URI)
		}
	}

	for _, expectedURI := range expectedResources {
		if !resourceURIs[expectedURI] {
			t.Errorf("Expected resource URI %s not found", expectedURI)
		}
	}
}

// Test resources/read handler for database/all resource
func TestHandleResourceRead_DatabaseAll(t *testing.T) {
	server, db, output := createTestServer()

	// Add test movies
	testMovies := []*database.Movie{
		{
			ID:       1,
			Title:    "Test Movie 1",
			Director: "Director 1",
			Year:     2020,
			Genre:    []string{"Action"},
			Rating:   sql.NullFloat64{Float64: 8.0, Valid: true},
		},
		{
			ID:       2,
			Title:    "Test Movie 2", 
			Director: "Director 2",
			Year:     2021,
			Genre:    []string{"Drama"},
			Rating:   sql.NullFloat64{Float64: 7.5, Valid: true},
		},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	// Test reading database/all resource
	params := models.ResourceReadRequest{
		URI: "movies://database/all",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	// Parse the result
	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		t.Fatalf("Failed to parse resource read result: %v", err)
	}

	if len(result.Contents) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(result.Contents))
	}

	content := result.Contents[0]
	if content.URI != "movies://database/all" {
		t.Errorf("Expected URI 'movies://database/all', got '%s'", content.URI)
	}

	if content.MimeType != "application/json" {
		t.Errorf("Expected mime type 'application/json', got '%s'", content.MimeType)
	}

	// Verify the content contains valid JSON with movies
	if content.Text == "" {
		t.Errorf("Content text should not be empty")
	}

	// Try to parse as JSON array of movies
	var movies []map[string]interface{}
	if err := json.Unmarshal([]byte(content.Text), &movies); err != nil {
		t.Errorf("Content should be valid JSON: %v", err)
	}

	if len(movies) != len(testMovies) {
		t.Errorf("Expected %d movies in content, got %d", len(testMovies), len(movies))
	}
}

// Test resources/read handler for database/stats resource
func TestHandleResourceRead_DatabaseStats(t *testing.T) {
	server, db, output := createTestServer()

	// Add test data for statistics
	testMovies := []*database.Movie{
		{
			ID:       1,
			Title:    "Action Movie",
			Director: "Director A",
			Year:     2020,
			Genre:    []string{"Action"},
			Rating:   sql.NullFloat64{Float64: 8.5, Valid: true},
		},
		{
			ID:       2,
			Title:    "Drama Movie",
			Director: "Director B", 
			Year:     2020,
			Genre:    []string{"Drama"},
			Rating:   sql.NullFloat64{Float64: 7.0, Valid: true},
		},
		{
			ID:         3,
			Title:      "Comedy Movie",
			Director:   "Director A",
			Year:       2021,
			Genre:      []string{"Comedy"},
			Rating:     sql.NullFloat64{Float64: 6.5, Valid: true},
			PosterData: []byte("fake poster data"),
			PosterType: "image/jpeg",
		},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	params := models.ResourceReadRequest{
		URI: "movies://database/stats",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &result)

	content := result.Contents[0]
	
	// Parse stats JSON
	var stats database.DatabaseStats
	if err := json.Unmarshal([]byte(content.Text), &stats); err != nil {
		t.Errorf("Stats content should be valid JSON: %v", err)
	}

	// Verify statistics
	if stats.TotalMovies != 3 {
		t.Errorf("Expected 3 total movies, got %d", stats.TotalMovies)
	}

	if stats.TotalPosters != 1 {
		t.Errorf("Expected 1 poster, got %d", stats.TotalPosters)
	}

	expectedAverage := (8.5 + 7.0 + 6.5) / 3
	if stats.AverageRating != expectedAverage {
		t.Errorf("Expected average rating %.2f, got %.2f", expectedAverage, stats.AverageRating)
	}
}

// Test resources/read handler for genres resource  
func TestHandleResourceRead_Genres(t *testing.T) {
	server, db, output := createTestServer()

	// Add test movies with different genres
	testMovies := []*database.Movie{
		{ID: 1, Title: "Movie 1", Genre: []string{"Action", "Thriller"}},
		{ID: 2, Title: "Movie 2", Genre: []string{"Drama"}},
		{ID: 3, Title: "Movie 3", Genre: []string{"Action", "Comedy"}},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	params := models.ResourceReadRequest{
		URI: "movies://database/genres",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read", 
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &result)

	content := result.Contents[0]

	// Parse genres JSON
	var genres []database.GenreCount
	if err := json.Unmarshal([]byte(content.Text), &genres); err != nil {
		t.Errorf("Genres content should be valid JSON: %v", err)
	}

	// Should have 4 unique genres: Action, Thriller, Drama, Comedy
	expectedGenres := map[string]int{
		"Action":   2,
		"Thriller": 1,
		"Drama":    1,
		"Comedy":   1,
	}

	genreMap := make(map[string]int)
	for _, g := range genres {
		genreMap[g.Genre] = g.Count
	}

	for expectedGenre, expectedCount := range expectedGenres {
		if genreMap[expectedGenre] != expectedCount {
			t.Errorf("Expected %d %s movies, got %d", expectedCount, expectedGenre, genreMap[expectedGenre])
		}
	}
}

// Test resources/read handler for directors resource
func TestHandleResourceRead_Directors(t *testing.T) {
	server, db, output := createTestServer()

	// Add test movies with different directors
	testMovies := []*database.Movie{
		{ID: 1, Title: "Movie 1", Director: "Director A", Rating: sql.NullFloat64{Float64: 8.0, Valid: true}},
		{ID: 2, Title: "Movie 2", Director: "Director A", Rating: sql.NullFloat64{Float64: 7.0, Valid: true}},
		{ID: 3, Title: "Movie 3", Director: "Director B", Rating: sql.NullFloat64{Float64: 9.0, Valid: true}},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	params := models.ResourceReadRequest{
		URI: "movies://database/directors",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &result)

	content := result.Contents[0]

	// Parse directors JSON
	var directors []database.DirectorCount
	if err := json.Unmarshal([]byte(content.Text), &directors); err != nil {
		t.Errorf("Directors content should be valid JSON: %v", err)
	}

	// Should have 2 directors
	if len(directors) != 2 {
		t.Errorf("Expected 2 directors, got %d", len(directors))
	}

	directorMap := make(map[string]database.DirectorCount)
	for _, d := range directors {
		directorMap[d.Director] = d
	}

	// Director A: 2 movies, average rating 7.5
	if directorA, exists := directorMap["Director A"]; exists {
		if directorA.MovieCount != 2 {
			t.Errorf("Expected Director A to have 2 movies, got %d", directorA.MovieCount)
		}
		expectedAvg := (8.0 + 7.0) / 2
		if directorA.AverageRating != expectedAvg {
			t.Errorf("Expected Director A average rating %.1f, got %.1f", expectedAvg, directorA.AverageRating)
		}
	} else {
		t.Errorf("Director A not found in results")
	}

	// Director B: 1 movie, average rating 9.0
	if directorB, exists := directorMap["Director B"]; exists {
		if directorB.MovieCount != 1 {
			t.Errorf("Expected Director B to have 1 movie, got %d", directorB.MovieCount)
		}
		if directorB.AverageRating != 9.0 {
			t.Errorf("Expected Director B average rating 9.0, got %.1f", directorB.AverageRating)
		}
	} else {
		t.Errorf("Director B not found in results")
	}
}

// Test resources/read handler for individual poster resource
func TestHandleResourceRead_IndividualPoster(t *testing.T) {
	server, db, output := createTestServer()

	// Add test movie with poster
	posterData := []byte("fake jpeg data")
	testMovie := &database.Movie{
		ID:         1,
		Title:      "Movie With Poster",
		PosterData: posterData,
		PosterType: "image/jpeg",
	}
	db.AddTestMovie(testMovie)

	params := models.ResourceReadRequest{
		URI: "movies://posters/1",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &result)

	if len(result.Contents) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(result.Contents))
	}

	content := result.Contents[0]

	if content.MimeType != "image/jpeg" {
		t.Errorf("Expected mime type 'image/jpeg', got '%s'", content.MimeType)
	}

	// Should have base64 encoded blob data
	if content.Blob == "" {
		t.Errorf("Expected blob data to be present")
	}

	if content.Text != "" {
		t.Errorf("Expected text to be empty for binary resource")
	}
}

// Test resources/read handler for poster collection
func TestHandleResourceRead_PosterCollection(t *testing.T) {
	server, db, output := createTestServer()

	// Add test movies with posters
	testMovies := []*database.Movie{
		{
			ID:         1,
			Title:      "Movie 1",
			PosterData: []byte("poster1"),
			PosterType: "image/jpeg",
		},
		{
			ID:         2,
			Title:      "Movie 2",
			PosterData: []byte("poster2"),
			PosterType: "image/png",
		},
		{
			ID:    3,
			Title: "Movie Without Poster",
		},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	params := models.ResourceReadRequest{
		URI: "movies://posters/collection",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	var result models.ResourceReadResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &result)

	content := result.Contents[0]

	if content.MimeType != "application/json" {
		t.Errorf("Expected mime type 'application/json', got '%s'", content.MimeType)
	}

	// Parse poster collection JSON
	var posterCollection []map[string]interface{}
	if err := json.Unmarshal([]byte(content.Text), &posterCollection); err != nil {
		t.Errorf("Poster collection should be valid JSON: %v", err)
	}

	// Should only include movies with posters (2 out of 3)
	if len(posterCollection) != 2 {
		t.Errorf("Expected 2 movies with posters, got %d", len(posterCollection))
	}

	// Check structure of poster collection items
	for _, item := range posterCollection {
		if _, hasID := item["movie_id"]; !hasID {
			t.Errorf("Poster collection item should have movie_id")
		}
		if _, hasTitle := item["title"]; !hasTitle {
			t.Errorf("Poster collection item should have title")
		}
		if _, hasPosterURI := item["poster_uri"]; !hasPosterURI {
			t.Errorf("Poster collection item should have poster_uri")
		}
		if _, hasMimeType := item["mime_type"]; !hasMimeType {
			t.Errorf("Poster collection item should have mime_type")
		}
	}
}

// Test error handling for invalid resource URIs
func TestHandleResourceRead_InvalidURI(t *testing.T) {
	server, _, output := createTestServer()

	invalidURIs := []string{
		"invalid://uri",
		"movies://invalid/resource",
		"movies://posters/invalid_id",
		"",
		"not-a-uri",
	}

	for _, uri := range invalidURIs {
		t.Run("InvalidURI_"+uri, func(t *testing.T) {
			output.buffer.Reset()

			params := models.ResourceReadRequest{
				URI: uri,
			}
			paramsBytes, _ := json.Marshal(params)

			req := &models.JSONRPCRequest{
				ID:     1,
				Method: "resources/read",
				Params: paramsBytes,
			}

			server.handleResourceRead(req)

			response := parseResponse(t, output.buffer.String())

			if response.Error == nil {
				t.Errorf("Expected error for invalid URI '%s' but got none", uri)
			}
		})
	}
}

// Test error handling when resource not found
func TestHandleResourceRead_NotFound(t *testing.T) {
	server, _, output := createTestServer()

	params := models.ResourceReadRequest{
		URI: "movies://posters/999", // Non-existent movie ID
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error == nil {
		t.Errorf("Expected error for non-existent poster but got none")
	}

	if !strings.Contains(response.Error.Message, "not found") {
		t.Errorf("Expected 'not found' error message, got: %s", response.Error.Message)
	}
}

// Test database error handling
func TestHandleResourceRead_DatabaseError(t *testing.T) {
	server, db, output := createTestServer()

	// Configure mock to return error
	db.SetError(true, "database connection failed")

	params := models.ResourceReadRequest{
		URI: "movies://database/stats",
	}
	paramsBytes, _ := json.Marshal(params)

	req := &models.JSONRPCRequest{
		ID:     1,
		Method: "resources/read",
		Params: paramsBytes,
	}

	server.handleResourceRead(req)

	response := parseResponse(t, output.buffer.String())

	if response.Error == nil {
		t.Errorf("Expected error due to database failure but got none")
	}
}