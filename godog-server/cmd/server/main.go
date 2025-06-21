package main

import (
	"shared-mcp/pkg/logging"

	"github.com/francknouama/movies-mcp-server/godog-server/internal/config"
	"github.com/francknouama/movies-mcp-server/godog-server/internal/server"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	godogServer := server.NewServer(cfg, logger)

	logger.Info("Starting Godog MCP Server...")

	if err := godogServer.Run(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
