package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ValidationServiceImpl implements ValidationService interface
type ValidationServiceImpl struct {
	configService ConfigService
}

// NewValidationService creates a new validation service
func NewValidationService(configService ConfigService) ValidationService {
	return &ValidationServiceImpl{
		configService: configService,
	}
}

// ValidateProject validates the entire project structure
func (s *ValidationServiceImpl) ValidateProject() (*ValidationResult, error) {
	result := &ValidationResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
		Checks:   0,
	}

	// Check required files using lo.Map and lo.Filter for functional patterns
	requiredFiles := []string{
		"go.mod",
		"README.md",
		"LICENSE",
	}

	result.Checks += len(requiredFiles)
	missingFiles := lo.Filter(requiredFiles, func(file string, _ int) bool {
		_, err := os.Stat(file)
		return os.IsNotExist(err)
	})

	if len(missingFiles) > 0 {
		result.Errors = append(result.Errors, lo.Map(missingFiles, func(file string, _ int) string {
			return fmt.Sprintf("Missing required file: %s", file)
		})...)
		result.Success = false
	}

	// Check required directories using functional patterns
	requiredDirs := []string{
		"cmd",
	}

	result.Checks += len(requiredDirs)
	missingDirs := lo.Filter(requiredDirs, func(dir string, _ int) bool {
		info, err := os.Stat(dir)
		return os.IsNotExist(err) || !info.IsDir()
	})

	if len(missingDirs) > 0 {
		result.Errors = append(result.Errors, lo.Map(missingDirs, func(dir string, _ int) string {
			return fmt.Sprintf("Missing required directory: %s", dir)
		})...)
		result.Success = false
	}

	// Check GoReleaser configuration files using functional patterns
	goreleaserFiles := []string{
		".goreleaser.yaml",
		".goreleaser.yml",
		".goreleaser.pro.yaml",
		".goreleaser.pro.yml",
		"goreleaser.yaml",
		"goreleaser.yml",
	}

	existingConfigs := lo.Filter(goreleaserFiles, func(file string, _ int) bool {
		_, err := os.Stat(file)
		return err == nil
	})

	result.Checks += len(existingConfigs) + 1

	if len(existingConfigs) == 0 {
		result.Errors = append(result.Errors, "No GoReleaser configuration files found")
		result.Success = false
	} else {
		// Validate each existing config
		for _, file := range existingConfigs {
			// Validate YAML syntax
			if !s.validateYAMLSyntax(file) {
				result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML syntax in %s", file))
				result.Success = false
			}

			// Validate with GoReleaser if available
			if _, err := exec.LookPath("goreleaser"); err == nil {
				if !s.validateGoReleaserConfig(file) {
					result.Warnings = append(result.Warnings, fmt.Sprintf("GoReleaser validation failed for %s", file))
				}
			}
		}
	}

	return result, nil
}

// ValidateEnvironment validates environment variables
func (s *ValidationServiceImpl) ValidateEnvironment() (*ValidationResult, error) {
	result := &ValidationResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
		Checks:   0,
	}

	// Check critical environment variables
	criticalVars := map[string]string{
		"GITHUB_TOKEN": "GitHub API access token for releases",
	}

	for varName, description := range criticalVars {
		result.Checks++
		if value := os.Getenv(varName); value == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Missing critical environment variable: %s (%s)", varName, description))
			result.Success = false
		} else if s.isPlaceholderValue(value) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Environment variable %s appears to be a placeholder", varName))
		}
	}

	// Check optional environment variables
	optionalVars := map[string]string{
		"DOCKER_USERNAME":           "Docker Hub username",
		"DOCKER_PASSWORD":           "Docker Hub password/token",
		"GORELEASER_KEY":            "GoReleaser Pro license key",
		"HOMEBREW_TAP_GITHUB_TOKEN": "GitHub token for Homebrew tap",
		"SCOOP_GITHUB_TOKEN":        "GitHub token for Scoop bucket",
	}

	for varName := range optionalVars {
		result.Checks++
		if value := os.Getenv(varName); value != "" {
			if s.isPlaceholderValue(value) {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Environment variable %s appears to be a placeholder", varName))
			}
		}
	}

	// Check Go environment
	result.Checks++
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		result.Errors = append(result.Errors, "go.mod file not found")
		result.Success = false
	}

	result.Checks++
	if _, err := exec.LookPath("go"); err != nil {
		result.Errors = append(result.Errors, "Go compiler not found in PATH")
		result.Success = false
	}

	// Check Git environment
	result.Checks++
	if info, err := os.Stat(".git"); os.IsNotExist(err) || !info.IsDir() {
		result.Errors = append(result.Errors, "Not a git repository")
		result.Success = false
	}

	result.Checks++
	if _, err := exec.LookPath("git"); err != nil {
		result.Errors = append(result.Errors, "Git not found in PATH")
		result.Success = false
	}

	return result, nil
}

// ValidateGoReleaser validates GoReleaser configuration files
func (s *ValidationServiceImpl) ValidateGoReleaser(configPath string) (*ValidationResult, error) {
	result := &ValidationResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
		Checks:   0,
	}

	// Check if file exists
	result.Checks++
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("GoReleaser config file not found: %s", configPath))
		result.Success = false
		return result, nil
	}

	// Validate YAML syntax
	result.Checks++
	if !s.validateYAMLSyntax(configPath) {
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML syntax in %s", configPath))
		result.Success = false
	}

	// Validate with GoReleaser if available
	result.Checks++
	if _, err := exec.LookPath("goreleaser"); err != nil {
		result.Warnings = append(result.Warnings, "GoReleaser not installed, skipping native validation")
	} else {
		if !s.validateGoReleaserConfig(configPath) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("GoReleaser validation failed for %s", configPath))
		}
	}

	return result, nil
}

// ValidateTools validates required tools are installed
func (s *ValidationServiceImpl) ValidateTools() (*ValidationResult, error) {
	result := &ValidationResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
		Checks:   0,
	}

	// Check required tools using functional patterns
	requiredTools := map[string]string{
		"go":         "Go compiler",
		"git":        "Git version control",
		"goreleaser": "GoReleaser binary",
	}

	result.Checks += len(requiredTools)
	missingRequiredTools := lo.PickBy(requiredTools, func(tool, description string) bool {
		return s.getToolPath(tool).IsAbsent()
	})

	if len(missingRequiredTools) > 0 {
		result.Errors = append(result.Errors, lo.MapToSlice(missingRequiredTools, func(tool, description string) string {
			return fmt.Sprintf("Required tool not found: %s (%s)", tool, description)
		})...)
		result.Success = false
	}

	// Check recommended tools using functional patterns
	recommendedTools := map[string]string{
		"docker": "Docker for container builds",
		"yq":     "YAML processor",
		"cosign": "Container signing",
		"syft":   "SBOM generation",
	}

	result.Checks += len(recommendedTools)
	missingRecommendedTools := lo.PickBy(recommendedTools, func(tool, description string) bool {
		return s.getToolPath(tool).IsAbsent()
	})

	if len(missingRecommendedTools) > 0 {
		result.Warnings = append(result.Warnings, lo.MapToSlice(missingRecommendedTools, func(tool, description string) string {
			return fmt.Sprintf("Recommended tool not found: %s (%s)", tool, description)
		})...)
	}

	return result, nil
}

// validateYAMLSyntax performs basic YAML syntax validation
func (s *ValidationServiceImpl) validateYAMLSyntax(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		return false
	}

	// Basic YAML syntax validation - check for common issues
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Check for tabs (YAML doesn't allow them)
		if strings.Contains(line, "\t") {
			return false
		}
	}
	return true
}

// validateGoReleaserConfig runs GoReleaser check command
func (s *ValidationServiceImpl) validateGoReleaserConfig(file string) bool {
	cmd := exec.Command("goreleaser", "check", "--config", file)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		// Filter out multiple token warnings as they're expected in templates
		if !strings.Contains(outputStr, "multiple tokens") {
			return false
		}
	}
	return true
}

// isPlaceholderValue checks if a value appears to be a placeholder using functional patterns
func (s *ValidationServiceImpl) isPlaceholderValue(value string) bool {
	placeholders := []string{"your-", "xxxx", "example", "changeme", "todo", "test-"}
	lowerValue := strings.ToLower(value)
	
	return len(value) < 3 || lo.SomeBy(placeholders, func(placeholder string) bool {
		return strings.HasPrefix(lowerValue, placeholder)
	})
}

// getToolPath returns an Option containing the tool path if found
func (s *ValidationServiceImpl) getToolPath(tool string) mo.Option[string] {
	if path, err := exec.LookPath(tool); err == nil {
		return mo.Some(path)
	}
	return mo.None[string]()
}

// checkFileExists returns an Option containing file info if file exists
func (s *ValidationServiceImpl) checkFileExists(path string) mo.Option[os.FileInfo] {
	if info, err := os.Stat(path); err == nil {
		return mo.Some(info)
	}
	return mo.None[os.FileInfo]()
}
