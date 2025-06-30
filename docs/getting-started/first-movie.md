# 🎬 Add Your First Movie

Learn the Movies MCP Server by adding your first movie and exploring the core features.

## Prerequisites

✅ **Server running** - Complete [Installation](./installation.md) or [Claude Desktop Integration](./claude-desktop.md)  
✅ **Database connected** - You should see "Connected to database" in server logs

## Method 1: Using Claude Desktop (Recommended)

If you've connected to Claude Desktop, just ask Claude naturally:

```
Add "The Matrix" to my movie database. It's a 1999 sci-fi film directed by The Wachowskis, with a rating of 8.7 out of 10.
```

Claude will use the MCP tools automatically! Skip to [Step 3: Explore Your Movie](#step-3-explore-your-movie).

---

## Method 2: Direct MCP Communication

### Step 1: Initialize the Connection

```bash
# Send initialization request
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
./build/movies-server-clean
```

**Expected Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {},
      "resources": {}
    },
    "serverInfo": {
      "name": "movies-mcp-server",
      "version": "0.2.0"
    }
  },
  "id": 1
}
```

### Step 2: See Available Tools

```bash
# List all available tools
echo '{"jsonrpc":"2.0","method":"tools/list","id":2}' | \
./build/movies-server-clean
```

**You'll see tools like:**
- `add_movie` - Add a new movie
- `search_movies` - Search the database
- `get_movie` - Get movie details
- `update_movie` - Update movie info
- `delete_movie` - Remove movies

### Step 3: Add Your First Movie

```bash
# Add "The Matrix" movie
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Matrix",
      "director": "The Wachowskis",
      "release_year": 1999,
      "genre": "Sci-Fi",
      "rating": 8.7,
      "description": "A computer hacker learns about the true nature of reality and his role in the war against its controllers."
    }
  },
  "id": 3
}' | ./build/movies-server-clean
```

**Expected Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Movie 'The Matrix' added successfully with ID: 1"
      }
    ]
  },
  "id": 3
}
```

🎉 **Congratulations!** You've added your first movie.

---

## Step 4: Explore Your Movie

### Search for Your Movie
```bash
# Search by title
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "query": "Matrix",
      "search_type": "title"
    }
  },
  "id": 4
}' | ./build/movies-server-clean
```

### Get Movie Details
```bash
# Get full details (replace ID with your movie's ID)
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "get_movie",
    "arguments": {
      "movie_id": 1
    }
  },
  "id": 5
}' | ./build/movies-server-clean
```

### Search by Genre
```bash
# Find all sci-fi movies
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies", 
    "arguments": {
      "query": "Sci-Fi",
      "search_type": "genre"
    }
  },
  "id": 6
}' | ./build/movies-server-clean
```

---

## Step 5: Add More Movies

Try adding a few more movies to build your collection:

### Add a Classic
```bash
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Casablanca",
      "director": "Michael Curtiz", 
      "release_year": 1942,
      "genre": "Drama",
      "rating": 8.5,
      "description": "A cynical American expatriate struggles to decide whether or not he should help his former lover and her fugitive husband escape French Morocco."
    }
  },
  "id": 7
}' | ./build/movies-server-clean
```

### Add a Recent Film
```bash
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Dune",
      "director": "Denis Villeneuve",
      "release_year": 2021,
      "genre": "Sci-Fi",
      "rating": 8.0,
      "description": "Feature adaptation of Frank Herbert's science fiction novel about the son of a noble family entrusted with the protection of the most valuable asset and most vital element in the galaxy."
    }
  },
  "id": 8
}' | ./build/movies-server-clean
```

---

## Step 6: Advanced Operations

### Update a Movie
```bash
# Add a plot twist to The Matrix description
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "update_movie",
    "arguments": {
      "movie_id": 1,
      "description": "A computer hacker learns about the true nature of reality and his role in the war against its controllers. Features groundbreaking visual effects and philosophical themes about reality and choice."
    }
  },
  "id": 9
}' | ./build/movies-server-clean
```

### Get Top Movies
```bash
# Find your highest-rated movies
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_top_movies",
    "arguments": {
      "limit": 5
    }
  },
  "id": 10
}' | ./build/movies-server-clean
```

---

## Working with Images (Optional)

### Add a Movie with Poster URL
```bash
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Inception",
      "director": "Christopher Nolan",
      "release_year": 2010,
      "genre": "Sci-Fi",
      "rating": 8.8,
      "description": "A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.",
      "poster_url": "https://example.com/inception-poster.jpg"
    }
  },
  "id": 11
}' | ./build/movies-server-clean
```

The server will download and store the poster automatically!

---

## Troubleshooting

### ❌ "Movie already exists"
```bash
# Search first to see if it exists
echo '{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "query": "Matrix",
      "search_type": "title"
    }
  },
  "id": 12
}' | ./build/movies-server-clean
```

### ❌ "Invalid rating"
Ratings must be between 0.0 and 10.0:
```json
{
  "rating": 8.7  // ✅ Valid
  "rating": 11   // ❌ Invalid
  "rating": -1   // ❌ Invalid
}
```

### ❌ "Database connection failed"
```bash
# Check database is running
docker ps | grep postgres

# Test connection
psql "$DATABASE_URL" -c "SELECT COUNT(*) FROM movies;"
```

---

## What You've Learned

✅ **MCP Protocol** - How to communicate with the server  
✅ **Core Tools** - `add_movie`, `search_movies`, `get_movie`  
✅ **Data Validation** - Proper movie data format  
✅ **Error Handling** - How to troubleshoot issues  
✅ **Advanced Features** - Updates, ratings, images

## Next Steps

🎯 **Ready to explore more?**

- **[User Guide](../guides/user-guide.md)** - Discover all features and tools
- **[Examples](../guides/examples.md)** - See real-world use cases
- **[API Reference](../reference/api.md)** - Complete tool documentation

## Quick Commands Reference

```bash
# Initialize
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | ./build/movies-server-clean

# List tools  
echo '{"jsonrpc":"2.0","method":"tools/list","id":2}' | ./build/movies-server-clean

# Add movie
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"add_movie","arguments":{"title":"Movie Title","director":"Director Name","release_year":2024,"genre":"Genre","rating":8.0}},"id":3}' | ./build/movies-server-clean

# Search movies
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_movies","arguments":{"query":"search term","search_type":"title"}},"id":4}' | ./build/movies-server-clean
```