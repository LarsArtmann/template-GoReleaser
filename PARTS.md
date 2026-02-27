# GoReleaser-Wizard: Component Analysis & Extraction Strategy

**Version:** 1.0 | **Updated:** February 26, 2026

---

## Executive Summary

This document identifies reusable components within GoReleaser-Wizard that could be extracted as standalone libraries/SDKs. Each component is evaluated against:

1. **Reusability** — Can other projects benefit?
2. **Independence** — Can it stand alone without project-specific logic?
3. **Value-add** — Does our abstraction provide something existing alternatives don't?
4. **Maintenance burden** — Is the extraction worth the overhead?

---

## Component Inventory

| Component           | Location                                              | Lines | Reusability | Independence | Extraction Priority           |
| ------------------- | ----------------------------------------------------- | ----- | ----------- | ------------ | ----------------------------- |
| Domain Errors       | `internal/errors/`                                    | ~327  | High        | High         | **1 - Extract Now**           |
| Typed Enums         | `internal/domain/enums*.go`                           | ~800+ | High        | High         | **1 - Extract Now**           |
| Job/Workflow System | `cmd/goreleaser-wizard/job_manager.go`, `workflow.go` | ~700  | Medium      | Medium       | **2 - Extract with Refactor** |
| Form Validator      | `internal/validation/form_validator.go`               | ~175  | Medium      | High         | **2 - Extract with Refactor** |
| SafeProjectConfig   | `internal/domain/config_core.go`                      | ~420  | Low\*       | High         | **3 - Domain-Specific**       |
| Template Generator  | `cmd/goreleaser-wizard/generators/`                   | ~400  | Medium      | Medium       | **3 - Context-Specific**      |

\*Low reusability as-is, but the pattern is highly reusable.

---

## Component 1: Structured Domain Errors

### Current State

```
internal/errors/domain_errors.go (~327 lines)
```

**Features:**

- Typed error codes (40+ predefined)
- Error levels (critical, high, medium, low)
- Rich context (field, details, cause, retryable)
- Builder pattern for error construction
- Caller information capture
- JSON serialization

### Alternatives

| Library                   | Stars | Our Advantage                                    |
| ------------------------- | ----- | ------------------------------------------------ |
| `cockroachdb/errors`      | 2.5k  | Our error codes + levels + retryable flag        |
| `larsartmann/uniflow`     | —     | Complementary (pipeline errors vs domain errors) |
| `pkg/errors` (deprecated) | 8.9k  | Our structured codes + modern patterns           |
| `emperror/errors`         | 1.5k  | Our simpler API + retryable semantics            |

### Value Proposition

Our `DomainError` provides:

1. **Typed Error Codes** — Enumerated codes for programmatic handling
2. **Severity Levels** — Built-in categorization for alerting/filtering
3. **Retryable Flag** — Explicit retry semantics (not in cockroachdb/errors)
4. **JSON Serialization** — Structured logging/API responses out of the box
5. **Zero Dependencies** — Pure Go, no external imports

### Recommendation: **EXTRACT → `github.com/LarsArtmann/go-domain-errors`**

```go
// Proposed API
import "github.com/LarsArtmann/go-domain-errors"

type ErrorCode string
type ErrorLevel string

type DomainError struct {
    Code      ErrorCode
    Message   string
    Details   string
    Level     ErrorLevel
    Retryable bool
    // ...builder methods
}

// Factory functions
func NewValidationError(code ErrorCode, message, details string) *DomainError
func NewSystemError(code ErrorCode, message, details string) *DomainError
// ...
```

**Migration Path:**

1. Extract to standalone repo
2. Add generic error codes (validation, system, network, etc.)
3. Keep GoReleaser-specific codes in the main project
4. Update imports

---

## Component 2: Type-Safe Enums

### Current State

```
internal/domain/enums*.go (~800+ lines total)
```

**Files:**

- `enums.go` — Core enum types (CGOStatus, SigningLevel)
- `enums_platform.go` — Platform, Architecture
- `enums_build.go` — Build-related enums
- `enums_project.go` — ProjectType, FeatureLevel
- `enums_actions.go` — ActionLevel, ActionTrigger

**Features:**

- String-based with compile-time safety
- Validation methods (`IsValid()`)
- Display methods (`String()`)
- Semantic methods (`IsEnabled()`, `IsRequired()`)
- Cross-type compatibility methods
- Smart conversion utilities

### Alternatives

| Library                 | Stars | Our Advantage                           |
| ----------------------- | ----- | --------------------------------------- |
| `alvaroloes/enumer`     | 800+  | Our semantic methods + compatibility    |
| `abice/go-enum`         | 400+  | Our zero-code approach + smart defaults |
| `json-iterator/go-enum` | —     | Our validation + display methods        |
| Go 1.22+ `go/types`     | —     | Our domain-specific behaviors           |

### Value Proposition

Our enum pattern provides:

1. **Semantic Methods** — `IsEnabled()`, `IsRequired()`, `ToBool()` for gradual typing
2. **Compatibility Checking** — `IsCompatibleWith()` for cross-enum validation
3. **Smart Defaults** — Context-aware recommendations
4. **Zero Generation** — No code generation step required
5. **Serialization** — Built-in JSON/YAML support

### Recommendation: **EXTRACT → `github.com/LarsArtmann/go-enumx`**

```go
// Proposed API
import "github.com/LarsArtmann/go-enumx"

// Interface for all enums
type Enum interface {
    IsValid() bool
    String() string
}

// Mixin for tri-state enums (disabled/enabled/required)
type TriState string

const (
    TriStateDisabled TriState = "disabled"
    TriStateEnabled  TriState = "enabled"
    TriStateRequired TriState = "required"
)

func (ts TriState) IsEnabled() bool
func (ts TriState) IsRequired() bool
func (ts TriState) ToBool() bool

// Level-based enum pattern
type Level string

func (l Level) IsEnabled() bool  // l != "none"
func (l Level) IsAtLeast(min Level) bool
```

**Key Innovation:** A library that provides the _pattern_ rather than just code generation. Users get:

- Generic tri-state enum (`disabled/enabled/required`)
- Generic level enum (`none/basic/advanced/enterprise`)
- Generic multi-select enum
- Compatibility matrix helpers

---

## Component 3: Job/Workflow Execution System

### Current State

```
cmd/goreleaser-wizard/
  job_manager.go    (~337 lines)
  workflow.go       (~429 lines)
  jobs.go           (~varies)
```

**Features:**

- Job interface with Execute/Rollback
- Sequential and parallel execution
- Context-aware cancellation
- Job result tracking
- Automatic rollback on failure
- Statistics and timing

### Alternatives

| Library                  | Stars | Our Advantage                     |
| ------------------------ | ----- | --------------------------------- |
| `hibiken/asynq`          | 9k    | Our simpler, synchronous approach |
| `machinery-v2/machinery` | 7.5k  | Our rollback support + zero infra |
| `go-co-op/gocron`        | 5k    | Our workflow composition          |
| `dagu-go/dagu`           | 1.5k  | Our embedded nature               |
| `temporalio/sdk-go`      | 5k    | Our simplicity (no server)        |

### Value Proposition

Our system provides:

1. **Rollback Support** — Automatic rollback on failure (rare in alternatives)
2. **Zero Infrastructure** — Embedded, no external services
3. **Context-Aware** — Proper cancellation propagation
4. **Type-Safe Results** — Structured job results
5. **Builder Pattern** — Fluent workflow construction

### Recommendation: **EXTRACT → `github.com/LarsArtmann/go-workflow`**

```go
// Proposed API
import "github.com/LarsArtmann/go-workflow"

type Job interface {
    ID() string
    Name() string
    Execute(ctx context.Context) error
    Rollback(ctx context.Context) error
}

type Workflow struct { ... }

func NewWorkflow(name string, opts ...Option) *Workflow
func (w *Workflow) AddJob(job Job) *Workflow
func (w *Workflow) SetParallel(maxJobs int) *Workflow
func (w *Workflow) SetTimeout(d time.Duration) *Workflow
func (w *Workflow) Execute(ctx context.Context) error
func (w *Workflow) GetResults() []JobResult

// Builder pattern
workflow := workflow.New("deploy").
    AddJob(validateJob).
    AddJob(buildJob).
    AddJob(deployJob).
    SetParallel(3).
    SetTimeout(10 * time.Minute).
    OnRollback(func(ctx context.Context, failed Job) {
        log.Warn("rolling back", "job", failed.Name())
    })
```

**Key Differentiator:** Rollback-first design. Most workflow libs focus on success paths; ours handles failure gracefully with automatic rollback.

---

## Component 4: Form Validator (Huh-Compatible)

### Current State

```
internal/validation/form_validator.go (~175 lines)
```

**Features:**

- `huh` form integration
- Field-level error collection
- Generic validation helpers
- Pattern-based validation
- Input sanitization

### Alternatives

| Library                   | Stars | Our Advantage           |
| ------------------------- | ----- | ----------------------- |
| `go-playground/validator` | 16k   | Our huh integration     |
| `charmbracelet/huh`       | 4k    | Our reusable validators |
| `asaskevich/govalidator`  | 6k    | Our form-specific API   |

### Value Proposition

1. **Huh-Native** — Direct integration with Charmbracelet's form library
2. **Error Aggregation** — Collects all errors, not just first
3. **Reusable Patterns** — Common validators pre-built
4. **Security Helpers** — Shell metacharacter, path traversal detection

### Recommendation: **EXTRACT → `github.com/LarsArtmann/go-huh-validate`**

```go
// Proposed API
import "github.com/LarsArtmann/go-huh-validate"

type Validator struct { ... }

func New() *Validator

// Huh-compatible validators
func (v *Validator) Required(field string) func(string) error
func (v *Validator) Email(field string) func(string) error
func (v *Validator) URL(field string) func(string) error
func (v *Validator) NoShellMeta(field string) func(string) error
func (v *Validator) NoPathTraversal(field string) func(string) error
func (v *Validator) Match(field string, pattern *regexp.Regexp) func(string) error

// Access collected errors
func (v *Validator) Errors() map[string]string
func (v *Validator) HasErrors() bool
```

**Key Differentiator:** First-class `huh` integration with security-focused validators.

---

## Component 5: SafeProjectConfig Pattern

### Current State

```
internal/domain/config_core.go (~420 lines)
```

**Features:**

- Single source of truth for configuration
- Smart defaults based on project analysis
- Multi-level validation (field, config, business rules)
- Serialization (JSON, YAML)
- Clone/deep copy support

### Alternatives

| Library                     | Purpose        | Our Advantage                   |
| --------------------------- | -------------- | ------------------------------- |
| `knadh/koanf`               | Config loading | Our typed config + validation   |
| `spf13/viper`               | Config loading | Our immutable, validated config |
| `kelseyhightower/envconfig` | Env binding    | Our rich validation             |

### Value Proposition

The _pattern_ is reusable, not the specific config:

1. **Immutable After Validation** — Once validated, safe to use
2. **Smart Defaults** — Context-aware initialization
3. **Business Rule Validation** — Cross-field validation
4. **Serialization Agnostic** — JSON, YAML, or custom

### Recommendation: **EXTRACT PATTERN → `github.com/LarsArtmann/go-safecfg`**

```go
// Proposed API - A framework for building safe configs
import "github.com/LarsArtmann/go-safecfg"

// Users implement this
type SafeConfig interface {
    Validate() error
    ApplyDefaults()
}

// Framework provides
type Builder[T SafeConfig] struct { ... }

func NewBuilder[T SafeConfig](defaults T) *Builder[T]
func (b *Builder[T]) Validate() error
func (b *Builder[T]) Build() (T, error)
func (b *Builder[T]) FromYAML(data []byte) (*Builder[T], error)
func (b *Builder[T]) FromJSON(data []byte) (*Builder[T], error)
func (b *Builder[T]) FromEnv(prefix string) (*Builder[T], error)
```

**Key Innovation:** A generic framework for building validated, immutable configuration types with smart defaults.

---

## Component 6: Template Generator Pattern

### Current State

```
cmd/goreleaser-wizard/generators/
  goreleaser.go
  github_actions.go
  dockerfile.go
  homebrew.go
```

**Features:**

- Type-safe template data
- Git-aware version injection
- Backup before overwrite
- Preview mode (no file write)
- Template validation

### Alternatives

| Library            | Purpose        | Our Advantage              |
| ------------------ | -------------- | -------------------------- |
| `a-h/templ`        | Type-safe HTML | Our file-based approach    |
| `text/template`    | Stdlib         | Our type-safe data structs |
| `plopjs/plop` (JS) | Generator      | Our Go-native approach     |

### Value Proposition

1. **Type-Safe Data** — Strongly typed template data structs
2. **Git Integration** — Automatic version/commit injection
3. **Safe Operations** — Backup, preview, validation
4. **Composable** — Chain generators

### Recommendation: **KEEP INTERNAL** (or minimal extraction)

The pattern is good but too context-specific. Consider extracting only:

```go
// Potential extraction: go-templatekit
type Generator interface {
    Generate(ctx context.Context) error
    Preview(ctx context.Context) (string, error)
    Validate() error
}

type GitAwareGenerator struct {
    // Mixin for git version injection
}
```

---

## Extraction Roadmap

### Phase 1: Immediate (Low Risk, High Value)

| Library            | Effort | Impact | Dependencies |
| ------------------ | ------ | ------ | ------------ |
| `go-domain-errors` | 2 days | High   | None         |
| `go-enumx`         | 3 days | High   | None         |

### Phase 2: With Refactoring (Medium Risk, High Value)

| Library           | Effort | Impact | Dependencies     |
| ----------------- | ------ | ------ | ---------------- |
| `go-workflow`     | 5 days | High   | go-domain-errors |
| `go-huh-validate` | 2 days | Medium | None             |

### Phase 3: Pattern Extraction (Higher Effort)

| Library      | Effort | Impact | Dependencies               |
| ------------ | ------ | ------ | -------------------------- |
| `go-safecfg` | 5 days | Medium | go-domain-errors, go-enumx |

---

## Recommended Library Stack

After extraction, the dependency graph becomes:

```
go-domain-errors (zero deps)
    ↓
go-enumx (zero deps)
    ↓
go-safecfg → depends on: go-domain-errors, go-enumx
go-workflow → depends on: go-domain-errors
go-huh-validate (zero deps)
    ↓
goreleaser-wizard-cli → depends on: all above
```

---

## Alternative: Monorepo Approach

If maintaining multiple repos is too much overhead, consider a monorepo:

```
github.com/LarsArtmann/goreleaser-tools/
├── pkg/
│   ├── errors/       → go-domain-errors
│   ├── enumx/        → go-enumx
│   ├── workflow/     → go-workflow
│   ├── huhvalidate/  → go-huh-validate
│   └── safecfg/      → go-safecfg
├── cmd/
│   └── goreleaser-wizard/
└── go.work
```

Each `pkg/` can still be imported independently:

```go
import "github.com/LarsArtmann/goreleaser-tools/pkg/errors"
```

---

## Decision Matrix

| Component          | Extract?    | Why                                   |
| ------------------ | ----------- | ------------------------------------- |
| Domain Errors      | **YES**     | Zero deps, high value, clear API      |
| Typed Enums        | **YES**     | Zero deps, reusable pattern, generic  |
| Job/Workflow       | **YES**     | Rollback support is unique, embedded  |
| Form Validator     | **YES**     | Huh integration is niche but valuable |
| SafeProjectConfig  | **PATTERN** | Extract the framework, not the config |
| Template Generator | **NO**      | Too context-specific                  |

---

## Summary

**Extract 4 libraries:**

1. `go-domain-errors` — Structured errors with codes, levels, retryable
2. `go-enumx` — Type-safe enum patterns with semantic methods
3. `go-workflow` — Job execution with automatic rollback
4. `go-huh-validate` — Huh-compatible form validators

**Extract 1 pattern:**

5. `go-safecfg` — Generic framework for validated, immutable configs

**Keep internal:**

- Template generators (too context-specific)
- GoReleaser-specific types (domain logic)

This extraction strategy maximizes reuse while minimizing maintenance burden, following the principle: **extract what's generic, keep what's specific**.
