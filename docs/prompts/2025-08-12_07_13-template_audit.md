# Template Project Audit and Improvement Prompt

## Name: Complete Template Audit with Ghost System Detection

## Purpose
Perform a comprehensive audit of a template project to identify and fix issues, remove ghost systems, and ensure all claimed functionality actually works.

## Prompt

```
I need you to perform a BRUTAL and HONEST audit of this template project. 

FIRST, verify the project context:
1. Run `pwd && git remote -v` to confirm location
2. Read README.md to understand project purpose
3. Check go.mod for actual dependencies
4. Run `git status` to see current state

THEN, identify ALL ghost systems:
1. Check for binaries in git (should be in .gitignore)
2. Find libraries in go.mod that aren't actually used properly
3. Identify config files without corresponding implementation
4. Look for tests that don't actually pass
5. Find documentation that doesn't match reality

For EACH issue found:
1. Explain what's wrong and why it's a problem
2. Determine if it should be fixed or removed
3. Provide the exact commands/code to fix it
4. Verify the fix actually works

REQUIREMENTS:
- Be BRUTALLY HONEST - no sugar-coating
- Test everything you claim works
- Remove or integrate all ghost systems
- Focus on customer value, not complexity
- Use established libraries properly (don't reinvent)
- Make small, focused commits

Libraries to properly integrate if claimed:
- gin-gonic/gin for HTTP
- samber/lo for functional programming
- samber/mo for error handling
- samber/do for dependency injection
- sqlc for database operations
- templ for HTML templates
- htmx for frontend interactivity

After fixes, verify:
1. All tests pass (with reasonable timeouts)
2. No binaries in git
3. All claimed libraries are actually used
4. Documentation matches implementation
5. No ghost systems remain

Report format:
1. Issues Found (with severity)
2. Actions Taken
3. Remaining Problems
4. Recommendations

Be prepared to explain WHY each decision was made and HOW it improves the template.
```

## Usage Instructions

Use this prompt when:
- Taking over a template project from someone else
- Performing periodic template maintenance
- Before releasing a new template version
- When users report disconnects between docs and reality

## Expected Outcomes

1. Complete list of ghost systems and issues
2. Fixed test suite that actually passes
3. Proper library integration
4. Accurate documentation
5. Clean git repository (no binaries)
6. Honest assessment of template quality

## Key Questions to Answer

- What's claimed vs what's real?
- What works vs what's broken?
- What's integrated vs what's a ghost system?
- What adds value vs what adds complexity?
- What should be in a template vs an application?

## Success Criteria

- [ ] All tests pass within 30 seconds
- [ ] No compiled binaries in git
- [ ] All libraries in go.mod are properly used
- [ ] Documentation accurately describes functionality
- [ ] No ghost systems remain
- [ ] Template can be used immediately for new projects