package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Validation repository (would be injected in real implementation)
var fileSystemRepo FileSystemRepository

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate GoReleaser configuration",
	Long: `Validate your GoReleaser configuration and check for common issues.

This command will:
- Check if .goreleaser.yaml exists and is valid YAML
- Run goreleaser check if available
- Verify project structure matches configuration
- Check for missing dependencies
- Suggest improvements`,
	Run: runValidate,
}

func init() {
	validateCmd.Flags().Bool("verbose", false, "show detailed validation output")
	validateCmd.Flags().Bool("fix", false, "attempt to fix common issues")
	validateCmd.Flags().Bool("project-only", false, "validate project structure only")
}

func runValidate(cmd *cobra.Command, args []string) {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("validate command")

	verbose, _ := cmd.Flags().GetBool("verbose")
	fix, _ := cmd.Flags().GetBool("fix")
	projectOnly, _ := cmd.Flags().GetBool("project-only")

	fmt.Println(titleStyle.Render("🔍 Validating GoReleaser Configuration"))
	fmt.Println()

	// Initialize dependencies (in real implementation, this would be injected)
	fileSystemRepo = &SimpleFileSystemRepository{}

	// Collect validation results
	results := &ValidationResults{}

	if !projectOnly {
		// Validate GoReleaser configuration
		if err := validateGoReleaserConfig(results); err != nil {
			displayError(err)
			return
		}

		// Validate GitHub Actions workflow
		if err := validateGitHubActions(results); err != nil {
			displayError(err)
			return
		}
	}

	// Validate project structure
	if err := validateProjectStructure(results); err != nil {
		displayError(err)
		return
	}

	// Display results
	displayValidationResults(results, verbose)

	// Attempt fixes if requested
	if fix && len(results.Errors) > 0 {
		if err := attemptFixes(results); err != nil {
			displayError(err)
			return
		}
	}

	// Exit with appropriate code
	os.Exit(results.GetExitCode())
}
