package domain

import (
	"context"
	"fmt"
	"time"
)

// ProgressReporter represents progress reporting interface.
type ProgressReporter interface {
	Start(total int64, message string) error
	Update(current int64) error
	Finish(message string) error
	Error(err error) error
	Increment() error
}

// Progress represents progress information.
type Progress struct {
	Current int64         `json:"current"`
	Total   int64         `json:"total"`
	Percent float64       `json:"percent"`
	Message string        `json:"message"`
	Elapsed time.Duration `json:"elapsed"`
	ETA     time.Duration `json:"eta,omitempty"`
}

// Event represents a domain event.
type Event interface {
	ID() string
	Type() EventType
	Data() any
	OccurredAt() time.Time
	AggregateID() string
	Version() int
}

// EventType represents the type of domain event.
type EventType string

const (
	// Configuration events.
	EventTypeConfigCreated   EventType = "CONFIG_CREATED"
	EventTypeConfigUpdated   EventType = "CONFIG_UPDATED"
	EventTypeConfigValidated EventType = "CONFIG_VALIDATED"
	EventTypeConfigGenerated EventType = "CONFIG_GENERATED"

	// Job execution events.
	EventTypeJobCreated   EventType = "JOB_CREATED"
	EventTypeJobStarted   EventType = "JOB_STARTED"
	EventTypeJobCompleted EventType = "JOB_COMPLETED"
	EventTypeJobFailed    EventType = "JOB_FAILED"
	EventTypeJobCancelled EventType = "JOB_CANCELLED"
	EventTypeJobRetried   EventType = "JOB_RETRIED"

	// Workflow execution events.
	EventTypeWorkflowCreated   EventType = "WORKFLOW_CREATED"
	EventTypeWorkflowStarted   EventType = "WORKFLOW_STARTED"
	EventTypeWorkflowCompleted EventType = "WORKFLOW_COMPLETED"
	EventTypeWorkflowFailed    EventType = "WORKFLOW_FAILED"
	EventTypeWorkflowCancelled EventType = "WORKFLOW_CANCELLED"

	// Template generation events.
	EventTypeTemplateGenerated EventType = "TEMPLATE_GENERATED"
	EventTypeTemplateRendered  EventType = "TEMPLATE_RENDERED"
	EventTypeTemplateValidated EventType = "TEMPLATE_VALIDATED"

	// File system events.
	EventTypeFileCreated EventType = "FILE_CREATED"
	EventTypeFileUpdated EventType = "FILE_UPDATED"
	EventTypeFileDeleted EventType = "FILE_DELETED"
	EventTypeDirCreated  EventType = "DIR_CREATED"
	EventTypeDirDeleted  EventType = "DIR_DELETED"

	// Git operations events.
	EventTypeGitOperation EventType = "GIT_OPERATION"
	EventTypeGitPush      EventType = "GIT_PUSH"
	EventTypeGitCommit    EventType = "GIT_COMMIT"
	EventTypeGitTag       EventType = "GIT_TAG"

	// External integration events.
	EventTypeGitHubAPI           EventType = "GITHUB_API"
	EventTypeDockerOperation     EventType = "DOCKER_OPERATION"
	EventTypeGoReleaserOperation EventType = "GORELEASER_OPERATION"
)

// EventPublisher represents event publishing interface.
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	PublishAsync(ctx context.Context, event Event)
	PublishBatch(ctx context.Context, events []Event) error
}

// EventHandler represents event handling interface.
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(eventType EventType) bool
}

// EventBus represents event bus interface.
type EventBus interface {
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
	Publish(ctx context.Context, event Event) error
	Close() error
}

// EventBase provides a base implementation for domain events.
type EventBase struct {
	id          IdID
	eventType   EventType
	data        any
	occurredAt  time.Time
	aggregateID AggregateID
	version     int
}

// NewEventBase creates a new base event.
func NewEventBase(eventType EventType, data any, aggregateID AggregateID) *EventBase {
	return &EventBase{
		id:          NewIdID(generateEventID()),
		eventType:   eventType,
		data:        data,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
		version:     1,
	}
}

func (e *EventBase) ID() string {
	return e.id.Get()
}

func (e *EventBase) Type() EventType {
	return e.eventType
}

func (e *EventBase) Data() any {
	return e.data
}

func (e *EventBase) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *EventBase) AggregateID() string {
	return e.aggregateID.Get()
}

func (e *EventBase) Version() int {
	return e.version
}

// WithData updates event data.
func (e *EventBase) WithData(data any) *EventBase {
	e.data = data

	return e
}

// WithVersion updates event version.
func (e *EventBase) WithVersion(version int) *EventBase {
	e.version = version

	return e
}

// generateEventID generates a unique event ID.
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// ConfigCreatedEvent represents a configuration creation event.
type ConfigCreatedEvent struct {
	*EventBase

	ConfigID    ConfigID `json:"config_id"`
	ConfigType  string   `json:"config_type"`
	CreatedBy   string   `json:"created_by"`
	ProjectName string   `json:"project_name"`
}

// NewConfigCreatedEvent creates a new config created event.
func NewConfigCreatedEvent(
	configID ConfigID, configType, createdBy, projectName string,
) *ConfigCreatedEvent {
	return &ConfigCreatedEvent{
		EventBase: NewEventBase(EventTypeConfigCreated, map[string]any{
			"config_id":    configID,
			"config_type":  configType,
			"created_by":   createdBy,
			"project_name": projectName,
		}, NewAggregateID(configID.Get())),
		ConfigID:    configID,
		ConfigType:  configType,
		CreatedBy:   createdBy,
		ProjectName: projectName,
	}
}

// JobStartedEvent represents a job started event.
type JobStartedEvent struct {
	*EventBase

	JobID     JobID     `json:"job_id"`
	JobName   string    `json:"job_name"`
	StartedBy string    `json:"started_by"`
	StartTime time.Time `json:"start_time"`
}

// NewJobStartedEvent creates a new job started event.
func NewJobStartedEvent(jobID JobID, jobName, startedBy string) *JobStartedEvent {
	startTime := time.Now()

	return &JobStartedEvent{
		EventBase: NewEventBase(EventTypeJobStarted, map[string]any{
			"job_id":     jobID,
			"job_name":   jobName,
			"started_by": startedBy,
			"start_time": startTime,
		}, NewAggregateID(jobID.Get())),
		JobID:     jobID,
		JobName:   jobName,
		StartedBy: startedBy,
		StartTime: startTime,
	}
}

// WorkflowCompletedEvent represents a workflow completed event.
type WorkflowCompletedEvent struct {
	*EventBase

	WorkflowID     WorkflowID    `json:"workflow_id"`
	WorkflowName   string        `json:"workflow_name"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	TotalJobs      int           `json:"total_jobs"`
	SuccessfulJobs int           `json:"successful_jobs"`
	FailedJobs     int           `json:"failed_jobs"`
}

// NewWorkflowCompletedEvent creates a new workflow completed event.
func NewWorkflowCompletedEvent(
	workflowID WorkflowID, workflowName string,
	startTime, endTime time.Time,
	totalJobs, successfulJobs, failedJobs int,
) *WorkflowCompletedEvent {
	duration := endTime.Sub(startTime)

	return &WorkflowCompletedEvent{
		EventBase: NewEventBase(EventTypeWorkflowCompleted, map[string]any{
			"workflow_id":     workflowID,
			"workflow_name":   workflowName,
			"start_time":      startTime,
			"end_time":        endTime,
			"duration":        duration,
			"total_jobs":      totalJobs,
			"successful_jobs": successfulJobs,
			"failed_jobs":     failedJobs,
		}, NewAggregateID(workflowID.Get())),
		WorkflowID:     workflowID,
		WorkflowName:   workflowName,
		StartTime:      startTime,
		EndTime:        endTime,
		Duration:       duration,
		TotalJobs:      totalJobs,
		SuccessfulJobs: successfulJobs,
		FailedJobs:     failedJobs,
	}
}
