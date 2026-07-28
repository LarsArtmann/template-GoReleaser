# Dedup Acceptance Log

Record of clone groups evaluated and intentionally accepted, so the next session
can skip re-evaluating them.

---

## Clone Group: "command setup boilerplate" (assignment + defer)

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

**Reason accepted:**

Structural similarity with entirely unique values. Each call site reads
different flags, uses a different panic-recovery context string, and prints a
different header. This is idiomatic Cobra command-function boilerplate, not
copy-paste. A higher-level abstraction (e.g. a `setupCommand` helper) would need
to accept variable flag sets and return them under different names — more
parameters than the duplicated code has lines, fighting the Cobra pattern rather
than improving it.

The underlying `cmd.Flags().GetBool/GetString` error-discarding was already
extracted into `getBoolFlag`/`getStringFlag` helpers in `flags.go`.
