package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/francknouama/movies-mcp-server/shared-mcp/pkg/protocol"
)

// Transport defines the interface for MCP communication transport
type Transport interface {
	Send(response *protocol.JSONRPCResponse) error
	Receive() (*protocol.JSONRPCRequest, error)
	Close() error
}

// StdioTransport implements MCP communication over stdin/stdout
type StdioTransport struct {
	reader  io.Reader
	writer  io.Writer
	scanner *bufio.Scanner
	encoder *json.Encoder
	mutex   sync.Mutex
}

// NewStdioTransport creates a new stdin/stdout transport
func NewStdioTransport(reader io.Reader, writer io.Writer) *StdioTransport {
	return &StdioTransport{
		reader:  reader,
		writer:  writer,
		scanner: bufio.NewScanner(reader),
		encoder: json.NewEncoder(writer),
	}
}

// Send sends a JSON-RPC response through the transport
func (t *StdioTransport) Send(response *protocol.JSONRPCResponse) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.encoder.Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}

// Receive receives a JSON-RPC request from the transport
func (t *StdioTransport) Receive() (*protocol.JSONRPCRequest, error) {
	if !t.scanner.Scan() {
		if err := t.scanner.Err(); err != nil {
			return nil, fmt.Errorf("scanner error: %w", err)
		}
		return nil, io.EOF
	}

	line := t.scanner.Bytes()
	if len(line) == 0 {
		return nil, fmt.Errorf("empty line received")
	}

	var request protocol.JSONRPCRequest
	if err := json.Unmarshal(line, &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return &request, nil
}

// Close closes the transport
func (t *StdioTransport) Close() error {
	// For stdio transport, we don't actually close stdin/stdout
	return nil
}

// BufferedTransport wraps a transport with buffering capabilities
type BufferedTransport struct {
	transport Transport
	buffer    []byte
	mutex     sync.RWMutex
}

// NewBufferedTransport creates a new buffered transport
func NewBufferedTransport(transport Transport) *BufferedTransport {
	return &BufferedTransport{
		transport: transport,
		buffer:    make([]byte, 0, 4096),
	}
}

// Send sends a response through the buffered transport
func (t *BufferedTransport) Send(response *protocol.JSONRPCResponse) error {
	return t.transport.Send(response)
}

// Receive receives a request from the buffered transport
func (t *BufferedTransport) Receive() (*protocol.JSONRPCRequest, error) {
	return t.transport.Receive()
}

// Close closes the buffered transport
func (t *BufferedTransport) Close() error {
	return t.transport.Close()
}

// MockTransport implements a mock transport for testing
type MockTransport struct {
	requests  chan *protocol.JSONRPCRequest
	responses chan *protocol.JSONRPCResponse
	closed    bool
	mutex     sync.RWMutex
}

// NewMockTransport creates a new mock transport
func NewMockTransport() *MockTransport {
	return &MockTransport{
		requests:  make(chan *protocol.JSONRPCRequest, 10),
		responses: make(chan *protocol.JSONRPCResponse, 10),
	}
}

// Send sends a response through the mock transport
func (t *MockTransport) Send(response *protocol.JSONRPCResponse) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	select {
	case t.responses <- response:
		return nil
	default:
		return fmt.Errorf("response channel is full")
	}
}

// Receive receives a request from the mock transport
func (t *MockTransport) Receive() (*protocol.JSONRPCRequest, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.closed {
		return nil, fmt.Errorf("transport is closed")
	}

	select {
	case request := <-t.requests:
		return request, nil
	default:
		return nil, fmt.Errorf("no requests available")
	}
}

// Close closes the mock transport
func (t *MockTransport) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.closed {
		t.closed = true
		close(t.requests)
		close(t.responses)
	}

	return nil
}

// SendRequest sends a request to the mock transport (for testing)
func (t *MockTransport) SendRequest(request *protocol.JSONRPCRequest) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	select {
	case t.requests <- request:
		return nil
	default:
		return fmt.Errorf("request channel is full")
	}
}

// GetResponse gets a response from the mock transport (for testing)
func (t *MockTransport) GetResponse() (*protocol.JSONRPCResponse, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.closed {
		return nil, fmt.Errorf("transport is closed")
	}

	select {
	case response := <-t.responses:
		return response, nil
	default:
		return nil, fmt.Errorf("no responses available")
	}
}