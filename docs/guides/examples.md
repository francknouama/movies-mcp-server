# 🎬 Real-World Examples & Use Cases

Discover practical ways to use the Movies MCP Server through detailed examples and common workflows. Each example includes complete JSON-RPC requests and expected responses.

## 📋 Table of Contents

1. [🏠 Personal Movie Collection](#-personal-movie-collection)
2. [🎥 Film Study & Research](#-film-study--research)
3. [🤖 AI-Powered Movie Recommendations](#-ai-powered-movie-recommendations)
4. [📊 Movie Database Analytics](#-movie-database-analytics)
5. [🎭 Cast & Crew Management](#-cast--crew-management)
6. [🔍 Advanced Search Scenarios](#-advanced-search-scenarios)
7. [📱 Application Integration](#-application-integration)
8. [🎪 Event & Collection Planning](#-event--collection-planning)

---

## 🏠 Personal Movie Collection

### Scenario: Building Your Personal Movie Library

**Goal:** Create and manage a personal movie collection with ratings, organize by genres, and track what you've watched.

#### Step 1: Add Your Favorite Movies

```json
// Add The Shawshank Redemption
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Shawshank Redemption",
      "director": "Frank Darabont",
      "year": 1994,
      "genres": ["Drama"],
      "rating": 9.3,
      "poster_url": "https://example.com/shawshank.jpg"
    }
  },
  "id": 1
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text", 
        "text": "Movie 'The Shawshank Redemption' added successfully with ID: 1"
      }
    ]
  },
  "id": 1
}
```

#### Step 2: Add Multiple Movies Quickly

```json
// Add Forrest Gump
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Forrest Gump",
      "director": "Robert Zemeckis",
      "year": 1994,
      "genres": ["Drama", "Romance"],
      "rating": 8.8
    }
  },
  "id": 2
}

// Add Pulp Fiction
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
  "id": 3
}
```

#### Step 3: Find Your Top-Rated Movies

```json
// Get your top 10 movies
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_top_movies",
    "arguments": {
      "limit": 10
    }
  },
  "id": 4
}
```

#### Step 4: Organize by Genre

```json
// Find all your drama movies
{
  "jsonrpc": "2.0", 
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Drama",
      "limit": 20
    }
  },
  "id": 5
}
```

**💡 Pro Tip:** Use Claude Desktop integration for natural language:
```
"Show me all my drama movies from the 1990s rated above 8.5"
```

---

## 🎥 Film Study & Research

### Scenario: Analyzing Christopher Nolan's Filmography

**Goal:** Create a comprehensive database of Christopher Nolan films for academic research.

#### Step 1: Add Christopher Nolan Movies

```json
// Add Inception
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Inception",
      "director": "Christopher Nolan", 
      "year": 2010,
      "genres": ["Action", "Sci-Fi", "Thriller"],
      "rating": 8.8,
      "poster_url": "https://example.com/inception.jpg"
    }
  },
  "id": 6
}

// Add The Dark Knight
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Dark Knight", 
      "director": "Christopher Nolan",
      "year": 2008,
      "genres": ["Action", "Crime", "Drama"],
      "rating": 9.0
    }
  },
  "id": 7
}

// Add Interstellar
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Interstellar",
      "director": "Christopher Nolan",
      "year": 2014,
      "genres": ["Adventure", "Drama", "Sci-Fi"],
      "rating": 8.6
    }
  },
  "id": 8
}
```

#### Step 2: Research Nolan's Complete Filmography

```json
// Find all Christopher Nolan movies
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "director": "Christopher Nolan"
    }
  },
  "id": 9
}
```

#### Step 3: Analyze Genre Evolution

```json
// Find Nolan's sci-fi films
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies", 
    "arguments": {
      "director": "Christopher Nolan",
      "genre": "Sci-Fi"
    }
  },
  "id": 10
}

// Find films by decade
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_decade",
    "arguments": {
      "decade": "2000s"
    }
  },
  "id": 11
}
```

#### Step 4: Add Key Actors

```json
// Add Leonardo DiCaprio 
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_actor",
    "arguments": {
      "name": "Leonardo DiCaprio",
      "birth_year": 1974,
      "bio": "American actor and producer known for his work in biographical and period films."
    }
  },
  "id": 12
}

// Link to Inception (assuming movie_id: 1, actor_id: 1)
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
  "id": 13
}
```

---

## 🤖 AI-Powered Movie Recommendations

### Scenario: Building a Smart Recommendation System

**Goal:** Use the MCP server with AI to create personalized movie recommendations.

#### Step 1: Identify User Preferences

```json
// Get user's highest-rated movies
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_top_movies",
    "arguments": {
      "limit": 5
    }
  },
  "id": 14
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Top Movies:\n1. The Shawshank Redemption (9.3) - Drama\n2. The Dark Knight (9.0) - Action, Crime, Drama\n3. Pulp Fiction (8.9) - Crime, Drama\n4. Inception (8.8) - Action, Sci-Fi, Thriller\n5. Forrest Gump (8.8) - Drama, Romance"
      }
    ]
  },
  "id": 14
}
```

#### Step 2: Find Similar Movies

```json
// Find movies similar to The Dark Knight (ID: 2)
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_similar_movies",
    "arguments": {
      "movie_id": 2,
      "limit": 5
    }
  },
  "id": 15
}
```

#### Step 3: Genre-Based Recommendations

```json
// Find high-rated crime dramas (user's apparent preference)
{
  "jsonrpc": "2.0", 
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Crime",
      "min_rating": 8.0,
      "limit": 10
    }
  },
  "id": 16
}
```

#### Step 4: Era-Based Discovery

```json
// Explore 1990s movies (user likes several from this era)
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "search_by_decade",
    "arguments": {
      "decade": "1990s"
    }
  },
  "id": 17
}
```

**🤖 Claude Desktop Example:**
```
User: "Based on my movie ratings, what should I watch next?"

Claude: "I can see you love highly-rated dramas and crime films from the 1990s and 2000s. Based on your top movies, here are some recommendations..."

[Claude automatically uses the MCP tools to analyze patterns and suggest movies]
```

---

## 📊 Movie Database Analytics

### Scenario: Database Statistics and Insights

**Goal:** Generate comprehensive analytics about your movie collection.

#### Step 1: Overall Database Statistics

```json
// Get database statistics
{
  "jsonrpc": "2.0",
  "method": "resources/read",
  "params": {
    "uri": "movies://database/stats"
  },
  "id": 18
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "contents": [
      {
        "uri": "movies://database/stats",
        "mimeType": "application/json",
        "text": "{\"total_movies\": 142, \"total_actors\": 387, \"genres\": [\"Action\", \"Drama\", \"Comedy\", \"Sci-Fi\"], \"year_range\": {\"earliest\": 1939, \"latest\": 2024}}"
      }
    ]
  },
  "id": 18
}
```

#### Step 2: Genre Distribution Analysis

```json
// Find all action movies
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Action"
    }
  },
  "id": 19
}

// Find all dramas
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Drama"
    }
  },
  "id": 20
}
```

#### Step 3: Rating Distribution

```json
// Movies rated 9.0 and above
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
  "id": 21
}

// Movies rated 8.0-8.9
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_rating_range", 
    "arguments": {
      "min_rating": 8.0,
      "max_rating": 8.9
    }
  },
  "id": 22
}
```

#### Step 4: Decade Analysis

```json
// Movies from each decade
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_by_decade",
    "arguments": {
      "decade": "2010s"
    }
  },
  "id": 23
}
```

---

## 🎭 Cast & Crew Management

### Scenario: Managing Actor Relationships and Filmographies

**Goal:** Build comprehensive actor profiles and track their collaborations.

#### Step 1: Add Major Actors

```json
// Add Robert De Niro
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_actor",
    "arguments": {
      "name": "Robert De Niro", 
      "birth_year": 1943,
      "bio": "American actor and director known for his method acting and collaborations with Martin Scorsese."
    }
  },
  "id": 24
}

// Add Al Pacino
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_actor",
    "arguments": {
      "name": "Al Pacino",
      "birth_year": 1940, 
      "bio": "American actor known for his intense performances in crime dramas."
    }
  },
  "id": 25
}
```

#### Step 2: Link Actors to Movies

```json
// Add The Godfather with cast
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
      "rating": 9.2
    }
  },
  "id": 26
}

// Link Al Pacino to The Godfather (assuming IDs)
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "link_actor_to_movie",
    "arguments": {
      "actor_id": 2,  // Al Pacino
      "movie_id": 5   // The Godfather  
    }
  },
  "id": 27
}
```

#### Step 3: Explore Actor Filmographies

```json
// Get Al Pacino's complete filmography
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor_movies",
    "arguments": {
      "actor_id": 2
    }
  },
  "id": 28
}

// Get The Godfather's complete cast
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_movie_cast",
    "arguments": {
      "movie_id": 5
    }
  },
  "id": 29
}
```

#### Step 4: Find Actor Collaborations

```json
// Search for movies with both De Niro and Pacino
// First get De Niro's movies
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor_movies", 
    "arguments": {
      "actor_id": 1  // Robert De Niro
    }
  },
  "id": 30
}

// Then get Pacino's movies  
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_actor_movies",
    "arguments": {
      "actor_id": 2  // Al Pacino
    }
  },
  "id": 31
}
```

---

## 🔍 Advanced Search Scenarios

### Scenario: Complex Movie Discovery Workflows

#### Finding Hidden Gems

```json
// High-rated movies from specific years with limited popularity
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "min_year": 2015,
      "max_year": 2020,
      "min_rating": 8.5,
      "limit": 20
    }
  },
  "id": 32
}
```

#### Director Deep Dives

```json
// Explore lesser-known Kubrick films
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "director": "Stanley Kubrick",
      "max_rating": 8.5  // Exclude the super famous ones
    }
  },
  "id": 33
}
```

#### Genre Evolution Studies

```json
// Horror movies by decade
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Horror",
      "min_year": 1970,
      "max_year": 1979
    }
  },
  "id": 34
}

{
  "jsonrpc": "2.0", 
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Horror",
      "min_year": 2010,
      "max_year": 2019
    }
  },
  "id": 35
}
```

---

## 📱 Application Integration

### Scenario: Building a Movie Recommendation App

#### JavaScript Integration Example

```javascript
class MovieRecommendationEngine {
  constructor(mcpClient) {
    this.mcp = mcpClient;
  }

  async getUserProfile(userId) {
    // Get user's top movies to understand preferences
    const topMovies = await this.mcp.call('list_top_movies', {limit: 10});
    
    // Analyze genre preferences
    const genres = this.extractGenres(topMovies);
    const avgRating = this.calculateAverageRating(topMovies);
    
    return {
      favoriteGenres: genres,
      averageRating: avgRating,
      topMovies: topMovies
    };
  }

  async generateRecommendations(userProfile) {
    const recommendations = [];
    
    // Find similar movies to user's top picks
    for (const movie of userProfile.topMovies.slice(0, 3)) {
      const similar = await this.mcp.call('search_similar_movies', {
        movie_id: movie.id,
        limit: 3
      });
      recommendations.push(...similar);
    }

    // Add genre-based recommendations
    for (const genre of userProfile.favoriteGenres) {
      const genreMovies = await this.mcp.call('search_movies', {
        genre: genre,
        min_rating: userProfile.averageRating - 0.5,
        limit: 5
      });
      recommendations.push(...genreMovies);
    }

    return this.deduplicateRecommendations(recommendations);
  }

  extractGenres(movies) {
    const genreCount = {};
    movies.forEach(movie => {
      movie.genres.forEach(genre => {
        genreCount[genre] = (genreCount[genre] || 0) + 1;
      });
    });
    
    return Object.keys(genreCount)
      .sort((a, b) => genreCount[b] - genreCount[a])
      .slice(0, 3);
  }
}

// Usage
const engine = new MovieRecommendationEngine(mcpClient);
const userProfile = await engine.getUserProfile('user123');
const recommendations = await engine.generateRecommendations(userProfile);
```

#### Python Integration Example

```python
import asyncio
import json
from typing import List, Dict

class MovieAnalyzer:
    def __init__(self, mcp_client):
        self.mcp = mcp_client
    
    async def analyze_decade_trends(self, decades: List[str]) -> Dict:
        """Analyze movie trends across decades"""
        trends = {}
        
        for decade in decades:
            movies = await self.mcp.call('search_by_decade', {
                'decade': decade
            })
            
            trends[decade] = {
                'total_movies': len(movies),
                'average_rating': self.calculate_avg_rating(movies),
                'top_genres': self.get_top_genres(movies),
                'top_directors': self.get_top_directors(movies)
            }
        
        return trends
    
    async def find_breakout_directors(self, min_movies: int = 3) -> List[Dict]:
        """Find directors with multiple high-rated movies"""
        directors = {}
        
        # Get all movies
        all_movies = await self.mcp.call('search_movies', {'limit': 1000})
        
        for movie in all_movies:
            director = movie['director']
            if director not in directors:
                directors[director] = []
            directors[director].append(movie)
        
        breakout_directors = []
        for director, movies in directors.items():
            if len(movies) >= min_movies:
                avg_rating = sum(m['rating'] for m in movies) / len(movies)
                if avg_rating >= 8.0:
                    breakout_directors.append({
                        'name': director,
                        'movie_count': len(movies),
                        'average_rating': avg_rating,
                        'movies': movies
                    })
        
        return sorted(breakout_directors, key=lambda x: x['average_rating'], reverse=True)

# Usage
analyzer = MovieAnalyzer(mcp_client)
trends = await analyzer.analyze_decade_trends(['1990s', '2000s', '2010s'])
directors = await analyzer.find_breakout_directors()
```

---

## 🎪 Event & Collection Planning

### Scenario: Planning a Film Festival or Movie Night

#### Step 1: Curate Themed Collections

```json
// Create a "Mind-Bending Movies" collection
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Sci-Fi",
      "min_rating": 8.0,
      "limit": 10
    }
  },
  "id": 36
}

// Find psychological thrillers
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Thriller",
      "min_rating": 8.0
    }
  },
  "id": 37
}
```

#### Step 2: Balance Collections

```json
// Get movies of different lengths/intensities
// High-intensity films
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Action",
      "min_rating": 8.0,
      "limit": 5
    }
  },
  "id": 38
}

// Lighter options  
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Comedy",
      "min_rating": 7.5,
      "limit": 5
    }
  },
  "id": 39
}
```

#### Step 3: Historical Progression

```json
// Show evolution of a genre over time
// 1970s Sci-Fi
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "genre": "Sci-Fi",
      "min_year": 1970,
      "max_year": 1979,
      "min_rating": 7.0
    }
  },
  "id": 40
}

// 2010s Sci-Fi
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies", 
    "arguments": {
      "genre": "Sci-Fi",
      "min_year": 2010,
      "max_year": 2019,
      "min_rating": 7.0
    }
  },
  "id": 41
}
```

---

## 🎯 Quick Reference: Common Workflows

### 1. New User Setup
```json
// 1. Add first movie
{"name": "add_movie", "arguments": {...}}

// 2. Search existing collection
{"name": "search_movies", "arguments": {"limit": 10}}

// 3. Get database stats
{"method": "resources/read", "params": {"uri": "movies://database/stats"}}
```

### 2. Daily Movie Discovery
```json
// 1. Check top movies
{"name": "list_top_movies", "arguments": {"limit": 5}}

// 2. Find similar to last watched
{"name": "search_similar_movies", "arguments": {"movie_id": X}}

// 3. Explore by mood/genre
{"name": "search_movies", "arguments": {"genre": "Comedy", "min_rating": 8.0}}
```

### 3. Collection Maintenance
```json
// 1. Update movie details
{"name": "update_movie", "arguments": {"id": X, ...}}

// 2. Add missing actors
{"name": "add_actor", "arguments": {...}}

// 3. Link relationships
{"name": "link_actor_to_movie", "arguments": {"actor_id": X, "movie_id": Y}}
```

### 4. Research & Analysis
```json
// 1. Director filmography
{"name": "search_movies", "arguments": {"director": "Christopher Nolan"}}

// 2. Genre evolution
{"name": "search_by_decade", "arguments": {"decade": "1990s"}}

// 3. Quality analysis
{"name": "search_by_rating_range", "arguments": {"min_rating": 9.0, "max_rating": 10.0}}
```

---

## 🚀 Next Steps

Ready to implement these examples in your own projects?

- **[User Guide](./user-guide.md)** - Master all the tools and features
- **[API Reference](../reference/api.md)** - Complete technical documentation  
- **[Troubleshooting](../reference/troubleshooting.md)** - Fix common issues
- **[Development Guide](../development/README.md)** - Extend and customize

---

*💡 **Pro Tip:** Start with simple examples and gradually build complexity. The MCP server is designed to handle both basic operations and sophisticated workflows!*