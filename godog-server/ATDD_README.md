# ATDD (Acceptance Test-Driven Development) for Movies MCP Server

This directory contains comprehensive ATDD tests using Cucumber/Godog for the Movies MCP Server.

## Features

- **307 step definitions** covering all MCP server functionality
- **39 test scenarios** across 4 feature areas:
  - MCP Protocol Communication
  - Movie Operations (CRUD, search, validation)
  - Actor Operations (CRUD, relationships)
  - Advanced Search and Integration
- **Dual testing modes**: Mock server and real server with database
- **Docker integration** for isolated database testing

## Test Structure

```
godog-server/
├── examples/features/           # Gherkin feature files
│   ├── mcp_protocol.feature    # MCP protocol tests
│   ├── movie_operations.feature # Movie CRUD operations
│   ├── actor_operations.feature # Actor operations and relationships
│   └── advanced_search.feature  # Complex search scenarios
├── step_definitions/           # Go step implementations
│   ├── mcp_protocol_steps.go   # MCP protocol step definitions
│   ├── movie_steps.go          # Movie operation step definitions
│   ├── actor_steps.go          # Actor operation step definitions
│   ├── test_context.go         # Test context and server communication
│   └── test_helpers.go         # Utility functions
├── docker-compose.test.yml     # Test database configuration
├── main_test.go               # Mock server tests
├── integration_test.go        # Real server tests
└── Makefile                   # Test automation
```

## Running Tests

### Option 1: Mock Server Tests (Fast)
Run tests with mocked MCP server responses - great for step definition validation:

```bash
# Run all tests with mock server
make test

# Or directly
go test -v
```

### Option 2: Real Server Tests (Complete)
Run tests with real MCP server and PostgreSQL database:

```bash
# Start test database
make test-db-up

# Run integration tests with real server
export DB_HOST=localhost
export DB_PORT=5433
export DB_USER=movies_user
export DB_PASSWORD=movies_password
export DB_NAME=movies_mcp_test
export DB_SSLMODE=disable
export USE_REAL_SERVER=true
go test -v -tags=integration

# Stop test database
make test-db-down
```

### Using Make Commands

```bash
# Start test database
make test-db-up

# Stop test database  
make test-db-down

# Clean database (remove volumes)
make test-db-clean

# Run specific scenario
make test-scenario SCENARIO="Initialize_MCP_connection"
```

## Environment Variables

For real server tests, configure these environment variables:

- `USE_REAL_SERVER=true` - Enable real MCP server mode
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5433)
- `DB_USER` - Database user (default: movies_user)
- `DB_PASSWORD` - Database password (default: movies_password)
- `DB_NAME` - Database name (default: movies_mcp_test)
- `DB_SSLMODE` - SSL mode (default: disable)

## Test Scenarios

### MCP Protocol
- Initialize MCP connection
- List available tools and resources
- Handle invalid methods and protocol versions

### Movie Operations
- Add, get, update, delete movies
- Search movies by title, director, rating
- Get top-rated movies
- Handle validation errors

### Actor Operations  
- Add, get, update, delete actors
- Link/unlink actors to movies
- Get movie cast and actor filmography
- Search actors by name and birth year

### Advanced Features
- Complex search queries
- Performance testing with large datasets
- Concurrent operations
- Data integrity validation

## Example Test Run

```bash
$ make test-db-up
Starting test database...
✅ Database is ready!

$ export USE_REAL_SERVER=true && go test -v -tags=integration -run "Initialize_MCP_connection"
=== RUN   TestFeaturesIntegration/Initialize_MCP_connection
Starting real MCP server...
[MCP Server STDERR]: Connected to database: movies_mcp_test
[MCP Server STDERR]: Starting Movies MCP Server with Clean Architecture...
✅ Scenario passed: Initialize MCP connection

$ make test-db-down
Stopping test database...
```

## Benefits

1. **Living Documentation**: Feature files serve as executable specifications
2. **Behavior-Driven Development**: Tests describe system behavior from user perspective
3. **Integration Testing**: Validates real MCP server with database
4. **Regression Prevention**: Comprehensive test coverage prevents breaking changes
5. **CI/CD Ready**: Docker-based setup works in any environment

## Dependencies

- Go 1.24+
- Docker and Docker Compose
- Godog (Cucumber for Go)
- PostgreSQL (via Docker)

The ATDD framework provides confidence that the Movies MCP Server works correctly end-to-end, from MCP protocol communication to database operations.