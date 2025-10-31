package context

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/francknouama/movies-mcp-server/pkg/client"
	"github.com/francknouama/movies-mcp-server/pkg/communication"
	"github.com/francknouama/movies-mcp-server/pkg/protocol"
)

// BDDContext provides a simplified test context for BDD scenarios
// This replaces the complex 1,191-line TestContext with a clean, focused design
type BDDContext struct {
	mcpClient     *client.MCPClient
	serverProcess *exec.Cmd
	testData      map[string]interface{}
	lastResponse  *protocol.ToolCallResponse
	lastError     error
	cleanup       []func() error
}

// NewBDDContext creates a new simplified BDD test context
func NewBDDContext() *BDDContext {
	return &BDDContext{
		testData: make(map[string]interface{}),
		cleanup:  make([]func() error, 0),
	}
}

// SetDatabaseEnvironment sets the database connection string for the MCP server
func (ctx *BDDContext) SetDatabaseEnvironment(connectionString string) error {
	// Parse the connection string to extract individual components
	// Expected format: postgresql://user:password@host:port/dbname?sslmode=disable
	parsed, err := parseConnectionString(connectionString)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Store the parsed components to be used when starting the server
	ctx.SetTestData("db_host", parsed.host)
	ctx.SetTestData("db_port", parsed.port)
	ctx.SetTestData("db_name", parsed.dbname)
	ctx.SetTestData("db_user", parsed.user)
	ctx.SetTestData("db_password", parsed.password)
	ctx.SetTestData("db_sslmode", parsed.sslmode)

	return nil
}

// dbConfig holds parsed database connection components
type dbConfig struct {
	host     string
	port     string
	dbname   string
	user     string
	password string
	sslmode  string
}

// parseConnectionString parses a PostgreSQL connection string
func parseConnectionString(connStr string) (*dbConfig, error) {
	// Simple parsing for PostgreSQL connection strings
	// Format: postgresql://user:password@host:port/dbname?sslmode=disable

	if !strings.HasPrefix(connStr, "postgresql://") && !strings.HasPrefix(connStr, "postgres://") {
		return nil, fmt.Errorf("invalid PostgreSQL connection string format")
	}

	// Remove protocol prefix
	connStr = strings.TrimPrefix(connStr, "postgresql://")
	connStr = strings.TrimPrefix(connStr, "postgres://")

	// Split on '@' to separate user:pass from host:port/db
	parts := strings.Split(connStr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid connection string format")
	}

	userPass := parts[0]
	hostDbQuery := parts[1]

	// Parse user:password
	userParts := strings.Split(userPass, ":")
	if len(userParts) != 2 {
		return nil, fmt.Errorf("invalid user:password format")
	}
	user := userParts[0]
	password := userParts[1]

	// Split query parameters
	hostDbParts := strings.Split(hostDbQuery, "?")
	hostDb := hostDbParts[0]
	sslmode := "disable" // default

	if len(hostDbParts) == 2 {
		// Parse query parameters for sslmode
		queryParams := hostDbParts[1]
		if strings.Contains(queryParams, "sslmode=") {
			for _, param := range strings.Split(queryParams, "&") {
				if strings.HasPrefix(param, "sslmode=") {
					sslmode = strings.TrimPrefix(param, "sslmode=")
					break
				}
			}
		}
	}

	// Parse host:port/dbname
	hostPortDb := strings.Split(hostDb, "/")
	if len(hostPortDb) != 2 {
		return nil, fmt.Errorf("invalid host:port/dbname format")
	}

	hostPort := hostPortDb[0]
	dbname := hostPortDb[1]

	// Parse host:port
	hostPortParts := strings.Split(hostPort, ":")
	if len(hostPortParts) != 2 {
		return nil, fmt.Errorf("invalid host:port format")
	}

	host := hostPortParts[0]
	port := hostPortParts[1]

	return &dbConfig{
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
		sslmode:  sslmode,
	}, nil
}

// getProjectRoot finds the project root directory by looking for go.mod file
func getProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree looking for go.mod file
	for {
		goModPath := filepath.Join(wd, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return wd, nil
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			// Reached filesystem root
			break
		}
		wd = parent
	}

	return "", fmt.Errorf("project root with go.mod file not found")
}

// getServerType determines which server implementation to test
// Returns "sdk" or "legacy" based on TEST_MCP_SERVER environment variable
// Defaults to "legacy" for backwards compatibility
func getServerType() string {
	serverType := os.Getenv("TEST_MCP_SERVER")
	if serverType == "sdk" {
		return "sdk"
	}
	// Default to legacy for backwards compatibility
	return "legacy"
}

// buildServerBinary builds the MCP server binary based on the specified server type
func buildServerBinary(projectRoot, serverType string) (string, error) {
	var serverBinary, serverMainPath string

	switch serverType {
	case "sdk":
		serverBinary = filepath.Join(projectRoot, "movies-mcp-server-sdk")
		serverMainPath = filepath.Join(projectRoot, "cmd", "server-sdk", "main.go")
	case "legacy":
		serverBinary = filepath.Join(projectRoot, "movies-mcp-server")
		serverMainPath = filepath.Join(projectRoot, "cmd", "server", "main.go")
	default:
		return "", fmt.Errorf("unknown server type: %s (expected 'sdk' or 'legacy')", serverType)
	}

	// Check if binary already exists
	if _, err := os.Stat(serverBinary); err == nil {
		return serverBinary, nil
	}

	// Build the server binary
	// #nosec G204 - Safe: building our own Go binary in test environment
	buildCmd := exec.Command("go", "build", "-o", serverBinary, serverMainPath)
	buildCmd.Dir = projectRoot // Set working directory to project root

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to build %s server binary: %w\nOutput: %s", serverType, err, string(output))
	}

	return serverBinary, nil
}

// StartMCPServer starts the real MCP server for testing
func (ctx *BDDContext) StartMCPServer() error {
	// Find project root and build server binary
	projectRoot, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Determine which server to test (sdk or legacy)
	serverType := getServerType()

	serverBinary, err := buildServerBinary(projectRoot, serverType)
	if err != nil {
		return fmt.Errorf("failed to build server binary: %w", err)
	}

	// Start the real MCP server (no mocks - Phase 1 remediation)
	// #nosec G204 - Safe: executing our own built binary in test environment
	ctx.serverProcess = exec.Command(serverBinary)

	// Set database environment if database configuration was provided
	env := os.Environ()
	if host, exists := ctx.GetTestData("db_host"); exists {
		if hostStr, ok := host.(string); ok {
			env = append(env, "DB_HOST="+hostStr)
		}
	}
	if port, exists := ctx.GetTestData("db_port"); exists {
		if portStr, ok := port.(string); ok {
			env = append(env, "DB_PORT="+portStr)
		}
	}
	if name, exists := ctx.GetTestData("db_name"); exists {
		if nameStr, ok := name.(string); ok {
			env = append(env, "DB_NAME="+nameStr)
		}
	}
	if user, exists := ctx.GetTestData("db_user"); exists {
		if userStr, ok := user.(string); ok {
			env = append(env, "DB_USER="+userStr)
		}
	}
	if password, exists := ctx.GetTestData("db_password"); exists {
		if passwordStr, ok := password.(string); ok {
			env = append(env, "DB_PASSWORD="+passwordStr)
		}
	}
	if sslmode, exists := ctx.GetTestData("db_sslmode"); exists {
		if sslmodeStr, ok := sslmode.(string); ok {
			env = append(env, "DB_SSLMODE="+sslmodeStr)
		}
	}

	// Apply the environment to the server process
	if len(env) > len(os.Environ()) {
		ctx.serverProcess.Env = env
	}

	// Get pipes for communication
	stdout, err := ctx.serverProcess.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stdin, err := ctx.serverProcess.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Create stdio transport for MCP communication
	// reader = server's stdout (we read from), writer = server's stdin (we write to)
	transport := communication.NewStdioTransport(stdout, stdin)

	// Create MCP client with proper options
	ctx.mcpClient = client.NewMCPClient(client.ClientOptions{
		Transport: transport,
		Timeout:   30 * time.Second,
		ClientInfo: protocol.ClientInfo{
			Name:    "bdd-test-client",
			Version: "1.0.0",
		},
		Capabilities: protocol.ClientCapabilities{
			Tools:     &protocol.ToolsCapability{},
			Resources: &protocol.ResourcesCapability{},
			Prompts:   &protocol.PromptsCapability{},
		},
	})

	// Start the server process
	err = ctx.serverProcess.Start()
	if err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Wait for server to be ready with retry logic
	time.Sleep(2 * time.Second) // Initial delay for server startup

	// Initialize MCP connection with retry logic
	var initErr error
	for attempts := 0; attempts < 5; attempts++ {
		initErr = ctx.mcpClient.Initialize(
			protocol.ClientInfo{
				Name:    "bdd-test-client",
				Version: "1.0.0",
			},
			protocol.ClientCapabilities{
				Tools:     &protocol.ToolsCapability{},
				Resources: &protocol.ResourcesCapability{},
				Prompts:   &protocol.PromptsCapability{},
			},
		)
		if initErr == nil {
			break // Success!
		}

		// Wait before retry
		time.Sleep(1 * time.Second)
	}

	if initErr != nil {
		return fmt.Errorf("failed to initialize MCP connection after retries: %w", initErr)
	}

	return nil
}

// CallTool calls an MCP tool directly on the real server using correct API
func (ctx *BDDContext) CallTool(toolName string, arguments map[string]interface{}) (*protocol.ToolCallResponse, error) {
	response, err := ctx.mcpClient.CallTool(toolName, arguments)
	ctx.lastResponse = response
	ctx.lastError = err

	return response, err
}

// GetLastResponse returns the last MCP response received
func (ctx *BDDContext) GetLastResponse() *protocol.ToolCallResponse {
	return ctx.lastResponse
}

// GetLastError returns the last error encountered
func (ctx *BDDContext) GetLastError() error {
	return ctx.lastError
}

// SetTestData stores test data by key for use across steps
func (ctx *BDDContext) SetTestData(key string, value interface{}) {
	ctx.testData[key] = value
}

// GetTestData retrieves test data by key
func (ctx *BDDContext) GetTestData(key string) (interface{}, bool) {
	value, exists := ctx.testData[key]
	return value, exists
}

// AddCleanup adds a cleanup function to be executed after the scenario
func (ctx *BDDContext) AddCleanup(fn func() error) {
	ctx.cleanup = append(ctx.cleanup, fn)
}

// Cleanup executes all registered cleanup functions
func (ctx *BDDContext) Cleanup() error {
	var errors []error

	// Execute cleanup functions in reverse order
	for i := len(ctx.cleanup) - 1; i >= 0; i-- {
		if err := ctx.cleanup[i](); err != nil {
			errors = append(errors, err)
		}
	}

	// Close MCP client
	if ctx.mcpClient != nil {
		if err := ctx.mcpClient.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	// Stop server process
	if ctx.serverProcess != nil && ctx.serverProcess.Process != nil {
		if err := ctx.serverProcess.Process.Kill(); err != nil {
			errors = append(errors, err)
		}
		if err := ctx.serverProcess.Wait(); err != nil {
			errors = append(errors, err)
		}
	}

	// Clear test data
	ctx.testData = make(map[string]interface{})
	ctx.cleanup = make([]func() error, 0)
	ctx.lastResponse = nil
	ctx.lastError = nil

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}
	return nil
}

// WaitForServer waits for the MCP server to be ready
func (ctx *BDDContext) WaitForServer(timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		if ctx.mcpClient != nil {
			// Try a simple tools list to see if server is responsive
			_, err := ctx.mcpClient.ListTools()
			if err == nil {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not ready after %v", timeout)
}

// ParseJSONResponse parses the last response as JSON into the provided structure
func (ctx *BDDContext) ParseJSONResponse(target interface{}) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response available to parse")
	}

	if len(ctx.lastResponse.Content) == 0 {
		return fmt.Errorf("response content is empty")
	}

	// Get the first content block and parse as JSON
	content := ctx.lastResponse.Content[0]
	if content.Type != "text" {
		return fmt.Errorf("expected text content, got %s", content.Type)
	}

	err := json.Unmarshal([]byte(content.Text), target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response into target: %w", err)
	}

	return nil
}

// HasError returns true if the last response contained an error
func (ctx *BDDContext) HasError() bool {
	return ctx.lastError != nil || (ctx.lastResponse != nil && ctx.lastResponse.IsError)
}

// GetErrorMessage returns the error message from the last response or error
func (ctx *BDDContext) GetErrorMessage() string {
	if ctx.lastError != nil {
		return ctx.lastError.Error()
	}
	if ctx.lastResponse != nil && ctx.lastResponse.IsError {
		return fmt.Sprintf("Tool call failed: %v", ctx.lastResponse.Content)
	}
	return ""
}

// GetMCPClient returns the MCP client for direct access
func (ctx *BDDContext) GetMCPClient() *client.MCPClient {
	return ctx.mcpClient
}
