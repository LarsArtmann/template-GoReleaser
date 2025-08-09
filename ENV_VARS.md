# Environment Variables

This document lists all environment variables used by the GoReleaser configurations.

## Critical Environment Variables

These variables are required for basic functionality:

### GITHUB_TOKEN
GitHub API access token for releases

**Format:** GitHub Personal Access Token (ghp_...) or Fine-grained token (github_pat_...)

### GITHUB_OWNER
GitHub repository owner/organization

### GITHUB_REPO
GitHub repository name

## Optional Environment Variables

These variables enable additional features:

### WEBHOOK_TOKEN
Custom webhook authentication token

### DISCORD_WEBHOOK_TOKEN
Discord webhook token

### COSIGN_PRIVATE_KEY
Cosign private key for signing

### VENDOR_NAME
Vendor/company name for packages

### AZURE_STORAGE_CONTAINER
Azure storage container name

### DISCORD_WEBHOOK_ID
Discord webhook ID for notifications

### SLACK_WEBHOOK_URL
Slack webhook URL for notifications

### PROJECT_DESCRIPTION
Project description for packages

### TEAMS_WEBHOOK_URL
Microsoft Teams webhook URL

### GOOGLE_APPLICATION_CREDENTIALS
GCS service account credentials path

### AUR_KEY
AUR SSH private key path

### DOCKER_USERNAME
Docker Hub username for container publishing

### ARTIFACTORY_REPO
Artifactory repository name

### SCOOP_GITHUB_TOKEN
GitHub token for Scoop bucket

### ARTIFACTORY_PASSWORD
Artifactory password/token

### ARTIFACTORY_HOST
Artifactory server hostname

### AWS_SECRET_ACCESS_KEY
AWS secret access key

### WEBHOOK_URL
Custom webhook URL

### SMTP_USERNAME
SMTP server username

### GITEA_TOKEN
Gitea API token (alternative to GitHub)

### AZURE_STORAGE_KEY
Azure storage access key

### FURY_TOKEN
Fury.io API token

### HOMEBREW_TAP_GITHUB_TOKEN
GitHub token for Homebrew tap

### AWS_ACCESS_KEY_ID
AWS access key ID

### COSIGN_PASSWORD
Cosign private key password

### MAINTAINER_EMAIL
Package maintainer email

### AZURE_STORAGE_ACCOUNT
Azure storage account name

### ARTIFACTORY_USERNAME
Artifactory username

### GCS_BUCKET
Google Cloud Storage bucket

### SMTP_PASSWORD
SMTP server password

### SMTP_TO
SMTP recipient email address

### FURY_ACCOUNT
Fury.io account name

### DOCKER_TOKEN
Docker Hub access token

### SMTP_FROM
SMTP sender email address

### GITLAB_TOKEN
GitLab API token (alternative to GitHub)

### S3_BUCKET
AWS S3 bucket for releases

### MAINTAINER_NAME
Package maintainer name

