package errors

import (
	"errors"
	"fmt"
	"runtime"
)

// ErrorCode represents typed error codes.
type ErrorCode string

const (
	// Validation Errors.
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrInvalidProject   ErrorCode = "INVALID_PROJECT"
	ErrInvalidBinary    ErrorCode = "INVALID_BINARY"
	ErrInvalidConfig    ErrorCode = "INVALID_CONFIG"
	ErrInvalidTemplate  ErrorCode = "INVALID_TEMPLATE"
	ErrInvalidField     ErrorCode = "INVALID_FIELD"
	ErrInvalidOperation ErrorCode = "INVALID_OPERATION"

	// Configuration Errors.
	ErrConfigGeneration          ErrorCode = "CONFIG_GENERATION"
	ErrConfigFound               ErrorCode = "CONFIG_FOUND"
	ErrConfigNotFound            ErrorCode = "CONFIG_NOT_FOUND"
	ErrConfigInvalid             ErrorCode = "CONFIG_INVALID"
	ErrInvalidMainPath           ErrorCode = "INVALID_MAIN_PATH"
	ErrInvalidProjectDescription ErrorCode = "INVALID_PROJECT_DESCRIPTION"
	ErrInvalidDockerImage        ErrorCode = "INVALID_DOCKER_IMAGE"
	ErrInvalidDockerRegistry     ErrorCode = "INVALID_DOCKER_REGISTRY"
	ErrInvalidVersion            ErrorCode = "INVALID_VERSION"
	ErrInvalidGitBranch          ErrorCode = "INVALID_GIT_BRANCH"
	ErrInvalidGitTag             ErrorCode = "INVALID_GIT_TAG"
	ErrInvalidBuildTag           ErrorCode = "INVALID_BUILD_TAG"
	ErrInvalidPort               ErrorCode = "INVALID_PORT"
	ErrConfigPermission          ErrorCode = "CONFIG_PERMISSION"

	// Template Errors.
	ErrTemplateNotFound  ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrTemplateParsing   ErrorCode = "TEMPLATE_PARSING"
	ErrTemplateRendering ErrorCode = "TEMPLATE_RENDERING"
	ErrTemplateExecution ErrorCode = "TEMPLATE_EXECUTION"

	// File System Errors.
	ErrFileOperation  ErrorCode = "FILE_OPERATION"
	ErrFileNotFound   ErrorCode = "FILE_NOT_FOUND"
	ErrFilePermission ErrorCode = "FILE_PERMISSION"
	ErrFileCorrupted  ErrorCode = "FILE_CORRUPTED"
	ErrDirNotFound    ErrorCode = "DIR_NOT_FOUND"
	ErrDirPermission  ErrorCode = "DIR_PERMISSION"

	// Git Errors.
	ErrGitOperation  ErrorCode = "GIT_OPERATION"
	ErrGitNotFound   ErrorCode = "GIT_NOT_FOUND"
	ErrGitPermission ErrorCode = "GIT_PERMISSION"
	ErrGitCommand    ErrorCode = "GIT_COMMAND"

	// Dependency Errors.
	ErrDependencyMissing  ErrorCode = "DEPENDENCY_MISSING"
	ErrDependencyVersion  ErrorCode = "DEPENDENCY_VERSION"
	ErrDependencyConflict ErrorCode = "DEPENDENCY_CONFLICT"

	// Job Execution Errors.
	ErrJobExecution ErrorCode = "JOB_EXECUTION"
	ErrJobTimeout   ErrorCode = "JOB_TIMEOUT"
	ErrJobCancelled ErrorCode = "JOB_CANCELLED"
	ErrJobFailed    ErrorCode = "JOB_FAILED"

	// Workflow Errors.
	ErrWorkflowExecution ErrorCode = "WORKFLOW_EXECUTION"
	ErrWorkflowTimeout   ErrorCode = "WORKFLOW_TIMEOUT"
	ErrWorkflowCancelled ErrorCode = "WORKFLOW_CANCELLED"
	ErrWorkflowFailed    ErrorCode = "WORKFLOW_FAILED"

	// Network Errors.
	ErrNetworkError       ErrorCode = "NETWORK_ERROR"
	ErrNetworkTimeout     ErrorCode = "NETWORK_TIMEOUT"
	ErrNetworkUnavailable ErrorCode = "NETWORK_UNAVAILABLE"

	// Permission Errors.
	ErrPermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrForbidden        ErrorCode = "FORBIDDEN"

	// System Errors.
	ErrSystemError ErrorCode = "SYSTEM_ERROR"
	ErrMemoryError ErrorCode = "MEMORY_ERROR"
	ErrDiskError   ErrorCode = "DISK_ERROR"
	ErrCPUError    ErrorCode = "CPU_ERROR"

	// Unknown Error.
	ErrUnknown ErrorCode = "UNKNOWN"
)

// ErrorLevel represents error severity levels.
type ErrorLevel string

const (
	ErrorLevelCritical ErrorLevel = "critical"
	ErrorLevelHigh     ErrorLevel = "high"
	ErrorLevelMedium   ErrorLevel = "medium"
	ErrorLevelLow      ErrorLevel = "low"
)

// DomainError represents a structured domain error.
type DomainError struct {
	Code      ErrorCode  `json:"code"`
	Message   string     `json:"message"`
	Details   string     `json:"details,omitempty"`
	Context   string     `json:"context,omitempty"`
	Field     string     `json:"field,omitempty"`
	Level     ErrorLevel `json:"level"`
	Cause     error      `json:"cause,omitempty"`
	Retryable bool       `json:"retryable"`
	File      string     `json:"file,omitempty"`
	Line      int        `json:"line,omitempty"`
	Function  string     `json:"function,omitempty"`
}

// Error implements the error interface.
func (de *DomainError) Error() string {
	if de.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", de.Code, de.Message, de.Details)
	}
	return fmt.Sprintf("[%s] %s", de.Code, de.Message)
}

// Unwrap returns the underlying cause.
func (de *DomainError) Unwrap() error {
	return de.Cause
}

// WithContext adds context to the error.
func (de *DomainError) WithContext(context string) *DomainError {
	de.Context = context
	return de
}

// WithField adds field information to the error.
func (de *DomainError) WithField(field string) *DomainError {
	de.Field = field
	return de
}

// WithLevel sets the error level.
func (de *DomainError) WithLevel(level ErrorLevel) *DomainError {
	de.Level = level
	return de
}

// WithRetryable sets if the error is retryable.
func (de *DomainError) WithRetryable(retryable bool) *DomainError {
	de.Retryable = retryable
	return de
}

// WithCause adds the underlying cause.
func (de *DomainError) WithCause(cause error) *DomainError {
	de.Cause = cause
	return de
}

// WithCaller adds caller information.
func (de *DomainError) WithCaller() *DomainError {
	if pc, file, line, ok := runtime.Caller(1); ok {
		de.Function = runtime.FuncForPC(pc).Name()
		de.File = file
		de.Line = line
	}
	return de
}

// WithSuggestion adds a suggestion to the error.
func (de *DomainError) WithSuggestion(suggestion string) *DomainError {
	de.Context = suggestion
	return de
}

// NewValidationError creates a new validation error.
func NewValidationError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelMedium,
		Retryable: false,
	}
	return err.WithCaller()
}

// NewConfigError creates a new configuration error.
func NewConfigError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelHigh,
		Retryable: false,
	}
	return err.WithCaller()
}

// NewFileError creates a new file system error.
func NewFileError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelMedium,
		Retryable: true,
	}
	return err.WithCaller()
}

// NewGitError creates a new git error.
func NewGitError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelMedium,
		Retryable: true,
	}
	return err.WithCaller()
}

// NewJobError creates a new job execution error.
func NewJobError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelHigh,
		Retryable: false,
	}
	return err.WithCaller()
}

// NewWorkflowError creates a new workflow error.
func NewWorkflowError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelHigh,
		Retryable: false,
	}
	return err.WithCaller()
}

// NewNetworkError creates a new network error.
func NewNetworkError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelMedium,
		Retryable: true,
	}
	return err.WithCaller()
}

// NewPermissionError creates a new permission error.
func NewPermissionError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelHigh,
		Retryable: false,
	}
	return err.WithCaller()
}

// NewSystemError creates a new system error.
func NewSystemError(code ErrorCode, message, details string) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Details:   details,
		Level:     ErrorLevelCritical,
		Retryable: false,
	}
	return err.WithCaller()
}

// WrapError wraps an existing error with domain context.
func WrapError(err error, code ErrorCode, message string) *DomainError {
	domainErr := &DomainError{
		Code:      code,
		Message:   message,
		Details:   err.Error(),
		Cause:     err,
		Level:     ErrorLevelMedium,
		Retryable: false,
	}
	return domainErr.WithCaller()
}

// IsRetryable checks if an error is retryable.
func IsRetryable(err error) bool {
	domainErr := &DomainError{}
	if errors.As(err, &domainErr) {
		return domainErr.Retryable
	}
	return false
}

// GetErrorCode extracts the error code from an error.
func GetErrorCode(err error) ErrorCode {
	domainErr := &DomainError{}
	if errors.As(err, &domainErr) {
		return domainErr.Code
	}
	return ErrUnknown
}

// GetErrorLevel extracts the error level from an error.
func GetErrorLevel(err error) ErrorLevel {
	domainErr := &DomainError{}
	if errors.As(err, &domainErr) {
		if domainErr.Level != "" {
			return domainErr.Level
		}
	}
	return ErrorLevelMedium
}
