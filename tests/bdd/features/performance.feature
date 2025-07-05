Feature: Performance Requirements
  As an MCP client
  I want the system to handle operations efficiently
  So that I can work with large datasets without performance degradation

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And the database is clean

  @performance @concurrency
  Scenario: Handle Concurrent Movie Searches
    Given I have 1000 movies in the database
    When I perform 50 concurrent searches for "action"
    Then all searches should complete within 2 seconds
    And no search should fail due to resource contention
    And each search result should be valid

  @performance @large-dataset
  Scenario: Large Database Operations
    Given I have 10000 movies in the database
    When I search for movies by genre "drama"
    Then the response should be returned within 500ms
    And the response should contain all matching movies
    And the memory usage should not exceed baseline by more than 50MB

  @performance @batch-operations
  Scenario: Batch Movie Creation Performance
    When I create 100 movies in batch
    Then the operation should complete within 1 second
    And all movies should be successfully created
    And the database should contain exactly 100 movies

  @performance @concurrent-writes
  Scenario: Concurrent Movie Creation
    When I create 20 movies concurrently
    Then all operations should complete within 3 seconds
    And no operations should fail due to conflicts
    And the database should contain exactly 20 movies
    And all movie IDs should be unique

  @performance @pagination
  Scenario: Paginated Results Performance
    Given I have 5000 movies in the database
    When I request page 50 with 100 movies per page
    Then the response should be returned within 200ms
    And the response should contain exactly 100 movies
    And the pagination metadata should be correct

  @performance @complex-search
  Scenario: Complex Search Performance
    Given I have 2000 movies with various attributes
    When I perform a complex search with multiple filters:
      | field    | value      |
      | genre    | Action     |
      | year_min | 2010       |
      | year_max | 2020       |
      | rating   | >8.0       |
    Then the response should be returned within 300ms
    And all returned movies should match the search criteria

  @performance @memory
  Scenario: Memory Usage During Large Operations
    Given I measure the baseline memory usage
    When I load 1000 movies with full details
    Then the memory increase should not exceed 100MB
    And the memory should be released after the operation

  @performance @concurrent-reads
  Scenario: High Read Concurrency
    Given I have 500 movies in the database
    When I perform 100 concurrent read operations
    Then all operations should complete within 1 second
    And all responses should be successful
    And no data corruption should occur

  @performance @stress-test
  Scenario: System Stress Test
    Given I have 1000 movies and 500 actors in the database
    When I perform mixed operations for 30 seconds:
      | operation    | percentage |
      | search       | 60%        |
      | create       | 20%        |
      | update       | 15%        |
      | delete       | 5%         |
    Then the system should maintain response times under 500ms
    And no operations should fail
    And the database should remain consistent

  @performance @resource-cleanup
  Scenario: Resource Cleanup Performance
    Given I have created 1000 temporary test records
    When I delete all test records
    Then the cleanup should complete within 2 seconds
    And all records should be properly removed
    And database space should be reclaimed