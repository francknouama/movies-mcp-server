Feature: [Feature with Structured Data]
  Tests that require table data for input or validation

  Background:
    Given I have a valid MCP client connection

  Scenario: [Test with table input]
    Given the following actors exist:
      | name              | birth_year | bio                    |
      | Tom Hanks         | 1956       | American actor         |
      | Meryl Streep      | 1949       | American actress       |
      | Leonardo DiCaprio | 1974       | American actor         |
    When I list all actors
    Then the response should contain 3 actors
    And the response should contain these actors:
      | name              | birth_year |
      | Tom Hanks         | 1956       |
      | Meryl Streep      | 1949       |
      | Leonardo DiCaprio | 1974       |

  Scenario: [Test with table validation]
    Given I have 5 movies in the database
    When I search for movies by rating above 8.0
    Then the response should include:
      | title          | rating |
      | Movie Alpha    | 8.5    |
      | Movie Beta     | 9.0    |
