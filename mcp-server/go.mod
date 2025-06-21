module github.com/francknouama/movies-mcp-server/mcp-server

go 1.24.4

require (
	github.com/francknouama/movies-mcp-server/shared-mcp v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

replace github.com/francknouama/movies-mcp-server/shared-mcp => ../shared-mcp
