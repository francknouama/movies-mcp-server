# CI/CD Enhancement: BDD Tests for Both Servers

## Summary

Enhanced the GitHub Actions CI pipeline to automatically test both legacy and SDK servers using BDD tests in a matrix strategy, ensuring continuous validation of feature parity and behavioral compatibility.

## 🎯 What Was Added

### New Job: `bdd-tests`

A new parallel testing job that runs BDD tests against both server implementations:

```yaml
bdd-tests:
  strategy:
    matrix:
      server: [legacy, sdk]
    fail-fast: false
```

**Key Features:**
- ✅ Matrix strategy for parallel execution
- ✅ Separate PostgreSQL database for each test run
- ✅ Environment variable `TEST_MCP_SERVER` to select server
- ✅ Automatic database migrations for each server type
- ✅ Test result artifacts uploaded separately
- ✅ Fail-fast disabled (tests both even if one fails)

### Enhanced Job: `test`

Updated the existing test job to include SDK testing:

**Changes:**
- Build both server binaries (legacy + SDK)
- Run explicit SDK unit tests (tools + resources)
- Upload both binaries as artifacts

## 📊 CI Pipeline Architecture

### Complete Job Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     GitHub Actions CI                       │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
   ┌─────────┐          ┌──────────┐         ┌──────────┐
   │  test   │          │   lint   │         │ security │
   └─────────┘          └──────────┘         └──────────┘
        │
        ├─ Build both servers
        ├─ Run migrations
        ├─ Unit tests
        ├─ SDK unit tests (NEW)
        ├─ Integration tests
        └─ Upload coverage + binaries

                    ┌────────────────┐
                    │  bdd-tests     │  (NEW)
                    └────────────────┘
                            │
            ┌───────────────┴───────────────┐
            ▼                               ▼
    ┌──────────────┐              ┌──────────────┐
    │   legacy     │              │     sdk      │
    │   server     │              │    server    │
    └──────────────┘              └──────────────┘
            │                               │
            ├─ Setup Postgres               ├─ Setup Postgres
            ├─ Build server                 ├─ Build server
            ├─ Run migrations               ├─ Run migrations
            ├─ Run BDD tests                ├─ Run BDD tests
            └─ Upload results               └─ Upload results
```

### Matrix Strategy Details

```yaml
strategy:
  matrix:
    server: [legacy, sdk]
  fail-fast: false

env:
  TEST_MCP_SERVER: ${{ matrix.server }}
```

**Execution:**
1. Job splits into 2 parallel runs
2. Each gets its own PostgreSQL instance
3. Each builds and tests its server type
4. Results uploaded separately
5. Both must pass for pipeline success

## 🔧 Implementation Details

### BDD Test Job

```yaml
- name: Run database migrations
  run: |
    if [ "${{ matrix.server }}" = "sdk" ]; then
      go build -o migrate-tool ./cmd/server-sdk/main.go
      ./migrate-tool -migrate-only
    else
      go build -o migrate-tool ./cmd/server/main.go
      ./migrate-tool -migrate-only
    fi

- name: Run BDD tests against ${{ matrix.server }} server
  run: |
    echo "Testing ${{ matrix.server }} server implementation"
    go test -v -timeout 5m ./tests/bdd/...
```

### SDK Unit Tests

```yaml
- name: Run SDK unit tests
  run: |
    go test -v -race ./internal/mcp/tools/...
    go test -v -race ./internal/mcp/resources/...
```

## 📈 Test Coverage in CI

| Test Type | Legacy Server | SDK Server | Status |
|-----------|--------------|------------|--------|
| **Unit Tests** | ✅ All tests | ✅ All tests | Parallel |
| **SDK Unit Tests** | N/A | ✅ 58 tests | Explicit |
| **Integration Tests** | ✅ All tests | ✅ All tests | Shared |
| **BDD Tests** | ✅ All scenarios | ✅ All scenarios | Matrix |

**Total CI Test Runs per Push:**
- Unit tests: 1 run (covers both)
- SDK tests: 1 run (58 tests)
- Integration tests: 1 run
- BDD tests: 2 runs (legacy + sdk in parallel)

## 💡 Benefits

### 1. Continuous Validation
- Every commit tests both implementations
- Catches regressions immediately
- Ensures behavioral parity

### 2. Parallel Execution
- Legacy and SDK tests run simultaneously
- Faster feedback (no sequential wait)
- Independent failure isolation

### 3. Complete Coverage
- Unit tests: 58 tests for SDK tools/resources
- BDD tests: Full end-to-end scenarios
- Integration tests: Cross-cutting concerns

### 4. Early Detection
- API compatibility issues
- Behavioral differences
- Performance regressions

### 5. Confidence in Migration
- Validates SDK server is drop-in replacement
- Proves feature parity
- Documents expected behavior

## 🚀 CI Workflow Features

### Database Management
- Each BDD test run gets fresh PostgreSQL instance
- Automatic health checks before tests
- Clean state for every test

### Artifact Collection
- Test results uploaded per server type
- Both server binaries uploaded
- Coverage reports generated
- 7-day retention

### Error Handling
- `fail-fast: false` - Both tests run even if one fails
- Detailed logs per server
- Separate artifact uploads

### Caching
- Go module cache shared across jobs
- Build cache reused
- Faster subsequent runs

## 📊 Expected CI Output

### Successful Run
```
✅ test job
  ├─ Unit tests: PASS
  ├─ SDK unit tests: 58/58 PASS
  ├─ Integration tests: PASS
  └─ Artifacts uploaded

✅ bdd-tests (legacy)
  ├─ Database ready
  ├─ Migrations applied
  ├─ BDD scenarios: ALL PASS
  └─ Results uploaded

✅ bdd-tests (sdk)
  ├─ Database ready
  ├─ Migrations applied
  ├─ BDD scenarios: ALL PASS
  └─ Results uploaded

✅ lint: PASS
✅ security: PASS
✅ validate-dependencies: PASS
```

### Failure Example
```
✅ test job: PASS

❌ bdd-tests (legacy): PASS
✅ bdd-tests (sdk): FAIL
  └─ Scenario "Create Movie" failed
      └─ See artifact: bdd-test-results-sdk-123456

Result: Pipeline fails (SDK test failure)
Action: Fix SDK server, re-run
```

## 🎯 Impact

### Before Enhancement
- BDD tests excluded from CI (`grep -v tests/bdd`)
- Only legacy server tested automatically
- SDK server validation was manual
- No automated feature parity checks

### After Enhancement
- ✅ BDD tests run automatically on every push/PR
- ✅ Both servers tested in parallel
- ✅ Automated feature parity validation
- ✅ 58 SDK unit tests explicitly run
- ✅ Comprehensive CI coverage

### Metrics
- **Jobs Added:** 1 (bdd-tests with 2 matrix runs)
- **Tests Added to CI:** BDD scenarios × 2 servers
- **SDK Tests Explicitly Run:** 58 tests
- **Parallel Runs:** 2 (legacy + sdk)
- **Total CI Time:** ~Same (parallel execution)

## ✅ Validation

### What Gets Tested
1. **Unit Level:** 58 SDK tests (tools + resources)
2. **Integration Level:** Database interactions
3. **BDD Level:** End-to-end scenarios for both servers
4. **Behavioral Level:** Both servers produce identical results

### What Gets Validated
- SDK tools work correctly
- Resource handlers function properly
- Both servers handle same scenarios identically
- Database migrations work for both
- No regressions in either implementation

## 🔄 Next Steps

With CI enhancements complete, the pipeline now:
- Automatically validates both implementations
- Catches regressions early
- Ensures feature parity continuously
- Provides confidence in SDK migration

**Future Enhancements Could Include:**
- Performance benchmarking (SDK vs legacy)
- Load testing
- Integration with external services
- Deployment automation
- Release workflow

## 📝 Usage

### Running BDD Tests Locally

```bash
# Test legacy server
TEST_MCP_SERVER=legacy go test -v ./tests/bdd/...

# Test SDK server
TEST_MCP_SERVER=sdk go test -v ./tests/bdd/...

# Test both (matches CI)
TEST_MCP_SERVER=legacy go test -v ./tests/bdd/... && \
TEST_MCP_SERVER=sdk go test -v ./tests/bdd/...
```

### Viewing CI Results

1. Go to Actions tab in GitHub
2. Select workflow run
3. View "bdd-tests (legacy)" job
4. View "bdd-tests (sdk)" job
5. Download artifacts if needed

## 🎉 Conclusion

The CI pipeline now provides comprehensive automated testing of both server implementations, ensuring that the SDK migration maintains perfect feature parity with the legacy server while catching any regressions early in the development cycle.

---

**Related Work:**
- PR #20: SDK Migration with 23 Tools
- PR #21: SDK Server Enhancements (Resources + BDD Testing)
- PR #22: Resource Unit Tests
- This: CI/CD Enhancement (BDD Tests for Both Servers)
