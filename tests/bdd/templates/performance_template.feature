Feature: [Performance Testing]
  Performance and load testing scenarios

  Background:
    Given I have a valid MCP client connection

  @slow @performance
  Scenario: [Concurrent operations]
    Given I have 100 movies in the database
    When I perform 20 concurrent searches for "Action"
    Then all searches should complete within 5 seconds
    And no search should fail due to resource contention
    And each search result should be valid

  @slow @performance
  Scenario: [Batch operations]
    When I create 100 movies in batch
    Then the operation should complete within 10 seconds
    And all movies should be successfully created

  @performance @memory
  Scenario: [Memory usage]
    Given I measure the baseline memory usage
    When I load 1000 movies with full details
    Then the operation should complete within 5 seconds
    And the memory increase should not exceed 50MB

  @performance @throughput
  Scenario: [System throughput]
    Given I have 500 actors and 1000 movies in the database
    When I perform mixed operations for 60 seconds:
      | operation      | percentage |
      | search         | 50%        |
      | create         | 30%        |
      | read           | 15%        |
      | update         | 5%         |
    Then the average response time should be under 100ms
    And the throughput should exceed 100 operations per second
    And the error rate should be under 1%
