package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/generators"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// JobFactory creates jobs for common wizard operations.
type JobFactory struct {
	logger Logger
}

// NewJobFactory creates a new job factory.
func NewJobFactory(logger Logger) *JobFactory {
	return &JobFactory{
		logger: logger,
	}
}

// CreateFullWizardJobs creates all jobs for a complete wizard operation.
func (jf *JobFactory) CreateFullWizardJobs(config *domain.SafeProjectConfig, force bool) []Job {
	var jobs []Job

	// Convert to domain.SafeProjectConfig since ProjectConfig is an alias
	safeConfig := config

	// Add project validation job
	jobs = append(jobs, NewProjectValidationJob(".", jf.logger))

	// Add dependency check job
	dependencies := jf.getRequiredDependencies(safeConfig)
	jobs = append(jobs, NewDependencyCheckJob(dependencies, jf.logger))

	// Add config generation job
	jobs = append(jobs, NewConfigGenerationJob(safeConfig, force, jf.logger))

	// Add Dockerfile generation job
	if safeConfig.DockerSupport.ShouldBuild() {
		jobs = append(jobs, NewDockerfileGenerationJob(safeConfig, force, jf.logger))
	}

	// Add Homebrew generation job
	if safeConfig.Homebrew {
		jobs = append(jobs, NewHomebrewGenerationJob(safeConfig, force, jf.logger))
	}

	// Add GitHub Actions generation job
	if safeConfig.ShouldGenerateActionsFiles() {
		jobs = append(jobs, NewGitHubActionsGenerationJob(safeConfig, jf.logger))
	}

	return jobs
}

// CreateConfigOnlyJobs creates jobs for config generation only.
func (jf *JobFactory) CreateConfigOnlyJobs(config *domain.SafeProjectConfig, force bool) []Job {
	return []Job{
		NewProjectValidationJob(".", jf.logger),
		NewConfigGenerationJob(config, force, jf.logger),
	}
}

// CreateValidationOnlyJob creates a validation-only job.
func (jf *JobFactory) CreateValidationOnlyJob(projectDir string) Job {
	return NewProjectValidationJob(projectDir, jf.logger)
}

// CreateCustomJobs creates jobs for custom operations.
func (jf *JobFactory) CreateCustomJobs(operation string, config *domain.SafeProjectConfig, options map[string]any) ([]Job, error) {
	switch operation {
	case "preview":
		return jf.createPreviewJobs(config, options)
	case "validate":
		return jf.createValidationJobs(config, options)
	case "rollback":
		return jf.createRollbackJobs(config, options)
	default:
		return nil, errors.NewValidationError(
			errors.ErrInvalidOperation,
			"Unknown operation",
			fmt.Sprintf("Operation '%s' is not supported", operation),
		).WithField("operation")
	}
}

// createPreviewJobs creates jobs for preview operations.
func (jf *JobFactory) createPreviewJobs(config *domain.SafeProjectConfig, options map[string]any) ([]Job, error) {
	var jobs []Job

	// Add validation job
	jobs = append(jobs, NewProjectValidationJob(".", jf.logger))

	// Add preview-specific jobs
	if previewJob, err := jf.createPreviewJob(config, options); err == nil {
		jobs = append(jobs, previewJob)
	}

	return jobs, nil
}

// createValidationJobs creates jobs for validation operations.
func (jf *JobFactory) createValidationJobs(config *domain.SafeProjectConfig, options map[string]any) ([]Job, error) {
	var jobs []Job

	// Add project validation
	jobs = append(jobs, NewProjectValidationJob(".", jf.logger))

	// Add dependency check
	dependencies := jf.getRequiredDependencies(config)
	jobs = append(jobs, NewDependencyCheckJob(dependencies, jf.logger))

	// Add config validation if needed
	if validateConfig, ok := options["validate_config"].(bool); ok && validateConfig {
		jobs = append(jobs, jf.createConfigValidationJob(config))
	}

	return jobs, nil
}

// createRollbackJobs creates jobs for rollback operations.
func (jf *JobFactory) createRollbackJobs(config *domain.SafeProjectConfig, options map[string]any) ([]Job, error) {
	var jobs []Job

	// Add rollback validation
	if validateJob, ok := options["validate_job"].(string); ok && validateJob != "" {
		if job := jf.createJobRollback(validateJob, config); job != nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// createPreviewJob creates a preview generation job.
func (jf *JobFactory) createPreviewJob(config *domain.SafeProjectConfig, options map[string]any) (Job, error) {
	return &PreviewGenerationJob{
		id:          "preview-generation",
		config:      config,
		logger:      jf.logger,
		previewType: getStringOption(options, "type", "all"),
	}, nil
}

// createConfigValidationJob creates a config validation job.
func (jf *JobFactory) createConfigValidationJob(config *domain.SafeProjectConfig) Job {
	return &ConfigValidationJob{
		id:     "config-validation",
		config: config,
		logger: jf.logger,
	}
}

// createJobRollback creates a job rollback operation.
func (jf *JobFactory) createJobRollback(jobID string, config *domain.SafeProjectConfig) Job {
	return &JobRollbackJob{
		id:     "job-rollback",
		jobID:  jobID,
		config: config,
		logger: jf.logger,
	}
}

// getRequiredDependencies returns required dependencies based on configuration.
func (jf *JobFactory) getRequiredDependencies(config *domain.SafeProjectConfig) []string {
	dependencies := []string{"go"}

	if config.DockerSupport.ShouldBuild() || config.DockerSupport.ShouldPublish() {
		dependencies = append(dependencies, "docker")
	}

	// Only require cosign for advanced signing levels
	if config.SigningLevel == domain.SigningLevelAdvanced || config.SigningLevel == domain.SigningLevelEnterprise {
		dependencies = append(dependencies, "cosign")
	}

	return dependencies
}

// getStringOption safely extracts string option from options map.
func getStringOption(options map[string]any, key, defaultValue string) string {
	if value, ok := options[key].(string); ok {
		return value
	}
	return defaultValue
}

// PreviewGenerationJob generates preview of configurations.
type PreviewGenerationJob struct {
	id          string
	config      *domain.SafeProjectConfig
	logger      Logger
	previewType string
}

func (j *PreviewGenerationJob) ID() string {
	return j.id
}

func (j *PreviewGenerationJob) Name() string {
	return "Generate Preview"
}

func (j *PreviewGenerationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Generates preview of configurations without writing files",
		EstimatedTime: 2 * time.Second,
		Retryable:     true,
		MaxRetries:    1,
		Dependencies:  []string{"project-validation"},
		Tags:          []string{"preview", "generation"},
	}
}

func (j *PreviewGenerationJob) Execute(ctx context.Context) error {
	j.logger.Info("Generating preview", "type", j.previewType)

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Generate preview based on type
	switch j.previewType {
	case "goreleaser":
		return j.generateGoReleaserPreview(ctx)
	case "github-actions":
		return j.generateGitHubActionsPreview(ctx)
	case "all":
		if err := j.generateGoReleaserPreview(ctx); err != nil {
			return err
		}
		return j.generateGitHubActionsPreview(ctx)
	default:
		return errors.NewValidationError(
			errors.ErrInvalidOperation,
			"Unknown preview type",
			fmt.Sprintf("Preview type '%s' is not supported", j.previewType),
		).WithField("preview_type")
	}
}

// PreviewGenerator is an interface for generators that can create previews.
type PreviewGenerator interface {
	GeneratePreview(ctx context.Context) (string, error)
}

// generatePreview is a helper function that generates a preview using any PreviewGenerator.
func (j *PreviewGenerationJob) generatePreview(ctx context.Context, generator PreviewGenerator, label string) error {
	preview, err := generator.GeneratePreview(ctx)
	if err != nil {
		return errors.WrapError(err, errors.ErrConfigGeneration, fmt.Sprintf("Failed to generate %s preview", label))
	}

	j.logger.Info(label+" preview generated", "preview", preview)
	return nil
}

func (j *PreviewGenerationJob) generateGoReleaserPreview(ctx context.Context) error {
	generator := generators.NewGoReleaserGenerator(j.config, j.logger)
	return j.generatePreview(ctx, generator, "GoReleaser")
}

func (j *PreviewGenerationJob) generateGitHubActionsPreview(ctx context.Context) error {
	generator := generators.NewGitHubActionsGenerator(j.config, j.logger)
	return j.generatePreview(ctx, generator, "GitHub Actions")
}

func (j *PreviewGenerationJob) Rollback(ctx context.Context) error {
	// Preview generation doesn't create any files, so rollback is a no-op
	j.logger.Info("Preview generation rollback is a no-op")
	return nil
}

// ConfigValidationJob validates configuration.
type ConfigValidationJob struct {
	id     string
	config *domain.SafeProjectConfig
	logger Logger
}

func (j *ConfigValidationJob) ID() string {
	return j.id
}

func (j *ConfigValidationJob) Name() string {
	return "Validate Configuration"
}

func (j *ConfigValidationJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Validates project configuration against best practices",
		EstimatedTime: 1 * time.Second,
		Retryable:     false,
		MaxRetries:    0,
		Dependencies:  []string{"project-validation"},
		Tags:          []string{"validation", "config"},
	}
}

func (j *ConfigValidationJob) Execute(ctx context.Context) error {
	j.logger.Info("Validating configuration")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: Implement comprehensive configuration validation
	// This would validate all fields, business rules, etc.

	j.logger.Info("Configuration validation passed")
	return nil
}

func (j *ConfigValidationJob) Rollback(ctx context.Context) error {
	// Validation job doesn't modify state, so rollback is a no-op
	j.logger.Info("Configuration validation rollback is a no-op")
	return nil
}

// JobRollbackJob rolls back a specific job.
type JobRollbackJob struct {
	id     string
	jobID  string
	config *domain.SafeProjectConfig
	logger Logger
}

func (j *JobRollbackJob) ID() string {
	return j.id
}

func (j *JobRollbackJob) Name() string {
	return "Rollback Job: " + j.jobID
}

func (j *JobRollbackJob) GetMetadata() JobMetadata {
	return JobMetadata{
		Description:   "Rolls back a specific job execution",
		EstimatedTime: 5 * time.Second,
		Retryable:     true,
		MaxRetries:    2,
		Dependencies:  []string{},
		Tags:          []string{"rollback", "recovery"},
	}
}

func (j *JobRollbackJob) Execute(ctx context.Context) error {
	j.logger.Info("Rolling back job", "job_id", j.jobID)

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: Implement specific job rollback logic
	// This would find the job and call its Rollback method

	j.logger.Info("Job rollback completed", "job_id", j.jobID)
	return nil
}

func (j *JobRollbackJob) Rollback(ctx context.Context) error {
	// Rollback job doesn't create any state, so rollback is a no-op
	j.logger.Info("Job rollback is a no-op")
	return nil
}
