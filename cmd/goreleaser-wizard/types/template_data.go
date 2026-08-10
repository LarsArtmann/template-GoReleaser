package types

import (
	"context"
	"os/exec"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/git"
)

// GoReleaserTemplateData represents strongly typed template data for GoReleaser configuration
// Eliminates map[string]any usage for type safety.
type GoReleaserTemplateData struct {
	// Project Information
	ProjectName string `json:"projectName"`
	BinaryName  string `json:"binaryName"`
	MainPath    string `json:"mainPath"`

	// Version Information
	Version    string `json:"version"`
	Tag        string `json:"tag"`
	Major      string `json:"major"`
	Date       string `json:"date"`
	FullCommit string `json:"fullCommit"`

	// Build Configuration
	CGOEnabled         string            `json:"cgoEnabled"`
	Platforms          []string          `json:"platforms"`
	Architectures      []string          `json:"architectures"`
	BuildTags          []string          `json:"buildTags,omitempty"`
	IgnoreCombinations []PlatformIgnored `json:"ignoreCombinations,omitempty"`

	// Release Configuration
	DockerEnabled  bool   `json:"dockerEnabled"`
	SigningEnabled bool   `json:"signingEnabled"`
	DockerRegistry string `json:"dockerRegistry,omitempty"`
	DockerImage    string `json:"dockerImage,omitempty"`

	// Environment Variables
	Env map[string]string `json:"env"`
}

// GitHubActionsTemplateData represents strongly typed template data for GitHub Actions.
type GitHubActionsTemplateData struct {
	ProjectName    string   `json:"projectName"`
	Triggers       []string `json:"triggers"`
	DockerEnabled  bool     `json:"dockerEnabled"`
	SigningEnabled bool     `json:"signingEnabled"`
	DockerRegistry string   `json:"dockerRegistry,omitempty"`
	DockerImage    string   `json:"dockerImage,omitempty"`
}

// PlatformIgnored represents platform/architecture combinations to ignore.
type PlatformIgnored struct {
	GoOS   string `json:"goos"`
	GoArch string `json:"goarch"`
}

// DockerConfig represents Docker-specific configuration.
type DockerConfig struct {
	Enabled  bool   `json:"enabled"`
	Registry string `json:"registry"`
	Image    string `json:"image"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// SigningConfig represents code signing configuration.
type SigningConfig struct {
	Enabled     bool   `json:"enabled"`
	Level       string `json:"level"`
	KeyID       string `json:"keyId,omitempty"`
	Certificate string `json:"certificate,omitempty"`
}

// JobExecutionResult represents the result of a job execution.
type JobExecutionResult struct {
	JobID    string      `json:"jobId"`
	JobName  string      `json:"jobName"`
	Status   JobStatus   `json:"status"`
	Duration string      `json:"duration"`
	Error    *JobError   `json:"error,omitempty"`
	Output   string      `json:"output,omitempty"`
	Metadata JobMetadata `json:"metadata"`
}

// JobStatus represents the status of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// JobError represents a job execution error.
type JobError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Context   string `json:"context,omitempty"`
	Retryable bool   `json:"retryable"`
}

// JobMetadata represents job execution metadata.
type JobMetadata struct {
	StartedAt   string            `json:"startedAt"`
	CompletedAt string            `json:"completedAt,omitempty"`
	Retries     int               `json:"retries"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// WorkflowState represents the state of a workflow.
type WorkflowState struct {
	WorkflowID   string                `json:"workflowId"`
	WorkflowName string                `json:"workflowName"`
	Status       WorkflowStatus        `json:"status"`
	CurrentStep  string                `json:"currentStep"`
	TotalSteps   int                   `json:"totalSteps"`
	Progress     float64               `json:"progress"`
	StartedAt    string                `json:"startedAt"`
	Results      []*JobExecutionResult `json:"results"`
	Metadata     WorkflowMetadata      `json:"metadata"`
}

// WorkflowStatus represents the status of a workflow.
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// WorkflowMetadata represents workflow metadata.
type WorkflowMetadata struct {
	CreatedBy   string            `json:"createdBy"`
	Timeout     string            `json:"timeout"`
	Parallel    bool              `json:"parallel"`
	Environment string            `json:"environment"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// NewGoReleaserTemplateData creates typed template data from SafeProjectConfig.
func NewGoReleaserTemplateData(config *domain.SafeProjectConfig) *GoReleaserTemplateData {
	data := &GoReleaserTemplateData{
		ProjectName: config.ProjectName,
		BinaryName:  config.BinaryName,
		MainPath:    config.MainPath,
		Env: map[string]string{
			"GITHUB_OWNER": GetGitHubOwner(),
			"GITHUB_REPO":  GetGitHubRepo(),
		},
		IgnoreCombinations: []PlatformIgnored{
			{GoOS: "darwin", GoArch: "386"},
			{GoOS: "windows", GoArch: "arm64"},
		},
	}

	// Convert platforms
	if len(config.Platforms) > 0 {
		platforms := make([]string, 0, len(config.Platforms))
		for _, platform := range config.Platforms {
			platforms = append(platforms, string(platform))
		}

		data.Platforms = platforms
	}

	// Convert architectures
	if len(config.Architectures) > 0 {
		architectures := make([]string, 0, len(config.Architectures))
		for _, arch := range config.Architectures {
			architectures = append(architectures, string(arch))
		}

		data.Architectures = architectures
	}

	// Convert build tags
	if len(config.BuildTags) > 0 {
		tags := make([]string, 0, len(config.BuildTags))
		for _, tag := range config.BuildTags {
			tags = append(tags, tag.String())
		}

		data.BuildTags = tags
	}

	// Set Docker configuration
	if config.DockerSupport.IsEnabled() {
		data.DockerEnabled = true
		data.DockerRegistry = config.DockerRegistry.GetURL()
		data.DockerImage = config.GetDockerImageName()
	}

	// Set signing configuration
	data.SigningEnabled = config.SigningLevel.IsEnabled()

	return data
}

// NewGitHubActionsTemplateData creates typed template data from SafeProjectConfig.
func NewGitHubActionsTemplateData(config *domain.SafeProjectConfig) *GitHubActionsTemplateData {
	data := &GitHubActionsTemplateData{
		ProjectName:    config.ProjectName,
		DockerEnabled:  config.DockerSupport.IsEnabled(),
		SigningEnabled: config.SigningLevel.IsEnabled(),
	}

	// Convert action triggers
	if len(config.ActionsOn) > 0 {
		triggers := make([]string, 0, len(config.ActionsOn))
		for _, trigger := range config.ActionsOn {
			triggers = append(triggers, trigger.GitHubPattern())
		}

		data.Triggers = triggers
	}

	// Set Docker configuration
	if config.DockerSupport.IsEnabled() {
		data.DockerRegistry = config.DockerRegistry.GetURL()
		data.DockerImage = config.GetDockerImageName()
	}

	return data
}

// parseGitHubRemote tries to get GitHub owner and repo from git remote.
func parseGitHubRemote() (string, string) {
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	if cmd != nil {
		output, err := cmd.Output()
		if err == nil {
			remote := strings.TrimSpace(string(output))

			return git.ParseGitHubURL(remote)
		}
	}

	return "owner", "repo" // fallbacks
}

// GetGitHubOwner tries to get GitHub owner from git remote.
func GetGitHubOwner() string {
	owner, _ := parseGitHubRemote()

	return owner
}

// GetGitHubRepo tries to get GitHub repo from git remote.
func GetGitHubRepo() string {
	_, repo := parseGitHubRemote()

	return repo
}
