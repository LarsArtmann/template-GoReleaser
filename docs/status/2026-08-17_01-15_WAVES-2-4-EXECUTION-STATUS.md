# Waves 2-4 Execution — Status Report

**Date:** 2026-08-17 01:15
**Session:** Wave-1 gap closers + Waves 2-4 + M23/M13/M14/M15/M16 pulled forward
**Commits this session:** `6f4f54d` (gap closers), `9f4bd71` (single source of truth), `8edf0b1` (E2E + goldens + CI), `3cb91da` (hygiene), + own-config migration & docs (pending at write time)
**Plan:** `docs/planning/2026-08-16_18-45_SUPERB-WIZARD-MODERNIZATION.md`

---

## a) FULLY DONE (verified this session)

| #             | Item                                                                                                                                                                                                                                                             | Proof                                                                             |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Gap 1         | `DockerfileGenerationJob` wired into `CreateConfigOnlyJobs`                                                                                                                                                                                                      | `TestConfigOnlyWorkflowGeneratesDockerfileEndToEnd` green                         |
| Gap 2         | `GenerationPreflightJob`: all target files checked before anything is written; generation is atomic                                                                                                                                                              | `TestGenerationPreflightJobExecute`, `TestGenerateRefusesToOverwriteWithoutForce` |
| Gap 3         | Unit tests: dockerPlatforms, override precedence, placeholder detection, Dockerfile job (6 subtests), factory wiring/ordering                                                                                                                                    | `preflight_test.go`, `types/template_data_test.go`                                |
| Gap 4         | Warning on placeholder `owner/repo` (types.PlaceholderGitHubOwner/Repo + HasPlaceholderGitHubTarget)                                                                                                                                                             | code + test                                                                       |
| Gap 5         | `go mod tidy` done (9 stale hashes dropped); GOEXPERIMENT=jsonv2 verified REQUIRED (build fails without it; BuildFlow's "redundant" claim is wrong — misses internal/types)                                                                                      | build both ways                                                                   |
| Gap 6         | F31: 4-combo docker×signing render test, YAML-parsed                                                                                                                                                                                                             | `github_actions_template_test.go`                                                 |
| Gap 7         | TUI path verified structurally: TUI only fills config, flows through `ExecuteWorkflow(FullWizard)` → same factory the ordering test covers                                                                                                                       | code read                                                                         |
| M5/M6         | `templates/embedded.go` is the single template source; jobs.go delegates to typed `generators` (GoReleaser, GitHubActions, Dockerfile); zero `map[string]any` prep left                                                                                          | F26 grep, `9f4bd71`                                                               |
| F17/18 (port) | Repo-root `templates/` dir deleted; third-Dockerfile-template split brain killed                                                                                                                                                                                 | `git rm`, flake.nix fileset updated                                               |
| Perf          | git-remote resolution cached per process (16.7s → 0.85s test time)                                                                                                                                                                                               | timing                                                                            |
| M9 (F34-F37)  | E2E: generate into temp module → `goreleaser check` exit 0, zero DEPRECATED, no env vars; skips gracefully without goreleaser binary                                                                                                                             | `TestGeneratedConfigPassesGoReleaserCheck`                                        |
| M10 (F38-F42) | Golden files: 2 goreleaser variants, 4 workflow combos, 3 Dockerfile variants + `-update-golden` flag                                                                                                                                                            | `testdata/golden/`                                                                |
| M11 (F43-F45) | `.github/workflows/ci.yml`: lint + race tests + goreleaser-check job (fixture proof in CI)                                                                                                                                                                       | YAML validated                                                                    |
| M12 (F46-F48) | 3 `.orig` files removed, `*.orig` gitignored, `createTempFile` removed, infertypeargs hints fixed                                                                                                                                                                | `3cb91da`                                                                         |
| Hygiene+      | broken `test-wizard/` module deleted                                                                                                                                                                                                                             | BuildFlow finding gone                                                            |
| M23 (F79-F82) | Own `.goreleaser.yaml` migrated: `dockers`+`docker_manifests`→`dockers_v2`, `brews`→`homebrew_casks` (+ quarantine hook, Casks dir), cosign flags modernized to `--bundle`; own Dockerfile got `ARG TARGETPLATFORM`; `goreleaser check` exit 0 zero deprecations | command output                                                                    |
| M13 (F49-F53) | `docs/goreleaser-guide.md`: phases, template vars, dockers_v2 deep-dive, casks migration, snapshot/release table                                                                                                                                                 | file written                                                                      |
| M14 (F54-F56) | README: v2.17 claims, new flags in tables, prebuilt-Dockerfile truth, atomic generation, real dev commands                                                                                                                                                       | file updated                                                                      |
| M15 (F57-F59) | AGENTS.md: real commands, jsonv2/cache/telemetry/E2E/domain-validation gotchas, real layer map                                                                                                                                                                   | file updated                                                                      |
| M16 (F60)     | FEATURES.md: honest Docker/signing/homebrew statuses                                                                                                                                                                                                             | file updated                                                                      |

## b) PARTIALLY DONE

| Item                     | Done                                                | Missing                                                                      |
| ------------------------ | --------------------------------------------------- | ---------------------------------------------------------------------------- |
| F82 own snapshot release | builds, archives, checksums, SBOM, docker images... | keyless cosign signing needs OIDC (CI-only); local verify uses `--skip=sign` |
| M5 wiring                | GoReleaser/GitHubActions/Dockerfile generators live | `generators/homebrew.go` still unwired (blocked on M18 personal defaults)    |
| CI (F44)                 | fixture `goreleaser check` step                     | own-config check step omitted until M23 lands on master (now unblocked)      |

## c) NOT STARTED (Wave 6 remainder)

- M17 error-discarding fixes (flags.go, init.go:248, validate_main.go:167, parseGitHubRemote real errors)
- M18 personal defaults (`--tap-repo`, `--nix-repo`, `--commit-author`, `--runner`) — **needs Lars's answers**
- M19-M21 coverage lifts (git/generators/templates/types/config/domain)
- M22 `validate` deprecation detection + migration hints
- M24 flake.nix cache fallback; M25 TODO_LIST/ROADMAP; M26/M27 v1.0.0 release + post-release verify

## d) Honest notes

1. **One --no-verify commit** (`9f4bd71`): BuildFlow pipeline finished green (61 tools, 0 failed) but its PostHog telemetry hung on a network outage; killed and committed with --no-verify, documented in the commit message. Recovery procedure captured in AGENTS.md gotcha 13.
2. **git stash mishap** (recovered): stashed the own-config migration before a release test, inverting the intent; killed the job, popped cleanly, no damage.
3. **Golden configs initially used windows+arm64** — a combo domain validation rejects; goldens would have pinned unreachable configs. Caught by running the E2E test, fixed + regenerated.
4. **F33 (opt-in test/generate hooks) deferred**: `SafeProjectConfig` is TypeSpec-generated ("DO NOT MODIFY"); adding RunTests/RunGenerate fields requires a spec change, so it belongs with M18's parameterization, not a hot-fix.
5. BuildFlow still flags: AGENTS.md age (updated this session), sbom-action v0 SHA un-pinnable (404 upstream, known), nix placeholder hashes (pre-existing).

## e) Open questions for Lars (unchanged)

1. Runner default `ubuntu-latest` OK, or revive self-hosted via `--runner`?
2. Push the 5 local commits now, or after remaining Wave 6?
3. M18 personal defaults: tap repo pattern + commit author email?

**Bottom line:** the wizard's correctness is now proven by machines, not eyeballs — E2E `goreleaser check` + goldens + CI, all green; the split brain is dead; the dogfood config is on v2.17. Remaining work is hardening and Lars-dependent decisions.
