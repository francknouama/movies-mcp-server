package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageProcessor_ValidateImage(t *testing.T) {
	cfg := &ImageConfig{
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
	cfg := &ImageConfig{
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
	cfg := &ImageConfig{}
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
			cfg := &ImageConfig{
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
		0x01,                   // Units
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

func TestNewImageProcessor(t *testing.T) {
	cfg := &ImageConfig{
		MaxSize:      1024 * 1024,
		AllowedTypes: []string{"image/jpeg"},
	}

	processor := NewImageProcessor(cfg)
	if processor == nil {
		t.Fatal("NewImageProcessor() returned nil")
	}

	if processor.config != cfg {
		t.Error("ImageProcessor.config is not set correctly")
	}
}

func TestDefaultImageConfig(t *testing.T) {
	cfg := DefaultImageConfig()

	if cfg == nil {
		t.Fatal("DefaultImageConfig() returned nil")
	}

	if cfg.MaxSize != 10*1024*1024 {
		t.Errorf("DefaultImageConfig() MaxSize = %d, want %d", cfg.MaxSize, 10*1024*1024)
	}

	if len(cfg.AllowedTypes) != 3 {
		t.Errorf("DefaultImageConfig() AllowedTypes length = %d, want 3", len(cfg.AllowedTypes))
	}

	if !cfg.EnableThumbnails {
		t.Error("DefaultImageConfig() EnableThumbnails should be true")
	}

	if cfg.ThumbnailSize != "200x200" {
		t.Errorf("DefaultImageConfig() ThumbnailSize = %s, want 200x200", cfg.ThumbnailSize)
	}
}

func TestImageProcessor_ValidateImageFormat_WebP(t *testing.T) {
	cfg := &ImageConfig{
		MaxSize:      1024 * 1024,
		AllowedTypes: []string{"image/webp"},
	}
	processor := NewImageProcessor(cfg)

	webpData := createTestWebPData()
	err := processor.ValidateImage(webpData, "image/webp")
	if err != nil {
		t.Errorf("ValidateImage() error = %v for valid WebP", err)
	}
}

func TestImageProcessor_DecodeFromBase64_InvalidData(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	// Test with invalid base64
	_, err := processor.DecodeFromBase64("!!!invalid base64!!!")
	if err == nil {
		t.Error("DecodeFromBase64() should return error for invalid base64")
	}
}

func TestImageProcessor_DetectMimeType_WebP(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	webpData := createTestWebPData()
	mimeType := processor.detectMimeType(webpData)
	if mimeType != "image/webp" {
		t.Errorf("detectMimeType() = %s, want image/webp", mimeType)
	}
}

func TestImageProcessor_GetImageInfo(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	tests := []struct {
		name        string
		data        []byte
		mimeType    string
		expectError bool
	}{
		{
			name:        "Valid JPEG",
			data:        createValidJPEGImage(),
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name:        "Valid PNG",
			data:        createValidPNGImage(),
			mimeType:    "image/png",
			expectError: false,
		},
		{
			name:        "Invalid image data",
			data:        []byte("not an image"),
			mimeType:    "image/jpeg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := processor.GetImageInfo(tt.data, tt.mimeType)
			if (err != nil) != tt.expectError {
				t.Errorf("GetImageInfo() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if info == nil {
					t.Fatal("GetImageInfo() returned nil info")
				}
				if info.Width <= 0 || info.Height <= 0 {
					t.Errorf("GetImageInfo() invalid dimensions: %dx%d", info.Width, info.Height)
				}
				if info.Size <= 0 {
					t.Errorf("GetImageInfo() invalid size: %d", info.Size)
				}
				if info.MimeType != tt.mimeType {
					t.Errorf("GetImageInfo() MimeType = %s, want %s", info.MimeType, tt.mimeType)
				}
			}
		})
	}
}

func TestImageProcessor_GetFormatFromMimeType(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	tests := []struct {
		mimeType       string
		expectedFormat string
	}{
		{"image/jpeg", "JPEG"},
		{"image/png", "PNG"},
		{"image/webp", "WebP"},
		{"image/unknown", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			format := processor.getFormatFromMimeType(tt.mimeType)
			if format != tt.expectedFormat {
				t.Errorf("getFormatFromMimeType(%s) = %s, want %s", tt.mimeType, format, tt.expectedFormat)
			}
		})
	}
}

func TestImageProcessor_GenerateThumbnail(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *ImageConfig
		data        []byte
		mimeType    string
		expectError bool
	}{
		{
			name: "Thumbnails disabled",
			cfg: &ImageConfig{
				EnableThumbnails: false,
			},
			data:        createValidJPEGImage(),
			mimeType:    "image/jpeg",
			expectError: true,
		},
		{
			name: "Invalid thumbnail size",
			cfg: &ImageConfig{
				EnableThumbnails: true,
				ThumbnailSize:    "invalid",
			},
			data:        createValidJPEGImage(),
			mimeType:    "image/jpeg",
			expectError: true,
		},
		{
			name: "Valid JPEG thumbnail",
			cfg: &ImageConfig{
				EnableThumbnails: true,
				ThumbnailSize:    "50x50",
			},
			data:        createValidJPEGImage(),
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name: "Valid PNG thumbnail",
			cfg: &ImageConfig{
				EnableThumbnails: true,
				ThumbnailSize:    "50x50",
			},
			data:        createValidPNGImage(),
			mimeType:    "image/png",
			expectError: false,
		},
		{
			name: "Invalid image data",
			cfg: &ImageConfig{
				EnableThumbnails: true,
				ThumbnailSize:    "50x50",
			},
			data:        []byte("not an image"),
			mimeType:    "image/jpeg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewImageProcessor(tt.cfg)
			thumbnail, err := processor.GenerateThumbnail(tt.data, tt.mimeType)

			if (err != nil) != tt.expectError {
				t.Errorf("GenerateThumbnail() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if thumbnail == nil || len(thumbnail) == 0 {
					t.Error("GenerateThumbnail() returned empty thumbnail")
				}
			}
		})
	}
}

func TestImageProcessor_DecodeImage(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	tests := []struct {
		name        string
		data        []byte
		mimeType    string
		expectError bool
	}{
		{
			name:        "Valid JPEG",
			data:        createValidJPEGImage(),
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name:        "Valid PNG",
			data:        createValidPNGImage(),
			mimeType:    "image/png",
			expectError: false,
		},
		{
			name:        "Unknown MIME type but valid PNG",
			data:        createValidPNGImage(),
			mimeType:    "image/unknown",
			expectError: false,
		},
		{
			name:        "Invalid image data",
			data:        []byte("not an image"),
			mimeType:    "image/jpeg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := processor.decodeImage(tt.data, tt.mimeType)
			if (err != nil) != tt.expectError {
				t.Errorf("decodeImage() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && img == nil {
				t.Error("decodeImage() returned nil image")
			}
		})
	}
}

func TestImageProcessor_ResizeImage(t *testing.T) {
	cfg := &ImageConfig{}
	processor := NewImageProcessor(cfg)

	// Decode a valid image first
	img, err := processor.decodeImage(createValidPNGImage(), "image/png")
	if err != nil {
		t.Fatalf("Failed to decode test image: %v", err)
	}

	// Test resizing
	resized := processor.resizeImage(img, 10, 10)
	if resized == nil {
		t.Fatal("resizeImage() returned nil")
	}

	bounds := resized.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("resizeImage() dimensions = %dx%d, want 10x10", bounds.Dx(), bounds.Dy())
	}
}

func TestImageProcessor_ParseThumbnailSize_NonNumericHeight(t *testing.T) {
	cfg := &ImageConfig{
		ThumbnailSize: "200xabc",
	}
	processor := NewImageProcessor(cfg)

	_, err := processor.parseThumbnailSize()
	if err == nil {
		t.Error("parseThumbnailSize() should return error for non-numeric height")
	}
}

func TestImageProcessor_DownloadImageFromURL_EmptyURL(t *testing.T) {
	cfg := &ImageConfig{
		MaxSize:      1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png"},
	}
	processor := NewImageProcessor(cfg)

	_, _, err := processor.DownloadImageFromURL("")
	if err == nil {
		t.Error("DownloadImageFromURL() should return error for empty URL")
	}
}

// Helper functions to create valid image data for actual decoding

func createValidJPEGImage() []byte {
	// Create a minimal valid 1x1 JPEG image using the jpeg encoder
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func createValidPNGImage() []byte {
	// Create a minimal valid 1x1 PNG image using the png encoder
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func createTestWebPData() []byte {
	return []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x1A, 0x00, 0x00, 0x00, // File size
		0x57, 0x45, 0x42, 0x50, // "WEBP"
		0x56, 0x50, 0x38, 0x20, // "VP8 "
		0x0E, 0x00, 0x00, 0x00, // Chunk size
		0x30, 0x01, 0x00, 0x9D, // Frame tag
		0x01, 0x2A, 0x01, 0x00, 0x01, 0x00, // Image dimensions
	}
}

func TestImageProcessor_DownloadImageFromURL_WithMockServer(t *testing.T) {
	// Test successful download
	t.Run("Successful download", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			w.Write(createValidJPEGImage())
		}))
		defer server.Close()

		cfg := &ImageConfig{
			MaxSize:      1024 * 1024,
			AllowedTypes: []string{"image/jpeg", "image/png"},
		}
		processor := NewImageProcessor(cfg)

		data, mimeType, err := processor.DownloadImageFromURL(server.URL)
		if err != nil {
			t.Fatalf("DownloadImageFromURL() error = %v", err)
		}

		if len(data) == 0 {
			t.Error("Downloaded data is empty")
		}

		if mimeType != "image/jpeg" {
			t.Errorf("MimeType = %s, want image/jpeg", mimeType)
		}
	})

	// Test HTTP error response
	t.Run("HTTP error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		cfg := &ImageConfig{
			MaxSize:      1024 * 1024,
			AllowedTypes: []string{"image/jpeg"},
		}
		processor := NewImageProcessor(cfg)

		_, _, err := processor.DownloadImageFromURL(server.URL)
		if err == nil {
			t.Error("DownloadImageFromURL() should return error for HTTP 404")
		}
	})

	// Test image too large
	t.Run("Image size exceeds limit", func(t *testing.T) {
		largeData := make([]byte, 2*1024*1024) // 2MB
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			w.Write(largeData)
		}))
		defer server.Close()

		cfg := &ImageConfig{
			MaxSize:      1024 * 1024, // 1MB limit
			AllowedTypes: []string{"image/jpeg"},
		}
		processor := NewImageProcessor(cfg)

		_, _, err := processor.DownloadImageFromURL(server.URL)
		if err == nil {
			t.Error("DownloadImageFromURL() should return error for oversized image")
		}
	})

	// Test invalid image data
	t.Run("Invalid image data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not an image"))
		}))
		defer server.Close()

		cfg := &ImageConfig{
			MaxSize:      1024 * 1024,
			AllowedTypes: []string{"image/jpeg"},
		}
		processor := NewImageProcessor(cfg)

		_, _, err := processor.DownloadImageFromURL(server.URL)
		if err == nil {
			t.Error("DownloadImageFromURL() should return error for invalid image data")
		}
	})

	// Test MIME type detection when header is missing
	t.Run("MIME type detection from data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Don't set Content-Type header
			w.WriteHeader(http.StatusOK)
			w.Write(createValidJPEGImage())
		}))
		defer server.Close()

		cfg := &ImageConfig{
			MaxSize:      1024 * 1024,
			AllowedTypes: []string{"image/jpeg"},
		}
		processor := NewImageProcessor(cfg)

		_, mimeType, err := processor.DownloadImageFromURL(server.URL)
		if err != nil {
			t.Fatalf("DownloadImageFromURL() error = %v", err)
		}

		if mimeType != "image/jpeg" {
			t.Errorf("Detected MimeType = %s, want image/jpeg", mimeType)
		}
	})

	// Test network error
	t.Run("Network error", func(t *testing.T) {
		cfg := &ImageConfig{
			MaxSize:      1024 * 1024,
			AllowedTypes: []string{"image/jpeg"},
		}
		processor := NewImageProcessor(cfg)

		// Use an invalid URL
		_, _, err := processor.DownloadImageFromURL("http://invalid-url-that-does-not-exist-12345.com")
		if err == nil {
			t.Error("DownloadImageFromURL() should return error for network failure")
		}
	})
}

func BenchmarkImageProcessor_Base64Encoding(b *testing.B) {
	cfg := &ImageConfig{}
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
