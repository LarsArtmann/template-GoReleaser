# Modularization Proposal — GoReleaser-Wizard

**Date:** 2026-05-13
**Status:** Draft — Pending Self-Review
**Scope:** Full modularization with prerequisite de-duplication

---

## 1. Executive Summary

GoReleaser-Wizard is a ~14k LOC Go CLI tool in a single-module monorepo. The codebase has significant structural issues that must be resolved **before** modularization yields value: three split brains (error types, validation results, validation functions) create inconsistent behavior and make boundary detection unreliable.

**Proposal:** Consolidate split brains first, then split into 3 sub-modules with `go.work` for local development. The core domain becomes a reusable library; git operations and CLI each get isolated modules.

**Expected benefits:**

- Domain layer reusable without CLI dependencies
- Clear compile-time enforced boundaries
- Independent versioning possible for the domain library
- CI can test modules in parallel
- Elimination of 3 critical split brains

---

## 2. Current State Analysis

### Module Landscape

| Module | Path | Internal Deps | External Deps | State |
|---|---|---|---|---|
| `github.com/LarsArtmann/GoReleaser-Wizard` | `./` | All internal | 14 direct | **Monolith** |
| `test-wizard` | `./test-wizard/` | None | None | **Isolated** |

### Codebase Size

| Area | Files | Lines |
|---|---|---|
| `internal/domain/` | 28 | 5,754 |
| `cmd/goreleaser-wizard/` | 24 | 7,030 |
| `cmd/.../generators/` | 6 | 933 |
| `cmd/.../types/` | 1 | ~120 |
| `cmd/.../templates/` | 1 | ~200 |
| `internal/errors/` | 1 | ~350 |
| `internal/types/` | 3 | ~350 |
| `internal/validation/` | 6 | ~700 |
| `internal/git/` | 1 | ~200 |
| `internal/config/` | 1 | ~200 |
| `internal/utils/` | 1 | ~150 |
| **Total** | **~73** | **~16,000** |

### Coupling Hotspots (Critical)

#### Split Brain #1: Two DomainError Systems

| Package | ErrorCode | DomainError | Used By |
|---|---|---|---|
| `internal/errors` | 40+ codes, infra-focused | With Level, Retryable, Caller | generators, git, validation, types |
| `internal/domain` | 35 codes, domain-focused | With Severity, RecoverySuggestion | main, init, validate |

Both define `ErrorCode string`, `DomainError struct`, `Error()`, `Unwrap()`, `WithContext()`, `NewValidationError()`. Different semantics (`WithContext` mutates vs copies).

#### Split Brain #2: Two ValidationResult Types

| Package | Structure | Features |
|---|---|---|
| `internal/types` | `[]*ValidationError` + `[]*ValidationWarning` + `ValidationSummary` | Rich grading, scoring, filtering |
| `internal/domain` | `[]*DomainError` + `[]*DomainError` + `ValidationRules` | Simple error/warning lists |

#### Split Brain #3: Duplicate Validation Functions

Four functions duplicated across `internal/domain/validators.go` and `internal/validation/basic.go` with **different limits** (project name: 63 vs 50 chars, binary name: 255 vs 30). Two `ValidateConfiguration` implementations.

### Self-Review Corrections (Phase 4)

**Critical finding:** `internal/validation/` is **dead code** — zero non-test files import it. Only `init_test.go` imports it. This eliminates Split Brain #3 from the migration path (the duplicate validation functions in `internal/validation/basic.go` are simply unused and should be deleted, not reconciled).

**Additional dead code:**
- `pkg/errors/` — zero imports, deprecated, should be deleted
- `internal/domain/enums_release_temp.go` — empty file, should be deleted

**`generators.Logger` vs `domain.Logger`:** Not compatible — `generators.Logger` is a 4-method subset, `domain.Logger` has 13 methods. Generators should adopt `domain.Logger` and the `LoggerAdapter` in main should be the single implementation.

**`cmd/.../types/template_data.go`** depends on `internal/git`, so it stays in CLI module. `cmd/.../generators/` also references `internal/git` (in `goreleaser.go`), so it stays in CLI module too. Both correctly placed in the CLI module.

**`internal/utils/recommendations.go`** depends on both `domain` and `git`. The git-dependent functions (`GetGitHubOwner`, `GetGitHubRepo`) should move into `gitutil` module. The domain-only recommendation functions can stay in core.

**Revised module contents:**

Module `core`:
- `domain/` — all domain types (unchanged)
- `types/` — ValidationResult types (consolidated)
- `utils/` — recommendation functions (domain-only, after git functions removed)

Module `gitutil`:
- `git/` — git commands
- `utils/` git-dependent functions → merged into gitutil package

Module `cli` (root):
- `cmd/goreleaser-wizard/` — CLI, generators, templates, types, workflow, jobs
- `config/` — koanf-based config

### God-Package Analysis

| Package | Files | Lines | Concerns |
|---|---|---|---|
| `internal/domain` | 28 | 5,754 | Core config, enums (7 types), errors, IDs, events, interfaces, validation |
| `cmd/goreleaser-wizard` (main) | 16 | 7,030 | CLI commands, workflow engine, job manager, validation display, TUI |

The domain package is large but cohesive — all files belong to the same domain concept (GoReleaser project configuration). The 7 enum files are naturally separate concerns within the domain. No further package-level splits recommended at this time.

The cmd package is the real god-package: CLI commands, workflow orchestration, job management, validation display, TUI wizard, and template generation all share the `main` package. However, sub-packages (`generators/`, `types/`, `templates/`) already provide internal separation, and further splitting would not improve the module boundary since all of this belongs to the CLI module.

---

## 3. Proposed Module Structure

### Prerequisites (Split Brain Resolution)

Before any module split, the following consolidations must happen:

1. **Merge error systems** → Keep `internal/domain/errors.go` as canonical, expand with best features from `internal/errors/domain_errors.go` (caller info, retryable, suggestion), delete `internal/errors/`
2. **Merge ValidationResult** → Keep `internal/types/validation.go` as canonical (richer), remove `ValidationResult` from `internal/domain/interfaces.go`
3. **Merge validation functions** → Keep `internal/domain/validators.go` as canonical (closer to domain), reconcile limits, delete duplicates from `internal/validation/basic.go`, keep `internal/validation/` for business rules and form validation only

### Proposed Modules

```
GoReleaser-Wizard/
├── go.work
├── core/                          ← Module 1: github.com/LarsArtmann/GoReleaser-Wizard/core
│   ├── go.mod
│   └── (contents of internal/ after consolidation)
│       ├── domain/                ← Core domain types, enums, config, validation
│       ├── types/                 ← ValidationResult, validation types
│       └── utils/                 ← Recommendation functions (domain-only, no git deps)
│
├── gitutil/                       ← Module 2: github.com/LarsArtmann/GoReleaser-Wizard/gitutil
│   ├── go.mod                     ← depends on: core
│   └── (contents of internal/git/ + git-dependent utils functions)
│
├── cmd/goreleaser-wizard/         ← Module 3: github.com/LarsArtmann/GoReleaser-Wizard (root)
│   ├── go.mod                     ← depends on: core, gitutil
│   ├── main.go
│   ├── generators/
│   ├── templates/
│   ├── types/
│   └── ...
│
└── test-wizard/                   ← Unchanged (isolated test binary)
    └── go.mod
```

### Module Definitions

#### Module 1: `core`

| Field | Content |
|---|---|
| **Path** | `./core/` |
| **Module** | `github.com/LarsArtmann/GoReleaser-Wizard/core` |
| **Purpose** | Pure domain types, validation, and configuration for GoReleaser projects |
| **Dependencies (prod)** | None (zero internal) |
| **Dependencies (test)** | None |
| **Public API** | `SafeProjectConfig`, all enum types, `DomainError`, `ValidationResult`, validators, `Logger` interface, `FileSystemRepository` interface |
| **External deps** | `github.com/go-faster/yaml`, `github.com/larsartmann/go-branded-id`, `github.com/stretchr/testify` (test) |
| **Packages** | `domain/`, `types/`, `utils/` (domain-only recommendations) |

#### Module 2: `gitutil`

| Field | Content |
|---|---|
| **Path** | `./gitutil/` |
| **Module** | `github.com/LarsArtmann/GoReleaser-Wizard/gitutil` |
| **Purpose** | Git command wrappers for repository introspection |
| **Dependencies (prod)** | `core` (for `DomainError`) |
| **Dependencies (test)** | None |
| **Public API** | `Command`, `RepositoryInfo`, `VersionInfo`, `IncPatchVersion`, `GetGitHubOwner`, etc. |
| **External deps** | None (stdlib `os/exec` only) |
| **Packages** | Single `gitutil` package |

#### Module 3: `cli` (root module)

| Field | Content |
|---|---|
| **Path** | `./` (root go.mod) |
| **Module** | `github.com/LarsArtmann/GoReleaser-Wizard` |
| **Purpose** | CLI application — interactive wizard, generators, workflow orchestration |
| **Dependencies (prod)** | `core`, `gitutil` |
| **Dependencies (test)** | `github.com/stretchr/testify` |
| **Public API** | CLI binary only (no public Go API) |
| **External deps** | `charm.land/*`, `github.com/spf13/cobra`, `github.com/knadh/koanf/v2` |
| **Packages** | `cmd/goreleaser-wizard/`, `internal/config/`, `internal/utils/`, `cmd/.../generators/`, `cmd/.../templates/`, `cmd/.../types/` |

### DAG Verification

```
core (zero deps) ← gitutil ← cli (root)
                  ← cli (root)
```

- `core` → no internal dependencies ✅
- `gitutil` → depends on `core` only ✅
- `cli` → depends on `core` and `gitutil` ✅
- No cycles ✅
- No upward dependencies from core ✅

---

## 4. Replace / Workspace Strategy

**Chosen: `go.work` at repo root.**

| Aspect | Decision |
|---|---|
| File | `go.work` at repo root |
| Entries | `./core`, `./gitutil`, `./cmd/goreleaser-wizard` (or `.` if root stays) |
| Replace directives | None — clean go.mod files |
| Consumer experience | `go.work` ignored by Go proxy; consumers import versioned modules |
| CI | `go work sync` before build; each module testable independently |

Rules:
- No `replace` directives in any `go.mod`
- `go.work` committed to repo (since all modules are in-repo)
- Each module's `go.mod` must be independently valid (no workspace required for published consumption)

---

## 5. Test Dependency Isolation

| Module | Production Deps | Test-Only Deps |
|---|---|---|
| `core` | `go-faster/yaml`, `go-branded-id` | `stretchr/testify` |
| `gitutil` | `core` | `stretchr/testify` (if tests added) |
| `cli` | `core`, `gitutil`, charm libs, cobra, koanf | `stretchr/testify` |

`stretchr/testify` appears in root `go.mod` as a direct require but is only used in `*_test.go` files. This is acceptable for a monorepo but should be excluded from production dependency scanning.

---

## 6. Interface Extraction Plan

The current codebase already uses interfaces well:

- `domain.Logger` — implemented by `LoggerAdapter` in main
- `domain.FileSystemRepository` — implemented by `SimpleFileSystemRepository` in main
- `generators.Logger` — duplicated interface, should reference `domain.Logger`

**Changes needed:**

1. Remove `generators.Logger` — use `core/domain.Logger` instead
2. Keep `domain.FileSystemRepository` and `domain.TemplateRepository` as-is
3. `gitutil` exposes concrete types (no interface needed — it wraps `os/exec`)

---

## 7. Versioning Strategy

**Chosen: Root-only versioning.**

| Strategy | Details |
|---|---|
| Tags | Single `v1.2.3` tags at repo root |
| Sub-modules | No independent tags initially |
| Rationale | Single team, single consumer (CLI binary), domain library not yet published externally |
| Future | If `core` is published independently, switch to per-module semver tags (`core/v1.2.3`) |

---

## 8. Migration Strategy

### Phase 0: Dead Code Removal + Split Brain Resolution (Prerequisites)

1. Delete `pkg/errors/` — unused, deprecated
2. Delete `internal/domain/enums_release_temp.go` — empty file
3. Delete `internal/validation/` — dead code (zero non-test imports; test references updated)
4. Merge `internal/errors/` into `internal/domain/errors.go` (consolidate ErrorCode, DomainError, add WithCaller/WithSuggestion from errors package)
5. Consolidate `ValidationResult`: keep `internal/types/validation.go` as canonical, remove `ValidationResult` from `internal/domain/interfaces.go`
6. Remove `generators.Logger` duplicate interface — adopt `domain.Logger` instead
7. Verify all tests pass

### Phase 1: Create `core` Module

6. Create `./core/` directory structure
7. Move `internal/domain/` → `core/domain/`
8. Move `internal/types/` → `core/types/`
9. Move git-only functions from `internal/utils/recommendations.go` → `gitutil/`
10. Move remaining `internal/utils/` → `core/utils/`
11. Create `core/go.mod`
12. Update all import paths (`internal/domain` → `core/domain`, etc.)
13. Verify build and tests

### Phase 2: Create `gitutil` Module

14. Create `./gitutil/` directory
15. Move `internal/git/` contents → `gitutil/`
16. Add git-dependent recommendation functions from old `internal/utils/`
17. Create `gitutil/go.mod` with `core` dependency
18. Update import paths
19. Verify build and tests

### Phase 3: Workspace Setup

20. Create `go.work` at repo root
21. Move `internal/config/` → root `internal/config/` (stays in CLI module)
22. Update root `go.mod` to depend on `core` and `gitutil`
23. Run `go work sync`
24. Full CI verification

### Phase 4: Cleanup

25. Remove empty `internal/` directories
26. Update `.go-arch-lint.yml` for new module paths
27. Update `.golangci.yml` for new paths
28. Update documentation (AGENTS.md, README.md)
29. Verify `nix build` / `just ci`

---

## 9. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Error system consolidation breaks behavior | Medium | High | Comprehensive test suite before and after; type aliases during transition |
| Import path migration misses references | Low | Medium | `go vet ./...` + `grep` verification at each step |
| `go-faster/yaml` in domain limits core reusability | Low | Low | yaml is only used for `ToYAML()`/`FromYAML()` — could be extracted later |
| CI breaks with go.work | Low | Medium | Test both workspace and non-workspace builds |
| `internal/config` (koanf) pulls heavy deps into CLI | None | None | Stays in CLI module where it belongs |
| Removing dead `internal/validation/` breaks test in `init_test.go` | High | Low | Update or remove the affected test; the test imports dead code |
| `generators.Logger` → `domain.Logger` migration breaks generators | Medium | Medium | Expand `LoggerAdapter` to satisfy `domain.Logger` (it already does) |

---

## 10. Build System Impact

| System | Changes Needed |
|---|---|
| `flake.nix` | Add per-module build targets, aggregate at root |
| `justfile` | Update build/test/lint commands for multi-module |
| `.go-arch-lint.yml` | Update component paths to reflect new module structure |
| `.golangci.yml` | Update path-specific rules |
| `.github/workflows/release.yml` | Add per-module test jobs, parallel execution |
| `AGENTS.md` | Update build commands, module structure documentation |

---

## Key Decisions

1. **Consolidate first, split second** — Split brains must be resolved before module boundaries are meaningful
2. **Three modules, not more** — The codebase is ~16k LOC; finer granularity adds overhead without benefit
3. **`go.work` over `replace`** — Cleaner than per-module replace directives
4. **Root-only versioning** — Single team, single consumer; independent versioning is premature
5. **`gitutil` as separate module** — Git operations are infrastructure, not domain; isolating them prevents `os/exec` from leaking into core
6. **`internal/validation/` is dead code** — Removed entirely, not migrated
7. **`generators.Logger` adopts `domain.Logger`** — Single logger interface across all modules
8. **`internal/config/` stays in CLI module** — koanf/pflag are CLI-layer concerns

---

## Self-Review (Phase 4)

### 12 Questions Answered

1. **What did you forget?** Initially proposed moving `internal/validation/` into core. Self-review revealed it's dead code (zero non-test imports). Corrected: delete, don't migrate.
2. **What could you have done better?** Should have checked import usage before proposing module contents. Corrected by checking all import sites.
3. **What could you still improve?** `internal/utils/` has mixed dependencies (domain + git). Splitting the git-dependent functions out would make core cleaner.
4. **Did you create split brains?** The proposal eliminates 3 existing split brains (errors, validation results, validation functions) and introduces no new ones.
5. **Are boundaries at the right granularity?** Yes — 3 modules for 16k LOC is appropriate. Further splitting the domain package would create overhead without benefit.
6. **Can existing code be reused?** `internal/config/` stays as-is (only 3 importers, all CLI). `generators/` stay in CLI (git dependency prevents core placement).
7. **Can type models be improved?** `generators.Logger` should use `domain.Logger` — this is already in the plan.
8. **Are you leveraging established libraries?** Yes — go-branded-id, go-faster/yaml, koanf, cobra, charm libs are all appropriate choices.
9. **Does the workspace strategy work?** `go.work` with 3 entries (core, gitutil, root) — verified that import paths resolve correctly in the proposed structure.
10. **Are test-only deps isolated?** `stretchr/testify` is in root go.mod direct requires but only used in tests. This is standard Go practice and acceptable.
11. **Will CI be faster?** Marginally — `core` can be tested independently. The real speedup comes from parallel module builds in CI.
12. **Is versioning realistic?** Root-only versioning is correct for this project — the domain library has no external consumers yet.
