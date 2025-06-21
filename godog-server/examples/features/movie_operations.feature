Feature: Movie Operations
  As an MCP client
  I want to manage movies in the database
  So that I can maintain a comprehensive movie collection

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And the database is clean

  @movies @crud
  Scenario: Add a new movie
    When I call the "add_movie" tool with:
      """
      {
        "title": "The Shawshank Redemption",
        "director": "Frank Darabont",
        "year": 1994,
        "rating": 9.3,
        "genre": "Drama"
      }
      """
    Then the response should be successful
    And the response should contain a movie with:
      | field    | value                     |
      | title    | The Shawshank Redemption  |
      | director | Frank Darabont            |
      | year     | 1994                      |
      | rating   | 9.3                       |
    And the movie should have an assigned ID

  @movies @crud
  Scenario: Get movie details by ID
    Given a movie exists with:
      | title    | The Godfather    |
      | director | Francis Coppola  |
      | year     | 1972             |
      | rating   | 9.2              |
    When I call the "get_movie" tool with the movie ID
    Then the response should be successful
    And the response should contain the movie details
    And the movie title should be "The Godfather"

  @movies @crud
  Scenario: Update movie information
    Given a movie exists with title "Test Movie"
    When I call the "update_movie" tool with:
      """
      {
        "movie_id": "{movie_id}",
        "title": "Updated Test Movie",
        "rating": 8.5
      }
      """
    Then the response should be successful
    And the movie title should be updated to "Updated Test Movie"
    And the movie rating should be updated to 8.5

  @movies @crud
  Scenario: Delete a movie
    Given a movie exists with title "Movie to Delete"
    When I call the "delete_movie" tool with the movie ID
    Then the response should be successful
    And the movie should no longer exist in the database

  @movies @search
  Scenario: Search movies by title
    Given the following movies exist:
      | title           | director        | year | rating |
      | Inception       | Christopher Nolan | 2010 | 8.8    |
      | Interstellar    | Christopher Nolan | 2014 | 8.6    |
      | The Dark Knight | Christopher Nolan | 2008 | 9.0    |
    When I call the "search_movies" tool with:
      """
      {
        "title": "Inter"
      }
      """
    Then the response should be successful
    And the response should contain 1 movie
    And the movie title should contain "Interstellar"

  @movies @search
  Scenario: Search movies by director
    Given the following movies exist:
      | title           | director        | year | rating |
      | Inception       | Christopher Nolan | 2010 | 8.8    |
      | Pulp Fiction    | Quentin Tarantino | 1994 | 8.9    |
      | Kill Bill       | Quentin Tarantino | 2003 | 8.1    |
    When I call the "search_movies" tool with:
      """
      {
        "director": "Quentin Tarantino"
      }
      """
    Then the response should be successful
    And the response should contain 2 movies
    And all movies should have director "Quentin Tarantino"

  @movies @search
  Scenario: Get top rated movies
    Given the following movies exist:
      | title           | rating |
      | Movie A         | 9.5    |
      | Movie B         | 8.0    |
      | Movie C         | 9.2    |
      | Movie D         | 7.5    |
    When I call the "list_top_movies" tool with limit 2
    Then the response should be successful
    And the response should contain 2 movies
    And the movies should be ordered by rating descending
    And the first movie should have rating 9.5

  @movies @error-handling
  Scenario: Get non-existent movie
    When I call the "get_movie" tool with movie ID 99999
    Then the response should contain an error
    And the error message should indicate movie not found

  @movies @error-handling
  Scenario: Add movie with invalid data
    When I call the "add_movie" tool with:
      """
      {
        "title": "",
        "year": "invalid_year",
        "rating": 15
      }
      """
    Then the response should contain an error
    And the error should contain validation errors for:
      | field  | issue                           |
      | title  | Title cannot be empty           |
      | year   | Year must be a valid number     |
      | rating | Rating must be between 0 and 10 |

  @movies @search
  Scenario Outline: Search movies by decade
    Given movies exist from various decades
    When I call the "search_by_decade" tool with decade "<decade>"
    Then the response should be successful
    And all movies should be from years <min_year> to <max_year>

    Examples:
      | decade | min_year | max_year |
      | 1990s  | 1990     | 1999     |
      | 2000s  | 2000     | 2009     |
      | 2010s  | 2010     | 2019     |