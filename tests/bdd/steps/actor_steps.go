package steps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/francknouama/movies-mcp-server/legacy/internal/interfaces/dto"
)

// MovieCastResponse represents the cast of a movie
type MovieCastResponse struct {
	MovieID int                  `json:"movie_id"`
	Cast    []*dto.ActorResponse `json:"cast"`
}

// ActorMoviesResponse represents movies associated with an actor
type ActorMoviesResponse struct {
	ActorID  int   `json:"actor_id"`
	MovieIDs []int `json:"movie_ids"`
}

// InitializeActorSteps registers actor-related step definitions
func InitializeActorSteps(ctx *godog.ScenarioContext) {
	stepContext := NewCommonStepContext()

	// Actor CRUD operations
	ctx.Step(`^the response should contain an actor with:$`, stepContext.theResponseShouldContainAnActorWith)
	ctx.Step(`^the actor should have an assigned ID$`, stepContext.theActorShouldHaveAnAssignedID)
	ctx.Step(`^an actor exists with:$`, stepContext.anActorExistsWith)
	ctx.Step(`^an actor exists with name "([^"]*)"$`, stepContext.anActorExistsWithName)
	ctx.Step(`^the response should contain the actor details$`, stepContext.theResponseShouldContainTheActorDetails)
	ctx.Step(`^the actor name should be "([^"]*)"$`, stepContext.theActorNameShouldBe)
	ctx.Step(`^the actor name should be updated to "([^"]*)"$`, stepContext.theActorNameShouldBeUpdatedTo)
	ctx.Step(`^the actor bio should be updated to "([^"]*)"$`, stepContext.theActorBioShouldBeUpdatedTo)
	ctx.Step(`^the actor should no longer exist in the database$`, stepContext.theActorShouldNoLongerExistInDatabase)
	ctx.Step(`^I call the "([^"]*)" tool with the actor ID$`, stepContext.iCallToolWithActorID)

	// Actor-Movie relationships
	ctx.Step(`^the actor should be linked to the movie$`, stepContext.theActorShouldBeLinkedToTheMovie)
	ctx.Step(`^the message should indicate successful linking$`, stepContext.theMessageShouldIndicateSuccessfulLinking)
	ctx.Step(`^the following actors are linked to the movie:$`, stepContext.theFollowingActorsAreLinkedToTheMovie)
	ctx.Step(`^the response should contain (\d+) actors?$`, stepContext.theResponseShouldContainNActors)
	ctx.Step(`^the cast should include:$`, stepContext.theCastShouldInclude)
	ctx.Step(`^the actor is linked to the following movies:$`, stepContext.theActorIsLinkedToTheFollowingMovies)
	ctx.Step(`^the response should contain (\d+) movie IDs$`, stepContext.theResponseShouldContainNMovieIDs)
	ctx.Step(`^the actor should be associated with all linked movies$`, stepContext.theActorShouldBeAssociatedWithAllLinkedMovies)
	ctx.Step(`^the actor is linked to the movie$`, stepContext.theActorIsLinkedToTheMovie)
	ctx.Step(`^the actor should no longer be linked to the movie$`, stepContext.theActorShouldNoLongerBeLinkedToTheMovie)
	ctx.Step(`^the actor is already linked to the movie$`, stepContext.theActorIsAlreadyLinkedToTheMovie)
	ctx.Step(`^I call the "([^"]*)" tool with the same actor and movie$`, stepContext.iCallToolWithTheSameActorAndMovie)

	// Actor search operations
	ctx.Step(`^the following actors exist:$`, stepContext.theFollowingActorsExist)
	ctx.Step(`^all actor names should contain "([^"]*)"$`, stepContext.allActorNamesShouldContain)
	ctx.Step(`^all actors should have birth year between (\d+) and (\d+)$`, stepContext.allActorsShouldHaveBirthYearBetween)

	// Error handling
	ctx.Step(`^the error should contain validation errors for:$`, stepContext.theErrorShouldContainValidationErrorsFor)
	ctx.Step(`^the error message should indicate the relationship already exists$`, stepContext.theErrorMessageShouldIndicateRelationshipAlreadyExists)
}

// theResponseShouldContainAnActorWith verifies actor response contains expected fields
func (c *CommonStepContext) theResponseShouldContainAnActorWith(table *godog.Table) error {
	var actor dto.ActorResponse
	if err := c.bddContext.ParseJSONResponse(&actor); err != nil {
		return fmt.Errorf("failed to parse actor response: %w", err)
	}

	for _, row := range table.Rows {
		field := row.Cells[0].Value
		expectedValue := row.Cells[1].Value

		switch field {
		case "name":
			if actor.Name != expectedValue {
				return fmt.Errorf("expected name '%s', got '%s'", expectedValue, actor.Name)
			}
		case "birth_year":
			expectedYear, err := strconv.Atoi(expectedValue)
			if err != nil {
				return fmt.Errorf("invalid birth_year value: %s", expectedValue)
			}
			if actor.BirthYear != expectedYear {
				return fmt.Errorf("expected birth_year %d, got %d", expectedYear, actor.BirthYear)
			}
		case "bio":
			if actor.Bio != expectedValue {
				return fmt.Errorf("expected bio '%s', got '%s'", expectedValue, actor.Bio)
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}

	return nil
}

// theActorShouldHaveAnAssignedID verifies the actor has a valid ID
func (c *CommonStepContext) theActorShouldHaveAnAssignedID() error {
	var actor dto.ActorResponse
	if err := c.bddContext.ParseJSONResponse(&actor); err != nil {
		return fmt.Errorf("failed to parse actor response: %w", err)
	}

	if actor.ID <= 0 {
		return fmt.Errorf("actor should have a positive ID, got %d", actor.ID)
	}

	return nil
}

// theResponseShouldContainNActors verifies the response contains the expected number of actors
func (c *CommonStepContext) theResponseShouldContainNActors(expectedCount int) error {
	var response dto.ActorsListResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a single actor list
		var actors []*dto.ActorResponse
		if err2 := c.bddContext.ParseJSONResponse(&actors); err2 != nil {
			return fmt.Errorf("failed to parse actors response: %w", err)
		}
		response.Actors = actors
	}

	actualCount := len(response.Actors)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d actors, got %d", expectedCount, actualCount)
	}

	return nil
}

// theFollowingActorsExist loads test actors from a data table
func (c *CommonStepContext) theFollowingActorsExist(table *godog.Table) error {
	// Load actors fixture first
	err := c.testDB.LoadFixtures("actors")
	if err != nil {
		return fmt.Errorf("failed to load actors fixture: %w", err)
	}

	// Create actors from the table data
	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}

		actorData := make(map[string]interface{})
		for j, cell := range row.Cells {
			header := table.Rows[0].Cells[j].Value
			switch header {
			case "birth_year":
				year, err := strconv.Atoi(cell.Value)
				if err != nil {
					return fmt.Errorf("invalid birth_year: %s", cell.Value)
				}
				actorData[header] = year
			default:
				actorData[header] = cell.Value
			}
		}

		// Set default values if not provided
		if _, exists := actorData["bio"]; !exists {
			actorData["bio"] = "Test actor biography"
		}

		// Call add_actor tool for each actor
		_, err := c.bddContext.CallTool("add_actor", actorData)
		if err != nil {
			return fmt.Errorf("failed to create actor %d: %w", i, err)
		}
	}

	return nil
}

// anActorExistsWith creates a single actor from a data table
func (c *CommonStepContext) anActorExistsWith(table *godog.Table) error {
	actorData := make(map[string]interface{})

	for _, row := range table.Rows {
		field := row.Cells[0].Value
		value := row.Cells[1].Value

		switch field {
		case "birth_year":
			year, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid birth_year: %s", value)
			}
			actorData[field] = year
		default:
			actorData[field] = value
		}
	}

	_, err := c.bddContext.CallTool("add_actor", actorData)
	if err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}

	// Store the actor ID
	if !c.bddContext.HasError() {
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr == nil {
			_ = c.dataManager.StoreIDFromResponse(responseData, "id", "actor_id")
		}
	}

	return nil
}

// anActorExistsWithName creates an actor with the specified name
func (c *CommonStepContext) anActorExistsWithName(name string) error {
	arguments := map[string]interface{}{
		"name":       name,
		"birth_year": 1980,
		"bio":        "Test actor biography",
	}

	_, err := c.bddContext.CallTool("add_actor", arguments)
	if err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}

	// Store the actor ID
	if !c.bddContext.HasError() {
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr == nil {
			_ = c.dataManager.StoreIDFromResponse(responseData, "id", "actor_id")
		}
	}

	return nil
}

// theResponseShouldContainTheActorDetails verifies the response contains actor details
func (c *CommonStepContext) theResponseShouldContainTheActorDetails() error {
	var actor dto.ActorResponse
	if err := c.bddContext.ParseJSONResponse(&actor); err != nil {
		return fmt.Errorf("failed to parse actor details: %w", err)
	}

	if actor.ID <= 0 {
		return fmt.Errorf("actor details should contain a valid ID, got %d", actor.ID)
	}

	if actor.Name == "" {
		return fmt.Errorf("actor details should contain a name")
	}

	return nil
}

// theActorNameShouldBe verifies the actor name matches expected value
func (c *CommonStepContext) theActorNameShouldBe(expectedName string) error {
	var actor dto.ActorResponse
	if err := c.bddContext.ParseJSONResponse(&actor); err != nil {
		return fmt.Errorf("failed to parse actor response: %w", err)
	}

	if actor.Name != expectedName {
		return fmt.Errorf("expected actor name '%s', got '%s'", expectedName, actor.Name)
	}

	return nil
}

// theActorNameShouldBeUpdatedTo verifies name after update
func (c *CommonStepContext) theActorNameShouldBeUpdatedTo(expectedName string) error {
	return c.theActorNameShouldBe(expectedName)
}

// theActorBioShouldBeUpdatedTo verifies bio after update
func (c *CommonStepContext) theActorBioShouldBeUpdatedTo(expectedBio string) error {
	var actor dto.ActorResponse
	if err := c.bddContext.ParseJSONResponse(&actor); err != nil {
		return fmt.Errorf("failed to parse actor response: %w", err)
	}

	if actor.Bio != expectedBio {
		return fmt.Errorf("expected actor bio '%s', got '%s'", expectedBio, actor.Bio)
	}

	return nil
}

// theActorShouldNoLongerExistInDatabase verifies actor deletion
func (c *CommonStepContext) theActorShouldNoLongerExistInDatabase() error {
	actorID := c.dataManager.GetLastActorID()
	if actorID == 0 {
		return fmt.Errorf("no actor ID to check")
	}

	// Try to get the actor - should fail
	_, err := c.bddContext.CallTool("get_actor", map[string]interface{}{
		"actor_id": actorID,
	})

	if err == nil && !c.bddContext.HasError() {
		return fmt.Errorf("actor still exists in database")
	}

	return nil
}

// iCallToolWithActorID calls a tool with the last stored actor ID
func (c *CommonStepContext) iCallToolWithActorID(toolName string) error {
	actorID := c.dataManager.GetLastActorID()
	if actorID == 0 {
		return fmt.Errorf("no actor ID available")
	}

	arguments := map[string]interface{}{
		"actor_id": actorID,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// theActorShouldBeLinkedToTheMovie verifies actor-movie relationship
func (c *CommonStepContext) theActorShouldBeLinkedToTheMovie() error {
	// This step verifies that the linking was successful
	if c.bddContext.HasError() {
		return fmt.Errorf("linking failed: %s", c.bddContext.GetErrorMessage())
	}
	return nil
}

// theMessageShouldIndicateSuccessfulLinking verifies success message
func (c *CommonStepContext) theMessageShouldIndicateSuccessfulLinking() error {
	if c.bddContext.HasError() {
		return fmt.Errorf("expected success message but got error: %s", c.bddContext.GetErrorMessage())
	}

	// The response should indicate success
	response := c.bddContext.GetLastResponse()
	if response == nil {
		return fmt.Errorf("no response received")
	}

	return nil
}

// theFollowingActorsAreLinkedToTheMovie creates actors and links them to the current movie
func (c *CommonStepContext) theFollowingActorsAreLinkedToTheMovie(table *godog.Table) error {
	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID available for linking")
	}

	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}

		actorName := row.Cells[0].Value

		// Create the actor
		actorData := map[string]interface{}{
			"name":       actorName,
			"birth_year": 1980,
			"bio":        "Test actor",
		}

		_, err := c.bddContext.CallTool("add_actor", actorData)
		if err != nil {
			return fmt.Errorf("failed to create actor %s: %w", actorName, err)
		}

		// Get the actor ID from response
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr != nil {
			return fmt.Errorf("failed to parse actor response: %w", parseErr)
		}

		actorID, err := c.dataManager.ParseIDFromResponse(responseData, "id")
		if err != nil {
			return fmt.Errorf("failed to get actor ID: %w", err)
		}

		// Link actor to movie
		linkData := map[string]interface{}{
			"actor_id": actorID,
			"movie_id": movieID,
		}

		_, err = c.bddContext.CallTool("link_actor_to_movie", linkData)
		if err != nil {
			return fmt.Errorf("failed to link actor %s to movie: %w", actorName, err)
		}
	}

	return nil
}

// theCastShouldInclude verifies the cast includes expected actors
func (c *CommonStepContext) theCastShouldInclude(table *godog.Table) error {
	var response MovieCastResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as actors array
		var actors []*dto.ActorResponse
		if err2 := c.bddContext.ParseJSONResponse(&actors); err2 != nil {
			return fmt.Errorf("failed to parse cast response: %w", err)
		}
		response.Cast = actors
	}

	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}

		expectedName := row.Cells[0].Value
		found := false

		for _, actor := range response.Cast {
			if actor.Name == expectedName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("expected actor '%s' not found in cast", expectedName)
		}
	}

	return nil
}

// theActorIsLinkedToTheFollowingMovies creates movies and links them to the current actor
func (c *CommonStepContext) theActorIsLinkedToTheFollowingMovies(table *godog.Table) error {
	actorID := c.dataManager.GetLastActorID()
	if actorID == 0 {
		return fmt.Errorf("no actor ID available for linking")
	}

	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}

		movieTitle := row.Cells[0].Value

		// Create the movie
		movieData := map[string]interface{}{
			"title":    movieTitle,
			"director": "Test Director",
			"year":     2023,
			"rating":   7.5,
		}

		_, err := c.bddContext.CallTool("add_movie", movieData)
		if err != nil {
			return fmt.Errorf("failed to create movie %s: %w", movieTitle, err)
		}

		// Get the movie ID from response
		var responseData map[string]interface{}
		if parseErr := c.bddContext.ParseJSONResponse(&responseData); parseErr != nil {
			return fmt.Errorf("failed to parse movie response: %w", parseErr)
		}

		movieID, err := c.dataManager.ParseIDFromResponse(responseData, "id")
		if err != nil {
			return fmt.Errorf("failed to get movie ID: %w", err)
		}

		// Link actor to movie
		linkData := map[string]interface{}{
			"actor_id": actorID,
			"movie_id": movieID,
		}

		_, err = c.bddContext.CallTool("link_actor_to_movie", linkData)
		if err != nil {
			return fmt.Errorf("failed to link actor to movie %s: %w", movieTitle, err)
		}
	}

	return nil
}

// theResponseShouldContainNMovieIDs verifies the response contains the expected number of movie IDs
func (c *CommonStepContext) theResponseShouldContainNMovieIDs(expectedCount int) error {
	var response ActorMoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		return fmt.Errorf("failed to parse actor movies response: %w", err)
	}

	actualCount := len(response.MovieIDs)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d movie IDs, got %d", expectedCount, actualCount)
	}

	return nil
}

// theActorShouldBeAssociatedWithAllLinkedMovies verifies all movies are linked
func (c *CommonStepContext) theActorShouldBeAssociatedWithAllLinkedMovies() error {
	var response ActorMoviesResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		return fmt.Errorf("failed to parse actor movies response: %w", err)
	}

	if len(response.MovieIDs) == 0 {
		return fmt.Errorf("actor should be associated with movies but found none")
	}

	// Verify all movie IDs are positive
	for _, movieID := range response.MovieIDs {
		if movieID <= 0 {
			return fmt.Errorf("invalid movie ID: %d", movieID)
		}
	}

	return nil
}

// allActorNamesShouldContain verifies all actor names contain the expected substring
func (c *CommonStepContext) allActorNamesShouldContain(expectedSubstring string) error {
	var response dto.ActorsListResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a simple actors array
		var actors []*dto.ActorResponse
		if err2 := c.bddContext.ParseJSONResponse(&actors); err2 != nil {
			return fmt.Errorf("failed to parse actors response: %w", err)
		}
		response.Actors = actors
	}

	if len(response.Actors) == 0 {
		return fmt.Errorf("no actors found in response")
	}

	for _, actor := range response.Actors {
		if !strings.Contains(actor.Name, expectedSubstring) {
			return fmt.Errorf("actor name '%s' does not contain '%s'", actor.Name, expectedSubstring)
		}
	}

	return nil
}

// allActorsShouldHaveBirthYearBetween verifies all actors have birth year in range
func (c *CommonStepContext) allActorsShouldHaveBirthYearBetween(minYear, maxYear int) error {
	var response dto.ActorsListResponse
	if err := c.bddContext.ParseJSONResponse(&response); err != nil {
		// Try parsing as a simple actors array
		var actors []*dto.ActorResponse
		if err2 := c.bddContext.ParseJSONResponse(&actors); err2 != nil {
			return fmt.Errorf("failed to parse actors response: %w", err)
		}
		response.Actors = actors
	}

	if len(response.Actors) == 0 {
		return fmt.Errorf("no actors found in response")
	}

	for _, actor := range response.Actors {
		if actor.BirthYear < minYear || actor.BirthYear > maxYear {
			return fmt.Errorf("actor '%s' has birth year %d, expected between %d and %d",
				actor.Name, actor.BirthYear, minYear, maxYear)
		}
	}

	return nil
}

// theActorIsLinkedToTheMovie links the current actor to the current movie
func (c *CommonStepContext) theActorIsLinkedToTheMovie() error {
	actorID := c.dataManager.GetLastActorID()
	if actorID == 0 {
		return fmt.Errorf("no actor ID available")
	}

	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID available")
	}

	linkData := map[string]interface{}{
		"actor_id": actorID,
		"movie_id": movieID,
	}

	_, err := c.bddContext.CallTool("link_actor_to_movie", linkData)
	if err != nil {
		return fmt.Errorf("failed to link actor to movie: %w", err)
	}

	return nil
}

// theActorShouldNoLongerBeLinkedToTheMovie verifies actor-movie relationship removal
func (c *CommonStepContext) theActorShouldNoLongerBeLinkedToTheMovie() error {
	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID available to check")
	}

	// Get the movie cast and verify actor is not in it
	_, err := c.bddContext.CallTool("get_movie_cast", map[string]interface{}{
		"movie_id": movieID,
	})

	if err != nil {
		return fmt.Errorf("failed to get movie cast: %w", err)
	}

	var castResponse MovieCastResponse
	if parseErr := c.bddContext.ParseJSONResponse(&castResponse); parseErr != nil {
		return fmt.Errorf("failed to parse cast response: %w", parseErr)
	}

	actorID := c.dataManager.GetLastActorID()
	for _, actor := range castResponse.Cast {
		if actor.ID == actorID {
			return fmt.Errorf("actor should no longer be linked to movie but was found in cast")
		}
	}

	return nil
}

// theActorIsAlreadyLinkedToTheMovie ensures actor-movie relationship exists
func (c *CommonStepContext) theActorIsAlreadyLinkedToTheMovie() error {
	// This is the same as theActorIsLinkedToTheMovie but used in different context
	return c.theActorIsLinkedToTheMovie()
}

// iCallToolWithTheSameActorAndMovie calls a tool with the same actor and movie IDs
func (c *CommonStepContext) iCallToolWithTheSameActorAndMovie(toolName string) error {
	actorID := c.dataManager.GetLastActorID()
	if actorID == 0 {
		return fmt.Errorf("no actor ID available")
	}

	movieID := c.dataManager.GetLastMovieID()
	if movieID == 0 {
		return fmt.Errorf("no movie ID available")
	}

	arguments := map[string]interface{}{
		"actor_id": actorID,
		"movie_id": movieID,
	}

	_, err := c.bddContext.CallTool(toolName, arguments)
	return err
}

// theErrorShouldContainValidationErrorsFor verifies validation errors
func (c *CommonStepContext) theErrorShouldContainValidationErrorsFor(table *godog.Table) error {
	if !c.bddContext.HasError() {
		return fmt.Errorf("expected validation errors but got successful response")
	}

	errorMessage := c.bddContext.GetErrorMessage()

	for _, row := range table.Rows {
		field := row.Cells[0].Value
		expectedIssue := row.Cells[1].Value

		// Check if the error message contains information about this field and issue
		if !strings.Contains(strings.ToLower(errorMessage), strings.ToLower(field)) {
			return fmt.Errorf("validation error should mention field '%s', got: %s", field, errorMessage)
		}

		if !strings.Contains(strings.ToLower(errorMessage), strings.ToLower(expectedIssue)) {
			return fmt.Errorf("validation error should mention issue '%s' for field '%s', got: %s",
				expectedIssue, field, errorMessage)
		}
	}

	return nil
}

// theErrorMessageShouldIndicateRelationshipAlreadyExists verifies duplicate relationship error
func (c *CommonStepContext) theErrorMessageShouldIndicateRelationshipAlreadyExists() error {
	if !c.bddContext.HasError() {
		return fmt.Errorf("expected error but got successful response")
	}

	errorMessage := c.bddContext.GetErrorMessage()
	if !strings.Contains(strings.ToLower(errorMessage), "already") &&
		!strings.Contains(strings.ToLower(errorMessage), "exists") &&
		!strings.Contains(strings.ToLower(errorMessage), "duplicate") {
		return fmt.Errorf("error message should indicate relationship already exists, got: %s", errorMessage)
	}

	return nil
}
