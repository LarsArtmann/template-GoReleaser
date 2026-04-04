package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// FormValidator provides validation functions for huh forms.
type FormValidator struct {
	errors map[string]string
}

// NewFormValidator creates a new form validator.
func NewFormValidator() *FormValidator {
	return &FormValidator{
		errors: make(map[string]string),
	}
}

// createValidator is a generic helper for creating validation closures.
// It handles error tracking and clearing for any value type.
func createValidator[T any](
	fv *FormValidator,
	fieldName string,
	validateFunc func(T) error,
) func(T) error {
	return func(value T) error {
		err := validateFunc(value)
		if err != nil {
			fv.errors[fieldName] = err.Error()

			return err
		}

		delete(fv.errors, fieldName)

		return nil
	}
}

// ValidateProjectName creates a huh-compatible validator for project names.
func (fv *FormValidator) ValidateProjectName() func(string) error {
	return createValidator(fv, "project_name", ValidateProjectName)
}

// ValidateBinaryName creates a huh-compatible validator for binary names.
func (fv *FormValidator) ValidateBinaryName() func(string) error {
	return createValidator(fv, "binary_name", ValidateBinaryName)
}

// ValidateMainPath creates a huh-compatible validator for main path.
func (fv *FormValidator) ValidateMainPath() func(string) error {
	return createValidator(fv, "main_path", ValidateMainPath)
}

// ValidateProjectDescription creates a huh-compatible validator for project description.
func (fv *FormValidator) ValidateProjectDescription() func(string) error {
	return createValidator(fv, "project_description", ValidateProjectDescription)
}

// ValidateDockerRegistry creates a huh-compatible validator for Docker registry.
func (fv *FormValidator) ValidateDockerRegistry() func(string) error {
	return createValidator(fv, "docker_registry", ValidateDockerRegistry)
}

// ValidateBuildTags creates a huh-compatible validator for build tags.
func (fv *FormValidator) ValidateBuildTags(tags []string) error {
	return createValidator(fv, "build_tags", ValidateBuildTags)(tags)
}

// GetErrors returns all current validation errors.
func (fv *FormValidator) GetErrors() map[string]string {
	return fv.errors
}

// HasErrors returns true if there are validation errors.
func (fv *FormValidator) HasErrors() bool {
	return len(fv.errors) > 0
}

// GetErrorSummary returns a formatted summary of all errors.
func (fv *FormValidator) GetErrorSummary() string {
	if !fv.HasErrors() {
		return ""
	}

	var errors []string
	for field, message := range fv.errors {
		errors = append(errors, fmt.Sprintf("• %s: %s", field, message))
	}

	return "Validation errors:\n" + strings.Join(errors, "\n")
}

// ClearErrors clears all current validation errors.
func (fv *FormValidator) ClearErrors() {
	fv.errors = make(map[string]string)
}

// SanitizeAndValidate sanitizes input and validates it.
func (fv *FormValidator) SanitizeAndValidate(
	input string,
	validator func(string) error,
) (string, error) {
	// First sanitize the input
	sanitized := SanitizeInput(input)

	// Then validate it
	err := validator(sanitized)
	if err != nil {
		return "", err
	}

	return sanitized, nil
}

// GetFieldError gets a specific field error.
func (fv *FormValidator) GetFieldError(field string) string {
	return fv.errors[field]
}

// SetFieldError manually sets a field error.
func (fv *FormValidator) SetFieldError(field, message string) {
	fv.errors[field] = message
}

// ValidateRequired validates that a field is not empty.
func (fv *FormValidator) ValidateRequired(fieldName string) func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			err := fmt.Errorf("%s is required", fieldName)
			fv.errors[fieldName] = err.Error()

			return err
		}

		delete(fv.errors, fieldName)

		return nil
	}
}

// ValidateLength validates string length within bounds.
func (fv *FormValidator) ValidateLength(minLen, maxLen int, fieldName string) func(string) error {
	return func(value string) error {
		length := len(strings.TrimSpace(value))
		if length < minLen || length > maxLen {
			err := fmt.Errorf("%s must be between %d and %d characters", fieldName, minLen, maxLen)
			fv.errors[fieldName] = err.Error()

			return err
		}

		delete(fv.errors, fieldName)

		return nil
	}
}

// validatePattern is a generic helper for pattern-based validation.
func (fv *FormValidator) validatePattern(
	fieldName string,
	pattern *regexp.Regexp,
	errorMessage string,
) func(string) error {
	return func(value string) error {
		if pattern.MatchString(value) {
			err := fmt.Errorf("%s %s", fieldName, errorMessage)
			fv.errors[fieldName] = err.Error()

			return err
		}

		delete(fv.errors, fieldName)

		return nil
	}
}

// ValidateNoShellMetacharacters validates that input doesn't contain shell metacharacters.
func (fv *FormValidator) ValidateNoShellMetacharacters(fieldName string) func(string) error {
	return fv.validatePattern(fieldName, shellMetacharPattern, "contains dangerous characters")
}

// ValidateNoPathTraversal validates that input doesn't contain path traversal.
func (fv *FormValidator) ValidateNoPathTraversal(fieldName string) func(string) error {
	return fv.validatePattern(fieldName, pathTraversalPattern, "contains path traversal attempts")
}
