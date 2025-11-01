# SQLite Migration Testing Guide

This guide will help you test the SQLite migration locally.

## Prerequisites

- Go 1.23.0 or later
- SQLite3 CLI (optional, for database inspection)

## Quick Start Testing

### 1. Download Dependencies

```bash
cd /path/to/movies-mcp-server
go mod download
```

### 2. Build the Server

```bash
go build -o movies-mcp-server cmd/server-sdk/main.go
```

### 3. Run Automated Tests

```bash
chmod +x test-sqlite-migration.sh
./test-sqlite-migration.sh
```

## Manual Testing Steps

### Test 1: SQLite (Default) - Basic Startup

```bash
# Clean start
rm -f movies.db

# Run migrations
./movies-mcp-server --migrate-only

# Check database created
ls -lh movies.db
# Expected output: movies.db file (~20KB after migrations)

# Verify tables
sqlite3 movies.db "SELECT name FROM sqlite_master WHERE type='table';"
# Expected output:
# schema_migrations
# movies
# actors
# movie_actors
```

### Test 2: Verify Migration Count

```bash
sqlite3 movies.db "SELECT version FROM schema_migrations ORDER BY version;"
# Expected output:
# 1
# 2
# 3
# 4
# 5
```

### Test 3: Check Schema

```bash
# Movies table structure
sqlite3 movies.db ".schema movies"

# Verify genre is TEXT (JSON), not an array
sqlite3 movies.db "SELECT type FROM pragma_table_info('movies') WHERE name='genre';"
# Expected: TEXT
```

### Test 4: Insert Test Data

```bash
sqlite3 movies.db <<EOF
INSERT INTO movies (title, director, year, genre, rating, created_at, updated_at)
VALUES ('The Matrix', 'Wachowski Sisters', 1999, '["Action", "Sci-Fi"]', 8.7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO movies (title, director, year, genre, rating, created_at, updated_at)
VALUES ('Inception', 'Christopher Nolan', 2010, '["Action", "Thriller", "Sci-Fi"]', 8.8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO actors (name, birth_year, created_at, updated_at)
VALUES ('Keanu Reeves', 1964, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
EOF
```

### Test 5: Query Test Data

```bash
# Count movies
sqlite3 movies.db "SELECT COUNT(*) FROM movies;"
# Expected: 2

# Test JSON genre search (this is the key SQLite feature!)
sqlite3 movies.db "SELECT title FROM movies WHERE EXISTS (SELECT 1 FROM json_each(genre) WHERE value = 'Sci-Fi');"
# Expected:
# The Matrix
# Inception

# Test case-insensitive search
sqlite3 movies.db "SELECT title FROM movies WHERE title LIKE '%matrix%' COLLATE NOCASE;"
# Expected: The Matrix
```

### Test 6: Start Server

```bash
# Start server (it will listen on stdin/stdout for MCP protocol)
./movies-mcp-server

# You should see:
# Connected to database: movies.db (driver: sqlite)
# Starting Movies MCP Server with Official SDK...
# ✓ Registered 23 tools successfully
# ✓ Registered 3 resources successfully
# Server ready - listening on stdin/stdout
```

### Test 7: Test with Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-mcp-server/movies-mcp-server"
    }
  }
}
```

Then in Claude Desktop, try:
- "Add a movie: Interstellar, directed by Christopher Nolan, year 2014"
- "Search for movies with Sci-Fi genre"
- "List top rated movies"

## PostgreSQL Compatibility Test

### Test with PostgreSQL (if you have it running)

```bash
# Set environment variables
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=movies_mcp
export DB_USER=movies_user
export DB_PASSWORD=movies_password

# Create database first
psql -U postgres -c "CREATE DATABASE movies_mcp;"
psql -U postgres -c "CREATE USER movies_user WITH PASSWORD 'movies_password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE movies_mcp TO movies_user;"

# Run migrations
./movies-mcp-server --migrate-only

# Start server
./movies-mcp-server
```

## Verify Migration Features

### Feature 1: JSON Genre Storage

```bash
# Insert with JSON array
sqlite3 movies.db "INSERT INTO movies (title, director, year, genre, rating, created_at, updated_at) VALUES ('Test Movie', 'Test Director', 2020, '[\"Drama\", \"Comedy\"]', 7.5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);"

# Query by genre
sqlite3 movies.db "SELECT title, genre FROM movies WHERE EXISTS (SELECT 1 FROM json_each(genre) WHERE value = 'Drama');"
```

### Feature 2: Case-Insensitive Search

```bash
# All these should return the same results
sqlite3 movies.db "SELECT title FROM movies WHERE title LIKE '%matrix%' COLLATE NOCASE;"
sqlite3 movies.db "SELECT title FROM movies WHERE title LIKE '%MATRIX%' COLLATE NOCASE;"
sqlite3 movies.db "SELECT title FROM movies WHERE title LIKE '%Matrix%' COLLATE NOCASE;"
```

### Feature 3: Triggers (Auto-update timestamps)

```bash
# Update a movie
sqlite3 movies.db "UPDATE movies SET rating = 9.0 WHERE title = 'The Matrix';"

# Check updated_at changed
sqlite3 movies.db "SELECT title, rating, updated_at FROM movies WHERE title = 'The Matrix';"
```

### Feature 4: Foreign Keys (CASCADE DELETE)

```bash
# Link actor to movie
sqlite3 movies.db "INSERT INTO movie_actors (movie_id, actor_id, created_at) VALUES (1, 1, CURRENT_TIMESTAMP);"

# Delete movie (should cascade to movie_actors)
sqlite3 movies.db "DELETE FROM movies WHERE id = 1;"

# Verify relationship deleted
sqlite3 movies.db "SELECT COUNT(*) FROM movie_actors WHERE movie_id = 1;"
# Expected: 0
```

## Performance Comparison

### SQLite File Size

```bash
# After migrations only
du -h movies.db
# Expected: ~20KB

# After 100 movies
# Expected: ~100KB

# After 1000 movies with actors
# Expected: ~1-2MB
```

### Query Performance

```bash
# Add timing
sqlite3 movies.db <<EOF
.timer on
SELECT COUNT(*) FROM movies;
SELECT title FROM movies WHERE EXISTS (SELECT 1 FROM json_each(genre) WHERE value = 'Sci-Fi');
EOF
```

## Troubleshooting

### Issue: "missing go.sum entry for modernc.org/sqlite"

```bash
go mod download modernc.org/sqlite
go mod tidy
```

### Issue: Database locked

```bash
# Make sure no other processes are accessing movies.db
fuser movies.db  # Linux
lsof movies.db   # macOS

# Close all connections and retry
```

### Issue: JSON search not working

```bash
# Verify JSON is valid
sqlite3 movies.db "SELECT genre FROM movies WHERE json_valid(genre) = 0;"
# Should return nothing (all genres should be valid JSON)
```

## Expected Results Summary

| Test | Expected Result |
|------|----------------|
| Build | ✓ Binary created successfully |
| Migrations | ✓ 5 migrations applied |
| Database file | ✓ movies.db created (~20KB) |
| Tables | ✓ 4 tables (schema_migrations, movies, actors, movie_actors) |
| Genre storage | ✓ TEXT column with JSON arrays |
| JSON search | ✓ Finds movies by genre using json_each() |
| Case-insensitive | ✓ LIKE COLLATE NOCASE works |
| Triggers | ✓ updated_at auto-updates |
| Foreign keys | ✓ CASCADE DELETE works |
| Server startup | ✓ Logs "driver: sqlite" |

## Next Steps

After successful testing:

1. ✅ Run BDD test suite: `go test ./tests/bdd/...`
2. ✅ Run unit tests: `go test ./...`
3. ✅ Test with Claude Desktop
4. ✅ Create PR for review
5. ✅ Update main README with SQLite instructions

## Need Help?

- Check server logs for detailed error messages
- Verify `sqlite3 --version` is installed
- Ensure Go version is 1.23.0+
- Review migration files in `migrations/` directory
