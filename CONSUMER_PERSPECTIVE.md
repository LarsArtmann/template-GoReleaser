# Consumer Perspective: What's Missing

A frank assessment of GoReleaser-Wizard from the perspective of someone who just installed it and wants to release their Go project.

---

## Critical Gaps (Blocks a Real User Today)

### 1. No Saved Preferences — You Re-Answer Everything Every Time

The wizard asks 20+ questions every single run. There's no `.goreleaser-wizard.yaml` that remembers your last answers. If you want to regenerate (new platform, different signing), you start from scratch.

**What a consumer expects:** "I ran this last month. Just re-use those answers and let me change what I need."

**Reality:** The config struct has two fields: `Debug` and `NoColor`. Everything else is ephemeral.

### 2. Templates Don't Adapt to Project Type

The wizard asks whether you're building a CLI, web service, library, gRPC service, daemon, etc. — 10 project types total. Then it generates **the same template regardless**. A library gets Docker publishing steps. A web service gets Homebrew formula sections. A daemon gets Scoop (Windows package manager) config.

**What a consumer expects:** "I said web service — give me a Docker-first config with health checks, no Homebrew."

**Reality:** Every project type gets the same kitchen-sink `.goreleaser.yaml` with conditional sections driven only by boolean flags, not by the project type itself.

### 3. Hardcoded Author Values in Template Files

The `templates/goreleaser.yaml.tmpl` file (used by some code paths) contains hardcoded references to `LarsArtmann/homebrew-tap`, `LarsArtmann/nur-packages`, and `goreleaser@larsartmann.com`. A new user running the wizard could end up with configs pointing to someone else's repositories.

**What a consumer expects:** "My config, my repos, my email."

**Note:** The embedded templates in `embedded.go` are clean — but the `.tmpl` files coexist and are reachable, creating confusion about which template is actually used.

### 4. No `--dry-run` or Preview

Every generator has a `GeneratePreview()` method implemented. The workflow engine supports `dryRun`. But there's no `--dry-run` flag exposed on any CLI command.

**What a consumer expects:** "Show me what you'd generate before writing files to my project."

**Reality:** Files are written immediately. The only safety net is `--force` being opt-in for overwrites.

---

## High-Impact Gaps (Hurts Adoption and Retention)

### 5. Only GitHub Actions — No GitLab, Bitbucket, or Gitea CI

The wizard lets you select your Git provider (GitHub, GitLab, Bitbucket, Gitea, self-hosted). Then it generates a GitHub Actions workflow regardless of what you picked.

**What a consumer expects:** "I picked GitLab — give me a `.gitlab-ci.yml`."

**Reality:** Selecting GitLab/Bitbucket/Gitea is cosmetic. Only GitHub Actions templates exist.

### 6. No Config Migration or Update Path

The wizard has `migrate` and `update` workflow types scaffolded in code. But no CLI subcommand exposes them. If GoReleaser changes its config format (it does regularly), there's no way to update your existing config through the wizard.

**What a consumer expects:** "GoReleaser v2 changed the config format. Run `goreleaser-wizard migrate` to update my `.goreleaser.yaml`."

**Reality:** Re-run `init` from scratch and hope the new output is compatible with what you had.

### 7. No Multi-Binary / Monorepo Support

Many Go projects have multiple binaries in `cmd/`. The wizard detects `cmd/*/main.go` but only picks the first one. There's no way to configure builds for multiple binaries.

**What a consumer expects:** "I have `cmd/server` and `cmd/cli` — generate builds for both."

**Reality:** You get one binary. Edit the `.goreleaser.yaml` manually for the rest.

### 8. Snap/AUR/Scoop Collections But No Generation

The TUI collects whether you want Snap packages. The domain model has Scoop and AUR fields. But no templates exist for any of these package formats. Homebrew is the only package manager with actual generation.

**What a consumer expects:** "I enabled Snap — where's my `snapcraft.yaml`?"

**Reality:** The flag does nothing.

---

## Quality-of-Life Gaps (Would Make the Tool Feel Professional)

### 9. No Shell Completions

Cobra has built-in completion support for bash, zsh, fish, and PowerShell. It's zero-effort to add. But no completion command is registered.

**Impact:** Tab completion for commands, flags, and enum values (project types, platforms, etc.) would make the CLI feel native.

### 10. No Example Outputs

The `examples/` directory contains generic GoReleaser documentation, not wizard-specific examples. A consumer can't see what the tool will produce for different project types without running it.

**What would help:** Example `.goreleaser.yaml` outputs for each project type in the README or an `examples/` directory with annotated outputs.

### 11. No Post-Setup Guidance

After running the wizard, the next-steps message says "Review generated files" and "Run validate." It doesn't mention:

- Setting up the `HOMEBREW_TAP_GITHUB_TOKEN` secret (required for Homebrew publishing)
- Creating a GitHub environment for releases
- Running `goreleaser release --snapshot --clean` to test locally
- What to do about Cosign key generation for signing
- How to create and push a first `v1.0.0` tag

**Impact:** A first-time user generates a config that references secrets and tools they don't have, then hits cryptic errors on their first release attempt.

### 12. No Config Diff or Change Summary

When re-running the wizard on an existing project, there's no diff or summary of what changed between the old and new config.

**What a consumer expects:** "You changed: added arm64 builds, removed FreeBSD, enabled Docker publishing. Apply? [y/N]"

**Reality:** Silent overwrite (with backup, but no visibility into what changed).

### 13. FEATURES.md Is Stale

FEATURES.md claims Dockerfile and Homebrew generation are "not implemented" when both are functional. This undermines trust — a consumer reading the docs would think the tool does less than it actually does.

---

## Summary: Prioritized Consumer Impact

| #  | Gap                                    | Impact                              | Effort |
| -- | -------------------------------------- | ----------------------------------- | ------ |
| 1  | Save/load wizard preferences           | Every user, every re-run            | Medium |
| 2  | Project-type-aware templates           | Misleading UX — 10 types, 1 output  | High   |
| 3  | Remove hardcoded author values         | Trust/ownership issue               | Low    |
| 4  | Expose `--dry-run` flag                | Already implemented, just not wired | Low    |
| 5  | GitLab/Bitbucket CI templates          | Non-GitHub users are stranded       | High   |
| 6  | Config migration command               | Already scaffolded                  | Medium |
| 7  | Multi-binary support                   | Common Go project pattern           | Medium |
| 8  | Generate what you collect (Snap, etc.) | Trust — flags that do nothing       | Medium |
| 9  | Shell completions                      | Professional polish, zero effort    | Low    |
| 10 | Example outputs in docs                | Helps evaluation before install     | Low    |
| 11 | Post-setup guidance                    | First release success rate          | Low    |
| 12 | Config diff on re-run                  | Confidence in changes               | Medium |
| 13 | Fix stale FEATURES.md                  | Docs accuracy                       | Low    |
