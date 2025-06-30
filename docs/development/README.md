# 🔧 Development Setup

Set up a complete development environment for the Movies MCP Server.

## Prerequisites

### Required Tools
- **Go 1.24.4+** - [Download here](https://golang.org/doc/install)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Make** - Usually pre-installed on macOS/Linux

### Recommended Tools
- **PostgreSQL client** - `psql` for database access
- **VS Code** with Go extension
- **golangci-lint** - `brew install golangci-lint` (linting)

## Quick Development Setup

### 1. Clone and Enter
```bash
git clone https://github.com/francknouama/movies-mcp-server.git
cd movies-mcp-server/mcp-server
```

### 2. Start Development Services
```bash
# Start development databases and tools
make docker-compose-up-dev

# Verify services are running
docker ps | grep -E "(postgres|redis|pgladmin)"
```

### 3. Set Environment Variables
```bash
# Development environment
export DATABASE_URL="postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"
export LOG_LEVEL=debug
```

### 4. Build and Test
```bash
# Download dependencies
make deps

# Build clean architecture version (recommended)
make build-clean

# Run tests
make test

# Run with test coverage
make test-coverage
```

### 5. Start Development Server
```bash
# Run clean architecture server
make run-clean

# Or run with auto-restart (requires entr)
find . -name "*.go" | entr -r make run-clean
```

## Development Workflow

### Code → Test → Build → Run Cycle

```bash
# 1. Make changes to code
vim internal/application/movie/service.go

# 2. Run relevant tests
go test ./internal/application/movie/...

# 3. Run all tests
make test

# 4. Check code quality
make check  # runs fmt, vet, lint

# 5. Build and test
make build-clean
make test-init
```

### Database Development

```bash
# Connect to development database
psql "postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev"

# Reset database (clean slate)
make db-migrate-reset

# Seed with sample data
make db-seed

# Check migration status
make db-migrate-version
```

### Running Tests

```bash
# Unit tests only (fast)
go test ./internal/domain/... ./internal/application/...

# Integration tests (requires database)
export TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5435/movies_mcp_test"
go test ./tests/integration/...

# All tests with coverage
make test-coverage
open coverage.html  # View coverage report
```

## Development Architecture

### Clean Architecture Structure
```
mcp-server/
├── 🏗️ cmd/
│   ├── server/       # Legacy entry point
│   └── server-new/   # Clean architecture entry (main.go missing, use cmd/server/)
├── 🧠 internal/
│   ├── domain/       # 💎 Business logic (pure, no dependencies)
│   ├── application/  # 🔄 Use cases & orchestration
│   ├── infrastructure/ # 🔌 Database, external APIs
│   └── interfaces/   # 🌐 MCP handlers, DTOs
├── 🧪 tests/
│   ├── integration/ # End-to-end tests
│   └── performance/ # Load & benchmark tests
└── 🛠️ tools/
    └── migrate/     # Custom migration tool
```

### Key Principles
1. **Domain layer** - Pure business logic, no external dependencies
2. **Application layer** - Orchestrates use cases, calls domain
3. **Infrastructure layer** - Database, file system, external APIs  
4. **Interface layer** - MCP protocol, DTOs, dependency injection

## Debugging

### Debug Server Communication
```bash
# Test MCP protocol manually
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}},"id":1}' | \
./build/movies-server-clean

# Debug with detailed logging
LOG_LEVEL=debug ./build/movies-server-clean
```

### Debug Database Issues
```bash
# Check database connectivity
psql "$DATABASE_URL" -c "SELECT version();"

# Check table structure
psql "$DATABASE_URL" -c "\dt"

# Check migration status
psql "$DATABASE_URL" -c "SELECT * FROM schema_migrations ORDER BY version;"

# View logs
docker logs postgres-dev
```

### Debug Build Issues
```bash
# Clean everything
make clean

# Check Go installation
go version
go env GOPATH GOROOT

# Update dependencies
make deps-update

# Rebuild
make build-clean
```

## Performance Profiling

### Memory and CPU Profiling
```bash
# Build with profiling
go build -o movies-server-debug cmd/server/main.go

# Run with CPU profiling  
./movies-server-debug -cpuprofile=cpu.prof

# Analyze profile
go tool pprof cpu.prof
```

### Database Performance
```bash
# Enable query logging in PostgreSQL
docker exec -it postgres-dev psql -U dev_user -d movies_mcp_dev -c "
  ALTER SYSTEM SET log_statement = 'all';
  SELECT pg_reload_conf();
"

# Check slow queries
docker logs postgres-dev | grep "slow query"
```

### Load Testing
```bash
# Run performance tests
go test -bench=. ./tests/integration/...

# Custom load test
for i in {1..100}; do
  echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_movies","arguments":{"query":"matrix"}},"id":'$i'}' | \
  ./build/movies-server-clean &
done
```

## Contributing Guidelines

### Code Standards
1. **Follow Clean Architecture** - Respect layer boundaries
2. **Write tests first** - TDD approach preferred
3. **Use value objects** - Type-safe domain models
4. **No external dependencies in domain** - Keep business logic pure
5. **Document public APIs** - Godoc comments for exported functions

### Git Workflow
```bash
# 1. Create feature branch
git checkout -b feature/add-actor-search

# 2. Make changes with tests
# 3. Ensure tests pass
make test

# 4. Check code quality
make check

# 5. Commit with clear message
git commit -m "Add actor search functionality

- Implement actor search by name
- Add fuzzy matching for partial names  
- Include unit and integration tests
- Update API documentation"

# 6. Push and create PR
git push origin feature/add-actor-search
```

### Before Submitting PR
- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] No security issues (`make security-check`)
- [ ] Documentation updated
- [ ] Clean commit history

## Advanced Development

### Custom Migration
```bash
# Create new migration files
mkdir -p migrations
cat > migrations/003_add_ratings_table.up.sql << 'EOF'
CREATE TABLE user_ratings (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER REFERENCES movies(id),
    user_id VARCHAR(255) NOT NULL,
    rating DECIMAL(3,1) CHECK (rating >= 0 AND rating <= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
EOF

cat > migrations/003_add_ratings_table.down.sql << 'EOF'
DROP TABLE user_ratings;
EOF

# Test migration
make db-migrate
```

### Add New MCP Tool
```bash
# 1. Define in domain layer
vim internal/domain/movie/repository.go

# 2. Implement in application layer  
vim internal/application/movie/service.go

# 3. Add infrastructure implementation
vim internal/infrastructure/postgres/movie_repository.go

# 4. Create interface handler
vim internal/interfaces/mcp/movie_handlers.go

# 5. Register in server
vim internal/server/mcp_server.go

# 6. Add tests
vim internal/application/movie/service_test.go
```

### Environment Configuration
```bash
# Create .env.development
cat > .env.development << 'EOF'
DATABASE_URL=postgres://dev_user:dev_password@localhost:5434/movies_mcp_dev
LOG_LEVEL=debug
MAX_IMAGE_SIZE=10485760
ENABLE_THUMBNAILS=true
THUMBNAIL_SIZE=300x300
EOF

# Load environment
source .env.development
```

## Troubleshooting Development Issues

### Port Already in Use
```bash
# Find what's using the port
lsof -i :5434

# Stop development services
make docker-compose-down-dev

# Kill specific process
kill -9 <PID>
```

### Go Module Issues
```bash
# Clear module cache
go clean -modcache

# Verify module integrity
go mod verify

# Update all dependencies
go get -u ./...
go mod tidy
```

### Docker Issues
```bash
# Clean up Docker
docker system prune -f

# Rebuild development containers
make docker-compose-down-dev
docker-compose -f docker-compose.dev.yml build --no-cache
make docker-compose-up-dev
```

## Development Resources

- **[Clean Architecture Guide](./architecture.md)** - Detailed architecture explanation
- **[API Reference](../reference/api.md)** - Complete MCP tool documentation  
- **[Testing Strategy](./testing.md)** - Testing approaches and patterns
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions

## Quick Reference

```bash
# Essential commands
make deps          # Download dependencies
make build-clean   # Build clean architecture version  
make test          # Run all tests
make run-clean     # Start development server
make check         # Code quality checks
make docker-compose-up-dev   # Start development services
make db-migrate    # Run database migrations
make db-seed       # Add sample data
```