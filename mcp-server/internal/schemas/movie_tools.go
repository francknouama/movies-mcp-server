package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

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
			Properties: map[string]dto.SchemaProperty{
				"movie_id": {
					Type:        "integer",
					Description: "The movie ID",
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
			Properties: map[string]dto.SchemaProperty{
				"title": {
					Type:        "string",
					Description: "Movie title",
				},
				"director": {
					Type:        "string",
					Description: "Movie director",
				},
				"year": {
					Type:        "integer",
					Description: "Release year",
				},
				"rating": {
					Type:        "number",
					Description: "Movie rating (0-10)",
					Minimum:     Float64Ptr(0),
					Maximum:     Float64Ptr(10),
				},
				"genres": {
					Type:        "array",
					Description: "List of genres",
					Items:       StringArrayItems(),
				},
				"poster_url": {
					Type:        "string",
					Description: "URL to movie poster",
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
			Properties: map[string]dto.SchemaProperty{
				"id": {
					Type:        "integer",
					Description: "Movie ID",
				},
				"title": {
					Type:        "string",
					Description: "Movie title",
				},
				"director": {
					Type:        "string",
					Description: "Movie director",
				},
				"year": {
					Type:        "integer",
					Description: "Release year",
				},
				"rating": {
					Type:        "number",
					Description: "Movie rating (0-10)",
					Minimum:     Float64Ptr(0),
					Maximum:     Float64Ptr(10),
				},
				"genres": {
					Type:        "array",
					Description: "List of genres",
					Items:       StringArrayItems(),
				},
				"poster_url": {
					Type:        "string",
					Description: "URL to movie poster",
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
			Properties: map[string]dto.SchemaProperty{
				"movie_id": {
					Type:        "integer",
					Description: "The movie ID to delete",
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
			Properties: map[string]dto.SchemaProperty{
				"limit": {
					Type:        "integer",
					Description: "Number of movies to return",
					Default:     10,
				},
			},
			Required: []string{},
		},
	}
}
