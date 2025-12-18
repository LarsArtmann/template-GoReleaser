package types

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// ValidationResult represents structured validation results
type ValidationResult struct {
	IsValid  bool                 `json:"is_valid"`
	Errors   []*ValidationError   `json:"errors"`
	Warnings []*ValidationWarning `json:"warnings"`
	Summary  ValidationSummary    `json:"summary"`
}

// AddError adds an error to the validation result
func (vr *ValidationResult) AddError(err *ValidationError) {
	vr.Errors = append(vr.Errors, err)
	vr.IsValid = false
}

// AddWarning adds a warning to the validation result
func (vr *ValidationResult) AddWarning(warn *ValidationWarning) {
	vr.Warnings = append(vr.Warnings, warn)
}

// HasErrors returns true if there are errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// GetErrorCount returns the number of errors
func (vr *ValidationResult) GetErrorCount() int {
	return len(vr.Errors)
}

// GetWarningCount returns the number of warnings
func (vr *ValidationResult) GetWarningCount() int {
	return len(vr.Warnings)
}

// GetCriticalErrors returns critical errors
func (vr *ValidationResult) GetCriticalErrors() []*ValidationError {
	var critical []*ValidationError
	for _, err := range vr.Errors {
		if err.Level == ErrorLevelCritical {
			critical = append(critical, err)
		}
	}
	return critical
}

// GetHighErrors returns high-severity errors
func (vr *ValidationResult) GetHighErrors() []*ValidationError {
	var high []*ValidationError
	for _, err := range vr.Errors {
		if err.Level == ErrorLevelHigh {
			high = append(high, err)
		}
	}
	return high
}

// GetMediumErrors returns medium-severity errors
func (vr *ValidationResult) GetMediumErrors() []*ValidationError {
	var medium []*ValidationError
	for _, err := range vr.Errors {
		if err.Level == ErrorLevelMedium {
			medium = append(medium, err)
		}
	}
	return medium
}

// GetLowErrors returns low-severity errors
func (vr *ValidationResult) GetLowErrors() []*ValidationError {
	var low []*ValidationError
	for _, err := range vr.Errors {
		if err.Level == ErrorLevelLow {
			low = append(low, err)
		}
	}
	return low
}

// GetHighWarnings returns high-severity warnings
func (vr *ValidationResult) GetHighWarnings() []*ValidationWarning {
	var high []*ValidationWarning
	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelHigh {
			high = append(high, warn)
		}
	}
	return high
}

// GetMediumWarnings returns medium-severity warnings
func (vr *ValidationResult) GetMediumWarnings() []*ValidationWarning {
	var medium []*ValidationWarning
	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelMedium {
			medium = append(medium, warn)
		}
	}
	return medium
}

// GetLowWarnings returns low-severity warnings
func (vr *ValidationResult) GetLowWarnings() []*ValidationWarning {
	var low []*ValidationWarning
	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelLow {
			low = append(low, warn)
		}
	}
	return low
}

// GetErrorsByField returns errors grouped by field
func (vr *ValidationResult) GetErrorsByField() map[string][]*ValidationError {
	fieldErrors := make(map[string][]*ValidationError)
	for _, err := range vr.Errors {
		field := err.Field
		if field == "" {
			field = "general"
		}
		fieldErrors[field] = append(fieldErrors[field], err)
	}
	return fieldErrors
}

// GetWarningsByField returns warnings grouped by field
func (vr *ValidationResult) GetWarningsByField() map[string][]*ValidationWarning {
	fieldWarnings := make(map[string][]*ValidationWarning)
	for _, warn := range vr.Warnings {
		field := warn.Field
		if field == "" {
			field = "general"
		}
		fieldWarnings[field] = append(fieldWarnings[field], warn)
	}
	return fieldWarnings
}

// GetErrorsByCode returns errors grouped by error code
func (vr *ValidationResult) GetErrorsByCode() map[errors.ErrorCode][]*ValidationError {
	codeErrors := make(map[errors.ErrorCode][]*ValidationError)
	for _, err := range vr.Errors {
		codeErrors[err.Code] = append(codeErrors[err.Code], err)
	}
	return codeErrors
}

// GetWarningsByCode returns warnings grouped by warning code
func (vr *ValidationResult) GetWarningsByCode() map[string][]*ValidationWarning {
	codeWarnings := make(map[string][]*ValidationWarning)
	for _, warn := range vr.Warnings {
		codeWarnings[warn.Code] = append(codeWarnings[warn.Code], warn)
	}
	return codeWarnings
}

// UpdateSummary updates the validation summary
func (vr *ValidationResult) UpdateSummary() {
	vr.Summary = ValidationSummary{
		TotalErrors:    len(vr.Errors),
		TotalWarnings:  len(vr.Warnings),
		Critical:       len(vr.GetCriticalErrors()),
		High:           len(vr.GetHighErrors()),
		Medium:         len(vr.GetMediumErrors()),
		Low:            len(vr.GetLowErrors()),
		HighWarnings:   len(vr.GetHighWarnings()),
		MediumWarnings: len(vr.GetMediumWarnings()),
		LowWarnings:    len(vr.GetLowWarnings()),
	}
}

// String returns a human-readable string representation
func (vr *ValidationResult) String() string {
	var result strings.Builder

	if vr.IsValid {
		result.WriteString("✅ Validation passed")
	} else {
		result.WriteString("❌ Validation failed")
	}

	if len(vr.Errors) > 0 {
		result.WriteString(fmt.Sprintf("\n\n🚨 Errors (%d):", len(vr.Errors)))
		for _, err := range vr.Errors {
			result.WriteString(fmt.Sprintf("\n  • [%s] %s: %s", err.Level, err.Field, err.Message))
			if err.Suggestion != "" {
				result.WriteString(fmt.Sprintf("\n    💡 %s", err.Suggestion))
			}
		}
	}

	if len(vr.Warnings) > 0 {
		result.WriteString(fmt.Sprintf("\n\n⚠️  Warnings (%d):", len(vr.Warnings)))
		for _, warn := range vr.Warnings {
			result.WriteString(fmt.Sprintf("\n  • [%s] %s: %s", warn.Level, warn.Field, warn.Message))
			if warn.Suggestion != "" {
				result.WriteString(fmt.Sprintf("\n    💡 %s", warn.Suggestion))
			}
		}
	}

	return result.String()
}

// ToJSON converts validation result to JSON
func (vr *ValidationResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(vr, "", "  ")
}

// Merge merges another validation result into this one
func (vr *ValidationResult) Merge(other *ValidationResult) {
	vr.Errors = append(vr.Errors, other.Errors...)
	vr.Warnings = append(vr.Warnings, other.Warnings...)

	if other.HasErrors() {
		vr.IsValid = false
	}

	vr.UpdateSummary()
}

// Clone creates a deep copy of the validation result
func (vr *ValidationResult) Clone() *ValidationResult {
	clone := &ValidationResult{
		IsValid:  vr.IsValid,
		Errors:   make([]*ValidationError, len(vr.Errors)),
		Warnings: make([]*ValidationWarning, len(vr.Warnings)),
		Summary:  vr.Summary,
	}

	// Deep copy errors
	for i, err := range vr.Errors {
		clone.Errors[i] = &ValidationError{
			Code:       err.Code,
			Field:      err.Field,
			Message:    err.Message,
			Details:    err.Details,
			Context:    err.Context,
			Level:      err.Level,
			Suggestion: err.Suggestion,
		}
	}

	// Deep copy warnings
	for i, warn := range vr.Warnings {
		clone.Warnings[i] = &ValidationWarning{
			Code:       warn.Code,
			Field:      warn.Field,
			Message:    warn.Message,
			Details:    warn.Details,
			Context:    warn.Context,
			Level:      warn.Level,
			Suggestion: warn.Suggestion,
		}
	}

	return clone
}

// ValidationError represents a structured validation error
type ValidationError struct {
	Code       errors.ErrorCode `json:"code"`
	Field      string           `json:"field"`
	Message    string           `json:"message"`
	Details    string           `json:"details,omitempty"`
	Context    string           `json:"context,omitempty"`
	Level      ErrorLevel       `json:"level"`
	Suggestion string           `json:"suggestion,omitempty"`
}

// Error implements the error interface
func (ve *ValidationError) Error() string {
	if ve.Context != "" {
		return fmt.Sprintf("[%s] %s: %s (%s)", ve.Code, ve.Field, ve.Message, ve.Context)
	}
	return fmt.Sprintf("[%s] %s: %s", ve.Code, ve.Field, ve.Message)
}

// WithContext adds context to the validation error
func (ve *ValidationError) WithContext(context string) *ValidationError {
	ve.Context = context
	return ve
}

// WithDetails adds details to the validation error
func (ve *ValidationError) WithDetails(details string) *ValidationError {
	ve.Details = details
	return ve
}

// WithSuggestion adds suggestion to the validation error
func (ve *ValidationError) WithSuggestion(suggestion string) *ValidationError {
	ve.Suggestion = suggestion
	return ve
}

// ValidationWarning represents a structured validation warning
type ValidationWarning struct {
	Code       string       `json:"code"`
	Field      string       `json:"field"`
	Message    string       `json:"message"`
	Details    string       `json:"details,omitempty"`
	Context    string       `json:"context,omitempty"`
	Level      WarningLevel `json:"level"`
	Suggestion string       `json:"suggestion,omitempty"`
}

// String returns a string representation of the warning
func (vw *ValidationWarning) String() string {
	if vw.Context != "" {
		return fmt.Sprintf("[%s] %s: %s (%s)", vw.Code, vw.Field, vw.Message, vw.Context)
	}
	return fmt.Sprintf("[%s] %s: %s", vw.Code, vw.Field, vw.Message)
}

// WithContext adds context to the validation warning
func (vw *ValidationWarning) WithContext(context string) *ValidationWarning {
	vw.Context = context
	return vw
}

// WithDetails adds details to the validation warning
func (vw *ValidationWarning) WithDetails(details string) *ValidationWarning {
	vw.Details = details
	return vw
}

// WithSuggestion adds suggestion to the validation warning
func (vw *ValidationWarning) WithSuggestion(suggestion string) *ValidationWarning {
	vw.Suggestion = suggestion
	return vw
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalErrors    int `json:"total_errors"`
	TotalWarnings  int `json:"total_warnings"`
	Critical       int `json:"critical"`
	High           int `json:"high"`
	Medium         int `json:"medium"`
	Low            int `json:"low"`
	HighWarnings   int `json:"high_warnings"`
	MediumWarnings int `json:"medium_warnings"`
	LowWarnings    int `json:"low_warnings"`
}

// GetScore returns a validation score (0-100)
func (vs *ValidationSummary) GetScore() int {
	if vs.TotalErrors > 0 {
		return 0
	}

	if vs.TotalWarnings == 0 {
		return 100
	}

	// Deduct points based on warning severity
	score := 100
	score -= vs.HighWarnings * 10
	score -= vs.MediumWarnings * 5
	score -= vs.LowWarnings * 2

	if score < 0 {
		score = 0
	}

	return score
}

// GetGrade returns a validation grade
func (vs *ValidationSummary) GetGrade() string {
	score := vs.GetScore()

	if score == 100 {
		return "A+"
	} else if score >= 95 {
		return "A"
	} else if score >= 90 {
		return "A-"
	} else if score >= 85 {
		return "B+"
	} else if score >= 80 {
		return "B"
	} else if score >= 75 {
		return "B-"
	} else if score >= 70 {
		return "C+"
	} else if score >= 65 {
		return "C"
	} else if score >= 60 {
		return "C-"
	} else if score >= 50 {
		return "D"
	} else {
		return "F"
	}
}

// GetStatus returns a validation status
func (vs *ValidationSummary) GetStatus() ValidationStatus {
	if vs.TotalErrors > 0 {
		return ValidationStatusFailed
	}

	if vs.Critical > 0 || vs.High > 0 {
		return ValidationStatusCritical
	}

	if vs.Medium > 0 {
		return ValidationStatusWarning
	}

	if vs.Low > 0 {
		return ValidationStatusNotice
	}

	return ValidationStatusPassed
}

// ValidationStatus represents validation status
type ValidationStatus string

const (
	ValidationStatusPassed   ValidationStatus = "passed"
	ValidationStatusNotice   ValidationStatus = "notice"
	ValidationStatusWarning  ValidationStatus = "warning"
	ValidationStatusCritical ValidationStatus = "critical"
	ValidationStatusFailed   ValidationStatus = "failed"
)

// String returns string representation
func (vs ValidationStatus) String() string {
	return string(vs)
}

// GetIcon returns an icon for the status
func (vs ValidationStatus) GetIcon() string {
	switch vs {
	case ValidationStatusPassed:
		return "✅"
	case ValidationStatusNotice:
		return "ℹ️"
	case ValidationStatusWarning:
		return "⚠️"
	case ValidationStatusCritical:
		return "🚨"
	case ValidationStatusFailed:
		return "❌"
	default:
		return "❓"
	}
}

// GetColor returns a color for the status
func (vs ValidationStatus) GetColor() string {
	switch vs {
	case ValidationStatusPassed:
		return "green"
	case ValidationStatusNotice:
		return "blue"
	case ValidationStatusWarning:
		return "yellow"
	case ValidationStatusCritical:
		return "red"
	case ValidationStatusFailed:
		return "red"
	default:
		return "gray"
	}
}

// ErrorLevel represents error severity levels
type ErrorLevel string

const (
	ErrorLevelCritical ErrorLevel = "critical"
	ErrorLevelHigh     ErrorLevel = "high"
	ErrorLevelMedium   ErrorLevel = "medium"
	ErrorLevelLow      ErrorLevel = "low"
)

// String returns string representation
func (el ErrorLevel) String() string {
	return string(el)
}

// GetPriority returns numeric priority for sorting
func (el ErrorLevel) GetPriority() int {
	switch el {
	case ErrorLevelCritical:
		return 4
	case ErrorLevelHigh:
		return 3
	case ErrorLevelMedium:
		return 2
	case ErrorLevelLow:
		return 1
	default:
		return 0
	}
}

// WarningLevel represents warning severity levels
type WarningLevel string

const (
	WarningLevelHigh   WarningLevel = "high"
	WarningLevelMedium WarningLevel = "medium"
	WarningLevelLow    WarningLevel = "low"
)

// String returns string representation
func (wl WarningLevel) String() string {
	return string(wl)
}

// GetPriority returns numeric priority for sorting
func (wl WarningLevel) GetPriority() int {
	switch wl {
	case WarningLevelHigh:
		return 3
	case WarningLevelMedium:
		return 2
	case WarningLevelLow:
		return 1
	default:
		return 0
	}
}

// ValidationFilter represents filters for validation results
type ValidationFilter struct {
	Fields        []string           `json:"fields,omitempty"`
	Levels        []ErrorLevel       `json:"error_levels,omitempty"`
	Codes         []errors.ErrorCode `json:"error_codes,omitempty"`
	WarningLevels []WarningLevel     `json:"warning_levels,omitempty"`
	WarningCodes  []string           `json:"warning_codes,omitempty"`
	Context       string             `json:"context,omitempty"`
	MinScore      int                `json:"min_score,omitempty"`
}

// Filter applies the filter to a validation result
func (f *ValidationFilter) Filter(result *ValidationResult) *ValidationResult {
	filtered := &ValidationResult{
		IsValid:  result.IsValid,
		Errors:   []*ValidationError{},
		Warnings: []*ValidationWarning{},
	}

	// Filter errors
	for _, err := range result.Errors {
		if f.matchesError(err) {
			filtered.Errors = append(filtered.Errors, err)
		}
	}

	// Filter warnings
	for _, warn := range result.Warnings {
		if f.matchesWarning(warn) {
			filtered.Warnings = append(filtered.Warnings, warn)
		}
	}

	filtered.UpdateSummary()
	return filtered
}

// matchesError checks if an error matches the filter criteria
func (f *ValidationFilter) matchesError(err *ValidationError) bool {
	// Check fields
	if len(f.Fields) > 0 && !contains(f.Fields, err.Field) {
		return false
	}

	// Check levels
	if len(f.Levels) > 0 && !containsErrorLevels(f.Levels, err.Level) {
		return false
	}

	// Check codes
	if len(f.Codes) > 0 && !containsErrorCodes(f.Codes, err.Code) {
		return false
	}

	// Check context
	if f.Context != "" && err.Context != f.Context {
		return false
	}

	return true
}

// matchesWarning checks if a warning matches the filter criteria
func (f *ValidationFilter) matchesWarning(warn *ValidationWarning) bool {
	// Check fields
	if len(f.Fields) > 0 && !contains(f.Fields, warn.Field) {
		return false
	}

	// Check warning levels
	if len(f.WarningLevels) > 0 && !containsWarningLevels(f.WarningLevels, warn.Level) {
		return false
	}

	// Check warning codes
	if len(f.WarningCodes) > 0 && !contains(f.WarningCodes, warn.Code) {
		return false
	}

	// Check context
	if f.Context != "" && warn.Context != f.Context {
		return false
	}

	return true
}

// Helper functions
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func containsErrorLevels(slice []ErrorLevel, item ErrorLevel) bool {
	return slices.Contains(slice, item)
}

func containsErrorCodes(slice []errors.ErrorCode, item errors.ErrorCode) bool {
	return slices.Contains(slice, item)
}

func containsWarningLevels(slice []WarningLevel, item WarningLevel) bool {
	return slices.Contains(slice, item)
}
