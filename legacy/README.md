# Legacy MCP Server - Archived

## ⚠️ Status: Archived (October 31, 2024)

This directory contains the **archived legacy MCP server code** that has been replaced by the official SDK-based implementation.

**DO NOT USE THIS CODE FOR NEW DEPLOYMENTS.**

## 📍 Use the SDK Server Instead

**Active Server:** [`cmd/server-sdk/`](../cmd/server-sdk/)

The SDK-based server provides:
- ✅ Official Golang MCP SDK v1.1.0 (maintained by Anthropic & Google)
- ✅ 26% less code with better type safety
- ✅ Automatic schema generation
- ✅ 100% feature parity (validated via tests)
- ✅ Active development and support

## 📂 What's Archived Here

This archive contains the deprecated custom MCP server implementation:

```
legacy/
├── cmd/
│   └── server/              # Legacy MCP server entrypoint
│       ├── main.go          # Server startup code
│       └── DEPRECATED.md    # Deprecation notice
│
├── internal/
│   ├── composition/         # Dependency injection for legacy server
│   │   ├── container.go
│   │   └── container_test.go
│   │
│   ├── interfaces/          # Legacy MCP protocol handlers
│   │   ├── dto/             # Data transfer objects
│   │   └── mcp/             # Handler implementations
│   │       ├── actor_handlers.go
│   │       ├── movie_handlers.go
│   │       ├── prompt_handlers.go
│   │       ├── compound_handlers.go
│   │       ├── context_manager.go
│   │       ├── tool_validator.go
│   │       └── *_test.go files
│   │
│   ├── schemas/             # Manual tool schema definitions
│   │   ├── movie_tools.go
│   │   ├── actor_tools.go
│   │   ├── compound_tools.go
│   │   ├── context_tools.go
│   │   ├── search_tools.go
│   │   └── helpers.go
│   │
│   └── server/              # Legacy MCP server core
│       ├── mcp_server.go    # Main server logic
│       ├── protocol.go      # MCP protocol implementation
│       ├── router.go        # Tool routing
│       ├── resources.go     # Resource handlers
│       ├── registry.go      # Tool registry
│       └── *_test.go files
│
├── tests/
│   └── integration/         # Legacy server integration tests
│
└── docs/
    └── (archived documentation)
```

**Total Archived:** ~4,000+ lines of legacy code

## 🔍 Why Was This Archived?

The legacy custom MCP server was replaced with an SDK-based implementation for several reasons:

### Problems with Legacy Implementation
- **Manual Schema Definitions** - Required hand-written JSON schemas
- **Custom Protocol Layer** - ~1,200 lines of custom MCP protocol code
- **High Maintenance** - More code to maintain and test
- **No Official Support** - Custom implementation, not maintained by SDK team

### Advantages of SDK Server
- **Official SDK** - Golang MCP SDK v1.1.0 by Anthropic & Google
- **Type Safety** - Compile-time validation instead of runtime errors
- **Auto Schemas** - Generated automatically from Go types
- **Less Code** - 26% reduction (~1,200 lines eliminated)
- **Better Testing** - 37% less test code with clearer patterns
- **Active Support** - Maintained by official SDK team

## 📅 Timeline

| Date | Event |
|------|-------|
| Sep 2024 | Legacy server operational |
| Oct 2024 | SDK migration started (PR #20) |
| Oct 31, 2024 | SDK migration complete, legacy archived |
| **Now** | **Legacy code in archive/** |

## 🚫 Do Not Use

**This code is archived and should NOT be used:**

- ❌ Not actively maintained
- ❌ Not tested in CI/CD
- ❌ Missing new features
- ❌ No security updates
- ❌ Deprecated and unsupported

## ✅ Migration Guide

If you're still using the legacy server, migrate to the SDK server:

### Step 1: Build SDK Server
```bash
go build -o movies-mcp-server-sdk ./cmd/server-sdk/main.go
```

### Step 2: Test Locally
```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=movies_user
export DB_PASSWORD=movies_password
export DB_NAME=movies_mcp
export DB_SSLMODE=disable

# Run SDK server
./movies-mcp-server-sdk
```

### Step 3: Update Configuration

**Claude Desktop Config:**
```json
{
  "mcpServers": {
    "movies": {
      "command": "/path/to/movies-mcp-server-sdk"
    }
  }
}
```

### Step 4: Verify

All 23 tools + 3 resources should work identically:
- 8 movie management tools
- 9 actor management tools
- 3 compound/intelligence tools
- 3 context management tools
- 3 MCP resources

## 📖 Documentation

**For SDK Server:**
- [SDK Server README](../cmd/server-sdk/README.md)
- [Main README](../README.md)
- [SDK Migration Comparison](../docs/SDK_MIGRATION_COMPARISON.md)
- [Migration Complete Summary](../docs/SDK_MIGRATION_COMPLETE.md)

**For Legacy (Archived):**
- [DEPRECATED.md](cmd/server/DEPRECATED.md) - Deprecation details
- [LEGACY_ARCHIVAL_PLAN.md](../LEGACY_ARCHIVAL_PLAN.md) - Archival plan

## 🔗 Related PRs

The SDK migration was completed through these PRs:

- **PR #20**: SDK Migration with 23 Tools
- **PR #21**: SDK Server Enhancements (Resources + BDD Testing)
- **PR #22**: Resource Unit Tests (58 total tests)
- **PR #23**: CI/CD Enhancement (BDD tests for both servers)
- **PR #24**: Deprecate Legacy Server and Update Documentation
- **PR #25**: Project Cleanup and Add CHANGELOG
- **PR #27**: Legacy Archival Plan for Q1 2025

## ❓ FAQ

### Q: Can I still run the legacy server?

**A:** Technically yes, but you shouldn't. The code is archived and unsupported. Use the SDK server instead.

### Q: Will the legacy code be deleted?

**A:** The code is preserved in this `legacy/` directory for historical reference. It won't be deleted but it's not maintained.

### Q: What if I find a bug in the legacy server?

**A:** Use the SDK server instead. Legacy code is not maintained. If the bug exists in the SDK server too, report it and it will be fixed there.

### Q: Can I contribute to the legacy code?

**A:** No. All development happens on the SDK server at `cmd/server-sdk/`. Legacy code is frozen.

### Q: How do I compare legacy vs SDK?

**A:** See [docs/SDK_MIGRATION_COMPARISON.md](../docs/SDK_MIGRATION_COMPARISON.md) for detailed before/after comparisons.

## 🎯 Summary

- **Status:** Archived and unsupported
- **Replacement:** SDK server at `cmd/server-sdk/`
- **Reason:** Official SDK provides better code quality and maintainability
- **Action Required:** Migrate to SDK server if still using legacy

---

**For current documentation, see the [main README](../README.md).**

**Last Updated:** October 31, 2024
**Archived:** October 31, 2024
**Status:** Deprecated and archived - use SDK server instead
