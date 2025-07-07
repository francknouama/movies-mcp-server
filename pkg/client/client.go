package client

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/francknouama/movies-mcp-server/pkg/communication"
	"github.com/francknouama/movies-mcp-server/pkg/protocol"
)

// MCPClient represents an MCP protocol client
type MCPClient struct {
	transport    communication.Transport
	initialized  bool
	capabilities *protocol.ServerCapabilities
	serverInfo   *protocol.ServerInfo
	requestID    int64
	mutex        sync.RWMutex
	timeout      time.Duration
}

// ClientOptions represents options for creating an MCP client
type ClientOptions struct {
	Transport    communication.Transport
	Timeout      time.Duration
	ClientInfo   protocol.ClientInfo
	Capabilities protocol.ClientCapabilities
}

// NewMCPClient creates a new MCP client
func NewMCPClient(options ClientOptions) *MCPClient {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &MCPClient{
		transport: options.Transport,
		timeout:   timeout,
		requestID: 1,
	}
}

// NewStdioMCPClient creates a new MCP client using stdin/stdout transport
func NewStdioMCPClient(reader io.Reader, writer io.Writer) *MCPClient {
	transport := communication.NewStdioTransport(reader, writer)
	return NewMCPClient(ClientOptions{
		Transport: transport,
		ClientInfo: protocol.ClientInfo{
			Name:    "shared-mcp-client",
			Version: "1.0.0",
		},
		Capabilities: protocol.ClientCapabilities{
			Tools:     &protocol.ToolsCapability{},
			Resources: &protocol.ResourcesCapability{},
			Prompts:   &protocol.PromptsCapability{},
		},
	})
}

// Initialize initializes the MCP connection
func (c *MCPClient) Initialize(clientInfo protocol.ClientInfo, capabilities protocol.ClientCapabilities) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.initialized {
		return fmt.Errorf("client already initialized")
	}

	request := &protocol.JSONRPCRequest{
		JSONRPC: protocol.JSONRPC2Version,
		ID:      c.nextRequestID(),
		Method:  protocol.MethodInitialize,
		Params: c.marshalParams(protocol.InitializeRequest{
			ProtocolVersion: protocol.MCPVersion,
			Capabilities:    capabilities,
			ClientInfo:      clientInfo,
		}),
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return fmt.Errorf("initialize request failed: %w", err)
	}

	if response.Error != nil {
		return fmt.Errorf("initialize error: %s", response.Error.Message)
	}

	var initResponse protocol.InitializeResponse
	if err := c.unmarshalResult(response.Result, &initResponse); err != nil {
		return fmt.Errorf("failed to unmarshal initialize response: %w", err)
	}

	c.capabilities = &initResponse.Capabilities
	c.serverInfo = &initResponse.ServerInfo
	c.initialized = true

	return nil
}

// CallTool calls an MCP tool
func (c *MCPClient) CallTool(name string, arguments map[string]interface{}) (*protocol.ToolCallResponse, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	request := &protocol.JSONRPCRequest{
		JSONRPC: protocol.JSONRPC2Version,
		ID:      c.nextRequestID(),
		Method:  protocol.MethodToolsCall,
		Params: c.marshalParams(protocol.ToolCallRequest{
			Name:      name,
			Arguments: arguments,
		}),
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("tool call error: %s", response.Error.Message)
	}

	var toolResponse protocol.ToolCallResponse
	if err := c.unmarshalResult(response.Result, &toolResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool response: %w", err)
	}

	return &toolResponse, nil
}

// ListTools lists available tools
func (c *MCPClient) ListTools() (*protocol.ToolsListResponse, error) {
	var toolsResponse protocol.ToolsListResponse
	if err := c.listResource(protocol.MethodToolsList, &toolsResponse, "list tools"); err != nil {
		return nil, err
	}
	return &toolsResponse, nil
}

// ListResources lists available resources
func (c *MCPClient) ListResources() (*protocol.ResourcesListResponse, error) {
	var resourcesResponse protocol.ResourcesListResponse
	if err := c.listResource(protocol.MethodResourcesList, &resourcesResponse, "list resources"); err != nil {
		return nil, err
	}
	return &resourcesResponse, nil
}

// ReadResource reads a resource
func (c *MCPClient) ReadResource(uri string) (*protocol.ResourceReadResponse, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	request := &protocol.JSONRPCRequest{
		JSONRPC: protocol.JSONRPC2Version,
		ID:      c.nextRequestID(),
		Method:  protocol.MethodResourcesRead,
		Params: c.marshalParams(protocol.ResourceReadRequest{
			URI: uri,
		}),
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("read resource failed: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("read resource error: %s", response.Error.Message)
	}

	var resourceResponse protocol.ResourceReadResponse
	if err := c.unmarshalResult(response.Result, &resourceResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource response: %w", err)
	}

	return &resourceResponse, nil
}

// Close closes the MCP client
func (c *MCPClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.transport != nil {
		err := c.transport.Close()
		c.transport = nil
		return err
	}

	return nil
}

// GetServerCapabilities returns the server capabilities
func (c *MCPClient) GetServerCapabilities() *protocol.ServerCapabilities {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.capabilities
}

// GetServerInfo returns the server info
func (c *MCPClient) GetServerInfo() *protocol.ServerInfo {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.serverInfo
}

// IsInitialized returns whether the client is initialized
func (c *MCPClient) IsInitialized() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.initialized
}

// Helper methods

func (c *MCPClient) nextRequestID() interface{} {
	id := c.requestID
	c.requestID++
	return id
}

func (c *MCPClient) marshalParams(params interface{}) json.RawMessage {
	data, err := json.Marshal(params)
	if err != nil {
		// Return empty JSON object if marshaling fails
		return []byte("{}")
	}
	return data
}

func (c *MCPClient) unmarshalResult(result interface{}, target interface{}) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (c *MCPClient) sendRequest(request *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	// Send request
	if err := c.transport.Send(&protocol.JSONRPCResponse{
		JSONRPC: request.JSONRPC,
		ID:      request.ID,
		Result:  request,
	}); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Receive response
	responseReq, err := c.transport.Receive()
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	// Convert to response (this is a simplification - in real implementation,
	// we'd need proper request/response correlation)
	response := &protocol.JSONRPCResponse{
		JSONRPC: responseReq.JSONRPC,
		ID:      responseReq.ID,
		Result:  responseReq.Params,
	}

	return response, nil
}

// listResource is a helper function to reduce code duplication in list operations
func (c *MCPClient) listResource(method string, result interface{}, errorPrefix string) error {
	if !c.initialized {
		return fmt.Errorf("client not initialized")
	}

	request := &protocol.JSONRPCRequest{
		JSONRPC: protocol.JSONRPC2Version,
		ID:      c.nextRequestID(),
		Method:  method,
		Params:  c.marshalParams(protocol.ListRequest{}),
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return fmt.Errorf("%s failed: %w", errorPrefix, err)
	}

	if response.Error != nil {
		return fmt.Errorf("%s error: %s", errorPrefix, response.Error.Message)
	}

	if err := c.unmarshalResult(response.Result, result); err != nil {
		return fmt.Errorf("failed to unmarshal %s response: %w", errorPrefix, err)
	}

	return nil
}
