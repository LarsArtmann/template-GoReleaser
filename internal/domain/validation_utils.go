package domain

import (
	"errors"
	"fmt"
)

// ValidatePlatforms validates slice of platforms.
func ValidatePlatforms(platforms []Platform) error {
	if len(platforms) == 0 {
		return errors.New("at least one platform is required")
	}

	for _, platform := range platforms {
		if !platform.IsValid() {
			return fmt.Errorf("invalid platform: %s", platform)
		}
	}

	return nil
}

// ValidatePlatformArchCompatibility validates platform-architecture compatibility.
func ValidatePlatformArchCompatibility(platforms []Platform, architectures []Architecture) error {
	for _, platform := range platforms {
		for _, arch := range architectures {
			if !arch.IsCompatibleWith(platform) {
				return fmt.Errorf("architecture %s is not compatible with platform %s", arch, platform)
			}
		}
	}

	return nil
}
