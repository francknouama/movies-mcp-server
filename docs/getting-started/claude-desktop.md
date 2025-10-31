# 🧠 Claude Desktop Integration

Connect your Movies MCP Server to Claude Desktop for seamless AI-powered movie and actor management.

> **🎉 Using the Official Golang MCP SDK!**
> This guide uses the SDK-based server (`movies-mcp-server-sdk`) which provides 23 type-safe tools with automatic schema generation.

## Prerequisites

✅ **Movies MCP Server installed** - Complete [Installation Guide](./installation.md) first
✅ **Claude Desktop app** - [Download here](https://claude.ai/download)
✅ **SDK server built** - You'll need the full path to `movies-mcp-server-sdk` executable

## Step 1: Build and Locate Your SDK Server Binary

### Build the SDK Server:
```bash
# Navigate to project root
cd /path/to/movies-mcp-server

# Build the SDK server
go build -o movies-mcp-server-sdk ./cmd/server-sdk/

# Get the absolute path
pwd
# Example: /home/user/movies-mcp-server

# Your binary is at: /home/user/movies-mcp-server/movies-mcp-server-sdk
```

### Make it Executable (if needed):
```bash
chmod +x movies-mcp-server-sdk
```

### Test the Binary:
```bash
# Test with --help flag
./movies-mcp-server-sdk --help

# Should show usage information
```

## Step 2: Configure Claude Desktop

### Find Claude Desktop Config

| OS | Config File Location |
|----|---------------------|
| **macOS** | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| **Windows** | `%APPDATA%\Claude\claude_desktop_config.json` |
| **Linux** | `~/.config/Claude/claude_desktop_config.json` |

### Add Server Configuration

Edit your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "movies": {
      "command": "/absolute/path/to/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

**🔧 Replace `/absolute/path/to/` with your actual path!**

### Environment Setup Examples

#### Local PostgreSQL:
```json
{
  "mcpServers": {
    "movies": {
      "command": "/home/user/movies-mcp-server/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

#### Docker PostgreSQL (Custom Port):
```json
{
  "mcpServers": {
    "movies": {
      "command": "/home/user/movies-mcp-server/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5433",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

#### With Skip Migrations:
```json
{
  "mcpServers": {
    "movies": {
      "command": "/home/user/movies-mcp-server/movies-mcp-server-sdk",
      "args": ["--skip-migrations"],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

## Step 3: Verify Setup

### Start Your Database
```bash
# If using Docker
cd movies-mcp-server
make docker-up

# Or start PostgreSQL locally
# PostgreSQL should be running on port 5432 (or your configured port)
```

### Test Binary Directly
```bash
# Test that your binary works
./movies-mcp-server-sdk --version

# Test help
./movies-mcp-server-sdk --help

# Test MCP protocol (should output initialization response)
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
./movies-mcp-server-sdk
```

### Restart Claude Desktop
1. **Quit Claude Desktop completely**
2. **Restart Claude Desktop**
3. **Look for the Movies server** in available tools

## Step 4: Test Integration

### Verify All 23 Tools Are Available

In Claude Desktop, you should see:
- **8 Movie Tools**: get_movie, add_movie, update_movie, delete_movie, list_top_movies, search_movies, search_by_decade, search_by_rating_range
- **9 Actor Tools**: get_actor, add_actor, update_actor, delete_actor, link_actor_to_movie, unlink_actor_from_movie, get_movie_cast, get_actor_movies, search_actors
- **3 Compound Tools**: bulk_movie_import, movie_recommendation_engine, director_career_analysis
- **3 Context Tools**: create_search_context, get_context_page, get_context_info

### Try These Commands:

**Movies:**
```
Can you help me search for movies in my database?
```

```
Add a new movie: "The Matrix" directed by "The Wachowskis" from 1999, rated 8.7
```

```
Show me all movies from the 1990s
```

**Actors:**
```
Add Leonardo DiCaprio as an actor born in 1974
```

```
Link Leonardo DiCaprio to Inception
```

```
Show me all movies that Tom Hanks appears in
```

**Smart Features:**
```
Give me movie recommendations based on sci-fi and action genres
```

```
Analyze Christopher Nolan's career
```

## Troubleshooting

### Connection Issues

#### ❌ "Command not found" or "Permission denied"
```bash
# Make binary executable
chmod +x ./movies-mcp-server-sdk

# Test the exact path from your config
/absolute/path/to/movies-mcp-server-sdk --help
```

#### ❌ "Database connection failed"
```bash
# Check database is running
docker ps | grep postgres

# Test database connection with your settings
psql -h localhost -p 5432 -U movies_user -d movies_mcp -c "SELECT 1;"

# Check port conflicts
lsof -i :5432

# Verify environment variables are set correctly in Claude config
```

#### ❌ "MCP server not responding"
```bash
# Check server starts correctly
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
./movies-mcp-server-sdk

# Check Claude Desktop logs
# macOS:
tail -f ~/Library/Logs/Claude/mcp*.log

# Linux:
tail -f ~/.config/Claude/logs/mcp*.log
```

### Performance Issues

#### ❌ "Server is slow"
```bash
# Check database performance
psql -h localhost -p 5432 -U movies_user -d movies_mcp -c "SELECT COUNT(*) FROM movies;"

# Monitor resource usage
ps aux | grep movies-mcp-server-sdk
```

### Configuration Issues

#### ❌ "Invalid JSON configuration"
```bash
# Validate JSON syntax
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .

# If jq not installed: copy-paste into https://jsonlint.com
```

## Advanced Configuration

### Skip Migrations (if database is already setup)
```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-mcp-server-sdk",
      "args": ["--skip-migrations"],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

### Multiple Environments
```json
{
  "mcpServers": {
    "movies-dev": {
      "command": "/path/to/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5434",
        "DB_USER": "dev_user",
        "DB_PASSWORD": "dev_password",
        "DB_NAME": "movies_mcp_dev",
        "DB_SSLMODE": "disable"
      }
    },
    "movies-prod": {
      "command": "/path/to/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "movies_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "require"
      }
    }
  }
}
```

### SSL/TLS Connection
```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "production-db.example.com",
        "DB_PORT": "5432",
        "DB_USER": "movies_user",
        "DB_PASSWORD": "secure_password",
        "DB_NAME": "movies_mcp",
        "DB_SSLMODE": "require"
      }
    }
  }
}
```

## What You Can Do Now

With Claude Desktop connected, you can leverage all 23 tools:

### 🎬 Movie Management (8 tools)
- 🔍 **Search & filter**: "Find all sci-fi movies from the 1990s"
- ➕ **Add movies**: "Add Inception to my database"
- ✏️ **Update movies**: "Update The Matrix rating to 8.7"
- 🗑️ **Delete movies**: "Remove movie ID 42"
- 📊 **Top movies**: "Show me the top 10 highest-rated movies"
- 📅 **Search by decade**: "Find movies from the 80s"
- ⭐ **Rating ranges**: "Show movies rated between 7.0 and 8.0"

### 🎭 Actor Management (9 tools)
- 👤 **Manage actors**: "Add Leonardo DiCaprio born in 1974"
- 🔗 **Link relationships**: "Link Tom Hanks to Forrest Gump"
- 📋 **Get cast**: "Show me all actors in The Matrix"
- 🎬 **Actor filmography**: "What movies has Brad Pitt appeared in?"
- 🔍 **Search actors**: "Find all actors born in the 1960s"

### 🧠 Smart Features (3 compound tools)
- 📥 **Bulk import**: "Import 10 movies at once"
- 💡 **Recommendations**: "Recommend movies based on sci-fi and thriller genres with min rating 7.5"
- 📊 **Career analysis**: "Analyze Christopher Nolan's career trajectory"

### 📑 Pagination (3 context tools)
- **Large result sets**: Automatically handles pagination for 100+ results
- **Context management**: Navigate through pages of search results
- **Performance**: Efficient handling of large datasets

## Next Steps

✅ **Integration Complete!**

- **[Add Your First Movie](./first-movie.md)** - Learn the tools hands-on
- **[User Guide](../guides/user-guide.md)** - Discover all features
- **[Examples](../guides/examples.md)** - See real-world use cases

---

## Quick Reference

### Essential Config Template
```json
{
  "mcpServers": {
    "movies": {
      "command": "/ABSOLUTE/PATH/TO/movies-mcp-server-sdk",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "USER",
        "DB_PASSWORD": "PASSWORD",
        "DB_NAME": "DATABASE",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

### Environment Variables Reference
| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `movies_user` |
| `DB_PASSWORD` | Database password | `movies_password` |
| `DB_NAME` | Database name | `movies_mcp` |
| `DB_SSLMODE` | SSL mode | `disable`, `require`, `verify-full` |

### Command Line Flags
| Flag | Description |
|------|-------------|
| `--version` | Show version information |
| `--help` | Show help message |
| `--skip-migrations` | Skip database migrations on startup |

### Port Reference
| Environment | PostgreSQL Port | Use Case |
|-------------|----------------|----------|
| Default | 5432 | Standard PostgreSQL |
| Docker (custom) | 5433 | Docker mapped port |
| Development | 5434 | Development environment |
| Test | 5435 | Testing only |

---

## SDK Migration Notes

This guide uses the **SDK-based server** (`movies-mcp-server-sdk`) which is built with the official Golang MCP SDK v1.1.0.

**Benefits over the legacy server:**
- ✅ Type-safe tool handlers
- ✅ Automatic JSON schema generation
- ✅ Better error handling
- ✅ Official support from Anthropic & Google
- ✅ 26% less code

**Legacy server:** If you need the old custom server (`movies-server-clean`), it's still available but not recommended for new deployments.

For more details, see [SDK Migration Complete](../SDK_MIGRATION_COMPLETE.md).