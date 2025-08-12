package validation

import (
	"encoding/json"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// UserFriendlyError wraps an error with a user-friendly message
type UserFriendlyError struct {
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	ErrorCode   string `json:"error_code,omitempty"`
	Cause       error  `json:"-"`
}

func (e *UserFriendlyError) Error() string {
	return e.Message
}

func NewUserFriendlyError(message, userMessage string, cause error) *UserFriendlyError {
	return &UserFriendlyError{
		Message:     message,
		UserMessage: userMessage,
		Cause:       cause,
	}
}

// ValidationSeverity defines the severity level of validation issues
type ValidationSeverity string

const (
	SeverityCritical ValidationSeverity = "critical"
	SeverityError    ValidationSeverity = "error"
	SeverityWarning  ValidationSeverity = "warning"
	SeverityInfo     ValidationSeverity = "info"
)

// ValidationIssue represents a single validation issue
type ValidationIssue struct {
	Field       string             `json:"field"`
	Message     string             `json:"message"`
	UserMessage string             `json:"user_message"`
	Severity    ValidationSeverity `json:"severity"`
	Code        string             `json:"code,omitempty"`
	Value       string             `json:"value,omitempty"` // Masked value for debugging
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid     bool              `json:"valid"`
	Issues    []ValidationIssue `json:"issues"`
	Warnings  []ValidationIssue `json:"warnings"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:     true,
		Issues:    []ValidationIssue{},
		Warnings:  []ValidationIssue{},
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// AddIssue adds a validation issue
func (vr *ValidationResult) AddIssue(issue ValidationIssue) {
	if issue.Severity == SeverityCritical || issue.Severity == SeverityError {
		vr.Valid = false
		vr.Issues = append(vr.Issues, issue)
	} else {
		vr.Warnings = append(vr.Warnings, issue)
	}
}

// AddError adds an error-level validation issue
func (vr *ValidationResult) AddError(field, message, userMessage string) {
	vr.AddIssue(ValidationIssue{
		Field:       field,
		Message:     message,
		UserMessage: userMessage,
		Severity:    SeverityError,
	})
}

// AddWarning adds a warning-level validation issue
func (vr *ValidationResult) AddWarning(field, message, userMessage string) {
	vr.AddIssue(ValidationIssue{
		Field:       field,
		Message:     message,
		UserMessage: userMessage,
		Severity:    SeverityWarning,
	})
}

// AddCritical adds a critical-level validation issue
func (vr *ValidationResult) AddCritical(field, message, userMessage string) {
	vr.AddIssue(ValidationIssue{
		Field:       field,
		Message:     message,
		UserMessage: userMessage,
		Severity:    SeverityCritical,
	})
}

// HasCriticalIssues returns true if there are critical validation issues
func (vr *ValidationResult) HasCriticalIssues() bool {
	return lo.SomeBy(vr.Issues, func(issue ValidationIssue) bool {
		return issue.Severity == SeverityCritical
	})
}

// GetIssuesBySeverity returns issues filtered by severity level
func (vr *ValidationResult) GetIssuesBySeverity(severity ValidationSeverity) []ValidationIssue {
	allIssues := append(vr.Issues, vr.Warnings...)
	return lo.Filter(allIssues, func(issue ValidationIssue, _ int) bool {
		return issue.Severity == severity
	})
}

// EnvironmentVariable represents a single environment variable definition
type EnvironmentVariable struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required"`
	Category    VariableCategory  `json:"category"`
	Validator   mo.Option[string] `json:"validator,omitempty"` // Name of validation function
	Example     mo.Option[string] `json:"example,omitempty"`
	Format      mo.Option[string] `json:"format,omitempty"`
}

// VariableCategory defines categories for environment variables
type VariableCategory string

const (
	CategoryGitHub       VariableCategory = "github"
	CategoryDocker       VariableCategory = "docker"
	CategoryCloud        VariableCategory = "cloud"
	CategorySigning      VariableCategory = "signing"
	CategoryNotification VariableCategory = "notification"
	CategoryGeneral      VariableCategory = "general"
	CategoryArtifacts    VariableCategory = "artifacts"
)

// EnvironmentValidationResult represents the result of environment validation
type EnvironmentValidationResult struct {
	Valid              bool                        `json:"valid"`
	CriticalMissing    []string                    `json:"critical_missing"`
	OptionalMissing    []string                    `json:"optional_missing"`
	Issues             []ValidationIssue           `json:"issues"`
	Warnings           []ValidationIssue           `json:"warnings"`
	ValidatedVariables map[string]VariableStatus   `json:"validated_variables"`
	Timestamp          time.Time                   `json:"timestamp"`
	ValidationStatus   EnvironmentValidationStatus `json:"validation_status"`
}

// VariableStatus represents the status of a single variable validation
type VariableStatus struct {
	Present bool   `json:"present"`
	Valid   bool   `json:"valid"`
	Masked  string `json:"masked_value,omitempty"`
	Issue   string `json:"issue,omitempty"`
}

// EnvironmentValidationStatus defines overall validation status
type EnvironmentValidationStatus string

const (
	StatusReady      EnvironmentValidationStatus = "ready"
	StatusNeedsSetup EnvironmentValidationStatus = "needs_setup"
	StatusHasIssues  EnvironmentValidationStatus = "has_issues"
)

// NewEnvironmentValidationResult creates a new environment validation result
func NewEnvironmentValidationResult() *EnvironmentValidationResult {
	return &EnvironmentValidationResult{
		Valid:              true,
		CriticalMissing:    []string{},
		OptionalMissing:    []string{},
		Issues:             []ValidationIssue{},
		Warnings:           []ValidationIssue{},
		ValidatedVariables: make(map[string]VariableStatus),
		Timestamp:          time.Now(),
		ValidationStatus:   StatusReady,
	}
}

// AddMissingCritical adds a missing critical variable
func (evr *EnvironmentValidationResult) AddMissingCritical(name string) {
	evr.Valid = false
	evr.CriticalMissing = append(evr.CriticalMissing, name)
	evr.ValidationStatus = StatusNeedsSetup
}

// AddMissingOptional adds a missing optional variable
func (evr *EnvironmentValidationResult) AddMissingOptional(name string) {
	evr.OptionalMissing = append(evr.OptionalMissing, name)
	if evr.ValidationStatus == StatusReady {
		evr.ValidationStatus = StatusHasIssues
	}
}

// ConfigAnalysisResult represents the result of analyzing GoReleaser configs
type ConfigAnalysisResult struct {
	ConfigFiles        []string          `json:"config_files"`
	ExtractedVariables []string          `json:"extracted_variables"`
	MissingInExample   []string          `json:"missing_in_example"`
	UnusedInExample    []string          `json:"unused_in_example"`
	Issues             []ValidationIssue `json:"issues"`
	Warnings           []ValidationIssue `json:"warnings"`
	Timestamp          time.Time         `json:"timestamp"`
}

// NewConfigAnalysisResult creates a new config analysis result
func NewConfigAnalysisResult() *ConfigAnalysisResult {
	return &ConfigAnalysisResult{
		ConfigFiles:        []string{},
		ExtractedVariables: []string{},
		MissingInExample:   []string{},
		UnusedInExample:    []string{},
		Issues:             []ValidationIssue{},
		Warnings:           []ValidationIssue{},
		Timestamp:          time.Now(),
	}
}

// ValidationReport represents a comprehensive validation report
type ValidationReport struct {
	Environment        *EnvironmentValidationResult `json:"environment"`
	ConfigAnalysis     *ConfigAnalysisResult        `json:"config_analysis"`
	OverallStatus      ValidationStatus             `json:"overall_status"`
	Summary            ValidationSummary            `json:"summary"`
	Timestamp          time.Time                    `json:"timestamp"`
	RecommendedActions []string                     `json:"recommended_actions"`
}

// ValidationStatus defines overall validation status
type ValidationStatus string

const (
	StatusOK             ValidationStatus = "ok"
	StatusWarnings       ValidationStatus = "warnings"
	StatusErrors         ValidationStatus = "errors"
	StatusCriticalErrors ValidationStatus = "critical_errors"
)

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalChecks     int `json:"total_checks"`
	CriticalIssues  int `json:"critical_issues"`
	Errors          int `json:"errors"`
	Warnings        int `json:"warnings"`
	MissingCritical int `json:"missing_critical"`
	MissingOptional int `json:"missing_optional"`
}

// ToJSON converts validation result to JSON
func (vr *ValidationResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(vr, "", "  ")
}

// ToJSON converts environment validation result to JSON
func (evr *EnvironmentValidationResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(evr, "", "  ")
}

// ToJSON converts validation report to JSON
func (report *ValidationReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}
