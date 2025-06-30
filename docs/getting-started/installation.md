# 📦 Installation Guide

Get Movies MCP Server running in 5 minutes with Docker, or 10 minutes building from source.

## Option 1: Docker (Recommended)

### Step 1: Clone and Start
```bash
# Clone the repository
git clone https://github.com/francknouama/movies-mcp-server.git
cd movies-mcp-server/mcp-server

# Start with clean architecture (recommended)
make docker-compose-up-clean
```

### Step 2: Verify Installation
```bash
# Check all services are running
docker-compose -f docker-compose.clean.yml ps

# Should show:
# ✅ movies-mcp-server-clean (healthy)
# ✅ postgres (healthy)  
# ✅ redis-clean (healthy)
```

### Step 3: Test the Server
```bash
# Test MCP protocol communication
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
docker exec -i movies-mcp-server-clean /usr/local/bin/movies-server-clean

# You should see: {"jsonrpc":"2.0","result":{"protocolVersion":"2024-11-05",...},"id":1}
```

**🎉 Success!** Your server is running. Skip to [Add Your First Movie](./first-movie.md).

---

## Option 2: Build from Source

### Prerequisites Check
```bash
# Verify Go version
go version  # Should be 1.24.4+

# Verify Docker
docker version

# ✅ That's it! No other tools needed for Clean Architecture
```

### Step 1: Setup Database
```bash
# From the mcp-server directory
cd movies-mcp-server/mcp-server

# Start development database
make docker-compose-up-dev

# Verify database is ready
psql "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev" -c "SELECT 1;"
```

### Step 2: Build and Run
```bash
# Build clean architecture version (recommended)
make build-clean

# Set environment variables
export DATABASE_URL="postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"
export LOG_LEVEL=info

# Run the server (migrations run automatically!)
./build/movies-server-clean
```

### Step 3: Verify in Another Terminal
```bash
# From the same mcp-server directory, test initialization
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
./build/movies-server-clean
```

---

## Database Migrations

### ✅ Clean Architecture (Recommended)
**Migrations run automatically** - no external tools or setup needed!

- Migrations run on server startup
- Uses built-in custom migration tool
- No external dependencies required

```bash
# Control migration behavior
./build/movies-server-clean                    # Run migrations automatically (default)
./build/movies-server-clean --skip-migrations  # Skip migrations 
./build/movies-server-clean --migrate-only     # Run migrations only, then exit
```

### ⚠️ Legacy Version Only
If using the legacy version, you need [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
# Install golang-migrate (legacy version only)
make install-migrate

# Run migrations manually (legacy version only)  
make db-migrate
```

---

## Configuration

### Environment Variables
```bash
# Required
DATABASE_URL=postgres://user:password@host:port/database

# Optional
LOG_LEVEL=info              # debug, info, warn, error
SERVER_TIMEOUT=30s          # Request timeout
MAX_IMAGE_SIZE=5242880      # 5MB max image size
```

---

## Troubleshooting

### Database Connection Issues
```bash
# Test database connectivity
psql "$DATABASE_URL" -c "SELECT version();"

# If connection fails:
# 1. Check database is running: docker ps
# 2. Verify credentials in DATABASE_URL
# 3. Check network connectivity: telnet host port
```

### Migration Issues (Clean Architecture)
```bash
# Check if migrations completed
./build/movies-server-clean --migrate-only

# If migrations fail, check database logs:
docker-compose -f docker-compose.clean.yml logs postgres
```

### Port Conflicts
If you see "port already in use":
```bash
# Check what's using the port
lsof -i :5433  # Clean architecture PostgreSQL port

# Stop conflicting services
docker-compose down  # stop other Docker services
brew services stop postgresql  # stop local PostgreSQL
```

### Build Issues
```bash
# Clear Go module cache
go clean -modcache

# Update dependencies  
go mod download
go mod tidy

# Rebuild
make clean build-clean
```

---

## Next Steps

✅ **Installation Complete!** 

Choose your next step:
- **[Add Your First Movie](./first-movie.md)** - Learn basic operations
- **[Claude Desktop Setup](./claude-desktop.md)** - Connect to AI assistant
- **[User Guide](../guides/user-guide.md)** - Explore all features

## Quick Reference

| Component | Port | URL | Credentials |
|-----------|------|-----|-------------|
| MCP Server | - | stdin/stdout | - |
| PostgreSQL | 5433 | localhost:5433 | movies_user:movies_password |
| pgAdmin | 5051 | http://localhost:5051 | admin@example.com:admin |
| Health Check | 8081 | http://localhost:8081/health | - |

## Architecture Versions

| Version | Migration Tool | Migration Approach | Best For |
|---------|---------------|-------------------|----------|
| **Clean Architecture** | ✅ Built-in custom tool | 🔄 Automatic on startup | Production, new projects |
| **Legacy** | ⚠️ External golang-migrate | 📋 Manual execution | Existing integrations |