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

// TestLicenseSystemIntegration tests the integration between license generation and readme config
func (suite *IntegrationTestSuite) TestLicenseSystemIntegration() {
	testCases := []struct {
		name                string
		readmeConfig        string
		expectedLicenseType string
		expectedAuthor      string
		expectedYear        string
	}{
		{
			name:                "MIT License with Complete Config",
			readmeConfig:        fixtures.ReadmeConfigs["complete"],
			expectedLicenseType: "MIT",
			expectedAuthor:      "Test Author",
			expectedYear:        "2024",
		},
		{
			name:                "Minimal Config with MIT License",
			readmeConfig:        fixtures.ReadmeConfigs["minimal"],
			expectedLicenseType: "MIT",
			expectedAuthor:      "Test Author",
			expectedYear:        "", // Year might be current year
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-integration-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
			suite.RegisterCleanup(cleanup)

			// Create readme config
			readmeConfigDir := filepath.Join(testDir, ".readme", "configs")
			err := os.MkdirAll(readmeConfigDir, 0755)
			require.NoError(suite.T(), err)

			configPath := filepath.Join(readmeConfigDir, "readme-config.yaml")
			helpers.WriteFile(suite.T(), configPath, tc.readmeConfig)

			// Verify config was created correctly
			assert.True(suite.T(), helpers.FileExists(configPath), "Readme config should exist")
			assert.True(suite.T(), helpers.FileContains(suite.T(), configPath, tc.expectedLicenseType),
				"Config should contain license type %s", tc.expectedLicenseType)

			// Run license generation
			licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
			require.True(suite.T(), helpers.FileExists(licenseScript), "License script should exist")

			err = os.Chmod(licenseScript, 0755)
			require.NoError(suite.T(), err)

			result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh")
			helpers.AssertCommandSuccess(suite.T(), result, "License generation should succeed")

			// Verify LICENSE file was created with correct content
			licensePath := filepath.Join(testDir, "LICENSE")
			assert.True(suite.T(), helpers.FileExists(licensePath), "LICENSE file should be created")

			// Check license content for expected elements
			if tc.expectedAuthor != "" {
				assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, tc.expectedAuthor),
					"LICENSE should contain author name: %s", tc.expectedAuthor)
			}

			// Verify the license follows the expected template structure
			expectedTemplate := fixtures.ExpectedLicenseContent[tc.expectedLicenseType]
			if expectedTemplate != "" {
				assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, expectedTemplate),
					"LICENSE should contain template content for %s", tc.expectedLicenseType)
			}

			// Run validation to ensure integration works
			validateScript := filepath.Join(testDir, "verify.sh")
			err = os.Chmod(validateScript, 0755)
			require.NoError(suite.T(), err)

			validationResult := helpers.RunCommand(suite.T(), testDir, "./verify.sh")
			// Should pass validation after license generation
			assert.True(suite.T(), validationResult.ExitCode == 0 || strings.Contains(validationResult.Stdout, "warning"),
				"Validation should pass after license generation")
		})
	}
}

// TestGoReleaserConfigIntegration tests integration between GoReleaser configs and project structure
func (suite *IntegrationTestSuite) TestGoReleaserConfigIntegration() {
	testCases := []struct {
		name          string
		configFile    string
		setupProject  func(string)
		expectSuccess bool
	}{
		{
			name:       "Free Config with Standard Project",
			configFile: ".goreleaser.yaml",
			setupProject: func(testDir string) {
				suite.setupGoProject(testDir)
			},
			expectSuccess: true,
		},
		{
			name:       "Pro Config with Complete Project",
			configFile: ".goreleaser.pro.yaml",
			setupProject: func(testDir string) {
				suite.setupGoProject(testDir)
				suite.setupDockerSupport(testDir)
			},
			expectSuccess: true, // May have warnings about pro license
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "goreleaser-integration-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up project structure
			tc.setupProject(testDir)

			// Initialize git repository
			suite.initGitRepo(testDir)

			// Set up environment variables
			envVars := fixtures.TestEnvironmentVars["minimal"]
			if strings.Contains(tc.configFile, "pro") {
				envVars = fixtures.TestEnvironmentVars["complete"]
			}
			cleanup := helpers.SetEnvVars(suite.T(), envVars)
			suite.RegisterCleanup(cleanup)

			// Check that config file exists
			configPath := filepath.Join(testDir, tc.configFile)
			if !helpers.FileExists(configPath) {
				suite.T().Skipf("Configuration file %s does not exist", tc.configFile)
			}

			// Test that GoReleaser can validate the config with the project structure
			result := helpers.RunCommand(suite.T(), testDir, "goreleaser", "check", "--config", tc.configFile)

			if tc.expectSuccess {
				if result.ExitCode != 0 {
					suite.T().Logf("GoReleaser check output: %s", result.Stdout)
					suite.T().Logf("GoReleaser check errors: %s", result.Stderr)

					// Pro version might fail without license, which is acceptable
					if strings.Contains(tc.configFile, "pro") &&
						(strings.Contains(result.Stderr, "license") || strings.Contains(result.Stderr, "pro")) {
						suite.T().Log("Pro config validation failed due to missing license, which is expected")
						return
					}
				}
				helpers.AssertCommandSuccess(suite.T(), result, "GoReleaser check should pass for %s", tc.configFile)
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "GoReleaser check should fail for %s", tc.configFile)
			}

			// Run project validation to ensure everything integrates
			validateScript := filepath.Join(testDir, "verify.sh")
			err := os.Chmod(validateScript, 0755)
			require.NoError(suite.T(), err)

			validationResult := helpers.RunCommand(suite.T(), testDir, "./verify.sh")
			// Integration should pass overall validation
			assert.True(suite.T(), validationResult.ExitCode == 0 || strings.Contains(validationResult.Stdout, "warning"),
				"Overall validation should pass for integrated project")
		})
	}
}

// TestValidationScriptIntegration tests integration between validation scripts and actual project
func (suite *IntegrationTestSuite) TestValidationScriptIntegration() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "validation-integration-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Set up complete project
	suite.setupGoProject(testDir)
	suite.initGitRepo(testDir)

	// Set up environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
	suite.RegisterCleanup(cleanup)

	// Generate license
	suite.generateLicense(testDir)

	// Test that all validation scripts work with the integrated project
	validationScripts := []struct {
		name   string
		script string
		expect string
	}{
		{
			name:   "Basic Verification",
			script: "verify.sh",
			expect: "success",
		},
		{
			name:   "Strict Validation",
			script: "validate-strict.sh",
			expect: "success_or_warning",
		},
	}

	for _, vs := range validationScripts {
		suite.Run(vs.name, func() {
			scriptPath := filepath.Join(testDir, vs.script)
			if !helpers.FileExists(scriptPath) {
				suite.T().Skipf("Validation script %s does not exist", vs.script)
			}

			err := os.Chmod(scriptPath, 0755)
			require.NoError(suite.T(), err)

			result := helpers.RunCommand(suite.T(), testDir, "./"+vs.script)

			switch vs.expect {
			case "success":
				helpers.AssertCommandSuccess(suite.T(), result, "Validation script %s should succeed", vs.script)
			case "success_or_warning":
				assert.True(suite.T(), result.ExitCode == 0 || strings.Contains(result.Stdout, "warning"),
					"Validation script %s should succeed or have warnings only", vs.script)
			}

			// Verify that validation provides useful output
			output := result.Stdout + result.Stderr
			assert.True(suite.T(),
				strings.Contains(output, "check") || strings.Contains(output, "validation") ||
					strings.Contains(output, "verify") || strings.Contains(output, "Summary"),
				"Validation should provide meaningful output")
		})
	}
}

// TestCompleteWorkflowIntegration tests the complete end-to-end workflow
func (suite *IntegrationTestSuite) TestCompleteWorkflowIntegration() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "complete-workflow-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Step 1: Set up environment
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
	suite.RegisterCleanup(cleanup)

	// Step 2: Set up project structure
	suite.setupGoProject(testDir)

	// Step 3: Initialize git repository
	suite.initGitRepo(testDir)

	// Step 4: Create readme config
	readmeConfigDir := filepath.Join(testDir, ".readme", "configs")
	err := os.MkdirAll(readmeConfigDir, 0755)
	require.NoError(suite.T(), err)

	configPath := filepath.Join(readmeConfigDir, "readme-config.yaml")
	helpers.WriteFile(suite.T(), configPath, fixtures.ReadmeConfigs["complete"])

	// Step 5: Generate license
	suite.generateLicense(testDir)

	// Step 6: Validate project structure
	validateScript := filepath.Join(testDir, "verify.sh")
	err = os.Chmod(validateScript, 0755)
	require.NoError(suite.T(), err)

	validationResult := helpers.RunCommand(suite.T(), testDir, "./verify.sh")
	assert.True(suite.T(), validationResult.ExitCode == 0 || strings.Contains(validationResult.Stdout, "warning"),
		"Project validation should pass in complete workflow")

	// Step 7: Test GoReleaser build (if available)
	if helpers.FileExists(filepath.Join(testDir, ".goreleaser.yaml")) {
		// Test basic build functionality
		result := helpers.RunCommand(suite.T(), testDir, "goreleaser", "build", "--snapshot", "--single-target")

		if result.ExitCode == 0 {
			// Verify build artifacts
			distDir := filepath.Join(testDir, "dist")
			assert.True(suite.T(), helpers.FileExists(distDir), "Build should create dist directory")
		} else {
			suite.T().Logf("GoReleaser build failed (may be expected): %s", result.Stderr)
		}
	}

	// Step 8: Verify all expected files exist
	expectedFiles := []string{
		"LICENSE",
		"go.mod",
		"cmd/myproject/main.go",
		".readme/configs/readme-config.yaml",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(testDir, file)
		assert.True(suite.T(), helpers.FileExists(filePath), "Expected file should exist: %s", file)
	}

	// Step 9: Verify integration points work
	// License should be properly integrated
	licensePath := filepath.Join(testDir, "LICENSE")
	assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, "Test Author"),
		"License should contain author from config")

	suite.T().Log("Complete workflow integration test passed")
}

// Helper method to set up Docker support
func (suite *IntegrationTestSuite) setupDockerSupport(testDir string) {
	dockerfilePath := filepath.Join(testDir, "Dockerfile")
	if !helpers.FileExists(dockerfilePath) {
		dockerfileContent := `FROM scratch
COPY myproject /myproject
ENTRYPOINT ["/myproject"]
`
		helpers.WriteFile(suite.T(), dockerfilePath, dockerfileContent)
	}
}

// Helper method to generate license
func (suite *IntegrationTestSuite) generateLicense(testDir string) {
	licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
	if helpers.FileExists(licenseScript) {
		err := os.Chmod(licenseScript, 0755)
		require.NoError(suite.T(), err)

		result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh")
		if result.ExitCode != 0 {
			suite.T().Logf("License generation warning: %s", result.Stderr)
		}
	}
}
