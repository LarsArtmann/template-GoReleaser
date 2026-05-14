# Execution Plan — GoReleaser-Wizard Modularization

**Date:** 2026-05-13
**Based on:** PROPOSAL.md v2 (post self-review)
**Total Steps:** 29 ordered tasks across 5 phases

---

## Pareto Impact Tiers

| Tier      | Impact        | Tasks             | Description                                                    |
| --------- | ------------- | ----------------- | -------------------------------------------------------------- |
| 1% → 51%  | Foundational  | Phase 0 (1-7)     | Eliminate split brains and dead code — enables everything else |
| 4% → 64%  | High leverage | Phase 1 (8-14)    | Create core module — enables independent domain testing        |
| 16% → 80% | Broad value   | Phase 2-3 (15-25) | gitutil module + go.work — enables parallel CI                 |
| Remaining | Polish        | Phase 4 (26-29)   | Cleanup and documentation                                      |

---

## Phase 0: Dead Code Removal + Split Brain Resolution

**Goal:** Eliminate all split brains and dead code. Project must build and test cleanly after each step.

### Task 1: Delete dead `pkg/errors/`

- **Action:** Remove `pkg/errors/errors.go` and `pkg/errors/` directory
- **Why:** Zero imports, deprecated
- **Verification:** `go build ./...` passes
- **Rollback:** `git revert HEAD`
- **Effort:** 5 min

### Task 2: Delete dead `internal/domain/enums_release_temp.go`

- **Action:** Remove empty file
- **Why:** Contains only `package domain` — no code
- **Verification:** `go build ./...` passes
- **Rollback:** `git revert HEAD`
- **Effort:** 2 min

### Task 3: Delete dead `internal/validation/`

- **Action:** Remove entire `internal/validation/` directory (6 files)
- **Why:** Zero non-test imports. Only `cmd/goreleaser-wizard/init_test.go` imports it
- **Pre-step:** Fix `init_test.go` to remove `internal/validation` import (the test likely uses `ValidateConfiguration` from domain instead)
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 15 min
- **Impact:** Eliminates Split Brain #3 (duplicate validation functions)

### Task 4: Consolidate error systems — expand `internal/domain/errors.go`

- **Action:** Add features from `internal/errors/domain_errors.go` into `internal/domain/errors.go`:
  - `WithCaller()` method (runtime caller info)
  - `WithSuggestion()` method
  - `WithRetryable()` method
  - `WithLevel()` method
  - `ErrorLevel` type (if not conflicting with existing `ErrorSeverity`)
  - Helper functions: `IsRetryable()`, `GetErrorCode()`, `GetErrorLevel()`
  - Missing error codes from `internal/errors/` that are used by generators/git
- **Why:** Two parallel `DomainError` systems with different capabilities
- **Verification:** `go build ./...` passes
- **Rollback:** `git revert HEAD`
- **Effort:** 30 min

### Task 5: Consolidate error systems — migrate `internal/errors/` consumers

- **Action:** Update all files importing `internal/errors` to import `internal/domain` instead:
  - `cmd/goreleaser-wizard/generators/goreleaser.go`
  - `cmd/goreleaser-wizard/generators/github_actions.go`
  - `cmd/goreleaser-wizard/generators/template_utils.go`
  - `cmd/goreleaser-wizard/generators/homebrew.go`
  - `internal/git/commands.go`
  - `internal/types/validation.go`
- **Action:** Delete `internal/errors/` directory
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 30 min
- **Impact:** Eliminates Split Brain #1 (two DomainError systems)

### Task 6: Consolidate `ValidationResult` — remove duplicate from domain

- **Action:** Remove `ValidationResult` type from `internal/domain/interfaces.go` (keep the richer version in `internal/types/validation.go`)
- **Action:** Update `internal/domain/validation.go`'s `ValidationUseCase.ValidateConfiguration()` to return `*types.ValidationResult` instead of domain's `ValidationResult`
- **Action:** Add `internal/types` import to `internal/domain/validation.go`
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 20 min
- **Impact:** Eliminates Split Brain #2 (two ValidationResult types)

### Task 7: Remove `generators.Logger` duplicate interface

- **Action:** Update `cmd/goreleaser-wizard/generators/` to use `domain.Logger` instead of local `Logger` interface
- **Action:** Update generators to use full `domain.Logger` signature (DebugContext, InfoContext, etc.)
- **Action:** Remove `Logger` interface from generators
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 20 min

---

## Phase 1: Create `core` Module

**Goal:** Extract pure domain logic into independently buildable `core` module.

### Task 8: Create `core/` directory structure

- **Action:** `mkdir -p core/domain core/types core/utils`
- **Verification:** Directory exists
- **Effort:** 2 min

### Task 9: Move `internal/domain/` → `core/domain/`

- **Action:** `git mv internal/domain/ core/domain/`
- **Verification:** Files moved
- **Effort:** 5 min

### Task 10: Move `internal/types/` → `core/types/`

- **Action:** `git mv internal/types/ core/types/`
- **Verification:** Files moved
- **Effort:** 5 min

### Task 11: Split `internal/utils/recommendations.go`

- **Action:** Extract git-dependent functions (`GetGitHubOwner`, `GetGitHubRepo`) into a temporary file
- **Action:** Keep domain-only functions (`GetRecommendedProjectType`, `GetRecommendedPlatforms`, `GetRecommendedArchitectures`, `GetRecommendedGitProvider`, `GetRecommendedDockerRegistry`, `IsDevelopmentEnvironment`, `IsProductionEnvironment`, `GetEnvironment`) in utils
- **Action:** `git mv internal/utils/ core/utils/`
- **Verification:** Files in place
- **Effort:** 15 min

### Task 12: Create `core/go.mod`

- **Action:** Create `core/go.mod`:

  ```
  module github.com/LarsArtmann/GoReleaser-Wizard/core

  go 1.26.2

  require (
      github.com/go-faster/yaml v0.4.6
      github.com/larsartmann/go-branded-id v0.1.0
  )
  ```

- **Action:** Run `cd core && go mod tidy`
- **Verification:** `cd core && go build ./...` passes
- **Effort:** 10 min

### Task 13: Update all import paths — `internal/` → `core/`

- **Action:** Find and replace across all files:
  - `github.com/LarsArtmann/GoReleaser-Wizard/internal/domain` → `github.com/LarsArtmann/GoReleaser-Wizard/core/domain`
  - `github.com/LarsArtmann/GoReleaser-Wizard/internal/types` → `github.com/LarsArtmann/GoReleaser-Wizard/core/types`
  - `github.com/LarsArtmann/GoReleaser-Wizard/internal/utils` → `github.com/LarsArtmann/GoReleaser-Wizard/core/utils`
- **Action:** Update root `go.mod` to add `require github.com/LarsArtmann/GoReleaser-Wizard/core v0.0.0` with replace directive (temporary, until go.work)
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 30 min

### Task 14: Verify core module independence

- **Action:** `cd core && go build ./...` and `cd core && go test ./...`
- **Action:** Verify `core/` has zero imports from `internal/`
- **Verification:** `cd core && go vet ./...` passes; `grep -r "github.com/LarsArtmann/GoReleaser-Wizard/internal" core/` returns nothing
- **Rollback:** `git revert HEAD`
- **Effort:** 10 min

---

## Phase 2: Create `gitutil` Module

**Goal:** Extract git operations into independently buildable module.

### Task 15: Create `gitutil/` directory

- **Action:** `mkdir -p gitutil`
- **Verification:** Directory exists
- **Effort:** 2 min

### Task 16: Move `internal/git/` → `gitutil/`

- **Action:** Rename package from `git` to `gitutil`
- **Action:** `git mv internal/git/commands.go gitutil/commands.go`
- **Action:** Update package declaration and any internal references
- **Verification:** Files moved
- **Effort:** 10 min

### Task 17: Add git-dependent recommendation functions to `gitutil`

- **Action:** Move `GetGitHubOwner()`, `GetGitHubRepo()` from `core/utils/` (or original `internal/utils/`) into `gitutil/`
- **Action:** Update import paths in these functions
- **Verification:** `go build ./...` passes
- **Effort:** 10 min

### Task 18: Create `gitutil/go.mod`

- **Action:** Create `gitutil/go.mod`:

  ```
  module github.com/LarsArtmann/GoReleaser-Wizard/gitutil

  go 1.26.2

  require github.com/LarsArtmann/GoReleaser-Wizard/core v0.0.0

  replace github.com/LarsArtmann/GoReleaser-Wizard/core => ../core
  ```

- **Action:** Run `cd gitutil && go mod tidy`
- **Verification:** `cd gitutil && go build ./...` passes
- **Effort:** 10 min

### Task 19: Update all import paths — `internal/git` → `gitutil`

- **Action:** Find and replace:
  - `github.com/LarsArtmann/GoReleaser-Wizard/internal/git` → `github.com/LarsArtmann/GoReleaser-Wizard/gitutil`
- **Affected files:**
  - `cmd/goreleaser-wizard/generators/goreleaser.go`
  - `cmd/goreleaser-wizard/jobs.go`
  - `cmd/goreleaser-wizard/types/template_data.go`
  - `cmd/goreleaser-wizard/init.go`
  - `core/utils/recommendations.go` (if git functions still referenced)
- **Verification:** `go build ./...` and `go test ./...` pass
- **Rollback:** `git revert HEAD`
- **Effort:** 20 min

### Task 20: Verify gitutil module independence

- **Action:** `cd gitutil && go build ./...` and `cd gitutil && go vet ./...`
- **Action:** Verify `gitutil/` only imports `core/`, not `internal/`
- **Verification:** `cd gitutil && go vet ./...` passes
- **Rollback:** `git revert HEAD`
- **Effort:** 10 min

---

## Phase 3: Workspace Setup + CLI Module Finalization

**Goal:** Wire everything together with `go.work`.

### Task 21: Create `go.work`

- **Action:** Create `go.work`:

  ```
  go 1.26.2

  use (
      ./core
      ./gitutil
      .
  )
  ```

- **Verification:** `go work sync` runs without errors
- **Effort:** 5 min

### Task 22: Update root `go.mod` for workspace

- **Action:** Add `core` and `gitutil` as `require` dependencies in root `go.mod`
- **Action:** Remove temporary `replace` directives (go.work handles local dev)
- **Action:** Run `go mod tidy`
- **Verification:** `go build ./...` passes at root
- **Effort:** 15 min

### Task 23: Remove remaining `internal/` references

- **Action:** Verify no files import `github.com/LarsArtmann/GoReleaser-Wizard/internal/` anymore
- **Action:** If `internal/config/` still exists, decide: keep as `internal/config/` in root module or move to top-level
- **Action:** Clean up empty `internal/` directories
- **Verification:** `grep -r "github.com/LarsArtmann/GoReleaser-Wizard/internal" --include="*.go" .` returns nothing (except test-wizard if applicable)
- **Rollback:** `git revert HEAD`
- **Effort:** 15 min

### Task 24: Full build verification

- **Action:** Run complete verification:
  - `go work sync`
  - `cd core && go build ./... && go test ./... && go vet ./...`
  - `cd gitutil && go build ./... && go vet ./...`
  - `go build ./... && go test ./... && go vet ./...`
  - `go mod tidy` (at root)
- **Verification:** All pass
- **Rollback:** `git revert HEAD`
- **Effort:** 10 min

### Task 25: Verify without workspace (published consumer experience)

- **Action:** Temporarily rename `go.work` to `go.work.bak`
- **Action:** Add `replace` directives in root `go.mod` for `core` and `gitutil`
- **Action:** Run `go build ./...`
- **Action:** Restore `go.work`, remove replace directives
- **Verification:** Both modes work
- **Effort:** 15 min

---

## Phase 4: Cleanup and Documentation

**Goal:** Update all tooling and documentation to reflect new structure.

### Task 26: Update `.go-arch-lint.yml`

- **Action:** Update component paths to reflect new module structure
- **Action:** Add `core/` and `gitutil/` as components with proper dependency rules
- **Verification:** `go-arch-lint` passes
- **Effort:** 20 min

### Task 27: Update `.golangci.yml`

- **Action:** Update path-specific rules for new directory layout
- **Action:** Add `core/` and `gitutil/` specific rules if needed
- **Verification:** `golangci-lint run` passes
- **Effort:** 15 min

### Task 28: Update documentation

- **Action:** Update `AGENTS.md`:
  - New module structure
  - Build/test commands per module (`cd core && go test ./...`, etc.)
  - Updated architecture description
- **Action:** Update `README.md` if it mentions internal structure
- **Action:** Update `FEATURES.md` if needed
- **Verification:** Docs are accurate
- **Effort:** 30 min

### Task 29: Update build system

- **Action:** Update `justfile` commands for multi-module:
  - `just test-core` → `cd core && go test ./...`
  - `just test-gitutil` → `cd gitutil && go test ./...`
  - `just test` → runs all modules
  - `just build` → builds CLI
  - `just lint` → lints all modules
- **Action:** Update `flake.nix` if present (per-module build targets)
- **Action:** Update `.github/workflows/release.yml` for parallel module testing
- **Verification:** `just ci` passes
- **Effort:** 30 min

---

## Dependency Graph Between Tasks

```
Task 1 ─┐
Task 2 ─┤
Task 3 ─┤  Phase 0 (sequential)
Task 4 ─┤
Task 5 ─┤  (4 depends on 3; 5 depends on 4)
Task 6 ─┤
Task 7 ─┘
         │
Task 8 ──┤  Phase 1 (sequential, depends on Phase 0)
Task 9 ──┤
Task 10 ─┤
Task 11 ─┤
Task 12 ─┤  (12 depends on 9,10,11)
Task 13 ─┤  (13 depends on 12)
Task 14 ─┘  (14 depends on 13)
         │
Task 15 ─┤  Phase 2 (sequential, depends on Phase 1)
Task 16 ─┤
Task 17 ─┤
Task 18 ─┤  (18 depends on 16,17)
Task 19 ─┤  (19 depends on 18)
Task 20 ─┘  (20 depends on 19)
         │
Task 21 ─┤  Phase 3 (sequential, depends on Phase 2)
Task 22 ─┤
Task 23 ─┤
Task 24 ─┤  (24 depends on 21,22,23)
Task 25 ─┘  (25 depends on 24)
         │
Task 26 ─┐
Task 27 ─┤  Phase 4 (mostly parallel, depends on Phase 3)
Task 28 ─┤
Task 29 ─┘
```

---

## Estimated Total Effort

| Phase                             | Tasks        | Time         | Cumulative |
| --------------------------------- | ------------ | ------------ | ---------- |
| Phase 0: Dead Code + Split Brains | 1-7          | ~2 hours     | 2h         |
| Phase 1: Core Module              | 8-14         | ~1.5 hours   | 3.5h       |
| Phase 2: Gitutil Module           | 15-20        | ~1 hour      | 4.5h       |
| Phase 3: Workspace Setup          | 21-25        | ~1 hour      | 5.5h       |
| Phase 4: Cleanup + Docs           | 26-29        | ~1.5 hours   | 7h         |
| **Total**                         | **29 tasks** | **~7 hours** |            |

---

## Per-Task Verification Checklist

After each task:

- [ ] `go build ./...` passes at root
- [ ] `go test ./...` passes at root
- [ ] `go vet ./...` passes at root
- [ ] No new compiler warnings
- [ ] Commit with descriptive message

After Phase 1+ (per module):

- [ ] `cd core && go build ./...` passes independently
- [ ] `cd gitutil && go build ./...` passes independently

After Phase 3 (workspace):

- [ ] `go work sync` succeeds
- [ ] `go mod tidy` at root changes nothing
- [ ] Each module builds independently
- [ ] Full test suite passes with workspace

After Phase 4 (final):

- [ ] `just ci` passes (or equivalent full pipeline)
- [ ] `golangci-lint run` passes
- [ ] Documentation is accurate
- [ ] No `internal/` references remain
