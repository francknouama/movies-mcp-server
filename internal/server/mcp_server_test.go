package server

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"movies-mcp-server/internal/composition"
	"movies-mcp-server/internal/interfaces/dto"
)

func TestNewMCPServer(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.input != input {
		t.Error("Expected input to be set correctly")
	}

	if server.output != output {
		t.Error("Expected output to be set correctly")
	}

	if server.logger != logger {
		t.Error("Expected logger to be set correctly")
	}

	if server.container != container {
		t.Error("Expected container to be set correctly")
	}

	if server.validator == nil {
		t.Error("Expected validator to be initialized")
	}

	if server.toolHandlers == nil {
		t.Error("Expected tool handlers to be initialized")
	}
}

func TestMCPServer_HandleInitialize(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	request := &dto.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      1,
		Params:  json.RawMessage(`{}`),
	}

	server.handleInitialize(request)

	// Parse the response
	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", response.JSONRPC)
	}

	if response.ID != float64(1) {
		t.Errorf("Expected ID 1, got %v", response.ID)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// Check the initialize response
	resultBytes, _ := json.Marshal(response.Result)
	var initResponse dto.InitializeResponse
	json.Unmarshal(resultBytes, &initResponse)

	if initResponse.ProtocolVersion != "2024-11-05" {
		t.Errorf("Expected protocol version 2024-11-05, got %s", initResponse.ProtocolVersion)
	}

	if initResponse.ServerInfo.Name != "movies-mcp-server" {
		t.Errorf("Expected server name movies-mcp-server, got %s", initResponse.ServerInfo.Name)
	}

	if initResponse.Capabilities.Tools == nil {
		t.Error("Expected tools capability to be present")
	}

	if initResponse.Capabilities.Resources == nil {
		t.Error("Expected resources capability to be present")
	}

	if initResponse.Capabilities.Prompts == nil {
		t.Error("Expected prompts capability to be present")
	}
}

func TestMCPServer_HandleToolsList(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	request := &dto.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      1,
	}

	server.handleToolsList(request)

	// Parse the response
	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// Check the tools list response
	resultBytes, _ := json.Marshal(response.Result)
	var toolsResponse dto.ToolsListResponse
	json.Unmarshal(resultBytes, &toolsResponse)

	// Should return tools from the test container
	if len(toolsResponse.Tools) == 0 {
		t.Log("No tools returned - this is expected for test container")
	}
}

func TestMCPServer_HandleResourcesList(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	request := &dto.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "resources/list",
		ID:      1,
	}

	server.handleResourcesList(request)

	// Parse the response
	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// Check the resources list response
	resultBytes, _ := json.Marshal(response.Result)
	var resourcesResponse dto.ResourcesListResponse
	json.Unmarshal(resultBytes, &resourcesResponse)

	if len(resourcesResponse.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resourcesResponse.Resources))
	}

	// Check that we have the expected resources
	expectedURIs := map[string]bool{
		"movies://database/all":   false,
		"movies://database/stats": false,
	}

	for _, resource := range resourcesResponse.Resources {
		if _, exists := expectedURIs[resource.URI]; exists {
			expectedURIs[resource.URI] = true
		}
	}

	for uri, found := range expectedURIs {
		if !found {
			t.Errorf("Expected resource %s not found", uri)
		}
	}
}

func TestMCPServer_HandlePromptsList(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	request := &dto.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "prompts/list",
		ID:      1,
	}

	server.handlePromptsList(request)

	// Parse the response
	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// Check the prompts list response
	resultBytes, _ := json.Marshal(response.Result)
	var promptsResponse dto.PromptsListResponse
	json.Unmarshal(resultBytes, &promptsResponse)

	// Should return prompts from the test container
	if len(promptsResponse.Prompts) == 0 {
		t.Log("No prompts returned - checking if PromptHandlers is available")
		if container.PromptHandlers == nil {
			t.Log("PromptHandlers is nil in test container")
		}
	}
}

func TestMCPServer_HandleUnknownMethod(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	request := &dto.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "unknown/method",
		ID:      1,
	}

	server.handleRequest(request)

	// Parse the response
	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error for unknown method")
	}

	if response.Error.Code != dto.MethodNotFound {
		t.Errorf("Expected error code %d, got %d", dto.MethodNotFound, response.Error.Code)
	}
}

func TestMCPServer_HandleInvalidJSON(t *testing.T) {
	input := strings.NewReader("invalid json\n")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	// Run the server for one iteration
	err := server.Run()
	if err != nil {
		t.Fatalf("Server run failed: %v", err)
	}

	// Parse the error response
	var response dto.JSONRPCResponse
	err = json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error for invalid JSON")
	}

	if response.Error.Code != dto.ParseError {
		t.Errorf("Expected error code %d, got %d", dto.ParseError, response.Error.Code)
	}
}

func TestMCPServer_SendResult(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	testResult := map[string]string{"test": "result"}
	server.sendResult(123, testResult)

	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", response.JSONRPC)
	}

	if response.ID != float64(123) {
		t.Errorf("Expected ID 123, got %v", response.ID)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	resultMap, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Error("Expected result to be a map")
	} else if resultMap["test"] != "result" {
		t.Errorf("Expected result test=result, got %v", resultMap)
	}
}

func TestMCPServer_SendError(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	server.sendError(456, dto.InvalidParams, "Test error", "error data")

	var response dto.JSONRPCResponse
	err := json.Unmarshal(output.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", response.JSONRPC)
	}

	if response.ID != float64(456) {
		t.Errorf("Expected ID 456, got %v", response.ID)
	}

	if response.Result != nil {
		t.Errorf("Expected no result, got %v", response.Result)
	}

	if response.Error == nil {
		t.Fatal("Expected error")
	}

	if response.Error.Code != dto.InvalidParams {
		t.Errorf("Expected error code %d, got %d", dto.InvalidParams, response.Error.Code)
	}

	if response.Error.Message != "Test error" {
		t.Errorf("Expected error message 'Test error', got %s", response.Error.Message)
	}

	if response.Error.Data != "error data" {
		t.Errorf("Expected error data 'error data', got %v", response.Error.Data)
	}
}

func TestMCPServer_InitToolHandlers(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)
	container := composition.NewTestContainer()

	server := NewMCPServer(input, output, logger, container)

	// Test that tool handlers map is initialized
	if server.toolHandlers == nil {
		t.Error("Expected tool handlers to be initialized")
	}

	// Since test container has minimal handlers, we expect limited tools
	// This is mainly testing that the initialization doesn't panic
	if len(server.toolHandlers) == 0 {
		t.Log("Tool handlers map is empty as expected for test container")
	}
}