# Docker Configuration - Clean Architecture

This document describes the Docker setup for the Movies MCP Server with Clean Architecture.

## Version Requirements

- **Go**: 1.24.4
- **PostgreSQL**: 17
- **Alpine Linux**: Latest (for minimal image size)
- **Redis**: 7-alpine

## Overview

The Docker configuration provides multiple environments:

- **Production**: `Dockerfile.clean` + `docker-compose.clean.yml`
- **Development**: `docker-compose.dev.yml`
- **Legacy**: `Dockerfile` + `docker-compose.yml` (original architecture)

## Quick Start

### Development Environment

```bash
# Start development databases
docker-compose -f docker-compose.dev.yml up -d

# Set environment variables
export TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5435/movies_mcp_test?sslmode=disable"

# Run the application locally
go run cmd/server-new/main.go

# Run tests
go test ./...
```

### Production Environment

```bash
# Build and start all services
docker-compose -f docker-compose.clean.yml up --build

# Or run in background
docker-compose -f docker-compose.clean.yml up -d --build
```

## Docker Files

### 1. Dockerfile.clean

**Purpose**: Production-ready multi-stage build for clean architecture

**Features**:
- Multi-stage build for minimal image size
- Uses Go 1.24.4 for building
- Builds both `movies-server` and `migrate` binaries
- Uses distroless base image for security
- Non-root user execution
- Includes health checks

**Build**:
```bash
docker build -f Dockerfile.clean -t movies-mcp-server:clean .
```

### 2. docker-compose.clean.yml

**Purpose**: Complete production environment with monitoring

**Services**:
- `movies-mcp-server-clean`: Main application (port 8081)
- `postgres`: PostgreSQL 17 database (port 5433)
- `migrations`: One-time migration runner
- `redis-clean`: Caching (port 6380)
- `prometheus-clean`: Metrics collection (port 9091)
- `grafana-clean`: Visualization (port 3001)
- `pgadmin-clean`: Database admin (port 5051)

**Usage**:
```bash
# Start all services
docker-compose -f docker-compose.clean.yml up -d

# View logs
docker-compose -f docker-compose.clean.yml logs -f movies-mcp-server-clean

# Stop all services
docker-compose -f docker-compose.clean.yml down
```

### 3. docker-compose.dev.yml

**Purpose**: Lightweight development environment

**Services**:
- `postgres-dev`: PostgreSQL 17 development database (port 5434)
- `postgres-test`: PostgreSQL 17 test database (port 5435)
- `redis-dev`: Development cache (port 6381)
- `pgadmin-dev`: Database admin (port 5052)

**Usage**:
```bash
# Start development services
docker-compose -f docker-compose.dev.yml up -d

# Connect to development database
psql "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"

# Connect to test database
psql "postgres://test_user:test_password@localhost:5435/movies_mcp_test"
```

## Environment Variables

### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `movies_mcp` | Database name |
| `DB_USER` | `movies_user` | Database user |
| `DB_PASSWORD` | `movies_password` | Database password |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_MAX_OPEN_CONNS` | `25` | Max open connections |
| `DB_MAX_IDLE_CONNS` | `5` | Max idle connections |
| `DB_CONN_MAX_LIFETIME` | `1h` | Connection max lifetime |

### Application Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `SERVER_TIMEOUT` | `30s` | Server timeout |
| `MAX_IMAGE_SIZE` | `5242880` | Max image size (5MB) |
| `ALLOWED_IMAGE_TYPES` | `image/jpeg,image/png,image/webp` | Allowed image types |
| `ENABLE_THUMBNAILS` | `true` | Enable thumbnail generation |
| `THUMBNAIL_SIZE` | `200x200` | Thumbnail size |

## Port Mapping

### Clean Architecture (Production)

| Service | Host Port | Container Port | Purpose |
|---------|-----------|----------------|---------|
| MCP Server | 8081 | 8080 | Health check |
| PostgreSQL | 5433 | 5432 | Database |
| Redis | 6380 | 6379 | Cache |
| Prometheus | 9091 | 9090 | Metrics |
| Grafana | 3001 | 3000 | Dashboards |
| pgAdmin | 5051 | 80 | DB Admin |

### Development

| Service | Host Port | Container Port | Purpose |
|---------|-----------|----------------|---------|
| PostgreSQL (dev) | 5434 | 5432 | Dev database |
| PostgreSQL (test) | 5435 | 5432 | Test database |
| Redis | 6381 | 6379 | Dev cache |
| pgAdmin | 5052 | 80 | DB Admin |

### Legacy (Original)

| Service | Host Port | Container Port | Purpose |
|---------|-----------|----------------|---------|
| MCP Server | 8080 | 8080 | Health check |
| PostgreSQL | 5432 | 5432 | Database |
| Redis | 6379 | 6379 | Cache |
| Prometheus | 9090 | 9090 | Metrics |
| Grafana | 3000 | 3000 | Dashboards |
| pgAdmin | 5050 | 80 | DB Admin |

## Volumes

### Clean Architecture

- `postgres_data_clean`: PostgreSQL data
- `redis_data_clean`: Redis data
- `prometheus_data_clean`: Prometheus data
- `grafana_data_clean`: Grafana data
- `pgladmin_data_clean`: pgAdmin data

### Development

- `postgres_dev_data`: Development PostgreSQL data
- `postgres_test_data`: Test PostgreSQL data
- `redis_dev_data`: Development Redis data
- `pgadmin_dev_data`: Development pgAdmin data

## Networks

- `movies-network-clean` (172.21.0.0/16): Production network
- `movies-dev-network` (172.22.0.0/16): Development network
- `movies-network` (172.20.0.0/16): Legacy network

## Health Checks

All services include health checks:

- **PostgreSQL**: `pg_isready` command
- **Redis**: `redis-cli ping`
- **MCP Server**: Process check (since it uses stdin/stdout)

## Migration Handling

The clean architecture includes a dedicated migration service:

```yaml
migrations:
  build:
    context: .
    dockerfile: Dockerfile.clean
  entrypoint: ["/usr/local/bin/migrate"]
  command: ["postgres://...", "/migrations", "up"]
  restart: "no"
```

This runs automatically when the stack starts and applies all pending migrations.

## Monitoring and Observability

### Prometheus Metrics

Access metrics at: http://localhost:9091

### Grafana Dashboards

Access dashboards at: http://localhost:3001
- Username: `admin`
- Password: `admin`

### Database Administration

Access pgAdmin at: http://localhost:5051 (prod) or http://localhost:5052 (dev)
- Email: `admin@example.com`
- Password: `admin`

## Build Arguments

When building the Docker image, you can set:

```bash
docker build -f Dockerfile.clean \
  --build-arg VERSION=0.2.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t movies-mcp-server:0.2.0 .
```

## Production Deployment

### 1. Build and Tag

```bash
# Build production image
docker build -f Dockerfile.clean -t movies-mcp-server:0.2.0 .

# Tag for registry
docker tag movies-mcp-server:0.2.0 your-registry/movies-mcp-server:0.2.0
```

### 2. Push to Registry

```bash
docker push your-registry/movies-mcp-server:0.2.0
```

### 3. Deploy

```bash
# Update image in docker-compose.clean.yml
# image: your-registry/movies-mcp-server:0.2.0

# Deploy
docker-compose -f docker-compose.clean.yml up -d
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Different environments use different ports
2. **Database connection**: Check health checks and environment variables
3. **Migration failures**: Check migration service logs
4. **Permission issues**: All services run as non-root users

### Debugging

```bash
# View service logs
docker-compose -f docker-compose.clean.yml logs [service-name]

# Execute shell in container
docker-compose -f docker-compose.clean.yml exec movies-mcp-server-clean sh

# Check service health
docker-compose -f docker-compose.clean.yml ps
```

### Database Access

```bash
# Development database
psql "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"

# Production database
psql "postgres://movies_user:movies_password@localhost:5433/movies_mcp"

# Test database
psql "postgres://test_user:test_password@localhost:5435/movies_mcp_test"
```

## Security Considerations

1. **Non-root execution**: All services run as non-root users
2. **Distroless base**: Minimal attack surface
3. **Environment isolation**: Separate networks for different environments
4. **Secret management**: Use Docker secrets or external secret management in production
5. **Image scanning**: Regularly scan images for vulnerabilities

## Performance Optimization

1. **Multi-stage builds**: Minimize image size
2. **Layer caching**: Optimize Dockerfile layer order
3. **Connection pooling**: Configured in environment variables
4. **Resource limits**: Set appropriate CPU and memory limits in production

```yaml
# Example resource limits
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
    reservations:
      cpus: '0.25'
      memory: 256M
```