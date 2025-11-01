# BDD Feature File Templates

This directory contains templates for creating new BDD test scenarios. Use these templates as starting points for writing your own feature files.

## Available Templates

### 1. Basic Feature Template
**File**: `basic_feature_template.feature`

Use this template for simple scenarios with straightforward Given-When-Then steps.

**When to use**:
- Simple CRUD operations
- Basic functionality tests
- Scenarios with clear linear flow

**Example**:
```gherkin
Feature: Actor Management
  As a movie database administrator
  I want to manage actor information
  So that I can maintain accurate actor records

  @smoke @actors
  Scenario: Add a new actor
    Given I have a valid MCP client connection
    When I add an actor with name "Tom Hanks" and birth year 1956
    Then the actor should be created successfully
```

### 2. Table Data Template
**File**: `table_data_template.feature`

Use this template when you need to work with multiple records or structured data.

**When to use**:
- Multiple entities need to be created
- Validation requires checking multiple records
- Structured input or output data

**Example**:
```gherkin
Scenario: Bulk actor creation
  Given the following actors exist:
    | name       | birth_year |
    | Tom Hanks  | 1956       |
    | Meryl Streep | 1949     |
```

### 3. Scenario Outline Template
**File**: `scenario_outline_template.feature`

Use this template to test the same scenario with different parameters.

**When to use**:
- Testing multiple variations of the same behavior
- Validation tests with different inputs
- Boundary value testing

**Example**:
```gherkin
Scenario Outline: Rating validation
  When I create a movie with rating <rating>
  Then the operation should <result>

  Examples:
    | rating | result  |
    | 9.5    | succeed |
    | 11.0   | fail    |
```

### 4. JSON Data Template
**File**: `json_data_template.feature`

Use this template for scenarios involving complex JSON structures.

**When to use**:
- Creating entities with nested data
- Complex request payloads
- API-style testing

**Example**:
```gherkin
Scenario: Create movie with complex data
  When I send the following movie data:
    """json
    {
      "title": "Inception",
      "genres": ["Action", "Sci-Fi"]
    }
    """
  Then the movie should be created
```

### 5. Error Handling Template
**File**: `error_handling_template.feature`

Use this template for testing error conditions and validation rules.

**When to use**:
- Validation testing
- Error code verification
- Negative testing scenarios

**Example**:
```gherkin
Scenario: Invalid birth year
  When I try to create an actor with birth year 2050
  Then the operation should fail
  And the error message should contain "invalid birth year"
```

### 6. Performance Template
**File**: `performance_template.feature`

Use this template for load, stress, and performance testing.

**When to use**:
- Load testing
- Concurrency testing
- Memory and resource usage testing
- Throughput measurement

**Example**:
```gherkin
@slow @performance
Scenario: Concurrent searches
  When I perform 20 concurrent searches
  Then all searches should complete within 5 seconds
```

### 7. MCP Protocol Template
**File**: `mcp_protocol_template.feature`

Use this template for testing MCP protocol compliance.

**When to use**:
- Protocol compliance testing
- Tool and resource discovery
- Protocol version compatibility
- Error code verification

**Example**:
```gherkin
Scenario: List MCP tools
  When I send a tools/list request
  Then the response should contain the tools
```

## Using Templates

### 1. Copy Template

```bash
# Copy a template to create a new feature file
cp tests/bdd/templates/basic_feature_template.feature \
   tests/bdd/features/my_new_feature.feature
```

### 2. Customize

Edit the new feature file:
- Replace `[placeholders]` with actual values
- Update tags to match your scenario
- Add or remove steps as needed
- Ensure step definitions exist or create new ones

### 3. Validate

Check your feature file:

```bash
# Run godog to check for undefined steps
cd tests/bdd
godog run features/my_new_feature.feature
```

### 4. Implement Missing Steps

If you see "undefined" step errors, implement them in the appropriate `steps/*.go` file:

```go
func (c *CommonStepContext) iAddAnActorWithNameAndBirthYear(name string, birthYear int) error {
    // Implementation
}
```

## Template Selection Guide

| Use Case | Recommended Template |
|----------|---------------------|
| Basic CRUD | Basic Feature Template |
| Multiple records | Table Data Template |
| Variations of same test | Scenario Outline Template |
| Complex JSON payloads | JSON Data Template |
| Validation/Error cases | Error Handling Template |
| Performance testing | Performance Template |
| Protocol compliance | MCP Protocol Template |

## Best Practices

### 1. Naming Conventions

```gherkin
# Good: Descriptive and specific
Feature: Movie Search Functionality
Scenario: Search movies by genre and year range

# Bad: Vague
Feature: Tests
Scenario: Test 1
```

### 2. Tags

Use consistent tags:
- `@smoke`: Critical functionality
- `@slow`: Long-running tests
- `@wip`: Work in progress
- `@domain`: e.g., @actors, @movies
- `@type`: e.g., @crud, @search, @validation

### 3. Background

Use Background for common setup:

```gherkin
Background:
  Given I have a valid MCP client connection
  And I have a clean database
```

### 4. Step Reusability

Write steps that can be reused:

```gherkin
# Good: Reusable
When I search for actors by name "Tom Hanks"

# Bad: Too specific
When I search for Tom Hanks who was born in 1956
```

### 5. Assertions

Be explicit about expected outcomes:

```gherkin
# Good: Clear expectations
Then the response should contain 5 movies
And each movie should have a valid rating

# Bad: Vague
Then it should work
```

## Adding New Templates

If you create a useful template pattern:

1. Create the template file in this directory
2. Document it in this README
3. Add examples showing when to use it
4. Include a selection guide entry

## Examples

See the existing feature files in `tests/bdd/features/` for real-world examples of these templates in use:

- `actor_operations.feature` - Uses Basic and Table templates
- `movie_operations.feature` - Uses Basic and JSON templates
- `error_handling.feature` - Uses Error Handling template
- `performance.feature` - Uses Performance template
- `mcp_protocol.feature` - Uses MCP Protocol template

## Questions?

See the main [BDD Testing Guide](../README.md) for comprehensive documentation on:
- Writing scenarios
- Step definitions
- Best practices
- Troubleshooting

---

**Last Updated**: 2025-11-01
