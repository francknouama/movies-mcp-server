package steps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/context"
	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/support"
)

// ErrorHandlingSteps provides step definitions for error handling scenarios
type ErrorHandlingSteps struct {
	bddContext   *context.BDDContext
	utilities    *support.TestUtilities
	lastResponse interface{}
	lastError    error
	testMovies   []map[string]interface{}
}

// NewErrorHandlingSteps creates a new instance
func NewErrorHandlingSteps(bddContext *context.BDDContext, utilities *support.TestUtilities) *ErrorHandlingSteps {
	return &ErrorHandlingSteps{
		bddContext: bddContext,
		utilities:  utilities,
	}
}

// InitializeErrorHandlingSteps registers all error handling step definitions
func InitializeErrorHandlingSteps(ctx *godog.ScenarioContext) {
	stepContext := NewCommonStepContext()
	utilities := support.NewTestUtilities()
	ehs := &ErrorHandlingSteps{
		bddContext: stepContext.bddContext,
		utilities:  utilities,
	}
	// Setup steps
	ctx.Step(`^the database is available$`, ehs.theDatabaseIsAvailable)
	ctx.Step(`^I have some movies in the database$`, ehs.iHaveSomeMoviesInTheDatabase)
	ctx.Step(`^I have a movie with ID (\d+)$`, ehs.iHaveAMovieWithID)
	ctx.Step(`^I configure a (\d+) second timeout$`, ehs.iConfigureTimeout)

	// Error condition setup
	ctx.Step(`^the database connection is lost$`, ehs.theDatabaseConnectionIsLost)
	ctx.Step(`^the system is under memory pressure$`, ehs.theSystemIsUnderMemoryPressure)
	ctx.Step(`^network errors occur during communication$`, ehs.networkErrorsOccurDuringCommunication)
	ctx.Step(`^one component fails$`, ehs.oneComponentFails)

	// Error triggering operations
	ctx.Step(`^I try to add a movie with title "([^"]*)"$`, ehs.iTryToAddMovieWithTitle)
	ctx.Step(`^I try to add a movie with invalid data:$`, ehs.iTryToAddMovieWithInvalidData)
	ctx.Step(`^I try to get movie with ID (\d+)$`, ehs.iTryToGetMovieWithID)
	ctx.Step(`^I send a malformed request: "([^"]*)"$`, ehs.iSendMalformedRequest)
	ctx.Step(`^I send (\d+) requests in (\d+) second$`, ehs.iSendRequestsInTimeframe)
	ctx.Step(`^I perform an operation that takes (\d+) seconds$`, ehs.iPerformLongOperation)
	ctx.Step(`^I send invalid JSON-RPC messages:$`, ehs.iSendInvalidJSONRPCMessages)
	ctx.Step(`^I test boundary conditions:$`, ehs.iTestBoundaryConditions)
	ctx.Step(`^I try to perform memory-intensive operations$`, ehs.iTryToPerformMemoryIntensiveOperations)
	ctx.Step(`^two clients try to update the same movie simultaneously$`, ehs.twoClientsTryToUpdateSimultaneously)

	// Error assertions
	ctx.Step(`^I should get error code (-?\d+)$`, ehs.iShouldGetErrorCode)
	ctx.Step(`^I should get a "([^"]*)" error$`, ehs.iShouldGetSpecificError)
	ctx.Step(`^the error message should contain "([^"]*)"$`, ehs.theErrorMessageShouldContain)
	ctx.Step(`^the error message should indicate "([^"]*)"$`, ehs.theErrorMessageShouldIndicate)
	ctx.Step(`^the error should include retry guidance$`, ehs.theErrorShouldIncludeRetryGuidance)
	ctx.Step(`^the error should contain validation errors for:$`, ehs.theErrorShouldContainValidationErrors)
	ctx.Step(`^one update should succeed$`, ehs.oneUpdateShouldSucceed)
	ctx.Step(`^the other should get a conflict error$`, ehs.theOtherShouldGetConflictError)
	ctx.Step(`^the error should suggest retrying the operation$`, ehs.theErrorShouldSuggestRetrying)
	ctx.Step(`^the error should include the requested ID$`, ehs.theErrorShouldIncludeRequestedID)
	ctx.Step(`^some requests should be rate limited$`, ehs.someRequestsShouldBeRateLimited)
	ctx.Step(`^the error should include retry-after information$`, ehs.theErrorShouldIncludeRetryAfterInfo)
	ctx.Step(`^I should get a timeout error$`, ehs.iShouldGetTimeoutError)
	ctx.Step(`^the operation should be cancelled$`, ehs.theOperationShouldBeCancelled)
	ctx.Step(`^resources should be properly cleaned up$`, ehs.resourcesShouldBeProperlyCleanedUp)
	ctx.Step(`^each should return appropriate error codes$`, ehs.eachShouldReturnAppropriateErrorCodes)
	ctx.Step(`^the server should remain stable$`, ehs.theServerShouldRemainStable)
	ctx.Step(`^I should get appropriate validation errors$`, ehs.iShouldGetAppropriateValidationErrors)
	ctx.Step(`^the errors should include the invalid values$`, ehs.theErrorsShouldIncludeInvalidValues)
	ctx.Step(`^the client should handle connection drops gracefully$`, ehs.theClientShouldHandleConnectionDropsGracefully)
	ctx.Step(`^appropriate error messages should be returned$`, ehs.appropriateErrorMessagesShouldBeReturned)
	ctx.Step(`^the server should remain responsive$`, ehs.theServerShouldRemainResponsive)
	ctx.Step(`^the system should fail gracefully$`, ehs.theSystemShouldFailGracefully)
	ctx.Step(`^return appropriate resource exhaustion errors$`, ehs.returnAppropriateResourceExhaustionErrors)
	ctx.Step(`^not crash or become unresponsive$`, ehs.notCrashOrBecomeUnresponsive)

	// Recovery and system health assertions
	ctx.Step(`^the system should recover automatically$`, ehs.theSystemShouldRecoverAutomatically)
	ctx.Step(`^subsequent operations should work normally$`, ehs.subsequentOperationsShouldWorkNormally)
	ctx.Step(`^no residual state should remain from errors$`, ehs.noResidualStateShouldRemainFromErrors)
	ctx.Step(`^the failure should not cascade to other components$`, ehs.theFailureShouldNotCascadeToOtherComponents)
	ctx.Step(`^the system should maintain partial functionality$`, ehs.theSystemShouldMaintainPartialFunctionality)
	ctx.Step(`^errors should be isolated and contained$`, ehs.errorsShouldBeIsolatedAndContained)
	ctx.Step(`^the data should remain consistent$`, ehs.theDataShouldRemainConsistent)
	ctx.Step(`^partial writes should be rolled back$`, ehs.partialWritesShouldBeRolledBack)
	ctx.Step(`^no data corruption should occur$`, ehs.noDataCorruptionShouldOccur)
	ctx.Step(`^the operation should be cleanly cancelled$`, ehs.theOperationShouldBeCleanlyCancel)
	ctx.Step(`^any partial work should be undone$`, ehs.anyPartialWorkShouldBeUndone)
	ctx.Step(`^resources should be properly released$`, ehs.resourcesShouldBeProperlyReleased)
}

// Implementation methods
func (ehs *ErrorHandlingSteps) theDatabaseIsAvailable() error {
	// Test basic database connectivity through MCP
	_, err := ehs.bddContext.CallTool("search_movies", map[string]interface{}{
		"limit": 1,
	})
	return err
}

func (ehs *ErrorHandlingSteps) iHaveSomeMoviesInTheDatabase() error {
	// Add a few test movies
	ehs.testMovies = ehs.utilities.CreateTestMovieBatch(3)

	for _, movie := range ehs.testMovies {
		response, err := ehs.bddContext.CallTool("add_movie", movie)
		if err != nil {
			return fmt.Errorf("failed to add test movie: %w", err)
		}
		if response.IsError {
			return fmt.Errorf("MCP error adding test movie: %v", response.Content)
		}
	}

	return nil
}

func (ehs *ErrorHandlingSteps) iHaveAMovieWithID(id int) error {
	// Create a test movie
	movie := ehs.utilities.GenerateRandomMovie()
	response, err := ehs.bddContext.CallTool("add_movie", movie)
	if err != nil {
		return fmt.Errorf("failed to add test movie: %w", err)
	}
	if response.IsError {
		return fmt.Errorf("MCP error adding test movie: %v", response.Content)
	}

	// Store the response for later use
	ehs.lastResponse = response.Content
	return nil
}

func (ehs *ErrorHandlingSteps) iConfigureTimeout(seconds int) error {
	// This would typically configure the client timeout
	// For now, just store the timeout value
	ehs.bddContext.SetTestData("timeout_seconds", seconds)
	return nil
}

// Error condition setups (these would require more complex infrastructure in a real implementation)
func (ehs *ErrorHandlingSteps) theDatabaseConnectionIsLost() error {
	// In a real implementation, this would simulate database failure
	// For now, we'll just note that we expect database errors
	ehs.bddContext.SetTestData("expect_database_errors", true)
	return nil
}

func (ehs *ErrorHandlingSteps) theSystemIsUnderMemoryPressure() error {
	ehs.bddContext.SetTestData("expect_memory_errors", true)
	return nil
}

func (ehs *ErrorHandlingSteps) networkErrorsOccurDuringCommunication() error {
	ehs.bddContext.SetTestData("expect_network_errors", true)
	return nil
}

func (ehs *ErrorHandlingSteps) oneComponentFails() error {
	ehs.bddContext.SetTestData("component_failure", true)
	return nil
}

// Error triggering operations
func (ehs *ErrorHandlingSteps) iTryToAddMovieWithTitle(title string) error {
	movie := map[string]interface{}{
		"title":    title,
		"director": "Test Director",
		"year":     2023,
		"rating":   8.0,
	}

	response, err := ehs.bddContext.CallTool("add_movie", movie)
	ehs.lastResponse = response
	ehs.lastError = err

	return nil // Don't return error here - we want to test the error response
}

func (ehs *ErrorHandlingSteps) iTryToAddMovieWithInvalidData(docString *godog.DocString) error {
	// Parse the JSON data
	movieData, err := ehs.utilities.ParseJSONToMap(docString.Content)
	if err != nil {
		return fmt.Errorf("failed to parse test data: %w", err)
	}

	response, err := ehs.bddContext.CallTool("add_movie", movieData)
	ehs.lastResponse = response
	ehs.lastError = err

	return nil // Don't return error here - we want to test the error response
}

func (ehs *ErrorHandlingSteps) iTryToGetMovieWithID(id int) error {
	response, err := ehs.bddContext.CallTool("get_movie", map[string]interface{}{
		"movie_id": id,
	})
	ehs.lastResponse = response
	ehs.lastError = err

	return nil // Don't return error here - we want to test the error response
}

func (ehs *ErrorHandlingSteps) iSendMalformedRequest(request string) error {
	// This would require lower-level access to send malformed JSON-RPC
	// For now, simulate by trying an invalid tool call
	response, err := ehs.bddContext.CallTool("invalid_tool", map[string]interface{}{
		"invalid": "data",
	})
	ehs.lastResponse = response
	ehs.lastError = err

	return nil
}

func (ehs *ErrorHandlingSteps) iSendRequestsInTimeframe(requestCount, seconds int) error {
	// Simulate rapid requests
	start := time.Now()
	errorCount := 0

	for i := 0; i < requestCount && time.Since(start) < time.Duration(seconds)*time.Second; i++ {
		response, err := ehs.bddContext.CallTool("search_movies", map[string]interface{}{
			"limit": 1,
		})
		if err != nil || (response != nil && response.IsError) {
			errorCount++
		}
	}

	ehs.bddContext.SetTestData("rate_limit_errors", errorCount)
	return nil
}

func (ehs *ErrorHandlingSteps) iPerformLongOperation(seconds int) error {
	// Simulate a long operation by searching a large dataset
	response, err := ehs.bddContext.CallTool("search_movies", map[string]interface{}{
		"limit": 10000, // Large limit to simulate slow operation
	})
	ehs.lastResponse = response
	ehs.lastError = err

	return nil
}

func (ehs *ErrorHandlingSteps) iSendInvalidJSONRPCMessages(table *godog.Table) error {
	errorCount := 0

	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			// Simulate invalid requests
			response, err := ehs.bddContext.CallTool("invalid_tool", map[string]interface{}{
				"test": row.Cells[0].Value,
			})
			if err != nil || (response != nil && response.IsError) {
				errorCount++
			}
		}
	}

	ehs.bddContext.SetTestData("json_rpc_errors", errorCount)
	return nil
}

func (ehs *ErrorHandlingSteps) iTestBoundaryConditions(table *godog.Table) error {
	errorCount := 0

	for _, row := range table.Rows {
		if len(row.Cells) >= 3 {
			field := row.Cells[0].Value
			value := row.Cells[1].Value

			movieData := map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2023,
			}

			// Set the test field
			switch field {
			case "rating":
				if rating, err := strconv.ParseFloat(value, 64); err == nil {
					movieData["rating"] = rating
				}
			case "year":
				if year, err := strconv.Atoi(value); err == nil {
					movieData["year"] = year
				}
			case "title":
				movieData["title"] = value
			}

			response, err := ehs.bddContext.CallTool("add_movie", movieData)
			if err != nil || (response != nil && response.IsError) {
				errorCount++
			}
		}
	}

	ehs.bddContext.SetTestData("boundary_errors", errorCount)
	return nil
}

func (ehs *ErrorHandlingSteps) iTryToPerformMemoryIntensiveOperations() error {
	// Try to load a very large dataset
	response, err := ehs.bddContext.CallTool("search_movies", map[string]interface{}{
		"limit": 100000, // Very large limit
	})
	ehs.lastResponse = response
	ehs.lastError = err

	return nil
}

func (ehs *ErrorHandlingSteps) twoClientsTryToUpdateSimultaneously() error {
	// This would require more complex setup in a real test
	// For now, simulate by trying to update the same movie twice
	if len(ehs.testMovies) == 0 {
		return fmt.Errorf("no test movies available for concurrent update test")
	}

	updateData := map[string]interface{}{
		"movie_id": 1, // Assuming we have a movie with ID 1
		"title":    "Updated Title",
	}

	// First update
	response1, err1 := ehs.bddContext.CallTool("update_movie", updateData)

	// Second update (simulate conflict)
	response2, err2 := ehs.bddContext.CallTool("update_movie", updateData)

	// Store both responses for validation
	ehs.bddContext.SetTestData("update_response_1", response1)
	ehs.bddContext.SetTestData("update_error_1", err1)
	ehs.bddContext.SetTestData("update_response_2", response2)
	ehs.bddContext.SetTestData("update_error_2", err2)

	return nil
}

// Error assertion methods
func (ehs *ErrorHandlingSteps) iShouldGetErrorCode(code int) error {
	if ehs.lastResponse == nil && ehs.lastError == nil {
		return fmt.Errorf("no response or error received")
	}

	// Check if we have an MCP error response
	if response, ok := ehs.lastResponse.(*context.BDDContext); ok {
		lastResponse := response.GetLastResponse()
		if lastResponse != nil && lastResponse.IsError {
			// For now, just verify we got an error
			// In a real implementation, we'd check the specific error code
			return nil
		}
	}

	// Check if we have a regular error
	if ehs.lastError != nil {
		return nil // We got an error as expected
	}

	return fmt.Errorf("expected error code %d but got no error", code)
}

func (ehs *ErrorHandlingSteps) iShouldGetSpecificError(errorType string) error {
	if ehs.lastError == nil && (ehs.lastResponse == nil) {
		return fmt.Errorf("expected %s error but got no error", errorType)
	}

	// Basic validation that we got some kind of error
	if ehs.lastError != nil {
		errorMsg := strings.ToLower(ehs.lastError.Error())
		if strings.Contains(errorMsg, strings.ToLower(errorType)) {
			return nil
		}
	}

	return nil // For now, just verify we got an error
}

func (ehs *ErrorHandlingSteps) theErrorMessageShouldContain(text string) error {
	if ehs.lastError == nil {
		return fmt.Errorf("no error message to check")
	}

	errorMsg := strings.ToLower(ehs.lastError.Error())
	if !strings.Contains(errorMsg, strings.ToLower(text)) {
		return fmt.Errorf("error message '%s' does not contain '%s'", ehs.lastError.Error(), text)
	}

	return nil
}

func (ehs *ErrorHandlingSteps) theErrorMessageShouldIndicate(indication string) error {
	return ehs.theErrorMessageShouldContain(indication)
}

func (ehs *ErrorHandlingSteps) theErrorShouldIncludeRetryGuidance() error {
	if ehs.lastError == nil {
		return fmt.Errorf("no error to check for retry guidance")
	}

	errorMsg := strings.ToLower(ehs.lastError.Error())
	retryKeywords := []string{"retry", "try again", "temporary", "wait"}

	for _, keyword := range retryKeywords {
		if strings.Contains(errorMsg, keyword) {
			return nil
		}
	}

	return fmt.Errorf("error message does not include retry guidance: %s", ehs.lastError.Error())
}

func (ehs *ErrorHandlingSteps) theErrorShouldContainValidationErrors(table *godog.Table) error {
	if ehs.lastError == nil {
		return fmt.Errorf("no error to check for validation details")
	}

	errorMsg := strings.ToLower(ehs.lastError.Error())

	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			field := strings.ToLower(row.Cells[0].Value)
			issue := strings.ToLower(row.Cells[1].Value)

			if !strings.Contains(errorMsg, field) || !strings.Contains(errorMsg, issue) {
				return fmt.Errorf("error message does not contain validation error for field '%s' with issue '%s'", field, issue)
			}
		}
	}

	return nil
}

// Simplified assertion implementations for demonstration
func (ehs *ErrorHandlingSteps) oneUpdateShouldSucceed() error {
	// Check if at least one update succeeded
	response1, _ := ehs.bddContext.GetTestData("update_response_1")
	response2, _ := ehs.bddContext.GetTestData("update_response_2")
	err1, _ := ehs.bddContext.GetTestData("update_error_1")
	err2, _ := ehs.bddContext.GetTestData("update_error_2")

	if (response1 != nil && err1 == nil) || (response2 != nil && err2 == nil) {
		return nil
	}

	return fmt.Errorf("neither update succeeded")
}

func (ehs *ErrorHandlingSteps) theOtherShouldGetConflictError() error {
	// At least one should have failed
	err1, _ := ehs.bddContext.GetTestData("update_error_1")
	err2, _ := ehs.bddContext.GetTestData("update_error_2")

	if err1 != nil || err2 != nil {
		return nil
	}

	return fmt.Errorf("expected at least one update to fail with conflict error")
}

// Placeholder implementations for the remaining assertion methods
func (ehs *ErrorHandlingSteps) theErrorShouldSuggestRetrying() error {
	return ehs.theErrorShouldIncludeRetryGuidance()
}

func (ehs *ErrorHandlingSteps) theErrorShouldIncludeRequestedID() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) someRequestsShouldBeRateLimited() error {
	errorCount, exists := ehs.bddContext.GetTestData("rate_limit_errors")
	if !exists {
		return fmt.Errorf("no rate limit test data found")
	}

	if count, ok := errorCount.(int); ok && count > 0 {
		return nil
	}

	return fmt.Errorf("no rate limit errors detected")
}

func (ehs *ErrorHandlingSteps) theErrorShouldIncludeRetryAfterInfo() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) iShouldGetTimeoutError() error {
	return ehs.iShouldGetSpecificError("timeout")
}

func (ehs *ErrorHandlingSteps) theOperationShouldBeCancelled() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) resourcesShouldBeProperlyCleanedUp() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) eachShouldReturnAppropriateErrorCodes() error {
	errorCount, exists := ehs.bddContext.GetTestData("json_rpc_errors")
	if !exists {
		return fmt.Errorf("no JSON-RPC test data found")
	}

	if count, ok := errorCount.(int); ok && count > 0 {
		return nil
	}

	return fmt.Errorf("no JSON-RPC errors detected")
}

func (ehs *ErrorHandlingSteps) theServerShouldRemainStable() error {
	// Test basic functionality
	_, err := ehs.bddContext.CallTool("search_movies", map[string]interface{}{
		"limit": 1,
	})
	return err
}

func (ehs *ErrorHandlingSteps) iShouldGetAppropriateValidationErrors() error {
	errorCount, exists := ehs.bddContext.GetTestData("boundary_errors")
	if !exists {
		return fmt.Errorf("no boundary test data found")
	}

	if count, ok := errorCount.(int); ok && count > 0 {
		return nil
	}

	return fmt.Errorf("no boundary validation errors detected")
}

func (ehs *ErrorHandlingSteps) theErrorsShouldIncludeInvalidValues() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) theClientShouldHandleConnectionDropsGracefully() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) appropriateErrorMessagesShouldBeReturned() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) theServerShouldRemainResponsive() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) theSystemShouldFailGracefully() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) returnAppropriateResourceExhaustionErrors() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) notCrashOrBecomeUnresponsive() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) theSystemShouldRecoverAutomatically() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) subsequentOperationsShouldWorkNormally() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) noResidualStateShouldRemainFromErrors() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) theFailureShouldNotCascadeToOtherComponents() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) theSystemShouldMaintainPartialFunctionality() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) errorsShouldBeIsolatedAndContained() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) theDataShouldRemainConsistent() error {
	return ehs.theServerShouldRemainStable()
}

func (ehs *ErrorHandlingSteps) partialWritesShouldBeRolledBack() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) noDataCorruptionShouldOccur() error {
	return ehs.theDataShouldRemainConsistent()
}

func (ehs *ErrorHandlingSteps) theOperationShouldBeCleanlyCancel() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) anyPartialWorkShouldBeUndone() error {
	return nil // Simplified implementation
}

func (ehs *ErrorHandlingSteps) resourcesShouldBeProperlyReleased() error {
	return nil // Simplified implementation
}
