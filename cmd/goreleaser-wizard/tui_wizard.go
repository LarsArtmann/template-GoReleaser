package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// Sentinel errors for TUI validation.
var (
	ErrProjectNameRequired  = errors.New("project name is required")
	ErrBinaryNameRequired   = errors.New("binary name is required")
	ErrMainPathRequired     = errors.New("main path is required")
	ErrPlatformRequired     = errors.New("at least one platform is required")
	ErrArchitectureRequired = errors.New("at least one architecture is required")
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

// tuiFormData holds the mutable form data for the TUI wizard.
type tuiFormData struct {
	projectName        string
	projectDesc        string
	binaryName         string
	mainPath           string
	projectType        domain.ProjectType
	selectedPlatforms  []string
	selectedArchitects []string
	cgoStatus          domain.CGOStatus
	dockerSupport      domain.DockerSupport
	gitProvider        domain.GitProvider
	includeLDFlags     bool
	enableSigning      bool
	generateSBOM       bool
	generateHomebrew   bool
	generateSnap       bool
	generateActions    bool
}

// RunTUIWizard runs the interactive TUI wizard using huh.
func RunTUIWizard(config *domain.SafeProjectConfig) error {
	formData := initializeFormData(config)
	platformOpts := getPlatformOptions()
	archOpts := getArchitectureOptions()

	form := buildTUIForm(&formData, platformOpts, archOpts)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	updateConfigFromForm(config, &formData)
	config.ApplyDefaults()

	return nil
}

func initializeFormData(config *domain.SafeProjectConfig) tuiFormData {
	selectedPlatforms := make([]string, 0, len(config.Platforms))
	for _, p := range config.Platforms {
		selectedPlatforms = append(selectedPlatforms, string(p))
	}

	selectedArchitects := make([]string, 0, len(config.Architectures))
	for _, a := range config.Architectures {
		selectedArchitects = append(selectedArchitects, string(a))
	}

	return tuiFormData{
		projectName:        config.ProjectName,
		projectDesc:        config.ProjectDescription,
		binaryName:         config.BinaryName,
		mainPath:           config.MainPath,
		projectType:        config.ProjectType,
		selectedPlatforms:  selectedPlatforms,
		selectedArchitects: selectedArchitects,
		cgoStatus:          config.CGOStatus,
		dockerSupport:      config.DockerSupport,
		gitProvider:        config.GitProvider,
		includeLDFlags:     true,
		enableSigning:      true,
		generateSBOM:       true,
		generateHomebrew:   true,
		generateSnap:       true,
		generateActions:    true,
	}
}

func getPlatformOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("Linux", string(domain.PlatformLinux)),
		huh.NewOption("macOS", string(domain.PlatformDarwin)),
		huh.NewOption("Windows", string(domain.PlatformWindows)),
		huh.NewOption("FreeBSD", string(domain.PlatformFreeBSD)),
	}
}

func getArchitectureOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("amd64 (x86_64)", string(domain.ArchitectureAMD64)),
		huh.NewOption("arm64 (ARM64)", string(domain.ArchitectureARM64)),
	}
}

func buildTUIForm(formData *tuiFormData, platformOpts, archOpts []huh.Option[string]) *huh.Form {
	theme := huh.ThemeFunc(huh.ThemeCharm)

	return huh.NewForm(
		buildBasicInfoGroup(formData),
		buildPlatformsGroup(formData, platformOpts, archOpts),
		buildBuildConfigGroup(formData),
		buildAdvancedOptionsGroup(formData),
		buildConfirmationGroup(formData),
	).WithTheme(theme)
}

func buildBasicInfoGroup(formData *tuiFormData) *huh.Group {
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

	return huh.NewGroup(
		huh.NewNote().
			Title("GoReleaser Wizard").
			Description("Let's configure your GoReleaser setup!\nThis wizard will guide you through the configuration process."),

		huh.NewInput().
			Title("Project Name").
			Description("The name of your project").
			Placeholder("my-awesome-project").
			Value(&formData.projectName).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf(
						"project name validation failed (input: %q): %w",
						formData.projectName,
						ErrProjectNameRequired,
					)
				}

				return domain.ValidateProjectName(s)
			}),

		huh.NewInput().
			Title("Project Description").
			Description("A brief description of your project").
			Placeholder("A CLI tool that does amazing things").
			Value(&formData.projectDesc).
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
			Value(&formData.binaryName).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf(
						"binary name validation failed (input: %q): %w",
						formData.binaryName,
						ErrBinaryNameRequired,
					)
				}

				return domain.ValidateBinaryName(s)
			}),

		huh.NewInput().
			Title("Main Package Path").
			Description("Path to your main.go file (e.g., ./cmd/myapp)").
			Placeholder("./cmd/myapp").
			Value(&formData.mainPath).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf(
						"main path validation failed (input: %q): %w",
						formData.mainPath,
						ErrMainPathRequired,
					)
				}

				return domain.ValidateMainPath(s)
			}),

		huh.NewSelect[domain.ProjectType]().
			Title("Project Type").
			Description("What type of project is this?").
			Options(projectTypes...).
			Value(&formData.projectType),
	)
}

func buildPlatformsGroup(
	formData *tuiFormData,
	platformOpts, archOpts []huh.Option[string],
) *huh.Group {
	return huh.NewGroup(
		huh.NewNote().
			Title("Target Platforms").
			Description("Select the platforms you want to build for."),

		huh.NewMultiSelect[string]().
			Title("Platforms").
			Description("Select target operating systems").
			Options(platformOpts...).
			Value(&formData.selectedPlatforms).
			Validate(func(s []string) error {
				if len(s) == 0 {
					return fmt.Errorf(
						"platforms validation failed (selected: %v): %w",
						formData.selectedPlatforms,
						ErrPlatformRequired,
					)
				}

				return nil
			}),

		huh.NewMultiSelect[string]().
			Title("Architectures").
			Description("Select target CPU architectures").
			Options(archOpts...).
			Value(&formData.selectedArchitects).
			Validate(func(s []string) error {
				if len(s) == 0 {
					return fmt.Errorf(
						"architectures validation failed: %w",
						ErrArchitectureRequired,
					)
				}

				return nil
			}),
	)
}

func buildBuildConfigGroup(formData *tuiFormData) *huh.Group {
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

	return huh.NewGroup(
		huh.NewNote().
			Title("Build Configuration").
			Description("Configure build options for your project."),

		huh.NewSelect[domain.CGOStatus]().
			Title("CGO Configuration").
			Description("Should CGO be enabled for builds?").
			Options(cgoOptions...).
			Value(&formData.cgoStatus),

		huh.NewSelect[domain.DockerSupport]().
			Title("Docker Support").
			Description("Configure Docker image building and publishing").
			Options(dockerOptions...).
			Value(&formData.dockerSupport),

		huh.NewSelect[domain.GitProvider]().
			Title("Git Provider").
			Description("Which Git hosting service do you use?").
			Options(gitProviderOptions...).
			Value(&formData.gitProvider),
	)
}

func buildAdvancedOptionsGroup(formData *tuiFormData) *huh.Group {
	return huh.NewGroup(
		huh.NewNote().
			Title("Advanced Options").
			Description("Additional features and integrations (optional)."),

		huh.NewConfirm().
			Title("Include LDFlags?").
			Description("Inject version information into the binary at build time").
			Value(&formData.includeLDFlags),

		huh.NewConfirm().
			Title("Enable Code Signing?").
			Description("Sign binaries for distribution").
			Value(&formData.enableSigning),

		huh.NewConfirm().
			Title("Generate SBOM?").
			Description("Software Bill of Materials for security compliance").
			Value(&formData.generateSBOM),

		huh.NewConfirm().
			Title("Generate Homebrew Formula?").
			Description("Create a Homebrew formula for macOS users").
			Value(&formData.generateHomebrew),

		huh.NewConfirm().
			Title("Generate Snap Package?").
			Description("Create a Snap package for Ubuntu/Debian").
			Value(&formData.generateSnap),

		huh.NewConfirm().
			Title("Generate GitHub Actions?").
			Description("Create a CI/CD workflow for automated releases").
			Value(&formData.generateActions),
	)
}

func buildConfirmationGroup(formData *tuiFormData) *huh.Group {
	return huh.NewGroup(
		huh.NewNote().
			Title("Configuration Summary").
			DescriptionFunc(func() string {
				return fmt.Sprintf(
					"Project: %s (%s)\nBinary: %s\nPlatforms: %v\nArchitectures: %v\nDocker: %s\nGitHub Actions: %v",
					formData.projectName,
					formData.projectType,
					formData.binaryName,
					formData.selectedPlatforms,
					formData.selectedArchitects,
					formData.dockerSupport,
					formData.generateActions,
				)
			}, nil),
	).WithHide(true)
}

func updateConfigFromForm(config *domain.SafeProjectConfig, formData *tuiFormData) {
	config.ProjectName = formData.projectName
	config.ProjectDescription = formData.projectDesc
	config.BinaryName = formData.binaryName
	config.MainPath = formData.mainPath
	config.ProjectType = formData.projectType

	convertPlatformSelections(config, formData.selectedPlatforms)
	convertArchitectureSelections(config, formData.selectedArchitects)

	config.CGOStatus = formData.cgoStatus
	config.DockerSupport = formData.dockerSupport
	config.GitProvider = formData.gitProvider
	config.LDFlags = formData.includeLDFlags

	if formData.enableSigning {
		config.SigningLevel = domain.SigningLevelBasic
	} else {
		config.SigningLevel = domain.SigningLevelNone
	}

	config.SBOM = formData.generateSBOM
	config.Homebrew = formData.generateHomebrew
	config.Snap = formData.generateSnap

	if formData.generateActions {
		config.ActionLevel = domain.ActionLevelBasic
	} else {
		config.ActionLevel = domain.ActionLevelNone
	}
}

func convertPlatformSelections(config *domain.SafeProjectConfig, selected []string) {
	config.Platforms = nil

	validPlatforms := []string{
		string(domain.PlatformLinux),
		string(domain.PlatformDarwin),
		string(domain.PlatformWindows),
		string(domain.PlatformFreeBSD),
	}

	for _, p := range selected {
		if slices.Contains(validPlatforms, p) {
			config.Platforms = append(config.Platforms, domain.Platform(p))
		}
	}
}

func convertArchitectureSelections(config *domain.SafeProjectConfig, selected []string) {
	config.Architectures = nil

	validArchs := []string{string(domain.ArchitectureAMD64), string(domain.ArchitectureARM64)}

	for _, a := range selected {
		if slices.Contains(validArchs, a) {
			config.Architectures = append(config.Architectures, domain.Architecture(a))
		}
	}
}
