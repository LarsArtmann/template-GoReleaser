package main

const exitCodeWarnings = 2

// ValidationResults holds all validation results.
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

// GetExitCode returns appropriate exit code.
func (vr *ValidationResults) GetExitCode() int {
	if len(vr.Errors) > 0 {
		return 1
	}

	if len(vr.Warnings) > 0 {
		return exitCodeWarnings
	}

	return 0
}
