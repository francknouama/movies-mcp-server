Feature: Contract Testing
  As an API consumer
  I want the MCP interface to remain stable
  So that my integrations don't break when the server is updated

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And the contract definitions are loaded

  @contract @movie-tools
  Scenario: Movie Management Tool Contracts
    When I validate the "add_movie" tool contract
    Then the tool should have required parameters: ["title", "director", "year"]
    And the tool should have optional parameters: ["genre", "rating", "description", "poster_url"]
    And the parameter constraints should be enforced:
      | parameter | constraint              |
      | rating    | float 0.0-10.0         |
      | year      | integer 1888-2030      |
      | title     | string max 255 chars   |
    And the success response should contain: ["id", "title", "director", "year", "created_at"]
    And the error codes should include: [-32602, -32603]

  @contract @actor-tools
  Scenario: Actor Management Tool Contracts
    When I validate the "add_actor" tool contract
    Then the tool should have required parameters: ["name", "birth_year"]
    And the tool should have optional parameters: ["bio", "death_year", "photo_url"]
    And the parameter constraints should be enforced:
      | parameter  | constraint              |
      | birth_year | integer 1800-2020      |
      | death_year | integer 1800-2030      |
      | name       | string max 100 chars   |
    And the success response should contain: ["id", "name", "birth_year", "created_at"]

  @contract @search-tools
  Scenario: Search Tool Contracts
    When I validate the "search_movies" tool contract
    Then the tool should have optional parameters: ["title", "director", "genre", "year_min", "year_max", "rating_min"]
    And the parameter constraints should be enforced:
      | parameter  | constraint              |
      | year_min   | integer 1888-2030      |
      | year_max   | integer 1888-2030      |
      | rating_min | float 0.0-10.0         |
    And the success response should contain an array of movies
    And each movie should have: ["id", "title", "director", "year", "rating"]

  @contract @resources
  Scenario: MCP Resource Contracts
    When I validate the MCP resources
    Then the "movies://database/stats" resource should be available
    And the "movies://database/all" resource should be available
    And the stats resource should return:
      | field       | type    |
      | movie_count | integer |
      | actor_count | integer |
      | total_size  | integer |
    And the all resource should return an array of movies

  @contract @error-responses
  Scenario: Error Response Contracts
    When I test error response contracts
    Then all errors should follow JSON-RPC 2.0 format
    And error responses should contain: ["jsonrpc", "error", "id"]
    And error objects should contain: ["code", "message"]
    And error codes should be consistent:
      | code  | meaning              |
      | -32700| Parse error          |
      | -32600| Invalid Request      |
      | -32601| Method not found     |
      | -32602| Invalid params       |
      | -32603| Internal error       |

  @contract @backward-compatibility
  Scenario: Backward Compatibility Validation
    Given I have a baseline contract from version 1.0
    When I compare with the current contract
    Then no required parameters should be removed
    And no response fields should be removed
    And parameter constraints should not be more restrictive
    And error codes should remain consistent

  @contract @tool-schema
  Scenario: Tool Schema Contract Validation
    When I request the tools list from the server
    Then each tool should have a valid JSON schema
    And the schema should include: ["type", "properties", "required"]
    And all required fields should be marked as required
    And all constraints should be properly defined in the schema

  @contract @regression
  Scenario: Contract Regression Detection
    Given I have contracts from the previous version
    When I run contract regression tests
    Then no breaking changes should be detected
    And any new features should be additive only
    And deprecated features should be properly marked
    And migration guides should be provided for any changes

  @contract @versioning
  Scenario: API Versioning Contract
    When I check the MCP protocol version
    Then the server should declare version "2024-11-05"
    And the protocol should remain compatible
    And version negotiation should work correctly
    And unsupported versions should be rejected gracefully

  @contract @data-types
  Scenario: Data Type Contract Validation
    When I validate data type contracts
    Then all dates should be in ISO 8601 format
    And all IDs should be positive integers
    And all ratings should be floats between 0.0 and 10.0
    And all years should be integers between 1888 and 2030
    And all strings should have defined maximum lengths

  @contract @performance-contracts
  Scenario: Performance Contract Validation
    When I validate performance contracts
    Then simple operations should complete within 100ms
    And search operations should complete within 500ms
    And batch operations should complete within 2 seconds
    And the server should handle 100 concurrent requests
    And memory usage should not exceed defined limits