a# MCP Image Support Documentation

## Overview
MCP servers can return images and other binary data through resources using base64 encoding.

## Resource Types

### Text Resources
- Plain text content
- JSON data
- Markdown documents
- Source code

### Binary Resources
- **Images** (PNG, JPEG, GIF, WebP, etc.)
- PDFs
- Audio files
- Video files
- Other non-text formats

## Implementation Example

### 1. Resource Definition
```go
// Resource with movie poster image
{
    "uri": "movies://posters/movie-123",
    "name": "The Matrix Poster",
    "description": "Movie poster for The Matrix",
    "mimeType": "image/jpeg"
}
```

### 2. Resource Response Structure
```json
{
    "contents": [
        {
            "uri": "movies://posters/movie-123",
            "mimeType": "image/jpeg",
            "blob": "base64_encoded_image_data_here..."
        }
    ]
}
```

### 3. Go Implementation
```go
type ResourceContent struct {
    URI      string `json:"uri"`
    MimeType string `json:"mimeType,omitempty"`
    Text     string `json:"text,omitempty"`
    Blob     string `json:"blob,omitempty"` // base64 encoded binary data
}

func (s *MoviesServer) handleImageResource(movieID int) (*ResourceContent, error) {
    // Load image from database or filesystem
    imageData, err := s.db.GetMoviePoster(movieID)
    if err != nil {
        return nil, err
    }
    
    // Encode to base64
    encodedImage := base64.StdEncoding.EncodeToString(imageData)
    
    return &ResourceContent{
        URI:      fmt.Sprintf("movies://posters/movie-%d", movieID),
        MimeType: "image/jpeg",
        Blob:     encodedImage,
    }, nil
}
```

## Practical Use Cases for Movies MCP Server

### 1. Movie Posters
- Resource: `movies://posters/{movie-id}`
- Returns movie poster images
- Useful for client applications to display visual content

### 2. Thumbnail Gallery
- Resource: `movies://thumbnails/all`
- Returns a collection of movie thumbnails
- Can include multiple images in one response

### 3. Actor/Director Photos
- Resource: `movies://people/{person-id}/photo`
- Returns photos of cast and crew

## Implementation Considerations

### 1. Storage Options
```go
// Option A: Store in database
type Movie struct {
    ID         int    `db:"id"`
    Title      string `db:"title"`
    PosterData []byte `db:"poster_data"` // Binary image data
    PosterType string `db:"poster_type"` // MIME type
}

// Option B: Store file paths
type Movie struct {
    ID         int    `db:"id"`
    Title      string `db:"title"`
    PosterPath string `db:"poster_path"` // Path to image file
}
```

### 2. Performance Considerations
- Base64 encoding increases data size by ~33%
- Consider caching frequently accessed images
- Implement size limits for images
- Use appropriate image compression

### 3. Example Resource Handler
```go
func (s *MoviesServer) handleResourceRead(uri string) (*ResourceContent, error) {
    parts := strings.Split(uri, "/")
    
    switch parts[2] { // Assuming format: movies://type/id
    case "posters":
        movieID, _ := strconv.Atoi(parts[3])
        return s.handleImageResource(movieID)
    case "database":
        // Handle text/JSON resources
        return s.handleDatabaseResource(parts[3])
    default:
        return nil, fmt.Errorf("unknown resource type")
    }
}
```

## Client Usage

MCP clients can:
1. Request the resource list to discover available images
2. Read specific image resources
3. Decode the base64 blob data
4. Display or save the image

## Size Limitations

Consider implementing:
- Maximum image size (e.g., 5MB)
- Image resizing/compression
- Thumbnail generation for large images
- Lazy loading strategies

## Example Resource Listing
```json
{
    "resources": [
        {
            "uri": "movies://database/all",
            "name": "All Movies",
            "mimeType": "application/json"
        },
        {
            "uri": "movies://posters/collection",
            "name": "Movie Posters Collection",
            "description": "All movie posters in the database",
            "mimeType": "application/json"
        },
        {
            "uri": "movies://posters/1",
            "name": "The Matrix Poster",
            "mimeType": "image/jpeg"
        }
    ]
}
```