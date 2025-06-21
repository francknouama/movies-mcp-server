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
