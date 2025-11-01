package mcp

import (
	"fmt"
	"strings"

	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

// PromptHandlers handles MCP prompt-related requests
type PromptHandlers struct{}

// NewPromptHandlers creates a new PromptHandlers instance
func NewPromptHandlers() *PromptHandlers {
	return &PromptHandlers{}
}

// GetPrompts returns all available prompt templates
func (h *PromptHandlers) GetPrompts() []dto.Prompt {
	return h.getPromptTemplates()
}

// getPromptTemplates returns all available prompt templates
func (h *PromptHandlers) getPromptTemplates() []dto.Prompt {
	return []dto.Prompt{
		{
			Name:        "movie_recommendation",
			Description: "Generate a movie recommendation based on preferences",
			Arguments: []dto.PromptArgument{
				{
					Name:        "genre",
					Description: "Preferred movie genre",
					Required:    true,
				},
				{
					Name:        "min_rating",
					Description: "Minimum rating (0-10)",
					Required:    false,
				},
				{
					Name:        "year_range",
					Description: "Year range (e.g., '2020-2024')",
					Required:    false,
				},
			},
		},
		{
			Name:        "movie_analysis",
			Description: "Analyze a movie's themes and characteristics",
			Arguments: []dto.PromptArgument{
				{
					Name:        "movie_title",
					Description: "Title of the movie to analyze",
					Required:    true,
				},
				{
					Name:        "aspects",
					Description: "Specific aspects to analyze (themes, characters, plot)",
					Required:    false,
				},
			},
		},
		{
			Name:        "director_filmography",
			Description: "Explore a director's filmography and style",
			Arguments: []dto.PromptArgument{
				{
					Name:        "director_name",
					Description: "Name of the director",
					Required:    true,
				},
				{
					Name:        "focus_period",
					Description: "Specific time period to focus on",
					Required:    false,
				},
			},
		},
		{
			Name:        "genre_exploration",
			Description: "Deep dive into a specific movie genre",
			Arguments: []dto.PromptArgument{
				{
					Name:        "genre",
					Description: "Genre to explore",
					Required:    true,
				},
				{
					Name:        "sub_genre",
					Description: "Specific sub-genre",
					Required:    false,
				},
			},
		},
		{
			Name:        "movie_comparison",
			Description: "Compare and contrast multiple movies",
			Arguments: []dto.PromptArgument{
				{
					Name:        "movie1",
					Description: "First movie title",
					Required:    true,
				},
				{
					Name:        "movie2",
					Description: "Second movie title",
					Required:    true,
				},
				{
					Name:        "comparison_aspects",
					Description: "Aspects to compare (themes, style, impact)",
					Required:    false,
				},
			},
		},
	}
}

// HandlePromptsList handles the prompts/list request
func (h *PromptHandlers) HandlePromptsList(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	prompts := h.getPromptTemplates()

	response := dto.PromptsListResponse{
		Prompts: prompts,
	}

	sendResult(id, response)
}

// HandlePromptGet handles the prompts/get request
func (h *PromptHandlers) HandlePromptGet(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	// Parse request
	name, ok := arguments["name"].(string)
	if !ok || name == "" {
		sendError(id, dto.InvalidParams, "Missing or invalid prompt name", nil)
		return
	}

	args, ok := arguments["arguments"].(map[string]interface{})
	if !ok {
		args = make(map[string]interface{}) // Default to empty if not provided
	}

	// Find the prompt template
	var promptTemplate *dto.Prompt
	for _, p := range h.getPromptTemplates() {
		if p.Name == name {
			promptTemplate = &p
			break
		}
	}

	if promptTemplate == nil {
		sendError(id, dto.InvalidParams, fmt.Sprintf("Unknown prompt: %s", name), nil)
		return
	}

	// Validate required arguments
	for _, arg := range promptTemplate.Arguments {
		if arg.Required {
			if _, exists := args[arg.Name]; !exists {
				sendError(id, dto.InvalidParams,
					fmt.Sprintf("Missing required argument: %s", arg.Name), nil)
				return
			}
		}
	}

	// Generate the prompt response based on the template
	response := h.generatePromptResponse(name, args)
	sendResult(id, response)
}

// generatePromptResponse generates a prompt response based on the template and arguments
func (h *PromptHandlers) generatePromptResponse(name string, args map[string]interface{}) dto.PromptGetResponse {
	switch name {
	case "movie_recommendation":
		return h.generateMovieRecommendationPrompt(args)
	case "movie_analysis":
		return h.generateMovieAnalysisPrompt(args)
	case "director_filmography":
		return h.generateDirectorFilmographyPrompt(args)
	case "genre_exploration":
		return h.generateGenreExplorationPrompt(args)
	case "movie_comparison":
		return h.generateMovieComparisonPrompt(args)
	default:
		return dto.PromptGetResponse{
			Description: "Unknown prompt template",
			Messages:    []dto.PromptMessage{},
		}
	}
}

// generateMovieRecommendationPrompt generates a movie recommendation prompt
func (h *PromptHandlers) generateMovieRecommendationPrompt(args map[string]interface{}) dto.PromptGetResponse {
	genre, ok := args["genre"].(string)
	if !ok {
		genre = "any"
	}
	minRating, ok := args["min_rating"].(float64)
	if !ok {
		minRating = 0.0
	}
	yearRange, ok := args["year_range"].(string)
	if !ok {
		yearRange = ""
	}

	var promptText strings.Builder
	promptText.WriteString(fmt.Sprintf("I need movie recommendations in the %s genre", genre))

	if minRating > 0 {
		promptText.WriteString(fmt.Sprintf(" with a minimum rating of %.1f", minRating))
	}

	if yearRange != "" {
		promptText.WriteString(fmt.Sprintf(" from the period %s", yearRange))
	}

	promptText.WriteString(". Please search for movies matching these criteria and provide detailed recommendations including plot summaries, key themes, and why each movie is worth watching.")

	return dto.PromptGetResponse{
		Description: "Movie recommendation prompt",
		Messages: []dto.PromptMessage{
			{
				Role: "user",
				Content: []dto.ContentBlock{
					{
						Type: "text",
						Text: promptText.String(),
					},
				},
			},
		},
	}
}

// generateMovieAnalysisPrompt generates a movie analysis prompt
func (h *PromptHandlers) generateMovieAnalysisPrompt(args map[string]interface{}) dto.PromptGetResponse {
	movieTitle, ok := args["movie_title"].(string)
	if !ok {
		movieTitle = ""
	}
	aspects, ok := args["aspects"].(string)
	if !ok {
		aspects = ""
	}

	var promptText strings.Builder
	promptText.WriteString(fmt.Sprintf("Please analyze the movie '%s'", movieTitle))

	if aspects != "" {
		promptText.WriteString(fmt.Sprintf(" focusing on %s", aspects))
	} else {
		promptText.WriteString(" covering themes, character development, cinematography, and cultural impact")
	}

	promptText.WriteString(". First, search for this movie in the database to get accurate information, then provide a comprehensive analysis.")

	return dto.PromptGetResponse{
		Description: "Movie analysis prompt",
		Messages: []dto.PromptMessage{
			{
				Role: "user",
				Content: []dto.ContentBlock{
					{
						Type: "text",
						Text: promptText.String(),
					},
				},
			},
		},
	}
}

// generatePromptResponse creates a standardized prompt response structure
func generatePromptResponse(description, promptText string) dto.PromptGetResponse {
	return dto.PromptGetResponse{
		Description: description,
		Messages: []dto.PromptMessage{
			{
				Role: "user",
				Content: []dto.ContentBlock{
					{
						Type: "text",
						Text: promptText,
					},
				},
			},
		},
	}
}

// generateDirectorFilmographyPrompt generates a director filmography prompt
func (h *PromptHandlers) generateDirectorFilmographyPrompt(args map[string]interface{}) dto.PromptGetResponse {
	directorName, ok := args["director_name"].(string)
	if !ok {
		directorName = ""
	}
	focusPeriod, ok := args["focus_period"].(string)
	if !ok {
		focusPeriod = ""
	}

	var promptText strings.Builder
	promptText.WriteString(fmt.Sprintf("Explore the filmography of director %s", directorName))

	if focusPeriod != "" {
		promptText.WriteString(fmt.Sprintf(" during %s", focusPeriod))
	}

	promptText.WriteString(". Search for all movies by this director, analyze their stylistic evolution, recurring themes, and significant contributions to cinema.")

	return generatePromptResponse("Director filmography exploration prompt", promptText.String())
}

// generateGenreExplorationPrompt generates a genre exploration prompt
func (h *PromptHandlers) generateGenreExplorationPrompt(args map[string]interface{}) dto.PromptGetResponse {
	genre, ok := args["genre"].(string)
	if !ok {
		genre = ""
	}
	subGenre, ok := args["sub_genre"].(string)
	if !ok {
		subGenre = ""
	}

	var promptText strings.Builder
	promptText.WriteString(fmt.Sprintf("Provide a comprehensive exploration of the %s genre", genre))

	if subGenre != "" {
		promptText.WriteString(fmt.Sprintf(", specifically the %s sub-genre", subGenre))
	}

	promptText.WriteString(". Search for representative movies in this genre, discuss its evolution, key characteristics, influential films, and notable directors who shaped it.")

	return generatePromptResponse("Genre exploration prompt", promptText.String())
}

// generateMovieComparisonPrompt generates a movie comparison prompt
func (h *PromptHandlers) generateMovieComparisonPrompt(args map[string]interface{}) dto.PromptGetResponse {
	movie1, ok := args["movie1"].(string)
	if !ok {
		movie1 = ""
	}
	movie2, ok := args["movie2"].(string)
	if !ok {
		movie2 = ""
	}
	aspects, ok := args["comparison_aspects"].(string)
	if !ok {
		aspects = ""
	}

	var promptText strings.Builder
	promptText.WriteString(fmt.Sprintf("Compare and contrast the movies '%s' and '%s'", movie1, movie2))

	if aspects != "" {
		promptText.WriteString(fmt.Sprintf(" focusing on %s", aspects))
	} else {
		promptText.WriteString(" examining themes, directorial style, performances, and cultural impact")
	}

	promptText.WriteString(". First search for both movies to get accurate information, then provide a detailed comparative analysis.")

	return dto.PromptGetResponse{
		Description: "Movie comparison prompt",
		Messages: []dto.PromptMessage{
			{
				Role: "user",
				Content: []dto.ContentBlock{
					{
						Type: "text",
						Text: promptText.String(),
					},
				},
			},
		},
	}
}
