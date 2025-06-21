package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"movies-mcp-server/pkg/errors"
)

// ValidationRule represents a validation rule
type ValidationRule func(value interface{}) error

// Validator provides request validation functionality
type Validator struct {
	rules map[string][]ValidationRule
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]ValidationRule),
	}
}

// AddRule adds a validation rule for a field
func (v *Validator) AddRule(field string, rule ValidationRule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate validates a map of values against registered rules
func (v *Validator) Validate(values map[string]interface{}) error {
	var validationErrors []string

	for field, rules := range v.rules {
		value, exists := values[field]
		
		// Check if required field is missing
		if !exists {
			for _, rule := range rules {
				if err := rule(nil); err != nil {
					if strings.Contains(err.Error(), "required") {
						validationErrors = append(validationErrors, fmt.Sprintf("%s: %s", field, err.Error()))
						break
					}
				}
			}
			continue
		}

		// Apply validation rules
		for _, rule := range rules {
			if err := rule(value); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("%s: %s", field, err.Error()))
			}
		}
	}

	if len(validationErrors) > 0 {
		// Create a more descriptive error message that includes field names
		mainMessage := fmt.Sprintf("Validation failed for fields: %s", strings.Join(getFieldNames(validationErrors), ", "))
		return errors.NewValidationError(
			mainMessage,
			"multiple_fields",
			map[string]interface{}{
				"validation_errors": validationErrors,
				"error_count":      len(validationErrors),
			},
		)
	}

	return nil
}

// ValidateStruct validates a struct using struct tags
func (v *Validator) ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return errors.NewValidationError("Value must be a struct", "type", val.Kind())
	}

	typ := val.Type()
	var validationErrors []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		
		// Check validation tags
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if err := v.validateFieldByTag(field.Name, fieldValue.Interface(), rule); err != nil {
				validationErrors = append(validationErrors, err.Error())
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError(
			"Struct validation failed",
			"struct",
			map[string]interface{}{
				"validation_errors": validationErrors,
			},
		)
	}

	return nil
}

// validateFieldByTag validates a field based on a tag rule
func (v *Validator) validateFieldByTag(fieldName string, value interface{}, rule string) error {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		return Required()(value)
	case "min":
		if minVal, err := strconv.Atoi(ruleValue); err == nil {
			// Check if it's a string/slice (use MinLength) or number (use Min)
			switch value.(type) {
			case string, []interface{}:
				return MinLength(minVal)(value)
			default:
				return Min(float64(minVal))(value)
			}
		}
	case "max":
		if maxVal, err := strconv.Atoi(ruleValue); err == nil {
			// Check if it's a string/slice (use MaxLength) or number (use Max)
			switch value.(type) {
			case string, []interface{}:
				return MaxLength(maxVal)(value)
			default:
				return Max(float64(maxVal))(value)
			}
		}
	case "email":
		return Email()(value)
	case "url":
		return URL()(value)
	case "alpha":
		return Alpha()(value)
	case "numeric":
		return Numeric()(value)
	case "alphanumeric":
		return AlphaNumeric()(value)
	}

	return nil
}

// Common validation rules

// Required validates that a field is not empty
func Required() ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return fmt.Errorf("field is required")
		}

		switch v := value.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				return fmt.Errorf("field is required")
			}
		case []interface{}:
			if len(v) == 0 {
				return fmt.Errorf("field is required")
			}
		case map[string]interface{}:
			// For maps, we just check if they exist, not if they're empty
			// Empty maps are valid for capabilities and clientInfo
		}

		return nil
	}
}

// MinLength validates minimum length for strings and slices
func MinLength(min int) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		var length int
		switch v := value.(type) {
		case string:
			length = len(v)
		case []interface{}:
			length = len(v)
		default:
			return fmt.Errorf("field type does not support length validation")
		}

		if length < min {
			return fmt.Errorf("minimum length is %d, got %d", min, length)
		}

		return nil
	}
}

// MaxLength validates maximum length for strings and slices
func MaxLength(max int) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		var length int
		switch v := value.(type) {
		case string:
			length = len(v)
		case []interface{}:
			length = len(v)
		default:
			return fmt.Errorf("field type does not support length validation")
		}

		if length > max {
			return fmt.Errorf("maximum length is %d, got %d", max, length)
		}

		return nil
	}
}

// Min validates minimum value for numbers
func Min(min float64) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		var numValue float64
		var err error

		switch v := value.(type) {
		case int:
			numValue = float64(v)
		case int64:
			numValue = float64(v)
		case float32:
			numValue = float64(v)
		case float64:
			numValue = v
		case string:
			numValue, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("field must be a number")
			}
		default:
			return fmt.Errorf("field must be a number")
		}

		if numValue < min {
			return fmt.Errorf("minimum value is %.2f, got %.2f", min, numValue)
		}

		return nil
	}
}

// Max validates maximum value for numbers
func Max(max float64) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		var numValue float64
		var err error

		switch v := value.(type) {
		case int:
			numValue = float64(v)
		case int64:
			numValue = float64(v)
		case float32:
			numValue = float64(v)
		case float64:
			numValue = v
		case string:
			numValue, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("field must be a number")
			}
		default:
			return fmt.Errorf("field must be a number")
		}

		if numValue > max {
			return fmt.Errorf("maximum value is %.2f, got %.2f", max, numValue)
		}

		return nil
	}
}

// Email validates email format
func Email() ValidationRule {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !emailRegex.MatchString(str) {
			return fmt.Errorf("field must be a valid email address")
		}

		return nil
	}
}

// URL validates URL format
func URL() ValidationRule {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !urlRegex.MatchString(str) {
			return fmt.Errorf("field must be a valid URL")
		}

		return nil
	}
}

// Alpha validates alphabetic characters only
func Alpha() ValidationRule {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !alphaRegex.MatchString(str) {
			return fmt.Errorf("field must contain only alphabetic characters")
		}

		return nil
	}
}

// Numeric validates numeric characters only
func Numeric() ValidationRule {
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !numericRegex.MatchString(str) {
			return fmt.Errorf("field must contain only numeric characters")
		}

		return nil
	}
}

// AlphaNumeric validates alphanumeric characters only
func AlphaNumeric() ValidationRule {
	alphaNumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !alphaNumericRegex.MatchString(str) {
			return fmt.Errorf("field must contain only alphanumeric characters")
		}

		return nil
	}
}

// OneOf validates that value is one of allowed values
func OneOf(allowedValues ...interface{}) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		for _, allowed := range allowedValues {
			if reflect.DeepEqual(value, allowed) {
				return nil
			}
		}

		return fmt.Errorf("field must be one of: %v", allowedValues)
	}
}

// Date validates date format
func Date(layout string) ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if _, err := time.Parse(layout, str); err != nil {
			return fmt.Errorf("field must be a valid date in format %s", layout)
		}

		return nil
	}
}

// UUID validates UUID format
func UUID() ValidationRule {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !uuidRegex.MatchString(strings.ToLower(str)) {
			return fmt.Errorf("field must be a valid UUID")
		}

		return nil
	}
}

// JSON validates JSON format
func JSON() ValidationRule {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		// Try to parse as JSON
		var js interface{}
		if err := json.Unmarshal([]byte(str), &js); err != nil {
			return fmt.Errorf("field must be valid JSON")
		}

		return nil
	}
}

// MCP-specific validation rules

// MCPMethod validates MCP method names
func MCPMethod() ValidationRule {
	validMethods := []string{
		"initialize",
		"ping",
		"tools/list",
		"tools/call",
		"resources/list",
		"resources/read",
		"resources/subscribe",
		"resources/unsubscribe",
		"prompts/list",
		"prompts/get",
		"notifications/subscribe",
		"notifications/unsubscribe",
	}

	return OneOf(interfaceSlice(validMethods)...)
}

// MCPProtocolVersion validates MCP protocol version
func MCPProtocolVersion() ValidationRule {
	validVersions := []string{"2024-11-05"}
	return OneOf(interfaceSlice(validVersions)...)
}

// MCPToolName validates tool names for this server
func MCPToolName() ValidationRule {
	validTools := []string{
		"get_movie",
		"add_movie",
		"update_movie",
		"delete_movie",
		"search_movies",
		"list_top_movies",
		"get_movie_poster",
		"add_actor",
		"link_actor_to_movie",
		"get_movie_cast",
		"get_actor_movies",
		"search_by_decade",
		"search_by_rating_range",
		"search_similar_movies",
	}

	return OneOf(interfaceSlice(validTools)...)
}

// MCPResourceURI validates resource URIs
func MCPResourceURI() ValidationRule {
	uriRegex := regexp.MustCompile(`^movies://(database|posters)/.+$`)
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if !uriRegex.MatchString(str) {
			return fmt.Errorf("field must be a valid movies:// URI")
		}

		return nil
	}
}

// MovieRating validates movie rating (1-10)
func MovieRating() ValidationRule {
	return func(value interface{}) error {
		rule := Min(1.0)
		if err := rule(value); err != nil {
			return err
		}
		
		rule = Max(10.0)
		if err := rule(value); err != nil {
			return err
		}

		return nil
	}
}

// getFieldNames extracts field names from validation error messages
func getFieldNames(validationErrors []string) []string {
	fieldNames := make([]string, 0, len(validationErrors))
	for _, err := range validationErrors {
		if colonIndex := strings.Index(err, ":"); colonIndex != -1 {
			fieldName := strings.TrimSpace(err[:colonIndex])
			fieldNames = append(fieldNames, fieldName)
		}
	}
	return fieldNames
}

// Helper function to convert string slice to interface slice
func interfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// RequestValidator provides validation for MCP requests
type RequestValidator struct {
	validator *Validator
}

// NewRequestValidator creates a new request validator
func NewRequestValidator() *RequestValidator {
	v := NewValidator()
	return &RequestValidator{validator: v}
}

// ValidateInitializeRequest validates initialize request parameters
func (rv *RequestValidator) ValidateInitializeRequest(params map[string]interface{}) error {
	validator := NewValidator()
	validator.AddRule("protocolVersion", Required())
	validator.AddRule("protocolVersion", MCPProtocolVersion())
	validator.AddRule("capabilities", Required())
	validator.AddRule("clientInfo", Required())

	return validator.Validate(params)
}

// ValidateToolsCallRequest validates tools/call request parameters
func (rv *RequestValidator) ValidateToolsCallRequest(params map[string]interface{}) error {
	validator := NewValidator()
	validator.AddRule("name", Required())
	validator.AddRule("name", MCPToolName())

	return validator.Validate(params)
}

// ValidateResourcesReadRequest validates resources/read request parameters
func (rv *RequestValidator) ValidateResourcesReadRequest(params map[string]interface{}) error {
	validator := NewValidator()
	validator.AddRule("uri", Required())
	validator.AddRule("uri", MCPResourceURI())

	return validator.Validate(params)
}

// ValidateMovieData validates movie data for add operations (requires title)
func (rv *RequestValidator) ValidateMovieData(args map[string]interface{}) error {
	return rv.validateMovieData(args, true)
}

// ValidateMovieUpdateData validates movie data for update operations (title optional)
func (rv *RequestValidator) ValidateMovieUpdateData(args map[string]interface{}) error {
	return rv.validateMovieData(args, false)
}

// validateMovieData validates movie data with optional title requirement
func (rv *RequestValidator) validateMovieData(args map[string]interface{}, requireTitle bool) error {
	validator := NewValidator()
	
	// Title validation (required for add, optional for update)
	if requireTitle {
		validator.AddRule("title", Required())
	}
	if _, exists := args["title"]; exists {
		validator.AddRule("title", MinLength(1))
		validator.AddRule("title", MaxLength(255))
	}
	
	// Optional but validated fields
	if _, exists := args["director"]; exists {
		validator.AddRule("director", MinLength(1))
		validator.AddRule("director", MaxLength(255))
	}
	
	if _, exists := args["year"]; exists {
		validator.AddRule("year", Min(1800))
		validator.AddRule("year", Max(float64(time.Now().Year()+10)))
	}
	
	if _, exists := args["rating"]; exists {
		validator.AddRule("rating", MovieRating())
	}
	
	if _, exists := args["genre"]; exists {
		validator.AddRule("genre", MinLength(1))
		validator.AddRule("genre", MaxLength(100))
	}

	// Validate poster_url if provided
	if _, exists := args["poster_url"]; exists {
		validator.AddRule("poster_url", MaxLength(500))
		// Allow empty string for removal
	}

	return validator.Validate(args)
}

// ValidateSearchQuery validates search query parameters
func (rv *RequestValidator) ValidateSearchQuery(args map[string]interface{}) error {
	validator := NewValidator()
	
	// At least one search parameter is required
	hasQuery := false
	if query, exists := args["query"]; exists && query != nil {
		if str, ok := query.(string); ok && strings.TrimSpace(str) != "" {
			hasQuery = true
			validator.AddRule("query", MinLength(1))
			validator.AddRule("query", MaxLength(500))
		}
	}
	
	if title, exists := args["title"]; exists && title != nil {
		if str, ok := title.(string); ok && strings.TrimSpace(str) != "" {
			hasQuery = true
			validator.AddRule("title", MinLength(1))
			validator.AddRule("title", MaxLength(255))
		}
	}
	
	if director, exists := args["director"]; exists && director != nil {
		if str, ok := director.(string); ok && strings.TrimSpace(str) != "" {
			hasQuery = true
			validator.AddRule("director", MinLength(1))
			validator.AddRule("director", MaxLength(255))
		}
	}
	
	if genre, exists := args["genre"]; exists && genre != nil {
		if str, ok := genre.(string); ok && strings.TrimSpace(str) != "" {
			hasQuery = true
			validator.AddRule("genre", MinLength(1))
			validator.AddRule("genre", MaxLength(100))
		}
	}
	
	if !hasQuery {
		return errors.NewValidationError(
			"At least one search parameter must be provided",
			"search_query",
			map[string]interface{}{
				"allowed_parameters": []string{"query", "title", "director", "genre"},
			},
		)
	}
	
	// Validate limit if provided
	if _, exists := args["limit"]; exists {
		validator.AddRule("limit", Min(1))
		validator.AddRule("limit", Max(1000))
	}

	return validator.Validate(args)
}