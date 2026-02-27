package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/charmbracelet/log"
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
	config := &domain.SafeProjectConfig{}
	if err := detectProjectInfo(config); err != nil {
		displayError(err)

		return
	}

	// Override with flags if provided
	if name, _ := cmd.Flags().GetString("project-name"); name != "" {
		config.ProjectName = name
	}

	if path, _ := cmd.Flags().GetString("main-path"); path != "" {
		config.MainPath = path
	}

	if binName, _ := cmd.Flags().GetString("binary-name"); binName != "" {
		config.BinaryName = binName
	}

	if projType, _ := cmd.Flags().GetString("project-type"); projType != "" {
		config.ProjectType = domain.ProjectType(projType)
	}

	// Create workflow based on options
	var workflowType WorkflowType
	if configOnly {
		workflowType = WorkflowTypeConfigOnly
	} else {
		workflowType = WorkflowTypeFullWizard
	}

	logger := log.New(os.Stderr)
	builder := NewWorkflowBuilder(logger)

	workflow, err := builder.BuildWorkflow(workflowType, config, force)
	if err != nil {
		displayError(err)

		return
	}

	// Execute workflow
	ctx := context.Background()
	if err := workflow.Execute(ctx); err != nil {
		displayError(err)

		return
	}

	// Display results
	displayJobResults(workflow.GetResults())

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
