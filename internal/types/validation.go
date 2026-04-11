package types

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// ValidationResult represents structured validation results.
type ValidationResult struct {
	IsValid  bool                 `json:"is_valid"`
	Errors   []*ValidationError   `json:"errors"`
	Warnings []*ValidationWarning `json:"warnings"`
	Summary  ValidationSummary    `json:"summary"`
}

// AddError adds an error to the validation result.
func (vr *ValidationResult) AddError(err *ValidationError) {
	vr.Errors = append(vr.Errors, err)
	vr.IsValid = false
}

// AddWarning adds a warning to the validation result.
func (vr *ValidationResult) AddWarning(warn *ValidationWarning) {
	vr.Warnings = append(vr.Warnings, warn)
}

// HasErrors returns true if there are errors.
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are warnings.
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// GetErrorCount returns the number of errors.
func (vr *ValidationResult) GetErrorCount() int {
	return len(vr.Errors)
}

// GetWarningCount returns the number of warnings.
func (vr *ValidationResult) GetWarningCount() int {
	return len(vr.Warnings)
}

// GetCriticalErrors returns critical errors.
func (vr *ValidationResult) GetCriticalErrors() []*ValidationError {
	var critical []*ValidationError

	for _, err := range vr.Errors {
		if err.Level == ErrorLevelCritical {
			critical = append(critical, err)
		}
	}

	return critical
}

// GetHighErrors returns high-severity errors.
func (vr *ValidationResult) GetHighErrors() []*ValidationError {
	var high []*ValidationError

	for _, err := range vr.Errors {
		if err.Level == ErrorLevelHigh {
			high = append(high, err)
		}
	}

	return high
}

// GetMediumErrors returns medium-severity errors.
func (vr *ValidationResult) GetMediumErrors() []*ValidationError {
	var medium []*ValidationError

	for _, err := range vr.Errors {
		if err.Level == ErrorLevelMedium {
			medium = append(medium, err)
		}
	}

	return medium
}

// GetLowErrors returns low-severity errors.
func (vr *ValidationResult) GetLowErrors() []*ValidationError {
	var low []*ValidationError

	for _, err := range vr.Errors {
		if err.Level == ErrorLevelLow {
			low = append(low, err)
		}
	}

	return low
}

// GetHighWarnings returns high-severity warnings.
func (vr *ValidationResult) GetHighWarnings() []*ValidationWarning {
	var high []*ValidationWarning

	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelHigh {
			high = append(high, warn)
		}
	}

	return high
}

// GetMediumWarnings returns medium-severity warnings.
func (vr *ValidationResult) GetMediumWarnings() []*ValidationWarning {
	var medium []*ValidationWarning

	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelMedium {
			medium = append(medium, warn)
		}
	}

	return medium
}

// GetLowWarnings returns low-severity warnings.
func (vr *ValidationResult) GetLowWarnings() []*ValidationWarning {
	var low []*ValidationWarning

	for _, warn := range vr.Warnings {
		if warn.Level == WarningLevelLow {
			low = append(low, warn)
		}
	}

	return low
}

// groupByField groups validation items by field name.
func groupByField[T validationItem](items []T) map[string][]T {
	fieldGroups := make(map[string][]T)

	for _, item := range items {
		field := item.GetField()
		if field == "" {
			field = "general"
		}

		fieldGroups[field] = append(fieldGroups[field], item)
	}

	return fieldGroups
}

// GetErrorsByField returns errors grouped by field.
func (vr *ValidationResult) GetErrorsByField() map[string][]*ValidationError {
	return groupByField(vr.Errors)
}

// GetWarningsByField returns warnings grouped by field.
func (vr *ValidationResult) GetWarningsByField() map[string][]*ValidationWarning {
	return groupByField(vr.Warnings)
}

// GetErrorsByCode returns errors grouped by error code.
func (vr *ValidationResult) GetErrorsByCode() map[errors.ErrorCode][]*ValidationError {
	codeErrors := make(map[errors.ErrorCode][]*ValidationError)
	for _, err := range vr.Errors {
		codeErrors[err.Code] = append(codeErrors[err.Code], err)
	}

	return codeErrors
}

// GetWarningsByCode returns warnings grouped by warning code.
func (vr *ValidationResult) GetWarningsByCode() map[string][]*ValidationWarning {
	codeWarnings := make(map[string][]*ValidationWarning)
	for _, warn := range vr.Warnings {
		codeWarnings[warn.Code] = append(codeWarnings[warn.Code], warn)
	}

	return codeWarnings
}

// UpdateSummary updates the validation summary.
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

// LevelGetter is an interface for types that have a Level field.
type LevelGetter interface {
	GetLevel() string
}

// validationItem represents the common fields of ValidationError and ValidationWarning.
type validationItem interface {
	LevelGetter
	GetField() string
	GetMessage() string
	GetSuggestion() string
}

// toValidationItems converts a slice of validation items to validationItem interface.
func toValidationItems[T validationItem](items []T) []validationItem {
	result := make([]validationItem, len(items))
	for i, item := range items {
		result[i] = item
	}

	return result
}

// GetLevel returns the level for ValidationError.
func (ve *ValidationError) GetLevel() string {
	return string(ve.Level)
}

// GetField returns the field for ValidationError.
func (ve *ValidationError) GetField() string {
	return ve.Field
}

// GetMessage returns the message for ValidationError.
func (ve *ValidationError) GetMessage() string {
	return ve.Message
}

// GetSuggestion returns the suggestion for ValidationError.
func (ve *ValidationError) GetSuggestion() string {
	return ve.Suggestion
}

// GetLevel returns the level for ValidationWarning.
func (vw *ValidationWarning) GetLevel() string {
	return string(vw.Level)
}

// GetField returns the field for ValidationWarning.
func (vw *ValidationWarning) GetField() string {
	return vw.Field
}

// GetMessage returns the message for ValidationWarning.
func (vw *ValidationWarning) GetMessage() string {
	return vw.Message
}

// GetSuggestion returns the suggestion for ValidationWarning.
func (vw *ValidationWarning) GetSuggestion() string {
	return vw.Suggestion
}

// appendValidationItems appends validation items (errors or warnings) to the result builder.
func appendValidationItems[T validationItem](
	result *strings.Builder,
	emoji string,
	label string,
	items []T,
) {
	if len(items) == 0 {
		return
	}

	fmt.Fprintf(result, "\n\n%s %s (%d):", emoji, label, len(items))

	for _, item := range items {
		fmt.Fprintf(
			result,
			"\n  • [%s] %s: %s",
			item.GetLevel(),
			item.GetField(),
			item.GetMessage(),
		)

		if item.GetSuggestion() != "" {
			result.WriteString("\n    💡 " + item.GetSuggestion())
		}
	}
}

// String returns a human-readable string representation.
func (vr *ValidationResult) String() string {
	var result strings.Builder

	if vr.IsValid {
		result.WriteString("✅ Validation passed")
	} else {
		result.WriteString("❌ Validation failed")
	}

	appendValidationItems(&result, "🚨", "Errors", toValidationItems(vr.Errors))
	appendValidationItems(&result, "⚠️ ", "Warnings", toValidationItems(vr.Warnings))

	return result.String()
}

// ToJSON converts validation result to JSON.
func (vr *ValidationResult) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(vr, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validation result to JSON: %w", err)
	}

	return data, nil
}

// Merge merges another validation result into this one.
func (vr *ValidationResult) Merge(other *ValidationResult) {
	vr.Errors = append(vr.Errors, other.Errors...)
	vr.Warnings = append(vr.Warnings, other.Warnings...)

	if other.HasErrors() {
		vr.IsValid = false
	}

	vr.UpdateSummary()
}

// Clone creates a deep copy of the validation result.
func (vr *ValidationResult) Clone() *ValidationResult {
	clone := &ValidationResult{
		IsValid:  vr.IsValid,
		Errors:   make([]*ValidationError, len(vr.Errors)),
		Warnings: make([]*ValidationWarning, len(vr.Warnings)),
		Summary:  vr.Summary,
	}

	// Deep copy errors
	for i, err := range vr.Errors {
		clone.Errors[i] = err.Clone()
	}

	// Deep copy warnings
	for i, warn := range vr.Warnings {
		clone.Warnings[i] = warn.Clone()
	}

	return clone
}

// ValidationError represents a structured validation error.
type ValidationError struct {
	Code       errors.ErrorCode `json:"code"`
	Field      string           `json:"field"`
	Message    string           `json:"message"`
	Details    string           `json:"details,omitempty"`
	Context    string           `json:"context,omitempty"`
	Level      ErrorLevel       `json:"level"`
	Suggestion string           `json:"suggestion,omitempty"`
}

// formatValidationItem formats a validation item with optional context.
func formatValidationItem(code, field, message, context string) string {
	if context != "" {
		return fmt.Sprintf("[%s] %s: %s (%s)", code, field, message, context)
	}

	return fmt.Sprintf("[%s] %s: %s", code, field, message)
}

// Error implements the error interface.
func (ve *ValidationError) Error() string {
	return formatValidationItem(string(ve.Code), ve.Field, ve.Message, ve.Context)
}

// WithContext adds context to the validation error.
func (ve *ValidationError) WithContext(context string) *ValidationError {
	ve.Context = context

	return ve
}

// WithDetails adds details to the validation error.
func (ve *ValidationError) WithDetails(details string) *ValidationError {
	ve.Details = details

	return ve
}

// WithSuggestion adds suggestion to the validation error.
func (ve *ValidationError) WithSuggestion(suggestion string) *ValidationError {
	ve.Suggestion = suggestion

	return ve
}

// cloneStruct creates a deep copy of a struct by value.
func cloneStruct[T any](item T) *T {
	clone := item

	return &clone
}

// Clone creates a deep copy of the validation error.
func (ve *ValidationError) Clone() *ValidationError {
	return cloneStruct(*ve)
}

// ValidationWarning represents a structured validation warning.
type ValidationWarning struct {
	Code       string       `json:"code"`
	Field      string       `json:"field"`
	Message    string       `json:"message"`
	Details    string       `json:"details,omitempty"`
	Context    string       `json:"context,omitempty"`
	Level      WarningLevel `json:"level"`
	Suggestion string       `json:"suggestion,omitempty"`
}

// String returns a string representation of the warning.
func (vw *ValidationWarning) String() string {
	return formatValidationItem(vw.Code, vw.Field, vw.Message, vw.Context)
}

// WithContext adds context to the validation warning.
func (vw *ValidationWarning) WithContext(context string) *ValidationWarning {
	vw.Context = context

	return vw
}

// WithDetails adds details to the validation warning.
func (vw *ValidationWarning) WithDetails(details string) *ValidationWarning {
	vw.Details = details

	return vw
}

// WithSuggestion adds suggestion to the validation warning.
func (vw *ValidationWarning) WithSuggestion(suggestion string) *ValidationWarning {
	vw.Suggestion = suggestion

	return vw
}

// Clone creates a deep copy of the validation warning.
func (vw *ValidationWarning) Clone() *ValidationWarning {
	return cloneStruct(*vw)
}

// ValidationSummary provides a summary of validation results.
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

// GetScore returns a validation score (0-100).
func (vs *ValidationSummary) GetScore() int {
	if vs.TotalErrors > 0 {
		return 0
	}

	if vs.TotalWarnings == 0 {
		return perfectScore
	}

	// Deduct points based on warning severity
	score := perfectScore
	score -= vs.HighWarnings * highWarningDeduction
	score -= vs.MediumWarnings * mediumWarningDeduction
	score -= vs.LowWarnings * lowWarningDeduction

	if score < 0 {
		score = 0
	}

	return score
}

// GetGrade returns a validation grade.
func (vs *ValidationSummary) GetGrade() string {
	score := vs.GetScore()

	// Special case: perfect score
	if score == gradeThresholdAPlus {
		return "A+"
	}

	// Grade thresholds in descending order for efficient lookup
	gradeThresholds := []struct {
		threshold int
		grade     string
	}{
		{gradeThresholdA, "A"},
		{gradeThresholdAMinus, "A-"},
		{gradeThresholdBPlus, "B+"},
		{gradeThresholdB, "B"},
		{gradeThresholdBMinus, "B-"},
		{gradeThresholdCPlus, "C+"},
		{gradeThresholdC, "C"},
		{gradeThresholdCMinus, "C-"},
		{gradeThresholdD, "D"},
	}

	for _, g := range gradeThresholds {
		if score >= g.threshold {
			return g.grade
		}
	}

	return "F"
}

// GetStatus returns a validation status.
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

// ValidationStatus represents validation status.
type ValidationStatus string

const (
	ValidationStatusPassed   ValidationStatus = "passed"
	ValidationStatusNotice   ValidationStatus = "notice"
	ValidationStatusWarning  ValidationStatus = "warning"
	ValidationStatusCritical ValidationStatus = "critical"
	ValidationStatusFailed   ValidationStatus = "failed"
)

// ValidationStatusInfo holds display information for a ValidationStatus.
type ValidationStatusInfo struct {
	Icon  string
	Color string
}

// validationStatusInfo is a map of ValidationStatus to its display information.
var validationStatusInfo = map[ValidationStatus]ValidationStatusInfo{
	ValidationStatusPassed:   {Icon: "✅", Color: "green"},
	ValidationStatusNotice:   {Icon: "ℹ️", Color: "blue"},
	ValidationStatusWarning:  {Icon: "⚠️", Color: "yellow"},
	ValidationStatusCritical: {Icon: "🚨", Color: "red"},
	ValidationStatusFailed:   {Icon: "❌", Color: "red"},
}

const validationStatusDefaultIcon = "❓"
const validationStatusDefaultColor = "gray"

// String returns string representation.
func (vs ValidationStatus) String() string {
	return string(vs)
}

// GetIcon returns an icon for the status.
func (vs ValidationStatus) GetIcon() string {
	if info, ok := validationStatusInfo[vs]; ok {
		return info.Icon
	}
	return validationStatusDefaultIcon
}

// GetColor returns a color for the status.
func (vs ValidationStatus) GetColor() string {
	if info, ok := validationStatusInfo[vs]; ok {
		return info.Color
	}
	return validationStatusDefaultColor
}

// ErrorLevel represents error severity levels.
type ErrorLevel string

const (
	ErrorLevelCritical ErrorLevel = "critical"
	ErrorLevelHigh     ErrorLevel = "high"
	ErrorLevelMedium   ErrorLevel = "medium"
	ErrorLevelLow      ErrorLevel = "low"
)

// String returns string representation.
func (el ErrorLevel) String() string {
	return string(el)
}

// GetPriority returns numeric priority for sorting.
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

// WarningLevel represents warning severity levels.
type WarningLevel string

const (
	WarningLevelHigh   WarningLevel = "high"
	WarningLevelMedium WarningLevel = "medium"
	WarningLevelLow    WarningLevel = "low"
)

// String returns string representation.
func (wl WarningLevel) String() string {
	return string(wl)
}

// GetPriority returns numeric priority for sorting.
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

// ValidationFilter represents filters for validation results.
type ValidationFilter struct {
	Fields        []string           `json:"fields,omitempty"`
	Levels        []ErrorLevel       `json:"error_levels,omitempty"`
	Codes         []errors.ErrorCode `json:"error_codes,omitempty"`
	WarningLevels []WarningLevel     `json:"warning_levels,omitempty"`
	WarningCodes  []string           `json:"warning_codes,omitempty"`
	Context       string             `json:"context,omitempty"`
	MinScore      int                `json:"min_score,omitempty"`
}

// Filter applies the filter to a validation result.
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

// filterItem represents a filterable validation item with common fields.
type filterItem interface {
	getField() string
	getContext() string
}

// matchesFilter checks if an item matches the filter criteria using the provided checkers.
func matchesFilter[T filterItem](
	item T,
	fields []string,
	context string,
	levelMatches func() bool,
	codeMatches func() bool,
) bool {
	// Check fields
	if len(fields) > 0 && !contains(fields, item.getField()) {
		return false
	}

	// Check levels
	if !levelMatches() {
		return false
	}

	// Check codes
	if !codeMatches() {
		return false
	}

	// Check context
	if context != "" && item.getContext() != context {
		return false
	}

	return true
}

// filterableItem represents an item that can be filtered with level and code.
type filterableItem interface {
	filterItem
	getLevel() any
	getCode() any
}

// matchesGeneric checks if an item matches the filter criteria.
func matchesGeneric[T filterableItem](
	f *ValidationFilter,
	item T,
	levels any,
	codes any,
	levelContains func(any, any) bool,
	codeContains func(any, any) bool,
) bool {
	return matchesFilter(
		item,
		f.Fields,
		f.Context,
		func() bool {
			return isZero(levels) || levelContains(levels, item.getLevel())
		},
		func() bool {
			return isZero(codes) || codeContains(codes, item.getCode())
		},
	)
}

// isZero checks if a value is zero/empty.
func isZero(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case []ErrorLevel:
		return len(val) == 0
	case []errors.ErrorCode:
		return len(val) == 0
	case []WarningLevel:
		return len(val) == 0
	case []string:
		return len(val) == 0
	default:
		return false
	}
}

// matchesError checks if an error matches the filter criteria.
func (f *ValidationFilter) matchesError(err *ValidationError) bool {
	return matchesGeneric(
		f,
		err,
		f.Levels,
		f.Codes,
		func(levels, level any) bool {
			levelsSlice, ok1 := levels.([]ErrorLevel)

			levelVal, ok2 := level.(ErrorLevel)
			if !ok1 || !ok2 {
				return false
			}

			return containsErrorLevels(levelsSlice, levelVal)
		},
		func(codes, code any) bool {
			codesSlice, ok1 := codes.([]errors.ErrorCode)

			codeVal, ok2 := code.(errors.ErrorCode)
			if !ok1 || !ok2 {
				return false
			}

			return containsErrorCodes(codesSlice, codeVal)
		},
	)
}

// matchesWarning checks if a warning matches the filter criteria.
func (f *ValidationFilter) matchesWarning(warn *ValidationWarning) bool {
	return matchesGeneric(
		f,
		warn,
		f.WarningLevels,
		f.WarningCodes,
		func(levels, level any) bool {
			levelsSlice, ok1 := levels.([]WarningLevel)

			levelVal, ok2 := level.(WarningLevel)
			if !ok1 || !ok2 {
				return false
			}

			return containsWarningLevels(levelsSlice, levelVal)
		},
		func(codes, code any) bool {
			codesSlice, ok1 := codes.([]string)

			codeVal, ok2 := code.(string)
			if !ok1 || !ok2 {
				return false
			}

			return contains(codesSlice, codeVal)
		},
	)
}

// getLevel implements filterableItem for ValidationError.
func (e *ValidationError) getLevel() any {
	return e.Level
}

// getCode implements filterableItem for ValidationError.
func (e *ValidationError) getCode() any {
	return e.Code
}

// getLevel implements filterableItem for ValidationWarning.
func (w *ValidationWarning) getLevel() any {
	return w.Level
}

// getCode implements filterableItem for ValidationWarning.
func (w *ValidationWarning) getCode() any {
	return w.Code
}

// getField implements filterItem for ValidationError.
func (e *ValidationError) getField() string {
	return e.Field
}

// getContext implements filterItem for ValidationError.
func (e *ValidationError) getContext() string {
	return e.Context
}

// getField implements filterItem for ValidationWarning.
func (w *ValidationWarning) getField() string {
	return w.Field
}

// getContext implements filterItem for ValidationWarning.
func (w *ValidationWarning) getContext() string {
	return w.Context
}

// Helper functions.
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
