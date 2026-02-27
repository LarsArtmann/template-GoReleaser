package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// validateGoReleaserConfig validates GoReleaser configuration.
func validateGoReleaserConfig(results *ValidationResults) error {
	configPath := ".goreleaser.yaml"

	// Check if config exists
	exists, err := fileSystemRepo.FileExists(context.Background(), configPath)
	if err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check configuration file",
				"Cannot access "+configPath,
				err,
			).WithContext(configPath))

		return nil
	}

	results.ConfigExists = exists
	if !exists {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_NOT_FOUND",
				"Configuration file not found",
				configPath+" does not exist",
				nil,
			).WithContext(configPath))
		results.Recommendations = append(results.Recommendations,
			"Run 'goreleaser-wizard init' to create configuration")

		return nil
	}

	// Check if config is valid YAML
	if err := validateYAML(configPath, results); err != nil {
		results.ConfigValid = false

		return nil
	}

	// Check if goreleaser is available
	if _, err := exec.LookPath("goreleaser"); err == nil {
		results.GoReleaserFound = true
		// Run goreleaser check
		cmd := exec.Command("goreleaser", "check", "--config", configPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			results.ConfigValid = false
			results.Errors = append(results.Errors,
				NewValidationError(
					"INVALID_CONFIGURATION",
					"GoReleaser configuration is invalid",
					string(output),
				).WithContext(configPath))
		} else {
			results.ConfigValid = true
		}
	} else {
		results.Warnings = append(results.Warnings,
			NewValidationError(
				"EXTERNAL_TOOL_NOT_FOUND",
				"GoReleaser not found",
				"Install GoReleaser for configuration validation",
			))
		// Assume valid if no goreleaser available
		results.ConfigValid = true
	}

	return nil
}

// validateGitHubActions validates GitHub Actions workflow.
func validateGitHubActions(results *ValidationResults) error {
	workflowPath := filepath.Join(".github", "workflows", "release.yml")

	// Check if workflow exists
	exists, err := fileSystemRepo.FileExists(context.Background(), workflowPath)
	if err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check workflow file",
				"Cannot access "+workflowPath,
				err,
			).WithContext(workflowPath))

		return nil
	}

	results.ActionsExists = exists
	if !exists {
		results.Recommendations = append(results.Recommendations,
			"Consider adding GitHub Actions workflow for automated releases")

		return nil
	}

	// Validate workflow content
	if err := validateWorkflowContent(workflowPath, results); err != nil {
		results.ActionsValid = false

		return nil
	}

	results.ActionsValid = true

	return nil
}

// validateProjectStructure validates project structure.
func validateProjectStructure(results *ValidationResults) error {
	// Check for go.mod
	if exists, err := fileSystemRepo.FileExists(context.Background(), "go.mod"); err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check go.mod",
				"Cannot access go.mod file",
				err,
			).WithContext("go.mod"))

		return nil
	} else if !exists {
		results.Errors = append(results.Errors,
			NewValidationError(
				"MISSING_DEPENDENCY",
				"go.mod not found",
				"Run 'go mod init' to initialize Go module",
			).WithContext("go.mod"))
		results.ProjectValid = false

		return nil
	}

	// Check for main.go directory
	mainPath := "cmd"
	if exists, err := fileSystemRepo.DirExists(context.Background(), mainPath); err != nil {
		results.Errors = append(results.Errors,
			NewSystemError(
				"FILE_READ_FAILED",
				"Failed to check main directory",
				"Cannot access "+mainPath,
				err,
			).WithContext(mainPath))

		return nil
	} else if !exists {
		// Try common alternatives
		alternatives := []string{"main", "src"}
		for _, alt := range alternatives {
			if exists, _ := fileSystemRepo.DirExists(context.Background(), alt); exists {
				results.Warnings = append(results.Warnings,
					NewValidationError(
						"INVALID_PROJECT_STRUCTURE",
						"Non-standard main directory",
						fmt.Sprintf("Found %s instead of %s", alt, mainPath),
					).WithContext(alt))
				results.ProjectValid = true

				return nil
			}
		}

		results.Errors = append(results.Errors,
			NewValidationError(
				"INVALID_PROJECT_STRUCTURE",
				"Main directory not found",
				"Create cmd directory with main package",
			).WithContext(mainPath))
		results.ProjectValid = false

		return nil
	}

	results.ProjectValid = true

	return nil
}

// validateYAML validates YAML file.
func validateYAML(filePath string, results *ValidationResults) error {
	// Simple YAML validation - try to read file content
	content, err := fileSystemRepo.ReadFile(context.Background(), filePath)
	if err != nil {
		return err
	}

	// Basic YAML structure check
	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 {
		results.Errors = append(results.Errors,
			NewValidationError(
				"INVALID_FILE_FORMAT",
				"Empty YAML file",
				filePath+" is empty",
			).WithContext(filePath))

		return nil
	}

	// Check for common YAML structure indicators
	hasYamlStructure := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "#") {
			hasYamlStructure = true

			break
		}
	}

	if !hasYamlStructure {
		results.Errors = append(results.Errors,
			NewValidationError(
				"INVALID_FILE_FORMAT",
				"Invalid YAML structure",
				filePath+" does not appear to be valid YAML",
			).WithContext(filePath))

		return nil
	}

	return nil
}
