package schemas

import "github.com/francknouama/movies-mcp-server/internal/interfaces/dto"

// Helper functions for schema creation

// Float64Ptr returns a pointer to a float64 value
func Float64Ptr(v float64) *float64 {
	return &v
}

// IntPtr returns a pointer to an int value
func IntPtr(v int) *int {
	return &v
}

// StringArrayItems returns a SchemaProperty for string array items
func StringArrayItems() *dto.SchemaProperty {
	return &dto.SchemaProperty{
		Type: "string",
	}
}

// IntegerArrayItems returns a SchemaProperty for integer array items
func IntegerArrayItems() *dto.SchemaProperty {
	return &dto.SchemaProperty{
		Type: "integer",
	}
}

// NumberArrayItems returns a SchemaProperty for number array items
func NumberArrayItems() *dto.SchemaProperty {
	return &dto.SchemaProperty{
		Type: "number",
	}
}
