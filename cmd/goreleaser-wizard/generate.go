package main

import (
	"fmt"
	"os"

	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/spf13/cobra"
)

// runGenerate runs the generate command.
func runGenerate(cmd *cobra.Command, args []string) {
	defer recoverFromPanic("generate command")

	force, _ := cmd.Flags().GetBool("force")
	configOnly, _ := cmd.Flags().GetBool("config-only")

	fmt.Println(titleStyle.Render("⚙️ Generating GoReleaser Configuration"))
	fmt.Println()

	// Detect project information
	projectConfig := &domain.SafeProjectConfig{}

	err := detectProjectInfo(projectConfig)
	if err != nil {
		displayError(err)

		return
	}

	// Override with flags if provided
	if name, _ := cmd.Flags().GetString("project-name"); name != "" {
		projectConfig.ProjectName = name
	}

	if path, _ := cmd.Flags().GetString("main-path"); path != "" {
		projectConfig.MainPath = path
	}

	if binName, _ := cmd.Flags().GetString("binary-name"); binName != "" {
		projectConfig.BinaryName = binName
	}

	if projType, _ := cmd.Flags().GetString("project-type"); projType != "" {
		projectConfig.ProjectType = domain.ProjectType(projType)
	}

	// Create workflow based on options
	var workflowType WorkflowType
	if configOnly {
		workflowType = WorkflowTypeConfigOnly
	} else {
		workflowType = WorkflowTypeFullWizard
	}

	logger := log.New(os.Stderr)
	if !ExecuteWorkflow(workflowType, projectConfig, force, logger) {
		return
	}

	fmt.Println()

	if configOnly {
		fmt.Println(successStyle.Render("📄 GoReleaser configuration generated successfully!"))
	} else {
		fmt.Println(
			successStyle.Render("🎉 GoReleaser configuration and workflow generated successfully!"),
		)
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Review generated files")
	fmt.Println("  • Run 'goreleaser-wizard validate' to check configuration")

	if !configOnly {
		fmt.Println("  • Commit .github/workflows/release.yml to enable CI/CD")
	}

	fmt.Println("  • Test configuration with 'goreleaser check'")
}

func init() {
	configureGenerateCommand(generateCmd)
}
