package types

import "time"

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
