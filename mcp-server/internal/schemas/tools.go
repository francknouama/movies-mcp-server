package schemas

import "github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"

// GetToolSchemas returns all tool schemas for the MCP server
func GetToolSchemas() []dto.Tool {
	var tools []dto.Tool
	
	// Movie Management Tools
	tools = append(tools, GetMovieTools()...)
	
	// Actor Management Tools
	tools = append(tools, GetActorTools()...)
	
	// Search Tools
	tools = append(tools, GetSearchTools()...)
	
	// Compound Tools
	tools = append(tools, GetCompoundTools()...)
	
	// Context Management Tools
	tools = append(tools, GetContextTools()...)
	
	// Validation Tools
	tools = append(tools, GetValidationTools()...)
	
	return tools
}