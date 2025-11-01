package schemas

import "github.com/francknouama/movies-mcp-server/internal/interfaces/dto"

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
			Properties: map[string]dto.SchemaProperty{
				"movies": {
					Type:        "array",
					Description: "Array of movies to import",
					Items: &dto.SchemaProperty{
						Type: "object",
						Properties: map[string]dto.SchemaProperty{
							"title": {
								Type: "string",
							},
							"director": {
								Type: "string",
							},
							"year": {
								Type: "integer",
							},
						},
						Required: []string{"title", "director", "year"},
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
			Properties: map[string]dto.SchemaProperty{
				"user_preferences": {
					Type:        "object",
					Description: "User preferences for recommendations",
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
			Properties: map[string]dto.SchemaProperty{
				"director_name": {
					Type:        "string",
					Description: "Director name to analyze",
				},
			},
			Required: []string{"director_name"},
		},
	}
}
