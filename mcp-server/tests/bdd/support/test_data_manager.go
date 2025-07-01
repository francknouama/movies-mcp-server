package support

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TestDataManager handles ID interpolation and test data management
type TestDataManager struct {
	storedIDs map[string]int
	lastMovieID int
	lastActorID int
}

// NewTestDataManager creates a new test data manager
func NewTestDataManager() *TestDataManager {
	return &TestDataManager{
		storedIDs: make(map[string]int),
	}
}

// StoreID stores an ID with a given key for later reference
func (tdm *TestDataManager) StoreID(key string, id int) {
	tdm.storedIDs[key] = id
	
	// Track last IDs for automatic reference
	if strings.Contains(key, "movie") {
		tdm.lastMovieID = id
	}
	if strings.Contains(key, "actor") {
		tdm.lastActorID = id
	}
}

// GetID retrieves a stored ID by key
func (tdm *TestDataManager) GetID(key string) (int, bool) {
	id, exists := tdm.storedIDs[key]
	return id, exists
}

// GetLastMovieID returns the last created movie ID
func (tdm *TestDataManager) GetLastMovieID() int {
	return tdm.lastMovieID
}

// GetLastActorID returns the last created actor ID
func (tdm *TestDataManager) GetLastActorID() int {
	return tdm.lastActorID
}

// InterpolateString replaces placeholders like {movie_id}, {actor_id} with actual values
func (tdm *TestDataManager) InterpolateString(input string) string {
	// Pattern to match {key} placeholders
	re := regexp.MustCompile(`\{([^}]+)\}`)
	
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract key from {key}
		key := match[1 : len(match)-1]
		
		// Handle special cases
		switch key {
		case "movie_id":
			if tdm.lastMovieID > 0 {
				return strconv.Itoa(tdm.lastMovieID)
			}
		case "actor_id":
			if tdm.lastActorID > 0 {
				return strconv.Itoa(tdm.lastActorID)
			}
		default:
			// Look up stored ID
			if id, exists := tdm.storedIDs[key]; exists {
				return strconv.Itoa(id)
			}
		}
		
		// Return original if no replacement found
		return match
	})
}

// InterpolateMap replaces placeholders in map values
func (tdm *TestDataManager) InterpolateMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range input {
		switch v := value.(type) {
		case string:
			result[key] = tdm.InterpolateString(v)
		case int:
			// Check if this looks like a placeholder value (negative numbers often used as placeholders)
			if v < 0 {
				switch key {
				case "movie_id":
					if tdm.lastMovieID > 0 {
						result[key] = tdm.lastMovieID
					} else {
						result[key] = v
					}
				case "actor_id":
					if tdm.lastActorID > 0 {
						result[key] = tdm.lastActorID
					} else {
						result[key] = v
					}
				default:
					result[key] = v
				}
			} else {
				result[key] = v
			}
		default:
			result[key] = value
		}
	}
	
	return result
}

// Clear resets all stored IDs
func (tdm *TestDataManager) Clear() {
	tdm.storedIDs = make(map[string]int)
	tdm.lastMovieID = 0
	tdm.lastActorID = 0
}

// ParseIDFromResponse extracts an ID from a response structure
func (tdm *TestDataManager) ParseIDFromResponse(response interface{}, idField string) (int, error) {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("response is not a map")
	}
	
	idValue, exists := responseMap[idField]
	if !exists {
		return 0, fmt.Errorf("field %s not found in response", idField)
	}
	
	// Handle different number types
	switch v := idValue.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("field %s is not a number: %T", idField, v)
	}
}

// StoreIDFromResponse extracts and stores an ID from response
func (tdm *TestDataManager) StoreIDFromResponse(response interface{}, idField, key string) error {
	id, err := tdm.ParseIDFromResponse(response, idField)
	if err != nil {
		return err
	}
	
	tdm.StoreID(key, id)
	return nil
}