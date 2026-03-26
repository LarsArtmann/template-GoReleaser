package jobs

import (
	"fmt"
	"slices"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// JobExecutionStatus represents the status of a job execution.
type JobExecutionStatus struct {
	JobID       domain.JobID           `json:"job_id"`
	JobName     string                 `json:"job_name"`
	Status      JobExecutionStatusType `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Error       *types.JobError        `json:"error,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    JobExecutionMetadata   `json:"metadata"`
}

// JobExecutionStatusType represents job status types.
type JobExecutionStatusType string

const (
	JobExecutionStatusPending   JobExecutionStatusType = "pending"
	JobExecutionStatusRunning   JobExecutionStatusType = "running"
	JobExecutionStatusCompleted JobExecutionStatusType = "completed"
	JobExecutionStatusFailed    JobExecutionStatusType = "failed"
	JobExecutionStatusCancelled JobExecutionStatusType = "cancelled"
)

// JobExecutionMetadata represents metadata for job execution.
type JobExecutionMetadata struct {
	Retries int               `json:"retries"`
	Tags    map[string]string `json:"tags,omitempty"`
	Context map[string]any    `json:"context,omitempty"`
}

// WorkflowExecution represents a workflow execution with multiple jobs.
type WorkflowExecution struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Status      WorkflowExecutionStatusType `json:"status"`
	StartTime   time.Time                   `json:"start_time"`
	EndTime     *time.Time                  `json:"end_time,omitempty"`
	Duration    time.Duration               `json:"duration"`
	Parallel    bool                        `json:"parallel"`
	Timeout     time.Duration               `json:"timeout"`
	CurrentStep int                         `json:"current_step"`
	TotalSteps  int                         `json:"total_steps"`
	Progress    float64                     `json:"progress"`
	Jobs        []*JobExecutionStatus       `json:"jobs"`
	Results     []*types.JobExecutionResult `json:"results"`
	Errors      []*types.JobError           `json:"errors,omitempty"`
	Metadata    WorkflowExecutionMetadata   `json:"metadata"`
}

// WorkflowExecutionStatusType represents workflow status types.
type WorkflowExecutionStatusType string

const (
	WorkflowExecutionStatusPending   WorkflowExecutionStatusType = "pending"
	WorkflowExecutionStatusRunning   WorkflowExecutionStatusType = "running"
	WorkflowExecutionStatusCompleted WorkflowExecutionStatusType = "completed"
	WorkflowExecutionStatusFailed    WorkflowExecutionStatusType = "failed"
	WorkflowExecutionStatusCancelled WorkflowExecutionStatusType = "cancelled"
)

// WorkflowExecutionMetadata represents workflow execution metadata.
type WorkflowExecutionMetadata struct {
	CreatedBy   string            `json:"created_by"`
	Environment string            `json:"environment"`
	Tags        map[string]string `json:"tags,omitempty"`
	Options     map[string]any    `json:"options,omitempty"`
}

// JobResult represents the result of a job execution.
type JobResult struct {
	Job      Job                    `json:"-"`
	Status   JobExecutionStatusType `json:"status"`
	Duration time.Duration          `json:"duration"`
	Error    error                  `json:"error,omitempty"`
	Output   string                 `json:"output,omitempty"`
	Metadata JobExecutionMetadata   `json:"metadata"`
}

// WorkflowOptions represents options for workflow execution.
type WorkflowOptions struct {
	Parallel        bool              `json:"parallel"`
	Timeout         time.Duration     `json:"timeout"`
	RetryFailed     bool              `json:"retry_failed"`
	MaxRetries      int               `json:"max_retries"`
	ContinueOnError bool              `json:"continue_on_error"`
	Environment     string            `json:"environment"`
	Tags            map[string]string `json:"tags,omitempty"`
	Context         map[string]any    `json:"context,omitempty"`
}

// DefaultWorkflowOptions returns default workflow options.
func DefaultWorkflowOptions() WorkflowOptions {
	return WorkflowOptions{
		Parallel:        false,
		Timeout:         defaultWorkflowTimeout,
		RetryFailed:     false,
		MaxRetries:      2,
		ContinueOnError: false,
		Environment:     "development",
		Tags:            make(map[string]string),
		Context:         make(map[string]any),
	}
}

// JobExecutionOptions represents options for job execution.
type JobExecutionOptions struct {
	Timeout       time.Duration     `json:"timeout"`
	RetryCount    int               `json:"retry_count"`
	MaxRetries    int               `json:"max_retries"`
	IgnoreError   bool              `json:"ignore_error"`
	SkipExecution bool              `json:"skip_execution"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// DefaultJobExecutionOptions returns default job execution options.
func DefaultJobExecutionOptions() JobExecutionOptions {
	return JobExecutionOptions{
		Timeout:       defaultJobTimeout,
		RetryCount:    0,
		MaxRetries:    2,
		IgnoreError:   false,
		SkipExecution: false,
		Metadata:      make(map[string]string),
	}
}

// JobExecutionPlan represents a plan for executing jobs.
type JobExecutionPlan struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Jobs         []Job               `json:"jobs"`
	Dependencies map[string][]string `json:"dependencies"`
	Options      WorkflowOptions     `json:"options"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// NewJobExecutionPlan creates a new job execution plan.
func NewJobExecutionPlan(
	id, name, description string,
	jobs []Job,
	options WorkflowOptions,
) *JobExecutionPlan {
	return &JobExecutionPlan{
		ID:           id,
		Name:         name,
		Description:  description,
		Jobs:         jobs,
		Dependencies: make(map[string][]string),
		Options:      options,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// AddDependency adds a dependency between jobs.
func (p *JobExecutionPlan) AddDependency(jobID, dependsOn string) {
	if p.Dependencies == nil {
		p.Dependencies = make(map[string][]string)
	}

	p.Dependencies[jobID] = append(p.Dependencies[jobID], dependsOn)
	p.UpdatedAt = time.Now()
}

// RemoveDependency removes a dependency between jobs.
func (p *JobExecutionPlan) RemoveDependency(jobID, dependsOn string) {
	if deps, ok := p.Dependencies[jobID]; ok {
		for i, dep := range deps {
			if dep == dependsOn {
				p.Dependencies[jobID] = append(deps[:i], deps[i+1:]...)

				break
			}
		}

		if len(p.Dependencies[jobID]) == 0 {
			delete(p.Dependencies, jobID)
		}

		p.UpdatedAt = time.Now()
	}
}

// HasDependency checks if a job has a dependency.
func (p *JobExecutionPlan) HasDependency(jobID, dependsOn string) bool {
	if deps, ok := p.Dependencies[jobID]; ok {
		if slices.Contains(deps, dependsOn) {
			return true
		}
	}

	return false
}

// GetDependencies returns all dependencies for a job.
func (p *JobExecutionPlan) GetDependencies(jobID string) []string {
	if deps, ok := p.Dependencies[jobID]; ok {
		result := make([]string, len(deps))
		copy(result, deps)

		return result
	}

	return nil
}

// ValidatePlan validates the execution plan.
func (p *JobExecutionPlan) ValidatePlan() error {
	// Check for circular dependencies
	err := p.validateCircularDependencies()
	if err != nil {
		return errors.NewValidationError(
			errors.ErrInvalidConfig,
			"Invalid job execution plan",
			err.Error(),
		)
	}

	// Check for missing dependencies
	err = p.validateMissingDependencies()
	if err != nil {
		return errors.NewValidationError(
			errors.ErrInvalidConfig,
			"Invalid job execution plan",
			err.Error(),
		)
	}

	// Check for self-dependencies
	err = p.validateSelfDependencies()
	if err != nil {
		return errors.NewValidationError(
			errors.ErrInvalidConfig,
			"Invalid job execution plan",
			err.Error(),
		)
	}

	return nil
}

// validateCircularDependencies checks for circular dependencies.
func (p *JobExecutionPlan) validateCircularDependencies() error {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for _, job := range p.Jobs {
		err := p.checkCircularDependency(job.ID(), visited, recursionStack)
		if err != nil {
			return err
		}
	}

	return nil
}

// checkCircularDependency recursively checks for circular dependencies.
func (p *JobExecutionPlan) checkCircularDependency(
	jobID string,
	visited, recursionStack map[string]bool,
) error {
	if recursionStack[jobID] {
		return errors.NewValidationError(
			errors.ErrInvalidConfig,
			"Circular dependency detected",
			fmt.Sprintf("Job %s is part of a circular dependency", jobID),
		).WithField("job_id")
	}

	if visited[jobID] {
		return nil
	}

	visited[jobID] = true
	recursionStack[jobID] = true

	for _, depID := range p.GetDependencies(jobID) {
		err := p.checkCircularDependency(depID, visited, recursionStack)
		if err != nil {
			return fmt.Errorf("checking circular dependency for job %s: %w", jobID, err)
		}
	}

	delete(recursionStack, jobID)

	return nil
}

// validateMissingDependencies checks for missing dependencies.
func (p *JobExecutionPlan) validateMissingDependencies() error {
	jobIDs := make(map[string]bool)
	for _, job := range p.Jobs {
		jobIDs[job.ID()] = true
	}

	for jobID, deps := range p.Dependencies {
		for _, depID := range deps {
			if !jobIDs[depID] {
				return errors.NewValidationError(
					errors.ErrInvalidConfig,
					"Missing dependency",
					fmt.Sprintf("Job %s depends on non-existent job %s", jobID, depID),
				).WithField("job_id").WithField("dependency")
			}
		}
	}

	return nil
}

// validateSelfDependencies checks for self-dependencies.
func (p *JobExecutionPlan) validateSelfDependencies() error {
	for jobID, deps := range p.Dependencies {
		if slices.Contains(deps, jobID) {
			return errors.NewValidationError(
				errors.ErrInvalidConfig,
				"Self-dependency detected",
				fmt.Sprintf("Job %s depends on itself", jobID),
			).WithField("job_id")
		}
	}

	return nil
}

// Clone creates a copy of the execution plan.
func (p *JobExecutionPlan) Clone() *JobExecutionPlan {
	clone := &JobExecutionPlan{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Jobs:         make([]Job, len(p.Jobs)),
		Dependencies: make(map[string][]string),
		Options:      p.Options,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	// Copy jobs
	copy(clone.Jobs, p.Jobs)

	// Copy dependencies
	for jobID, deps := range p.Dependencies {
		clone.Dependencies[jobID] = make([]string, len(deps))
		copy(clone.Dependencies[jobID], deps)
	}

	return clone
}
