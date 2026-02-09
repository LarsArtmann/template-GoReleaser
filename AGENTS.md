# GoReleaser-Wizard Agent Guide

**Last Updated**: February 9, 2026
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

```bash
# Build the application
just build
# or: go build -o goreleaser-wizard ./cmd/goreleaser-wizard

# Run tests with coverage
just test
# or: go test ./...

# Format code
just fmt
# or: gofmt -w . && goimports -w .

# Full CI pipeline
just ci  # Runs: fmt, test, build, verify, check

# Clean build artifacts
just clean
```

### Advanced Commands (dev/arch-lint.just)

The project uses an enterprise-grade linting justfile with comprehensive tooling:

```bash
# View all available commands
just --list

# Complete linting pipeline
just lint  # Includes architecture, code quality, security, cycles

# Security audit
just security-audit

# Code coverage with threshold (default 80%)
just coverage 90

# Generate architecture dependency graphs
just graph
just graph-all  # Generate all graph types

# Run pre-commit checks
just check-pre-commit

# Format with enhanced formatters
just format  # Uses gofumpt + goimports
```

### Project-Specific Commands

```bash
# Validate GoReleaser configuration
goreleaser-wizard validate

# Initialize interactive wizard
goreleaser-wizard init

# Generate configuration non-interactively
goreleaser-wizard generate --name my-project --binary my-app

# Run GoReleaser checks
just check
just snapshot
```

---

## Project Architecture

### Directory Structure

```
GoReleaser-Wizard/
├── cmd/
│   └── goreleaser-wizard/         # Main CLI application
│       ├── main.go                # Entry point with Cobra setup
│       ├── init.go                # Interactive wizard implementation
│       ├── validate.go            # Configuration validation
│       ├── workflow.go            # Workflow orchestration
│       ├── jobs.go                # Job management system
│       ├── generators/            # Configuration generators
│       │   ├── goreleaser.go
│       │   ├── github_actions.go
│       │   ├── dockerfile.go
│       │   └── homebrew.go
│       ├── jobs/                 # Job types and factory
│       │   ├── types.go
│       │   ├── factory.go
│       │   └── implementations.go
│       └── types/                # Template data structures
│           └── template_data.go
├── internal/
│   ├── domain/                   # Domain layer (pure business logic)
│   │   ├── config_core.go        # Core configuration types
│   │   ├── config_defaults.go    # Smart defaults
│   │   ├── enums.go             # Domain enumerations
│   │   ├── validation.go         # Validation logic
│   │   ├── architecture.go       # Architecture configuration
│   │   ├── interfaces.go        # Domain interfaces
│   │   └── errors.go            # Domain error types
│   ├── validation/              # Validation utilities
│   │   ├── basic.go
│   │   ├── business_rules.go
│   │   ├── form_validator.go
│   │   └── template_escaping.go
│   ├── errors/                  # Centralized error definitions
│   │   └── domain_errors.go     # Typed error codes
│   ├── git/                     # Git operations
│   │   └── commands.go
│   └── utils/                   # Utility functions
│       └── recommendations.go
├── templates/                   # Go templates for configs
│   ├── goreleaser.yaml.tmpl
│   ├── github-actions.yml.tmpl
│   ├── Dockerfile.tmpl
│   └── homebrew.rb.tmpl
├── dev/
│   └── arch-lint.just           # Enterprise-grade linting commands
├── .go-arch-lint.yml           # Architecture rules
├── .golangci.yml              # Linting configuration
├── justfile                    # Primary task runner
└── go.mod                      # Go module definition
```

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

```go
// Error codes are strongly typed
type ErrorCode string

const (
    ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
    ErrInvalidProject   ErrorCode = "INVALID_PROJECT"
    ErrFileNotFound     ErrorCode = "FILE_NOT_FOUND"
    // ... 40+ defined error codes
)

// DomainError provides rich error context
type DomainError struct {
    Code     ErrorCode   `json:"code"`
    Message  string      `json:"message"`
    Details  string      `json:"details,omitempty"`
    Context  string      `json:"context,omitempty"`
    Cause    error       `json:"cause,omitempty"`
}

// Builder pattern for error construction
func NewValidationError(code ErrorCode, message, details string) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

func (de *DomainError) WithField(field string) *DomainError { }
func (de *DomainError) WithContext(ctx string) *DomainError { }
func (de *DomainError) WithCause(err error) *DomainError { }
```

**Error Usage Pattern**:
```go
// Validation errors
if name == "" {
    return errors.NewValidationError(
        errors.ErrInvalidProject,
        "Project name is required",
        "Project name cannot be empty",
    ).WithField("project_name").WithSuggestion("Choose a valid project name")
}

// System errors
if err := os.ReadFile(path); err != nil {
    return errors.NewSystemError(
        errors.ErrFileReadFailed,
        "Failed to read file",
        "Cannot read configuration",
        err,
    ).WithContext(path)
}
```

**Recovery from Panics**:
```go
func recoverFromPanic(context string) {
    if r := recover(); r != nil {
        logger.Error("Panic recovered", "context", context, "panic", r)

        err := errors.NewSystemError(
            errors.ErrTemplateExecutionFailed,
            "Unexpected error occurred",
            fmt.Sprintf("The wizard encountered an unexpected problem: %v", r),
            fmt.Errorf("panic: %v", r),
        ).WithContext(context)

        displayError(err)
        os.Exit(1)
    }
}
```

### Workflow and Job System

The project uses a sophisticated workflow-based execution system:

**Workflow Pattern**:
```go
// Workflow orchestrates multiple jobs
type Workflow struct {
    Name        string
    JobManager  *JobManager
    Factory     *JobFactory
    Timeout     time.Duration
}

// Execute runs all jobs with rollback support
func (w *Workflow) Execute(ctx context.Context) error {
    // Execute jobs through JobManager
    err := w.JobManager.ExecuteJobs(ctx)
    if err != nil {
        // Automatic rollback on failure
        rollbackErr := w.JobManager.RollbackFailedJobs(ctx)
        // ...
    }
    return nil
}
```

**Job Factory Pattern**:
```go
// JobFactory creates jobs based on configuration
type JobFactory struct {
    logger *log.Logger
}

func (jf *JobFactory) CreateConfigJob(config *domain.SafeProjectConfig) Job {
    return &ConfigGenerationJob{
        Config: config,
        Logger: jf.logger,
    }
}
```

**Job Implementation**:
```go
type ConfigGenerationJob struct {
    Config *domain.SafeProjectConfig
    Logger *log.Logger
}

func (j *ConfigGenerationJob) Execute(ctx context.Context) error {
    j.Logger.Info("Generating GoReleaser configuration")

    generator := generators.NewGoReleaserGenerator(j.Config, j.Logger)
    return generator.Generate(ctx)
}

func (j *ConfigGenerationJob) Rollback(ctx context.Context) error {
    // Clean up generated files
    return os.Remove(".goreleaser.yaml")
}
```

### Validation Patterns

**Multi-Level Validation**:
```go
// 1. Field-level validation (internal/validation/basic.go)
func ValidateProjectName(name string) error {
    if name == "" {
        return errors.NewValidationError(...)
    }
    if len(name) > 50 {
        return errors.NewValidationError(...)
    }
    if !projectNamePattern.MatchString(name) {
        return errors.NewValidationError(...)
    }
    return nil
}

// 2. Configuration-level validation (internal/domain/validation.go)
func (spc *SafeProjectConfig) Validate() error {
    if err := ValidateProjectName(spc.ProjectName); err != nil {
        return err
    }
    if err := ValidateBinaryName(spc.BinaryName); err != nil {
        return err
    }
    // ... validate all fields
    return nil
}

// 3. Business rule validation (internal/validation/business_rules.go)
func ValidateBusinessRules(config *domain.SafeProjectConfig) error {
    if config.DockerSupport && config.DockerRegistry == "" {
        return errors.NewValidationError(
            errors.ErrInvalidConfig,
            "Docker registry required",
            "When Docker support is enabled, a registry must be specified",
        )
    }
    return nil
}
```

### Template Generation Pattern

```go
// Generator encapsulates template rendering
type GoReleaserGenerator struct {
    templateData *types.GoReleaserTemplateData
    logger       Logger
}

// Generate renders template and writes to file
func (g *GoReleaserGenerator) Generate(ctx context.Context) error {
    // 1. Parse template
    tmpl := template.New("goreleaser").Funcs(template.FuncMap{
        "incpatch": git.IncPatchVersion,
    })
    tmpl, err := tmpl.Parse(templates.GoReleaserTemplate)
    if err != nil {
        return errors.NewConfigError(
            errors.ErrTemplateParsing,
            "Failed to parse template",
            err.Error(),
        ).WithCause(err)
    }

    // 2. Prepare template data
    data, err := g.prepareTemplateData(ctx)
    if err != nil {
        return err
    }

    // 3. Execute template
    var output bytes.Buffer
    if err := tmpl.Execute(&output, data); err != nil {
        return errors.NewConfigError(
            errors.ErrTemplateRendering,
            "Failed to execute template",
            err.Error(),
        ).WithCause(err)
    }

    // 4. Write with backup
    if err := g.createBackup(".goreleaser.yaml"); err != nil {
        return err
    }
    return os.WriteFile(".goreleaser.yaml", output.Bytes(), 0o644)
}
```

### Testing Patterns

**Table-Driven Tests**:
```go
func TestValidateProjectName(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantErr     bool
        expectedErr *domain.DomainError
    }{
        {
            name:    "valid project name",
            input:   "my-awesome-project",
            wantErr: false,
        },
        {
            name:        "empty project name",
            input:       "",
            wantErr:     true,
            expectedErr: errors.NewValidationError(...),
        },
        {
            name:    "project name too long",
            input:   strings.Repeat("a", 51),
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validation.ValidateProjectName(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErr != nil {
                    var domainErr *domain.DomainError
                    assert.True(t, errors.As(err, &domainErr))
                    assert.Equal(t, tt.expectedErr.Code, domainErr.Code)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Test Setup Pattern**:
```go
func TestInitCommand(t *testing.T) {
    tests := []struct {
        name        string
        setupFunc   func() string
        expectError bool
    }{
        {
            name: "basic_init_command",
            setupFunc: func() string {
                dir, _ := os.MkdirTemp("", "wizard-test")
                os.WriteFile(dir+"/go.mod", []byte("module test"), 0o644)
                return dir
            },
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testDir := tt.setupFunc()
            defer os.RemoveAll(testDir)

            originalDir, _ := os.Getwd()
            os.Chdir(testDir)
            defer os.Chdir(originalDir)

            // Test execution...
        })
    }
}
```

### CLI Pattern (Cobra)

```go
// Command registration in main.go
func init() {
    cobra.OnInitialize(initConfig)

    // Global flags
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
    rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

    // Bind flags to viper
    viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

    // Add commands
    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(validateCmd)
    rootCmd.AddCommand(generateCmd)
}

// Command implementation
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize GoReleaser configuration",
    Run:   runInitWizard,
}

func runInitWizard(cmd *cobra.Command, args []string) {
    defer recoverFromPanic("init wizard")

    // Parse flags
    force, _ := cmd.Flags().GetBool("force")
    interactive, _ := cmd.Flags().GetBool("interactive")

    // Execute workflow
    // ...
}
```

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

```go
// 1. Standard library
import (
    "context"
    "errors"
    "fmt"
)

// 2. Internal packages
import (
    "github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
    "github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// 3. External dependencies
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/spf13/cobra"
)
```

### Function Patterns

**Early returns preferred** over nested conditionals:
```go
// Good
func validateConfig(config *Config) error {
    if config.Name == "" {
        return errors.New("name required")
    }
    if config.Path == "" {
        return errors.New("path required")
    }
    return nil
}

// Avoid
func validateConfig(config *Config) error {
    if config.Name != "" {
        if config.Path != "" {
            return nil
        }
        return errors.New("path required")
    }
    return errors.New("name required")
}
```

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
```go
func TestMain(m *testing.M) {
    // Setup
    setupTestEnvironment()

    // Run tests
    code := m.Run()

    // Teardown
    cleanupTestEnvironment()

    os.Exit(code)
}
```

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

**11. Go Module Management**
- Use `go mod tidy` to clean dependencies
- Verify with `go mod verify`
- Check for vulnerabilities with `just lint-vulns`

**12. Tool Versions**
- Go 1.25.6 specified in go.mod
- golangci-lint v2.6.0
- go-arch-lint v1.14.0
- Use `just install` to install tools

### Linting and Quality

**13. Pre-commit Hooks**
- Install with `just install-hooks` (fast) or `just install-hooks-full` (comprehensive)
- Fast version: formatting only
- Full version: includes architecture graph validation

**14. Linting is Strict**
- All linters in `.golangci.yml` are enabled
- Fix issues before committing
- Use `just fix` for automatic fixes

**15. Unused Code Detection**
- `unusedparams` linter catches unused parameters
- `unusedfunc` linter catches unused functions
- Remove unused code or use `_` prefix

### Git Workflow

**16. Never Use `git reset --hard`**
- Use `git mv` for file moves (preserves history)
- Never use plain `mv` in git repos
- Check with memory files before destructive operations

**17. Commit Message Format**
- Follow Conventional Commits spec
- Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
- Example: `feat: add Docker multi-stage build support`

---

## Development Workflow

### Setup New Feature

1. **Create feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Implement feature**
   - Add domain types in `internal/domain/`
   - Implement business logic in domain layer
   - Create workflow/jobs in `cmd/goreleaser-wizard/`
   - Add tests alongside implementation

3. **Run tests and linting**
   ```bash
   just test
   just lint
   ```

4. **Commit with proper message**
   ```bash
   git add cmd/goreleaser-wizard/my_feature.go internal/domain/types.go
   git commit -m "feat: add my feature"
   ```

5. **Run full CI**
   ```bash
   just ci
   ```

### Debugging

**Enable debug logging**:
```bash
# Global flag
goreleaser-wizard --debug init

# Environment variable
GORELEASER_WIZARD_DEBUG=true goreleaser-wizard init
```

**Check architecture violations**:
```bash
just lint-arch  # Detailed architecture output
just graph      # Visualize dependencies
```

**Profile performance**:
```bash
just profile-cpu      # CPU profiling
just profile-heap     # Memory profiling
just profile-trace    # Execution trace
```

### Common Issues

**Import cycle detected**:
- Check architecture with `just lint-arch -v`
- Use dependency injection to break cycle
- Move shared code to domain layer

**Template rendering fails**:
- Check template escaping with `internal/validation/template_escaping.go`
- Verify template functions are registered
- Ensure template data matches expected structure

**Validation errors**:
- Check error code in `internal/errors/domain_errors.go`
- Verify validation logic in `internal/validation/`
- Check business rules in `internal/validation/business_rules.go`

**Tests failing in CI**:
- Run locally with `just ci`
- Check test coverage with `just coverage`
- Verify architecture with `just lint-arch`

---

## Project-Specific Context

### GoReleaser Configuration Generation

The wizard generates four main configuration files:

1. **`.goreleaser.yaml`** - Main GoReleaser configuration
2. **`.github/workflows/release.yml`** - GitHub Actions workflow
3. **`Dockerfile`** - Multi-stage Docker build
4. **`Homebrew formula`** - Package manager formula

### Supported Project Types

- **CLI Application** - Single binary with version info
- **Web Service** - Includes Docker configuration
- **Library with CLI** - Focuses on CLI component
- **Multiple Binaries** - Configures multiple build targets

### Smart Defaults

The wizard analyzes the project to provide intelligent defaults:
- Detects project structure from `go.mod`
- Recommends platforms (Linux, macOS, Windows)
- Suggests architectures (amd64, arm64)
- Identifies Git provider (GitHub, GitLab, etc.)
- Detects Docker registry preference

### Validation Workflow

1. **Pre-generation validation** - Check configuration is valid
2. **Post-generation validation** - Verify generated files
3. **GoReleaser check** - Run `goreleaser check`
4. **File system validation** - Ensure files are writable

---

## Configuration Files

### `.go-arch-lint.yml`
- Defines architectural components and rules
- Enforces Clean Architecture boundaries
- Enables deep scanning for AST analysis
- Excludes test files from validation

### `.golangci.yml`
- Comprehensive linter configuration
- 100+ linters enabled
- Strict error/warning thresholds
- Custom rules for code quality

### `justfile`
- Primary task runner
- Core commands: build, test, fmt, clean, ci
- Integration with go-arch-lint commands

### `dev/arch-lint.just`
- Enterprise-grade linting commands
- Advanced features: profiling, benchmarks, reporting
- Security scanning: capslock, gosec, govulncheck
- Architecture graph generation

---

## Key Dependencies

### Production Dependencies

- **github.com/spf13/cobra** - CLI framework
- **github.com/spf13/viper** - Configuration management
- **github.com/charmbracelet/lipgloss** - Terminal styling
- **github.com/charmbracelet/log** - Logging
- **gopkg.in/yaml.v3** - YAML parsing

### Development Dependencies

- **github.com/golangci/golangci-lint** - Linting
- **github.com/fe3dback/go-arch-lint** - Architecture enforcement
- **github.com/google/capslock** - Capability analysis
- **github.com/stretchr/testify** - Testing assertions

---

## Documentation Resources

### Project Documentation
- **README.md** - Project overview and usage
- **CONTRIBUTING.md** - Contribution guidelines
- **CHANGELOG.md** - Version history
- **SECURITY.md** - Security policy

### Internal Documentation
- **docs/status/** - Status reports and analysis
- **docs/planning/** - Architecture and planning documents
- **docs/github/** - GitHub issue analysis

### External Resources
- **GoReleaser Documentation** - https://goreleaser.com
- **Cobra Documentation** - https://github.com/spf13/cobra
- **Charm Libraries** - https://charm.sh

---

## Performance Considerations

### Optimization Guidelines

1. **Template Compilation**
   - Compile templates once, reuse multiple times
   - Use `template.Must` for panic on error
   - Register custom functions in `template.FuncMap`

2. **File I/O**
   - Use buffered I/O for large files
   - Batch file operations when possible
   - Clean up temporary files

3. **Concurrency**
   - Use goroutines for independent jobs
   - Limit concurrency with worker pools
   - Use context for cancellation

4. **Memory**
   - Reuse buffers and byte slices
   - Avoid unnecessary allocations in hot paths
   - Use `sync.Pool` for frequently allocated objects

---

## Security Considerations

### Security Best Practices

1. **Input Validation**
   - Validate all user inputs
   - Use regex patterns for format validation
   - Sanitize template inputs

2. **File Operations**
   - Check file permissions before writing
   - Use explicit permissions (0o644, 0o755)
   - Validate file paths (prevent directory traversal)

3. **External Commands**
   - Validate command arguments
   - Use context for timeout
   - Capture and validate output

4. **Secret Management**
   - Never log sensitive information
   - Use environment variables for secrets
   - Check for secrets before committing

### Security Tools

```bash
# Run security audit
just security-audit

# Scan for vulnerabilities
just lint-vulns

# Check for privilege escalation
just lint-capslock
```

---

## Troubleshooting

### Common Errors

**"go-arch-lint not found"**
- Run `just install` to install tools
- Or install manually: `go install github.com/fe3dback/go-arch-lint@v1.14.0`

**"architecture violations found"**
- Check `just lint-arch -v` for details
- Verify imports follow dependency rules
- Use `just graph` to visualize dependencies

**"template parsing failed"**
- Check template escaping (`{{` vs `{{"{{"}}`)
- Verify template syntax
- Ensure all variables are defined

**"validation failed"**
- Check error code and message
- Review validation logic in `internal/validation/`
- Verify business rules in `internal/validation/business_rules.go`

### Getting Help

- Check existing documentation: `just --list`
- Search issues: https://github.com/LarsArtmann/GoReleaser-Wizard/issues
- Review diagnostic messages from linters
- Enable debug logging: `--debug` flag

---

## Summary Checklist

Before committing changes:
- [ ] All tests pass (`just test`)
- [ ] Code formatted (`just fmt`)
- [ ] Linting passes (`just lint`)
- [ ] Architecture validated (`just lint-arch`)
- [ ] Security scan passes (`just security-audit`)
- [ ] New features have tests
- [ ] Error handling uses domain errors
- [ ] Documentation updated (if needed)
- [ ] Commit message follows Conventional Commits

Before pushing changes:
- [ ] Full CI passes (`just ci`)
- [ ] No architecture violations
- [ ] Code coverage >= 80%
- [ ] No security vulnerabilities
- [ ] GoReleaser configuration validates

---

**Remember**: This project enforces strict architecture and code quality standards. Always run `just ci` before committing to ensure all checks pass.
