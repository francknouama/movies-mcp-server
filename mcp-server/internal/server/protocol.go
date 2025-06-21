package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

// Protocol handles the low-level MCP JSON-RPC communication
type Protocol struct {
	input  io.Reader
	output io.Writer
	logger *log.Logger
}

// NewProtocol creates a new protocol handler
func NewProtocol(input io.Reader, output io.Writer, logger *log.Logger) *Protocol {
	return &Protocol{
		input:  input,
		output: output,
		logger: logger,
	}
}

// Listen starts listening for incoming requests and routes them to the handler
func (p *Protocol) Listen(handler RequestHandler) error {
	if p.logger != nil {
		p.logger.Println("Starting MCP Protocol listener...")
	}

	scanner := bufio.NewScanner(p.input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if p.logger != nil {
			p.logger.Printf("Received: %s", line)
		}

		var request dto.JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			p.SendError(nil, dto.ParseError, "Parse error", err.Error())
			continue
		}

		handler.HandleRequest(&request, p)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// SendResult sends a successful JSON-RPC response
func (p *Protocol) SendResult(id any, result any) {
	response := dto.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	p.sendResponse(response)
}

// SendError sends an error JSON-RPC response
func (p *Protocol) SendError(id any, code int, message string, data interface{}) {
	response := dto.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &dto.JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	p.sendResponse(response)
}

// sendResponse sends a JSON-RPC response
func (p *Protocol) sendResponse(response dto.JSONRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		if p.logger != nil {
			p.logger.Printf("Failed to marshal response: %v", err)
		}
		return
	}

	if p.logger != nil {
		p.logger.Printf("Sending: %s", string(data))
	}

	p.output.Write(data)
	p.output.Write([]byte("\n"))
}

// UnmarshalParams unmarshals request parameters
func (p *Protocol) UnmarshalParams(params interface{}, target interface{}) error {
	if params == nil {
		return fmt.Errorf("missing parameters")
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// RequestHandler defines the interface for handling MCP requests
type RequestHandler interface {
	HandleRequest(req *dto.JSONRPCRequest, protocol *Protocol)
}

// ResponseSender provides methods for sending responses
type ResponseSender interface {
	SendResult(id any, result any)
	SendError(id any, code int, message string, data interface{})
}
