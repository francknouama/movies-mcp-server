Feature: Advanced Search and Integration
  As an MCP client
  I want to perform complex searches and operations
  So that I can efficiently find and manage movie data

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And the database contains sample movie data

  @search @advanced
  Scenario: Search movies by rating range
    Given the database contains movies with various ratings
    When I call the "search_by_rating_range" tool with:
      """
      {
        "min_rating": 8.5,
        "max_rating": 9.5
      }
      """
    Then the response should be successful
    And all returned movies should have rating between 8.5 and 9.5
    And the movies should be ordered by rating descending

  @search @advanced
  Scenario: Search similar movies
    Given a movie exists with:
      | title    | Inception         |
      | director | Christopher Nolan |
      | genre    | Sci-Fi           |
      | year     | 2010             |
    When I call the "search_similar_movies" tool with the movie ID
    Then the response should be successful
    And the response should contain movies with similar characteristics
    And the original movie should not be included in results

  @search @integration
  Scenario: Complex movie and cast search
    Given the following movies with cast exist:
      | movie_title | director      | actor_name    |
      | Titanic     | James Cameron | Leonardo DiCaprio |
      | Titanic     | James Cameron | Kate Winslet  |
      | Inception   | Christopher Nolan | Leonardo DiCaprio |
    When I search for movies with "Leonardo DiCaprio"
    Then the response should contain 2 movies
    And the movies should be "Titanic" and "Inception"

  @search @performance
  Scenario: Large dataset search performance
    Given the database contains 1000+ movies
    When I call the "search_movies" tool with:
      """
      {
        "genre": "Action",
        "min_year": 2000,
        "limit": 50
      }
      """
    Then the response should be successful
    And the response time should be under 2 seconds
    And the response should contain up to 50 movies

  @resources @integration
  Scenario: Read database statistics resource
    Given the database contains sample data
    When I call the "resources/read" method with URI "movies://database/info"
    Then the response should be successful
    And the response should contain:
      | field        | type    |
      | total_movies | number  |
      | total_actors | number  |
      | genres       | array   |
      | date_range   | object  |

  @resources @integration
  Scenario: Read poster information resource
    When I call the "resources/read" method with URI "movies://posters/info"
    Then the response should be successful
    And the response should contain poster storage statistics

  @integration @workflow
  Scenario: Complete movie management workflow
    # Add a new movie
    When I call the "add_movie" tool with:
      """
      {
        "title": "Test Movie Workflow",
        "director": "Test Director",
        "year": 2024,
        "rating": 8.0,
        "genre": "Drama"
      }
      """
    Then the response should be successful
    And I store the movie ID as "workflow_movie_id"
    
    # Add actors
    When I call the "add_actor" tool with:
      """
      {
        "name": "Test Actor 1",
        "birth_year": 1980
      }
      """
    Then the response should be successful
    And I store the actor ID as "workflow_actor1_id"
    
    When I call the "add_actor" tool with:
      """
      {
        "name": "Test Actor 2",
        "birth_year": 1985
      }
      """
    Then the response should be successful
    And I store the actor ID as "workflow_actor2_id"
    
    # Link actors to movie
    When I link actor "workflow_actor1_id" to movie "workflow_movie_id"
    And I link actor "workflow_actor2_id" to movie "workflow_movie_id"
    
    # Verify cast
    When I call the "get_movie_cast" tool with movie "workflow_movie_id"
    Then the response should contain 2 actors
    
    # Update movie
    When I call the "update_movie" tool with:
      """
      {
        "movie_id": "{workflow_movie_id}",
        "rating": 8.5
      }
      """
    Then the movie rating should be updated to 8.5
    
    # Search for the movie
    When I search for movies with title "Test Movie Workflow"
    Then the response should contain the created movie
    
    # Clean up
    When I delete movie "workflow_movie_id"
    And I delete actor "workflow_actor1_id"
    And I delete actor "workflow_actor2_id"
    Then all test data should be removed

  @error-handling @edge-cases
  Scenario: Concurrent operations handling
    Given multiple clients are connected
    When client 1 and client 2 simultaneously try to:
      | operation   | parameters                    |
      | add_movie   | same movie title and director |
      | add_actor   | same actor name and birth year |
    Then one operation should succeed
    And one operation should handle the conflict gracefully
    And data integrity should be maintained

  @search @pagination
  Scenario: Paginated search results
    Given the database contains 100+ movies
    When I call the "search_movies" tool with:
      """
      {
        "limit": 10,
        "offset": 20
      }
      """
    Then the response should be successful
    And the response should contain 10 movies
    And the response should indicate total available results
    And the results should start from the 21st movie