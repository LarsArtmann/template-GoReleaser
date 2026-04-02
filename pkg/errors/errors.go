// Package errors provides custom error types and error handling utilities.
//
// This package defines application-specific error types that can be used
// throughout the codebase for consistent error handling and reporting.
package errors

// Error is the foundation for all custom errors in this project.
type Error struct {
	Message string
	Code    string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// New creates a new Error with the given message and code.
func New(message, code string) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}
