package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	// Check required files
	requiredFiles := []string{
		"go.mod",
		"README.md",
		"LICENSE",
	}

	for _, file := range requiredFiles {
		result.Checks++
		if _, err := os.Stat(file); os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required file: %s", file))
			result.Success = false
		}
	}

	// Check required directories
	requiredDirs := []string{
		"cmd",
	}

	for _, dir := range requiredDirs {
		result.Checks++
		if info, err := os.Stat(dir); os.IsNotExist(err) || !info.IsDir() {
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required directory: %s", dir))
			result.Success = false
		}
	}

	// Check GoReleaser configuration files
	goreleaserFiles := []string{
		".goreleaser.yaml",
		".goreleaser.yml",
		".goreleaser.pro.yaml",
		".goreleaser.pro.yml",
		"goreleaser.yaml",
		"goreleaser.yml",
	}

	foundGoReleaserConfig := false
	for _, file := range goreleaserFiles {
		if _, err := os.Stat(file); err == nil {
			foundGoReleaserConfig = true
			result.Checks++

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

	result.Checks++
	if !foundGoReleaserConfig {
		result.Errors = append(result.Errors, "No GoReleaser configuration files found")
		result.Success = false
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

	// Check required tools
	requiredTools := map[string]string{
		"go":         "Go compiler",
		"git":        "Git version control",
		"goreleaser": "GoReleaser binary",
	}

	for tool, description := range requiredTools {
		result.Checks++
		if _, err := exec.LookPath(tool); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Required tool not found: %s (%s)", tool, description))
			result.Success = false
		}
	}

	// Check recommended tools
	recommendedTools := map[string]string{
		"docker": "Docker for container builds",
		"yq":     "YAML processor",
		"cosign": "Container signing",
		"syft":   "SBOM generation",
	}

	for tool, description := range recommendedTools {
		result.Checks++
		if _, err := exec.LookPath(tool); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Recommended tool not found: %s (%s)", tool, description))
		}
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

// isPlaceholderValue checks if a value appears to be a placeholder
func (s *ValidationServiceImpl) isPlaceholderValue(value string) bool {
	placeholders := []string{"your-", "xxxx", "example", "changeme", "todo", "test-"}
	lowerValue := strings.ToLower(value)
	for _, placeholder := range placeholders {
		if strings.HasPrefix(lowerValue, placeholder) {
			return true
		}
	}
	return len(value) < 3
}
