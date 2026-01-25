package generators

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// DockerfileGenerator handles Dockerfile generation.
type DockerfileGenerator struct {
	templateData *DockerfileTemplateData
	logger       Logger
}

// DockerfileTemplateData contains data for Dockerfile template.
type DockerfileTemplateData struct {
	ProjectName   string
	GoVersion     string
	AlpineVersion string
	CGOEnabled    string
	TargetOS      string
	TargetArch    string
	LDFlags       string
	MainPath      string
	ConfigFiles   []string
	ExposePort    string
	EnvVars       []EnvVar
	HealthCheck   *HealthCheck
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Key   string
	Value string
}

// HealthCheck represents container health check configuration.
type HealthCheck struct {
	Interval    string
	Timeout     string
	StartPeriod string
	Retries     int
	Commands    []string
}

// NewDockerfileGenerator creates a new Dockerfile generator.
func NewDockerfileGenerator(config *domain.SafeProjectConfig, logger Logger) *DockerfileGenerator {
	return &DockerfileGenerator{
		templateData: createDockerfileTemplateData(config),
		logger:       logger,
	}
}

// createDockerfileTemplateData creates template data from project config.
func createDockerfileTemplateData(config *domain.SafeProjectConfig) *DockerfileTemplateData {
	data := &DockerfileTemplateData{
		ProjectName:   config.ProjectName,
		GoVersion:     "1.21", // Default Go version
		AlpineVersion: "latest",
		CGOEnabled:    "false",
		TargetOS:      "linux",
		TargetArch:    "amd64",
		LDFlags:       "-s -w",
		MainPath:      config.MainPath,
		ConfigFiles:   []string{},
		ExposePort:    "",
		EnvVars:       []EnvVar{},
		HealthCheck:   nil,
	}

	// Set CGO based on configuration
	if config.CGOStatus.IsEnabled() {
		data.CGOEnabled = "true"
	}

	// Add common config files
	if config.ProjectType == domain.ProjectTypeWebAPI {
		data.ConfigFiles = []string{"config.yaml", ".env.example"}
		data.ExposePort = "8080"
		data.EnvVars = []EnvVar{
			{Key: "PORT", Value: "8080"},
			{Key: "LOG_LEVEL", Value: "info"},
		}
		data.HealthCheck = &HealthCheck{
			Interval:    "30s",
			Timeout:     "10s",
			StartPeriod: "5s",
			Retries:     3,
			Commands:    []string{"CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"},
		}
	}

	return data
}

// Generate generates Dockerfile.
func (g *DockerfileGenerator) Generate(ctx context.Context) error {
	g.logger.Info("Generating Dockerfile")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Create template
	tmpl := template.New("dockerfile")

	tmpl, err := tmpl.Parse(templates.DockerfileTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse Dockerfile template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute Dockerfile template",
			err.Error(),
		).WithCause(err)
	}

	// Write Dockerfile
	if err := os.WriteFile("Dockerfile", output.Bytes(), 0o644); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to write Dockerfile",
			err.Error(),
		).WithCause(err)
	}

	g.logger.Info("Dockerfile generated successfully")
	return nil
}

// ValidateTemplate validates Dockerfile template.
func (g *DockerfileGenerator) ValidateTemplate() error {
	return ValidateTemplate("dockerfile", templates.DockerfileTemplate)
}

// GeneratePreview generates a preview without writing to file.
func (g *DockerfileGenerator) GeneratePreview(ctx context.Context) (string, error) {
	return GeneratePreview(ctx, g.logger, "dockerfile", templates.DockerfileTemplate, "Generating Dockerfile preview", g.templateData)
}

// Rollback removes generated Dockerfile.
func (g *DockerfileGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back Dockerfile generation")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated Dockerfile
	if _, err := os.Stat("Dockerfile"); err == nil {
		err := os.Remove("Dockerfile")
		if err != nil {
			return errors.NewFileError(
				errors.ErrFileOperation,
				"Failed to remove generated Dockerfile",
				err.Error(),
			).WithCause(err)
		}
		g.logger.Info("Removed generated Dockerfile")
	}

	return nil
}

// UpdateConfig updates generator configuration.
func (g *DockerfileGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = createDockerfileTemplateData(config)
}
