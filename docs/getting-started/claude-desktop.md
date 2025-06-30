# 🧠 Claude Desktop Integration

Connect your Movies MCP Server to Claude Desktop for seamless AI-powered movie management.

## Prerequisites

✅ **Movies MCP Server installed** - Complete [Installation Guide](./installation.md) first  
✅ **Claude Desktop app** - [Download here](https://claude.ai/download)  
✅ **Server binary built** - You'll need the full path to your executable

## Step 1: Locate Your Server Binary

### If you used Docker:
```bash
# Find the container path (you'll use the binary path inside container)
docker ps | grep movies-mcp-server-clean

# The binary is at: /usr/local/bin/movies-server-clean (inside container)
# But for Claude Desktop, you need the host machine binary
```

### If you built from source:
```bash
# From your mcp-server directory
cd path/to/your/movies-mcp-server/mcp-server

# Get the absolute path to your binary
pwd
ls build/movies-server-clean

# Example result: /home/user/movies-mcp-server/mcp-server/build/movies-server-clean
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
    "movies-mcp-server": {
      "command": "/FULL/PATH/TO/YOUR/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5433/movies_mcp"
      }
    }
  }
}
```

**🔧 Replace `/FULL/PATH/TO/YOUR/` with your actual path!**

### Environment Setup

Choose your database configuration:

#### Docker Setup (Clean Architecture):
```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/home/user/movies-mcp-server/mcp-server/build/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5433/movies_mcp"
      }
    }
  }
}
```

#### Development Setup:
```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/home/user/movies-mcp-server/mcp-server/build/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

## Step 3: Verify Setup

### Start Your Database
```bash
# For Docker setup
cd movies-mcp-server/mcp-server
make docker-compose-up-clean

# For development setup  
make docker-compose-up-dev
```

### Test Binary Directly
```bash
# Test that your binary works
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
/FULL/PATH/TO/YOUR/movies-server-clean

# Should return MCP initialization response
```

### Restart Claude Desktop
1. **Quit Claude Desktop completely**
2. **Restart Claude Desktop**
3. **Look for connection confirmation**

## Step 4: Test Integration

### In Claude Desktop, try these commands:

```
Can you help me search for movies in my database?
```

```
Add a new movie: "The Matrix" directed by "The Wachowskis" from 1999, rated 8.7
```

```
Show me all the movies in my database
```

## Troubleshooting

### Connection Issues

#### ❌ "Command not found" or "Permission denied"
```bash
# Make binary executable
chmod +x /path/to/your/movies-server-clean

# Test the exact path from your config
/FULL/PATH/TO/YOUR/movies-server-clean --help
```

#### ❌ "Database connection failed"
```bash
# Check database is running
docker ps | grep postgres

# Test database connection
psql "postgres://movies_user:movies_password@localhost:5433/movies_mcp" -c "SELECT 1;"

# Check port conflicts
lsof -i :5433
```

#### ❌ "MCP server not responding"
```bash
# Check server starts correctly
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
/path/to/your/movies-server-clean

# Check Claude Desktop logs (macOS)
tail -f ~/Library/Logs/Claude/claude_desktop.log
```

### Performance Issues

#### ❌ "Server is slow"
```bash
# Check database performance
psql "$DATABASE_URL" -c "SELECT COUNT(*) FROM movies;"

# Monitor resource usage
top -p $(pgrep movies-server-clean)
```

### Configuration Issues

#### ❌ "Invalid JSON configuration"
```bash
# Validate JSON syntax
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .

# If jq not installed: copy-paste into https://jsonlint.com
```

## Advanced Configuration

### Custom Log Level
```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/path/to/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://...",
        "LOG_LEVEL": "debug"
      }
    }
  }
}
```

### Skip Migrations (if database is already setup)
```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/path/to/movies-server-clean",
      "args": ["--skip-migrations"],
      "env": {
        "DATABASE_URL": "postgres://..."
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
      "command": "/path/to/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"
      }
    },
    "movies-prod": {
      "command": "/path/to/movies-server-clean", 
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5433/movies_mcp"
      }
    }
  }
}
```

## What You Can Do Now

With Claude Desktop connected, you can:

- 🔍 **Search movies**: "Find all sci-fi movies from the 1990s"
- ➕ **Add movies**: "Add Inception to my database" 
- 📊 **Get statistics**: "How many movies do I have?"
- 🎭 **Manage actors**: "Add Leonardo DiCaprio to Inception"
- 🖼️ **Handle posters**: "Show me the poster for The Matrix"
- 📈 **Get insights**: "What are my top-rated comedies?"

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
    "movies-mcp-server": {
      "command": "/ABSOLUTE/PATH/TO/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://USER:PASS@HOST:PORT/DATABASE"
      }
    }
  }
}
```

### Port Reference
| Environment | PostgreSQL Port | Use Case |
|-------------|----------------|----------|
| Clean Architecture | 5433 | Production-like |
| Development | 5434 | Development |
| Test | 5435 | Testing only |