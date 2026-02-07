package jobs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/generators"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// Job represents a unit of work that can be executed.
type Job interface {
	ID() string
	Name() string
	Execute(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetMetadata() JobMetadata
}

// JobMetadata represents job metadata.
type JobMetadata struct {
	Description   string        `json:"description"`
	EstimatedTime time.Duration `json:"estimated_time"`
	Retryable     bool          `json:"retryable"`
	MaxRetries    int           `json:"max_retries"`
	Dependencies  []string      `json:"dependencies"`
	Tags          []string      `json:"tags"`
}

// Logger interface for jobs.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Rollbacker interface for types that can be rolled back.
type Rollbacker interface {
	Rollback(ctx context.Context) error
}

// ConfigGenerationJob generates GoReleaser configuration.
type ConfigGenerationJob struct {
	id        string
	config    *domain.SafeProjectConfig
	force     bool
	logger    Logger
	generator *generators.GoReleaserGenerator
}

// NewConfigGenerationJob creates a new config generation job.
func NewConfigGenerationJob(config *domain.SafeProjectConfig, force bool, logger Logger) *ConfigGenerationJob {
	return &ConfigGenerationJob{
		id:        "config-generation",
		config:    config,
		force:     force,
		logger:    logger,
		generator: generators.NewGoReleaserGenerator(config, logger),
	}
}

func (j *ConfigGenerationJob) ID() string {
	return j.id
}

func (j *ConfigGenerationJob) Name() string {
	return "Generate GoReleaser Configuration"
}

func (j *ConfigGenerationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Generates GoReleaser configuration from project settings",
		EstimatedTime: 5 * time.Second,
		Retryable:     false,
		MaxRetries:    0,
		Dependencies:  []string{},
		Tags:          []string{"config", "goreleaser"},
	}
}

func (j *ConfigGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating GoReleaser configuration")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Transition to Processing state
	j.config.State = domain.ConfigStateProcessing

	// Validate config
	if j.config.ProjectName == "" {
		j.config.State = domain.ConfigStateInvalid
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Project name is required",
			"Project name cannot be empty",
		).WithField("project_name")
	}

	// Check existing files
	if !j.force {
		if _, err := os.Stat(".goreleaser.yaml"); err == nil {
			j.config.State = domain.ConfigStateDraft
			return errors.NewConfigError(
				errors.ErrConfigFound,
				".goreleaser.yaml already exists",
				"Use --force to overwrite existing configuration",
			).WithRetryable(false)
		}
	}

	// Apply defaults to ensure config is complete
	j.config.ApplyDefaults()

	// Generate configuration
	err := j.generator.Generate(ctx)
	if err != nil {
		j.config.State = domain.ConfigStateInvalid
		return errors.WrapError(err, errors.ErrConfigGeneration, "Failed to generate GoReleaser config")
	}

	// Transition to Generated state
	j.config.State = domain.ConfigStateGenerated

	j.logger.Info("GoReleaser configuration generated successfully")
	return nil
}

func (j *ConfigGenerationJob) Rollback(ctx context.Context) error {
	j.logger.Info("Rolling back GoReleaser configuration generation")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated config
	if _, err := os.Stat(".goreleaser.yaml"); err == nil {
		// Check if backup exists
		if _, err := os.Stat(".goreleaser.yaml.backup"); err == nil {
			// Restore backup
			err := os.Rename(".goreleaser.yaml.backup", ".goreleaser.yaml")
			if err != nil {
				j.logger.Error("Failed to restore backup", "error", err)
				return errors.NewFileError(
					errors.ErrFileOperation,
					"Failed to restore backup",
					err.Error(),
				).WithCause(err)
			}
			j.logger.Info("Restored backup configuration")
		} else {
			// Remove generated file
			err := os.Remove(".goreleaser.yaml")
			if err != nil {
				j.logger.Error("Failed to remove generated config", "error", err)
				return errors.NewFileError(
					errors.ErrFileOperation,
					"Failed to remove generated config",
					err.Error(),
				).WithCause(err)
			}
			j.logger.Info("Removed generated configuration")
		}
	}

	return nil
}

// GitHubActionsGenerationJob generates GitHub Actions workflow.
type GitHubActionsGenerationJob struct {
	id        string
	config    *domain.SafeProjectConfig
	logger    Logger
	generator *generators.GitHubActionsGenerator
}

// NewGitHubActionsGenerationJob creates a new GitHub Actions generation job.
func NewGitHubActionsGenerationJob(config *domain.SafeProjectConfig, logger Logger) *GitHubActionsGenerationJob {
	return &GitHubActionsGenerationJob{
		id:        "github-actions-generation",
		config:    config,
		logger:    logger,
		generator: generators.NewGitHubActionsGenerator(config, logger),
	}
}

func (j *GitHubActionsGenerationJob) ID() string {
	return j.id
}

func (j *GitHubActionsGenerationJob) Name() string {
	return "Generate GitHub Actions Workflow"
}

func (j *GitHubActionsGenerationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Generates GitHub Actions workflow for automated releases",
		EstimatedTime: 3 * time.Second,
		Retryable:     true,
		MaxRetries:    2,
		Dependencies:  []string{"config-generation"},
		Tags:          []string{"workflow", "github", "ci-cd"},
	}
}

func (j *GitHubActionsGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating GitHub Actions workflow")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Generate workflow
	err := j.generator.Generate(ctx)
	if err != nil {
		return errors.WrapError(err, errors.ErrConfigGeneration, "Failed to generate GitHub Actions workflow")
	}

	j.logger.Info("GitHub Actions workflow generated successfully")
	return nil
}

func (j *GitHubActionsGenerationJob) Rollback(ctx context.Context) error {
	return rollbackGeneration(j.logger, j.generator, ctx, "GitHub Actions workflow generation")
}

// ProjectValidationJob validates project structure.
type ProjectValidationJob struct {
	id         string
	projectDir string
	logger     Logger
}

// NewProjectValidationJob creates a new project validation job.
func NewProjectValidationJob(projectDir string, logger Logger) *ProjectValidationJob {
	return &ProjectValidationJob{
		id:         "project-validation",
		projectDir: projectDir,
		logger:     logger,
	}
}

func (j *ProjectValidationJob) ID() string {
	return j.id
}

func (j *ProjectValidationJob) Name() string {
	return "Validate Project Structure"
}

func (j *ProjectValidationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Validates project structure and dependencies",
		EstimatedTime: 2 * time.Second,
		Retryable:     true,
		MaxRetries:    1,
		Dependencies:  []string{},
		Tags:          []string{"validation", "project"},
	}
}

func (j *ProjectValidationJob) Execute(ctx context.Context) error {
	j.logger.Info("Validating project structure")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Check if project directory exists
	if _, err := os.Stat(j.projectDir); os.IsNotExist(err) {
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Project directory does not exist",
			"Directory not found: "+j.projectDir,
		).WithField("project_dir")
	}

	// TODO: Implement comprehensive project validation
	// This would include checking for go.mod, main.go, etc.

	j.logger.Info("Project structure validation passed")
	return nil
}

func (j *ProjectValidationJob) Rollback(ctx context.Context) error {
	// Validation job doesn't create any files, so rollback is a no-op
	j.logger.Info("Project validation rollback is a no-op")
	return nil
}

// DependencyCheckJob checks for required dependencies.
type DependencyCheckJob struct {
	id           string
	dependencies []string
	logger       Logger
}

// NewDependencyCheckJob creates a new dependency check job.
func NewDependencyCheckJob(dependencies []string, logger Logger) *DependencyCheckJob {
	return &DependencyCheckJob{
		id:           "dependency-check",
		dependencies: dependencies,
		logger:       logger,
	}
}

func (j *DependencyCheckJob) ID() string {
	return j.id
}

func (j *DependencyCheckJob) Name() string {
	return "Check Dependencies"
}

func (j *DependencyCheckJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Checks if required dependencies are available",
		EstimatedTime: 1 * time.Second,
		Retryable:     false,
		MaxRetries:    0,
		Dependencies:  []string{},
		Tags:          []string{"dependencies", "validation"},
	}
}

func (j *DependencyCheckJob) Execute(ctx context.Context) error {
	j.logger.Info("Checking dependencies")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var missingDeps []string

	for _, dep := range j.dependencies {
		path, err := exec.LookPath(dep)
		if err != nil {
			missingDeps = append(missingDeps, dep)
			j.logger.Warn("Dependency not found", "dependency", dep)
		} else {
			j.logger.Debug("Found dependency", "dependency", dep, "path", path)
		}
	}

	if len(missingDeps) > 0 {
		return errors.NewValidationError(
			errors.ErrDependencyMissing,
			"Missing required dependencies",
			fmt.Sprintf("Required dependencies not found: %v", missingDeps),
		)
	}

	j.logger.Info("All dependencies are available")
	return nil
}

func (j *DependencyCheckJob) Rollback(ctx context.Context) error {
	// Dependency check doesn't modify state, so rollback is a no-op
	j.logger.Info("Dependency check rollback is a no-op")
	return nil
}

// DockerfileGenerationJob generates Dockerfile.
type DockerfileGenerationJob struct {
	id        string
	config    *domain.SafeProjectConfig
	force     bool
	logger    Logger
	generator *generators.DockerfileGenerator
}

// NewDockerfileGenerationJob creates a new Dockerfile generation job.
func NewDockerfileGenerationJob(config *domain.SafeProjectConfig, force bool, logger Logger) *DockerfileGenerationJob {
	return &DockerfileGenerationJob{
		id:        "dockerfile-generation",
		config:    config,
		force:     force,
		logger:    logger,
		generator: generators.NewDockerfileGenerator(config, logger),
	}
}

func (j *DockerfileGenerationJob) ID() string {
	return j.id
}

func (j *DockerfileGenerationJob) Name() string {
	return "Generate Dockerfile"
}

func (j *DockerfileGenerationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Generates Dockerfile for container builds",
		EstimatedTime: 3 * time.Second,
		Retryable:     false,
		MaxRetries:    0,
		Dependencies:  []string{},
		Tags:          []string{"docker", "container"},
	}
}

func (j *DockerfileGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating Dockerfile")

	// Skip if Docker support is disabled
	if !j.config.DockerSupport.ShouldBuild() {
		j.logger.Info("Docker support disabled, skipping Dockerfile generation")
		return nil
	}

	// Check context
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Generate Dockerfile
	err := j.generator.Generate(ctx)
	if err != nil {
		return errors.WrapError(err, errors.ErrConfigGeneration, "Failed to generate Dockerfile")
	}

	j.logger.Info("Dockerfile generated successfully")
	return nil
}

func (j *DockerfileGenerationJob) Rollback(ctx context.Context) error {
	return rollbackGeneration(j.logger, j.generator, ctx, "Dockerfile generation")
}

// HomebrewGenerationJob generates Homebrew formula.
type HomebrewGenerationJob struct {
	id        string
	config    *domain.SafeProjectConfig
	force     bool
	logger    Logger
	generator *generators.HomebrewGenerator
}

// NewHomebrewGenerationJob creates a new Homebrew generation job.
func NewHomebrewGenerationJob(config *domain.SafeProjectConfig, force bool, logger Logger) *HomebrewGenerationJob {
	return &HomebrewGenerationJob{
		id:        "homebrew-generation",
		config:    config,
		force:     force,
		logger:    logger,
		generator: generators.NewHomebrewGenerator(config, logger),
	}
}

func (j *HomebrewGenerationJob) ID() string {
	return j.id
}

func (j *HomebrewGenerationJob) Name() string {
	return "Generate Homebrew Formula"
}

func (j *HomebrewGenerationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Generates Homebrew formula for package distribution",
		EstimatedTime: 3 * time.Second,
		Retryable:     false,
		MaxRetries:    0,
		Dependencies:  []string{},
		Tags:          []string{"homebrew", "package", "distribution"},
	}
}

func (j *HomebrewGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating Homebrew formula")

	// Skip if Homebrew support is disabled
	if !j.config.Homebrew {
		j.logger.Info("Homebrew support disabled, skipping formula generation")
		return nil
	}

	// Check context
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Generate formula
	err := j.generator.Generate(ctx)
	if err != nil {
		return errors.WrapError(err, errors.ErrConfigGeneration, "Failed to generate Homebrew formula")
	}

	j.logger.Info("Homebrew formula generated successfully")
	return nil
}

func (j *HomebrewGenerationJob) Rollback(ctx context.Context) error {
	return rollbackGeneration(j.logger, j.generator, ctx, "Homebrew formula generation")
}

// rollbackGeneration is a helper function that handles the common rollback pattern for generators.
func rollbackGeneration(logger Logger, rollbacker Rollbacker, ctx context.Context, actionName string) error {
	logger.Info("Rolling back " + actionName)

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Rollback generator
	err := rollbacker.Rollback(ctx)
	if err != nil {
		return errors.WrapError(err, errors.ErrWorkflowExecution, "Failed to rollback "+actionName)
	}

	return nil
}
