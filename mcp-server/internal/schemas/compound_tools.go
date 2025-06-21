package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

// GetCompoundTools returns all compound operation tool schemas
func GetCompoundTools() []dto.Tool {
	return []dto.Tool{
		bulkMovieImportTool(),
		movieRecommendationEngineTool(),
		directorCareerAnalysisTool(),
	}
}

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