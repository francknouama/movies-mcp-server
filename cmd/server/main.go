package main

import (
	"flag"
	"fmt"
	"os"

	"movies-mcp-server/internal/config"
	"movies-mcp-server/internal/database"
	"movies-mcp-server/internal/server"
)

const (
	version = "0.1.0"
	name    = "movies-mcp-server"
)

func main() {
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
	)
	
	flag.Parse()
	
	if *showVersion {
		fmt.Printf("%s version %s\n", name, version)
		os.Exit(0)
	}
	
	if *showHelp {
		fmt.Printf("Movies MCP Server - A Model Context Protocol server for movie database\n\n")
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nThe server communicates via stdin/stdout using the MCP protocol.\n")
		os.Exit(0)
	}
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Connect to database
	db, err := database.NewPostgresDatabase(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	
	// Test database connection
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Fprintf(os.Stderr, "Connected to database: %s\n", cfg.Database.Name)
	
	// Create and run the server
	srv := server.New(db)
	
	if err := srv.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}