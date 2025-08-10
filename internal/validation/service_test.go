package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	assert.NotEmpty(t, service.envVars)
}

func TestValidateEnvironment(t *testing.T) {
	// Save original environment
	originalGitHubToken := os.Getenv("GITHUB_TOKEN")
	originalDockerToken := os.Getenv("DOCKER_TOKEN")
	defer func() {
		if originalGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", originalGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
		if originalDockerToken != "" {
			os.Setenv("DOCKER_TOKEN", originalDockerToken)
		} else {
			os.Unsetenv("DOCKER_TOKEN")
		}
	}()

	service := NewService()

	t.Run("missing critical variables", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_OWNER")
		os.Unsetenv("GITHUB_REPO")
		
		result := service.ValidateEnvironment()
		
		assert.False(t, result.Valid)
		assert.Equal(t, StatusNeedsSetup, result.ValidationStatus)
		assert.Contains(t, result.CriticalMissing, "GITHUB_TOKEN")
		assert.Contains(t, result.CriticalMissing, "GITHUB_OWNER")
		assert.Contains(t, result.CriticalMissing, "GITHUB_REPO")
	})

	t.Run("valid critical variables", func(t *testing.T) {
		// Set valid critical variables
		os.Setenv("GITHUB_TOKEN", "ghp_validtokenformathere1234567890123456")
		os.Setenv("GITHUB_OWNER", "testowner")
		os.Setenv("GITHUB_REPO", "testrepo")
		
		result := service.ValidateEnvironment()
		
		assert.True(t, result.Valid)
		assert.Empty(t, result.CriticalMissing)
		assert.Contains(t, result.ValidatedVariables, "GITHUB_TOKEN")
		assert.True(t, result.ValidatedVariables["GITHUB_TOKEN"].Valid)
	})

	t.Run("invalid token format", func(t *testing.T) {
		// Set invalid token
		os.Setenv("GITHUB_TOKEN", "invalid_token")
		os.Setenv("GITHUB_OWNER", "testowner")
		os.Setenv("GITHUB_REPO", "testrepo")
		
		result := service.ValidateEnvironment()
		
		assert.False(t, result.Valid)
		assert.Len(t, result.Issues, 1)
		assert.Equal(t, "VALIDATION_FAILED", result.Issues[0].Code)
	})

	t.Run("placeholder values", func(t *testing.T) {
		// Set placeholder value for a variable without a validator (GITHUB_OWNER)
		os.Setenv("GITHUB_TOKEN", "ghp_validtokenformathere1234567890123456")
		os.Setenv("GITHUB_OWNER", "your-owner-here")
		os.Setenv("GITHUB_REPO", "testrepo")
		
		result := service.ValidateEnvironment()
		
		assert.True(t, result.Valid) // Placeholder warnings don't fail validation
		
		// Find the placeholder warning
		found := false
		for _, warning := range result.Warnings {
			if warning.Code == "PLACEHOLDER_VALUE" && warning.Field == "GITHUB_OWNER" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find placeholder warning for GITHUB_OWNER")
		// There will be many warnings for missing optional variables
		assert.Greater(t, len(result.Warnings), 0)
	})
}

func TestExtractEnvironmentVariablesFromConfigs(t *testing.T) {
	tmpDir := t.TempDir()
	service := NewService()

	// Create test GoReleaser config files
	goreleaserContent := `
project_name: test
env:
  - GO111MODULE=on
builds:
  - binary: test
    goos:
      - linux
    env:
      - CGO_ENABLED=0
archives:
  - name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
release:
  github:
    owner: "{{ .Env.GITHUB_OWNER }}"
    name: "{{ .Env.GITHUB_REPO }}"
  name_template: "Release {{ .Version }}"
docker:
  - image: "{{ .Env.DOCKER_USERNAME }}/test"
    build_flag_templates:
      - "--build-arg=TOKEN={{ .Env.DOCKER_TOKEN }}"
`

	configFile := filepath.Join(tmpDir, ".goreleaser.yaml")
	err := os.WriteFile(configFile, []byte(goreleaserContent), 0644)
	require.NoError(t, err)

	t.Run("extract variables from config", func(t *testing.T) {
		result := service.ExtractEnvironmentVariablesFromConfigs(configFile)
		
		assert.Contains(t, result.ConfigFiles, configFile)
		assert.Contains(t, result.ExtractedVariables, "GITHUB_OWNER")
		assert.Contains(t, result.ExtractedVariables, "GITHUB_REPO")
		assert.Contains(t, result.ExtractedVariables, "DOCKER_USERNAME")
		assert.Contains(t, result.ExtractedVariables, "DOCKER_TOKEN")
		assert.Empty(t, result.Issues)
	})

	t.Run("non-existent config file", func(t *testing.T) {
		result := service.ExtractEnvironmentVariablesFromConfigs("/nonexistent/file.yaml")
		
		assert.Empty(t, result.ConfigFiles)
		assert.Empty(t, result.ExtractedVariables)
		assert.Empty(t, result.Issues) // Non-existent files are just skipped
	})
}

func TestValidateEnvExampleSync(t *testing.T) {
	tmpDir := t.TempDir()
	service := NewService()

	// Create test .env.example file
	envExampleContent := `# GitHub configuration
GITHUB_TOKEN=your_github_token
GITHUB_OWNER=your_username
GITHUB_REPO=your_repository

# Docker configuration
DOCKER_USERNAME=your_docker_username
DOCKER_TOKEN=your_docker_token

# Unused variable
UNUSED_VAR=unused_value
`

	envExampleFile := filepath.Join(tmpDir, ".env.example")
	err := os.WriteFile(envExampleFile, []byte(envExampleContent), 0644)
	require.NoError(t, err)

	// Create test GoReleaser config
	goreleaserContent := `
release:
  github:
    owner: "{{ .Env.GITHUB_OWNER }}"
    name: "{{ .Env.GITHUB_REPO }}"
docker:
  - image: "{{ .Env.DOCKER_USERNAME }}/test"
    build_flag_templates:
      - "--build-arg=TOKEN={{ .Env.DOCKER_TOKEN }}"
      - "--build-arg=MISSING={{ .Env.MISSING_VAR }}"
`

	configFile := filepath.Join(tmpDir, ".goreleaser.yaml")
	err = os.WriteFile(configFile, []byte(goreleaserContent), 0644)
	require.NoError(t, err)

	t.Run("sync validation", func(t *testing.T) {
		result := service.ValidateEnvExampleSync(envExampleFile, configFile)
		
		// Should find MISSING_VAR in config but not in .env.example
		assert.Contains(t, result.MissingInExample, "MISSING_VAR")
		
		// Should find UNUSED_VAR in .env.example but not used in config
		assert.Contains(t, result.UnusedInExample, "UNUSED_VAR")
		
		// Should have issues for missing variables
		assert.Len(t, result.Issues, 1)
		assert.Equal(t, "MISSING_IN_ENV_EXAMPLE", result.Issues[0].Code)
		
		// Should have warnings for unused variables
		assert.Len(t, result.Warnings, 1)
		assert.Equal(t, "UNUSED_IN_ENV_EXAMPLE", result.Warnings[0].Code)
	})
}

func TestGenerateValidationReport(t *testing.T) {
	// Save original environment
	originalGitHubToken := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", originalGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	service := NewService()

	t.Run("complete validation report", func(t *testing.T) {
		// Set up environment for testing
		os.Setenv("GITHUB_TOKEN", "ghp_validtokenformathere1234567890123456")
		os.Setenv("GITHUB_OWNER", "testowner")
		os.Setenv("GITHUB_REPO", "testrepo")
		
		report := service.GenerateValidationReport()
		
		assert.NotNil(t, report.Environment)
		assert.NotNil(t, report.ConfigAnalysis)
		assert.NotEmpty(t, report.Summary)
		assert.NotEmpty(t, report.RecommendedActions)
		
		// Environment should be valid since we have all critical variables
		assert.True(t, report.Environment.Valid)
		
		// Environment status should be Ready or HasIssues since critical vars are set
		assert.Contains(t, []EnvironmentValidationStatus{StatusReady, StatusHasIssues}, report.Environment.ValidationStatus)
	})

	t.Run("report with missing critical variables", func(t *testing.T) {
		// Clear critical variables
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_OWNER")
		os.Unsetenv("GITHUB_REPO")
		
		report := service.GenerateValidationReport()
		
		assert.Equal(t, StatusCriticalErrors, report.OverallStatus)
		assert.False(t, report.Environment.Valid)
		assert.Greater(t, report.Summary.CriticalIssues, 0)
		assert.Contains(t, report.RecommendedActions[0], "Set required environment variables")
	})
}

func TestGetVariableDocumentation(t *testing.T) {
	service := NewService()
	
	t.Run("existing variable", func(t *testing.T) {
		doc := service.GetVariableDocumentation("GITHUB_TOKEN")
		assert.True(t, doc.IsPresent())
		
		envVar, _ := doc.Get()
		assert.Equal(t, "GITHUB_TOKEN", envVar.Name)
		assert.True(t, envVar.Required)
	})
	
	t.Run("non-existent variable", func(t *testing.T) {
		doc := service.GetVariableDocumentation("NONEXISTENT_VAR")
		assert.False(t, doc.IsPresent())
	})
}

func TestServiceGetVariablesByCategory(t *testing.T) {
	service := NewService()
	
	githubVars := service.GetVariablesByCategory(CategoryGitHub)
	assert.Contains(t, githubVars, "GITHUB_TOKEN")
	assert.Contains(t, githubVars, "GITHUB_OWNER")
	assert.Contains(t, githubVars, "GITHUB_REPO")
	
	dockerVars := service.GetVariablesByCategory(CategoryDocker)
	assert.Contains(t, dockerVars, "DOCKER_USERNAME")
	assert.Contains(t, dockerVars, "DOCKER_TOKEN")
}

func TestListAllVariables(t *testing.T) {
	service := NewService()
	
	allVars := service.ListAllVariables()
	assert.NotEmpty(t, allVars)
	assert.Contains(t, allVars, "GITHUB_TOKEN")
	assert.Contains(t, allVars, "DOCKER_TOKEN")
	assert.Contains(t, allVars, "MAINTAINER_EMAIL")
}

func TestDetermineOverallStatus(t *testing.T) {
	service := NewService()
	
	tests := []struct {
		name           string
		envResult      *EnvironmentValidationResult
		configResult   *ConfigAnalysisResult
		expectedStatus ValidationStatus
	}{
		{
			name: "critical errors",
			envResult: &EnvironmentValidationResult{
				CriticalMissing: []string{"GITHUB_TOKEN"},
			},
			configResult:   &ConfigAnalysisResult{},
			expectedStatus: StatusCriticalErrors,
		},
		{
			name: "errors",
			envResult: &EnvironmentValidationResult{
				Issues: []ValidationIssue{{Severity: SeverityError}},
			},
			configResult:   &ConfigAnalysisResult{},
			expectedStatus: StatusErrors,
		},
		{
			name: "warnings",
			envResult: &EnvironmentValidationResult{
				Warnings: []ValidationIssue{{Severity: SeverityWarning}},
			},
			configResult:   &ConfigAnalysisResult{},
			expectedStatus: StatusWarnings,
		},
		{
			name: "ok",
			envResult: &EnvironmentValidationResult{},
			configResult: &ConfigAnalysisResult{},
			expectedStatus: StatusOK,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := service.determineOverallStatus(tt.envResult, tt.configResult)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestGenerateRecommendedActions(t *testing.T) {
	service := NewService()
	
	t.Run("missing critical variables", func(t *testing.T) {
		envResult := &EnvironmentValidationResult{
			CriticalMissing: []string{"GITHUB_TOKEN", "GITHUB_OWNER"},
		}
		configResult := &ConfigAnalysisResult{}
		
		actions := service.generateRecommendedActions(envResult, configResult)
		assert.Len(t, actions, 1)
		assert.Contains(t, actions[0], "Set required environment variables")
		assert.Contains(t, actions[0], "GITHUB_TOKEN")
		assert.Contains(t, actions[0], "GITHUB_OWNER")
	})
	
	t.Run("no issues", func(t *testing.T) {
		envResult := &EnvironmentValidationResult{}
		configResult := &ConfigAnalysisResult{}
		
		actions := service.generateRecommendedActions(envResult, configResult)
		assert.Len(t, actions, 1)
		assert.Contains(t, actions[0], "ready for GoReleaser execution")
	})
}