package domain

import (
	"fmt"
)

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
