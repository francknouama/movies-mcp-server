# 🔧 Troubleshooting Guide

Comprehensive troubleshooting guide for the Movies MCP Server. Find solutions to common issues, understand error codes, and get back up and running quickly.

## 📋 Table of Contents

1. [🚨 Quick Diagnostics](#-quick-diagnostics)
2. [🔗 Connection Issues](#-connection-issues)
3. [📊 Database Problems](#-database-problems)
4. [🛠️ Tool Execution Errors](#-tool-execution-errors)
5. [🖼️ Image & Resource Issues](#-image--resource-issues)
6. [⚡ Performance Problems](#-performance-problems)
7. [🔍 Error Code Reference](#-error-code-reference)
8. [🧪 Debug Tools & Techniques](#-debug-tools--techniques)
9. [💡 Prevention & Best Practices](#-prevention--best-practices)

---

## 🚨 Quick Diagnostics

### Server Health Check

**Step 1: Verify Server is Running**
```bash
# Check if process is running
ps aux | grep movies-mcp-server

# Check logs
tail -f /var/log/movies-mcp-server.log
```

**Step 2: Test Basic Connectivity**
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {}
  },
  "id": 1
}
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

**Step 3: Database Connection Test**
```bash
# Test database connectivity
psql "$DATABASE_URL" -c "SELECT COUNT(*) FROM movies;"
```

---

## 🔗 Connection Issues

### Problem: "Connection Refused" / "Server Not Responding"

**Symptoms:**
- Cannot connect to MCP server
- Timeout errors
- "Connection refused" messages

**Solutions:**

#### For Claude Desktop Integration

**1. Check Configuration File**
```bash
# macOS
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

# Windows  
cat %APPDATA%\Claude\claude_desktop_config.json

# Linux
cat ~/.config/Claude/claude_desktop_config.json
```

**Common Configuration Issues:**
```json
{
  "mcpServers": {
    "movies-mcp-server": {
      // ❌ WRONG: Relative path
      "command": "./movies-mcp-server",
      
      // ✅ CORRECT: Absolute path
      "command": "/absolute/path/to/movies-mcp-server/build/movies-server-clean",
      
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5432/movies_db"
      }
    }
  }
}
```

**2. Verify Binary Permissions**
```bash
# Check if executable
ls -la /path/to/movies-mcp-server
# Should show: -rwxr-xr-x

# Fix permissions if needed
chmod +x /path/to/movies-mcp-server
```

**3. Test Binary Directly**
```bash
# Test server startup
/absolute/path/to/movies-mcp-server/build/movies-server-clean
# Should not exit immediately and should accept stdin
```

#### For Direct MCP Communication

**1. Check Server Process**
```bash
# Start server manually
./build/movies-server-clean

# In another terminal, test
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | ./build/movies-server-clean
```

**2. Environment Variables**
```bash
# Check required environment variables
echo $DATABASE_URL
# Should output: postgres://user:password@host:port/database

# Set if missing
export DATABASE_URL="postgres://movies_user:movies_password@localhost:5432/movies_db"
```

### Problem: "Protocol Version Mismatch"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Unsupported protocol version"
  },
  "id": 1
}
```

**Solution:**
Use the correct protocol version in initialization:

```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",  // ✅ Correct version
    "capabilities": {}
  },
  "id": 1
}
```

---

## 📊 Database Problems

### Problem: "Database Connection Failed"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "Database connection failed"
  },
  "id": 1
}
```

**Diagnosis Steps:**

**1. Check Database Server**
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql
# or on macOS with Homebrew:
brew services list | grep postgresql

# Start if not running
sudo systemctl start postgresql
# or on macOS:
brew services start postgresql
```

**2. Verify Database Exists**
```bash
# Connect to PostgreSQL
psql -U postgres -h localhost

# List databases
\l

# Should see 'movies_db' in the list
```

**3. Check Connection String**
```bash
# Test connection string manually
psql "postgres://movies_user:movies_password@localhost:5432/movies_db"
```

**Common Connection String Issues:**
```bash
# ❌ WRONG: Missing password
DATABASE_URL="postgres://movies_user@localhost:5432/movies_db"

# ❌ WRONG: Wrong port
DATABASE_URL="postgres://movies_user:movies_password@localhost:3306/movies_db"

# ✅ CORRECT
DATABASE_URL="postgres://movies_user:movies_password@localhost:5432/movies_db"
```

### Problem: "Table Does Not Exist"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "relation \"movies\" does not exist"
  },
  "id": 1
}
```

**Solution: Run Database Migrations**

**For Clean Architecture Version:**
```bash
cd mcp-server
./build/movies-server-clean -migrate-only
```

**For Legacy Version:**
```bash
cd mcp-server
make migrate-up
```

**Manual Migration (if needed):**
```bash
# Check current schema
psql "$DATABASE_URL" -c "\dt"

# Run migrations manually
psql "$DATABASE_URL" -f migrations/001_create_movies_table.up.sql
psql "$DATABASE_URL" -f migrations/002_add_indexes.up.sql
psql "$DATABASE_URL" -f migrations/003_add_search_indexes.up.sql
psql "$DATABASE_URL" -f migrations/004_create_actors_tables.up.sql
psql "$DATABASE_URL" -f migrations/005_align_schema_with_domain.up.sql
```

### Problem: "Permission Denied" Database Errors

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "permission denied for table movies"
  },
  "id": 1
}
```

**Solution: Fix Database Permissions**
```sql
-- Connect as superuser
psql -U postgres

-- Grant permissions to movies_user
GRANT ALL PRIVILEGES ON DATABASE movies_db TO movies_user;
\c movies_db
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO movies_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO movies_user;
```

---

## 🛠️ Tool Execution Errors

### Problem: "Tool Not Found"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32601,
    "message": "Method not found: unknown_tool"
  },
  "id": 1
}
```

**Solution:**
Check available tools first:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 1
}
```

**Available Tools:**
- `add_movie`, `get_movie`, `update_movie`, `delete_movie`, `list_top_movies`
- `add_actor`, `get_actor`, `update_actor`, `delete_actor`, `search_actors`
- `link_actor_to_movie`, `unlink_actor_from_movie`, `get_movie_cast`, `get_actor_movies`
- `search_movies`, `search_by_decade`, `search_by_rating_range`, `search_similar_movies`

### Problem: "Invalid Parameters"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid parameters: missing required field 'title'"
  },
  "id": 1
}
```

**Common Parameter Issues:**

#### Movie Parameters
```json
// ❌ WRONG: Missing required fields
{
  "name": "add_movie",
  "arguments": {
    "title": "The Matrix"
    // Missing director and year
  }
}

// ✅ CORRECT: All required fields
{
  "name": "add_movie", 
  "arguments": {
    "title": "The Matrix",
    "director": "The Wachowskis",
    "year": 1999
  }
}
```

#### Rating Validation
```json
// ❌ WRONG: Rating out of range
{
  "name": "add_movie",
  "arguments": {
    "title": "Movie",
    "director": "Director",
    "year": 2024,
    "rating": 11  // Must be 0-10
  }
}

// ✅ CORRECT: Valid rating
{
  "name": "add_movie",
  "arguments": {
    "title": "Movie",
    "director": "Director", 
    "year": 2024,
    "rating": 8.5
  }
}
```

#### ID Parameters
```json
// ❌ WRONG: String ID instead of integer
{
  "name": "get_movie",
  "arguments": {
    "movie_id": "1"  // Should be integer
  }
}

// ✅ CORRECT: Integer ID
{
  "name": "get_movie",
  "arguments": {
    "movie_id": 1
  }
}
```

### Problem: "Movie Not Found"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Movie not found with ID: 999"
  },
  "id": 1
}
```

**Solutions:**

**1. Verify Movie Exists**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "limit": 10
    }
  },
  "id": 1
}
```

**2. Check Movie ID**
Look for the correct ID in the search results or when adding movies.

### Problem: "Actor Already Linked"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Actor is already linked to this movie"
  },
  "id": 1
}
```

**Solution:**
Check current cast before linking:

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
  "id": 1
}
```

---

## 🖼️ Image & Resource Issues

### Problem: "Failed to Download Poster"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "Failed to download poster from URL"
  },
  "id": 1
}
```

**Solutions:**

**1. Check URL Accessibility**
```bash
# Test URL manually
curl -I "https://example.com/poster.jpg"
# Should return 200 OK
```

**2. Common URL Issues**
```json
// ❌ WRONG: Not a direct image URL
{
  "poster_url": "https://example.com/movie-page"
}

// ❌ WRONG: Unsupported format
{
  "poster_url": "https://example.com/poster.gif"  
}

// ✅ CORRECT: Direct image URL with supported format
{
  "poster_url": "https://example.com/poster.jpg"
}
```

**Supported Image Formats:**
- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)

### Problem: "Resource Not Found"

**Error Message:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Resource not found: movies://invalid/path"
  },
  "id": 1
}
```

**Solution:**
Check available resources:

```json
{
  "jsonrpc": "2.0",
  "method": "resources/list",
  "id": 1
}
```

**Available Resource URIs:**
- `movies://database/all` - Complete movie database
- `movies://database/stats` - Database statistics  
- `movies://posters/collection` - All posters collection
- `movies://posters/{id}` - Individual movie poster

---

## ⚡ Performance Problems

### Problem: "Slow Query Response"

**Symptoms:**
- Search queries taking > 5 seconds
- Database timeouts
- Memory usage growing

**Solutions:**

**1. Check Database Indexes**
```sql
-- Connect to database
psql "$DATABASE_URL"

-- Check if indexes exist
\d movies
\d actors

-- Should see indexes on commonly queried fields
```

**2. Optimize Search Parameters**
```json
// ❌ SLOW: Very broad search
{
  "name": "search_movies",
  "arguments": {
    "limit": 10000  // Too large
  }
}

// ✅ FAST: Specific search with reasonable limit
{
  "name": "search_movies",
  "arguments": {
    "genre": "Sci-Fi",
    "min_rating": 8.0,
    "limit": 50
  }
}
```

**3. Database Maintenance**
```sql
-- Run database maintenance
ANALYZE movies;
ANALYZE actors;
VACUUM movies;
VACUUM actors;
```

### Problem: "Memory Leak" / "High Memory Usage"

**Diagnosis:**
```bash
# Monitor memory usage
top -p $(pgrep movies-mcp-server)

# Check for memory leaks
valgrind --leak-check=full ./build/movies-server-clean
```

**Solutions:**
1. Restart server regularly in production
2. Use reasonable result limits
3. Update to latest version (memory optimizations)

---

## 🔍 Error Code Reference

### JSON-RPC Standard Errors

| Code | Name | Description | Solution |
|------|------|-------------|----------|
| `-32700` | Parse Error | Invalid JSON | Check JSON syntax |
| `-32600` | Invalid Request | Invalid JSON-RPC | Verify request format |
| `-32601` | Method Not Found | Unknown tool name | Check available tools |
| `-32602` | Invalid Params | Parameter validation failed | Verify parameter format |
| `-32603` | Internal Error | Server error | Check logs and database |

### Movies MCP Server Specific Errors

| Error Pattern | Cause | Solution |
|---------------|-------|----------|
| `"Movie not found with ID: X"` | Movie doesn't exist | Verify movie ID with search |
| `"Actor not found with ID: X"` | Actor doesn't exist | Verify actor ID with search |
| `"Movie already exists: Title"` | Duplicate movie | Check existing movies first |
| `"Actor already linked to movie"` | Duplicate relationship | Check current cast |
| `"Rating must be between 0 and 10"` | Invalid rating value | Use rating 0.0-10.0 |
| `"Year must be a valid integer"` | Invalid year format | Use 4-digit year |
| `"Failed to download poster"` | Image URL issue | Check URL accessibility |
| `"Database connection failed"` | DB connectivity issue | Check database status |
| `"Permission denied for table"` | DB permissions issue | Fix database permissions |

---

## 🧪 Debug Tools & Techniques

### Enable Debug Logging

**Environment Variable:**
```bash
export LOG_LEVEL=debug
./build/movies-server-clean
```

**Expected Debug Output:**
```
2024-01-15T10:30:00Z DEBUG: Received request: {"jsonrpc":"2.0",...}
2024-01-15T10:30:00Z DEBUG: Executing tool: add_movie
2024-01-15T10:30:00Z DEBUG: Database query: INSERT INTO movies...
2024-01-15T10:30:00Z DEBUG: Query result: Success, ID: 123
2024-01-15T10:30:00Z DEBUG: Sending response: {"jsonrpc":"2.0",...}
```

### Manual Testing with curl

**Test Tool Execution:**
```bash
# Create test request file
cat > test_request.json << EOF
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "Test Movie",
      "director": "Test Director",
      "year": 2024
    }
  },
  "id": 1
}
EOF

# Send via stdin
cat test_request.json | ./build/movies-server-clean
```

### Database Query Debugging

**Enable PostgreSQL Query Logging:**
```sql
-- Connect as superuser
psql -U postgres

-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
SELECT pg_reload_conf();

-- Check logs
tail -f /var/log/postgresql/postgresql-main.log
```

### Network Debugging (Claude Desktop)

**Check Claude Desktop Logs:**
```bash
# macOS
tail -f ~/Library/Logs/Claude/app.log

# Windows
tail -f %APPDATA%\Claude\logs\app.log

# Linux  
tail -f ~/.config/Claude/logs/app.log
```

---

## 💡 Prevention & Best Practices

### Regular Maintenance

**1. Database Maintenance**
```bash
# Weekly maintenance script
#!/bin/bash
psql "$DATABASE_URL" << EOF
ANALYZE movies;
ANALYZE actors;
VACUUM movies;
VACUUM actors;
REINDEX TABLE movies;
REINDEX TABLE actors;
EOF
```

**2. Log Rotation**
```bash
# Setup logrotate for server logs
sudo tee /etc/logrotate.d/movies-mcp-server << EOF
/var/log/movies-mcp-server.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
}
EOF
```

### Monitoring Setup

**1. Health Check Script**
```bash
#!/bin/bash
# health_check.sh

# Test server response
response=$(echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | timeout 5 ./build/movies-server-clean)

if [[ $? -eq 0 ]] && [[ "$response" == *"tools"* ]]; then
    echo "✅ Server healthy"
    exit 0
else
    echo "❌ Server unhealthy"
    exit 1
fi
```

**2. Database Connection Check**
```bash
#!/bin/bash
# db_check.sh

if psql "$DATABASE_URL" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Database healthy"
    exit 0
else
    echo "❌ Database unhealthy"  
    exit 1
fi
```

### Configuration Best Practices

**1. Environment Configuration**
```bash
# Use a .env file
cat > .env << EOF
# Database Configuration
DATABASE_URL=postgres://movies_user:secure_password@localhost:5432/movies_db

# Logging Configuration  
LOG_LEVEL=info
LOG_FILE=/var/log/movies-mcp-server.log

# Performance Configuration
MAX_CONNECTIONS=100
QUERY_TIMEOUT=30s
EOF

# Load in startup script
source .env
./build/movies-server-clean
```

**2. Security Configuration**
```bash
# Restrict file permissions
chmod 600 .env
chmod 700 ./build/movies-server-clean

# Use non-root user
sudo useradd -m -s /bin/bash movies-server
sudo chown -R movies-server:movies-server /path/to/movies-mcp-server
```

---

## 🆘 Emergency Recovery

### Complete Database Reset

**⚠️ WARNING: This will delete all data**

```bash
# 1. Stop server
pkill movies-mcp-server

# 2. Drop and recreate database
psql -U postgres << EOF
DROP DATABASE IF EXISTS movies_db;
CREATE DATABASE movies_db;
GRANT ALL PRIVILEGES ON DATABASE movies_db TO movies_user;
EOF

# 3. Run migrations
cd mcp-server
./build/movies-server-clean -migrate-only

# 4. Restart server
./build/movies-server-clean
```

### Configuration Reset

**Claude Desktop Configuration Reset:**
```bash
# Backup current config
cp ~/Library/Application\ Support/Claude/claude_desktop_config.json ~/claude_config_backup.json

# Reset to minimal config
cat > ~/Library/Application\ Support/Claude/claude_desktop_config.json << EOF
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/absolute/path/to/movies-mcp-server/build/movies-server-clean",
      "args": [],
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5432/movies_db"
      }
    }
  }
}
EOF

# Restart Claude Desktop
```

---

## 📞 Getting Additional Help

### Before Seeking Help

1. ✅ Check this troubleshooting guide
2. ✅ Review error messages carefully
3. ✅ Test with minimal examples
4. ✅ Check server and database logs
5. ✅ Try the emergency recovery steps

### Information to Include

When reporting issues, include:

- **Server Version:** Check with `./build/movies-server-clean --version`
- **Database Version:** `psql --version`
- **Operating System:** `uname -a`
- **Error Message:** Complete JSON-RPC error response
- **Steps to Reproduce:** Minimal example that causes the issue
- **Server Logs:** Recent log entries around the error
- **Configuration:** Sanitized configuration (remove passwords)

### Where to Get Help

- **GitHub Issues:** [Report bugs and request features](https://github.com/francknouama/movies-mcp-server/issues)
- **Documentation:** [Complete user guides](../guides/user-guide.md)
- **Examples:** [Common use cases](../guides/examples.md)

---

*🔧 **Remember:** Most issues are caused by configuration problems or missing dependencies. Start with the basics and work your way up to more complex diagnostics.*