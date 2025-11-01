package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// ===== BulkMovieImport Tests =====

func TestBulkMovieImport_AllSuccess(t *testing.T) {
	callCount := 0
	mockService := &MockMovieService{
		CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
			callCount++
			return &movieApp.MovieDTO{
				ID:       callCount,
				Title:    cmd.Title,
				Director: cmd.Director,
				Year:     cmd.Year,
				Rating:   cmd.Rating,
				Genres:   cmd.Genres,
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.BulkMovieImport(context.Background(), nil, BulkMovieImportInput{
		Movies: []MovieImportItem{
			{Title: "The Matrix", Director: "Wachowskis", Year: 1999, Rating: 8.7, Genres: []string{"Sci-Fi"}},
			{Title: "Inception", Director: "Nolan", Year: 2010, Rating: 8.8, Genres: []string{"Sci-Fi", "Thriller"}},
			{Title: "Pulp Fiction", Director: "Tarantino", Year: 1994, Rating: 8.9, Genres: []string{"Crime"}},
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Imported != 3 {
		t.Errorf("Expected 3 imported, got: %d", output.Imported)
	}

	if output.Failed != 0 {
		t.Errorf("Expected 0 failed, got: %d", output.Failed)
	}

	if output.Total != 3 {
		t.Errorf("Expected 3 total, got: %d", output.Total)
	}

	if output.SuccessRate != "100.0%" {
		t.Errorf("Expected success rate '100.0%%', got: %s", output.SuccessRate)
	}

	if len(output.Results) != 3 {
		t.Errorf("Expected 3 results, got: %d", len(output.Results))
	}

	if len(output.Errors) != 0 {
		t.Errorf("Expected 0 errors, got: %d", len(output.Errors))
	}
}

func TestBulkMovieImport_PartialSuccess(t *testing.T) {
	callCount := 0
	mockService := &MockMovieService{
		CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
			callCount++
			// Fail the second movie
			if callCount == 2 {
				return nil, errors.New("invalid rating")
			}
			return &movieApp.MovieDTO{
				ID:       callCount,
				Title:    cmd.Title,
				Director: cmd.Director,
				Year:     cmd.Year,
				Rating:   cmd.Rating,
				Genres:   cmd.Genres,
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.BulkMovieImport(context.Background(), nil, BulkMovieImportInput{
		Movies: []MovieImportItem{
			{Title: "Movie 1", Director: "Director 1", Year: 2000, Rating: 7.0},
			{Title: "Movie 2", Director: "Director 2", Year: 2001, Rating: 8.0},
			{Title: "Movie 3", Director: "Director 3", Year: 2002, Rating: 9.0},
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Imported != 2 {
		t.Errorf("Expected 2 imported, got: %d", output.Imported)
	}

	if output.Failed != 1 {
		t.Errorf("Expected 1 failed, got: %d", output.Failed)
	}

	if output.SuccessRate != "66.7%" {
		t.Errorf("Expected success rate '66.7%%', got: %s", output.SuccessRate)
	}

	if len(output.Errors) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(output.Errors))
	}

	if output.Errors[0].Index != 1 {
		t.Errorf("Expected error at index 1, got: %d", output.Errors[0].Index)
	}

	if output.Errors[0].Title != "Movie 2" {
		t.Errorf("Expected error for 'Movie 2', got: %s", output.Errors[0].Title)
	}
}

func TestBulkMovieImport_AllFailed(t *testing.T) {
	mockService := &MockMovieService{
		CreateMovieFunc: func(ctx context.Context, cmd movieApp.CreateMovieCommand) (*movieApp.MovieDTO, error) {
			return nil, errors.New("database error")
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.BulkMovieImport(context.Background(), nil, BulkMovieImportInput{
		Movies: []MovieImportItem{
			{Title: "Movie 1", Director: "Director 1", Year: 2000, Rating: 7.0},
			{Title: "Movie 2", Director: "Director 2", Year: 2001, Rating: 8.0},
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Imported != 0 {
		t.Errorf("Expected 0 imported, got: %d", output.Imported)
	}

	if output.Failed != 2 {
		t.Errorf("Expected 2 failed, got: %d", output.Failed)
	}

	if output.SuccessRate != "0.0%" {
		t.Errorf("Expected success rate '0.0%%', got: %s", output.SuccessRate)
	}
}

// ===== MovieRecommendationEngine Tests =====

func TestMovieRecommendationEngine_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "The Matrix",
					Director: "Wachowskis",
					Year:     1999,
					Rating:   8.7,
					Genres:   []string{"Sci-Fi", "Action"},
				},
				{
					ID:       2,
					Title:    "Blade Runner",
					Director: "Ridley Scott",
					Year:     1982,
					Rating:   8.1,
					Genres:   []string{"Sci-Fi"},
				},
				{
					ID:       3,
					Title:    "Inception",
					Director: "Christopher Nolan",
					Year:     2010,
					Rating:   8.8,
					Genres:   []string{"Sci-Fi", "Thriller"},
				},
				{
					ID:       4,
					Title:    "The Godfather",
					Director: "Francis Ford Coppola",
					Year:     1972,
					Rating:   9.2,
					Genres:   []string{"Crime", "Drama"},
				},
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.MovieRecommendationEngine(context.Background(), nil, MovieRecommendationInput{
		Preferences: UserPreferences{
			Genres:    []string{"Sci-Fi", "Action"},
			MinRating: 8.0,
		},
		Limit: 3,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Recommendations) == 0 {
		t.Fatal("Expected recommendations, got none")
	}

	// Should be sorted by rank
	for i := 0; i < len(output.Recommendations)-1; i++ {
		if output.Recommendations[i].Rank > output.Recommendations[i+1].Rank {
			t.Errorf("Recommendations not sorted by rank: %d > %d",
				output.Recommendations[i].Rank, output.Recommendations[i+1].Rank)
		}
	}

	// Verify The Matrix is in recommendations (has both Sci-Fi and Action)
	foundMatrix := false
	for _, rec := range output.Recommendations {
		if rec.Title == "The Matrix" {
			foundMatrix = true
			// Should have a match score
			if rec.MatchScore == "" {
				t.Error("Expected match score for The Matrix")
			}
			// Should have a recommendation reason
			if rec.RecommendationReason == "" {
				t.Error("Expected recommendation reason for The Matrix")
			}
		}
	}

	if !foundMatrix {
		t.Error("Expected The Matrix in recommendations")
	}
}

func TestMovieRecommendationEngine_MinRatingFilter(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "High Rated Movie",
					Director: "Director",
					Year:     2000,
					Rating:   9.0,
					Genres:   []string{"Drama"},
				},
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.MovieRecommendationEngine(context.Background(), nil, MovieRecommendationInput{
		Preferences: UserPreferences{
			Genres:    []string{"Drama"},
			MinRating: 8.5,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}

	// Verify high rated movie meets criteria
	if output.Recommendations[0].Rating < 8.5 {
		t.Errorf("Expected recommendation with rating >= 8.5, got: %f", output.Recommendations[0].Rating)
	}
}

func TestMovieRecommendationEngine_LimitResults(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Return 10 movies
			movies := make([]*movieApp.MovieDTO, 10)
			for i := 0; i < 10; i++ {
				movies[i] = &movieApp.MovieDTO{
					ID:       i + 1,
					Title:    fmt.Sprintf("Movie %d", i+1),
					Director: "Director",
					Year:     2000 + i,
					Rating:   8.0 + float64(i)*0.1,
					Genres:   []string{"Drama"},
				}
			}
			return movies, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.MovieRecommendationEngine(context.Background(), nil, MovieRecommendationInput{
		Preferences: UserPreferences{
			Genres: []string{"Drama"},
		},
		Limit: 5,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Recommendations) != 5 {
		t.Errorf("Expected 5 recommendations (limit), got: %d", len(output.Recommendations))
	}
}

func TestMovieRecommendationEngine_NoResults(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.MovieRecommendationEngine(context.Background(), nil, MovieRecommendationInput{
		Preferences: UserPreferences{
			Genres:    []string{"NonExistentGenre"},
			MinRating: 9.5,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(output.Recommendations) != 0 {
		t.Errorf("Expected 0 recommendations, got: %d", len(output.Recommendations))
	}

	if output.TotalFound != 0 {
		t.Errorf("Expected total found 0, got: %d", output.TotalFound)
	}
}

// ===== DirectorCareerAnalysis Tests =====

func TestDirectorCareerAnalysis_Success(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			// Return movies spanning director's career
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "Early Work",
					Director: "Christopher Nolan",
					Year:     1998,
					Rating:   7.0,
					Genres:   []string{"Thriller"},
				},
				{
					ID:       2,
					Title:    "Memento",
					Director: "Christopher Nolan",
					Year:     2000,
					Rating:   8.4,
					Genres:   []string{"Thriller", "Mystery"},
				},
				{
					ID:       3,
					Title:    "The Dark Knight",
					Director: "Christopher Nolan",
					Year:     2008,
					Rating:   9.0,
					Genres:   []string{"Action", "Crime"},
				},
				{
					ID:       4,
					Title:    "Inception",
					Director: "Christopher Nolan",
					Year:     2010,
					Rating:   8.8,
					Genres:   []string{"Sci-Fi", "Thriller"},
				},
				{
					ID:       5,
					Title:    "Interstellar",
					Director: "Christopher Nolan",
					Year:     2014,
					Rating:   8.6,
					Genres:   []string{"Sci-Fi", "Drama"},
				},
				{
					ID:       6,
					Title:    "Dunkirk",
					Director: "Christopher Nolan",
					Year:     2017,
					Rating:   7.8,
					Genres:   []string{"War", "Drama"},
				},
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.DirectorCareerAnalysis(context.Background(), nil, DirectorCareerAnalysisInput{
		Director: "Christopher Nolan",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.Director != "Christopher Nolan" {
		t.Errorf("Expected director 'Christopher Nolan', got: %s", output.Director)
	}

	if output.CareerOverview.TotalMovies != 6 {
		t.Errorf("Expected 6 total movies, got: %d", output.CareerOverview.TotalMovies)
	}

	if !strings.HasPrefix(output.CareerOverview.CareerSpan, "1998-2017") {
		t.Errorf("Expected career span to start with '1998-2017', got: %s", output.CareerOverview.CareerSpan)
	}

	// Should have career phases
	if output.CareerPhases.Early.MovieCount == 0 {
		t.Error("Expected early phase to have movies")
	}

	if output.CareerPhases.Mid.MovieCount == 0 {
		t.Error("Expected mid phase to have movies")
	}

	if output.CareerPhases.Late.MovieCount == 0 {
		t.Error("Expected late phase to have movies")
	}

	// Should have trajectory
	if output.CareerTrajectory == "" {
		t.Error("Expected career trajectory")
	}
}

func TestDirectorCareerAnalysis_NoMovies(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, _, err := tools.DirectorCareerAnalysis(context.Background(), nil, DirectorCareerAnalysisInput{
		Director: "NonExistent Director",
	})

	if err == nil {
		t.Fatal("Expected error for director with no movies, got nil")
	}

	if !strings.Contains(err.Error(), "no movies found") {
		t.Errorf("Expected 'no movies found' error, got: %v", err)
	}
}

func TestDirectorCareerAnalysis_SingleMovie(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return []*movieApp.MovieDTO{
				{
					ID:       1,
					Title:    "First and Only",
					Director: "New Director",
					Year:     2024,
					Rating:   7.5,
					Genres:   []string{"Drama"},
				},
			}, nil
		},
	}

	tools := NewCompoundTools(mockService)
	_, output, err := tools.DirectorCareerAnalysis(context.Background(), nil, DirectorCareerAnalysisInput{
		Director: "New Director",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.CareerOverview.TotalMovies != 1 {
		t.Errorf("Expected 1 total movie, got: %d", output.CareerOverview.TotalMovies)
	}

	// Career span should be single year
	if !strings.HasPrefix(output.CareerOverview.CareerSpan, "2024-2024") {
		t.Errorf("Expected career span to start with '2024-2024', got: %s", output.CareerOverview.CareerSpan)
	}
}

func TestDirectorCareerAnalysis_ServiceError(t *testing.T) {
	mockService := &MockMovieService{
		SearchMoviesFunc: func(ctx context.Context, query movieApp.SearchMoviesQuery) ([]*movieApp.MovieDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}

	tools := NewCompoundTools(mockService)
	_, _, err := tools.DirectorCareerAnalysis(context.Background(), nil, DirectorCareerAnalysisInput{
		Director: "Some Director",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to search movies") {
		t.Errorf("Expected 'failed to search movies' error, got: %v", err)
	}
}

// ===== Helper Function Tests =====

func TestCalculateYearScore(t *testing.T) {
	tests := []struct {
		name       string
		movieYear  float64
		yearFrom   float64
		yearTo     float64
		wantScore  float64
		expectDesc string
	}{
		{
			name:       "movie within range",
			movieYear:  2000,
			yearFrom:   1990,
			yearTo:     2010,
			wantScore:  1.0,
			expectDesc: "should return 1.0 for movie within year range",
		},
		{
			name:       "movie at lower bound",
			movieYear:  1990,
			yearFrom:   1990,
			yearTo:     2010,
			wantScore:  1.0,
			expectDesc: "should return 1.0 for movie at lower bound",
		},
		{
			name:       "movie at upper bound",
			movieYear:  2010,
			yearFrom:   1990,
			yearTo:     2010,
			wantScore:  1.0,
			expectDesc: "should return 1.0 for movie at upper bound",
		},
		{
			name:       "movie before range - 10 years",
			movieYear:  1980,
			yearFrom:   1990,
			yearTo:     2010,
			wantScore:  0.8,
			expectDesc: "should return 0.8 for movie 10 years before range (1.0 - 10/50)",
		},
		{
			name:       "movie after range - 25 years",
			movieYear:  2035,
			yearFrom:   1990,
			yearTo:     2010,
			wantScore:  0.5,
			expectDesc: "should return 0.5 for movie 25 years after range (1.0 - 25/50)",
		},
		{
			name:       "movie with default yearFrom",
			movieYear:  1950,
			yearFrom:   0,
			yearTo:     2000,
			wantScore:  1.0,
			expectDesc: "should use default yearFrom=1900 when yearFrom is 0",
		},
		{
			name:       "movie with default yearTo",
			movieYear:  2050,
			yearFrom:   2000,
			yearTo:     0,
			wantScore:  1.0,
			expectDesc: "should use default yearTo=2100 when yearTo is 0",
		},
		{
			name:       "movie with both defaults",
			movieYear:  2000,
			yearFrom:   0,
			yearTo:     0,
			wantScore:  1.0,
			expectDesc: "should use default range 1900-2100 when both are 0",
		},
		{
			name:       "very old movie before default range",
			movieYear:  1850,
			yearFrom:   0,
			yearTo:     0,
			wantScore:  0.0,
			expectDesc: "should return 0.0 for movie 50 years before default yearFrom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateYearScore(tt.movieYear, tt.yearFrom, tt.yearTo)
			if score != tt.wantScore {
				t.Errorf("%s: got score %.2f, want %.2f", tt.expectDesc, score, tt.wantScore)
			}
		})
	}
}

func TestDetermineTrajectory(t *testing.T) {
	tests := []struct {
		name        string
		early       float64
		mid         float64
		late        float64
		wantPattern string
	}{
		{
			name:        "ascending career",
			early:       7.0,
			mid:         8.0,
			late:        9.0,
			wantPattern: "Ascending - Consistent improvement over career",
		},
		{
			name:        "descending career",
			early:       9.0,
			mid:         8.0,
			late:        7.0,
			wantPattern: "Descending - Ratings declined over time",
		},
		{
			name:        "peak in mid-career",
			early:       7.5,
			mid:         9.0,
			late:        8.0,
			wantPattern: "Peak in mid-career",
		},
		{
			name:        "late career resurgence - ascending",
			early:       7.0,
			mid:         7.5,
			late:        8.5,
			wantPattern: "Ascending - Consistent improvement over career",
		},
		{
			name:        "peak in mid-career with slight decline",
			early:       8.0,
			mid:         8.1,
			late:        7.9,
			wantPattern: "Peak in mid-career",
		},
		{
			name:        "consistent quality throughout career",
			early:       9.0,
			mid:         8.0,
			late:        8.0,
			wantPattern: "Consistent quality throughout career",
		},
		{
			name:        "late career resurgence after mid decline",
			early:       7.0,
			mid:         6.5,
			late:        8.0,
			wantPattern: "Late career resurgence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trajectory := determineTrajectory(tt.early, tt.mid, tt.late)
			if !strings.Contains(trajectory, tt.wantPattern) && trajectory != tt.wantPattern {
				t.Errorf("Expected trajectory to contain '%s', got: %s", tt.wantPattern, trajectory)
			}
		})
	}
}

func TestFindWorstMovie(t *testing.T) {
	tests := []struct {
		name      string
		movies    []*movieApp.MovieDTO
		wantTitle string
		wantNil   bool
	}{
		{
			name:    "empty list",
			movies:  []*movieApp.MovieDTO{},
			wantNil: true,
		},
		{
			name: "single movie",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Only Movie", Rating: 7.5},
			},
			wantTitle: "Only Movie",
		},
		{
			name: "multiple movies with different ratings",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Great Movie", Rating: 9.0},
				{ID: 2, Title: "Bad Movie", Rating: 5.0},
				{ID: 3, Title: "Average Movie", Rating: 7.0},
			},
			wantTitle: "Bad Movie",
		},
		{
			name: "movies with zero ratings ignored in worst selection",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "Rated Movie", Rating: 6.0},
				{ID: 2, Title: "Unrated Movie", Rating: 0},
			},
			wantTitle: "Rated Movie",
		},
		{
			name: "all movies unrated - returns first",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "First Unrated", Rating: 0},
				{ID: 2, Title: "Second Unrated", Rating: 0},
			},
			wantTitle: "First Unrated",
		},
		{
			name: "movies with same low rating - returns first",
			movies: []*movieApp.MovieDTO{
				{ID: 1, Title: "First Bad", Rating: 5.0},
				{ID: 2, Title: "Second Bad", Rating: 5.0},
				{ID: 3, Title: "Good Movie", Rating: 8.0},
			},
			wantTitle: "First Bad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worst := findWorstMovie(tt.movies)
			if tt.wantNil {
				if worst != nil {
					t.Errorf("Expected nil, got: %v", worst)
				}
			} else {
				if worst == nil {
					t.Fatal("Expected movie, got nil")
				}
				if worst.Title != tt.wantTitle {
					t.Errorf("Expected worst movie '%s', got: '%s'", tt.wantTitle, worst.Title)
				}
			}
		})
	}
}
