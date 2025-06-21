package server

import (
	"io"
	"os"

	"shared-mcp/pkg/logging"

	"github.com/francknouama/movies-mcp-server/godog-server/internal/config"
	"github.com/francknouama/movies-mcp-server/godog-server/internal/godog"
)

// GodogServer is the main MCP server for Godog testing
type GodogServer struct {
	mcpServer   *MCPServer
	godogRunner *godog.Runner
	config      *config.Config
	logger      *logging.Logger
}

// NewServer creates a new GodogServer instance
func NewServer(cfg *config.Config, logger *logging.Logger) *GodogServer {
	// Create Godog runner
	godogRunner := godog.NewRunner(cfg, logger)

	// Create MCP server components
	mcpServer := NewMCPServer(os.Stdin, os.Stdout, logger, godogRunner)

	return &GodogServer{
		mcpServer:   mcpServer,
		godogRunner: godogRunner,
		config:      cfg,
		logger:      logger,
	}
}

// NewServerWithIO creates a new GodogServer with custom IO streams
func NewServerWithIO(cfg *config.Config, logger *logging.Logger, input io.Reader, output io.Writer) *GodogServer {
	// Create Godog runner
	godogRunner := godog.NewRunner(cfg, logger)

	// Create MCP server components
	mcpServer := NewMCPServer(input, output, logger, godogRunner)

	return &GodogServer{
		mcpServer:   mcpServer,
		godogRunner: godogRunner,
		config:      cfg,
		logger:      logger,
	}
}

// Run starts the server and handles incoming requests
func (s *GodogServer) Run() error {
	s.logger.Info("Starting Godog MCP Server")

	// Check Godog availability
	if err := s.godogRunner.CheckAvailability(); err != nil {
		s.logger.WithField("error", err).Error("Godog not available")
		return err
	}

	s.logger.WithField("godog_binary", s.config.GodogBinary).Info("Godog binary detected")

	// Start the MCP server
	return s.mcpServer.Run()
}

// GetConfig returns the server configuration
func (s *GodogServer) GetConfig() *config.Config {
	return s.config
}

// GetGodogRunner returns the Godog runner instance
func (s *GodogServer) GetGodogRunner() *godog.Runner {
	return s.godogRunner
}
