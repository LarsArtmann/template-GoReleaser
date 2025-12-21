package generators

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// HomebrewGenerator handles Homebrew formula generation
type HomebrewGenerator struct {
	templateData *HomebrewTemplateData
	logger       Logger
}

// HomebrewTemplateData contains data for Homebrew formula template
type HomebrewTemplateData struct {
	FormulaName        string
	ProjectName        string
	ProjectDescription string
	Homepage           string
	ArchiveURL         string
	Checksum           string
	License            string
	MainPath           string
	LDFlags            string
	BinaryName         string
	TestOutput         string
	Service            bool
}

// NewHomebrewGenerator creates a new Homebrew generator
func NewHomebrewGenerator(config *domain.SafeProjectConfig, logger Logger) *HomebrewGenerator {
	return &HomebrewGenerator{
		templateData: createHomebrewTemplateData(config),
		logger:       logger,
	}
}

// createHomebrewTemplateData creates template data from project config
func createHomebrewTemplateData(config *domain.SafeProjectConfig) *HomebrewTemplateData {
	// Convert project name to formula name (CamelCase)
	formulaName := toCamelCase(config.ProjectName)

	// Default values
	description := config.ProjectDescription
	if description == "" {
		description = fmt.Sprintf("%s is a Go application", config.ProjectName)
	}

	return &HomebrewTemplateData{
		FormulaName:        formulaName,
		ProjectName:        config.ProjectName,
		ProjectDescription: description,
		Homepage:           "https://github.com/user/" + config.ProjectName,
		ArchiveURL:         fmt.Sprintf("https://github.com/user/%s/archive/v{{version}}.tar.gz", config.ProjectName),
		Checksum:           "{{sha256}}", // Will be filled by GoReleaser
		License:            "MIT",
		MainPath:           config.MainPath,
		LDFlags:            "-s -w -X main.version={{version}}",
		BinaryName:         config.BinaryName,
		TestOutput:         fmt.Sprintf("%s version", config.BinaryName),
		Service:            config.ProjectType == domain.ProjectTypeWebAPI,
	}
}

// toCamelCase converts snake_case or kebab-case to CamelCase
func toCamelCase(s string) string {
	// Handle kebab-case and snake_case
	words := splitWords(s)

	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

// splitWords splits a string into words by common separators
func splitWords(s string) []string {
	var words []string
	current := ""

	for _, r := range s {
		if r == '-' || r == '_' || r == ' ' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}

	if current != "" {
		words = append(words, current)
	}

	return words
}

// Generate generates Homebrew formula
func (g *HomebrewGenerator) Generate(ctx context.Context) error {
	g.logger.Info("Generating Homebrew formula")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Create template
	tmpl := template.New("homebrew")

	tmpl, err := tmpl.Parse(templates.HomebrewTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse Homebrew template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute Homebrew template",
			err.Error(),
		).WithCause(err)
	}

	// Ensure homebrew directory exists
	formulaDir := "homebrew"
	if err := os.MkdirAll(formulaDir, 0o755); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to create homebrew directory",
			err.Error(),
		).WithCause(err)
	}

	// Write formula file
	formulaPath := fmt.Sprintf("%s/%s.rb", formulaDir, g.templateData.FormulaName)
	if err := os.WriteFile(formulaPath, output.Bytes(), 0o644); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to write Homebrew formula",
			err.Error(),
		).WithCause(err)
	}

	g.logger.Info("Homebrew formula generated successfully", "path", formulaPath)
	return nil
}

// ValidateTemplate validates Homebrew template
func (g *HomebrewGenerator) ValidateTemplate() error {
	tmpl := template.New("homebrew")

	_, err := tmpl.Parse(templates.HomebrewTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrInvalidTemplate,
			"Homebrew template validation failed",
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// GeneratePreview generates a preview without writing to file
func (g *HomebrewGenerator) GeneratePreview(ctx context.Context) (string, error) {
	g.logger.Debug("Generating Homebrew formula preview")

	// Check context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Create template
	tmpl := template.New("homebrew")

	tmpl, err := tmpl.Parse(templates.HomebrewTemplate)
	if err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse Homebrew template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute Homebrew template preview",
			err.Error(),
		).WithCause(err)
	}

	return output.String(), nil
}

// Rollback removes generated Homebrew formula
func (g *HomebrewGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back Homebrew formula generation")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated formula
	formulaPath := fmt.Sprintf("homebrew/%s.rb", g.templateData.FormulaName)
	if _, err := os.Stat(formulaPath); err == nil {
		err := os.Remove(formulaPath)
		if err != nil {
			return errors.NewFileError(
				errors.ErrFileOperation,
				"Failed to remove generated formula",
				err.Error(),
			).WithCause(err)
		}
		g.logger.Info("Removed generated formula", "path", formulaPath)
	}

	// Try to remove homebrew directory if empty
	if files, err := os.ReadDir("homebrew"); err == nil && len(files) == 0 {
		os.Remove("homebrew")
		g.logger.Info("Removed empty homebrew directory")
	}

	return nil
}

// UpdateConfig updates generator configuration
func (g *HomebrewGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = createHomebrewTemplateData(config)
}
