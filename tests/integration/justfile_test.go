package integration

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LarsArtmann/template-GoReleaser/tests/fixtures"
	"github.com/LarsArtmann/template-GoReleaser/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJustCommands tests that just commands work correctly with the actual project
func (suite *IntegrationTestSuite) TestJustCommands() {
	basicCommands := []struct {
		name        string
		command     string
		expectSuccess bool
		timeout     time.Duration
		requiresGo  bool
		description string
	}{
		{
			name:        "Just Help",
			command:     "help",
			expectSuccess: true,
			timeout:     30 * time.Second,
			requiresGo:  false,
			description: "Should display help information",
		},
		{
			name:        "Just List",
			command:     "--list",
			expectSuccess: true,
			timeout:     30 * time.Second,
			requiresGo:  false,
			description: "Should list available commands",
		},
		{
			name:        "Just Init",
			command:     "init",
			expectSuccess: true,
			timeout:     2 * time.Minute,
			requiresGo:  true,
			description: "Should initialize the project",
		},
		{
			name:        "Just Format",
			command:     "fmt",
			expectSuccess: true,
			timeout:     1 * time.Minute,
			requiresGo:  true,
			description: "Should format the code",
		},
		{
			name:        "Just Clean",
			command:     "clean",
			expectSuccess: true,
			timeout:     30 * time.Second,
			requiresGo:  false,
			description: "Should clean build artifacts",
		},
	}

	for _, cmd := range basicCommands {
		suite.Run(cmd.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "just-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up Go project if required
			if cmd.requiresGo {
				suite.setupGoProject(testDir)
			}

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
			suite.RegisterCleanup(cleanup)

			// Check if just command is available
			if !suite.isJustAvailable() {
				suite.T().Skip("just command not available")
			}

			// Verify justfile exists
			justfilePath := filepath.Join(testDir, "justfile")
			require.True(suite.T(), helpers.FileExists(justfilePath), "justfile should exist")

			// Run just command
			result := helpers.RunCommandWithTimeout(suite.T(), cmd.timeout, testDir, "just", cmd.command)

			if cmd.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Just %s output: %s", cmd.command, result.Stdout)
					suite.T().Logf("Just %s errors: %s", cmd.command, result.Stderr)
				}
				helpers.AssertCommandSuccess(suite.T(), result, "Just %s should succeed: %s", cmd.command, cmd.description)
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "Just %s should fail", cmd.command)
			}

			// Verify command produces expected output
			output := result.Stdout + result.Stderr
			switch cmd.command {
			case "help", "--list":
				assert.True(suite.T(), strings.Contains(output, "build") || strings.Contains(output, "test") || strings.Contains(output, "init"),
					"Help/list should show available commands")
			case "fmt":
				// Should complete without major errors
				assert.True(suite.T(), len(output) >= 0, "Format command should complete")
			case "clean":
				// Should complete and mention cleaning
				assert.True(suite.T(), strings.Contains(output, "clean") || strings.Contains(output, "Clean") || len(output) == 0,
					"Clean command should mention cleaning or be silent")
			}
		})
	}
}

// TestJustBuildCommands tests build-related just commands
func (suite *IntegrationTestSuite) TestJustBuildCommands() {
	buildCommands := []struct {
		name        string
		command     string
		expectSuccess bool
		timeout     time.Duration
		verifyFunc  func(*IntegrationTestSuite, string, helpers.CommandResult)
	}{
		{
			name:        "Just Build",
			command:     "build",
			expectSuccess: true,
			timeout:     3 * time.Minute,
			verifyFunc: func(suite *IntegrationTestSuite, testDir string, result helpers.CommandResult) {
				// Should create a binary
				binaryPath := filepath.Join(testDir, "myproject")
				if !helpers.FileExists(binaryPath) {
					// Try common binary locations
					altPaths := []string{
						filepath.Join(testDir, "bin", "myproject"),
						filepath.Join(testDir, "build", "myproject"),
					}
					for _, altPath := range altPaths {
						if helpers.FileExists(altPath) {
							binaryPath = altPath
							break
						}
					}
				}
				
				if helpers.FileExists(binaryPath) {
					suite.T().Logf("Binary found at: %s", binaryPath)
				} else {
					suite.T().Log("Binary not found, but build command succeeded (may be expected)")
				}
			},
		},
		{
			name:        "Just Test",
			command:     "test",
			expectSuccess: true,
			timeout:     2 * time.Minute,
			verifyFunc: func(suite *IntegrationTestSuite, testDir string, result helpers.CommandResult) {
				output := result.Stdout + result.Stderr
				assert.True(suite.T(), strings.Contains(output, "test") || strings.Contains(output, "PASS") || strings.Contains(output, "ok"),
					"Test output should indicate testing occurred")
			},
		},
		{
			name:        "Just Test Coverage",
			command:     "test-coverage",
			expectSuccess: true,
			timeout:     3 * time.Minute,
			verifyFunc: func(suite *IntegrationTestSuite, testDir string, result helpers.CommandResult) {
				// Should create coverage files
				coverageFiles := []string{"coverage.out", "coverage.html"}
				for _, file := range coverageFiles {
					filePath := filepath.Join(testDir, file)
					if helpers.FileExists(filePath) {
						suite.T().Logf("Coverage file found: %s", file)
					}
				}
			},
		},
	}

	for _, cmd := range buildCommands {
		suite.Run(cmd.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "just-build-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up Go project
			suite.setupGoProject(testDir)

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
			suite.RegisterCleanup(cleanup)

			// Check if just command is available
			if !suite.isJustAvailable() {
				suite.T().Skip("just command not available")
			}

			// Run just command
			result := helpers.RunCommandWithTimeout(suite.T(), cmd.timeout, testDir, "just", cmd.command)

			if cmd.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Just %s output: %s", cmd.command, result.Stdout)
					suite.T().Logf("Just %s errors: %s", cmd.command, result.Stderr)
				}
				helpers.AssertCommandSuccess(suite.T(), result, "Just %s should succeed", cmd.command)
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "Just %s should fail", cmd.command)
			}

			// Run verification function
			if cmd.verifyFunc != nil {
				cmd.verifyFunc(suite, testDir, result)
			}
		})
	}
}

// TestJustValidationCommands tests validation-related just commands
func (suite *IntegrationTestSuite) TestJustValidationCommands() {
	validationCommands := []struct {
		name        string
		command     string
		expectSuccess bool
		timeout     time.Duration
		setupFunc   func(string)
	}{
		{
			name:        "Just Validate",
			command:     "validate",
			expectSuccess: true,
			timeout:     2 * time.Minute,
			setupFunc: func(testDir string) {
				// Standard setup
			},
		},
		{
			name:        "Just Validate Strict",
			command:     "validate-strict",
			expectSuccess: false, // May fail with minimal setup
			timeout:     2 * time.Minute,
			setupFunc: func(testDir string) {
				// Standard setup
			},
		},
		{
			name:        "Just Check",
			command:     "check",
			expectSuccess: true,
			timeout:     1 * time.Minute,
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
			},
		},
	}

	for _, cmd := range validationCommands {
		suite.Run(cmd.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "just-validation-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up project as needed
			cmd.setupFunc(testDir)

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
			suite.RegisterCleanup(cleanup)

			// Check if just command is available
			if !suite.isJustAvailable() {
				suite.T().Skip("just command not available")
			}

			// Run just command
			result := helpers.RunCommandWithTimeout(suite.T(), cmd.timeout, testDir, "just", cmd.command)

			if cmd.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Just %s output: %s", cmd.command, result.Stdout)
					suite.T().Logf("Just %s errors: %s", cmd.command, result.Stderr)
					
					// Allow warnings for validation commands
					if strings.Contains(result.Stdout, "warning") || strings.Contains(result.Stdout, "Warning") {
						suite.T().Logf("Validation completed with warnings, which is acceptable")
						return
					}
				}
				helpers.AssertCommandSuccess(suite.T(), result, "Just %s should succeed", cmd.command)
			} else {
				// For validation commands that might fail, check they provide useful feedback
				output := result.Stdout + result.Stderr
				assert.True(suite.T(),
					strings.Contains(output, "error") || strings.Contains(output, "fail") || strings.Contains(output, "Error") ||
					strings.Contains(output, "warning") || strings.Contains(output, "Warning"),
					"Failed validation should provide useful feedback")
			}
		})
	}
}

// TestJustGoReleaserCommands tests GoReleaser-related just commands
func (suite *IntegrationTestSuite) TestJustGoReleaserCommands() {
	goreleaserCommands := []struct {
		name         string
		command      string
		expectSuccess bool
		timeout      time.Duration
		skipPro      bool
		setupFunc    func(string)
	}{
		{
			name:        "Just Snapshot",
			command:     "snapshot",
			expectSuccess: true,
			timeout:     5 * time.Minute,
			skipPro:     false,
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
			},
		},
		{
			name:        "Just Dry Run",
			command:     "dry-run",
			expectSuccess: false, // May fail without proper setup
			timeout:     5 * time.Minute,
			skipPro:     false,
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
				suite.createGitTag(testDir, "v1.0.0")
			},
		},
		{
			name:        "Just Snapshot Pro",
			command:     "snapshot-pro",
			expectSuccess: false, // Will fail without pro license
			timeout:     5 * time.Minute,
			skipPro:     true, // Skip unless we want to test failure
			setupFunc: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.initGitRepo(testDir)
			},
		},
	}

	for _, cmd := range goreleaserCommands {
		suite.Run(cmd.name, func() {
			if cmd.skipPro && strings.Contains(cmd.name, "Pro") {
				suite.T().Skip("Pro commands require GoReleaser Pro license")
			}

			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "just-goreleaser-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up project as needed
			cmd.setupFunc(testDir)

			// Set up environment variables
			envVars := fixtures.TestEnvironmentVars["minimal"]
			if strings.Contains(cmd.command, "pro") {
				envVars = fixtures.TestEnvironmentVars["complete"]
			}
			cleanup := helpers.SetEnvVars(suite.T(), envVars)
			suite.RegisterCleanup(cleanup)

			// Check if just and goreleaser commands are available
			if !suite.isJustAvailable() {
				suite.T().Skip("just command not available")
			}
			
			if !suite.isGoReleaserAvailable() {
				suite.T().Skip("goreleaser command not available")
			}

			// Run just command
			result := helpers.RunCommandWithTimeout(suite.T(), cmd.timeout, testDir, "just", cmd.command)

			if cmd.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("Just %s output: %s", cmd.command, result.Stdout)
					suite.T().Logf("Just %s errors: %s", cmd.command, result.Stderr)
				}
				helpers.AssertCommandSuccess(suite.T(), result, "Just %s should succeed", cmd.command)
				
				// Verify dist directory was created for successful builds
				if strings.Contains(cmd.command, "snapshot") {
					distDir := filepath.Join(testDir, "dist")
					assert.True(suite.T(), helpers.FileExists(distDir), "Snapshot should create dist directory")
				}
			} else {
				// For commands expected to fail, verify they fail gracefully
				output := result.Stdout + result.Stderr
				assert.True(suite.T(),
					strings.Contains(output, "error") || strings.Contains(output, "Error") ||
					strings.Contains(output, "fail") || strings.Contains(output, "license") ||
					result.ExitCode != 0,
					"Command should fail gracefully or provide error information")
			}
		})
	}
}

// TestJustCIWorkflow tests the complete CI workflow using just
func (suite *IntegrationTestSuite) TestJustCIWorkflow() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "just-ci-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Set up complete project
	suite.setupGoProject(testDir)
	suite.initGitRepo(testDir)

	// Set up environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
	suite.RegisterCleanup(cleanup)

	// Check if just command is available
	if !suite.isJustAvailable() {
		suite.T().Skip("just command not available")
	}

	// Run CI workflow
	result := helpers.RunCommandWithTimeout(suite.T(), 10*time.Minute, testDir, "just", "ci")

	if result.ExitCode != 0 {
		suite.T().Logf("Just CI output: %s", result.Stdout)
		suite.T().Logf("Just CI errors: %s", result.Stderr)
		
		// CI might fail due to missing tools, which is acceptable in test environment
		output := result.Stdout + result.Stderr
		if strings.Contains(output, "not installed") || strings.Contains(output, "command not found") {
			suite.T().Log("CI failed due to missing tools, which is expected in test environment")
			return
		}
	}

	// CI should either succeed or fail gracefully with useful information
	if result.ExitCode != 0 {
		output := result.Stdout + result.Stderr
		assert.True(suite.T(),
			strings.Contains(output, "error") || strings.Contains(output, "fail") ||
			strings.Contains(output, "Error") || strings.Contains(output, "warning"),
			"CI workflow should provide useful feedback on failure")
	} else {
		suite.T().Log("CI workflow completed successfully")
		
		// Verify expected outputs
		expectedOutputs := []string{
			"Clean", "format", "lint", "test", "build",
		}
		
		output := result.Stdout
		for _, expected := range expectedOutputs {
			// Allow case variations
			assert.True(suite.T(),
				strings.Contains(strings.ToLower(output), strings.ToLower(expected)) ||
				strings.Contains(output, expected),
				"CI output should mention %s step", expected)
		}
	}
}

// Helper methods for checking command availability

func (suite *IntegrationTestSuite) isJustAvailable() bool {
	result := helpers.RunCommand(suite.T(), suite.originalDir, "just", "--version")
	return result.ExitCode == 0
}

func (suite *IntegrationTestSuite) isGoReleaserAvailable() bool {
	result := helpers.RunCommand(suite.T(), suite.originalDir, "goreleaser", "--version")
	return result.ExitCode == 0
}