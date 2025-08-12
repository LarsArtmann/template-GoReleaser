# GoReleaser Template

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)
[![Just](https://img.shields.io/badge/Just-Task%20Runner-blue.svg)](https://github.com/casey/just)

A comprehensive GoReleaser template for Go projects with extensive automation, validation scripts, and both free and pro configurations.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Release Process](#release-process)
- [Available Commands](#available-commands)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Features

### Core Features
- **Complete GoReleaser Setup**: Pre-configured for both free and pro versions
- **Automated Validation**: Comprehensive verification scripts for configurations
- **Just Integration**: Extensive task automation with Just commands
- **Multi-Platform Builds**: Support for multiple OS/architecture combinations
- **Docker Support**: Container builds and multi-platform images
- **Security Features**: Code signing, SBOM generation, and security scanning
- **License Templates**: Multiple license templates (MIT, Apache-2.0, BSD-3-Clause, EUPL-1.2)

### Development Tools
- **Code Quality**: Linting, formatting, and security scanning
- **Testing**: Comprehensive test automation with coverage reports
- **CI/CD Ready**: GitHub Actions compatible
- **Environment Management**: `.env` file setup and management
- **Dependency Management**: Automated updates and security checks

### Advanced Features (Pro)
- **Enhanced Security**: Advanced signing and verification
- **Template Generation**: Support for templ templates
- **Advanced Compression**: UPX compression support
- **Extended Validation**: Strict validation with additional checks

## Quick Start

1. **Use as GitHub Template**:
   ```bash
   gh repo create my-project --template LarsArtmann/template-GoReleaser
   cd my-project
   ```

2. **Or Clone Directly**:
   ```bash
   git clone https://github.com/LarsArtmann/template-GoReleaser.git my-project
   cd my-project
   rm -rf .git && git init
   ```

3. **Initialize Project**:
   ```bash
   just setup-env
   just install-tools
   just init
   ```

4. **Validate Setup**:
   ```bash
   just validate
   ```

## Installation

### Prerequisites

- **Go 1.24+**: [Download Go](https://golang.org/dl/)
- **Git**: Version control
- **Docker** (optional): For container builds
- **Just** (recommended): Task runner

### Install Just (Task Runner)

```bash
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ~/bin

# Windows
scoop install just
```

### Install Development Tools

```bash
just install-tools
```

This installs:
- GoReleaser
- golangci-lint
- gosec
- cosign
- syft
- templ

## Usage

### Environment Setup

1. **Create Environment File**:
   ```bash
   just setup-env
   ```

2. **Edit `.env` with Your Configuration**:
   ```bash
   # Example .env contents
   GITHUB_TOKEN=your_github_token
   DOCKER_REGISTRY=ghcr.io
   PROJECT_NAME=my-project
   ```

### Basic Development Workflow

```bash
# Build the project
just build

# Run tests
just test

# Run with coverage
just test-coverage

# Format code
just fmt

# Run linters
just lint

# Full CI pipeline
just ci
```

### Release Workflow

```bash
# Validate configuration
just validate

# Create a snapshot build (no release)
just snapshot

# Dry run (test release process)
just dry-run

# Create version tag
just tag v1.0.0

# Push tag to trigger release
git push origin v1.0.0
```

## Configuration

### Project Structure

```
├── cmd/myproject/          # Application entry point
│   └── main.go
├── assets/                 # License templates and assets
│   └── licenses/
├── templates/              # Additional templates
├── scripts/                # Build and utility scripts
├── .readme/configs/        # README generator config
├── justfile                # Task automation
├── verify.sh              # Configuration verifier
├── validate-strict.sh     # Strict validation
├── Dockerfile             # Container configuration
└── go.mod                 # Go module definition
```

### GoReleaser Configurations

The template supports two GoReleaser configurations:

- **`.goreleaser.yaml`**: Free version with basic features
- **`.goreleaser.pro.yaml`**: Pro version with advanced features

### License Templates

Choose from multiple license templates in `assets/licenses/`:
- MIT License
- Apache License 2.0
- BSD 3-Clause License
- EUPL 1.2 License

## Development

### Running the Application

```bash
# Run directly
just run

# Run with specific flags
go run ./cmd/myproject -version
go run ./cmd/myproject -health

# Build and run
just build
./myproject
```

### Testing

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# View coverage report
open coverage.html
```

### Code Quality

```bash
# Format code
just fmt

# Run linters
just lint

# Security scan
just security-scan

# Complete quality check
just ci
```

## Release Process

### Automated (Recommended)

1. **Prepare Release**:
   ```bash
   just ci                    # Ensure everything passes
   just validate             # Validate GoReleaser config
   ```

2. **Create Release**:
   ```bash
   just tag v1.0.0          # Create version tag
   git push origin v1.0.0   # Trigger release
   ```

3. **GitHub Actions** handles the rest automatically.

### Manual Release

```bash
# For manual releases (requires tag)
just release              # Free version
just release-pro          # Pro version (requires license)
```

## Available Commands

### Core Commands

| Command | Description |
|---------|-------------|
| `just init` | Initialize project dependencies |
| `just build` | Build the application |
| `just test` | Run tests with coverage |
| `just fmt` | Format code |
| `just lint` | Run linters |
| `just clean` | Clean build artifacts |

### Validation Commands

| Command | Description |
|---------|-------------|
| `just validate` | Validate GoReleaser configuration |
| `just validate-strict` | Run strict validation checks |
| `just check` | Check free GoReleaser config |
| `just check-pro` | Check pro GoReleaser config |

### Release Commands

| Command | Description |
|---------|-------------|
| `just snapshot` | Build snapshot (no release) |
| `just snapshot-pro` | Build pro snapshot |
| `just dry-run` | Test release process |
| `just tag <version>` | Create version tag |
| `just release` | Create release |
| `just release-pro` | Create pro release |

### Development Tools

| Command | Description |
|---------|-------------|
| `just setup-env` | Create .env from template |
| `just install-tools` | Install development tools |
| `just docker-build` | Build Docker image |
| `just security-scan` | Run security analysis |
| `just update-deps` | Update dependencies |
| `just ci` | Complete CI pipeline |

### Utility Commands

| Command | Description |
|---------|-------------|
| `just version` | Show application version |
| `just health` | Health check |
| `just changelog` | Generate changelog |
| `just help` | Show detailed help |

## Troubleshooting

### Common Issues

#### GoReleaser Validation Fails
```bash
# Check configuration syntax
just validate

# Run with verbose output
goreleaser check --config .goreleaser.yaml

# Validate project structure
./verify.sh
```

#### Build Fails
```bash
# Check Go version
go version  # Should be 1.23+

# Update dependencies
just update-deps

# Clean and rebuild
just clean
just build
```

#### Missing Tools
```bash
# Install all required tools
just install-tools

# Check specific tool
which goreleaser
which golangci-lint
```

#### Docker Issues
```bash
# Test Docker build
just docker-build

# Check Docker daemon
docker info
```

### Environment Issues

#### Missing Environment Variables
```bash
# Check .env file exists
ls -la .env

# Recreate from template
just setup-env
```

#### Git Configuration
```bash
# Check Git setup
git remote -v
git status

# Verify tags
git tag
```

### Getting Help

1. **Check Logs**: Most commands provide verbose output
2. **Run Validation**: Use `just validate` and `./verify.sh`
3. **Check Documentation**: Run `just help` for command details
4. **Inspect Configuration**: Review `.goreleaser.yaml` files

### Debug Mode

```bash
# Run any command with verbose output
goreleaser --debug

# Check environment
env | grep -E "(GITHUB|DOCKER)"

# Validate step by step
just validate-strict
```

## Contributing

1. **Fork the Repository**
2. **Create Feature Branch**: `git checkout -b feature/amazing-feature`
3. **Make Changes**: Follow the existing code style
4. **Add Tests**: Ensure new features are tested
5. **Run CI**: `just ci` must pass
6. **Submit Pull Request**: With clear description

### Development Guidelines

- Follow Go conventions and gofmt
- Add tests for new features
- Update documentation as needed
- Ensure all linters pass
- Validate GoReleaser configurations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Additional License Templates

This template includes several license options in `assets/licenses/`:
- MIT License (default)
- Apache License 2.0
- BSD 3-Clause License
- EUPL 1.2 License

## Support

- **Documentation**: Check the `justfile` and validation scripts
- **Issues**: [GitHub Issues](https://github.com/LarsArtmann/template-GoReleaser/issues)
- **Discussions**: [GitHub Discussions](https://github.com/LarsArtmann/template-GoReleaser/discussions)

---

**⚡ Quick Commands Reference**:
```bash
just setup-env && just install-tools  # Initial setup
just ci                                # Full validation
just tag v1.0.0 && git push origin v1.0.0  # Release
```