package mcp

import (
	"testing"

	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

func createTestToolSchemas() []dto.Tool {
	return []dto.Tool{
		{
			Name:        "add_movie",
			Description: "Add a new movie",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"title": {
						Type:        "string",
						Description: "Movie title",
					},
					"director": {
						Type:        "string",
						Description: "Director name",
					},
					"year": {
						Type:        "integer",
						Description: "Release year",
					},
					"rating": {
						Type:        "number",
						Description: "Movie rating",
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
				Properties: map[string]dto.SchemaProperty{
					"query": {
						Type:        "string",
						Description: "Search query",
					},
					"limit": {
						Type:    "integer",
						Default: 10,
					},
				},
				Required: []string{"query"},
			},
		},
	}
}

func TestValidateToolCall(t *testing.T) {
	tools := createTestToolSchemas()
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		toolName  string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name:     "valid add_movie call",
			toolName: "add_movie",
			arguments: map[string]interface{}{
				"title":    "The Matrix",
				"director": "The Wachowskis",
				"year":     1999,
				"rating":   8.7,
			},
			wantValid: true,
		},
		{
			name:     "missing required field",
			toolName: "add_movie",
			arguments: map[string]interface{}{
				"title": "The Matrix",
				"year":  1999,
			},
			wantValid: false,
		},
		{
			name:     "invalid tool name",
			toolName: "nonexistent_tool",
			arguments: map[string]interface{}{
				"param": "value",
			},
			wantValid: false,
		},
		{
			name:     "valid search_movies call",
			toolName: "search_movies",
			arguments: map[string]interface{}{
				"query": "action movies",
				"limit": 5,
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall(tt.toolName, tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestNewToolValidator(t *testing.T) {
	tools := createTestToolSchemas()
	validator := NewToolValidator(tools)

	if validator == nil {
		t.Error("NewToolValidator returned nil")
	}

	// Test that schemas are stored correctly
	schemas := validator.GetSchemas()
	if len(schemas) != len(tools) {
		t.Errorf("Expected %d schemas, got %d", len(tools), len(schemas))
	}
}

func TestValidateString(t *testing.T) {
	minLen := 3
	maxLen := 10
	tools := []dto.Tool{
		{
			Name:        "test_string_validation",
			Description: "Test string validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"email": {
						Type:        "string",
						Format:      "email",
						Description: "Email address",
					},
					"url": {
						Type:        "string",
						Format:      "uri",
						Description: "Website URL",
					},
					"date": {
						Type:        "string",
						Format:      "date",
						Description: "Date in YYYY-MM-DD format",
					},
					"datetime": {
						Type:        "string",
						Format:      "date-time",
						Description: "ISO 8601 datetime",
					},
					"pattern_field": {
						Type:        "string",
						Pattern:     "^[A-Z]{2,4}$",
						Description: "Uppercase code",
					},
					"enum_field": {
						Type:        "string",
						Enum:        []interface{}{"action", "comedy", "drama"},
						Description: "Movie genre",
					},
					"length_field": {
						Type:        "string",
						MinLength:   &minLen,
						MaxLength:   &maxLen,
						Description: "Username",
					},
				},
				Required: []string{},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid email",
			arguments: map[string]interface{}{
				"email": "test@example.com",
			},
			wantValid: true,
		},
		{
			name: "invalid email",
			arguments: map[string]interface{}{
				"email": "not-an-email",
			},
			wantValid: false,
		},
		{
			name: "valid URL",
			arguments: map[string]interface{}{
				"url": "https://example.com",
			},
			wantValid: true,
		},
		{
			name: "invalid URL",
			arguments: map[string]interface{}{
				"url": "not a url",
			},
			wantValid: false,
		},
		{
			name: "valid date",
			arguments: map[string]interface{}{
				"date": "2023-12-25",
			},
			wantValid: true,
		},
		{
			name: "invalid date format",
			arguments: map[string]interface{}{
				"date": "25-12-2023",
			},
			wantValid: false,
		},
		{
			name: "valid datetime",
			arguments: map[string]interface{}{
				"datetime": "2023-12-25T10:30:00Z",
			},
			wantValid: true,
		},
		{
			name: "valid pattern match",
			arguments: map[string]interface{}{
				"pattern_field": "ABC",
			},
			wantValid: true,
		},
		{
			name: "invalid pattern match",
			arguments: map[string]interface{}{
				"pattern_field": "abc",
			},
			wantValid: false,
		},
		{
			name: "valid enum value",
			arguments: map[string]interface{}{
				"enum_field": "action",
			},
			wantValid: true,
		},
		{
			name: "invalid enum value",
			arguments: map[string]interface{}{
				"enum_field": "horror",
			},
			wantValid: false,
		},
		{
			name: "valid string length",
			arguments: map[string]interface{}{
				"length_field": "valid",
			},
			wantValid: true,
		},
		{
			name: "string too short",
			arguments: map[string]interface{}{
				"length_field": "ab",
			},
			wantValid: false,
		},
		{
			name: "string too long",
			arguments: map[string]interface{}{
				"length_field": "verylongusername",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_string_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
	min := float64(0)
	max := float64(150)
	tools := []dto.Tool{
		{
			Name:        "test_integer_validation",
			Description: "Test integer validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"age": {
						Type:        "integer",
						Minimum:     &min,
						Maximum:     &max,
						Description: "Person's age",
					},
				},
				Required: []string{},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid age",
			arguments: map[string]interface{}{
				"age": 25,
			},
			wantValid: true,
		},
		{
			name: "age too low",
			arguments: map[string]interface{}{
				"age": -5,
			},
			wantValid: false,
		},
		{
			name: "age too high",
			arguments: map[string]interface{}{
				"age": 200,
			},
			wantValid: false,
		},
		{
			name: "invalid type - float instead of integer",
			arguments: map[string]interface{}{
				"age": 25.5,
			},
			wantValid: false,
		},
		{
			name: "invalid type - string instead of integer",
			arguments: map[string]interface{}{
				"age": "25",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_integer_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestValidateNumber(t *testing.T) {
	min := float64(0)
	max := float64(1000)
	tools := []dto.Tool{
		{
			Name:        "test_number_validation",
			Description: "Test number validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"price": {
						Type:        "number",
						Minimum:     &min,
						Maximum:     &max,
						Description: "Product price",
					},
				},
				Required: []string{},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid price",
			arguments: map[string]interface{}{
				"price": 99.99,
			},
			wantValid: true,
		},
		{
			name: "price too low",
			arguments: map[string]interface{}{
				"price": -10.0,
			},
			wantValid: false,
		},
		{
			name: "price too high",
			arguments: map[string]interface{}{
				"price": 1500.0,
			},
			wantValid: false,
		},
		{
			name: "valid integer as number",
			arguments: map[string]interface{}{
				"price": 100,
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_number_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestValidateBoolean(t *testing.T) {
	tools := []dto.Tool{
		{
			Name:        "test_boolean_validation",
			Description: "Test boolean validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"is_active": {
						Type:        "boolean",
						Description: "Active status",
					},
				},
				Required: []string{"is_active"},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid boolean true",
			arguments: map[string]interface{}{
				"is_active": true,
			},
			wantValid: true,
		},
		{
			name: "valid boolean false",
			arguments: map[string]interface{}{
				"is_active": false,
			},
			wantValid: true,
		},
		{
			name: "invalid boolean - string",
			arguments: map[string]interface{}{
				"is_active": "true",
			},
			wantValid: false,
		},
		{
			name: "invalid boolean - number",
			arguments: map[string]interface{}{
				"is_active": 1,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_boolean_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestValidateArray(t *testing.T) {
	tools := []dto.Tool{
		{
			Name:        "test_array_validation",
			Description: "Test array validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"tags": {
						Type:        "array",
						Description: "List of tags",
						Items: &dto.SchemaProperty{
							Type: "string",
						},
					},
					"numbers": {
						Type:        "array",
						Description: "List of numbers",
						Items: &dto.SchemaProperty{
							Type: "integer",
						},
					},
				},
				Required: []string{},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid tags array",
			arguments: map[string]interface{}{
				"tags": []interface{}{"action", "thriller"},
			},
			wantValid: true,
		},
		{
			name: "empty tags array",
			arguments: map[string]interface{}{
				"tags": []interface{}{},
			},
			wantValid: true,
		},
		{
			name: "valid numbers array",
			arguments: map[string]interface{}{
				"numbers": []interface{}{1, 2, 3},
			},
			wantValid: true,
		},
		{
			name: "invalid numbers array - contains string",
			arguments: map[string]interface{}{
				"numbers": []interface{}{1, "two", 3},
			},
			wantValid: false,
		},
		{
			name: "invalid - not an array",
			arguments: map[string]interface{}{
				"tags": "not an array",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_array_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestValidateObject(t *testing.T) {
	tools := []dto.Tool{
		{
			Name:        "test_object_validation",
			Description: "Test object validation",
			InputSchema: dto.InputSchema{
				Type: "object",
				Properties: map[string]dto.SchemaProperty{
					"config": {
						Type:        "object",
						Description: "Configuration object",
						Properties: map[string]dto.SchemaProperty{
							"host": {
								Type:        "string",
								Description: "Host name",
							},
							"port": {
								Type:        "integer",
								Description: "Port number",
							},
						},
						Required: []string{"host"},
					},
				},
				Required: []string{},
			},
		},
	}
	validator := NewToolValidator(tools)

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid config object",
			arguments: map[string]interface{}{
				"config": map[string]interface{}{
					"host": "localhost",
					"port": 8080,
				},
			},
			wantValid: true,
		},
		{
			name: "missing required property in nested object",
			arguments: map[string]interface{}{
				"config": map[string]interface{}{
					"port": 8080,
				},
			},
			wantValid: false,
		},
		{
			name: "invalid property type in nested object",
			arguments: map[string]interface{}{
				"config": map[string]interface{}{
					"host": "localhost",
					"port": "8080", // Should be integer
				},
			},
			wantValid: false,
		},
		{
			name: "invalid - not an object",
			arguments: map[string]interface{}{
				"config": "not an object",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateToolCall("test_object_validation", tt.arguments)
			if result.Valid != tt.wantValid {
				t.Errorf("ValidateToolCall() result.Valid = %v, wantValid %v, errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestHandleValidateToolCall(t *testing.T) {
	tools := createTestToolSchemas()
	validator := NewToolValidator(tools)

	tests := []struct {
		name       string
		id         interface{}
		arguments  map[string]interface{}
		wantValid  bool
		wantError  bool
		errorCode  int
		errorMsg   string
	}{
		{
			name: "valid tool validation request",
			id:   1,
			arguments: map[string]interface{}{
				"tool_name": "add_movie",
				"tool_arguments": map[string]interface{}{
					"title":    "The Matrix",
					"director": "The Wachowskis",
					"year":     1999,
				},
			},
			wantValid: true,
			wantError: false,
		},
		{
			name: "missing tool_name",
			id:   2,
			arguments: map[string]interface{}{
				"tool_arguments": map[string]interface{}{
					"title": "The Matrix",
				},
			},
			wantValid: false,
			wantError: true,
			errorCode: dto.InvalidParams,
			errorMsg:  "Tool name is required",
		},
		{
			name: "missing arguments",
			id:   3,
			arguments: map[string]interface{}{
				"tool_name": "add_movie",
			},
			wantValid: false,
			wantError: true,
			errorCode: dto.InvalidParams,
			errorMsg:  "Tool arguments must be an object",
		},
		{
			name: "invalid tool validation request",
			id:   4,
			arguments: map[string]interface{}{
				"tool_name": "add_movie",
				"tool_arguments": map[string]interface{}{
					"title": "The Matrix",
					// Missing required fields
				},
			},
			wantValid: false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resultReceived interface{}
			var errorReceived bool
			var errorCode int
			var errorMsg string

			sendResult := func(id interface{}, result interface{}) {
				resultReceived = result
			}

			sendError := func(id interface{}, code int, msg string, data interface{}) {
				errorReceived = true
				errorCode = code
				errorMsg = msg
			}

			validator.HandleValidateToolCall(tt.id, tt.arguments, sendResult, sendError)

			if errorReceived != tt.wantError {
				t.Errorf("HandleValidateToolCall() error = %v, wantError %v", errorReceived, tt.wantError)
				return
			}

			if errorReceived {
				if errorCode != tt.errorCode {
					t.Errorf("HandleValidateToolCall() errorCode = %v, want %v", errorCode, tt.errorCode)
				}
				if errorMsg != tt.errorMsg {
					t.Errorf("HandleValidateToolCall() errorMsg = %v, want %v", errorMsg, tt.errorMsg)
				}
			} else {
				resultMap, ok := resultReceived.(map[string]interface{})
				if !ok {
					t.Errorf("HandleValidateToolCall() result is not a map")
					return
				}
				valid, ok := resultMap["valid"].(bool)
				if !ok {
					t.Errorf("HandleValidateToolCall() result.valid is not a bool")
					return
				}
				if valid != tt.wantValid {
					t.Errorf("HandleValidateToolCall() result.valid = %v, want %v", valid, tt.wantValid)
				}
			}
		})
	}
}

func TestIsValidDateFormat(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid date", "2023-12-25", true},
		{"invalid date format", "25-12-2023", false},
		{"invalid date - too long", "2023-12-250", false},
		{"invalid date - too short", "2023-12-2", false},
		{"invalid date - month out of range", "2023-13-25", false},
		{"invalid date - day out of range", "2023-12-32", false},
		{"invalid date - not numeric", "abcd-ef-gh", false},
		{"invalid date - year too small", "999-12-25", false},
		{"invalid date - year too large", "10000-12-25", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidDateFormat(tt.value); got != tt.want {
				t.Errorf("isValidDateFormat(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsValidURI(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid http URL", "http://example.com", true},
		{"valid https URL", "https://example.com/path", true},
		{"valid URL with query", "https://example.com/path?query=value", true},
		{"valid absolute path", "/path/to/resource", true},
		{"valid mailto", "mailto:test@example.com", true},
		{"invalid URL - no scheme or path", "example.com", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURI(tt.value); got != tt.want {
				t.Errorf("isValidURI(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestSchemaPropertyToMap(t *testing.T) {
	minLen := 5
	maxLen := 10
	min := float64(0)
	max := float64(100)
	
	tests := []struct {
		name     string
		prop     dto.SchemaProperty
		wantKeys []string
	}{
		{
			name: "string property with all constraints",
			prop: dto.SchemaProperty{
				Type:        "string",
				Description: "Test string",
				Format:      "email",
				Pattern:     "^[a-z]+$",
				MinLength:   &minLen,
				MaxLength:   &maxLen,
				Enum:        []interface{}{"a", "b", "c"},
			},
			wantKeys: []string{"type", "description", "format", "pattern", "minLength", "maxLength", "enum"},
		},
		{
			name: "number property with constraints",
			prop: dto.SchemaProperty{
				Type:        "number",
				Description: "Test number",
				Minimum:     &min,
				Maximum:     &max,
			},
			wantKeys: []string{"type", "description", "minimum", "maximum"},
		},
		{
			name: "array property with items",
			prop: dto.SchemaProperty{
				Type:        "array",
				Description: "Test array",
				Items: &dto.SchemaProperty{
					Type: "string",
				},
			},
			wantKeys: []string{"type", "description", "items"},
		},
		{
			name: "object property with nested properties",
			prop: dto.SchemaProperty{
				Type:        "object",
				Description: "Test object",
				Properties: map[string]dto.SchemaProperty{
					"nested": {
						Type: "string",
					},
				},
				Required: []string{"nested"},
			},
			wantKeys: []string{"type", "description", "properties", "required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schemaPropertyToMap(tt.prop)
			
			// Check that all expected keys are present
			for _, key := range tt.wantKeys {
				if _, ok := result[key]; !ok {
					t.Errorf("schemaPropertyToMap() missing expected key %q", key)
				}
			}
			
			// Check type is always present
			if _, ok := result["type"]; !ok {
				t.Errorf("schemaPropertyToMap() missing required 'type' field")
			}
		})
	}
}