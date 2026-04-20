package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"charm.land/log/v2"
)

// checkContextCancellation checks if context is cancelled and returns an error if so.
func checkContextCancellation(ctx context.Context, msg string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", msg, ctx.Err())
	default:
		return nil
	}
}

// Job represents a wizard operation job.
type Job interface {
	ID() string
	Name() string
	Execute(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// JobStatus represents the status of a job.
type JobStatus int

const (
	JobStatusPending JobStatus = iota
	JobStatusRunning
	JobStatusCompleted
	JobStatusFailed
	JobStatusRolledBack
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusPending:
		return "pending"
	case JobStatusRunning:
		return "running"
	case JobStatusCompleted:
		return "completed"
	case JobStatusFailed:
		return "failed"
	case JobStatusRolledBack:
		return "rolled_back"
	default:
		return "unknown"
	}
}

// JobResult represents the result of a job execution.
type JobResult struct {
	Job      Job
	Status   JobStatus
	Error    error
	Duration time.Duration
	Started  time.Time
	Finished time.Time
	Output   string
}

// JobManager manages the execution of wizard jobs.
type JobManager struct {
	jobs        []Job
	results     []JobResult
	mu          sync.Mutex
	logger      *log.Logger
	parallel    bool
	maxJobs     int
	currentJobs int
}

// NewJobManager creates a new job manager.
func NewJobManager(logger *log.Logger) *JobManager {
	return &JobManager{
		jobs:        make([]Job, 0),
		results:     make([]JobResult, 0),
		logger:      logger,
		parallel:    false,
		maxJobs:     3, // Default max parallel jobs
		currentJobs: 0,
	}
}

// SetParallel sets whether jobs should run in parallel.
func (jm *JobManager) SetParallel(parallel bool) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.parallel = parallel
}

// SetMaxJobs sets the maximum number of parallel jobs.
func (jm *JobManager) SetMaxJobs(maxJobs int) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.maxJobs = maxJobs
}

// AddJob adds a job to the manager.
func (jm *JobManager) AddJob(job Job) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.jobs = append(jm.jobs, job)
}

// ExecuteJobs executes all jobs according to the manager settings.
func (jm *JobManager) ExecuteJobs(ctx context.Context) error {
	if jm.parallel {
		return jm.executeParallel(ctx)
	}

	return jm.executeSequential(ctx)
}

// executeSequential executes jobs one by one.
func (jm *JobManager) executeSequential(ctx context.Context) error {
	for _, job := range jm.jobs {
		err := checkContext(ctx)
		if err != nil {
			return err
		}

		result := jm.executeJob(ctx, job)
		jm.addResult(result)

		if result.Status == JobStatusFailed {
			return fmt.Errorf("job %s failed: %w", job.Name(), result.Error)
		}
	}

	return nil
}

// executeParallel executes jobs in parallel with concurrency limits.
func (jm *JobManager) executeParallel(ctx context.Context) error {
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, jm.maxJobs)
	errChan := make(chan error, len(jm.jobs))

	for _, job := range jm.jobs {
		if err := checkContextCancellation(ctx, "context cancelled"); err != nil {
			return err
		}

		wg.Add(1)

		go func(j Job) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}

			defer func() { <-semaphore }()

			result := jm.executeJob(ctx, j)
			jm.addResult(result)

			if result.Status == JobStatusFailed {
				errChan <- fmt.Errorf("job %s failed: %w", j.Name(), result.Error)
			}
		}(job)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// executeJob executes a single job and records the result.
func (jm *JobManager) executeJob(ctx context.Context, job Job) JobResult {
	start := time.Now()

	// Update job status
	jm.updateJobStatus(job.ID(), JobStatusRunning)
	jm.logger.Infof("Executing job: %s", job.Name())

	// Execute the job
	err := job.Execute(ctx)
	duration := time.Since(start)
	finished := time.Now()

	status := JobStatusCompleted
	if err != nil {
		status = JobStatusFailed

		jm.logger.Errorf("Job %s failed: %v", job.Name(), err)
	} else {
		jm.logger.Infof("Job %s completed successfully", job.Name())
	}

	result := JobResult{
		Job:      job,
		Status:   status,
		Error:    err,
		Duration: duration,
		Started:  start,
		Finished: finished,
		Output:   fmt.Sprintf("Job %s %s", job.Name(), status),
	}

	// Update job status
	jm.updateJobStatus(job.ID(), status)

	return result
}

// updateJobStatus updates the status of a job (for UI updates).
func (jm *JobManager) updateJobStatus(jobID string, status JobStatus) {
	// This could be extended to update a UI or event system
	jm.logger.Debugf("Job %s status: %s", jobID, status)
}

// addResult adds a result to the results list.
func (jm *JobManager) addResult(result JobResult) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.results = append(jm.results, result)
}

// GetResults returns all job results.
func (jm *JobManager) GetResults() []JobResult {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	results := make([]JobResult, len(jm.results))
	copy(results, jm.results)

	return results
}

// getResultsByStatus returns job results filtered by status.
func (jm *JobManager) getResultsByStatus(status JobStatus) []JobResult {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	filtered := make([]JobResult, 0)

	for _, result := range jm.results {
		if result.Status == status {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// GetCompletedResults returns only completed job results.
func (jm *JobManager) GetCompletedResults() []JobResult {
	return jm.getResultsByStatus(JobStatusCompleted)
}

// GetFailedResults returns only failed job results.
func (jm *JobManager) GetFailedResults() []JobResult {
	return jm.getResultsByStatus(JobStatusFailed)
}

// RollbackFailedJobs rolls back all failed jobs.
func (jm *JobManager) RollbackFailedJobs(ctx context.Context) error {
	failed := jm.GetFailedResults()

	jm.logger.Infof("Rolling back %d failed jobs", len(failed))

	for i := len(failed) - 1; i >= 0; i-- { // Rollback in reverse order
		result := failed[i]

		if err := checkContextCancellation(ctx, "context cancelled during rollback"); err != nil {
			return err
		}

		jm.logger.Infof("Rolling back job: %s", result.Job.Name())

		err := result.Job.Rollback(ctx)
		if err != nil {
			jm.logger.Errorf("Failed to rollback job %s: %v", result.Job.Name(), err)
			// Continue with other rollbacks
		} else {
			jm.updateJobStatus(result.Job.ID(), JobStatusRolledBack)
			jm.logger.Infof("Successfully rolled back job: %s", result.Job.Name())
		}
	}

	return nil
}

// Clear clears all jobs and results.
func (jm *JobManager) Clear() {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.jobs = make([]Job, 0)
	jm.results = make([]JobResult, 0)
}

// GetStatistics returns job execution statistics.
func (jm *JobManager) GetStatistics() map[string]any {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	stats := map[string]any{
		"total_jobs":     len(jm.jobs),
		"total_results":  len(jm.results),
		"completed":      0,
		"failed":         0,
		"total_duration": time.Duration(0),
	}

	var totalDuration time.Duration
	for _, result := range jm.results {
		totalDuration += result.Duration
		switch result.Status {
		case JobStatusCompleted:
			if v, ok := stats["completed"].(int); ok {
				stats["completed"] = v + 1
			}
		case JobStatusFailed:
			if v, ok := stats["failed"].(int); ok {
				stats["failed"] = v + 1
			}
		case JobStatusRolledBack:
			if v, ok := stats["rolled_back"].(int); ok {
				stats["rolled_back"] = v + 1
			}
		case JobStatusPending, JobStatusRunning:
			// These statuses don't contribute to final stats
		}
	}

	stats["total_duration"] = totalDuration

	stats["average_duration"] = time.Duration(0)
	if len(jm.results) > 0 {
		stats["average_duration"] = totalDuration / time.Duration(len(jm.results))
	}

	return stats
}

// displayJobResults displays job execution results to the user.
func displayJobResults(results []JobResult) {
	for _, result := range results {
		switch result.Status {
		case JobStatusCompleted:
			fmt.Printf(
				"%s %s completed successfully\n",
				successStyle.Render("✅"),
				result.Job.Name(),
			)
		case JobStatusFailed:
			fmt.Printf(
				"%s %s failed: %v\n",
				errorStyle.Render("❌"),
				result.Job.Name(),
				result.Error,
			)
		case JobStatusRolledBack:
			fmt.Printf(
				"%s %s rolled back\n",
				infoStyle.Render("↩️"),
				result.Job.Name(),
			)
		case JobStatusPending, JobStatusRunning:
			// These statuses are not displayed in final results
		}
	}
}
