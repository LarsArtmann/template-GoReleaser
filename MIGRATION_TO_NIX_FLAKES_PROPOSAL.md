# Migration to Nix Flakes — Proposal

**Project:** GoReleaser-Wizard  
**Date:** 2026-04-09  
**Status:** Draft — Pending Approval  
**Author:** Generated from tooling audit

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State Analysis](#2-current-state-analysis)
3. [Why Nix Flakes](#3-why-nix-flakes)
4. [Proposed Architecture](#4-proposed-architecture)
5. [Migration Plan](#5-migration-plan)
6. [Risk Assessment](#6-risk-assessment)
7. [Open Questions](#7-open-questions)

---

## 1. Executive Summary

This proposal outlines migrating GoReleaser-Wizard's development toolchain from ad-hoc shell scripts, imperative `go install` commands, and manual dependency management to **Nix Flakes** — a declarative, reproducible, and locked approach to defining the entire development environment, build pipeline, and CI toolchain.

The project currently manages **30+ external tools** across two justfiles and seven shell scripts, with no lockfile or reproducibility guarantees. A Nix Flake would replace this sprawl with a single `flake.nix` that provides:

- **Reproducible dev shells** — every contributor gets identical tool versions
- **Declarative builds** — `nix build` produces the binary without Go installed on the host
- **Locked dependencies** — `flake.lock` ensures bit-for-bit reproducibility
- **CI parity** — same flake drives local dev and CI/CD

---

## 2. Current State Analysis

### 2.1 Build & Task System

| File                 | Purpose                                               | Lines |
| -------------------- | ----------------------------------------------------- | ----- |
| `justfile`           | Core build/test/fmt/clean/verify commands             | 78    |
| `dev/arch-lint.just` | Enterprise linting, security, profiling, benchmarking | 1,476 |
| `verify.sh`          | Project verification (tools + structure + build)      | 40    |

### 2.2 External Tools Currently Required

The project depends on a significant toolchain, all installed imperatively:

#### Core Go Toolchain

| Tool          | Used In                    | Install Method          |
| ------------- | -------------------------- | ----------------------- |
| `go` (1.26.1) | everywhere                 | manual / system package |
| `goreleaser`  | justfile (snapshot, check) | manual                  |
| `git`         | verify.sh, justfile        | system                  |

#### Linting & Code Quality

| Tool            | Used In        | Version Pinning                                      |
| --------------- | -------------- | ---------------------------------------------------- |
| `golangci-lint` | arch-lint.just | **none** — uses `$(go env GOPATH)/bin/golangci-lint` |
| `go-arch-lint`  | arch-lint.just | `v1.14.0` (constant only)                            |
| `gofumpt`       | arch-lint.just | **none** — `go install -tool ... @latest`            |
| `goimports`     | arch-lint.just | **none** — `go install -tool ... @latest`            |
| `dupl`          | arch-lint.just | **none** — `go install -tool ... @latest`            |
| `govulncheck`   | arch-lint.just | **none** — `go get -tool ... @latest`                |
| `capslock`      | arch-lint.just | **none** — `go install -tool ... @latest`            |
| `nilaway`       | arch-lint.just | **none** — `go get -tool ... @latest`                |
| `go-licenses`   | arch-lint.just | **none** — `go get -tool ... @latest`                |
| `goleak`        | arch-lint.just | **none** — `go get -tool ... @latest`                |
| `jscpd`         | arch-lint.just | **none** — `bun install -g jscpd`                    |
| `jq`            | arch-lint.just | system dependency                                    |
| `bc`            | arch-lint.just | system dependency                                    |

#### Development & Profiling

| Tool    | Used In                    | Version Pinning                           |
| ------- | -------------------------- | ----------------------------------------- |
| `templ` | arch-lint.just (build)     | **none** — `go install -tool ... @latest` |
| `air`   | arch-lint.just (dev)       | **none** — `go install -tool ... @latest` |
| `pprof` | arch-lint.just (profiling) | via `go tool`                             |

#### GitHub Management Scripts

| Script                                 | Dependencies      |
| -------------------------------------- | ----------------- |
| `scripts/create_milestones.sh`         | `gh` CLI          |
| `scripts/assign_milestones.sh`         | `gh` CLI          |
| `scripts/create_critical_issue.sh`     | `gh` CLI          |
| `scripts/organize_github_complete.sh`  | all above scripts |
| `scripts/add_architecture_comments.sh` | `gh` CLI          |
| `scripts/generate-license.sh`          | `yq`, `sed`       |

### 2.3 Key Problems Identified

| Problem                               | Impact                                                                              | Severity |
| ------------------------------------- | ----------------------------------------------------------------------------------- | -------- |
| **No tool version locking**           | "Works on my machine" — 15+ tools install `@latest`                                 | High     |
| **Imperative installs in recipes**    | Tools auto-install mid-recipe, causing side effects and nondeterminism              | High     |
| **Missing scripts**                   | `scripts/check-cmd-single.sh` referenced but doesn't exist (arch-lint.just:365)     | Medium   |
| **Stale Dockerfile**                  | `FROM golang:1.24-alpine` while `go.mod` requires Go 1.26.1                         | Medium   |
| **No `.golangci.yml` on disk**        | Referenced in arch-lint.just but file doesn't exist in repo                         | Medium   |
| **No `.go-arch-lint.yml` on disk**    | Referenced in arch-lint.just but file doesn't exist in repo                         | Medium   |
| **No CI/CD pipeline**                 | No `.github/workflows/` at all                                                      | Medium   |
| **Bootstrap downloads from internet** | `bootstrap.sh` fetched via `curl` at runtime — supply chain risk                    | High     |
| **Mixed tool naming**                 | Some tools use `$(go env GOPATH)/bin/`, others use `command -v`                     | Low      |
| **External scripts from remote**      | `bootstrap.sh` and `test-bootstrap-simple-bdd.sh` downloaded from GitHub at runtime | High     |

### 2.4 Tool Version Inventory

Only **1 of 15+ tools** has a version constant in the codebase:

```
GO_ARCH_LINT_VERSION := "v1.14.0"
CAPSLOCK_VERSION := "latest"   # ← "latest" is not a version
```

Everything else installs `@latest` — meaning two developers running `just lint` on different days get different toolchains.

---

## 3. Why Nix Flakes

### 3.1 What Nix Flakes Provide

| Capability                    | Current State                                              | With Nix Flakes                                       |
| ----------------------------- | ---------------------------------------------------------- | ----------------------------------------------------- |
| **Reproducible environments** | No — tools install `@latest`                               | Yes — `flake.lock` pins exact versions                |
| **Declarative toolchain**     | Spread across 2 justfiles, 7 scripts                       | Single `flake.nix`                                    |
| **Hermetic builds**           | No — depends on host Go, system tools                      | Yes — `nix build` works without Go installed          |
| **Onboarding**                | Install Go, just, goreleaser, gh, yq, bc, jq, ... manually | `nix develop` — done                                  |
| **CI caching**                | Manual layer caching                                       | Built-in Nix cache (`cachix` or GitHub Actions cache) |
| **Cross-platform**            | Assumed macOS/Linux                                        | Explicit `systems` declaration                        |
| **Rollback**                  | Manual                                                     | `nix rollback` — instant                              |

### 3.2 Why Not Alternatives

| Alternative               | Why Not                                                                            |
| ------------------------- | ---------------------------------------------------------------------------------- |
| **Docker dev containers** | Heavy, slow startup, poor editor integration (LSP/gopls), no nested Docker support |
| **asdf/mise**             | Only manages tool versions, not the build itself; no package-level caching         |
| **nix-shell (legacy)**    | No lockfile, no standardized inputs, no `nix build` integration                    |
| **Devbox**                | Built on Nix but adds abstraction; less control, smaller package set               |
| **Bazel/Buck**            | Massive complexity for a Go CLI project; Go's build system is already fast         |

### 3.3 Nix Flakes Specifically (vs. Classic Nix)

- **`flake.lock`** — Automatic lockfile (like `go.sum` for the dev environment)
- **Pure evaluation** — No reading from `$NIX_PATH` or arbitrary files
- **Composable inputs** — Reference other flakes directly (nixpkgs, flake-utils, etc.)
- **Standardized structure** — `packages`, `devShells`, `apps`, `checks`, `formatter`
- **`nix develop`** — Drop-in replacement for `shell.nix` with better isolation
- **`nix flake check`** — Run all checks (tests, linting, formatting) in parallel

---

## 4. Proposed Architecture

### 4.1 File Structure

```
GoReleaser-Wizard/
├── flake.nix                    # Main flake definition
├── flake.lock                   # Auto-generated lockfile (committed)
├── nix/
│   ├── packages/
│   │   └── goreleaser-wizard.nix  # Build derivation for the binary
│   └── devshell/
│       └── default.nix            # Dev shell extras (shell hooks, env vars)
├── justfile                     # Simplified — delegates to nix for tool paths
├── dev/
│   └── arch-lint.just           # Simplified — tools guaranteed by dev shell
├── .envrc                       # Optional: direnv auto-loading
└── ... (existing files unchanged)
```

### 4.2 Proposed `flake.nix`

```nix
{
  description = "GoReleaser-Wizard — Interactive GoReleaser configuration wizard";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Go version matching go.mod
        go = pkgs.go_1_26;
        buildGoModule = pkgs.buildGo126Module;

        # Project metadata
        version = "0.1.0";

        # The main package
        goreleaser-wizard = buildGoModule {
          pname = "goreleaser-wizard";
          inherit version;

          src = ./.;

          # Hash of the vendor directory — update with `nix build` after go.mod changes
          # Run: nix build && copy the expected hash from the error
          vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

          subPackages = [ "cmd/goreleaser-wizard" ];

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
            "-X main.commit=${self.rev or "dirty"}"
            "-X main.date=1970-01-01T00:00:00Z"
          ];

          meta = with pkgs.lib; {
            description = "Interactive GoReleaser configuration wizard";
            homepage = "https://github.com/LarsArtmann/GoReleaser-Wizard";
            license = licenses.unlicense;
            mainProgram = "goreleaser-wizard";
          };
        };

        # Linting tools — all pinned to specific Go versions
        lintTools = with pkgs; [
          # Core Go toolchain
          go
          gotools          # goimports, godoc, etc.

          # Linting
          golangci-lint
          gofumpt
          govulncheck

          # Architecture
          go-arch-lint

          # Security
          capslock

          # Duplication detection
          dupl

          # Formatting
          nixpkgs-fmt      # For formatting the flake itself

          # Task runner
          just

          # General utilities
          jq
          yq-go
          git
          gh
        ];

      in
      {
        # Packages — `nix build .#default` or `nix build`
        packages = {
          default = goreleaser-wizard;
          goreleaser-wizard = goreleaser-wizard;
        };

        # Dev shells — `nix develop`
        devShells = {
          default = pkgs.mkShell {
            # Inherit build dependencies
            inputsFrom = [ goreleaser-wizard ];

            # Additional development tools
            packages = lintTools;

            # Environment variables
            GOROOT = "${go}/share/go";
            GOFLAGS = "-mod=mod";

            shellHook = ''
              echo "🧙 GoReleaser-Wizard Development Environment"
              echo "  Go:        $(go version)"
              echo "  Tools:     golangci-lint, gofumpt, go-arch-lint, govulncheck, just"
              echo ""
              echo "Run 'just --list' to see available commands."
            '';
          };

          # Minimal shell — only Go and just (for CI)
          ci = pkgs.mkShell {
            inputsFrom = [ goreleaser-wizard ];
            packages = with pkgs; [
              go
              gotools
              golangci-lint
              gofumpt
              govulncheck
              go-arch-lint
              just
              jq
            ];
          };
        };

        # Apps — `nix run .#default`
        apps = {
          default = {
            type = "app";
            program = "${goreleaser-wizard}/bin/goreleaser-wizard";
          };
        };

        # Checks — `nix flake check`
        checks = {
          build = goreleaser-wizard;

          # Verify the binary runs
          version-check = pkgs.runCommandLocal "version-check" { } ''
            ${goreleaser-wizard}/bin/goreleaser-wizard --help
            touch $out
          '';
        };

        # Formatter — `nix fmt`
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}
```

### 4.3 Proposed `.envrc` (Optional — for direnv users)

```bash
# .envrc
use flake
```

This enables automatic shell activation when entering the project directory.

### 4.4 Justfile Changes

The justfiles remain **largely unchanged** — Nix provides the tools, `just` orchestrates them. The key difference:

**Before (current):**

```just
lint-code:
    @if command -v golangci-lint >/dev/null 2>&1; then \
        golangci-lint run --config .golangci.yml; \
    else \
        echo "❌ golangci-lint not installed. Run 'just install' first."; \
        exit 1; \
    fi
```

**After (with Nix):**

```just
lint-code:
    golangci-lint run --config .golangci.yml
```

All the `command -v` guards, auto-install fallbacks, and `@latest` installs disappear because the dev shell guarantees the tool exists at a pinned version.

### 4.5 What Each Nix Command Replaces

| Current Command                                         | Nix Equivalent                  | Notes                                                |
| ------------------------------------------------------- | ------------------------------- | ---------------------------------------------------- |
| `just install` (imperative tool installs)               | `nix develop` (declarative)     | Tools are in the shell, not in `$GOPATH/bin`         |
| `go build -o goreleaser-wizard ./cmd/goreleaser-wizard` | `nix build`                     | Hermetic, cached, reproducible                       |
| `./verify.sh` (tool existence checks)                   | `nix develop` guarantees        | All tools available or build fails                   |
| `just bootstrap` (downloads scripts from internet)      | Not needed                      | Nix handles everything                               |
| Manual Go installation                                  | `nix develop`                   | Go is a flake dependency                             |
| `goreleaser build --snapshot --clean`                   | `nix build .#goreleaser-wizard` | For quick builds; Goreleaser still used for releases |

---

## 5. Migration Plan

### Phase 1: Foundation (Day 1)

**Goal:** Add `flake.nix` without breaking existing workflow.

- [ ] Create `flake.nix` with `devShells.default` containing all linting tools
- [ ] Create `flake.lock` via `nix flake lock`
- [ ] Test: `nix develop` provides all tools
- [ ] Test: Existing justfile recipes work inside `nix develop`
- [ ] Commit `flake.nix` and `flake.lock`
- [ ] **No changes to justfiles or scripts yet**

**Acceptance Criteria:**

- `nix develop -c just ci` passes
- `nix develop -c just lint` passes
- All tools available: `go`, `golangci-lint`, `gofumpt`, `go-arch-lint`, `govulncheck`, `just`

### Phase 2: Build Package (Day 2)

**Goal:** Add `nix build` support.

- [ ] Add `packages.goreleaser-wizard` derivation to `flake.nix`
- [ ] Compute initial `vendorHash` (empty → error → copy from error)
- [ ] Test: `nix build` produces a working binary
- [ ] Test: `nix run .` executes the wizard
- [ ] Add `nix run` as `apps.default`
- [ ] Update `Dockerfile` to use `nix build` output (optional)

**Acceptance Criteria:**

- `nix build && ./result/bin/goreleaser-wizard --help` works
- `nix run .` works
- `vendorHash` is correct and committed

### Phase 3: Cleanup (Day 3)

**Goal:** Simplify justfiles by removing tool-installation logic.

- [ ] Remove all `command -v ... >/dev/null 2>&1` guards from justfiles
- [ ] Remove all `go install -tool ... @latest` fallback installs
- [ ] Remove `install` recipe from arch-lint.just (tools come from Nix)
- [ ] Remove `bootstrap`, `bootstrap-diagnose`, `bootstrap-fix`, `bootstrap-verbose`, `bootstrap-quick`, `bootstrap-test` recipes (all download remote scripts)
- [ ] Fix missing `scripts/check-cmd-single.sh` or remove the recipe
- [ ] Simplify error messages (no more "Run 'just install' first")
- [ ] Update `verify.sh` to check `nix develop` availability instead of individual tools

**Acceptance Criteria:**

- No `go install` commands remain in any recipe
- No `curl` downloads remain in any recipe
- All recipes assume tools are available (Nix provides them)

### Phase 4: CI Integration (Day 4)

**Goal:** Use the same flake in CI.

- [ ] Create `.github/workflows/ci.yml` using `nix develop .#ci`
- [ ] Use `cachix` or GitHub Actions Nix cache for build caching
- [ ] Run `nix flake check` in CI (builds, checks, tests)
- [ ] Remove imperative CI setup steps (manual Go installs, etc.)

**Example CI workflow:**

```yaml
name: CI
on: [push, pull_request]
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: cachix/install-nix-action@v30
        with:
          extra_nix_config: |
            experimental-features = nix-command flakes
      - uses: cachix/cachix-action@v15
        with:
          name: goreleaser-wizard
          signingKey: "${{ secrets.CACHIX_SIGNING_KEY }}"
      - run: nix flake check
      - run: nix develop .#ci -c just test
```

**Acceptance Criteria:**

- CI runs entirely from the flake
- No manual tool installation steps in CI
- Build cache works across runs

### Phase 5: Polish & Documentation (Day 5)

**Goal:** Finalize the migration.

- [ ] Update `README.md` with Nix-based setup instructions
- [ ] Update `CONTRIBUTING.md` (if exists)
- [ ] Add `.envrc` for direnv users (optional, committed)
- [ ] Add `nix fmt` to `just format` recipe
- [ ] Verify all arch-lint.just recipes work in `nix develop`
- [ ] Create migration guide for contributors without Nix
- [ ] Add `flake.nix` and `flake.lock` to `.gitignore` exceptions if needed

**Acceptance Criteria:**

- New contributor can run `nix develop` and be fully set up
- All existing functionality preserved
- Documentation reflects new workflow

---

## 6. Risk Assessment

### 6.1 Risks

| Risk                            | Likelihood | Impact | Mitigation                                                                                                                          |
| ------------------------------- | ---------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------- |
| **Nix learning curve**          | Medium     | Medium | Phased migration; justfiles remain the primary interface                                                                            |
| **Missing packages in nixpkgs** | Low        | Medium | `go-arch-lint`, `capslock` may need overlay or `buildGoModule`                                                                      |
| **CI slower initially**         | Low        | Low    | Nix cache (cachix) mitigates; first run downloads, subsequent are cached                                                            |
| **Contributor resistance**      | Medium     | Low    | Keep justfiles as primary interface; Nix is optional (tools can still be installed manually)                                        |
| **vendorHash updates**          | Medium     | Low    | Documented process: change go.mod → `nix build` → copy hash from error                                                              |
| **Replace directive in go.mod** | Low        | Medium | `replace github.com/larsartmann/go-composable-business-types => /Users/...` breaks hermetic builds — must be removed or conditional |
| **Platform-specific issues**    | Low        | Medium | `flake-utils.lib.eachDefaultSystem` covers Linux/macOS; test on both                                                                |

### 6.2 Critical Blocker: `go.mod` Replace Directive

The current `go.mod` contains:

```
replace github.com/larsartmann/go-composable-business-types => /Users/larsartmann/projects/go-composable-business-types
```

This **breaks hermetic Nix builds** because it references an absolute local path. Options:

1. **Remove it** — if the dependency is published/publishable
2. **Use a Nix overlay** — override the Go module with a local path via `goSum`
3. **Keep it for dev only** — Nix build uses the real module; local replace only in dev shell
4. **Make it conditional** — use a separate `go.mod` for development

**Recommendation:** Option 1 (publish the dependency) or Option 3 (dev-only replace).

### 6.3 Tool Availability in nixpkgs

| Tool            | Available in nixpkgs   | Notes                                         |
| --------------- | ---------------------- | --------------------------------------------- |
| `go`            | Yes                    | `go_1_26`                                     |
| `golangci-lint` | Yes                    | May need overlay for Go 1.26 compatibility    |
| `gofumpt`       | Yes                    | `mvdan.cc/gofumpt`                            |
| `goimports`     | Yes                    | Via `gotools` package                         |
| `go-arch-lint`  | **Check needed**       | May require custom `buildGoModule` derivation |
| `govulncheck`   | Yes                    | `golang.org/x/vuln/cmd/govulncheck`           |
| `capslock`      | **Check needed**       | May require custom derivation                 |
| `nilaway`       | **Check needed**       | May require custom derivation                 |
| `dupl`          | Yes                    | `github.com/mibk/dupl`                        |
| `go-licenses`   | Yes                    | `github.com/google/go-licenses`               |
| `goleak`        | Yes (as library)       | `go.uber.org/goleak` — used in test code      |
| `goreleaser`    | Yes                    | `goreleaser` package                          |
| `just`          | Yes                    | `just` package                                |
| `gh`            | Yes                    | `gh` package                                  |
| `yq`            | Yes                    | `yq-go` package                               |
| `jq`            | Yes                    | `jq` package                                  |
| `jscpd`         | Yes (via nodePackages) | `nodePackages.jscpd`                          |
| `templ`         | Yes                    | `templ` package                               |
| `air`           | Yes                    | `cosmtrek/air`                                |

Tools marked **"Check needed"** will require either:

- Finding the correct nixpkgs attribute name
- Creating a small custom `buildGoModule` derivation
- Falling back to `buildGoModule` with `src = pkgs.fetchFromGitHub { ... }`

---

## 7. Open Questions

| #   | Question                                                                                                  | Decision Needed By |
| --- | --------------------------------------------------------------------------------------------------------- | ------------------ |
| 1   | Should `flake-parts` be used for modularity, or plain `flake-utils`?                                      | Phase 1            |
| 2   | Is `github.com/larsartmann/go-composable-business-types` publishable?                                     | Phase 2            |
| 3   | Should `direnv` + `.envrc` be committed to the repo?                                                      | Phase 5            |
| 4   | Should `cachix` be set up for CI caching?                                                                 | Phase 4            |
| 5   | What is the target `vendorHash` strategy? (computed vs. fake)                                             | Phase 2            |
| 6   | Should the Dockerfile be updated to use Nix-built binary?                                                 | Phase 2+           |
| 7   | Should `.golangci.yml` and `.go-arch-lint.yml` be created as part of this migration?                      | Phase 3            |
| 8   | Is macOS ARM (aarch64-darwin) the primary dev platform?                                                   | Phase 1            |
| 9   | Should the missing `scripts/check-cmd-single.sh` be created or should the recipe be removed?              | Phase 3            |
| 10  | Should the `bootstrap*.sh` remote-download recipes be removed entirely or replaced with Nix alternatives? | Phase 3            |

---

## Appendix A: Quick Reference — Nix Commands for This Project

| Command                                 | Purpose                        |
| --------------------------------------- | ------------------------------ |
| `nix develop`                           | Enter dev shell with all tools |
| `nix develop -c just ci`                | Run CI pipeline in Nix shell   |
| `nix develop .#ci`                      | Enter minimal CI shell         |
| `nix build`                             | Build the binary               |
| `nix build .#goreleaser-wizard`         | Build explicitly by name       |
| `nix run .`                             | Run the wizard directly        |
| `nix flake check`                       | Run all checks (build + tests) |
| `nix flake update`                      | Update all flake inputs        |
| `nix fmt`                               | Format `.nix` files            |
| `nix flake lock --update-input nixpkgs` | Update only nixpkgs            |

## Appendix B: Tool Version Comparison — Before and After

| Tool            | Before (Current)                              | After (Nix Flake)                             |
| --------------- | --------------------------------------------- | --------------------------------------------- |
| `go`            | 1.26.1 (manual install)                       | 1.26.x (from nixpkgs, locked in `flake.lock`) |
| `golangci-lint` | `@latest` (nondeterministic)                  | Specific nixpkgs revision                     |
| `gofumpt`       | `@latest`                                     | Specific nixpkgs revision                     |
| `goimports`     | `@latest`                                     | Specific nixpkgs revision                     |
| `go-arch-lint`  | `v1.14.0` (declared, but installs `@v1.14.0`) | Specific nixpkgs revision                     |
| `govulncheck`   | `@latest`                                     | Specific nixpkgs revision                     |
| `capslock`      | `latest` (literal string "latest")            | Specific nixpkgs revision                     |
| `nilaway`       | `@latest`                                     | Specific nixpkgs revision                     |
| `dupl`          | `@latest`                                     | Specific nixpkgs revision                     |
| `go-licenses`   | `@latest`                                     | Specific nixpkgs revision                     |
| `just`          | manual install                                | Specific nixpkgs revision                     |
| `gh`            | manual install                                | Specific nixpkgs revision                     |
| `yq`            | manual install                                | Specific nixpkgs revision                     |
| `jq`            | system package                                | Specific nixpkgs revision                     |
| `bc`            | system package                                | Specific nixpkgs revision                     |
| `jscpd`         | `bun install -g`                              | nixpkgs `nodePackages.jscpd`                  |
| `templ`         | `@latest`                                     | Specific nixpkgs revision                     |
| `air`           | `@latest`                                     | Specific nixpkgs revision                     |
| `goreleaser`    | manual install                                | Specific nixpkgs revision                     |

**Summary:** All 19 tools move from nondeterministic/manual to fully locked.

## Appendix C: Dependency Map

```
flake.nix
├── inputs.nixpkgs → Go toolchain, all linting tools, utilities
├── inputs.flake-utils → System matrix (linux/macOS × amd64/arm64)
│
├── outputs.packages.goreleaser-wizard
│   └── buildGo126Module ← go.mod + vendor hash
│
├── outputs.devShells.default
│   ├── inputsFrom: goreleaser-wizard (build deps)
│   └── packages: lintTools (19 tools)
│
├── outputs.devShells.ci
│   └── packages: minimal CI toolset
│
├── outputs.apps.default → goreleaser-wizard binary
│
├── outputs.checks → build + version-check
│
└── outputs.formatter → nixpkgs-fmt
```

---

_This proposal is a living document. Update as decisions are made during migration._
