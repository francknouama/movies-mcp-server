package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	// Config will need to be passed as parameters or made configurable
)

// ImageProcessor handles image operations for the MCP server
type ImageProcessor struct {
	config *ImageConfig
}

// NewImageProcessor creates a new image processor with the given configuration
func NewImageProcessor(cfg *ImageConfig) *ImageProcessor {
	return &ImageProcessor{
		config: cfg,
	}
}

// ValidateImage checks if the image data is valid and within size limits
func (p *ImageProcessor) ValidateImage(data []byte, mimeType string) error {
	// Check size
	if int64(len(data)) > p.config.MaxSize {
		return fmt.Errorf("image size %d bytes exceeds maximum allowed size %d bytes",
			len(data), p.config.MaxSize)
	}

	if len(data) == 0 {
		return fmt.Errorf("image data cannot be empty")
	}

	// Check MIME type
	if !p.isValidMimeType(mimeType) {
		return fmt.Errorf("unsupported image type: %s. Allowed types: %v",
			mimeType, p.config.AllowedTypes)
	}

	// Validate image format matches MIME type
	if !p.validateImageFormat(data, mimeType) {
		return fmt.Errorf("image data does not match declared MIME type %s", mimeType)
	}

	return nil
}

// isValidMimeType checks if the MIME type is in the allowed list
func (p *ImageProcessor) isValidMimeType(mimeType string) bool {
	for _, allowed := range p.config.AllowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// validateImageFormat checks if the image data matches the declared MIME type
func (p *ImageProcessor) validateImageFormat(data []byte, mimeType string) bool {
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

// EncodeToBase64 encodes image data to base64 string
func (p *ImageProcessor) EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeFromBase64 decodes base64 string to image data
func (p *ImageProcessor) DecodeFromBase64(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image data: %w", err)
	}
	return data, nil
}

// DownloadImageFromURL downloads an image from a URL and returns the data and MIME type
func (p *ImageProcessor) DownloadImageFromURL(url string) ([]byte, string, error) {
	if url == "" {
		return nil, "", fmt.Errorf("URL cannot be empty")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image from %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	// Read response body with size limit
	limitedReader := io.LimitReader(resp.Body, p.config.MaxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Check if size limit was exceeded
	if int64(len(data)) > p.config.MaxSize {
		return nil, "", fmt.Errorf("downloaded image size exceeds maximum allowed size")
	}

	// Determine MIME type from Content-Type header
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		// Try to detect from data
		mimeType = p.detectMimeType(data)
	}

	// Validate the downloaded image
	if err := p.ValidateImage(data, mimeType); err != nil {
		return nil, "", fmt.Errorf("downloaded image validation failed: %w", err)
	}

	return data, mimeType, nil
}

// detectMimeType attempts to detect MIME type from image data
func (p *ImageProcessor) detectMimeType(data []byte) string {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}

	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if len(data) >= len(pngHeader) && bytes.Equal(data[:len(pngHeader)], pngHeader) {
		return "image/png"
	}

	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}

	return "application/octet-stream"
}

// GenerateThumbnail generates a thumbnail of the specified size
func (p *ImageProcessor) GenerateThumbnail(data []byte, mimeType string) ([]byte, error) {
	if !p.config.EnableThumbnails {
		return nil, fmt.Errorf("thumbnail generation is disabled")
	}

	// Parse thumbnail size
	size, err := p.parseThumbnailSize()
	if err != nil {
		return nil, fmt.Errorf("invalid thumbnail size configuration: %w", err)
	}

	// Decode image
	img, err := p.decodeImage(data, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image for thumbnail: %w", err)
	}

	// Resize image (simplified - in real implementation would use image/draw or external library)
	thumbnail := p.resizeImage(img, size.width, size.height)

	// Encode thumbnail as JPEG (for smaller file size)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

// thumbnailSize represents width and height
type thumbnailSize struct {
	width, height int
}

// parseThumbnailSize parses the thumbnail size string (e.g., "200x200")
func (p *ImageProcessor) parseThumbnailSize() (thumbnailSize, error) {
	parts := strings.Split(p.config.ThumbnailSize, "x")
	if len(parts) != 2 {
		return thumbnailSize{}, fmt.Errorf("invalid thumbnail size format: %s", p.config.ThumbnailSize)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return thumbnailSize{}, fmt.Errorf("invalid width: %s", parts[0])
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return thumbnailSize{}, fmt.Errorf("invalid height: %s", parts[1])
	}

	if width <= 0 || height <= 0 {
		return thumbnailSize{}, fmt.Errorf("thumbnail dimensions must be positive")
	}

	return thumbnailSize{width: width, height: height}, nil
}

// decodeImage decodes image data based on MIME type
func (p *ImageProcessor) decodeImage(data []byte, mimeType string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch mimeType {
	case "image/jpeg":
		return jpeg.Decode(reader)
	case "image/png":
		return png.Decode(reader)
	default:
		// Try generic decode
		img, _, err := image.Decode(reader)
		return img, err
	}
}

// resizeImage resizes an image to the specified dimensions (simplified implementation)
func (p *ImageProcessor) resizeImage(img image.Image, width, height int) image.Image {
	// This is a simplified implementation
	// In a production system, you would use a proper image resizing library
	// like "github.com/nfnt/resize" or "golang.org/x/image/draw"

	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple nearest-neighbor scaling
	xRatio := float64(bounds.Dx()) / float64(width)
	yRatio := float64(bounds.Dy()) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			newImg.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return newImg
}

// GetImageInfo returns information about an image
func (p *ImageProcessor) GetImageInfo(data []byte, mimeType string) (*ImageInfo, error) {
	img, err := p.decodeImage(data, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	return &ImageInfo{
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Size:     int64(len(data)),
		MimeType: mimeType,
		Format:   p.getFormatFromMimeType(mimeType),
	}, nil
}

// ImageInfo contains information about an image
type ImageInfo struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	Format   string `json:"format"`
}

// getFormatFromMimeType extracts format from MIME type
func (p *ImageProcessor) getFormatFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return "JPEG"
	case "image/png":
		return "PNG"
	case "image/webp":
		return "WebP"
	default:
		return "Unknown"
	}
}
