# Examples and Usage Guide

This document provides comprehensive examples and usage patterns for the GoReleaser CLI template project.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Advanced Configuration](#advanced-configuration)
- [Integration Examples](#integration-examples)
- [Custom Command Examples](#custom-command-examples)
- [CI/CD Integration](#cicd-integration)
- [Docker Examples](#docker-examples)
- [Development Workflows](#development-workflows)
- [Automation Scripts](#automation-scripts)

## Basic Usage

### Quick Start

```bash
# Install the CLI
go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest

# Initialize configuration
goreleaser-cli config init

# Validate your project
goreleaser-cli validate

# Generate a license
goreleaser-cli license generate MIT "Your Name"

# Start the web interface
goreleaser-cli server
```

### Basic Commands

```bash
# Check version information
goreleaser-cli version

# Show help for all commands
goreleaser-cli --help

# Show help for specific command
goreleaser-cli validate --help
```

## Advanced Configuration

### Custom Configuration File

Create a custom configuration file `.goreleaser-cli.yaml`:

```yaml
# Project configuration
project:
  name: "my-awesome-project"
  description: "An amazing Go project built with GoReleaser"
  version: "1.0.0"
  
# Author information
author:
  name: "John Doe"
  email: "john@example.com"
  
# License settings
license:
  type: "Apache-2.0"
  
# CLI behavior
cli:
  verbose: true
  colors: true
  
# Server configuration
server:
  host: "0.0.0.0"
  port: 3000
  debug: true
  static_dir: "./custom/static"
  
# Development settings
development:
  hot_reload: true
  debug_templates: true
  
# Validation settings
validation:
  strict_mode: true
  required_files:
    - "README.md"
    - "LICENSE"
    - "CHANGELOG.md"
  required_env_vars:
    - "GITHUB_TOKEN"
    - "DOCKER_REGISTRY"
```

### Environment-Based Configuration

Set up different configurations for different environments:

```bash
# Development
export GORELEASER_CLI_ENV=development
export GORELEASER_CLI_CONFIG=.goreleaser-cli.dev.yaml
export GORELEASER_CLI_LOG_LEVEL=debug

# Staging
export GORELEASER_CLI_ENV=staging
export GORELEASER_CLI_CONFIG=.goreleaser-cli.staging.yaml
export GORELEASER_CLI_LOG_LEVEL=info

# Production
export GORELEASER_CLI_ENV=production
export GORELEASER_CLI_CONFIG=.goreleaser-cli.prod.yaml
export GORELEASER_CLI_LOG_LEVEL=warn
```

## Integration Examples

### GitHub Actions Integration

Create `.github/workflows/goreleaser.yml`:

```yaml
name: GoReleaser Build and Release

on:
  push:
    tags:
      - 'v*'
  pull_request:
    branches: [ main ]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install goreleaser-cli
      run: go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest
    
    - name: Validate project
      run: |
        goreleaser-cli validate --all
        goreleaser-cli license validate
    
    - name: Run verification
      run: goreleaser-cli verify
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  release:
    if: startsWith(github.ref, 'refs/tags/')
    needs: validate
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        distribution: goreleaser
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### GitLab CI Integration

Create `.gitlab-ci.yml`:

```yaml
stages:
  - validate
  - build
  - release

variables:
  GO_VERSION: "1.21"

before_script:
  - apt-get update -qq && apt-get install -y -qq git curl
  - curl -fsSL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz
  - export PATH=/usr/local/go/bin:$PATH
  - go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest

validate:
  stage: validate
  script:
    - goreleaser-cli validate --all
    - goreleaser-cli license validate
  only:
    - merge_requests
    - main

build:
  stage: build
  script:
    - go build -o build/ ./...
    - goreleaser-cli verify
  artifacts:
    paths:
      - build/
  only:
    - main

release:
  stage: release
  script:
    - curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
    - ./bin/goreleaser release --clean
  only:
    - tags
  variables:
    GITLAB_TOKEN: $CI_JOB_TOKEN
```

### Jenkins Pipeline

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
        GORELEASER_CLI_ENV = 'ci'
    }
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    # Install Go
                    wget -O go.tar.gz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
                    tar -C /usr/local -xzf go.tar.gz
                    export PATH=/usr/local/go/bin:$PATH
                    
                    # Install goreleaser-cli
                    go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest
                '''
            }
        }
        
        stage('Validate') {
            steps {
                sh '''
                    export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH
                    goreleaser-cli validate --all
                    goreleaser-cli license validate
                '''
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    export PATH=/usr/local/go/bin:$PATH
                    go build -o build/ ./...
                '''
                archiveArtifacts artifacts: 'build/*', fingerprint: true
            }
        }
        
        stage('Release') {
            when {
                tag 'v*'
            }
            steps {
                sh '''
                    export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH
                    goreleaser-cli verify
                    
                    # Install and run GoReleaser
                    curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
                    ./bin/goreleaser release --clean
                '''
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}
```

## Custom Command Examples

### Batch License Generation

Create licenses for multiple projects:

```bash
#!/bin/bash
# generate-licenses.sh

projects=(
    "project1:MIT:John Doe"
    "project2:Apache-2.0:Jane Smith"
    "project3:BSD-3-Clause:Acme Corp"
)

for project_info in "${projects[@]}"; do
    IFS=':' read -r project license author <<< "$project_info"
    
    echo "Generating license for $project..."
    mkdir -p "$project"
    cd "$project"
    
    goreleaser-cli license generate "$license" "$author"
    
    cd ..
    echo "✅ License generated for $project"
done
```

### Automated Validation Script

```bash
#!/bin/bash
# validate-all.sh

set -e

echo "🔍 Running comprehensive validation..."

# Set validation options
export GORELEASER_CLI_ENV=validation
export GORELEASER_CLI_LOG_LEVEL=debug

# Run all validations
echo "📁 Validating project structure..."
goreleaser-cli validate --structure

echo "🌍 Validating environment..."
goreleaser-cli validate --env

echo "🛠️ Validating configuration..."
goreleaser-cli validate --config

echo "📄 Validating license..."
goreleaser-cli license validate

echo "🔒 Running security verification..."
goreleaser-cli verify --skip-dry-run

echo "🎉 All validations passed!"
```

### Configuration Management

```bash
#!/bin/bash
# setup-project.sh

PROJECT_NAME="$1"
AUTHOR_NAME="$2"
AUTHOR_EMAIL="$3"
LICENSE_TYPE="${4:-MIT}"

if [ $# -lt 3 ]; then
    echo "Usage: $0 <project-name> <author-name> <author-email> [license-type]"
    exit 1
fi

echo "🚀 Setting up project: $PROJECT_NAME"

# Initialize configuration
goreleaser-cli config init --force

# Set project information
goreleaser-cli config set project.name "$PROJECT_NAME"
goreleaser-cli config set project.description "Auto-generated project description"

# Set author information
goreleaser-cli config set author.name "$AUTHOR_NAME"
goreleaser-cli config set author.email "$AUTHOR_EMAIL"

# Set license
goreleaser-cli config set license.type "$LICENSE_TYPE"

# Generate license file
goreleaser-cli license generate "$LICENSE_TYPE" "$AUTHOR_NAME"

# Enable verbose output for development
goreleaser-cli config set cli.verbose true

echo "✅ Project setup complete!"
goreleaser-cli config show
```

## CI/CD Integration

### Pre-commit Hook Setup

Create `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: goreleaser-validate
        name: GoReleaser CLI Validation
        entry: goreleaser-cli validate --all
        language: system
        files: '\.(go|yaml|yml)$'
        pass_filenames: false
      
      - id: license-validate
        name: License Validation
        entry: goreleaser-cli license validate
        language: system
        files: '^LICENSE$'
        pass_filenames: false
```

### Makefile Integration

Create a `Makefile`:

```makefile
.PHONY: install validate build test release clean help

# Variables
BINARY_NAME=goreleaser-cli
GO_VERSION=1.21

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install goreleaser-cli
	go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest

validate: ## Run all validations
	$(BINARY_NAME) validate --all
	$(BINARY_NAME) license validate

build: validate ## Build the project
	go build -o build/ ./...

test: ## Run tests
	go test -v -race ./...

verify: ## Run comprehensive verification
	$(BINARY_NAME) verify

release: verify ## Create a release (requires tag)
	goreleaser release --clean

clean: ## Clean build artifacts
	rm -rf build/ dist/

setup: ## Initialize project configuration
	$(BINARY_NAME) config init
	$(BINARY_NAME) license generate MIT "$(shell git config user.name)"

dev: ## Start development server
	$(BINARY_NAME) server --debug --port 3000

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(BINARY_NAME)
```

### Docker Integration

Create a multi-stage `Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install goreleaser-cli
RUN go install github.com/LarsArtmann/template-GoReleaser/cmd/goreleaser-cli@latest

# Copy project files
COPY . .

# Validate and build
RUN goreleaser-cli validate --all
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git
WORKDIR /root/

# Copy binary and goreleaser-cli
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goreleaser-cli /usr/local/bin/

# Expose port
EXPOSE 8080

# Health check using goreleaser-cli
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD goreleaser-cli --health || exit 1

CMD ["./main"]
```

### Docker Compose Example

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GORELEASER_CLI_ENV=production
      - GORELEASER_CLI_LOG_LEVEL=info
    volumes:
      - ./config:/app/config:ro
    healthcheck:
      test: ["CMD", "goreleaser-cli", "--health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped

  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - app
    restart: unless-stopped
```

## Development Workflows

### Local Development Setup

```bash
#!/bin/bash
# dev-setup.sh

echo "🛠️ Setting up development environment..."

# Install dependencies
echo "📦 Installing dependencies..."
go mod download

# Install development tools
echo "🔧 Installing development tools..."
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup pre-commit hooks
echo "🪝 Setting up pre-commit hooks..."
pip install pre-commit
pre-commit install

# Initialize configuration
echo "⚙️ Initializing configuration..."
goreleaser-cli config init --force

# Set development-friendly defaults
goreleaser-cli config set cli.verbose true
goreleaser-cli config set cli.colors true

# Generate development license
goreleaser-cli license generate MIT "$(git config user.name)"

# Validate setup
echo "✅ Validating setup..."
goreleaser-cli validate --all

echo "🎉 Development environment ready!"
echo "💡 Start development with: make dev"
```

### Hot Reload Development

Create `air.toml`:

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = ["server", "--debug", "--port", "3000"]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/goreleaser-cli"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "templ"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
```

## Automation Scripts

### Release Automation

```bash
#!/bin/bash
# release.sh

set -e

VERSION="$1"
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

echo "🚀 Preparing release $VERSION..."

# Validate project
echo "🔍 Running validations..."
goreleaser-cli validate --all
goreleaser-cli license validate

# Update version in configuration
echo "📝 Updating version..."
goreleaser-cli config set project.version "${VERSION#v}"

# Run tests
echo "🧪 Running tests..."
go test -v ./...

# Create and push tag
echo "🏷️ Creating tag..."
git add .
git commit -m "chore: prepare release $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"
git push origin main --tags

# Verify before release
echo "✅ Running final verification..."
goreleaser-cli verify

echo "🎉 Release $VERSION prepared successfully!"
echo "💡 Run 'goreleaser release --clean' to publish"
```

### Multi-Environment Deployment

```bash
#!/bin/bash
# deploy.sh

ENVIRONMENT="$1"
VERSION="${2:-latest}"

if [ -z "$ENVIRONMENT" ]; then
    echo "Usage: $0 <environment> [version]"
    echo "Environments: dev, staging, prod"
    exit 1
fi

case "$ENVIRONMENT" in
    "dev")
        CONFIG_FILE=".goreleaser-cli.dev.yaml"
        PORT="3000"
        ;;
    "staging")
        CONFIG_FILE=".goreleaser-cli.staging.yaml"
        PORT="8080"
        ;;
    "prod")
        CONFIG_FILE=".goreleaser-cli.prod.yaml"
        PORT="80"
        ;;
    *)
        echo "❌ Unknown environment: $ENVIRONMENT"
        exit 1
        ;;
esac

echo "🚀 Deploying to $ENVIRONMENT (version: $VERSION)..."

# Set environment
export GORELEASER_CLI_CONFIG="$CONFIG_FILE"
export GORELEASER_CLI_ENV="$ENVIRONMENT"

# Validate configuration
echo "✅ Validating configuration for $ENVIRONMENT..."
goreleaser-cli validate --all

# Build Docker image
echo "🐳 Building Docker image..."
docker build -t "goreleaser-cli:$ENVIRONMENT-$VERSION" .

# Deploy
echo "📦 Deploying..."
docker run -d \
    --name "goreleaser-cli-$ENVIRONMENT" \
    -p "$PORT:8080" \
    -e "GORELEASER_CLI_ENV=$ENVIRONMENT" \
    --restart unless-stopped \
    "goreleaser-cli:$ENVIRONMENT-$VERSION"

echo "🎉 Deployment to $ENVIRONMENT completed!"
echo "🌐 Access at: http://localhost:$PORT"
```

### Batch Project Setup

```bash
#!/bin/bash
# batch-setup.sh

# Configuration file with project definitions
PROJECTS_FILE="projects.txt"

if [ ! -f "$PROJECTS_FILE" ]; then
    cat > "$PROJECTS_FILE" << 'EOF'
# Format: project-name|author-name|author-email|license-type|description
my-api|John Doe|john@example.com|MIT|REST API service
my-cli|Jane Smith|jane@example.com|Apache-2.0|Command line tool
my-lib|Bob Johnson|bob@example.com|BSD-3-Clause|Utility library
EOF
    echo "📝 Created $PROJECTS_FILE template. Edit it and run again."
    exit 0
fi

while IFS='|' read -r project author email license description; do
    # Skip comments and empty lines
    [[ "$project" =~ ^#.*$ ]] && continue
    [[ -z "$project" ]] && continue
    
    echo "🚀 Setting up project: $project"
    
    mkdir -p "$project"
    cd "$project"
    
    # Initialize Git if needed
    if [ ! -d ".git" ]; then
        git init
        git config user.name "$author"
        git config user.email "$email"
    fi
    
    # Setup goreleaser-cli
    goreleaser-cli config init --force
    goreleaser-cli config set project.name "$project"
    goreleaser-cli config set project.description "$description"
    goreleaser-cli config set author.name "$author"
    goreleaser-cli config set author.email "$email"
    goreleaser-cli config set license.type "$license"
    
    # Generate license
    goreleaser-cli license generate "$license" "$author"
    
    # Initial commit
    git add .
    git commit -m "feat: initial project setup with goreleaser-cli"
    
    cd ..
    echo "✅ Project $project setup complete"
done < "$PROJECTS_FILE"

echo "🎉 All projects setup completed!"
```

### Health Check and Monitoring

```bash
#!/bin/bash
# monitor.sh

ENDPOINTS=(
    "http://localhost:8080/health"
    "http://localhost:8080/metrics"
)

SLACK_WEBHOOK="${SLACK_WEBHOOK_URL}"
LOG_FILE="/var/log/goreleaser-cli-monitor.log"

check_endpoint() {
    local endpoint="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    if curl -sf "$endpoint" > /dev/null 2>&1; then
        echo "[$timestamp] ✅ $endpoint is healthy" | tee -a "$LOG_FILE"
        return 0
    else
        echo "[$timestamp] ❌ $endpoint is down" | tee -a "$LOG_FILE"
        
        # Send alert to Slack
        if [ -n "$SLACK_WEBHOOK" ]; then
            curl -X POST -H 'Content-type: application/json' \
                --data "{\"text\":\"🚨 Alert: $endpoint is down\"}" \
                "$SLACK_WEBHOOK"
        fi
        
        return 1
    fi
}

main() {
    echo "🔍 Starting health check monitoring..."
    
    while true; do
        all_healthy=true
        
        for endpoint in "${ENDPOINTS[@]}"; do
            if ! check_endpoint "$endpoint"; then
                all_healthy=false
            fi
        done
        
        if [ "$all_healthy" = true ]; then
            echo "$(date '+%Y-%m-%d %H:%M:%S') 🎉 All endpoints healthy"
        fi
        
        # Wait 30 seconds before next check
        sleep 30
    done
}

# Handle signals
trap 'echo "Stopping monitor..."; exit 0' SIGTERM SIGINT

main
```

---

These examples provide comprehensive coverage of how to use the GoReleaser CLI template in various scenarios. For more specific use cases or custom integrations, refer to the [API documentation](../docs/API.md) or examine the source code.