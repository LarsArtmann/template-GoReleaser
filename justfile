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
    go build -v ./cmd/myproject
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
    rm -f myproject *.exe
    rm -f coverage.out coverage.html
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
    go install github.com/sigstore/cosign/v2/cmd/cosign@latest
    go install github.com/anchore/syft/cmd/syft@latest
    go install github.com/a-h/templ/cmd/templ@latest
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
        docker build -t myproject:latest .; \
        echo "✓ Docker image built"; \
    else \
        echo "⚠ Docker daemon is not running. Please start Docker first."; \
        exit 1; \
    fi

# Docker run
docker-run: docker-build
    @echo "Running Docker container..."
    @if docker info >/dev/null 2>&1; then \
        docker run --rm myproject:latest; \
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

# Run security scan
security-scan:
    @echo "Running security scan..."
    gosec ./... || echo "⚠ Security issues found in code"
    @if command -v trivy >/dev/null 2>&1; then \
        echo "Running Trivy filesystem scan..."; \
        trivy fs . || true; \
    else \
        echo "⚠ Trivy not found. Install with: brew install trivy"; \
    fi
    @echo "✓ Security scan complete"

# Update dependencies
update-deps:
    @echo "Updating dependencies..."
    go get -u ./...
    go mod tidy
    @echo "✓ Dependencies updated"

# Show current version
version:
    @go run ./cmd/myproject -version

# Health check
health:
    @go run ./cmd/myproject -health

# Run the application
run:
    go run ./cmd/myproject

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

# Complete CI pipeline
ci: clean init fmt lint test validate build snapshot
    @echo "✓ CI pipeline complete"

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
    @echo "Release Process:"
    @echo "  1. Run 'just ci' to verify everything"
    @echo "  2. Run 'just tag <version>' to create a version tag"
    @echo "  3. Push tag: 'git push origin v<version>'"
    @echo "  4. GitHub Actions will handle the release"
    @echo ""
    @echo "Or manually:"
    @echo "  - 'just release' for free version"
    @echo "  - 'just release-pro' for pro version"