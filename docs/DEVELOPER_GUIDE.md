# Developer Guide

This comprehensive guide will help you set up the development environment and contribute to the GoReleaser CLI template project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [System Requirements](#system-requirements)
- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [IDE Configuration](#ide-configuration)
- [Development Workflow](#development-workflow)
- [Testing Guidelines](#testing-guidelines)
- [Debugging](#debugging)
- [Performance](#performance)
- [Common Issues](#common-issues)
- [Contributing](#contributing)

## Prerequisites

### Required Tools

| Tool | Minimum Version | Purpose | Installation |
|------|----------------|---------|--------------|
| **Go** | 1.21+ | Core language | [golang.org](https://golang.org/dl/) |
| **Git** | 2.30+ | Version control | [git-scm.com](https://git-scm.com/) |
| **Make** | 4.0+ | Build automation | Usually pre-installed |
| **Docker** | 20.0+ | Containerization | [docker.com](https://docker.com/) |
| **Just** | 1.0+ | Command runner | `cargo install just` |

### Recommended Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| **pre-commit** | Git hooks | `pip install pre-commit` |
| **golangci-lint** | Go linting | [golangci-lint.run](https://golangci-lint.run/usage/install/) |
| **gosec** | Security scanner | `go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest` |
| **air** | Live reload | `go install github.com/cosmtrek/air@latest` |
| **delve** | Go debugger | `go install github.com/go-delve/delve/cmd/dlv@latest` |

## System Requirements

### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Disk**: 2GB free space
- **OS**: Linux, macOS, or Windows with WSL2

### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Disk**: 5GB+ free space
- **OS**: Linux or macOS for optimal performance

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/LarsArtmann/template-GoReleaser.git
cd template-GoReleaser
```

### 2. Initialize Development Environment

```bash
# Install Go dependencies
go mod download

# Install pre-commit hooks
pre-commit install

# Set up git hooks (if not using pre-commit)
just setup-hooks

# Verify installation
just verify
```

### 3. Build the Project

```bash
# Build for development
just build

# Build for production
just build-release

# Build with version info
just build-version
```

### 4. Run Tests

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# Run integration tests
just test-integration

# Run benchmarks
just bench
```

## Development Environment Setup

### Environment Variables

Create a `.env` file in the project root:

```bash
# Development configuration
GORELEASER_CLI_ENV=development
GORELEASER_CLI_LOG_LEVEL=debug
GORELEASER_CLI_CONFIG=.goreleaser-cli.dev.yaml

# Optional: Custom paths
GORELEASER_CLI_DATA_DIR=./data
GORELEASER_CLI_CACHE_DIR=./cache

# Testing
TEST_TIMEOUT=30s
INTEGRATION_TEST_TIMEOUT=300s
```

### Go Environment

```bash
# Set Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org

# Enable Go modules
go env -w GO111MODULE=on
```

### Development Configuration

Create `.goreleaser-cli.dev.yaml`:

```yaml
project_name: "goreleaser-cli-dev"
go_version: "1.21"
author: "Your Name"
license: "MIT"
git_repository: "github.com/yourusername/template-GoReleaser"
description: "Development instance of GoReleaser CLI"
enable_docker: false
enable_git_hooks: false
development_mode: true

logging:
  level: "debug"
  format: "text"
  
server:
  host: "localhost"
  port: 8080
  debug: true
```

## IDE Configuration

### Visual Studio Code

#### Required Extensions

1. **Go** (`golang.Go`) - Official Go extension
2. **Test Explorer UI** (`hbenl.vscode-test-explorer`) - Test runner
3. **Go Test Explorer** (`premparihar.gotestexplorer`) - Go test integration

#### Recommended Extensions

1. **GitLens** (`eamodio.gitlens`) - Git integration
2. **REST Client** (`humao.rest-client`) - API testing
3. **Thunder Client** (`rangav.vscode-thunder-client`) - API client
4. **Docker** (`ms-azuretools.vscode-docker`) - Docker integration
5. **YAML** (`redhat.vscode-yaml`) - YAML support

#### Settings Configuration

Create `.vscode/settings.json`:

```json
{
    "go.toolsManagement.autoUpdate": true,
    "go.useLanguageServer": true,
    "go.lintOnSave": "package",
    "go.vetOnSave": "package",
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v", "-race"],
    "go.testTimeout": "30s",
    "go.coverOnSave": true,
    "go.coverageDecorator": "gutter",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
    "files.associations": {
        "*.templ": "html"
    }
}
```

#### Launch Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch CLI",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "cmd/goreleaser-cli/main.go",
            "args": ["--help"]
        },
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "cmd/goreleaser-cli/main.go",
            "args": ["server", "--debug"]
        },
        {
            "name": "Debug Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/internal/validation"
        }
    ]
}
```

### GoLand/IntelliJ IDEA

#### Required Plugins
1. **Go** - Built-in Go support
2. **Docker** - Container integration

#### Configuration

1. **Go Modules**: Enable Go modules support
2. **Code Style**: Import Go code style settings
3. **Run/Debug Configurations**:
   - Create run configuration for `cmd/goreleaser-cli/main.go`
   - Set working directory to project root
   - Add program arguments as needed

## Development Workflow

### 1. Branch Management

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Keep branch up to date
git fetch origin
git rebase origin/main

# Push changes
git push -u origin feature/your-feature-name
```

### 2. Code Development

```bash
# Start development server with live reload
air

# Or use just command
just dev

# Run specific command during development
go run cmd/goreleaser-cli/main.go validate --help
```

### 3. Testing Cycle

```bash
# Run tests during development
just test-watch

# Run specific test
go test -v ./internal/validation -run TestValidateEnvironment

# Run with coverage
just test-coverage

# View coverage report
just coverage-report
```

### 4. Code Quality

```bash
# Format code
just fmt

# Lint code
just lint

# Security scan
just security-scan

# Pre-commit checks
pre-commit run --all-files
```

### 5. Build and Verify

```bash
# Build for current platform
just build

# Build for all platforms
just build-all

# Verify build
just verify

# Run integration tests
just test-integration
```

## Testing Guidelines

### Unit Tests

- Place tests in the same package as the code being tested
- Use `_test.go` suffix for test files
- Follow table-driven test patterns
- Mock external dependencies using interfaces

Example:
```go
func TestValidationService_ValidateEnvironment(t *testing.T) {
    tests := []struct {
        name    string
        setup   func()
        want    *ValidationResult
        wantErr bool
    }{
        {
            name: "valid environment",
            setup: func() {
                os.Setenv("REQUIRED_VAR", "value")
            },
            want: &ValidationResult{Success: true},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
            }
            
            service := NewValidationService()
            got, err := service.ValidateEnvironment()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Integration Tests

- Place in `tests/integration/` directory
- Test complete workflows
- Use real implementations where possible
- Clean up resources after tests

### Benchmark Tests

- Place in `benchmarks/` directory
- Use `Benchmark` prefix for functions
- Report allocations with `b.ReportAllocs()`
- Run with: `go test -bench=. -benchmem ./benchmarks/`

### Test Commands

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# Run specific test package
go test -v ./internal/validation

# Run specific test
go test -v -run TestSpecificFunction ./internal/validation

# Run tests with race detector
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./benchmarks/

# Generate test coverage report
just coverage-html
```

## Debugging

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main application
dlv debug cmd/goreleaser-cli/main.go -- validate

# Debug tests
dlv test ./internal/validation -- -test.run TestValidateEnvironment

# Attach to running process
dlv attach <pid>
```

### Debug Commands

```bash
# Enable debug logging
export GORELEASER_CLI_LOG_LEVEL=debug

# Run with verbose output
go run cmd/goreleaser-cli/main.go validate --verbose

# Use debug build
just build-debug
./goreleaser-cli-debug validate
```

### Common Debug Scenarios

1. **Configuration Issues**:
   ```bash
   # Check configuration loading
   go run cmd/goreleaser-cli/main.go config show --debug
   ```

2. **Validation Problems**:
   ```bash
   # Run validation with verbose output
   go run cmd/goreleaser-cli/main.go validate --verbose
   ```

3. **Server Issues**:
   ```bash
   # Start server with debug logging
   go run cmd/goreleaser-cli/main.go server --debug --port 8080
   ```

## Performance

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./benchmarks/
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./benchmarks/
go tool pprof mem.prof

# Trace profiling
go test -trace=trace.out ./benchmarks/
go tool trace trace.out
```

### Optimization Tips

1. **Use build tags** for optional features
2. **Implement connection pooling** for external services
3. **Cache expensive operations**
4. **Use sync.Pool** for object reuse
5. **Profile before optimizing**

## Common Issues

### 1. Go Module Issues

**Problem**: `go: module not found`
**Solution**:
```bash
go clean -modcache
go mod download
go mod tidy
```

### 2. Build Failures

**Problem**: `build constraints exclude all Go files`
**Solution**:
```bash
# Check build tags
go list -f '{{.GoFiles}}' ./cmd/goreleaser-cli
go build -tags=dev ./cmd/goreleaser-cli
```

### 3. Test Failures

**Problem**: Tests fail in CI but pass locally
**Solution**:
```bash
# Run tests with same flags as CI
go test -race -v ./...

# Check for time-dependent tests
TZ=UTC go test ./...
```

### 4. Import Cycle Issues

**Problem**: `import cycle not allowed`
**Solution**:
- Extract common interfaces to separate package
- Use dependency injection
- Refactor to remove circular dependencies

### 5. Performance Issues

**Problem**: Slow startup or execution
**Solution**:
```bash
# Profile the application
go build -o goreleaser-cli cmd/goreleaser-cli/main.go
time ./goreleaser-cli validate

# Use benchmarks to identify bottlenecks
go test -bench=. -benchmem ./benchmarks/
```

### 6. Docker Issues

**Problem**: Container won't start
**Solution**:
```bash
# Check logs
docker logs goreleaser-cli

# Debug container
docker run -it --entrypoint=/bin/sh goreleaser-cli

# Rebuild image
just docker-build-clean
```

## Contributing

### Code Style

1. **Follow Go conventions**: Use `gofmt` and `goimports`
2. **Write clear variable names**: Prefer clarity over brevity
3. **Keep functions small**: Aim for < 30 lines
4. **Document public APIs**: Include examples in godoc
5. **Handle errors explicitly**: Don't ignore errors

### Commit Guidelines

1. **Use conventional commits**: `feat:`, `fix:`, `docs:`, etc.
2. **Write clear commit messages**: Explain what and why
3. **Keep commits atomic**: One logical change per commit
4. **Reference issues**: Include issue numbers in commits

### Pull Request Process

1. **Create feature branch** from main
2. **Write comprehensive tests** for new functionality
3. **Update documentation** as needed
4. **Run full test suite** before submitting
5. **Request review** from maintainers

### Development Commands Reference

```bash
# Setup
just setup           # Initial project setup
just install         # Install all dependencies
just update          # Update dependencies

# Development
just dev             # Start development mode
just build           # Build for development
just build-release   # Build for production

# Testing
just test            # Run all tests
just test-unit       # Run unit tests only
just test-integration # Run integration tests only
just bench           # Run benchmarks
just coverage        # Generate coverage report

# Code Quality
just fmt             # Format code
just lint            # Run linters
just security-scan   # Security analysis
just check           # Run all checks

# Documentation
just docs            # Generate documentation
just serve-docs      # Serve docs locally

# Deployment
just docker-build    # Build Docker image
just docker-run      # Run Docker container
just release         # Create release
```

For more specific command details, run `just --list` or check the `justfile`.

---

**Need help?** 
- Check the [API documentation](API.md)
- Review [existing issues](https://github.com/LarsArtmann/template-GoReleaser/issues)
- Join our discussions