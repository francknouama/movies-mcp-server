package mcp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

// ToolValidator provides enhanced validation for MCP tool calls
type ToolValidator struct {
	toolSchemas map[string]dto.Tool
}

// ValidationError represents a detailed validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationResult contains validation results and detailed errors
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// NewToolValidator creates a new tool validator
func NewToolValidator(toolSchemas []dto.Tool) *ToolValidator {
	schemaMap := make(map[string]dto.Tool)
	for _, tool := range toolSchemas {
		schemaMap[tool.Name] = tool
	}
	
	return &ToolValidator{
		toolSchemas: schemaMap,
	}
}

// GetSchemas returns all tool schemas
func (tv *ToolValidator) GetSchemas() []dto.Tool {
	schemas := make([]dto.Tool, 0, len(tv.toolSchemas))
	for _, schema := range tv.toolSchemas {
		schemas = append(schemas, schema)
	}
	return schemas
}

// ValidateToolCall validates a tool call against its schema
func (tv *ToolValidator) ValidateToolCall(toolName string, arguments map[string]interface{}) ValidationResult {
	schema, exists := tv.toolSchemas[toolName]
	if !exists {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Field:   "tool_name",
				Value:   toolName,
				Message: fmt.Sprintf("Unknown tool: %s", toolName),
				Code:    "UNKNOWN_TOOL",
			}},
		}
	}

	return tv.validateArguments(arguments, schema.InputSchema)
}

// validateArguments validates arguments against a schema
func (tv *ToolValidator) validateArguments(arguments map[string]interface{}, schema dto.InputSchema) ValidationResult {
	var errors []ValidationError

	// Check required fields
	for _, required := range schema.Required {
		if _, exists := arguments[required]; !exists {
			errors = append(errors, ValidationError{
				Field:   required,
				Value:   "",
				Message: fmt.Sprintf("Required field '%s' is missing", required),
				Code:    "REQUIRED_FIELD_MISSING",
			})
		}
	}

	// Validate each provided argument
	for fieldName, value := range arguments {
		fieldSchema, exists := schema.Properties[fieldName]
		if !exists {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("%v", value),
				Message: fmt.Sprintf("Unknown field '%s'", fieldName),
				Code:    "UNKNOWN_FIELD",
			})
			continue
		}

		fieldErrors := tv.validateField(fieldName, value, fieldSchema)
		errors = append(errors, fieldErrors...)
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

// validateField validates a single field against its schema
func (tv *ToolValidator) validateField(fieldName string, value interface{}, fieldSchema interface{}) []ValidationError {
	var errors []ValidationError

	schema, ok := fieldSchema.(map[string]interface{})
	if !ok {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: "Invalid field schema",
			Code:    "INVALID_SCHEMA",
		}}
	}

	fieldType, exists := schema["type"].(string)
	if !exists {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: "Field schema missing type",
			Code:    "MISSING_TYPE",
		}}
	}

	// Type validation
	switch fieldType {
	case "string":
		errors = append(errors, tv.validateString(fieldName, value, schema)...)
	case "integer":
		errors = append(errors, tv.validateInteger(fieldName, value, schema)...)
	case "number":
		errors = append(errors, tv.validateNumber(fieldName, value, schema)...)
	case "boolean":
		errors = append(errors, tv.validateBoolean(fieldName, value, schema)...)
	case "array":
		errors = append(errors, tv.validateArray(fieldName, value, schema)...)
	case "object":
		errors = append(errors, tv.validateObject(fieldName, value, schema)...)
	default:
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Unsupported field type: %s", fieldType),
			Code:    "UNSUPPORTED_TYPE",
		})
	}

	return errors
}

// validateString validates string fields
func (tv *ToolValidator) validateString(fieldName string, value interface{}, schema map[string]interface{}) []ValidationError {
	var errors []ValidationError

	str, ok := value.(string)
	if !ok {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected string, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}

	// Enum validation
	if enum, exists := schema["enum"]; exists {
		if enumSlice, ok := enum.([]string); ok {
			valid := false
			for _, enumValue := range enumSlice {
				if str == enumValue {
					valid = true
					break
				}
			}
			if !valid {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Value:   str,
					Message: fmt.Sprintf("Value must be one of: %s", strings.Join(enumSlice, ", ")),
					Code:    "INVALID_ENUM_VALUE",
				})
			}
		}
	}

	// Length validation
	if minLength, exists := schema["minLength"]; exists {
		if min, ok := minLength.(float64); ok && len(str) < int(min) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   str,
				Message: fmt.Sprintf("String length must be at least %d characters", int(min)),
				Code:    "STRING_TOO_SHORT",
			})
		}
	}

	if maxLength, exists := schema["maxLength"]; exists {
		if max, ok := maxLength.(float64); ok && len(str) > int(max) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   str,
				Message: fmt.Sprintf("String length must be at most %d characters", int(max)),
				Code:    "STRING_TOO_LONG",
			})
		}
	}

	// Pattern validation (basic)
	if pattern, exists := schema["pattern"]; exists {
		if patternStr, ok := pattern.(string); ok {
			// For this implementation, we'll just check some common patterns
			switch patternStr {
			case "^\\d{4}-\\d{2}-\\d{2}$": // Date format YYYY-MM-DD
				if !isValidDateFormat(str) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Value:   str,
						Message: "Invalid date format, expected YYYY-MM-DD",
						Code:    "INVALID_DATE_FORMAT",
					})
				}
			}
		}
	}

	// Format validation
	if format, exists := schema["format"]; exists {
		if formatStr, ok := format.(string); ok {
			switch formatStr {
			case "date":
				if !isValidDateFormat(str) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Value:   str,
						Message: "Invalid date format, expected YYYY-MM-DD",
						Code:    "INVALID_DATE_FORMAT",
					})
				}
			case "uri":
				if !isValidURI(str) {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Value:   str,
						Message: "Invalid URI format",
						Code:    "INVALID_URI_FORMAT",
					})
				}
			}
		}
	}

	return errors
}

// validateInteger validates integer fields
func (tv *ToolValidator) validateInteger(fieldName string, value interface{}, schema map[string]interface{}) []ValidationError {
	var errors []ValidationError

	var intVal int64
	switch v := value.(type) {
	case int:
		intVal = int64(v)
	case int64:
		intVal = v
	case float64:
		if v != float64(int64(v)) {
			return []ValidationError{{
				Field:   fieldName,
				Value:   fmt.Sprintf("%v", value),
				Message: "Expected integer, got decimal number",
				Code:    "NOT_INTEGER",
			}}
		}
		intVal = int64(v)
	default:
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected integer, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}

	// Range validation
	if minimum, exists := schema["minimum"]; exists {
		if min, ok := minimum.(float64); ok && intVal < int64(min) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("%d", intVal),
				Message: fmt.Sprintf("Value must be at least %d", int64(min)),
				Code:    "VALUE_TOO_SMALL",
			})
		}
	}

	if maximum, exists := schema["maximum"]; exists {
		if max, ok := maximum.(float64); ok && intVal > int64(max) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("%d", intVal),
				Message: fmt.Sprintf("Value must be at most %d", int64(max)),
				Code:    "VALUE_TOO_LARGE",
			})
		}
	}

	return errors
}

// validateNumber validates number fields
func (tv *ToolValidator) validateNumber(fieldName string, value interface{}, schema map[string]interface{}) []ValidationError {
	var errors []ValidationError

	var numVal float64
	switch v := value.(type) {
	case int:
		numVal = float64(v)
	case int64:
		numVal = float64(v)
	case float64:
		numVal = v
	default:
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected number, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}

	// Range validation
	if minimum, exists := schema["minimum"]; exists {
		if min, ok := minimum.(float64); ok && numVal < min {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("%g", numVal),
				Message: fmt.Sprintf("Value must be at least %g", min),
				Code:    "VALUE_TOO_SMALL",
			})
		}
	}

	if maximum, exists := schema["maximum"]; exists {
		if max, ok := maximum.(float64); ok && numVal > max {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("%g", numVal),
				Message: fmt.Sprintf("Value must be at most %g", max),
				Code:    "VALUE_TOO_LARGE",
			})
		}
	}

	return errors
}

// validateBoolean validates boolean fields
func (tv *ToolValidator) validateBoolean(fieldName string, value interface{}, _ map[string]interface{}) []ValidationError {
	if _, ok := value.(bool); !ok {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected boolean, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}
	return nil
}

// validateArray validates array fields
func (tv *ToolValidator) validateArray(fieldName string, value interface{}, schema map[string]interface{}) []ValidationError {
	var errors []ValidationError

	arr, ok := value.([]interface{})
	if !ok {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected array, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}

	// Length validation
	if minItems, exists := schema["minItems"]; exists {
		if min, ok := minItems.(float64); ok && len(arr) < int(min) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("array with %d items", len(arr)),
				Message: fmt.Sprintf("Array must have at least %d items", int(min)),
				Code:    "ARRAY_TOO_SHORT",
			})
		}
	}

	if maxItems, exists := schema["maxItems"]; exists {
		if max, ok := maxItems.(float64); ok && len(arr) > int(max) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   fmt.Sprintf("array with %d items", len(arr)),
				Message: fmt.Sprintf("Array must have at most %d items", int(max)),
				Code:    "ARRAY_TOO_LONG",
			})
		}
	}

	// Item validation
	if items, exists := schema["items"]; exists {
		if itemSchema, ok := items.(map[string]interface{}); ok {
			for i, item := range arr {
				itemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
				itemErrors := tv.validateField(itemFieldName, item, itemSchema)
				errors = append(errors, itemErrors...)
			}
		}
	}

	return errors
}

// validateObject validates object fields
func (tv *ToolValidator) validateObject(fieldName string, value interface{}, schema map[string]interface{}) []ValidationError {
	var errors []ValidationError

	obj, ok := value.(map[string]interface{})
	if !ok {
		return []ValidationError{{
			Field:   fieldName,
			Value:   fmt.Sprintf("%v", value),
			Message: fmt.Sprintf("Expected object, got %T", value),
			Code:    "TYPE_MISMATCH",
		}}
	}

	// Validate nested properties if schema is provided
	if properties, exists := schema["properties"]; exists {
		if propSchema, ok := properties.(map[string]interface{}); ok {
			for propName, propValue := range obj {
				if propFieldSchema, exists := propSchema[propName]; exists {
					nestedFieldName := fmt.Sprintf("%s.%s", fieldName, propName)
					propErrors := tv.validateField(nestedFieldName, propValue, propFieldSchema)
					errors = append(errors, propErrors...)
				}
			}
		}
	}

	// Check required properties
	if required, exists := schema["required"]; exists {
		if requiredFields, ok := required.([]string); ok {
			for _, req := range requiredFields {
				if _, exists := obj[req]; !exists {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("%s.%s", fieldName, req),
						Value:   "",
						Message: fmt.Sprintf("Required property '%s' is missing", req),
						Code:    "REQUIRED_PROPERTY_MISSING",
					})
				}
			}
		}
	}

	return errors
}

// Helper functions

func isValidDateFormat(dateStr string) bool {
	if len(dateStr) != 10 {
		return false
	}
	
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return false
	}
	
	year, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	day, err3 := strconv.Atoi(parts[2])
	
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	
	return year >= 1000 && year <= 9999 && 
		   month >= 1 && month <= 12 && 
		   day >= 1 && day <= 31
}

func isValidURI(uri string) bool {
	// Basic URI validation - just check it's not empty and has a reasonable format
	if len(uri) == 0 {
		return false
	}
	
	// Check for common URI patterns
	return strings.Contains(uri, "://") || strings.HasPrefix(uri, "/") || strings.HasPrefix(uri, "mailto:")
}

// HandleValidateToolCall handles MCP tool call for validating other tool calls
func (tv *ToolValidator) HandleValidateToolCall(
	id interface{},
	arguments map[string]interface{},
	sendResult func(interface{}, interface{}),
	sendError func(interface{}, int, string, interface{}),
) {
	toolName, ok := arguments["tool_name"].(string)
	if !ok || toolName == "" {
		sendError(id, dto.InvalidParams, "Tool name is required", nil)
		return
	}

	toolArgs, ok := arguments["tool_arguments"].(map[string]interface{})
	if !ok {
		sendError(id, dto.InvalidParams, "Tool arguments must be an object", nil)
		return
	}

	result := tv.ValidateToolCall(toolName, toolArgs)
	sendResult(id, result)
}