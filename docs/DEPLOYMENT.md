# Movies MCP Server - Deployment Guide

This guide covers deploying the Movies MCP Server to various production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Docker Deployment](#docker-deployment)
- [Cloud Platform Deployment](#cloud-platform-deployment)
- [Database Setup](#database-setup)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

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
DATABASE_URL=postgres://user:password@host:5432/movies_db
DB_HOST=localhost
DB_PORT=5432
DB_NAME=movies_db
DB_USER=movies_user
DB_PASSWORD=secure_password
DB_SSL_MODE=require

# Server Configuration
PORT=8080
HTTP_PORT=8080  # For MCP over HTTP testing
LOG_LEVEL=info
LOG_FORMAT=json

# Security
JWT_SECRET=your-super-secure-jwt-secret-here
API_KEY=your-api-key-for-authentication

# External Services
REDIS_URL=redis://localhost:6379
METRICS_PORT=9090

# Performance
MAX_CONNECTIONS=100
QUERY_TIMEOUT=30s
IDLE_TIMEOUT=5m
```

### Optional Environment Variables

```bash
# Monitoring
PROMETHEUS_ENABLED=true
GRAFANA_ENABLED=true
HEALTH_CHECK_INTERVAL=30s

# Caching
CACHE_TTL=1h
CACHE_MAX_SIZE=1000

# Development
DEBUG=false
PROFILE=false
TRACE_ENABLED=false
```

## Docker Deployment

### Using Docker Compose (Recommended)

1. **Clone the repository:**
```bash
git clone <repository-url>
cd movies-mcp-server
```

2. **Create environment file:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Deploy with Docker Compose:**
```bash
docker-compose up -d
```

4. **Verify deployment:**
```bash
docker-compose ps
docker-compose logs movies-mcp-server
```

### Using Production Dockerfile

1. **Build production image:**
```bash
docker build -f Dockerfile.production -t movies-mcp-server:latest .
```

2. **Run with external database:**
```bash
docker run -d \
  --name movies-mcp-server \
  --restart unless-stopped \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/movies_db \
  -e LOG_LEVEL=info \
  movies-mcp-server:latest
```

### Docker Swarm Deployment

1. **Initialize swarm:**
```bash
docker swarm init
```

2. **Deploy stack:**
```bash
docker stack deploy -c docker-compose.yml movies-stack
```

## Cloud Platform Deployment

### AWS ECS

1. **Create task definition:**
```json
{
  "family": "movies-mcp-server",
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskRole",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "containerDefinitions": [
    {
      "name": "movies-mcp-server",
      "image": "YOUR_ECR_REPO/movies-mcp-server:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:movies-db-url"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/movies-mcp-server",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

2. **Create ECS service:**
```bash
aws ecs create-service \
  --cluster movies-cluster \
  --service-name movies-mcp-service \
  --task-definition movies-mcp-server \
  --desired-count 2 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-12345],securityGroups=[sg-12345],assignPublicIp=ENABLED}"
```

### Google Cloud Run

1. **Build and push image:**
```bash
gcloud builds submit --tag gcr.io/PROJECT_ID/movies-mcp-server
```

2. **Deploy to Cloud Run:**
```bash
gcloud run deploy movies-mcp-server \
  --image gcr.io/PROJECT_ID/movies-mcp-server \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars LOG_LEVEL=info \
  --set-secrets DATABASE_URL=movies-db-url:latest
```

### Azure Container Instances

```bash
az container create \
  --resource-group movies-rg \
  --name movies-mcp-server \
  --image movies-mcp-server:latest \
  --cpu 1 \
  --memory 2 \
  --ports 8080 \
  --environment-variables LOG_LEVEL=info \
  --secure-environment-variables DATABASE_URL=$DATABASE_URL
```

## Database Setup

### PostgreSQL Production Setup

1. **Create database and user:**
```sql
CREATE DATABASE movies_db;
CREATE USER movies_user WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE movies_db TO movies_user;
```

2. **Run migrations:**
```bash
# Using the application
./movies-mcp-server --migrate

# Or using psql
psql -h hostname -U movies_user -d movies_db -f migrations/001_initial.sql
```

3. **Configure connection pooling:**
```bash
# Using PgBouncer
echo "movies_db = host=db_host port=5432 dbname=movies_db" >> pgbouncer.ini
```

### Cloud Database Options

#### AWS RDS
```bash
aws rds create-db-instance \
  --db-instance-identifier movies-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 13.7 \
  --allocated-storage 20 \
  --storage-encrypted \
  --master-username movies_user \
  --master-user-password secure_password
```

#### Google Cloud SQL
```bash
gcloud sql instances create movies-db \
  --database-version POSTGRES_13 \
  --tier db-f1-micro \
  --region us-central1 \
  --storage-auto-increase
```

## Monitoring and Logging

### Prometheus Metrics

The server exposes metrics on `/metrics` endpoint:
- Request duration histograms
- Request count counters
- Active connections gauge
- Database operation metrics

### Grafana Dashboard

Import the provided dashboard:
```bash
curl -X POST \
  http://grafana:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @monitoring/grafana-dashboard.json
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

### Network Security
- Use TLS/SSL for all connections
- Implement network segmentation
- Configure firewall rules
- Use VPN for database access

### Application Security
- Enable CORS restrictions
- Implement rate limiting
- Use strong JWT secrets
- Validate all inputs
- Implement audit logging

### Container Security
- Use non-root user in containers
- Scan images for vulnerabilities
- Use minimal base images
- Keep dependencies updated
- Implement security policies

### Example Security Headers
```go
// Add to your HTTP server
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000")
```

## Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Test database connectivity
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;"

# Check connection pool
curl http://localhost:8080/health/db
```

#### Performance Issues
```bash
# Check resource usage
docker stats movies-mcp-server

# View application metrics
curl http://localhost:9090/metrics | grep movies_
```

#### Log Analysis
```bash
# View recent logs
docker logs --tail 100 movies-mcp-server

# Filter error logs
docker logs movies-mcp-server 2>&1 | grep ERROR
```

### Debug Mode

Enable debug mode for troubleshooting:
```bash
export DEBUG=true
export LOG_LEVEL=debug
./movies-mcp-server
```

### Health Checks

The server provides multiple health check endpoints:
- `/health` - Overall health status
- `/health/db` - Database connectivity
- `/health/redis` - Redis connectivity (if enabled)
- `/ready` - Readiness for traffic

### Backup and Recovery

#### Database Backup
```bash
# Create backup
pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > movies_backup.sql

# Restore backup
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < movies_backup.sql
```

#### Automated Backups
```bash
# Add to crontab
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | gzip > /backups/movies_$(date +\%Y\%m\%d).sql.gz
```

## Scaling Considerations

### Horizontal Scaling
- Deploy multiple instances behind a load balancer
- Use external cache (Redis) for session storage
- Implement connection pooling for database

### Vertical Scaling
- Monitor CPU and memory usage
- Adjust container resource limits
- Optimize database queries

### Load Balancing
```nginx
# nginx.conf
upstream movies_backend {
    server movies-1:8080;
    server movies-2:8080;
    server movies-3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://movies_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

This deployment guide provides comprehensive instructions for deploying the Movies MCP Server in various production environments with proper security, monitoring, and scaling considerations.