package validation

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// Service provides validation functionality
type Service struct {
	envVars map[string]EnvironmentVariable
}

// NewService creates a new validation service
func NewService() *Service {
	return &Service{
		envVars: GetEnvironmentVariables(),
	}
}

// ValidateEnvironment validates all environment variables
func (s *Service) ValidateEnvironment() *EnvironmentValidationResult {
	result := NewEnvironmentValidationResult()

	criticalVars := GetCriticalVariables()
	optionalVars := GetOptionalVariables()

	// Validate critical variables
	for name, envVar := range criticalVars {
		value := os.Getenv(name)
		issue := ValidateEnvironmentVariable(name, value, envVar)
		
		if IsValidationIssueEmpty(issue) {
			// Variable is valid
			result.ValidatedVariables[name] = VariableStatus{
				Present: true,
				Valid:   true,
				Masked:  maskValue(value),
			}
		} else {
			// Variable has issues
			if issue.Severity == SeverityCritical {
				result.AddMissingCritical(name)
				result.Issues = append(result.Issues, issue)
			} else if issue.Severity == SeverityError {
				result.Issues = append(result.Issues, issue)
			} else {
				result.Warnings = append(result.Warnings, issue)
			}
			
			result.ValidatedVariables[name] = VariableStatus{
				Present: value != "",
				Valid:   issue.Severity != SeverityCritical && issue.Severity != SeverityError,
				Issue:   issue.UserMessage,
				Masked:  maskValue(value),
			}
		}
	}

	// Validate optional variables
	for name, envVar := range optionalVars {
		value := os.Getenv(name)
		issue := ValidateEnvironmentVariable(name, value, envVar)
		
		if IsValidationIssueEmpty(issue) {
			if value != "" {
				// Variable is valid and present
				result.ValidatedVariables[name] = VariableStatus{
					Present: true,
					Valid:   true,
					Masked:  maskValue(value),
				}
			}
		} else {
			// Variable has issues
			if issue.Severity == SeverityWarning && value == "" {
				result.AddMissingOptional(name)
			}
			
			if issue.Severity == SeverityCritical || issue.Severity == SeverityError {
				result.Issues = append(result.Issues, issue)
			} else {
				result.Warnings = append(result.Warnings, issue)
			}
			
			result.ValidatedVariables[name] = VariableStatus{
				Present: value != "",
				Valid:   issue.Severity != SeverityCritical && issue.Severity != SeverityError,
				Issue:   issue.UserMessage,
				Masked:  maskValue(value),
			}
		}
	}

	return result
}

// ExtractEnvironmentVariablesFromConfigs extracts environment variables from GoReleaser configs
func (s *Service) ExtractEnvironmentVariablesFromConfigs(configFiles ...string) *ConfigAnalysisResult {
	result := NewConfigAnalysisResult()
	
	if len(configFiles) == 0 {
		configFiles = []string{".goreleaser.yaml", ".goreleaser.pro.yaml"}
	}

	var allVars []string
	
	// Extract variables from each config file
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			continue // Skip non-existent files
		}
		
		result.ConfigFiles = append(result.ConfigFiles, configFile)
		
		file, err := os.Open(configFile)
		if err != nil {
			result.Issues = append(result.Issues, ValidationIssue{
				Field:       configFile,
				Message:     fmt.Sprintf("Cannot read config file: %v", err),
				UserMessage: fmt.Sprintf("Unable to read configuration file %s", configFile),
				Severity:    SeverityError,
				Code:        "CONFIG_READ_ERROR",
			})
			continue
		}
		defer file.Close()

		// Scan file for environment variable patterns
		scanner := bufio.NewScanner(file)
		envVarPattern := regexp.MustCompile(`\{\{\s*\.Env\.([A-Z_]+)\s*\}\}`)
		
		for scanner.Scan() {
			line := scanner.Text()
			matches := envVarPattern.FindAllStringSubmatch(line, -1)
			
			for _, match := range matches {
				if len(match) > 1 {
					varName := match[1]
					allVars = append(allVars, varName)
				}
			}
		}
		
		if err := scanner.Err(); err != nil {
			result.Issues = append(result.Issues, ValidationIssue{
				Field:       configFile,
				Message:     fmt.Sprintf("Error scanning config file: %v", err),
				UserMessage: fmt.Sprintf("Error reading configuration file %s", configFile),
				Severity:    SeverityError,
				Code:        "CONFIG_SCAN_ERROR",
			})
		}
	}
	
	// Remove duplicates and sort
	uniqueVars := lo.Uniq(allVars)
	sort.Strings(uniqueVars)
	result.ExtractedVariables = uniqueVars
	
	return result
}

// ValidateEnvExampleSync validates that .env.example is in sync with config usage
func (s *Service) ValidateEnvExampleSync(envExamplePath string, configFiles ...string) *ConfigAnalysisResult {
	result := s.ExtractEnvironmentVariablesFromConfigs(configFiles...)
	
	if envExamplePath == "" {
		envExamplePath = ".env.example"
	}
	
	// Read .env.example file
	envExampleVars, err := s.readEnvExampleVariables(envExamplePath)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:       envExamplePath,
			Message:     fmt.Sprintf("Cannot read .env.example: %v", err),
			UserMessage: "Unable to read .env.example file",
			Severity:    SeverityError,
			Code:        "ENV_EXAMPLE_READ_ERROR",
		})
		return result
	}
	
	// Find variables in configs but missing from .env.example
	configVarsSet := lo.SliceToMap(result.ExtractedVariables, func(v string) (string, bool) {
		return v, true
	})
	
	for _, configVar := range result.ExtractedVariables {
		if !lo.Contains(envExampleVars, configVar) {
			result.MissingInExample = append(result.MissingInExample, configVar)
		}
	}
	
	// Find variables in .env.example but not used in configs
	for _, exampleVar := range envExampleVars {
		if !configVarsSet[exampleVar] {
			result.UnusedInExample = append(result.UnusedInExample, exampleVar)
		}
	}
	
	// Add warnings for sync issues
	if len(result.MissingInExample) > 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:       envExamplePath,
			Message:     fmt.Sprintf("Variables used in configs but missing from .env.example: %s", strings.Join(result.MissingInExample, ", ")),
			UserMessage: "Some environment variables used in GoReleaser configs are not documented in .env.example",
			Severity:    SeverityError,
			Code:        "MISSING_IN_ENV_EXAMPLE",
		})
	}
	
	if len(result.UnusedInExample) > 0 {
		result.Warnings = append(result.Warnings, ValidationIssue{
			Field:       envExamplePath,
			Message:     fmt.Sprintf("Variables in .env.example but not used in configs: %s", strings.Join(result.UnusedInExample, ", ")),
			UserMessage: "Some variables in .env.example are not used in any GoReleaser configuration",
			Severity:    SeverityWarning,
			Code:        "UNUSED_IN_ENV_EXAMPLE",
		})
	}
	
	return result
}

// readEnvExampleVariables reads variable names from .env.example file
func (s *Service) readEnvExampleVariables(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var vars []string
	scanner := bufio.NewScanner(file)
	envVarPattern := regexp.MustCompile(`^([A-Z_]+)=`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		
		matches := envVarPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			vars = append(vars, matches[1])
		}
	}
	
	return vars, scanner.Err()
}

// GenerateValidationReport generates a comprehensive validation report
func (s *Service) GenerateValidationReport(configFiles ...string) *ValidationReport {
	report := &ValidationReport{
		Timestamp: time.Now(),
		RecommendedActions: []string{},
	}
	
	// Validate environment variables
	report.Environment = s.ValidateEnvironment()
	
	// Analyze configuration files
	report.ConfigAnalysis = s.ValidateEnvExampleSync("", configFiles...)
	
	// Determine overall status
	report.OverallStatus = s.determineOverallStatus(report.Environment, report.ConfigAnalysis)
	
	// Generate summary
	report.Summary = s.generateSummary(report.Environment, report.ConfigAnalysis)
	
	// Generate recommended actions
	report.RecommendedActions = s.generateRecommendedActions(report.Environment, report.ConfigAnalysis)
	
	return report
}

// determineOverallStatus determines the overall validation status
func (s *Service) determineOverallStatus(env *EnvironmentValidationResult, config *ConfigAnalysisResult) ValidationStatus {
	hasCritical := len(env.CriticalMissing) > 0 || 
		lo.SomeBy(env.Issues, func(issue ValidationIssue) bool { return issue.Severity == SeverityCritical }) ||
		lo.SomeBy(config.Issues, func(issue ValidationIssue) bool { return issue.Severity == SeverityCritical })
	
	if hasCritical {
		return StatusCriticalErrors
	}
	
	hasErrors := len(env.Issues) > 0 || len(config.Issues) > 0
	if hasErrors {
		return StatusErrors
	}
	
	hasWarnings := len(env.Warnings) > 0 || len(config.Warnings) > 0 || len(env.OptionalMissing) > 0
	if hasWarnings {
		return StatusWarnings
	}
	
	return StatusOK
}

// generateSummary generates a summary of validation results
func (s *Service) generateSummary(env *EnvironmentValidationResult, config *ConfigAnalysisResult) ValidationSummary {
	return ValidationSummary{
		TotalChecks:     len(s.envVars),
		CriticalIssues: len(lo.Filter(env.Issues, func(issue ValidationIssue, _ int) bool {
			return issue.Severity == SeverityCritical
		})) + len(lo.Filter(config.Issues, func(issue ValidationIssue, _ int) bool {
			return issue.Severity == SeverityCritical
		})),
		Errors:         len(env.Issues) + len(config.Issues),
		Warnings:       len(env.Warnings) + len(config.Warnings),
		MissingCritical: len(env.CriticalMissing),
		MissingOptional: len(env.OptionalMissing),
	}
}

// generateRecommendedActions generates recommended actions based on validation results
func (s *Service) generateRecommendedActions(env *EnvironmentValidationResult, config *ConfigAnalysisResult) []string {
	var actions []string
	
	if len(env.CriticalMissing) > 0 {
		actions = append(actions, "Set required environment variables: "+strings.Join(env.CriticalMissing, ", "))
	}
	
	if len(config.MissingInExample) > 0 {
		actions = append(actions, "Add missing variables to .env.example: "+strings.Join(config.MissingInExample, ", "))
	}
	
	if len(env.OptionalMissing) > 0 && len(env.OptionalMissing) <= 5 {
		actions = append(actions, "Consider setting optional variables for enhanced functionality: "+strings.Join(env.OptionalMissing[:5], ", "))
	} else if len(env.OptionalMissing) > 5 {
		actions = append(actions, fmt.Sprintf("Consider setting optional variables for enhanced functionality (%d available)", len(env.OptionalMissing)))
	}
	
	hasValidationErrors := lo.SomeBy(env.Issues, func(issue ValidationIssue) bool {
		return strings.Contains(issue.Code, "VALIDATION_FAILED")
	})
	
	if hasValidationErrors {
		actions = append(actions, "Fix environment variable format issues")
	}
	
	hasPlaceholderValues := lo.SomeBy(env.Warnings, func(issue ValidationIssue) bool {
		return issue.Code == "PLACEHOLDER_VALUE"
	})
	
	if hasPlaceholderValues {
		actions = append(actions, "Replace placeholder values with actual configuration")
	}
	
	if len(actions) == 0 {
		actions = append(actions, "Environment validation passed - ready for GoReleaser execution")
	}
	
	return actions
}

// GetVariableDocumentation returns documentation for a specific variable
func (s *Service) GetVariableDocumentation(name string) mo.Option[EnvironmentVariable] {
	if envVar, exists := s.envVars[name]; exists {
		return mo.Some(envVar)
	}
	return mo.None[EnvironmentVariable]()
}

// GetVariablesByCategory returns variables filtered by category
func (s *Service) GetVariablesByCategory(category VariableCategory) map[string]EnvironmentVariable {
	return GetVariablesByCategory(category)
}

// ListAllVariables returns all environment variable definitions
func (s *Service) ListAllVariables() map[string]EnvironmentVariable {
	return s.envVars
}