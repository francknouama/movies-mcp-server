package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/francknouama/movies-mcp-server/mcp-server/internal/composition"
	"github.com/francknouama/movies-mcp-server/mcp-server/internal/config"
)

// MoviesServer is a compatibility wrapper around MCPServer
// This maintains backward compatibility while using the new clean architecture
type MoviesServer struct {
	mcpServer *MCPServer
	db        *sql.DB
}

// NewServer creates a new MoviesServer instance with clean architecture
func NewServer(db *sql.DB) *MoviesServer {
	logger := log.New(os.Stderr, "[movies-mcp] ", log.LstdFlags)
	container := composition.NewContainer(db)
	mcpServer := NewMCPServer(os.Stdin, os.Stdout, logger, container)

	return &MoviesServer{
		mcpServer: mcpServer,
		db:        db,
	}
}

// NewServerWithConfig creates a new MoviesServer instance with custom config and clean architecture
func NewServerWithConfig(db *sql.DB, cfg *config.Config) *MoviesServer {
	logger := log.New(os.Stderr, "[movies-mcp] ", log.LstdFlags)
	container := composition.NewContainer(db)
	mcpServer := NewMCPServer(os.Stdin, os.Stdout, logger, container)

	return &MoviesServer{
		mcpServer: mcpServer,
		db:        db,
	}
}

// Run starts the server and handles incoming requests
func (s *MoviesServer) Run() error {
	return s.mcpServer.Run()
}
