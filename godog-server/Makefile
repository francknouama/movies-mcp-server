# Godog MCP Server Makefile

.PHONY: help build test clean run install-godog test-setup test-teardown test-db-up test-db-down test-atdd

# Default target
help:
	@echo "Godog MCP Server Commands:"
	@echo "  make build         - Build the server"
	@echo "  make test          - Run tests (mocked)"
	@echo "  make test-atdd     - Run ATDD tests with real database"
	@echo "  make test-db-up    - Start test database"
	@echo "  make test-db-down  - Stop test database"
	@echo "  make test-setup    - Set up test environment"
	@echo "  make test-teardown - Tear down test environment"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make run           - Build and run the server"
	@echo "  make install-godog - Install Godog CLI"

# Build the server
build:
	@echo "Building godog-server..."
	@go build -o godog-server cmd/server/main.go
	@echo "✅ Build complete: ./godog-server"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f godog-server coverage.out coverage.html
	@echo "✅ Clean complete"

# Build and run the server
run: build
	@echo "Starting godog-server..."
	@./godog-server

# Install Godog CLI
install-godog:
	@echo "Installing Godog CLI..."
	@go install github.com/cucumber/godog/cmd/godog@latest
	@echo "✅ Godog installed to $$GOPATH/bin/godog"

# Verify Godog installation
verify-godog:
	@echo "Checking Godog installation..."
	@which godog > /dev/null && godog --version || echo "❌ Godog not found. Run 'make install-godog'"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./... || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

# Database setup for ATDD tests
test-db-up:
	@echo "Starting test database..."
	docker-compose -f docker-compose.test.yml up -d postgres-test
	@echo "Waiting for database to be ready..."
	@timeout 30 sh -c 'until docker-compose -f docker-compose.test.yml exec postgres-test pg_isready -U movies_user -d movies_mcp_test; do sleep 1; done'
	@echo "✅ Database is ready!"

test-db-down:
	@echo "Stopping test database..."
	docker-compose -f docker-compose.test.yml down

test-db-clean:
	@echo "Cleaning test database..."
	docker-compose -f docker-compose.test.yml down -v

# Test setup and teardown
test-setup: test-db-up
	@echo "Setting up test environment..."

test-teardown: test-db-down
	@echo "Tearing down test environment..."

# Run ATDD tests with real database
test-atdd: test-setup
	@echo "Running ATDD tests with real database..."
	@export DB_HOST=localhost && \
	export DB_PORT=5433 && \
	export DB_USER=movies_user && \
	export DB_PASSWORD=movies_password && \
	export DB_NAME=movies_mcp_test && \
	export DB_SSLMODE=disable && \
	export USE_REAL_SERVER=true && \
	go test -v -tags=integration
	@$(MAKE) test-teardown

# Run specific test scenario with real database
test-scenario: test-setup
	@echo "Running specific test scenario: $(SCENARIO)"
	@export DB_HOST=localhost && \
	export DB_PORT=5433 && \
	export DB_USER=movies_user && \
	export DB_PASSWORD=movies_password && \
	export DB_NAME=movies_mcp_test && \
	export DB_SSLMODE=disable && \
	export USE_REAL_SERVER=true && \
	go test -v -run "$(SCENARIO)" -tags=integration