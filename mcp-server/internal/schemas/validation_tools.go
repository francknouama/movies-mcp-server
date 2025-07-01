package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

// GetValidationTools returns all validation tool schemas
func GetValidationTools() []dto.Tool {
	return []dto.Tool{
		validateToolCallTool(),
	}
}

func validateToolCallTool() dto.Tool {
	return dto.Tool{
		Name:        "validate_tool_call",
		Description: "Validate a tool call against its schema",
		InputSchema: dto.InputSchema{
			Type: "object",
			Properties: map[string]dto.SchemaProperty{
				"tool_name": {
					Type:        "string",
					Description: "Tool name to validate",
				},
				"arguments": {
					Type:        "object",
					Description: "Arguments to validate",
				},
			},
			Required: []string{"tool_name", "arguments"},
		},
	}
}
