# Learnings from Template Improvement Session

Date: 2025-08-12T07:13:33+02:00

## Key Learnings

### 1. Always Verify Project Context First
**Mistake**: Started working on Kernovia code in a GoReleaser template project.
**Learning**: Always run `pwd && git remote -v` immediately upon starting work.
**Action**: Check project README and purpose before making any changes.

### 2. Templates Should Remain Generic
**Mistake**: Added project-specific code (pkg/types/generated) to a template.
**Learning**: Template projects should provide structure, not implementation.
**Action**: Keep templates minimal and extensible, not opinionated about specific features.

### 3. Binary Files Don't Belong in Git
**Mistake**: Left compiled binaries (goreleaser-cli) in version control.
**Learning**: Always check for and exclude binary files from git.
**Action**: Add appropriate patterns to .gitignore immediately.

### 4. Test What You Claim
**Mistake**: Claimed "all tests passing" when they were timing out.
**Learning**: Actually run tests and verify results, don't assume.
**Action**: Use shorter timeouts and verify test completion status.

### 5. Ghost Systems Are Real
**Identified Ghost Systems**:
- Binaries in git that should be built
- Test scripts that don't complete
- Libraries in go.mod but not fully utilized
- Configuration files without implementation

**Learning**: Regularly audit for disconnected components.
**Action**: Remove or integrate ghost systems immediately.

### 6. Functional Programming Patterns Work Well in Go
**Success**: samber/lo and samber/mo significantly improved code clarity.
**Learning**: Functional patterns reduce boilerplate and improve readability.
**Action**: Use functional libraries consistently throughout the codebase.

### 7. Context Cancellation Prevents Hangs
**Problem**: Tests hanging due to goroutine leaks.
**Solution**: Use context.Context with timeouts everywhere.
**Learning**: Always use context for cancellable operations.
**Action**: Make context the first parameter in all async operations.

### 8. Library Integration Should Be Complete
**Problem**: Libraries mentioned but not fully integrated.
**Learning**: If you claim to use a library, use it properly throughout.
**Action**: Either fully integrate or remove from documentation.

## Best Practices Established

1. **Project Verification Checklist**:
   - Run `pwd && git remote -v`
   - Read README.md
   - Check go.mod for actual dependencies
   - Run tests to verify current state

2. **Template Project Guidelines**:
   - Keep generic and extensible
   - Provide structure, not implementation
   - Include examples, not production code
   - Document clearly what's included and why

3. **Testing Best Practices**:
   - Use short timeouts (5 seconds default)
   - Clean up resources in defer/teardown
   - Use context for cancellation
   - Verify actual results, don't assume

4. **Code Quality Standards**:
   - No binaries in git
   - Functional patterns where appropriate
   - Complete library integration
   - Remove ghost systems immediately

## Mistakes to Avoid

1. **Don't assume project context** - Always verify
2. **Don't claim false success** - Test and verify
3. **Don't leave ghost systems** - Integrate or remove
4. **Don't partially integrate libraries** - Go all in or don't
5. **Don't ignore test failures** - Fix them properly

## Recommended Workflow

1. **Start**: Verify project context
2. **Analyze**: Identify ghost systems and issues
3. **Plan**: Create comprehensive but realistic plans
4. **Execute**: Fix critical issues first
5. **Test**: Verify all changes work
6. **Document**: Update docs to match reality
7. **Commit**: Small, self-contained changes
8. **Review**: Check for ghost systems again

## Tools That Helped

- `git status` - Constant state awareness
- `go test -timeout` - Prevent hanging tests
- `grep -r` - Find actual usage patterns
- `gh issue list` - Stay aligned with project goals
- Context-based cancellation - Prevent resource leaks

## Future Improvements Needed

1. Better integration test architecture
2. Complete htmx example implementation
3. Full TypeSpec integration (if needed)
4. OpenTelemetry monitoring setup
5. Comprehensive BDD test suite with Ginkgo

## Conclusion

The session revealed significant disconnects between claimed functionality and actual implementation. The key learning is to always verify reality before making claims, and to maintain consistency between documentation, configuration, and code. Ghost systems are a real problem that accumulates over time and should be addressed immediately when found.