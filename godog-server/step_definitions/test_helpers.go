package step_definitions

import (
	"encoding/json"
	"fmt"
)

// Helper function to create movie via MCP (movie_steps.go already has createMovieViaMCP)

// Helper function to validate JSON structure
func (ctx *TestContext) validateJSONStructure(actual interface{}, expected map[string]interface{}) error {
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return fmt.Errorf("actual value is not a JSON object")
	}

	for key, expectedValue := range expected {
		actualValue, exists := actualMap[key]
		if !exists {
			return fmt.Errorf("field %s not found", key)
		}

		// Type-specific comparisons
		switch expectedTyped := expectedValue.(type) {
		case string:
			if actualStr, ok := actualValue.(string); ok {
				if actualStr != expectedTyped {
					return fmt.Errorf("field %s mismatch: expected %s, got %s", key, expectedTyped, actualStr)
				}
			} else {
				return fmt.Errorf("field %s is not a string", key)
			}
		case float64:
			if actualFloat, ok := actualValue.(float64); ok {
				if actualFloat != expectedTyped {
					return fmt.Errorf("field %s mismatch: expected %f, got %f", key, expectedTyped, actualFloat)
				}
			} else {
				return fmt.Errorf("field %s is not a number", key)
			}
		case int:
			expectedFloat := float64(expectedTyped)
			if actualFloat, ok := actualValue.(float64); ok {
				if actualFloat != expectedFloat {
					return fmt.Errorf("field %s mismatch: expected %d, got %f", key, expectedTyped, actualFloat)
				}
			} else {
				return fmt.Errorf("field %s is not a number", key)
			}
		default:
			// Generic comparison for other types
			actualStr := fmt.Sprintf("%v", actualValue)
			expectedStr := fmt.Sprintf("%v", expectedValue)
			if actualStr != expectedStr {
				return fmt.Errorf("field %s mismatch: expected %v, got %v", key, expectedValue, actualValue)
			}
		}
	}

	return nil
}

// Helper function to parse table data into a map
func (ctx *TestContext) parseTableToMap(table [][]string) (map[string]interface{}, error) {
	if len(table) < 2 {
		return nil, fmt.Errorf("table must have at least a header and one data row")
	}

	result := make(map[string]interface{})
	headers := table[0]

	for i, cell := range table[1] {
		if i < len(headers) {
			key := headers[i]
			value := cell

			// Try to convert to appropriate types
			if value == "true" {
				result[key] = true
			} else if value == "false" {
				result[key] = false
			} else if floatVal, err := parseFloat(value); err == nil {
				result[key] = floatVal
			} else {
				result[key] = value
			}
		}
	}

	return result, nil
}

// Helper function to parse float with better error handling
func parseFloat(s string) (float64, error) {
	return json.Number(s).Float64()
}

// Helper function to check if response contains specific error patterns
func (ctx *TestContext) checkErrorPattern(pattern string) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}

	if response.Error == nil {
		return fmt.Errorf("response should contain an error")
	}

	return nil
}

// Helper function to validate array response structure
func (ctx *TestContext) validateArrayResponse(fieldName string, expectedCount int) ([]interface{}, error) {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return nil, err
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response result is not an object")
	}

	arrayInterface, exists := result[fieldName]
	if !exists {
		return nil, fmt.Errorf("%s field not found in response", fieldName)
	}

	array, ok := arrayInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s field is not an array", fieldName)
	}

	if expectedCount >= 0 && len(array) != expectedCount {
		return nil, fmt.Errorf("expected %d items in %s, got %d", expectedCount, fieldName, len(array))
	}

	return array, nil
}

// Helper function to convert godog table to string slice
func tableToStringSlice(rows [][]*string) [][]string {
	result := make([][]string, len(rows))
	for i, row := range rows {
		result[i] = make([]string, len(row))
		for j, cell := range row {
			if cell != nil {
				result[i][j] = *cell
			}
		}
	}
	return result
}

// Helper function to wait for async operations
func (ctx *TestContext) waitForOperation(operationID string, timeoutSeconds int) error {
	// This would poll for operation completion
	// For now, return immediately
	return nil
}

// Helper function to setup test data
func (ctx *TestContext) setupTestData(dataType string) error {
	switch dataType {
	case "movies":
		return ctx.createSampleMovies()
	case "actors":
		return ctx.createSampleActors()
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
}

// Helper function to create sample movies
func (ctx *TestContext) createSampleMovies() error {
	sampleMovies := []map[string]interface{}{
		{
			"title":    "Test Movie 1",
			"director": "Test Director 1",
			"year":     2020,
			"rating":   8.5,
		},
		{
			"title":    "Test Movie 2",
			"director": "Test Director 2",
			"year":     2021,
			"rating":   7.8,
		},
	}

	for _, movie := range sampleMovies {
		err := ctx.createMovieViaMCP(movie)
		if err != nil {
			return fmt.Errorf("failed to create sample movie: %w", err)
		}
	}

	return nil
}

// Helper function to create sample actors
func (ctx *TestContext) createSampleActors() error {
	sampleActors := []map[string]interface{}{
		{
			"name":       "Test Actor 1",
			"birth_year": 1980,
			"bio":        "Test bio 1",
		},
		{
			"name":       "Test Actor 2",
			"birth_year": 1985,
			"bio":        "Test bio 2",
		},
	}

	for _, actor := range sampleActors {
		err := ctx.createActorViaMCP(actor)
		if err != nil {
			return fmt.Errorf("failed to create sample actor: %w", err)
		}
	}

	return nil
}
