package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

// GetContextTools returns all context management tool schemas
func GetContextTools() []dto.Tool {
	return []dto.Tool{
		createSearchContextTool(),
		getContextPageTool(),
		getContextInfoTool(),
	}
}

func createSearchContextTool() dto.Tool {
	return dto.Tool{
		Name:        "create_search_context",
		Description: "Create a search context for large result sets",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]dto.SchemaProperty{
				"search_criteria": {
					Type:        "object",
					Description: "Search criteria",
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
			Properties: map[string]dto.SchemaProperty{
				"context_id": {
					Type:        "string",
					Description: "Context ID",
				},
				"page": {
					Type:        "integer",
					Description: "Page number",
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
			Properties: map[string]dto.SchemaProperty{
				"context_id": {
					Type:        "string",
					Description: "Context ID",
				},
			},
			Required: []string{"context_id"},
		},
	}
}
