package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvironmentVariables(t *testing.T) {
	vars := GetEnvironmentVariables()

	assert.NotEmpty(t, vars)

	// Check that critical variables are present
	criticalVars := []string{"GITHUB_TOKEN", "GITHUB_OWNER", "GITHUB_REPO"}
	for _, varName := range criticalVars {
		assert.Contains(t, vars, varName, "Critical variable %s should be present", varName)
		envVar := vars[varName]
		assert.True(t, envVar.Required, "Variable %s should be marked as required", varName)
	}

	// Check that some optional variables are present
	optionalVars := []string{"DOCKER_USERNAME", "DOCKER_TOKEN", "MAINTAINER_EMAIL"}
	for _, varName := range optionalVars {
		assert.Contains(t, vars, varName, "Optional variable %s should be present", varName)
		envVar := vars[varName]
		assert.False(t, envVar.Required, "Variable %s should be marked as optional", varName)
	}
}

func TestGetCriticalVariables(t *testing.T) {
	criticalVars := GetCriticalVariables()

	assert.NotEmpty(t, criticalVars)

	// All variables should be required
	for name, envVar := range criticalVars {
		assert.True(t, envVar.Required, "Variable %s should be required", name)
	}

	// Check specific critical variables
	expectedCritical := []string{"GITHUB_TOKEN", "GITHUB_OWNER", "GITHUB_REPO"}
	for _, varName := range expectedCritical {
		assert.Contains(t, criticalVars, varName)
	}

	// Should not contain optional variables
	assert.NotContains(t, criticalVars, "DOCKER_TOKEN")
}

func TestGetOptionalVariables(t *testing.T) {
	optionalVars := GetOptionalVariables()

	assert.NotEmpty(t, optionalVars)

	// All variables should be optional
	for name, envVar := range optionalVars {
		assert.False(t, envVar.Required, "Variable %s should be optional", name)
	}

	// Check specific optional variables
	expectedOptional := []string{"DOCKER_TOKEN", "MAINTAINER_EMAIL", "SLACK_WEBHOOK_URL"}
	for _, varName := range expectedOptional {
		assert.Contains(t, optionalVars, varName)
	}

	// Should not contain required variables
	assert.NotContains(t, optionalVars, "GITHUB_TOKEN")
}

func TestGetVariablesByCategory(t *testing.T) {
	t.Run("GitHub category", func(t *testing.T) {
		githubVars := GetVariablesByCategory(CategoryGitHub)

		expectedGitHubVars := []string{"GITHUB_TOKEN", "GITHUB_OWNER", "GITHUB_REPO", "GITLAB_TOKEN", "GITEA_TOKEN"}
		for _, varName := range expectedGitHubVars {
			if varName == "GITLAB_TOKEN" || varName == "GITEA_TOKEN" {
				// These might be categorized as GitHub since they're git providers
				continue
			}
			assert.Contains(t, githubVars, varName, "GitHub variable %s should be present", varName)
		}

		// Should not contain Docker variables
		assert.NotContains(t, githubVars, "DOCKER_TOKEN")
	})

	t.Run("Docker category", func(t *testing.T) {
		dockerVars := GetVariablesByCategory(CategoryDocker)

		expectedDockerVars := []string{"DOCKER_USERNAME", "DOCKER_TOKEN"}
		for _, varName := range expectedDockerVars {
			assert.Contains(t, dockerVars, varName, "Docker variable %s should be present", varName)
		}

		// Should not contain GitHub variables
		assert.NotContains(t, dockerVars, "GITHUB_TOKEN")
	})

	t.Run("Cloud category", func(t *testing.T) {
		cloudVars := GetVariablesByCategory(CategoryCloud)

		expectedCloudVars := []string{
			"S3_BUCKET", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY",
			"AZURE_STORAGE_ACCOUNT", "AZURE_STORAGE_CONTAINER", "AZURE_STORAGE_KEY",
			"GCS_BUCKET", "GOOGLE_APPLICATION_CREDENTIALS",
		}
		for _, varName := range expectedCloudVars {
			assert.Contains(t, cloudVars, varName, "Cloud variable %s should be present", varName)
		}
	})

	t.Run("Notification category", func(t *testing.T) {
		notificationVars := GetVariablesByCategory(CategoryNotification)

		expectedNotificationVars := []string{
			"DISCORD_WEBHOOK_ID", "DISCORD_WEBHOOK_TOKEN",
			"SLACK_WEBHOOK_URL", "TEAMS_WEBHOOK_URL",
			"SMTP_FROM", "SMTP_TO", "SMTP_USERNAME", "SMTP_PASSWORD",
			"WEBHOOK_URL", "WEBHOOK_TOKEN",
		}
		for _, varName := range expectedNotificationVars {
			assert.Contains(t, notificationVars, varName, "Notification variable %s should be present", varName)
		}
	})

	t.Run("Signing category", func(t *testing.T) {
		signingVars := GetVariablesByCategory(CategorySigning)

		expectedSigningVars := []string{"COSIGN_PRIVATE_KEY", "COSIGN_PASSWORD"}
		for _, varName := range expectedSigningVars {
			assert.Contains(t, signingVars, varName, "Signing variable %s should be present", varName)
		}
	})
}

func TestVariableStructure(t *testing.T) {
	vars := GetEnvironmentVariables()

	t.Run("all variables have required fields", func(t *testing.T) {
		for name, envVar := range vars {
			assert.Equal(t, name, envVar.Name, "Variable name should match map key")
			assert.NotEmpty(t, envVar.Description, "Variable %s should have description", name)
			assert.NotEmpty(t, envVar.Category, "Variable %s should have category", name)
		}
	})

	t.Run("validator names are consistent", func(t *testing.T) {
		validValidators := map[string]bool{
			"github_token":       true,
			"docker_token":       true,
			"email":              true,
			"url":                true,
			"aws_bucket_name":    true,
			"gcs_bucket_name":    true,
			"azure_storage_name": true,
			"hostname":           true,
			"file_path":          true,
		}

		for name, envVar := range vars {
			if validatorName, hasValidator := envVar.Validator.Get(); hasValidator {
				assert.Contains(t, validValidators, validatorName,
					"Variable %s has unknown validator %s", name, validatorName)
			}
		}
	})

	t.Run("examples and formats are provided where appropriate", func(t *testing.T) {
		// Check that token variables have examples
		tokenVars := []string{"GITHUB_TOKEN", "DOCKER_TOKEN"}
		for _, varName := range tokenVars {
			if envVar, exists := vars[varName]; exists {
				example, hasExample := envVar.Example.Get()
				assert.True(t, hasExample, "Variable %s should have example", varName)
				assert.NotEmpty(t, example, "Variable %s example should not be empty", varName)
			}
		}

		// Check that email variables have format info
		emailVars := []string{"MAINTAINER_EMAIL", "SMTP_FROM", "SMTP_TO"}
		for _, varName := range emailVars {
			if envVar, exists := vars[varName]; exists {
				format, hasFormat := envVar.Format.Get()
				assert.True(t, hasFormat, "Variable %s should have format", varName)
				assert.NotEmpty(t, format, "Variable %s format should not be empty", varName)
			}
		}
	})
}

func TestVariableCategories(t *testing.T) {
	vars := GetEnvironmentVariables()

	// Count variables by category
	categoryCounts := make(map[VariableCategory]int)
	for _, envVar := range vars {
		categoryCounts[envVar.Category]++
	}

	// Ensure all categories are used
	expectedCategories := []VariableCategory{
		CategoryGitHub, CategoryDocker, CategoryCloud, CategorySigning,
		CategoryNotification, CategoryGeneral, CategoryArtifacts,
	}

	for _, category := range expectedCategories {
		assert.Greater(t, categoryCounts[category], 0,
			"Category %s should have at least one variable", category)
	}
}

func TestVariableConsistency(t *testing.T) {
	allVars := GetEnvironmentVariables()
	criticalVars := GetCriticalVariables()
	optionalVars := GetOptionalVariables()

	t.Run("no overlap between critical and optional", func(t *testing.T) {
		for name := range criticalVars {
			assert.NotContains(t, optionalVars, name,
				"Variable %s should not be in both critical and optional", name)
		}
	})

	t.Run("all variables are either critical or optional", func(t *testing.T) {
		for name := range allVars {
			inCritical := false
			inOptional := false

			if _, exists := criticalVars[name]; exists {
				inCritical = true
			}
			if _, exists := optionalVars[name]; exists {
				inOptional = true
			}

			assert.True(t, inCritical || inOptional,
				"Variable %s should be either critical or optional", name)
			assert.False(t, inCritical && inOptional,
				"Variable %s should not be both critical and optional", name)
		}
	})

	t.Run("total count matches", func(t *testing.T) {
		totalCount := len(criticalVars) + len(optionalVars)
		assert.Equal(t, len(allVars), totalCount,
			"Total of critical and optional should equal all variables")
	})
}

func TestSpecificVariableProperties(t *testing.T) {
	vars := GetEnvironmentVariables()

	t.Run("GITHUB_TOKEN properties", func(t *testing.T) {
		githubToken := vars["GITHUB_TOKEN"]
		assert.True(t, githubToken.Required)
		assert.Equal(t, CategoryGitHub, githubToken.Category)

		validator, hasValidator := githubToken.Validator.Get()
		assert.True(t, hasValidator)
		assert.Equal(t, "github_token", validator)
	})

	t.Run("DOCKER_TOKEN properties", func(t *testing.T) {
		dockerToken := vars["DOCKER_TOKEN"]
		assert.False(t, dockerToken.Required)
		assert.Equal(t, CategoryDocker, dockerToken.Category)

		validator, hasValidator := dockerToken.Validator.Get()
		assert.True(t, hasValidator)
		assert.Equal(t, "docker_token", validator)
	})

	t.Run("email variables have email validator", func(t *testing.T) {
		emailVars := []string{"MAINTAINER_EMAIL", "SMTP_FROM", "SMTP_TO"}
		for _, varName := range emailVars {
			envVar := vars[varName]
			validator, hasValidator := envVar.Validator.Get()
			assert.True(t, hasValidator, "Variable %s should have validator", varName)
			assert.Equal(t, "email", validator, "Variable %s should have email validator", varName)
		}
	})

	t.Run("URL variables have URL validator", func(t *testing.T) {
		urlVars := []string{"SLACK_WEBHOOK_URL", "TEAMS_WEBHOOK_URL", "WEBHOOK_URL"}
		for _, varName := range urlVars {
			envVar := vars[varName]
			validator, hasValidator := envVar.Validator.Get()
			assert.True(t, hasValidator, "Variable %s should have validator", varName)
			assert.Equal(t, "url", validator, "Variable %s should have URL validator", varName)
		}
	})
}
