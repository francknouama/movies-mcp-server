Feature: [MCP Protocol Compliance]
  Tests for MCP protocol implementation and compliance

  Background:
    Given I have a valid MCP client connection

  @mcp @protocol
  Scenario: [Initialize request]
    When I send an initialize request with:
      """json
      {
        "protocolVersion": "2024-11-05",
        "capabilities": {
          "tools": {}
        },
        "clientInfo": {
          "name": "test-client",
          "version": "1.0.0"
        }
      }
      """
    Then the response should be successful
    And the response should contain server capabilities
    And the protocol version should be "2024-11-05"

  @mcp @tools
  Scenario: [List available tools]
    When I send a tools/list request
    Then the response should contain the following tools:
      | name         | description                |
      | add_movie    | Add a new movie            |
      | search_movies| Search for movies          |
      | add_actor    | Add a new actor            |

  @mcp @resources
  Scenario: [List available resources]
    When I send a resources/list request
    Then the response should contain the following resources:
      | uri              | name           | description                    |
      | movie://list     | Movies List    | List all available movies      |
      | actor://list     | Actors List    | List all available actors      |

  @mcp @error
  Scenario: [Invalid method error]
    When I send a request with invalid method "unknown_method"
    Then the operation should fail
    And the error code should be -32601
    And the error message should contain "method not found"
