package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"movies-mcp-server/internal/models"
)

// Integration tests for the full MCP server
// These tests start an actual server process and communicate via JSON-RPC

func TestMCPServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Start the server
	cmd := exec.Command("go", "run", "../../cmd/server/main.go")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	scanner := bufio.NewScanner(stdout)
	requestID := 1

	// Helper function to send request and read response
	sendRequest := func(method string, params interface{}) models.JSONRPCResponse {
		req := models.JSONRPCRequest{
			JSONRPC: "2.0",
			Method:  method,
			ID:      requestID,
		}
		
		if params != nil {
			paramsBytes, _ := json.Marshal(params)
			req.Params = paramsBytes
		}

		reqBytes, _ := json.Marshal(req)
		fmt.Fprintf(stdin, "%s\n", reqBytes)
		requestID++

		// Read response
		if scanner.Scan() {
			var response models.JSONRPCResponse
			json.Unmarshal([]byte(scanner.Text()), &response)
			return response
		}

		t.Fatalf("Failed to read response")
		return models.JSONRPCResponse{}
	}

	// Give server time to start
	time.Sleep(2 * time.Second)

	// Test 1: Initialize the server
	t.Run("Initialize", func(t *testing.T) {
		initParams := map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		}

		response := sendRequest("initialize", initParams)

		if response.Error != nil {
			t.Errorf("Initialize failed: %v", response.Error)
		}

		// Parse initialization response
		var initResp models.InitializeResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &initResp)

		if initResp.ProtocolVersion != "2024-11-05" {
			t.Errorf("Expected protocol version 2024-11-05, got %s", initResp.ProtocolVersion)
		}

		if initResp.ServerInfo.Name != "movies-mcp-server" {
			t.Errorf("Expected server name 'movies-mcp-server', got %s", initResp.ServerInfo.Name)
		}

		// Check capabilities
		if initResp.Capabilities.Tools == nil {
			t.Errorf("Server should support tools capability")
		}
		if initResp.Capabilities.Resources == nil {
			t.Errorf("Server should support resources capability")
		}
	})

	// Test 2: List available tools
	t.Run("ToolsList", func(t *testing.T) {
		response := sendRequest("tools/list", nil)

		if response.Error != nil {
			t.Errorf("Tools list failed: %v", response.Error)
		}

		var toolsResp models.ToolsListResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &toolsResp)

		expectedTools := []string{
			"get_movie",
			"add_movie", 
			"update_movie",
			"delete_movie",
			"search_movies",
			"list_top_movies",
		}

		toolNames := make(map[string]bool)
		for _, tool := range toolsResp.Tools {
			toolNames[tool.Name] = true
		}

		for _, expectedTool := range expectedTools {
			if !toolNames[expectedTool] {
				t.Errorf("Expected tool %s not found", expectedTool)
			}
		}
	})

	// Test 3: List available resources
	t.Run("ResourcesList", func(t *testing.T) {
		response := sendRequest("resources/list", nil)

		if response.Error != nil {
			t.Errorf("Resources list failed: %v", response.Error)
		}

		var resourcesResp models.ResourcesListResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &resourcesResp)

		expectedResources := []string{
			"movies://database/all",
			"movies://database/stats",
		}

		resourceURIs := make(map[string]bool)
		for _, resource := range resourcesResp.Resources {
			resourceURIs[resource.URI] = true
		}

		for _, expectedResource := range expectedResources {
			if !resourceURIs[expectedResource] {
				t.Errorf("Expected resource %s not found", expectedResource)
			}
		}
	})

	// Test 4: Add a movie
	t.Run("AddMovie", func(t *testing.T) {
		addParams := models.ToolCallRequest{
			Name: "add_movie",
			Arguments: map[string]interface{}{
				"title":       "Integration Test Movie",
				"director":    "Test Director",
				"year":        2023,
				"genre":       []string{"Test", "Integration"},
				"rating":      8.5,
				"description": "A movie created during integration testing",
			},
		}

		response := sendRequest("tools/call", addParams)

		if response.Error != nil {
			t.Errorf("Add movie failed: %v", response.Error)
		}

		var toolResp models.ToolCallResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &toolResp)

		if len(toolResp.Content) == 0 || !strings.Contains(toolResp.Content[0].Text, "Successfully created") {
			t.Errorf("Expected success message, got: %v", toolResp.Content)
		}
	})

	// Test 5: Search for the added movie
	t.Run("SearchMovie", func(t *testing.T) {
		searchParams := models.ToolCallRequest{
			Name: "search_movies",
			Arguments: map[string]interface{}{
				"query": "Integration Test Movie",
				"type":  "title",
			},
		}

		response := sendRequest("tools/call", searchParams)

		if response.Error != nil {
			t.Errorf("Search movie failed: %v", response.Error)
		}

		var toolResp models.ToolCallResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &toolResp)

		if len(toolResp.Content) == 0 || !strings.Contains(toolResp.Content[0].Text, "Integration Test Movie") {
			t.Errorf("Expected to find the test movie, got: %v", toolResp.Content)
		}
	})

	// Test 6: Read database stats resource
	t.Run("ReadDatabaseStats", func(t *testing.T) {
		readParams := models.ResourceReadRequest{
			URI: "movies://database/stats",
		}

		response := sendRequest("resources/read", readParams)

		if response.Error != nil {
			t.Errorf("Read database stats failed: %v", response.Error)
		}

		var resourceResp models.ResourceReadResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &resourceResp)

		if len(resourceResp.Contents) == 0 {
			t.Errorf("Expected stats content")
		}

		content := resourceResp.Contents[0]
		if content.MimeType != "application/json" {
			t.Errorf("Expected JSON mime type, got %s", content.MimeType)
		}

		// Verify stats content is valid JSON
		var stats map[string]interface{}
		if err := json.Unmarshal([]byte(content.Text), &stats); err != nil {
			t.Errorf("Stats content should be valid JSON: %v", err)
		}

		// Check for expected stats fields
		expectedFields := []string{"total_movies", "average_rating"}
		for _, field := range expectedFields {
			if _, exists := stats[field]; !exists {
				t.Errorf("Expected stats field %s not found", field)
			}
		}
	})

	// Test 7: Read database/all resource
	t.Run("ReadDatabaseAll", func(t *testing.T) {
		readParams := models.ResourceReadRequest{
			URI: "movies://database/all",
		}

		response := sendRequest("resources/read", readParams)

		if response.Error != nil {
			t.Errorf("Read database all failed: %v", response.Error)
		}

		var resourceResp models.ResourceReadResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &resourceResp)

		if len(resourceResp.Contents) == 0 {
			t.Errorf("Expected movies content")
		}

		content := resourceResp.Contents[0]
		if content.MimeType != "application/json" {
			t.Errorf("Expected JSON mime type, got %s", content.MimeType)
		}

		// Verify content is valid JSON array
		var movies []map[string]interface{}
		if err := json.Unmarshal([]byte(content.Text), &movies); err != nil {
			t.Errorf("Movies content should be valid JSON array: %v", err)
		}

		// Should include our test movie
		foundTestMovie := false
		for _, movie := range movies {
			if title, exists := movie["title"]; exists && title == "Integration Test Movie" {
				foundTestMovie = true
				break
			}
		}

		if !foundTestMovie {
			t.Errorf("Test movie not found in database/all resource")
		}
	})

	// Test 8: List top movies
	t.Run("ListTopMovies", func(t *testing.T) {
		listParams := models.ToolCallRequest{
			Name: "list_top_movies",
			Arguments: map[string]interface{}{
				"limit": 5,
			},
		}

		response := sendRequest("tools/call", listParams)

		if response.Error != nil {
			t.Errorf("List top movies failed: %v", response.Error)
		}

		var toolResp models.ToolCallResponse
		respBytes, _ := json.Marshal(response.Result)
		json.Unmarshal(respBytes, &toolResp)

		if len(toolResp.Content) == 0 || !strings.Contains(toolResp.Content[0].Text, "Top") {
			t.Errorf("Expected top movies list, got: %v", toolResp.Content)
		}
	})

	// Test 9: Error handling - invalid resource URI
	t.Run("InvalidResourceURI", func(t *testing.T) {
		readParams := models.ResourceReadRequest{
			URI: "invalid://resource/uri",
		}

		response := sendRequest("resources/read", readParams)

		if response.Error == nil {
			t.Errorf("Expected error for invalid resource URI")
		}
	})

	// Test 10: Error handling - invalid tool call
	t.Run("InvalidToolCall", func(t *testing.T) {
		toolParams := models.ToolCallRequest{
			Name: "nonexistent_tool",
			Arguments: map[string]interface{}{
				"param": "value",
			},
		}

		response := sendRequest("tools/call", toolParams)

		if response.Error == nil {
			t.Errorf("Expected error for invalid tool name")
		}
	})
}

// Test server performance under load
func TestMCPServerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test")
	}

	// This is a placeholder for performance testing
	// In a real implementation, this would:
	// 1. Start the server
	// 2. Send many concurrent requests
	// 3. Measure response times
	// 4. Verify the server handles load gracefully

	t.Run("ConcurrentRequests", func(t *testing.T) {
		// Test concurrent search requests
		// Test concurrent resource reads
		// Measure average response time
		// Verify no errors under normal load
		t.Skip("Performance test placeholder - implement when needed")
	})

	t.Run("LargeDatasets", func(t *testing.T) {
		// Test with databases containing thousands of movies
		// Test pagination performance
		// Test search performance with large datasets
		t.Skip("Large dataset test placeholder - implement when needed")
	})
}

// Test server error recovery
func TestMCPServerErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping error recovery test")
	}

	t.Run("DatabaseConnectionLoss", func(t *testing.T) {
		// Test server behavior when database connection is lost
		// Test graceful error handling
		// Test recovery when connection is restored
		t.Skip("Database error recovery test placeholder - implement when needed")
	})

	t.Run("InvalidJSONRequests", func(t *testing.T) {
		// Test server handling of malformed JSON
		// Test server handling of invalid method names
		// Test server handling of missing required parameters
		t.Skip("Invalid JSON handling test placeholder - implement when needed")
	})
}

// Test MCP protocol compliance
func TestMCPProtocolCompliance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping protocol compliance test")
	}

	t.Run("JSONRPCCompliance", func(t *testing.T) {
		// Test JSON-RPC 2.0 compliance
		// Test proper error code usage
		// Test proper response formatting
		t.Skip("JSON-RPC compliance test placeholder - implement when needed")
	})

	t.Run("MCPSpecCompliance", func(t *testing.T) {
		// Test MCP specification compliance
		// Test proper capability negotiation
		// Test proper resource URI formatting
		t.Skip("MCP spec compliance test placeholder - implement when needed")
	})
}