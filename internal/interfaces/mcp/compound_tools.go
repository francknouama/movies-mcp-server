package mcp

import (
	"context"
	"fmt"
	"sort"
	"strings"

	movieApp "github.com/francknouama/movies-mcp-server/internal/application/movie"
)

// CompoundToolHandlers handles compound tool operations
type CompoundToolHandlers struct {
	movieService *movieApp.Service
}

// NewCompoundToolHandlers creates a new CompoundToolHandlers instance
func NewCompoundToolHandlers(movieService *movieApp.Service) *CompoundToolHandlers {
	return &CompoundToolHandlers{
		movieService: movieService,
	}
}

// HandleBulkMovieImport handles importing multiple movies at once
func (h *CompoundToolHandlers) HandleBulkMovieImport(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	// Parse movies array
	moviesData, ok := arguments["movies"].([]interface{})
	if !ok {
		sendError(id, -32602, "Invalid movies array", nil)
		return
	}

	ctx := context.Background()
	var results []map[string]interface{}
	var errors []map[string]interface{}

	for i, movieData := range moviesData {
		movie, ok := movieData.(map[string]interface{})
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid movie data format",
			})
			continue
		}

		// Extract movie fields
		title, ok := movie["title"].(string)
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid title field",
			})
			continue
		}
		director, ok := movie["director"].(string)
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid director field",
			})
			continue
		}
		year, ok := movie["year"].(float64)
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid year field",
			})
			continue
		}
		rating, ok := movie["rating"].(float64)
		if !ok {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": "Invalid rating field",
			})
			continue
		}
		genres, ok := movie["genres"].([]interface{})
		if !ok {
			genres = []interface{}{} // Default to empty if not provided
		}
		posterURL, ok := movie["poster_url"].(string)
		if !ok {
			posterURL = "" // Default to empty if not provided
		}

		// Convert genres
		genreStrings := make([]string, 0, len(genres))
		for _, g := range genres {
			if genreStr, ok := g.(string); ok {
				genreStrings = append(genreStrings, genreStr)
			}
		}

		// Create movie
		cmd := movieApp.CreateMovieCommand{
			Title:     title,
			Director:  director,
			Year:      int(year),
			Rating:    rating,
			Genres:    genreStrings,
			PosterURL: posterURL,
		}

		movieDTO, err := h.movieService.CreateMovie(ctx, cmd)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"title": title,
				"error": err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"index": i,
				"id":    movieDTO.ID,
				"title": movieDTO.Title,
			})
		}
	}

	response := map[string]interface{}{
		"imported":     len(results),
		"failed":       len(errors),
		"total":        len(moviesData),
		"success_rate": fmt.Sprintf("%.1f%%", float64(len(results))/float64(len(moviesData))*100),
		"results":      results,
		"errors":       errors,
	}

	sendResult(id, response)
}

// HandleMovieRecommendationEngine provides intelligent movie recommendations
func (h *CompoundToolHandlers) HandleMovieRecommendationEngine(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	ctx := context.Background()

	// Parse parameters
	userPreferences, ok := arguments["preferences"].(map[string]interface{})
	if !ok {
		userPreferences = make(map[string]interface{}) // Default to empty if not provided
	}
	limit := 10
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	// Extract preferences
	genres := extractStringArray(userPreferences["genres"])
	minRating, ok := userPreferences["min_rating"].(float64)
	if !ok {
		minRating = 0.0 // Default to 0.0 if not provided
	}
	yearFrom, ok := userPreferences["year_from"].(float64)
	if !ok {
		yearFrom = 0.0 // Default to 0.0 if not provided
	}
	yearTo, ok := userPreferences["year_to"].(float64)
	if !ok {
		yearTo = 0.0 // Default to 0.0 if not provided
	}
	excludeMovies := extractStringArray(userPreferences["exclude_movies"])

	// Build recommendation query
	query := movieApp.SearchMoviesQuery{
		Limit: limit * 3, // Get more to filter
	}

	// Get all movies (in production, would use more sophisticated queries)
	movies, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, -32603, "Failed to search movies", err.Error())
		return
	}

	// Score and filter movies
	type scoredMovie struct {
		movie *movieApp.MovieDTO
		score float64
	}

	var scoredMovies []scoredMovie
	excludeMap := make(map[string]bool)
	for _, title := range excludeMovies {
		excludeMap[strings.ToLower(title)] = true
	}

	for _, movie := range movies {
		// Skip excluded movies
		if excludeMap[strings.ToLower(movie.Title)] {
			continue
		}

		// Calculate recommendation score
		score := 0.0

		// Genre matching (40% weight)
		if len(genres) > 0 {
			genreScore := calculateGenreScore(movie.Genres, genres)
			score += genreScore * 0.4
		} else {
			score += 0.4 // No genre preference
		}

		// Rating score (30% weight)
		if movie.Rating >= minRating {
			ratingScore := movie.Rating / 10.0
			score += ratingScore * 0.3
		}

		// Year relevance (20% weight)
		if yearFrom > 0 || yearTo > 0 {
			yearScore := calculateYearScore(float64(movie.Year), yearFrom, yearTo)
			score += yearScore * 0.2
		} else {
			score += 0.2 // No year preference
		}

		// Popularity boost (10% weight) - using rating as proxy
		if movie.Rating >= 8.0 {
			score += 0.1
		}

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
	recommendations := []map[string]interface{}{}
	for i, sm := range scoredMovies {
		if i >= limit {
			break
		}

		recommendations = append(recommendations, map[string]interface{}{
			"rank":                  i + 1,
			"movie_id":              sm.movie.ID,
			"title":                 sm.movie.Title,
			"director":              sm.movie.Director,
			"year":                  sm.movie.Year,
			"rating":                sm.movie.Rating,
			"genres":                sm.movie.Genres,
			"match_score":           fmt.Sprintf("%.1f%%", sm.score*100),
			"recommendation_reason": generateRecommendationReason(sm.movie, userPreferences, sm.score),
		})
	}

	response := map[string]interface{}{
		"recommendations": recommendations,
		"total_found":     len(recommendations),
		"preferences_used": map[string]interface{}{
			"genres":         genres,
			"min_rating":     minRating,
			"year_range":     fmt.Sprintf("%d-%d", int(yearFrom), int(yearTo)),
			"excluded_count": len(excludeMovies),
		},
	}

	sendResult(id, response)
}

// HandleDirectorCareerAnalysis analyzes a director's career trajectory
func (h *CompoundToolHandlers) HandleDirectorCareerAnalysis(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	ctx := context.Background()

	directorName, ok := arguments["director"].(string)
	if !ok || directorName == "" {
		sendError(id, -32602, "Director name is required", nil)
		return
	}

	// Search for all movies by this director
	query := movieApp.SearchMoviesQuery{
		Director: directorName,
		Limit:    100,
	}

	movies, err := h.movieService.SearchMovies(ctx, query)
	if err != nil {
		sendError(id, -32603, "Failed to search movies", err.Error())
		return
	}

	if len(movies) == 0 {
		sendError(id, -32602, "No movies found for director", directorName)
		return
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
	earlyCareer := []*movieApp.MovieDTO{}
	midCareer := []*movieApp.MovieDTO{}
	lateCareer := []*movieApp.MovieDTO{}

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

	response := map[string]interface{}{
		"director": directorName,
		"career_overview": map[string]interface{}{
			"total_movies":   totalMovies,
			"career_span":    fmt.Sprintf("%d-%d (%d years)", firstYear, lastYear, careerSpan),
			"average_rating": fmt.Sprintf("%.1f", totalRating/float64(ratingsCount)),
		},
		"career_phases": map[string]interface{}{
			"early": map[string]interface{}{
				"period":         fmt.Sprintf("%d-%d", firstYear, firstYear+careerSpan/3),
				"movie_count":    len(earlyCareer),
				"average_rating": fmt.Sprintf("%.1f", earlyAvg),
			},
			"mid": map[string]interface{}{
				"period":         fmt.Sprintf("%d-%d", firstYear+careerSpan/3, firstYear+2*careerSpan/3),
				"movie_count":    len(midCareer),
				"average_rating": fmt.Sprintf("%.1f", midAvg),
			},
			"late": map[string]interface{}{
				"period":         fmt.Sprintf("%d-%d", firstYear+2*careerSpan/3, lastYear),
				"movie_count":    len(lateCareer),
				"average_rating": fmt.Sprintf("%.1f", lateAvg),
			},
		},
		"career_trajectory":    determineTrajectory(earlyAvg, midAvg, lateAvg),
		"genre_specialization": primaryGenres,
		"notable_works": map[string]interface{}{
			"highest_rated": map[string]interface{}{
				"title":  bestMovie.Title,
				"year":   bestMovie.Year,
				"rating": bestMovie.Rating,
			},
			"lowest_rated": map[string]interface{}{
				"title":  worstMovie.Title,
				"year":   worstMovie.Year,
				"rating": worstMovie.Rating,
			},
		},
		"filmography": formatFilmography(movies),
	}

	sendResult(id, response)
}

// Helper functions

func extractStringArray(data interface{}) []string {
	result := []string{}
	if arr, ok := data.([]interface{}); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
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
	} else {
		diff := movieYear - yearTo
		return 1.0 - (diff / 50.0)
	}
}

func generateRecommendationReason(movie *movieApp.MovieDTO, preferences map[string]interface{}, score float64) string {
	reasons := []string{}

	if score > 0.8 {
		reasons = append(reasons, "Excellent match")
	} else if score > 0.6 {
		reasons = append(reasons, "Good match")
	}

	if movie.Rating >= 8.0 {
		reasons = append(reasons, "Highly rated")
	}

	if genres, ok := preferences["genres"].([]interface{}); ok && len(genres) > 0 {
		for _, genre := range movie.Genres {
			for _, prefGenre := range genres {
				if strings.EqualFold(genre, prefGenre.(string)) {
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

func findTopGenres(genreCount map[string]int, limit int) []map[string]interface{} {
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

	result := []map[string]interface{}{}
	for i := 0; i < limit && i < len(frequencies); i++ {
		result = append(result, map[string]interface{}{
			"genre": frequencies[i].genre,
			"count": frequencies[i].count,
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
	} else {
		return "Consistent quality throughout career"
	}
}

func formatFilmography(movies []*movieApp.MovieDTO) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, m := range movies {
		result = append(result, map[string]interface{}{
			"year":   m.Year,
			"title":  m.Title,
			"rating": m.Rating,
			"genres": m.Genres,
		})
	}
	return result
}
