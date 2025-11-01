Feature: [Feature Name]
  As a [user role]
  I want to [action]
  So that [business value]

  Background:
    Given I have a valid MCP client connection

  @smoke @[domain]
  Scenario: [Simple scenario title]
    Given [initial context]
    When [action/event]
    Then [expected outcome]

  @[tag]
  Scenario: [Another scenario title]
    Given [setup step 1]
    And [setup step 2]
    When [action is performed]
    Then [verify outcome 1]
    And [verify outcome 2]
    But [verify negative case]
