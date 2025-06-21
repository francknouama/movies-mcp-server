# Movies MCP Server - Deployment Guide

This guide covers deploying the Movies MCP Server, a Model Context Protocol (MCP) server for movie database management built with Clean Architecture and Domain-Driven Design.

## Table of Contents

- [MCP Server Overview](#mcp-server-overview)
- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Deployment Methods](#deployment-methods)
- [Database Setup](#database-setup)
- [MCP Client Integration](#mcp-client-integration)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

## MCP Server Overview

The Movies MCP Server is a **Model Context Protocol (MCP)** server that communicates via **stdin/stdout**, not HTTP. It's designed to be integrated into MCP-compatible clients (like Claude Desktop) rather than deployed as a standalone web service.

### Key Characteristics

- **Communication**: stdin/stdout using MCP protocol
- **Architecture**: Clean Architecture with Domain-Driven Design
- **Database**: PostgreSQL with full-text search capabilities
- **Image Support**: Binary image storage with base64 encoding
- **Monitoring**: Optional HTTP endpoints for health checks and metrics only

## Prerequisites

### System Requirements
- **CPU**: 2+ cores recommended
- **Memory**: 1GB+ RAM
- **Storage**: 10GB+ available disk space
- **Network**: Outbound internet access for dependencies

### Software Dependencies
- **Docker**: 20.10+ (for containerized deployment)
- **PostgreSQL**: 13+ (database)
- **Redis**: 6+ (caching, optional)
- **Git**: 2.30+ (for source code)

## Environment Configuration

### Required Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=movies_mcp
DB_USER=movies_user
DB_PASSWORD=movies_password
DB_SSLMODE=disable

# Alternative: Single connection string
DATABASE_URL=postgres://movies_user:movies_password@localhost:5432/movies_mcp?sslmode=disable

# Server Configuration
LOG_LEVEL=info
SERVER_TIMEOUT=30s
```

### Optional Environment Variables

```bash
# Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=1h

# Image Processing
MAX_IMAGE_SIZE=5242880  # 5MB
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp
ENABLE_THUMBNAILS=true
THUMBNAIL_SIZE=200x200

# Monitoring (HTTP endpoints)
METRICS_ENABLED=true
METRICS_INTERVAL=30s
HEALTH_CHECK_PORT=8080
```

## Deployment Methods

### 1. Docker Deployment (Recommended)

#### Using Docker Compose - Clean Architecture

```bash
# Clone the repository
git clone <repository-url>
cd movies-mcp-server

# Start the clean architecture stack
docker-compose -f docker-compose.clean.yml up -d

# Verify deployment
docker-compose -f docker-compose.clean.yml ps
docker-compose -f docker-compose.clean.yml logs movies-server-clean
```

#### Using Docker Compose - Full Production Stack

```bash
# Start with monitoring stack
docker-compose up -d

# This includes:
# - PostgreSQL database
# - Movies MCP Server
# - Redis (caching)
# - Prometheus (metrics)
# - Grafana (visualization)
# - pgAdmin (database management)
```

#### Using Production Dockerfile

```bash
# Build production image
docker build -f Dockerfile.production -t movies-mcp-server:latest .

# Run with external database
docker run -d \
  --name movies-mcp-server \
  --restart unless-stopped \
  -e DATABASE_URL=postgres://user:pass@host:5432/movies_mcp \
  -e LOG_LEVEL=info \
  movies-mcp-server:latest
```

### 2. Binary Deployment

#### Build from Source

```bash
# Clone repository
git clone <repository-url>
cd movies-mcp-server

# Build clean architecture version
make build-clean

# Set up database
make db-setup
make db-migrate

# Run the server
./build/movies-server-clean
```

#### Cross-Platform Builds

```bash
# Build for multiple platforms
make build-all

# Output:
# - build/movies-server-clean-linux-amd64
# - build/movies-server-clean-linux-arm64  
# - build/movies-server-clean-darwin-amd64
# - build/movies-server-clean-darwin-arm64
# - build/movies-server-clean-windows-amd64.exe
```

### 3. Local Development Deployment

**Note**: MCP servers are not deployed to the cloud as standalone services. They run locally as processes controlled by MCP clients.

```bash
# For local development with Docker
make docker-compose-up-dev

# For local binary development
make build-clean test

# For production-like local environment
make docker-compose-up-clean
```

## Database Setup

### PostgreSQL Configuration

```sql
-- Create database and user
CREATE DATABASE movies_mcp;
CREATE USER movies_user WITH ENCRYPTED PASSWORD 'movies_password';
GRANT ALL PRIVILEGES ON DATABASE movies_mcp TO movies_user;

-- Enable required extensions
\c movies_mcp;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
```

### Migration System

The server includes a custom migration tool:

```bash
# Automatic migrations (default)
./movies-server-clean

# Skip migrations
./movies-server-clean --skip-migrations

# Run migrations manually
./movies-server-clean --migrations ./migrations
```

### Database Schema

The schema includes:
- **movies** table with full-text search indexes
- **actors** table with biography support
- **movie_actors** many-to-many relationships
- **Binary image storage** with MIME type support
- **Automatic timestamps** and audit triggers

## MCP Client Integration

### Claude Desktop Configuration

Add to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/path/to/movies-server-clean",
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5432/movies_mcp"
      }
    }
  }
}
```

### Alternative MCP Clients

For other MCP clients, ensure they can:
1. Execute the binary with proper environment variables
2. Communicate via stdin/stdout using MCP protocol
3. Handle MCP resources (for image data)

### Available MCP Tools

The server provides these tools:
- `get_movie` - Retrieve movie details
- `add_movie` - Add new movies
- `update_movie` - Update existing movies
- `delete_movie` - Remove movies
- `search_movies` - Full-text search
- `list_top_movies` - Get top-rated movies
- Actor management tools
- Compound operations for complex queries

### Available MCP Resources

- `movies://database/all` - All movies (JSON)
- `movies://database/stats` - Database statistics
- `movies://posters/{id}` - Movie poster images (base64)
- `movies://posters/collection` - Poster gallery

## Monitoring and Logging

### Health Check Endpoints

The server optionally provides HTTP endpoints:

```bash
# Health checks
curl http://localhost:8080/health        # Overall health
curl http://localhost:8080/health/db     # Database connectivity
curl http://localhost:8080/ready         # Readiness probe

# Metrics (Prometheus format)
curl http://localhost:8080/metrics
```

### Prometheus Metrics

Available metrics:
- `movies_requests_total` - Total requests by tool
- `movies_request_duration_seconds` - Request duration histogram
- `movies_db_connections_active` - Active database connections
- `movies_db_operations_total` - Database operations count

### Grafana Dashboard

Import the provided dashboard:

```bash
# Start Grafana with Docker Compose
docker-compose up -d grafana

# Access at http://localhost:3000 (admin/admin)
# Import dashboard from monitoring/grafana-dashboard.json
```

### Structured Logging

The server uses structured JSON logging:

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "message": "Tool executed successfully",
  "tool": "get_movie",
  "duration": "15ms",
  "movie_id": 123
}
```

### Log Aggregation

#### ELK Stack
```yaml
# docker-compose.override.yml
version: '3.8'
services:
  filebeat:
    image: elastic/filebeat:7.17.0
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
```

#### Fluentd
```yaml
fluentd:
  image: fluent/fluentd:v1.14
  volumes:
    - ./fluentd.conf:/fluentd/etc/fluent.conf
    - /var/log:/var/log:ro
```

## Security Considerations

### MCP Security Model

- **Process Isolation**: MCP servers run as separate processes
- **No Network Exposure**: Communication only via stdin/stdout
- **Client-Controlled**: MCP client manages server lifecycle

### Application Security

- **Input Validation**: All inputs validated at domain boundaries
- **SQL Injection Prevention**: Prepared statements only
- **Type Safety**: Domain models prevent invalid states
- **Image Size Limits**: Configurable max image sizes
- **MIME Type Validation**: Only allowed image types

### Container Security

- **Non-root Execution**: Containers run as non-root user
- **Distroless Base Images**: Minimal attack surface
- **Security Scanning**: Images scanned for vulnerabilities
- **Read-only Root Filesystem**: Immutable containers

### Database Security

- **Connection Encryption**: SSL/TLS for database connections
- **Credential Management**: Environment-based secrets
- **Connection Pooling**: Prevents connection exhaustion
- **Query Timeout**: Prevents long-running queries

## Troubleshooting

### Common Issues

#### Database Connection Issues

```bash
# Test database connectivity
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;"

# Check server logs
docker logs movies-mcp-server 2>&1 | grep -E "(ERROR|FATAL)"

# Verify environment variables
env | grep -E "(DB_|DATABASE_)"
```

#### MCP Communication Issues

```bash
# Test MCP protocol manually
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | ./movies-server-clean

# Check client configuration
# Verify binary path and environment variables in MCP client config
```

#### Migration Issues

```bash
# Check migration status
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT * FROM schema_migrations;"

# Run migrations manually
./movies-server-clean --migrations ./migrations

# Reset database (development only)
make db-reset db-migrate
```

#### Performance Issues

```bash
# Check resource usage
docker stats movies-mcp-server

# Analyze slow queries
docker exec -it postgres psql -U movies_user -d movies_mcp -c "
SELECT query, mean_exec_time, calls 
FROM pg_stat_statements 
ORDER BY mean_exec_time DESC 
LIMIT 10;"
```

#### Log Analysis
```bash
# View recent logs
docker logs --tail 100 movies-mcp-server

# Filter error logs
docker logs movies-mcp-server 2>&1 | grep ERROR
```

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=debug
./movies-server-clean
```


### Backup and Recovery

### Backup and Recovery

```bash
# Create backup
pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > movies_backup.sql

# Restore backup
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < movies_backup.sql

# Automated backups
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | gzip > /backups/movies_$(date +\%Y\%m\%d).sql.gz
```

## Scaling Considerations

### Horizontal Scaling

MCP servers are typically deployed per client instance:
- Each MCP client runs its own server process
- No shared state between instances
- Database connection pooling handles concurrent access

### Vertical Scaling

- Monitor memory usage (especially for image processing)
- Adjust database connection pool sizes
- Optimize PostgreSQL configuration for workload

### Database Scaling

```sql
-- Optimize for read-heavy workloads
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '16MB';
ALTER SYSTEM SET maintenance_work_mem = '256MB';
```

## Quick Start Commands

```bash
# Docker development setup
make docker-compose-up-dev

# Build and test locally
make build-clean test

# Production deployment
make docker-compose-up-clean

# Database operations
make db-setup db-migrate db-seed

# Monitoring
make monitoring-up
```

This deployment guide provides comprehensive instructions for deploying the Movies MCP Server with proper security, monitoring, and scaling considerations tailored to the MCP protocol architecture.