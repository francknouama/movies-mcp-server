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
