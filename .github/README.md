# CI/CD Workflows

This directory contains GitHub Actions workflows for continuous integration and deployment of the Movies MCP Server project.

## 🚀 Workflows Overview

### 1. CI Pipeline (`ci.yml`)
**Triggers:** Push to `main`/`develop`, Pull Requests

- **Path-based job execution** - Only runs jobs for changed modules
- **Multi-module testing** - Tests `shared-mcp`, `mcp-server`, and `godog-server` independently
- **Database integration** - Uses PostgreSQL for realistic testing
- **ATDD testing** - Runs both mock and real server scenarios
- **Code quality** - Linting, security scans, and dependency validation
- **Go 1.24.4** - Latest Go version with performance optimizations

**Jobs:**
- `changes` - Detects which modules changed
- `test-shared-mcp` - Tests shared utilities and components
- `test-mcp-server` - Full MCP server testing with database
- `test-godog-server` - ATDD scenarios in both mock and real modes
- `lint` - Code quality and style checks
- `security` - Security vulnerability scanning
- `validate-dependencies` - Dependency consistency and updates

### 2. CD Pipeline (`cd.yml`)
**Triggers:** Push to `main`, Tags, Manual dispatch

- **Multi-stage deployment** - Staging → Production
- **Container building** - Multi-platform Docker images
- **Security scanning** - SBOM generation and vulnerability assessment
- **Smoke testing** - Automated deployment verification
- **Blue-green deployment** - Zero-downtime production deployments

**Environments:**
- **Staging** - Auto-deploys from `main` branch
- **Production** - Deploys from tags or manual trigger

### 3. Performance Testing (`performance.yml`)
**Triggers:** Push to `main`, Pull Requests, Daily schedule, Manual

- **Load testing** with k6
- **Realistic scenarios** - MCP protocol operations, database queries
- **Configurable parameters** - Duration, concurrent users
- **Performance thresholds** - 95% requests < 2s, error rate < 5%
- **Database seeding** - 10K movies, 5K actors, 50K relationships

**Scenarios tested:**
- MCP protocol initialization and tool discovery
- Movie CRUD operations and search
- Actor management and movie linking
- Complex queries and cast information

### 4. Dependency Management (`dependency-update.yml`)
**Triggers:** Weekly schedule, Manual dispatch

- **Automated Go module updates**
- **Security vulnerability scanning**
- **Dependency consistency validation**
- **Automated PR creation** for updates
- **Integration with Renovate** for GitHub Actions updates

### 5. Release Workflow (`release.yml`)
**Triggers:** Manual dispatch with version input

- **Semantic version validation**
- **Multi-platform binary building** (Linux, macOS, Windows)
- **Docker image publishing** to GitHub Container Registry
- **Automated changelog generation**
- **GitHub release creation** with assets
- **Digital signatures** and checksums

**Supported platforms:**
- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (amd64)

## 🔧 Configuration

### Environment Variables
```yaml
GO_VERSION: '1.24.4'  # Latest Go version
REGISTRY: ghcr.io      # Container registry
```

### Required Secrets
- `GITHUB_TOKEN` - Automatically provided, used for container registry and releases
- Additional secrets may be needed for production deployments

### Repository Settings
**Environments configured:**
- `staging` - Auto-deployment from main
- `production` - Protected, requires approval

**Branch Protection:**
- `main` - Requires CI checks to pass
- `develop` - Integration branch for feature development

## 📊 Monitoring and Metrics

### Coverage Reports
- Uploaded as artifacts for each module
- HTML reports available in workflow runs

### Performance Metrics
- k6 performance results with detailed metrics
- PR comments with performance impact analysis
- Daily performance trend tracking

### Security Scanning
- Gosec static analysis
- Nancy vulnerability scanning
- Container image security assessment
- SBOM (Software Bill of Materials) generation

## 🛠 Development Workflow

### Feature Development
1. Create feature branch from `develop`
2. Push commits trigger CI checks
3. Create PR to `develop`
4. CI validates all modules and runs ATDD tests
5. Merge to `develop` after approval

### Release Process
1. Merge `develop` to `main`
2. Tag with semantic version (`v1.2.3`)
3. Automated release workflow builds and publishes
4. Production deployment (manual approval required)

### Hotfix Process
1. Create hotfix branch from `main`
2. Apply critical fixes
3. Create PR to `main`
4. Emergency release after CI validation

## 🔍 Troubleshooting

### Common Issues

**1. Go Module Dependencies**
```bash
# Fix dependency issues locally
go work sync
go mod tidy
```

**2. Database Connection Issues**
- Check PostgreSQL service health in workflow logs
- Verify environment variables match service configuration

**3. Performance Test Failures**
- Check if server started properly
- Verify database seeding completed
- Review k6 test configuration

**4. Docker Build Issues**
- Ensure Dockerfile.production exists in mcp-server/
- Check multi-platform build support

### Workflow Debugging

**Enable debug logging:**
```yaml
env:
  RUNNER_DEBUG: 1
  ACTIONS_STEP_DEBUG: 1
```

**Check specific job logs:**
- Navigate to Actions tab
- Select failing workflow run
- Expand job and step details

## 📈 Performance Targets

| Metric | Target | Measurement |
|--------|---------|-------------|
| Response Time | 95% < 2000ms | k6 load testing |
| Error Rate | < 5% | HTTP status codes |
| Throughput | > 100 RPS | Concurrent operations |
| Database | < 100ms queries | PostgreSQL metrics |

## 🚀 Deployment Targets

| Environment | Trigger | Approval | Monitoring |
|-------------|---------|----------|------------|
| Staging | Auto (main) | None | Basic health checks |
| Production | Manual/Tags | Required | Full monitoring |

---

This CI/CD setup ensures high code quality, comprehensive testing, and reliable deployments for the Movies MCP Server project. All workflows are optimized for the Go 1.24.4 ecosystem and multi-module architecture.