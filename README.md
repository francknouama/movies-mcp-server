# 🎬 Movies MCP Server
> A comprehensive movie database for AI assistants via Model Context Protocol

Transform your AI assistant into a powerful movie database manager with full CRUD operations, advanced search, and image support.

## ⚡ Quick Start

**New to MCP?** → [5-minute setup](docs/getting-started/README.md)  
**Claude Desktop user?** → [Integration guide](docs/getting-started/claude-desktop.md)  
**Developer?** → [Development setup](docs/development/README.md)  
**Production deployment?** → [Deployment guide](docs/deployment/README.md)

## 🎯 What You Can Do

- 🔍 **Search & Browse**: Find movies by title, genre, director, or plot keywords
- 📝 **Manage Collection**: Add, update, and organize your personal movie database  
- 🖼️ **Handle Images**: Store and retrieve movie posters with automatic processing
- 📊 **Get Insights**: Database statistics, top-rated films, and smart recommendations
- 🎭 **Track People**: Manage actors, directors, and their filmographies

## 🏗️ Architecture Options

| Version | Status | Best For | Migration Tool |
|---------|--------|----------|----------------|
| [**Clean Architecture**](docs/development/architecture.md) | ✅ **Recommended** | Production, new projects | 🔄 Built-in (automatic) |
| [Legacy](docs/appendices/migration-guide.md) | 🔄 Maintenance | Existing integrations | ⚠️ External (manual) |

## 🚀 Repository Structure

This workspace contains multiple Model Context Protocol (MCP) servers:

```
📁 movies-mcp-server/
├── 📁 mcp-server/           # 🎬 Movies database server (main)
├── 📁 godog-server/         # 🧪 Cucumber/BDD testing server  
├── 📁 shared-mcp/           # 📚 Shared MCP utilities
└── 📁 docs/                 # 📖 User-centered documentation
    ├── 🚀 getting-started/  # Quick setup guides
    ├── 📖 guides/           # User manuals & examples
    ├── 🔧 development/      # Developer resources
    ├── 🚢 deployment/       # Production deployment
    ├── 🔍 reference/        # API & configuration
    └── 📊 appendices/       # FAQ, migration, performance
```

### 🎬 Movies MCP Server (Primary)
Full-featured movie database with clean architecture, PostgreSQL, and comprehensive MCP tool suite.

**→ [Complete Guide](docs/getting-started/README.md)**

### 🧪 Godog MCP Server  
Cucumber BDD testing integration for AI-driven test management and execution.

**→ [Godog Documentation](./godog-server/README.md)**

## Development

This repository uses Go workspaces (introduced in Go 1.18) to manage multiple modules. Each server is a separate Go module that can be developed and versioned independently. The shared-mcp module provides common functionality used across servers.

### Prerequisites
- Go 1.23 or later
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