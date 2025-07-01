package support

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strconv"
	"strings"
	"time"
)

// TestUtilities provides helper functions for BDD tests
type TestUtilities struct {
	random *mathrand.Rand
}

// NewTestUtilities creates a new test utilities instance
func NewTestUtilities() *TestUtilities {
	return &TestUtilities{
		random: mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateRandomString generates a random string of specified length
func (tu *TestUtilities) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[tu.random.Intn(len(charset))]
	}
	return string(b)
}

// GenerateSecureRandomString generates a cryptographically secure random string
func (tu *TestUtilities) GenerateSecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GenerateRandomMovie generates a random movie for testing
func (tu *TestUtilities) GenerateRandomMovie() map[string]interface{} {
	directors := []string{
		"Steven Spielberg", "Martin Scorsese", "Christopher Nolan", "Quentin Tarantino",
		"Ridley Scott", "David Fincher", "James Cameron", "Peter Jackson",
	}

	genres := []string{"Action", "Drama", "Comedy", "Sci-Fi", "Horror", "Romance", "Thriller", "Adventure"}

	return map[string]interface{}{
		"title":    "Test Movie " + tu.GenerateRandomString(8),
		"director": directors[tu.random.Intn(len(directors))],
		"year":     2000 + tu.random.Intn(24),     // 2000-2023
		"rating":   5.0 + tu.random.Float64()*5.0, // 5.0-10.0
		"genre":    genres[tu.random.Intn(len(genres))],
	}
}

// GenerateRandomActor generates a random actor for testing
func (tu *TestUtilities) GenerateRandomActor() map[string]interface{} {
	firstNames := []string{
		"John", "Jane", "Michael", "Sarah", "David", "Emily", "Chris", "Emma",
		"Robert", "Lisa", "James", "Anna", "William", "Maria", "Richard", "Jessica",
	}

	lastNames := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
		"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	}

	firstName := firstNames[tu.random.Intn(len(firstNames))]
	lastName := lastNames[tu.random.Intn(len(lastNames))]

	return map[string]interface{}{
		"name":       firstName + " " + lastName,
		"birth_year": 1940 + tu.random.Intn(60), // 1940-1999
		"bio":        "Test actor biography for " + firstName + " " + lastName,
	}
}

// ValidateMovieResponse validates a movie response structure
func (tu *TestUtilities) ValidateMovieResponse(response map[string]interface{}) error {
	requiredFields := []string{"id", "title", "director", "year", "rating"}

	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("movie response missing required field: %s", field)
		}
	}

	// Validate field types
	if id, ok := response["id"]; ok {
		switch v := id.(type) {
		case float64:
			if v <= 0 {
				return fmt.Errorf("movie ID should be positive, got %v", v)
			}
		case int:
			if v <= 0 {
				return fmt.Errorf("movie ID should be positive, got %v", v)
			}
		default:
			return fmt.Errorf("movie ID should be a number, got %T", v)
		}
	}

	if title, ok := response["title"]; ok {
		if titleStr, ok := title.(string); !ok || titleStr == "" {
			return fmt.Errorf("movie title should be a non-empty string, got %T: %v", title, title)
		}
	}

	if year, ok := response["year"]; ok {
		switch v := year.(type) {
		case float64:
			if v < 1800 || v > 2100 {
				return fmt.Errorf("movie year should be reasonable (1800-2100), got %v", v)
			}
		case int:
			if v < 1800 || v > 2100 {
				return fmt.Errorf("movie year should be reasonable (1800-2100), got %v", v)
			}
		default:
			return fmt.Errorf("movie year should be a number, got %T", v)
		}
	}

	if rating, ok := response["rating"]; ok {
		switch v := rating.(type) {
		case float64:
			if v < 0 || v > 10 {
				return fmt.Errorf("movie rating should be between 0-10, got %v", v)
			}
		case int:
			if v < 0 || v > 10 {
				return fmt.Errorf("movie rating should be between 0-10, got %v", v)
			}
		default:
			return fmt.Errorf("movie rating should be a number, got %T", v)
		}
	}

	return nil
}

// ValidateActorResponse validates an actor response structure
func (tu *TestUtilities) ValidateActorResponse(response map[string]interface{}) error {
	requiredFields := []string{"id", "name", "birth_year"}

	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("actor response missing required field: %s", field)
		}
	}

	// Validate field types
	if id, ok := response["id"]; ok {
		switch v := id.(type) {
		case float64:
			if v <= 0 {
				return fmt.Errorf("actor ID should be positive, got %v", v)
			}
		case int:
			if v <= 0 {
				return fmt.Errorf("actor ID should be positive, got %v", v)
			}
		default:
			return fmt.Errorf("actor ID should be a number, got %T", v)
		}
	}

	if name, ok := response["name"]; ok {
		if nameStr, ok := name.(string); !ok || nameStr == "" {
			return fmt.Errorf("actor name should be a non-empty string, got %T: %v", name, name)
		}
	}

	if birthYear, ok := response["birth_year"]; ok {
		switch v := birthYear.(type) {
		case float64:
			if v < 1800 || v > 2020 {
				return fmt.Errorf("actor birth year should be reasonable (1800-2020), got %v", v)
			}
		case int:
			if v < 1800 || v > 2020 {
				return fmt.Errorf("actor birth year should be reasonable (1800-2020), got %v", v)
			}
		default:
			return fmt.Errorf("actor birth year should be a number, got %T", v)
		}
	}

	return nil
}

// ParseJSONToMap parses JSON string to map[string]interface{}
func (tu *TestUtilities) ParseJSONToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// ExtractIDFromResponse extracts ID from response and converts to int
func (tu *TestUtilities) ExtractIDFromResponse(response map[string]interface{}, idField string) (int, error) {
	idValue, exists := response[idField]
	if !exists {
		return 0, fmt.Errorf("field %s not found in response", idField)
	}

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

// WaitForCondition waits for a condition to be true with timeout
func (tu *TestUtilities) WaitForCondition(condition func() bool, timeout time.Duration, checkInterval time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return nil
		}
		time.Sleep(checkInterval)
	}
	return fmt.Errorf("condition not met within timeout of %v", timeout)
}

// RetryOperation retries an operation with exponential backoff
func (tu *TestUtilities) RetryOperation(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries-1 {
			// Limit the shift to prevent overflow
			shift := i
			if shift > 10 { // Cap at 2^10 = 1024x
				shift = 10
			}
			delay := baseDelay * time.Duration(1<<uint(shift)) // Exponential backoff
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("operation failed after %d retries, last error: %w", maxRetries, lastErr)
}

// NormalizeString normalizes a string for comparison (lowercase, trimmed)
func (tu *TestUtilities) NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ContainsSubstring checks if a string contains a substring (case-insensitive)
func (tu *TestUtilities) ContainsSubstring(s, substr string) bool {
	return strings.Contains(tu.NormalizeString(s), tu.NormalizeString(substr))
}

// ValidateArrayField validates that a field is an array with expected properties
func (tu *TestUtilities) ValidateArrayField(response map[string]interface{}, fieldName string, minLength, maxLength int) error {
	field, exists := response[fieldName]
	if !exists {
		return fmt.Errorf("field %s not found in response", fieldName)
	}

	array, ok := field.([]interface{})
	if !ok {
		return fmt.Errorf("field %s should be an array, got %T", fieldName, field)
	}

	length := len(array)
	if length < minLength {
		return fmt.Errorf("field %s should have at least %d items, got %d", fieldName, minLength, length)
	}

	if maxLength >= 0 && length > maxLength {
		return fmt.Errorf("field %s should have at most %d items, got %d", fieldName, maxLength, length)
	}

	return nil
}

// ValidateNumberField validates that a field is a number within expected range
func (tu *TestUtilities) ValidateNumberField(response map[string]interface{}, fieldName string, min, max float64) error {
	field, exists := response[fieldName]
	if !exists {
		return fmt.Errorf("field %s not found in response", fieldName)
	}

	var value float64
	switch v := field.(type) {
	case int:
		value = float64(v)
	case float64:
		value = v
	default:
		return fmt.Errorf("field %s should be a number, got %T", fieldName, field)
	}

	if value < min {
		return fmt.Errorf("field %s should be at least %v, got %v", fieldName, min, value)
	}

	if value > max {
		return fmt.Errorf("field %s should be at most %v, got %v", fieldName, max, value)
	}

	return nil
}

// ValidateStringField validates that a field is a non-empty string
func (tu *TestUtilities) ValidateStringField(response map[string]interface{}, fieldName string, required bool) error {
	field, exists := response[fieldName]
	if !exists {
		if required {
			return fmt.Errorf("required field %s not found in response", fieldName)
		}
		return nil
	}

	str, ok := field.(string)
	if !ok {
		return fmt.Errorf("field %s should be a string, got %T", fieldName, field)
	}

	if required && strings.TrimSpace(str) == "" {
		return fmt.Errorf("required field %s should not be empty", fieldName)
	}

	return nil
}

// CreateTestMovieBatch creates a batch of test movies for performance testing
func (tu *TestUtilities) CreateTestMovieBatch(count int) []map[string]interface{} {
	movies := make([]map[string]interface{}, count)

	for i := 0; i < count; i++ {
		movie := tu.GenerateRandomMovie()
		// Ensure unique titles
		movie["title"] = fmt.Sprintf("Batch Movie %d - %s", i+1, tu.GenerateRandomString(6))
		movies[i] = movie
	}

	return movies
}

// CreateTestActorBatch creates a batch of test actors for performance testing
func (tu *TestUtilities) CreateTestActorBatch(count int) []map[string]interface{} {
	actors := make([]map[string]interface{}, count)

	for i := 0; i < count; i++ {
		actor := tu.GenerateRandomActor()
		// Ensure unique names
		actor["name"] = fmt.Sprintf("Batch Actor %d %s", i+1, tu.GenerateRandomString(6))
		actors[i] = actor
	}

	return actors
}

// MeasureExecutionTime measures the execution time of a function
func (tu *TestUtilities) MeasureExecutionTime(operation func() error) (time.Duration, error) {
	start := time.Now()
	err := operation()
	duration := time.Since(start)
	return duration, err
}

// ValidatePerformance validates that an operation completes within expected time
func (tu *TestUtilities) ValidatePerformance(operation func() error, maxDuration time.Duration) error {
	duration, err := tu.MeasureExecutionTime(operation)
	if err != nil {
		return fmt.Errorf("operation failed: %w", err)
	}

	if duration > maxDuration {
		return fmt.Errorf("operation took %v, expected under %v", duration, maxDuration)
	}

	return nil
}

// DeepCopyMap creates a deep copy of a map[string]interface{}
func (tu *TestUtilities) DeepCopyMap(original map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})

	for key, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[key] = tu.DeepCopyMap(v)
		case []interface{}:
			copySlice := make([]interface{}, len(v))
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					copySlice[i] = tu.DeepCopyMap(itemMap)
				} else {
					copySlice[i] = item
				}
			}
			copy[key] = copySlice
		default:
			copy[key] = value
		}
	}

	return copy
}
