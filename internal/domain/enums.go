package domain

const unknownValue = "Unknown"

// CGOStatus represents CGO compilation status with compile-time safety
// Replaces bool CGOEnabled for better type safety and semantic clarity.
type CGOStatus string

const (
	// CGOStatusDisabled disables CGO compilation completely.
	CGOStatusDisabled CGOStatus = "disabled"
	// CGOStatusEnabled enables CGO compilation when available.
	CGOStatusEnabled CGOStatus = "enabled"
	// CGOStatusRequired requires CGO compilation and will fail if not available.
	CGOStatusRequired CGOStatus = "required"
)

// IsValid returns true if CGOStatus is valid.
func (cs CGOStatus) IsValid() bool {
	switch cs {
	case CGOStatusDisabled, CGOStatusEnabled, CGOStatusRequired:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (cs CGOStatus) String() string {
	switch cs {
	case CGOStatusDisabled:
		return "Disabled"
	case CGOStatusEnabled:
		return "Enabled"
	case CGOStatusRequired:
		return "Required"
	default:
		return unknownValue
	}
}

// IsEnabled returns true if CGO is enabled (enabled or required).
func (cs CGOStatus) IsEnabled() bool {
	return cs == CGOStatusEnabled || cs == CGOStatusRequired
}

// IsDisabled returns true if CGO is disabled.
func (cs CGOStatus) IsDisabled() bool {
	return cs == CGOStatusDisabled
}

// IsRequired returns true if CGO is required.
func (cs CGOStatus) IsRequired() bool {
	return cs == CGOStatusRequired
}

// ToBool converts to legacy boolean for compatibility.
func (cs CGOStatus) ToBool() bool {
	return cs.IsEnabled()
}

// ValidateCGOStatus validates a CGO status.
func ValidateCGOStatus(status CGOStatus) error {
	return validateEnum("CGO status", string(status), status.IsValid())
}

// SigningLevel represents code signing level with compile-time safety
// Replaces bool Signing for better type safety and semantic clarity.
type SigningLevel string

const (
	// SigningLevelNone disables code signing completely.
	SigningLevelNone SigningLevel = "none"
	// SigningLevelBasic enables basic code signing.
	SigningLevelBasic SigningLevel = "basic"
	// SigningLevelAdvanced enables advanced code signing with additional verification.
	SigningLevelAdvanced SigningLevel = "advanced"
	// SigningLevelEnterprise enables enterprise-level code signing with full compliance.
	SigningLevelEnterprise SigningLevel = "enterprise"
)

// IsValid returns true if SigningLevel is valid.
func (sl SigningLevel) IsValid() bool {
	switch sl {
	case SigningLevelNone, SigningLevelBasic, SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
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
		return unknownValue
	}
}

// IsEnabled returns true if signing is enabled (basic, advanced, or enterprise).
func (sl SigningLevel) IsEnabled() bool {
	return sl != SigningLevelNone
}

// RequiresCosign returns true if level requires cosign.
func (sl SigningLevel) RequiresCosign() bool {
	switch sl {
	case SigningLevelNone, SigningLevelBasic:
		return false
	case SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// RequiresKeyManagement returns true if level requires key management.
func (sl SigningLevel) RequiresKeyManagement() bool {
	switch sl {
	case SigningLevelNone:
		return false
	case SigningLevelBasic, SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// GetRequiredTools returns tools required for this signing level.
func (sl SigningLevel) GetRequiredTools() []string {
	switch sl {
	case SigningLevelNone:
		return []string{}
	case SigningLevelBasic:
		return []string{"gpg"}
	case SigningLevelAdvanced:
		return []string{"gpg", "cosign"}
	case SigningLevelEnterprise:
		return []string{"gpg", "cosign", "openssl"}
	default:
		return []string{}
	}
}

// ToBool converts to legacy boolean for compatibility.
func (sl SigningLevel) ToBool() bool {
	return sl.IsEnabled()
}

// ValidateSigningLevel validates a signing level.
func ValidateSigningLevel(level SigningLevel) error {
	return validateEnum("signing level", string(level), level.IsValid())
}

// Enum migration utilities for backward compatibility

// CGOStatusFromBool converts legacy boolean to CGOStatus.
func CGOStatusFromBool(enabled bool) CGOStatus {
	if enabled {
		return CGOStatusEnabled
	}

	return CGOStatusDisabled
}

// DockerSupportFromBool converts legacy boolean to DockerSupport.
func DockerSupportFromBool(enabled bool) DockerSupport {
	if enabled {
		return DockerSupportBoth
	}

	return DockerSupportNone
}

// SigningLevelFromBool converts legacy boolean to SigningLevel.
func SigningLevelFromBool(enabled bool) SigningLevel {
	if enabled {
		return SigningLevelBasic
	}

	return SigningLevelNone
}

// ActionLevelFromBool converts legacy boolean to ActionLevel.
func ActionLevelFromBool(enabled bool) ActionLevel {
	if enabled {
		return ActionLevelBasic
	}

	return ActionLevelNone
}

// FeatureLevelFromBool converts legacy boolean to FeatureLevel.
func FeatureLevelFromBool(enabled bool) FeatureLevel {
	if enabled {
		return FeatureLevelStandard
	}

	return FeatureLevelBasic
}

// Smart conversion functions based on project type context

// GetDefaultCGOStatus returns smart CGO status based on project type.
func GetDefaultCGOStatus(projectType ProjectType) CGOStatus {
	if projectType.DefaultCGOEnabled() {
		return CGOStatusEnabled
	}

	return CGOStatusDisabled
}

// GetDefaultDockerSupport returns smart Docker support based on project type.
func GetDefaultDockerSupport(projectType ProjectType) DockerSupport {
	if projectType.DockerSupported() {
		return DockerSupportBuild
	}

	return DockerSupportNone
}

// GetRecommendedActionLevel returns recommended action level based on project type.
func GetRecommendedActionLevel(projectType ProjectType) ActionLevel {
	switch projectType {
	case ProjectTypeCLI, ProjectTypeWebAPI:
		return ActionLevelAdvanced
	case ProjectTypeMicroservice:
		return ActionLevelBasic
	case ProjectTypeLibrary:
		return ActionLevelBasic
	case ProjectTypeDesktop:
		return ActionLevelAdvanced
	case ProjectTypeGRPCService:
		return ActionLevelAdvanced
	case ProjectTypeMobile, ProjectTypePlugin, ProjectTypeDaemon, ProjectTypeTool:
		return ActionLevelBasic
	default:
		return ActionLevelBasic
	}
}

// GetRecommendedSigningLevel returns recommended signing level based on project type.
func GetRecommendedSigningLevel(projectType ProjectType) SigningLevel {
	switch projectType {
	case ProjectTypeCLI:
		return SigningLevelBasic
	case ProjectTypeWebAPI, ProjectTypeGRPCService:
		return SigningLevelAdvanced
	case ProjectTypeDesktop:
		return SigningLevelEnterprise
	case ProjectTypeLibrary:
		return SigningLevelNone
	case ProjectTypeMicroservice,
		ProjectTypeMobile,
		ProjectTypePlugin,
		ProjectTypeDaemon,
		ProjectTypeTool:
		return SigningLevelBasic
	default:
		return SigningLevelBasic
	}
}

// GetRecommendedFeatureLevel returns recommended feature level based on project type.
func GetRecommendedFeatureLevel(projectType ProjectType) FeatureLevel {
	switch projectType {
	case ProjectTypeWebAPI, ProjectTypeGRPCService:
		return FeatureLevelStandard
	case ProjectTypeDesktop:
		return FeatureLevelEnterprise
	case ProjectTypeCLI, ProjectTypeLibrary:
		return FeatureLevelBasic
	case ProjectTypeMicroservice,
		ProjectTypeMobile,
		ProjectTypePlugin,
		ProjectTypeDaemon,
		ProjectTypeTool:
		return FeatureLevelBasic
	default:
		return FeatureLevelBasic
	}
}
