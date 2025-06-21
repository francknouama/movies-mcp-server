package image

// ImageConfig contains configuration for image processing
type ImageConfig struct {
	MaxSize          int64    `json:"max_size"`
	AllowedTypes     []string `json:"allowed_types"`
	EnableThumbnails bool     `json:"enable_thumbnails"`
	ThumbnailSize    string   `json:"thumbnail_size"`
}

// DefaultImageConfig returns default image configuration
func DefaultImageConfig() *ImageConfig {
	return &ImageConfig{
		MaxSize:          10 * 1024 * 1024, // 10MB
		AllowedTypes:     []string{"image/jpeg", "image/png", "image/webp"},
		EnableThumbnails: true,
		ThumbnailSize:    "200x200",
	}
}