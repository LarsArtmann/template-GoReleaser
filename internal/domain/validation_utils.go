package domain

import (
	"fmt"
)

// validateEnum validates a single enum value and returns a domain validation error.
// Use this for any single-value enum with an IsValid() method.
//
// Parameters:
//   - itemName: human-readable label (e.g. "CGO status", "platform")
//   - value: the stringified enum value
//   - isValid: the type's IsValid() method (must return false for invalid values)
//
// Returns nil when valid, otherwise a *DomainError that callers can use directly
// or wrap via WithContext/WithField.
func validateEnum(itemName, value string, isValid bool) error {
	if isValid {
		return nil
	}

	return NewValidationError(
		ErrInvalidCharacters,
		"Invalid "+itemName,
		fmt.Sprintf("'%s' is not a valid %s", value, itemName),
	)
}

// validateEnumSlice validates a slice of enum values with a single IsValid() method.
func validateEnumSlice[T ~string](values []T, itemName string, isValid func(T) bool) error {
	if len(values) == 0 {
		return fmt.Errorf("at least one %s is required", itemName)
	}

	for _, value := range values {
		if !isValid(value) {
			return fmt.Errorf("invalid %s: %s", itemName, value)
		}
	}

	return nil
}

// ValidatePlatforms validates slice of platforms.
func ValidatePlatforms(platforms []Platform) error {
	return validateEnumSlice(platforms, "platform", Platform.IsValid)
}

// ValidatePlatformArchCompatibility validates platform-architecture compatibility.
func ValidatePlatformArchCompatibility(platforms []Platform, architectures []Architecture) error {
	for _, platform := range platforms {
		for _, arch := range architectures {
			if !arch.IsCompatibleWith(platform) {
				return fmt.Errorf(
					"architecture %s is not compatible with platform %s",
					arch,
					platform,
				)
			}
		}
	}

	return nil
}
