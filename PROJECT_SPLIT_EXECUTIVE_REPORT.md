# GoReleaser-Wizard Project Split Executive Report

## Introduction

This report outlines a strategic proposal to split the existing `GoReleaser-Wizard` project into multiple, highly focused, and independently manageable projects. This approach aims to enhance modularity, reusability, maintainability, and allow for more independent development and deployment cycles.

## Current Project Structure Overview

The current `GoReleaser-Wizard` project effectively implements a Clean Architecture with Domain-Driven Design principles. It comprises:

- **`cmd/goreleaser-wizard/`**: The main CLI application, including workflow orchestration, job management, and configuration generators.
- **`internal/domain/`**: The core business logic, configuration types (`SafeProjectConfig`), enumerations, validation rules, and domain interfaces. This layer is designed to be pure, with no external dependencies.
- **`internal/validation/`**: Utilities for various validation checks.
- **`internal/errors/`**: A centralized, structured error handling system.
- **`internal/git/`**: Utilities for Git operations.
- **`internal/utils/`**: General utility functions.
- **`templates/`**: Go templates for generating configuration files.
- **Development Tooling**: `justfile`, `.go-arch-lint.yml`, `.golangci.yml`, `dev/arch-lint.just` for build, test, lint, and architecture enforcement.

## Proposed Project Split

Based on the current architecture and principles, the project can be logically divided into two primary, highly focused projects:

### 1. `goreleaser-wizard-core` (Go Module/Library)

This project would encapsulate the fundamental, business-agnostic core logic that defines the GoReleaser configuration domain.

- **Responsibility**:
  - Definition of `SafeProjectConfig` and related domain types.
  - All domain enumerations (e.g., `ProjectType`).
  - Core, pure validation logic inherent to the domain.
  - The structured error handling system (`internal/errors/domain_errors.go`).
  - Domain interfaces and abstractions.
- **Current Location Mapping**: Primarily `internal/domain/` and `internal/errors/`.
- **Dependencies**: Zero external dependencies, ensuring its purity and maximum reusability.
- **Output**: A Go module (e.g., `github.com/LarsArtmann/goreleaser-wizard-core`).

### 2. `goreleaser-wizard-cli` (Go Application)

This project would house the interactive command-line application, responsible for user interaction, workflow orchestration, and leveraging the core library to generate configurations.

- **Responsibility**:
  - The main Cobra CLI application and all its commands (`init`, `validate`, `generate`).
  - Workflow orchestration and job execution system.
  - Configuration generators (GoReleaser, GitHub Actions, Dockerfile, Homebrew).
  - All template handling (`templates/`).
  - Infrastructure concerns like Git operations (`internal/git/`) and general utilities (`internal/utils/`).
  - Specific validation utilities that might depend on external contexts or the core types.
- **Current Location Mapping**: `cmd/goreleaser-wizard/`, `internal/validation/`, `internal/git/`, `internal/utils/`, `templates/`.
- **Dependencies**:
  - **`goreleaser-wizard-core`**: This would be the primary dependency, consumed as a Go module.
  - Existing external dependencies like Cobra, Viper, Charm libraries.
- **Output**: An executable CLI application.

## Advantages of this Split

1. **Enhanced Modularity and Reusability**:
   - The `goreleaser-wizard-core` library can be independently imported and used by other Go projects, APIs, or services that need to interact with GoReleaser configurations without the overhead of the CLI.
   - Enables a clear separation between domain logic and application-specific concerns.

2. **Independent Development and Release Cycles**:
   - Each project can evolve at its own pace. Changes in the core configuration definitions or validation rules can be released as a new version of `goreleaser-wizard-core` without immediately requiring a new `goreleaser-wizard-cli` release, and vice versa.
   - Facilitates faster iteration on specific components.

3. **Improved Maintainability and Focus**:
   - Smaller, more focused codebases are easier to understand, test, and maintain.
   - Developers can concentrate on specific areas (e.g., core domain logic vs. CLI user experience) without navigating a larger codebase.

4. **Clearer Architectural Boundaries**:
   - Reinforces the Clean Architecture principles, making the dependency flow explicit (CLI depends on Core, not the other way around).
   - Simplifies architecture linting and enforcement within each project.

5. **Reduced Cognitive Load**:
   - New contributors can more easily onboard to a smaller, single-purpose project.

## Implementation Considerations

- **Go Module Management**: `goreleaser-wizard-cli` would include `goreleaser-wizard-core` as a Go module dependency (`go.mod` entry).
- **Refactoring**: Careful refactoring would be required to move files and adjust import paths.
- **Testing**: Each new project would have its own comprehensive test suite, ensuring the integrity of its specific functionalities.
- **CI/CD**: Independent CI/CD pipelines would be established for each project, potentially leading to faster feedback loops.
- **Monorepo vs. Polyrepo**: Initially, these could reside in a monorepo for easier management of shared tooling (like `justfile`, linting configs) but still be treated as distinct Go modules. A full polyrepo split could follow if deemed beneficial.

## Conclusion

This proposed split into `goreleaser-wizard-core` and `goreleaser-wizard-cli` provides a robust, scalable, and maintainable architecture for the GoReleaser Wizard ecosystem. It aligns with best practices for Go project organization and sets a strong foundation for future development and expansion.
