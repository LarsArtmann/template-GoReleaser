# Report about missing/under-specified/confusing information

Date: 2025-08-10T11:09:33+02:00

I was asked to perform:
- Replace bash scripts with proper Go applications
- Integrate recommended libraries (gin, viper, templ, htmx, samber/*, etc.)
- Follow architecture patterns (DDD, CQRS, Event-Sourcing, etc.)
- Fix ghost systems and eliminate legacy code
- Create comprehensive documentation

I was given these context information's:
- List of recommended libraries to use
- Architecture patterns to follow
- Request to use existing libraries instead of reinventing
- Context7.com for library documentation
- Requirement for brutal honesty about issues

I was missing these information:
- **Clarity on business domain**: What is the actual business purpose of this template? Is it just a template or should it be a working tool?
- **User personas**: Who are the target users? Developers? DevOps? Both?
- **Deployment expectations**: Should this run as a service, CLI tool, or both?
- **Event sourcing specifics**: What events should be sourced in a GoReleaser template project?
- **CQRS boundaries**: Where should command/query separation apply in this context?
- **Performance requirements**: What are the actual performance targets?
- **Scale expectations**: Single user tool or multi-tenant service?

I was confused by:
- **Template vs Tool conflict**: The project is called "template-GoReleaser" but we're building a full application
- **Architecture pattern applicability**: Event-sourcing and CQRS seem over-engineered for a GoReleaser configuration tool
- **Library context**: Some recommended libraries (sqlc, event-sourcing) don't fit the use case
- **Ghost system definition**: Unclear boundary between "template code" and "ghost system"
- **Integration expectations**: Should this integrate with existing GoReleaser or replace it?

What I wish for the future is:
- **Clear product vision**: Define if this is a template, tool, or service
- **Architecture decision records**: Document why certain patterns are chosen
- **User journey maps**: Show how users will interact with the system
- **Integration specifications**: Clear boundaries with GoReleaser itself
- **Success metrics**: Define what "done" looks like beyond code completion
- **Domain model diagrams**: Visual representation of the business domain
- **Use case priorities**: Ranked list of user scenarios to implement

Best regards,
Claude (Claude Code Assistant)