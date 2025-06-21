package schemas

import "movies-mcp-server/internal/interfaces/dto"

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