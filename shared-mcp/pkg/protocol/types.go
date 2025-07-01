package protocol

import (
	"encoding/json"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCPRequest represents an MCP protocol request (higher-level wrapper)
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents an MCP protocol response (higher-level wrapper)
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP protocol error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToJSONRPC converts an MCPRequest to a JSONRPCRequest
func (r *MCPRequest) ToJSONRPC() (*JSONRPCRequest, error) {
	params, err := json.Marshal(r.Params)
	if err != nil {
		return nil, err
	}

	return &JSONRPCRequest{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Method:  r.Method,
		Params:  params,
	}, nil
}

// ToMCP converts a JSONRPCRequest to an MCPRequest
func (r *JSONRPCRequest) ToMCP() (*MCPRequest, error) {
	var params interface{}
	if len(r.Params) > 0 {
		if err := json.Unmarshal(r.Params, &params); err != nil {
			return nil, err
		}
	}

	return &MCPRequest{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Method:  r.Method,
		Params:  params,
	}, nil
}

// ToJSONRPC converts an MCPResponse to a JSONRPCResponse
func (r *MCPResponse) ToJSONRPC() *JSONRPCResponse {
	var err *JSONRPCError
	if r.Error != nil {
		err = &JSONRPCError{
			Code:    r.Error.Code,
			Message: r.Error.Message,
			Data:    r.Error.Data,
		}
	}

	return &JSONRPCResponse{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Result:  r.Result,
		Error:   err,
	}
}

// ToMCP converts a JSONRPCResponse to an MCPResponse
func (r *JSONRPCResponse) ToMCP() *MCPResponse {
	var err *MCPError
	if r.Error != nil {
		err = &MCPError{
			Code:    r.Error.Code,
			Message: r.Error.Message,
			Data:    r.Error.Data,
		}
	}

	return &MCPResponse{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Result:  r.Result,
		Error:   err,
	}
}
