package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/template-GoReleaser/internal/domain"
)

// validateWorkflowContent validates GitHub Actions workflow content
func validateWorkflowContent(workflowPath string, results *ValidationResults) error {
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to read workflow file",
			fmt.Sprintf("Cannot read %s", workflowPath),
			err,
		)
	}

	content := string(data)

	// Check for required workflow elements
	requiredElements := []string{"name:", "on:", "jobs:"}
	for _, element := range requiredElements {
		if !strings.Contains(content, element) {
			results.Warnings = append(results.Warnings,
				domain.NewTemplateError(
					domain.ErrTemplateExecutionFailed,
					"Missing workflow element",
					fmt.Sprintf("Workflow missing required element: %s", element),
				).WithContext(workflowPath))
		}
	}

	return nil
}