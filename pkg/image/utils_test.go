package image

import (
	"testing"

	"movies-mcp-server/internal/config"
)

func TestImageProcessor_ValidateImage(t *testing.T) {
	cfg := &config.ImageConfig{
		MaxSize:      1024 * 1024, // 1MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
	}
	processor := NewImageProcessor(cfg)

	tests := []struct {
		name        string
		data        []byte
		mimeType    string
		expectError bool
	}{
		{
			name:        "Valid JPEG",
			data:        createTestJPEGData(),
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name:        "Valid PNG",
			data:        createTestPNGData(),
			mimeType:    "image/png",
			expectError: false,
		},
		{
			name:        "Empty data",
			data:        []byte{},
			mimeType:    "image/jpeg",
			expectError: true,
		},
		{
			name:        "Invalid MIME type",
			data:        createTestJPEGData(),
			mimeType:    "image/gif",
			expectError: true,
		},
		{
			name:        "MIME type mismatch",
			data:        createTestJPEGData(),
			mimeType:    "image/png",
			expectError: true,
		},
		{
			name:        "Data too large",
			data:        make([]byte, 2*1024*1024), // 2MB
			mimeType:    "image/jpeg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.ValidateImage(tt.data, tt.mimeType)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateImage() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestImageProcessor_Base64EncodingDecoding(t *testing.T) {
	cfg := &config.ImageConfig{
		MaxSize:      1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png"},
	}
	processor := NewImageProcessor(cfg)

	originalData := []byte("test image data")

	// Test encoding
	encoded := processor.EncodeToBase64(originalData)
	if encoded == "" {
		t.Errorf("EncodeToBase64() returned empty string")
	}

	// Test decoding
	decoded, err := processor.DecodeFromBase64(encoded)
	if err != nil {
		t.Errorf("DecodeFromBase64() error = %v", err)
	}

	if string(decoded) != string(originalData) {
		t.Errorf("Decoded data doesn't match original: got %s, want %s", string(decoded), string(originalData))
	}
}

func TestImageProcessor_DetectMimeType(t *testing.T) {
	cfg := &config.ImageConfig{}
	processor := NewImageProcessor(cfg)

	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "JPEG data",
			data:     createTestJPEGData(),
			expected: "image/jpeg",
		},
		{
			name:     "PNG data",
			data:     createTestPNGData(),
			expected: "image/png",
		},
		{
			name:     "Unknown data",
			data:     []byte("not an image"),
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.detectMimeType(tt.data)
			if result != tt.expected {
				t.Errorf("detectMimeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestImageProcessor_ParseThumbnailSize(t *testing.T) {
	tests := []struct {
		name        string
		sizeStr     string
		expectError bool
		width       int
		height      int
	}{
		{
			name:        "Valid size",
			sizeStr:     "200x200",
			expectError: false,
			width:       200,
			height:      200,
		},
		{
			name:        "Different dimensions",
			sizeStr:     "150x100",
			expectError: false,
			width:       150,
			height:      100,
		},
		{
			name:        "Invalid format",
			sizeStr:     "200",
			expectError: true,
		},
		{
			name:        "Non-numeric width",
			sizeStr:     "abcx200",
			expectError: true,
		},
		{
			name:        "Zero dimensions",
			sizeStr:     "0x0",
			expectError: true,
		},
		{
			name:        "Negative dimensions",
			sizeStr:     "-100x200",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ImageConfig{
				ThumbnailSize: tt.sizeStr,
			}
			processor := NewImageProcessor(cfg)

			size, err := processor.parseThumbnailSize()
			if (err != nil) != tt.expectError {
				t.Errorf("parseThumbnailSize() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if size.width != tt.width || size.height != tt.height {
					t.Errorf("parseThumbnailSize() = %dx%d, want %dx%d", 
						size.width, size.height, tt.width, tt.height)
				}
			}
		})
	}
}

// Helper functions to create test image data
func createTestJPEGData() []byte {
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

func BenchmarkImageProcessor_Base64Encoding(b *testing.B) {
	cfg := &config.ImageConfig{}
	processor := NewImageProcessor(cfg)
	data := make([]byte, 100*1024) // 100KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoded := processor.EncodeToBase64(data)
		_, err := processor.DecodeFromBase64(encoded)
		if err != nil {
			b.Fatalf("Decode error: %v", err)
		}
	}
}