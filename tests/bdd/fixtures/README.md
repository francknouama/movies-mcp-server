# BDD Test Fixtures

This directory contains test data fixtures for BDD scenarios in JSON format. Fixtures provide consistent, reusable test data for different testing scenarios.

## Directory Structure

```
fixtures/
├── actors/              # Actor-related test data
│   ├── basic_actors.json
│   ├── search_test_actors.json
│   └── relationship_actors.json
├── movies/              # Movie-related test data
│   ├── basic_movies.json
│   ├── decade_movies.json
│   ├── rating_test_movies.json
│   └── genre_test_movies.json
├── scenarios/           # Complex multi-entity scenarios
│   ├── movie_cast_scenario.json
│   ├── director_analysis_scenario.json
│   └── search_performance_scenario.json
└── README.md           # This file
```

## Fixture Files

### Actors

- **basic_actors.json**: 5 well-known actors for basic CRUD operations
- **search_test_actors.json**: 8 actors for testing search and filtering
- **relationship_actors.json**: 5 actors with known movie relationships

### Movies

- **basic_movies.json**: 5 classic movies for basic CRUD operations
- **decade_movies.json**: Movies organized by decade (1970s-2010s)
- **rating_test_movies.json**: Movies grouped by rating ranges
- **genre_test_movies.json**: Movies organized by genre

### Scenarios

- **movie_cast_scenario.json**: Complete Interstellar cast and crew
- **director_analysis_scenario.json**: Christopher Nolan filmography with recurring actors
- **search_performance_scenario.json**: Large dataset for performance testing

## Using Fixtures in Tests

### Loading Fixtures in Go Tests

```go
import (
    "encoding/json"
    "os"
)

// Load a fixture file
func loadFixture(filename string, target interface{}) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, target)
}

// Example usage
type ActorsFixture struct {
    Actors []Actor `json:"actors"`
}

func setupActorData() error {
    var fixture ActorsFixture
    err := loadFixture("tests/bdd/fixtures/actors/basic_actors.json", &fixture)
    if err != nil {
        return err
    }

    // Use fixture.Actors in your test
    for _, actor := range fixture.Actors {
        // Add to database
    }
    return nil
}
```

### Using with SQLiteTestDatabase

The `SQLiteTestDatabase` helper includes a `LoadFixtures` method:

```go
// Load fixtures by name
err := testDB.LoadFixtures("basic_actors")
if err != nil {
    return fmt.Errorf("failed to load fixtures: %w", err)
}
```

### Using in Feature Files

Reference fixtures in your Gherkin scenarios:

```gherkin
Scenario: Search actors from basic fixture
  Given I have loaded the "basic_actors" fixture
  When I search for actors by name "Tom Hanks"
  Then the response should contain 1 actor
```

## Creating New Fixtures

### Guidelines

1. **Use realistic data**: Movie titles, actor names should be real or realistic
2. **Include all required fields**: Ensure all mandatory fields are present
3. **Organize logically**: Group related data together
4. **Document purpose**: Add a "description" field explaining the fixture's purpose
5. **Keep it focused**: Each fixture should test specific functionality

### Example Template

```json
{
  "description": "Brief description of what this fixture tests",
  "entities": [
    {
      "field1": "value1",
      "field2": "value2"
    }
  ]
}
```

### Validation

Before committing new fixtures:

1. Validate JSON syntax
2. Ensure all required fields are present
3. Test loading the fixture in a scenario
4. Document the fixture's purpose

## Fixture Data Ranges

### Actors
- Birth years: 1933-1984
- Total count: 18 unique actors across all fixtures

### Movies
- Release years: 1957-2020
- Ratings: 6.8-9.3
- Genres: 15+ unique genres
- Total count: 30+ unique movies across all fixtures

## Best Practices

1. **Isolation**: Each test scenario should work with its own fixture or clean state
2. **Consistency**: Use the same fixture for similar tests
3. **Completeness**: Include edge cases (min/max values, special characters)
4. **Maintenance**: Update fixtures when schema changes
5. **Documentation**: Document any special considerations or relationships

## Troubleshooting

### Fixture Won't Load

- Check JSON syntax with a validator
- Verify file path is correct
- Ensure all required fields are present
- Check file permissions

### Data Conflicts

- Ensure unique identifiers (titles, names)
- Clean database between scenarios
- Use `CleanupAfterScenario()` in teardown

### Performance Issues

- Limit fixture size for unit tests
- Use performance fixtures only for load testing
- Consider lazy loading for large datasets
