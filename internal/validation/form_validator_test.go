package validation

import (
	"testing"
)

// runValidationTest is a generic helper function that tests a validator with valid and invalid inputs
func runValidationTest[T any](t *testing.T, fv *FormValidator, validator func(T) error, field string, validInput, invalidInput T, checkFieldError bool) {
	// Test valid input
	err := validator(validInput)
	if err != nil {
		t.Errorf("%s valid input error = %v", field, err)
	}

	// Clear errors before testing invalid input
	fv.ClearErrors()

	// Test invalid input
	err = validator(invalidInput)
	if err == nil {
		t.Errorf("%s should error for invalid input", field)
	}

	if checkFieldError && fv.GetFieldError(field) == "" {
		t.Errorf("%s should set %s error", field, field)
	}
}

// runValidatorWithParamsTest is a helper function that tests a validator created with parameters
func runValidatorWithParamsTest(t *testing.T, fv *FormValidator, validator func(string) error, testName string, validInput, invalidInput string) {
	runValidationTest(t, fv, validator, testName, validInput, invalidInput, false)
}

func TestFormValidator(t *testing.T) {
	fv := NewFormValidator()

	// Test initial state
	if fv.HasErrors() {
		t.Errorf("NewFormValidator() should not have errors initially")
	}

	if len(fv.GetErrors()) != 0 {
		t.Errorf("NewFormValidator() should return empty errors map")
	}
}

func TestFormValidatorValidateProjectName(t *testing.T) {
	fv := NewFormValidator()

	// Test valid name and invalid name with dots
	runValidationTest(t, fv, fv.ValidateProjectName(), "project_name", "myproject", "invalid..name", true)

	// Clear errors and test
	fv.ClearErrors()
	if fv.HasErrors() {
		t.Errorf("ClearErrors() should clear all errors")
	}
}

func TestFormValidatorValidateBinaryName(t *testing.T) {
	fv := NewFormValidator()
	runValidationTest(t, fv, fv.ValidateBinaryName(), "binary_name", "myapp", "my;app", true)
}

func TestFormValidatorValidateMainPath(t *testing.T) {
	fv := NewFormValidator()
	runValidationTest(t, fv, fv.ValidateMainPath(), "main_path", "./cmd/app", "../../../etc/passwd", true)
}

func TestFormValidatorValidateProjectDescription(t *testing.T) {
	fv := NewFormValidator()
	runValidationTest(t, fv, fv.ValidateProjectDescription(), "project_description", "A great app", "<script>alert('xss')</script>", true)
}

func TestFormValidatorValidateDockerRegistry(t *testing.T) {
	fv := NewFormValidator()
	runValidationTest(t, fv, fv.ValidateDockerRegistry(), "docker_registry", "ghcr.io/username/app", "http://registry.example.com/app", true)
}

func TestFormValidatorValidateBuildTags(t *testing.T) {
	fv := NewFormValidator()
	runValidationTest(t, fv, fv.ValidateBuildTags, "build_tags", []string{"prod", "linux"}, []string{"prod;rm"}, true)
}

func TestFormValidatorErrorSummary(t *testing.T) {
	fv := NewFormValidator()

	// Empty summary
	summary := fv.GetErrorSummary()
	if summary != "" {
		t.Errorf("GetErrorSummary() should return empty string when no errors")
	}

	// Add some errors
	_ = fv.ValidateProjectName()("invalid..name")
	_ = fv.ValidateBinaryName()("con")

	summary = fv.GetErrorSummary()
	if summary == "" {
		t.Errorf("GetErrorSummary() should return non-empty string when errors exist")
	}

	if !contains(summary, "project_name") {
		t.Errorf("GetErrorSummary() should contain project_name error")
	}

	if !contains(summary, "binary_name") {
		t.Errorf("GetErrorSummary() should contain binary_name error")
	}
}

func TestFormValidatorSanitizeAndValidate(t *testing.T) {
	fv := NewFormValidator()

	// Test valid input
	sanitized, err := fv.SanitizeAndValidate("  myproject  ", fv.ValidateProjectName())
	if err != nil {
		t.Errorf("SanitizeAndValidate() error = %v", err)
	}

	if sanitized != "myproject" {
		t.Errorf("SanitizeAndValidate() = %v, want %v", sanitized, "myproject")
	}

	// Test invalid input
	_, err = fv.SanitizeAndValidate("con", fv.ValidateProjectName())
	if err == nil {
		t.Errorf("SanitizeAndValidate() should error for invalid input")
	}
}

func TestFormValidatorValidateRequired(t *testing.T) {
	fv := NewFormValidator()
	validator := fv.ValidateRequired("Test Field")
	runValidatorWithParamsTest(t, fv, validator, "ValidateRequired", "test", "")
}

func TestFormValidatorValidateLength(t *testing.T) {
	fv := NewFormValidator()

	validator := fv.ValidateLength(3, 10, "Test Field")

	// Test too short
	err := validator("ab")
	if err == nil {
		t.Errorf("ValidateLength() should error for too short input")
	}

	// Test too long
	err = validator("abcdefghijk")
	if err == nil {
		t.Errorf("ValidateLength() should error for too long input")
	}

	// Test valid length
	err = validator("abcde")
	if err != nil {
		t.Errorf("ValidateLength() should not error for valid length input")
	}
}

func TestFormValidatorValidateNoShellMetacharacters(t *testing.T) {
	fv := NewFormValidator()
	validator := fv.ValidateNoShellMetacharacters("Test Field")
	runValidatorWithParamsTest(t, fv, validator, "ValidateNoShellMetacharacters", "test-safe", "test;rm")
}

func TestFormValidatorValidateNoPathTraversal(t *testing.T) {
	fv := NewFormValidator()
	validator := fv.ValidateNoPathTraversal("Test Field")
	runValidatorWithParamsTest(t, fv, validator, "ValidateNoPathTraversal", "./cmd/app", "../../../etc")
}

func TestFormValidatorSetFieldError(t *testing.T) {
	fv := NewFormValidator()

	fv.SetFieldError("test_field", "test error message")

	if !fv.HasErrors() {
		t.Errorf("SetFieldError() should result in HasErrors() returning true")
	}

	if fv.GetFieldError("test_field") != "test error message" {
		t.Errorf("GetFieldError() = %v, want %v", fv.GetFieldError("test_field"), "test error message")
	}
}

// Helper function to check if string contains substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
