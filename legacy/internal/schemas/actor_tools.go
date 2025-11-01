package schemas

import "github.com/francknouama/movies-mcp-server/internal/interfaces/dto"

// GetActorTools returns all actor management tool schemas
func GetActorTools() []dto.Tool {
	return []dto.Tool{
		addActorTool(),
		getActorTool(),
		updateActorTool(),
		deleteActorTool(),
		linkActorToMovieTool(),
		unlinkActorFromMovieTool(),
		getMovieCastTool(),
		getActorMoviesTool(),
		searchActorsTool(),
	}
}

func addActorTool() dto.Tool {
	return dto.Tool{
		Name:        "add_actor",
		Description: "Add a new actor to the database",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]dto.SchemaProperty{
				"name": {
					Type:        "string",
					Description: "Actor name",
				},
				"birth_year": {
					Type:        "integer",
					Description: "Birth year",
				},
				"bio": {
					Type:        "string",
					Description: "Actor biography",
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
			Properties: map[string]dto.SchemaProperty{
				"actor_id": {
					Type:        "integer",
					Description: "The actor ID",
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
			Properties: map[string]dto.SchemaProperty{
				"id": {
					Type:        "integer",
					Description: "Actor ID",
				},
				"name": {
					Type:        "string",
					Description: "Actor name",
				},
				"birth_year": {
					Type:        "integer",
					Description: "Birth year",
				},
				"bio": {
					Type:        "string",
					Description: "Actor biography",
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
			Properties: map[string]dto.SchemaProperty{
				"actor_id": {
					Type:        "integer",
					Description: "The actor ID to delete",
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
			Properties: map[string]dto.SchemaProperty{
				"actor_id": {
					Type:        "integer",
					Description: "Actor ID",
				},
				"movie_id": {
					Type:        "integer",
					Description: "Movie ID",
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
			Properties: map[string]dto.SchemaProperty{
				"actor_id": {
					Type:        "integer",
					Description: "Actor ID",
				},
				"movie_id": {
					Type:        "integer",
					Description: "Movie ID",
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
			Properties: map[string]dto.SchemaProperty{
				"movie_id": {
					Type:        "integer",
					Description: "Movie ID",
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
			Properties: map[string]dto.SchemaProperty{
				"actor_id": {
					Type:        "integer",
					Description: "Actor ID",
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
			Properties: map[string]dto.SchemaProperty{
				"name": {
					Type:        "string",
					Description: "Actor name to search for",
				},
			},
			Required: []string{"name"},
		},
	}
}
