# Integration Testing Framework

This document describes the comprehensive integration testing framework implemented for the GoReleaser template project.

## Overview

The integration testing framework provides end-to-end validation of all critical workflows in the GoReleaser template, ensuring that:

- License generation workflow functions correctly
- GoReleaser configurations work with actual projects
- Environment variable validation catches configuration issues
- All components integrate properly
- Just commands work as expected
- Bash validation scripts provide accurate feedback

## Framework Architecture

### Testing Framework: Testify

The integration tests use **Testify** as the testing framework, chosen for:

- **Perfect fit for this project**: Straightforward integration testing needs
- **Lower complexity**: Right balance of power and simplicity
- **Excellent integration**: Works seamlessly with Go's standard testing package
- **Team-friendly**: Easier to adopt and maintain for template users
- **Comprehensive features**: Provides assertions, mocking, and suite-based testing

### Project Structure

```
tests/
├── integration/           # Integration test suites
│   ├── main_test.go      # Main test suite setup
│   ├── license_test.go   # License generation tests
│   ├── goreleaser_test.go # GoReleaser workflow tests
│   ├── validation_test.go # Environment validation tests
│   ├── components_test.go # Component integration tests
│   ├── justfile_test.go  # Just command tests
│   └── bash_validation_test.go # Bash script tests
├── helpers/              # Test helper utilities
│   └── test_helpers.go   # Common test functions
├── fixtures/             # Test data and configurations
│   └── test_configs.go   # Test fixtures and data
└── integration_test.go   # Package-level test entry point
```

## Test Categories

### 1. License Generation Tests (`license_test.go`)

Tests the complete license generation workflow:

- **MIT, Apache-2.0, BSD-3-Clause** license generation
- Integration with readme configuration
- License script help and list functionality
- Template validation and content verification
- Backup and restore functionality
- Error handling for invalid license types

### 2. GoReleaser Workflow Tests (`goreleaser_test.go`)

Validates GoReleaser configurations and build processes:

- **Configuration validation** for free and pro versions
- **Snapshot builds** with timeout handling
- **Dry-run testing** with proper Git setup
- **Single-target builds** for performance
- **Docker support validation** when configured
- **Build artifact verification**

### 3. Environment Variable Validation Tests (`validation_test.go`)

Ensures environment variable validation works correctly:

- **Complete vs minimal** environment setups
- **Missing critical variables** detection
- **Strict validation mode** testing
- **Configuration file validation**
- **Project structure validation**
- **Tool dependency checking**

### 4. Component Integration Tests (`components_test.go`)

Tests integration between different system components:

- **License system integration** with readme configs
- **GoReleaser config integration** with project structure
- **Validation script integration** with actual projects
- **Complete workflow integration** end-to-end
- **Cross-component compatibility**

### 5. Just Command Tests (`justfile_test.go`)

Validates that Just commands work with actual projects:

- **Basic commands**: help, list, init, format, clean
- **Build commands**: build, test, test-coverage
- **Validation commands**: validate, validate-strict, check
- **GoReleaser commands**: snapshot, dry-run, pro versions
- **CI workflow**: complete pipeline testing

### 6. Bash Validation Script Tests (`bash_validation_test.go`)

Converts existing bash validation to proper test suite:

- **Script functionality testing**
- **Edge case and error handling**
- **Tool detection validation**
- **Environment variable detection**
- **Output format verification**
- **Structured validation feedback**

## Running Tests

### Quick Start

```bash
# Run all integration tests
just integration-test

# Run integration tests with coverage
just integration-test-coverage

# Run all tests (unit + integration)
just test-all

# Run all tests with coverage reports
just test-all-coverage
```

### Individual Test Suites

```bash
# Run specific test suite
go test -v ./tests/integration/ -run TestLicenseGeneration

# Run with coverage
go test -v -race -coverprofile=coverage.out ./tests/integration/

# Run with timeout for long-running tests
go test -v -race -timeout=20m ./tests/integration/
```

### Environment Variables for Testing

The tests use predefined test environments:

```bash
# Minimal environment (basic functionality)
export GITHUB_TOKEN="test-token"
export DOCKER_USERNAME="testuser" 
export DOCKER_PASSWORD="testpass"

# Complete environment (all features)
export GITHUB_TOKEN="test-token"
export GORELEASER_KEY="test-pro-key"
export COSIGN_PRIVATE_KEY="test-cosign-key"
# ... (see fixtures/test_configs.go for complete list)
```

## CI/CD Integration

### GitHub Actions Workflow

The integration tests run automatically in GitHub Actions:

- **On push** to main/develop branches
- **On pull requests** to main
- **Daily schedule** at 2 AM UTC
- **Matrix testing** across OS and Go versions

### Coverage Reporting

- **Integration test coverage** reports generated
- **Coverage artifacts** uploaded for analysis
- **PR comments** with coverage summaries
- **Coverage thresholds** monitored

### Test Results

- **Structured output** with clear pass/fail indicators
- **Detailed logs** for debugging failures
- **Artifact uploads** for build outputs
- **Timeout protection** for long-running tests

## Test Development Guidelines

### Writing New Tests

1. **Use the TestSuite pattern** for related tests
2. **Create temporary directories** for each test
3. **Clean up resources** using RegisterCleanup
4. **Use test fixtures** for consistent data
5. **Validate both success and failure cases**

### Test Helper Functions

Available helper functions in `tests/helpers/test_helpers.go`:

```go
// Command execution
RunCommand(t, dir, command, args...)
RunCommandWithTimeout(t, timeout, dir, command, args...)

// File operations
CopyDir(t, src, dst)
CreateTestProject(t, templateDir, projectName)
FileExists(path)
FileContains(t, path, content)
WriteFile(t, path, content)

// Environment management
SetEnvVars(t, vars map[string]string) func()

// Assertions
AssertCommandSuccess(t, result, msgAndArgs...)
```

### Test Fixtures

Test fixtures are defined in `tests/fixtures/test_configs.go`:

```go
// Environment variable sets
TestEnvironmentVars["minimal"]
TestEnvironmentVars["complete"]

// Configuration templates
ReadmeConfigs["minimal"]
ReadmeConfigs["complete"]

// Expected license content
ExpectedLicenseContent["MIT"]
ExpectedLicenseContent["Apache-2.0"]
```

## Performance Considerations

### Test Optimization

- **Parallel execution** where safe
- **Timeout controls** for long-running operations
- **Minimal test environments** to reduce overhead
- **Cached dependencies** in CI/CD
- **Focused test runs** for development

### Resource Management

- **Temporary directories** automatically cleaned up
- **Environment isolation** between tests
- **Process cleanup** for spawned commands
- **Memory management** for large outputs

## Troubleshooting

### Common Issues

1. **Timeout errors**: Increase timeout for slow operations
2. **Permission errors**: Ensure scripts are executable
3. **Missing tools**: Install required dependencies
4. **Environment conflicts**: Clear environment variables
5. **File conflicts**: Use unique test directories

### Debug Mode

```bash
# Run with verbose output
go test -v ./tests/integration/

# Run specific failing test
go test -v ./tests/integration/ -run TestSpecificFunction

# Enable race detection
go test -v -race ./tests/integration/
```

### Log Analysis

- **Structured output** in test logs
- **Command outputs** captured in failures
- **Environment state** logged for debugging
- **Artifact inspection** in CI/CD

## Best Practices

### Test Design

1. **Test behavior, not implementation**
2. **Use descriptive test names**
3. **Test both happy and error paths**
4. **Isolate test environments**
5. **Make tests deterministic**

### Maintenance

1. **Keep tests up-to-date** with template changes
2. **Review test coverage** regularly
3. **Update fixtures** when features change
4. **Monitor test performance**
5. **Document test intentions**

### Integration

1. **Run tests before releases**
2. **Include in CI/CD pipelines**
3. **Monitor test results**
4. **Fix flaky tests immediately**
5. **Use tests as documentation**

## Future Enhancements

### Planned Improvements

- [ ] **Performance benchmarking** tests
- [ ] **Security validation** integration
- [ ] **Multi-platform testing** expansion
- [ ] **Docker integration** testing
- [ ] **Release workflow** validation

### Metrics and Monitoring

- [ ] **Test execution time** tracking
- [ ] **Coverage trend** analysis
- [ ] **Flaky test detection**
- [ ] **Success rate monitoring**
- [ ] **Resource usage** optimization

This integration testing framework ensures comprehensive validation of all critical workflows, providing confidence that the GoReleaser template works correctly for users across different environments and use cases.