package types

import (
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// validationItemEx extends validationItem with Code, Details, and Context accessors.
type validationItemEx interface {
	validationItem
	GetCode() string
	GetDetails() string
	GetContext() string
}

// GetCode for ValidationError.
func (ve *ValidationError) GetCode() string { return string(ve.Code) }

// GetDetails for ValidationError.
func (ve *ValidationError) GetDetails() string { return ve.Details }

// GetContext for ValidationError.
func (ve *ValidationError) GetContext() string { return ve.Context }

// GetCode for ValidationWarning.
func (vw *ValidationWarning) GetCode() string { return vw.Code }

// GetDetails for ValidationWarning.
func (vw *ValidationWarning) GetDetails() string { return vw.Details }

// GetContext for ValidationWarning.
func (vw *ValidationWarning) GetContext() string { return vw.Context }

// verifyCloneFields verifies all fields match between original and clone.
func verifyCloneFields(t *testing.T, name string, original, clone validationItemEx) {
	checkField(t, name, "Code", clone.GetCode(), original.GetCode())
	checkField(t, name, "Field", clone.GetField(), original.GetField())
	checkField(t, name, "Message", clone.GetMessage(), original.GetMessage())
	checkField(t, name, "Details", clone.GetDetails(), original.GetDetails())
	checkField(t, name, "Context", clone.GetContext(), original.GetContext())
	checkField(t, name, "Level", clone.GetLevel(), original.GetLevel())
	checkField(t, name, "Suggestion", clone.GetSuggestion(), original.GetSuggestion())
}

// verifyIndependence verifies modifying clone doesn't affect original.
func verifyIndependence(t *testing.T, name string, original, clone validationItemEx) {
	if original.GetCode() == clone.GetCode() {
		t.Errorf("%s: Clone modification affected original Code", name)
	}

	if original.GetField() == clone.GetField() {
		t.Errorf("%s: Clone modification affected original Field", name)
	}

	if original.GetMessage() == clone.GetMessage() {
		t.Errorf("%s: Clone modification affected original Message", name)
	}

	if original.GetDetails() == clone.GetDetails() {
		t.Errorf("%s: Clone modification affected original Details", name)
	}

	if original.GetContext() == clone.GetContext() {
		t.Errorf("%s: Clone modification affected original Context", name)
	}

	if original.GetLevel() == clone.GetLevel() {
		t.Errorf("%s: Clone modification affected original Level", name)
	}

	if original.GetSuggestion() == clone.GetSuggestion() {
		t.Errorf("%s: Clone modification affected original Suggestion", name)
	}
}

func checkField(t *testing.T, name, field, got, want string) {
	if got != want {
		t.Errorf("%s %s mismatch: got %v, want %v", name, field, got, want)
	}
}

func TestValidationError_Clone(t *testing.T) {
	original := &ValidationError{
		Code:       errors.ErrInvalidField,
		Field:      "name",
		Message:    "Invalid name",
		Details:    "Name contains invalid characters",
		Context:    "validation",
		Level:      ErrorLevelHigh,
		Suggestion: "Use only lowercase letters and numbers",
	}

	clone := original.Clone()

	verifyCloneFields(t, "ValidationError", original, clone)

	// Modify clone
	clone.Code = errors.ErrInvalidConfig
	clone.Field = "modified"
	clone.Message = "modified message"
	clone.Details = "modified details"
	clone.Context = "modified context"
	clone.Level = ErrorLevelCritical
	clone.Suggestion = "modified suggestion"

	verifyIndependence(t, "ValidationError", original, clone)

	// Verify they are different instances
	if clone == original {
		t.Error("Clone returns same instance, should be a different pointer")
	}
}

func TestValidationWarning_Clone(t *testing.T) {
	original := &ValidationWarning{
		Code:       "WARN_001",
		Field:      "description",
		Message:    "Description too short",
		Details:    "Description should be at least 10 characters",
		Context:    "validation",
		Level:      WarningLevelMedium,
		Suggestion: "Add more details to the description",
	}

	clone := original.Clone()

	verifyCloneFields(t, "ValidationWarning", original, clone)

	// Modify clone
	clone.Code = "WARN_002"
	clone.Field = "modified"
	clone.Message = "modified message"
	clone.Details = "modified details"
	clone.Context = "modified context"
	clone.Level = WarningLevelHigh
	clone.Suggestion = "modified suggestion"

	verifyIndependence(t, "ValidationWarning", original, clone)

	// Verify they are different instances
	if clone == original {
		t.Error("Clone returns same instance, should be a different pointer")
	}
}
