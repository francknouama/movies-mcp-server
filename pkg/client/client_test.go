package client

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/francknouama/movies-mcp-server/pkg/communication"
	"github.com/francknouama/movies-mcp-server/pkg/protocol"
)

func TestNewMCPClient(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	options := ClientOptions{
		Transport: mockTransport,
		Timeout:   10 * time.Second,
		ClientInfo: protocol.ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
	}

	client := NewMCPClient(options)

	if client == nil {
		t.Fatal("NewMCPClient() returned nil")
	}

	if client.transport != mockTransport {
		t.Error("Client transport not set correctly")
	}

	if client.timeout != 10*time.Second {
		t.Errorf("Client timeout = %v, want %v", client.timeout, 10*time.Second)
	}

	if client.requestID != 1 {
		t.Errorf("Client requestID = %d, want 1", client.requestID)
	}
}

func TestNewMCPClient_DefaultTimeout(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	options := ClientOptions{
		Transport: mockTransport,
		// Timeout not specified
	}

	client := NewMCPClient(options)

	if client.timeout != 30*time.Second {
		t.Errorf("Default timeout = %v, want %v", client.timeout, 30*time.Second)
	}
}

func TestNewStdioMCPClient(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}

	client := NewStdioMCPClient(reader, writer)

	if client == nil {
		t.Fatal("NewStdioMCPClient() returned nil")
	}

	if client.transport == nil {
		t.Error("Expected non-nil transport")
	}
}

func TestMCPClient_IsInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	if client.IsInitialized() {
		t.Error("Expected client to not be initialized")
	}

	// Manually set initialized for testing
	client.initialized = true

	if !client.IsInitialized() {
		t.Error("Expected client to be initialized")
	}
}

func TestMCPClient_GetServerCapabilities(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	caps := client.GetServerCapabilities()
	if caps != nil {
		t.Errorf("Expected nil capabilities before initialization, got %v", caps)
	}

	// Set capabilities
	testCaps := &protocol.ServerCapabilities{
		Tools: &protocol.ToolsCapability{
			ListChanged: true,
		},
	}
	client.capabilities = testCaps

	caps = client.GetServerCapabilities()
	if caps == nil {
		t.Fatal("Expected non-nil capabilities")
	}

	if caps.Tools == nil || !caps.Tools.ListChanged {
		t.Error("Expected Tools capability with ListChanged = true")
	}
}

func TestMCPClient_GetServerInfo(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	info := client.GetServerInfo()
	if info != nil {
		t.Errorf("Expected nil server info before initialization, got %v", info)
	}

	// Set server info
	testInfo := &protocol.ServerInfo{
		Name:    "test-server",
		Version: "1.0.0",
	}
	client.serverInfo = testInfo

	info = client.GetServerInfo()
	if info == nil {
		t.Fatal("Expected non-nil server info")
	}

	if info.Name != "test-server" {
		t.Errorf("Server name = %s, want test-server", info.Name)
	}
}

func TestMCPClient_Close(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	err := client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if client.transport != nil {
		t.Error("Expected transport to be nil after Close()")
	}
}

func TestMCPClient_Close_NilTransport(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Close once
	client.Close()

	// Close again should not error
	err := client.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}

func TestMCPClient_CallTool_NotInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	_, err := client.CallTool("test-tool", nil)
	if err == nil {
		t.Error("Expected error when calling tool before initialization")
	}

	if err.Error() != "client not initialized" {
		t.Errorf("Error message = %v, want 'client not initialized'", err.Error())
	}
}

func TestMCPClient_ListTools_NotInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	_, err := client.ListTools()
	if err == nil {
		t.Error("Expected error when listing tools before initialization")
	}

	if err.Error() != "client not initialized" {
		t.Errorf("Error message = %v, want 'client not initialized'", err.Error())
	}
}

func TestMCPClient_ListResources_NotInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	_, err := client.ListResources()
	if err == nil {
		t.Error("Expected error when listing resources before initialization")
	}

	if err.Error() != "client not initialized" {
		t.Errorf("Error message = %v, want 'client not initialized'", err.Error())
	}
}

func TestMCPClient_ReadResource_NotInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	_, err := client.ReadResource("test://resource")
	if err == nil {
		t.Error("Expected error when reading resource before initialization")
	}

	if err.Error() != "client not initialized" {
		t.Errorf("Error message = %v, want 'client not initialized'", err.Error())
	}
}

func TestMCPClient_Initialize_AlreadyInitialized(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Mark as initialized
	client.initialized = true

	err := client.Initialize(protocol.ClientInfo{}, protocol.ClientCapabilities{})
	if err == nil {
		t.Error("Expected error when initializing already initialized client")
	}

	if err.Error() != "client already initialized" {
		t.Errorf("Error message = %v, want 'client already initialized'", err.Error())
	}
}

func TestMCPClient_nextRequestID(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Initial ID should be 1
	id1 := client.nextRequestID()
	if id1 != int64(1) {
		t.Errorf("First request ID = %v, want 1", id1)
	}

	// Second ID should be 2
	id2 := client.nextRequestID()
	if id2 != int64(2) {
		t.Errorf("Second request ID = %v, want 2", id2)
	}

	// Third ID should be 3
	id3 := client.nextRequestID()
	if id3 != int64(3) {
		t.Errorf("Third request ID = %v, want 3", id3)
	}
}

func TestMCPClient_marshalParams(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	tests := []struct {
		name   string
		params interface{}
		valid  bool
	}{
		{
			name:   "valid struct",
			params: struct{ Name string }{Name: "test"},
			valid:  true,
		},
		{
			name:   "valid map",
			params: map[string]string{"key": "value"},
			valid:  true,
		},
		{
			name:   "nil params",
			params: nil,
			valid:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.marshalParams(tt.params)
			if result == nil {
				t.Error("marshalParams() returned nil")
			}

			// Verify it's valid JSON
			var decoded interface{}
			if err := json.Unmarshal(result, &decoded); err != nil {
				t.Errorf("marshalParams() produced invalid JSON: %v", err)
			}
		})
	}
}

func TestMCPClient_marshalParams_InvalidData(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Channels cannot be marshaled to JSON
	params := make(chan int)

	result := client.marshalParams(params)

	// Should return empty JSON object on error
	expected := []byte("{}")
	if string(result) != string(expected) {
		t.Errorf("marshalParams() with invalid data = %s, want %s", result, expected)
	}
}

func TestMCPClient_unmarshalResult(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	result := map[string]interface{}{
		"name":    "test",
		"version": "1.0.0",
	}

	var target struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	err := client.unmarshalResult(result, &target)
	if err != nil {
		t.Fatalf("unmarshalResult() error = %v", err)
	}

	if target.Name != "test" {
		t.Errorf("Target Name = %s, want test", target.Name)
	}

	if target.Version != "1.0.0" {
		t.Errorf("Target Version = %s, want 1.0.0", target.Version)
	}
}

func TestMCPClient_unmarshalResult_Error(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Channel cannot be marshaled
	result := make(chan int)

	var target interface{}

	err := client.unmarshalResult(result, &target)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid data")
	}
}

func TestClientOptions(t *testing.T) {
	mockTransport := communication.NewMockTransport()

	options := ClientOptions{
		Transport: mockTransport,
		Timeout:   5 * time.Second,
		ClientInfo: protocol.ClientInfo{
			Name:    "test",
			Version: "1.0",
		},
		Capabilities: protocol.ClientCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: true,
			},
		},
	}

	if options.Transport != mockTransport {
		t.Error("ClientOptions Transport not set correctly")
	}

	if options.Timeout != 5*time.Second {
		t.Errorf("ClientOptions Timeout = %v, want %v", options.Timeout, 5*time.Second)
	}

	if options.ClientInfo.Name != "test" {
		t.Errorf("ClientOptions ClientInfo.Name = %s, want test", options.ClientInfo.Name)
	}

	if options.Capabilities.Tools == nil || !options.Capabilities.Tools.ListChanged {
		t.Error("ClientOptions Capabilities not set correctly")
	}
}

func TestMCPClient_ConcurrentAccess(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Set server info and capabilities
	client.capabilities = &protocol.ServerCapabilities{}
	client.serverInfo = &protocol.ServerInfo{Name: "test", Version: "1.0"}

	const goroutines = 10
	done := make(chan bool, goroutines)

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		go func() {
			_ = client.GetServerCapabilities()
			_ = client.GetServerInfo()
			_ = client.IsInitialized()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func TestMCPClient_RequestIDIncrement(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Generate multiple IDs
	ids := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		ids[i] = client.nextRequestID()
	}

	// Verify all IDs are unique and sequential
	for i, id := range ids {
		expected := int64(i + 1)
		if id != expected {
			t.Errorf("Request ID at index %d = %v, want %d", i, id, expected)
		}
	}
}

func TestMCPClient_ZeroValueFields(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	if client.initialized {
		t.Error("Expected initialized to be false")
	}

	if client.capabilities != nil {
		t.Error("Expected capabilities to be nil")
	}

	if client.serverInfo != nil {
		t.Error("Expected serverInfo to be nil")
	}
}

// Integration tests for methods requiring transport interaction

func TestMCPClient_Initialize_Success(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	clientInfo := protocol.ClientInfo{
		Name:    "test-client",
		Version: "1.0.0",
	}
	capabilities := protocol.ClientCapabilities{}

	// Pre-queue the response that will be returned
	initResponse := &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64(1), // First request ID
		Params: json.RawMessage(`{
			"protocolVersion": "2024-11-05",
			"capabilities": {"tools": {"listChanged": true}},
			"serverInfo": {"name": "test-server", "version": "1.0.0"}
		}`),
	}
	mockTransport.SendRequest(initResponse)

	err := client.Initialize(clientInfo, capabilities)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if !client.IsInitialized() {
		t.Error("Client should be initialized after successful Initialize()")
	}

	if client.GetServerInfo() == nil {
		t.Error("ServerInfo should be set after Initialize()")
	}

	if client.GetServerCapabilities() == nil {
		t.Error("ServerCapabilities should be set after Initialize()")
	}
}

func TestMCPClient_Initialize_TransportError(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Close transport to cause send error
	mockTransport.Close()

	err := client.Initialize(protocol.ClientInfo{}, protocol.ClientCapabilities{})
	if err == nil {
		t.Error("Initialize() should return error when transport fails")
	}
}

func TestMCPClient_Initialize_ResponseError(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})

	// Set up mock error response
	go func() {
		resp, err := mockTransport.GetResponse()
		if err != nil {
			return
		}

		// Send back error response
		errorResponse := &protocol.JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      resp.ID,
			Params: json.RawMessage(`{
				"error": {"code": -32600, "message": "Invalid request"}
			}`),
		}
		mockTransport.SendRequest(errorResponse)
	}()

	err := client.Initialize(protocol.ClientInfo{}, protocol.ClientCapabilities{})
	if err == nil {
		t.Error("Initialize() should return error when server returns error")
	}
}

func TestMCPClient_CallTool_Success(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true // Manually set for testing

	// Pre-queue the response
	toolResponse := &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64(1),
		Params: json.RawMessage(`{
			"content": [{"type": "text", "text": "Tool result"}]
		}`),
	}
	mockTransport.SendRequest(toolResponse)

	result, err := client.CallTool("test-tool", map[string]interface{}{"arg": "value"})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}

	if result == nil {
		t.Error("CallTool() should return result")
	}
}

func TestMCPClient_CallTool_ErrorResponse(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Set up mock error response
	go func() {
		resp, err := mockTransport.GetResponse()
		if err != nil {
			return
		}

		errorResponse := &protocol.JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      resp.ID,
			Params: json.RawMessage(`{
				"error": {"code": -32000, "message": "Tool not found"}
			}`),
		}
		mockTransport.SendRequest(errorResponse)
	}()

	_, err := client.CallTool("nonexistent-tool", nil)
	if err == nil {
		t.Error("CallTool() should return error when server returns error")
	}
}

func TestMCPClient_ListTools_Success(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Pre-queue the response
	listResponse := &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64(1),
		Params: json.RawMessage(`{
			"tools": [
				{"name": "tool1", "description": "First tool"},
				{"name": "tool2", "description": "Second tool"}
			]
		}`),
	}
	mockTransport.SendRequest(listResponse)

	result, err := client.ListTools()
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}

	if result == nil {
		t.Fatal("ListTools() should return result")
	}

	if len(result.Tools) != 2 {
		t.Errorf("ListTools() returned %d tools, want 2", len(result.Tools))
	}
}

func TestMCPClient_ListTools_ErrorResponse(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Set up mock error response
	go func() {
		resp, err := mockTransport.GetResponse()
		if err != nil {
			return
		}

		errorResponse := &protocol.JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      resp.ID,
			Params: json.RawMessage(`{
				"error": {"code": -32603, "message": "Internal error"}
			}`),
		}
		mockTransport.SendRequest(errorResponse)
	}()

	_, err := client.ListTools()
	if err == nil {
		t.Error("ListTools() should return error when server returns error")
	}
}

func TestMCPClient_ListResources_Success(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Pre-queue the response
	listResponse := &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64(1),
		Params: json.RawMessage(`{
			"resources": [
				{"uri": "resource://1", "name": "Resource 1"},
				{"uri": "resource://2", "name": "Resource 2"}
			]
		}`),
	}
	mockTransport.SendRequest(listResponse)

	result, err := client.ListResources()
	if err != nil {
		t.Fatalf("ListResources() error = %v", err)
	}

	if result == nil {
		t.Fatal("ListResources() should return result")
	}

	if len(result.Resources) != 2 {
		t.Errorf("ListResources() returned %d resources, want 2", len(result.Resources))
	}
}

func TestMCPClient_ListResources_ErrorResponse(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Set up mock error response
	go func() {
		resp, err := mockTransport.GetResponse()
		if err != nil {
			return
		}

		errorResponse := &protocol.JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      resp.ID,
			Params: json.RawMessage(`{
				"error": {"code": -32603, "message": "Internal error"}
			}`),
		}
		mockTransport.SendRequest(errorResponse)
	}()

	_, err := client.ListResources()
	if err == nil {
		t.Error("ListResources() should return error when server returns error")
	}
}

func TestMCPClient_ReadResource_Success(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Pre-queue the response
	readResponse := &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64(1),
		Params: json.RawMessage(`{
			"contents": [
				{"uri": "resource://test", "mimeType": "text/plain", "text": "Resource content"}
			]
		}`),
	}
	mockTransport.SendRequest(readResponse)

	result, err := client.ReadResource("resource://test")
	if err != nil {
		t.Fatalf("ReadResource() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadResource() should return result")
	}

	if len(result.Contents) != 1 {
		t.Errorf("ReadResource() returned %d contents, want 1", len(result.Contents))
	}
}

func TestMCPClient_ReadResource_ErrorResponse(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Set up mock error response
	go func() {
		resp, err := mockTransport.GetResponse()
		if err != nil {
			return
		}

		errorResponse := &protocol.JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      resp.ID,
			Params: json.RawMessage(`{
				"error": {"code": -32010, "message": "Resource not found"}
			}`),
		}
		mockTransport.SendRequest(errorResponse)
	}()

	_, err := client.ReadResource("resource://nonexistent")
	if err == nil {
		t.Error("ReadResource() should return error when server returns error")
	}
}

func TestMCPClient_ReadResource_TransportError(t *testing.T) {
	mockTransport := communication.NewMockTransport()
	client := NewMCPClient(ClientOptions{Transport: mockTransport})
	client.initialized = true

	// Close transport to cause error
	mockTransport.Close()

	_, err := client.ReadResource("resource://test")
	if err == nil {
		t.Error("ReadResource() should return error when transport fails")
	}
}
