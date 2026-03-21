package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/charmbracelet/huh"
)

// NonInteractiveHelp is the help message shown when interactive mode is requested
// but the command is not running in a terminal.
const NonInteractiveHelp = `Interactive mode requires a terminal.

To run in non-interactive mode, add --interactive=false:
  goreleaser-wizard init --interactive=false

This will use detected project defaults without prompting.

Or run in a terminal to use the interactive TUI wizard.`

// IsTerminal returns true if stdout is connected to a terminal.
func IsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

// RunTUIWizard runs the interactive TUI wizard using huh.
func RunTUIWizard(config *domain.SafeProjectConfig) error {
	// Basic project information
	var projectName, projectDesc, binaryName, mainPath string
	var projectType domain.ProjectType

	projectName = config.ProjectName
	projectDesc = config.ProjectDescription
	binaryName = config.BinaryName
	mainPath = config.MainPath
	projectType = config.ProjectType

	// Build project type options
	projectTypes := []huh.Option[domain.ProjectType]{
		huh.NewOption("CLI Application", domain.ProjectTypeCLI),
		huh.NewOption("Web API", domain.ProjectTypeWebAPI),
		huh.NewOption("Library", domain.ProjectTypeLibrary),
		huh.NewOption("gRPC Service", domain.ProjectTypeGRPCService),
		huh.NewOption("Microservice", domain.ProjectTypeMicroservice),
		huh.NewOption("Desktop Application", domain.ProjectTypeDesktop),
		huh.NewOption("Daemon/Service", domain.ProjectTypeDaemon),
		huh.NewOption("Command Line Tool", domain.ProjectTypeTool),
	}

	// Platform and architecture selections
	var selectedPlatforms []string
	var selectedArchitectures []string

	platformOptions := []huh.Option[string]{
		huh.NewOption("Linux", string(domain.PlatformLinux)),
		huh.NewOption("macOS", string(domain.PlatformDarwin)),
		huh.NewOption("Windows", string(domain.PlatformWindows)),
		huh.NewOption("FreeBSD", string(domain.PlatformFreeBSD)),
	}

	archOptions := []huh.Option[string]{
		huh.NewOption("amd64 (x86_64)", string(domain.ArchitectureAMD64)),
		huh.NewOption("arm64 (ARM64)", string(domain.ArchitectureARM64)),
	}

	// Initialize with current values
	for _, p := range config.Platforms {
		selectedPlatforms = append(selectedPlatforms, string(p))
	}

	for _, a := range config.Architectures {
		selectedArchitectures = append(selectedArchitectures, string(a))
	}

	// CGO, Docker, Git Provider
	var cgoStatus domain.CGOStatus
	var dockerSupport domain.DockerSupport
	var gitProvider domain.GitProvider

	cgoStatus = config.CGOStatus
	dockerSupport = config.DockerSupport
	gitProvider = config.GitProvider

	cgoOptions := []huh.Option[domain.CGOStatus]{
		huh.NewOption("Disabled (recommended)", domain.CGOStatusDisabled),
		huh.NewOption("Enabled when available", domain.CGOStatusEnabled),
		huh.NewOption("Required", domain.CGOStatusRequired),
	}

	dockerOptions := []huh.Option[domain.DockerSupport]{
		huh.NewOption("None", domain.DockerSupportNone),
		huh.NewOption("Build images only", domain.DockerSupportBuild),
		huh.NewOption("Publish images only", domain.DockerSupportDeploy),
		huh.NewOption("Build and publish", domain.DockerSupportBoth),
	}

	gitProviderOptions := []huh.Option[domain.GitProvider]{
		huh.NewOption("GitHub", domain.GitProviderGitHub),
		huh.NewOption("GitLab", domain.GitProviderGitLab),
		huh.NewOption("Bitbucket", domain.GitProviderBitbucket),
	}

	// Advanced options (default to true for richer configuration)
	var includeLDFlags, enableSigning, generateSBOM, generateHomebrew, generateSnap, generateActions bool

	includeLDFlags = true
	enableSigning = true
	generateSBOM = true
	generateHomebrew = true
	generateSnap = true
	generateActions = true

	// Form theme with lipgloss
	theme := huh.ThemeCharm()

	// Build the form
	form := huh.NewForm(
		// Group 1: Basic Project Information
		huh.NewGroup(
			huh.NewNote().
				Title("GoReleaser Wizard").
				Description("Let's configure your GoReleaser setup!\nThis wizard will guide you through the configuration process."),

			huh.NewInput().
				Title("Project Name").
				Description("The name of your project").
				Placeholder("my-awesome-project").
				Value(&projectName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name is required")
					}

					return domain.ValidateProjectName(s)
				}),

			huh.NewInput().
				Title("Project Description").
				Description("A brief description of your project").
				Placeholder("A CLI tool that does amazing things").
				Value(&projectDesc).
				Validate(func(s string) error {
					if s == "" {
						return nil // Optional
					}

					return domain.ValidateProjectDescription(s)
				}),

			huh.NewInput().
				Title("Binary Name").
				Description("The name of the compiled binary").
				Placeholder("myapp").
				Value(&binaryName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("binary name is required")
					}

					return domain.ValidateBinaryName(s)
				}),

			huh.NewInput().
				Title("Main Package Path").
				Description("Path to your main.go file (e.g., ./cmd/myapp)").
				Placeholder("./cmd/myapp").
				Value(&mainPath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("main path is required")
					}

					return domain.ValidateMainPath(s)
				}),

			huh.NewSelect[domain.ProjectType]().
				Title("Project Type").
				Description("What type of project is this?").
				Options(projectTypes...).
				Value(&projectType),
		),

		// Group 2: Target Platforms
		huh.NewGroup(
			huh.NewNote().
				Title("Target Platforms").
				Description("Select the platforms you want to build for."),

			huh.NewMultiSelect[string]().
				Title("Platforms").
				Description("Select target operating systems").
				Options(platformOptions...).
				Value(&selectedPlatforms).
				Validate(func(s []string) error {
					if len(s) == 0 {
						return fmt.Errorf("at least one platform is required")
					}

					return nil
				}),

			huh.NewMultiSelect[string]().
				Title("Architectures").
				Description("Select target CPU architectures").
				Options(archOptions...).
				Value(&selectedArchitectures).
				Validate(func(s []string) error {
					if len(s) == 0 {
						return fmt.Errorf("at least one architecture is required")
					}

					return nil
				}),
		),

		// Group 3: Build Configuration
		huh.NewGroup(
			huh.NewNote().
				Title("Build Configuration").
				Description("Configure build options for your project."),

			huh.NewSelect[domain.CGOStatus]().
				Title("CGO Configuration").
				Description("Should CGO be enabled for builds?").
				Options(cgoOptions...).
				Value(&cgoStatus),

			huh.NewSelect[domain.DockerSupport]().
				Title("Docker Support").
				Description("Configure Docker image building and publishing").
				Options(dockerOptions...).
				Value(&dockerSupport),

			huh.NewSelect[domain.GitProvider]().
				Title("Git Provider").
				Description("Which Git hosting service do you use?").
				Options(gitProviderOptions...).
				Value(&gitProvider),
		),

		// Group 4: Advanced Options
		huh.NewGroup(
			huh.NewNote().
				Title("Advanced Options").
				Description("Additional features and integrations (optional)."),

			huh.NewConfirm().
				Title("Include LDFlags?").
				Description("Inject version information into the binary at build time").
				Value(&includeLDFlags),

			huh.NewConfirm().
				Title("Enable Code Signing?").
				Description("Sign binaries for distribution").
				Value(&enableSigning),

			huh.NewConfirm().
				Title("Generate SBOM?").
				Description("Software Bill of Materials for security compliance").
				Value(&generateSBOM),

			huh.NewConfirm().
				Title("Generate Homebrew Formula?").
				Description("Create a Homebrew formula for macOS users").
				Value(&generateHomebrew),

			huh.NewConfirm().
				Title("Generate Snap Package?").
				Description("Create a Snap package for Ubuntu/Debian").
				Value(&generateSnap),

			huh.NewConfirm().
				Title("Generate GitHub Actions?").
				Description("Create a CI/CD workflow for automated releases").
				Value(&generateActions),
		),

		// Group 5: Confirmation
		huh.NewGroup(
			huh.NewNote().
				Title("Configuration Summary").
				DescriptionFunc(func() string {
					return fmt.Sprintf(
						"Project: %s (%s)\nBinary: %s\nPlatforms: %v\nArchitectures: %v\nDocker: %s\nGitHub Actions: %v",
						projectName,
						projectType,
						binaryName,
						selectedPlatforms,
						selectedArchitectures,
						dockerSupport,
						generateActions,
					)
				}, nil),
		).WithHide(true),
	).WithTheme(theme)

	// Run the form
	err := form.Run()
	if err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	// Update config with form values
	config.ProjectName = projectName
	config.ProjectDescription = projectDesc
	config.BinaryName = binaryName
	config.MainPath = mainPath
	config.ProjectType = projectType

	// Convert platform selections
	config.Platforms = nil
	for _, p := range selectedPlatforms {
		if slices.Contains(
			[]string{
				string(domain.PlatformLinux),
				string(domain.PlatformDarwin),
				string(domain.PlatformWindows),
				string(domain.PlatformFreeBSD),
			},
			p,
		) {
			config.Platforms = append(config.Platforms, domain.Platform(p))
		}
	}

	// Convert architecture selections
	config.Architectures = nil
	for _, a := range selectedArchitectures {
		if slices.Contains(
			[]string{string(domain.ArchitectureAMD64), string(domain.ArchitectureARM64)},
			a,
		) {
			config.Architectures = append(config.Architectures, domain.Architecture(a))
		}
	}

	config.CGOStatus = cgoStatus
	config.DockerSupport = dockerSupport
	config.GitProvider = gitProvider
	config.LDFlags = includeLDFlags

	if enableSigning {
		config.SigningLevel = domain.SigningLevelBasic
	} else {
		config.SigningLevel = domain.SigningLevelNone
	}

	config.SBOM = generateSBOM
	config.Homebrew = generateHomebrew
	config.Snap = generateSnap

	if generateActions {
		config.ActionLevel = domain.ActionLevelBasic
	} else {
		config.ActionLevel = domain.ActionLevelNone
	}

	// Apply any needed defaults based on selections
	config.ApplyDefaults()

	return nil
}
