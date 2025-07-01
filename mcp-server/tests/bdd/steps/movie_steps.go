package steps

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// MovieResponse represents the structure of a movie response
type MovieResponse struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Director string   `json:"director"`
	Year     int      `json:"year"`
	Rating   float64  `json:"rating"`
	Genres   []string `json:"genres"`
}

// MoviesResponse represents a list of movies response
type MoviesResponse struct {
	Movies []MovieResponse `json:"movies"`
	Total  int             `json:"total"`
}

// InitializeMovieSteps registers movie-related step definitions
func InitializeMovieSteps(ctx *godog.ScenarioContext) {
	stepContext := NewCommonStepContext()

	// Movie management steps
	ctx.Step(`^I call the "([^"]*)" tool with:$`, stepContext.iCallTheToolWith)
	ctx.Step(`^the response should contain a movie with:$`, stepContext.theResponseShouldContainAMovieWith)
	ctx.Step(`^the movie should have an assigned ID$`, stepContext.theMovieShouldHaveAnAssignedID)
	ctx.Step(`^the response should contain (\d+) movies?$`, stepContext.theResponseShouldContainNMovies)
	ctx.Step(`^the following movies exist:$`, stepContext.theFollowingMoviesExist)
	ctx.Step(`^a movie exists with:$`, stepContext.aMovieExistsWith)
	ctx.Step(`^the movie title should contain "([^"]*)"$`, stepContext.theMovieTitleShouldContain)
	ctx.Step(`^all movies should have director "([^"]*)"$`, stepContext.allMoviesShouldHaveDirector)
	
	// Additional missing movie step definitions
	ctx.Step(`^I call the "([^"]*)" tool with the movie ID$`, stepContext.iCallToolWithMovieID)
	ctx.Step(`^I call the "([^"]*)" tool with movie ID (\d+)$`, stepContext.iCallToolWithSpecificMovieID)
	ctx.Step(`^the movie title should be "([^"]*)"$`, stepContext.theMovieTitleShouldBe)
	ctx.Step(`^a movie exists with title "([^"]*)"$`, stepContext.aMovieExistsWithTitle)
	ctx.Step(`^the movie title should be updated to "([^"]*)"$`, stepContext.theMovieTitleShouldBeUpdatedTo)
	ctx.Step(`^the movie rating should be updated to ([0-9.]+)$`, stepContext.theMovieRatingShouldBeUpdatedTo)
	ctx.Step(`^the movie should no longer exist in the database$`, stepContext.theMovieShouldNoLongerExistInDatabase)
	ctx.Step(`^I call the "([^"]*)" tool with limit (\d+)$`, stepContext.iCallToolWithLimit)
	ctx.Step(`^the movies should be ordered by rating descending$`, stepContext.theMoviesShouldBeOrderedByRatingDescending)
	ctx.Step(`^the first movie should have rating ([0-9.]+)$`, stepContext.theFirstMovieShouldHaveRating)
	ctx.Step(`^the error message should indicate movie not found$`, stepContext.theErrorMessageShouldIndicateMovieNotFound)
	
	// Missing step definitions for scenario outline
	ctx.Step(`^movies exist from various decades$`, stepContext.moviesExistFromVariousDecades)
	ctx.Step(`^I call the "([^"]*)" tool with decade "([^"]*)"$`, stepContext.iCallToolWithDecade)
	ctx.Step(`^all movies should be from years (\d+) to (\d+)$`, stepContext.allMoviesShouldBeFromYearsTo)
	ctx.Step(`^the response should contain the movie details$`, stepContext.theResponseShouldContainTheMovieDetails)
}

// iCallTheToolWith calls an MCP tool with JSON arguments
func (c *CommonStepContext) iCallTheToolWith(toolName string, docString *godog.DocString) error {
	var arguments map[string]interface{}
	
	if err := json.Unmarshal([]byte(docString.Content), &arguments); err != nil {
		return fmt.Errorf("failed to parse JSON arguments: %w", err)
	}

	// Apply ID interpolation
	arguments = c.dataManager.InterpolateMap(arguments)

	_, err := c.bddContext.CallTool(toolName, arguments)
	if err != nil {
		return err
	}

	// Store IDs from successful responses
	if !c.bddContext.HasError() {
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr == nil {
			// Store movie ID if this was an add_movie call
			if toolName == "add_movie" {
				if err := c.dataManager.StoreIDFromResponse(responseData, "id", "movie_id"); err == nil {
					// Also store as last_movie_id for reference
					c.dataManager.StoreID("last_movie_id", c.dataManager.GetLastMovieID())
				}
			}
			// Store actor ID if this was an add_actor call
			if toolName == "add_actor" {
				if err := c.dataManager.StoreIDFromResponse(responseData, "id", "actor_id"); err == nil {
					// Also store as last_actor_id for reference
					c.dataManager.StoreID("last_actor_id", c.dataManager.GetLastActorID())
				}
			}
		}
	}

	return nil
}

// theResponseShouldContainAMovieWith verifies movie response contains expected fields
func (c *CommonStepContext) theResponseShouldContainAMovieWith(table *godog.Table) error {
	var movie MovieResponse
	if err := c.bddContext.ParseJSONResponse(&movie); err != nil {
		return fmt.Errorf("failed to parse movie response: %w", err)
	}

	for _, row := range table.Rows {
		field := row.Cells[0].Value
		expectedValue := row.Cells[1].Value

		switch field {
		case "title":
			if movie.Title != expectedValue {
				return fmt.Errorf("expected title '%s', got '%s'", expectedValue, movie.Title)
			}
		case "director":
			if movie.Director != expectedValue {
				return fmt.Errorf("expected director '%s', got '%s'", expectedValue, movie.Director)
			}
		case "year":
			expectedYear, err := strconv.Atoi(expectedValue)
			if err != nil {
				return fmt.Errorf("invalid year value: %s", expectedValue)
			}
			if movie.Year != expectedYear {
				return fmt.Errorf("expected year %d, got %d", expectedYear, movie.Year)
			}
		case "rating":
			expectedRating, err := strconv.ParseFloat(expectedValue, 64)
			if err != nil {
				return fmt.Errorf("invalid rating value: %s", expectedValue)
			}
			if movie.Rating != expectedRating {
				return fmt.Errorf("expected rating %.1f, got %.1f", expectedRating, movie.Rating)
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}

	return nil
}

// theMovieShouldHaveAnAssignedID verifies the movie has a valid ID
func (c *CommonStepContext) theMovieShouldHaveAnAssignedID() error {
	var movie MovieResponse
	if err := c.bddContext.ParseJSONResponse(&movie); err != nil {
		return fmt.Errorf("failed to parse movie response: %w", err)
	}

	if movie.ID <= 0 {
		return fmt.Errorf("movie should have a positive ID, got %d", movie.ID)
	}

	return nil
}

// theResponseShouldContainNMovies verifies the response contains the expected number of movies
func (c *CommonStepContext) theResponseShouldContainNMovies(expectedCount int) error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single movie list
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	actualCount := len(response.Movies)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d movies, got %d", expectedCount, actualCount)
	}

	return nil
}

// theFollowingMoviesExist loads test movies from a data table
func (c *CommonStepContext) theFollowingMoviesExist(table *godog.Table) error {
	// Load movies fixture first
	err := c.testDB.LoadFixtures("movies")
	if err != nil {
		return fmt.Errorf("failed to load movies fixture: %w", err)
	}

	// Create movies from the table data
	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}

		movieData := make(map[string]interface{})
		for j, cell := range row.Cells {
			header := table.Rows[0].Cells[j].Value
			switch header {
			case "year":
				year, err := strconv.Atoi(cell.Value)
				if err != nil {
					return fmt.Errorf("invalid year: %s", cell.Value)
				}
				movieData[header] = year
			case "rating":
				rating, err := strconv.ParseFloat(cell.Value, 64)
				if err != nil {
					return fmt.Errorf("invalid rating: %s", cell.Value)
				}
				movieData[header] = rating
			default:
				movieData[header] = cell.Value
			}
		}

		// Call add_movie tool for each movie
		_, err := c.bddContext.CallTool("add_movie", movieData)
		if err != nil {
			return fmt.Errorf("failed to create movie %d: %w", i, err)
		}
	}

	return nil
}

// aMovieExistsWith creates a single movie from a data table
func (c *CommonStepContext) aMovieExistsWith(table *godog.Table) error {
	movieData := make(map[string]interface{})
	
	for _, row := range table.Rows {
		field := row.Cells[0].Value
		value := row.Cells[1].Value
		
		switch field {
		case "year":
			year, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid year: %s", value)
			}
			movieData[field] = year
		case "rating":
			rating, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid rating: %s", value)
			}
			movieData[field] = rating
		default:
			movieData[field] = value
		}
	}

	_, err := c.bddContext.CallTool("add_movie", movieData)
	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	return nil
}

// theMovieTitleShouldContain verifies movie titles contain the expected substring
func (c *CommonStepContext) theMovieTitleShouldContain(expectedSubstring string) error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single movie list
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) == 0 {
		return fmt.Errorf("no movies found in response")
	}

	for _, movie := range response.Movies {
		if !strings.Contains(movie.Title, expectedSubstring) {
			return fmt.Errorf("movie title '%s' does not contain '%s'", movie.Title, expectedSubstring)
		}
	}

	return nil
}

// allMoviesShouldHaveDirector verifies all movies have the expected director
func (c *CommonStepContext) allMoviesShouldHaveDirector(expectedDirector string) error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single movie list
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) == 0 {
		return fmt.Errorf("no movies found in response")
	}

	for _, movie := range response.Movies {
		if movie.Director != expectedDirector {
			return fmt.Errorf("movie '%s' has director '%s', expected '%s'", movie.Title, movie.Director, expectedDirector)
		}
	}

	return nil
}

// iCallToolWithMovieID calls a tool with the last stored movie ID
func (c *CommonStepContext) iCallToolWithMovieID(toolName string) error {
	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID available")
	}

	arguments := map[string]interface{}{
		"movie_id": movieID,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// iCallToolWithSpecificMovieID calls a tool with a specific movie ID
func (c *CommonStepContext) iCallToolWithSpecificMovieID(toolName string, movieID int) error {
	arguments := map[string]interface{}{
		"movie_id": movieID,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// theMovieTitleShouldBe verifies the movie title matches expected value
func (c *CommonStepContext) theMovieTitleShouldBe(expectedTitle string) error {
	var movie MovieResponse
	if err := c.bddContext.ParseJSONResponse(&movie); err != nil {
		return fmt.Errorf("failed to parse movie response: %w", err)
	}

	if movie.Title != expectedTitle {
		return fmt.Errorf("expected movie title '%s', got '%s'", expectedTitle, movie.Title)
	}

	return nil
}

// aMovieExistsWithTitle creates a movie with the specified title
func (c *CommonStepContext) aMovieExistsWithTitle(title string) error {
	arguments := map[string]interface{}{
		"title":    title,
		"director": "Test Director",
		"year":     2023,
		"rating":   7.5,
	}

	_, err := c.bddContext.CallTool("add_movie", arguments)
	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	// Store the movie ID
	if !c.bddContext.HasError() {
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr == nil {
			c.dataManager.StoreIDFromResponse(responseData, "id", "movie_id")
		}
	}

	return nil
}

// theMovieTitleShouldBeUpdatedTo verifies title after update
func (c *CommonStepContext) theMovieTitleShouldBeUpdatedTo(expectedTitle string) error {
	return c.theMovieTitleShouldBe(expectedTitle)
}

// theMovieRatingShouldBeUpdatedTo verifies rating after update
func (c *CommonStepContext) theMovieRatingShouldBeUpdatedTo(expectedRating float64) error {
	var movie MovieResponse
	if err := c.bddContext.ParseJSONResponse(&movie); err != nil {
		return fmt.Errorf("failed to parse movie response: %w", err)
	}

	if movie.Rating != expectedRating {
		return fmt.Errorf("expected movie rating %.1f, got %.1f", expectedRating, movie.Rating)
	}

	return nil
}

// theMovieShouldNoLongerExistInDatabase verifies movie deletion
func (c *CommonStepContext) theMovieShouldNoLongerExistInDatabase() error {
	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID to check")
	}

	// Try to get the movie - should fail
	_, err := c.bddContext.CallTool("get_movie", map[string]interface{}{
		"movie_id": movieID,
	})

	if err == nil && !c.bddContext.HasError() {
		return fmt.Errorf("movie still exists in database")
	}

	return nil
}

// iCallToolWithLimit calls a tool with a limit parameter
func (c *CommonStepContext) iCallToolWithLimit(toolName string, limit int) error {
	arguments := map[string]interface{}{
		"limit": limit,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// theMoviesShouldBeOrderedByRatingDescending verifies movies are ordered by rating
func (c *CommonStepContext) theMoviesShouldBeOrderedByRatingDescending() error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single movie list
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) < 2 {
		return nil // Can't verify ordering with less than 2 movies
	}

	for i := 1; i < len(response.Movies); i++ {
		if response.Movies[i-1].Rating < response.Movies[i].Rating {
			return fmt.Errorf("movies not ordered by rating descending: %.1f < %.1f", 
				response.Movies[i-1].Rating, response.Movies[i].Rating)
		}
	}

	return nil
}

// theFirstMovieShouldHaveRating verifies the first movie's rating
func (c *CommonStepContext) theFirstMovieShouldHaveRating(expectedRating float64) error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single movie list
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) == 0 {
		return fmt.Errorf("no movies in response")
	}

	if response.Movies[0].Rating != expectedRating {
		return fmt.Errorf("expected first movie rating %.1f, got %.1f", 
			expectedRating, response.Movies[0].Rating)
	}

	return nil
}

// theErrorMessageShouldIndicateMovieNotFound verifies movie not found error
func (c *CommonStepContext) theErrorMessageShouldIndicateMovieNotFound() error {
	if !c.bddContext.HasError() {
		return fmt.Errorf("expected error but got successful response")
	}

	errorMessage := c.bddContext.GetErrorMessage()
	if !strings.Contains(strings.ToLower(errorMessage), "not found") && 
	   !strings.Contains(strings.ToLower(errorMessage), "movie") {
		return fmt.Errorf("error message should indicate movie not found, got: %s", errorMessage)
	}

	return nil
}

// moviesExistFromVariousDecades creates movies from different decades
func (c *CommonStepContext) moviesExistFromVariousDecades() error {
	// Create movies from different decades for testing
	decadeMovies := []map[string]interface{}{
		{"title": "90s Movie 1", "director": "Director A", "year": 1993, "rating": 7.5},
		{"title": "90s Movie 2", "director": "Director B", "year": 1997, "rating": 8.0},
		{"title": "00s Movie 1", "director": "Director C", "year": 2003, "rating": 7.8},
		{"title": "00s Movie 2", "director": "Director D", "year": 2007, "rating": 8.2},
		{"title": "10s Movie 1", "director": "Director E", "year": 2013, "rating": 7.9},
		{"title": "10s Movie 2", "director": "Director F", "year": 2017, "rating": 8.5},
		{"title": "80s Movie", "director": "Director G", "year": 1985, "rating": 7.0},
		{"title": "20s Movie", "director": "Director H", "year": 2021, "rating": 8.1},
	}

	for _, movieData := range decadeMovies {
		_, err := c.bddContext.CallTool("add_movie", movieData)
		if err != nil {
			return fmt.Errorf("failed to create movie from year %v: %w", movieData["year"], err)
		}
	}

	return nil
}

// iCallToolWithDecade calls a tool with decade parameter
func (c *CommonStepContext) iCallToolWithDecade(toolName, decade string) error {
	arguments := map[string]interface{}{
		"decade": decade,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// allMoviesShouldBeFromYearsTo verifies all movies are within year range
func (c *CommonStepContext) allMoviesShouldBeFromYearsTo(minYear, maxYear int) error {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a simple movies array
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) == 0 {
		return fmt.Errorf("no movies found in response")
	}

	for _, movie := range response.Movies {
		if movie.Year < minYear || movie.Year > maxYear {
			return fmt.Errorf("movie '%s' has year %d, expected between %d and %d",
				movie.Title, movie.Year, minYear, maxYear)
		}
	}

	return nil
}

// theResponseShouldContainTheMovieDetails verifies movie details in response
func (c *CommonStepContext) theResponseShouldContainTheMovieDetails() error {
	var movie MovieResponse
	if err := c.bddContext.ParseJSONResponse(&movie); err != nil {
		return fmt.Errorf("failed to parse movie details: %w", err)
	}

	if movie.ID <= 0 {
		return fmt.Errorf("movie details should contain a valid ID, got %d", movie.ID)
	}

	if movie.Title == "" {
		return fmt.Errorf("movie details should contain a title")
	}

	if movie.Director == "" {
		return fmt.Errorf("movie details should contain a director")
	}

	if movie.Year <= 0 {
		return fmt.Errorf("movie details should contain a valid year, got %d", movie.Year)
	}

	return nil
}