# Scripts Directory

This directory contains utility scripts for the Movies MCP Server project.

## Available Scripts

### 1. deploy.sh
**Comprehensive deployment script for multiple environments**

```bash
./scripts/deploy.sh [OPTIONS] COMMAND

# Examples:
./scripts/deploy.sh docker                    # Deploy with Docker Compose (dev)
./scripts/deploy.sh build -e staging          # Build staging images
./scripts/deploy.sh test                      # Run deployment tests
./scripts/deploy.sh status -n movies-prod     # Check production status
```

**Commands:**
- `docker` - Deploy using Docker Compose
- `build` - Build Docker images
- `test` - Run deployment tests
- `clean` - Clean up resources
- `status` - Check deployment status

**Options:**
- `-e, --env` - Environment (development, staging, production)
- `-n, --namespace` - Kubernetes namespace
- `-f, --force` - Force deployment (skip confirmations)
- `-v, --verbose` - Verbose output

### 2. seed.sh
**Database seeding script with sample movie data**

```bash
./scripts/seed.sh [OPTIONS]

# Examples:
./scripts/seed.sh                             # Seed development database
./scripts/seed.sh -e production -f            # Force seed production
./scripts/seed.sh -c                          # Clear existing data only
./scripts/seed.sh -v                          # Verbose output
```

**Options:**
- `-e, --env` - Environment (development, staging, production)
- `-f, --force` - Force seeding (clear existing data)
- `-c, --clear-only` - Only clear existing data
- `-v, --verbose` - Verbose output

**Features:**
- 40+ sample movies from classics to modern hits
- Complete movie data (title, director, year, genre, rating, description, poster URLs)
- Data validation and integrity checks
- Environment-specific database support
- Safe operation with confirmation prompts

### 3. seed_data.sql
**SQL file containing sample movie data**

This file contains INSERT statements for 40 popular and classic movies including:
- The Shawshank Redemption (1994) - 9.3⭐
- The Godfather (1972) - 9.2⭐
- The Dark Knight (2008) - 9.0⭐
- And many more...

**Data includes:**
- Movie titles and directors
- Release years (1942-2019)
- Genres (Drama, Action, Crime, etc.)
- IMDB-style ratings
- Detailed descriptions
- TMDB poster URLs

## Usage Instructions

### Prerequisites

**For deploy.sh:**
- Docker and Docker Compose (for Docker deployments)
- Access to target environment

**For seed.sh:**
- PostgreSQL client (`psql`)
- Database connection credentials
- Initialized database with movies table

### Environment Variables

Scripts read configuration from:
1. Command line arguments
2. `.env` file in project root
3. Environment variables
4. Default values

**Key environment variables:**
```bash
# Database connection
DATABASE_URL=postgres://user:pass@host:port/db
DB_HOST=localhost
DB_PORT=5432
DB_NAME=movies_db
DB_USER=movies_user
DB_PASSWORD=password

# Deployment
ENVIRONMENT=development
PROD_URL=https://movies-mcp.yourdomain.com
```

### Quick Start

1. **Set up environment:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

2. **Deploy locally:**
```bash
./scripts/deploy.sh docker
```

3. **Seed database:**
```bash
./scripts/seed.sh
```

4. **Check status:**
```bash
./scripts/deploy.sh status
```

### Advanced Usage

#### Multi-environment Deployment
```bash
# Development
./scripts/deploy.sh docker

# Staging with Docker
./scripts/deploy.sh -e staging docker

# Production (with force flag)
./scripts/deploy.sh -e production -f docker
```

#### Database Management
```bash
# Seed development database
./scripts/seed.sh -e development

# Force reseed staging (clears existing data)
./scripts/seed.sh -e staging -f

# Clear production data only
./scripts/seed.sh -e production -c
```

#### Build and Test Pipeline
```bash
# Build images for all environments
for env in development staging production; do
    ./scripts/deploy.sh -e $env build
done

# Test deployments
./scripts/deploy.sh test
./scripts/deploy.sh -e staging test
```

### Error Handling

Scripts include comprehensive error handling:
- Dependency checks
- Connection validation
- Rollback on failure
- Detailed error messages
- Exit codes for CI/CD integration

### Logging

All scripts provide colored, structured logging:
- **INFO** (Blue): General information
- **SUCCESS** (Green): Successful operations
- **WARNING** (Yellow): Non-critical issues
- **ERROR** (Red): Critical failures

### Integration with CI/CD

Scripts are designed for CI/CD integration:

**GitHub Actions example:**
```yaml
- name: Deploy to staging
  run: ./scripts/deploy.sh -e staging -f k8s

- name: Seed database
  run: ./scripts/seed.sh -e staging -f

- name: Run tests
  run: ./scripts/deploy.sh -e staging test
```

**GitLab CI example:**
```yaml
deploy:staging:
  script:
    - ./scripts/deploy.sh -e staging -f k8s
    - ./scripts/seed.sh -e staging -f
    - ./scripts/deploy.sh -e staging test
```

### Troubleshooting

#### Common Issues

1. **Permission denied:**
```bash
chmod +x scripts/*.sh
```

2. **Database connection failed:**
```bash
# Check environment variables
echo $DATABASE_URL

# Test connection manually
psql $DATABASE_URL -c "SELECT 1;"
```


4. **Docker deployment failed:**
```bash
# Check Docker daemon
docker info

# Verify images
docker images | grep movies-mcp

# Check logs
docker-compose logs
```

#### Debug Mode

Enable verbose output for troubleshooting:
```bash
./scripts/deploy.sh -v docker
./scripts/seed.sh -v
```

### Script Customization

Scripts are designed to be customizable:

1. **Modify default values** in script headers
2. **Add custom environment logic** in environment-specific sections
3. **Extend validation functions** for additional checks
4. **Add custom deployment targets** in command handling

### Security Considerations

- Scripts mask sensitive data in logs
- Use environment variables for secrets
- Validate inputs and connections
- Implement least-privilege access
- Support secure secret management (K8s secrets, etc.)

### Contributing

When modifying scripts:
1. Test in all supported environments
2. Update help documentation
3. Maintain backward compatibility
4. Add appropriate error handling
5. Update this README