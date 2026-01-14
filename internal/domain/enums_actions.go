package domain

import "fmt"

// ActionLevel represents GitHub Actions complexity levels
// This enum replaces string-based action levels for type safety.
type ActionLevel string

const (
	ActionLevelNone       ActionLevel = "none"
	ActionLevelBasic      ActionLevel = "basic"
	ActionLevelStandard   ActionLevel = "standard"
	ActionLevelAdvanced   ActionLevel = "advanced"
	ActionLevelEnterprise ActionLevel = "enterprise"
)

// IsValid returns true if ActionLevel is valid.
func (al ActionLevel) IsValid() bool {
	switch al {
	case ActionLevelNone, ActionLevelBasic, ActionLevelStandard,
		ActionLevelAdvanced, ActionLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (al ActionLevel) String() string {
	switch al {
	case ActionLevelNone:
		return "None"
	case ActionLevelBasic:
		return "Basic"
	case ActionLevelStandard:
		return "Standard"
	case ActionLevelAdvanced:
		return "Advanced"
	case ActionLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if actions are enabled.
func (al ActionLevel) IsEnabled() bool {
	return al != ActionLevelNone
}

// IsProductionReady returns true if actions are production-ready.
func (al ActionLevel) IsProductionReady() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesSecurity returns true if level includes security features.
func (al ActionLevel) IncludesSecurity() bool {
	return al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesMonitoring returns true if level includes monitoring features.
func (al ActionLevel) IncludesMonitoring() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesTesting returns true if level includes testing features.
func (al ActionLevel) IncludesTesting() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// GetRecommendedTriggers returns recommended triggers for level.
func (al ActionLevel) GetRecommendedTriggers() []ActionTrigger {
	switch al {
	case ActionLevelNone:
		return []ActionTrigger{}
	case ActionLevelBasic:
		return []ActionTrigger{ActionTriggerManual, ActionTriggerVersionTags}
	case ActionLevelStandard:
		return []ActionTrigger{ActionTriggerManual, ActionTriggerVersionTags, ActionTriggerMainPush}
	case ActionLevelAdvanced:
		return []ActionTrigger{ActionTriggerManual, ActionTriggerVersionTags, ActionTriggerMainPush, ActionTriggerPullRequest}
	case ActionLevelEnterprise:
		return []ActionTrigger{ActionTriggerManual, ActionTriggerVersionTags, ActionTriggerMainPush, ActionTriggerPullRequest, ActionTriggerSchedule}
	default:
		return []ActionTrigger{}
	}
}

// GetRequiredPermissions returns required GitHub permissions.
func (al ActionLevel) GetRequiredPermissions() []string {
	switch al {
	case ActionLevelNone:
		return []string{}
	case ActionLevelBasic:
		return []string{"contents: write"}
	case ActionLevelStandard:
		return []string{"contents: write", "packages: write"}
	case ActionLevelAdvanced:
		return []string{"contents: write", "packages: write", "id-token: write"}
	case ActionLevelEnterprise:
		return []string{"contents: write", "packages: write", "id-token: write", "issues: write", "pull-requests: write"}
	default:
		return []string{}
	}
}

// GetEnvironmentCount returns recommended number of environments.
func (al ActionLevel) GetEnvironmentCount() int {
	switch al {
	case ActionLevelNone:
		return 0
	case ActionLevelBasic:
		return 1
	case ActionLevelStandard:
		return 2
	case ActionLevelAdvanced:
		return 3
	case ActionLevelEnterprise:
		return 4
	default:
		return 0
	}
}

// DockerSupport represents Docker build and deployment options
// This enum replaces boolean docker flags for type safety.
type DockerSupport string

const (
	DockerSupportNone   DockerSupport = "none"
	DockerSupportBuild  DockerSupport = "build"
	DockerSupportDeploy DockerSupport = "deploy"
	DockerSupportBoth   DockerSupport = "both"
)

// IsValid returns true if DockerSupport is valid.
func (ds DockerSupport) IsValid() bool {
	switch ds {
	case DockerSupportNone, DockerSupportBuild, DockerSupportDeploy, DockerSupportBoth:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (ds DockerSupport) String() string {
	switch ds {
	case DockerSupportNone:
		return "None"
	case DockerSupportBuild:
		return "Build Only"
	case DockerSupportDeploy:
		return "Deploy Only"
	case DockerSupportBoth:
		return "Build & Deploy"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if Docker is enabled.
func (ds DockerSupport) IsEnabled() bool {
	return ds != DockerSupportNone
}

// IsBuildEnabled returns true if Docker build is enabled.
func (ds DockerSupport) IsBuildEnabled() bool {
	return ds == DockerSupportBuild || ds == DockerSupportBoth
}

// IsDeployEnabled returns true if Docker deployment is enabled.
func (ds DockerSupport) IsDeployEnabled() bool {
	return ds == DockerSupportDeploy || ds == DockerSupportBoth
}

// ShouldBuild returns true if Docker images should be built.
func (ds DockerSupport) ShouldBuild() bool {
	return ds.IsBuildEnabled()
}

// ShouldPublish returns true if Docker images should be published.
func (ds DockerSupport) ShouldPublish() bool {
	return ds.IsDeployEnabled()
}

// ToBool converts to legacy boolean for compatibility.
func (ds DockerSupport) ToBool() bool {
	return ds.IsEnabled()
}

// ValidateDockerSupport validates a Docker support level.
func ValidateDockerSupport(support DockerSupport) error {
	if !support.IsValid() {
		return fmt.Errorf("invalid Docker support level: %s", support)
	}
	return nil
}
