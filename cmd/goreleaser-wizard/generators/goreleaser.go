package generators

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
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

	if err := CheckContext(ctx); err != nil {
		return err
	}

	// Create template with custom functions
	tmpl := newTemplateWithFuncs("goreleaser", template.FuncMap{
		"incpatch": git.IncPatchVersion,
	})

	// Parse template
	tmpl, err := parseTemplateWithError(tmpl, templates.GoReleaserTemplate, "GoReleaser")
	if err != nil {
		return err
	}

	// Execute template
	output, err := executeTemplateWithError(tmpl, g.templateData, "GoReleaser")
	if err != nil {
		return err
	}

	// Create backup if file exists
	if err := g.createBackup(goreleaserConfigFilename); err != nil {
		return err
	}

	// Write configuration file
	if err := WriteFile(goreleaserConfigFilename, output, filePermission); err != nil {
		return WrapFileError(err, "Failed to write GoReleaser config")
	}

	g.logger.Info("GoReleaser configuration generated successfully")

	return nil
}

// goreleaserConfigFilename is the file the GoReleaser generator writes.
const goreleaserConfigFilename = ".goreleaser.yaml"

// createBackup creates a backup of existing file.
func (g *GoReleaserGenerator) createBackup(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		backupPath := filename + ".backup"

		err := os.Rename(filename, backupPath)
		if err != nil {
			return fmt.Errorf("failed to create backup for filename=%q: %w", filename, err)
		}

		g.logger.Info("Created backup file", "backup", backupPath)
	}

	return nil
}

// ValidateTemplate validates the GoReleaser template.
func (g *GoReleaserGenerator) ValidateTemplate() error {
	tmpl := template.New("goreleaser").Funcs(template.FuncMap{
		"incpatch": git.IncPatchVersion,
	})

	_, err := tmpl.Parse(templates.GoReleaserTemplate)
	if err != nil {
		return domain.NewConfigError(
			domain.ErrInvalidTemplate,
			"GoReleaser template validation failed",
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// GeneratePreview generates a preview without writing to file.
func (g *GoReleaserGenerator) GeneratePreview(ctx context.Context) (string, error) {
	g.logger.Debug("Generating GoReleaser configuration preview")

	if err := CheckContext(ctx); err != nil {
		return "", err
	}

	output, err := GeneratePreview(
		ctx,
		g.logger,
		"goreleaser",
		templates.GoReleaserTemplate,
		"Generating GoReleaser configuration preview",
		g.templateData,
	)
	if err != nil {
		return "", err
	}

	return output, nil
}
