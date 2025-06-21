# Workspace-level Makefile for MCP Servers

.PHONY: help build-all test-all clean-all movies-server godog-server

# Default target
help:
	@echo "MCP Servers Workspace Commands:"
	@echo "  make build-all      - Build all servers"
	@echo "  make test-all       - Run tests for all servers"
	@echo "  make clean-all      - Clean all build artifacts"
	@echo ""
	@echo "Individual server commands:"
	@echo "  make movies-mcp-server  - Build movies MCP server"
	@echo "  make godog-server       - Build godog server"
	@echo ""
	@echo "For server-specific commands, use:"
	@echo "  cd movies-mcp-server && make help"
	@echo "  cd godog-server && make help"

# Build all servers
build-all: movies-mcp-server godog-server
	@echo "✅ All servers built successfully"

# Build individual servers
movies-mcp-server:
	@echo "Building movies-mcp-server..."
	@cd movies-mcp-server && go build -o movies-mcp-server cmd/server/main.go
	@echo "✅ movies-mcp-server built"

godog-server:
	@echo "Building godog-server..."
	@cd godog-server && go build -o godog-server cmd/server/main.go
	@echo "✅ godog-server built"

# Test all servers
test-all:
	@echo "Running tests for all servers..."
	@echo "\n📦 Testing movies-mcp-server..."
	@cd movies-mcp-server && go test ./...
	@echo "\n📦 Testing godog-server..."
	@cd godog-server && go test ./...
	@echo "\n✅ All tests completed"

# Clean all build artifacts
clean-all:
	@echo "Cleaning all build artifacts..."
	@cd movies-mcp-server && rm -f movies-mcp-server movies-server server *.test
	@cd godog-server && rm -f godog-server *.test
	@echo "✅ Clean completed"

# Verify workspace setup
verify-workspace:
	@echo "Verifying workspace setup..."
	@go work sync
	@echo "✅ Workspace is properly configured"

# Initialize workspace (useful after cloning)
init-workspace:
	@echo "Initializing workspace..."
	@go work sync
	@cd movies-mcp-server && go mod download
	@cd godog-server && go mod download
	@echo "✅ Workspace initialized"