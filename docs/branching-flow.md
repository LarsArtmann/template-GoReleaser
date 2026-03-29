# GoReleaser-Wizard Branching Flow

## Overview

This document defines the Git branching and workflow strategy for GoReleaser-Wizard.

---

## Branch Types

| Prefix | Purpose | Example |
|--------|---------|---------|
| `feature/` | New features | `feature/docker-multi-stage` |
| `fix/` | Bug fixes | `fix/validation-timeout` |
| `refactor/` | Code refactoring | `refactor/workflow-split` |
| `docs/` | Documentation | `docs/api-reference` |
| `chore/` | Maintenance | `chore/update-deps` |
| `test/` | Test improvements | `test/add-integration` |

---

## Workflow Steps

### 1. Start Work

```bash
# Ensure master is up-to-date
git checkout master
git pull origin master

# Create feature branch
git checkout -b feature/my-feature
```

### 2. Make Changes

```bash
# Make your changes
# ... edit files ...

# Check status
git status
```

### 3. Commit

```bash
# Stage specific files (never git add .)
git add cmd/... internal/...
git add docs/...  # if applicable

# Commit with conventional message
git commit -m "feat(wizard): add multi-platform support

- Add arm64 platform detection
- Update templates for cross-platform builds
- Add validation for platform-specific settings

Closes #42"
```

### 4. Push & Create PR

```bash
# Push branch
git push origin feature/my-feature

# Create PR
gh pr create --title "feat: add multi-platform support" \
  --body "$(cat <<'EOF'
## Summary

- Add arm64 platform detection
- Update templates for cross-platform builds
- Add validation for platform-specific settings

## Testing

- [x] Unit tests pass
- [x] Integration tests pass
- [x] CI pipeline passes

EOF
)"
```

### 5. After PR Merge

```bash
# Switch back to master
git checkout master
git pull origin master

# Clean up local branch
git branch -d feature/my-feature
```

---

## Quality Gates

Before committing/pushing, ensure:

- [ ] Tests pass: `just test`
- [ ] Code formatted: `just fmt`
- [ ] Linting passes: `golangci-lint run`
- [ ] Architecture valid: `go-arch-lint check`
- [ ] CI passes: `just ci`

---

## Commit Message Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting (no code change)
- `refactor` - Code restructuring
- `perf` - Performance improvement
- `test` - Adding/updating tests
- `chore` - Maintenance tasks

### Example

```
feat(validation): add project name format validation

- Add regex pattern for project names
- Add unit tests for edge cases
- Update validation error messages

Closes #123
```

---

## Important Rules

1. **Never commit directly to master**
2. **Use `git add <files>` not `git add .`**
3. **One logical change per commit**
4. **Run CI before merging**
5. **Delete branches after merge**
6. **Use Conventional Commits format**
