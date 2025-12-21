package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// InteractivePrompter handles user interaction during init
type InteractivePrompter struct {
	scanner *bufio.Scanner
}

// NewInteractivePrompter creates a new interactive prompter
func NewInteractivePrompter() *InteractivePrompter {
	return &InteractivePrompter{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// confirmDetectedInfo prompts user to confirm detected project information
func (ip *InteractivePrompter) confirmDetectedInfo(config *domain.SafeProjectConfig) (*domain.SafeProjectConfig, error) {
	fmt.Println(titleStyle.Render("📋 Project Information Detected:"))
	fmt.Printf("  Project Name: %s\n", config.ProjectName)
	fmt.Printf("  Main Path: %s\n", config.MainPath)
	fmt.Printf("  Binary Name: %s\n", config.BinaryName)
	fmt.Printf("  Project Type: %s\n", config.ProjectType)
	fmt.Println()

	// Ask if user wants to modify
	confirm, err := ip.promptYesNo("Do you want to modify any of these values?", true)
	if err != nil {
		return nil, err
	}

	if !confirm {
		return config, nil
	}

	// Allow modifications
	return ip.promptProjectInfo(config)
}

// promptProjectInfo prompts for project information
func (ip *InteractivePrompter) promptProjectInfo(config *domain.SafeProjectConfig) (*domain.SafeProjectConfig, error) {
	// Project name
	name, err := ip.promptString("Project Name", config.ProjectName)
	if err != nil {
		return nil, err
	}
	config.ProjectName = name

	// Project description
	desc, err := ip.promptString("Project Description", config.ProjectDescription)
	if err != nil {
		return nil, err
	}
	config.ProjectDescription = desc

	// Binary name
	binName, err := ip.promptString("Binary Name", config.BinaryName)
	if err != nil {
		return nil, err
	}
	config.BinaryName = binName

	// Main path
	mainPath, err := ip.promptString("Main Package Path", config.MainPath)
	if err != nil {
		return nil, err
	}
	config.MainPath = mainPath

	// Project type
	projectType, err := ip.promptProjectType(config.ProjectType)
	if err != nil {
		return nil, err
	}
	config.ProjectType = projectType

	return config, nil
}

// promptProjectType prompts for project type selection
func (ip *InteractivePrompter) promptProjectType(current domain.ProjectType) (domain.ProjectType, error) {
	fmt.Println(infoStyle.Render("\n🎯 Select Project Type:"))

	types := []domain.ProjectType{
		domain.ProjectTypeCLI,
		domain.ProjectTypeWebAPI,
		domain.ProjectTypeLibrary,
	}

	for i, pt := range types {
		marker := " "
		if pt == current {
			marker = "•"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, pt)
	}

	selection, err := ip.promptInt("Select project type (number)", 1, len(types))
	if err != nil {
		return "", err
	}

	return types[selection-1], nil
}

// promptPlatforms prompts for target platforms
func (ip *InteractivePrompter) promptPlatforms(current []domain.Platform) ([]domain.Platform, error) {
	fmt.Println(infoStyle.Render("\n🖥️  Select Target Platforms:"))

	available := []domain.Platform{
		domain.PlatformLinux,
		domain.PlatformDarwin,
		domain.PlatformWindows,
	}

	for i, platform := range available {
		marker := " "
		if containsPlatform(current, platform) {
			marker = "✓"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, platform)
	}

	selections, err := ip.promptMultiInt("Select platforms (comma-separated numbers)", 1, len(available))
	if err != nil {
		return nil, err
	}

	var selected []domain.Platform
	for _, sel := range selections {
		selected = append(selected, available[sel-1])
	}

	return selected, nil
}

// promptArchitectures prompts for target architectures
func (ip *InteractivePrompter) promptArchitectures(current []domain.Architecture) ([]domain.Architecture, error) {
	fmt.Println(infoStyle.Render("\n🏗️  Select Target Architectures:"))

	available := []domain.Architecture{
		domain.ArchitectureAMD64,
		domain.ArchitectureARM64,
	}

	for i, arch := range available {
		marker := " "
		if containsArchitecture(current, arch) {
			marker = "✓"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, arch)
	}

	selections, err := ip.promptMultiInt("Select architectures (comma-separated numbers)", 1, len(available))
	if err != nil {
		return nil, err
	}

	var selected []domain.Architecture
	for _, sel := range selections {
		selected = append(selected, available[sel-1])
	}

	return selected, nil
}

// promptCGO prompts for CGO configuration
func (ip *InteractivePrompter) promptCGO(current domain.CGOStatus) (domain.CGOStatus, error) {
	fmt.Println(infoStyle.Render("\n⚙️  CGO Configuration:"))

	statuses := []domain.CGOStatus{
		domain.CGOStatusDisabled,
		domain.CGOStatusEnabled,
		domain.CGOStatusRequired,
	}

	for i, status := range statuses {
		marker := " "
		if status == current {
			marker = "•"
		}
		fmt.Printf("  %s %d. %s - %s\n", marker, i+1, status, ip.getCGODescription(status))
	}

	selection, err := ip.promptInt("Select CGO status (number)", 1, len(statuses))
	if err != nil {
		return "", err
	}

	return statuses[selection-1], nil
}

// promptDocker prompts for Docker configuration
func (ip *InteractivePrompter) promptDocker(current domain.DockerSupport) (domain.DockerSupport, error) {
	fmt.Println(infoStyle.Render("\n🐳 Docker Support:"))

	levels := []domain.DockerSupport{
		domain.DockerSupportNone,
		domain.DockerSupportBuild,
		domain.DockerSupportDeploy,
		domain.DockerSupportBoth,
	}

	for i, level := range levels {
		marker := " "
		if level == current {
			marker = "•"
		}
		fmt.Printf("  %s %d. %s - %s\n", marker, i+1, level, ip.getDockerDescription(level))
	}

	selection, err := ip.promptInt("Select Docker support level (number)", 1, len(levels))
	if err != nil {
		return "", err
	}

	return levels[selection-1], nil
}

// promptGitProvider prompts for Git provider
func (ip *InteractivePrompter) promptGitProvider(current domain.GitProvider) (domain.GitProvider, error) {
	fmt.Println(infoStyle.Render("\n🔗 Git Provider:"))

	providers := []domain.GitProvider{
		domain.GitProviderGitHub,
		domain.GitProviderGitLab,
		domain.GitProviderBitbucket,
	}

	for i, provider := range providers {
		marker := " "
		if provider == current {
			marker = "•"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, provider)
	}

	selection, err := ip.promptInt("Select Git provider (number)", 1, len(providers))
	if err != nil {
		return "", err
	}

	return providers[selection-1], nil
}

// promptAdvancedOptions prompts for advanced configuration options
func (ip *InteractivePrompter) promptAdvancedOptions(config *domain.SafeProjectConfig) error {
	fmt.Println(infoStyle.Render("\n🔧 Advanced Options:"))

	// LDFlags
	ldflags, err := ip.promptYesNo("Include LDFlags for version information?", config.LDFlags)
	if err != nil {
		return err
	}
	config.LDFlags = ldflags

	// Code signing
	signing, err := ip.promptYesNo("Enable code signing?", config.SigningLevel != domain.SigningLevelNone)
	if err != nil {
		return err
	}
	if signing {
		config.SigningLevel = domain.SigningLevelBasic
	} else {
		config.SigningLevel = domain.SigningLevelNone
	}

	// SBOM
	sbom, err := ip.promptYesNo("Generate SBOM?", config.SBOM)
	if err != nil {
		return err
	}
	config.SBOM = sbom

	// Homebrew
	homebrew, err := ip.promptYesNo("Generate Homebrew formula?", config.Homebrew)
	if err != nil {
		return err
	}
	config.Homebrew = homebrew

	// Snap
	snap, err := ip.promptYesNo("Generate Snap package?", config.Snap)
	if err != nil {
		return err
	}
	config.Snap = snap

	// GitHub Actions
	actions, err := ip.promptYesNo("Generate GitHub Actions workflow?", config.ActionLevel != domain.ActionLevelNone)
	if err != nil {
		return err
	}
	if actions {
		config.ActionLevel = domain.ActionLevelBasic
	} else {
		config.ActionLevel = domain.ActionLevelNone
	}

	return nil
}

// promptString prompts for a string value
func (ip *InteractivePrompter) promptString(label, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("\n%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("\n%s: ", label)
	}

	if !ip.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	value := strings.TrimSpace(ip.scanner.Text())
	if value == "" && defaultValue != "" {
		return defaultValue, nil
	}
	return value, nil
}

// promptYesNo prompts for a yes/no answer
func (ip *InteractivePrompter) promptYesNo(label string, defaultValue bool) (bool, error) {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("\n%s [%s]? ", label, defaultStr)

	if !ip.scanner.Scan() {
		return false, fmt.Errorf("failed to read input")
	}

	response := strings.ToLower(strings.TrimSpace(ip.scanner.Text()))
	if response == "" {
		return defaultValue, nil
	}

	return response == "y" || response == "yes", nil
}

// promptInt prompts for an integer within a range
func (ip *InteractivePrompter) promptInt(label string, min, max int) (int, error) {
	for {
		fmt.Printf("\n%s (%d-%d): ", label, min, max)

		if !ip.scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}

		value, err := strconv.Atoi(strings.TrimSpace(ip.scanner.Text()))
		if err != nil {
			fmt.Println(errorStyle.Render("Please enter a valid number"))
			continue
		}

		if value < min || value > max {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		return value, nil
	}
}

// promptMultiInt prompts for multiple integers within a range
func (ip *InteractivePrompter) promptMultiInt(label string, min, max int) ([]int, error) {
	for {
		fmt.Printf("\n%s (%d-%d): ", label, min, max)

		if !ip.scanner.Scan() {
			return nil, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(ip.scanner.Text())
		if input == "" {
			return []int{1}, nil // Default to first option
		}

		parts := strings.Split(input, ",")
		var selections []int

		for _, part := range parts {
			value, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				fmt.Println(errorStyle.Render("Please enter valid numbers separated by commas"))
				selections = nil
				break
			}

			if value < min || value > max {
				fmt.Printf("Please enter numbers between %d and %d\n", min, max)
				selections = nil
				break
			}

			selections = append(selections, value)
		}

		if selections != nil {
			return selections, nil
		}
	}
}

// getCGODescription returns description for CGO status
func (ip *InteractivePrompter) getCGODescription(status domain.CGOStatus) string {
	switch status {
	case domain.CGOStatusDisabled:
		return "Disable CGO compilation (recommended for Go projects)"
	case domain.CGOStatusEnabled:
		return "Enable CGO compilation when available"
	case domain.CGOStatusRequired:
		return "Require CGO compilation and fail if not available"
	default:
		return "Unknown"
	}
}

// getDockerDescription returns description for Docker support level
func (ip *InteractivePrompter) getDockerDescription(level domain.DockerSupport) string {
	switch level {
	case domain.DockerSupportNone:
		return "No Docker support"
	case domain.DockerSupportBuild:
		return "Build Docker images only"
	case domain.DockerSupportDeploy:
		return "Publish Docker images only"
	case domain.DockerSupportBoth:
		return "Build and publish Docker images"
	default:
		return "Unknown"
	}
}

// containsPlatform checks if platform is in the list
func containsPlatform(platforms []domain.Platform, target domain.Platform) bool {
	return slices.Contains(platforms, target)
}

// containsArchitecture checks if architecture is in the list
func containsArchitecture(architectures []domain.Architecture, target domain.Architecture) bool {
	return slices.Contains(architectures, target)
}
