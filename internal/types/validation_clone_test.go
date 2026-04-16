package types

import (
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// verifyValidationErrorCloneFields verifies all fields match between original and clone.
func verifyValidationErrorCloneFields(t *testing.T, original, clone *ValidationError) {
	if clone.Code != original.Code {
		t.Errorf("Code mismatch: got %v, want %v", clone.Code, original.Code)
	}

	if clone.Field != original.Field {
		t.Errorf("Field mismatch: got %v, want %v", clone.Field, original.Field)
	}

	if clone.Message != original.Message {
		t.Errorf("Message mismatch: got %v, want %v", clone.Message, original.Message)
	}

	if clone.Details != original.Details {
		t.Errorf("Details mismatch: got %v, want %v", clone.Details, original.Details)
	}

	if clone.Context != original.Context {
		t.Errorf("Context mismatch: got %v, want %v", clone.Context, original.Context)
	}

	if clone.Level != original.Level {
		t.Errorf("Level mismatch: got %v, want %v", clone.Level, original.Level)
	}

	if clone.Suggestion != original.Suggestion {
		t.Errorf("Suggestion mismatch: got %v, want %v", clone.Suggestion, original.Suggestion)
	}
}

// verifyValidationWarningCloneFields verifies all fields match between original and clone.
func verifyValidationWarningCloneFields(t *testing.T, original, clone *ValidationWarning) {
	if clone.Code != original.Code {
		t.Errorf("Code mismatch: got %v, want %v", clone.Code, original.Code)
	}

	if clone.Field != original.Field {
		t.Errorf("Field mismatch: got %v, want %v", clone.Field, original.Field)
	}

	if clone.Message != original.Message {
		t.Errorf("Message mismatch: got %v, want %v", clone.Message, original.Message)
	}

	if clone.Details != original.Details {
		t.Errorf("Details mismatch: got %v, want %v", clone.Details, original.Details)
	}

	if clone.Context != original.Context {
		t.Errorf("Context mismatch: got %v, want %v", clone.Context, original.Context)
	}

	if clone.Level != original.Level {
		t.Errorf("Level mismatch: got %v, want %v", clone.Level, original.Level)
	}

	if clone.Suggestion != original.Suggestion {
		t.Errorf("Suggestion mismatch: got %v, want %v", clone.Suggestion, original.Suggestion)
	}
}

// verifyValidationErrorIndependence verifies modifying clone doesn't affect original.
func verifyValidationErrorIndependence(t *testing.T, original, clone *ValidationError) {
	if original.Code == clone.Code {
		t.Error("Clone modification affected original Code")
	}

	if original.Field == clone.Field {
		t.Error("Clone modification affected original Field")
	}

	if original.Message == clone.Message {
		t.Error("Clone modification affected original Message")
	}

	if original.Details == clone.Details {
		t.Error("Clone modification affected original Details")
	}

	if original.Context == clone.Context {
		t.Error("Clone modification affected original Context")
	}

	if original.Level == clone.Level {
		t.Error("Clone modification affected original Level")
	}

	if original.Suggestion == clone.Suggestion {
		t.Error("Clone modification affected original Suggestion")
	}
}

// verifyValidationWarningIndependence verifies modifying clone doesn't affect original.
func verifyValidationWarningIndependence(t *testing.T, original, clone *ValidationWarning) {
	if original.Code == clone.Code {
		t.Error("Clone modification affected original Code")
	}

	if original.Field == clone.Field {
		t.Error("Clone modification affected original Field")
	}

	if original.Message == clone.Message {
		t.Error("Clone modification affected original Message")
	}

	if original.Details == clone.Details {
		t.Error("Clone modification affected original Details")
	}

	if original.Context == clone.Context {
		t.Error("Clone modification affected original Context")
	}

	if original.Level == clone.Level {
		t.Error("Clone modification affected original Level")
	}

	if original.Suggestion == clone.Suggestion {
		t.Error("Clone modification affected original Suggestion")
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

	verifyValidationErrorCloneFields(t, original, clone)

	// Modify clone
	clone.Code = errors.ErrInvalidConfig
	clone.Field = "modified"
	clone.Message = "modified message"
	clone.Details = "modified details"
	clone.Context = "modified context"
	clone.Level = ErrorLevelCritical
	clone.Suggestion = "modified suggestion"

	verifyValidationErrorIndependence(t, original, clone)

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

	verifyValidationWarningCloneFields(t, original, clone)

	// Modify clone
	clone.Code = "WARN_002"
	clone.Field = "modified"
	clone.Message = "modified message"
	clone.Details = "modified details"
	clone.Context = "modified context"
	clone.Level = WarningLevelHigh
	clone.Suggestion = "modified suggestion"

	verifyValidationWarningIndependence(t, original, clone)

	// Verify they are different instances
	if clone == original {
		t.Error("Clone returns same instance, should be a different pointer")
	}
}
