package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/template-GoReleaser/tests/helpers"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite defines the integration test suite
type IntegrationTestSuite struct {
	suite.Suite
	originalDir  string
	tempDir      string
	cleanupFuncs []func()
}

// SetupSuite runs once before all tests in the suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Store original working directory and find project root
	var err error
	cwd, err := os.Getwd()
	suite.Require().NoError(err)

	// Navigate to project root (go up from tests/integration to project root)
	suite.originalDir = filepath.Join(cwd, "..", "..")

	// Verify we're at the project root by checking for key files
	goModPath := filepath.Join(suite.originalDir, "go.mod")
	suite.Require().True(helpers.FileExists(goModPath), "Could not find project root - go.mod not found at %s", goModPath)

	// Create temporary test directory
	suite.tempDir, err = os.MkdirTemp("", "goreleaser-template-test-*")
	suite.Require().NoError(err)

	suite.T().Logf("Running integration tests in: %s", suite.tempDir)
	suite.T().Logf("Project root: %s", suite.originalDir)
}

// SetupTest runs before each individual test
func (suite *IntegrationTestSuite) SetupTest() {
	// Change to the original directory before each test
	err := os.Chdir(suite.originalDir)
	suite.Require().NoError(err)
}

// TearDownTest runs after each individual test
func (suite *IntegrationTestSuite) TearDownTest() {
	// Run any cleanup functions registered during the test
	for _, cleanup := range suite.cleanupFuncs {
		cleanup()
	}
	suite.cleanupFuncs = nil
}

// TearDownSuite runs once after all tests in the suite
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}

	// Restore original working directory
	if suite.originalDir != "" {
		os.Chdir(suite.originalDir)
	}
}

// RegisterCleanup adds a cleanup function to be run after the current test
func (suite *IntegrationTestSuite) RegisterCleanup(cleanup func()) {
	suite.cleanupFuncs = append(suite.cleanupFuncs, cleanup)
}

// TestIntegrationSuite runs the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
