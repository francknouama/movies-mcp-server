# ⚠️ DEPRECATED: Legacy MCP Server

## Status: Deprecated

This legacy MCP server implementation has been **deprecated** in favor of the new SDK-based server.

**New Server Location:** [`cmd/server-sdk/`](../server-sdk/)

## Why Deprecated?

The legacy server has been replaced with an implementation using the official [Golang MCP SDK v1.1.0](https://github.com/modelcontextprotocol/go-sdk), which provides:

- ✅ Official SDK integration
- ✅ Type-safe tool handlers
- ✅ Automatic schema generation
- ✅ 26% less code for same functionality
- ✅ Better maintainability
- ✅ 100% feature parity validated

## Migration Status

| Component | Legacy | SDK | Status |
|-----------|--------|-----|--------|
| **Tools** | 23 | 23 | ✅ Migrated |
| **Resources** | 3 | 3 | ✅ Migrated |
| **Tests** | Partial | 58 tests | ✅ Complete |
| **CI/CD** | Tested | Tested | ✅ Both validated |

## Should I Use This?

**No.** You should use the SDK-based server instead.

### Use SDK Server Instead

```bash
# Build and run SDK server
go build -o movies-mcp-server-sdk ./cmd/server-sdk/main.go
./movies-mcp-server-sdk

# Or use with Claude Desktop
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-mcp-server-sdk"
    }
  }
}
```

See [`cmd/server-sdk/README.md`](../server-sdk/README.md) for complete documentation.

## When Will This Be Removed?

The legacy server is maintained for:

1. **BDD Testing**: CI runs tests against both servers to ensure feature parity
2. **Backwards Compatibility**: Users can still run legacy server if needed
3. **Comparison**: Developers can compare implementations

**Timeline:**
- **Now**: Deprecated (use SDK server for new deployments)
- **Q1 2025**: May be archived after SDK server proven in production
- **Future**: Will be removed when no longer needed for testing

## Feature Comparison

### SDK Server Advantages

| Feature | Legacy | SDK |
|---------|--------|-----|
| **Code Size** | Baseline | -26% |
| **Test Code** | Baseline | -37% |
| **Type Safety** | Manual | Automatic |
| **Schema Generation** | Manual | Automatic |
| **Maintenance** | High | Low |
| **SDK Version** | Custom | Official v1.1.0 |
| **Test Coverage** | Partial | 58 tests (100%) |

### Identical Functionality

Both servers provide:
- ✅ 23 MCP tools (movie, actor, compound, context)
- ✅ 3 MCP resources (database, statistics, posters)
- ✅ PostgreSQL integration
- ✅ Database migrations
- ✅ Clean Architecture
- ✅ BDD test validated

## Migration Guide

### For Developers

If you're maintaining or modifying this codebase:

1. **New Features**: Add to SDK server only (`cmd/server-sdk/`)
2. **Bug Fixes**: Fix in SDK server first, legacy if needed
3. **Refactoring**: Focus on SDK server
4. **Documentation**: Update SDK server docs

### For Deployments

If you're deploying this server:

1. **New Deployments**: Use SDK server
2. **Existing Deployments**: Plan migration to SDK server
3. **Testing**: Validate with BDD tests (`TEST_MCP_SERVER=sdk`)

### Migration Steps

```bash
# 1. Test SDK server locally
TEST_MCP_SERVER=sdk go test -v ./tests/bdd/...

# 2. Build SDK server
go build -o movies-mcp-server-sdk ./cmd/server-sdk/main.go

# 3. Update configuration
# Replace: movies-mcp-server
# With:    movies-mcp-server-sdk

# 4. Deploy and monitor
./movies-mcp-server-sdk

# 5. Verify functionality
# All 23 tools + 3 resources should work identically
```

## CI/CD Status

The CI pipeline validates both implementations:

```yaml
bdd-tests:
  strategy:
    matrix:
      server: [legacy, sdk]  # Both tested in parallel
```

This ensures the SDK server maintains 100% behavioral compatibility with the legacy server.

## Support

- **SDK Server Support**: ✅ Active development
- **Legacy Server Support**: ⚠️ Maintenance only (critical bugs only)

## Additional Resources

- [SDK Server Documentation](../server-sdk/README.md)
- [SDK Migration Comparison](../../docs/SDK_MIGRATION_COMPARISON.md)
- [SDK Migration Complete Summary](../../docs/SDK_MIGRATION_COMPLETE.md)
- [Testing Comparison](../../docs/TESTING_COMPARISON.md)
- [CI/CD Enhancement](../../docs/CI_CD_ENHANCEMENT.md)

## Questions?

See the main [README](../../README.md) for current documentation pointing to the SDK server.

---

**TL;DR:** This legacy server is deprecated. Use the SDK server at `cmd/server-sdk/` instead. It provides the same functionality with better code quality, official SDK integration, and comprehensive testing.
