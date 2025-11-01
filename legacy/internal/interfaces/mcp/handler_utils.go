package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/francknouama/movies-mcp-server/internal/interfaces/dto"
)

// parseIDArgument extracts and validates an ID argument from the arguments map
func parseIDArgument(arguments map[string]any, paramName string) (int, error) {
	idFloat, ok := arguments[paramName].(float64)
	if !ok {
		return 0, fmt.Errorf("%s is required and must be a number", paramName)
	}
	return int(idFloat), nil
}

// handleDeleteOperation performs a generic delete operation with standardized error handling
func handleDeleteOperation(
	id any,
	arguments map[string]any,
	paramName string,
	entityName string,
	deleteFunc func(context.Context, int) error,
	sendResult func(any, any),
	sendError func(any, int, string, any),
) {
	// Parse ID
	entityID, err := parseIDArgument(arguments, paramName)
	if err != nil {
		sendError(id, dto.InvalidParams, err.Error(), nil)
		return
	}

	// Delete entity
	ctx := context.Background()
	err = deleteFunc(ctx, entityID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, fmt.Sprintf("%s not found", entityName), nil)
		} else {
			sendError(id, dto.InternalError, fmt.Sprintf("Failed to delete %s", strings.ToLower(entityName)), err.Error())
		}
		return
	}

	sendResult(id, map[string]string{"message": fmt.Sprintf("%s deleted successfully", entityName)})
}

// handleGetOperation performs a generic get operation with standardized error handling
func handleGetOperation[T any, R any](
	id any,
	arguments map[string]any,
	paramName string,
	entityName string,
	getFunc func(context.Context, int) (T, error),
	responseFunc func(T) R,
	sendResult func(any, any),
	sendError func(any, int, string, any),
) {
	// Parse ID
	entityID, err := parseIDArgument(arguments, paramName)
	if err != nil {
		sendError(id, dto.InvalidParams, err.Error(), nil)
		return
	}

	// Get entity from service
	ctx := context.Background()
	entity, err := getFunc(ctx, entityID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendError(id, dto.InvalidParams, fmt.Sprintf("%s not found", entityName), nil)
		} else {
			sendError(id, dto.InternalError, fmt.Sprintf("Failed to get %s", strings.ToLower(entityName)), err.Error())
		}
		return
	}

	// Convert to response format
	response := responseFunc(entity)
	sendResult(id, response)
}