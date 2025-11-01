package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// CompoundTools provides SDK-based MCP handlers for compound operations
type CompoundTools struct {
	movieService MovieService
}

// NewCompoundTools creates a new compound tools instance
func NewCompoundTools(movieService MovieService) *CompoundTools {
	return &CompoundTools{
		movieService: movieService,
	}
}

// ===== bulk_movie_import Tool =====

// BulkMovieImportInput defines the input schema for bulk_movie_import tool
type BulkMovieImportInput struct {
	Movies []MovieImportItem `json:"movies" jsonschema:"required,description=Array of movies to import"`
}

// MovieImportItem defines a single movie for bulk import
type MovieImportItem struct {
	Title     string   `json:"title" jsonschema:"required,description=Movie title"`
	Director  string   `json:"director" jsonschema:"required,description=Movie director"`
	Year      int      `json:"year" jsonschema:"required,description=Release year"`
	Rating    float64  `json:"rating" jsonschema:"required,description=Movie rating (0-10)"`
	Genres    []string `json:"genres,omitempty" jsonschema:"description=List of genres"`
	PosterURL string   `json:"poster_url,omitempty" jsonschema:"description=URL to movie poster"`
}

// BulkMovieImportOutput defines the output schema for bulk_movie_import tool
type BulkMovieImportOutput struct {
	Imported    int            `json:"imported" jsonschema:"description=Number of successfully imported movies"`
	Failed      int            `json:"failed" jsonschema:"description=Number of failed imports"`
	Total       int            `json:"total" jsonschema:"description=Total movies attempted"`
	SuccessRate string         `json:"success_rate" jsonschema:"description=Success rate percentage"`
	Results     []ImportResult `json:"results" jsonschema:"description=Successful import results"`
	Errors      []ImportError  `json:"errors" jsonschema:"description=Failed import errors"`
}

// ImportResult represents a successful import
type ImportResult struct {
	Index int    `json:"index" jsonschema:"description=Index in original array"`
	ID    int    `json:"id" jsonschema:"description=Created movie ID"`
	Title string `json:"title" jsonschema:"description=Movie title"`
}

// ImportError represents a failed import
type ImportError struct {
	Index int    `json:"index" jsonschema:"description=Index in original array"`
	Title string `json:"title,omitempty" jsonschema:"description=Movie title if available"`
	Error string `json:"error" jsonschema:"description=Error message"`
}

// BulkMovieImport handles the bulk_movie_import tool call
func (t *CompoundTools) BulkMovieImport(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BulkMovieImportInput,
) (*mcp.CallToolResult, BulkMovieImportOutput, error) {
	var results []ImportResult
	var errors []ImportError

	for i, movie := range input.Movies {
		// Create movie command
		cmd := movieApp.CreateMovieCommand{
			Title:     movie.Title,
			Director:  movie.Director,
			Year:      movie.Year,
			Rating:    movie.Rating,
			Genres:    movie.Genres,
			PosterURL: movie.PosterURL,
		}

		// Create movie
		movieDTO, err := t.movieService.CreateMovie(ctx, cmd)
		if err != nil {
			errors = append(errors, ImportError{
				Index: i,
				Title: movie.Title,
				Error: err.Error(),
			})
		} else {
			results = append(results, ImportResult{
				Index: i,
				ID:    movieDTO.ID,
				Title: movieDTO.Title,
			})
		}
	}

	successRate := 0.0
	if len(input.Movies) > 0 {
		successRate = float64(len(results)) / float64(len(input.Movies)) * 100
	}

	output := BulkMovieImportOutput{
		Imported:    len(results),
		Failed:      len(errors),
		Total:       len(input.Movies),
		SuccessRate: fmt.Sprintf("%.1f%%", successRate),
		Results:     results,
		Errors:      errors,
	}

	return nil, output, nil
}

// ===== movie_recommendation_engine Tool =====

// MovieRecommendationInput defines the input schema for movie_recommendation_engine tool
type MovieRecommendationInput struct {
	Preferences UserPreferences `json:"preferences,omitempty" jsonschema:"description=User preferences for recommendations"`
	Limit       int             `json:"limit,omitempty" jsonschema:"description=Maximum number of recommendations,default=10"`
}

// UserPreferences defines user preferences for recommendations
type UserPreferences struct {
	Genres        []string `json:"genres,omitempty" jsonschema:"description=Preferred genres"`
	MinRating     float64  `json:"min_rating,omitempty" jsonschema:"description=Minimum rating"`
	YearFrom      int      `json:"year_from,omitempty" jsonschema:"description=Start of year range"`
	YearTo        int      `json:"year_to,omitempty" jsonschema:"description=End of year range"`
	ExcludeMovies []string `json:"exclude_movies,omitempty" jsonschema:"description=Movie titles to exclude"`
}

// MovieRecommendationOutput defines the output schema for movie_recommendation_engine tool
type MovieRecommendationOutput struct {
	Recommendations []Recommendation  `json:"recommendations" jsonschema:"description=List of recommended movies"`
	TotalFound      int               `json:"total_found" jsonschema:"description=Total recommendations found"`
	PreferencesUsed PreferenceSummary `json:"preferences_used" jsonschema:"description=Summary of preferences used"`
}

// Recommendation represents a single movie recommendation
type Recommendation struct {
	Rank                 int      `json:"rank" jsonschema:"description=Recommendation rank"`
	MovieID              int      `json:"movie_id" jsonschema:"description=Movie ID"`
	Title                string   `json:"title" jsonschema:"description=Movie title"`
	Director             string   `json:"director" jsonschema:"description=Director name"`
	Year                 int      `json:"year" jsonschema:"description=Release year"`
	Rating               float64  `json:"rating" jsonschema:"description=Movie rating"`
	Genres               []string `json:"genres" jsonschema:"description=List of genres"`
	MatchScore           string   `json:"match_score" jsonschema:"description=Match percentage"`
	RecommendationReason string   `json:"recommendation_reason" jsonschema:"description=Why this was recommended"`
}

// PreferenceSummary summarizes the preferences used
type PreferenceSummary struct {
	Genres        []string `json:"genres" jsonschema:"description=Genres used"`
	MinRating     float64  `json:"min_rating" jsonschema:"description=Minimum rating used"`
	YearRange     string   `json:"year_range" jsonschema:"description=Year range used"`
	ExcludedCount int      `json:"excluded_count" jsonschema:"description=Number of excluded movies"`
}

// MovieRecommendationEngine handles the movie_recommendation_engine tool call
func (t *CompoundTools) MovieRecommendationEngine(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MovieRecommendationInput,
) (*mcp.CallToolResult, MovieRecommendationOutput, error) {
	// Set default limit
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	// Build query to get candidate movies
	query := movieApp.SearchMoviesQuery{
		Limit: limit * 3, // Get more to filter
	}

	// Get all movies
	movies, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, MovieRecommendationOutput{}, fmt.Errorf("failed to search movies: %w", err)
	}

	// Score and filter movies
	type scoredMovie struct {
		movie *movieApp.MovieDTO
		score float64
	}

	var scoredMovies []scoredMovie
	excludeMap := make(map[string]bool)
	for _, title := range input.Preferences.ExcludeMovies {
		excludeMap[strings.ToLower(title)] = true
	}

	for _, movie := range movies {
		// Skip excluded movies
		if excludeMap[strings.ToLower(movie.Title)] {
			continue
		}

		// Calculate recommendation score
		score := calculateRecommendationScore(movie, input.Preferences)

		if score > 0.3 { // Minimum threshold
			scoredMovies = append(scoredMovies, scoredMovie{
				movie: movie,
				score: score,
			})
		}
	}

	// Sort by score
	sort.Slice(scoredMovies, func(i, j int) bool {
		return scoredMovies[i].score > scoredMovies[j].score
	})

	// Prepare recommendations
	recommendations := []Recommendation{}
	for i, sm := range scoredMovies {
		if i >= limit {
			break
		}

		recommendations = append(recommendations, Recommendation{
			Rank:                 i + 1,
			MovieID:              sm.movie.ID,
			Title:                sm.movie.Title,
			Director:             sm.movie.Director,
			Year:                 sm.movie.Year,
			Rating:               sm.movie.Rating,
			Genres:               sm.movie.Genres,
			MatchScore:           fmt.Sprintf("%.1f%%", sm.score*100),
			RecommendationReason: generateRecommendationReason(sm.movie, input.Preferences, sm.score),
		})
	}

	output := MovieRecommendationOutput{
		Recommendations: recommendations,
		TotalFound:      len(recommendations),
		PreferencesUsed: PreferenceSummary{
			Genres:        input.Preferences.Genres,
			MinRating:     input.Preferences.MinRating,
			YearRange:     fmt.Sprintf("%d-%d", input.Preferences.YearFrom, input.Preferences.YearTo),
			ExcludedCount: len(input.Preferences.ExcludeMovies),
		},
	}

	return nil, output, nil
}

// ===== director_career_analysis Tool =====

// DirectorCareerAnalysisInput defines the input schema for director_career_analysis tool
type DirectorCareerAnalysisInput struct {
	Director string `json:"director" jsonschema:"required,description=Director name to analyze"`
}

// DirectorCareerAnalysisOutput defines the output schema for director_career_analysis tool
type DirectorCareerAnalysisOutput struct {
	Director            string             `json:"director" jsonschema:"description=Director name"`
	CareerOverview      CareerOverview     `json:"career_overview" jsonschema:"description=Overall career statistics"`
	CareerPhases        CareerPhases       `json:"career_phases" jsonschema:"description=Career broken into phases"`
	CareerTrajectory    string             `json:"career_trajectory" jsonschema:"description=Description of career trajectory"`
	GenreSpecialization []GenreFrequency   `json:"genre_specialization" jsonschema:"description=Top genres by count"`
	NotableWorks        NotableWorks       `json:"notable_works" jsonschema:"description=Highest and lowest rated works"`
	Filmography         []FilmographyEntry `json:"filmography" jsonschema:"description=Complete filmography"`
}

// CareerOverview provides overall career stats
type CareerOverview struct {
	TotalMovies   int    `json:"total_movies" jsonschema:"description=Total number of movies"`
	CareerSpan    string `json:"career_span" jsonschema:"description=Career span in years"`
	AverageRating string `json:"average_rating" jsonschema:"description=Average rating across all movies"`
}

// CareerPhases breaks career into early/mid/late phases
type CareerPhases struct {
	Early PhaseInfo `json:"early" jsonschema:"description=Early career phase"`
	Mid   PhaseInfo `json:"mid" jsonschema:"description=Mid career phase"`
	Late  PhaseInfo `json:"late" jsonschema:"description=Late career phase"`
}

// PhaseInfo provides info about a career phase
type PhaseInfo struct {
	Period        string `json:"period" jsonschema:"description=Year range of this phase"`
	MovieCount    int    `json:"movie_count" jsonschema:"description=Number of movies in this phase"`
	AverageRating string `json:"average_rating" jsonschema:"description=Average rating in this phase"`
}

// GenreFrequency represents a genre and its count
type GenreFrequency struct {
	Genre string `json:"genre" jsonschema:"description=Genre name"`
	Count int    `json:"count" jsonschema:"description=Number of movies in this genre"`
}

// NotableWorks highlights best and worst movies
type NotableWorks struct {
	HighestRated MovieSummary `json:"highest_rated" jsonschema:"description=Highest rated movie"`
	LowestRated  MovieSummary `json:"lowest_rated" jsonschema:"description=Lowest rated movie"`
}

// MovieSummary provides brief movie info
type MovieSummary struct {
	Title  string  `json:"title" jsonschema:"description=Movie title"`
	Year   int     `json:"year" jsonschema:"description=Release year"`
	Rating float64 `json:"rating" jsonschema:"description=Movie rating"`
}

// FilmographyEntry represents one movie in filmography
type FilmographyEntry struct {
	Year   int      `json:"year" jsonschema:"description=Release year"`
	Title  string   `json:"title" jsonschema:"description=Movie title"`
	Rating float64  `json:"rating" jsonschema:"description=Movie rating"`
	Genres []string `json:"genres" jsonschema:"description=Movie genres"`
}

// DirectorCareerAnalysis handles the director_career_analysis tool call
func (t *CompoundTools) DirectorCareerAnalysis(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DirectorCareerAnalysisInput,
) (*mcp.CallToolResult, DirectorCareerAnalysisOutput, error) {
	// Search for all movies by this director
	query := movieApp.SearchMoviesQuery{
		Director: input.Director,
		Limit:    100,
	}

	movies, err := t.movieService.SearchMovies(ctx, query)
	if err != nil {
		return nil, DirectorCareerAnalysisOutput{}, fmt.Errorf("failed to search movies: %w", err)
	}

	if len(movies) == 0 {
		return nil, DirectorCareerAnalysisOutput{}, fmt.Errorf("no movies found for director: %s", input.Director)
	}

	// Sort movies by year
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Year < movies[j].Year
	})

	// Analyze career metrics
	totalMovies := len(movies)
	var totalRating float64
	var ratingsCount int
	genreCount := make(map[string]int)

	// Career phases
	var earlyCareer, midCareer, lateCareer []*movieApp.MovieDTO

	firstYear := movies[0].Year
	lastYear := movies[len(movies)-1].Year
	careerSpan := lastYear - firstYear

	for _, movie := range movies {
		// Rating analysis
		if movie.Rating > 0 {
			totalRating += movie.Rating
			ratingsCount++
		}

		// Genre analysis
		for _, genre := range movie.Genres {
			genreCount[genre]++
		}

		// Career phase categorization
		yearInCareer := movie.Year - firstYear
		if float64(yearInCareer) < float64(careerSpan)*0.33 {
			earlyCareer = append(earlyCareer, movie)
		} else if float64(yearInCareer) < float64(careerSpan)*0.66 {
			midCareer = append(midCareer, movie)
		} else {
			lateCareer = append(lateCareer, movie)
		}
	}

	// Calculate phase averages
	earlyAvg := calculateAverageRating(earlyCareer)
	midAvg := calculateAverageRating(midCareer)
	lateAvg := calculateAverageRating(lateCareer)

	// Find best and worst movies
	bestMovie := findBestMovie(movies)
	worstMovie := findWorstMovie(movies)

	// Genre evolution
	primaryGenres := findTopGenres(genreCount, 3)

	// Build filmography
	filmography := []FilmographyEntry{}
	for _, m := range movies {
		filmography = append(filmography, FilmographyEntry{
			Year:   m.Year,
			Title:  m.Title,
			Rating: m.Rating,
			Genres: m.Genres,
		})
	}

	output := DirectorCareerAnalysisOutput{
		Director: input.Director,
		CareerOverview: CareerOverview{
			TotalMovies:   totalMovies,
			CareerSpan:    fmt.Sprintf("%d-%d (%d years)", firstYear, lastYear, careerSpan),
			AverageRating: fmt.Sprintf("%.1f", totalRating/float64(ratingsCount)),
		},
		CareerPhases: CareerPhases{
			Early: PhaseInfo{
				Period:        fmt.Sprintf("%d-%d", firstYear, firstYear+careerSpan/3),
				MovieCount:    len(earlyCareer),
				AverageRating: fmt.Sprintf("%.1f", earlyAvg),
			},
			Mid: PhaseInfo{
				Period:        fmt.Sprintf("%d-%d", firstYear+careerSpan/3, firstYear+2*careerSpan/3),
				MovieCount:    len(midCareer),
				AverageRating: fmt.Sprintf("%.1f", midAvg),
			},
			Late: PhaseInfo{
				Period:        fmt.Sprintf("%d-%d", firstYear+2*careerSpan/3, lastYear),
				MovieCount:    len(lateCareer),
				AverageRating: fmt.Sprintf("%.1f", lateAvg),
			},
		},
		CareerTrajectory:    determineTrajectory(earlyAvg, midAvg, lateAvg),
		GenreSpecialization: primaryGenres,
		NotableWorks: NotableWorks{
			HighestRated: MovieSummary{
				Title:  bestMovie.Title,
				Year:   bestMovie.Year,
				Rating: bestMovie.Rating,
			},
			LowestRated: MovieSummary{
				Title:  worstMovie.Title,
				Year:   worstMovie.Year,
				Rating: worstMovie.Rating,
			},
		},
		Filmography: filmography,
	}

	return nil, output, nil
}

// Helper functions

func calculateRecommendationScore(movie *movieApp.MovieDTO, prefs UserPreferences) float64 {
	score := 0.0

	// Genre matching (40% weight)
	if len(prefs.Genres) > 0 {
		genreScore := calculateGenreScore(movie.Genres, prefs.Genres)
		score += genreScore * 0.4
	} else {
		score += 0.4 // No genre preference
	}

	// Rating score (30% weight)
	if movie.Rating >= prefs.MinRating {
		ratingScore := movie.Rating / 10.0
		score += ratingScore * 0.3
	}

	// Year relevance (20% weight)
	if prefs.YearFrom > 0 || prefs.YearTo > 0 {
		yearScore := calculateYearScore(float64(movie.Year), float64(prefs.YearFrom), float64(prefs.YearTo))
		score += yearScore * 0.2
	} else {
		score += 0.2 // No year preference
	}

	// Popularity boost (10% weight) - using rating as proxy
	if movie.Rating >= 8.0 {
		score += 0.1
	}

	return score
}

func calculateGenreScore(movieGenres, preferredGenres []string) float64 {
	if len(movieGenres) == 0 || len(preferredGenres) == 0 {
		return 0
	}

	matches := 0
	for _, mg := range movieGenres {
		for _, pg := range preferredGenres {
			if strings.EqualFold(mg, pg) {
				matches++
				break
			}
		}
	}

	return float64(matches) / float64(len(preferredGenres))
}

func calculateYearScore(movieYear, yearFrom, yearTo float64) float64 {
	if yearFrom == 0 {
		yearFrom = 1900
	}
	if yearTo == 0 {
		yearTo = 2100
	}

	if movieYear >= yearFrom && movieYear <= yearTo {
		return 1.0
	}

	// Gradual decrease for movies outside range
	if movieYear < yearFrom {
		diff := yearFrom - movieYear
		return 1.0 - (diff / 50.0) // -0.02 per year
	}
	diff := movieYear - yearTo
	return 1.0 - (diff / 50.0)
}

func generateRecommendationReason(movie *movieApp.MovieDTO, prefs UserPreferences, score float64) string {
	reasons := []string{}

	if score > 0.8 {
		reasons = append(reasons, "Excellent match")
	} else if score > 0.6 {
		reasons = append(reasons, "Good match")
	}

	if movie.Rating >= 8.0 {
		reasons = append(reasons, "Highly rated")
	}

	if len(prefs.Genres) > 0 {
		for _, genre := range movie.Genres {
			for _, prefGenre := range prefs.Genres {
				if strings.EqualFold(genre, prefGenre) {
					reasons = append(reasons, fmt.Sprintf("Matches your interest in %s", genre))
					break
				}
			}
		}
	}

	return strings.Join(reasons, "; ")
}

func calculateAverageRating(movies []*movieApp.MovieDTO) float64 {
	if len(movies) == 0 {
		return 0
	}

	var total float64
	var count int
	for _, m := range movies {
		if m.Rating > 0 {
			total += m.Rating
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func findBestMovie(movies []*movieApp.MovieDTO) *movieApp.MovieDTO {
	if len(movies) == 0 {
		return nil
	}

	best := movies[0]
	for _, m := range movies {
		if m.Rating > best.Rating {
			best = m
		}
	}
	return best
}

func findWorstMovie(movies []*movieApp.MovieDTO) *movieApp.MovieDTO {
	if len(movies) == 0 {
		return nil
	}

	worst := movies[0]
	for _, m := range movies {
		if m.Rating > 0 && (worst.Rating == 0 || m.Rating < worst.Rating) {
			worst = m
		}
	}
	return worst
}

func findTopGenres(genreCount map[string]int, limit int) []GenreFrequency {
	type genreFreq struct {
		genre string
		count int
	}

	frequencies := []genreFreq{}
	for g, c := range genreCount {
		frequencies = append(frequencies, genreFreq{g, c})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].count > frequencies[j].count
	})

	result := []GenreFrequency{}
	for i := 0; i < limit && i < len(frequencies); i++ {
		result = append(result, GenreFrequency{
			Genre: frequencies[i].genre,
			Count: frequencies[i].count,
		})
	}

	return result
}

func determineTrajectory(early, mid, late float64) string {
	if late > mid && mid > early {
		return "Ascending - Consistent improvement over career"
	} else if late < mid && mid < early {
		return "Descending - Ratings declined over time"
	} else if mid > early && mid > late {
		return "Peak in mid-career"
	} else if late > early {
		return "Late career resurgence"
	}
	return "Consistent quality throughout career"
}
