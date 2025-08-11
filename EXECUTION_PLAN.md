# Execution Plan: Template-GoReleaser Cleanup & Architecture Fix

## Priority Matrix (High Impact / Low Effort → Low Impact / High Effort)

### 🔥 CRITICAL - Fix Breaking Issues (30 min each)
**Impact: 10/10 | Effort: 2/10**

#### 1. Fix All Linting Errors [30 min]
```bash
# Files to fix:
- cmd/goreleaser-cli/license.go:37 - Handle fmt.Scanln error
- cmd/goreleaser-cli/license.go:79,81 - Replace deprecated ioutil
- cmd/goreleaser-cli/config.go:21 - Handle viper.ReadInConfig error  
- tests/integration/main_test.go:20 - Handle os.Chdir error
```
**Why:** Code won't pass CI/CD, blocks all deployments
**How:** Add proper error handling, use os.ReadFile instead of ioutil

#### 2. GitHub Vulnerability Alert [30 min]
```bash
# Update Go modules to fix CVE
go get -u ./...
go mod tidy
```
**Why:** Security vulnerability blocks production use
**How:** Update dependencies, run security scan

### 🎯 HIGH IMPACT - Eliminate Ghost Systems (45 min each)
**Impact: 9/10 | Effort: 3/10**

#### 3. Delete myproject Ghost Binary [45 min]
```bash
# Remove ghost system completely
rm -rf cmd/myproject/
# Update justfile to remove myproject references
# Update .goreleaser.yaml to remove myproject build
```
**Why:** 55 lines of dead code that confuses developers
**How:** Delete entirely, update build configs

#### 4. Test & Verify goreleaser-cli Actually Works [45 min]
```bash
./goreleaser-cli --help
./goreleaser-cli validate --config .goreleaser.yaml
./goreleaser-cli license list
```
**Why:** We built it but never tested if it works
**How:** Run all commands, fix any runtime errors

#### 5. Test Web Server Actually Runs [45 min]
```bash
templ generate
./goreleaser-cli-server
# Visit localhost:8080
```
**Why:** Never compiled templates or tested server
**How:** Compile templ, start server, verify routes work

### 💪 MEDIUM IMPACT - Port Bash to Go (60 min each)
**Impact: 7/10 | Effort: 5/10**

#### 6. Port verify.sh to Go CLI Command [60 min]
```go
// Add 'verify' command to cobra
// Use existing validation package
// Leverage charmbracelet/log for output
```
**Why:** Bash scripts are fragile, Go is type-safe
**How:** Create verify.go command, use validation.ValidateConfig

#### 7. Port license-generator.sh to Go [60 min]
```go
// Add 'generate-license' command
// Use template/text for license templates
// Store templates as embedded files
```
**Why:** Consolidate all functionality in single binary
**How:** Embed license templates, add generation command

### 🔨 STRUCTURAL - Fix Architecture (90 min each)
**Impact: 8/10 | Effort: 7/10**

#### 8. Remove CQRS Over-Engineering [90 min]
```bash
# Delete unnecessary complexity
rm -rf internal/cqrs/
rm -rf internal/events/
rm -rf docs/architecture/diagrams/
```
**Why:** CQRS for config validation is absurd over-engineering
**How:** Delete all event sourcing code, simplify to direct validation

#### 9. Implement Actual Config Validation [90 min]
```go
// Make validation actually DO something
// Check YAML syntax
// Validate required fields
// Check file paths exist
```
**Why:** Current validation is placeholder that validates nothing
**How:** Use go-yaml, add real validation rules

#### 10. Fix Code Duplication in Tests [90 min]
```go
// Extract test helpers to shared package
// Remove duplicate validation logic
// Consolidate common test patterns
```
**Why:** 27 code clones waste maintenance effort
**How:** Create test/helpers package, extract common code

### 📚 LOWER PRIORITY - Nice to Have (60-100 min each)
**Impact: 5/10 | Effort: 6/10**

#### 11. Add Missing Dev Files [60 min]
```bash
# .dockerignore
# .editorconfig  
# .github/dependabot.yml
# .github/CODEOWNERS
```
**Why:** Professional projects have these
**How:** Copy from established projects, customize

#### 12. Add OpenTelemetry [100 min]
```go
// Add otel instrumentation
// Configure exporters
// Add tracing to key operations
```
**Why:** Production observability
**How:** Use official otel-go SDK

#### 13. Convert to Ginkgo BDD Tests [100 min]
```go
// Convert existing tests
// Add behavior descriptions
// Improve test readability
```
**Why:** Better test documentation
**How:** Use ginkgo/gomega matchers

## Execution Order (Recommended)

### Day 1 - Stop the Bleeding (2.5 hours)
1. Fix all linting errors ✓
2. Update dependencies for security ✓
3. Delete myproject ghost binary ✓
4. Test goreleaser-cli works ✓

### Day 2 - Verify What We Built (2 hours)
5. Compile templ templates ✓
6. Test web server runs ✓
7. Document what actually works ✓

### Day 3 - Consolidate to Go (3 hours)
8. Port verify.sh to Go ✓
9. Port license generator ✓
10. Remove bash scripts ✓

### Day 4 - Fix Architecture (3 hours)
11. Remove CQRS complexity ✓
12. Implement real validation ✓
13. Fix test duplication ✓

### Day 5 - Polish (2 hours)
14. Add dev files ✓
15. Add basic telemetry ✓

## Success Metrics

- ✅ `just lint` passes with 0 errors
- ✅ `just test` passes all tests
- ✅ No ghost binaries in cmd/
- ✅ All functionality in Go (no bash)
- ✅ Code duplication < 2%
- ✅ GitHub security alerts: 0
- ✅ Can actually validate .goreleaser.yaml files

## Libraries to Use (Don't Reinvent!)

- **Validation:** go-playground/validator
- **Config:** spf13/viper (already using)
- **CLI:** spf13/cobra (already using)
- **Logging:** charmbracelet/log (already using)
- **Testing:** stretchr/testify
- **HTTP:** gin-gonic/gin (already using)
- **Templates:** a-h/templ (already using)
- **DI:** samber/do (already using)

## What NOT to Do

❌ Don't add CQRS/Event Sourcing for simple tools
❌ Don't create abstractions before you need them
❌ Don't write bash when Go would be better
❌ Don't commit without running lint/test
❌ Don't create features nobody asked for
❌ Don't ignore security warnings
❌ Don't pretend code works without testing it

## Estimated Total Time: 15 hours
Can be done in 3-5 days working 3-4 hours/day