# 🔧 Complete MCP Tools API Reference

Comprehensive technical documentation for all Movies MCP Server tools, parameters, responses, and resources.

## 📋 Table of Contents

1. [🚀 Protocol Overview](#-protocol-overview)
2. [🎬 Movie Management Tools](#-movie-management-tools)
3. [🎭 Actor Management Tools](#-actor-management-tools)
4. [🔗 Relationship Management Tools](#-relationship-management-tools)
5. [🔍 Search & Discovery Tools](#-search--discovery-tools)
6. [📊 Resource Endpoints](#-resource-endpoints)
7. [🎯 Quick Reference](#-quick-reference)
8. [🛠️ Error Handling](#-error-handling)

---

## 🚀 Protocol Overview

### JSON-RPC 2.0 Specification

The Movies MCP Server implements the MCP (Model Context Protocol) over JSON-RPC 2.0.

**Standard Request Format:**
```json
{
  "jsonrpc": "2.0",
  "method": "method_name",
  "params": { ... },
  "id": 1
}
```

**Standard Response Format:**
```json
{
  "jsonrpc": "2.0",
  "result": { ... },
  "id": 1
}
```

**Error Response Format:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Error description"
  },
  "id": 1
}
```

### MCP Protocol Methods

| Method | Purpose | Parameters |
|--------|---------|------------|
| `initialize` | Establish connection | `protocolVersion`, `capabilities` |
| `tools/list` | Get available tools | None |
| `tools/call` | Execute a tool | `name`, `arguments` |
| `resources/list` | Get available resources | None |
| `resources/read` | Read resource content | `uri` |

### Server Information

**Server Details:**
- **Name:** `movies-mcp-server`
- **Version:** `0.2.0`
- **Protocol Version:** `2024-11-05`
- **Capabilities:** `tools`, `resources`

---

## 🎬 Movie Management Tools

### `add_movie`

Add a new movie to the database.

**Parameters:**
| Parameter | Type | Required | Description | Constraints |
|-----------|------|----------|-------------|-------------|
| `title` | string | ✅ | Movie title | Max 255 chars |
| `director` | string | ✅ | Director name | Max 255 chars |
| `year` | integer | ✅ | Release year | 1888-2030 |
| `genres` | array[string] | ❌ | List of genres | Max 10 genres |
| `rating` | number | ❌ | Movie rating | 0.0-10.0 |
| `poster_url` | string | ❌ | Poster image URL | Valid HTTP/HTTPS URL |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Matrix",
      "director": "The Wachowskis",
      "year": 1999,
      "genres": ["Action", "Sci-Fi"],
      "rating": 8.7,
      "poster_url": "https://example.com/matrix-poster.jpg"
    }
  },
  "id": 1
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie 'The Matrix' added successfully with ID: 42"
      }
    ]
  },
  "id": 1
}
```

**Error Cases:**
- **Duplicate Movie:** Returns `-32602` if movie with same title, director, and year exists
- **Invalid Rating:** Returns `-32602` if rating outside 0.0-10.0 range
- **Invalid Year:** Returns `-32602` if year outside valid range
- **Poster Download Failed:** Returns `-32603` if poster URL is inaccessible

---

### `get_movie`

Retrieve movie details by ID.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `movie_id` | integer | ✅ | Movie ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_movie",
    "arguments": {
      "movie_id": 42
    }
  },
  "id": 2
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie Details:\nID: 42\nTitle: The Matrix\nDirector: The Wachowskis\nYear: 1999\nGenres: Action, Sci-Fi\nRating: 8.7\nPoster: Available"
      }
    ]
  },
  "id": 2
}
```

**Error Cases:**
- **Movie Not Found:** Returns `-32602` if movie ID doesn't exist

---

### `update_movie`

Update an existing movie's information.

**Parameters:**
| Parameter | Type | Required | Description | Constraints |
|-----------|------|----------|-------------|-------------|
| `id` | integer | ✅ | Movie ID | Must exist |
| `title` | string | ✅ | Movie title | Max 255 chars |
| `director` | string | ✅ | Director name | Max 255 chars |
| `year` | integer | ✅ | Release year | 1888-2030 |
| `genres` | array[string] | ❌ | List of genres | Max 10 genres |
| `rating` | number | ❌ | Movie rating | 0.0-10.0 |
| `poster_url` | string | ❌ | Poster image URL | Valid HTTP/HTTPS URL |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "update_movie",
    "arguments": {
      "id": 42,
      "title": "The Matrix",
      "director": "The Wachowskis",
      "year": 1999,
      "genres": ["Action", "Sci-Fi", "Thriller"],
      "rating": 8.8
    }
  },
  "id": 3
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie 'The Matrix' updated successfully"
      }
    ]
  },
  "id": 3
}
```

**Error Cases:**
- **Movie Not Found:** Returns `-32602` if movie ID doesn't exist
- **Validation Errors:** Same as `add_movie`

---

### `delete_movie`

Remove a movie from the database.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `movie_id` | integer | ✅ | Movie ID to delete |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "delete_movie",
    "arguments": {
      "movie_id": 42
    }
  },
  "id": 4
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie deleted successfully"
      }
    ]
  },
  "id": 4
}
```

**Error Cases:**
- **Movie Not Found:** Returns `-32602` if movie ID doesn't exist

**Side Effects:**
- Removes all actor-movie relationships
- Deletes associated poster images
- Removes movie from resource endpoints

---

### `list_top_movies`

Get highest-rated movies.

**Parameters:**
| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `limit` | integer | ❌ | Number of movies to return | 10 |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_top_movies",
    "arguments": {
      "limit": 5
    }
  },
  "id": 5
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Top 5 Movies:\n1. The Shawshank Redemption (9.3) - Drama\n2. The Godfather (9.2) - Crime, Drama\n3. The Dark Knight (9.0) - Action, Crime, Drama\n4. Pulp Fiction (8.9) - Crime, Drama\n5. The Matrix (8.8) - Action, Sci-Fi, Thriller"
      }
    ]
  },
  "id": 5
}
```

---

## 🎭 Actor Management Tools

### `add_actor`

Add a new actor to the database.

**Parameters:**
| Parameter | Type | Required | Description | Constraints |
|-----------|------|----------|-------------|-------------|
| `name` | string | ✅ | Actor's full name | Max 255 chars |
| `birth_year` | integer | ✅ | Year of birth | 1800-2030 |
| `bio` | string | ❌ | Actor biography | Max 2000 chars |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_actor",
    "arguments": {
      "name": "Keanu Reeves",
      "birth_year": 1964,
      "bio": "Canadian actor known for his roles in action films and his humble personality."
    }
  },
  "id": 6
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Actor 'Keanu Reeves' added successfully with ID: 15"
      }
    ]
  },
  "id": 6
}
```

**Error Cases:**
- **Duplicate Actor:** Returns `-32602` if actor with same name and birth year exists
- **Invalid Birth Year:** Returns `-32602` if birth year outside valid range

---

### `get_actor`

Retrieve actor details by ID.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actor_id` | integer | ✅ | Actor ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor",
    "arguments": {
      "actor_id": 15
    }
  },
  "id": 7
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Actor Details:\nID: 15\nName: Keanu Reeves\nBirth Year: 1964\nBio: Canadian actor known for his roles in action films and his humble personality."
      }
    ]
  },
  "id": 7
}
```

---

### `update_actor`

Update an existing actor's information.

**Parameters:**
| Parameter | Type | Required | Description | Constraints |
|-----------|------|----------|-------------|-------------|
| `id` | integer | ✅ | Actor ID | Must exist |
| `name` | string | ✅ | Actor's full name | Max 255 chars |
| `birth_year` | integer | ✅ | Year of birth | 1800-2030 |
| `bio` | string | ❌ | Actor biography | Max 2000 chars |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "update_actor",
    "arguments": {
      "id": 15,
      "name": "Keanu Reeves",
      "birth_year": 1964,
      "bio": "Canadian actor, director, and producer known for his roles in The Matrix trilogy, John Wick series, and Speed."
    }
  },
  "id": 8
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Actor 'Keanu Reeves' updated successfully"
      }
    ]
  },
  "id": 8
}
```

---

### `delete_actor`

Remove an actor from the database.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actor_id` | integer | ✅ | Actor ID to delete |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "delete_actor",
    "arguments": {
      "actor_id": 15
    }
  },
  "id": 9
}
```

**Side Effects:**
- Removes all actor-movie relationships

---

### `search_actors`

Search for actors by name.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | ✅ | Actor name to search |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_actors",
    "arguments": {
      "name": "Keanu"
    }
  },
  "id": 10
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Found 1 actor(s):\n1. Keanu Reeves (1964) - Canadian actor, director, and producer..."
      }
    ]
  },
  "id": 10
}
```

---

## 🔗 Relationship Management Tools

### `link_actor_to_movie`

Create a relationship between an actor and a movie.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actor_id` | integer | ✅ | Actor ID |
| `movie_id` | integer | ✅ | Movie ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "link_actor_to_movie",
    "arguments": {
      "actor_id": 15,
      "movie_id": 42
    }
  },
  "id": 11
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Actor 'Keanu Reeves' linked to movie 'The Matrix' successfully"
      }
    ]
  },
  "id": 11
}
```

**Error Cases:**
- **Actor Not Found:** Returns `-32602` if actor ID doesn't exist
- **Movie Not Found:** Returns `-32602` if movie ID doesn't exist
- **Already Linked:** Returns `-32602` if relationship already exists

---

### `unlink_actor_from_movie`

Remove a relationship between an actor and a movie.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actor_id` | integer | ✅ | Actor ID |
| `movie_id` | integer | ✅ | Movie ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "unlink_actor_from_movie",
    "arguments": {
      "actor_id": 15,
      "movie_id": 42
    }
  },
  "id": 12
}
```

---

### `get_movie_cast`

Get all actors in a specific movie.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `movie_id` | integer | ✅ | Movie ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_movie_cast",
    "arguments": {
      "movie_id": 42
    }
  },
  "id": 13
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Cast of 'The Matrix':\n1. Keanu Reeves (1964)\n2. Laurence Fishburne (1961)\n3. Carrie-Anne Moss (1967)"
      }
    ]
  },
  "id": 13
}
```

---

### `get_actor_movies`

Get all movies an actor has appeared in.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actor_id` | integer | ✅ | Actor ID |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor_movies",
    "arguments": {
      "actor_id": 15
    }
  },
  "id": 14
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movies featuring 'Keanu Reeves':\n1. The Matrix (1999) - Action, Sci-Fi\n2. John Wick (2014) - Action, Crime\n3. Speed (1994) - Action, Thriller"
      }
    ]
  },
  "id": 14
}
```

---

## 🔍 Search & Discovery Tools

### `search_movies`

Search movies by multiple criteria.

**Parameters:**
| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `title` | string | ❌ | Search by title | - |
| `director` | string | ❌ | Search by director | - |
| `genre` | string | ❌ | Filter by genre | - |
| `min_year` | integer | ❌ | Minimum release year | - |
| `max_year` | integer | ❌ | Maximum release year | - |
| `min_rating` | number | ❌ | Minimum rating | - |
| `max_rating` | number | ❌ | Maximum rating | - |
| `limit` | integer | ❌ | Maximum results | 50 |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Sci-Fi",
      "min_year": 1990,
      "max_year": 2010,
      "min_rating": 8.0,
      "limit": 10
    }
  },
  "id": 15
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Found 3 movies:\n1. The Matrix (1999) - Action, Sci-Fi - 8.7\n2. Blade Runner 2049 (2017) - Sci-Fi - 8.0\n3. Minority Report (2002) - Action, Sci-Fi - 8.1"
      }
    ]
  },
  "id": 15
}
```

**Search Features:**
- **Full-text search** on titles and descriptions
- **Fuzzy matching** for typos and variations
- **Multiple criteria** can be combined
- **Relevance ranking** for text searches

---

### `search_by_decade`

Search movies by decade.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `decade` | string | ✅ | Decade to search |

**Valid Decades:**
- `"1920s"`, `"1930s"`, `"1940s"`, `"1950s"`
- `"1960s"`, `"1970s"`, `"1980s"`, `"1990s"`
- `"2000s"`, `"2010s"`, `"2020s"`

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_decade",
    "arguments": {
      "decade": "1990s"
    }
  },
  "id": 16
}
```

---

### `search_by_rating_range`

Search movies within a specific rating range.

**Parameters:**
| Parameter | Type | Required | Description | Constraints |
|-----------|------|----------|-------------|-------------|
| `min_rating` | number | ✅ | Minimum rating | 0.0-10.0 |
| `max_rating` | number | ✅ | Maximum rating | 0.0-10.0 |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_rating_range",
    "arguments": {
      "min_rating": 9.0,
      "max_rating": 10.0
    }
  },
  "id": 17
}
```

---

### `search_similar_movies`

Find movies similar to a given movie.

**Parameters:**
| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `movie_id` | integer | ✅ | Reference movie ID | - |
| `limit` | integer | ❌ | Number of results | 5 |

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_similar_movies",
    "arguments": {
      "movie_id": 42,
      "limit": 3
    }
  },
  "id": 18
}
```

**Similarity Algorithm:**
- **Genre matching** (weighted by number of shared genres)
- **Director matching** (same director gets high similarity)
- **Year proximity** (movies from similar time periods)
- **Rating similarity** (movies with similar ratings)

---

## 📊 Resource Endpoints

### Available Resources

| Resource URI | Description | Content Type |
|-------------|-------------|--------------|
| `movies://database/all` | Complete movie database | `application/json` |
| `movies://database/stats` | Database statistics | `application/json` |
| `movies://posters/collection` | Movie poster collection | `application/json` |
| `movies://posters/{id}` | Individual movie poster | `image/jpeg` |

### `movies://database/all`

Complete movie database in JSON format.

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "resources/read",
  "params": {
    "uri": "movies://database/all"
  },
  "id": 19
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "contents": [
      {
        "uri": "movies://database/all",
        "mimeType": "application/json",
        "text": "{\"movies\":[{\"id\":42,\"title\":\"The Matrix\",\"director\":\"The Wachowskis\",\"year\":1999,\"genres\":[\"Action\",\"Sci-Fi\"],\"rating\":8.7}]}"
      }
    ]
  },
  "id": 19
}
```

### `movies://database/stats`

Database statistics and analytics.

**Response Structure:**
```json
{
  "total_movies": 150,
  "total_actors": 420,
  "genres": ["Action", "Drama", "Comedy", "Sci-Fi", "Horror"],
  "year_range": {
    "earliest": 1927,
    "latest": 2024
  },
  "rating_distribution": {
    "0-5": 10,
    "5-7": 45,
    "7-8": 65,
    "8-9": 25,
    "9-10": 5
  },
  "top_directors": [
    {"name": "Christopher Nolan", "count": 8},
    {"name": "Steven Spielberg", "count": 12}
  ],
  "movies_by_decade": {
    "1990s": 25,
    "2000s": 40,
    "2010s": 55,
    "2020s": 30
  }
}
```

### `movies://posters/collection`

Collection of all movie posters.

**Response Structure:**
```json
{
  "posters": [
    {
      "movie_id": 42,
      "title": "The Matrix",
      "poster_uri": "movies://posters/42",
      "format": "image/jpeg",
      "size": 245760
    }
  ],
  "total": 1
}
```

### `movies://posters/{id}`

Individual movie poster image.

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "resources/read",
  "params": {
    "uri": "movies://posters/42"
  },
  "id": 20
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "contents": [
      {
        "uri": "movies://posters/42",
        "mimeType": "image/jpeg",
        "blob": "base64-encoded-image-data"
      }
    ]
  },
  "id": 20
}
```

---

## 🎯 Quick Reference

### Tool Categories

**Movie Operations:**
```bash
add_movie      # Create new movie
get_movie      # Retrieve movie by ID
update_movie   # Modify movie details
delete_movie   # Remove movie
list_top_movies # Get highest rated
```

**Actor Operations:**
```bash
add_actor      # Create new actor
get_actor      # Retrieve actor by ID
update_actor   # Modify actor details
delete_actor   # Remove actor
search_actors  # Find actors by name
```

**Relationships:**
```bash
link_actor_to_movie    # Connect actor to movie
unlink_actor_from_movie # Disconnect actor from movie
get_movie_cast         # Get movie's actors
get_actor_movies       # Get actor's filmography
```

**Search & Discovery:**
```bash
search_movies          # Multi-criteria search
search_by_decade       # Find movies by decade
search_by_rating_range # Find movies by rating
search_similar_movies  # Find similar movies
```

### Common Parameter Patterns

**ID Parameters:**
```json
{"movie_id": 42}
{"actor_id": 15}
```

**Search Parameters:**
```json
{
  "title": "Matrix",
  "director": "Nolan",
  "genre": "Sci-Fi",
  "min_year": 2000,
  "max_year": 2010,
  "min_rating": 8.0,
  "limit": 10
}
```

**Movie Data:**
```json
{
  "title": "Movie Title",
  "director": "Director Name",
  "year": 2024,
  "genres": ["Genre1", "Genre2"],
  "rating": 8.5,
  "poster_url": "https://example.com/poster.jpg"
}
```

### Response Patterns

**Success Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Operation successful message"
      }
    ]
  },
  "id": 1
}
```

**List Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Item 1\nItem 2\nItem 3"
      }
    ]
  },
  "id": 1
}
```

---

## 🛠️ Error Handling

### Standard JSON-RPC Errors

| Code | Name | Description |
|------|------|-------------|
| `-32700` | Parse Error | Invalid JSON |
| `-32600` | Invalid Request | Malformed JSON-RPC |
| `-32601` | Method Not Found | Unknown method |
| `-32602` | Invalid Params | Parameter validation failed |
| `-32603` | Internal Error | Server error |

### Application-Specific Errors

**Movie Errors:**
```json
{
  "code": -32602,
  "message": "Movie not found with ID: 42"
}

{
  "code": -32602,
  "message": "Movie already exists: The Matrix (1999) by The Wachowskis"
}

{
  "code": -32602,
  "message": "Rating must be between 0 and 10, got: 11"
}
```

**Actor Errors:**
```json
{
  "code": -32602,
  "message": "Actor not found with ID: 15"
}

{
  "code": -32602,
  "message": "Actor already exists: Keanu Reeves (1964)"
}
```

**Relationship Errors:**
```json
{
  "code": -32602,
  "message": "Actor is already linked to this movie"
}

{
  "code": -32602,
  "message": "Relationship not found"
}
```

**Database Errors:**
```json
{
  "code": -32603,
  "message": "Database connection failed"
}

{
  "code": -32603,
  "message": "Database query failed: constraint violation"
}
```

**Image Errors:**
```json
{
  "code": -32603,
  "message": "Failed to download poster from URL: https://example.com/poster.jpg"
}

{
  "code": -32602,
  "message": "Invalid image format. Supported formats: JPEG, PNG, WebP"
}
```

### Error Response Format

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Detailed error message",
    "data": {
      "field": "parameter_name",
      "received": "invalid_value",
      "expected": "valid_format"
    }
  },
  "id": 1
}
```

### Handling Errors in Client Code

**JavaScript Example:**
```javascript
async function callTool(toolName, args) {
  const response = await mcpClient.call(toolName, args);
  
  if (response.error) {
    switch (response.error.code) {
      case -32602:
        console.error('Validation error:', response.error.message);
        break;
      case -32603:
        console.error('Server error:', response.error.message);
        break;
      default:
        console.error('Unexpected error:', response.error);
    }
    throw new Error(response.error.message);
  }
  
  return response.result;
}
```

**Python Example:**
```python
def handle_mcp_response(response):
    if 'error' in response:
        error_code = response['error']['code']
        error_message = response['error']['message']
        
        if error_code == -32602:
            raise ValueError(f"Validation error: {error_message}")
        elif error_code == -32603:
            raise RuntimeError(f"Server error: {error_message}")
        else:
            raise Exception(f"MCP error {error_code}: {error_message}")
    
    return response['result']
```

---

## 🔗 Related Documentation

- **[User Guide](../guides/user-guide.md)** - Complete feature walkthrough
- **[Examples](../guides/examples.md)** - Real-world usage scenarios
- **[Troubleshooting](./troubleshooting.md)** - Common issues and solutions
- **[Getting Started](../getting-started/README.md)** - Quick setup guide

---

*⚡ **Pro Tip:** Use the `tools/list` method to discover all available tools and their current schemas. The server provides complete OpenAPI-style documentation for each tool.*