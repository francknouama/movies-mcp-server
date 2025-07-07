package steps

import (
	"fmt"
)

// parseMoviesResponse parses the JSON response into a MoviesResponse struct
// It handles both direct movie arrays and wrapped responses
func parseMoviesResponse(c *CommonStepContext) (MoviesResponse, error) {
	var response MoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a simple movies array
		var movies []MovieResponse
		if err2 := c.bddContext.ParseJSONResponse(&movies); err2 != nil {
			return MoviesResponse{}, fmt.Errorf("failed to parse movies response: %w", err)
		}
		response.Movies = movies
	}

	if len(response.Movies) == 0 {
		return response, fmt.Errorf("no movies found in response")
	}

	return response, nil
}

// validateMovieRange validates that all movies in the response meet a specific condition
func validateMovieRange[T comparable](movies []MovieResponse, getValue func(MovieResponse) T, 
	isValid func(T, T, T) bool, min, max T, fieldName string) error {
	
	for _, movie := range movies {
		value := getValue(movie)
		if !isValid(value, min, max) {
			return fmt.Errorf("movie '%s' has %s %v, expected between %v and %v",
				movie.Title, fieldName, value, min, max)
		}
	}
	return nil
}