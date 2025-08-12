# Report about missing/under-specified/confusing information

Date: 2025-08-12T07:13:33+02:00

I was asked to perform:
Comprehensive assessment and improvement of the GoReleaser template project with brutal honesty about issues and implementation of missing features.

I was given these context information's:
- Instructions to be brutally honest and never lie
- List of libraries to use (gin, viper, templ, htmx, samber/*, sqlc, etc.)
- Architecture patterns to follow (DDD, CQRS, Event Sourcing, etc.)
- Requirement to focus on customer value and avoid ghost systems
- Multiple GitHub issues to address

I was missing these information:
- Clear indication of which project I was working on initially (confused with Kernovia)
- Whether the integration tests are expected to work with the current architecture
- Specific database schema requirements for sqlc integration
- Clear examples of how htmx should be integrated in a Go template project
- Whether TypeSpec is actually needed in a template project

I was confused by:
- The mismatch between claimed libraries and actual usage patterns
- Why integration tests have persistent timeout issues even after fixes
- Whether this template should include full example implementations or just structure
- The scope of what belongs in a "template" vs a full application
- Why there are references to Kernovia in the instructions for a GoReleaser template

What I wish for the future is:
- Clear project context verification at the start of each session
- Explicit scope definition for template projects vs applications
- Better separation between template improvements and feature additions
- Clearer expectations about test success criteria
- More specific guidance on which libraries are must-have vs nice-to-have

Best regards,
Claude