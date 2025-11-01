Feature: [Feature with Multiple Test Cases]
  Use scenario outlines to test multiple variations of the same behavior

  Background:
    Given I have a valid MCP client connection

  @parameterized
  Scenario Outline: [Test with different parameters]
    Given I have a movie with title "<title>"
    And the movie has rating <rating>
    When I search for movies rated above <threshold>
    Then the movie should be <included_or_not>

    Examples:
      | title           | rating | threshold | included_or_not |
      | High Rated      | 9.0    | 8.0       | included        |
      | Medium Rated    | 7.5    | 8.0       | not included    |
      | Low Rated       | 6.0    | 8.0       | not included    |

  @validation
  Scenario Outline: [Validation with examples]
    When I try to create an actor with name "<name>" and birth year <year>
    Then the operation should <result>
    And the response should contain "<message>"

    Examples:
      | name        | year | result  | message            |
      | Valid Name  | 1980 | succeed | created            |
      |             | 1980 | fail    | name is required   |
      | Valid Name  | 2030 | fail    | future birth year  |
      | Valid Name  | 1800 | fail    | unrealistic year   |
