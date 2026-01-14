package main

import (
	"fmt"
	"os"
	"strings"
)

// validateWorkflowContent validates GitHub Actions workflow content.
func validateWorkflowContent(workflowPath string, results *ValidationResults) error {
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return NewSystemError(
			"FILE_READ_FAILED",
			"Failed to read workflow file",
			"Cannot read "+workflowPath,
			err,
		)
	}

	content := string(data)

	// Check for required workflow elements
	requiredElements := []string{"name:", "on:", "jobs:"}
	for _, element := range requiredElements {
		if !strings.Contains(content, element) {
			results.Warnings = append(results.Warnings,
				NewTemplateError(
					"TEMPLATE_EXECUTION_FAILED",
					"Missing workflow element",
					"Workflow missing required element: "+element,
				).WithContext(workflowPath))
		}
	}

	return nil
}
