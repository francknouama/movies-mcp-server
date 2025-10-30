# Movies MCP Server

## Welcome to the Movies MCP Server 🎥✨

The Movies MCP Server is your go-to solution for managing movie databases with speed, flexibility, and scalability. Designed for AI-assisted environments, it provides advanced search capabilities, image handling, and integration with Claude UI, all while being production-ready and developer-friendly.

---

## Why Choose Movies MCP Server?

- 🔄 **Complete Movie Management**: Add, edit, delete, and search movies effortlessly.
- 🔍 **Advanced Search**: Find movies by title, genre, director, or any text field.
- 🌟 **AI-Optimized**: Built for seamless integration with AI systems like Claude UI.
- 🖼️ **Image Handling**: Store and retrieve movie posters with ease.
- 🚀 **Production-Ready**: Health checks, metrics, and monitoring included.
- 🐳 **Docker-Friendly**: Simplified deployment with Docker.
- 🛡️ **Secure & Scalable**: Enterprise-level standards for reliability and performance.

---

## Key Features

| Feature                | Description                                               |
|------------------------|-----------------------------------------------------------|
| **Full CRUD Operations** | Manage movie records with Create, Read, Update, Delete.   |
| **Advanced Search**      | Search movies by multiple criteria (title, genre, etc.).  |
| **Image Support**        | Upload and retrieve posters in base64 format.            |
| **API Endpoints**        | Access database stats, genres, directors, and more.      |
| **Monitoring**           | Prometheus metrics and Grafana dashboard support.        |

---

## Quick Start 🚀

### Prerequisites

Ensure you have the following:

- Go 1.24.4 or later
- Docker and Docker Compose
- PostgreSQL 17 (or use Docker-based setup)
- Make (optional, for easier commands)

### Installation Steps

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/francknouama/movies-mcp-server.git
   cd movies-mcp-server
   ```

2. **Prepare Environment**:
   ```bash
   cp .env.example .env
   ```

3. **Start the Database**:
   ```bash
   make docker-up
   ```

4. **Set Up the Database**:
   ```bash
   make db-setup
   make db-migrate
   make db-seed
   ```

5. **Build the Server**:
   ```bash
   make build
   ```

6. **Run the Server**:
   ```bash
   ./build/movies-server
   ```

---

## Integration with Claude UI 🤖

Easily integrate the Movies MCP Server with Claude UI:

1. **Update Claude Configuration**:
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

2. **Restart Claude UI** to apply changes.

3. **Enjoy Full Features**: Use Claude's interface to query, add, and manage movies.

---

## Developer Notes 🛠️

### Testing

Run tests to ensure everything works as expected:
```bash
make test
make test-coverage
make test-integration
```

### Database Migrations

Manage database migrations with ease:
```bash
make db-migrate
make db-migrate-down
make db-migrate-reset
```

### Build Options

Create builds for deployment:
```bash
make build
make build-all
make docker-build
```

---

## Need Help? 🤔

- **Found a bug?** [Report an Issue](https://github.com/francknouama/movies-mcp-server/issues)
- **Have questions?** Check out the [FAQ](docs/appendices/faq.md) or [User Guide](docs/guides/user-guide.md).
- **Want to contribute?** See the [Contributing Guide](docs/development/README.md).

---

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments 💡

Special thanks to:

- [Model Context Protocol](https://modelcontextprotocol.io) for the MCP ecosystem.
- PostgreSQL for its robust database capabilities.
- The Go community for providing excellent tools and libraries.
