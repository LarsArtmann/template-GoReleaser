# GoReleaser v2.17 — A Practical Guide

**Audience:** you, generating releases for Go projects (and maintaining this wizard).
**Verified against:** goreleaser v2.17.1 (docs + `goreleaser check` behavior), 2026-08-17.
**Companion:** the wizard encodes everything here — `goreleaser-wizard generate` produces a config that passes `goreleaser check` with zero deprecations and no environment variables.

---

## 1. Mental model: the pipeline

GoReleaser runs your config through ordered phases. Understanding the order explains most behavior:

```text
before hooks ─> builds ─> package (archives, nfpm, ...) ─> checksum ─> sign/SBOM
            ─> publish (release, dockers_v2, homebrew_casks, nix, scoop, ...)
            ─> metadata
```

Key consequences:

- **Docker images are built from artifacts, not source.** `dockers_v2` (and the old `dockers`) copy the binaries that the `builds` phase already produced into the Docker build context. Your Dockerfile must `COPY` the prebuilt binary — never rebuild inside the image.
- **`before.hooks` run before everything.** Keep them minimal (`go mod tidy`). Running `go test`/`go generate` here slows every release and surprises contributors; CI should test, the release pipeline should release.
- **Publishers are independent.** A failing cask publish does not stop the Docker publish — but it fails the release overall.

## 2. Config anatomy (what the wizard generates)

```yaml
version: 2 # required since v2

project_name: myapp # defaults to git repo name

before:
  hooks:
    - go mod tidy

builds:
  - id: myapp # referenced by dockers_v2/archives via `ids`
    main: ./cmd/myapp
    binary: myapp
    env: [CGO_ENABLED=0]
    goos: [linux, darwin]
    goarch: [amd64, arm64]

archives:
  - id: default
    format_overrides:
      - goos: windows
        formats: [zip] # note: `formats` (list), not the old `format` (string)

checksum:
  algorithm: sha256

snapshot:
  version_template: "{{ incpatch .Version }}-next"

changelog:
  use: github # use the GitHub API for the changelog

release:
  github:
    owner: octocat # literal values → local `goreleaser check` needs no env vars
    name: myapp
```

## 3. Template variables (evaluated by GoReleaser, not by you)

Inside GoReleaser string fields, `{{ ... }}` is evaluated at release time with this context:

| Variable                                     | Meaning                                 |
| -------------------------------------------- | --------------------------------------- |
| `.Version`                                   | clean tag without `v` (e.g. `1.2.3`)    |
| `.Tag`                                       | full tag (e.g. `v1.2.3`)                |
| `.Major`, `.Minor`, `.Patch`                 | version components                      |
| `.Os`, `.Arch`, `.Arm`                       | target platform of the current artifact |
| `.ProjectName`                               | from `project_name`                     |
| `.Env.FOO`                                   | environment variable (fails if unset!)  |
| `.IsNightly`                                 | true for nightly builds                 |
| `.Date`, `.Commit`, `.FullCommit`, `.GitURL` | git metadata                            |
| `.ArtifactName`, `.ConventionalFileName`     | artifact-scoped                         |

Functions include `incpatch`, `title`, `eq`, `if`, `time` — full list: <https://goreleaser.com/customization/templates/>.

**Gotcha:** `{{.Env.FOO}}` hard-fails when `FOO` is unset (that is why `goreleaser check` used to fail locally before the wizard baked literal owner/repo in). Prefer literal values for things known at generate time.

## 4. `dockers_v2` deep-dive

The modern Docker section. One entry builds a multi-platform image via buildx:

```yaml
dockers_v2:
  - id: myapp
    images:
      - ghcr.io/owner/myapp # tags appended from `tags`
    tags:
      - "{{.Tag}}"
      - "v{{.Major}}"
      - "{{if not .IsNightly}}latest{{end}}"
    dockerfile: Dockerfile
    platforms:
      - linux/amd64
      - linux/arm64
    annotations: # become OCI image labels
      org.opencontainers.image.title: "{{.ProjectName}}"
      org.opencontainers.image.revision: "{{.FullCommit}}"
```

Facts that matter:

- **The Dockerfile must be prebuilt-pattern.** GoReleaser places each platform's binary in the build context under `$TARGETPLATFORM`:

  ```dockerfile
  FROM scratch
  ARG TARGETPLATFORM
  COPY $TARGETPLATFORM/myapp /usr/bin/myapp
  ENTRYPOINT ["/usr/bin/myapp"]
  ```

- **Publish phase.** `dockers_v2` images are pushed during the publish phase; tags like `latest` are guarded with `{{if not .IsNightly}}` so nightlies never steal `latest`.
- **Snapshot mode** builds per-arch images with suffixes (`-amd64`, `-arm64`) and skips the manifest push.
- **`ids`** selects which `builds` entries feed the context; empty means all.
- Other knobs: `build_args`, `extra_files`, `labels`, `hooks`, `retry`, `sbom`, `flags`, `disable`.
- The old `dockers` (one entry per arch, `build_flag_templates`, `--platform` flags) and `docker_manifests` are **deprecated** in favor of this.

## 5. Homebrew: `brews` → `homebrew_casks`

`brews` (formula tap) is deprecated; `homebrew_casks` is the maintained path:

```yaml
homebrew_casks:
  - name: myapp
    binaries: [myapp] # binaries installed into PATH
    repository:
      owner: octocat
      name: homebrew-tap
    directory: Casks # default; formulas lived in Formula/
    commit_author:
      name: goreleaserbot
      email: bot@example.com
    hooks:
      post:
        install: | # remove quarantine bit for unsigned binaries
          if OS.mac?
            system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/myapp"]
          end
```

Casks can install plain binaries (`binaries:`) — they are not only for `.app` bundles. Unsigned binaries need the quarantine-removal hook or macOS users see "app is damaged".

## 6. Snapshot vs release vs nightly

| Mode     | Command                                 | Version      | Publishes?                              |
| -------- | --------------------------------------- | ------------ | --------------------------------------- |
| Check    | `goreleaser check`                      | —            | validates config only                   |
| Snapshot | `goreleaser release --snapshot --clean` | `X.Y.Z-next` | **no pushes**; docker images stay local |
| Release  | `goreleaser release --clean`            | from git tag | yes: GitHub release, images, taps       |
| Nightly  | scheduled CI + `nightly` settings       | timestamp    | yes, `.IsNightly` guards tags           |

`--clean` removes the `dist/` folder first — almost always what you want.

## 7. Signing, SBOMs, provenance

- **Keyless signing** (cosign + OIDC in CI) needs `id-token: write` permission and the cosign installer step; the wizard emits both when signing is enabled.
- **SBOMs** via `sboms: [{artifacts: archive}]` require `syft` in PATH at release time (`nix shell nixpkgs#syft`).
- **Attestations** (SLSA provenance) are GitHub-Actions-only and need `actions: read`, `id-token: write`.

## 8. Root `retry` and other modern keys

- `retry` moved to the root level (per-section retry is gone).
- `archives.format_overrides.format` (string) → `formats` (list).
- Full deprecation list: <https://goreleaser.com/deprecations/>.

## 9. Recipe: release a project from zero

```bash
goreleaser-wizard init                                   # or `generate --config cfg.yaml`
goreleaser check                                         # must be green, zero deprecations
git tag v0.1.0 && git push --tags                        # release trigger
# CI: goreleaser release --clean (GITHUB_TOKEN is enough for public repos)
```

Docker pushes need `DOCKER_USERNAME`/`DOCKER_PASSWORD` secrets; ghcr.io works with `GITHUB_TOKEN` plus `packages: write`.

---

_This guide describes what the wizard automates. If the wizard and this guide disagree, the golden files in `cmd/goreleaser-wizard/testdata/golden/` are the tiebreaker — and a bug report._
