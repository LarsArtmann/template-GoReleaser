package generators

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/git"
)

// GoReleaserGenerator handles GoReleaser configuration generation.
type GoReleaserGenerator struct {
	templateData *types.GoReleaserTemplateData
	logger       Logger
}

// Logger interface for dependency injection.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewGoReleaserGenerator creates a new GoReleaser generator.
func NewGoReleaserGenerator(config *domain.SafeProjectConfig, logger Logger) *GoReleaserGenerator {
	return &GoReleaserGenerator{
		templateData: types.NewGoReleaserTemplateData(config),
		logger:       logger,
	}
}

// Generate generates the GoReleaser configuration.
func (g *GoReleaserGenerator) Generate(ctx context.Context) error {
	g.logger.Info("Generating GoReleaser configuration")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Create template with custom functions
	tmpl := template.New("goreleaser").Funcs(template.FuncMap{
		"incpatch": git.IncPatchVersion,
	})

	tmpl, err := tmpl.Parse(templates.GoReleaserTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse GoReleaser template",
			err.Error(),
		).WithCause(err)
	}

	// Generate template data with git information
	data, err := g.prepareTemplateData(ctx)
	if err != nil {
		return err
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute GoReleaser template",
			err.Error(),
		).WithCause(err)
	}

	// Create backup if file exists
	if err := g.createBackup(".goreleaser.yaml"); err != nil {
		return err
	}

	// Write configuration file
	if err := os.WriteFile(".goreleaser.yaml", output.Bytes(), 0o644); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to write GoReleaser config",
			err.Error(),
		).WithCause(err)
	}

	g.logger.Info("GoReleaser configuration generated successfully")

	return nil
}

// createBackup creates a backup of existing file.
func (g *GoReleaserGenerator) createBackup(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		backupPath := filename + ".backup"
		err := os.Rename(filename, backupPath)
		if err != nil {
			return errors.NewFileError(
				errors.ErrFileOperation,
				"Failed to create backup",
				err.Error(),
			).WithCause(err)
		}

		g.logger.Info("Created backup file", "backup", backupPath)
	}

	return nil
}

// prepareTemplateData prepares complete template data including git information.
func (g *GoReleaserGenerator) prepareTemplateData(
	ctx context.Context,
) (*types.GoReleaserTemplateData, error) {
	// Get version information from git
	versionInfo, err := git.GetVersionInfo(ctx)
	if err != nil {
		g.logger.Warn("Failed to get git version info, using defaults", "error", err)

		versionInfo = &git.VersionInfo{
			Version:    "v0.1.0",
			CommitHash: "unknown",
		}
	}

	// Update template data with git information
	data := *g.templateData // Copy to avoid mutation
	data.Version = versionInfo.Version
	data.Tag = versionInfo.Version
	data.Major = git.GetMajorVersion(versionInfo.Version)
	data.Date = git.GetCurrentDate()
	data.FullCommit = versionInfo.CommitHash

	// Update environment variables
	if data.Env == nil {
		data.Env = make(map[string]string)
	}

	data.Env["GITHUB_OWNER"] = versionInfo.Owner
	data.Env["GITHUB_REPO"] = versionInfo.Repo

	return &data, nil
}

// ValidateTemplate validates the GoReleaser template.
func (g *GoReleaserGenerator) ValidateTemplate() error {
	tmpl := template.New("goreleaser").Funcs(template.FuncMap{
		"incpatch": git.IncPatchVersion,
	})

	_, err := tmpl.Parse(templates.GoReleaserTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrInvalidTemplate,
			"GoReleaser template validation failed",
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// GeneratePreview generates a preview without writing to file.
func (g *GoReleaserGenerator) GeneratePreview(ctx context.Context) (string, error) {
	g.logger.Debug("Generating GoReleaser configuration preview")

	// Check context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Create template
	tmpl := template.New("goreleaser").Funcs(template.FuncMap{
		"incpatch": git.IncPatchVersion,
	})

	tmpl, err := tmpl.Parse(templates.GoReleaserTemplate)
	if err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse GoReleaser template",
			err.Error(),
		).WithCause(err)
	}

	// Prepare template data
	data, err := g.prepareTemplateData(ctx)
	if err != nil {
		return "", err
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute GoReleaser template preview",
			err.Error(),
		).WithCause(err)
	}

	return output.String(), nil
}
