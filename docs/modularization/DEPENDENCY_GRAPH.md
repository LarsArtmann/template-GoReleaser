# Dependency Graph — GoReleaser-Wizard

**Generated:** 2026-05-13
**Module:** `github.com/LarsArtmann/GoReleaser-Wizard`

---

## Current State: Monolith with Partial Split

- **Root go.mod:** `github.com/LarsArtmann/GoReleaser-Wizard` (Go 1.26.2)
- **Secondary go.mod:** `test-wizard/` (standalone test binary, Go 1.26.0, no internal deps)
- **go.work:** Not present

---

## Internal Package Dependency Graph

```
                    ┌─────────────────────────────────┐
                    │         cmd/goreleaser-wizard     │
                    │            (package main)         │
                    └──────┬──────────────┬─────────────┘
                           │              │
              ┌────────────┘              └──────────────┐
              │                                          │
    ┌─────────▼──────────┐                   ┌───────────▼──────────┐
    │  internal/config    │                   │    internal/domain    │
    │  (koanf, pflag)     │                   │  (pure, zero deps)   │
    └─────────────────────┘                   └──┬──────┬──────┬────┘
                                                  │      │      │
                           ┌──────────────────────┘      │      └─────────────────┐
                           │                             │                        │
              ┌────────────▼──────────┐    ┌────────────▼─────┐    ┌──────────────▼──────────┐
              │    internal/errors     │    │   internal/git    │    │  internal/types          │
              │  (DomainError, codes)  │    │  (git commands)   │    │  (ValidationResult)      │
              └────────────────────────┘    └───────────────────┘    └─────────────────────────┘
                           ▲                             ▲
              ┌────────────┴──────┐           ┌─────────┘
              │                   │           │
    ┌─────────▼──────────┐  ┌────▼────────────▼──────┐
    │ internal/validation │  │    internal/utils       │
    │ (basic, business)   │  │  (recommendations)     │
    └─────────────────────┘  └────────────────────────┘

    Sub-packages under cmd/goreleaser-wizard:

    main ──→ generators ──→ templates (embedded template strings)
                    │
                    └──→ types (template data structs)
```

---

## Internal Dependency Matrix

| Package | Depends On |
|---|---|
| `internal/domain` | *(none — pure)* |
| `internal/errors` | *(none — stdlib only)* |
| `internal/config` | *(none — koanf/pflag only)* |
| `internal/git` | `internal/errors` |
| `internal/types` | `internal/errors` |
| `internal/utils` | `internal/domain`, `internal/git` |
| `internal/validation` | `internal/domain`, `internal/errors`, `internal/types` |
| `cmd/.../generators` | `cmd/.../templates`, `cmd/.../types`, `internal/domain`, `internal/errors`, `internal/git` |
| `cmd/.../types` | `internal/domain`, `internal/git` |
| `cmd/.../templates` | *(none — string constants)* |
| `cmd/...` (main) | `internal/config`, `internal/domain`, `internal/git` |

---

## External Dependency Usage by Package

| Package | External Dependencies |
|---|---|
| `internal/domain` | `github.com/go-faster/yaml`, `github.com/larsartmann/go-branded-id` |
| `internal/config` | `github.com/knadh/koanf/v2`, `github.com/spf13/pflag` |
| `internal/errors` | *(none — stdlib only)* |
| `internal/git` | *(none — stdlib `os/exec` only)* |
| `internal/types` | *(none — stdlib only)* |
| `internal/utils` | *(none — stdlib only)* |
| `internal/validation` | *(none — stdlib only)* |
| `cmd/...` (main) | `charm.land/lipgloss/v2`, `charm.land/log/v2`, `github.com/charmbracelet/fang`, `github.com/spf13/cobra` |
| `cmd/.../generators` | *(none beyond internal)* |
| `cmd/.../types` | *(none beyond internal)* |
| `cmd/.../templates` | *(none)* |

---

## Coupling Hotspots

### 1. Split Brain: Two DomainError Systems

Two packages define nearly identical error infrastructure:

| Aspect | `internal/errors/domain_errors.go` | `internal/domain/errors.go` |
|---|---|---|
| `ErrorCode` type | ✅ | ✅ |
| `DomainError` struct | ✅ (with Level, Retryable, File, Line) | ✅ (with Severity, RecoverySuggestion) |
| `WithContext()` | Mutates receiver | Returns new copy |
| Error codes | ~40 codes (infra-focused) | ~35 codes (domain-focused) |
| Used by | generators, git, validation, types | main, init, validate |

**Impact:** Callers must know which `DomainError` they're dealing with. Type aliases in `cmd/...` (`type DomainError = domain.DomainError`) bridge the gap, but the split creates confusion.

### 2. Split Brain: Two ValidationResult Types

| `internal/types/validation.go` | `internal/domain/interfaces.go` |
|---|---|
| `ValidationResult` with `[]*ValidationError`, `[]*ValidationWarning`, `ValidationSummary` | `ValidationResult` with `[]*DomainError`, `[]*DomainError`, `ValidationRules` |
| Rich grading/score system | Simple error list |
| Used by `internal/validation` | Used by `internal/domain` |

### 3. Duplicate Validation Logic

Same validation functions exist in two places with **different rules**:

| Function | `internal/domain/validators.go` | `internal/validation/basic.go` |
|---|---|---|
| `ValidateProjectName` | Max 63 chars, regex `^[a-zA-Z0-9][a-zA-Z0-9\-._]*[a-zA-Z0-9]$` | Max 50 chars, regex `^[a-zA-Z0-9][a-zA-Z0-9._-]{0,49}$` |
| `ValidateBinaryName` | Max 255 chars | Max 30 chars |
| `ValidateConfiguration` | In `domain/validation.go` | In `validation/business_rules.go` |

### 4. `internal/utils` Depends on Two Packages

`internal/utils/recommendations.go` imports both `internal/domain` and `internal/git`, coupling recommendation logic to git infrastructure.

### 5. `cmd/goreleaser-wizard/generators` Has Wide Import Surface

Generators import: templates, types, domain, errors, git — essentially the full stack.

---

## Proposed Module Boundaries

```
Module: core    ← internal/domain + consolidated errors + consolidated types
Module: cli     ← cmd/goreleaser-wizard (depends on core)
```

Rationale documented in PROPOSAL.md.
