# ❓ Frequently Asked Questions

Answers to commonly asked questions about the Movies MCP Server, from installation to advanced usage.

## 🚀 Getting Started

### What is the Movies MCP Server?

The Movies MCP Server is a **Model Context Protocol (MCP)** server that provides movie database management capabilities to MCP-compatible clients like Claude Desktop. It allows you to:

- 🎬 **Manage Movies**: Add, update, delete, and search movies
- 👥 **Handle Actors**: Manage actor information and relationships
- 🖼️ **Store Images**: Upload and manage movie posters
- 🔍 **Advanced Search**: Full-text search with PostgreSQL
- 📊 **Get Statistics**: Database insights and analytics

### Why use MCP instead of a REST API?

MCP (Model Context Protocol) offers several advantages:

- **🔒 Secure**: No network exposure, communicates via stdin/stdout
- **🎯 AI-Optimized**: Designed specifically for AI assistant integration
- **📝 Self-Describing**: Tools and resources are automatically discoverable
- **🔄 Stateful**: Maintains context across interactions
- **⚡ Efficient**: Direct process communication, no HTTP overhead

### Is this production-ready?

**Yes!** The Movies MCP Server is built with production considerations:

- ✅ **Clean Architecture**: Domain-driven design with proper separation
- ✅ **Comprehensive Testing**: Unit, integration, and ATDD tests
- ✅ **Monitoring**: Prometheus metrics and Grafana dashboards
- ✅ **Security**: Input validation, SQL injection prevention
- ✅ **Containerized**: Docker deployment with health checks
- ✅ **Migration System**: Database schema versioning

## 🔧 Installation & Setup

### Which architecture should I use - Legacy or Clean?

**Use Clean Architecture** (recommended):

| Aspect | Clean Architecture | Legacy |
|--------|-------------------|--------|
| **Code Quality** | ✅ High maintainability | ⚠️ Tightly coupled |
| **Testing** | ✅ Comprehensive test suite | ⚠️ Limited tests |
| **Performance** | ✅ Optimized queries | ⚠️ Less efficient |
| **Documentation** | ✅ Well documented | ⚠️ Minimal docs |
| **Future Support** | ✅ Active development | ❌ Maintenance mode |

### Can I use SQLite instead of PostgreSQL?

**Currently no.** The server is optimized for PostgreSQL features:

- **Full-text search** with `pg_trgm` extension
- **Binary data storage** for images
- **Advanced indexing** for performance
- **Connection pooling** for scalability

PostgreSQL is lightweight and can run in Docker, making setup straightforward.

### Do I need Redis?

**Redis is optional** but recommended for:

- **🚀 Performance**: Caching frequently accessed data
- **📊 Session Storage**: Maintaining request context
- **🔄 Rate Limiting**: Preventing abuse in production

For development or small installations, Redis can be disabled:

```bash
export REDIS_ENABLED=false
```

### How much resources does it need?

**Minimum Requirements**:
- **CPU**: 1 core
- **Memory**: 512MB RAM
- **Storage**: 2GB disk space

**Recommended Production**:
- **CPU**: 2+ cores
- **Memory**: 1GB+ RAM  
- **Storage**: 10GB+ SSD
- **Database**: Additional 1GB+ for PostgreSQL

## 🐳 Docker & Deployment

### Which Docker Compose file should I use?

Choose based on your needs:

| File | Purpose | When to Use |
|------|---------|-------------|
| `docker-compose.clean.yml` | **Production** | Full stack with monitoring |
| `docker-compose.dev.yml` | **Development** | Local development databases |
| `docker-compose.yml` | **Legacy** | Original architecture (deprecated) |

### Why are there different port numbers?

Different environments use different ports to avoid conflicts:

| Environment | MCP Health | PostgreSQL | Grafana | Purpose |
|-------------|------------|------------|---------|---------|
| **Production** | 8081 | 5433 | 3001 | Clean architecture |
| **Development** | - | 5434/5435 | - | Dev/test databases |
| **Legacy** | 8080 | 5432 | 3000 | Original architecture |

### Can I deploy to cloud platforms?

**Absolutely!** The server supports:

- **☁️ AWS**: ECS, Fargate, Lambda
- **🌐 Google Cloud**: Cloud Run, GKE
- **🔷 Azure**: Container Instances, AKS
- **🐋 Docker**: Any Docker-compatible platform

See the **[Deployment Guide](../deployment/README.md)** for detailed instructions.

### How do I update to a new version?

**Docker Deployment**:
```bash
# Pull latest images
docker-compose -f docker-compose.clean.yml pull

# Restart with new version
docker-compose -f docker-compose.clean.yml up -d
```

**Binary Deployment**:
```bash
# Backup database first
pg_dump movies_mcp > backup.sql

# Build new version
git pull origin main
make build-clean

# Replace binary
sudo systemctl stop movies-mcp
sudo cp build/movies-server-clean /opt/movies-mcp/bin/
sudo systemctl start movies-mcp
```

## 🔌 MCP Integration

### How do I configure Claude Desktop?

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-server-clean",
      "env": {
        "DATABASE_URL": "postgres://user:pass@localhost:5432/movies_mcp"
      }
    }
  }
}
```

### Can I use it with other MCP clients?

**Yes!** The server implements the standard MCP protocol and works with any compatible client:

- **Claude Desktop** (most popular)
- **Continue.dev** (VS Code extension)
- **Custom implementations** (following MCP spec)
- **Other AI assistants** (as they add MCP support)

### What tools are available?

**Movie Management**:
- `get_movie` - Retrieve movie details
- `add_movie` - Add new movies  
- `update_movie` - Update existing movies
- `delete_movie` - Remove movies
- `search_movies` - Full-text search
- `list_top_movies` - Get top-rated movies

**Actor Management**:
- `get_actor` - Retrieve actor details
- `add_actor` - Add new actors
- `update_actor` - Update actor information
- `delete_actor` - Remove actors
- `search_actors` - Find actors

**Advanced Operations**:
- `add_movie_with_actors` - Create movie with cast
- `get_movie_with_cast` - Retrieve complete movie info
- `search_movies_by_actor` - Find movies by actor

### What resources are provided?

**Database Resources**:
- `movies://database/all` - All movies (JSON)
- `movies://database/stats` - Database statistics
- `movies://database/schema` - Schema information

**Image Resources**:
- `movies://posters/{id}` - Movie poster images (base64)
- `movies://posters/collection` - Poster gallery
- `movies://images/stats` - Image storage statistics

## 🗄️ Database & Data

### How do I import existing movie data?

**Using seed scripts**:
```bash
# Use provided sample data
make db-seed

# Import custom data
psql -h localhost -U movies_user -d movies_mcp < your_data.sql
```

**Using MCP tools** (via Claude Desktop):
```
Add these movies to the database:
1. The Matrix (1999) - Directed by Wachowski Sisters
2. Inception (2010) - Directed by Christopher Nolan
```

### Can I backup my data?

**Database backup**:
```bash
# Create backup
pg_dump -h localhost -U movies_user movies_mcp > movies_backup.sql

# Restore backup  
psql -h localhost -U movies_user -d movies_mcp < movies_backup.sql
```

**Full system backup** (Docker):
```bash
# Backup volumes
docker run --rm -v movies_postgres_data_clean:/data -v $(pwd):/backup alpine tar czf /backup/movies_data.tar.gz -C /data .

# Restore volumes
docker run --rm -v movies_postgres_data_clean:/data -v $(pwd):/backup alpine tar xzf /backup/movies_data.tar.gz -C /data
```

### How do I handle image storage?

**Image upload via MCP**:
- Images are stored as **binary data** in PostgreSQL
- **Base64 encoding** for MCP transport
- **MIME type validation** (JPEG, PNG, WebP)
- **Size limits** (configurable, default 5MB)

**Storage optimization**:
```bash
# Enable thumbnail generation
export ENABLE_THUMBNAILS=true
export THUMBNAIL_SIZE=200x200

# Adjust size limits
export MAX_IMAGE_SIZE=10485760  # 10MB
```

### Can I use external image storage?

**Currently no**, but planned for future versions:
- **S3-compatible storage** (AWS S3, MinIO)
- **Cloud storage** (Google Cloud Storage)
- **CDN integration** for performance

Track progress in the GitHub repository issues.

## 🔍 Troubleshooting

### Server won't start - "database connection failed"

**Check database status**:
```bash
# Docker deployment
docker-compose -f docker-compose.clean.yml logs postgres

# Manual connection test
psql -h localhost -U movies_user -d movies_mcp -c "SELECT 1;"
```

**Common fixes**:
1. **Database not ready**: Wait for PostgreSQL to fully start
2. **Wrong credentials**: Check `DATABASE_URL` environment variable
3. **Network issues**: Verify Docker network configuration
4. **Port conflicts**: Ensure PostgreSQL port is not in use

### Claude Desktop says "Server not found"

**Check binary path**:
```bash
# Verify binary exists and is executable
ls -la /path/to/movies-server-clean
./movies-server-clean --version
```

**Check configuration**:
```json
{
  "mcpServers": {
    "movies": {
      "command": "/absolute/path/to/movies-server-clean",
      "env": {
        "DATABASE_URL": "postgres://user:pass@localhost:5432/movies_mcp"
      }
    }
  }
}
```

**Debug MCP communication**:
```bash
# Test MCP protocol manually
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | ./movies-server-clean
```

### High memory usage

**Monitor resource usage**:
```bash
# Docker
docker stats movies-server-clean

# System
top -p $(pgrep movies-server-clean)
```

**Optimization strategies**:
1. **Reduce connection pool size**:
   ```bash
   export DB_MAX_OPEN_CONNS=10
   export DB_MAX_IDLE_CONNS=2
   ```

2. **Limit image processing**:
   ```bash
   export MAX_IMAGE_SIZE=2097152  # 2MB
   export ENABLE_THUMBNAILS=false
   ```

3. **Tune Go garbage collector**:
   ```bash
   export GOGC=100
   export GOMEMLIMIT=512MiB
   ```

### Migrations fail to run

**Check migration status**:
```bash
psql -h localhost -U movies_user -d movies_mcp -c "SELECT * FROM schema_migrations;"
```

**Force migration reset** (⚠️ **CAUTION: This deletes data**):
```bash
# Drop and recreate database
make db-reset db-migrate

# Or manually
psql -h localhost -U postgres -c "DROP DATABASE movies_mcp;"
psql -h localhost -U postgres -c "CREATE DATABASE movies_mcp OWNER movies_user;"
```

**Run migrations manually**:
```bash
./movies-server-clean --migrations ./migrations
```

## 🎯 Usage & Best Practices

### How should I structure my movie data?

**Recommended approach**:
```bash
# Add movie with complete information
{
  "title": "The Matrix",
  "year": 1999,
  "genre": "Science Fiction",
  "director": "Wachowski Sisters",
  "description": "A computer hacker learns from mysterious rebels about the true nature of his reality.",
  "rating": 8.7,
  "image_url": "https://example.com/matrix_poster.jpg"
}
```

**Best practices**:
- ✅ **Use consistent genres** (e.g., "Science Fiction" not "Sci-Fi")
- ✅ **Include year for disambiguation** (multiple versions exist)
- ✅ **Add detailed descriptions** for better search
- ✅ **Use decimal ratings** (0.0-10.0 scale)
- ✅ **Provide high-quality poster images**

### How do I optimize search performance?

**Database optimization**:
```sql
-- Check search index usage
EXPLAIN (ANALYZE, BUFFERS) 
SELECT * FROM movies 
WHERE search_vector @@ plainto_tsquery('english', 'matrix');

-- Rebuild search indexes if needed
REINDEX INDEX movies_search_idx;
```

**Query optimization**:
- Use **specific terms** rather than broad searches
- **Combine filters** (year + genre + rating) for better results
- **Use actor names** in movie searches for cast-based queries

### Can I extend the server with custom tools?

**Currently no**, but the architecture supports extensions:

**Planned features**:
- 🔌 **Plugin system** for custom tools
- 📦 **Extension packages** via Go modules
- 🎨 **Custom schemas** for specialized domains
- 🔗 **API integrations** (TMDB, IMDB)

**Contributing**:
- Check the **[Contributing Guide](../development/contributing.md)**
- Submit **feature requests** via GitHub issues
- Join the **community discussions**

## 📊 Performance & Scaling

### How many movies can it handle?

**Performance benchmarks**:
- **Database**: 1M+ movies with good performance
- **Search**: Sub-second response for complex queries
- **Images**: Optimized binary storage
- **Concurrent Users**: 100+ simultaneous MCP clients

**Scaling considerations**:
- **Vertical scaling**: Increase CPU/memory for single instance
- **Database optimization**: Connection pooling, query optimization
- **Image optimization**: Enable thumbnails, consider CDN

### How do I monitor performance?

**Built-in monitoring**:
```bash
# Health checks
curl http://localhost:8081/health

# Prometheus metrics
curl http://localhost:8081/metrics

# Grafana dashboards
open http://localhost:3001
```

**Key metrics to watch**:
- **Request rate**: `movies_requests_total`
- **Response time**: `movies_request_duration_seconds`
- **Database connections**: `movies_db_connections_active`
- **Error rate**: `movies_errors_total`

### Can I run multiple instances?

**MCP deployment pattern**:
- Each **MCP client** runs its own server instance
- **No shared state** between instances required
- **Database handles** concurrent connections
- **Stateless design** enables horizontal scaling

**Load balancing** (if needed):
- Use **database connection pooling**
- **Shared Redis** for caching (optional)
- **Separate read replicas** for heavy read workloads

## 🔒 Security

### Is it secure for production use?

**Yes, with proper configuration**:

**Built-in security**:
- ✅ **No network exposure** (stdin/stdout only)
- ✅ **Input validation** at all layers
- ✅ **SQL injection prevention** (prepared statements)
- ✅ **Type safety** (Go's type system)
- ✅ **Container security** (non-root, distroless)

**Production checklist**:
- 🔐 **Secure database credentials**
- 🛡️ **Enable SSL/TLS** for database connections
- 🔒 **Restrict file permissions** on binary and config
- 📝 **Enable audit logging**
- 🔄 **Regular security updates**

### How do I secure the database?

**PostgreSQL security**:
```bash
# Use strong passwords
export DB_PASSWORD="$(openssl rand -base64 32)"

# Enable SSL connections
export DB_SSLMODE=require

# Restrict network access
# In pg_hba.conf:
host    movies_mcp    movies_user    127.0.0.1/32    md5
```

**Container security**:
```yaml
# docker-compose.clean.yml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password
```

### How do I handle sensitive data?

**Environment variables**:
```bash
# Use secrets management
export DATABASE_URL="$(cat /run/secrets/db_url)"

# Or use .env files with restricted permissions
chmod 600 .env.production
```

**Image data**:
- Images are stored as **binary data** in database
- **No file system exposure**
- **MIME type validation** prevents malicious uploads
- **Size limits** prevent abuse

## 📞 Support & Community

### Where can I get help?

**Documentation**:
- 📖 **[User Guide](../guides/user-guide.md)** - Comprehensive usage guide
- 🔧 **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
- 🐳 **[Docker Guide](../deployment/docker.md)** - Docker deployment details

**Community**:
- 💬 **GitHub Discussions** - Questions and community support
- 🐛 **GitHub Issues** - Bug reports and feature requests
- 📧 **Maintainer Contact** - For complex issues

### How do I report bugs?

**Before reporting**:
1. Check **existing issues** on GitHub
2. Review **troubleshooting guide**
3. Test with **latest version**

**Bug report should include**:
- 🖥️ **Environment details** (OS, Docker version, etc.)
- 🔄 **Steps to reproduce** the issue
- 📋 **Expected vs actual behavior**
- 📝 **Relevant logs** (with sensitive data removed)
- 🏗️ **Architecture version** (Clean vs Legacy)

### How do I contribute?

**Ways to contribute**:
- 🐛 **Report bugs** and suggest improvements
- 📝 **Improve documentation** and examples
- 🧪 **Write tests** for better coverage
- ⭐ **Add features** following clean architecture
- 🎨 **Design improvements** for user experience

**Getting started**:
1. Read the **[Contributing Guide](../development/contributing.md)**
2. Check **"good first issue"** labels
3. Join **community discussions**
4. Submit **pull requests** with tests

### What's the roadmap?

**Upcoming features**:
- 🔌 **Plugin system** for extensibility
- ☁️ **Cloud storage** integration
- 🔍 **Advanced search** with AI embeddings
- 📊 **Analytics dashboard** for insights
- 🌐 **Multi-language** support
- 🔗 **External APIs** (TMDB, IMDB)

**Track progress**:
- 📋 **GitHub Projects** - Development roadmap
- 🏷️ **Release tags** - Version history
- 📰 **Release notes** - Feature announcements

---

## 💡 Still Have Questions?

Can't find what you're looking for? Here are more resources:

- 📖 **[User Guide](../guides/user-guide.md)** - Detailed usage instructions
- 🔧 **[Troubleshooting Guide](../reference/troubleshooting.md)** - Solve common problems
- 🐳 **[Docker Deployment](../deployment/docker.md)** - Container deployment details
- 🏗️ **[Migration Guide](migration-guide.md)** - Upgrade between architectures

**Still stuck?** Open a **[GitHub Discussion](https://github.com/francknouama/movies-mcp-server/discussions)** and the community will help! 🚀