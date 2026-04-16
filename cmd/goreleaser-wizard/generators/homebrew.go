package generators

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// HomebrewGenerator handles Homebrew formula generation.
type HomebrewGenerator struct {
	templateData *HomebrewTemplateData
	logger       Logger
}

// HomebrewTemplateData contains data for Homebrew formula template.
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

// NewHomebrewGenerator creates a new Homebrew generator.
func NewHomebrewGenerator(config *domain.SafeProjectConfig, logger Logger) *HomebrewGenerator {
	return &HomebrewGenerator{
		templateData: createHomebrewTemplateData(config),
		logger:       logger,
	}
}

// createHomebrewTemplateData creates template data from project config.
func createHomebrewTemplateData(config *domain.SafeProjectConfig) *HomebrewTemplateData {
	// Convert project name to formula name (CamelCase)
	formulaName := toCamelCase(config.ProjectName)

	// Default values
	description := config.ProjectDescription
	if description == "" {
		description = config.ProjectName + " is a Go application"
	}

	return &HomebrewTemplateData{
		FormulaName:        formulaName,
		ProjectName:        config.ProjectName,
		ProjectDescription: description,
		Homepage:           "https://github.com/user/" + config.ProjectName,
		ArchiveURL: fmt.Sprintf(
			"https://github.com/user/%s/archive/v{{version}}.tar.gz",
			config.ProjectName,
		),
		Checksum:   "{{sha256}}", // Will be filled by GoReleaser
		License:    "MIT",
		MainPath:   config.MainPath,
		LDFlags:    "-s -w -X main.version={{version}}",
		BinaryName: config.BinaryName,
		TestOutput: config.BinaryName + " version",
		Service:    config.ProjectType == domain.ProjectTypeWebAPI,
	}
}

// toCamelCase converts snake_case or kebab-case to CamelCase.
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

// splitWords splits a string into words by common separators.
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

// Generate generates Homebrew formula.
func (g *HomebrewGenerator) Generate(ctx context.Context) error {
	output, err := GenerateTemplate(
		ctx,
		g.logger,
		"homebrew",
		templates.HomebrewTemplate,
		"Generating Homebrew formula",
		g.templateData,
	)
	if err != nil {
		return err
	}

	// Ensure homebrew directory exists
	formulaDir := "homebrew"
	if err := os.MkdirAll(formulaDir, directoryPermission); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to create homebrew directory",
			err.Error(),
		).WithCause(err)
	}

	// Write formula file
	formulaPath := fmt.Sprintf("%s/%s.rb", formulaDir, g.templateData.FormulaName)
	if err := WriteFile(formulaPath, output, filePermission); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to write Homebrew formula",
			err.Error(),
		).WithCause(err)
	}

	g.logger.Info("Homebrew formula generated successfully", "path", formulaPath)

	return nil
}

// ValidateTemplate validates Homebrew template.
func (g *HomebrewGenerator) ValidateTemplate() error {
	return ValidateTemplate("homebrew", templates.HomebrewTemplate)
}

// GeneratePreview generates a preview without writing to file.
func (g *HomebrewGenerator) GeneratePreview(ctx context.Context) (string, error) {
	return GeneratePreview(
		ctx,
		g.logger,
		"homebrew",
		templates.HomebrewTemplate,
		"Generating Homebrew formula preview",
		g.templateData,
	)
}

// Rollback removes generated Homebrew formula.
func (g *HomebrewGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back Homebrew formula generation")

	if err := CheckContext(ctx); err != nil {
		return err
	}

	// Remove generated formula
	formulaPath := fmt.Sprintf("homebrew/%s.rb", g.templateData.FormulaName)

	err := removeGeneratedFile(
		g.logger,
		formulaPath,
		"Failed to remove generated formula",
		"Removed generated formula",
	)
	if err != nil {
		return err
	}

	removeEmptyDirectories(g.logger, []string{"homebrew"})

	return nil
}

// UpdateConfig updates generator configuration.
func (g *HomebrewGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = createHomebrewTemplateData(config)
}
