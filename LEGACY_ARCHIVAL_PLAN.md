# Legacy Code Archival Plan

## Status: Planned for Q1 2025

This document outlines the plan for archiving the deprecated legacy MCP server code once the SDK server is proven in production.

## 📅 Timeline

| Phase | Target Date | Status |
|-------|-------------|--------|
| **SDK Server Proven** | Dec 2024 - Jan 2025 | 🔄 In Progress |
| **Legacy Code Archival** | Q1 2025 (Jan-Mar) | 📋 Planned |
| **Legacy Code Removal** | Q2 2025 (Apr-Jun) | 📋 Future |

## 🎯 Objectives

1. **Simplify Codebase** - Remove deprecated legacy code
2. **Reduce Maintenance** - Focus entirely on SDK server
3. **Preserve History** - Archive legacy code for reference
4. **Maintain Stability** - Ensure SDK server is battle-tested first

## 📊 Current State (As of Oct 2024)

### Active Code (SDK-Based)
- ✅ `cmd/server-sdk/` - Official SDK server
- ✅ `internal/mcp/` - SDK tools and resources (58 tests)
- ✅ `internal/domain/` - Business logic (shared)
- ✅ `internal/application/` - Use cases (shared)
- ✅ `internal/infrastructure/` - Database layer (shared)

### Deprecated Code (Legacy)
- ⚠️ `cmd/server/` - Legacy custom server
- ⚠️ `internal/interfaces/` - Legacy MCP handlers
- ⚠️ `internal/schemas/` - Manual tool schemas
- ⚠️ `internal/server/` - Legacy server core
- ⚠️ `internal/composition/container.go` - Uses legacy handlers

### Dependencies on Legacy Code

**Files Using Legacy Code:**
1. `cmd/server/main.go` - Legacy server entrypoint
2. `tests/integration/integration_test.go` - Integration tests
3. `tests/bdd/steps/actor_steps.go` - BDD test steps (DTOs)
4. `internal/composition/container.go` - DI container

**Why Legacy Code Is Still Needed:**
- ✅ CI/CD tests both servers in parallel (feature parity validation)
- ✅ BDD tests validate legacy server behavior
- ✅ Backwards compatibility for existing deployments
- ✅ Comparison reference for SDK implementation

## 🗂️ Archival Strategy

### Phase 1: SDK Server Proven (Dec 2024 - Jan 2025)

**Objectives:**
- Prove SDK server in production
- Gather real-world feedback
- Monitor for regressions
- Build confidence

**Success Criteria:**
- [ ] SDK server running in production for 30+ days
- [ ] Zero critical bugs specific to SDK server
- [ ] All users migrated to SDK server
- [ ] Performance meets or exceeds legacy server
- [ ] No feature requests requiring legacy server

**Actions:**
- Deploy SDK server to production
- Monitor metrics (latency, errors, throughput)
- Collect user feedback
- Document any issues
- Verify all 23 tools + 3 resources work correctly

### Phase 2: Prepare for Archival (Q1 2025)

**Objectives:**
- Plan archival execution
- Update CI/CD to SDK-only
- Prepare documentation
- Communicate to users

**Tasks:**
1. **Update CI/CD Pipeline**
   - Remove legacy server from BDD matrix
   - Update workflow to only test SDK server
   - Remove legacy server build step

2. **Update Tests**
   - Remove `TEST_MCP_SERVER=legacy` tests
   - Keep only SDK server tests
   - Update test documentation

3. **Move DTOs (if needed)**
   - Move shared DTOs to `internal/dto/` if BDD tests need them
   - Or update BDD tests to not depend on legacy DTOs

4. **Create Archive Structure**
   ```
   legacy/
   ├── README.md (explains archival)
   ├── cmd/
   │   └── server/ (legacy server)
   ├── internal/
   │   ├── interfaces/ (legacy handlers)
   │   ├── schemas/ (manual schemas)
   │   └── server/ (legacy server core)
   └── docs/
       └── LEGACY_MIGRATION.md
   ```

5. **Update Documentation**
   - Remove legacy references from README
   - Update project structure diagram
   - Add archive location to docs
   - Update CHANGELOG

**Success Criteria:**
- [ ] CI/CD runs without legacy server
- [ ] All tests pass with SDK server only
- [ ] Documentation updated
- [ ] Users notified of archival

### Phase 3: Execute Archival (Q1 2025)

**Step-by-Step Process:**

1. **Create Archive Directory**
   ```bash
   mkdir -p legacy/cmd
   mkdir -p legacy/internal
   mkdir -p legacy/docs
   ```

2. **Move Legacy Code**
   ```bash
   # Move legacy server
   git mv cmd/server/ legacy/cmd/server/

   # Move legacy internal packages
   git mv internal/interfaces/ legacy/internal/interfaces/
   git mv internal/schemas/ legacy/internal/schemas/
   git mv internal/server/ legacy/internal/server/
   ```

3. **Update Composition**
   - Remove legacy handlers from `internal/composition/container.go`
   - Keep only SDK-related dependencies
   - Or move container to legacy if only used by legacy server

4. **Update CI/CD**
   ```yaml
   # Remove from .github/workflows/ci.yml
   bdd-tests:
     strategy:
       matrix:
         server: [sdk]  # Remove 'legacy'
   ```

5. **Update Tests**
   - Remove or update `tests/integration/integration_test.go`
   - Update `tests/bdd/` if needed
   - Remove legacy-specific test code

6. **Create Archive Documentation**
   ```markdown
   legacy/README.md:
   # Legacy MCP Server - Archived

   This code has been archived as of [DATE].
   Use the SDK server at `cmd/server-sdk/` instead.

   See main README for migration guide.
   ```

7. **Update Main Documentation**
   - Remove legacy server from README
   - Update architecture diagrams
   - Remove deprecation warnings (no longer needed)
   - Update CHANGELOG

8. **Commit and PR**
   ```bash
   git commit -m "chore: archive legacy server code"
   git push
   # Create PR with detailed explanation
   ```

**Verification:**
- [ ] Project builds successfully
- [ ] All tests pass
- [ ] CI/CD pipeline succeeds
- [ ] Documentation is accurate
- [ ] No broken imports

### Phase 4: Complete Removal (Q2 2025 - Optional)

**Only if archive is no longer needed:**

- Move legacy code to separate archive branch
- Remove from main branch entirely
- Keep archive branch for historical reference
- Update documentation to point to archive branch

## 📋 Detailed Task Checklist

### Pre-Archival (Now - Dec 2024)

- [x] SDK migration complete (PRs #20-#25)
- [x] Deprecation notices added
- [x] CI/CD tests both servers
- [ ] SDK server deployed to production
- [ ] Production monitoring established
- [ ] User migration complete

### Preparation (Jan 2025)

- [ ] Review production metrics
- [ ] Confirm zero legacy server usage
- [ ] Update CI/CD workflows
- [ ] Test SDK-only CI pipeline
- [ ] Prepare archival PR
- [ ] Create GitHub issue for tracking

### Execution (Feb-Mar 2025)

- [ ] Create archive directory structure
- [ ] Move legacy code to archive/
- [ ] Update all imports
- [ ] Remove legacy from CI/CD
- [ ] Update documentation
- [ ] Run full test suite
- [ ] Create and merge archival PR
- [ ] Verify main branch works

### Post-Archival (After Archival)

- [ ] Monitor for issues
- [ ] Update CHANGELOG
- [ ] Announce to community
- [ ] Archive GitHub issues related to legacy
- [ ] Clean up project board

## 🔍 What Will Be Archived

### Code to Archive (~4,000+ lines)

```
legacy/
├── cmd/
│   └── server/
│       ├── main.go (~200 lines)
│       └── DEPRECATED.md
├── internal/
│   ├── interfaces/
│   │   └── mcp/
│   │       ├── actor_handlers.go (~500 lines)
│   │       ├── movie_handlers.go (~800 lines)
│   │       ├── prompt_handlers.go (~300 lines)
│   │       ├── compound_handlers.go (~400 lines)
│   │       ├── context_manager.go (~300 lines)
│   │       ├── tool_validator.go (~200 lines)
│   │       ├── handler_utils.go (~100 lines)
│   │       └── *_test.go files
│   ├── schemas/
│   │   ├── movie_tools.go (~400 lines)
│   │   ├── actor_tools.go (~400 lines)
│   │   ├── compound_tools.go (~300 lines)
│   │   ├── context_tools.go (~200 lines)
│   │   ├── search_tools.go (~300 lines)
│   │   ├── validation_tools.go (~100 lines)
│   │   ├── helpers.go (~100 lines)
│   │   └── tools.go (~100 lines)
│   └── server/
│       ├── mcp_server.go (~600 lines)
│       ├── protocol.go (~300 lines)
│       ├── router.go (~200 lines)
│       ├── resources.go (~200 lines)
│       ├── registry.go (~100 lines)
│       └── *_test.go files
└── docs/
    └── LEGACY_MIGRATION.md
```

### Code to Keep (Active)

```
cmd/
└── server-sdk/ (SDK server)

internal/
├── mcp/ (SDK tools + resources)
├── domain/ (business logic)
├── application/ (use cases)
├── infrastructure/ (database)
├── config/ (configuration)
└── composition/ (DI - may need update)
```

## 📊 Impact Analysis

### Benefits of Archival

1. **Simplified Codebase**
   - ~4,000 fewer lines to maintain
   - Focus on single implementation
   - Easier for new contributors

2. **Reduced Complexity**
   - No dual implementations
   - Clearer project structure
   - Simpler CI/CD

3. **Better Performance**
   - Faster CI/CD (no legacy tests)
   - Smaller binary builds
   - Reduced dependency tree

4. **Lower Maintenance**
   - No legacy bug fixes
   - No dual documentation
   - Focus on SDK improvements

### Risks and Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Users still using legacy | Low | Medium | Communicate early, provide migration guide |
| Undiscovered SDK bugs | Low | High | 3+ months production validation first |
| Lost reference code | Low | Low | Archive preserved, git history intact |
| CI/CD breaks | Medium | Medium | Test SDK-only pipeline before archival |
| Breaking changes | Low | Low | Only internal code affected, APIs unchanged |

## 🔗 Related Documentation

- [DEPRECATED.md](cmd/server/DEPRECATED.md) - Legacy server deprecation notice
- [CHANGELOG.md](CHANGELOG.md) - Project history including SDK migration
- [SDK_MIGRATION_COMPARISON.md](docs/SDK_MIGRATION_COMPARISON.md) - SDK vs Legacy comparison
- [CI_CD_ENHANCEMENT.md](docs/CI_CD_ENHANCEMENT.md) - Current dual-server testing

## 📞 Communication Plan

### Users to Notify

1. **Existing Users** - Email/announcement about archival timeline
2. **Contributors** - GitHub issue + discussion
3. **Community** - Blog post about SDK migration success
4. **Documentation** - Update all public docs

### Key Messages

1. **Legacy server is deprecated** - Use SDK server
2. **Archival planned for Q1 2025** - After production validation
3. **Migration guide available** - See DEPRECATED.md
4. **SDK server is production-ready** - 100% feature parity validated

### Communication Timeline

- **Now (Oct 2024)**: Deprecation notices in place
- **Dec 2024**: "SDK server in production" announcement
- **Jan 2025**: "Planning legacy archival" notice
- **Feb 2025**: "Archival PR created" notification
- **Mar 2025**: "Legacy code archived" announcement

## ✅ Success Metrics

### SDK Server Readiness

- [ ] 30+ days in production
- [ ] <0.1% error rate
- [ ] Response time < 100ms (p95)
- [ ] All 23 tools working correctly
- [ ] All 3 resources working correctly
- [ ] Zero critical bugs
- [ ] User feedback positive

### Archival Execution

- [ ] All code moved to legacy/
- [ ] CI/CD passes with SDK only
- [ ] All tests pass
- [ ] Documentation updated
- [ ] No broken imports
- [ ] Git history preserved

## 🚀 Next Steps

1. **Immediate (Now)**
   - ✅ Document archival plan (this file)
   - Create GitHub issue for tracking
   - Update CHANGELOG with timeline

2. **Short-term (Nov-Dec 2024)**
   - Deploy SDK server to production
   - Monitor performance and stability
   - Gather user feedback

3. **Mid-term (Q1 2025)**
   - Execute archival when ready
   - Update CI/CD to SDK-only
   - Archive legacy code

4. **Long-term (Q2 2025+)**
   - Optional: Complete removal
   - Move to archive branch
   - Keep for historical reference

## 📝 Notes

- This plan is flexible - dates may adjust based on production validation
- Legacy code will only be archived when SDK server is proven stable
- Git history will preserve all legacy code
- Archive will remain accessible for reference
- No rush - stability is more important than timeline

---

**Last Updated:** October 31, 2024
**Status:** Plan approved, awaiting production validation
**Target:** Q1 2025 (flexible based on SDK server stability)
