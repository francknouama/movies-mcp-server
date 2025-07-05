# Movies MCP Server Makefile

# Variables
BINARY_NAME=movies-server
BINARY_NAME_CLEAN=movies-server-clean
GO_VERSION=1.23
MAIN_PATH=cmd/server/main.go
MAIN_PATH_CLEAN=cmd/server-new/main.go
MIGRATE_PATH=tools/migrate/main.go
BUILD_DIR=build
DOCKER_IMAGE=movies-mcp-server
DOCKER_IMAGE_CLEAN=movies-mcp-server-clean
DOCKER_TAG=latest

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

# Colors for output
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: all build clean test run fmt vet lint deps help

# Default target
all: clean build test

# Build the binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "$(GREEN)Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-darwin:
	@echo "$(GREEN)Building for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows:
	@echo "$(GREEN)Building for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run the application
run: build
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f movies-server
	@echo "$(GREEN)Clean complete$(NC)"

# Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GOTEST) -v ./...

# Run integration tests with testcontainers
test-integration:
	@echo "$(GREEN)Running integration tests with testcontainers...$(NC)"
	@$(GOTEST) -v -tags=integration ./internal/infrastructure/postgres/...

# Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Run integration tests with coverage
test-integration-coverage:
	@echo "$(GREEN)Running integration tests with coverage...$(NC)"
	@$(GOTEST) -v -tags=integration -coverprofile=coverage-integration.out ./internal/infrastructure/postgres/...
	@$(GOCMD) tool cover -html=coverage-integration.out -o coverage-integration.html
	@echo "$(GREEN)Integration coverage report generated: coverage-integration.html$(NC)"

# Run the basic init test
test-init: build
	@echo "$(GREEN)Testing MCP initialization...$(NC)"
	@./test_init.sh

# Run all integration tests
test-all: build
	@echo "$(GREEN)Running all tests...$(NC)"
	@./test_all.sh

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GOFMT) ./...

# Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GOVET) ./...

# Run linter (requires golangci-lint)
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with: brew install golangci-lint$(NC)"; \
	fi

# Check code quality (fmt, vet, lint)
check: fmt vet lint

# Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy

# Update dependencies
deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

# Clean Architecture Targets
build-clean:
	@echo "$(GREEN)Building $(BINARY_NAME_CLEAN) (Clean Architecture)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_CLEAN) $(MAIN_PATH_CLEAN)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME_CLEAN)$(NC)"

build-migrate:
	@echo "$(GREEN)Building migration tool...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/migrate $(MIGRATE_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/migrate$(NC)"

run-clean: build-clean
	@echo "$(GREEN)Running $(BINARY_NAME_CLEAN) (Clean Architecture)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME_CLEAN)

# Docker targets
docker-build:
	@echo "$(GREEN)Building Docker image (Legacy)...$(NC)"
	@docker build -f Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-build-clean:
	@echo "$(GREEN)Building Docker image (Clean Architecture)...$(NC)"
	@docker build -f Dockerfile.clean -t $(DOCKER_IMAGE_CLEAN):$(DOCKER_TAG) .

docker-run:
	@echo "$(GREEN)Running Docker container (Legacy)...$(NC)"
	@docker run -it --rm $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-clean:
	@echo "$(GREEN)Running Docker container (Clean Architecture)...$(NC)"
	@docker run -it --rm $(DOCKER_IMAGE_CLEAN):$(DOCKER_TAG)

# Docker Compose targets
docker-compose-up:
	@echo "$(GREEN)Starting services with docker-compose (Legacy)...$(NC)"
	@docker-compose up -d

docker-compose-up-clean:
	@echo "$(GREEN)Starting services with docker-compose (Clean Architecture)...$(NC)"
	@docker-compose -f docker-compose.clean.yml up -d --build

docker-compose-up-dev:
	@echo "$(GREEN)Starting development services...$(NC)"
	@docker-compose -f docker-compose.dev.yml up -d

docker-compose-down:
	@echo "$(GREEN)Stopping services (Legacy)...$(NC)"
	@docker-compose down

docker-compose-down-clean:
	@echo "$(GREEN)Stopping services (Clean Architecture)...$(NC)"
	@docker-compose -f docker-compose.clean.yml down

docker-compose-down-dev:
	@echo "$(GREEN)Stopping development services...$(NC)"
	@docker-compose -f docker-compose.dev.yml down

docker-compose-logs:
	@echo "$(GREEN)Showing logs (Clean Architecture)...$(NC)"
	@docker-compose -f docker-compose.clean.yml logs -f

docker-compose-logs-dev:
	@echo "$(GREEN)Showing development logs...$(NC)"
	@docker-compose -f docker-compose.dev.yml logs -f

# Migration tool installation
install-migrate:
	@echo "$(GREEN)Installing golang-migrate CLI tool...$(NC)"
	@if ! command -v migrate >/dev/null; then \
		echo "Installing golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	else \
		echo "golang-migrate already installed"; \
	fi

# Database targets
db-setup: install-migrate
	@echo "$(GREEN)Setting up database...$(NC)"
	@./scripts/setup_db.sh

db-migrate: install-migrate
	@echo "$(GREEN)Running database migrations...$(NC)"
	@./scripts/migrate.sh

db-migrate-down: install-migrate
	@echo "$(GREEN)Rolling back migrations...$(NC)"
	@./scripts/migrate_down.sh

db-migrate-reset: install-migrate
	@echo "$(GREEN)Resetting database...$(NC)"
	@./scripts/migrate_reset.sh

db-migrate-version: install-migrate
	@echo "$(GREEN)Checking migration version...$(NC)"
	@./scripts/migrate_version.sh

db-seed: db-migrate
	@echo "$(GREEN)Seeding database...$(NC)"
	@./scripts/seed.sh

# Development helpers
dev: deps build
	@echo "$(GREEN)Development environment ready$(NC)"

# Install the binary to $GOPATH/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

# Create a release
release: clean build-all
	@echo "$(GREEN)Creating release artifacts...$(NC)"
	@mkdir -p releases
	@tar -czf releases/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	@tar -czf releases/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	@tar -czf releases/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	@zip -j releases/$(BINARY_NAME)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "$(GREEN)Release artifacts created in releases/$(NC)"

# Help target
help:
	@echo "$(GREEN)Movies MCP Server - Available targets:$(NC)"
	@echo ""
	@echo "$(YELLOW)Basic Targets:$(NC)"
	@echo "  $(YELLOW)make$(NC)              - Build and test (default)"
	@echo "  $(YELLOW)make build$(NC)        - Build the binary (legacy)"
	@echo "  $(YELLOW)make build-clean$(NC)  - Build clean architecture binary"
	@echo "  $(YELLOW)make build-migrate$(NC) - Build migration tool"
	@echo "  $(YELLOW)make run$(NC)          - Build and run the server (legacy)"
	@echo "  $(YELLOW)make run-clean$(NC)    - Build and run clean architecture"
	@echo "  $(YELLOW)make clean$(NC)        - Clean build artifacts"
	@echo ""
	@echo "$(YELLOW)Testing:$(NC)"
	@echo "  $(YELLOW)make test$(NC)         - Run unit tests"
	@echo "  $(YELLOW)make test-integration$(NC) - Run integration tests with testcontainers"
	@echo "  $(YELLOW)make test-coverage$(NC) - Run tests with coverage"
	@echo "  $(YELLOW)make test-integration-coverage$(NC) - Run integration tests with coverage"
	@echo "  $(YELLOW)make test-init$(NC)    - Test MCP initialization"
	@echo "  $(YELLOW)make test-all$(NC)     - Run all integration tests"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  $(YELLOW)make fmt$(NC)          - Format code"
	@echo "  $(YELLOW)make vet$(NC)          - Run go vet"
	@echo "  $(YELLOW)make lint$(NC)         - Run linter"
	@echo "  $(YELLOW)make check$(NC)        - Run all code quality checks"
	@echo ""
	@echo "$(YELLOW)Dependencies:$(NC)"
	@echo "  $(YELLOW)make deps$(NC)         - Download dependencies"
	@echo "  $(YELLOW)make deps-update$(NC)  - Update dependencies"
	@echo ""
	@echo "$(YELLOW)Docker (Legacy):$(NC)"
	@echo "  $(YELLOW)make docker-build$(NC) - Build Docker image (legacy)"
	@echo "  $(YELLOW)make docker-run$(NC)   - Run Docker container (legacy)"
	@echo "  $(YELLOW)make docker-compose-up$(NC) - Start legacy services"
	@echo "  $(YELLOW)make docker-compose-down$(NC) - Stop legacy services"
	@echo ""
	@echo "$(YELLOW)Docker (Clean Architecture):$(NC)"
	@echo "  $(YELLOW)make docker-build-clean$(NC) - Build clean architecture image"
	@echo "  $(YELLOW)make docker-run-clean$(NC) - Run clean architecture container"
	@echo "  $(YELLOW)make docker-compose-up-clean$(NC) - Start clean architecture stack"
	@echo "  $(YELLOW)make docker-compose-down-clean$(NC) - Stop clean architecture stack"
	@echo "  $(YELLOW)make docker-compose-logs$(NC) - Show clean architecture logs"
	@echo ""
	@echo "$(YELLOW)Docker (Development):$(NC)"
	@echo "  $(YELLOW)make docker-compose-up-dev$(NC) - Start development databases"
	@echo "  $(YELLOW)make docker-compose-down-dev$(NC) - Stop development databases"
	@echo "  $(YELLOW)make docker-compose-logs-dev$(NC) - Show development logs"
	@echo ""
	@echo "$(YELLOW)Database:$(NC)"
	@echo "  $(YELLOW)make db-setup$(NC)     - Set up database"
	@echo "  $(YELLOW)make db-migrate$(NC)   - Run migrations"
	@echo "  $(YELLOW)make db-migrate-down$(NC) - Rollback migrations"
	@echo "  $(YELLOW)make db-migrate-reset$(NC) - Reset database"
	@echo "  $(YELLOW)make db-seed$(NC)      - Seed database with sample data"
	@echo ""
	@echo "$(YELLOW)Release:$(NC)"
	@echo "  $(YELLOW)make install$(NC)      - Install binary to GOPATH"
	@echo "  $(YELLOW)make release$(NC)      - Create release artifacts"
	@echo "  $(YELLOW)make help$(NC)         - Show this help"

# Version info
version:
	@echo "$(GREEN)Movies MCP Server$(NC)"
	@echo "Legacy Version: 0.1.0"
	@echo "Clean Architecture Version: 0.2.0"
	@echo "Go Version Required: $(GO_VERSION)"