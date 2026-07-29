# Dedup Acceptance Log

Record of clone groups evaluated and intentionally accepted, so the next session
can skip re-evaluating them.

Each entry records the **location**, the **shape** of the duplication, and the
**reason** it was kept. Reasons fall into a small set:

- _Idiomatic Go_ — standard language or library pattern; extracting would obscure.
- _Different rules, same shape_ — the bodies are similar but encode distinct
  business rules that must evolve independently.
- _Generic refactor would add more code than it removes_ — see "abstraction
  trade-off" in the dedup skill.
- _Split-brain_: parallel implementations in different packages that should be
  reconciled by a future task (out of scope here).

---

## Group: "command setup boilerplate" (assignment + defer)

**Locations:**

- `cmd/goreleaser-wizard/generate.go:14-19`
- `cmd/goreleaser-wizard/validate.go:48-54`
- `cmd/goreleaser-wizard/init.go:17-24`

**Pattern:**

```go
defer recoverFromPanic("<unique command name>")

<flag> := getBoolFlag(cmd, "<unique flag>")

printCommandHeader("<unique header>")
```

**Reason:** Structural similarity with entirely unique values. Each call site reads
different flags, uses a different panic-recovery context string, and prints a
different header. This is idiomatic Cobra command-function boilerplate, not
copy-paste. A higher-level abstraction (e.g. a `setupCommand` helper) would need
to accept variable flag sets and return them under different names — more
parameters than the duplicated code has lines, fighting the Cobra pattern rather
than improving it.

The underlying `cmd.Flags().GetBool/GetString` error-discarding was already
extracted into `getBoolFlag`/`getStringFlag` helpers in `flags.go`.

---

## Group: "empty-input early-return in template escapers"

**Locations:**

- `internal/validation/template_escaping.go:47` (EscapeYAML)
- `internal/validation/template_escaping.go:79` (EscapeShell)
- `internal/validation/template_escaping.go:100` (EscapeGitHubActions)
- `internal/validation/template_escaping.go:121` (EscapeJSON — returns `""`, not `""`)
- `internal/validation/template_escaping.go:140` (EscapeDockerLabel)

**Pattern:**

```go
if value == "" {
    return ""  // EscapeJSON returns `""` (the JSON null literal)
}
```

**Reason:** 2-line guard clause with **unique fallback values** for each escaper
(JSON returns the literal 2-character string `""`). `SanitizeInput("")` already
returns `""` so the early return is purely a skip-the-work optimization. An
extracted `earlyReturnOrSanitize(value, fallback)` would shave 1 line per site
(10 lines total) at the cost of a 4-line helper and a harder-to-follow call
sequence — net negative. The duplication is the shape of the canonical
"early-return then transform" pattern, not maintenance burden.

---

## Group: "shell-quote escape" vs "JSON escape" `strings.ReplaceAll`

**Locations:**

- `internal/validation/template_escaping.go:92` — POSIX shell: `'` → `''`
- `internal/validation/template_escaping.go:133` — JSON control char: `\t` → `\t`

**Reason:** Same library call, completely different escape rules with different
surrounding context. Each call site operates on a different escaping scheme for
a different output format. Not duplication — coincidence.

---

## Group: "idempotent fluent builder setters"

**Locations:**

- `internal/domain/errors.go:144` — `de.Suggestion = suggestion`
- `internal/types/validation.go:395` — `ve.Suggestion = suggestion`
- `internal/types/validation.go:444` — `vw.Suggestion = suggestion`
- `internal/types/validation.go:381` — `ve.Context = context`
- `internal/types/validation.go:430` — `vw.Context = context`
- `internal/types/validation.go:388` — `ve.Details = details`
- `internal/types/validation.go:437` — `vw.Details = details`

**Reason:** Each setter is a 4-line method that mutates exactly one field and
returns the receiver. Three types × three setters = 9 methods, 36 lines. The
public API (`ve.WithContext("...")`) is called from ~50+ sites across
`cmd/goreleaser-wizard/` and `internal/`. Replacing with Go generics requires:

1. Renaming the public methods to lowercase (`setContext`).
2. Adding new public wrappers that delegate to the generic helpers.
3. Net change: 9 methods × ~5 lines of generic helper = ~45 lines added.

No reduction, more layers, public API churn. **DomainError.WithContext is also
notably different** — it returns a new immutable copy rather than mutating, so
it would need a different generic helper anyway. Different patterns in
different packages = intentional, not duplication.

---

## Group: "command/job entry preamble"

**Locations:**

- `cmd/goreleaser-wizard/jobs.go:332,417,450` — 3 sites, charm.Logger
- `cmd/goreleaser-wizard/generators/github_actions.go:125`
- `cmd/goreleaser-wizard/generators/goreleaser.go:40`
- `cmd/goreleaser-wizard/generators/homebrew.go:171`

**Pattern (jobs.go):**

```go
j.logger.Info("XYZ")
if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
    return err
}
```

**Reason:** Tried to extract a generic `logAndCheckCtx` helper. Blocked by
heterogeneous logger types: `*charm.Logger.Info(any, ...any)` vs
`domain.Logger.Info(string, ...any)` vs `generators.Logger.Info(...)`. A Go
generic with constraint `interface{ Info(string, ...any) }` rejects charm's
`Info(any, ...any)`. A type-parameterised function value (`func(string)`)
requires wrapping the call site with `func(msg string) { j.logger.Info(msg) }`
— that closure is **longer than the 3-line preamble it replaces**.

The preamble is the canonical "log + check ctx" idiom for any cancellable unit
of work. Accept as idiomatic.

---

## Group: "stderr-style user-facing output headers"

**Locations:**

- `cmd/goreleaser-wizard/validation_display.go:9` — `fmt.Println("📋 Validation Summary:")`
- `cmd/goreleaser-wizard/validation_fixes.go:22` — `fmt.Println("🔧 Attempting to fix common issues...")`

**Reason:** Different messages, different files. Single-line `fmt.Println`
calls with unique strings. Not duplication.

---

## Group: "ExecuteTemplate buffer pattern across packages"

**Locations:**

- `cmd/goreleaser-wizard/generators/template_utils.go:65-82` —
  `executeTemplateWithError` (uses domain error wrapping)
- `cmd/goreleaser-wizard/jobs.go:824-834` —
  `executeTemplate` (uses plain `fmt.Errorf`)

**Reason:** Two parallel implementations of `var output bytes.Buffer; tmpl.Execute(&output, data); return output.Bytes(), nil`
in **different packages with intentionally different error wrapping strategies**
(generators returns `domain.ErrTemplateRendering`; jobs uses `fmt.Errorf`).
The `generators` package is currently **not imported anywhere** — it's effectively
dead code or slated for migration. Reconciling it requires deciding which
package owns template execution (probably the `generators` package should go
away, with jobs.go promoted to be the canonical entry point). Out of scope for
dedup; flagged for future work.

---

## Group: "split-brain project-name pattern"

**Locations:**

- `internal/domain/validators.go:14` — `^[a-zA-Z0-9][a-zA-Z0-9\-._]*[a-zA-Z0-9]$` (1–63 chars, start AND end alphanumeric)
- `internal/validation/basic.go:25` — `^[a-zA-Z0-9][a-zA-Z0-9._-]{0,49}$` (1–50 chars, start alphanumeric)

**Reason:** Same name, **different rules**. The domain version is stricter on
trailing characters and uses a longer max length; the validation version uses
a 50-char max (matches user input form) and allows `.`/`-` at the end (handled
separately by the `strings.HasSuffix` checks). Different consumers
(`tui_wizard.go` vs `init_test.go`) need different rules. This is a split-brain
that _should_ be reconciled — but reconciliation is a product decision (do we
tighten the validation rules or loosen the domain rules?), not a dedup task.

---

## Group: "name-empty guard in domain validators"

**Locations:**

- `internal/domain/validators.go:56` (ValidateProjectName)
- `internal/domain/validators.go:97` (ValidateBinaryName)

**Pattern:**

```go
if name == "" {
    return errors.New("<noun> cannot be empty")  // unique message per validator
}
```

**Reason:** 2-line guard with unique error messages. Extracting a
`requireNonEmpty(name, noun)` helper saves 1 line per site (2 lines total) at
the cost of indirection. Net negative. Canonical guard-clause idiom.

---

## Group: "early return on empty section" (after `isEmpty` extraction)

**Locations:**

- `cmd/goreleaser-wizard/validation_display.go:73` — `if isEmpty(warnings)`
- `cmd/goreleaser-wizard/validation_display.go:91` — `if isEmpty(recommendations)`

**Reason:** Already extracted into generic `isEmpty[T any](items []T) bool` helper
in `validation_display.go`. The art-dupl hit is on the _call sites_ (each
passing a different slice type), not on the function body. Two sites of a
3-line helper call is not worth further extraction.

---

## Group: "if err != nil idiomatic Go boilerplate" (simple_filesystem.go)

**Locations:**

- `cmd/goreleaser-wizard/simple_filesystem.go:16,30,42,77,90,115,127,141,149,158,175`
  (all 4-line `if err != nil { return …wrapFSError(...) }` shapes)

**Reason:** Extracted the value-transforming part (`wrapFSError(op, path, err)`
helper), which removed the `fmt.Errorf("failed to X %q: %w", path, err)`
duplication. The remaining `if err != nil { return … }` shape is the canonical
Go error-handling idiom — extracting it via a closure or generic would obscure
without saving meaningful lines. Each remaining hit has a different operation
name, path, and return-value shape.

---

## Summary

Extracted (real harmful duplication, this session):

- `validateEnum(itemName, value, isValid)` in `internal/domain/validation_utils.go` —
  unified 8 single-value enum validators
  (`CGOStatus`/`SigningLevel`/`DockerSupport`/`BuildLevel`/`Platform`/`ConfigState`/`DockerRegistry`/`GitProvider`).
  Previously used inconsistent `NewValidationError`/`fmt.Errorf`.
- `wrapFSError(op, path, err)` in `cmd/goreleaser-wizard/simple_filesystem.go` —
  unified 12 file-system error-wrapping sites.
- `isEmpty[T any](items []T) bool` in
  `cmd/goreleaser-wizard/validation_display.go` — generic empty-slice guard.
- `stripVersionPrefix(version)` in `internal/git/commands.go` — unified
  `GetMajorVersion`/`IncPatchVersion` prefix handling.
- `checkReservedName(name, label)` in `internal/domain/validators.go` —
  unified case-insensitive reserved-name check in `ValidateProjectName`/`ValidateBinaryName`.
- `mapToStrings[T any](items []T, convert func(T) string) []string` in
  `cmd/goreleaser-wizard/jobs.go` — unified 4 typed-enum-slice → []string
  conversions in `prepareGoReleaserData` and `prepareGitHubActionsData`
  (`Platforms`/`Architectures`/`BuildTags`/`Triggers`).

Final `art-dupl --sort total-tokens -t 2` report: **17 groups, all accepted
above**. Each remaining hit is documented with its reason — split-brain
parallel implementations, idiomatic Go shapes that resist generic extraction,
or 2-line guard clauses whose extraction would add more code than it removes.
