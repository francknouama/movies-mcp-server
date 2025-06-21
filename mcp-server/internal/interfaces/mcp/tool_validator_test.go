package mcp

import (
	"testing"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

func createTestToolSchemas() []dto.Tool {
	return []dto.Tool{
		{
			Name:        "add_movie",
			Description: "Add a new movie",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Movie title",
						"minLength":   1.0,
						"maxLength":   200.0,
					},
					"director": map[string]interface{}{
						"type":        "string",
						"description": "Director name",
					},
					"year": map[string]interface{}{
						"type":        "integer",
						"description": "Release year",
						"minimum":     1888.0,
						"maximum":     2100.0,
					},
					"rating": map[string]interface{}{
						"type":        "number",
						"description": "Movie rating",
						"minimum":     0.0,
						"maximum":     10.0,
					},
					"genres": map[string]interface{}{
						"type":        "array",
						"description": "Movie genres",
						"items": map[string]interface{}{
							"type": "string",
						},
						"minItems": 1.0,
						"maxItems": 5.0,
					},
					"country": map[string]interface{}{
						"type":        "string",
						"description": "Country of origin",
						"enum":        []string{"US", "UK", "FR", "DE", "JP", "Other"},
					},
					"release_date": map[string]interface{}{
						"type":        "string",
						"description": "Release date",
						"format":      "date",
					},
					"poster_url": map[string]interface{}{
						"type":        "string",
						"description": "Poster URL",
						"format":      "uri",
					},
				},
				Required: []string{"title", "director", "year"},
			},
		},
		{
			Name:        "search_movies",
			Description: "Search for movies",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"include_adult": map[string]interface{}{
						"type":        "boolean",
						"description": "Include adult content",
					},
				},
				Required: []string{"query"},
			},
		},
	}
}

func TestToolValidator_ValidToolCall(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	// Valid movie addition
	result := validator.ValidateToolCall("add_movie", map[string]interface{}{
		"title":    "The Matrix",
		"director": "The Wachowskis",
		"year":     1999.0,
		"rating":   8.7,
		"genres":   []interface{}{"Action", "Sci-Fi"},
		"country":  "US",
	})

	if !result.Valid {
		t.Errorf("Expected valid result, got errors: %v", result.Errors)
	}

	if len(result.Errors) != 0 {
		t.Errorf("Expected no errors, got %d errors", len(result.Errors))
	}
}

func TestToolValidator_UnknownTool(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	result := validator.ValidateToolCall("unknown_tool", map[string]interface{}{})

	if result.Valid {
		t.Error("Expected invalid result for unknown tool")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if result.Errors[0].Code != "UNKNOWN_TOOL" {
		t.Errorf("Expected UNKNOWN_TOOL error code, got %s", result.Errors[0].Code)
	}
}

func TestToolValidator_MissingRequiredFields(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	// Missing required fields: title, director, year
	result := validator.ValidateToolCall("add_movie", map[string]interface{}{
		"rating": 8.5,
	})

	if result.Valid {
		t.Error("Expected invalid result for missing required fields")
	}

	expectedMissing := []string{"title", "director", "year"}
	if len(result.Errors) != len(expectedMissing) {
		t.Errorf("Expected %d errors for missing fields, got %d", len(expectedMissing), len(result.Errors))
	}

	for _, err := range result.Errors {
		if err.Code != "REQUIRED_FIELD_MISSING" {
			t.Errorf("Expected REQUIRED_FIELD_MISSING error code, got %s", err.Code)
		}
	}
}

func TestToolValidator_TypeMismatch(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	testCases := []struct {
		name          string
		arguments     map[string]interface{}
		expectedField string
		expectedCode  string
	}{
		{
			name: "string_as_integer",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     "1999", // Should be integer
			},
			expectedField: "year",
			expectedCode:  "TYPE_MISMATCH",
		},
		{
			name: "integer_as_string",
			arguments: map[string]interface{}{
				"title":    123, // Should be string
				"director": "Test Director",
				"year":     1999,
			},
			expectedField: "title",
			expectedCode:  "TYPE_MISMATCH",
		},
		{
			name: "string_as_boolean",
			arguments: map[string]interface{}{
				"query":         "test",
				"include_adult": "true", // Should be boolean
			},
			expectedField: "include_adult",
			expectedCode:  "TYPE_MISMATCH",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var toolName string
			if tc.name == "string_as_boolean" {
				toolName = "search_movies"
			} else {
				toolName = "add_movie"
			}

			result := validator.ValidateToolCall(toolName, tc.arguments)

			if result.Valid {
				t.Error("Expected invalid result for type mismatch")
			}

			foundError := false
			for _, err := range result.Errors {
				if err.Field == tc.expectedField && err.Code == tc.expectedCode {
					foundError = true
					break
				}
			}

			if !foundError {
				t.Errorf("Expected error with field %s and code %s, got errors: %v",
					tc.expectedField, tc.expectedCode, result.Errors)
			}
		})
	}
}

func TestToolValidator_RangeValidation(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	testCases := []struct {
		name         string
		arguments    map[string]interface{}
		expectedCode string
	}{
		{
			name: "year_too_small",
			arguments: map[string]interface{}{
				"title":    "Ancient Movie",
				"director": "Ancient Director",
				"year":     1800.0, // Below minimum of 1888
			},
			expectedCode: "VALUE_TOO_SMALL",
		},
		{
			name: "year_too_large",
			arguments: map[string]interface{}{
				"title":    "Future Movie",
				"director": "Future Director",
				"year":     2200.0, // Above maximum of 2100
			},
			expectedCode: "VALUE_TOO_LARGE",
		},
		{
			name: "rating_too_small",
			arguments: map[string]interface{}{
				"title":    "Bad Movie",
				"director": "Bad Director",
				"year":     2000.0,
				"rating":   -1.0, // Below minimum of 0
			},
			expectedCode: "VALUE_TOO_SMALL",
		},
		{
			name: "rating_too_large",
			arguments: map[string]interface{}{
				"title":    "Perfect Movie",
				"director": "Perfect Director",
				"year":     2000.0,
				"rating":   11.0, // Above maximum of 10
			},
			expectedCode: "VALUE_TOO_LARGE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateToolCall("add_movie", tc.arguments)

			if result.Valid {
				t.Error("Expected invalid result for range violation")
			}

			foundError := false
			for _, err := range result.Errors {
				if err.Code == tc.expectedCode {
					foundError = true
					break
				}
			}

			if !foundError {
				t.Errorf("Expected error with code %s, got errors: %v", tc.expectedCode, result.Errors)
			}
		})
	}
}

func TestToolValidator_StringValidation(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	testCases := []struct {
		name         string
		arguments    map[string]interface{}
		expectedCode string
	}{
		{
			name: "string_too_short",
			arguments: map[string]interface{}{
				"title":    "", // Below minLength of 1
				"director": "Test Director",
				"year":     2000.0,
			},
			expectedCode: "STRING_TOO_SHORT",
		},
		{
			name: "string_too_long",
			arguments: map[string]interface{}{
				"title":    string(make([]byte, 201)), // Above maxLength of 200
				"director": "Test Director",
				"year":     2000.0,
			},
			expectedCode: "STRING_TOO_LONG",
		},
		{
			name: "invalid_enum_value",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2000.0,
				"country":  "INVALID", // Not in enum values
			},
			expectedCode: "INVALID_ENUM_VALUE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateToolCall("add_movie", tc.arguments)

			if result.Valid {
				t.Error("Expected invalid result for string validation error")
			}

			foundError := false
			for _, err := range result.Errors {
				if err.Code == tc.expectedCode {
					foundError = true
					break
				}
			}

			if !foundError {
				t.Errorf("Expected error with code %s, got errors: %v", tc.expectedCode, result.Errors)
			}
		})
	}
}

func TestToolValidator_ArrayValidation(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	testCases := []struct {
		name         string
		arguments    map[string]interface{}
		expectedCode string
	}{
		{
			name: "array_too_short",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2000.0,
				"genres":   []interface{}{}, // Below minItems of 1
			},
			expectedCode: "ARRAY_TOO_SHORT",
		},
		{
			name: "array_too_long",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2000.0,
				"genres":   []interface{}{"Action", "Drama", "Comedy", "Thriller", "Romance", "Horror"}, // Above maxItems of 5
			},
			expectedCode: "ARRAY_TOO_LONG",
		},
		{
			name: "invalid_array_item_type",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2000.0,
				"genres":   []interface{}{"Action", 123}, // Second item should be string
			},
			expectedCode: "TYPE_MISMATCH",
		},
		{
			name: "non_array_value",
			arguments: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2000.0,
				"genres":   "Action", // Should be array
			},
			expectedCode: "TYPE_MISMATCH",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateToolCall("add_movie", tc.arguments)

			if result.Valid {
				t.Error("Expected invalid result for array validation error")
			}

			foundError := false
			for _, err := range result.Errors {
				if err.Code == tc.expectedCode {
					foundError = true
					break
				}
			}

			if !foundError {
				t.Errorf("Expected error with code %s, got errors: %v", tc.expectedCode, result.Errors)
			}
		})
	}
}

func TestToolValidator_FormatValidation(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	testCases := []struct {
		name         string
		arguments    map[string]interface{}
		expectedCode string
	}{
		{
			name: "invalid_date_format",
			arguments: map[string]interface{}{
				"title":        "Test Movie",
				"director":     "Test Director",
				"year":         2000.0,
				"release_date": "2023/12/25", // Should be YYYY-MM-DD format
			},
			expectedCode: "INVALID_DATE_FORMAT",
		},
		{
			name: "invalid_uri_format",
			arguments: map[string]interface{}{
				"title":      "Test Movie",
				"director":   "Test Director",
				"year":       2000.0,
				"poster_url": "not-a-valid-uri", // Should be valid URI
			},
			expectedCode: "INVALID_URI_FORMAT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateToolCall("add_movie", tc.arguments)

			if result.Valid {
				t.Error("Expected invalid result for format validation error")
			}

			foundError := false
			for _, err := range result.Errors {
				if err.Code == tc.expectedCode {
					foundError = true
					break
				}
			}

			if !foundError {
				t.Errorf("Expected error with code %s, got errors: %v", tc.expectedCode, result.Errors)
			}
		})
	}
}

func TestToolValidator_UnknownField(t *testing.T) {
	validator := NewToolValidator(createTestToolSchemas())

	result := validator.ValidateToolCall("add_movie", map[string]interface{}{
		"title":         "Test Movie",
		"director":      "Test Director",
		"year":          2000.0,
		"unknown_field": "some value", // Not in schema
	})

	if result.Valid {
		t.Error("Expected invalid result for unknown field")
	}

	foundError := false
	for _, err := range result.Errors {
		if err.Code == "UNKNOWN_FIELD" && err.Field == "unknown_field" {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected UNKNOWN_FIELD error for 'unknown_field', got errors: %v", result.Errors)
	}
}

func TestToolValidator_ValidDateFormats(t *testing.T) {
	validDates := []string{
		"2023-12-25",
		"1999-01-01",
		"2000-02-29", // Leap year
	}

	for _, date := range validDates {
		if !isValidDateFormat(date) {
			t.Errorf("Expected %s to be valid date format", date)
		}
	}
}

func TestToolValidator_InvalidDateFormats(t *testing.T) {
	invalidDates := []string{
		"2023/12/25",
		"23-12-25",
		"2023-13-01", // Invalid month
		"2023-12-32", // Invalid day
		"not-a-date",
		"",
	}

	for _, date := range invalidDates {
		if isValidDateFormat(date) {
			t.Errorf("Expected %s to be invalid date format", date)
		}
	}
}

func TestToolValidator_ValidURIFormats(t *testing.T) {
	validURIs := []string{
		"https://example.com",
		"http://example.com/path",
		"/relative/path",
		"mailto:test@example.com",
		"ftp://files.example.com",
	}

	for _, uri := range validURIs {
		if !isValidURI(uri) {
			t.Errorf("Expected %s to be valid URI format", uri)
		}
	}
}

func TestToolValidator_InvalidURIFormats(t *testing.T) {
	invalidURIs := []string{
		"",
		"not-a-uri",
		"just-text",
	}

	for _, uri := range invalidURIs {
		if isValidURI(uri) {
			t.Errorf("Expected %s to be invalid URI format", uri)
		}
	}
}
