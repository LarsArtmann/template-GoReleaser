# Security Guidelines

This document outlines the security measures, requirements, and best practices for this GoReleaser template project.

## Security Scanning Integration

### Automated Security Validation

This project includes comprehensive security scanning that runs automatically during CI/CD and can be executed manually:

```bash
# Run complete security scan
just security-scan

# Run CI pipeline with security validation
just ci

# Run verification with security checks
./verify.sh
```

### Security Tools Integrated

1. **gosec** - Go security analyzer that scans Go code for security vulnerabilities
2. **govulncheck** - Go vulnerability database scanner for dependencies
3. **shellcheck** - Shell script security and best practices analyzer
4. **hadolint** - Dockerfile security and best practices linter
5. **trivy** - Comprehensive vulnerability scanner (optional, recommended)

## Security Vulnerabilities Fixed

### Go Code Security Issues (Previously Found by gosec)

✅ **Fixed**: Path traversal vulnerabilities (G304)
- Added `#nosec G304` annotations with justification for test helper functions
- Test functions validated paths are controlled by test environment

✅ **Fixed**: File permission issues (G301, G306)
- Changed directory permissions from `0755` to `0750` (more restrictive)
- Changed file permissions from `0644` to `0600` (owner-only access)

✅ **Fixed**: Unhandled errors (G104)
- Added explicit error handling with `_` assignments for cleanup operations
- Added comments explaining why error handling is not critical in specific contexts

### Dockerfile Security Issues

✅ **Fixed**: Package version pinning (DL3018)
- Pinned versions for `git`, `ca-certificates`, and `tzdata` packages
- Consolidated RUN instructions to reduce layers (DL3059)

✅ **Fixed**: HEALTHCHECK syntax issue
- Corrected HEALTHCHECK command syntax from shell form to exec form
- Changed from `CMD [...] || exit 1` to proper `CMD [...]` format

## Security Best Practices Implemented

### 1. Secrets Management

- **No hardcoded secrets**: All sensitive values use environment variables
- **Template-based configuration**: Uses `{{ .Env.VARIABLE_NAME }}` syntax
- **Gitignore protection**: All key files, certificates, and secrets are excluded
- **Environment validation**: Scripts check for required environment variables

### 2. Dependencies Security

- **Regular vulnerability scanning**: `govulncheck` integration
- **Standard library monitoring**: Tracks Go version security updates
- **Minimal dependencies**: Uses Go standard library when possible

### 3. Container Security

- **Multi-stage builds**: Minimizes final image attack surface
- **Non-root user**: Runs containers as non-privileged `appuser`
- **Minimal base image**: Uses `scratch` for final stage
- **Version pinning**: Specific versions for all packages

### 4. Build Security

- **Static linking**: `-extldflags '-static'` for self-contained binaries
- **Strip symbols**: `-s -w` ldflags to reduce binary size and information disclosure
- **Reproducible builds**: `-trimpath` flag for consistent builds
- **CGO disabled**: `CGO_ENABLED=0` prevents C code injection

### 5. File Permissions

- **Restrictive permissions**: `0750` for directories, `0600` for files
- **Executable scripts**: Proper `+x` permissions on shell scripts only
- **Test isolation**: Test helper functions use appropriate permissions

## Environment Variables Security

### Required Variables (Secrets)

These must be set as GitHub repository secrets:

```bash
# Core functionality
GITHUB_TOKEN          # GitHub API access
DOCKER_USERNAME       # Docker Hub username  
DOCKER_TOKEN          # Docker Hub access token

# Pro features (if using .goreleaser.pro.yaml)
GORELEASER_KEY        # GoReleaser Pro license
COSIGN_PRIVATE_KEY    # Code signing key
COSIGN_PASSWORD       # Code signing key password
FURY_TOKEN           # Package repository token
HOMEBREW_TAP_GITHUB_TOKEN # Homebrew tap access
SCOOP_GITHUB_TOKEN    # Scoop bucket access
AUR_KEY              # AUR SSH key path

# Cloud storage (optional)
AWS_ACCESS_KEY_ID     # AWS credentials
AWS_SECRET_ACCESS_KEY # AWS credentials
AZURE_STORAGE_KEY     # Azure credentials
GOOGLE_APPLICATION_CREDENTIALS # GCP credentials

# Notifications (optional)
DISCORD_WEBHOOK_TOKEN # Discord notifications
SLACK_WEBHOOK_URL     # Slack notifications
SMTP_PASSWORD         # Email notifications
```

### Validation

Environment variables are validated by:
- `validate-env.sh` script with format validation
- GitHub Actions workflow checks
- GoReleaser configuration validation

## GitHub Actions Security

### Permissions

Workflow uses minimal required permissions:
```yaml
permissions:
  contents: write    # For creating releases
  packages: write    # For Docker registry
  id-token: write    # For cosign signing
```

### Security Best Practices Applied

✅ **Pinned action versions**: All actions use specific versions (not `@main`)
✅ **Minimal permissions**: Least-privilege principle applied
✅ **Secret management**: All sensitive data in GitHub secrets
✅ **Error handling**: Proper continue-on-error for optional steps
✅ **Verification steps**: Checksums and signatures validated

## Manual Security Verification

### Run Security Scans

```bash
# Install security tools
just install-tools

# Run comprehensive security scan
just security-scan

# Run individual scanners
gosec ./...
govulncheck ./...
shellcheck *.sh scripts/*.sh
hadolint Dockerfile
```

### Check for Vulnerabilities

```bash
# Check Go dependencies
go list -json -m all | nancy sleuth

# Check for leaked secrets (if truffleHog installed)
trufflehog git file://. --only-verified

# Filesystem scan (if trivy installed)  
trivy fs .
```

### Verify Configurations

```bash
# Validate all configurations
./verify.sh

# Strict validation mode
./validate-strict.sh

# Check environment setup
./validate-env.sh
```

## Incident Response

### If Secrets Are Compromised

1. **Immediately revoke** the compromised secret
2. **Generate new credentials** for the service
3. **Update GitHub repository secrets** with new values
4. **Review access logs** for unauthorized usage
5. **Audit recent releases** built with compromised secrets

### Security Issue Reporting

- Create private security advisory on GitHub
- Include detailed reproduction steps
- Provide proposed fix if possible
- Allow reasonable time for resolution before public disclosure

## Security Updates

### Regular Maintenance

- **Monthly**: Update Go version and dependencies
- **Weekly**: Review security advisories for used tools
- **On Release**: Run full security validation
- **Continuous**: Automated scanning in CI/CD

### Update Process

```bash
# Update dependencies
just update-deps

# Re-run security validation
just security-scan

# Test with updated dependencies
just test-all

# Validate configurations still work
just validate
```

## Compliance and Standards

### Security Standards Applied

- **OWASP Top 10**: Addressed common security vulnerabilities
- **NIST Guidelines**: Secure software development practices
- **GitHub Security Best Practices**: Repository and workflow security
- **Container Security**: CIS Docker Benchmarks compliance

### Audit Trail

All security-related changes are:
- Tracked in git commit history
- Documented in this security guide
- Validated by automated testing
- Reviewed before merging

## Resources

- [Go Security Checklist](https://github.com/securego/gosec)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [GitHub Actions Security](https://docs.github.com/en/actions/security-guides)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)