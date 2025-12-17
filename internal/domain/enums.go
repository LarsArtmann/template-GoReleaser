// CRITICAL ARCHITECTURE TODO: This file is 429 lines - SPLIT IMMEDIATELY by entity:
// 1. enums_platform.go - Platform and Architecture enums
// 2. enums_build.go - CGOStatus, BuildTag, BuildLevel enums
// 3. enums_release.go - GitProvider, DockerRegistry, SigningLevel enums
// 4. enums_project.go - ProjectType, FeatureLevel enums
// 5. enums_actions.go - ActionLevel, ActionTrigger enums
// 6. enums_state.go - ConfigState and state-related enums
//
// TODO: Replace string-based enums with proper typed enums
// TODO: Add enum validation methods
// TODO: Create enum conversion utilities
// TODO: Implement proper enum serialization/deserialization
// TODO: Add enum compatibility checking methods
package domain

import (
	"fmt"
)

// CGOStatus represents CGO compilation status with compile-time safety
// Replaces bool CGOEnabled for better type safety and semantic clarity
type CGOStatus string

const (
	// CGOStatusDisabled disables CGO compilation completely
	CGOStatusDisabled CGOStatus = "disabled"
	// CGOStatusEnabled enables CGO compilation when available
	CGOStatusEnabled CGOStatus = "enabled"
	// CGOStatusRequired requires CGO compilation and will fail if not available
	CGOStatusRequired CGOStatus = "required"
)

// IsValid returns true if CGOStatus is valid
func (cs CGOStatus) IsValid() bool {
	switch cs {
	case CGOStatusDisabled, CGOStatusEnabled, CGOStatusRequired:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (cs CGOStatus) String() string {
	switch cs {
	case CGOStatusDisabled:
		return "Disabled"
	case CGOStatusEnabled:
		return "Enabled"
	case CGOStatusRequired:
		return "Required"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if CGO is enabled (enabled or required)
func (cs CGOStatus) IsEnabled() bool {
	return cs == CGOStatusEnabled || cs == CGOStatusRequired
}

// IsRequired returns true if CGO is required
func (cs CGOStatus) IsRequired() bool {
	return cs == CGOStatusRequired
}

// ToBool converts to legacy boolean for compatibility
func (cs CGOStatus) ToBool() bool {
	return cs.IsEnabled()
}

// ValidateCGOStatus validates a CGO status
func ValidateCGOStatus(status CGOStatus) error {
	if !status.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid CGO status",
			fmt.Sprintf("'%s' is not a valid CGO status", status),
		)
	}
	return nil
}

// SigningLevel represents code signing level with compile-time safety
// Replaces bool Signing for better type safety and semantic clarity
type SigningLevel string

const (
	// SigningLevelNone disables code signing completely
	SigningLevelNone SigningLevel = "none"
	// SigningLevelBasic enables basic code signing
	SigningLevelBasic SigningLevel = "basic"
	// SigningLevelAdvanced enables advanced code signing with additional verification
	SigningLevelAdvanced SigningLevel = "advanced"
	// SigningLevelEnterprise enables enterprise-level code signing with full compliance
	SigningLevelEnterprise SigningLevel = "enterprise"
)

// IsValid returns true if SigningLevel is valid
func (sl SigningLevel) IsValid() bool {
	switch sl {
	case SigningLevelNone, SigningLevelBasic, SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (sl SigningLevel) String() string {
	switch sl {
	case SigningLevelNone:
		return "None"
	case SigningLevelBasic:
		return "Basic"
	case SigningLevelAdvanced:
		return "Advanced"
	case SigningLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if signing is enabled (basic, advanced, or enterprise)
func (sl SigningLevel) IsEnabled() bool {
	return sl != SigningLevelNone
}

// ToBool converts to legacy boolean for compatibility
func (sl SigningLevel) ToBool() bool {
	return sl.IsEnabled()
}

// ValidateSigningLevel validates a signing level
func ValidateSigningLevel(level SigningLevel) error {
	if !level.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid signing level",
			fmt.Sprintf("'%s' is not a valid signing level", level),
		)
	}
	return nil
}



// Enum migration utilities for backward compatibility

// CGOStatusFromBool converts legacy boolean to CGOStatus
func CGOStatusFromBool(enabled bool) CGOStatus {
	if enabled {
		return CGOStatusEnabled
	}
	return CGOStatusDisabled
}

// DockerSupportFromBool converts legacy boolean to DockerSupport
func DockerSupportFromBool(enabled bool) DockerSupport {
	if enabled {
		return DockerSupportBoth
	}
	return DockerSupportNone
}

// SigningLevelFromBool converts legacy boolean to SigningLevel
func SigningLevelFromBool(enabled bool) SigningLevel {
	if enabled {
		return SigningLevelBasic
	}
	return SigningLevelNone
}

// ActionLevelFromBool converts legacy boolean to ActionLevel
func ActionLevelFromBool(enabled bool) ActionLevel {
	if enabled {
		return ActionLevelBasic
	}
	return ActionLevelNone
}

// FeatureLevelFromBool converts legacy boolean to FeatureLevel
func FeatureLevelFromBool(enabled bool) FeatureLevel {
	if enabled {
		return FeatureLevelProfessional
	}
	return FeatureLevelBasic
}

// Smart conversion functions based on project type context

// GetDefaultCGOStatus returns smart CGO status based on project type
func GetDefaultCGOStatus(projectType ProjectType) CGOStatus {
	if projectType.DefaultCGOEnabled() {
		return CGOStatusEnabled
	}
	return CGOStatusDisabled
}

// GetDefaultDockerSupport returns smart Docker support based on project type
func GetDefaultDockerSupport(projectType ProjectType) DockerSupport {
	if projectType.DockerSupported() {
		return DockerSupportBuild
	}
	return DockerSupportNone
}

// GetRecommendedActionLevel returns recommended action level based on project type
func GetRecommendedActionLevel(projectType ProjectType) ActionLevel {
	switch projectType {
	case ProjectTypeCLI, ProjectTypeAPI:
		return ActionLevelAdvanced
	case ProjectTypeWeb:
		return ActionLevelBasic
	case ProjectTypeLibrary:
		return ActionLevelBasic
	case ProjectTypeDesktop:
		return ActionLevelAdvanced
	default:
		return ActionLevelBasic
	}
}

// GetRecommendedSigningLevel returns recommended signing level based on project type
func GetRecommendedSigningLevel(projectType ProjectType) SigningLevel {
	switch projectType {
	case ProjectTypeCLI:
		return SigningLevelBasic
	case ProjectTypeWeb, ProjectTypeAPI:
		return SigningLevelAdvanced
	case ProjectTypeDesktop:
		return SigningLevelEnterprise
	case ProjectTypeLibrary:
		return SigningLevelNone
	default:
		return SigningLevelBasic
	}
}

// GetRecommendedFeatureLevel returns recommended feature level based on project type
func GetRecommendedFeatureLevel(projectType ProjectType) FeatureLevel {
	switch projectType {
	case ProjectTypeAPI, ProjectTypeWeb:
		return FeatureLevelProfessional
	case ProjectTypeDesktop:
		return FeatureLevelEnterprise
	case ProjectTypeCLI, ProjectTypeLibrary:
		return FeatureLevelBasic
	default:
		return FeatureLevelBasic
	}
}
