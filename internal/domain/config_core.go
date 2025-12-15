package config

import (
	"encoding/json"
	"exec"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/validation"
)

// SafeProjectConfig represents single source of truth for project configuration
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type SafeProjectConfig struct {
	// Basic Information
	ProjectName        string             `json:"project_name" yaml:"project_name"`
	ProjectDescription string             `json:"project_description,omitempty" yaml:"project_description,omitempty"`
	ProjectType        domain.ProjectType `json:"project_type" yaml:"project_type"`
	BinaryName         string             `json:"binary_name" yaml:"binary_name"`
	MainPath           string             `json:"main_path" yaml:"main_path"`

	// Build Configuration
	Platforms     []domain.Platform     `json:"platforms" yaml:"platforms"`
	Architectures []domain.Architecture `json:"architectures" yaml:"architectures"`
	CGOStatus     domain.CGOStatus      `json:"cgo_status" yaml:"cgo_status"`
	BuildTags     []domain.BuildTag     `json:"build_tags,omitempty" yaml:"build_tags,omitempty"`
	LDFlags       bool                  `json:"ldflags" yaml:"ldflags"`

	// Release Configuration
	GitProvider    domain.GitProvider    `json:"git_provider" yaml:"git_provider"`
	DockerSupport  domain.DockerSupport  `json:"docker_support" yaml:"docker_support"`
	DockerRegistry domain.DockerRegistry `json:"docker_registry" yaml:"docker_registry"`
	DockerImage    string                `json:"docker_image,omitempty" yaml:"docker_image,omitempty"`
	SigningLevel   domain.SigningLevel   `json:"signing_level" yaml:"signing_level"`
	Homebrew       bool                  `json:"homebrew" yaml:"homebrew"`
	Snap           bool                  `json:"snap" yaml:"snap"`
	SBOM           bool                  `json:"sbom" yaml:"sbom"`

	// CI/CD Configuration
	ActionLevel domain.ActionLevel     `json:"action_level" yaml:"action_level"`
	ActionsOn   []domain.ActionTrigger `json:"actions_on" yaml:"actions_on"`

	// Advanced Features
	FeatureLevel domain.FeatureLevel `json:"feature_level" yaml:"feature_level"`

	// State Management
	State domain.ConfigState `json:"state" yaml:"state"`
}

// NewSafeProjectConfig creates a new safe configuration with smart defaults
func NewSafeProjectConfig() *SafeProjectConfig {
	return &SafeProjectConfig{
		// Smart defaults based on project analysis
		ProjectType:    GetRecommendedProjectType(),
		Platforms:      GetRecommendedPlatforms(),
		Architectures:  GetRecommendedArchitectures(),
		GitProvider:    GetRecommendedGitProvider(),
		DockerRegistry: GetRecommendedDockerRegistry(),
		CGOStatus:      domain.CGOStatusDisabled,
		ActionLevel:    domain.ActionLevelBasic,
		SigningLevel:   domain.SigningLevelNone,
		FeatureLevel:   domain.FeatureLevelBasic,
		State:          domain.ConfigStateDraft,
		LDFlags:        true,
		Homebrew:       false,
		Snap:           false,
		SBOM:           false,
	}
}

// Validate validates the configuration
func (spc *SafeProjectConfig) Validate() error {
	// Validate basic fields
	if err := validation.ValidateProjectName(spc.ProjectName); err != nil {
		return err
	}

	if err := validation.ValidateBinaryName(spc.BinaryName); err != nil {
		return err
	}

	if err := validation.ValidateMainPath(spc.MainPath); err != nil {
		return err
	}

	if spc.ProjectDescription != "" {
		if err := validation.ValidateProjectDescription(spc.ProjectDescription); err != nil {
			return err
		}
	}

	// Validate enums
	if !spc.ProjectType.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidProject,
			"Invalid project type",
			string(spc.ProjectType),
		).WithField("project_type")
	}

	if !spc.GitProvider.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid Git provider",
			string(spc.GitProvider),
		).WithField("git_provider")
	}

	if !spc.CGOStatus.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid CGO status",
			string(spc.CGOStatus),
		).WithField("cgo_status")
	}

	if !spc.DockerSupport.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid Docker support",
			string(spc.DockerSupport),
		).WithField("docker_support")
	}

	if !spc.DockerRegistry.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid Docker registry",
			string(spc.DockerRegistry),
		).WithField("docker_registry")
	}

	if !spc.SigningLevel.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid signing level",
			string(spc.SigningLevel),
		).WithField("signing_level")
	}

	if !spc.ActionLevel.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid action level",
			string(spc.ActionLevel),
		).WithField("action_level")
	}

	if !spc.FeatureLevel.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid feature level",
			string(spc.FeatureLevel),
		).WithField("feature_level")
	}

	if !spc.State.IsValid() {
		return domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Invalid config state",
			string(spc.State),
		).WithField("state")
	}

	// Validate platforms
	for _, platform := range spc.Platforms {
		if !platform.IsValid() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Invalid platform",
				string(platform),
			).WithField("platforms")
		}
	}

	// Validate architectures
	for _, arch := range spc.Architectures {
		if !arch.IsValid() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Invalid architecture",
				string(arch),
			).WithField("architectures")
		}
	}

	// Validate build tags
	for _, tag := range spc.BuildTags {
		if !tag.IsValid() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Invalid build tag",
				string(tag),
			).WithField("build_tags")
		}
	}

	// Validate action triggers
	for _, trigger := range spc.ActionsOn {
		if !trigger.IsValid() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Invalid action trigger",
				string(trigger),
			).WithField("actions_on")
		}
	}

	// Validate business rules
	return spc.validateBusinessRules()
}

// validateBusinessRules validates business logic rules
func (spc *SafeProjectConfig) validateBusinessRules() error {
	// Validate platform/architecture compatibility
	if err := spc.validatePlatformArchCompatibility(); err != nil {
		return err
	}

	// Validate Docker configuration
	if err := spc.validateDockerConfiguration(); err != nil {
		return err
	}

	// Validate signing configuration
	if err := spc.validateSigningConfiguration(); err != nil {
		return err
	}

	// Validate Actions configuration
	if err := spc.validateActionsConfiguration(); err != nil {
		return err
	}

	return nil
}

// validatePlatformArchCompatibility validates platform/architecture compatibility
func (spc *SafeProjectConfig) validatePlatformArchCompatibility() error {
	for _, platform := range spc.Platforms {
		for _, arch := range spc.Architectures {
			if !arch.IsCompatibleWith(platform) {
				return domain.NewValidationError(
					domain.ErrInvalidConfig,
					"Platform/architecture combination not supported",
					fmt.Sprintf("%s + %s", platform.String(), arch.String()),
				).WithField("platform_arch_compatibility")
			}
		}
	}

	return nil
}

// validateDockerConfiguration validates Docker configuration
func (spc *SafeProjectConfig) validateDockerConfiguration() error {
	if spc.DockerSupport.IsEnabled() {
		// Check if Docker is supported by project type
		if !spc.ProjectType.DockerSupported() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Docker not supported for project type",
				spc.ProjectType.String(),
			).WithField("docker_support")
		}

		// Validate Docker image name
		if spc.DockerSupport.IsDeployEnabled() && spc.GetDockerImageName() == "" {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Docker image name is required for deployment",
				"Docker image name cannot be empty when deployment is enabled",
			).WithField("docker_image")
		}
	}

	return nil
}

// validateSigningConfiguration validates signing configuration
func (spc *SafeProjectConfig) validateSigningConfiguration() error {
	if spc.SigningLevel.IsEnabled() {
		// Check if required tools are available
		requiredTools := spc.SigningLevel.GetRequiredTools()
		for _, tool := range requiredTools {
			if _, err := exec.LookPath(tool); err != nil {
				return domain.NewValidationError(
					domain.ErrDependencyMissing,
					"Required signing tool not found",
					tool,
				).WithField("signing_tools")
			}
		}
	}

	return nil
}

// validateActionsConfiguration validates Actions configuration
func (spc *SafeProjectConfig) validateActionsConfiguration() error {
	if spc.ActionLevel.IsEnabled() {
		// Check if Git provider supports actions
		if !spc.GitProvider.ActionsSupported() {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Git provider doesn't support actions",
				spc.GitProvider.String(),
			).WithField("git_provider")
		}

		// Validate action triggers
		if len(spc.ActionsOn) == 0 {
			return domain.NewValidationError(
				domain.ErrInvalidConfig,
				"Action triggers are required when actions are enabled",
				"At least one action trigger must be specified",
			).WithField("actions_on")
		}
	}

	return nil
}

// Clone creates a deep copy of the configuration
func (spc *SafeProjectConfig) Clone() *SafeProjectConfig {
	clone := *spc

	// Deep copy slices
	if spc.Platforms != nil {
		clone.Platforms = make([]domain.Platform, len(spc.Platforms))
		copy(clone.Platforms, spc.Platforms)
	}

	if spc.Architectures != nil {
		clone.Architectures = make([]domain.Architecture, len(spc.Architectures))
		copy(clone.Architectures, spc.Architectures)
	}

	if spc.BuildTags != nil {
		clone.BuildTags = make([]domain.BuildTag, len(spc.BuildTags))
		copy(clone.BuildTags, spc.BuildTags)
	}

	if spc.ActionsOn != nil {
		clone.ActionsOn = make([]domain.ActionTrigger, len(spc.ActionsOn))
		copy(clone.ActionsOn, spc.ActionsOn)
	}

	return &clone
}

// ToJSON converts configuration to JSON string
func (spc *SafeProjectConfig) ToJSON() (string, error) {
	data, err := json.MarshalIndent(spc, "", "  ")
	if err != nil {
		return "", domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Failed to serialize configuration to JSON",
			err.Error(),
		)
	}
	return string(data), nil
}

// ToYAML converts configuration to YAML string
func (spc *SafeProjectConfig) ToYAML() (string, error) {
	data, err := yaml.Marshal(spc)
	if err != nil {
		return "", domain.NewValidationError(
			domain.ErrInvalidConfig,
			"Failed to serialize configuration to YAML",
			err.Error(),
		)
	}
	return string(data), nil
}

// FromJSON loads configuration from JSON string
func (spc *SafeProjectConfig) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), spc)
}

// FromYAML loads configuration from YAML string
func (spc *SafeProjectConfig) FromYAML(yamlStr string) error {
	return yaml.Unmarshal([]byte(yamlStr), spc)
}

// GetSummary returns a summary of the configuration
func (spc *SafeProjectConfig) GetSummary() ConfigSummary {
	return ConfigSummary{
		ProjectName:    spc.ProjectName,
		ProjectType:    spc.ProjectType.String(),
		BinaryName:     spc.BinaryName,
		PlatformCount:  len(spc.Platforms),
		ArchCount:      len(spc.Architectures),
		CGOEnabled:     spc.CGOStatus.IsEnabled(),
		DockerEnabled:  spc.DockerSupport.IsEnabled(),
		SigningEnabled: spc.SigningLevel.IsEnabled(),
		ActionsEnabled: spc.ActionLevel.IsEnabled(),
		State:          spc.State.String(),
		FeatureLevel:   spc.FeatureLevel.String(),
	}
}

// ConfigSummary represents a summary of configuration
type ConfigSummary struct {
	ProjectName    string `json:"project_name"`
	ProjectType    string `json:"project_type"`
	BinaryName     string `json:"binary_name"`
	PlatformCount  int    `json:"platform_count"`
	ArchCount      int    `json:"arch_count"`
	CGOEnabled     bool   `json:"cgo_enabled"`
	DockerEnabled  bool   `json:"docker_enabled"`
	SigningEnabled bool   `json:"signing_enabled"`
	ActionsEnabled bool   `json:"actions_enabled"`
	State          string `json:"state"`
	FeatureLevel   string `json:"feature_level"`
}
