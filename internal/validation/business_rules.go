package validation

import (
	"fmt"
	"os/exec"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// ValidateConfiguration performs comprehensive validation of project configuration
func ValidateConfiguration(config *domain.SafeProjectConfig) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []*ValidationError{},
		Warnings: []*ValidationWarning{},
	}

	// Step 1: Basic field validation
	if err := validateBasicFields(config, result); err != nil {
		return nil, err
	}

	// Step 2: Type validation
	if err := validateTypes(config, result); err != nil {
		return nil, err
	}

	// Step 3: Platform-architecture compatibility
	if err := validatePlatformArchCompatibility(config, result); err != nil {
		return nil, err
	}

	// Step 4: Business rule validation
	if err := validateBusinessRules(config, result); err != nil {
		return nil, err
	}

	// Step 5: Security validation
	if err := validateSecurity(config, result); err != nil {
		return nil, err
	}

	// Step 6: Generate warnings
	generateWarnings(config, result)

	// Update result status
	result.UpdateSummary()

	return result, nil
}

// validateBasicFields validates basic required fields
func validateBasicFields(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Project name validation
	if err := ValidateProjectName(config.ProjectName); err != nil {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidProject,
			Field:      "project_name",
			Message:    err.Error(),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a descriptive project name following Go conventions",
		})
	}

	// Binary name validation
	if err := ValidateBinaryName(config.BinaryName); err != nil {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidBinary,
			Field:      "binary_name",
			Message:    err.Error(),
			Level:      ErrorLevelMedium,
			Suggestion: "Use a lowercase binary name with hyphens for multi-word names",
		})
	}

	// Main path validation
	if err := ValidateMainPath(config.MainPath); err != nil {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidMainPath,
			Field:      "main_path",
			Message:    err.Error(),
			Level:      ErrorLevelMedium,
			Suggestion: "Use a relative path like '.', './cmd/app', or './main.go'",
		})
	}

	// Project description validation (optional)
	if config.ProjectDescription != "" {
		if err := ValidateProjectDescription(config.ProjectDescription); err != nil {
			result.AddError(&ValidationError{
				Code:       errors.ErrInvalidProjectDescription,
				Field:      "project_description",
				Message:    err.Error(),
				Level:      ErrorLevelLow,
				Suggestion: "Provide a clear, concise description of your project",
			})
		}
	}

	return nil
}

// validateTypes validates enum types
func validateTypes(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Project type validation
	if !config.ProjectType.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidProject,
			Field:      "project_type",
			Message:    fmt.Sprintf("Invalid project type: %s", config.ProjectType),
			Level:      ErrorLevelHigh,
			Suggestion: "Choose a valid project type from the supported options",
		})
	}

	// Git provider validation
	if !config.GitProvider.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "git_provider",
			Message:    fmt.Sprintf("Invalid Git provider: %s", config.GitProvider),
			Level:      ErrorLevelHigh,
			Suggestion: "Choose a supported Git provider (GitHub, GitLab, etc.)",
		})
	}

	// CGO status validation
	if !config.CGOStatus.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "cgo_status",
			Message:    fmt.Sprintf("Invalid CGO status: %s", config.CGOStatus),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid CGO status (disabled, enabled, required)",
		})
	}

	// Docker support validation
	if !config.DockerSupport.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "docker_support",
			Message:    fmt.Sprintf("Invalid Docker support: %s", config.DockerSupport),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid Docker support level (none, build, deploy, both)",
		})
	}

	// Docker registry validation
	if !config.DockerRegistry.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "docker_registry",
			Message:    fmt.Sprintf("Invalid Docker registry: %s", config.DockerRegistry),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid Docker registry (dockerhub, github, gitlab, etc.)",
		})
	}

	// Signing level validation
	if !config.SigningLevel.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "signing_level",
			Message:    fmt.Sprintf("Invalid signing level: %s", config.SigningLevel),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid signing level (none, basic, standard, advanced, enterprise)",
		})
	}

	// Action level validation
	if !config.ActionLevel.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "action_level",
			Message:    fmt.Sprintf("Invalid action level: %s", config.ActionLevel),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid action level (none, basic, standard, advanced, enterprise)",
		})
	}

	// Feature level validation
	if !config.FeatureLevel.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "feature_level",
			Message:    fmt.Sprintf("Invalid feature level: %s", config.FeatureLevel),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid feature level (none, basic, standard, advanced, enterprise)",
		})
	}

	// Config state validation
	if !config.State.IsValid() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "state",
			Message:    fmt.Sprintf("Invalid config state: %s", config.State),
			Level:      ErrorLevelMedium,
			Suggestion: "Choose a valid config state (draft, processing, generated, etc.)",
		})
	}

	return nil
}

// validatePlatformArchCompatibility validates platform-architecture combinations
func validatePlatformArchCompatibility(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Validate platforms
	for i, platform := range config.Platforms {
		if !platform.IsValid() {
			result.AddError(&ValidationError{
				Code:       errors.ErrInvalidConfig,
				Field:      "platforms",
				Message:    fmt.Sprintf("Invalid platform at index %d: %s", i, platform),
				Level:      ErrorLevelHigh,
				Suggestion: "Choose from supported platforms (linux, darwin, windows, etc.)",
			})
		}
	}

	// Validate architectures
	for i, arch := range config.Architectures {
		if !arch.IsValid() {
			result.AddError(&ValidationError{
				Code:       errors.ErrInvalidConfig,
				Field:      "architectures",
				Message:    fmt.Sprintf("Invalid architecture at index %d: %s", i, arch),
				Level:      ErrorLevelHigh,
				Suggestion: "Choose from supported architectures (amd64, arm64, 386, etc.)",
			})
		}
	}

	// Validate compatibility
	for _, platform := range config.Platforms {
		for _, arch := range config.Architectures {
			if !arch.IsCompatibleWith(platform) {
				result.AddError(&ValidationError{
					Code:       errors.ErrInvalidConfig,
					Field:      "platform_arch_compatibility",
					Message:    fmt.Sprintf("Incompatible platform/architecture: %s + %s", platform.String(), arch.String()),
					Level:      ErrorLevelHigh,
					Suggestion: "Choose compatible platform-architecture combinations",
				})
			}
		}
	}

	// Check for known incompatible combinations
	incompatibleCombos := []struct {
		platform domain.Platform
		arch     domain.Architecture
	}{
		{domain.PlatformDarwin, domain.Architecture386},
		{domain.PlatformDarwin, domain.ArchitectureARM},
		{domain.PlatformWindows, domain.ArchitectureARM64},
	}

	for _, combo := range incompatibleCombos {
		platformExists := containsPlatform(config.Platforms, combo.platform)
		archExists := containsArch(config.Architectures, combo.arch)

		if platformExists && archExists {
			result.AddWarning(&ValidationWarning{
				Code:       "PLATFORM_ARCH_WARNING",
				Field:      "platform_arch_compatibility",
				Message:    fmt.Sprintf("Potentially unsupported combination: %s + %s", combo.platform.String(), combo.arch.String()),
				Level:      WarningLevelMedium,
				Suggestion: "Test this combination thoroughly or consider using supported alternatives",
			})
		}
	}

	return nil
}

// validateBusinessRules validates business logic rules
func validateBusinessRules(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Docker configuration validation
	if err := validateDockerBusinessRules(config, result); err != nil {
		return err
	}

	// Signing configuration validation
	if err := validateSigningBusinessRules(config, result); err != nil {
		return err
	}

	// Actions configuration validation
	if err := validateActionsBusinessRules(config, result); err != nil {
		return err
	}

	// Project type-specific validation
	if err := validateProjectTypeBusinessRules(config, result); err != nil {
		return err
	}

	return nil
}

// validateDockerBusinessRules validates Docker business rules
func validateDockerBusinessRules(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Docker support must be compatible with project type
	if config.DockerSupport.IsEnabled() && !config.ProjectType.DockerSupported() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "docker_support",
			Message:    fmt.Sprintf("Docker not supported for project type: %s", config.ProjectType.String()),
			Level:      ErrorLevelHigh,
			Suggestion: "Disable Docker support or choose a Docker-compatible project type",
		})
	}

	// Docker registry must be compatible with Git provider
	if config.DockerSupport.IsDeployEnabled() {
		if config.GitProvider != "" && config.DockerRegistry != "" {
			compatibleRegistry := config.GitProvider.DefaultRegistry()
			if config.DockerRegistry != compatibleRegistry {
				result.AddWarning(&ValidationWarning{
					Code:  "DOCKER_GIT_COMPATIBILITY",
					Field: "docker_registry",
					Message: fmt.Sprintf("Docker registry (%s) may not be optimal for Git provider (%s)",
						config.DockerRegistry.String(), config.GitProvider.String()),
					Level: WarningLevelMedium,
					Suggestion: fmt.Sprintf("Consider using %s registry for %s",
						compatibleRegistry.String(), config.GitProvider.String()),
				})
			}
		}

		// Docker image name is required for deployment
		if config.DockerImage == "" {
			result.AddError(&ValidationError{
				Code:       errors.ErrInvalidConfig,
				Field:      "docker_image",
				Message:    "Docker image name is required when deployment is enabled",
				Level:      ErrorLevelHigh,
				Suggestion: "Specify a Docker image name or disable deployment",
			})
		} else {
			// Validate Docker image name
			if err := ValidateDockerImageName(config.DockerImage); err != nil {
				result.AddError(&ValidationError{
					Code:       errors.ErrInvalidDockerImage,
					Field:      "docker_image",
					Message:    err.Error(),
					Level:      ErrorLevelHigh,
					Suggestion: "Use a valid Docker image name (lowercase, alphanumeric with separators)",
				})
			}
		}
	}

	return nil
}

// validateSigningBusinessRules validates signing business rules
func validateSigningBusinessRules(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Signing tools availability
	if config.SigningLevel.IsEnabled() {
		requiredTools := config.SigningLevel.GetRequiredTools()
		for _, tool := range requiredTools {
			if _, err := exec.LookPath(tool); err != nil {
				result.AddError(&ValidationError{
					Code:       errors.ErrDependencyMissing,
					Field:      "signing_tools",
					Message:    fmt.Sprintf("Required signing tool not found: %s", tool),
					Level:      ErrorLevelHigh,
					Suggestion: fmt.Sprintf("Install %s or lower signing level", tool),
				})
			}
		}
	}

	// Enterprise signing requires project to be ready
	if config.SigningLevel == domain.SigningLevelEnterprise {
		if config.FeatureLevel != domain.FeatureLevelEnterprise {
			result.AddWarning(&ValidationWarning{
				Code:       "SIGNING_FEATURE_LEVEL",
				Field:      "signing_level",
				Message:    "Enterprise signing requires enterprise feature level for optimal security",
				Level:      WarningLevelHigh,
				Suggestion: "Upgrade to enterprise feature level or lower signing level",
			})
		}
	}

	return nil
}

// validateActionsBusinessRules validates Actions business rules
func validateActionsBusinessRules(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Actions must be compatible with Git provider
	if config.ActionLevel.IsEnabled() && !config.GitProvider.ActionsSupported() {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "action_level",
			Message:    fmt.Sprintf("Actions not supported by Git provider: %s", config.GitProvider.String()),
			Level:      ErrorLevelHigh,
			Suggestion: "Choose an Actions-compatible Git provider or disable Actions",
		})
	}

	// Action triggers are required when Actions are enabled
	if config.ActionLevel.IsEnabled() && len(config.ActionsOn) == 0 {
		result.AddError(&ValidationError{
			Code:       errors.ErrInvalidConfig,
			Field:      "actions_on",
			Message:    "Action triggers are required when Actions are enabled",
			Level:      ErrorLevelHigh,
			Suggestion: "Specify at least one action trigger or disable Actions",
		})
	}

	// Validate action triggers
	for i, trigger := range config.ActionsOn {
		if !trigger.IsValid() {
			result.AddError(&ValidationError{
				Code:       errors.ErrInvalidConfig,
				Field:      "actions_on",
				Message:    fmt.Sprintf("Invalid action trigger at index %d: %s", i, trigger),
				Level:      ErrorLevelMedium,
				Suggestion: "Choose from supported action triggers",
			})
		}
	}

	return nil
}

// validateProjectTypeBusinessRules validates project type-specific business rules
func validateProjectTypeBusinessRules(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// CGO requirements based on project type
	if config.ProjectType.DefaultCGOEnabled() && config.CGOStatus.IsDisabled() {
		result.AddWarning(&ValidationWarning{
			Code:       "CGO_PROJECT_TYPE",
			Field:      "cgo_status",
			Message:    fmt.Sprintf("Project type %s typically requires CGO", config.ProjectType.String()),
			Level:      WarningLevelMedium,
			Suggestion: "Consider enabling CGO or verify your project works without CGO",
		})
	}

	// Platform recommendations based on project type
	recommendedPlatforms := config.ProjectType.RecommendedPlatforms()
	for _, platform := range config.Platforms {
		if !containsPlatform(recommendedPlatforms, platform) {
			result.AddWarning(&ValidationWarning{
				Code:  "PLATFORM_RECOMMENDATION",
				Field: "platforms",
				Message: fmt.Sprintf("Platform %s may not be optimal for project type %s",
					platform.String(), config.ProjectType.String()),
				Level:      WarningLevelLow,
				Suggestion: "Consider using recommended platforms for your project type",
			})
		}
	}

	return nil
}

// validateSecurity validates security-related configurations
func validateSecurity(config *domain.SafeProjectConfig, result *ValidationResult) error {
	// Security level consistency
	if config.FeatureLevel.IncludesAdvanced() && config.SigningLevel == domain.SigningLevelNone {
		result.AddWarning(&ValidationWarning{
			Code:       "SECURITY_SIGNING",
			Field:      "signing_level",
			Message:    "Advanced features should include code signing for security",
			Level:      WarningLevelHigh,
			Suggestion: "Enable code signing or lower feature level",
		})
	}

	// Docker security
	if config.DockerSupport.IsEnabled() && !config.DockerRegistry.RequiresAuth() {
		result.AddWarning(&ValidationWarning{
			Code:       "DOCKER_SECURITY",
			Field:      "docker_registry",
			Message:    "Consider using a registry that requires authentication for better security",
			Level:      WarningLevelMedium,
			Suggestion: "Use GitHub Container Registry, GitLab Registry, or other authenticated registries",
		})
	}

	// Actions security
	if config.ActionLevel.IsEnabled() {
		requiredPerms := config.ActionLevel.GetRequiredPermissions()
		if len(requiredPerms) > 2 {
			result.AddWarning(&ValidationWarning{
				Code:       "ACTIONS_PERMISSIONS",
				Field:      "action_level",
				Message:    "High permission levels increase security risks in GitHub Actions",
				Level:      WarningLevelMedium,
				Suggestion: "Review and minimize required permissions for security",
			})
		}
	}

	return nil
}

// generateWarnings generates helpful warnings based on configuration
func generateWarnings(config *domain.SafeProjectConfig, result *ValidationResult) {
	// Performance warnings
	if len(config.Platforms)*len(config.Architectures) > 8 {
		result.AddWarning(&ValidationWarning{
			Code:  "BUILD_PERFORMANCE",
			Field: "platform_arch_count",
			Message: fmt.Sprintf("Large build matrix (%d combos) may significantly increase build time",
				len(config.Platforms)*len(config.Architectures)),
			Level:      WarningLevelLow,
			Suggestion: "Consider reducing platform/architecture combinations for faster builds",
		})
	}

	// Best practice warnings
	if config.ProjectDescription == "" {
		result.AddWarning(&ValidationWarning{
			Code:       "PROJECT_DESCRIPTION",
			Field:      "project_description",
			Message:    "Project description is recommended for better documentation",
			Level:      WarningLevelLow,
			Suggestion: "Add a clear project description to improve documentation",
		})
	}

	// Feature utilization warnings
	if config.FeatureLevel == domain.FeatureLevelNone {
		result.AddWarning(&ValidationWarning{
			Code:       "FEATURE_UTILIZATION",
			Field:      "feature_level",
			Message:    "No advanced features enabled - consider enabling for better functionality",
			Level:      WarningLevelLow,
			Suggestion: "Enable basic features like Docker support or Actions if appropriate",
		})
	}

	// Maintenance warnings
	if config.SigningLevel == domain.SigningLevelNone && config.DockerSupport.IsEnabled() {
		result.AddWarning(&ValidationWarning{
			Code:       "SECURITY_MAINTENANCE",
			Field:      "signing_level",
			Message:    "Docker images should be signed for security and trust",
			Level:      WarningLevelMedium,
			Suggestion: "Enable code signing for Docker images",
		})
	}
}

// Helper functions
func containsPlatform(platforms []domain.Platform, target domain.Platform) bool {
	for _, p := range platforms {
		if p == target {
			return true
		}
	}
	return false
}

func containsArch(archs []domain.Architecture, target domain.Architecture) bool {
	for _, a := range archs {
		if a == target {
			return true
		}
	}
	return false
}
