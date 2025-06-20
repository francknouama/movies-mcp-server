# Movies MCP Server - Bruno Collection

This Bruno collection provides comprehensive testing for the Movies MCP Server, covering all available tools, resources, and error scenarios.

## Collection Structure

### 01 - Protocol Initialization
- **Initialize.bru**: MCP protocol initialization handshake
- **Ping.bru**: Server connectivity and health check
- **List Tools.bru**: Discover all available tools

### 02 - Tool Discovery  
- Tests MCP protocol compliance
- Validates tool descriptions and schemas
- Ensures proper capability advertising

### 03 - Movie Operations
- **Add Movie.bru**: Create new movie entries
- **Get Movie.bru**: Retrieve movies by ID
- **Update Movie.bru**: Modify existing movie data
- **Delete Movie.bru**: Remove movies from database

### 04 - Search and Query
- **Search Movies.bru**: Full-text and filtered search
- **List Top Movies.bru**: Retrieve highest-rated movies

### 05 - Resources
- **List Resources.bru**: Discover available movie resources
- **Read Resource.bru**: Access individual movie resources

### 06 - Error Cases
- **Invalid Tool Call.bru**: Non-existent tool handling
- **Invalid Parameters.bru**: Parameter validation testing
- **Movie Not Found.bru**: Missing resource error handling

### 07 - Batch Operations
- **Bulk Add Movies.bru**: Add multiple movies in sequence
- **Add Second Movie.bru**: Continuation of bulk operations
- **Verify Bulk Addition.bru**: Validate bulk operation results

## Environment Configuration

### Development Environment
- **Base URL**: `http://localhost:8080`
- **Protocol**: MCP over HTTP (for testing convenience)
- **Database**: Local PostgreSQL instance

### Production Environment  
- **Base URL**: `https://your-production-domain.com`
- **Protocol**: HTTPS with proper certificates
- **Database**: Production PostgreSQL instance

## Variables

The collection uses several dynamic variables:

- `{{base_url}}`: Server endpoint (from environment)
- `{{movie_id}}`: Movie ID captured from responses
- `{{inception_id}}`: Specific ID for Inception movie
- `{{dark_knight_id}}`: Specific ID for Dark Knight movie

## Running the Tests

### Prerequisites
1. Movies MCP Server running on configured port
2. Database initialized with proper schema
3. Bruno API client installed

### Execution Order
1. **Protocol Initialization**: Run first to establish connection
2. **Tool Discovery**: Verify available capabilities  
3. **Movie Operations**: Test CRUD functionality
4. **Search and Query**: Test search capabilities
5. **Resources**: Test MCP resource protocol
6. **Error Cases**: Validate error handling
7. **Batch Operations**: Test complex workflows

### Test Validation
Each request includes:
- **Status Checks**: HTTP 200 responses
- **Protocol Validation**: JSON-RPC 2.0 compliance
- **Data Validation**: Response content verification
- **Business Logic**: Domain-specific validations

## MCP Protocol Notes

This collection tests the MCP (Model Context Protocol) over HTTP for convenience. In production, MCP typically runs over stdin/stdout. The protocol mapping:

- **tools/list** → List available tools
- **tools/call** → Execute specific tool
- **resources/list** → List available resources  
- **resources/read** → Read resource content

## Error Handling

The server implements standard JSON-RPC error codes:
- **-32601**: Method not found (invalid tool)
- **-32602**: Invalid parameters (validation errors)
- **-32603**: Internal error (server issues)

## Sample Data

The collection includes sample data for popular movies:
- The Matrix (1999) - Sci-Fi classic
- Inception (2010) - Christopher Nolan film
- The Dark Knight (2008) - Batman sequel

## Customization

To adapt this collection for your environment:

1. Update environment variables (base_url, credentials)
2. Modify sample data in request bodies
3. Adjust validation expectations in tests
4. Add custom error scenarios as needed

## Troubleshooting

Common issues:
- **Connection refused**: Check server is running on correct port
- **404 errors**: Verify MCP endpoint path is correct
- **Validation errors**: Check request parameter formats
- **Database errors**: Ensure database is accessible and initialized