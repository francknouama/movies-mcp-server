package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

// GetSearchTools returns all search and query tool schemas
func GetSearchTools() []dto.Tool {
	return []dto.Tool{
		searchMoviesTool(),
		searchByDecadeTool(),
		searchByRatingRangeTool(),
		searchSimilarMoviesTool(),
	}
}

func searchMoviesTool() dto.Tool {
	return dto.Tool{
		Name:        "search_movies",
		Description: "Search for movies by various criteria",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]dto.SchemaProperty{
				"title": {
					Type:        "string",
					Description: "Search by title",
				},
				"director": {
					Type:        "string",
					Description: "Search by director",
				},
				"genre": {
					Type:        "string",
					Description: "Search by genre",
				},
				"min_year": {
					Type:        "integer",
					Description: "Minimum release year",
				},
				"max_year": {
					Type:        "integer",
					Description: "Maximum release year",
				},
				"min_rating": {
					Type:        "number",
					Description: "Minimum rating",
				},
				"max_rating": {
					Type:        "number",
					Description: "Maximum rating",
				},
				"limit": {
					Type:        "integer",
					Description: "Maximum number of results",
					Default:     50,
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
			Properties: map[string]dto.SchemaProperty{
				"decade": {
					Type:        "string",
					Description: "Decade to search (e.g., '1990s', '2000s')",
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
			Properties: map[string]dto.SchemaProperty{
				"min_rating": {
					Type:        "number",
					Description: "Minimum rating (0-10)",
					Minimum:     Float64Ptr(0),
					Maximum:     Float64Ptr(10),
				},
				"max_rating": {
					Type:        "number",
					Description: "Maximum rating (0-10)",
					Minimum:     Float64Ptr(0),
					Maximum:     Float64Ptr(10),
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
			Properties: map[string]dto.SchemaProperty{
				"movie_id": {
					Type:        "integer",
					Description: "Reference movie ID",
				},
				"limit": {
					Type:        "integer",
					Description: "Number of similar movies to return",
					Default:     5,
				},
			},
			Required: []string{"movie_id"},
		},
	}
}
