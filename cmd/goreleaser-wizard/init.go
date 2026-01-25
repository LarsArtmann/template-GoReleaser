package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// runInitWizard runs the init command.
func runInitWizard(cmd *cobra.Command, args []string) {
	defer recoverFromPanic("init wizard")

	force, _ := cmd.Flags().GetBool("force")
	interactive, _ := cmd.Flags().GetBool("interactive")

	fmt.Println(titleStyle.Render("🧙 Initializing GoReleaser Configuration"))
	fmt.Println()

	// Detect project information
	config := &domain.SafeProjectConfig{}
	if err := detectProjectInfo(config); err != nil {
		displayError(err)
		return
	}

	// Apply defaults to ensure complete configuration
	config.ApplyDefaults()

	// If interactive mode, prompt for confirmation/changes
	if interactive {
		if err := runInteractiveWizard(config); err != nil {
			displayError(err)
			return
		}
	}

	// Create and execute workflow
	logger := log.New(os.Stderr)
	builder := NewWorkflowBuilder(logger)

	workflow, err := builder.BuildWorkflow(WorkflowTypeFullWizard, config, force)
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
	fmt.Println(successStyle.Render("🎉 GoReleaser configuration initialized successfully!"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Review the generated .goreleaser.yaml")
	fmt.Println("  • Run 'goreleaser-wizard validate' to check the configuration")
	fmt.Println("  • Commit the configuration to your repository")
}

// detectProjectInfo detects project information from the current directory.
func detectProjectInfo(config *domain.SafeProjectConfig) error {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to get working directory",
			"Could not determine current directory",
			err,
		).WithContext("detect_project")
	}

	// Check for go.mod file
	goModPath := filepath.Join(wd, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return domain.NewValidationError(
			domain.ErrFileNotFound,
			"Not a Go project",
			"go.mod file not found. Please run this command in a Go project directory.",
		).WithContext("detect_project")
	}

	// Read go.mod to get module name
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to read go.mod",
			"Could not read go.mod file",
			err,
		).WithContext(goModPath)
	}

	// Parse module name from go.mod
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	var moduleName string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				moduleName = parts[1]
			}
			break
		}
	}

	if moduleName == "" {
		return domain.NewValidationError(
			domain.ErrInvalidFileFormat,
			"Invalid go.mod format",
			"Could not find module declaration in go.mod",
		).WithContext(goModPath)
	}

	// Extract project name from module path
	parts := strings.Split(moduleName, "/")
	projectName := parts[len(parts)-1]

	// Detect main path
	mainPath, binaryName, projectType := detectMainStructure(wd)

	// Set configuration values
	config.ProjectName = projectName
	config.MainPath = mainPath
	config.BinaryName = binaryName
	config.ProjectType = domain.ProjectType(projectType)

	return nil
}

// detectMainStructure detects the main.go structure and determines project type.
func detectMainStructure(wd string) (mainPath, binaryName, projectType string) {
	// Check for main.go in root
	if _, err := os.Stat(filepath.Join(wd, "main.go")); err == nil {
		return ".", filepath.Base(wd), "cli"
	}

	// Check for cmd/*/main.go structure
	cmdDir := filepath.Join(wd, "cmd")
	if entries, err := os.ReadDir(cmdDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				mainFile := filepath.Join(cmdDir, entry.Name(), "main.go")
				if _, err := os.Stat(mainFile); err == nil {
					return "./cmd/" + entry.Name(), entry.Name(), "cli"
				}
			}
		}
	}

	// Check for other common structures
	patterns := []struct {
		path     string
		binName  string
		projType string
	}{
		{"src/*/main.go", "", "library"},
		{"pkg/*/main.go", "", "library"},
		{"*.go", filepath.Base(wd), "cli"},
	}

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(wd, pattern.path))
		if len(matches) > 0 {
			if pattern.binName == "" {
				// Extract binary name from path
				parts := strings.Split(matches[0], string(filepath.Separator))
				for i, part := range parts {
					if part == "src" || part == "pkg" && i+1 < len(parts) {
						binaryName = filepath.Base(matches[0])
						binaryName = strings.TrimSuffix(binaryName, ".go")
						return pattern.path, binaryName, pattern.projType
					}
				}
			}
			return pattern.path, pattern.binName, pattern.projType
		}
	}

	// Default fallback
	return ".", filepath.Base(wd), "cli"
}

// runInteractiveWizard runs an interactive wizard for configuration.
func runInteractiveWizard(config *domain.SafeProjectConfig) error {
	prompter := NewInteractivePrompter()

	// Confirm and modify detected information
	var err error
	config, err = prompter.confirmDetectedInfo(config)
	if err != nil {
		return err
	}

	// Prompt for platforms
	config.Platforms, err = prompter.promptPlatforms(config.Platforms)
	if err != nil {
		return err
	}

	// Prompt for architectures
	config.Architectures, err = prompter.promptArchitectures(config.Architectures)
	if err != nil {
		return err
	}

	// Prompt for CGO configuration
	config.CGOStatus, err = prompter.promptCGO(config.CGOStatus)
	if err != nil {
		return err
	}

	// Prompt for Docker support
	config.DockerSupport, err = prompter.promptDocker(config.DockerSupport)
	if err != nil {
		return err
	}

	// Prompt for Git provider
	config.GitProvider, err = prompter.promptGitProvider(config.GitProvider)
	if err != nil {
		return err
	}

	// Prompt for advanced options
	if err := prompter.promptAdvancedOptions(config); err != nil {
		return err
	}

	// Apply any needed defaults based on selections
	config.ApplyDefaults()

	// Show final configuration
	fmt.Println(titleStyle.Render("\n📋 Final Configuration:"))
	fmt.Printf("  Project: %s (%s)\n", config.ProjectName, config.ProjectType)
	fmt.Printf("  Binary: %s\n", config.BinaryName)
	fmt.Printf("  Main: %s\n", config.MainPath)
	fmt.Printf("  Platforms: %v\n", config.Platforms)
	fmt.Printf("  Architectures: %v\n", config.Architectures)
	fmt.Printf("  CGO: %s\n", config.CGOStatus)
	fmt.Printf("  Docker: %s\n", config.DockerSupport)
	fmt.Printf("  Git Provider: %s\n", config.GitProvider)
	fmt.Printf("  Code Signing: %s\n", config.SigningLevel)
	fmt.Printf("  Homebrew: %v\n", config.Homebrew)
	fmt.Printf("  GitHub Actions: %v\n", config.ActionLevel != domain.ActionLevelNone)
	fmt.Println()

	return nil
}

func init() {
	configureInitCommand(initCmd)
}
