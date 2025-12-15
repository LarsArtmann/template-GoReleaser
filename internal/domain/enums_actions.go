package interfaces

// ActionLevel represents GitHub Actions complexity levels
// This enum replaces string-based action levels for type safety
type ActionLevel string

const (
	ActionLevelNone       ActionLevel = "none"
	ActionLevelBasic      ActionLevel = "basic"
	ActionLevelStandard   ActionLevel = "standard"
	ActionLevelAdvanced   ActionLevel = "advanced"
	ActionLevelEnterprise ActionLevel = "enterprise"
)

// IsValid returns true if ActionLevel is valid
func (al ActionLevel) IsValid() bool {
	switch al {
	case ActionLevelNone, ActionLevelBasic, ActionLevelStandard,
		ActionLevelAdvanced, ActionLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
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

// IsEnabled returns true if actions are enabled
func (al ActionLevel) IsEnabled() bool {
	return al != ActionLevelNone
}

// IsProductionReady returns true if actions are production-ready
func (al ActionLevel) IsProductionReady() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesSecurity returns true if level includes security features
func (al ActionLevel) IncludesSecurity() bool {
	return al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesMonitoring returns true if level includes monitoring features
func (al ActionLevel) IncludesMonitoring() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// IncludesTesting returns true if level includes testing features
func (al ActionLevel) IncludesTesting() bool {
	return al == ActionLevelStandard || al == ActionLevelAdvanced || al == ActionLevelEnterprise
}

// GetRecommendedTriggers returns recommended triggers for level
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

// GetRequiredPermissions returns required GitHub permissions
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

// GetEnvironmentCount returns recommended number of environments
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

// ActionTrigger represents GitHub Actions triggers
// This enum replaces string-based triggers for type safety
type ActionTrigger string

const (
	ActionTriggerManual      ActionTrigger = "manual"
	ActionTriggerVersionTags ActionTrigger = "version-tags"
	ActionTriggerAllTags     ActionTrigger = "all-tags"
	ActionTriggerMainPush    ActionTrigger = "main-push"
	ActionTriggerAllPush     ActionTrigger = "all-push"
	ActionTriggerPullRequest ActionTrigger = "pull-request"
	ActionTriggerSchedule    ActionTrigger = "schedule"
	ActionTriggerWebhook     ActionTrigger = "webhook"
	ActionTriggerWorkflow    ActionTrigger = "workflow"
)

// IsValid returns true if ActionTrigger is valid
func (at ActionTrigger) IsValid() bool {
	switch at {
	case ActionTriggerManual, ActionTriggerVersionTags, ActionTriggerAllTags,
		ActionTriggerMainPush, ActionTriggerAllPush, ActionTriggerPullRequest,
		ActionTriggerSchedule, ActionTriggerWebhook, ActionTriggerWorkflow:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (at ActionTrigger) String() string {
	switch at {
	case ActionTriggerManual:
		return "Manual Dispatch"
	case ActionTriggerVersionTags:
		return "Version Tags"
	case ActionTriggerAllTags:
		return "All Tags"
	case ActionTriggerMainPush:
		return "Main Branch Push"
	case ActionTriggerAllPush:
		return "All Branches Push"
	case ActionTriggerPullRequest:
		return "Pull Request"
	case ActionTriggerSchedule:
		return "Scheduled"
	case ActionTriggerWebhook:
		return "Webhook"
	case ActionTriggerWorkflow:
		return "Workflow Dispatch"
	default:
		return "Unknown"
	}
}

// GitHubPattern returns the GitHub Actions syntax for this trigger
func (at ActionTrigger) GitHubPattern() string {
	switch at {
	case ActionTriggerManual:
		return "workflow_dispatch:"
	case ActionTriggerVersionTags:
		return "push:\n  tags:\n    - 'v*'"
	case ActionTriggerAllTags:
		return "push:\n  tags:\n    - '*'"
	case ActionTriggerMainPush:
		return "push:\n  branches:\n    - 'main'"
	case ActionTriggerAllPush:
		return "push:"
	case ActionTriggerPullRequest:
		return "pull_request:"
	case ActionTriggerSchedule:
		return "schedule:\n  - cron: '0 2 * * 0'" // Weekly on Sunday 2 AM
	case ActionTriggerWebhook:
		return "repository_dispatch:"
	case ActionTriggerWorkflow:
		return "workflow_run:"
	default:
		return ""
	}
}

// IsManual returns true if trigger is manual
func (at ActionTrigger) IsManual() bool {
	return at == ActionTriggerManual || at == ActionTriggerWorkflow
}

// IsAutomated returns true if trigger is automated
func (at ActionTrigger) IsAutomated() bool {
	return !at.IsManual()
}

// IsSecuritySensitive returns true if trigger is security-sensitive
func (at ActionTrigger) IsSecuritySensitive() bool {
	switch at {
	case ActionTriggerWebhook, ActionTriggerSchedule:
		return true
	default:
		return false
	}
}

// RequiresApproval returns true if trigger requires approval
func (at ActionTrigger) RequiresApproval() bool {
	switch at {
	case ActionTriggerSchedule, ActionTriggerWebhook:
		return true
	default:
		return false
	}
}

// CanTriggerOnPush returns true if trigger can be triggered by push
func (at ActionTrigger) CanTriggerOnPush() bool {
	switch at {
	case ActionTriggerVersionTags, ActionTriggerAllTags, ActionTriggerMainPush, ActionTriggerAllPush:
		return true
	default:
		return false
	}
}

// CanTriggerOnTag returns true if trigger can be triggered by tag
func (at ActionTrigger) CanTriggerOnTag() bool {
	switch at {
	case ActionTriggerVersionTags, ActionTriggerAllTags:
		return true
	default:
		return false
	}
}

// CanTriggerOnBranch returns true if trigger can be triggered by branch push
func (at ActionTrigger) CanTriggerOnBranch() bool {
	switch at {
	case ActionTriggerMainPush, ActionTriggerAllPush:
		return true
	default:
		return false
	}
}

// CanTriggerOnSchedule returns true if trigger can be scheduled
func (at ActionTrigger) CanTriggerOnSchedule() bool {
	return at == ActionTriggerSchedule
}

// CanTriggerOnWebhook returns true if trigger can be triggered by webhook
func (at ActionTrigger) CanTriggerOnWebhook() bool {
	return at == ActionTriggerWebhook
}

// GetRecommendedBranches returns recommended branches for this trigger
func (at ActionTrigger) GetRecommendedBranches() []string {
	switch at {
	case ActionTriggerMainPush:
		return []string{"main", "master"}
	case ActionTriggerAllPush:
		return []string{"main", "master", "develop", "feature/*", "hotfix/*"}
	case ActionTriggerPullRequest:
		return []string{"main", "master", "develop"}
	default:
		return []string{}
	}
}

// GetRecommendedTags returns recommended tag patterns for this trigger
func (at ActionTrigger) GetRecommendedTags() []string {
	switch at {
	case ActionTriggerVersionTags:
		return []string{"v*", "release/*"}
	case ActionTriggerAllTags:
		return []string{"*"}
	default:
		return []string{}
	}
}

// DockerSupport represents Docker build and deployment options
// This enum replaces boolean docker flags for type safety
type DockerSupport string

const (
	DockerSupportNone   DockerSupport = "none"
	DockerSupportBuild  DockerSupport = "build"
	DockerSupportDeploy DockerSupport = "deploy"
	DockerSupportBoth   DockerSupport = "both"
)

// IsValid returns true if DockerSupport is valid
func (ds DockerSupport) IsValid() bool {
	switch ds {
	case DockerSupportNone, DockerSupportBuild, DockerSupportDeploy, DockerSupportBoth:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
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

// IsEnabled returns true if Docker is enabled
func (ds DockerSupport) IsEnabled() bool {
	return ds != DockerSupportNone
}

// IsBuildEnabled returns true if Docker build is enabled
func (ds DockerSupport) IsBuildEnabled() bool {
	return ds == DockerSupportBuild || ds == DockerSupportBoth
}

// IsDeployEnabled returns true if Docker deployment is enabled
func (ds DockerSupport) IsDeployEnabled() bool {
	return ds == DockerSupportDeploy || ds == DockerSupportBoth
}

// RequiresRegistry returns true if requires container registry
func (ds DockerSupport) RequiresRegistry() bool {
	return ds.IsDeployEnabled()
}

// RequiresDockerfile returns true if requires Dockerfile
func (ds DockerSupport) RequiresDockerfile() bool {
	return ds.IsEnabled()
}

// GetRequiredDockerfile returns type of Dockerfile required
func (ds DockerSupport) GetRequiredDockerfile() string {
	switch ds {
	case DockerSupportNone:
		return ""
	case DockerSupportBuild, DockerSupportDeploy, DockerSupportBoth:
		return "Dockerfile"
	default:
		return ""
	}
}

// GetBuildContext returns recommended build context
func (ds DockerSupport) GetBuildContext() string {
	switch ds {
	case DockerSupportBuild, DockerSupportDeploy, DockerSupportBoth:
		return "."
	default:
		return ""
	}
}

// GetRequiredArgs returns required Docker build arguments
func (ds DockerSupport) GetRequiredArgs() []string {
	var args []string

	if ds.IsBuildEnabled() {
		args = append(args, "--pull", "--label=org.opencontainers.image.created={{.Date}}")
	}

	if ds.IsDeployEnabled() {
		args = append(args, "--platform=linux/amd64")
	}

	return args
}
