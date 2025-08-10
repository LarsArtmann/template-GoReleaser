package validation

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ValidatorFunction represents a validation function
type ValidatorFunction func(value string) mo.Result[string]

// GetValidators returns a map of all available validators
func GetValidators() map[string]ValidatorFunction {
	return map[string]ValidatorFunction{
		"github_token":        ValidateGitHubToken,
		"docker_token":        ValidateDockerToken,
		"email":               ValidateEmail,
		"url":                 ValidateURL,
		"aws_bucket_name":     ValidateAWSBucketName,
		"gcs_bucket_name":     ValidateGCSBucketName,
		"azure_storage_name":  ValidateAzureStorageName,
		"hostname":            ValidateHostname,
		"file_path":           ValidateFilePath,
	}
}

// ValidateGitHubToken validates GitHub token format
func ValidateGitHubToken(token string) mo.Result[string] {
	if token == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty GitHub token",
			"GitHub token cannot be empty",
			nil,
		))
	}

	// GitHub Personal Access Token (classic): ghp_[36 chars]
	classicPattern := regexp.MustCompile(`^ghp_[A-Za-z0-9]{36}$`)
	
	// GitHub Fine-grained Personal Access Token: github_pat_[82 chars including underscores]
	fineGrainedPattern := regexp.MustCompile(`^github_pat_[A-Za-z0-9_]{82}$`)

	if classicPattern.MatchString(token) || fineGrainedPattern.MatchString(token) {
		return mo.Ok("Valid GitHub token format")
	}

	return mo.Err[string](NewUserFriendlyError(
		"invalid GitHub token format",
		"GitHub token must be a valid Personal Access Token (ghp_...) or Fine-grained token (github_pat_...)",
		nil,
	))
}

// ValidateDockerToken validates Docker Hub token format
func ValidateDockerToken(token string) mo.Result[string] {
	if token == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty Docker token",
			"Docker token cannot be empty",
			nil,
		))
	}

	// Docker Personal Access Token: dckr_pat_[at least 30 chars including underscores and hyphens]
	pattern := regexp.MustCompile(`^dckr_pat_[A-Za-z0-9_-]{30,}$`)

	if pattern.MatchString(token) {
		return mo.Ok("Valid Docker token format")
	}

	return mo.Err[string](NewUserFriendlyError(
		"invalid Docker token format",
		"Docker token must be a valid Personal Access Token (dckr_pat_...)",
		nil,
	))
}

// ValidateEmail validates email format
func ValidateEmail(email string) mo.Result[string] {
	if email == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty email",
			"Email address cannot be empty",
			nil,
		))
	}

	// Basic email validation regex
	pattern := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

	if pattern.MatchString(email) {
		return mo.Ok("Valid email format")
	}

	return mo.Err[string](NewUserFriendlyError(
		"invalid email format",
		"Email must be a valid email address (e.g., user@example.com)",
		nil,
	))
}

// ValidateURL validates URL format
func ValidateURL(urlStr string) mo.Result[string] {
	if urlStr == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty URL",
			"URL cannot be empty",
			nil,
		))
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return mo.Err[string](NewUserFriendlyError(
			"invalid URL format",
			"URL must be a valid HTTP or HTTPS URL",
			err,
		))
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return mo.Err[string](NewUserFriendlyError(
			"invalid URL scheme",
			"URL must start with http:// or https://",
			nil,
		))
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return mo.Err[string](NewUserFriendlyError(
			"invalid URL host",
			"URL must have a valid hostname",
			nil,
		))
	}

	return mo.Ok("Valid URL format")
}

// ValidateAWSBucketName validates AWS S3 bucket name format
func ValidateAWSBucketName(bucketName string) mo.Result[string] {
	if bucketName == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty bucket name",
			"S3 bucket name cannot be empty",
			nil,
		))
	}

	// AWS S3 bucket naming rules:
	// - 3-63 characters long
	// - Can contain lowercase letters, numbers, and hyphens
	// - Must start and end with a letter or number
	// - Cannot contain consecutive periods or hyphens
	// - Cannot be formatted as an IP address

	if len(bucketName) < 3 || len(bucketName) > 63 {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name length",
			"S3 bucket name must be 3-63 characters long",
			nil,
		))
	}

	// Check format
	pattern := regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	if !pattern.MatchString(bucketName) {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name format",
			"S3 bucket name can only contain lowercase letters, numbers, and hyphens, and must start/end with alphanumeric characters",
			nil,
		))
	}

	// Check for consecutive periods or hyphens
	if strings.Contains(bucketName, "..") || strings.Contains(bucketName, "--") {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name format",
			"S3 bucket name cannot contain consecutive periods or hyphens",
			nil,
		))
	}

	// Check if it looks like an IP address
	ipPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	if ipPattern.MatchString(bucketName) {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name format",
			"S3 bucket name cannot be formatted as an IP address",
			nil,
		))
	}

	return mo.Ok("Valid S3 bucket name")
}

// ValidateGCSBucketName validates Google Cloud Storage bucket name format
func ValidateGCSBucketName(bucketName string) mo.Result[string] {
	if bucketName == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty bucket name",
			"GCS bucket name cannot be empty",
			nil,
		))
	}

	// GCS bucket naming rules:
	// - 3-63 characters long
	// - Can contain lowercase letters, numbers, hyphens, and underscores
	// - Must start and end with a letter or number

	if len(bucketName) < 3 || len(bucketName) > 63 {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name length",
			"GCS bucket name must be 3-63 characters long",
			nil,
		))
	}

	// Check format
	pattern := regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*[a-z0-9]$`)
	if !pattern.MatchString(bucketName) {
		return mo.Err[string](NewUserFriendlyError(
			"invalid bucket name format",
			"GCS bucket name can only contain lowercase letters, numbers, dots, hyphens, and underscores, and must start/end with alphanumeric characters",
			nil,
		))
	}

	return mo.Ok("Valid GCS bucket name")
}

// ValidateAzureStorageName validates Azure storage account name format
func ValidateAzureStorageName(name string) mo.Result[string] {
	if name == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty storage account name",
			"Azure storage account name cannot be empty",
			nil,
		))
	}

	// Azure storage account naming rules:
	// - 3-24 characters long
	// - Can contain only lowercase letters and numbers
	// - Must be unique across Azure

	if len(name) < 3 || len(name) > 24 {
		return mo.Err[string](NewUserFriendlyError(
			"invalid storage account name length",
			"Azure storage account name must be 3-24 characters long",
			nil,
		))
	}

	// Check format - only lowercase alphanumeric
	pattern := regexp.MustCompile(`^[a-z0-9]{3,24}$`)
	if !pattern.MatchString(name) {
		return mo.Err[string](NewUserFriendlyError(
			"invalid storage account name format",
			"Azure storage account name can only contain lowercase letters and numbers",
			nil,
		))
	}

	return mo.Ok("Valid Azure storage account name")
}

// ValidateHostname validates hostname format
func ValidateHostname(hostname string) mo.Result[string] {
	if hostname == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty hostname",
			"Hostname cannot be empty",
			nil,
		))
	}

	// Basic hostname validation
	pattern := regexp.MustCompile(`^[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	if !pattern.MatchString(hostname) {
		return mo.Err[string](NewUserFriendlyError(
			"invalid hostname format",
			"Hostname must be a valid domain name (e.g., example.com)",
			nil,
		))
	}

	return mo.Ok("Valid hostname")
}

// ValidateFilePath validates that a file exists at the given path
func ValidateFilePath(path string) mo.Result[string] {
	if path == "" {
		return mo.Err[string](NewUserFriendlyError(
			"empty file path",
			"File path cannot be empty",
			nil,
		))
	}

	// Expand tilde if present
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return mo.Err[string](NewUserFriendlyError(
				"cannot expand home directory",
				"Unable to expand ~ in file path",
				err,
			))
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return mo.Err[string](NewUserFriendlyError(
			"file does not exist",
			fmt.Sprintf("File does not exist: %s", path),
			err,
		))
	} else if err != nil {
		return mo.Err[string](NewUserFriendlyError(
			"cannot access file",
			fmt.Sprintf("Cannot access file: %s", path),
			err,
		))
	}

	return mo.Ok("File exists and is accessible")
}

// ValidateEnvironmentVariable validates a single environment variable
func ValidateEnvironmentVariable(name, value string, envVar EnvironmentVariable) ValidationIssue {
	// Check if value is empty
	if value == "" {
		if envVar.Required {
			return ValidationIssue{
				Field:       name,
				Message:     fmt.Sprintf("%s is required but not set", name),
				UserMessage: fmt.Sprintf("%s is required. %s", name, envVar.Description),
				Severity:    SeverityCritical,
				Code:        "MISSING_REQUIRED",
			}
		} else {
			return ValidationIssue{
				Field:       name,
				Message:     fmt.Sprintf("%s is not set", name),
				UserMessage: fmt.Sprintf("%s is not configured. %s", name, envVar.Description),
				Severity:    SeverityWarning,
				Code:        "MISSING_OPTIONAL",
			}
		}
	}

	// Check for placeholder values
	placeholderPatterns := []string{"your-", "xxxx", "example", "changeme", "todo", "replace"}
	lowerValue := strings.ToLower(value)
	
	if lo.SomeBy(placeholderPatterns, func(pattern string) bool {
		return strings.HasPrefix(lowerValue, pattern)
	}) || len(value) < 3 {
		return ValidationIssue{
			Field:       name,
			Message:     fmt.Sprintf("%s appears to be a placeholder value", name),
			UserMessage: fmt.Sprintf("%s looks like a placeholder. Please set a real value.", name),
			Severity:    SeverityWarning,
			Code:        "PLACEHOLDER_VALUE",
			Value:       maskValue(value),
		}
	}

	// Run custom validator if available
	if validatorName, hasValidator := envVar.Validator.Get(); hasValidator {
		validators := GetValidators()
		if validator, exists := validators[validatorName]; exists {
			if result := validator(value); result.IsError() {
				err := result.Error()
				if userErr, ok := err.(*UserFriendlyError); ok {
					severity := SeverityError
					if envVar.Required {
						severity = SeverityCritical
					}
					
					return ValidationIssue{
						Field:       name,
						Message:     userErr.Message,
						UserMessage: userErr.UserMessage,
						Severity:    severity,
						Code:        "VALIDATION_FAILED",
						Value:       maskValue(value),
					}
				}
				
				return ValidationIssue{
					Field:       name,
					Message:     err.Error(),
					UserMessage: fmt.Sprintf("%s has an invalid value", name),
					Severity:    SeverityError,
					Code:        "VALIDATION_FAILED",
					Value:       maskValue(value),
				}
			}
		}
	}

	// If we get here, the variable is valid
	return ValidationIssue{} // Empty issue means no problem
}

// maskValue masks sensitive values for logging/display
func maskValue(value string) string {
	if len(value) <= 4 {
		return "***"
	}
	
	// Show first 2 and last 2 characters with stars in between
	return fmt.Sprintf("%s***%s", value[:2], value[len(value)-2:])
}

// IsValidationIssueEmpty checks if a validation issue is empty (no problem)
func IsValidationIssueEmpty(issue ValidationIssue) bool {
	return issue.Field == "" && issue.Message == ""
}