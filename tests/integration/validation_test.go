package integration

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/template-GoReleaser/tests/fixtures"
	"github.com/LarsArtmann/template-GoReleaser/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnvironmentVariableValidation tests the environment variable validation workflow
func (suite *IntegrationTestSuite) TestEnvironmentVariableValidation() {
	testCases := []struct {
		name            string
		envVars         map[string]string
		expectSuccess   bool
		expectedMissing []string
	}{
		{
			name:            "All Required Variables Present",
			envVars:         fixtures.TestEnvironmentVars["complete"],
			expectSuccess:   true,
			expectedMissing: []string{},
		},
		{
			name:            "Minimal Variables Present",
			envVars:         fixtures.TestEnvironmentVars["minimal"],
			expectSuccess:   true, // Should pass with warnings
			expectedMissing: []string{},
		},
		{
			name: "Missing Critical Variables",
			envVars: map[string]string{
				"PROJECT_NAME": "test-project",
			},
			expectSuccess:   false,
			expectedMissing: []string{"GITHUB_TOKEN"},
		},
		{
			name:            "No Variables Set",
			envVars:         map[string]string{},
			expectSuccess:   false,
			expectedMissing: []string{"GITHUB_TOKEN"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "env-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Clear existing environment variables that might interfere
			cleanup := suite.clearRelevantEnvVars()
			suite.RegisterCleanup(cleanup)

			// Set test environment variables
			envCleanup := helpers.SetEnvVars(suite.T(), tc.envVars)
			suite.RegisterCleanup(envCleanup)

			// Run validation script
			validateScript := filepath.Join(testDir, "verify.sh")
			require.True(suite.T(), helpers.FileExists(validateScript), "Validation script should exist")

			// Make script executable
			err := os.Chmod(validateScript, 0755)
			require.NoError(suite.T(), err)

			// Execute validation
			result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")

			if tc.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Validation output: %s", result.Stdout)
					suite.T().Logf("Validation errors: %s", result.Stderr)
				}
				// Allow warnings but require overall success
				assert.True(suite.T(), result.ExitCode == 0 || strings.Contains(result.Stdout, "with warnings"),
					"Validation should succeed or complete with warnings for %s", tc.name)
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "Validation should fail for %s", tc.name)

				// Check that expected missing variables are mentioned
				for _, missing := range tc.expectedMissing {
					assert.True(suite.T(),
						strings.Contains(result.Stdout, missing) || strings.Contains(result.Stderr, missing),
						"Validation output should mention missing variable %s", missing)
				}
			}
		})
	}
}

// TestStrictValidation tests the strict validation workflow
func (suite *IntegrationTestSuite) TestStrictValidation() {
	testCases := []struct {
		name          string
		envVars       map[string]string
		expectSuccess bool
		description   string
	}{
		{
			name:          "Complete Environment",
			envVars:       fixtures.TestEnvironmentVars["complete"],
			expectSuccess: true,
			description:   "Should pass strict validation with all variables",
		},
		{
			name:          "Minimal Environment",
			envVars:       fixtures.TestEnvironmentVars["minimal"],
			expectSuccess: false, // Strict mode should be more demanding
			description:   "Should fail strict validation with minimal variables",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "strict-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Clear existing environment variables
			cleanup := suite.clearRelevantEnvVars()
			suite.RegisterCleanup(cleanup)

			// Set test environment variables
			envCleanup := helpers.SetEnvVars(suite.T(), tc.envVars)
			suite.RegisterCleanup(envCleanup)

			// Run strict validation script
			strictScript := filepath.Join(testDir, "validate-strict.sh")
			require.True(suite.T(), helpers.FileExists(strictScript), "Strict validation script should exist")

			// Make script executable
			err := os.Chmod(strictScript, 0755)
			require.NoError(suite.T(), err)

			// Execute strict validation
			result := helpers.RunCommand(suite.T(), testDir, "./validate-strict.sh")

			if tc.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Strict validation output: %s", result.Stdout)
					suite.T().Logf("Strict validation errors: %s", result.Stderr)
				}
				helpers.AssertCommandSuccess(suite.T(), result, "Strict validation should succeed: %s", tc.description)
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "Strict validation should fail: %s", tc.description)
			}
		})
	}
}

// TestValidationWithMissingDependencies tests validation when tools are missing
func (suite *IntegrationTestSuite) TestValidationWithMissingDependencies() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "missing-deps-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Set up environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
	suite.RegisterCleanup(cleanup)

	// Run validation in an environment where some tools might be missing
	validateScript := filepath.Join(testDir, "verify.sh")
	require.True(suite.T(), helpers.FileExists(validateScript), "Validation script should exist")

	// Make script executable
	err := os.Chmod(validateScript, 0755)
	require.NoError(suite.T(), err)

	// Execute validation
	result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")

	// Validation should handle missing tools gracefully
	// Either succeed with warnings or fail with helpful messages
	if result.ExitCode != 0 {
		// Check that error messages are helpful
		output := result.Stdout + result.Stderr
		isHelpful := strings.Contains(output, "not installed") ||
			strings.Contains(output, "command not found") ||
			strings.Contains(output, "tool") ||
			strings.Contains(output, "dependency")

		assert.True(suite.T(), isHelpful,
			"Validation should provide helpful messages about missing dependencies. Output: %s", output)
	}

	// Check that validation provides clear summary
	output := result.Stdout + result.Stderr
	assert.True(suite.T(),
		strings.Contains(output, "Summary") || strings.Contains(output, "Errors") || strings.Contains(output, "Warnings"),
		"Validation should provide a clear summary")
}

// TestConfigurationFileValidation tests validation of GoReleaser configuration files
func (suite *IntegrationTestSuite) TestConfigurationFileValidation() {
	configFiles := []string{
		".goreleaser.yaml",
		".goreleaser.pro.yaml",
	}

	for _, configFile := range configFiles {
		suite.Run("Validate "+configFile, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "config-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			configPath := filepath.Join(testDir, configFile)
			if !helpers.FileExists(configPath) {
				suite.T().Skipf("Configuration file %s does not exist", configFile)
			}

			// Set up environment variables (use complete set for pro config)
			envVars := fixtures.TestEnvironmentVars["minimal"]
			if strings.Contains(configFile, "pro") {
				envVars = fixtures.TestEnvironmentVars["complete"]
			}
			cleanup := helpers.SetEnvVars(suite.T(), envVars)
			suite.RegisterCleanup(cleanup)

			// Run validation script
			validateScript := filepath.Join(testDir, "verify.sh")
			require.True(suite.T(), helpers.FileExists(validateScript), "Validation script should exist")

			// Make script executable
			err := os.Chmod(validateScript, 0755)
			require.NoError(suite.T(), err)

			// Execute validation
			result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")

			// Log output for debugging
			if result.ExitCode != 0 {
				suite.T().Logf("Config validation output for %s: %s", configFile, result.Stdout)
				suite.T().Logf("Config validation errors for %s: %s", configFile, result.Stderr)
			}

			// Validation should either succeed or provide clear errors about the config
			output := result.Stdout + result.Stderr
			if result.ExitCode != 0 {
				// Check that config-related issues are mentioned
				isConfigRelated := strings.Contains(output, configFile) ||
					strings.Contains(output, "yaml") ||
					strings.Contains(output, "goreleaser") ||
					strings.Contains(output, "config")

				assert.True(suite.T(), isConfigRelated,
					"If validation fails, it should mention config-related issues for %s. Output: %s",
					configFile, output)
			}
		})
	}
}

// TestProjectStructureValidation tests validation of project structure
func (suite *IntegrationTestSuite) TestProjectStructureValidation() {
	structureTests := []struct {
		name           string
		setupFunc      func(string)
		expectIssues   bool
		expectedIssues []string
	}{
		{
			name: "Complete Project Structure",
			setupFunc: func(testDir string) {
				// Already has complete structure from template
			},
			expectIssues:   false,
			expectedIssues: []string{},
		},
		{
			name: "Missing LICENSE File",
			setupFunc: func(testDir string) {
				licensePath := filepath.Join(testDir, "LICENSE")
				if helpers.FileExists(licensePath) {
					os.Remove(licensePath)
				}
			},
			expectIssues:   true,
			expectedIssues: []string{"LICENSE"},
		},
		{
			name: "Missing go.mod File",
			setupFunc: func(testDir string) {
				goModPath := filepath.Join(testDir, "go.mod")
				if helpers.FileExists(goModPath) {
					os.Remove(goModPath)
				}
			},
			expectIssues:   true,
			expectedIssues: []string{"go.mod"},
		},
		{
			name: "Missing cmd Directory",
			setupFunc: func(testDir string) {
				cmdPath := filepath.Join(testDir, "cmd")
				if helpers.FileExists(cmdPath) {
					os.RemoveAll(cmdPath)
				}
			},
			expectIssues:   true,
			expectedIssues: []string{"cmd", "main.go"},
		},
	}

	for _, tt := range structureTests {
		suite.Run(tt.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "structure-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Apply test-specific modifications
			tt.setupFunc(testDir)

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
			suite.RegisterCleanup(cleanup)

			// Run validation script
			validateScript := filepath.Join(testDir, "verify.sh")
			require.True(suite.T(), helpers.FileExists(validateScript), "Validation script should exist")

			// Make script executable
			err := os.Chmod(validateScript, 0755)
			require.NoError(suite.T(), err)

			// Execute validation
			result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")

			output := result.Stdout + result.Stderr

			if tt.expectIssues {
				// Should detect issues
				for _, expectedIssue := range tt.expectedIssues {
					assert.True(suite.T(),
						strings.Contains(output, expectedIssue),
						"Validation should detect missing %s. Output: %s", expectedIssue, output)
				}
			} else {
				// Should pass or have only minor warnings
				if result.ExitCode != 0 {
					assert.True(suite.T(),
						strings.Contains(output, "warning"),
						"If validation fails, it should only be due to warnings for complete structure. Output: %s", output)
				}
			}
		})
	}
}

// clearRelevantEnvVars clears environment variables that might interfere with testing
func (suite *IntegrationTestSuite) clearRelevantEnvVars() func() {
	// List of environment variables that might interfere with tests
	relevantVars := []string{
		"GITHUB_TOKEN", "DOCKER_USERNAME", "DOCKER_PASSWORD", "GORELEASER_KEY",
		"COSIGN_PRIVATE_KEY", "FURY_TOKEN", "CHOCOLATEY_API_KEY", "AUR_KEY",
		"HOMEBREW_TAP_GITHUB_TOKEN", "DISCORD_WEBHOOK_ID", "DISCORD_WEBHOOK_TOKEN",
		"SLACK_WEBHOOK", "TWITTER_CONSUMER_KEY", "TWITTER_CONSUMER_SECRET",
		"TWITTER_ACCESS_TOKEN", "TWITTER_ACCESS_TOKEN_SECRET", "MASTODON_CLIENT_ID",
		"MASTODON_CLIENT_SECRET", "MASTODON_ACCESS_TOKEN", "MASTODON_SERVER",
		"REDDIT_APPLICATION_ID", "REDDIT_USERNAME", "REDDIT_PASSWORD",
		"TELEGRAM_TOKEN", "TELEGRAM_CHAT_ID", "LINKEDIN_ACCESS_TOKEN",
		"PROJECT_NAME", "PROJECT_DESCRIPTION", "PROJECT_AUTHOR", "PROJECT_EMAIL",
		"PROJECT_URL", "LICENSE_TYPE",
	}

	originalVars := make(map[string]string)
	for _, varName := range relevantVars {
		if value, exists := os.LookupEnv(varName); exists {
			originalVars[varName] = value
		}
		os.Unsetenv(varName)
	}

	// Return cleanup function
	return func() {
		for _, varName := range relevantVars {
			if originalValue, existed := originalVars[varName]; existed {
				os.Setenv(varName, originalValue)
			} else {
				os.Unsetenv(varName)
			}
		}
	}
}
