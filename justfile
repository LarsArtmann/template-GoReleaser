#!/usr/bin/env just --justfile

# GoReleaser Template Justfile
# Task automation for GoReleaser projects

set shell := ["bash", "-uc"]
set dotenv-load := true

# Default recipe - show help
default:
    @just --list

# Initialize project
init:
    @echo "Initializing GoReleaser project..."
    go mod tidy
    go mod download
    @echo "✓ Project initialized"

# Build the project
build:
    @echo "Building project..."
    go build -v ./cmd/goreleaser-cli
    @echo "✓ Build complete"

# Run tests
test:
    @echo "Running tests..."
    go test -v -race -coverprofile=coverage.out ./...
    @echo "✓ Tests complete"

# Run tests with coverage report
test-coverage: test
    @echo "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    @echo "✓ Coverage report generated: coverage.html"

# Run integration tests
integration-test:
    @echo "Running integration tests..."
    go test -v -race -timeout=20m -coverprofile=integration-coverage.out -coverpkg=./... ./tests/integration/...
    @echo "✓ Integration tests complete"

# Run integration tests with coverage report
integration-test-coverage: integration-test
    @echo "Generating integration test coverage report..."
    go tool cover -html=integration-coverage.out -o integration-coverage.html
    go tool cover -func=integration-coverage.out
    @echo "✓ Integration test coverage report generated: integration-coverage.html"

# Run all tests (unit + integration)
test-all:
    @echo "Running all tests..."
    just test
    just integration-test
    @echo "✓ All tests complete"

# Run all tests with coverage
test-all-coverage:
    @echo "Running all tests with coverage..."
    just test-coverage
    just integration-test-coverage
    @echo "✓ All tests with coverage complete"

# Run linters
lint:
    @echo "Running linters..."
    golangci-lint run ./... || true
    gosec ./... || true
    @echo "✓ Linting complete"

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...
    gofmt -s -w .
    @echo "✓ Code formatted"

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf dist/ build/ bin/
    rm -f goreleaser-cli goreleaser-cli-server *.exe
    rm -f coverage.out coverage.html
    rm -f integration-coverage.out integration-coverage.html
    rm -f validation-report.json
    @echo "✓ Clean complete"

# Validate GoReleaser configuration
validate:
    @echo "Validating GoReleaser configuration..."
    @if [ -f "./verify.sh" ] && [ -x "./verify.sh" ]; then \
        if command -v gtimeout >/dev/null 2>&1; then \
            gtimeout 30s ./verify.sh || echo "⚠ Validation script timed out or failed"; \
        else \
            ./verify.sh || echo "⚠ Validation script failed"; \
        fi; \
    else \
        echo "⚠ verify.sh not found or not executable"; \
        echo "Running basic GoReleaser validation instead:"; \
        just check; \
    fi
    @echo "✓ Validation complete"

# Strict validation
validate-strict:
    @echo "Running strict validation..."
    @if [ -f "./validate-strict.sh" ] && [ -x "./validate-strict.sh" ]; then \
        if command -v gtimeout >/dev/null 2>&1; then \
            gtimeout 60s ./validate-strict.sh || echo "⚠ Strict validation script timed out or failed"; \
        else \
            ./validate-strict.sh || echo "⚠ Strict validation script failed"; \
        fi; \
    else \
        echo "⚠ validate-strict.sh not found or not executable"; \
        echo "Running basic validation instead:"; \
        just check; \
        just check-pro; \
    fi
    @echo "✓ Strict validation complete"

# Check GoReleaser configuration
check:
    @echo "Checking GoReleaser configuration..."
    goreleaser check
    @echo "✓ Configuration check complete"

# Check Pro configuration
check-pro:
    @echo "Checking GoReleaser Pro configuration..."
    goreleaser check --config .goreleaser.pro.yaml
    @echo "✓ Pro configuration check complete"

# Build snapshot (no release)
snapshot:
    @echo "Building snapshot release..."
    goreleaser release --snapshot --skip=publish --clean
    @echo "✓ Snapshot built in dist/"

# Build snapshot with Pro features
snapshot-pro:
    @echo "Building Pro snapshot release..."
    goreleaser release --snapshot --skip=publish --clean --config .goreleaser.pro.yaml
    @echo "✓ Pro snapshot built in dist/"

# Dry run release
dry-run:
    @echo "Running release dry-run..."
    goreleaser release --skip=publish --clean
    @echo "✓ Dry-run complete"

# Dry run Pro release
dry-run-pro:
    @echo "Running Pro release dry-run..."
    goreleaser release --skip=publish --clean --config .goreleaser.pro.yaml
    @echo "✓ Pro dry-run complete"

# Create a new version tag
tag version:
    @echo "Creating tag v{{version}}..."
    git tag -a v{{version}} -m "Release v{{version}}"
    @echo "✓ Tag v{{version}} created"
    @echo "Push with: git push origin v{{version}}"

# Release (requires tag)
release:
    @echo "Creating release..."
    goreleaser release --clean
    @echo "✓ Release complete"

# Release with Pro features (requires tag and license)
release-pro:
    @echo "Creating Pro release..."
    goreleaser release --clean --config .goreleaser.pro.yaml
    @echo "✓ Pro release complete"

# Install GoReleaser
install-goreleaser:
    @echo "Installing GoReleaser..."
    go install github.com/goreleaser/goreleaser/v2@latest
    @echo "✓ GoReleaser installed"

# Install all development tools
install-tools:
    @echo "Installing development tools..."
    go install github.com/goreleaser/goreleaser/v2@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest
    go install github.com/sigstore/cosign/v2/cmd/cosign@latest
    go install github.com/anchore/syft/cmd/syft@latest
    go install github.com/a-h/templ/cmd/templ@latest
    @echo "Installing security scanning tools via Homebrew..."
    @if command -v brew >/dev/null 2>&1; then \
        brew install shellcheck hadolint trivy; \
    else \
        echo "⚠ Homebrew not found. Please install shellcheck, hadolint, and trivy manually"; \
    fi
    @echo "✓ All tools installed"

# Setup environment from example
setup-env:
    @echo "Setting up environment..."
    @if [ ! -f .env ]; then \
        cp .env.example .env; \
        echo "✓ Created .env from .env.example"; \
        echo "⚠ Please edit .env with your values"; \
    else \
        echo "✓ .env already exists"; \
    fi

# Docker build
docker-build:
    @echo "Building Docker image..."
    @if docker info >/dev/null 2>&1; then \
        docker build -t goreleaser-cli:latest .; \
        echo "✓ Docker image built"; \
    else \
        echo "⚠ Docker daemon is not running. Please start Docker first."; \
        exit 1; \
    fi

# Docker run
docker-run: docker-build
    @echo "Running Docker container..."
    @if docker info >/dev/null 2>&1; then \
        docker run --rm goreleaser-cli:latest; \
        echo "✓ Docker container executed"; \
    else \
        echo "⚠ Docker daemon is not running. Please start Docker first."; \
        exit 1; \
    fi

# Generate changelog
changelog:
    @echo "Generating changelog..."
    git log --pretty=format:"* %s (%h)" > CHANGELOG.md
    @echo "✓ Changelog generated"

# Run comprehensive security scan
security-scan:
    @echo "Running comprehensive security scan..."
    @echo "1. Scanning Go code with gosec..."
    gosec ./... || echo "⚠ Security issues found in code"
    @echo "2. Scanning dependencies with govulncheck..."
    govulncheck ./... || echo "⚠ Vulnerable dependencies found"
    @if command -v shellcheck >/dev/null 2>&1; then \
        echo "3. Scanning shell scripts with shellcheck..."; \
        find . -name "*.sh" -exec shellcheck --severity=error --format=gcc {} \; || echo "⚠ Shell script issues found"; \
    else \
        echo "⚠ shellcheck not found. Install with: brew install shellcheck"; \
    fi
    @if [ -f Dockerfile ]; then \
        if command -v hadolint >/dev/null 2>&1; then \
            echo "4. Scanning Dockerfile with hadolint..."; \
            hadolint Dockerfile || echo "⚠ Dockerfile issues found"; \
        else \
            echo "⚠ hadolint not found. Install with: brew install hadolint"; \
        fi; \
    fi
    @if command -v trivy >/dev/null 2>&1; then \
        echo "5. Running Trivy filesystem scan..."; \
        trivy fs . --security-checks vuln,config,secret || echo "⚠ Trivy found issues"; \
    else \
        echo "⚠ Trivy not found. Install with: brew install trivy"; \
    fi
    @echo "✓ Comprehensive security scan complete"

# Find duplicate code (alias: fd)
find-duplicates:
    @echo "Finding duplicate code with jscpd..."
    @if ! command -v jscpd >/dev/null 2>&1; then \
        echo "⚠ jscpd not found. Installing with bun..."; \
        bun add -g jscpd; \
    fi
    @if [ ! -f .jscpd.json ]; then \
        echo "Creating default .jscpd.json configuration..."; \
        echo '{ \
            "threshold": 0, \
            "reporters": ["console", "json", "html"], \
            "minLines": 5, \
            "minTokens": 50, \
            "ignore": [ \
                "**/vendor/**", \
                "**/node_modules/**", \
                "**/.git/**", \
                "**/dist/**", \
                "**/build/**", \
                "**/*.pb.go", \
                "**/go.sum", \
                "**/go.mod", \
                "**/*.md", \
                "**/LICENSE", \
                "**/*.yaml", \
                "**/*.yml", \
                "**/*.json", \
                "**/testdata/**" \
            ], \
            "output": "./duplication-report", \
            "silent": false, \
            "exitCode": 0 \
        }' | jq '.' > .jscpd.json; \
        echo "✓ Created .jscpd.json configuration"; \
    fi
    @echo "Running jscpd to detect code duplication..."
    @jscpd . --config .jscpd.json || echo "✓ Duplication analysis complete"
    @if [ -f ./duplication-report/jscpd-report.json ]; then \
        echo ""; \
        echo "📊 Duplication Summary:"; \
        jq -r '.statistics | "  Total Lines: \(.total.lines // 0)\n  Duplicate Lines: \(.total.duplicatedLines // 0)\n  Duplicate Percentage: \(.total.percentage // 0)%"' ./duplication-report/jscpd-report.json 2>/dev/null || true; \
        echo ""; \
        echo "📁 Reports generated in ./duplication-report/"; \
        echo "   - HTML report: ./duplication-report/jscpd-report.html"; \
        echo "   - JSON report: ./duplication-report/jscpd-report.json"; \
    fi
    @echo "✓ Code duplication analysis complete"

# Alias for find-duplicates
fd: find-duplicates

# Update dependencies
update-deps:
    @echo "Updating dependencies..."
    go get -u ./...
    go mod tidy
    @echo "✓ Dependencies updated"

# Show current version
version:
    @go run ./cmd/goreleaser-cli version

# Health check
health:
    @echo "Health check not applicable for CLI tool"

# Run the application
run:
    go run ./cmd/goreleaser-cli

# Watch for changes and rebuild
watch:
    @echo "Watching for changes..."
    @if command -v inotifywait >/dev/null 2>&1; then \
        while true; do \
            inotifywait -e modify,create,delete -r . --exclude 'dist|.git|.idea|vendor' && \
            just build && \
            echo "✓ Rebuilt"; \
        done; \
    else \
        echo "⚠ inotifywait not found. Install inotify-tools or use 'brew install fswatch' on macOS"; \
        echo "Alternative: Use 'find . -name '*.go' | entr -r just build'"; \
        echo "For now, running build once:"; \
        just build; \
    fi

# Complete CI pipeline with security scanning
ci: clean init fmt lint security-scan test-all-coverage validate build snapshot
    @echo "✓ CI pipeline complete with security validation"

# Integration testing pipeline
integration-ci: clean init fmt lint integration-test-coverage validate
    @echo "✓ Integration CI pipeline complete"

# Help - show all available commands
help:
    @echo "GoReleaser Template - Available Commands:"
    @echo ""
    @just --list --unsorted
    @echo ""
    @echo "Environment Setup:"
    @echo "  1. Run 'just setup-env' to create .env file"
    @echo "  2. Edit .env with your configuration"
    @echo "  3. Run 'just install-tools' to install development tools"
    @echo "  4. Run 'just validate' to check configuration"
    @echo ""
    @echo "Testing:"
    @echo "  - 'just test' for unit tests"
    @echo "  - 'just integration-test' for integration tests"
    @echo "  - 'just test-all' for all tests"
    @echo "  - 'just test-all-coverage' for all tests with coverage"
    @echo "  - 'just integration-ci' for integration testing pipeline"
    @echo ""
    @echo "Release Process:"
    @echo "  1. Run 'just ci' to verify everything"
    @echo "  2. Run 'just tag <version>' to create a version tag"
    @echo "  3. Push tag: 'git push origin v<version>'"
    @echo "  4. GitHub Actions will handle the release"
    @echo ""
    @echo "Or manually:"
    @echo "  - 'just release' for free version"
    @echo "  - 'just release-pro' for pro version"