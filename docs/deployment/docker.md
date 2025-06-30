# 🐳 Docker Deployment Guide

Comprehensive guide for deploying the Movies MCP Server using Docker with Clean Architecture implementation.

## 📋 Version Requirements

- **Go**: 1.24.4
- **PostgreSQL**: 17
- **Alpine Linux**: Latest (for minimal image size)
- **Redis**: 7-alpine
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

## 🏗️ Architecture Overview

The Docker configuration provides multiple environments tailored for different deployment scenarios:

| Environment | Purpose | Best For |
|-------------|---------|----------|
| **Production** | `Dockerfile.clean` + `docker-compose.clean.yml` | Production deployment |
| **Development** | `docker-compose.dev.yml` | Local development |
| **Legacy** | `Dockerfile` + `docker-compose.yml` | Original architecture |

## 🚀 Quick Start

### Development Environment

Perfect for local development with minimal overhead:

```bash
# Start development databases
docker-compose -f docker-compose.dev.yml up -d

# Set environment variables
export TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5435/movies_mcp_test?sslmode=disable"

# Run the application locally
cd mcp-server
go run cmd/server/main.go

# Run tests
go test ./...
```

### Production Environment

Full production stack with monitoring and observability:

```bash
# Build and start all services
docker-compose -f docker-compose.clean.yml up --build

# Or run in background
docker-compose -f docker-compose.clean.yml up -d --build

# View logs
docker-compose -f docker-compose.clean.yml logs -f movies-server-clean
```

## 🐳 Docker Files Detailed

### 1. Dockerfile.clean (Recommended)

**Purpose**: Production-ready multi-stage build for clean architecture

**Key Features**:
- **Multi-stage build** for minimal image size
- **Go 1.24.4** for building
- Builds both `movies-server` and `migrate` binaries
- **Distroless base image** for security
- **Non-root user** execution
- Built-in **health checks**

**Build Commands**:
```bash
# Standard build
docker build -f Dockerfile.clean -t movies-mcp-server:clean .

# Build with version information
docker build -f Dockerfile.clean \
  --build-arg VERSION=0.2.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t movies-mcp-server:0.2.0 .
```

### 2. docker-compose.clean.yml

**Purpose**: Complete production environment with full observability stack

**Included Services**:

| Service | Port | Purpose | Health Check |
|---------|------|---------|---------------|
| `movies-mcp-server-clean` | 8081 | Main application | Process check |
| `postgres` | 5433 | PostgreSQL 17 database | `pg_isready` |
| `migrations` | - | One-time migration runner | N/A |
| `redis-clean` | 6380 | Caching layer | `redis-cli ping` |
| `prometheus-clean` | 9091 | Metrics collection | HTTP check |
| `grafana-clean` | 3001 | Visualization dashboards | HTTP check |
| `pgadmin-clean` | 5051 | Database administration | HTTP check |

**Usage Examples**:
```bash
# Start all services
docker-compose -f docker-compose.clean.yml up -d

# View specific service logs
docker-compose -f docker-compose.clean.yml logs -f movies-mcp-server-clean

# Scale the main service
docker-compose -f docker-compose.clean.yml up -d --scale movies-mcp-server-clean=3

# Stop all services
docker-compose -f docker-compose.clean.yml down

# Stop with volume cleanup
docker-compose -f docker-compose.clean.yml down -v
```

### 3. docker-compose.dev.yml

**Purpose**: Lightweight development environment focused on database services

**Included Services**:

| Service | Port | Purpose | Credentials |
|---------|------|---------|-------------|
| `postgres-dev` | 5434 | Development database | `dev_user:dev_password` |
| `postgres-test` | 5435 | Test database | `test_user:test_password` |
| `redis-dev` | 6381 | Development cache | No auth |
| `pgadmin-dev` | 5052 | Database admin | `admin@example.com:admin` |

**Development Workflow**:
```bash
# Start development services
docker-compose -f docker-compose.dev.yml up -d

# Connect to development database
psql "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"

# Connect to test database  
psql "postgres://test_user:test_password@localhost:5435/movies_mcp_test"

# Run application with development database
export DATABASE_URL="postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev?sslmode=disable"
go run cmd/server/main.go
```

## ⚙️ Environment Configuration

### Production Environment Variables

Create a `.env.production` file:

```bash
# Database Configuration
POSTGRES_USER=movies_user
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=movies_mcp
DB_HOST=postgres
DB_PORT=5432

# Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=1h

# Application Configuration
LOG_LEVEL=info
SERVER_TIMEOUT=30s
MAX_IMAGE_SIZE=5242880  # 5MB
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp
ENABLE_THUMBNAILS=true

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_PORT=8080
PROMETHEUS_ENABLED=true
```

### Development Environment Variables

```bash
# Development Database
DATABASE_URL=postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev?sslmode=disable

# Test Database
TEST_DATABASE_URL=postgres://test_user:test_password@localhost:5435/movies_mcp_test?sslmode=disable

# Development Settings
LOG_LEVEL=debug
METRICS_ENABLED=false
```

## 🌐 Port Mapping Reference

### Clean Architecture (Production)

| Service | Host Port | Container Port | External Access |
|---------|-----------|----------------|-----------------|
| MCP Server Health | 8081 | 8080 | http://localhost:8081/health |
| PostgreSQL | 5433 | 5432 | localhost:5433 |
| Redis | 6380 | 6379 | localhost:6380 |
| Prometheus | 9091 | 9090 | http://localhost:9091 |
| Grafana | 3001 | 3000 | http://localhost:3001 |
| pgAdmin | 5051 | 80 | http://localhost:5051 |

### Development Environment

| Service | Host Port | Container Port | External Access |
|---------|-----------|----------------|-----------------|
| PostgreSQL (dev) | 5434 | 5432 | localhost:5434 |
| PostgreSQL (test) | 5435 | 5432 | localhost:5435 |
| Redis | 6381 | 6379 | localhost:6381 |
| pgAdmin | 5052 | 80 | http://localhost:5052 |

## 💾 Volume Management

### Production Volumes

```bash
# View all volumes
docker volume ls | grep movies

# Inspect volume details
docker volume inspect movies_postgres_data_clean

# Backup volume
docker run --rm -v movies_postgres_data_clean:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# Restore volume
docker run --rm -v movies_postgres_data_clean:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /data
```

### Volume Configuration

- **`postgres_data_clean`**: PostgreSQL data persistence
- **`redis_data_clean`**: Redis data persistence  
- **`prometheus_data_clean`**: Prometheus metrics history
- **`grafana_data_clean`**: Grafana dashboards and settings
- **`pgadmin_data_clean`**: pgAdmin configuration

## 🌐 Networking

### Network Architecture

```bash
# View networks
docker network ls | grep movies

# Inspect network
docker network inspect movies-network-clean
```

### Network Configuration

- **`movies-network-clean`** (172.21.0.0/16): Production network
- **`movies-dev-network`** (172.22.0.0/16): Development network  
- **`movies-network`** (172.20.0.0/16): Legacy network

### Service Discovery

Services communicate via Docker's internal DNS:

```bash
# From movies-server-clean container
ping postgres      # Resolves to PostgreSQL container
ping redis-clean   # Resolves to Redis container

# Connection strings use service names
DATABASE_URL=postgres://movies_user:password@postgres:5432/movies_mcp
REDIS_URL=redis://redis-clean:6379
```

## 🏥 Health Checks

### Built-in Health Checks

All services include comprehensive health checks:

```yaml
# PostgreSQL Health Check
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U movies_user -d movies_mcp"]
  interval: 10s
  timeout: 5s
  retries: 5

# Redis Health Check  
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 3s
  retries: 3

# MCP Server Health Check
healthcheck:
  test: ["CMD-SHELL", "pgrep movies-server-clean"]
  interval: 30s
  timeout: 10s
  retries: 3
```

### Manual Health Verification

```bash
# Check all service health
docker-compose -f docker-compose.clean.yml ps

# Detailed health status
docker inspect movies-server-clean --format='{{.State.Health.Status}}'

# Service-specific checks
curl http://localhost:8081/health        # MCP Server
curl http://localhost:9091/-/healthy     # Prometheus
curl http://localhost:3001/api/health    # Grafana
```

## 🔄 Migration Handling

### Automatic Migration Service

The clean architecture includes a dedicated migration service that runs automatically:

```yaml
migrations:
  build:
    context: .
    dockerfile: Dockerfile.clean
  entrypoint: ["/usr/local/bin/migrate"]
  command: ["postgres://movies_user:password@postgres:5432/movies_mcp", "/migrations", "up"]
  restart: "no"
  depends_on:
    postgres:
      condition: service_healthy
```

### Manual Migration Management

```bash
# Run migrations manually
docker-compose -f docker-compose.clean.yml run --rm migrations

# Check migration status
docker exec postgres psql -U movies_user -d movies_mcp -c "SELECT * FROM schema_migrations;"

# Rollback migrations (development only)
docker-compose -f docker-compose.clean.yml run --rm migrations \
  /usr/local/bin/migrate postgres://movies_user:password@postgres:5432/movies_mcp /migrations down 1
```

## 📊 Monitoring and Observability

### Prometheus Configuration

Access metrics at: **http://localhost:9091**

**Key Metrics Available**:
- `movies_requests_total` - Total requests by tool
- `movies_request_duration_seconds` - Request duration histogram
- `movies_db_connections_active` - Active database connections
- `movies_db_operations_total` - Database operations count

### Grafana Dashboards

Access dashboards at: **http://localhost:3001**
- **Username**: `admin`
- **Password**: `admin` (change on first login)

**Pre-configured Dashboards**:
- **MCP Server Overview**: Request rates, response times, error rates
- **Database Performance**: Connection pools, query performance
- **System Metrics**: CPU, memory, disk usage
- **Image Processing**: Upload statistics, processing times

### Database Administration

Access pgAdmin at: **http://localhost:5051** (prod) or **http://localhost:5052** (dev)
- **Email**: `admin@example.com`
- **Password**: `admin`

**Pre-configured Connections**:
- Production database (postgres:5432)
- Development database (postgres-dev:5432)
- Test database (postgres-test:5432)

## 🏭 Production Deployment

### 1. Image Registry Setup

```bash
# Build production image
docker build -f Dockerfile.clean -t movies-mcp-server:0.2.0 .

# Tag for your registry
docker tag movies-mcp-server:0.2.0 your-registry.com/movies-mcp-server:0.2.0
docker tag movies-mcp-server:0.2.0 your-registry.com/movies-mcp-server:latest

# Push to registry
docker push your-registry.com/movies-mcp-server:0.2.0
docker push your-registry.com/movies-mcp-server:latest
```

### 2. Production Deployment Script

```bash
#!/bin/bash
# deploy.sh

set -e

VERSION=${1:-latest}
ENVIRONMENT=${2:-production}

echo "🚀 Deploying Movies MCP Server $VERSION to $ENVIRONMENT"

# Pull latest images
docker-compose -f docker-compose.clean.yml pull

# Update environment
cp .env.$ENVIRONMENT .env

# Deploy with zero downtime
docker-compose -f docker-compose.clean.yml up -d --no-deps movies-server-clean

# Wait for health check
echo "⏳ Waiting for health check..."
for i in {1..30}; do
    if curl -f http://localhost:8081/health > /dev/null 2>&1; then
        echo "✅ Deployment successful!"
        exit 0
    fi
    sleep 10
done

echo "❌ Deployment failed - rolling back"
docker-compose -f docker-compose.clean.yml restart movies-server-clean
exit 1
```

### 3. Blue-Green Deployment

```bash
# Blue-green deployment script
#!/bin/bash

BLUE_PORT=8081
GREEN_PORT=8082

# Deploy to green
docker-compose -f docker-compose.green.yml up -d

# Health check green
if curl -f http://localhost:$GREEN_PORT/health; then
    # Switch traffic (update load balancer/reverse proxy)
    echo "✅ Green deployment healthy, switching traffic"
    
    # Stop blue after successful switch
    docker-compose -f docker-compose.clean.yml down
else
    echo "❌ Green deployment failed"
    docker-compose -f docker-compose.green.yml down
    exit 1
fi
```

## 🔧 Troubleshooting

### Common Issues

#### ❌ "Port already in use"

```bash
# Find process using port
lsof -i :5433
netstat -tlnp | grep :5433

# Stop conflicting services
docker-compose -f docker-compose.clean.yml down
docker-compose -f docker-compose.dev.yml down

# Clean up orphaned containers
docker container prune
```

#### ❌ "Database connection refused"

```bash
# Check database health
docker-compose -f docker-compose.clean.yml logs postgres

# Test connection manually
docker exec movies-server-clean pg_isready -h postgres -p 5432 -U movies_user

# Restart database service
docker-compose -f docker-compose.clean.yml restart postgres
```

#### ❌ "Migration failed"

```bash
# Check migration logs
docker-compose -f docker-compose.clean.yml logs migrations

# Reset database (CAUTION: This deletes all data)
docker-compose -f docker-compose.clean.yml down -v
docker volume rm movies_postgres_data_clean
docker-compose -f docker-compose.clean.yml up -d
```

#### ❌ "High memory usage"

```bash
# Monitor resource usage
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# Adjust resource limits in docker-compose.clean.yml
services:
  movies-server-clean:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### Debug Commands

```bash
# View detailed container information
docker inspect movies-server-clean

# Execute shell in container
docker-compose -f docker-compose.clean.yml exec movies-server-clean sh

# View container filesystem
docker diff movies-server-clean

# Export container logs
docker-compose -f docker-compose.clean.yml logs movies-server-clean > server.log
```

## 🛡️ Security Best Practices

### Container Security

```bash
# Scan images for vulnerabilities
docker scan movies-mcp-server:latest

# Use specific versions (avoid :latest in production)
image: movies-mcp-server:0.2.0

# Run security benchmarks
docker run --rm -it --pid host --userns host --cap-add audit_control \
  -e DOCKER_CONTENT_TRUST=$DOCKER_CONTENT_TRUST \
  -v /etc:/etc:ro \
  -v /usr/bin/containerd:/usr/bin/containerd:ro \
  -v /usr/bin/runc:/usr/bin/runc:ro \
  -v /usr/lib/systemd:/usr/lib/systemd:ro \
  docker/docker-bench-security
```

### Network Security

```bash
# Restrict external access
# Only expose necessary ports
ports:
  - "127.0.0.1:8081:8080"  # Only localhost

# Use custom networks
networks:
  movies-internal:
    driver: bridge
    internal: true  # No external connectivity
```

### Secrets Management

```bash
# Use Docker secrets in production
echo "secure_password" | docker secret create db_password -

# Reference in compose file
services:
  postgres:
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

secrets:
  db_password:
    external: true
```

## ⚡ Performance Optimization

### Build Optimization

```dockerfile
# Multi-stage build optimization
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
RUN go build -ldflags="-w -s" ./cmd/server
```

### Runtime Optimization

```yaml
# docker-compose.clean.yml optimization
services:
  movies-server-clean:
    # Set resource limits
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    
    # Optimize for production
    environment:
      GOGC: 100
      GOMEMLIMIT: 900MiB
```

### Database Optimization

```yaml
postgres:
  environment:
    # Optimize PostgreSQL for containers
    POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
  command: |
    postgres
    -c shared_buffers=256MB
    -c effective_cache_size=1GB
    -c maintenance_work_mem=64MB
    -c checkpoint_completion_target=0.9
    -c wal_buffers=16MB
    -c default_statistics_target=100
```

## 📚 Next Steps

After setting up Docker deployment:

1. **[Production Deployment](README.md)** - Complete production deployment guide
2. **[Monitoring Setup](../reference/monitoring.md)** - Advanced monitoring configuration
3. **[Performance Tuning](../reference/performance.md)** - Optimize for your workload
4. **[Security Hardening](../reference/security.md)** - Production security checklist

Your Docker-based Movies MCP Server deployment is now ready for production! 🎉