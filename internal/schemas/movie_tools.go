package schemas

import "movies-mcp-server/internal/interfaces/dto"

// GetMovieTools returns all movie management tool schemas
func GetMovieTools() []dto.Tool {
	return []dto.Tool{
		getMovieTool(),
		addMovieTool(),
		updateMovieTool(),
		deleteMovieTool(),
		listTopMoviesTool(),
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