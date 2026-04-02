package main

import (
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// Re-export domain error types for backward compatibility.
// These are type aliases to the domain package's error types.

type DomainError = domain.DomainError

// Error constructors - delegate to domain package.
func NewValidationError(code, message, details string) *DomainError {
	return domain.NewValidationError(domain.ErrorCode(code), message, details)
}

func NewSystemError(code, message, details string, cause error) *DomainError {
	return domain.NewSystemError(domain.ErrorCode(code), message, details, cause)
}

func NewTemplateError(code, message, details string) *DomainError {
	return domain.NewTemplateError(domain.ErrorCode(code), message, details)
}
