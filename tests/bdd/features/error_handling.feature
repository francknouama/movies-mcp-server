Feature: Error Handling
  As an MCP client
  I want proper error handling for all edge cases
  So that I can handle failures gracefully and get meaningful feedback

  Background:
    Given the MCP server is running
    And the MCP connection is initialized

  @error-handling @database
  Scenario: Database Connection Lost
    Given the database is available
    And I have some movies in the database
    When the database connection is lost
    And I try to add a movie with title "Test Movie"
    Then I should get error code -32603
    And the error message should contain "database unavailable"
    And the error should include retry guidance

  @error-handling @validation
  Scenario: Invalid Input Validation
    Given the database is clean
    When I try to add a movie with invalid data:
      """
      {
        "title": "",
        "director": "",
        "year": "invalid_year",
        "rating": 15,
        "genre": null
      }
      """
    Then I should get error code -32602
    And the error should contain validation errors for:
      | field    | issue                           |
      | title    | Title cannot be empty           |
      | director | Director cannot be empty        |
      | year     | Year must be a valid number     |
      | rating   | Rating must be between 0 and 10 |

  @error-handling @concurrency
  Scenario: Concurrent Modification Conflicts
    Given I have a movie with ID 1
    When two clients try to update the same movie simultaneously
    Then one update should succeed
    And the other should get a conflict error
    And the error should suggest retrying the operation

  @error-handling @not-found
  Scenario: Resource Not Found Errors
    Given the database is clean
    When I try to get movie with ID 99999
    Then I should get error code -32602
    And the error message should indicate "movie not found"
    And the error should include the requested ID

  @error-handling @malformed-requests
  Scenario Outline: Malformed Request Handling
    When I send a malformed request: "<request>"
    Then I should get error code <error_code>
    And the error message should indicate "<error_type>"

    Examples:
      | request                           | error_code | error_type        |
      | {"invalid": "json"}               | -32700     | parse error       |
      | {"jsonrpc": "1.0"}                | -32600     | invalid request   |
      | {"method": "unknown_tool"}        | -32601     | method not found  |
      | {"method": "add_movie"}           | -32602     | invalid params    |

  @error-handling @rate-limiting
  Scenario: Rate Limiting Protection
    When I send 1000 requests in 1 second
    Then some requests should be rate limited
    And I should get error code -32099
    And the error message should indicate "rate limit exceeded"
    And the error should include retry-after information

  @error-handling @timeout
  Scenario: Request Timeout Handling
    Given I configure a 1 second timeout
    When I perform an operation that takes 2 seconds
    Then I should get a timeout error
    And the operation should be cancelled
    And resources should be properly cleaned up

  @error-handling @invalid-json
  Scenario: Invalid JSON-RPC Protocol
    When I send invalid JSON-RPC messages:
      | message                                      | expected_error |
      | not json at all                              | parse error    |
      | {"method": "test"}                           | invalid request|
      | {"jsonrpc": "2.0", "method": 123}            | invalid request|
      | {"jsonrpc": "2.0", "id": "test"}             | invalid request|
    Then each should return appropriate error codes
    And the server should remain stable

  @error-handling @boundary-conditions
  Scenario: Boundary Condition Errors
    When I test boundary conditions:
      | field     | value              | expected_error                    |
      | rating    | -1                 | Rating must be between 0 and 10  |
      | rating    | 11                 | Rating must be between 0 and 10  |
      | year      | 1800               | Year must be after 1888           |
      | year      | 2100               | Year must be before 2030          |
      | title     | 500 character long | Title must be under 255 characters|
    Then I should get appropriate validation errors
    And the errors should include the invalid values

  @error-handling @network
  Scenario: Network Error Simulation
    Given the MCP server is running
    When network errors occur during communication
    Then the client should handle connection drops gracefully
    And appropriate error messages should be returned
    And the server should remain responsive

  @error-handling @memory-pressure
  Scenario: Memory Pressure Handling
    When the system is under memory pressure
    And I try to perform memory-intensive operations
    Then the system should fail gracefully
    And return appropriate resource exhaustion errors
    And not crash or become unresponsive

  @error-handling @recovery
  Scenario: Error Recovery Testing
    Given the system has encountered various errors
    When the error conditions are resolved
    Then the system should recover automatically
    And subsequent operations should work normally
    And no residual state should remain from errors

  @error-handling @cascading-failures
  Scenario: Cascading Failure Prevention
    Given multiple system components
    When one component fails
    Then the failure should not cascade to other components
    And the system should maintain partial functionality
    And errors should be isolated and contained

  @error-handling @data-integrity
  Scenario: Data Integrity During Errors
    Given I have important data in the database
    When errors occur during write operations
    Then the data should remain consistent
    And partial writes should be rolled back
    And no data corruption should occur

  @error-handling @long-running-operations
  Scenario: Long Running Operation Interruption
    Given I start a long-running operation
    When the operation is interrupted by an error
    Then the operation should be cleanly cancelled
    And any partial work should be undone
    And resources should be properly released