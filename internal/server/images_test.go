package server

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"bytes"
	"strings"
	"testing"

	"movies-mcp-server/internal/database"
	"movies-mcp-server/internal/models"
)

// Test image encoding/decoding functionality
func TestImageEncodingDecoding(t *testing.T) {
	// Create a simple test image (1x1 pixel)
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // Red pixel

	tests := []struct {
		name     string
		encoder  func(*bytes.Buffer, image.Image) error
		mimeType string
	}{
		{
			name:     "JPEG encoding",
			encoder:  func(b *bytes.Buffer, img image.Image) error { return jpeg.Encode(b, img, nil) },
			mimeType: "image/jpeg",
		},
		{
			name:     "PNG encoding", 
			encoder:  func(b *bytes.Buffer, img image.Image) error { return png.Encode(b, img) },
			mimeType: "image/png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode image to bytes
			var buf bytes.Buffer
			if err := tt.encoder(&buf, img); err != nil {
				t.Fatalf("Failed to encode image: %v", err)
			}

			imageData := buf.Bytes()

			// Test base64 encoding
			encoded := base64.StdEncoding.EncodeToString(imageData)
			if encoded == "" {
				t.Errorf("Base64 encoding should not be empty")
			}

			// Test base64 decoding
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				t.Errorf("Failed to decode base64: %v", err)
			}

			if !bytes.Equal(imageData, decoded) {
				t.Errorf("Decoded data doesn't match original")
			}

			// Test image validation
			if !isValidImageData(imageData, tt.mimeType) {
				t.Errorf("Image data should be valid for mime type %s", tt.mimeType)
			}
		})
	}
}

// Helper function to validate image data (to be implemented in actual code)
func isValidImageData(data []byte, mimeType string) bool {
	switch mimeType {
	case "image/jpeg":
		// JPEG files start with FF D8 FF
		return len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
	case "image/png":
		// PNG files start with specific signature
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if len(data) < len(pngHeader) {
			return false
		}
		return bytes.Equal(data[:len(pngHeader)], pngHeader)
	case "image/webp":
		// WebP files start with "RIFF" and contain "WEBP"
		return len(data) >= 12 && 
			   string(data[0:4]) == "RIFF" && 
			   string(data[8:12]) == "WEBP"
	default:
		return false
	}
}

// Test image size validation
func TestImageSizeValidation(t *testing.T) {
	tests := []struct {
		name      string
		imageSize int64
		maxSize   int64
		shouldPass bool
	}{
		{
			name:      "Image within size limit",
			imageSize: 1024,
			maxSize:   5 * 1024 * 1024, // 5MB
			shouldPass: true,
		},
		{
			name:      "Image exceeds size limit",
			imageSize: 10 * 1024 * 1024, // 10MB
			maxSize:   5 * 1024 * 1024,  // 5MB
			shouldPass: false,
		},
		{
			name:      "Image at exact size limit",
			imageSize: 5 * 1024 * 1024, // 5MB
			maxSize:   5 * 1024 * 1024, // 5MB
			shouldPass: true,
		},
		{
			name:      "Zero size image",
			imageSize: 0,
			maxSize:   5 * 1024 * 1024,
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateImageSize(tt.imageSize, tt.maxSize)
			if valid != tt.shouldPass {
				t.Errorf("Expected validation result %v, got %v", tt.shouldPass, valid)
			}
		})
	}
}

// Helper function for image size validation
func validateImageSize(imageSize, maxSize int64) bool {
	return imageSize > 0 && imageSize <= maxSize
}

// Test MIME type validation
func TestImageMimeTypeValidation(t *testing.T) {
	allowedTypes := []string{"image/jpeg", "image/png", "image/webp"}

	tests := []struct {
		mimeType   string
		shouldPass bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/webp", true},
		{"image/gif", false},
		{"image/bmp", false},
		{"text/plain", false},
		{"application/json", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run("MimeType_"+tt.mimeType, func(t *testing.T) {
			valid := isValidMimeType(tt.mimeType, allowedTypes)
			if valid != tt.shouldPass {
				t.Errorf("Expected mime type %s validation to be %v, got %v", 
					tt.mimeType, tt.shouldPass, valid)
			}
		})
	}
}

// Helper function for MIME type validation
func isValidMimeType(mimeType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// Test add_movie tool with poster URL handling
func TestHandleAddMovieWithPoster(t *testing.T) {
	server, _, output := createTestServer()

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
		description string
	}{
		{
			name: "Add movie with valid poster URL",
			args: map[string]interface{}{
				"title":      "Movie with Poster",
				"director":   "Test Director",
				"year":       2023,
				"poster_url": "https://example.com/poster.jpg",
			},
			expectError: false,
			description: "Should accept valid poster URL",
		},
		{
			name: "Add movie without poster",
			args: map[string]interface{}{
				"title":    "Movie without Poster",
				"director": "Test Director",
				"year":     2023,
			},
			expectError: false,
			description: "Should work without poster URL",
		},
		{
			name: "Add movie with invalid poster URL",
			args: map[string]interface{}{
				"title":      "Movie with Bad Poster",
				"director":   "Test Director", 
				"year":       2023,
				"poster_url": "not-a-valid-url",
			},
			expectError: false, // Should not fail, just ignore invalid URL
			description: "Should handle invalid poster URL gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			server.handleAddMovie(1, tt.args)

			response := parseResponse(t, output.buffer.String())

			if tt.expectError && response.Error == nil {
				t.Errorf("Expected error but got none: %s", tt.description)
			} else if !tt.expectError && response.Error != nil {
				t.Errorf("Unexpected error: %v (%s)", response.Error, tt.description)
			}

			if !tt.expectError && response.Result != nil {
				// Verify movie was created successfully
				var result models.ToolCallResponse
				resultBytes, _ := json.Marshal(response.Result)
				json.Unmarshal(resultBytes, &result)

				if len(result.Content) == 0 || !strings.Contains(result.Content[0].Text, "Successfully created") {
					t.Errorf("Expected success message, got: %v", result.Content)
				}
			}
		})
	}
}

// Test update_movie tool with poster updates
func TestHandleUpdateMovieWithPoster(t *testing.T) {
	server, db, output := createTestServer()

	// Create initial movie
	testMovie := &database.Movie{
		ID:       1,
		Title:    "Original Movie",
		Director: "Original Director",
		Year:     2020,
	}
	db.AddTestMovie(testMovie)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
		description string
	}{
		{
			name: "Update movie with new poster URL",
			args: map[string]interface{}{
				"id":         1,
				"poster_url": "https://example.com/new-poster.jpg",
			},
			expectError: false,
			description: "Should accept poster URL updates",
		},
		{
			name: "Update movie removing poster",
			args: map[string]interface{}{
				"id":         1,
				"poster_url": "",
			},
			expectError: false,
			description: "Should handle empty poster URL (removal)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			server.handleUpdateMovie(1, tt.args)

			response := parseResponse(t, output.buffer.String())

			if tt.expectError && response.Error == nil {
				t.Errorf("Expected error but got none: %s", tt.description)
			} else if !tt.expectError && response.Error != nil {
				t.Errorf("Unexpected error: %v (%s)", response.Error, tt.description)
			}
		})
	}
}

// Test poster resource access with different formats
func TestPosterResourceFormats(t *testing.T) {
	server, db, output := createTestServer()

	// Create test images in different formats
	testMovies := []*database.Movie{
		{
			ID:         1,
			Title:      "JPEG Movie",
			PosterData: createTestJPEGData(),
			PosterType: "image/jpeg",
		},
		{
			ID:         2,
			Title:      "PNG Movie", 
			PosterData: createTestPNGData(),
			PosterType: "image/png",
		},
		{
			ID:         3,
			Title:      "Movie without poster",
		},
	}

	for _, movie := range testMovies {
		db.AddTestMovie(movie)
	}

	tests := []struct {
		movieID     int
		expectError bool
		expectedMime string
	}{
		{1, false, "image/jpeg"},
		{2, false, "image/png"},
		{3, true, ""},  // No poster
		{999, true, ""}, // Non-existent movie
	}

	for _, tt := range tests {
		t.Run("Movie_"+string(rune(tt.movieID)), func(t *testing.T) {
			output.buffer.Reset()

			params := models.ResourceReadRequest{
				URI: "movies://posters/" + string(rune(tt.movieID+'0')),
			}
			paramsBytes, _ := json.Marshal(params)

			req := &models.JSONRPCRequest{
				ID:     1,
				Method: "resources/read",
				Params: paramsBytes,
			}

			server.handleResourceRead(req)

			response := parseResponse(t, output.buffer.String())

			if tt.expectError {
				if response.Error == nil {
					t.Errorf("Expected error for movie %d but got none", tt.movieID)
				}
			} else {
				if response.Error != nil {
					t.Errorf("Unexpected error for movie %d: %v", tt.movieID, response.Error)
				} else {
					var result models.ResourceReadResponse
					resultBytes, _ := json.Marshal(response.Result)
					json.Unmarshal(resultBytes, &result)

					if len(result.Contents) > 0 {
						content := result.Contents[0]
						if content.MimeType != tt.expectedMime {
							t.Errorf("Expected mime type %s, got %s", tt.expectedMime, content.MimeType)
						}

						if content.Blob == "" {
							t.Errorf("Expected blob data to be present")
						}

						// Verify base64 decoding works
						if _, err := base64.StdEncoding.DecodeString(content.Blob); err != nil {
							t.Errorf("Blob should be valid base64: %v", err)
						}
					}
				}
			}
		})
	}
}

// Helper functions to create test image data
func createTestJPEGData() []byte {
	// Create a minimal valid JPEG header
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, // JPEG SOI + APP0
		0x00, 0x10, // APP0 length
		0x4A, 0x46, 0x49, 0x46, 0x00, // "JFIF\0"
		0x01, 0x01, // Version
		0x01, // Units
		0x00, 0x01, 0x00, 0x01, // X/Y density
		0x00, 0x00, // Thumbnail dimensions
		0xFF, 0xD9, // EOI
	}
}

func createTestPNGData() []byte {
	// Create PNG signature
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, // IHDR chunk length
		0x49, 0x48, 0x44, 0x52, // "IHDR"
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, // Bit depth, color type, etc
		0x90, 0x77, 0x53, 0xDE, // CRC
		0x00, 0x00, 0x00, 0x00, // IEND chunk length
		0x49, 0x45, 0x4E, 0x44, // "IEND"
		0xAE, 0x42, 0x60, 0x82, // CRC
	}
}

// Test thumbnail generation (placeholder for future implementation)
func TestThumbnailGeneration(t *testing.T) {
	tests := []struct {
		name           string
		originalSize   string
		thumbnailSize  string
		expectedResult bool
	}{
		{
			name:           "Generate 200x200 thumbnail from 800x600 image",
			originalSize:   "800x600",
			thumbnailSize:  "200x200",
			expectedResult: true,
		},
		{
			name:           "Generate 150x150 thumbnail from 1920x1080 image",
			originalSize:   "1920x1080",
			thumbnailSize:  "150x150",
			expectedResult: true,
		},
		{
			name:           "Invalid thumbnail size",
			originalSize:   "800x600",
			thumbnailSize:  "0x0",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a placeholder test for thumbnail generation
			// The actual implementation would generate thumbnails
			result := shouldGenerateThumbnail(tt.originalSize, tt.thumbnailSize)
			if result != tt.expectedResult {
				t.Errorf("Expected thumbnail generation result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

// Helper function for thumbnail generation validation
func shouldGenerateThumbnail(originalSize, thumbnailSize string) bool {
	return thumbnailSize != "0x0" && originalSize != "" && thumbnailSize != ""
}

// Benchmark image encoding/decoding performance
func BenchmarkImageBase64Encoding(b *testing.B) {
	// Create test image data (100KB)
	imageData := make([]byte, 100*1024)
	for i := range imageData {
		imageData[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoded := base64.StdEncoding.EncodeToString(imageData)
		_, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Decode error: %v", err)
		}
	}
}

func BenchmarkImageValidation(b *testing.B) {
	jpegData := createTestJPEGData()
	pngData := createTestPNGData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isValidImageData(jpegData, "image/jpeg")
		isValidImageData(pngData, "image/png")
	}
}