package validation

import "github.com/samber/mo"

// GetEnvironmentVariables returns all defined environment variables
func GetEnvironmentVariables() map[string]EnvironmentVariable {
	return map[string]EnvironmentVariable{
		// Critical variables (required for basic functionality)
		"GITHUB_TOKEN": {
			Name:        "GITHUB_TOKEN",
			Description: "GitHub API access token for releases",
			Required:    true,
			Category:    CategoryGitHub,
			Validator:   mo.Some("github_token"),
			Example:     mo.Some("ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
			Format:      mo.Some("GitHub Personal Access Token (ghp_...) or Fine-grained token (github_pat_...)"),
		},
		"GITHUB_OWNER": {
			Name:        "GITHUB_OWNER",
			Description: "GitHub repository owner/organization",
			Required:    true,
			Category:    CategoryGitHub,
			Example:     mo.Some("myorg"),
			Format:      mo.Some("GitHub username or organization name"),
		},
		"GITHUB_REPO": {
			Name:        "GITHUB_REPO",
			Description: "GitHub repository name",
			Required:    true,
			Category:    CategoryGitHub,
			Example:     mo.Some("myproject"),
			Format:      mo.Some("Repository name (without owner)"),
		},

		// Docker variables
		"DOCKER_USERNAME": {
			Name:        "DOCKER_USERNAME",
			Description: "Docker Hub username for container publishing",
			Required:    false,
			Category:    CategoryDocker,
			Example:     mo.Some("mydockerusername"),
			Format:      mo.Some("Docker Hub username"),
		},
		"DOCKER_TOKEN": {
			Name:        "DOCKER_TOKEN",
			Description: "Docker Hub access token",
			Required:    false,
			Category:    CategoryDocker,
			Validator:   mo.Some("docker_token"),
			Example:     mo.Some("dckr_pat_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
			Format:      mo.Some("Docker Personal Access Token (dckr_pat_...)"),
		},

		// Project metadata
		"PROJECT_DESCRIPTION": {
			Name:        "PROJECT_DESCRIPTION",
			Description: "Project description for packages",
			Required:    false,
			Category:    CategoryGeneral,
			Example:     mo.Some("A sample Go application"),
		},
		"VENDOR_NAME": {
			Name:        "VENDOR_NAME",
			Description: "Vendor/company name for packages",
			Required:    false,
			Category:    CategoryGeneral,
			Example:     mo.Some("My Company"),
		},
		"MAINTAINER_NAME": {
			Name:        "MAINTAINER_NAME",
			Description: "Package maintainer name",
			Required:    false,
			Category:    CategoryGeneral,
			Example:     mo.Some("John Doe"),
		},
		"MAINTAINER_EMAIL": {
			Name:        "MAINTAINER_EMAIL",
			Description: "Package maintainer email",
			Required:    false,
			Category:    CategoryGeneral,
			Validator:   mo.Some("email"),
			Example:     mo.Some("maintainer@example.com"),
			Format:      mo.Some("Valid email address"),
		},

		// Package managers
		"HOMEBREW_TAP_GITHUB_TOKEN": {
			Name:        "HOMEBREW_TAP_GITHUB_TOKEN",
			Description: "GitHub token for Homebrew tap",
			Required:    false,
			Category:    CategoryArtifacts,
			Validator:   mo.Some("github_token"),
			Format:      mo.Some("GitHub Personal Access Token for tap repository"),
		},
		"SCOOP_GITHUB_TOKEN": {
			Name:        "SCOOP_GITHUB_TOKEN",
			Description: "GitHub token for Scoop bucket",
			Required:    false,
			Category:    CategoryArtifacts,
			Validator:   mo.Some("github_token"),
			Format:      mo.Some("GitHub Personal Access Token for bucket repository"),
		},
		"AUR_KEY": {
			Name:        "AUR_KEY",
			Description: "AUR SSH private key path",
			Required:    false,
			Category:    CategoryArtifacts,
			Validator:   mo.Some("file_path"),
			Example:     mo.Some("/path/to/aur_rsa"),
			Format:      mo.Some("Path to SSH private key file"),
		},

		// Package repositories
		"FURY_TOKEN": {
			Name:        "FURY_TOKEN",
			Description: "Fury.io API token",
			Required:    false,
			Category:    CategoryArtifacts,
			Example:     mo.Some("your_fury_token"),
		},
		"FURY_ACCOUNT": {
			Name:        "FURY_ACCOUNT",
			Description: "Fury.io account name",
			Required:    false,
			Category:    CategoryArtifacts,
			Example:     mo.Some("myaccount"),
		},

		// Cloud storage - AWS S3
		"S3_BUCKET": {
			Name:        "S3_BUCKET",
			Description: "AWS S3 bucket for releases",
			Required:    false,
			Category:    CategoryCloud,
			Validator:   mo.Some("aws_bucket_name"),
			Example:     mo.Some("my-release-bucket"),
			Format:      mo.Some("Valid S3 bucket name (3-63 chars, lowercase alphanumeric and hyphens)"),
		},
		"AWS_ACCESS_KEY_ID": {
			Name:        "AWS_ACCESS_KEY_ID",
			Description: "AWS access key ID",
			Required:    false,
			Category:    CategoryCloud,
			Example:     mo.Some("AKIAIOSFODNN7EXAMPLE"),
		},
		"AWS_SECRET_ACCESS_KEY": {
			Name:        "AWS_SECRET_ACCESS_KEY",
			Description: "AWS secret access key",
			Required:    false,
			Category:    CategoryCloud,
			Example:     mo.Some("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
		},

		// Cloud storage - Azure
		"AZURE_STORAGE_ACCOUNT": {
			Name:        "AZURE_STORAGE_ACCOUNT",
			Description: "Azure storage account name",
			Required:    false,
			Category:    CategoryCloud,
			Validator:   mo.Some("azure_storage_name"),
			Example:     mo.Some("myazurestorage"),
			Format:      mo.Some("3-24 lowercase alphanumeric characters"),
		},
		"AZURE_STORAGE_CONTAINER": {
			Name:        "AZURE_STORAGE_CONTAINER",
			Description: "Azure storage container name",
			Required:    false,
			Category:    CategoryCloud,
			Example:     mo.Some("releases"),
		},
		"AZURE_STORAGE_KEY": {
			Name:        "AZURE_STORAGE_KEY",
			Description: "Azure storage access key",
			Required:    false,
			Category:    CategoryCloud,
			Example:     mo.Some("your_azure_storage_key"),
		},

		// Cloud storage - Google Cloud
		"GCS_BUCKET": {
			Name:        "GCS_BUCKET",
			Description: "Google Cloud Storage bucket",
			Required:    false,
			Category:    CategoryCloud,
			Validator:   mo.Some("gcs_bucket_name"),
			Example:     mo.Some("my-gcs-bucket"),
			Format:      mo.Some("Valid GCS bucket name (3-63 chars, lowercase alphanumeric, dots, hyphens)"),
		},
		"GOOGLE_APPLICATION_CREDENTIALS": {
			Name:        "GOOGLE_APPLICATION_CREDENTIALS",
			Description: "GCS service account credentials path",
			Required:    false,
			Category:    CategoryCloud,
			Validator:   mo.Some("file_path"),
			Example:     mo.Some("/path/to/service-account.json"),
			Format:      mo.Some("Path to Google service account JSON file"),
		},

		// Artifactory
		"ARTIFACTORY_HOST": {
			Name:        "ARTIFACTORY_HOST",
			Description: "Artifactory server hostname",
			Required:    false,
			Category:    CategoryArtifacts,
			Validator:   mo.Some("hostname"),
			Example:     mo.Some("mycompany.jfrog.io"),
			Format:      mo.Some("Valid hostname"),
		},
		"ARTIFACTORY_REPO": {
			Name:        "ARTIFACTORY_REPO",
			Description: "Artifactory repository name",
			Required:    false,
			Category:    CategoryArtifacts,
			Example:     mo.Some("my-generic-repo"),
		},
		"ARTIFACTORY_USERNAME": {
			Name:        "ARTIFACTORY_USERNAME",
			Description: "Artifactory username",
			Required:    false,
			Category:    CategoryArtifacts,
			Example:     mo.Some("myusername"),
		},
		"ARTIFACTORY_PASSWORD": {
			Name:        "ARTIFACTORY_PASSWORD",
			Description: "Artifactory password/token",
			Required:    false,
			Category:    CategoryArtifacts,
			Example:     mo.Some("mypassword_or_token"),
		},

		// Notifications - Discord
		"DISCORD_WEBHOOK_ID": {
			Name:        "DISCORD_WEBHOOK_ID",
			Description: "Discord webhook ID for notifications",
			Required:    false,
			Category:    CategoryNotification,
			Example:     mo.Some("1234567890123456789"),
		},
		"DISCORD_WEBHOOK_TOKEN": {
			Name:        "DISCORD_WEBHOOK_TOKEN",
			Description: "Discord webhook token",
			Required:    false,
			Category:    CategoryNotification,
			Example:     mo.Some("your_discord_webhook_token"),
		},

		// Notifications - Slack
		"SLACK_WEBHOOK_URL": {
			Name:        "SLACK_WEBHOOK_URL",
			Description: "Slack webhook URL for notifications",
			Required:    false,
			Category:    CategoryNotification,
			Validator:   mo.Some("url"),
			Example:     mo.Some("https://hooks.slack.com/services/..."),
			Format:      mo.Some("Valid HTTPS URL"),
		},

		// Notifications - Teams
		"TEAMS_WEBHOOK_URL": {
			Name:        "TEAMS_WEBHOOK_URL",
			Description: "Microsoft Teams webhook URL",
			Required:    false,
			Category:    CategoryNotification,
			Validator:   mo.Some("url"),
			Example:     mo.Some("https://company.webhook.office.com/..."),
			Format:      mo.Some("Valid HTTPS URL"),
		},

		// SMTP notifications
		"SMTP_FROM": {
			Name:        "SMTP_FROM",
			Description: "SMTP sender email address",
			Required:    false,
			Category:    CategoryNotification,
			Validator:   mo.Some("email"),
			Example:     mo.Some("noreply@example.com"),
			Format:      mo.Some("Valid email address"),
		},
		"SMTP_TO": {
			Name:        "SMTP_TO",
			Description: "SMTP recipient email address",
			Required:    false,
			Category:    CategoryNotification,
			Validator:   mo.Some("email"),
			Example:     mo.Some("releases@example.com"),
			Format:      mo.Some("Valid email address"),
		},
		"SMTP_USERNAME": {
			Name:        "SMTP_USERNAME",
			Description: "SMTP server username",
			Required:    false,
			Category:    CategoryNotification,
			Example:     mo.Some("smtp_user"),
		},
		"SMTP_PASSWORD": {
			Name:        "SMTP_PASSWORD",
			Description: "SMTP server password",
			Required:    false,
			Category:    CategoryNotification,
			Example:     mo.Some("smtp_password"),
		},

		// Custom webhooks
		"WEBHOOK_URL": {
			Name:        "WEBHOOK_URL",
			Description: "Custom webhook URL",
			Required:    false,
			Category:    CategoryNotification,
			Validator:   mo.Some("url"),
			Example:     mo.Some("https://api.example.com/webhook"),
			Format:      mo.Some("Valid HTTP/HTTPS URL"),
		},
		"WEBHOOK_TOKEN": {
			Name:        "WEBHOOK_TOKEN",
			Description: "Custom webhook authentication token",
			Required:    false,
			Category:    CategoryNotification,
			Example:     mo.Some("your_webhook_token"),
		},

		// Signing
		"COSIGN_PRIVATE_KEY": {
			Name:        "COSIGN_PRIVATE_KEY",
			Description: "Cosign private key for signing",
			Required:    false,
			Category:    CategorySigning,
			Validator:   mo.Some("file_path"),
			Example:     mo.Some("/path/to/cosign.key"),
			Format:      mo.Some("Path to Cosign private key file"),
		},
		"COSIGN_PASSWORD": {
			Name:        "COSIGN_PASSWORD",
			Description: "Cosign private key password",
			Required:    false,
			Category:    CategorySigning,
			Example:     mo.Some("your_cosign_password"),
		},

		// Alternative Git providers
		"GITLAB_TOKEN": {
			Name:        "GITLAB_TOKEN",
			Description: "GitLab API token (alternative to GitHub)",
			Required:    false,
			Category:    CategoryGitHub, // Using GitHub category for Git providers
			Example:     mo.Some("glpat-xxxxxxxxxxxxxxxxxxxx"),
		},
		"GITEA_TOKEN": {
			Name:        "GITEA_TOKEN",
			Description: "Gitea API token (alternative to GitHub)",
			Required:    false,
			Category:    CategoryGitHub, // Using GitHub category for Git providers
			Example:     mo.Some("your_gitea_token"),
		},
	}
}

// GetCriticalVariables returns only the critical environment variables
func GetCriticalVariables() map[string]EnvironmentVariable {
	allVars := GetEnvironmentVariables()
	criticalVars := make(map[string]EnvironmentVariable)
	
	for name, envVar := range allVars {
		if envVar.Required {
			criticalVars[name] = envVar
		}
	}
	
	return criticalVars
}

// GetOptionalVariables returns only the optional environment variables
func GetOptionalVariables() map[string]EnvironmentVariable {
	allVars := GetEnvironmentVariables()
	optionalVars := make(map[string]EnvironmentVariable)
	
	for name, envVar := range allVars {
		if !envVar.Required {
			optionalVars[name] = envVar
		}
	}
	
	return optionalVars
}

// GetVariablesByCategory returns variables filtered by category
func GetVariablesByCategory(category VariableCategory) map[string]EnvironmentVariable {
	allVars := GetEnvironmentVariables()
	categoryVars := make(map[string]EnvironmentVariable)
	
	for name, envVar := range allVars {
		if envVar.Category == category {
			categoryVars[name] = envVar
		}
	}
	
	return categoryVars
}