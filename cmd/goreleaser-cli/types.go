package main

import (
	"time"
)

// Config represents the complete configuration structure for the CLI
type Config struct {
	License LicenseConfig `yaml:"license" json:"license" mapstructure:"license"`
	Author  AuthorConfig  `yaml:"author" json:"author" mapstructure:"author"`
	Project ProjectConfig `yaml:"project" json:"project" mapstructure:"project"`
	CLI     CLIConfig     `yaml:"cli" json:"cli" mapstructure:"cli"`
}

// LicenseConfig represents license-related configuration
type LicenseConfig struct {
	Type         string `yaml:"type" json:"type" mapstructure:"type"`
	Year         int    `yaml:"year,omitempty" json:"year,omitempty" mapstructure:"year"`
	TemplatePath string `yaml:"template_path,omitempty" json:"template_path,omitempty" mapstructure:"template_path"`
}

// AuthorConfig represents author/copyright information
type AuthorConfig struct {
	Name  string `yaml:"name" json:"name" mapstructure:"name"`
	Email string `yaml:"email,omitempty" json:"email,omitempty" mapstructure:"email"`
	URL   string `yaml:"url,omitempty" json:"url,omitempty" mapstructure:"url"`
}

// ProjectConfig represents project-specific settings
type ProjectConfig struct {
	Name        string `yaml:"name,omitempty" json:"name,omitempty" mapstructure:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty" mapstructure:"description"`
	Version     string `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version"`
	Repository  string `yaml:"repository,omitempty" json:"repository,omitempty" mapstructure:"repository"`
}

// CLIConfig represents CLI behavior settings
type CLIConfig struct {
	Verbose bool `yaml:"verbose" json:"verbose" mapstructure:"verbose"`
	Colors  bool `yaml:"colors" json:"colors" mapstructure:"colors"`
	Timeout int  `yaml:"timeout,omitempty" json:"timeout,omitempty" mapstructure:"timeout"` // in seconds
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Category string        `json:"category"`
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Message  string        `json:"message,omitempty"`
	Duration time.Duration `json:"duration"`
	Details  []string      `json:"details,omitempty"`
}

// ValidationReport represents a collection of validation results
type ValidationReport struct {
	Timestamp   time.Time          `json:"timestamp"`
	TotalChecks int                `json:"total_checks"`
	Passed      int                `json:"passed"`
	Failed      int                `json:"failed"`
	Warnings    int                `json:"warnings"`
	Duration    time.Duration      `json:"duration"`
	Results     []ValidationResult `json:"results"`
	Environment EnvironmentInfo    `json:"environment"`
}

// EnvironmentInfo represents information about the execution environment
type EnvironmentInfo struct {
	OS         string            `json:"os"`
	Arch       string            `json:"arch"`
	GoVersion  string            `json:"go_version"`
	CLIVersion string            `json:"cli_version"`
	WorkingDir string            `json:"working_dir"`
	GitCommit  string            `json:"git_commit,omitempty"`
	GitBranch  string            `json:"git_branch,omitempty"`
	EnvVars    map[string]string `json:"env_vars,omitempty"`
}

// GoReleaserConfig represents configuration for GoReleaser validation
type GoReleaserConfig struct {
	Free GoReleaserSettings `yaml:"free" json:"free" mapstructure:"free"`
	Pro  GoReleaserSettings `yaml:"pro" json:"pro" mapstructure:"pro"`
}

// GoReleaserSettings represents settings for a specific GoReleaser version
type GoReleaserSettings struct {
	ConfigFile string            `yaml:"config_file" json:"config_file" mapstructure:"config_file"`
	EnvVars    map[string]string `yaml:"env_vars,omitempty" json:"env_vars,omitempty" mapstructure:"env_vars"`
	Required   bool              `yaml:"required" json:"required" mapstructure:"required"`
	Timeout    int               `yaml:"timeout,omitempty" json:"timeout,omitempty" mapstructure:"timeout"`
}

// LicenseTemplate represents a license template
type LicenseTemplate struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Path        string            `json:"path"`
	Variables   []string          `json:"variables"`
	Description string            `json:"description"`
	Content     string            `json:"content,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// VerificationResult represents the result of a verification operation
type VerificationResult struct {
	Category    string        `json:"category"`
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration"`
	Command     string        `json:"command,omitempty"`
	Output      string        `json:"output,omitempty"`
	ExitCode    int           `json:"exit_code,omitempty"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// VerificationReport represents a collection of verification results
type VerificationReport struct {
	Timestamp   time.Time            `json:"timestamp"`
	TotalChecks int                  `json:"total_checks"`
	Passed      int                  `json:"passed"`
	Failed      int                  `json:"failed"`
	Warnings    int                  `json:"warnings"`
	Duration    time.Duration        `json:"duration"`
	Results     []VerificationResult `json:"results"`
	Environment EnvironmentInfo      `json:"environment"`
}

// Status constants for validation and verification results
const (
	StatusPassed  = "passed"
	StatusFailed  = "failed"
	StatusWarning = "warning"
	StatusSkipped = "skipped"
)

// Category constants for organization
const (
	CategoryConfig      = "config"
	CategoryStructure   = "structure"
	CategoryEnvironment = "environment"
	CategoryBuild       = "build"
	CategoryTest        = "test"
	CategoryFormat      = "format"
	CategoryLint        = "lint"
	CategoryDependency  = "dependency"
	CategorySecurity    = "security"
	CategoryLicense     = "license"
)

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		License: LicenseConfig{
			Type: "MIT",
			Year: time.Now().Year(),
		},
		Author: AuthorConfig{
			Name: "",
		},
		Project: ProjectConfig{},
		CLI: CLIConfig{
			Verbose: false,
			Colors:  true,
			Timeout: 300, // 5 minutes
		},
	}
}

// Validate validates the configuration and returns any errors
func (c *Config) Validate() []string {
	var errors []string

	// Validate license
	if c.License.Type == "" {
		errors = append(errors, "license type is required")
	}

	// Validate author
	if c.Author.Name == "" {
		errors = append(errors, "author name is required")
	}

	// Validate CLI settings
	if c.CLI.Timeout <= 0 {
		errors = append(errors, "CLI timeout must be positive")
	}

	return errors
}

// IsComplete returns true if all required fields are filled
func (c *Config) IsComplete() bool {
	return c.License.Type != "" &&
		c.Author.Name != "" &&
		c.Project.Name != ""
}
