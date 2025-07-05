Feature: Advanced Resource Testing
  As an MCP client
  I want comprehensive resource endpoint coverage
  So that all MCP resources work correctly and efficiently

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And I have test data in the database

  @resources @stats
  Scenario: Database Statistics Resource
    Given I have 5 movies and 3 actors in the database
    When I request the "movies://database/stats" resource
    Then the response should be successful
    And the response should contain movie count 5
    And the response should contain actor count 3
    And the response should include storage usage information
    And the response should have last_updated timestamp
    And the statistics should be accurate

  @resources @collection
  Scenario: All Movies Resource
    Given I have movies with various attributes in the database
    When I request the "movies://database/all" resource
    Then the response should be successful
    And the response should contain all movies in the database
    And each movie should have complete information
    And the response should include metadata with total count
    And the metadata should have generated timestamp

  @resources @actors
  Scenario: All Actors Resource
    Given I have actors in the database
    When I request the "movies://actors/all" resource
    Then the response should be successful
    And the response should contain all actors in the database
    And each actor should have complete information
    And the response should include metadata

  @resources @recent
  Scenario: Recent Movies Resource
    Given I have movies created in the last 7 days
    And I have movies created more than 7 days ago
    When I request the "movies://search/recent" resource
    Then the response should be successful
    And the response should contain only recent movies
    And the search criteria should be included in the response
    And the cutoff date should be correct

  @resources @recent @parameterized
  Scenario: Recent Movies with Custom Parameters
    Given I have movies created in the last 30 days
    When I request the "movies://search/recent" resource with parameters:
      | parameter | value |
      | limit     | 5     |
      | days      | 14    |
    Then the response should be successful
    And the response should contain at most 5 movies
    And all movies should be from the last 14 days
    And the search criteria should reflect the parameters

  @resources @posters
  Scenario: Movie Posters Collection Resource
    Given I have movies with poster images
    When I request the "movies://posters/collection" resource
    Then the response should be successful
    And the response should contain poster data for each movie
    And each poster should be properly encoded in base64
    And the metadata should include format and size information
    And the total size should not exceed reasonable limits

  @resources @posters @formats
  Scenario: Movie Posters with Different Formats
    Given I have movies with poster images
    When I request the "movies://posters/collection" resource with format "thumbnail"
    Then the response should be successful
    And all posters should be in thumbnail format
    And the file sizes should be appropriate for thumbnails
    And the response time should be fast

  @resources @performance
  Scenario: Resource Response Time Requirements
    Given I have a populated database
    When I request each MCP resource
    Then the "movies://database/stats" resource should respond within 200ms
    And the "movies://database/all" resource should respond within 1000ms
    And the "movies://actors/all" resource should respond within 500ms
    And the "movies://search/recent" resource should respond within 300ms
    And the "movies://posters/collection" resource should respond within 2000ms

  @resources @error-handling
  Scenario: Resource Error Handling
    When I request a non-existent resource "movies://invalid/resource"
    Then I should get a resource not found error
    And the error should include the invalid resource URI
    And the server should remain stable

  @resources @invalid-uri
  Scenario: Invalid Resource URI Format
    When I request an invalid URI format "not-a-valid-uri"
    Then I should get an invalid URI format error
    And the error should explain the correct URI format
    And the server should remain stable

  @resources @parameter-validation
  Scenario: Resource Parameter Validation
    When I request the "movies://search/recent" resource with invalid parameters:
      | parameter | value | expected_error                    |
      | limit     | -1    | Limit must be positive          |
      | limit     | 1000  | Limit must not exceed 100       |
      | days      | 0     | Days must be positive           |
      | days      | 500   | Days must not exceed 365        |
    Then each request should return appropriate parameter validation errors
    And the errors should include the invalid parameter values

  @resources @caching
  Scenario: Resource Caching Behavior
    Given I have data in the database
    When I request the "movies://database/stats" resource twice
    Then the first request should populate the cache
    And the second request should be served from cache
    And the second request should be faster than the first
    And the cache should respect the 60-second TTL

  @resources @concurrent-access
  Scenario: Concurrent Resource Access
    Given I have data in the database
    When I make 10 concurrent requests to "movies://database/stats"
    Then all requests should succeed
    And all responses should be consistent
    And no race conditions should occur
    And the server should remain stable

  @resources @large-datasets
  Scenario: Large Dataset Resource Performance
    Given I have 10000 movies in the database
    When I request the "movies://database/all" resource
    Then the response should be returned within acceptable time limits
    And the memory usage should not exceed 100MB
    And the response should be complete and accurate
    And the server should remain responsive during the operation

  @resources @content-types
  Scenario: Resource Content Type Validation
    When I request each MCP resource
    Then all text resources should return JSON content
    And all binary resources should have appropriate content types
    And the content encoding should be specified correctly
    And the character encoding should be UTF-8 for text content

  @resources @security
  Scenario: Resource Security Validation
    When I request MCP resources
    Then no sensitive information should be exposed
    And access should be properly controlled
    And no unauthorized data should be accessible
    And audit trails should be maintained for resource access

  @resources @versioning
  Scenario: Resource Versioning Support
    Given the MCP server supports resource versioning
    When I request resources with version parameters
    Then the correct version of each resource should be returned
    And version negotiation should work properly
    And unsupported versions should be handled gracefully
    And backward compatibility should be maintained