package server

import (
	"encoding/json"
	"testing"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

func TestResourceSchemasIntegration(t *testing.T) {
	// Test that resource schemas are properly formatted JSON Schema
	resources := []dto.Resource{
		{
			URI:         "movies://database/all",
			Name:        "All Movies",
			Description: "Complete list of all movies in the database",
			MimeType:    "application/json",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"movies": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{
									"type": "integer",
									"description": "Unique movie identifier",
								},
								"title": map[string]interface{}{
									"type": "string",
									"description": "Movie title",
								},
							},
							"required": []string{"id", "title"},
						},
					},
				},
				"required": []string{"movies"},
			},
		},
	}

	// Test serialization
	for _, resource := range resources {
		t.Run("resource_"+resource.Name, func(t *testing.T) {
			// Should be able to marshal to JSON
			data, err := json.Marshal(resource)
			if err != nil {
				t.Fatalf("Failed to marshal resource: %v", err)
			}

			// Should be able to unmarshal back
			var unmarshaled dto.Resource
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal resource: %v", err)
			}

			// Check essential fields
			if unmarshaled.URI != resource.URI {
				t.Errorf("URI mismatch: got %s, want %s", unmarshaled.URI, resource.URI)
			}

			if unmarshaled.Name != resource.Name {
				t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, resource.Name)
			}

			if unmarshaled.MimeType != resource.MimeType {
				t.Errorf("MimeType mismatch: got %s, want %s", unmarshaled.MimeType, resource.MimeType)
			}

			// Schema should exist and be non-empty
			if unmarshaled.Schema == nil {
				t.Error("Schema is nil")
			}

			if len(unmarshaled.Schema) == 0 {
				t.Error("Schema is empty")
			}

			// Schema should have correct structure
			schemaType, ok := unmarshaled.Schema["type"].(string)
			if !ok || schemaType != "object" {
				t.Errorf("Schema type is not 'object': %v", schemaType)
			}

			properties, ok := unmarshaled.Schema["properties"].(map[string]interface{})
			if !ok {
				t.Error("Schema properties is not a map")
			}

			if len(properties) == 0 {
				t.Error("Schema properties is empty")
			}
		})
	}
}

func TestResourceSchemaValidation(t *testing.T) {
	// Test schema structure follows JSON Schema specification
	testCases := []struct {
		name   string
		schema map[string]interface{}
		valid  bool
	}{
		{
			name: "valid_object_schema",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type": "integer",
					},
				},
				"required": []string{"id"},
			},
			valid: true,
		},
		{
			name: "valid_array_schema",
			schema: map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			valid: true,
		},
		{
			name: "invalid_schema_no_type",
			schema: map[string]interface{}{
				"properties": map[string]interface{}{},
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := dto.Resource{
				URI:      "test://resource",
				Name:     "Test Resource",
				MimeType: "application/json",
				Schema:   tc.schema,
			}

			// Should always be able to marshal valid JSON
			data, err := json.Marshal(resource)
			if err != nil {
				t.Fatalf("Failed to marshal resource: %v", err)
			}

			// Should always be able to unmarshal back
			var unmarshaled dto.Resource
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal resource: %v", err)
			}

			// For valid schemas, check type exists
			if tc.valid {
				schemaType, ok := unmarshaled.Schema["type"].(string)
				if !ok || schemaType == "" {
					t.Error("Valid schema should have a type")
				}
			}
		})
	}
}