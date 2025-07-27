# CI/CD Pipeline Documentation

## Overview

This repository uses GitHub Actions for automated testing, linting, security scanning, and build verification. The pipeline is designed to ensure code quality and security while providing comprehensive test coverage across multiple platforms.

## Workflow Structure

### Main Workflow: `.github/workflows/test.yml`

The primary workflow consists of several jobs that run in parallel:

1. **Access Control Check** - Validates permissions for manual triggers
2. **Unit Tests** - Fast tests with comprehensive coverage
3. **Integration Tests** - Real PE file testing
4. **Linting** - Code quality and security scanning
5. **Build Verification** - Cross-platform binary builds
6. **Test Summary** - Aggregates all job results

## Trigger Conditions

### Automatic Triggers
```yaml
on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]
```

- **Push to master**: Full pipeline execution on main branch
- **Pull requests**: Validation before merging to master

### Manual Triggers (Access Controlled)
```yaml
workflow_dispatch:
  inputs:
    test_type:
      description: 'Type of tests to run'
      type: choice
      options: [all, unit-only, integration-only]
    go_version:
      description: 'Go version to test with'
      default: '1.21'
```

**Access Control**: Only repository owners and collaborators with write/admin/maintain permissions can manually trigger workflows.

## Jobs Detail

### 1. Access Control Check
- **Purpose**: Validates user permissions for manual triggers
- **Technology**: GitHub API via `actions/github-script`
- **Permissions Required**: `write`, `admin`, or `maintain`
- **Scope**: Only applies to `workflow_dispatch` events

### 2. Unit Tests
- **Matrix**: 
  - OS: Ubuntu, Windows, macOS
  - Go: 1.21, 1.22
- **Features**:
  - Race condition detection (`-race`)
  - Code coverage generation
  - Fast execution (`-short`)
  - Codecov integration
- **Artifacts**: Coverage reports (HTML + raw data)

### 3. Integration Tests
- **Matrix**: 
  - OS: Ubuntu, Windows (reduced matrix for efficiency)
  - Go: 1.21
- **Features**:
  - Real PE file testing
  - Automatic test file setup
  - Benchmark execution
  - Performance monitoring
- **Artifacts**: Test files and results

### 4. Linting & Code Quality
- **Tools**:
  - `golangci-lint` (comprehensive Go linting) - v6
  - `go vet` (Go static analysis)
  - `gofmt` (code formatting)
  - `gosec` (security scanning) - latest via direct installation
- **Configuration**: `.golangci.yml`

### 5. Build Verification
- **Matrix**: Ubuntu, Windows, macOS
- **Output**: Cross-platform binaries
- **Testing**: Binary execution verification
- **Artifacts**: Built binaries (30-day retention)

## Access Control Implementation

The manual trigger access control works as follows:

1. **Permission Check**: Query GitHub API for user's repository permissions
2. **Validation**: Allow only `admin`, `maintain`, `write` permissions
3. **Enforcement**: Fail workflow if insufficient permissions
4. **Logging**: Clear audit trail of who triggered what

```javascript
const { data: user } = await github.rest.repos.getCollaboratorPermissionLevel({
  owner: context.repo.owner,
  repo: context.repo.repo,
  username: context.actor
});

if (!['admin', 'maintain', 'write'].includes(user.permission)) {
  throw new Error(`Access denied. User ${context.actor} does not have sufficient permissions.`);
}
```

## Configuration Files

### `.golangci.yml`
Comprehensive linting configuration with:
- 15+ enabled linters
- Security-focused rules
- Test file exemptions
- Complexity thresholds

### `.github/dependabot.yml`
Automated dependency updates for:
- Go modules (weekly)
- GitHub Actions (weekly)
- Automatic PR creation with labels

### `.github/SECURITY.md`
Security policy covering:
- Vulnerability reporting process
- Response timelines
- Security best practices

## Usage Examples

### For Repository Owners/Collaborators

1. **Manual Workflow Trigger**:
   - Go to Actions tab → "Tests" workflow
   - Click "Run workflow"
   - Select test type and Go version
   - Click "Run workflow"

2. **Test Type Options**:
   - `all`: Run both unit and integration tests
   - `unit-only`: Fast unit tests only
   - `integration-only`: Real PE file tests only

### For Contributors

1. **Pull Request Testing**:
   - Push to branch → Open PR → Automatic testing
   - All checks must pass before merge

2. **Local Development**:
   ```bash
   # Run the same checks locally
   go test -short -v -race ./...           # Unit tests
   go test -v -run Integration ./...       # Integration tests
   golangci-lint run                       # Linting
   ```

## Monitoring & Debugging

### Workflow Status
- **GitHub UI**: Actions tab shows all workflow runs
- **Status Badges**: Can be added to README
- **Notifications**: Configure in repository settings

### Common Issues
1. **Permission Denied**: User lacks repository permissions
2. **Integration Test Failures**: Test file download issues
3. **Lint Failures**: Code quality issues
4. **Build Failures**: Cross-platform compatibility issues

### Debugging Steps
1. Check workflow logs in GitHub Actions
2. Reproduce locally with same Go version
3. Verify test file availability
4. Run linting locally

## Security Considerations

1. **Access Control**: Strict permission checking for manual triggers
2. **Dependency Management**: Automated updates via Dependabot
3. **Security Scanning**: Multiple tools (gosec, golangci-lint)
4. **Artifact Management**: Controlled retention periods
5. **Secrets**: No secrets required for current workflow

## Latest GitHub Actions Versions (2024)

The workflow uses the most current versions of GitHub Actions to avoid deprecation issues and gain performance improvements:

1. **actions/checkout@v4** - Latest stable version
2. **actions/setup-go@v5** - Built-in caching enabled, faster setup
3. **actions/upload-artifact@v4** - Significant performance improvements (up to 90% faster)
4. **codecov/codecov-action@v5** - Latest features with faster updates via Codecov Wrapper
5. **actions/github-script@v7** - Updated to Node 20 runtime
6. **golangci/golangci-lint-action@v6** - Latest version with golangci-lint v2 support

### Performance Improvements

- **Caching**: setup-go v5 includes automatic module and build caching
- **Upload Speed**: artifact actions v4 provides up to 90% upload speed improvement
- **No Manual Cache**: Removed explicit cache actions since setup-go handles it automatically
- **Security**: Direct gosec installation ensures latest security scanning capabilities

## Best Practices

1. **Performance**: Matrix optimization for different test types
2. **Caching**: Go module caching for faster builds
3. **Artifacts**: Meaningful retention periods
4. **Monitoring**: Comprehensive job status tracking
5. **Documentation**: Clear workflow documentation and comments