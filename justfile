#!/usr/bin/env just --justfile

# Simple GoReleaser Template Justfile

set shell := ["bash", "-uc"]

# Show available commands
default:
    @just --list

# Initialize project dependencies
init:
    @echo "Initializing project..."
    go mod tidy
    @echo "✓ Project initialized"

# Build the application
build:
    @echo "Building application..."
    go build -o goreleaser-wizard ./cmd/goreleaser-wizard
    @echo "✓ Build complete"

# Run tests
test:
    @echo "Running tests..."
    go test ./...
    @echo "✓ Tests complete"

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...
    @echo "✓ Code formatted"

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    go clean
    rm -f goreleaser-wizard goreleaser-cli
    @echo "✓ Clean complete"

# Verify project setup
verify:
    @echo "Verifying project setup..."
    ./verify.sh
    @echo "✓ Verification complete"

# Run GoReleaser in snapshot mode
snapshot:
    @echo "Building snapshot..."
    goreleaser build --snapshot --clean
    @echo "✓ Snapshot complete"

# Run GoReleaser check
check:
    @echo "Checking GoReleaser configuration..."
    goreleaser check
    @echo "✓ Configuration check complete"

# Full CI pipeline
ci: fmt test build verify check
    @echo "✓ CI pipeline complete"