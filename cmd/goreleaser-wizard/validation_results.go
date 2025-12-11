package main

import (
	"fmt"
)

// ValidationResults holds all validation results
type ValidationResults struct {
	ConfigExists    bool
	ConfigValid     bool
	ActionsExists   bool
	ActionsValid    bool
	ProjectValid    bool
	GoReleaserFound bool
	Errors          []*DomainError
	Warnings        []*DomainError
	Recommendations []string
}

// GetExitCode returns appropriate exit code
func (vr *ValidationResults) GetExitCode() int {
	if len(vr.Errors) > 0 {
		return 1
	}
	if len(vr.Warnings) > 0 {
		return 2
	}
	return 0
}

// DomainError represents domain error for validation results
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Context string `json:"context,omitempty"`
}

// Error implements error interface
func (de *DomainError) Error() string {
	if de.Context != "" {
		return fmt.Sprintf("[%s] %s (context: %s)", de.Code, de.Message, de.Context)
	}
	return fmt.Sprintf("[%s] %s", de.Code, de.Message)
}