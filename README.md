# MCP Servers Workspace

This repository contains multiple Model Context Protocol (MCP) servers organized as a Go workspace.

## Structure

```
.
├── go.work                 # Go workspace configuration
├── movies-mcp-server/      # Movies database MCP server
│   ├── go.mod
│   ├── cmd/               # Application entrypoints
│   ├── internal/          # Private application code
│   └── pkg/               # Public packages
└── godog-server/          # Cucumber/Godog testing MCP server
    ├── go.mod
    ├── cmd/               # Application entrypoints
    ├── internal/          # Private application code
    └── pkg/               # Public packages
```

## Servers

### Movies MCP Server
A comprehensive movie database server with advanced search, CRUD operations, and image support. Built with Go and PostgreSQL.

[Full documentation →](./movies-mcp-server/README.md)

### Godog MCP Server
A Cucumber BDD testing server that enables AI assistants to run and manage Godog tests through the MCP protocol.

[Full documentation →](./godog-server/README.md)

## Development

This repository uses Go workspaces (introduced in Go 1.18) to manage multiple modules. Each server is a separate Go module that can be developed and versioned independently.

### Prerequisites
- Go 1.24.4 or later
- Docker and Docker Compose (for movies-server)
- Godog CLI (for godog-server)

### Building

Build all servers:
```bash
make build-all
```

Build specific server:
```bash
# Movies MCP server
cd movies-mcp-server && go build -o movies-mcp-server cmd/server/main.go

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

## Contributing

See individual server READMEs for specific contribution guidelines.

## License

Each server may have its own license. Check the LICENSE file in each server directory.