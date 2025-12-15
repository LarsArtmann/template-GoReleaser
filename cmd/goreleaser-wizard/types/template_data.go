package types

import (
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// GoReleaserTemplateData represents strongly typed template data for GoReleaser configuration
// Eliminates map[string]any usage for type safety
type GoReleaserTemplateData struct {
	// Project Information
	ProjectName string `json:"project_name"`
	BinaryName  string `json:"binary_name"`
	MainPath    string `json:"main_path"`

	// Version Information
	Version    string `json:"version"`
	Tag        string `json:"tag"`
	Major      string `json:"major"`
	Date       string `json:"date"`
	FullCommit string `json:"full_commit"`

	// Build Configuration
	CGOEnabled         string            `json:"cgo_enabled"`
	Platforms          []string          `json:"platforms"`
	Architectures      []string          `json:"architectures"`
	BuildTags          []string          `json:"build_tags,omitempty"`
	IgnoreCombinations []PlatformIgnored `json:"ignore_combinations,omitempty"`

	// Release Configuration
	DockerEnabled  bool   `json:"docker_enabled"`
	SigningEnabled bool   `json:"signing_enabled"`
	DockerRegistry string `json:"docker_registry,omitempty"`
	DockerImage    string `json:"docker_image,omitempty"`

	// Environment Variables
	Env map[string]string `json:"env"`
}

// GitHubActionsTemplateData represents strongly typed template data for GitHub Actions
type GitHubActionsTemplateData struct {
	ProjectName    string   `json:"project_name"`
	Triggers       []string `json:"triggers"`
	DockerEnabled  bool     `json:"docker_enabled"`
	SigningEnabled bool     `json:"signing_enabled"`
	DockerRegistry string   `json:"docker_registry,omitempty"`
	DockerImage    string   `json:"docker_image,omitempty"`
}

// PlatformIgnored represents platform/architecture combinations to ignore
type PlatformIgnored struct {
	GoOS   string `json:"goos"`
	GoArch string `json:"goarch"`
}

// DockerConfig represents Docker-specific configuration
type DockerConfig struct {
	Enabled  bool   `json:"enabled"`
	Registry string `json:"registry"`
	Image    string `json:"image"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// SigningConfig represents code signing configuration
type SigningConfig struct {
	Enabled     bool   `json:"enabled"`
	Level       string `json:"level"`
	KeyID       string `json:"key_id,omitempty"`
	Certificate string `json:"certificate,omitempty"`
}

// ValidationResult represents structured validation results
type ValidationResult struct {
	IsValid  bool                 `json:"is_valid"`
	Errors   []*ValidationError   `json:"errors"`
	Warnings []*ValidationWarning `json:"warnings"`
	Summary  ValidationSummary    `json:"summary"`
}

// ValidationError represents a structured validation error
type ValidationError struct {
	Code       string `json:"code"`
	Field      string `json:"field"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	Context    string `json:"context,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ValidationWarning represents a structured validation warning
type ValidationWarning struct {
	Code       string `json:"code"`
	Field      string `json:"field"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	Context    string `json:"context,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalErrors   int `json:"total_errors"`
	TotalWarnings int `json:"total_warnings"`
	Critical      int `json:"critical"`
	High          int `json:"high"`
	Medium        int `json:"medium"`
	Low           int `json:"low"`
}

// JobExecutionResult represents the result of a job execution
type JobExecutionResult struct {
	JobID    string      `json:"job_id"`
	JobName  string      `json:"job_name"`
	Status   JobStatus   `json:"status"`
	Duration string      `json:"duration"`
	Error    *JobError   `json:"error,omitempty"`
	Output   string      `json:"output,omitempty"`
	Metadata JobMetadata `json:"metadata"`
}

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// JobError represents a job execution error
type JobError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Context   string `json:"context,omitempty"`
	Retryable bool   `json:"retryable"`
}

// JobMetadata represents job execution metadata
type JobMetadata struct {
	StartedAt   string            `json:"started_at"`
	CompletedAt string            `json:"completed_at,omitempty"`
	Retries     int               `json:"retries"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// WorkflowState represents the state of a workflow
type WorkflowState struct {
	WorkflowID   string                `json:"workflow_id"`
	WorkflowName string                `json:"workflow_name"`
	Status       WorkflowStatus        `json:"status"`
	CurrentStep  string                `json:"current_step"`
	TotalSteps   int                   `json:"total_steps"`
	Progress     float64               `json:"progress"`
	StartedAt    string                `json:"started_at"`
	Results      []*JobExecutionResult `json:"results"`
	Metadata     WorkflowMetadata      `json:"metadata"`
}

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// WorkflowMetadata represents workflow metadata
type WorkflowMetadata struct {
	CreatedBy   string            `json:"created_by"`
	Timeout     string            `json:"timeout"`
	Parallel    bool              `json:"parallel"`
	Environment string            `json:"environment"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// NewGoReleaserTemplateData creates typed template data from SafeProjectConfig
func NewGoReleaserTemplateData(config *domain.SafeProjectConfig) *GoReleaserTemplateData {
	data := &GoReleaserTemplateData{
		ProjectName: config.ProjectName,
		BinaryName:  config.BinaryName,
		MainPath:    config.MainPath,
		Env: map[string]string{
			"GITHUB_OWNER": getGitHubOwner(),
			"GITHUB_REPO":  getGitHubRepo(),
		},
		IgnoreCombinations: []PlatformIgnored{
			{GoOS: "darwin", GoArch: "386"},
			{GoOS: "windows", GoArch: "arm64"},
		},
	}

	// Convert platforms
	if len(config.Platforms) > 0 {
		platforms := make([]string, len(config.Platforms))
		for i, platform := range config.Platforms {
			platforms[i] = string(platform)
		}
		data.Platforms = platforms
	}

	// Convert architectures
	if len(config.Architectures) > 0 {
		architectures := make([]string, len(config.Architectures))
		for i, arch := range config.Architectures {
			architectures[i] = string(arch)
		}
		data.Architectures = architectures
	}

	// Convert build tags
	if len(config.BuildTags) > 0 {
		tags := make([]string, len(config.BuildTags))
		for i, tag := range config.BuildTags {
			tags[i] = tag.String()
		}
		data.BuildTags = tags
	}

	// Set Docker configuration
	if config.DockerSupport.IsEnabled() {
		data.DockerEnabled = true
		data.DockerRegistry = config.DockerRegistry.String()
		data.DockerImage = config.GetDockerImageName()
	}

	// Set signing configuration
	data.SigningEnabled = config.SigningLevel.IsEnabled()

	return data
}

// NewGitHubActionsTemplateData creates typed template data from SafeProjectConfig
func NewGitHubActionsTemplateData(config *domain.SafeProjectConfig) *GitHubActionsTemplateData {
	data := &GitHubActionsTemplateData{
		ProjectName:    config.ProjectName,
		DockerEnabled:  config.DockerSupport.IsEnabled(),
		SigningEnabled: config.SigningLevel.IsEnabled(),
	}

	// Convert action triggers
	if len(config.ActionsOn) > 0 {
		triggers := make([]string, len(config.ActionsOn))
		for i, trigger := range config.ActionsOn {
			triggers[i] = trigger.GitHubPattern()
		}
		data.Triggers = triggers
	}

	// Set Docker configuration
	if config.DockerSupport.IsEnabled() {
		data.DockerRegistry = config.DockerRegistry.String()
		data.DockerImage = config.GetDockerImageName()
	}

	return data
}
