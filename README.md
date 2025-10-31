# Movies MCP Server

A production-ready **Model Context Protocol (MCP) server** for intelligent movie database management, built with Clean Architecture principles and optimized for AI-assisted environments.

## What is Movies MCP Server?

Movies MCP Server is a sophisticated movie database management system that communicates via the **Model Context Protocol**—designed specifically for integration with AI assistants like Claude. Unlike traditional HTTP APIs, it uses JSON-RPC over stdin/stdout to provide seamless, intelligent movie and actor data operations.

**Perfect for:**
- AI-powered movie recommendation systems
- Claude Desktop integrations
- Intelligent film analysis and exploration
- Director career research
- Movie database management with AI assistance

---

## Why Choose Movies MCP Server?

- **MCP Protocol Native**: Built specifically for the Model Context Protocol, not an HTTP API wrapper
- **Clean Architecture**: Exemplary separation of concerns with domain-driven design
- **Intelligent Features**: AI-powered recommendations, director career analysis, and similarity searches
- **Comprehensive Actor Management**: Full actor database with movie associations and career tracking
- **Production-Ready**: Health checks, Prometheus metrics, Grafana dashboards, and comprehensive monitoring
- **Type-Safe Domain Models**: Value objects prevent invalid states at compile time
- **Advanced Search**: Full-text search, decade filtering, rating ranges, genre matching, and similarity scoring
- **Image Support**: Store and retrieve movie posters via MCP resources with base64 encoding
- **BDD Testing**: Comprehensive test coverage with Cucumber/Godog behavior scenarios
- **Docker-Optimized**: Multi-stage builds, distroless images, non-root execution

---

## Key Performance Metrics

- **Throughput**: >50 operations/second under load
- **Concurrency**: Safely handles 50+ concurrent requests
- **Response Time**: <100ms for typical operations
- **Test Coverage**: Comprehensive unit and integration tests with BDD scenarios

---

## MCP Capabilities

### 22 Available Tools

#### Movie Management (5 tools)
- `get_movie` - Retrieve movie by ID
- `add_movie` - Create movie with title, director, year, rating, genres, poster
- `update_movie` - Update existing movie details
- `delete_movie` - Delete movie by ID
- `list_top_movies` - Get top-rated movies with configurable limit

#### Actor Management (9 tools)
- `add_actor` - Create actor with name, birth year, biography
- `get_actor` - Retrieve actor by ID
- `update_actor` - Update actor information
- `delete_actor` - Delete actor
- `link_actor_to_movie` - Associate actor with movie
- `unlink_actor_from_movie` - Remove actor-movie association
- `get_movie_cast` - Get all actors in a movie
- `get_actor_movies` - Get all movies for an actor
- `search_actors` - Search actors by name with birth year filtering

#### Advanced Search (4 tools)
- `search_movies` - Multi-criteria search (title, director, genre, year range, rating)
- `search_by_decade` - Find movies from specific decades (1990s, 2000s, etc.)
- `search_by_rating_range` - Filter movies by rating boundaries
- `search_similar_movies` - Find similar movies based on genres and ratings

#### Intelligence & Analysis (4 tools)
- `bulk_movie_import` - Import multiple movies with error tracking
- `movie_recommendation_engine` - AI-powered recommendations with preference scoring
- `director_career_analysis` - Career trajectory with early/mid/late phase analysis
- Pagination context management for large datasets

### 5 Built-in Prompts

- **movie_recommendation** - Generate personalized recommendations based on preferences
- **movie_analysis** - Analyze themes, cinematography, and characteristics
- **director_filmography** - Explore director's body of work and evolution
- **genre_exploration** - Deep dive into genre history and influential films
- **movie_comparison** - Compare two movies across multiple dimensions

### 3 MCP Resources

- `movies://database/all` - Complete movie database in JSON format
- `movies://database/stats` - Database statistics and analytics
- `movies://posters/collection` - All movie posters (base64 encoded)
- Dynamic: `movies://posters/{movie-id}` - Individual movie posters

---

## Architecture & Technology

### Clean Architecture Implementation

Built with strict separation of concerns:

```
internal/
├── domain/          # Pure business logic (entities, value objects)
├── application/     # Use cases and orchestration
├── infrastructure/  # Database and external integrations
├── interfaces/      # MCP protocol adapters
└── composition/     # Dependency injection
```

**Benefits:**
- Framework independence
- Testable business logic
- Database agnostic (currently PostgreSQL)
- Easy to maintain and extend

### Technology Stack

**Core:**
- Go 1.23.0+ with Go 1.24.4 toolchain
- PostgreSQL 17 with advanced indexing
- Model Context Protocol (MCP) via JSON-RPC

**Key Libraries:**
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/cucumber/godog` - BDD testing
- `github.com/testcontainers/testcontainers-go` - Integration testing
- `github.com/sirupsen/logrus` - Structured logging
- OpenTelemetry - Distributed tracing

**Database Features:**
- Full-text search (GIN indexes)
- Array-based genre filtering
- Many-to-many actor-movie relationships
- Automatic timestamp management
- Image storage (BYTEA columns)

---

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- Docker and Docker Compose
- PostgreSQL 17 (or use Docker setup)
- Make (optional, for convenience)

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/francknouama/movies-mcp-server.git
   cd movies-mcp-server
   ```

2. **Set Up Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Database with Docker**:
   ```bash
   make docker-up
   ```

4. **Initialize Database**:
   ```bash
   make db-setup      # Create database
   make db-migrate    # Run migrations
   make db-seed       # Load sample data
   ```

5. **Build the Server**:
   ```bash
   make build
   ```

6. **Run the Server**:
   ```bash
   ./build/movies-server
   ```

### Docker Deployment

**Development (databases only):**
```bash
docker-compose -f docker-compose.dev.yml up
```

**Production (with monitoring):**
```bash
docker-compose -f docker-compose.clean.yml up
```

**Included Services:**
- PostgreSQL 17 (port 5432)
- Movies MCP Server
- Grafana (port 3000)
- pgAdmin (port 5050)
- Prometheus (port 9090)

---

## Integration with Claude Desktop

Configure Claude Desktop to use Movies MCP Server:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "movies-mcp-server": {
      "command": "/path/to/movies-server",
      "env": {
        "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5432/movies_db"
      }
    }
  }
}
```

**Restart Claude Desktop** to activate the integration.

### What You Can Do with Claude

- "Find me thriller movies from the 1990s with ratings above 8"
- "Add a new movie: Inception, directed by Christopher Nolan, released in 2010"
- "Show me all movies starring Leonardo DiCaprio"
- "Analyze Quentin Tarantino's career trajectory"
- "Recommend movies similar to The Godfather"
- "Import this list of movies in bulk"

---

## Advanced Features

### Intelligent Recommendation Engine

Multi-factor scoring algorithm:
- **Genre matching** (40% weight)
- **Rating score** (30% weight)
- **Year relevance** (20% weight)
- **Popularity boost** (10% weight)

Returns ranked recommendations with match scores and reasoning.

### Director Career Analysis

Automatic analysis includes:
- Career phase detection (early/mid/late)
- Average rating per phase
- Genre specialization tracking
- Career trajectory (ascending/descending/peak/resurgence)
- Notable works (best and worst rated)

### Bulk Import Operations

Import multiple movies at once with:
- Per-item error tracking
- Success/failure statistics
- Partial success handling
- Detailed error reporting

### Advanced Search Capabilities

- **Full-text search**: Title, director, description using PostgreSQL GIN indexes
- **Decade parsing**: Intelligently handles "1990s", "90s", "1990" formats
- **Similarity scoring**: Genre and rating-based recommendations
- **Multi-criteria filtering**: Combine title, genre, year range, rating range
- **Pagination support**: Handle large result sets efficiently

---

## Monitoring & Observability

### Prometheus Metrics

Available at port **9090** with comprehensive metrics:
- Request/response times
- Concurrent operations
- Database connection pool stats
- Query performance
- Memory and CPU utilization

### Grafana Dashboards

Access Grafana at port **3000** for:
- Real-time performance monitoring
- Database health visualization
- Custom alerting rules
- System resource tracking

### Health Checks

Built-in health checks with:
- Configurable intervals (default: 30s)
- Database connectivity verification
- Graceful degradation
- Status reporting

### Alert Rules

Pre-configured alerts for:
- High error rates
- Slow query performance
- Database connection issues
- Memory/CPU thresholds

Configuration: `monitoring/alert_rules.yml`

---

## Developer Guide

### Testing

**Run All Tests:**
```bash
make test                  # Unit tests
make test-integration      # Integration tests with testcontainers
make test-coverage         # Coverage report
make test-bdd              # BDD scenarios with Godog
```

**BDD Feature Tests:**
- 40+ behavior scenarios in Gherkin
- Real PostgreSQL via testcontainers
- Contract testing for MCP protocol
- Performance and load testing

### Database Migrations

```bash
make db-migrate            # Apply migrations
make db-migrate-down       # Rollback last migration
make db-migrate-reset      # Reset database
make db-create-migration   # Create new migration
```

### Code Quality

```bash
make fmt                   # Format code
make vet                   # Run go vet
make lint                  # Run golangci-lint
```

### Build Options

```bash
make build                 # Build for current platform
make build-all             # Multi-platform builds
make docker-build          # Build Docker image
make release               # Create release with goreleaser
```

### Environment Variables

**Database:**
- `DATABASE_URL` - Full connection string
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `DB_MAX_CONNECTIONS=100`, `DB_MAX_IDLE_CONNECTIONS=10`

**Server:**
- `PORT=8080`, `METRICS_PORT=9090`
- `READ_TIMEOUT=30s`, `WRITE_TIMEOUT=30s`
- `LOG_LEVEL` (debug/info/warn/error)

**Security:**
- `JWT_SECRET`, `API_KEY`
- `RATE_LIMIT=1000` (per minute per IP)
- `TLS_ENABLED`, `TLS_CERT_FILE`, `TLS_KEY_FILE`

**Monitoring:**
- `PROMETHEUS_ENABLED=true`
- `HEALTH_CHECK_INTERVAL=30s`

See `.env.example` for complete configuration options.

---

## Documentation

Comprehensive documentation available in the `/docs` directory:

**Getting Started:**
- [Installation Guide](docs/getting-started/installation.md)
- [Your First Movie](docs/getting-started/first-movie.md)
- [Claude Desktop Integration](docs/getting-started/claude-desktop.md)

**Guides:**
- [User Guide](docs/guides/user-guide.md) - Feature walkthrough
- [Examples](docs/guides/examples.md) - Code examples

**Architecture:**
- [Architecture Overview](ARCHITECTURE.md) - Clean Architecture details
- [Docker Guide](DOCKER.md) - Docker setup and configuration
- [Deployment Guide](DEPLOYMENT.md) - Production deployment
- [Image Support](IMAGE_SUPPORT.md) - Image handling via MCP

**Reference:**
- [API Reference](docs/reference/api.md) - Complete API documentation
- [Troubleshooting](docs/reference/troubleshooting.md) - Common issues
- [FAQ](docs/appendices/faq.md) - Frequently asked questions

---

## Project Structure

```
movies-mcp-server/
├── cmd/server/              # CLI entry points
├── internal/
│   ├── domain/              # Business logic (entities, value objects)
│   ├── application/         # Use cases and services
│   ├── infrastructure/      # Database and integrations
│   ├── interfaces/          # MCP protocol handlers
│   ├── schemas/             # Tool definitions
│   └── server/              # MCP server core
├── migrations/              # Database migrations
├── tests/                   # Comprehensive test suites
│   └── bdd/                # BDD feature files
├── docs/                    # Documentation
├── monitoring/              # Prometheus and Grafana configs
└── docker/                  # Docker configurations
```

---

## Contributing

We welcome contributions! Please see the [Contributing Guide](docs/development/README.md) for:
- Code of conduct
- Development setup
- Pull request process
- Coding standards
- Testing requirements

---

## Support & Community

- **Found a bug?** [Report an Issue](https://github.com/francknouama/movies-mcp-server/issues)
- **Have questions?** Check the [FAQ](docs/appendices/faq.md)
- **Need help?** See the [Troubleshooting Guide](docs/reference/troubleshooting.md)
- **Want to contribute?** Read the [Development Guide](docs/development/README.md)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

Special thanks to:

- [Model Context Protocol](https://modelcontextprotocol.io) for the MCP ecosystem
- [Anthropic](https://www.anthropic.com) for Claude and MCP development
- PostgreSQL community for the robust database
- Go community for excellent tools and libraries
- All contributors and users of this project

---

## What's Next?

See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for the roadmap including:
- GraphQL integration
- Advanced caching strategies
- Enhanced recommendation algorithms
- Multi-language support
- Real-time notifications

---

**Built with Clean Architecture principles for maintainability, testability, and scalability.**
