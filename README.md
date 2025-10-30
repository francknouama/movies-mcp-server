# Movies MCP Server

A Model Context Protocol (MCP) server that provides a comprehensive movie database with advanced search, CRUD operations, and image support. Built with Go and PostgreSQL, this server enables AI-assisted movie management.

## Features

- **Full CRUD Operations**: Create, read, update, and delete movies
- **Advanced Search**: Search by title, genre, director, actors, or any text field
- **Image Support**: Store and retrieve movie posters with base64 encoding
- **Resource Endpoints**: Access database statistics, genre lists, and bulk data
- **Production Ready**: Includes health checks, metrics, logging, and monitoring
- **Docker Support**: Complete Docker setup for easy deployment
- **Comprehensive Testing**: Unit and integration tests with high coverage

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- Docker and Docker Compose
- PostgreSQL 17 (or use Docker)
- Make (optional but recommended)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/francknouama/movies-mcp-server.git
cd movies-mcp-server
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Start the database:
```bash
make docker-up
```

4. Set up the database:
```bash
make db-setup
make db-migrate
make db-seed
```

5. Build the server:
```bash
make build
```

6. Run the server:
```bash
./build/movies-server
```

## MCP Tools

The server implements the following MCP tools:

### Movie Operations
- `get_movie` - Retrieve a movie by ID
- `add_movie` - Add a new movie with optional poster
- `update_movie` - Update movie details including poster
- `delete_movie` - Remove a movie from the database

### Search and Query
- `search_movies` - Advanced search with multiple criteria
- `list_top_movies` - Get top-rated movies with filtering

## MCP Resources

Access movie data through these resource URIs:

- `movies://database/all` - All movies in JSON format
- `movies://database/stats` - Database statistics
- `movies://database/genres` - List of all genres
- `movies://database/directors` - List of all directors
- `movies://posters/{movie-id}` - Individual movie poster
- `movies://posters/collection` - Gallery of all posters

## Usage Examples

### Initialize Connection
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "example-client",
      "version": "1.0.0"
    }
  }
}
```

### Add a Movie
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "add_movie",
    "arguments": {
      "title": "The Matrix",
      "director": "The Wachowskis",
      "release_year": 1999,
      "genre": "Sci-Fi",
      "rating": 8.7,
      "description": "A computer hacker learns about the true nature of reality",
      "poster_url": "https://example.com/matrix-poster.jpg"
    }
  }
}
```

### Search Movies
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "search_movies",
    "arguments": {
      "query": "sci-fi",
      "search_type": "genre"
    }
  }
}
```

## Configuration

The server uses environment variables for configuration. See `.env.example` for all available options:

```bash
# Database configuration
DATABASE_URL=postgres://movies_user:movies_password@localhost:5432/movies_db?sslmode=disable

# Server configuration
LOG_LEVEL=info
SERVER_TIMEOUT=30s
MAX_CONNECTIONS=100

# Image handling
MAX_IMAGE_SIZE_MB=10
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp
```

## Development

### Running Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
DATABASE_URL=... make test-integration
```

### Database Migrations
```bash
# Apply migrations
make db-migrate

# Rollback one migration
make db-migrate-down

# Reset database
make db-migrate-reset
```

### Building
```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all

# Build Docker image
make docker-build
```

## Monitoring

The server includes Prometheus metrics and health endpoints:

- `/health` - Health check endpoint
- `/metrics` - Prometheus metrics

A Grafana dashboard is included in `monitoring/grafana-dashboard.json`.

## Bruno Collection

Interactive API testing is available through the included Bruno collection in `bruno-collection/`. This provides pre-configured requests for all MCP operations.

## Architecture

The server follows clean architecture principles:

```
├── cmd/               # Application entrypoints
├── internal/          # Private application code
│   ├── config/       # Configuration management
│   ├── database/     # Database layer
│   ├── models/       # Domain models
│   └── server/       # MCP server implementation
├── pkg/              # Public packages
│   ├── errors/       # Error handling
│   ├── health/       # Health checks
│   ├── logging/      # Structured logging
│   └── metrics/      # Prometheus metrics
├── migrations/       # Database migrations
└── scripts/          # Utility scripts
```

## Integration with Claude UI

To use the Movies MCP Server with Claude UI:

1. **Add the Server to Claude's Configuration**:
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

2. **Restart Claude UI** to apply the new configuration.

3. **Access Server Features**:
   - Use Claude's UI to make requests to `movies-mcp-server`.
   - All MCP tools, such as `search_movies` and `add_movie`, are supported.

For more details, refer to the [Integration Guide](docs/getting-started/claude-desktop.md).

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built for the [Model Context Protocol](https://modelcontextprotocol.io/) ecosystem
- PostgreSQL for reliable data storage
- The Go community for excellent libraries and tools
