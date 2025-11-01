package communication

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/francknouama/movies-mcp-server/pkg/protocol"
)

// Test helper to create a sample request.
func createTestRequest(id interface{}, method string) *protocol.JSONRPCRequest {
	return &protocol.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  json.RawMessage(`{"test":"data"}`),
	}
}

// Test helper to create a sample response.
func createTestResponse(id interface{}, result interface{}) *protocol.JSONRPCResponse {
	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// Test helper to create an error response.
func createTestErrorResponse(id interface{}, code int, message string) *protocol.JSONRPCResponse {
	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &protocol.JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}

func TestNewStdioTransport(t *testing.T) {
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}

	transport := NewStdioTransport(reader, writer)

	if transport == nil {
		t.Fatal("NewStdioTransport() returned nil")
	}

	if transport.reader != reader {
		t.Error("StdioTransport.reader is not set correctly")
	}

	if transport.writer != writer {
		t.Error("StdioTransport.writer is not set correctly")
	}

	if transport.scanner == nil {
		t.Error("StdioTransport.scanner should not be nil")
	}

	if transport.encoder == nil {
		t.Error("StdioTransport.encoder should not be nil")
	}
}

func TestStdioTransport_Send(t *testing.T) {
	tests := []struct {
		name     string
		response *protocol.JSONRPCResponse
		wantErr  bool
	}{
		{
			name:     "send success response",
			response: createTestResponse(1, map[string]string{"status": "ok"}),
			wantErr:  false,
		},
		{
			name:     "send error response",
			response: createTestErrorResponse(2, -32600, "Invalid Request"),
			wantErr:  false,
		},
		{
			name:     "send response with string id",
			response: createTestResponse("req-123", "success"),
			wantErr:  false,
		},
		{
			name:     "send response with nil result",
			response: createTestResponse(3, nil),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			transport := NewStdioTransport(nil, writer)

			err := transport.Send(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the response was written
				if writer.Len() == 0 {
					t.Error("Send() did not write any data")
				}

				// Verify it's valid JSON with newline
				output := writer.String()
				if !strings.HasSuffix(output, "\n") {
					t.Error("Send() output should end with newline")
				}

				// Verify we can decode the response
				var decoded protocol.JSONRPCResponse
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &decoded); err != nil {
					t.Errorf("Send() output is not valid JSON: %v", err)
				}

				if decoded.JSONRPC != tt.response.JSONRPC {
					t.Errorf("Send() JSONRPC = %v, want %v", decoded.JSONRPC, tt.response.JSONRPC)
				}
			}
		})
	}
}

func TestStdioTransport_Send_ConcurrentAccess(t *testing.T) {
	writer := &bytes.Buffer{}
	transport := NewStdioTransport(nil, writer)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			response := createTestResponse(id, "result")
			if err := transport.Send(response); err != nil {
				t.Errorf("Send() from goroutine %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify we got all responses
	lines := strings.Split(strings.TrimSpace(writer.String()), "\n")
	if len(lines) != goroutines {
		t.Errorf("Expected %d responses, got %d", goroutines, len(lines))
	}
}

func TestStdioTransport_Receive(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMethod string
		wantID     interface{}
		wantErr    bool
	}{
		{
			name:       "receive valid request with integer id",
			input:      `{"jsonrpc":"2.0","id":1,"method":"test/method","params":{"key":"value"}}`,
			wantMethod: "test/method",
			wantID:     float64(1), // JSON numbers decode as float64
			wantErr:    false,
		},
		{
			name:       "receive valid request with string id",
			input:      `{"jsonrpc":"2.0","id":"req-123","method":"another/method"}`,
			wantMethod: "another/method",
			wantID:     "req-123",
			wantErr:    false,
		},
		{
			name:       "receive notification (no id)",
			input:      `{"jsonrpc":"2.0","method":"notification/method"}`,
			wantMethod: "notification/method",
			wantID:     nil,
			wantErr:    false,
		},
		{
			name:    "receive invalid json",
			input:   `{invalid json}`,
			wantErr: true,
		},
		{
			name:    "receive empty line",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input + "\n")
			transport := NewStdioTransport(reader, nil)

			request, err := transport.Receive()
			if (err != nil) != tt.wantErr {
				t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if request == nil {
					t.Fatal("Receive() returned nil request")
				}

				if request.Method != tt.wantMethod {
					t.Errorf("Receive() method = %v, want %v", request.Method, tt.wantMethod)
				}

				if request.ID != tt.wantID {
					t.Errorf("Receive() id = %v, want %v", request.ID, tt.wantID)
				}

				if request.JSONRPC != "2.0" {
					t.Errorf("Receive() jsonrpc = %v, want 2.0", request.JSONRPC)
				}
			}
		})
	}
}

func TestStdioTransport_Receive_MultipleRequests(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"first"}
{"jsonrpc":"2.0","id":2,"method":"second"}
{"jsonrpc":"2.0","id":3,"method":"third"}
`
	reader := strings.NewReader(input)
	transport := NewStdioTransport(reader, nil)

	expectedMethods := []string{"first", "second", "third"}

	for i, expectedMethod := range expectedMethods {
		request, err := transport.Receive()
		if err != nil {
			t.Fatalf("Receive() request %d error = %v", i+1, err)
		}

		if request.Method != expectedMethod {
			t.Errorf("Receive() request %d method = %v, want %v", i+1, request.Method, expectedMethod)
		}
	}

	// Fourth receive should fail (EOF)
	_, err := transport.Receive()
	if err == nil {
		t.Error("Receive() should return error after all data is consumed")
	}
}

func TestStdioTransport_Close(t *testing.T) {
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}
	transport := NewStdioTransport(reader, writer)

	err := transport.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestNewBufferedTransport(t *testing.T) {
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}
	stdioTransport := NewStdioTransport(reader, writer)

	bufferedTransport := NewBufferedTransport(stdioTransport)

	if bufferedTransport == nil {
		t.Fatal("NewBufferedTransport() returned nil")
	}

	if bufferedTransport.transport != stdioTransport {
		t.Error("BufferedTransport.transport is not set correctly")
	}

	if bufferedTransport.buffer == nil {
		t.Error("BufferedTransport.buffer should not be nil")
	}

	if cap(bufferedTransport.buffer) != 4096 {
		t.Errorf("BufferedTransport.buffer capacity = %d, want 4096", cap(bufferedTransport.buffer))
	}
}

func TestBufferedTransport_Send(t *testing.T) {
	writer := &bytes.Buffer{}
	stdioTransport := NewStdioTransport(nil, writer)
	bufferedTransport := NewBufferedTransport(stdioTransport)

	response := createTestResponse(1, "test")
	err := bufferedTransport.Send(response)

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}

	// Verify the response was written to the underlying transport
	if writer.Len() == 0 {
		t.Error("Send() did not write any data to underlying transport")
	}
}

func TestBufferedTransport_Receive(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"test/method"}`
	reader := strings.NewReader(input + "\n")
	stdioTransport := NewStdioTransport(reader, nil)
	bufferedTransport := NewBufferedTransport(stdioTransport)

	request, err := bufferedTransport.Receive()
	if err != nil {
		t.Fatalf("Receive() error = %v", err)
	}

	if request.Method != "test/method" {
		t.Errorf("Receive() method = %v, want test/method", request.Method)
	}
}

func TestBufferedTransport_Close(t *testing.T) {
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}
	stdioTransport := NewStdioTransport(reader, writer)
	bufferedTransport := NewBufferedTransport(stdioTransport)

	err := bufferedTransport.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Verify the underlying transport was closed
	// (StdioTransport.Close() is a no-op, so we just verify no error)
}

func TestNewMockTransport(t *testing.T) {
	mockTransport := NewMockTransport()

	if mockTransport == nil {
		t.Fatal("NewMockTransport() returned nil")
	}

	if mockTransport.requests == nil {
		t.Error("MockTransport.requests should not be nil")
	}

	if cap(mockTransport.requests) != 10 {
		t.Errorf("MockTransport.requests capacity = %d, want 10", cap(mockTransport.requests))
	}

	if mockTransport.responses == nil {
		t.Error("MockTransport.responses should not be nil")
	}

	if cap(mockTransport.responses) != 10 {
		t.Errorf("MockTransport.responses capacity = %d, want 10", cap(mockTransport.responses))
	}

	if mockTransport.closed {
		t.Error("MockTransport should not be closed initially")
	}
}

func TestMockTransport_Send(t *testing.T) {
	mockTransport := NewMockTransport()

	response := createTestResponse(1, "test result")
	err := mockTransport.Send(response)

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}

	// Verify the response was buffered
	select {
	case received := <-mockTransport.responses:
		if received.ID != response.ID {
			t.Errorf("Buffered response ID = %v, want %v", received.ID, response.ID)
		}
		if received.Result != response.Result {
			t.Errorf("Buffered response Result = %v, want %v", received.Result, response.Result)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Response was not buffered")
	}
}

func TestMockTransport_Send_AfterClose(t *testing.T) {
	mockTransport := NewMockTransport()
	mockTransport.Close()

	response := createTestResponse(1, "test")
	err := mockTransport.Send(response)

	if err == nil {
		t.Error("Send() should return error after Close()")
	}

	expectedErr := "transport is closed"
	if err.Error() != expectedErr {
		t.Errorf("Send() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestMockTransport_Receive(t *testing.T) {
	mockTransport := NewMockTransport()

	request := createTestRequest(1, "test/method")

	// Send a request using SendRequest method
	err := mockTransport.SendRequest(request)
	if err != nil {
		t.Fatalf("SendRequest() error = %v", err)
	}

	received, err := mockTransport.Receive()
	if err != nil {
		t.Fatalf("Receive() error = %v", err)
	}

	if received.ID != request.ID {
		t.Errorf("Receive() ID = %v, want %v", received.ID, request.ID)
	}

	if received.Method != request.Method {
		t.Errorf("Receive() Method = %v, want %v", received.Method, request.Method)
	}
}

func TestMockTransport_Receive_AfterClose(t *testing.T) {
	mockTransport := NewMockTransport()
	mockTransport.Close()

	_, err := mockTransport.Receive()

	if err == nil {
		t.Error("Receive() should return error after Close()")
	}

	expectedErr := "transport is closed"
	if err.Error() != expectedErr {
		t.Errorf("Receive() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestMockTransport_Close(t *testing.T) {
	mockTransport := NewMockTransport()

	err := mockTransport.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if !mockTransport.closed {
		t.Error("MockTransport.closed should be true after Close()")
	}

	// Verify channels are closed
	select {
	case _, ok := <-mockTransport.requests:
		if ok {
			t.Error("requests channel should be closed")
		}
	default:
		t.Error("requests channel should be closed")
	}

	select {
	case _, ok := <-mockTransport.responses:
		if ok {
			t.Error("responses channel should be closed")
		}
	default:
		t.Error("responses channel should be closed")
	}
}

func TestMockTransport_Close_Idempotent(t *testing.T) {
	mockTransport := NewMockTransport()

	// Close multiple times
	err1 := mockTransport.Close()
	err2 := mockTransport.Close()
	err3 := mockTransport.Close()

	if err1 != nil {
		t.Errorf("First Close() error = %v", err1)
	}

	if err2 != nil {
		t.Errorf("Second Close() error = %v", err2)
	}

	if err3 != nil {
		t.Errorf("Third Close() error = %v", err3)
	}
}

func TestMockTransport_GetResponse(t *testing.T) {
	mockTransport := NewMockTransport()

	// Send a response
	response := createTestResponse(1, "test result")
	mockTransport.Send(response)

	// Get the response
	received, err := mockTransport.GetResponse()
	if err != nil {
		t.Fatalf("GetResponse() error = %v", err)
	}
	if received == nil {
		t.Fatal("GetResponse() returned nil")
	}

	if received.ID != response.ID {
		t.Errorf("GetResponse() ID = %v, want %v", received.ID, response.ID)
	}

	if received.Result != response.Result {
		t.Errorf("GetResponse() Result = %v, want %v", received.Result, response.Result)
	}
}

func TestMockTransport_GetResponse_Error(t *testing.T) {
	mockTransport := NewMockTransport()

	// Don't send any response
	received, err := mockTransport.GetResponse()
	if err == nil {
		t.Error("GetResponse() should return error when no response is available")
	}
	if received != nil {
		t.Error("GetResponse() should return nil when no response is available")
	}
}

func TestMockTransport_SendRequest(t *testing.T) {
	mockTransport := NewMockTransport()

	request := createTestRequest(1, "test/method")
	err := mockTransport.SendRequest(request)

	if err != nil {
		t.Errorf("SendRequest() error = %v", err)
	}

	// Verify the request was buffered
	received, err := mockTransport.Receive()
	if err != nil {
		t.Fatalf("Receive() error = %v", err)
	}

	if received.ID != request.ID {
		t.Errorf("Received request ID = %v, want %v", received.ID, request.ID)
	}

	if received.Method != request.Method {
		t.Errorf("Received request Method = %v, want %v", received.Method, request.Method)
	}
}

func TestMockTransport_SendRequest_AfterClose(t *testing.T) {
	mockTransport := NewMockTransport()
	mockTransport.Close()

	request := createTestRequest(1, "test/method")
	err := mockTransport.SendRequest(request)

	if err == nil {
		t.Error("SendRequest() should return error after Close()")
	}

	expectedErr := "transport is closed"
	if err.Error() != expectedErr {
		t.Errorf("SendRequest() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestMockTransport_ConcurrentOperations(t *testing.T) {
	mockTransport := NewMockTransport()
	const operations = 5 // Limited by channel buffer size of 10

	var wg sync.WaitGroup
	wg.Add(operations * 2) // Send and Receive for each operation

	// Concurrent sends
	for i := 0; i < operations; i++ {
		go func(id int) {
			defer wg.Done()
			response := createTestResponse(id, "result")
			if err := mockTransport.Send(response); err != nil {
				t.Errorf("Concurrent Send(%d) error = %v", id, err)
			}
		}(i)
	}

	// Concurrent receives
	for i := 0; i < operations; i++ {
		go func(id int) {
			defer wg.Done()
			request := createTestRequest(id, "method")
			if err := mockTransport.SendRequest(request); err != nil {
				t.Errorf("Concurrent SendRequest(%d) error = %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify we got all responses
	responsesReceived := 0
	timeout := time.After(1 * time.Second)
	for responsesReceived < operations {
		select {
		case <-mockTransport.responses:
			responsesReceived++
		case <-timeout:
			t.Fatalf("Timeout waiting for responses, received %d/%d", responsesReceived, operations)
		}
	}

	// Verify we got all requests
	requestsReceived := 0
	timeout = time.After(1 * time.Second)
	for requestsReceived < operations {
		select {
		case <-mockTransport.requests:
			requestsReceived++
		case <-timeout:
			t.Fatalf("Timeout waiting for requests, received %d/%d", requestsReceived, operations)
		}
	}
}

func TestTransport_InterfaceCompliance(t *testing.T) {
	// Verify all implementations satisfy the Transport interface
	var _ Transport = (*StdioTransport)(nil)
	var _ Transport = (*BufferedTransport)(nil)
	var _ Transport = (*MockTransport)(nil)
}

// Test error scenarios for StdioTransport.
func TestStdioTransport_Send_WriterError(t *testing.T) {
	// Create a writer that always returns an error
	errorWriter := &errorWriter{err: errors.New("write error")}
	transport := NewStdioTransport(nil, errorWriter)

	response := createTestResponse(1, "test")
	err := transport.Send(response)

	if err == nil {
		t.Error("Send() should return error when writer fails")
	}
}

// errorWriter is a helper type that always returns an error on Write.
type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}
