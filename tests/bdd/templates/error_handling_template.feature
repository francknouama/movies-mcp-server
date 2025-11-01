Feature: [Error Handling and Validation]
  Tests for error scenarios and validation rules

  Background:
    Given I have a valid MCP client connection

  @error @validation
  Scenario: [Missing required field]
    When I try to create a movie without a title
    Then the operation should fail
    And the error message should contain "title is required"
    And the error code should be 400

  @error @notfound
  Scenario: [Entity not found]
    When I try to retrieve actor with ID 99999
    Then the operation should fail
    And the error code should be 404
    And the error message should contain "not found"

  @error @duplicate
  Scenario: [Duplicate entry]
    Given I have a movie with title "The Matrix"
    When I try to create another movie with title "The Matrix"
    Then the operation should fail
    And the error message should contain "already exists"

  @error @validation
  Scenario Outline: [Invalid data validation]
    When I try to create an actor with <field> set to <invalid_value>
    Then the operation should fail
    And the error message should contain "<expected_message>"

    Examples:
      | field      | invalid_value | expected_message       |
      | name       | ""            | name is required       |
      | birth_year | 2050          | invalid birth year     |
      | birth_year | -100          | invalid birth year     |
