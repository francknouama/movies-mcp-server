package tools

import (
	"context"
	"errors"
	"testing"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// ===== CreateSearchContext Tests =====

func TestCreateSearchContext_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Return 100 movies for pagination testing
			movies := make([]*movieApp.MovieDTO, 100)
			for i := 0; i < 100; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    "Movie " + string(rune(i+1)),
					Director: "Director",
					Year:     2000,
					Rating:   8.0,
					Genres:   []string{"Drama"},
				}
			}
			return movies, nil
		},
	}

	tools := NewContextTools(mockService)
	_, output, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "Movie",
		},
		PageSize: 25,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.ContextID == "" {
		t.Error("Expected context ID to be generated")
	}

	if output.Total != 100 {
		t.Errorf("Expected total 100, got: %d", output.Total)
	}

	if output.PageSize != 25 {
		t.Errorf("Expected page size 25, got: %d", output.PageSize)
	}

	if output.TotalPages != 4 {
		t.Errorf("Expected 4 total pages (100/25), got: %d", output.TotalPages)
	}

	if output.CreatedAt == "" {
		t.Error("Expected created_at timestamp")
	}

	if output.ExpiresAt == "" {
		t.Error("Expected expires_at timestamp")
	}
}

func TestCreateSearchContext_DefaultPageSize(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{ID: 1, Title: "Movie 1", Director: "Director", Year: 2000, Rating: 8.0},
			}, nil
		},
	}

	tools := NewContextTools(mockService)
	_, output, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "Movie",
		},
		// No PageSize specified - should default to 50
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.PageSize != 50 {
		t.Errorf("Expected default page size 50, got: %d", output.PageSize)
	}
}

func TestCreateSearchContext_MaxPageSize(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{ID: 1, Title: "Movie 1", Director: "Director", Year: 2000, Rating: 8.0},
			}, nil
		},
	}

	tools := NewContextTools(mockService)
	_, output, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "Movie",
		},
		PageSize: 2000, // Above max - should be capped at 1000
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.PageSize != 1000 {
		t.Errorf("Expected max page size 1000, got: %d", output.PageSize)
	}
}

func TestCreateSearchContext_NoResults(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewContextTools(mockService)
	_, output, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "NonExistent",
		},
		PageSize: 25,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Total != 0 {
		t.Errorf("Expected total 0, got: %d", output.Total)
	}

	if output.TotalPages != 0 {
		t.Errorf("Expected 0 pages, got: %d", output.TotalPages)
	}
}

func TestCreateSearchContext_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}

	tools := NewContextTools(mockService)
	_, _, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "Movie",
		},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// ===== GetContextPage Tests =====

func TestGetContextPage_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Return 100 movies
			movies := make([]*movieApp.MovieDTO, 100)
			for i := 0; i < 100; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    "Movie " + string(rune(i+1)),
					Director: "Director",
					Year:     2000 + i,
					Rating:   8.0,
					Genres:   []string{"Drama"},
				}
			}
			return movies, nil
		},
	}

	tools := NewContextTools(mockService)

	// First create a context
	_, createOutput, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query: SearchMoviesInput{
			Title: "Movie",
		},
		PageSize: 25,
	})

	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// Now get page 1
	_, pageOutput, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      1,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(pageOutput.Data) != 25 {
		t.Errorf("Expected 25 results on page 1, got: %d", len(pageOutput.Data))
	}

	if pageOutput.Page != 1 {
		t.Errorf("Expected page 1, got: %d", pageOutput.Page)
	}

	if pageOutput.TotalPages != 4 {
		t.Errorf("Expected 4 total pages, got: %d", pageOutput.TotalPages)
	}

	if pageOutput.HasNext != true {
		t.Error("Expected has_next to be true for page 1 of 4")
	}

	if pageOutput.HasPrevious != false {
		t.Error("Expected has_previous to be false for page 1")
	}
}

func TestGetContextPage_LastPage(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Return 100 movies
			movies := make([]*movieApp.MovieDTO, 100)
			for i := 0; i < 100; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    "Movie",
					Director: "Director",
					Year:     2000,
					Rating:   8.0,
				}
			}
			return movies, nil
		},
	}

	tools := NewContextTools(mockService)

	// Create context with 25 items per page
	_, createOutput, _ := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query:    SearchMoviesInput{Title: "Movie"},
		PageSize: 25,
	})

	// Get last page (page 4 of 4)
	_, pageOutput, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      4,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if pageOutput.HasNext != false {
		t.Error("Expected has_next to be false for last page")
	}

	if pageOutput.HasPrevious != true {
		t.Error("Expected has_previous to be true for page 4")
	}
}

func TestGetContextPage_InvalidContextID(t *testing.T) {
	mockService := &MockMovieService{}
	tools := NewContextTools(mockService)

	_, _, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: "invalid-context-id",
		Page:      1,
	})

	if err == nil {
		t.Fatal("Expected error for invalid context ID, got nil")
	}
}

func TestGetContextPage_InvalidPageNumber(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{ID: 1, Title: "Movie", Director: "Director", Year: 2000, Rating: 8.0},
			}, nil
		},
	}

	tools := NewContextTools(mockService)

	// Create context
	_, createOutput, _ := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query:    SearchMoviesInput{Title: "Movie"},
		PageSize: 25,
	})

	// Try to get page 0 - implementation clamps to page 1
	_, output, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      0,
	})

	if err != nil {
		t.Fatalf("Expected no error (page clamped to 1), got: %v", err)
	}

	if output.Page != 1 {
		t.Errorf("Expected page to be clamped to 1, got: %d", output.Page)
	}

	// Try to get page beyond total pages - implementation clamps to last page
	_, output, err = tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      100,
	})

	if err != nil {
		t.Fatalf("Expected no error (page clamped to last), got: %v", err)
	}

	if output.Page > output.TotalPages {
		t.Errorf("Expected page to be clamped to total pages, got page %d of %d", output.Page, output.TotalPages)
	}
}

// ===== GetContextInfo Tests =====

func TestGetContextInfo_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			movies := make([]*movieApp.MovieDTO, 50)
			for i := 0; i < 50; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    "Movie",
					Director: "Director",
					Year:     2000,
					Rating:   8.0,
				}
			}
			return movies, nil
		},
	}

	tools := NewContextTools(mockService)

	// Create context
	_, createOutput, _ := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query:    SearchMoviesInput{Title: "Movie"},
		PageSize: 10,
	})

	// Get context info
	_, infoOutput, err := tools.GetContextInfo(context.Background(), nil, GetContextInfoInput{
		ContextID: createOutput.ContextID,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if infoOutput.ContextID != createOutput.ContextID {
		t.Errorf("Expected context ID %s, got: %s", createOutput.ContextID, infoOutput.ContextID)
	}

	if infoOutput.Total != 50 {
		t.Errorf("Expected total 50, got: %d", infoOutput.Total)
	}

	if infoOutput.PageSize != 10 {
		t.Errorf("Expected page size 10, got: %d", infoOutput.PageSize)
	}

	if infoOutput.TotalPages != 5 {
		t.Errorf("Expected 5 total pages (50/10), got: %d", infoOutput.TotalPages)
	}

	if infoOutput.CreatedAt == "" {
		t.Error("Expected created_at timestamp")
	}

	if infoOutput.ExpiresAt == "" {
		t.Error("Expected expires_at timestamp")
	}
}

func TestGetContextInfo_InvalidContextID(t *testing.T) {
	mockService := &MockMovieService{}
	tools := NewContextTools(mockService)

	_, _, err := tools.GetContextInfo(context.Background(), nil, GetContextInfoInput{
		ContextID: "invalid-context-id",
	})

	if err == nil {
		t.Fatal("Expected error for invalid context ID, got nil")
	}
}

// ===== Integration Test: Full Workflow =====

func TestContextTools_FullWorkflow(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			movies := make([]*movieApp.MovieDTO, 75)
			for i := 0; i < 75; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    "Movie " + string(rune(i+1)),
					Director: "Director",
					Year:     2000 + (i % 20),
					Rating:   7.0 + float64(i%30)/10.0,
					Genres:   []string{"Drama"},
				}
			}
			return movies, nil
		},
	}

	tools := NewContextTools(mockService)

	// Step 1: Create search context
	_, createOutput, err := tools.CreateSearchContext(context.Background(), nil, CreateSearchContextInput{
		Query:    SearchMoviesInput{Title: "Movie"},
		PageSize: 20,
	})

	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// Step 2: Get context info
	_, infoOutput, err := tools.GetContextInfo(context.Background(), nil, GetContextInfoInput{
		ContextID: createOutput.ContextID,
	})

	if err != nil {
		t.Fatalf("Failed to get context info: %v", err)
	}

	if infoOutput.Total != 75 {
		t.Errorf("Expected total 75, got: %d", infoOutput.Total)
	}

	if infoOutput.TotalPages != 4 {
		t.Errorf("Expected 4 pages (ceil(75/20)), got: %d", infoOutput.TotalPages)
	}

	// Step 3: Get page 1
	_, page1, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      1,
	})

	if err != nil {
		t.Fatalf("Failed to get page 1: %v", err)
	}

	if len(page1.Data) != 20 {
		t.Errorf("Expected 20 results on page 1, got: %d", len(page1.Data))
	}

	// Step 4: Get page 2
	_, page2, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      2,
	})

	if err != nil {
		t.Fatalf("Failed to get page 2: %v", err)
	}

	if len(page2.Data) != 20 {
		t.Errorf("Expected 20 results on page 2, got: %d", len(page2.Data))
	}

	// Step 5: Get last page (should have 15 items: 75 - 20 - 20 - 20 = 15)
	_, page4, err := tools.GetContextPage(context.Background(), nil, GetContextPageInput{
		ContextID: createOutput.ContextID,
		Page:      4,
	})

	if err != nil {
		t.Fatalf("Failed to get page 4: %v", err)
	}

	if len(page4.Data) != 15 {
		t.Errorf("Expected 15 results on last page, got: %d", len(page4.Data))
	}

	if page4.HasNext != false {
		t.Error("Expected has_next to be false on last page")
	}
}
