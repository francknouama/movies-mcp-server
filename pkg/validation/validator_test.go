package validation

import (
	"strings"
	"testing"
	"time"
)

func TestRequired(t *testing.T) {
	rule := Required()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, true},
		{"empty string", "", true},
		{"whitespace string", "   ", true},
		{"valid string", "hello", false},
		{"empty slice", []interface{}{}, true},
		{"valid slice", []interface{}{1, 2, 3}, false},
		{"empty map", map[string]interface{}{}, false}, // Empty maps are allowed
		{"valid map", map[string]interface{}{"key": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	rule := MinLength(3)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"short string", "ab", true},
		{"exact length string", "abc", false},
		{"long string", "abcd", false},
		{"short slice", []interface{}{1, 2}, true},
		{"exact length slice", []interface{}{1, 2, 3}, false},
		{"long slice", []interface{}{1, 2, 3, 4}, false},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	rule := MaxLength(5)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"short string", "abc", false},
		{"exact length string", "abcde", false},
		{"long string", "abcdef", true},
		{"short slice", []interface{}{1, 2, 3}, false},
		{"exact length slice", []interface{}{1, 2, 3, 4, 5}, false},
		{"long slice", []interface{}{1, 2, 3, 4, 5, 6}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMin(t *testing.T) {
	rule := Min(10.0)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"low int", 5, true},
		{"exact int", 10, false},
		{"high int", 15, false},
		{"low float", 5.5, true},
		{"exact float", 10.0, false},
		{"high float", 15.5, false},
		{"string number low", "5", true},
		{"string number exact", "10", false},
		{"string number high", "15", false},
		{"invalid string", "abc", true},
		{"invalid type", []int{1, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Min() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMax(t *testing.T) {
	rule := Max(100.0)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"low int", 50, false},
		{"exact int", 100, false},
		{"high int", 150, true},
		{"low float", 50.5, false},
		{"exact float", 100.0, false},
		{"high float", 150.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Max() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	rule := Email()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "user@mail.example.org", false},
		{"invalid email - no @", "testexample.com", true},
		{"invalid email - no domain", "test@", true},
		{"invalid email - no TLD", "test@example", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL(t *testing.T) {
	rule := URL()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid http URL", "http://example.com", false},
		{"valid https URL", "https://example.com/path", false},
		{"invalid URL - no protocol", "example.com", true},
		{"invalid URL - wrong protocol", "ftp://example.com", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlpha(t *testing.T) {
	rule := Alpha()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid alpha", "abcDEF", false},
		{"invalid - numbers", "abc123", true},
		{"invalid - spaces", "abc def", true},
		{"invalid - special chars", "abc@def", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Alpha() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	rule := Numeric()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid numeric", "12345", false},
		{"invalid - letters", "123abc", true},
		{"invalid - spaces", "123 456", true},
		{"invalid - special chars", "123.45", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Numeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlphaNumeric(t *testing.T) {
	rule := AlphaNumeric()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid alphanumeric", "abc123DEF", false},
		{"invalid - spaces", "abc 123", true},
		{"invalid - special chars", "abc@123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaNumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	rule := OneOf("apple", "banana", "cherry")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid value", "apple", false},
		{"invalid value", "grape", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("OneOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate(t *testing.T) {
	rule := Date("2006-01-02")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid date", "2023-12-25", false},
		{"invalid date format", "25-12-2023", true},
		{"invalid date", "2023-13-45", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	rule := UUID()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", false},
		{"invalid UUID - wrong format", "550e8400-e29b-41d4-a716", true},
		{"invalid UUID - wrong chars", "550e8400-e29b-41d4-a716-44665544000g", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	rule := JSON()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"valid JSON object", `{"key": "value"}`, false},
		{"valid JSON array", `[1, 2, 3]`, false},
		{"valid JSON string", `"hello"`, false},
		{"invalid JSON", `{"key": }`, true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPMethod(t *testing.T) {
	rule := MCPMethod()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid method", "initialize", false},
		{"valid method", "tools/list", false},
		{"invalid method", "invalid/method", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MCPMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPProtocolVersion(t *testing.T) {
	rule := MCPProtocolVersion()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid version", "2024-11-05", false},
		{"invalid version", "1.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MCPProtocolVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPToolName(t *testing.T) {
	rule := MCPToolName()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid tool", "get_movie", false},
		{"valid tool", "search_movies", false},
		{"invalid tool", "invalid_tool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MCPToolName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPResourceURI(t *testing.T) {
	rule := MCPResourceURI()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid database URI", "movies://database/stats", false},
		{"valid poster URI", "movies://posters/123", false},
		{"invalid scheme", "http://database/stats", true},
		{"invalid path", "movies://invalid/path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MCPResourceURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieRating(t *testing.T) {
	rule := MovieRating()

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid rating", 7.5, false},
		{"min rating", 1.0, false},
		{"max rating", 10.0, false},
		{"too low", 0.5, true},
		{"too high", 10.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRating() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()
	validator.AddRule("name", Required())
	validator.AddRule("name", MinLength(2))
	validator.AddRule("email", Email())

	tests := []struct {
		name    string
		values  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid data",
			values: map[string]interface{}{
				"name":  "John",
				"email": "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			values: map[string]interface{}{
				"email": "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			values: map[string]interface{}{
				"name":  "John",
				"email": "invalid-email",
			},
			wantErr: true,
		},
		{
			name: "short name",
			values: map[string]interface{}{
				"name":  "J",
				"email": "john@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestValidator_ValidateInitializeRequest(t *testing.T) {
	rv := NewRequestValidator()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid initialize request",
			params: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"clientInfo":      map[string]interface{}{"name": "test", "version": "1.0"},
			},
			wantErr: false,
		},
		{
			name: "missing protocol version",
			params: map[string]interface{}{
				"capabilities": map[string]interface{}{},
				"clientInfo":   map[string]interface{}{"name": "test", "version": "1.0"},
			},
			wantErr: true,
		},
		{
			name: "invalid protocol version",
			params: map[string]interface{}{
				"protocolVersion": "1.0.0",
				"capabilities":    map[string]interface{}{},
				"clientInfo":      map[string]interface{}{"name": "test", "version": "1.0"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rv.ValidateInitializeRequest(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInitializeRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestValidator_ValidateMovieData(t *testing.T) {
	rv := NewRequestValidator()

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid movie data",
			args: map[string]interface{}{
				"title":    "Test Movie",
				"director": "Test Director",
				"year":     2023,
				"rating":   7.5,
				"genre":    "Action",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			args: map[string]interface{}{
				"director": "Test Director",
				"year":     2023,
			},
			wantErr: true,
		},
		{
			name: "invalid rating",
			args: map[string]interface{}{
				"title":  "Test Movie",
				"rating": 15.0,
			},
			wantErr: true,
		},
		{
			name: "invalid year",
			args: map[string]interface{}{
				"title": "Test Movie",
				"year":  1500,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rv.ValidateMovieData(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMovieData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestValidator_ValidateSearchQuery(t *testing.T) {
	rv := NewRequestValidator()

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid search with query",
			args: map[string]interface{}{
				"query": "action movie",
				"limit": 10,
			},
			wantErr: false,
		},
		{
			name: "valid search with title",
			args: map[string]interface{}{
				"title": "Inception",
			},
			wantErr: false,
		},
		{
			name: "no search parameters",
			args: map[string]interface{}{
				"limit": 10,
			},
			wantErr: true,
		},
		{
			name: "empty query",
			args: map[string]interface{}{
				"query": "",
			},
			wantErr: true,
		},
		{
			name: "invalid limit",
			args: map[string]interface{}{
				"query": "action",
				"limit": 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rv.ValidateSearchQuery(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStructWithTags(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required,min=2"`
		Email string `validate:"email"`
		Age   int    `validate:"min=0,max=150"`
	}

	validator := NewValidator()

	tests := []struct {
		name    string
		data    TestStruct
		wantErr bool
	}{
		{
			name: "valid struct",
			data: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			data: TestStruct{
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			data: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationErrorMessages(t *testing.T) {
	validator := NewValidator()
	validator.AddRule("email", Required())
	validator.AddRule("email", Email())

	err := validator.Validate(map[string]interface{}{
		"email": "invalid-email",
	})

	if err == nil {
		t.Fatal("Expected validation error")
	}

	if !strings.Contains(err.Error(), "email") {
		t.Errorf("Expected error message to contain field name, got: %v", err.Error())
	}
}

func TestValidationPerformance(t *testing.T) {
	validator := NewValidator()
	validator.AddRule("field1", Required())
	validator.AddRule("field1", MinLength(1))
	validator.AddRule("field1", MaxLength(100))
	validator.AddRule("field2", Email())
	validator.AddRule("field3", Min(0))
	validator.AddRule("field3", Max(1000))

	data := map[string]interface{}{
		"field1": "valid string",
		"field2": "test@example.com",
		"field3": 500,
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		validator.Validate(data)
	}
	duration := time.Since(start)

	if duration > time.Second {
		t.Errorf("Validation took too long: %v", duration)
	}
}