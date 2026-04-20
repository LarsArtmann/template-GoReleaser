package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeGoModFile writes the go.mod file content.
func writeGoModFile(t *testing.T, dir, moduleName string) {
	t.Helper()

	goMod := "module " + moduleName + "\ngo 1.21\n"

	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
}

// createBasicGoProject creates a basic Go project structure for testing.
func createBasicGoProject(t *testing.T, dir, moduleName string) {
	t.Helper()

	writeGoModFile(t, dir, moduleName)

	mainContent := "package main\n\nfunc main() {}\n"

	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainContent), 0o644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}
}

// createGoProjectWithFiles creates a Go project with custom files.
func createGoProjectWithFiles(t *testing.T, dir, moduleName string, files map[string]string) {
	t.Helper()

	writeGoModFile(t, dir, moduleName)

	for path, content := range files {
		fullPath := filepath.Join(dir, path)

		err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", path, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("failed to create %s: %v", path, err)
		}
	}
}

// inTestDir executes the test function after changing to the test directory.
// It automatically restores the original directory after the test completes.
func inTestDir(t *testing.T, testDir string, testFunc func()) {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	defer func() {
		chdirErr := os.Chdir(originalDir)
		if chdirErr != nil {
			t.Logf("failed to restore directory: %v", chdirErr)
		}
	}()

	testFunc()
}

// verifyFileContains checks if a file contains expected strings.
func verifyFileContains(t *testing.T, filePath string, checks []string) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", filePath, err)
	}

	contentStr := string(content)
	for _, check := range checks {
		if !strings.Contains(contentStr, check) {
			t.Errorf("file %s missing expected string: %q", filePath, check)
		}
	}
}

// setupMinimalProject creates a minimal project for testing.
func setupMinimalProject(pattern string) (string, func()) {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		return "", func() {}
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}
