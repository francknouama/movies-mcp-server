package protocol

import (
	"encoding/json"
	"testing"
)

func TestJSONRPCRequest(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test/method",
		Params:  json.RawMessage(`{"key":"value"}`),
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC = 2.0, got %s", req.JSONRPC)
	}

	if req.ID != 1 {
		t.Errorf("Expected ID = 1, got %v", req.ID)
	}

	if req.Method != "test/method" {
		t.Errorf("Expected Method = test/method, got %s", req.Method)
	}
}

func TestJSONRPCRequest_JSONMarshaling(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "req-123",
		Method:  "test/method",
		Params:  json.RawMessage(`{"param1":"value1"}`),
	}

	// Marshal to JSON
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Unmarshal back
	var decoded JSONRPCRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.JSONRPC != req.JSONRPC {
		t.Errorf("Decoded JSONRPC = %s, want %s", decoded.JSONRPC, req.JSONRPC)
	}

	if decoded.Method != req.Method {
		t.Errorf("Decoded Method = %s, want %s", decoded.Method, req.Method)
	}
}

func TestJSONRPCResponse(t *testing.T) {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
		Error:   nil,
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC = 2.0, got %s", resp.JSONRPC)
	}

	if resp.Result == nil {
		t.Error("Expected non-nil Result")
	}

	if resp.Error != nil {
		t.Error("Expected nil Error")
	}
}

func TestJSONRPCResponse_WithError(t *testing.T) {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  nil,
		Error: &JSONRPCError{
			Code:    -32600,
			Message: "Invalid Request",
			Data:    "additional data",
		},
	}

	if resp.Error == nil {
		t.Fatal("Expected non-nil Error")
	}

	if resp.Error.Code != -32600 {
		t.Errorf("Expected Error Code = -32600, got %d", resp.Error.Code)
	}

	if resp.Error.Message != "Invalid Request" {
		t.Errorf("Expected Error Message = 'Invalid Request', got %s", resp.Error.Message)
	}
}

func TestJSONRPCResponse_JSONMarshaling(t *testing.T) {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      42,
		Result:  map[string]interface{}{"success": true},
	}

	// Marshal to JSON
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Unmarshal back
	var decoded JSONRPCResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if decoded.JSONRPC != resp.JSONRPC {
		t.Errorf("Decoded JSONRPC = %s, want %s", decoded.JSONRPC, resp.JSONRPC)
	}
}

func TestJSONRPCError(t *testing.T) {
	err := &JSONRPCError{
		Code:    -32700,
		Message: "Parse error",
		Data:    map[string]string{"detail": "invalid JSON"},
	}

	if err.Code != -32700 {
		t.Errorf("Expected Code = -32700, got %d", err.Code)
	}

	if err.Message != "Parse error" {
		t.Errorf("Expected Message = 'Parse error', got %s", err.Message)
	}

	if err.Data == nil {
		t.Error("Expected non-nil Data")
	}
}

func TestMCPRequest(t *testing.T) {
	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "test/method",
		Params:  map[string]interface{}{"key": "value"},
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC = 2.0, got %s", req.JSONRPC)
	}

	if req.Method != "test/method" {
		t.Errorf("Expected Method = test/method, got %s", req.Method)
	}
}

func TestMCPRequest_ToJSONRPC(t *testing.T) {
	mcpReq := &MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test/method",
		Params:  map[string]interface{}{"key": "value"},
	}

	jsonrpcReq, err := mcpReq.ToJSONRPC()
	if err != nil {
		t.Fatalf("ToJSONRPC() error = %v", err)
	}

	if jsonrpcReq == nil {
		t.Fatal("ToJSONRPC() returned nil")
	}

	if jsonrpcReq.JSONRPC != mcpReq.JSONRPC {
		t.Errorf("JSONRPC = %s, want %s", jsonrpcReq.JSONRPC, mcpReq.JSONRPC)
	}

	if jsonrpcReq.Method != mcpReq.Method {
		t.Errorf("Method = %s, want %s", jsonrpcReq.Method, mcpReq.Method)
	}

	// Verify params can be unmarshaled
	var params map[string]interface{}
	if err := json.Unmarshal(jsonrpcReq.Params, &params); err != nil {
		t.Fatalf("Failed to unmarshal params: %v", err)
	}

	if params["key"] != "value" {
		t.Errorf("Params key = %v, want value", params["key"])
	}
}

func TestMCPRequest_ToJSONRPC_WithNilParams(t *testing.T) {
	mcpReq := &MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test/method",
		Params:  nil,
	}

	jsonrpcReq, err := mcpReq.ToJSONRPC()
	if err != nil {
		t.Fatalf("ToJSONRPC() error = %v", err)
	}

	if jsonrpcReq == nil {
		t.Fatal("ToJSONRPC() returned nil")
	}
}

func TestJSONRPCRequest_ToMCP(t *testing.T) {
	jsonrpcReq := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test/method",
		Params:  json.RawMessage(`{"key":"value"}`),
	}

	mcpReq, err := jsonrpcReq.ToMCP()
	if err != nil {
		t.Fatalf("ToMCP() error = %v", err)
	}

	if mcpReq == nil {
		t.Fatal("ToMCP() returned nil")
	}

	if mcpReq.JSONRPC != jsonrpcReq.JSONRPC {
		t.Errorf("JSONRPC = %s, want %s", mcpReq.JSONRPC, jsonrpcReq.JSONRPC)
	}

	if mcpReq.Method != jsonrpcReq.Method {
		t.Errorf("Method = %s, want %s", mcpReq.Method, jsonrpcReq.Method)
	}

	// Verify params
	params, ok := mcpReq.Params.(map[string]interface{})
	if !ok {
		t.Fatal("Expected params to be map[string]interface{}")
	}

	if params["key"] != "value" {
		t.Errorf("Params key = %v, want value", params["key"])
	}
}

func TestJSONRPCRequest_ToMCP_WithEmptyParams(t *testing.T) {
	jsonrpcReq := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test/method",
		Params:  json.RawMessage(``),
	}

	mcpReq, err := jsonrpcReq.ToMCP()
	if err != nil {
		t.Fatalf("ToMCP() error = %v", err)
	}

	if mcpReq.Params != nil {
		t.Errorf("Expected nil Params, got %v", mcpReq.Params)
	}
}

func TestMCPResponse(t *testing.T) {
	resp := &MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  "success",
		Error:   nil,
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC = 2.0, got %s", resp.JSONRPC)
	}

	if resp.Result != "success" {
		t.Errorf("Expected Result = success, got %v", resp.Result)
	}
}

func TestMCPResponse_ToJSONRPC(t *testing.T) {
	mcpResp := &MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
		Error:   nil,
	}

	jsonrpcResp := mcpResp.ToJSONRPC()
	if jsonrpcResp == nil {
		t.Fatal("ToJSONRPC() returned nil")
	}

	if jsonrpcResp.JSONRPC != mcpResp.JSONRPC {
		t.Errorf("JSONRPC = %s, want %s", jsonrpcResp.JSONRPC, mcpResp.JSONRPC)
	}

	if jsonrpcResp.Error != nil {
		t.Errorf("Expected nil Error, got %v", jsonrpcResp.Error)
	}
}

func TestMCPResponse_ToJSONRPC_WithError(t *testing.T) {
	mcpResp := &MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  nil,
		Error: &MCPError{
			Code:    -32600,
			Message: "Invalid Request",
			Data:    "extra info",
		},
	}

	jsonrpcResp := mcpResp.ToJSONRPC()
	if jsonrpcResp == nil {
		t.Fatal("ToJSONRPC() returned nil")
	}

	if jsonrpcResp.Error == nil {
		t.Fatal("Expected non-nil Error")
	}

	if jsonrpcResp.Error.Code != mcpResp.Error.Code {
		t.Errorf("Error Code = %d, want %d", jsonrpcResp.Error.Code, mcpResp.Error.Code)
	}

	if jsonrpcResp.Error.Message != mcpResp.Error.Message {
		t.Errorf("Error Message = %s, want %s", jsonrpcResp.Error.Message, mcpResp.Error.Message)
	}
}

func TestJSONRPCResponse_ToMCP(t *testing.T) {
	jsonrpcResp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
		Error:   nil,
	}

	mcpResp := jsonrpcResp.ToMCP()
	if mcpResp == nil {
		t.Fatal("ToMCP() returned nil")
	}

	if mcpResp.JSONRPC != jsonrpcResp.JSONRPC {
		t.Errorf("JSONRPC = %s, want %s", mcpResp.JSONRPC, jsonrpcResp.JSONRPC)
	}

	if mcpResp.Error != nil {
		t.Errorf("Expected nil Error, got %v", mcpResp.Error)
	}
}

func TestJSONRPCResponse_ToMCP_WithError(t *testing.T) {
	jsonrpcResp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  nil,
		Error: &JSONRPCError{
			Code:    -32700,
			Message: "Parse error",
			Data:    "details",
		},
	}

	mcpResp := jsonrpcResp.ToMCP()
	if mcpResp == nil {
		t.Fatal("ToMCP() returned nil")
	}

	if mcpResp.Error == nil {
		t.Fatal("Expected non-nil Error")
	}

	if mcpResp.Error.Code != jsonrpcResp.Error.Code {
		t.Errorf("Error Code = %d, want %d", mcpResp.Error.Code, jsonrpcResp.Error.Code)
	}

	if mcpResp.Error.Message != jsonrpcResp.Error.Message {
		t.Errorf("Error Message = %s, want %s", mcpResp.Error.Message, jsonrpcResp.Error.Message)
	}
}

func TestMCPError(t *testing.T) {
	err := &MCPError{
		Code:    -32601,
		Message: "Method not found",
		Data:    "methodName",
	}

	if err.Code != -32601 {
		t.Errorf("Expected Code = -32601, got %d", err.Code)
	}

	if err.Message != "Method not found" {
		t.Errorf("Expected Message = 'Method not found', got %s", err.Message)
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test MCP -> JSONRPC -> MCP roundtrip
	originalMCP := &MCPRequest{
		JSONRPC: "2.0",
		ID:      123,
		Method:  "test/roundtrip",
		Params:  map[string]interface{}{"test": "data"},
	}

	// Convert to JSONRPC
	jsonrpcReq, err := originalMCP.ToJSONRPC()
	if err != nil {
		t.Fatalf("ToJSONRPC() error = %v", err)
	}

	// Convert back to MCP
	mcpReq, err := jsonrpcReq.ToMCP()
	if err != nil {
		t.Fatalf("ToMCP() error = %v", err)
	}

	// Verify data integrity
	if mcpReq.JSONRPC != originalMCP.JSONRPC {
		t.Errorf("Roundtrip JSONRPC = %s, want %s", mcpReq.JSONRPC, originalMCP.JSONRPC)
	}

	if mcpReq.Method != originalMCP.Method {
		t.Errorf("Roundtrip Method = %s, want %s", mcpReq.Method, originalMCP.Method)
	}
}

func TestDifferentIDTypes(t *testing.T) {
	tests := []struct {
		name string
		id   interface{}
	}{
		{"integer", 1},
		{"string", "req-123"},
		{"float", 3.14},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      tt.id,
				Method:  "test",
			}

			// Should be able to marshal and unmarshal
			data, err := json.Marshal(req)
			if err != nil {
				t.Fatalf("Marshal error = %v", err)
			}

			var decoded JSONRPCRequest
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Unmarshal error = %v", err)
			}
		})
	}
}
