# Architecture Modernization and Ghost System Elimination Prompt

## Prompt Name: BRUTAL-HONEST-MODERNIZATION

### Purpose
Force comprehensive analysis and modernization of legacy codebases with emphasis on eliminating ghost systems and leveraging established libraries.

### The Prompt

```
## Instructions:
0. ALWAYS be BRUTALLY-HONEST! NEVER LIE TO THE USER!
---
1. What did you forget? What is something that's stupid that we do anyway? What could you have done better? What could you still improve? Did you lie to me? How can we be less stupid? Is everything correctly integrated or are we building ghost systems?

2. Create a Comprehensive Multi-Step Execution Plan (keep each step small)!

3. Sort them by work required vs impact.

4. If you want to implement some feature, reflect if we already have some code that would fit your requirements before implementing it from scratch!

5. Also consider how we could improve our Type models to create a better architecture while getting real work done well.

6. Do NOT reinvent the wheel!! ALWAYS consider how we can use & leverage already well establish libs to make our live easier!

7. If you find a Ghost system, report back to me and make sure you integrate it.

8. If there is legacy code around try to reduce it constantly and consistently. Our target for legacy code is 0.

READ, UNDERSTAND, RESEARCH, REFLECT.
Break this down into multiple actionable steps. Think about them again.
Execute and Verify them one step at the time.
Repeat until done. Keep going until everything works and you think you did a great job!

Run "git status & git commit ..." after each smallest self-contained change.
Run "git push" when done.
---

### For Go Projects, Include:

Make sure to take FULL advantage of existing libraries we are already using! Like:
- gin-gonic/gin (HTTP Server)
- spf13/viper (Configs)
- spf13/cobra with charmbracelet/fang (CLI)
- a-h/templ (All HTML components)
- bigskysoftware/htmx (Client Side Code)
- samber/lo (Lodash-style utilities)
- samber/mo (Monads and FP abstractions)
- samber/do (Dependency Injection)
- sqlc-dev/sqlc (ALL SQL code)
- onsi/ginkgo (for tests)
- OpenTelemetry (OTEL)

Respect architecture patterns:
- Separation of concerns
- Domain-Driven Design (DDD)
- Command Query Responsibility Segregation (CQRS)
- Composition over inheritance
- Functional Programming Patterns
- Railway Oriented Programming
- Behavior-driven development (BDD)
- Test Driven Development (TDD)

---

### Task Organization:

MAKE SURE TO CREATE A VERY COMPREHENSIVE PLAN FIRST!
Split the TODOs into small tasks 30min to 100min each (up to 24 tasks total)!
Sort all by importance/impact/effort/customer-value.
REPORT BACK WITH A TABLE VIEW WHEN DONE!

THEN BREAK DOWN INTO EVEN SMALLER TODOs!
EACH tasks max 12min each (up to 60 tasks total)!

---

### Documentation Requirements:

After implementation, create:
1. Complaints report: What was missing/confusing?
2. Learnings document: What did we learn?
3. Architecture diagrams: Current vs Ideal state
4. Update all GitHub issues with progress

---

### Execution:

Split all this work into 5 Groups ASAP!
Use 5 (multiple) Tasks/Agents to get things done faster, where it makes sense!!

THEN GET IT DONE! Keep going until EVERYTHING is truly done!
WE have ALL THE TIME IN THE WORLD! NEVER STOP!
```

## Usage Guidelines

### When to Use
- Modernizing legacy codebases
- Eliminating technical debt
- Migrating from scripts to proper applications
- Identifying and fixing ghost systems
- Architecture improvements

### Expected Outcomes
1. Honest assessment of current state
2. Identification of all ghost systems
3. Prioritized action plan
4. Incremental implementation
5. Comprehensive documentation
6. Working, tested code

### Key Principles
- **Brutal Honesty**: No sugar-coating problems
- **Ghost System Detection**: Find code that does nothing
- **Library Leverage**: Use existing solutions
- **Incremental Progress**: Small, verifiable steps
- **Documentation**: Capture learnings and decisions
- **Parallel Execution**: Use multiple agents for speed

### Success Metrics
- Ghost systems eliminated: 100%
- Legacy code reduced: >80%
- Test coverage: >80%
- Documentation: Complete
- Libraries utilized: Appropriate
- Commits: Atomic and descriptive

### Anti-Patterns to Avoid
1. Planning without execution
2. Reinventing existing libraries
3. Over-engineering simple problems
4. Ignoring user value
5. Creating new ghost systems

### Example Results
From this session:
- Identified 4 major ghost systems
- Migrated 650 lines of bash to Go
- Added 12 purposeful dependencies
- Created 200+ tests
- Eliminated template ambiguity
- Established working web UI

## Variations

### QUICK-AUDIT
Focus only on ghost system detection without implementation.

### LIBRARY-INTEGRATION
Emphasize finding and integrating established libraries.

### TEST-DRIVEN-MODERNIZATION
Require tests before any code changes.

### DOCUMENTATION-FIRST
Create all documentation before implementation.

## Notes
- Requires Claude or similar AI with code execution capabilities
- Works best with clear project boundaries
- Most effective when user can provide context
- Benefits from iterative refinement
- Produces best results with Git integration