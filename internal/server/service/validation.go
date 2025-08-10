package service

import (
	"os"
	"os/exec"
	"strings"

	"github.com/samber/do/v2"
)

type ValidationService struct {
	configService *ConfigService
}

type ValidationResult struct {
	Valid     bool     `json:"valid"`
	Issues    []string `json:"issues"`
	Warnings  []string `json:"warnings"`
	Timestamp string   `json:"timestamp"`
}

type EnvironmentResult struct {
	Valid       bool              `json:"valid"`
	MissingVars []string          `json:"missing_vars"`
	Warnings    []string          `json:"warnings"`
	Issues      []string          `json:"issues"`
	EnvVars     map[string]string `json:"env_vars,omitempty"`
	Timestamp   string            `json:"timestamp"`
}

func NewValidationService(injector do.Injector) (*ValidationService, error) {
	configService, err := do.InvokeAs[*ConfigService](injector)
	if err != nil {
		return nil, err
	}
	
	return &ValidationService{
		configService: configService,
	}, nil
}

func (s *ValidationService) ValidateConfig(configPath string, content string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Issues:   []string{},
		Warnings: []string{},
	}
	
	// If content is provided, validate syntax
	if content != "" {
		if err := s.configService.ValidateSyntax(content); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, "YAML syntax error: "+err.Error())
			return result
		}
	}
	
	// If configPath is provided, validate file
	if configPath != "" {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			result.Valid = false
			result.Issues = append(result.Issues, "Config file does not exist: "+configPath)
			return result
		}
		
		// Try to load and validate structure
		if _, err := s.configService.LoadConfig(configPath); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, "Config validation failed: "+err.Error())
			return result
		}
	}
	
	// Check for GoReleaser binary and validate with it
	if _, err := exec.LookPath("goreleaser"); err != nil {
		result.Warnings = append(result.Warnings, "GoReleaser binary not found - skipping native validation")
	} else if configPath != "" {
		if err := s.validateWithGoReleaser(configPath); err != nil {
			// Don't fail validation for GoReleaser issues, just warn
			result.Warnings = append(result.Warnings, "GoReleaser validation: "+err.Error())
		}
	}
	
	return result
}

func (s *ValidationService) ValidateEnvironment() *EnvironmentResult {
	result := &EnvironmentResult{
		Valid:       true,
		MissingVars: []string{},
		Warnings:    []string{},
		Issues:      []string{},
		EnvVars:     make(map[string]string),
	}
	
	// Critical environment variables
	criticalVars := []string{
		"GITHUB_TOKEN",
	}
	
	// Optional but recommended variables
	optionalVars := []string{
		"DOCKER_USERNAME",
		"DOCKER_PASSWORD",
		"GORELEASER_KEY",
		"HOMEBREW_TAP_GITHUB_TOKEN",
		"SCOOP_GITHUB_TOKEN",
	}
	
	// Check critical variables
	for _, varName := range criticalVars {
		value := os.Getenv(varName)
		if value == "" {
			result.Valid = false
			result.MissingVars = append(result.MissingVars, varName)
		} else {
			result.EnvVars[varName] = "***" // Hide actual values
		}
	}
	
	// Check optional variables
	for _, varName := range optionalVars {
		value := os.Getenv(varName)
		if value == "" {
			result.Warnings = append(result.Warnings, varName+" is not set")
		} else {
			result.EnvVars[varName] = "***" // Hide actual values
		}
	}
	
	// Check for required tools
	requiredTools := []string{"go", "git"}
	var issues []string
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			result.Valid = false
			issues = append(issues, tool+" not found in PATH")
		}
	}
	
	// Initialize the Issues field if it doesn't exist
	if result.Issues == nil {
		result.Issues = []string{}
	}
	result.Issues = append(result.Issues, issues...)
	
	// Check for optional tools
	optionalTools := []string{"docker", "goreleaser", "cosign", "syft"}
	for _, tool := range optionalTools {
		if _, err := exec.LookPath(tool); err != nil {
			result.Warnings = append(result.Warnings, tool+" not found in PATH")
		}
	}
	
	return result
}

func (s *ValidationService) GetValidationStatus() map[string]interface{} {
	configPath := s.configService.findConfigFile()
	
	configValid := configPath != ""
	if configValid {
		result := s.ValidateConfig(configPath, "")
		configValid = result.Valid
	}
	
	envResult := s.ValidateEnvironment()
	
	return map[string]interface{}{
		"config_valid": configValid,
		"env_valid":    envResult.Valid,
		"config_path":  configPath,
		"issues":       s.collectAllIssues(configPath),
	}
}

func (s *ValidationService) validateWithGoReleaser(configPath string) error {
	cmd := exec.Command("goreleaser", "check", "--config", configPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		// Filter out some common warnings that aren't critical
		if !strings.Contains(outputStr, "multiple tokens") {
			return err
		}
	}
	return nil
}

func (s *ValidationService) collectAllIssues(configPath string) []string {
	var allIssues []string
	
	if configPath != "" {
		result := s.ValidateConfig(configPath, "")
		allIssues = append(allIssues, result.Issues...)
	} else {
		allIssues = append(allIssues, "No GoReleaser config file found")
	}
	
	envResult := s.ValidateEnvironment()
	for _, missing := range envResult.MissingVars {
		allIssues = append(allIssues, "Missing environment variable: "+missing)
	}
	
	return allIssues
}