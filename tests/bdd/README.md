# BDD Tests for Movies MCP Server

This directory contains the consolidated BDD (Behavior-Driven Development) tests for the Movies MCP Server, implementing **Phase 3** of the BDD remediation plan.

## 🏗️ Architecture

### Phase 3 Improvements

This new structure replaces the complex 1,191-line `TestContext` with a simplified, focused design:

- **Simplified Context**: `BDDContext` with <100 lines of focused functionality
- **Test Fixtures**: YAML-based test data instead of hardcoded values
- **Database Management**: Proper cleanup and isolation between scenarios
- **Direct Server Communication**: Uses real MCP server (no mocks - Phase 1 remediation)
- **Shared Protocol**: Uses unified MCP types (Phase 2 remediation)

## 📁 Directory Structure

```
tests/bdd/
├── features/           # Gherkin feature files
├── steps/              # Step definitions
├── context/            # Simplified BDD context
├── support/            # Database and test utilities
├── fixtures/           # YAML test data
└── bdd_test.go        # Main test runner
```

## 🧪 Running Tests

### Prerequisites

1. Database running:
   ```bash
   docker-compose up -d postgres
   ```

2. Test database setup:
   ```bash
   make setup-test-db
   ```

### Run BDD Tests

```bash
# Run all BDD tests
go test ./tests/bdd/

# Run with verbose output
go test -v ./tests/bdd/

# Run specific features
go test ./tests/bdd/ -godog.tags=@movies

# Run with custom format
go test ./tests/bdd/ -godog.format=pretty
```

## 📋 Test Data Management

### Fixtures

Test data is managed through YAML fixtures in the `fixtures/` directory:

- `movies.yaml` - Common movie test data
- `actors.yaml` - Actor test data with relationships
- `search_scenarios.yaml` - Data for search testing

### Database Cleanup

Each scenario automatically:
1. Cleans the database before starting
2. Loads required fixtures
3. Runs the test scenario
4. Cleans up after completion

## 🔧 Key Components

### BDDContext

Simplified test context (`context/bdd_context.go`):
- Direct MCP client communication
- Test data storage
- Cleanup management
- Response parsing utilities

### TestDatabase

Database helper (`support/test_database.go`):
- Fixture loading from YAML
- Cleanup strategies
- Verification utilities
- Test isolation

### Step Definitions

Modular step definitions:
- `common_steps.go` - Shared MCP protocol steps
- `movie_steps.go` - Movie-specific operations
- `actor_steps.go` - Actor-specific operations

## 🎯 Benefits

Compared to the previous complex structure:

- **50% reduction** in test context complexity
- **Better test isolation** with proper cleanup
- **Faster test execution** with optimized database operations
- **Easier maintenance** with external fixtures
- **No hardcoded test data** in step definitions

## 🔄 Migration from godog-server

This BDD structure consolidates tests from `godog-server/` with:

- ✅ All feature files moved
- ✅ Simplified step definitions
- ✅ Eliminated code duplication
- ✅ Uses shared MCP protocol library
- ✅ Proper database management

## 🐛 Debugging

For debugging failed tests:

1. Check server logs in `/tmp/mcp-server.log`
2. Verify database connectivity
3. Check fixture loading in test output
4. Use `-v` flag for verbose test output

## 📈 Next Steps

This is **Phase 3** of the BDD remediation plan:

- ✅ Phase 1: Eliminate code duplication
- ✅ Phase 2: Extract shared MCP protocol
- ✅ Phase 3: Enhanced test infrastructure (this)
- 🔄 Phase 4: Advanced testing capabilities
- 🔄 Phase 5: Migration coordination