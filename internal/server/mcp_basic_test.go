package server

import (
	"encoding/json"
	"testing"

	"movies-mcp-server/internal/interfaces/mcp"
	"movies-mcp-server/internal/interfaces/dto"
	"movies-mcp-server/internal/schemas"
)

func TestMCPBasicCapabilities(t *testing.T) {
	// Test tool toolSchemas are properly defined
	toolSchemas := schemas.GetToolSchemas()
	
	if len(toolSchemas) == 0 {
		t.Fatal("Expected at least one tool schema")
	}
	
	// Verify each schema has required fields
	for _, schema := range toolSchemas {
		if schema.Name == "" {
			t.Error("Tool schema missing name")
		}
		
		if schema.Description == "" {
			t.Error("Tool schema missing description")
		}
		
		if schema.InputSchema.Type != "object" {
			t.Errorf("Tool %s input schema type should be 'object', got %s", schema.Name, schema.InputSchema.Type)
		}
		
		if len(schema.InputSchema.Properties) == 0 {
			t.Errorf("Tool %s has no properties", schema.Name)
		}
	}
	
	// Test that key tools exist
	expectedTools := []string{
		"get_movie",
		"add_movie", 
		"search_movies",
		"validate_tool_call",
	}
	
	toolMap := make(map[string]bool)
	for _, schema := range toolSchemas {
		toolMap[schema.Name] = true
	}
	
	for _, expectedTool := range expectedTools {
		if !toolMap[expectedTool] {
			t.Errorf("Expected tool %s not found", expectedTool)
		}
	}
}

func TestMCPServerCapabilities(t *testing.T) {
	// Test server capabilities structure
	capabilities := dto.ServerCapabilities{
		Tools: &dto.ToolsCapability{},
		Resources: &dto.ResourcesCapability{
			Subscribe: false,
		},
		Prompts: &dto.PromptsCapability{},
	}
	
	// Should be serializable
	data, err := json.Marshal(capabilities)
	if err != nil {
		t.Fatalf("Failed to marshal server capabilities: %v", err)
	}
	
	// Should be deserializable
	var unmarshaled dto.ServerCapabilities
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal server capabilities: %v", err)
	}
	
	// Check structure
	if unmarshaled.Tools == nil {
		t.Error("Expected tools capability")
	}
	
	if unmarshaled.Resources == nil {
		t.Error("Expected resources capability")
	}
	
	if unmarshaled.Prompts == nil {
		t.Error("Expected prompts capability")
	}
}

func TestMCPRequestResponseStructure(t *testing.T) {
	// Test JSON-RPC request structure
	request := dto.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  json.RawMessage(`{}`),
	}
	
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	
	var unmarshaled dto.JSONRPCRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}
	
	if unmarshaled.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC 2.0, got %s", unmarshaled.JSONRPC)
	}
	
	if unmarshaled.Method != "tools/list" {
		t.Errorf("Expected method tools/list, got %s", unmarshaled.Method)
	}
	
	// Test JSON-RPC response structure
	response := dto.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]interface{}{"test": "result"},
	}
	
	data, err = json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	
	var unmarshaledResponse dto.JSONRPCResponse
	if err := json.Unmarshal(data, &unmarshaledResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if unmarshaledResponse.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC 2.0, got %s", unmarshaledResponse.JSONRPC)
	}
}

func TestMCPToolValidatorIntegration(t *testing.T) {
	// Test that tool validator works with actual toolSchemas
	toolSchemas := schemas.GetToolSchemas()
	validator := mcp.NewToolValidator(toolSchemas)
	
	// Test valid tool call
	result := validator.ValidateToolCall("get_movie", map[string]interface{}{
		"movie_id": 1.0,
	})
	
	if !result.Valid {
		t.Errorf("Expected valid result, got errors: %v", result.Errors)
	}
	
	// Test invalid tool call
	result = validator.ValidateToolCall("get_movie", map[string]interface{}{
		"movie_id": "not_a_number",
	})
	
	if result.Valid {
		t.Error("Expected invalid result for wrong type")
	}
	
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
	
	// Test unknown tool
	result = validator.ValidateToolCall("unknown_tool", map[string]interface{}{})
	
	if result.Valid {
		t.Error("Expected invalid result for unknown tool")
	}
	
	foundUnknownToolError := false
	for _, err := range result.Errors {
		if err.Code == "UNKNOWN_TOOL" {
			foundUnknownToolError = true
			break
		}
	}
	
	if !foundUnknownToolError {
		t.Error("Expected UNKNOWN_TOOL error")
	}
}

func TestMCPResourceSchemas(t *testing.T) {
	// Test that resources have proper toolSchemas
	// This tests the schema structure without requiring a full server setup
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
						},
					},
				},
			},
		},
	}
	
	// Verify resource structure
	for _, resource := range resources {
		if resource.URI == "" {
			t.Error("Resource missing URI")
		}
		
		if resource.Name == "" {
			t.Error("Resource missing name")
		}
		
		if resource.MimeType == "" {
			t.Error("Resource missing MIME type")
		}
		
		if resource.Schema == nil {
			t.Error("Resource missing schema")
		} else {
			if schemaType, ok := resource.Schema["type"].(string); !ok || schemaType == "" {
				t.Error("Resource schema missing type")
			}
		}
	}
}

func TestMCPPromptTemplates(t *testing.T) {
	// Test prompt structure
	prompts := []dto.Prompt{
		{
			Name:        "movie_recommendation",
			Description: "Generate movie recommendations",
			Arguments: []dto.PromptArgument{
				{
					Name:        "genre",
					Description: "Preferred genre",
					Required:    true,
				},
			},
		},
	}
	
	// Verify prompt structure
	for _, prompt := range prompts {
		if prompt.Name == "" {
			t.Error("Prompt missing name")
		}
		
		if prompt.Description == "" {
			t.Error("Prompt missing description")
		}
		
		// Check arguments structure
		for _, arg := range prompt.Arguments {
			if arg.Name == "" {
				t.Error("Prompt argument missing name")
			}
		}
	}
	
	// Test serialization
	data, err := json.Marshal(prompts)
	if err != nil {
		t.Fatalf("Failed to marshal prompts: %v", err)
	}
	
	var unmarshaled []dto.Prompt
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal prompts: %v", err)
	}
	
	if len(unmarshaled) != len(prompts) {
		t.Errorf("Prompt count mismatch after serialization")
	}
}

func TestMCPErrorCodeConstants(t *testing.T) {
	// Test that error codes are properly defined
	expectedCodes := map[string]int{
		"ParseError":     dto.ParseError,
		"InvalidRequest": dto.InvalidRequest,
		"MethodNotFound": dto.MethodNotFound,
		"InvalidParams":  dto.InvalidParams,
		"InternalError":  dto.InternalError,
	}
	
	for name, code := range expectedCodes {
		if code == 0 {
			t.Errorf("Error code %s is zero", name)
		}
	}
	
	// Test error creation
	err := dto.NewJSONRPCError(dto.InvalidParams, "Test error", "test data")
	if err.Code != dto.InvalidParams {
		t.Errorf("Expected code %d, got %d", dto.InvalidParams, err.Code)
	}
	
	if err.Message != "Test error" {
		t.Errorf("Expected message 'Test error', got %s", err.Message)
	}
}