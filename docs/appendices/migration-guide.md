# 🔄 Migration Guide: Legacy to Clean Architecture

Complete guide for migrating from the Legacy architecture to the new Clean Architecture implementation of the Movies MCP Server.

## 📋 Migration Overview

The Movies MCP Server has evolved from a simple Legacy implementation to a robust Clean Architecture design. This guide helps you migrate seamlessly while preserving your data and configurations.

### 🎯 Why Migrate?

| Aspect | Legacy Architecture | Clean Architecture |
|--------|-------------------|-------------------|
| **🏗️ Code Structure** | Monolithic, tightly coupled | Layered, separation of concerns |
| **🧪 Testing** | Limited unit tests | Comprehensive test suite |
| **📊 Performance** | Basic queries, limited optimization | Optimized repositories, efficient queries |
| **🔧 Maintainability** | Difficult to extend | Easy to modify and extend |
| **📖 Documentation** | Minimal documentation | Well-documented codebase |
| **🚀 Future Support** | Maintenance mode only | Active development |

### 🚨 Migration Impact

**Breaking Changes**:
- ⚠️ **Binary name change**: `movies-server` → `movies-server-clean`
- ⚠️ **Port changes**: 8080 → 8081 (health check)
- ⚠️ **Database ports**: 5432 → 5433 (Docker)
- ⚠️ **Some environment variable names**

**Preserved**:
- ✅ **Database schema**: Fully compatible
- ✅ **MCP protocol**: Same tools and resources
- ✅ **Data integrity**: No data loss
- ✅ **Basic configuration**: Most settings preserved

## 🗂️ Pre-Migration Checklist

### 1. Backup Current Setup

**Database Backup**:
```bash
# Create complete database backup
pg_dump -h localhost -U movies_user -p 5432 movies_mcp > legacy_backup_$(date +%Y%m%d).sql

# Verify backup integrity
psql -h localhost -U movies_user -d postgres -c "CREATE DATABASE movies_test_restore;"
psql -h localhost -U movies_user -d movies_test_restore < legacy_backup_$(date +%Y%m%d).sql
psql -h localhost -U movies_user -d postgres -c "DROP DATABASE movies_test_restore;"
```

**Configuration Backup**:
```bash
# Backup current environment files
cp .env .env.legacy.backup
cp docker-compose.yml docker-compose.legacy.backup

# Backup Claude Desktop configuration
cp ~/Library/Application\ Support/Claude/claude_desktop_config.json claude_config.legacy.backup
```

**Docker Volume Backup** (if using Docker):
```bash
# Create volume backups
docker run --rm -v movies_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/legacy_postgres_data.tar.gz -C /data .
docker run --rm -v movies_redis_data:/data -v $(pwd):/backup alpine tar czf /backup/legacy_redis_data.tar.gz -C /data .
```

### 2. Document Current Configuration

**Record current settings**:
```bash
# Document environment variables
env | grep -E "(DB_|DATABASE_|LOG_|SERVER_)" > current_config.txt

# Document running services
docker-compose ps > current_services.txt

# Document port usage
netstat -tlnp | grep -E "(5432|8080|6379|9090|3000)" > current_ports.txt
```

### 3. Verify System Requirements

**Check versions**:
```bash
# Docker version (required: 20.10+)
docker --version

# Docker Compose version (required: 2.0+)
docker-compose --version

# Available disk space (recommended: 2GB+ free)
df -h

# Available memory (recommended: 2GB+ RAM)
free -h
```

## 🚀 Migration Process

### Step 1: Stop Legacy Services

**Docker Deployment**:
```bash
# Stop legacy services gracefully  
docker-compose down

# Verify all containers stopped
docker ps | grep -E "(movies|postgres|redis)"
```

**Binary Deployment**:
```bash
# Stop system service (if configured)
sudo systemctl stop movies-mcp

# Or kill process manually
pkill -f movies-server
```

### Step 2: Get Clean Architecture Code

```bash
# If not already done, clone/update repository
git clone https://github.com/francknouama/movies-mcp-server.git
cd movies-mcp-server

# Or update existing repository
git fetch origin
git checkout main
git pull origin main

# Verify clean architecture files exist
ls -la mcp-server/docker-compose.clean.yml
ls -la mcp-server/Dockerfile.clean
```

### Step 3: Migrate Database

#### Option A: Data Migration (Recommended)

**Preserve existing data**:
```bash
cd mcp-server

# Start only the clean database
docker-compose -f docker-compose.clean.yml up -d postgres

# Wait for database to be ready
sleep 30

# Restore data to clean architecture database
psql -h localhost -U movies_user -p 5433 -d movies_mcp < ../legacy_backup_$(date +%Y%m%d).sql

# Verify data migration
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM movies;"
```

#### Option B: Fresh Start (Clean Install)

**Start with empty database**:
```bash
cd mcp-server

# Start clean architecture stack
docker-compose -f docker-compose.clean.yml up -d

# Database will be initialized automatically
# Seed with sample data (optional)
make db-seed
```

### Step 4: Update Configuration

**Environment Variables**:
```bash
# Create new configuration file
cat > .env.production << 'EOF'
# Database Configuration (Clean Architecture)
POSTGRES_USER=movies_user
POSTGRES_PASSWORD=movies_password
POSTGRES_DB=movies_mcp
POSTGRES_PORT=5433

# Server Configuration
LOG_LEVEL=info
SERVER_TIMEOUT=30s

# Image Processing
MAX_IMAGE_SIZE=5242880
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp
ENABLE_THUMBNAILS=true

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_PORT=8080
EOF
```

**Port Mapping Updates**:
```bash
# Update any scripts or configurations that reference old ports
sed -i 's/5432:/5433:/g' *.sh
sed -i 's/8080:/8081:/g' *.sh
sed -i 's/3000:/3001:/g' *.sh
```

### Step 5: Deploy Clean Architecture

**Docker Deployment**:
```bash
cd mcp-server

# Deploy full clean architecture stack
docker-compose -f docker-compose.clean.yml up -d --build

# Verify all services are running
docker-compose -f docker-compose.clean.yml ps

# Check service health
docker-compose -f docker-compose.clean.yml logs movies-server-clean
```

**Binary Deployment**:
```bash
cd mcp-server

# Build clean architecture binary
make build-clean

# Install binary
sudo cp build/movies-server-clean /opt/movies-mcp/bin/
sudo chmod +x /opt/movies-mcp/bin/movies-server-clean

# Update systemd service (if used)
sudo sed -i 's/movies-server/movies-server-clean/g' /etc/systemd/system/movies-mcp.service
sudo systemctl daemon-reload
sudo systemctl start movies-mcp
```

### Step 6: Update MCP Client Configuration

**Claude Desktop Configuration**:
```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-server-clean",
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5433/movies_mcp"
      }
    }
  }
}
```

**Key changes**:
- Binary name: `movies-server` → `movies-server-clean`
- Database port: `:5432/` → `:5433/` (Docker deployment)
- Health check port: `8080` → `8081`

## ✅ Post-Migration Verification

### 1. Service Health Checks

**Docker Services**:
```bash
# Check all services are healthy
docker-compose -f docker-compose.clean.yml ps

# Verify health endpoints
curl http://localhost:8081/health        # MCP Server
curl http://localhost:9091/-/healthy     # Prometheus
curl http://localhost:3001/api/health    # Grafana
```

**Database Connectivity**:
```bash
# Test database connection
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM movies;"

# Test search functionality
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT title FROM movies WHERE search_vector @@ plainto_tsquery('english', 'matrix');"
```

### 2. MCP Protocol Testing

**Manual Protocol Test**:
```bash
# Test MCP initialization
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | ./build/movies-server-clean

# Test tool listing
echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":2}' | ./build/movies-server-clean
```

**Claude Desktop Integration**:
1. Restart Claude Desktop
2. Open new conversation
3. Try: "List all movies in the database"
4. Try: "Add a test movie to verify functionality"

### 3. Data Integrity Verification

**Compare record counts**:
```bash
# Movies count
echo "Legacy:" && psql -h localhost -U movies_user -p 5432 -d movies_mcp -c "SELECT COUNT(*) FROM movies;" 2>/dev/null || echo "Legacy DB not accessible"
echo "Clean:" && psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM movies;"

# Actors count (if applicable)
echo "Actors - Clean:" && psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM actors;" 2>/dev/null
```

**Spot check specific records**:
```bash
# Check specific movie exists
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT title, year, rating FROM movies WHERE title ILIKE '%matrix%';"

# Check image data (if applicable)
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT id, title, length(poster_data) as image_size FROM movies WHERE poster_data IS NOT NULL LIMIT 5;"
```

### 4. Performance Comparison

**Response time testing**:
```bash
# Time a search query
time psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM movies WHERE search_vector @@ plainto_tsquery('english', 'action');"

# Monitor resource usage
docker stats movies-server-clean --no-stream
```

**Load testing** (optional):
```bash
# Install hey tool for load testing
go install github.com/rakyll/hey@latest

# Test health endpoint
hey -n 100 -c 10 http://localhost:8081/health
```

## 🔧 Troubleshooting Migration Issues

### Database Connection Issues

**Problem**: Can't connect to new database
```bash
# Check if database is running
docker-compose -f docker-compose.clean.yml logs postgres

# Verify port binding
netstat -tlnp | grep 5433

# Test connection manually
psql -h localhost -U movies_user -p 5433 -d movies_mcp
```

**Solution**:
```bash
# Restart database service
docker-compose -f docker-compose.clean.yml restart postgres

# Check database logs for errors
docker-compose -f docker-compose.clean.yml logs postgres | tail -50
```

### Port Conflicts

**Problem**: Port already in use
```bash
# Find what's using the port
lsof -i :5433
netstat -tlnp | grep 5433
```

**Solution**:
```bash
# Stop conflicting services
docker-compose down  # Stop legacy stack
docker-compose -f docker-compose.dev.yml down  # Stop dev stack

# Or change ports in docker-compose.clean.yml
ports:
  - "5434:5432"  # Use alternative port
```

### MCP Client Configuration Issues

**Problem**: Claude Desktop can't find server
```bash
# Verify binary exists and is executable
ls -la /path/to/movies-server-clean
./movies-server-clean --version
```

**Solution**:
```bash
# Use absolute path
which movies-server-clean

# Update Claude Desktop config with full path
{
  "mcpServers": {
    "movies": {
      "command": "/full/absolute/path/to/movies-server-clean",
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5433/movies_mcp"
      }
    }
  }
}
```

### Data Migration Issues

**Problem**: Data missing after migration
```bash
# Check if data was imported correctly
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "\dt"  # List tables
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT COUNT(*) FROM movies;"
```

**Solution**:
```bash
# Re-import from backup
psql -h localhost -U movies_user -p 5433 -d movies_mcp < legacy_backup_$(date +%Y%m%d).sql

# Or import specific data
pg_dump -h localhost -U movies_user -p 5432 movies_mcp --data-only --table=movies | psql -h localhost -U movies_user -p 5433 -d movies_mcp
```

### Performance Issues

**Problem**: Server running slowly
```bash
# Check resource usage
docker stats movies-server-clean --no-stream

# Check database connection pool
psql -h localhost -U movies_user -p 5433 -d movies_mcp -c "SELECT count(*) as connections, state FROM pg_stat_activity GROUP BY state;"
```

**Solution**:
```bash
# Tune connection pool settings
export DB_MAX_OPEN_CONNS=10
export DB_MAX_IDLE_CONNS=5

# Restart with new settings
docker-compose -f docker-compose.clean.yml restart movies-server-clean
```

## 🔄 Rollback Plan

If migration fails and you need to rollback:

### 1. Stop Clean Architecture Services

```bash
# Stop clean services
docker-compose -f docker-compose.clean.yml down

# Verify ports are freed
netstat -tlnp | grep -E "(5433|8081|9091|3001)"
```

### 2. Restore Legacy Services

```bash
# Restore legacy configuration
cp .env.legacy.backup .env
cp docker-compose.legacy.backup docker-compose.yml

# Start legacy services
docker-compose up -d

# Restore legacy volumes if needed
docker run --rm -v movies_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/legacy_postgres_data.tar.gz -C /data
```

### 3. Restore Claude Desktop Configuration

```bash
# Restore legacy MCP configuration
cp claude_config.legacy.backup ~/Library/Application\ Support/Claude/claude_desktop_config.json

# Restart Claude Desktop
```

## 📈 Post-Migration Optimization

### 1. Clean Up Legacy Resources

**After successful migration**:
```bash
# Remove legacy Docker images
docker image rm movies-mcp-server:legacy 2>/dev/null

# Remove legacy volumes (CAUTION: Only after confirming migration success)
docker volume rm movies_postgres_data 2>/dev/null
docker volume rm movies_redis_data 2>/dev/null

# Archive legacy backups
mkdir -p backups/legacy_$(date +%Y%m%d)
mv legacy_backup_*.sql backups/legacy_$(date +%Y%m%d)/
mv *.legacy.backup backups/legacy_$(date +%Y%m%d)/
```

### 2. Optimize Clean Architecture

**Database optimization**:
```sql
-- Connect to clean database
\c movies_mcp

-- Update table statistics
ANALYZE;

-- Rebuild indexes for optimal performance
REINDEX INDEX CONCURRENTLY movies_search_idx;
REINDEX INDEX CONCURRENTLY movies_year_idx;
```

**Docker resource limits**:
```yaml
# Add to docker-compose.clean.yml
services:
  movies-server-clean:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### 3. Enable Advanced Features

**Monitoring setup**:
```bash
# Access Grafana dashboards
open http://localhost:3001

# Import additional dashboards
curl -X POST http://admin:admin@localhost:3001/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @monitoring/grafana-dashboard.json
```

**Performance metrics**:
```bash
# Enable detailed metrics
export DB_METRICS_ENABLED=true
export REQUEST_METRICS_ENABLED=true

# Restart to apply changes
docker-compose -f docker-compose.clean.yml restart movies-server-clean
```

## ✅ Migration Success Checklist

Verify all items before considering migration complete:

### Technical Verification
- [ ] **Clean architecture services running** (`docker-compose ps` shows all healthy)
- [ ] **Database connection working** (can query movies table)
- [ ] **MCP protocol responding** (manual protocol test passes)
- [ ] **Claude Desktop integration** (can execute MCP tools)
- [ ] **All data migrated** (record counts match)
- [ ] **Images preserved** (poster data accessible)
- [ ] **Search functionality** (full-text search working)
- [ ] **Health endpoints responding** (8081/health returns OK)

### Monitoring & Observability
- [ ] **Prometheus collecting metrics** (9091 accessible)
- [ ] **Grafana dashboards loaded** (3001 shows data)
- [ ] **pgAdmin database access** (5051 connects to database)
- [ ] **Structured logging active** (JSON logs in container)

### Performance & Security
- [ ] **Response times acceptable** (sub-second for basic queries)
- [ ] **Resource usage reasonable** (memory < 1GB, CPU < 50%)
- [ ] **Security features enabled** (non-root containers, restricted permissions)
- [ ] **Backup strategy confirmed** (automated or manual process)

### Documentation & Cleanup
- [ ] **Documentation updated** (internal docs reflect new setup)
- [ ] **Team members notified** (if applicable)
- [ ] **Legacy resources cleaned** (old containers/volumes removed)
- [ ] **Monitoring alerts configured** (if using external monitoring)

## 🎉 Congratulations!

You've successfully migrated from Legacy to Clean Architecture! Your Movies MCP Server now benefits from:

- 🏗️ **Better Architecture**: Maintainable, testable, and extensible code
- 🚀 **Improved Performance**: Optimized queries and efficient resource usage
- 📊 **Better Observability**: Comprehensive monitoring and logging
- 🔒 **Enhanced Security**: Container security and input validation
- 📈 **Future-Proof**: Active development and feature additions

## 📚 Next Steps

After successful migration:

1. **[User Guide](../guides/user-guide.md)** - Explore new features and capabilities
2. **[Performance Tuning](../reference/performance.md)** - Optimize for your workload
3. **[Monitoring Setup](../reference/monitoring.md)** - Configure advanced monitoring
4. **[Contributing Guide](../development/contributing.md)** - Help improve the project

Need help with migration? Check the **[FAQ](faq.md)** or open a **[GitHub Discussion](https://github.com/francknouama/movies-mcp-server/discussions)** for community support! 🚀