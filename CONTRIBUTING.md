# Contributing to GoReleaser Template

Thank you for your interest in contributing to the GoReleaser Template project! This document provides guidelines and information for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Testing Guidelines](#testing-guidelines)
- [Commit Message Format](#commit-message-format)
- [Pull Request Process](#pull-request-process)
- [Code Review Guidelines](#code-review-guidelines)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

## Code of Conduct

This project follows a simple code of conduct:

### Our Pledge

We are committed to providing a friendly, safe, and welcoming environment for all contributors, regardless of level of experience, gender identity and expression, sexual orientation, disability, personal appearance, body size, race, ethnicity, age, religion, nationality, or other similar characteristics.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

- Being respectful and inclusive in discussions
- Focusing on constructive feedback
- Gracefully accepting constructive criticism
- Showing empathy toward other community members
- Being patient with newcomers

Examples of unacceptable behavior include:

- Harassment, discrimination, or hostile behavior
- Publishing others' private information without permission
- Any conduct that could reasonably be considered inappropriate

## How to Contribute

There are many ways to contribute to this project:

### Reporting Issues

- **Bug Reports**: Use the GitHub issue tracker to report bugs
- **Feature Requests**: Suggest new features or improvements
- **Documentation Issues**: Report unclear or missing documentation

When reporting issues:
- Use a clear and descriptive title
- Provide detailed steps to reproduce the problem
- Include relevant system information (Go version, OS, etc.)
- Add relevant logs or error messages

### Contributing Code

1. **Fork the Repository**
   ```bash
   gh repo fork LarsArtmann/template-GoReleaser
   cd template-GoReleaser
   ```

2. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number
   ```

3. **Make Your Changes**
   - Follow existing code style and patterns
   - Add tests for new functionality
   - Update documentation as needed

4. **Test Your Changes**
   ```bash
   just ci
   just validate
   ```

5. **Commit Your Changes**
   - Follow the commit message format below
   - Make atomic commits with clear messages

6. **Push and Create Pull Request**
   ```bash
   git push origin feature/your-feature-name
   gh pr create --title "feat: add new feature" --body "Description of changes"
   ```

### Contributing Documentation

- Improve existing documentation
- Add examples and use cases
- Fix typos and clarify instructions
- Update outdated information

## Development Setup

### Prerequisites

- Go 1.24+
- Git
- Docker (for container builds)
- Just (task runner)

### Initial Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/LarsArtmann/template-GoReleaser.git
   cd template-GoReleaser
   ```

2. **Install Development Tools**
   ```bash
   just install-tools
   ```

3. **Set Up Environment**
   ```bash
   just setup-env
   # Edit .env file with your configuration
   ```

4. **Initialize Project**
   ```bash
   just init
   ```

5. **Validate Setup**
   ```bash
   just validate
   ```

### Development Workflow

```bash
# Start development
just build        # Build the project
just test         # Run tests
just lint         # Run linters
just fmt          # Format code

# Full CI pipeline
just ci           # Run complete validation

# Validate configurations
just validate     # Validate GoReleaser configs
just validate-strict  # Run strict validation
```

### Available Commands

Run `just --list` to see all available commands or `just help` for detailed descriptions.

## Testing Guidelines

### Test Structure

- **Unit Tests**: Test individual functions and components
- **Integration Tests**: Test component interactions
- **Validation Tests**: Test configuration validation

### Writing Tests

1. **Use Table-Driven Tests** for multiple test cases:
   ```go
   func TestFunction(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected string
           wantErr  bool
       }{
           // test cases
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // test logic
           })
       }
   }
   ```

2. **Use testify/assert** for assertions:
   ```go
   import "github.com/stretchr/testify/assert"
   
   assert.Equal(t, expected, actual)
   assert.NoError(t, err)
   assert.True(t, condition)
   ```

3. **Test File Naming**:
   - Test files: `*_test.go`
   - Test fixtures: `tests/fixtures/`
   - Test helpers: `tests/helpers/`

### Running Tests

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# Run specific test
go test ./internal/validation -run TestValidateConfig

# Run integration tests
go test ./tests/integration
```

### Test Coverage

- Maintain minimum 80% test coverage
- Include tests for error conditions
- Test edge cases and boundary conditions
- Use coverage reports to identify untested code

## Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring without feature changes
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **ci**: CI/CD changes
- **build**: Build system changes

### Examples

```bash
# Simple feature
feat: add Docker multi-stage build support

# Bug fix with scope
fix(validation): handle empty environment variables

# Breaking change
feat!: redesign configuration structure

BREAKING CHANGE: Configuration format has changed.
See migration guide for details.

# With issue reference
fix: resolve memory leak in config parser

Closes #123
```

### Guidelines

- Use present tense ("add feature" not "added feature")
- Keep first line under 50 characters
- Reference issues and PRs when relevant
- Include breaking change notes when applicable

## Pull Request Process

### Before Submitting

1. **Ensure Tests Pass**
   ```bash
   just ci
   ```

2. **Validate Configurations**
   ```bash
   just validate
   just validate-strict
   ```

3. **Check Documentation**
   - Update relevant documentation
   - Add examples if needed
   - Update CHANGELOG.md if applicable

### PR Description Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to break)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if applicable)
- [ ] Conventional commit messages used
```

### Review Process

1. **Automated Checks**: CI must pass
2. **Code Review**: At least one maintainer review required
3. **Testing**: Comprehensive test coverage
4. **Documentation**: Updated as needed

## Code Review Guidelines

### For Authors

- **Keep PRs Small**: Focus on single features or fixes
- **Provide Context**: Clear description and rationale
- **Respond Promptly**: Address reviewer feedback quickly
- **Test Thoroughly**: Ensure changes work as expected

### For Reviewers

- **Be Constructive**: Provide helpful, actionable feedback
- **Focus on Important Issues**: Don't nitpick minor style issues
- **Suggest Improvements**: Offer specific recommendations
- **Test Changes**: Verify functionality when possible

### Review Checklist

- [ ] Code follows Go conventions and project patterns
- [ ] Tests are comprehensive and pass
- [ ] Documentation is updated
- [ ] No security vulnerabilities introduced
- [ ] Performance impact considered
- [ ] Backward compatibility maintained (or breaking changes documented)

## Documentation

### Types of Documentation

1. **Code Documentation**
   - Go doc comments for exported functions
   - Inline comments for complex logic
   - README files for packages

2. **User Documentation**
   - README.md updates
   - Command documentation
   - Configuration examples

3. **Developer Documentation**
   - Architecture decisions
   - Setup instructions
   - Contributing guidelines

### Documentation Standards

- Use clear, concise language
- Include practical examples
- Keep documentation up-to-date with code changes
- Use consistent formatting and style

## Getting Help

### Resources

- **Documentation**: Check README.md and inline documentation
- **Issues**: Search existing GitHub issues
- **Discussions**: Use GitHub Discussions for questions
- **Commands**: Run `just help` for available commands

### Asking Questions

When asking for help:

1. **Search First**: Check existing documentation and issues
2. **Provide Context**: Include relevant system information
3. **Show What You Tried**: Include commands run and error messages
4. **Be Specific**: Clear, focused questions get better answers

### Community Guidelines

- Be patient and respectful
- Help others when you can
- Share knowledge and experiences
- Follow up on your questions with solutions

## Recognition

Contributors are recognized in:

- **README.md**: Contributors section
- **Release Notes**: Major contributions highlighted
- **GitHub**: Contributor graphs and statistics

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to the GoReleaser Template project!**

For additional help or questions, feel free to:
- Open an issue on GitHub
- Start a discussion in GitHub Discussions
- Review existing documentation and examples