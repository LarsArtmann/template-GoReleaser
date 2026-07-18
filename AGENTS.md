# GoReleaser-Wizard Agent Guide

**Project**: GoReleaser-Wizard - Interactive GoReleaser configuration wizard
**Repository**: https://github.com/LarsArtmann/GoReleaser-Wizard

---

## Overview

GoReleaser-Wizard is an interactive CLI tool that generates production-ready GoReleaser configurations. It follows Domain-Driven Design (DDD) principles with Clean Architecture, using a workflow-based job execution system for configuration generation and validation.

**Key Technologies**:

- Go 1.25.6
- Cobra CLI framework
- Viper for configuration
- Charm libraries (lipgloss, log) for terminal UI
- golangci-lint for code quality
- go-arch-lint for architecture enforcement

---

## Essential Commands

### Primary Development Commands

### Advanced Commands (dev/arch-lint.just)

The project uses an enterprise-grade linting justfile with comprehensive tooling:

### Project-Specific Commands

---

## Project Architecture

### Directory Structure

### Architectural Layers

The project implements Clean Architecture with DDD principles:

1. **Domain Layer** (`internal/domain/`)
   - Pure business logic with zero external dependencies
   - Core types: `SafeProjectConfig`, enumerations, validation
   - All business rules and domain constraints
   - Interfaces for repository and service abstractions

2. **Application Layer** (`cmd/goreleaser-wizard/`)
   - Workflow orchestration and job execution
   - CLI commands (Cobra-based)
   - Template generation
   - File system operations

3. **Infrastructure Layer** (`internal/git/`, `internal/utils/`)
   - External system integrations
   - Git command wrappers
   - Utility functions

4. **Configuration Layer** (`templates/`)
   - Go templates for file generation
   - Template data structures
   - Template escaping utilities

### Architecture Enforcement

The project uses **go-arch-lint** to enforce architectural boundaries:

- **Deep scanning enabled**: AST-based method call analysis (v1.14.0+)
- **Domain isolation**: Domain layer cannot import application/infrastructure
- **Test file exclusions**: `*_test.go` files excluded from architecture validation
- **Vendor support**: Can validate with or without vendor dependencies

**Key rules** (from `.go-arch-lint.yml`):

- Domain components can only depend on other domain components
- Application layer can depend on domain but not infrastructure
- Infrastructure implements domain interfaces
- Test files have no architectural constraints

---

## Code Patterns and Conventions

### Domain-Driven Design Patterns

**1. Safe Configuration Types**

```go
// Single source of truth for configuration
type SafeProjectConfig struct {
    ProjectName        string      `json:"project_name" yaml:"project_name"`
    ProjectType        ProjectType `json:"project_type" yaml:"project_type"`
    // ... with validation methods
}

func (spc *SafeProjectConfig) Validate() error { }
func (spc *SafeProjectConfig) ApplyDefaults() { }
```

**2. Typed Enumerations**

```go
type ProjectType string

const (
    ProjectTypeCLI           ProjectType = "cli"
    ProjectTypeWebService    ProjectType = "web-service"
    ProjectTypeLibrary       ProjectType = "library"
)

func (pt ProjectType) IsValid() bool { }
```

**3. Smart Defaults**

```go
func NewSafeProjectConfig() *SafeProjectConfig {
    return &SafeProjectConfig{
        ProjectType:    GetRecommendedProjectType(),
        Platforms:      GetRecommendedPlatforms(),
        CGOStatus:      CGOStatusDisabled,
        // ... intelligent defaults based on project analysis
    }
}
```

### Error Handling

**Centralized Error System** (`internal/errors/domain_errors.go`)

All errors use typed error codes for structured error handling:

**Error Usage Pattern**:

**Recovery from Panics**:

### Workflow and Job System

The project uses a sophisticated workflow-based execution system:

**Workflow Pattern**:

**Job Factory Pattern**:

**Job Implementation**:

### Validation Patterns

**Multi-Level Validation**:

### Template Generation Pattern

### Testing Patterns

**Table-Driven Tests**:

**Test Setup Pattern**:

### CLI Pattern (Cobra)

---

## Code Style and Conventions

### Naming Conventions

**Packages**: Lowercase, single word when possible

- `domain`, `validation`, `errors`, `generators`, `jobs`

**Types**: PascalCase

- `SafeProjectConfig`, `Workflow`, `JobFactory`, `DomainError`

**Interfaces**: PascalCase, usually end with "er"

- `Logger`, `Job`, `FileSystem`

**Functions**: PascalCase (exported), camelCase (unexported)

- `ValidateProjectName()`, `prepareTemplateData()`, `createBackup()`

**Constants**: PascalCase

- `ErrValidationFailed`, `ProjectTypeCLI`, `JobExecutionStatusRunning`

**Variables**: camelCase

- `projectName`, `templateData`, `jobManager`

**Acronyms**: Capitalized in PascalCase, lowercase in camelCase

- `HTTPResponse` (not HttpResponse)
- `httpClient` (not HTTPClient)
- `CGOStatus`, `SBOM`, `LDFlags` (well-known acronyms all caps)

### File Organization

**Keep files under 300 lines** - Split immediately when exceeded

Example file splitting pattern (from workflow.go):

```
workflow_core.go         - Core Workflow struct and basic methods
workflow_execution.go    - Execution logic and state management
workflow_types.go        - Workflow-related types and enums
workflow_templates.go    - Workflow templates and factories
workflow_validation.go    - Workflow validation logic
```

**Related code together** - Group by functionality, not by type

### Import Order

### Function Patterns

**Early returns preferred** over nested conditionals:

**Small, focused functions** - Prefer <30 lines

**Explicit over implicit** - Clear function signatures

---

## Testing Guidelines

### Test Structure

- **Unit tests**: `*_test.go` in same package
- **Integration tests**: `integration_test.go`
- **Performance tests**: `*_test.go` with `Benchmark*` functions

### Coverage Requirements

- **Minimum 80% coverage** enforced by `just coverage`
- Test both success and error paths
- Test edge cases and boundary conditions

### Test Utilities

**Use testify/assert**:

```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)
assert.Error(t, err)
assert.Contains(t, str, substr)
```

**Setup/Teardown**:

---

## Important Gotchas

### Architecture Constraints

**1. Domain Layer Purity**

- Domain layer MUST NOT import application or infrastructure
- Domain cannot use external dependencies (except through interfaces)
- Test domain layer in isolation - no external services

**2. Test File Exclusions**

- `*_test.go` files are excluded from architecture validation
- This allows testing patterns like package imports
- Don't rely on this for production code

**3. Circular Dependencies**

- go-arch-lint will catch circular dependencies
- Import cycles prevent compilation
- Use dependency injection to break cycles

### Error Handling

**4. Always Use Domain Errors**

- Never use `errors.New()` or `fmt.Errorf()` directly
- Always create errors with `errors.New*Error()`
- Provides structured error context and recovery suggestions

**5. Context Propagation**

- All async operations should accept `context.Context`
- Check `ctx.Err()` for cancellation
- Pass context through to downstream operations

### Template System

**6. Template Escaping**

- Go templates use `{{` and `}}` delimiters
- Must escape literal braces with `{{"{{"}}` and `{{"}}"}}`
- Use `internal/validation/template_escaping.go` utilities

**7. Template Functions**

- Custom template functions registered in `template.FuncMap`
- Example: `"incpatch": git.IncPatchVersion`
- Functions must be exported and error-safe

### File Operations

**8. Use `os.WriteFile` with Explicit Permissions**

```go
// Good
os.WriteFile("config.yaml", data, 0o644)

// Avoid
ioutil.WriteFile("config.yaml", data, 0o644)  // Deprecated
```

**9. Backup Before Overwrite**

- Always create backups before overwriting existing files
- Pattern used in generators: `createBackup()` method

### Build and Dependencies

**10. Justfile Priority**

- Prefer `just` commands over manual commands
- `just test` instead of `go test ./...`
- `just build` instead of `go build ./...`

> This file was auto-trimmed. Full history in git.
