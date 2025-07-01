package mcp

import (
	"testing"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/interfaces/dto"
)

func TestPromptHandlers_HandlePromptsList(t *testing.T) {
	handlers := NewPromptHandlers()

	var result interface{}
	var errorCalled bool

	sendResult := func(id interface{}, res interface{}) {
		result = res
	}

	sendError := func(id interface{}, code int, message string, data interface{}) {
		errorCalled = true
		t.Errorf("Unexpected error: %s", message)
	}

	// Test prompts list
	handlers.HandlePromptsList(1, nil, sendResult, sendError)

	if errorCalled {
		t.Fatal("Handler returned an error")
	}

	response, ok := result.(dto.PromptsListResponse)
	if !ok {
		t.Fatal("Expected PromptsListResponse")
	}

	// Verify we have prompts
	if len(response.Prompts) == 0 {
		t.Error("Expected at least one prompt template")
	}

	// Verify prompt structure
	expectedPrompts := []string{
		"movie_recommendation",
		"movie_analysis",
		"director_filmography",
		"genre_exploration",
		"movie_comparison",
	}

	promptMap := make(map[string]bool)
	for _, prompt := range response.Prompts {
		promptMap[prompt.Name] = true

		// Verify each prompt has required fields
		if prompt.Name == "" {
			t.Error("Prompt missing name")
		}
		if prompt.Description == "" {
			t.Error("Prompt missing description")
		}
	}

	// Verify all expected prompts are present
	for _, expected := range expectedPrompts {
		if !promptMap[expected] {
			t.Errorf("Missing expected prompt: %s", expected)
		}
	}
}

func TestPromptHandlers_HandlePromptGet(t *testing.T) {
	handlers := NewPromptHandlers()

	tests := []struct {
		name      string
		args      map[string]interface{}
		expectErr bool
		errMsg    string
	}{
		{
			name:      "Missing prompt name",
			args:      map[string]interface{}{},
			expectErr: true,
			errMsg:    "Missing or invalid prompt name",
		},
		{
			name: "Unknown prompt",
			args: map[string]interface{}{
				"name": "unknown_prompt",
			},
			expectErr: true,
			errMsg:    "Unknown prompt",
		},
		{
			name: "Valid movie recommendation - minimal",
			args: map[string]interface{}{
				"name": "movie_recommendation",
				"arguments": map[string]interface{}{
					"genre": "Action",
				},
			},
			expectErr: false,
		},
		{
			name: "Valid movie recommendation - full",
			args: map[string]interface{}{
				"name": "movie_recommendation",
				"arguments": map[string]interface{}{
					"genre":      "Sci-Fi",
					"min_rating": 8.0,
					"year_range": "2020-2024",
				},
			},
			expectErr: false,
		},
		{
			name: "Missing required argument",
			args: map[string]interface{}{
				"name":      "movie_recommendation",
				"arguments": map[string]interface{}{},
			},
			expectErr: true,
			errMsg:    "Missing required argument: genre",
		},
		{
			name: "Valid movie analysis",
			args: map[string]interface{}{
				"name": "movie_analysis",
				"arguments": map[string]interface{}{
					"movie_title": "The Matrix",
				},
			},
			expectErr: false,
		},
		{
			name: "Valid movie comparison",
			args: map[string]interface{}{
				"name": "movie_comparison",
				"arguments": map[string]interface{}{
					"movie1": "The Matrix",
					"movie2": "Inception",
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			var errorCode int
			var errorMsg string
			errorCalled := false

			sendResult := func(id interface{}, res interface{}) {
				result = res
			}

			sendError := func(id interface{}, code int, message string, data interface{}) {
				errorCalled = true
				errorCode = code
				errorMsg = message
			}

			handlers.HandlePromptGet(1, tt.args, sendResult, sendError)

			if tt.expectErr {
				if !errorCalled {
					t.Error("Expected error but none occurred")
				}
				if errorCode != dto.InvalidParams {
					t.Errorf("Expected InvalidParams error code, got %d", errorCode)
				}
				if !containsString(errorMsg, tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errMsg, errorMsg)
				}
			} else {
				if errorCalled {
					t.Errorf("Unexpected error: %s", errorMsg)
				}

				response, ok := result.(dto.PromptGetResponse)
				if !ok {
					t.Fatal("Expected PromptGetResponse")
				}

				// Verify response has messages
				if len(response.Messages) == 0 {
					t.Error("Expected at least one message in prompt response")
				}

				// Verify message structure
				for _, msg := range response.Messages {
					if msg.Role == "" {
						t.Error("Message missing role")
					}
					if len(msg.Content) == 0 || msg.Content[0].Type != "text" {
						t.Errorf("Expected text content type, got %v", msg.Content)
					}
					if len(msg.Content) == 0 || msg.Content[0].Text == "" {
						t.Error("Message missing text content")
					}
				}
			}
		})
	}
}

func TestPromptHandlers_PromptGeneration(t *testing.T) {
	handlers := NewPromptHandlers()

	// Test movie recommendation prompt generation
	t.Run("Movie recommendation prompt", func(t *testing.T) {
		args := map[string]interface{}{
			"genre":      "Action",
			"min_rating": 8.5,
			"year_range": "2020-2024",
		}

		response := handlers.generateMovieRecommendationPrompt(args)

		if len(response.Messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(response.Messages))
		}

		msg := response.Messages[0]
		if msg.Role != "user" {
			t.Errorf("Expected user role, got %s", msg.Role)
		}

		// Verify prompt contains expected elements
		text := msg.Content[0].Text
		expectedStrings := []string{"Action", "8.5", "2020-2024"}
		for _, expected := range expectedStrings {
			if !containsString(text, expected) {
				t.Errorf("Expected prompt to contain '%s'", expected)
			}
		}
	})

	// Test director filmography prompt generation
	t.Run("Director filmography prompt", func(t *testing.T) {
		args := map[string]interface{}{
			"director_name": "Christopher Nolan",
			"focus_period":  "2010-2020",
		}

		response := handlers.generateDirectorFilmographyPrompt(args)

		text := response.Messages[0].Content[0].Text
		if !containsString(text, "Christopher Nolan") {
			t.Error("Expected prompt to contain director name")
		}
		if !containsString(text, "2010-2020") {
			t.Error("Expected prompt to contain focus period")
		}
	})
}

// Helper function
func containsString(text, substr string) bool {
	return len(substr) > 0 && len(text) >= len(substr) &&
		(text == substr || containsSubstring(text, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
