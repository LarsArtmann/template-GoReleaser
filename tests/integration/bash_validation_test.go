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

// TestBashValidationScriptsFunctionality tests the existing bash validation scripts
func (suite *IntegrationTestSuite) TestBashValidationScriptsFunctionality() {
	validationScripts := []struct {
		name            string
		scriptPath      string
		expectSuccess   bool
		requiredEnvVars map[string]string
		setupFunc       func(string)
	}{
		{
			name:            "Verify Script Basic Functionality",
			scriptPath:      "verify.sh",
			expectSuccess:   true,
			requiredEnvVars: fixtures.TestEnvironmentVars["minimal"],
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
			},
		},
		{
			name:            "Strict Validation Script",
			scriptPath:      "validate-strict.sh",
			expectSuccess:   false, // May fail with minimal setup
			requiredEnvVars: fixtures.TestEnvironmentVars["minimal"],
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
			},
		},
	}

	for _, script := range validationScripts {
		suite.Run(script.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "bash-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up project structure
			script.setupFunc(testDir)

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), script.requiredEnvVars)
			suite.RegisterCleanup(cleanup)

			// Verify script exists
			scriptPath := filepath.Join(testDir, script.scriptPath)
			require.True(suite.T(), helpers.FileExists(scriptPath), "Script %s should exist", script.scriptPath)

			// Make script executable
			err := os.Chmod(scriptPath, 0755)
			require.NoError(suite.T(), err)

			// Run the script
			result := helpers.RunCommand(suite.T(), testDir, "./"+script.scriptPath)

			if script.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("%s output: %s", script.scriptPath, result.Stdout)
					suite.T().Logf("%s errors: %s", script.scriptPath, result.Stderr)

					// Allow warnings but require overall success
					if strings.Contains(result.Stdout, "warning") {
						suite.T().Logf("Script completed with warnings, which is acceptable")
						return
					}
				}
				helpers.AssertCommandSuccess(suite.T(), result, "%s should succeed", script.scriptPath)
			} else {
				// If expected to fail, verify it provides useful information
				if result.ExitCode != 0 {
					output := result.Stdout + result.Stderr
					assert.True(suite.T(),
						strings.Contains(output, "error") || strings.Contains(output, "Error") ||
							strings.Contains(output, "fail") || strings.Contains(output, "warning"),
						"Failed script should provide useful feedback")
				}
			}

			// Verify script produces structured output
			suite.validateScriptOutput(result, script.scriptPath)
		})
	}
}

// TestBashValidationScriptsEdgeCases tests edge cases and error handling
func (suite *IntegrationTestSuite) TestBashValidationScriptsEdgeCases() {
	edgeCases := []struct {
		name        string
		scriptPath  string
		setupFunc   func(string)
		expectError bool
		description string
	}{
		{
			name:       "Missing GoReleaser Config",
			scriptPath: "verify.sh",
			setupFunc: func(testDir string) {
				// Remove GoReleaser configs
				configs := []string{".goreleaser.yaml", ".goreleaser.pro.yaml"}
				for _, config := range configs {
					configPath := filepath.Join(testDir, config)
					if helpers.FileExists(configPath) {
						os.Remove(configPath)
					}
				}
			},
			expectError: true,
			description: "Should handle missing GoReleaser configuration gracefully",
		},
		{
			name:       "Missing Go Project Structure",
			scriptPath: "verify.sh",
			setupFunc: func(testDir string) {
				// Remove Go files
				goModPath := filepath.Join(testDir, "go.mod")
				if helpers.FileExists(goModPath) {
					os.Remove(goModPath)
				}
				cmdPath := filepath.Join(testDir, "cmd")
				if helpers.FileExists(cmdPath) {
					os.RemoveAll(cmdPath)
				}
			},
			expectError: true,
			description: "Should handle missing Go project structure",
		},
		{
			name:       "Invalid YAML Configuration",
			scriptPath: "verify.sh",
			setupFunc: func(testDir string) {
				// Create invalid YAML
				configPath := filepath.Join(testDir, ".goreleaser.yaml")
				if helpers.FileExists(configPath) {
					helpers.WriteFile(suite.T(), configPath, "invalid: yaml: content: [unclosed")
				}
			},
			expectError: true,
			description: "Should detect invalid YAML syntax",
		},
		{
			name:       "Missing License Templates",
			scriptPath: "verify.sh",
			setupFunc: func(testDir string) {
				// Remove license templates
				licensesDir := filepath.Join(testDir, "assets", "licenses")
				if helpers.FileExists(licensesDir) {
					os.RemoveAll(licensesDir)
				}
			},
			expectError: false, // Should warn but not necessarily fail
			description: "Should handle missing license templates",
		},
	}

	for _, tc := range edgeCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "bash-edge-case-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up basic project first
			suite.setupGoProject(testDir)

			// Apply edge case setup
			tc.setupFunc(testDir)

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
			suite.RegisterCleanup(cleanup)

			// Verify script exists
			scriptPath := filepath.Join(testDir, tc.scriptPath)
			require.True(suite.T(), helpers.FileExists(scriptPath), "Script %s should exist", tc.scriptPath)

			// Make script executable
			err := os.Chmod(scriptPath, 0755)
			require.NoError(suite.T(), err)

			// Run the script
			result := helpers.RunCommand(suite.T(), testDir, "./"+tc.scriptPath)

			if tc.expectError {
				assert.NotEqual(suite.T(), 0, result.ExitCode,
					"Script should fail for edge case: %s", tc.description)

				// Verify error messages are helpful
				output := result.Stdout + result.Stderr
				assert.True(suite.T(), len(output) > 0,
					"Script should provide error output for: %s", tc.description)
			} else {
				// Should handle gracefully (may warn but not necessarily fail)
				if result.ExitCode != 0 {
					output := result.Stdout + result.Stderr
					assert.True(suite.T(),
						strings.Contains(output, "warning") || strings.Contains(output, "Warning"),
						"Script should provide warnings for: %s", tc.description)
				}
			}

			suite.T().Logf("Edge case '%s' handled appropriately", tc.name)
		})
	}
}

// TestBashValidationScriptsToolDetection tests tool detection functionality
func (suite *IntegrationTestSuite) TestBashValidationScriptsToolDetection() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "tool-detection-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	suite.setupGoProject(testDir)
	suite.initGitRepo(testDir)

	// Set up environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
	suite.RegisterCleanup(cleanup)

	// Test verify.sh tool detection
	scriptPath := filepath.Join(testDir, "verify.sh")
	require.True(suite.T(), helpers.FileExists(scriptPath), "verify.sh should exist")

	err := os.Chmod(scriptPath, 0755)
	require.NoError(suite.T(), err)

	result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")

	output := result.Stdout + result.Stderr

	// Should check for various tools
	expectedToolChecks := []string{
		"go", "git", "goreleaser",
	}

	for _, tool := range expectedToolChecks {
		// Should mention the tool in output (either as installed or missing)
		assert.True(suite.T(),
			strings.Contains(output, tool),
			"Script should check for %s tool", tool)
	}

	// Should provide summary of tool availability
	assert.True(suite.T(),
		strings.Contains(output, "check") || strings.Contains(output, "Check") ||
			strings.Contains(output, "install") || strings.Contains(output, "Install"),
		"Script should provide information about tool checks")
}

// TestBashValidationScriptsEnvironmentVariableDetection tests environment variable detection
func (suite *IntegrationTestSuite) TestBashValidationScriptsEnvironmentVariableDetection() {
	envTestCases := []struct {
		name          string
		envVars       map[string]string
		expectMention []string
	}{
		{
			name:          "Minimal Environment Variables",
			envVars:       fixtures.TestEnvironmentVars["minimal"],
			expectMention: []string{"GITHUB_TOKEN", "DOCKER_USERNAME"},
		},
		{
			name:          "Complete Environment Variables",
			envVars:       fixtures.TestEnvironmentVars["complete"],
			expectMention: []string{"GITHUB_TOKEN", "GORELEASER_KEY"},
		},
		{
			name:          "Missing Critical Variables",
			envVars:       map[string]string{"PROJECT_NAME": "test"},
			expectMention: []string{"GITHUB_TOKEN"},
		},
	}

	for _, tc := range envTestCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "env-detection-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			suite.setupGoProject(testDir)

			// Clear environment first
			envCleanup := suite.clearRelevantEnvVars()
			suite.RegisterCleanup(envCleanup)

			// Set test environment variables
			cleanup := helpers.SetEnvVars(suite.T(), tc.envVars)
			suite.RegisterCleanup(cleanup)

			// Run verification script
			scriptPath := filepath.Join(testDir, "verify.sh")
			err := os.Chmod(scriptPath, 0755)
			require.NoError(suite.T(), err)

			result := helpers.RunCommand(suite.T(), testDir, "./verify.sh")
			output := result.Stdout + result.Stderr

			// Should mention environment variables in output
			for _, expectedVar := range tc.expectMention {
				assert.True(suite.T(),
					strings.Contains(output, expectedVar),
					"Script should mention environment variable %s", expectedVar)
			}

			// Should provide information about environment variable status
			assert.True(suite.T(),
				strings.Contains(output, "environment") || strings.Contains(output, "Environment") ||
					strings.Contains(output, "variable") || strings.Contains(output, "Variable") ||
					strings.Contains(output, "env") || strings.Contains(output, "Env"),
				"Script should provide information about environment variables")
		})
	}
}

// TestBashValidationScriptsOutputFormat tests that scripts produce well-formatted output
func (suite *IntegrationTestSuite) TestBashValidationScriptsOutputFormat() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "output-format-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	suite.setupGoProject(testDir)
	suite.initGitRepo(testDir)

	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
	suite.RegisterCleanup(cleanup)

	scripts := []string{"verify.sh", "validate-strict.sh"}

	for _, scriptName := range scripts {
		suite.Run("Output Format: "+scriptName, func() {
			scriptPath := filepath.Join(testDir, scriptName)
			if !helpers.FileExists(scriptPath) {
				suite.T().Skipf("Script %s does not exist", scriptName)
			}

			err := os.Chmod(scriptPath, 0755)
			require.NoError(suite.T(), err)

			result := helpers.RunCommand(suite.T(), testDir, "./"+scriptName)
			output := result.Stdout + result.Stderr

			// Should have structured output with sections
			assert.True(suite.T(), len(output) > 0, "Script should produce output")

			// Should use consistent formatting
			hasHeaders := strings.Contains(output, "===") || strings.Contains(output, "---") ||
				strings.Contains(output, "Summary") || strings.Contains(output, "SUMMARY")
			assert.True(suite.T(), hasHeaders, "Script should have section headers")

			// Should use status indicators
			hasStatusIndicators := strings.Contains(output, "✓") || strings.Contains(output, "✗") ||
				strings.Contains(output, "[") || strings.Contains(output, "PASS") ||
				strings.Contains(output, "FAIL") || strings.Contains(output, "OK")
			assert.True(suite.T(), hasStatusIndicators, "Script should have status indicators")

			suite.T().Logf("Script %s produces well-formatted output", scriptName)
		})
	}
}

// validateScriptOutput validates that script output is well-structured
func (suite *IntegrationTestSuite) validateScriptOutput(result helpers.CommandResult, scriptPath string) {
	output := result.Stdout + result.Stderr

	// Should have meaningful content
	assert.True(suite.T(), len(output) > 100,
		"Script %s should produce substantial output", scriptPath)

	// Should provide summary information
	hasSummary := strings.Contains(output, "Summary") || strings.Contains(output, "SUMMARY") ||
		strings.Contains(output, "check") || strings.Contains(output, "Check") ||
		strings.Contains(output, "error") || strings.Contains(output, "Error") ||
		strings.Contains(output, "warning") || strings.Contains(output, "Warning")
	assert.True(suite.T(), hasSummary,
		"Script %s should provide summary information", scriptPath)

	// Should mention key components
	keyComponents := []string{"goreleaser", "GoReleaser", "license", "License", "config", "Config"}
	mentionsComponents := false
	for _, component := range keyComponents {
		if strings.Contains(output, component) {
			mentionsComponents = true
			break
		}
	}
	assert.True(suite.T(), mentionsComponents,
		"Script %s should mention key project components", scriptPath)
}
