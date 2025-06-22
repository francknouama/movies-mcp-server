# MCP Servers Workspace

This repository contains multiple Model Context Protocol (MCP) servers organized as a Go workspace.

## Structure

```
.
├── go.work                 # Go workspace configuration
├── mcp-server/            # Movies database MCP server
│   ├── go.mod
│   ├── cmd/               # Application entrypoints
│   ├── internal/          # Private application code
│   └── migrations/        # Database migrations
├── godog-server/          # Cucumber/Godog testing MCP server
│   ├── go.mod
│   ├── cmd/               # Application entrypoints
│   ├── internal/          # Private application code
│   └── step_definitions/  # Godog step definitions
└── shared-mcp/            # Shared MCP utilities and libraries
    ├── go.mod
    └── pkg/               # Shared packages
```

## Servers

### Movies MCP Server
A comprehensive movie database server with advanced search, CRUD operations, and image support. Built with Go and PostgreSQL.

[Full documentation →](./mcp-server/README.md)

### Godog MCP Server
A Cucumber BDD testing server that enables AI assistants to run and manage Godog tests through the MCP protocol.

[Full documentation →](./godog-server/README.md)

### Shared MCP
A shared library containing common utilities, database abstractions, and shared functionality used across multiple MCP servers.

## Development

This repository uses Go workspaces (introduced in Go 1.18) to manage multiple modules. Each server is a separate Go module that can be developed and versioned independently. The shared-mcp module provides common functionality used across servers.

### Prerequisites
- Go 1.24.4 or later
- Docker and Docker Compose (for mcp-server)
- Godog CLI (for godog-server)

### Building

Build all servers:
```bash
make build-all
```

Build specific server:
```bash
# Movies MCP server
cd mcp-server && go build -o movies-mcp-server cmd/server/main.go

# Godog server
cd godog-server && go build -o godog-server cmd/server/main.go
```

### Testing

Run all tests:
```bash
make test-all
```

### Adding a New Server

1. Create a new directory for your server
2. Initialize a Go module: `cd new-server && go mod init your-module-name`
3. Add it to the workspace: `go work use ./new-server`
4. Follow the existing server patterns for structure

## Claude Desktop Integration

To use the Movies MCP Server with Claude Desktop, you need to configure it in your Claude Desktop settings.

### Configuration

1. **Build the Movies MCP Server**:
   ```bash
   cd mcp-server && go build -o movies-mcp-server cmd/server/main.go
   ```

2. **Set up the database** (see [mcp-server/README.md](./mcp-server/README.md) for detailed instructions):
   ```bash
   cd mcp-server
   docker-compose up -d  # Start PostgreSQL
   make migrate-up       # Run database migrations
   ```

3. **Add to Claude Desktop config** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):
   ```json
   {
     "mcpServers": {
       "movies-mcp-server": {
         "command": "/absolute/path/to/movies-mcp-server/mcp-server/movies-mcp-server",
         "args": [],
         "env": {
           "DATABASE_URL": "postgres://movies_user:movies_password@localhost:5432/movies_db"
         }
       }
     }
   }
   ```

4. **Restart Claude Desktop** to load the new configuration.

### Usage

Once configured, you can interact with Claude Desktop to:
- Search and manage movies in the database
- Add, update, and delete movie records
- Search by title, genre, director, or year
- Manage actors and their relationships with movies
- Get movie recommendations and statistics

### Troubleshooting

- Ensure the server binary has execute permissions: `chmod +x mcp-server/movies-mcp-server`
- Verify the absolute path in the configuration is correct
- Check that PostgreSQL is running and accessible
- Ensure the DATABASE_URL environment variable matches your database setup
- Check Claude Desktop logs for connection issues

## Contributing

See individual server READMEs for specific contribution guidelines.

## License

Each server may have its own license. Check the LICENSE file in each server directory.