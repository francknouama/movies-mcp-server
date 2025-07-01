package server

import (
	"encoding/json"
	"fmt"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

// ResourceManager handles resource operations
type ResourceManager struct {
	registry *Registry
}

// NewResourceManager creates a new resource manager
func NewResourceManager(registry *Registry) *ResourceManager {
	return &ResourceManager{
		registry: registry,
	}
}

// RegisterDefaultResources registers the default movie database resources
func (rm *ResourceManager) RegisterDefaultResources() {
	// Register database resources
	rm.registry.RegisterResource(
		"movies://database/all",
		rm.handleAllMovies,
		dto.Resource{
			URI:         "movies://database/all",
			Name:        "All Movies",
			Description: "Complete movie database in JSON format",
			MimeType:    "application/json",
		},
	)

	rm.registry.RegisterResource(
		"movies://database/stats",
		rm.handleDatabaseStats,
		dto.Resource{
			URI:         "movies://database/stats",
			Name:        "Database Statistics",
			Description: "Movie database statistics and analytics",
			MimeType:    "application/json",
		},
	)

	// Register poster collection resource
	rm.registry.RegisterResource(
		"movies://posters/collection",
		rm.handlePosterCollection,
		dto.Resource{
			URI:         "movies://posters/collection",
			Name:        "Movie Posters Collection",
			Description: "Collection of all movie posters",
			MimeType:    "application/json",
		},
	)
}

// handleAllMovies handles the movies://database/all resource
func (rm *ResourceManager) handleAllMovies(uri string, sender ResponseSender) {
	// This would typically fetch all movies from the database
	// For now, return a placeholder response
	content := dto.ResourceContent{
		URI:      uri,
		MimeType: "application/json",
		Text:     `{"message": "All movies resource - implementation pending"}`,
	}

	response := dto.ResourceReadResponse{
		Contents: []dto.ResourceContent{content},
	}

	// We need to extract the ID from the original request
	// This is a limitation of the current design that we'll need to address
	sender.SendResult(nil, response)
}

// handleDatabaseStats handles the movies://database/stats resource
func (rm *ResourceManager) handleDatabaseStats(uri string, sender ResponseSender) {
	stats := map[string]interface{}{
		"total_movies": 0,
		"total_actors": 0,
		"genres":       []string{},
		"year_range": map[string]interface{}{
			"earliest": nil,
			"latest":   nil,
		},
	}

	statsJSON, _ := json.Marshal(stats)
	content := dto.ResourceContent{
		URI:      uri,
		MimeType: "application/json",
		Text:     string(statsJSON),
	}

	response := dto.ResourceReadResponse{
		Contents: []dto.ResourceContent{content},
	}

	sender.SendResult(nil, response)
}

// handlePosterCollection handles the movies://posters/collection resource
func (rm *ResourceManager) handlePosterCollection(uri string, sender ResponseSender) {
	collection := map[string]interface{}{
		"posters": []map[string]interface{}{},
		"total":   0,
	}

	collectionJSON, _ := json.Marshal(collection)
	content := dto.ResourceContent{
		URI:      uri,
		MimeType: "application/json",
		Text:     string(collectionJSON),
	}

	response := dto.ResourceReadResponse{
		Contents: []dto.ResourceContent{content},
	}

	sender.SendResult(nil, response)
}

// RegisterPosterResource registers a dynamic poster resource
func (rm *ResourceManager) RegisterPosterResource(movieID int, movieTitle string) {
	uri := fmt.Sprintf("movies://posters/%d", movieID)

	rm.registry.RegisterResource(
		uri,
		func(uri string, sender ResponseSender) {
			rm.handleMoviePoster(uri, movieID, sender)
		},
		dto.Resource{
			URI:         uri,
			Name:        fmt.Sprintf("%s Poster", movieTitle),
			Description: fmt.Sprintf("Movie poster for %s", movieTitle),
			MimeType:    "image/jpeg",
		},
	)
}

// handleMoviePoster handles individual movie poster resources
func (rm *ResourceManager) handleMoviePoster(uri string, movieID int, sender ResponseSender) {
	// This would typically fetch the poster from the database
	// For now, return a placeholder
	content := dto.ResourceContent{
		URI:      uri,
		MimeType: "image/jpeg",
		Blob:     nil, // Base64 encoded image data would go here
	}

	response := dto.ResourceReadResponse{
		Contents: []dto.ResourceContent{content},
	}

	sender.SendResult(nil, response)
}
