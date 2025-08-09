# Environment Variables Guide

This document provides comprehensive information about environment variables used in the GoReleaser template configurations.

## Overview

The GoReleaser template uses environment variables to:
- Keep sensitive information out of configuration files
- Allow customization without modifying YAML files
- Enable different configurations for different environments
- Support both GitHub Actions and local development

## Quick Setup

1. **Copy the example file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit with your values:**
   ```bash
   editor .env
   ```

3. **Source the environment:**
   ```bash
   source .env
   ```

4. **Validate your setup:**
   ```bash
   ./validate-env.sh
   ```

## Environment Variable Categories

### 🚨 Critical Variables (Required)

These variables are essential for basic GoReleaser functionality:

#### GITHUB_TOKEN
- **Description:** GitHub API access token for releases
- **Required by:** Both free and pro versions
- **Format:** `ghp_...` (Personal Access Token) or `github_pat_...` (Fine-grained token)
- **Permissions needed:** `contents:write`, `metadata:read`
- **Example:** `ghp_abcdefghijklmnopqrstuvwxyz123456789012`
- **How to get:** GitHub Settings → Developer settings → Personal access tokens

#### GITHUB_OWNER
- **Description:** GitHub repository owner/organization name
- **Required by:** Both free and pro versions
- **Format:** Username or organization name
- **Example:** `LarsArtmann`
- **Auto-set in Actions:** `${{ github.repository_owner }}`

#### GITHUB_REPO
- **Description:** GitHub repository name
- **Required by:** Both free and pro versions
- **Format:** Repository name without owner
- **Example:** `template-GoReleaser`
- **Auto-set in Actions:** `${{ github.event.repository.name }}`

### 📦 Docker Configuration (Optional)

Required only if Docker features are enabled:

#### DOCKER_USERNAME
- **Description:** Docker Hub username for container publishing
- **Required by:** Pro version with Docker enabled
- **Format:** Docker Hub username
- **Example:** `larsartmann`

#### DOCKER_TOKEN
- **Description:** Docker Hub access token
- **Required by:** Pro version with Docker enabled
- **Format:** `dckr_pat_...`
- **Example:** `dckr_pat_abcdefghijklmnopqrstuvwxyz12345678`
- **How to get:** Docker Hub → Account Settings → Security → New Access Token

### 📋 Project Information (Recommended)

Used for package metadata:

#### PROJECT_DESCRIPTION
- **Description:** Brief description of your project
- **Used in:** Package descriptions, release notes
- **Example:** `"A powerful CLI tool for managing infrastructure"`

#### VENDOR_NAME
- **Description:** Company or vendor name
- **Used in:** Package metadata, about information
- **Example:** `"Acme Corporation"`

#### MAINTAINER_NAME
- **Description:** Package maintainer's name
- **Used in:** Package metadata, contact information
- **Example:** `"Lars Artmann"`

#### MAINTAINER_EMAIL
- **Description:** Package maintainer's email address
- **Used in:** Package metadata, contact information
- **Format:** Valid email address
- **Example:** `"lars@example.com"`

### 🍺 Package Managers (Pro Features)

#### Homebrew
- **HOMEBREW_TAP_GITHUB_TOKEN:** GitHub token for Homebrew tap repository
- **Format:** Same as GITHUB_TOKEN
- **Permissions:** Access to your Homebrew tap repository

#### Scoop (Windows)
- **SCOOP_GITHUB_TOKEN:** GitHub token for Scoop bucket repository
- **Format:** Same as GITHUB_TOKEN
- **Permissions:** Access to your Scoop bucket repository

#### AUR (Arch Linux)
- **AUR_KEY:** Path to SSH private key for AUR access
- **Format:** File path (e.g., `~/.ssh/aur`)
- **Requirements:** SSH key registered with AUR account

### ☁️ Cloud Storage (Pro Features)

#### Amazon S3
- **S3_BUCKET:** S3 bucket name for release artifacts
- **AWS_ACCESS_KEY_ID:** AWS access key ID
- **AWS_SECRET_ACCESS_KEY:** AWS secret access key
- **Format:** Standard AWS credentials

#### Google Cloud Storage
- **GCS_BUCKET:** GCS bucket name for release artifacts
- **GOOGLE_APPLICATION_CREDENTIALS:** Path to service account JSON file

#### Azure Blob Storage
- **AZURE_STORAGE_ACCOUNT:** Azure storage account name
- **AZURE_STORAGE_CONTAINER:** Container name
- **AZURE_STORAGE_KEY:** Storage account access key

### 📦 Package Repositories (Pro Features)

#### Fury.io
- **FURY_TOKEN:** Fury.io API token
- **FURY_ACCOUNT:** Fury.io account name
- **Used for:** Publishing packages to Fury.io

#### Artifactory
- **ARTIFACTORY_HOST:** Artifactory server hostname
- **ARTIFACTORY_REPO:** Repository name
- **ARTIFACTORY_USERNAME:** Username for authentication
- **ARTIFACTORY_PASSWORD:** Password or API token

### 📢 Notifications (Pro Features)

#### Discord
- **DISCORD_WEBHOOK_ID:** Discord webhook ID
- **DISCORD_WEBHOOK_TOKEN:** Discord webhook token
- **Format:** From Discord webhook URL: `https://discord.com/api/webhooks/{ID}/{TOKEN}`

#### Slack
- **SLACK_WEBHOOK_URL:** Slack webhook URL
- **Format:** `https://hooks.slack.com/services/...`

#### Microsoft Teams
- **TEAMS_WEBHOOK_URL:** Teams webhook URL
- **Format:** `https://outlook.office.com/webhook/...`

#### Email (SMTP)
- **SMTP_FROM:** Sender email address
- **SMTP_TO:** Recipient email address
- **SMTP_USERNAME:** SMTP server username
- **SMTP_PASSWORD:** SMTP server password

#### Custom Webhook
- **WEBHOOK_URL:** Custom webhook endpoint
- **WEBHOOK_TOKEN:** Authentication token for webhook

### 🔐 Security (Pro Features)

#### Code Signing
- **COSIGN_PRIVATE_KEY:** Path to Cosign private key file
- **COSIGN_PASSWORD:** Password for Cosign private key
- **Used for:** Signing release artifacts with Cosign

### 🔄 Alternative Git Providers

#### GitLab
- **GITLAB_TOKEN:** GitLab API token
- **Format:** `glpat_...`
- **Use instead of:** GitHub tokens when using GitLab

#### Gitea
- **GITEA_TOKEN:** Gitea API token
- **Use instead of:** GitHub tokens when using Gitea

## Validation

### Using the Validation Script

Run the environment validation script to check your setup:

```bash
# Full validation
./validate-env.sh

# List all variables
./validate-env.sh --list

# Generate documentation
./validate-env.sh --docs

# Test GoReleaser loading
./validate-env.sh --test
```

### Common Issues and Solutions

#### Invalid Token Formats
```
✗ GITHUB_TOKEN: set but invalid - Invalid GitHub token format
```
**Solution:** Ensure your token starts with `ghp_` (classic) or `github_pat_` (fine-grained).

#### Placeholder Values
```
⚠ PROJECT_DESCRIPTION: appears to be a placeholder value
```
**Solution:** Replace placeholder values with actual project information.

#### Missing Critical Variables
```
✗ GITHUB_TOKEN: MISSING (CRITICAL)
```
**Solution:** Set all critical environment variables before running GoReleaser.

### GitHub Actions Setup

For GitHub Actions, set environment variables as repository secrets:

1. Go to your repository on GitHub
2. Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add each required variable

**Required secrets for GitHub Actions:**
- `GITHUB_TOKEN` (usually auto-provided)
- `DOCKER_USERNAME` (if using Docker)
- `DOCKER_TOKEN` (if using Docker)
- `PROJECT_DESCRIPTION`
- `VENDOR_NAME`
- `MAINTAINER_NAME`
- `MAINTAINER_EMAIL`

## Security Best Practices

### ✅ Do:
- Use `.env` files for local development
- Add `.env` to your `.gitignore`
- Use GitHub secrets for CI/CD
- Rotate tokens regularly
- Use fine-grained tokens when possible
- Validate tokens before use

### ❌ Don't:
- Commit tokens to version control
- Share tokens in plain text
- Use overly broad permissions
- Hardcode tokens in configurations
- Ignore validation warnings

## Environment File Template

Create a `.env` file with this template:

```bash
# Critical Variables
GITHUB_TOKEN=ghp_your_github_token_here
GITHUB_OWNER=your-github-username
GITHUB_REPO=your-repo-name

# Docker (if using containers)
DOCKER_USERNAME=your-docker-username
DOCKER_TOKEN=dckr_pat_your_docker_token

# Project Information
PROJECT_DESCRIPTION="Your amazing project description"
VENDOR_NAME="Your Company Name"
MAINTAINER_NAME="Your Name"
MAINTAINER_EMAIL="your.email@example.com"

# Pro Features (uncomment and configure as needed)
# HOMEBREW_TAP_GITHUB_TOKEN=ghp_your_homebrew_token
# SCOOP_GITHUB_TOKEN=ghp_your_scoop_token
# COSIGN_PRIVATE_KEY=/path/to/cosign.key
# COSIGN_PASSWORD=your-cosign-password
```

## Testing Your Configuration

### Manual Testing
```bash
# Source environment
source .env

# Test with GoReleaser
goreleaser check --config .goreleaser.yaml
goreleaser build --snapshot --single-target --clean
```

### Automated Testing
```bash
# Run all validation scripts
./validate-env.sh
./verify.sh
./validate-strict.sh
```

## Troubleshooting

### Common Error Messages

#### "template: :1:15: executing \"\" at <.Env.VAR_NAME>: map has no entry for key \"VAR_NAME\""
**Cause:** Environment variable is referenced in config but not set.
**Solution:** Set the missing variable or remove the reference.

#### "error: failed to create release: POST https://api.github.com/repos/.../releases: 401"
**Cause:** Invalid or insufficient GitHub token permissions.
**Solution:** Check token validity and permissions.

#### "error: docker login failed"
**Cause:** Invalid Docker credentials.
**Solution:** Verify Docker username and token.

### Getting Help

1. **Run validation:** `./validate-env.sh`
2. **Check logs:** Review GoReleaser output for specific errors
3. **Verify tokens:** Test tokens manually via API
4. **Check permissions:** Ensure tokens have required scopes

## Advanced Configuration

### Environment-Specific Configs

For different environments, you can use multiple env files:

```bash
# Development
cp .env.example .env.dev
# Configure for development

# Production  
cp .env.example .env.prod
# Configure for production

# Load specific environment
source .env.dev  # or .env.prod
```

### Dynamic Environment Variables

Some variables can be set dynamically in GitHub Actions:

```yaml
env:
  GITHUB_OWNER: ${{ github.repository_owner }}
  GITHUB_REPO: ${{ github.event.repository.name }}
  BUILD_DATE: ${{ github.event.head_commit.timestamp }}
```

### Conditional Features

Use environment variables to enable/disable features:

```yaml
# Only include Docker if credentials are available
dockers:
  {{ if .Env.DOCKER_USERNAME }}
  - image_templates:
    - "{{ .Env.DOCKER_USERNAME }}/{{ .ProjectName }}:latest"
  {{ end }}
```

---

For more information, see:
- [GoReleaser Environment Variables Documentation](https://goreleaser.com/environment/)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Docker Access Tokens](https://docs.docker.com/docker-hub/access-tokens/)