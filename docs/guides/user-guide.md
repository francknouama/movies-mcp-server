# 🎬 Complete User Guide to Movies MCP Server

Welcome to the comprehensive user guide for the Movies MCP Server! This guide covers everything you need to know to master your movie database management through AI assistants.

## 📋 Table of Contents

1. [🚀 Getting Started](#-getting-started)
2. [🎭 Movie Management](#-movie-management)
3. [🎪 Actor Management](#-actor-management)
4. [🔍 Search & Discovery](#-search--discovery)
5. [📊 Resources & Analytics](#-resources--analytics)
6. [🛠️ Advanced Features](#-advanced-features)
7. [💡 Pro Tips & Best Practices](#-pro-tips--best-practices)
8. [🔄 Integration Examples](#-integration-examples)

---

## 🚀 Getting Started

### Prerequisites
- ✅ Movies MCP Server installed and running
- ✅ Database connection established
- ✅ MCP client configured (Claude Desktop, direct MCP communication, etc.)

**New here?** → Check out the [Quick Start Guide](../getting-started/README.md)

### Understanding MCP Tools

The Movies MCP Server provides **17 powerful tools** organized into categories:

| Category | Tools | Purpose |
|----------|-------|---------|
| **Movie Management** | `add_movie`, `get_movie`, `update_movie`, `delete_movie`, `list_top_movies` | Core CRUD operations |
| **Actor Management** | `add_actor`, `get_actor`, `update_actor`, `delete_actor`, `search_actors` | People & cast management |
| **Relationships** | `link_actor_to_movie`, `unlink_actor_from_movie`, `get_movie_cast`, `get_actor_movies` | Connect actors to films |
| **Search & Discovery** | `search_movies`, `search_by_decade`, `search_by_rating_range`, `search_similar_movies` | Find and explore content |

---

## 🎭 Movie Management

### Adding Movies

The `add_movie` tool is your primary way to build your collection.

**Basic Movie Addition:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Godfather",
      "director": "Francis Ford Coppola",
      "year": 1972,
      "genres": ["Crime", "Drama"],
      "rating": 9.2,
      "poster_url": "https://example.com/godfather-poster.jpg"
    }
  },
  "id": 1
}
```

**Required Fields:**
- `title` (string) - Movie title
- `director` (string) - Director name  
- `year` (integer) - Release year

**Optional Fields:**
- `genres` (array) - List of genre strings
- `rating` (number, 0-10) - Movie rating
- `poster_url` (string) - URL to movie poster image

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie 'The Godfather' added successfully with ID: 1"
      }
    ]
  },
  "id": 1
}
```

### Retrieving Movies

**Get Movie by ID:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_movie",
    "arguments": {
      "movie_id": 1
    }
  },
  "id": 2
}
```

**Response includes:**
- Complete movie details
- Genre information
- Poster data (if available)
- Cast list (if linked)

### Updating Movies

**Update Movie Details:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "update_movie",
    "arguments": {
      "id": 1,
      "title": "The Godfather",
      "director": "Francis Ford Coppola",
      "year": 1972,
      "rating": 9.3,
      "genres": ["Crime", "Drama", "Thriller"]
    }
  },
  "id": 3
}
```

**🔸 Update Rules:**
- Must include `id` and all required fields
- Partial updates not supported (provide complete data)
- Genres array replaces existing genres entirely

### Deleting Movies

**Remove Movie:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "delete_movie",
    "arguments": {
      "movie_id": 1
    }
  },
  "id": 4
}
```

**⚠️ Warning:** Deletion is permanent and will also remove:
- All actor-movie relationships
- Associated poster images
- Any resource references

### Top Movies

**Get Highest Rated Movies:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_top_movies",
    "arguments": {
      "limit": 10
    }
  },
  "id": 5
}
```

---

## 🎪 Actor Management

### Adding Actors

**Create New Actor:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_actor",
    "arguments": {
      "name": "Al Pacino",
      "birth_year": 1940,
      "bio": "American actor and filmmaker known for his intense and passionate performances."
    }
  },
  "id": 6
}
```

**Required Fields:**
- `name` (string) - Actor's full name
- `birth_year` (integer) - Year of birth

**Optional Fields:**
- `bio` (string) - Actor biography/description

### Searching Actors

**Find Actors by Name:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_actors",
    "arguments": {
      "name": "Al Pacino"
    }
  },
  "id": 7
}
```

**💡 Search Tips:**
- Partial names work: "Al" will find "Al Pacino"
- Case-insensitive search
- Searches both first and last names

### Managing Actor-Movie Relationships

**Link Actor to Movie:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "link_actor_to_movie",
    "arguments": {
      "actor_id": 1,
      "movie_id": 1
    }
  },
  "id": 8
}
```

**Get Movie Cast:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_movie_cast",
    "arguments": {
      "movie_id": 1
    }
  },
  "id": 9
}
```

**Get Actor's Filmography:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor_movies",
    "arguments": {
      "actor_id": 1
    }
  },
  "id": 10
}
```

---

## 🔍 Search & Discovery

### Basic Movie Search

**Multi-Criteria Search:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "title": "Matrix",
      "director": "Wachowski",
      "genre": "Sci-Fi",
      "min_year": 1990,
      "max_year": 2010,
      "min_rating": 8.0,
      "limit": 20
    }
  },
  "id": 11
}
```

**🔍 Search Parameters (all optional):**
- `title` - Search movie titles
- `director` - Find by director name
- `genre` - Filter by genre
- `min_year` / `max_year` - Year range
- `min_rating` / `max_rating` - Rating range
- `limit` - Max results (default: 50)

### Advanced Search Tools

**Search by Decade:**
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
  "id": 12
}
```

**Supported Decades:**
- "1920s", "1930s", "1940s", etc.
- "2000s", "2010s", "2020s"

**Search by Rating Range:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_rating_range",
    "arguments": {
      "min_rating": 8.5,
      "max_rating": 10.0
    }
  },
  "id": 13
}
```

**Find Similar Movies:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_similar_movies",
    "arguments": {
      "movie_id": 1,
      "limit": 5
    }
  },
  "id": 14
}
```

**🎯 Similarity Algorithm:**
- Matches by genre overlap
- Considers director and year proximity
- Accounts for rating similarity

---

## 📊 Resources & Analytics

### Available Resources

The server provides rich data through MCP resources:

| Resource URI | Description | Content Type |
|-------------|-------------|--------------|
| `movies://database/all` | Complete movie database | `application/json` |
| `movies://database/stats` | Database statistics & analytics | `application/json` |
| `movies://posters/collection` | All movie posters | `application/json` |
| `movies://posters/{id}` | Individual movie poster | `image/jpeg` |

### Accessing Resources

**List Available Resources:**
```json
{
  "jsonrpc": "2.0",
  "method": "resources/list",
  "id": 15
}
```

**Read Database Statistics:**
```json
{
  "jsonrpc": "2.0",
  "method": "resources/read",
  "params": {
    "uri": "movies://database/stats"
  },
  "id": 16
}
```

**Example Stats Response:**
```json
{
  "total_movies": 1247,
  "total_actors": 2891,
  "genres": ["Action", "Comedy", "Drama", "Horror", "Sci-Fi"],
  "year_range": {
    "earliest": 1925,
    "latest": 2024
  },
  "top_directors": [
    {"name": "Steven Spielberg", "count": 12},
    {"name": "Christopher Nolan", "count": 8}
  ]
}
```

---

## 🛠️ Advanced Features

### Image Handling

**Automatic Poster Processing:**
- Provide any `poster_url` when adding movies
- Server downloads and stores images automatically
- Images are resized and optimized for storage
- Access via `movies://posters/{movie_id}` resource

**Supported Image Formats:**
- JPEG, PNG, WebP
- Automatic format conversion
- Thumbnail generation

### Database Full-Text Search

**Search Movie Descriptions:**
The search tools perform full-text search across:
- Movie titles
- Director names
- Genre information
- Plot descriptions (when available)

**🔍 Search Enhancement:**
- Fuzzy matching for typos
- Stemming for word variations
- Relevance-based ranking

### Batch Operations

**Adding Multiple Movies:**
```json
// Movie 1
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Pulp Fiction",
      "director": "Quentin Tarantino",
      "year": 1994,
      "genres": ["Crime", "Drama"],
      "rating": 8.9
    }
  },
  "id": 17
}

// Movie 2
{
  "jsonrpc": "2.0", 
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Kill Bill",
      "director": "Quentin Tarantino", 
      "year": 2003,
      "genres": ["Action", "Crime"],
      "rating": 8.2
    }
  },
  "id": 18
}
```

---

## 💡 Pro Tips & Best Practices

### Data Quality

**🎯 Movie Titles:**
- Use official titles, avoid abbreviations
- Include subtitle/part numbers: "Kill Bill: Vol. 1"
- For foreign films, use the most recognized title

**🎭 Director Names:**
- Use full names: "Christopher Nolan", not "C. Nolan"
- For multiple directors: "Anthony Russo, Joe Russo"
- Be consistent with name formatting

**📅 Years:**
- Use theatrical release year, not production year
- For festival films, use wide release year
- Be consistent across your database

**⭐ Ratings:**
- Use a consistent scale (we recommend 0-10)
- Consider using IMDb or Metacritic as reference
- Be consistent with rating sources

### Genre Management

**🏷️ Recommended Genres:**
- Action, Adventure, Animation, Biography, Comedy
- Crime, Documentary, Drama, Family, Fantasy
- History, Horror, Musical, Mystery, Romance
- Sci-Fi, Sport, Thriller, War, Western

**Genre Best Practices:**
- Limit to 2-4 genres per movie
- Use primary genre first in array
- Maintain consistent genre naming

### Search Optimization

**🔍 Effective Search Strategies:**

1. **Start Broad, Then Narrow:**
   ```json
   // First: Find all Sci-Fi
   {"genre": "Sci-Fi"}
   
   // Then: Add year range
   {"genre": "Sci-Fi", "min_year": 2000}
   
   // Finally: Add rating filter
   {"genre": "Sci-Fi", "min_year": 2000, "min_rating": 8.0}
   ```

2. **Use Multiple Search Types:**
   - Title search for specific movies
   - Director search for filmographies
   - Genre + year for era exploration
   - Rating range for quality filtering

3. **Leverage Similar Movies:**
   - Start with a known favorite
   - Use `search_similar_movies` for recommendations
   - Build themed collections

### Performance Optimization

**🚀 Database Performance:**
- Use specific searches vs. broad queries
- Implement reasonable limits (default: 50)
- Cache frequently accessed movie details
- Consider pagination for large result sets

---

## 🔄 Integration Examples

### Claude Desktop Usage

**Natural Language → MCP Tools:**

```
User: "Find me some great sci-fi movies from the 2000s with high ratings"

Claude automatically translates to:
{
  "name": "search_movies",
  "arguments": {
    "genre": "Sci-Fi",
    "min_year": 2000,
    "max_year": 2009,
    "min_rating": 8.0
  }
}
```

**Building Collections:**
```
User: "Add all the Christopher Nolan movies to my database"

Claude uses multiple tool calls:
1. search_movies(director="Christopher Nolan") - check existing
2. add_movie(...) - for each missing movie
3. search_actors(name="Christopher Nolan") - if not exists
4. add_actor(...) - add director as actor entry
```

### API Integration Patterns

**Movie Recommendation Workflow:**
```javascript
// 1. Get user's top-rated movies
const topMovies = await mcpCall('list_top_movies', {limit: 5});

// 2. Find similar movies for each
const recommendations = [];
for (const movie of topMovies) {
  const similar = await mcpCall('search_similar_movies', {
    movie_id: movie.id,
    limit: 3
  });
  recommendations.push(...similar);
}

// 3. Remove duplicates and user's existing movies
const uniqueRecs = filterUniqueRecommendations(recommendations);
```

**Collection Building:**
```javascript
// Import from external API (e.g., IMDb)
async function importMovieCollection(imdbIds) {
  for (const imdbId of imdbIds) {
    const movieData = await fetchFromIMDb(imdbId);
    
    // Add movie
    const movieResult = await mcpCall('add_movie', {
      title: movieData.title,
      director: movieData.director,
      year: movieData.year,
      genres: movieData.genres,
      rating: movieData.rating,
      poster_url: movieData.poster
    });
    
    // Add actors
    for (const actor of movieData.cast) {
      const actorResult = await mcpCall('add_actor', {
        name: actor.name,
        birth_year: actor.birthYear,
        bio: actor.bio
      });
      
      // Link actor to movie
      await mcpCall('link_actor_to_movie', {
        actor_id: actorResult.id,
        movie_id: movieResult.id
      });
    }
  }
}
```

---

## 🎊 Congratulations!

You've mastered the Movies MCP Server! You now know how to:

✅ **Manage Movies** - Add, update, search, and organize your collection  
✅ **Handle Actors** - Build comprehensive cast databases  
✅ **Advanced Search** - Find exactly what you're looking for  
✅ **Use Resources** - Access analytics and media content  
✅ **Best Practices** - Maintain a high-quality database  
✅ **Integration** - Connect with AI assistants and applications

## 🔗 What's Next?

- **[Examples Guide](./examples.md)** - See real-world use cases and workflows
- **[API Reference](../reference/api.md)** - Complete technical documentation
- **[Troubleshooting](../reference/troubleshooting.md)** - Fix common issues
- **[Development Guide](../development/README.md)** - Extend and customize the server

---

*💡 **Need help?** Check our [Troubleshooting Guide](../reference/troubleshooting.md) or [GitHub Issues](https://github.com/francknouama/movies-mcp-server/issues)*