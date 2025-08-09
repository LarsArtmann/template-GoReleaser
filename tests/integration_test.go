// Package tests provides integration tests for the GoReleaser template project.
// This file serves as the main test entry point and provides package-level test configuration.
package tests

import (
	"testing"
	
	// Import the integration test package to run the test suite
	_ "github.com/LarsArtmann/template-GoReleaser/tests/integration"
)

// TestMain provides package-level test setup and teardown if needed
func TestMain(m *testing.M) {
	// Run all tests
	m.Run()
}