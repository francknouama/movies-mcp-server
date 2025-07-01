Feature: MCP Protocol Communication
  As an MCP client
  I want to communicate with the MCP server
  So that I can use the movie database functionality

  Background:
    Given the MCP server is running
    And I have a valid MCP client connection

  @mcp @smoke
  Scenario: Initialize MCP connection
    When I send an initialize request with:
      """
      {
        "method": "initialize",
        "params": {
          "protocolVersion": "2024-11-05",
          "capabilities": {},
          "clientInfo": {
            "name": "test-client",
            "version": "1.0.0"
          }
        }
      }
      """
    Then the response should be successful
    And the response should contain server capabilities
    And the protocol version should be "2024-11-05"

  @mcp
  Scenario: List available tools
    Given the MCP connection is initialized
    When I send a tools/list request
    Then the response should be successful
    And the response should contain the following tools:
      | tool_name              | description                        |
      | get_movie             | Get movie details by ID            |
      | add_movie             | Add a new movie to the database    |
      | update_movie          | Update an existing movie           |
      | delete_movie          | Delete a movie by ID               |
      | search_movies         | Search movies by various criteria  |
      | list_top_movies       | Get top rated movies              |
      | add_actor             | Add a new actor to the database    |
      | link_actor_to_movie   | Link an actor to a movie          |
      | get_movie_cast        | Get cast for a specific movie      |
      | get_actor_movies      | Get movies for a specific actor    |

  @mcp
  Scenario: List available resources
    Given the MCP connection is initialized
    When I send a resources/list request
    Then the response should be successful
    And the response should contain the following resources:
      | uri                    | name              | description                 |
      | movies://database/info | Database Info     | Database statistics         |
      | movies://posters/info  | Poster Info       | Poster storage information  |

  @mcp @error-handling
  Scenario: Invalid method call
    Given the MCP connection is initialized
    When I send a request with invalid method "invalid_method"
    Then the response should contain an error
    And the error code should be -32601
    And the error message should contain "Method not found"

  @mcp @error-handling
  Scenario: Invalid protocol version
    When I send an initialize request with protocol version "1.0.0"
    Then the response should contain an error
    And the error should indicate unsupported protocol version