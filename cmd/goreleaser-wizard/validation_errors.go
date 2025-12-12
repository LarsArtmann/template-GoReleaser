package main

import (
	"fmt"
)

// DomainError represents domain error for validation results
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Context string `json:"context,omitempty"`
	Cause   error  `json:"cause,omitempty"`
}

// Error implements error interface
func (de *DomainError) Error() string {
	if de.Context != "" {
		return fmt.Sprintf("[%s] %s (context: %s)", de.Code, de.Message, de.Context)
	}
	return fmt.Sprintf("[%s] %s", de.Code, de.Message)
}

// Unwrap returns the underlying cause
func (de *DomainError) Unwrap() error {
	return de.Cause
}

// NewValidationError creates a validation domain error
func NewValidationError(code, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewSystemError creates a system domain error
func NewSystemError(code, message, details string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
		Cause:   cause,
	}
}

// NewTemplateError creates a template domain error
func NewTemplateError(code, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WithContext adds context to the error
func (de *DomainError) WithContext(context string) *DomainError {
	return &DomainError{
		Code:    de.Code,
		Message: de.Message,
		Details: de.Details,
		Context: context,
		Cause:   de.Cause,
	}
}
