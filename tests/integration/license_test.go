package integration

import (
	"os"
	"path/filepath"

	"github.com/LarsArtmann/template-GoReleaser/tests/fixtures"
	"github.com/LarsArtmann/template-GoReleaser/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLicenseGeneration tests the complete license generation workflow
func (suite *IntegrationTestSuite) TestLicenseGeneration() {
	testCases := []struct {
		name            string
		licenseType     string
		readmeConfig    string
		envVars         map[string]string
		expectSuccess   bool
		expectedContent string
	}{
		{
			name:            "MIT License Generation",
			licenseType:     "MIT",
			readmeConfig:    fixtures.ReadmeConfigs["complete"],
			envVars:         fixtures.TestEnvironmentVars["complete"],
			expectSuccess:   true,
			expectedContent: fixtures.ExpectedLicenseContent["MIT"],
		},
		{
			name:            "Apache-2.0 License Generation",
			licenseType:     "Apache-2.0",
			readmeConfig:    fixtures.ReadmeConfigs["complete"],
			envVars:         fixtures.TestEnvironmentVars["complete"],
			expectSuccess:   true,
			expectedContent: fixtures.ExpectedLicenseContent["Apache-2.0"],
		},
		{
			name:            "BSD-3-Clause License Generation",
			licenseType:     "BSD-3-Clause",
			readmeConfig:    fixtures.ReadmeConfigs["complete"],
			envVars:         fixtures.TestEnvironmentVars["complete"],
			expectSuccess:   true,
			expectedContent: fixtures.ExpectedLicenseContent["BSD-3-Clause"],
		},
		{
			name:            "Invalid License Type",
			licenseType:     "INVALID",
			readmeConfig:    fixtures.ReadmeConfigs["minimal"],
			envVars:         fixtures.TestEnvironmentVars["minimal"],
			expectSuccess:   false,
			expectedContent: "",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create test project
			testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-test")
			suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

			// Set up environment variables
			cleanup := helpers.SetEnvVars(suite.T(), tc.envVars)
			suite.RegisterCleanup(cleanup)

			// Create readme config with the specific license type
			readmeConfigDir := filepath.Join(testDir, ".readme", "configs")
			err := os.MkdirAll(readmeConfigDir, 0755)
			require.NoError(suite.T(), err)

			// Modify readme config to use the test license type
			modifiedConfig := tc.readmeConfig
			if tc.licenseType != "" {
				// Replace the license type in the config
				modifiedConfig = `project:
  name: test-project
  description: Test project for integration testing

author:
  name: Test Author
  email: test@example.com

license:
  type: ` + tc.licenseType + `

repository:
  url: https://github.com/testuser/test-project
`
			}

			configPath := filepath.Join(readmeConfigDir, "readme-config.yaml")
			helpers.WriteFile(suite.T(), configPath, modifiedConfig)

			// Check if license generation script exists and is executable
			licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
			require.True(suite.T(), helpers.FileExists(licenseScript), "License generation script should exist")

			// Make script executable
			err = os.Chmod(licenseScript, 0755)
			require.NoError(suite.T(), err)

			// Run license generation
			result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh")

			if tc.expectSuccess {
				helpers.AssertCommandSuccess(suite.T(), result, "License generation should succeed for %s", tc.licenseType)

				// Verify LICENSE file was created
				licensePath := filepath.Join(testDir, "LICENSE")
				assert.True(suite.T(), helpers.FileExists(licensePath), "LICENSE file should be created")

				// Verify LICENSE file contains expected content
				if tc.expectedContent != "" {
					assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, tc.expectedContent),
						"LICENSE file should contain expected content for %s", tc.licenseType)
				}

				// Verify license file is not empty
				stat, err := os.Stat(licensePath)
				require.NoError(suite.T(), err)
				assert.Greater(suite.T(), stat.Size(), int64(100), "LICENSE file should have substantial content")
			} else {
				assert.NotEqual(suite.T(), 0, result.ExitCode, "License generation should fail for invalid license type")
			}
		})
	}
}

// TestLicenseScriptHelp tests the license script help functionality
func (suite *IntegrationTestSuite) TestLicenseScriptHelp() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-help-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Check if license generation script exists
	licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
	require.True(suite.T(), helpers.FileExists(licenseScript), "License generation script should exist")

	// Make script executable
	err := os.Chmod(licenseScript, 0755)
	require.NoError(suite.T(), err)

	// Test help functionality
	result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh", "--help")
	helpers.AssertCommandSuccess(suite.T(), result, "License script help should work")

	// Verify help output contains expected information
	assert.Contains(suite.T(), result.Stdout, "Usage:", "Help output should contain usage information")
	assert.Contains(suite.T(), result.Stdout, "license", "Help output should mention license")
}

// TestLicenseScriptList tests the license script list functionality
func (suite *IntegrationTestSuite) TestLicenseScriptList() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-list-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Check if license generation script exists
	licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
	require.True(suite.T(), helpers.FileExists(licenseScript), "License generation script should exist")

	// Make script executable
	err := os.Chmod(licenseScript, 0755)
	require.NoError(suite.T(), err)

	// Test list functionality
	result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh", "--list")
	helpers.AssertCommandSuccess(suite.T(), result, "License script list should work")

	// Verify list output contains expected license types
	expectedLicenses := []string{"MIT", "Apache-2.0", "BSD-3-Clause", "EUPL-1.2"}
	for _, license := range expectedLicenses {
		assert.Contains(suite.T(), result.Stdout, license, "License list should contain %s", license)
	}
}

// TestLicenseTemplatesExist tests that all license templates exist
func (suite *IntegrationTestSuite) TestLicenseTemplatesExist() {
	expectedTemplates := []string{
		"MIT.template",
		"Apache-2.0.template",
		"BSD-3-Clause.template",
		"EUPL-1.2.template",
	}

	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-templates-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	licensesDir := filepath.Join(testDir, "assets", "licenses")
	require.True(suite.T(), helpers.FileExists(licensesDir), "License templates directory should exist")

	for _, template := range expectedTemplates {
		templatePath := filepath.Join(licensesDir, template)
		assert.True(suite.T(), helpers.FileExists(templatePath), "License template %s should exist", template)

		// Verify template file is not empty
		stat, err := os.Stat(templatePath)
		require.NoError(suite.T(), err)
		assert.Greater(suite.T(), stat.Size(), int64(100), "License template %s should have substantial content", template)
	}
}

// TestLicenseGenerationWithoutReadmeConfig tests license generation without readme config
func (suite *IntegrationTestSuite) TestLicenseGenerationWithoutReadmeConfig() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-no-config-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Set up minimal environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["minimal"])
	suite.RegisterCleanup(cleanup)

	// Ensure no readme config exists
	readmeConfigPath := filepath.Join(testDir, ".readme", "configs", "readme-config.yaml")
	if helpers.FileExists(readmeConfigPath) {
		err := os.Remove(readmeConfigPath)
		require.NoError(suite.T(), err)
	}

	// Check if license generation script exists
	licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
	require.True(suite.T(), helpers.FileExists(licenseScript), "License generation script should exist")

	// Make script executable
	err := os.Chmod(licenseScript, 0755)
	require.NoError(suite.T(), err)

	// Try to run license generation without config
	result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh")

	// Should fail gracefully or use defaults
	if result.ExitCode != 0 {
		assert.Contains(suite.T(), result.Stderr, "config", "Error should mention missing config")
	} else {
		// If it succeeds, verify it created some kind of license
		licensePath := filepath.Join(testDir, "LICENSE")
		assert.True(suite.T(), helpers.FileExists(licensePath), "LICENSE file should be created even without config")
	}
}

// TestLicenseBackupAndRestore tests that license generation properly handles existing LICENSE files
func (suite *IntegrationTestSuite) TestLicenseBackupAndRestore() {
	// Create test project
	testDir := helpers.CreateTestProject(suite.T(), suite.originalDir, "license-backup-test")
	suite.RegisterCleanup(func() { os.RemoveAll(testDir) })

	// Set up environment variables
	cleanup := helpers.SetEnvVars(suite.T(), fixtures.TestEnvironmentVars["complete"])
	suite.RegisterCleanup(cleanup)

	// Create readme config
	readmeConfigDir := filepath.Join(testDir, ".readme", "configs")
	err := os.MkdirAll(readmeConfigDir, 0755)
	require.NoError(suite.T(), err)

	configPath := filepath.Join(readmeConfigDir, "readme-config.yaml")
	helpers.WriteFile(suite.T(), configPath, fixtures.ReadmeConfigs["complete"])

	// Create an existing LICENSE file
	existingLicenseContent := "This is an existing LICENSE file that should be preserved if generation fails"
	licensePath := filepath.Join(testDir, "LICENSE")
	helpers.WriteFile(suite.T(), licensePath, existingLicenseContent)

	// Verify the existing LICENSE file
	assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, existingLicenseContent))

	// Check if license generation script exists
	licenseScript := filepath.Join(testDir, "scripts", "generate-license.sh")
	require.True(suite.T(), helpers.FileExists(licenseScript), "License generation script should exist")

	// Make script executable
	err = os.Chmod(licenseScript, 0755)
	require.NoError(suite.T(), err)

	// Run license generation
	result := helpers.RunCommand(suite.T(), testDir, "./scripts/generate-license.sh")
	helpers.AssertCommandSuccess(suite.T(), result, "License generation should succeed")

	// Verify LICENSE file was updated (should not contain the old content)
	assert.False(suite.T(), helpers.FileContains(suite.T(), licensePath, existingLicenseContent),
		"LICENSE file should be updated with new content")

	// Verify LICENSE file contains new MIT content
	assert.True(suite.T(), helpers.FileContains(suite.T(), licensePath, fixtures.ExpectedLicenseContent["MIT"]),
		"LICENSE file should contain MIT license content")
}
