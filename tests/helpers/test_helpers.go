package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// CommandResult holds the result of a command execution
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// RunCommand executes a command and returns the result
func RunCommand(t *testing.T, dir, command string, args ...string) CommandResult {
	t.Helper()

	start := time.Now()
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
}

// RunCommandWithTimeout executes a command with a timeout
func RunCommandWithTimeout(t *testing.T, timeout time.Duration, dir, command string, args ...string) CommandResult {
	t.Helper()

	start := time.Now()
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	var err error
	select {
	case err = <-done:
		// Command completed normally
	case <-time.After(timeout):
		// Command timed out
		if cmd.Process != nil {
			_ = cmd.Process.Kill() // Error handling not critical for cleanup
		}
		err = fmt.Errorf("command timed out after %v", timeout)
	}

	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
}

// CopyDir recursively copies a directory
func CopyDir(t *testing.T, src, dst string) {
	t.Helper()

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip certain directories that shouldn't be copied for tests
		if strings.Contains(relPath, "tests") || 
		   strings.Contains(relPath, ".git") ||
		   strings.Contains(relPath, "dist") ||
		   strings.Contains(relPath, "build") ||
		   strings.Contains(relPath, "vendor") ||
		   strings.Contains(relPath, "node_modules") ||
		   strings.Contains(relPath, ".DS_Store") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dstPath := filepath.Join(dst, relPath)

		// Create directory if it's a directory
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file if it's a regular file
		return copyFile(path, dstPath, info.Mode())
	})

	require.NoError(t, err, "Failed to copy directory from %s to %s", src, dst)
}

// copyFile copies a single file
func copyFile(src, dst string, mode os.FileMode) error {
	// #nosec G304 - src is validated by caller in test environment
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return err
	}

	// #nosec G304 - dst is validated by caller in test environment
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}

// CreateTestProject creates a temporary test project based on the template
func CreateTestProject(t *testing.T, templateDir, projectName string) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", fmt.Sprintf("test-%s-*", projectName))
	require.NoError(t, err, "Failed to create temporary directory")

	// Copy template to temporary directory
	CopyDir(t, templateDir, tempDir)

	return tempDir
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileContains checks if a file contains a specific string
func FileContains(t *testing.T, path, content string) bool {
	t.Helper()

	// #nosec G304 - path is validated by caller in test environment
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), content)
}

// WriteFile writes content to a file
func WriteFile(t *testing.T, path, content string) {
	t.Helper()

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0750)
	require.NoError(t, err, "Failed to create directory for file %s", path)

	err = os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err, "Failed to write file %s", path)
}

// SetEnvVars sets multiple environment variables and returns a cleanup function
func SetEnvVars(t *testing.T, vars map[string]string) func() {
	t.Helper()

	originalVars := make(map[string]string)

	// Store original values and set new ones
	for key, value := range vars {
		if original, exists := os.LookupEnv(key); exists {
			originalVars[key] = original
		}
		_ = os.Setenv(key, value) // Error handling not critical in tests
	}

	// Return cleanup function
	return func() {
		for key := range vars {
			if original, exists := originalVars[key]; exists {
				_ = os.Setenv(key, original) // Error handling not critical in tests
			} else {
				os.Unsetenv(key)
			}
		}
	}
}

// AssertCommandSuccess asserts that a command executed successfully
func AssertCommandSuccess(t *testing.T, result CommandResult, msgAndArgs ...interface{}) {
	t.Helper()

	if result.ExitCode != 0 {
		t.Errorf("Command failed with exit code %d\nStdout: %s\nStderr: %s",
			result.ExitCode, result.Stdout, result.Stderr)
		if len(msgAndArgs) > 0 {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		}
		t.FailNow()
	}
}
