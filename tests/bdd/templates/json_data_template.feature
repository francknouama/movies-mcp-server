Feature: [Feature with JSON Data]
  Tests that work with JSON request/response data

  Background:
    Given I have a valid MCP client connection

  @json
  Scenario: [Create entity with JSON]
    When I send the following movie data:
      """json
      {
        "title": "Inception",
        "director": "Christopher Nolan",
        "year": 2010,
        "rating": 8.8,
        "genres": ["Action", "Sci-Fi", "Thriller"]
      }
      """
    Then the movie should be created successfully
    And the response should have the following structure:
      """json
      {
        "id": 1,
        "title": "Inception",
        "director": "Christopher Nolan",
        "year": 2010,
        "rating": 8.8,
        "genres": ["Action", "Sci-Fi", "Thriller"]
      }
      """

  @json
  Scenario: [Complex JSON validation]
    Given I have a movie database
    When I request movie details with complex filters:
      """json
      {
        "filters": {
          "year": {"min": 2000, "max": 2020},
          "rating": {"min": 8.0},
          "genres": ["Sci-Fi", "Thriller"]
        },
        "sort": {"field": "rating", "order": "desc"},
        "limit": 10
      }
      """
    Then the response should match the criteria
