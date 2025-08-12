#!/usr/bin/env bash
set -euo pipefail

# Simple GoReleaser Template Verifier

echo "Verifying GoReleaser template..."

# Check required tools
echo "Checking required tools:"
for tool in go git goreleaser; do
    if command -v "$tool" &> /dev/null; then
        echo "✓ $tool is installed"
    else
        echo "✗ $tool is missing - please install it"
        exit 1
    fi
done

# Check project structure
echo "Checking project structure:"
for file in go.mod cmd/goreleaser-cli/main.go; do
    if [ -f "$file" ]; then
        echo "✓ $file exists"
    else
        echo "✗ $file is missing"
        exit 1
    fi
done

# Test build
echo "Testing build:"
if go build -o /tmp/test-build ./cmd/goreleaser-cli; then
    echo "✓ Project builds successfully"
    rm -f /tmp/test-build
else
    echo "✗ Build failed"
    exit 1
fi

echo "✓ Template verification complete!"