# BDD Testing Guide for Movies MCP Server

![BDD Tests](https://github.com/francknouama/movies-mcp-server/workflows/CI%20Pipeline/badge.svg)
![BDD Smoke Tests](https://github.com/francknouama/movies-mcp-server/workflows/BDD%20Smoke%20Tests/badge.svg)
![Advanced Tests](https://github.com/francknouama/movies-mcp-server/workflows/Advanced%20Testing%20Pipeline/badge.svg)

This guide provides comprehensive documentation for Behavior-Driven Development (BDD) testing using Godog (Cucumber for Go) in the Movies MCP Server project.

## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Running Tests](#running-tests)
4. [Writing New Scenarios](#writing-new-scenarios)
5. [Step Definitions](#step-definitions)
6. [Test Data and Fixtures](#test-data-and-fixtures)
7. [Architecture](#architecture)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [CI/CD Integration](#cicd-integration)

## Overview

### What is BDD?

Behavior-Driven Development (BDD) is a software development approach that emphasizes collaboration between developers, QA, and non-technical stakeholders. Tests are written in natural language (Gherkin) that describes the behavior of the system.

### Why BDD for MCP Server?

- **Readable**: Tests are written in plain English using Given-When-Then syntax
- **Collaborative**: Business stakeholders can understand and contribute to tests
- **Comprehensive**: Covers end-to-end scenarios including MCP protocol, tools, and resources
- **Maintainable**: Separates test scenarios from implementation details

### Test Coverage

Current BDD test suite includes **88 scenarios** across **8 feature files**:

- **Actor Operations** (13 scenarios): CRUD operations, search, relationships
- **Movie Operations** (9 scenarios): CRUD operations, genre management
- **MCP Protocol** (5 scenarios): Initialization, capabilities, protocol compliance
- **Error Handling** (14 scenarios): Validation, fault injection, resilience
- **Advanced Search** (9 scenarios): Complex queries, pagination, concurrency
- **Advanced Resources** (17 scenarios): Resource discovery, templates, dynamic resources
- **Performance** (10 scenarios): Load testing, memory monitoring, throughput
- **Contract Testing** (11 scenarios): Schema validation, version compatibility

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite support (modernc.org/sqlite)
- Godog testing framework

### Installation

```bash
# Install Godog
go install github.com/cucumber/godog/cmd/godog@latest

# Verify installation
godog --version
```

### Project Structure

```
tests/bdd/
├── features/              # Feature files (Gherkin scenarios)
│   ├── actor_operations.feature
│   ├── movie_operations.feature
│   ├── mcp_protocol.feature
│   ├── error_handling.feature
│   ├── advanced_search.feature
│   ├── advanced_resources.feature
│   ├── performance.feature
│   └── contract_testing.feature
├── steps/                 # Step definitions (Go code)
│   ├── actor_steps.go
│   ├── movie_steps.go
│   ├── mcp_protocol_steps.go
│   ├── error_handling_steps.go
│   ├── advanced_search_steps.go
│   ├── performance_steps.go
│   ├── contract_testing_steps.go
│   └── common_steps.go
├── context/               # Test context and state management
│   └── bdd_context.go
├── support/               # Test utilities and helpers
│   ├── sqlite_database.go
│   ├── test_utilities.go
│   ├── test_data_manager.go
│   └── fault_injection_stub.go
├── fixtures/              # Test data fixtures
│   ├── actors/
│   ├── movies/
│   └── scenarios/
├── types/                 # Public response types
│   └── responses.go
└── bdd_test.go           # Test runner
```

## Running Tests

### Run All Tests

```bash
# Run all BDD tests
go test -v ./tests/bdd/...

# Or using godog directly
cd tests/bdd
godog run features/
```

### Run Specific Feature

```bash
# Run only actor operations
godog run features/actor_operations.feature

# Run only performance tests
godog run features/performance.feature
```

### Run Tests by Tag

```bash
# Run tests tagged with @smoke
godog run --tags @smoke

# Run tests excluding @slow tag
godog run --tags "~@slow"

# Run tests with multiple tags
godog run --tags "@smoke && @actors"
```

### Run Specific Scenario

```bash
# Run scenario by name
godog run --name "Add a new actor with complete information"

# Run scenario by line number
godog run features/actor_operations.feature:10
```

### Test Output Options

```bash
# Pretty format (default)
godog run --format pretty

# Progress dots
godog run --format progress

# JSON output for CI/CD
godog run --format json > test-results.json

# JUnit XML for CI/CD
godog run --format junit > test-results.xml
```

### Parallel Execution

```bash
# Run tests in parallel (4 workers)
godog run --concurrency 4

# Note: Ensure tests are isolated for parallel execution
```

### With Timeout

```bash
# Set timeout for long-running tests
timeout 120s go test -v -timeout=2m ./tests/bdd/...
```

### Testing SDK vs Legacy Server

```bash
# Test with SDK server (default)
go test ./tests/bdd/

# Test with legacy server
TEST_MCP_SERVER=legacy go test ./tests/bdd/

# Test both implementations
TEST_MCP_SERVER=sdk go test ./tests/bdd/ && \
TEST_MCP_SERVER=legacy go test ./tests/bdd/
```

## Writing New Scenarios

### Feature File Template

Create a new feature file in `tests/bdd/features/`:

```gherkin
Feature: [Feature Name]
  As a [user role]
  I want to [action]
  So that [business value]

  Background:
    Given [common setup for all scenarios]

  @tag1 @tag2
  Scenario: [Scenario title]
    Given [initial context]
    When [action/event]
    Then [expected outcome]
    And [additional expectations]

  Scenario Outline: [Scenario with examples]
    Given I have a <parameter>
    When I perform <action>
    Then the result should be <expected>

    Examples:
      | parameter | action | expected |
      | value1    | act1   | result1  |
      | value2    | act2   | result2  |
```

### Gherkin Syntax Guidelines

#### Given Steps
Define the initial state or context:

```gherkin
Given I have a valid MCP client connection
Given I have 10 movies in the database
Given the following actors exist:
  | name       | birth_year |
  | Tom Hanks  | 1956       |
```

#### When Steps
Describe the action or event:

```gherkin
When I send a tools/list request
When I search for actors by name "Tom Hanks"
When I add a new movie with the following details:
  """json
  {
    "title": "Inception",
    "director": "Christopher Nolan",
    "year": 2010
  }
  """
```

#### Then Steps
Verify the expected outcome:

```gherkin
Then the response should be successful
Then the response should contain 5 movies
Then the actor should have ID 1 and name "Tom Hanks"
Then the error message should contain "not found"
```

### Using Tables

For structured data:

```gherkin
Scenario: Search multiple actors
  Given the following actors exist:
    | name              | birth_year |
    | Tom Hanks         | 1956       |
    | Meryl Streep      | 1949       |
  When I list all actors
  Then the response should contain these actors:
    | name              | birth_year |
    | Tom Hanks         | 1956       |
    | Meryl Streep      | 1949       |
```

### Using Doc Strings

For JSON or large text blocks:

```gherkin
Scenario: Create movie with JSON
  When I send the following movie data:
    """json
    {
      "title": "The Matrix",
      "director": "Wachowski Sisters",
      "year": 1999,
      "rating": 8.7,
      "genres": ["Action", "Sci-Fi"]
    }
    """
  Then the movie should be created successfully
```

### Tags

Organize and filter tests:

```gherkin
@smoke @actors @crud
Scenario: Basic actor CRUD operations
  # Scenario steps...

@slow @performance
Scenario: Load test with 1000 movies
  # Scenario steps...

@skip @wip
Scenario: Feature under development
  # Scenario steps...
```

Common tags:
- `@smoke`: Critical tests for smoke testing
- `@slow`: Long-running tests
- `@performance`: Performance tests
- `@integration`: Integration tests
- `@skip` or `@wip`: Tests to skip (work in progress)

## Step Definitions

### Creating New Step Definitions

Step definitions map Gherkin steps to Go functions.

#### Basic Structure

Create or update a file in `tests/bdd/steps/`:

```go
package steps

import (
    "fmt"
    "github.com/cucumber/godog"
    "github.com/francknouama/movies-mcp-server/tests/bdd/context"
)

// Initialize registers step definitions
func InitializeMySteps(ctx *godog.ScenarioContext) {
    stepContext := NewCommonStepContext()

    // Register steps
    ctx.Step(`^I perform some action$`, stepContext.iPerformSomeAction)
    ctx.Step(`^I perform action with parameter "([^"]*)"$`,
        stepContext.iPerformActionWithParameter)
}

// Step implementation
func (c *CommonStepContext) iPerformSomeAction() error {
    // Perform the action
    response, err := c.bddContext.CallTool("tool_name", map[string]interface{}{
        "param": "value",
    })

    if err != nil {
        return fmt.Errorf("action failed: %w", err)
    }

    if response.IsError {
        return fmt.Errorf("MCP error: %v", response.Content)
    }

    // Store result for later assertions
    c.bddContext.SetTestData("action_result", response.Content)

    return nil
}

func (c *CommonStepContext) iPerformActionWithParameter(param string) error {
    // Use the captured parameter
    fmt.Printf("Parameter: %s\n", param)
    return nil
}
```

#### Regex Patterns

Common regex patterns for step matching:

```go
// Match quoted strings
`^I search for "([^"]*)"$`  // Matches: I search for "Tom Hanks"

// Match numbers
`^I have (\d+) movies$`      // Matches: I have 10 movies
`^the rating should be ([\d.]+)$`  // Matches: the rating should be 8.5

// Match optional words
`^the operation should complete within (\d+) seconds?$`

// Match tables
func (c *CommonStepContext) iHaveActors(table *godog.Table) error {
    for i, row := range table.Rows {
        if i == 0 {
            continue // Skip header
        }
        name := row.Cells[0].Value
        birthYear := row.Cells[1].Value
        // Process row...
    }
    return nil
}

// Match doc strings
func (c *CommonStepContext) iSendJSON(docString *godog.DocString) error {
    jsonData := docString.Content
    // Parse and process JSON...
    return nil
}
```

### Step Context and State

The `CommonStepContext` provides shared state:

```go
type CommonStepContext struct {
    bddContext  *context.BDDContext     // Test context
    testDB      DatabaseInterface        // Database access
    dataManager *support.TestDataManager // Test data tracking
}

// Store data for later steps
c.bddContext.SetTestData("key", value)

// Retrieve data
value, exists := c.bddContext.GetTestData("key")

// Check for errors
if c.bddContext.HasError() {
    message := c.bddContext.GetErrorMessage()
}
```

### Calling MCP Tools

```go
// Call a tool
response, err := c.bddContext.CallTool("add_movie", map[string]interface{}{
    "title":    "Inception",
    "director": "Christopher Nolan",
    "year":     2010,
    "rating":   8.8,
})

if err != nil {
    return fmt.Errorf("tool call failed: %w", err)
}

if response.IsError {
    return fmt.Errorf("MCP error: %v", response.Content)
}

// Parse response
var movie types.MovieResponse
if err := response.ParseContent(&movie); err != nil {
    return fmt.Errorf("failed to parse response: %w", err)
}
```

### Registering Steps

Add your initialization function to `bdd_test.go`:

```go
func InitializeScenario(ctx *godog.ScenarioContext) {
    steps.InitializeCommonSteps(ctx)
    steps.InitializeActorSteps(ctx)
    steps.InitializeMovieSteps(ctx)
    steps.InitializeMCPProtocolSteps(ctx)
    steps.InitializeMyNewSteps(ctx)  // Add your new steps
}
```

## Test Data and Fixtures

### Using Fixtures

Fixtures provide reusable test data. See [fixtures/README.md](fixtures/README.md) for details.

#### Load Fixtures in Scenarios

```gherkin
Scenario: Test with fixture data
  Given I have loaded the "basic_actors" fixture
  When I search for actors
  Then the response should contain 5 actors
```

#### Load Fixtures in Steps

```go
// Load a fixture programmatically
err := c.sqliteDB.LoadFixtures("basic_movies")
if err != nil {
    return fmt.Errorf("failed to load fixtures: %w", err)
}
```

### Test Data Management

The `TestDataManager` tracks created entities:

```go
// Track created entity
c.dataManager.AddCreatedID("movies", movieID)

// Get all created IDs
movieIDs := c.dataManager.GetCreatedIDs("movies")

// Clear after scenario
c.dataManager.Clear()
```

### Database Operations

Direct database access when needed:

```go
// Verify entity exists
exists, err := c.testDB.VerifyMovieExists("Inception", "Christopher Nolan", 2010)

// Count rows
count, err := c.testDB.CountRows("actors", "birth_year > ?", 1980)

// Execute custom query
result, err := c.testDB.ExecuteQuery(
    "SELECT * FROM movies WHERE rating > ?", 8.0)
```

### Cleanup

Tests automatically clean up between scenarios:

```go
// Automatic cleanup in common_steps.go
func (c *CommonStepContext) teardownScenario() error {
    // Stop MCP server
    c.bddContext.StopMCPServer()

    // Cleanup database
    if c.sqliteDB != nil {
        c.sqliteDB.Cleanup()
    }

    // Clear test data
    c.dataManager.Clear()

    return nil
}
```

## Architecture

### Test Flow

```
1. Feature File (Gherkin)
   ↓
2. Godog parses scenarios
   ↓
3. Matches steps to definitions
   ↓
4. Step functions execute
   ↓
5. BDDContext manages state
   ↓
6. MCP Client calls tools
   ↓
7. MCP Server processes requests
   ↓
8. SQLite Database stores data
   ↓
9. Response validated
   ↓
10. Cleanup and teardown
```

### Component Responsibilities

#### BDDContext (`context/bdd_context.go`)
- Manages test execution context
- Starts/stops MCP server
- Stores test data and state
- Handles MCP client communication

#### Steps (`steps/*.go`)
- Implement Gherkin step definitions
- Parse step parameters
- Call MCP tools
- Validate responses
- Assert expectations

#### Support (`support/*.go`)
- **SQLiteTestDatabase**: Database operations and migrations
- **TestUtilities**: Helper functions for test data
- **TestDataManager**: Track created entities
- **FaultInjector**: Chaos engineering (stubbed)

#### Types (`types/responses.go`)
- Public response types
- Avoid internal package imports
- Type-safe response parsing

## Best Practices

### Writing Scenarios

1. **Keep scenarios focused**: One scenario tests one behavior
2. **Use descriptive names**: Clearly state what is being tested
3. **Follow Given-When-Then**: Maintain clear structure
4. **Avoid technical details**: Focus on business behavior
5. **Use Background wisely**: Extract common setup

**Good Example:**
```gherkin
Scenario: Search for actors by birth decade
  Given I have 10 actors born in different decades
  When I search for actors born in the 1970s
  Then the response should contain only actors born between 1970 and 1979
```

**Bad Example:**
```gherkin
Scenario: Database query test
  Given I insert 10 rows into actors table
  When I execute SELECT * FROM actors WHERE birth_year BETWEEN 1970 AND 1979
  Then I should get some results
```

### Step Definitions

1. **Reuse existing steps**: Don't duplicate step definitions
2. **Keep steps simple**: Complex logic belongs in helper functions
3. **Use clear error messages**: Help debugging failed tests
4. **Avoid sleeps**: Use proper synchronization
5. **Clean up resources**: Ensure teardown is complete

### Test Data

1. **Use fixtures**: For consistent, reusable data
2. **Isolate tests**: Each scenario should be independent
3. **Clean state**: Reset database between scenarios
4. **Realistic data**: Use real movie/actor names
5. **Cover edge cases**: Include boundary values

### Tags and Organization

1. **Tag consistently**: Use standard tags across features
2. **Group related tests**: Keep features cohesive
3. **Separate slow tests**: Mark performance tests as `@slow`
4. **Enable selective runs**: Use tags for CI/CD pipelines

## Troubleshooting

### Common Issues

#### 1. Step Definition Not Found

**Error**: `Step is undefined`

**Solution**:
- Ensure step is registered in `InitializeScenario`
- Check regex pattern matches exactly
- Verify step function signature

#### 2. Database Connection Failed

**Error**: `failed to initialize SQLite test database`

**Solution**:
- Check SQLite is installed and accessible
- Verify database path permissions
- Ensure migrations are valid SQL

#### 3. MCP Server Won't Start

**Error**: `failed to start MCP server`

**Solution**:
- Check server binary exists and is executable
- Verify environment variables are set correctly
- Check for port conflicts
- Review server logs for errors

#### 4. Test Hangs/Timeout

**Error**: `test timeout exceeded`

**Solution**:
- Check for missing assertions that cause blocking
- Verify MCP server is responding
- Increase timeout: `go test -timeout=5m`
- Look for goroutine leaks

#### 5. Fixture Load Failed

**Error**: `failed to load fixtures`

**Solution**:
- Verify fixture file path is correct
- Check JSON syntax is valid
- Ensure required fields are present
- Verify database schema matches fixture

#### 6. Response Parse Error

**Error**: `failed to parse response`

**Solution**:
- Check response type matches expected type
- Verify JSON structure is correct
- Ensure all fields are exported (capitalized)
- Add logging to inspect actual response

### Debugging Tips

#### Enable Verbose Output

```bash
go test -v ./tests/bdd/...
```

#### Add Debug Logging

```go
import "log"

func (c *CommonStepContext) iSearchMovies(genre string) error {
    log.Printf("Searching for movies with genre: %s", genre)

    response, err := c.bddContext.CallTool("search_movies", map[string]interface{}{
        "genre": genre,
    })

    log.Printf("Response: %+v", response)
    log.Printf("Error: %v", err)

    return err
}
```

#### Inspect Database State

```go
// In step definition
count, _ := c.testDB.CountRows("movies", "1=1")
log.Printf("Total movies in database: %d", count)
```

#### Check MCP Server Logs

```bash
# Set environment variable for server logs
export MCP_LOG_LEVEL=debug
go test -v ./tests/bdd/...
```

### Performance Issues

If tests are slow:

1. **Use in-memory SQLite**: Already default with temporary databases
2. **Reduce data volume**: Use smaller fixtures for non-performance tests
3. **Skip slow tests locally**: `godog run --tags "~@slow"`
4. **Run in parallel**: `godog run --concurrency 4`
5. **Profile tests**: `go test -cpuprofile=cpu.prof`

## CI/CD Integration

### GitHub Actions

Example workflow file (`.github/workflows/bdd-tests.yml`):

```yaml
name: BDD Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  bdd-tests:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: |
        go mod download
        go install github.com/cucumber/godog/cmd/godog@latest

    - name: Run BDD tests
      run: |
        cd tests/bdd
        godog run --format=pretty --format=json:../../test-results.json

    - name: Upload test results
      if: always()
      uses: actions/upload-artifact@v3
      with:
        name: test-results
        path: test-results.json
```

### Running Specific Suites

```yaml
- name: Run smoke tests
  run: godog run --tags @smoke

- name: Run full test suite
  run: godog run

- name: Run performance tests (nightly)
  if: github.event.schedule
  run: godog run --tags @performance
```

### Test Reporting

Generate and publish reports:

```bash
# HTML report
godog run --format=pretty --format=html:report.html

# JUnit XML for CI integration
godog run --format=junit:report.xml

# Cucumber JSON for advanced reporting
godog run --format=cucumber:report.json
```

## Additional Resources

- [Godog Documentation](https://github.com/cucumber/godog)
- [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)
- [BDD Best Practices](https://cucumber.io/docs/bdd/)
- [MCP Protocol Specification](https://github.com/modelcontextprotocol/specification)

## Contributing

When adding new BDD tests:

1. Follow the existing patterns in feature files
2. Reuse step definitions when possible
3. Add new fixtures for complex scenarios
4. Document any new test utilities
5. Update this README if adding new concepts
6. Ensure all tests pass before committing

---

**Last Updated**: 2025-11-01
**Maintainer**: Movies MCP Server Team