package interfaces

import (
	"context"
	"io"
	"os"
)

// FileSystemRepository handles all file system operations
// This interface follows the Repository pattern and Interface Segregation Principle
type FileSystemRepository interface {
	// File operations
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error
	CreateFile(ctx context.Context, path string) (io.WriteCloser, error)
	DeleteFile(ctx context.Context, path string) error
	FileExists(ctx context.Context, path string) (bool, error)

	// Directory operations
	CreateDir(ctx context.Context, path string, perm os.FileMode) error
	CreateDirAll(ctx context.Context, path string, perm os.FileMode) error
	DirExists(ctx context.Context, path string) (bool, error)
	ReadDir(ctx context.Context, path string) ([]os.DirEntry, error)

	// Permission and metadata operations
	GetFileInfo(ctx context.Context, path string) (os.FileInfo, error)
	CheckPermissions(ctx context.Context, path string) (bool, error)

	// Utility operations
	AbsPath(path string) (string, error)
	RelPath(base, target string) (string, error)
	CleanPath(path string) string
	JoinPath(elem ...string) string
	TempDir(dir, pattern string) (string, error)
}

// TemplateRepository handles template rendering and management
// This interface follows the Repository pattern and Interface Segregation Principle
type TemplateRepository interface {
	// Template operations
	LoadTemplate(ctx context.Context, name string) (string, error)
	RenderTemplate(ctx context.Context, templateContent string, data any) (string, error)
	ValidateTemplate(ctx context.Context, content string) error

	// Template discovery and management
	ListTemplates(ctx context.Context) ([]string, error)
	TemplateExists(ctx context.Context, name string) (bool, error)
	GetTemplatePath(ctx context.Context, name string) (string, error)

	// Template categories
	GetTemplatesForType(ctx context.Context, projectType ProjectType) ([]string, error)
	GetTemplateMetadata(ctx context.Context, name string) (TemplateMetadata, error)
}

// GoReleaserRepository handles GoReleaser integration
// This interface follows the Repository pattern and Interface Segregation Principle
type GoReleaserRepository interface {
	// Configuration management
	ValidateConfig(ctx context.Context, config *SafeProjectConfig) error
	CheckConfig(ctx context.Context, configPath string) error
	GenerateConfig(ctx context.Context, config *SafeProjectConfig) (string, error)

	// Build operations
	BuildSnapshot(ctx context.Context, config *SafeProjectConfig) error
	BuildRelease(ctx context.Context, config *SafeProjectConfig, version string) error

	// Release operations
	ReleaseDryRun(ctx context.Context, config *SafeProjectConfig) error
	Release(ctx context.Context, config *SafeProjectConfig, version string) error

	// Utility operations
	GetSupportedPlatforms() ([]Platform, error)
	GetSupportedArchitectures() ([]Architecture, error)
	GetVersion() (string, error)
	CheckInstallation() (bool, error)
}

// GitHubRepository handles GitHub API operations
// This interface follows the Repository pattern and Interface Segregation Principle
type GitHubRepository interface {
	// Repository operations
	GetRepo(ctx context.Context, owner, name string) (*GitHubRepo, error)
	CreateRepo(ctx context.Context, repo *GitHubRepoRequest) (*GitHubRepo, error)
	UpdateRepo(ctx context.Context, owner, name string, updates *GitHubRepoUpdate) (*GitHubRepo, error)

	// Branch and tag operations
	GetBranches(ctx context.Context, owner, name string) ([]*GitHubBranch, error)
	GetTags(ctx context.Context, owner, name string) ([]*GitHubTag, error)
	CreateTag(ctx context.Context, owner, name string, tag *GitHubTagRequest) (*GitHubTag, error)

	// Release operations
	GetReleases(ctx context.Context, owner, name string) ([]*GitHubRelease, error)
	CreateRelease(ctx context.Context, owner, name string, release *GitHubReleaseRequest) (*GitHubRelease, error)
	UploadReleaseAsset(ctx context.Context, owner, name, releaseID string, asset *GitHubAsset) (*GitHubAsset, error)

	// Workflow operations
	GetWorkflows(ctx context.Context, owner, name string) ([]*GitHubWorkflow, error)
	TriggerWorkflow(ctx context.Context, owner, name, workflowID string, inputs map[string]any) error

	// Authentication and permissions
	ValidateToken(ctx context.Context, token string) error
	GetUser(ctx context.Context, token string) (*GitHubUser, error)
	CheckPermissions(ctx context.Context, token, owner, repo string) (*GitHubPermissions, error)
}

// JobRepository handles job execution and management
// This interface follows the Repository pattern and Interface Segregation Principle
type JobRepository interface {
	// Job operations
	CreateJob(ctx context.Context, job Job) (*Job, error)
	GetJob(ctx context.Context, id string) (*Job, error)
	UpdateJob(ctx context.Context, job Job) (*Job, error)
	DeleteJob(ctx context.Context, id string) error
	ListJobs(ctx context.Context, filter JobFilter) ([]*Job, error)

	// Job execution operations
	ExecuteJob(ctx context.Context, id string, options JobExecutionOptions) (*JobExecutionResult, error)
	CancelJob(ctx context.Context, id string) error
	RetryJob(ctx context.Context, id string) error
	GetJobHistory(ctx context.Context, id string) ([]*JobExecutionResult, error)

	// Job dependency operations
	AddDependency(ctx context.Context, jobID, dependsOn string) error
	RemoveDependency(ctx context.Context, jobID, dependsOn string) error
	GetDependencies(ctx context.Context, jobID string) ([]string, error)
}

// WorkflowRepository handles workflow execution and management
// This interface follows the Repository pattern and Interface Segregation Principle
type WorkflowRepository interface {
	// Workflow operations
	CreateWorkflow(ctx context.Context, workflow Workflow) (*Workflow, error)
	GetWorkflow(ctx context.Context, id string) (*Workflow, error)
	UpdateWorkflow(ctx context.Context, workflow Workflow) (*Workflow, error)
	DeleteWorkflow(ctx context.Context, id string) error
	ListWorkflows(ctx context.Context, filter WorkflowFilter) ([]*Workflow, error)

	// Workflow execution operations
	ExecuteWorkflow(ctx context.Context, id string, options WorkflowOptions) (*WorkflowExecution, error)
	CancelWorkflow(ctx context.Context, id string) error
	RetryWorkflow(ctx context.Context, id string) error
	GetWorkflowHistory(ctx context.Context, id string) ([]*WorkflowExecution, error)

	// Workflow job operations
	AddJobToWorkflow(ctx context.Context, workflowID, jobID string) error
	RemoveJobFromWorkflow(ctx context.Context, workflowID, jobID string) error
	GetWorkflowJobs(ctx context.Context, workflowID string) ([]string, error)
}

// ValidationRepository handles validation operations
// This interface follows the Repository pattern and Interface Segregation Principle
type ValidationRepository interface {
	// Configuration validation
	ValidateProjectConfig(ctx context.Context, config *SafeProjectConfig) (*ValidationResult, error)
	ValidateTemplate(ctx context.Context, template string) (*ValidationResult, error)
	ValidateDependencies(ctx context.Context, dependencies []string) (*ValidationResult, error)

	// Project validation
	ValidateProjectStructure(ctx context.Context, projectDir string) (*ValidationResult, error)
	ValidateGitRepository(ctx context.Context, repoDir string) (*ValidationResult, error)
	ValidateGoModule(ctx context.Context, modulePath string) (*ValidationResult, error)

	// Platform/architecture validation
	ValidatePlatformArchCompatibility(ctx context.Context, platforms []Platform, archs []Architecture) (*ValidationResult, error)
	ValidateBuildTags(ctx context.Context, tags []BuildTag) (*ValidationResult, error)

	// Security validation
	ValidateSecurity(ctx context.Context, config *SafeProjectConfig) (*ValidationResult, error)
	ValidatePermissions(ctx context.Context, paths []string) (*ValidationResult, error)
}
