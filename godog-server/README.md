# Godog MCP Server

A Model Context Protocol (MCP) server that enables AI assistants to run and manage Cucumber/Godog BDD tests.

## Features

### Phase 1 (Current)
- ✅ MCP protocol implementation with initialize/list_tools handlers
- ✅ Godog CLI detection and validation
- ✅ Feature file parsing and validation
- ✅ Core domain models (Feature, Scenario, TestResult)
- ✅ Configuration management with environment variables
- ✅ Structured logging and error handling

### Planned Features
- **Phase 2**: Feature file management (list, read, validate)
- **Phase 3**: Test execution (run suite, run feature, run scenario)
- **Phase 4**: Reporting (JSON, HTML, JUnit formats)
- **Phase 5**: Step definition management
- **Phase 6**: Test data and environment management
- **Phase 7**: Production readiness (monitoring, metrics)

## Quick Start

### Prerequisites
- Go 1.24.4 or later
- Godog CLI installed (`go install github.com/cucumber/godog/cmd/godog@latest`)

### Installation

```bash
# Clone the workspace repository
git clone https://github.com/francknouama/movies-mcp-server.git
cd movies-mcp-server/godog-server

# Build the server
go build -o godog-server cmd/server/main.go

# Run the server
./godog-server
```

### Configuration

Configure via environment variables:

```bash
# Server configuration
export LOG_LEVEL=info                    # Log level (debug, info, warn, error)
export SERVER_TIMEOUT=30s                # Server timeout

# Godog configuration
export GODOG_BINARY=godog                # Path to godog binary
export FEATURES_DIR=./features           # Feature files directory
export STEP_DEFS_DIR=./step_definitions  # Step definitions directory
export REPORTS_DIR=./reports             # Reports output directory

# Test execution
export MAX_PARALLEL=4                    # Max parallel test execution
export DEFAULT_TIMEOUT=5m                # Default test timeout
export RETRY_COUNT=0                     # Test retry count
```

## Usage

### Initialize Connection
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "example-client",
      "version": "1.0.0"
    }
  }
}
```

### List Available Tools
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### Validate Feature File
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "validate_feature",
    "arguments": {
      "file_path": "features/example.feature"
    }
  }
}
```

## Development

### Project Structure
```
godog-server/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── godog/          # Godog runner and integration
│   ├── models/         # Domain models
│   └── server/         # MCP server implementation
└── pkg/
    ├── errors/         # Custom error types
    └── logging/        # Logging utilities
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/godog/...
```

### Building
```bash
# Build for current platform
go build -o godog-server cmd/server/main.go

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o godog-server cmd/server/main.go
```

## MCP Protocol

The server implements the Model Context Protocol specification (version 2024-11-05).

### Supported Methods
- `initialize` - Initialize server connection
- `tools/list` - List available tools
- `tools/call` - Execute a tool
- `resources/list` - List available resources
- `resources/read` - Read resource content

### Available Tools
- `validate_feature` - Parse and validate Gherkin feature files
- `list_features` - List all available feature files

### Available Resources
- `godog://features/all` - All feature files
- `godog://reports/latest` - Latest test results

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.