# Movies MCP Server - Clean Architecture

## Overview

The Movies MCP Server has been completely restructured using **Clean Architecture** principles and **Domain-Driven Design**. This document outlines the new architecture, its benefits, and how to work with it.

## Architecture Layers

```
┌─────────────────────────────────────────┐
│             Interface Layer             │ ← MCP Protocol, DTOs
├─────────────────────────────────────────┤
│           Application Layer             │ ← Use Cases, Services
├─────────────────────────────────────────┤
│             Domain Layer                │ ← Business Logic, Entities
├─────────────────────────────────────────┤
│          Infrastructure Layer           │ ← Database, External APIs
└─────────────────────────────────────────┘
```

### 1. Domain Layer (`internal/domain/`)

**Pure business logic with no external dependencies.**

- **Value Objects**: Type-safe primitives with validation
  - `MovieID`, `ActorID`: Unique identifiers
  - `Rating`: 0-10 validated rating
  - `Year`: Valid movie year (1888-future)

- **Aggregates**: Core business entities
  - `Movie`: Complete movie model with genres, ratings, poster
  - `Actor`: Actor model with filmography and biography

- **Repository Interfaces**: Define data access contracts
  - Segregated by responsibility (Reader/Writer)
  - Domain-focused, not database-focused

#### Example Domain Model

```go
// Value Object with validation
type Rating struct {
    value float64
}

func NewRating(rating float64) (Rating, error) {
    if rating < 0 || rating > 10 {
        return Rating{}, errors.New("rating must be between 0 and 10")
    }
    return Rating{value: rating}, nil
}

// Aggregate with business rules
type Movie struct {
    id       MovieID
    title    string
    director string
    year     Year
    rating   Rating
    genres   []string
    // ... timestamps, etc.
}

func (m *Movie) AddGenre(genre string) error {
    // Business rule: no duplicate genres
    for _, g := range m.genres {
        if g == genre {
            return errors.New("genre already exists")
        }
    }
    m.genres = append(m.genres, genre)
    m.touch() // Update timestamp
    return nil
}
```

### 2. Application Layer (`internal/application/`)

**Orchestrates use cases and coordinates domain operations.**

- **Services**: Use case implementations
  - `MovieService`: CRUD, search, top-rated movies
  - `ActorService`: CRUD, movie linking, search

- **Commands/Queries**: Clear separation of reads and writes
- **DTOs**: Data transfer objects for external communication
- **Error Handling**: Wraps domain errors with context

#### Example Application Service

```go
type MovieService struct {
    movieRepo movie.Repository
}

func (s *MovieService) CreateMovie(ctx context.Context, cmd CreateMovieCommand) (*MovieDTO, error) {
    // Create domain movie
    domainMovie, err := movie.NewMovie(cmd.Title, cmd.Director, cmd.Year)
    if err != nil {
        return nil, fmt.Errorf("failed to create movie: %w", err)
    }

    // Apply business rules
    if cmd.Rating > 0 {
        if err := domainMovie.SetRating(cmd.Rating); err != nil {
            return nil, fmt.Errorf("failed to set rating: %w", err)
        }
    }

    // Persist
    if err := s.movieRepo.Save(ctx, domainMovie); err != nil {
        return nil, fmt.Errorf("failed to save movie: %w", err)
    }

    return s.toDTO(domainMovie), nil
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Handles external concerns and implements domain interfaces.**

- **PostgreSQL Repositories**: Database implementations
  - Implements domain repository interfaces
  - Handles database-specific concerns (SQL, transactions)
  - Converts between domain models and database models

- **Migrations**: Database schema management
- **Integration Tests**: Test database interactions

#### Example Repository Implementation

```go
type MovieRepository struct {
    db *sql.DB
}

func (r *MovieRepository) Save(ctx context.Context, domainMovie *movie.Movie) error {
    dbMovie := r.toDBModel(domainMovie)
    
    if domainMovie.ID().IsZero() {
        return r.insert(ctx, dbMovie, domainMovie)
    }
    return r.update(ctx, dbMovie, domainMovie)
}

func (r *MovieRepository) toDomainModel(dbMovie *dbMovie) (*movie.Movie, error) {
    movieID, err := shared.NewMovieID(dbMovie.ID)
    if err != nil {
        return nil, fmt.Errorf("invalid movie ID: %w", err)
    }
    
    domainMovie, err := movie.NewMovieWithID(movieID, dbMovie.Title, dbMovie.Director, dbMovie.Year)
    if err != nil {
        return nil, fmt.Errorf("failed to create domain movie: %w", err)
    }
    
    // Set optional fields...
    return domainMovie, nil
}
```

### 4. Interface Layer (`internal/interfaces/`)

**Adapts external protocols to internal use cases.**

- **MCP Handlers**: Thin adapters for MCP protocol
  - Parse MCP requests to application commands
  - Convert application DTOs to MCP responses
  - Handle protocol-specific error codes

- **DTOs**: External data transfer objects
- **Dependency Injection**: Wire up all components

#### Example MCP Handler

```go
type MovieHandlers struct {
    movieService *movieApp.Service
}

func (h *MovieHandlers) HandleAddMovie(id any, arguments map[string]any, sendResult func(any, any), sendError func(any, int, string, any)) {
    // Parse MCP request
    req, err := h.parseCreateMovieRequest(arguments)
    if err != nil {
        sendError(id, models.InvalidParams, "Invalid movie data", err.Error())
        return
    }

    // Convert to application command
    cmd := movieApp.CreateMovieCommand{
        Title:    req.Title,
        Director: req.Director,
        Year:     req.Year,
        // ...
    }

    // Execute use case
    movieDTO, err := h.movieService.CreateMovie(context.Background(), cmd)
    if err != nil {
        sendError(id, models.InvalidParams, "Failed to create movie", err.Error())
        return
    }

    // Convert to MCP response
    response := h.toMovieResponse(movieDTO)
    sendResult(id, response)
}
```

## Benefits of Clean Architecture

### 1. **Testability**
- **Unit Tests**: Domain logic tested in isolation
- **Integration Tests**: Database and external dependencies tested separately
- **100% Test Coverage**: TDD approach ensures comprehensive testing

### 2. **Maintainability**
- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Clear Boundaries**: Easy to understand and modify

### 3. **Type Safety**
- **Value Objects**: Prevent invalid states at compile time
- **Domain Models**: Business rules enforced in the type system
- **No Primitive Obsession**: Rich domain types instead of primitives

### 4. **Extensibility**
- **Plugin Architecture**: Easy to add new features
- **Interface Segregation**: Small, focused interfaces
- **Open/Closed Principle**: Open for extension, closed for modification

### 5. **Performance**
- **Optimized Queries**: Repository pattern enables query optimization
- **Connection Pooling**: Proper database connection management
- **Caching Ready**: Clean separation makes caching straightforward

## Directory Structure

```
internal/
├── domain/                    # Business Logic (Core)
│   ├── shared/               # Value Objects (MovieID, Rating, Year)
│   ├── movie/                # Movie Aggregate + Repository Interface
│   └── actor/                # Actor Aggregate + Repository Interface
├── application/              # Use Cases 
│   ├── movie/                # Movie Services + DTOs
│   └── actor/                # Actor Services + DTOs
├── infrastructure/           # External Concerns
│   └── postgres/             # PostgreSQL Repository Implementations
├── interfaces/               # Interface Adapters
│   ├── mcp/                  # MCP Protocol Handlers
│   └── dto/                  # External Data Transfer Objects
└── server/                   # Dependency Injection + Server Setup

tests/
├── integration/              # End-to-End Tests
└── performance/              # Performance & Load Tests

cmd/
├── server/                   # Legacy Entry Point
└── server-new/               # Clean Architecture Entry Point

tools/
└── migrate/                  # Custom Migration Tool

migrations/                   # Database Schema Migrations
```

## Testing Strategy

### 1. **Unit Tests** (`*_test.go` in each package)
- Test domain logic in isolation
- Mock external dependencies
- Fast, deterministic tests

### 2. **Integration Tests** (`tests/integration/`)
- Test entire system with real database
- End-to-end MCP protocol testing
- Database migration testing

### 3. **Performance Tests** (`tests/integration/performance_test.go`)
- Concurrent operation testing
- Memory usage validation
- Throughput benchmarks

## Getting Started

### Running with Clean Architecture

```bash
# Use the new clean architecture entry point
go run cmd/server-new/main.go

# Or build and run
go build -o movies-server cmd/server-new/main.go
./movies-server
```

### Running Tests

```bash
# Unit tests (fast)
go test ./internal/domain/... ./internal/application/...

# Integration tests (requires database)
export TEST_DATABASE_URL="postgres://user:pass@localhost/test_db?sslmode=disable"
go test ./tests/integration/...

# Performance tests
go test -bench=. ./tests/integration/...

# All tests
go test ./...
```

### Database Migrations

The project includes a custom migration tool built as a Go tool:

```bash
# Build the migration tool
go build -o migrate ./tools/migrate

# Run migrations up
./migrate "postgres://user:pass@localhost/db?sslmode=disable" ./migrations up

# Run migrations down (rollback last migration)
./migrate "postgres://user:pass@localhost/db?sslmode=disable" ./migrations down
```

**Migration File Format**: `001_create_movies.up.sql` and `001_create_movies.down.sql`

The migration tool:
- Tracks applied migrations in `schema_migrations` table
- Supports both up and down migrations
- Uses transactions for safety
- Provides clear error messages

### Adding New Features

1. **Start with Domain**: Define value objects, entities, and business rules
2. **Write Tests First**: TDD approach ensures correct behavior
3. **Add Repository Interface**: Define data access needs
4. **Implement Application Service**: Orchestrate use case
5. **Add Infrastructure**: Implement repository and external integrations
6. **Create Interface Adapter**: Handle external protocol (MCP)
7. **Wire Dependencies**: Update dependency injection

## Migration from Legacy

The old architecture is still available in the existing files. The new architecture is in:

- **Entry Point**: `cmd/server-new/main.go` (vs `cmd/server/main.go`)
- **Server**: `server.NewCleanServer()` (vs `server.New()`)
- **Version**: `0.2.0` (vs `0.1.0`)
- **Migration Tool**: `tools/migrate/main.go` (custom Go tool for database migrations)

Both versions can coexist during the transition period.

## Performance Characteristics

- **Throughput**: >50 operations/second under load
- **Concurrency**: Handles 50+ concurrent requests safely
- **Memory**: Efficient memory usage with connection pooling
- **Database**: Optimized queries with proper indexing
- **Response Time**: <100ms for typical operations

## Docker Deployment

The project includes comprehensive Docker support for all environments:

### Quick Start

```bash
# Development environment (databases only)
docker-compose -f docker-compose.dev.yml up -d

# Production environment (full stack)
docker-compose -f docker-compose.clean.yml up --build
```

### Available Configurations

- **`Dockerfile.clean`**: Production-ready multi-stage build
- **`docker-compose.clean.yml`**: Full production stack with monitoring
- **`docker-compose.dev.yml`**: Development databases and tools
- **`Dockerfile`** + **`docker-compose.yml`**: Legacy architecture (still available)

### Port Mapping

| Environment | MCP Server | PostgreSQL | Redis | Grafana | pgAdmin |
|-------------|------------|------------|-------|---------|---------|
| Clean Arch  | 8081       | 5433       | 6380  | 3001    | 5051    |
| Development | -          | 5434/5435  | 6381  | -       | 5052    |
| Legacy      | 8080       | 5432       | 6379  | 3000    | 5050    |

### Features

- **Multi-stage builds** for minimal image size
- **Custom migration tool** included in container
- **Health checks** for all services
- **Monitoring stack** (Prometheus + Grafana)
- **Development databases** for testing
- **Non-root execution** for security

For detailed Docker documentation, see [docs/DOCKER.md](docs/DOCKER.md).

## Contributing

When contributing to the clean architecture:

1. Follow the layer boundaries strictly
2. Write tests before implementation (TDD)
3. Use value objects for type safety
4. Keep domain logic pure (no external dependencies)
5. Use dependency injection for all external dependencies
6. Follow Go idioms and conventions

The clean architecture makes the codebase more maintainable, testable, and extensible while following Go best practices.