package generators

import (
	"context"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// dockerfileBaseImageScratch is the base for pure-Go static binaries.
const dockerfileBaseImageScratch = "scratch"

// dockerfileBaseImageAlpine is the base for CGO binaries, which are
// dynamically linked and need a C runtime.
const dockerfileBaseImageAlpine = "alpine:latest"

// webAPIExposePort is the conventional port exposed for Web API projects.
const webAPIExposePort = "8080"

// DockerfileGenerator handles Dockerfile generation.
type DockerfileGenerator struct {
	templateData *DockerfileTemplateData
	logger       Logger
}

// DockerfileTemplateData is the typed data for the prebuilt-pattern
// Dockerfile template.
type DockerfileTemplateData struct {
	ProjectName string
	BinaryName  string
	BaseImage   string
	ExposePort  string
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
		ProjectName: config.ProjectName,
		BinaryName:  config.BinaryName,
		BaseImage:   dockerfileBaseImageScratch,
		ExposePort:  "",
	}

	// CGO binaries are dynamically linked and need a C runtime.
	if config.CGOStatus.IsEnabled() {
		data.BaseImage = dockerfileBaseImageAlpine
	}

	if config.ProjectType == domain.ProjectTypeWebAPI {
		data.ExposePort = webAPIExposePort
	}

	return data
}

// Generate generates Dockerfile.
func (g *DockerfileGenerator) Generate(ctx context.Context) error {
	output, err := GenerateTemplate(
		ctx,
		g.logger,
		"dockerfile",
		templates.DockerfileTemplate,
		"Generating Dockerfile",
		g.templateData,
	)
	if err != nil {
		return err
	}

	// Write Dockerfile
	err = WriteFile("Dockerfile", output, filePermission)
	if err != nil {
		return WrapFileError(err, "Failed to write Dockerfile")
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
	return GeneratePreview(
		ctx,
		g.logger,
		"dockerfile",
		templates.DockerfileTemplate,
		"Generating Dockerfile preview",
		g.templateData,
	)
}

// Rollback removes generated Dockerfile.
func (g *DockerfileGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back Dockerfile generation")

	if err := CheckContext(ctx); err != nil {
		return err
	}

	// Remove generated Dockerfile
	return removeGeneratedFile(
		g.logger,
		"Dockerfile",
		"Failed to remove generated Dockerfile",
		"Removed generated Dockerfile",
	)
}

// UpdateConfig updates generator configuration.
func (g *DockerfileGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = createDockerfileTemplateData(config)
}
