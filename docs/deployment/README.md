# 🚀 Production Deployment Guide

Deploy the Movies MCP Server to production with confidence using Docker, cloud platforms, or traditional server setups.

## 📋 Overview

The Movies MCP Server is a **Model Context Protocol (MCP)** server that communicates via **stdin/stdout**, not HTTP. It's designed to be integrated into MCP-compatible clients (like Claude Desktop) rather than deployed as a standalone web service.

### 🏗️ Architecture Characteristics
- **Communication**: stdin/stdout using MCP protocol
- **Architecture**: Clean Architecture with Domain-Driven Design
- **Database**: PostgreSQL with full-text search capabilities
- **Image Support**: Binary image storage with base64 encoding
- **Monitoring**: Optional HTTP endpoints for health checks and metrics only

## 🎯 Deployment Options

| Method | Complexity | Scalability | Best For |
|--------|------------|-------------|----------|
| **[Docker Compose](#-docker-compose-deployment)** | ⭐ Low | ⭐⭐ Medium | Single server, quick start |
| **[Cloud Platforms](#-cloud-platform-deployment)** | ⭐⭐ Medium | ⭐⭐⭐ High | Production, scaling needs |
| **[Kubernetes](#-kubernetes-deployment)** | ⭐⭐⭐ High | ⭐⭐⭐ High | Enterprise, complex orchestration |
| **[Binary Deployment](#-binary-deployment)** | ⭐⭐ Medium | ⭐ Low | Legacy infrastructure |

## 🔧 Prerequisites

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

## ⚙️ Environment Configuration

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

## 🐳 Docker Compose Deployment

**Best for**: Quick production deployment, single server setups

### Clean Architecture Stack (Recommended)

```bash
# Clone repository
git clone <repository-url>
cd movies-mcp-server/mcp-server

# Start the clean architecture stack
docker-compose -f docker-compose.clean.yml up -d

# Verify deployment
docker-compose -f docker-compose.clean.yml ps
docker-compose -f docker-compose.clean.yml logs movies-server-clean
```

### Full Production Stack

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

For detailed Docker configuration, see **[Docker Deployment Guide](docker.md)**.

## ☁️ Cloud Platform Deployment

### AWS ECS Deployment

```bash
# Build and tag image
docker build -f Dockerfile.clean -t movies-mcp-server:latest .

# Push to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker tag movies-mcp-server:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/movies-mcp-server:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/movies-mcp-server:latest
```

### Google Cloud Run

```bash
# Build and push to Container Registry
gcloud builds submit --tag gcr.io/PROJECT-ID/movies-mcp-server

# Deploy to Cloud Run
gcloud run deploy movies-mcp-server \
  --image gcr.io/PROJECT-ID/movies-mcp-server \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8081 \
  --memory 1Gi \
  --cpu 1 \
  --set-env-vars DATABASE_URL="postgres://user:pass@/movies_mcp?host=/cloudsql/PROJECT-ID:REGION:INSTANCE-ID"
```

### Azure Container Instances

```bash
# Create resource group
az group create --name movies-mcp-rg --location eastus

# Create container instance
az container create \
  --resource-group movies-mcp-rg \
  --name movies-mcp-server \
  --image your-registry/movies-mcp-server:latest \
  --cpu 1 \
  --memory 1 \
  --ports 8081 \
  --environment-variables DATABASE_URL="postgres://user:pass@host:5432/movies_mcp"
```

## ⚙️ Kubernetes Deployment

**Best for**: Enterprise environments, high scalability

### Basic Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: movies-mcp-server
  namespace: movies-mcp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: movies-mcp-server
  template:
    metadata:
      labels:
        app: movies-mcp-server
    spec:
      containers:
      - name: movies-server
        image: movies-mcp-server:latest
        ports:
        - containerPort: 8081
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: movies-mcp-secrets
              key: DATABASE_URL
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
```

### Deploy to Kubernetes

```bash
# Apply configurations
kubectl apply -f namespace.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Verify deployment
kubectl get pods -n movies-mcp
kubectl logs -f deployment/movies-mcp-server -n movies-mcp
```

## 🔧 Binary Deployment

**Best for**: Existing infrastructure, bare metal servers

### Build from Source

```bash
# Clone repository
git clone <repository-url>
cd movies-mcp-server/mcp-server

# Build clean architecture version
make build-clean

# Set up database
make db-setup
make db-migrate

# Run the server
./build/movies-server-clean
```

### Cross-Platform Builds

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

## 🗄️ Database Setup

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

## 🔌 MCP Client Integration

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

## 📊 Monitoring and Logging

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

## 🔒 Security Considerations

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

## 🔧 Troubleshooting

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

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=debug
./movies-server-clean
```

### Backup and Recovery

```bash
# Create backup
pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > movies_backup.sql

# Restore backup
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < movies_backup.sql

# Automated backups
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | gzip > /backups/movies_$(date +\%Y\%m\%d).sql.gz
```

## 📈 Scaling Considerations

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

## 🚀 Quick Start Commands

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

## 📚 Next Steps

After successful deployment:

1. **[Docker Configuration](docker.md)** - Detailed Docker setup and optimization
2. **[Troubleshooting Guide](../reference/troubleshooting.md)** - Common issues and solutions
3. **[User Guide](../guides/user-guide.md)** - Using the MCP server effectively
4. **[FAQ](../appendices/faq.md)** - Frequently asked questions

Your Movies MCP Server is now production-ready! 🎉