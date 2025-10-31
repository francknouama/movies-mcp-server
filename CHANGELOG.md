# Changelog

All notable changes to the Movies MCP Server project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CHANGELOG.md to track project changes

## [2.0.0] - 2024-10-31

### 🎉 SDK Migration Complete

This major release completes the migration from a custom MCP protocol implementation to the official Golang MCP SDK v1.1.0.

### Added

#### SDK-Based Server (PR #20, #21)
- **New SDK server** at `cmd/server-sdk/` using official Golang MCP SDK v1.1.0
- **23 MCP tools** migrated to SDK-based handlers:
  - 8 movie management tools
  - 9 actor management tools
  - 3 compound/intelligence tools
  - 3 context management tools
- **3 MCP resources** implemented with SDK resource handlers:
  - `movies://database/all` - Complete movie database
  - `movies://database/stats` - Database statistics
  - `movies://posters/collection` - All movie posters

#### Comprehensive Testing (PR #22, #23)
- **58 unit tests** for SDK implementation (46 tool + 12 resource tests)
- **BDD test support** for both legacy and SDK servers via `TEST_MCP_SERVER` env var
- **CI/CD matrix testing** - Automated parallel testing of both servers in GitHub Actions
- **Feature parity validation** on every push/PR

#### Documentation (PR #24)
- **Deprecation notices** for legacy server at `cmd/server/DEPRECATED.md`
- **Migration guides** for developers and deployments
- **Updated README** with prominent SDK server recommendations
- **SDK migration documentation** in `/docs`:
  - SDK Migration Comparison
  - Testing Comparison
  - Migration Complete Summary
  - CI/CD Enhancement docs

### Changed

#### Code Quality
- **26% code reduction** - Eliminated ~1,200 lines of custom protocol layer
- **37% test code reduction** - Simplified testing with SDK patterns
- **Type-safe handlers** with compile-time validation
- **Automatic schema generation** replacing manual JSON definitions

#### Testing Infrastructure
- BDD tests now test both legacy and SDK servers in parallel
- GitHub Actions CI validates both implementations on every push
- Automated feature parity checks ensure 100% compatibility
- Enhanced test coverage with comprehensive SDK unit tests

#### Documentation
- All examples now use SDK server by default
- README prominently recommends SDK server
- Legacy server marked as deprecated throughout
- Clear migration path documented

### Deprecated

#### Legacy Server (`cmd/server/`)
- **Custom MCP server** at `cmd/server/` officially deprecated
- **Legacy protocol handlers** in `internal/interfaces/mcp/`
- **Manual schemas** in `internal/schemas/`
- **Legacy server core** in `internal/server/`

**Deprecation Timeline:**
- **Now**: Deprecated - Use SDK server for all new deployments
- **Q1 2025**: May be archived after SDK server proven in production
- **Future**: Will be removed when no longer needed for testing

**Why deprecated:**
- Official SDK provides better type safety and maintainability
- 26% less code with automatic schema generation
- Maintained by Anthropic and Google
- 100% feature parity validated via CI/CD

**Why kept (for now):**
- BDD testing validates both implementations
- Backwards compatibility for existing users
- Comparison reference for developers

### Migration Path

**For New Users:**
- Use SDK server at `cmd/server-sdk/`
- Follow Quick Start in README
- Configure Claude Desktop with SDK server

**For Existing Users:**
1. Test SDK server locally: `TEST_MCP_SERVER=sdk go test -v ./tests/bdd/...`
2. Build SDK server: `go build -o movies-mcp-server-sdk ./cmd/server-sdk/main.go`
3. Update configuration to use `movies-mcp-server-sdk`
4. Deploy and monitor
5. Verify all 23 tools + 3 resources work identically

See `cmd/server/DEPRECATED.md` for detailed migration guide.

### Technical Details

#### SDK Migration Statistics
- **Tools Migrated:** 23/23 (100%)
- **Resources Migrated:** 3/3 (100%)
- **Unit Tests:** 58 (46 tools + 12 resources)
- **Code Reduction:** 26% (~1,200 lines eliminated)
- **Test Code Reduction:** 37%
- **Feature Parity:** 100% validated via BDD + CI/CD

#### Architecture Preserved
- Clean Architecture principles maintained
- Domain layer unchanged (zero business logic changes)
- Application layer unchanged (use cases preserved)
- Infrastructure layer unchanged (PostgreSQL integration)
- Only MCP protocol layer migrated to SDK

#### Dependencies
- Added: `github.com/modelcontextprotocol/go-sdk` v1.1.0
- Maintained: All existing application dependencies
- No breaking changes to business logic or database

### Pull Requests

- **PR #20**: SDK Migration with 23 Tools - Merged
- **PR #21**: SDK Server Enhancements (Resources + BDD Testing) - Merged
- **PR #22**: Resource Unit Tests - Merged
- **PR #23**: CI/CD Enhancement (BDD tests for both servers) - Merged
- **PR #24**: Deprecate Legacy Server and Update Documentation - Merged

### Contributors

Special thanks to all contributors who helped with the SDK migration!

## [1.0.0] - 2024-09-15

### Initial Release

- Custom MCP protocol implementation
- 23 MCP tools for movie and actor management
- PostgreSQL database integration
- Clean Architecture implementation
- Docker support with monitoring (Prometheus + Grafana)
- BDD testing with Cucumber/Godog
- Comprehensive documentation

---

## Migration Journey Timeline

### Phase 1: Foundation (PR #20)
**Date:** October 2024
- Migrated all 23 tools to SDK
- Created SDK-based server
- Established migration patterns

### Phase 2: Resources & Testing (PR #21)
**Date:** October 2024
- Added 3 SDK resource handlers
- Implemented BDD test support for both servers
- Validated feature parity

### Phase 3: Complete Testing (PR #22)
**Date:** October 2024
- Created 12 resource unit tests
- Achieved 58 total unit tests
- 100% SDK component coverage

### Phase 4: CI/CD Automation (PR #23)
**Date:** October 2024
- Added matrix testing to GitHub Actions
- Automated validation of both servers
- Continuous feature parity checks

### Phase 5: Deprecation (PR #24)
**Date:** October 31, 2024
- Official deprecation of legacy server
- Comprehensive migration documentation
- SDK server now recommended for all deployments

### Migration Complete! 🎉
**Status:** SDK migration 100% complete
**Recommendation:** Use SDK server for all new deployments
**Legacy Support:** Maintained for testing and backwards compatibility

---

For more details, see:
- [SDK Migration Comparison](docs/SDK_MIGRATION_COMPARISON.md)
- [Testing Comparison](docs/TESTING_COMPARISON.md)
- [Migration Complete Summary](docs/SDK_MIGRATION_COMPLETE.md)
- [CI/CD Enhancement](docs/CI_CD_ENHANCEMENT.md)
- [Legacy Server Deprecation](cmd/server/DEPRECATED.md)

[Unreleased]: https://github.com/francknouama/movies-mcp-server/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/francknouama/movies-mcp-server/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/francknouama/movies-mcp-server/releases/tag/v1.0.0
