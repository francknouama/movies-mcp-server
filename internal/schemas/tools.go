package schemas

import "movies-mcp-server/internal/interfaces/dto"

// GetToolSchemas returns all tool schemas for the MCP server
func GetToolSchemas() []dto.Tool {
	return []dto.Tool{
		// Movie Management Tools
		getMovieTool(),
		addMovieTool(),
		updateMovieTool(),
		deleteMovieTool(),
		searchMoviesTool(),
		listTopMoviesTool(),
		
		// Enhanced Search Tools
		searchByDecadeTool(),
		searchByRatingRangeTool(),
		searchSimilarMoviesTool(),
		
		// Actor Management Tools
		addActorTool(),
		getActorTool(),
		updateActorTool(),
		deleteActorTool(),
		linkActorToMovieTool(),
		unlinkActorFromMovieTool(),
		getMovieCastTool(),
		getActorMoviesTool(),
		searchActorsTool(),
		
		// Compound Tools
		bulkMovieImportTool(),
		movieRecommendationEngineTool(),
		directorCareerAnalysisTool(),
		
		// Context Management Tools
		createSearchContextTool(),
		getContextPageTool(),
		getContextInfoTool(),
		
		// Validation Tools
		validateToolCallTool(),
	}
}

func getMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "get_movie",
		Description: "Get a movie by ID",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "The movie ID",
				},
			},
			Required: []string{"movie_id"},
		},
	}
}

func addMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "add_movie",
		Description: "Add a new movie to the database",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Movie title",
				},
				"director": map[string]interface{}{
					"type":        "string",
					"description": "Movie director",
				},
				"year": map[string]interface{}{
					"type":        "integer",
					"description": "Release year",
				},
				"rating": map[string]interface{}{
					"type":        "number",
					"description": "Movie rating (0-10)",
					"minimum":     0,
					"maximum":     10,
				},
				"genres": map[string]interface{}{
					"type":        "array",
					"description": "List of genres",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"poster_url": map[string]interface{}{
					"type":        "string",
					"description": "URL to movie poster",
				},
			},
			Required: []string{"title", "director", "year"},
		},
	}
}

func updateMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "update_movie",
		Description: "Update an existing movie",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "Movie ID",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Movie title",
				},
				"director": map[string]interface{}{
					"type":        "string",
					"description": "Movie director",
				},
				"year": map[string]interface{}{
					"type":        "integer",
					"description": "Release year",
				},
				"rating": map[string]interface{}{
					"type":        "number",
					"description": "Movie rating (0-10)",
					"minimum":     0,
					"maximum":     10,
				},
				"genres": map[string]interface{}{
					"type":        "array",
					"description": "List of genres",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"poster_url": map[string]interface{}{
					"type":        "string",
					"description": "URL to movie poster",
				},
			},
			Required: []string{"id", "title", "director", "year"},
		},
	}
}

func deleteMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "delete_movie",
		Description: "Delete a movie by ID",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "The movie ID to delete",
				},
			},
			Required: []string{"movie_id"},
		},
	}
}

func searchMoviesTool() dto.Tool {
	return dto.Tool{
		Name:        "search_movies",
		Description: "Search for movies by various criteria",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Search by title",
				},
				"director": map[string]interface{}{
					"type":        "string",
					"description": "Search by director",
				},
				"genre": map[string]interface{}{
					"type":        "string",
					"description": "Search by genre",
				},
				"min_year": map[string]interface{}{
					"type":        "integer",
					"description": "Minimum release year",
				},
				"max_year": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum release year",
				},
				"min_rating": map[string]interface{}{
					"type":        "number",
					"description": "Minimum rating",
				},
				"max_rating": map[string]interface{}{
					"type":        "number",
					"description": "Maximum rating",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results",
					"default":     50,
				},
			},
			Required: []string{},
		},
	}
}

func listTopMoviesTool() dto.Tool {
	return dto.Tool{
		Name:        "list_top_movies",
		Description: "Get top-rated movies",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Number of movies to return",
					"default":     10,
				},
			},
			Required: []string{},
		},
	}
}

func searchByDecadeTool() dto.Tool {
	return dto.Tool{
		Name:        "search_by_decade",
		Description: "Search movies by decade (e.g., '1990s', '2000s')",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"decade": map[string]interface{}{
					"type":        "string",
					"description": "Decade to search (e.g., '1990s', '2000s')",
				},
			},
			Required: []string{"decade"},
		},
	}
}

func searchByRatingRangeTool() dto.Tool {
	return dto.Tool{
		Name:        "search_by_rating_range",
		Description: "Search movies within a specific rating range",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"min_rating": map[string]interface{}{
					"type":        "number",
					"description": "Minimum rating (0-10)",
					"minimum":     0,
					"maximum":     10,
				},
				"max_rating": map[string]interface{}{
					"type":        "number",
					"description": "Maximum rating (0-10)",
					"minimum":     0,
					"maximum":     10,
				},
			},
			Required: []string{"min_rating", "max_rating"},
		},
	}
}

func searchSimilarMoviesTool() dto.Tool {
	return dto.Tool{
		Name:        "search_similar_movies",
		Description: "Find movies similar to a given movie",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "Reference movie ID",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Number of similar movies to return",
					"default":     5,
				},
			},
			Required: []string{"movie_id"},
		},
	}
}

// Actor tools
func addActorTool() dto.Tool {
	return dto.Tool{
		Name:        "add_actor",
		Description: "Add a new actor to the database",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Actor name",
				},
				"birth_year": map[string]interface{}{
					"type":        "integer",
					"description": "Birth year",
				},
				"bio": map[string]interface{}{
					"type":        "string",
					"description": "Actor biography",
				},
			},
			Required: []string{"name", "birth_year"},
		},
	}
}

func getActorTool() dto.Tool {
	return dto.Tool{
		Name:        "get_actor",
		Description: "Get an actor by ID",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"actor_id": map[string]interface{}{
					"type":        "integer",
					"description": "The actor ID",
				},
			},
			Required: []string{"actor_id"},
		},
	}
}

func updateActorTool() dto.Tool {
	return dto.Tool{
		Name:        "update_actor",
		Description: "Update an existing actor",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "Actor ID",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Actor name",
				},
				"birth_year": map[string]interface{}{
					"type":        "integer",
					"description": "Birth year",
				},
				"bio": map[string]interface{}{
					"type":        "string",
					"description": "Actor biography",
				},
			},
			Required: []string{"id", "name", "birth_year"},
		},
	}
}

func deleteActorTool() dto.Tool {
	return dto.Tool{
		Name:        "delete_actor",
		Description: "Delete an actor by ID",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"actor_id": map[string]interface{}{
					"type":        "integer",
					"description": "The actor ID to delete",
				},
			},
			Required: []string{"actor_id"},
		},
	}
}

func linkActorToMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "link_actor_to_movie",
		Description: "Link an actor to a movie",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"actor_id": map[string]interface{}{
					"type":        "integer",
					"description": "Actor ID",
				},
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "Movie ID",
				},
			},
			Required: []string{"actor_id", "movie_id"},
		},
	}
}

func unlinkActorFromMovieTool() dto.Tool {
	return dto.Tool{
		Name:        "unlink_actor_from_movie",
		Description: "Unlink an actor from a movie",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"actor_id": map[string]interface{}{
					"type":        "integer",
					"description": "Actor ID",
				},
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "Movie ID",
				},
			},
			Required: []string{"actor_id", "movie_id"},
		},
	}
}

func getMovieCastTool() dto.Tool {
	return dto.Tool{
		Name:        "get_movie_cast",
		Description: "Get all actors in a movie",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"movie_id": map[string]interface{}{
					"type":        "integer",
					"description": "Movie ID",
				},
			},
			Required: []string{"movie_id"},
		},
	}
}

func getActorMoviesTool() dto.Tool {
	return dto.Tool{
		Name:        "get_actor_movies",
		Description: "Get all movies an actor appeared in",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"actor_id": map[string]interface{}{
					"type":        "integer",
					"description": "Actor ID",
				},
			},
			Required: []string{"actor_id"},
		},
	}
}

func searchActorsTool() dto.Tool {
	return dto.Tool{
		Name:        "search_actors",
		Description: "Search for actors by name",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Actor name to search for",
				},
			},
			Required: []string{"name"},
		},
	}
}

// Compound tools
func bulkMovieImportTool() dto.Tool {
	return dto.Tool{
		Name:        "bulk_movie_import",
		Description: "Import multiple movies at once",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"movies": map[string]interface{}{
					"type":        "array",
					"description": "Array of movies to import",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type": "string",
							},
							"director": map[string]interface{}{
								"type": "string",
							},
							"year": map[string]interface{}{
								"type": "integer",
							},
						},
						"required": []string{"title", "director", "year"},
					},
				},
			},
			Required: []string{"movies"},
		},
	}
}

func movieRecommendationEngineTool() dto.Tool {
	return dto.Tool{
		Name:        "movie_recommendation_engine",
		Description: "Get movie recommendations based on preferences",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"user_preferences": map[string]interface{}{
					"type":        "object",
					"description": "User preferences for recommendations",
				},
			},
			Required: []string{"user_preferences"},
		},
	}
}

func directorCareerAnalysisTool() dto.Tool {
	return dto.Tool{
		Name:        "director_career_analysis",
		Description: "Analyze a director's career trajectory",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"director_name": map[string]interface{}{
					"type":        "string",
					"description": "Director name to analyze",
				},
			},
			Required: []string{"director_name"},
		},
	}
}

// Context management tools
func createSearchContextTool() dto.Tool {
	return dto.Tool{
		Name:        "create_search_context",
		Description: "Create a search context for large result sets",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"search_criteria": map[string]interface{}{
					"type":        "object",
					"description": "Search criteria",
				},
			},
			Required: []string{"search_criteria"},
		},
	}
}

func getContextPageTool() dto.Tool {
	return dto.Tool{
		Name:        "get_context_page",
		Description: "Get a page of results from a search context",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"context_id": map[string]interface{}{
					"type":        "string",
					"description": "Context ID",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number",
				},
			},
			Required: []string{"context_id", "page"},
		},
	}
}

func getContextInfoTool() dto.Tool {
	return dto.Tool{
		Name:        "get_context_info",
		Description: "Get information about a search context",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"context_id": map[string]interface{}{
					"type":        "string",
					"description": "Context ID",
				},
			},
			Required: []string{"context_id"},
		},
	}
}

// Validation tools
func validateToolCallTool() dto.Tool {
	return dto.Tool{
		Name:        "validate_tool_call",
		Description: "Validate a tool call against its schema",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tool_name": map[string]interface{}{
					"type":        "string",
					"description": "Tool name to validate",
				},
				"arguments": map[string]interface{}{
					"type":        "object",
					"description": "Arguments to validate",
				},
			},
			Required: []string{"tool_name", "arguments"},
		},
	}
}