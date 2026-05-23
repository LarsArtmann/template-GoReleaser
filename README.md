# GoReleaser Wizard

**Stop copy-pasting GoReleaser configs. Get a production-ready setup with one command.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.26+-blue.svg)](https://golang.org)
[![GoReleaser](https://img.shields.io/badge/powered%20by-GoReleaser-blue.svg)](https://goreleaser.com)

An interactive CLI tool that generates production-ready GoReleaser configurations for Go projects. It auto-detects your project structure, guides you through options with a rich TUI, and outputs `.goreleaser.yaml`, GitHub Actions workflows, Dockerfiles, and Homebrew formulas -- all following best practices.

## Quick Start

```bash
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest
goreleaser-wizard init
```

That's it. Your `.goreleaser.yaml` is ready.

## Demo

```
$ goreleaser-wizard init

  GoReleaser Wizard
  Let's configure your GoReleaser setup!

? Project Name › my-awesome-cli
? Binary Name › my-awesome-cli
? Main Package Path › ./cmd/my-awesome-cli
? Project Type › CLI Application
? Platforms › ✓ Linux  ✓ macOS  ✓ Windows  ✓ FreeBSD
? Architectures › ✓ amd64  ✓ arm64
? CGO Configuration › Disabled (recommended)
? Docker Support › Build and publish
? Git Provider › GitHub
? Include LDFlags? › Yes
? Enable Code Signing? › Yes
? Generate SBOM? › Yes
? Generate Homebrew Formula? › Yes
? Generate GitHub Actions? › Yes

Created  .goreleaser.yaml
Created  .github/workflows/release.yml

GoReleaser configuration initialized successfully!
```

## Installation

### Go Install

```bash
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest
```

### From Source

```bash
git clone https://github.com/LarsArtmann/GoReleaser-Wizard.git
cd GoReleaser-Wizard
go build -o goreleaser-wizard ./cmd/goreleaser-wizard
```

### Pre-built Binary

Download the latest release from the [releases page](https://github.com/LarsArtmann/GoReleaser-Wizard/releases).

### Homebrew

```bash
brew tap LarsArtmann/tap
brew install goreleaser-wizard
```

### Docker

```bash
docker pull ghcr.io/larsartmann/goreleaser-wizard:latest
```

## Commands

### `init` -- Interactive Wizard

Guides you through creating a complete configuration with a rich terminal UI:

```bash
goreleaser-wizard init
```

| Flag             | Default       | Description                                   |
| ---------------- | ------------- | --------------------------------------------- |
| `--force`        | `false`       | Overwrite existing configuration              |
| `--interactive`  | `true`        | Run the TUI wizard (set `false` for headless) |
| `--project-name` | auto-detected | Override project name                         |
| `--binary-name`  | auto-detected | Override binary name                          |
| `--main-path`    | auto-detected | Override main package path                    |
| `--project-type` | auto-detected | Override project type                         |

When stdout is not a terminal, `init` prints a helpful message with the non-interactive flag.

### `generate` -- Non-Interactive Generation

For CI pipelines and scripting:

```bash
goreleaser-wizard generate \
  --project-name my-project \
  --binary-name my-app \
  --project-type cli \
  --config-only=false
```

| Flag             | Default       | Description                                          |
| ---------------- | ------------- | ---------------------------------------------------- |
| `--force`        | `false`       | Overwrite existing configuration                     |
| `--config-only`  | `false`       | Generate only `.goreleaser.yaml` (no GitHub Actions) |
| `--project-name` | auto-detected | Override project name                                |
| `--binary-name`  | auto-detected | Override binary name                                 |
| `--main-path`    | auto-detected | Override main package path                           |
| `--project-type` | auto-detected | Override project type                                |

### `validate` -- Check Configuration

Validates your existing GoReleaser configuration:

```bash
goreleaser-wizard validate
goreleaser-wizard validate --verbose
goreleaser-wizard validate --fix
goreleaser-wizard validate --project-only
```

| Flag             | Default | Description                                |
| ---------------- | ------- | ------------------------------------------ |
| `--verbose`      | `false` | Show detailed validation output            |
| `--fix`          | `false` | Attempt to fix common issues automatically |
| `--project-only` | `false` | Validate project structure only            |

### `version` -- Version Info

```bash
goreleaser-wizard version
```

## What It Generates

### `.goreleaser.yaml`

Optimized build configuration with multi-platform support, archive generation, checksums, changelog, code signing (cosign), SBOM, Homebrew, Docker images, nFPM packages, Nix, and Scoop.

### `.github/workflows/release.yml`

Automated release pipeline triggered on tags -- multi-platform builds, Docker image publishing, code signing, and SBOM generation.

### `Dockerfile`

Multi-stage production Docker build with non-root user, health checks, and minimal final image.

### Homebrew Formula

Package manager formula with automatic CamelCase naming and proper install/test directives.

## Supported Project Types

The wizard adapts configuration to your project:

| Type                    | Description                                    |
| ----------------------- | ---------------------------------------------- |
| **CLI Application**     | Single binary with version info via ldflags    |
| **Web API**             | Includes Docker and health check configuration |
| **Library**             | Focuses on the CLI component if present        |
| **gRPC Service**        | Docker with multi-platform images              |
| **Microservice**        | Docker publishing, minimal cross-compilation   |
| **Desktop Application** | Platform-specific builds (Linux, macOS)        |
| **Daemon/Service**      | System service with nFPM packages              |
| **Command Line Tool**   | Lightweight single-binary distribution         |

## Platforms and Architectures

| Platform | Architectures |
| -------- | ------------- |
| Linux    | amd64, arm64  |
| macOS    | amd64, arm64  |
| Windows  | amd64, arm64  |
| FreeBSD  | amd64         |

## Docker Support

Four levels of Docker integration:

| Level             | Behavior                                              |
| ----------------- | ----------------------------------------------------- |
| None              | No Docker configuration                               |
| Build only        | Generates Dockerfile, builds images locally           |
| Publish only      | Pushes pre-built images to registry                   |
| Build and Publish | Full pipeline: build + push with multi-arch manifests |

Supported registries: GitHub Container Registry (ghcr.io), Docker Hub, GitLab Registry, Quay.io, and custom registries.

## Git Providers

| Provider    | CI/CD Support       |
| ----------- | ------------------- |
| GitHub      | GitHub Actions      |
| GitLab      | GitLab CI           |
| Bitbucket   | Bitbucket Pipelines |
| Gitea       | Manual              |
| Self-hosted | Manual              |

## Architecture

GoReleaser Wizard follows Domain-Driven Design with Clean Architecture:

```
internal/domain/       Pure business logic, zero external dependencies
internal/validation/   Field-level and business rule validation
internal/errors/       Typed error codes with recovery suggestions
internal/git/          Git operations
internal/config/       Configuration management (koanf)
cmd/goreleaser-wizard/ CLI layer (Cobra), TUI (huh), generators, workflows
templates/             Go templates for all generated files
```

The domain layer uses type-safe enums for all configuration options (platforms, architectures, project types, Docker support levels, signing levels, etc.) and enforces validation through a centralized error system with structured diagnostics.

## Development

```bash
just build          # Build the binary
just test           # Run tests
just fmt            # Format code
just ci             # Full CI pipeline (fmt + test + build + verify + check)
just clean          # Remove build artifacts
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

[MIT](LICENSE) -- Lars Artmann
