Feature: Actor Operations
  As an MCP client
  I want to manage actors and their movie relationships
  So that I can maintain comprehensive cast information

  Background:
    Given the MCP server is running
    And the MCP connection is initialized
    And the database is clean

  @actors @crud
  Scenario: Add a new actor
    When I call the "add_actor" tool with:
      """
      {
        "name": "Leonardo DiCaprio",
        "birth_year": 1974,
        "bio": "American actor and film producer"
      }
      """
    Then the response should be successful
    And the response should contain an actor with:
      | field      | value                          |
      | name       | Leonardo DiCaprio              |
      | birth_year | 1974                           |
      | bio        | American actor and film producer |
    And the actor should have an assigned ID

  @actors @crud
  Scenario: Get actor details by ID
    Given an actor exists with:
      | name       | Morgan Freeman |
      | birth_year | 1937          |
      | bio        | American actor |
    When I call the "get_actor" tool with the actor ID
    Then the response should be successful
    And the response should contain the actor details
    And the actor name should be "Morgan Freeman"

  @actors @crud
  Scenario: Update actor information
    Given an actor exists with name "Test Actor"
    When I call the "update_actor" tool with:
      """
      {
        "actor_id": "{actor_id}",
        "name": "Updated Test Actor",
        "bio": "Updated biography"
      }
      """
    Then the response should be successful
    And the actor name should be updated to "Updated Test Actor"
    And the actor bio should be updated to "Updated biography"

  @actors @crud
  Scenario: Delete an actor
    Given an actor exists with name "Actor to Delete"
    When I call the "delete_actor" tool with the actor ID
    Then the response should be successful
    And the actor should no longer exist in the database

  @actors @relationships
  Scenario: Link actor to movie
    Given an actor exists with name "Tom Hanks"
    And a movie exists with title "Forrest Gump"
    When I call the "link_actor_to_movie" tool with:
      """
      {
        "actor_id": "{actor_id}",
        "movie_id": "{movie_id}"
      }
      """
    Then the response should be successful
    And the actor should be linked to the movie
    And the message should indicate successful linking

  @actors @relationships
  Scenario: Get movie cast
    Given a movie exists with title "The Avengers"
    And the following actors are linked to the movie:
      | name           |
      | Robert Downey Jr. |
      | Chris Evans    |
      | Scarlett Johansson |
    When I call the "get_movie_cast" tool with the movie ID
    Then the response should be successful
    And the response should contain 3 actors
    And the cast should include:
      | name           |
      | Robert Downey Jr. |
      | Chris Evans    |
      | Scarlett Johansson |

  @actors @relationships
  Scenario: Get actor's movies
    Given an actor exists with name "Brad Pitt"
    And the actor is linked to the following movies:
      | title        |
      | Fight Club   |
      | Ocean's Eleven |
      | Moneyball    |
    When I call the "get_actor_movies" tool with the actor ID
    Then the response should be successful
    And the response should contain 3 movie IDs
    And the actor should be associated with all linked movies

  @actors @search
  Scenario: Search actors by name
    Given the following actors exist:
      | name              | birth_year |
      | Chris Hemsworth   | 1983       |
      | Chris Evans       | 1981       |
      | Chris Pratt       | 1979       |
      | Christian Bale    | 1974       |
    When I call the "search_actors" tool with:
      """
      {
        "name": "Chris"
      }
      """
    Then the response should be successful
    And the response should contain 4 actors
    And all actor names should contain "Chris"

  @actors @search
  Scenario: Search actors by birth year range
    Given the following actors exist:
      | name           | birth_year |
      | Actor A        | 1970       |
      | Actor B        | 1980       |
      | Actor C        | 1990       |
      | Actor D        | 2000       |
    When I call the "search_actors" tool with:
      """
      {
        "min_birth_year": 1975,
        "max_birth_year": 1995
      }
      """
    Then the response should be successful
    And the response should contain 2 actors
    And all actors should have birth year between 1975 and 1995

  @actors @relationships
  Scenario: Unlink actor from movie
    Given an actor exists with name "Test Actor"
    And a movie exists with title "Test Movie"
    And the actor is linked to the movie
    When I call the "unlink_actor_from_movie" tool with:
      """
      {
        "actor_id": "{actor_id}",
        "movie_id": "{movie_id}"
      }
      """
    Then the response should be successful
    And the actor should no longer be linked to the movie

  @actors @error-handling
  Scenario: Link actor to non-existent movie
    Given an actor exists with name "Test Actor"
    When I call the "link_actor_to_movie" tool with:
      """
      {
        "actor_id": "{actor_id}",
        "movie_id": 99999
      }
      """
    Then the response should contain an error
    And the error message should indicate movie not found

  @actors @error-handling
  Scenario: Add actor with invalid data
    When I call the "add_actor" tool with:
      """
      {
        "name": "",
        "birth_year": "invalid_year"
      }
      """
    Then the response should contain an error
    And the error should contain validation errors for:
      | field      | issue                       |
      | name       | Name cannot be empty        |
      | birth_year | Birth year must be a number |

  @actors @error-handling
  Scenario: Duplicate actor-movie relationship
    Given an actor exists with name "Test Actor"
    And a movie exists with title "Test Movie"
    And the actor is already linked to the movie
    When I call the "link_actor_to_movie" tool with the same actor and movie
    Then the response should contain an error
    And the error message should indicate the relationship already exists