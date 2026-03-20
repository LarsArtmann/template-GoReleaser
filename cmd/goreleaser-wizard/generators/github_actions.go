package generators

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// GitHubActionsGenerator handles GitHub Actions workflow generation.
type GitHubActionsGenerator struct {
	templateData *types.GitHubActionsTemplateData
	logger       Logger
}

// NewGitHubActionsGenerator creates a new GitHub Actions generator.
func NewGitHubActionsGenerator(
	config *domain.SafeProjectConfig,
	logger Logger,
) *GitHubActionsGenerator {
	return &GitHubActionsGenerator{
		templateData: types.NewGitHubActionsTemplateData(config),
		logger:       logger,
	}
}

// Generate generates GitHub Actions workflow.
func (g *GitHubActionsGenerator) Generate(ctx context.Context) error {
	g.logger.Info("Generating GitHub Actions workflow")

	// Check if GitHub Actions generation is enabled
	if !g.shouldGenerate() {
		g.logger.Info("GitHub Actions generation is disabled, skipping")

		return nil
	}

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Create template with custom delimiters to avoid GitHub Actions syntax conflict
	tmpl := template.New("github-actions").Delims("[[", "]]")

	tmpl, err := tmpl.Parse(templates.GitHubActionsTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse GitHub Actions template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute GitHub Actions template",
			err.Error(),
		).WithCause(err)
	}

	// Ensure .github/workflows directory exists
	workflowDir := filepath.Join(".github", "workflows")
	if err := os.MkdirAll(workflowDir, directoryPermission); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to create workflow directory",
			err.Error(),
		).WithCause(err)
	}

	// Write workflow file
	workflowPath := filepath.Join(workflowDir, "release.yml")
	if err := os.WriteFile(workflowPath, output.Bytes(), filePermission); err != nil {
		return errors.NewFileError(
			errors.ErrFileOperation,
			"Failed to write GitHub Actions workflow",
			err.Error(),
		).WithCause(err)
	}

	g.logger.Info("GitHub Actions workflow generated successfully", "path", workflowPath)

	return nil
}

// shouldGenerate checks if GitHub Actions should be generated.
func (g *GitHubActionsGenerator) shouldGenerate() bool {
	// Check if we have triggers defined
	return len(g.templateData.Triggers) > 0
}

// ValidateTemplate validates GitHub Actions template.
func (g *GitHubActionsGenerator) ValidateTemplate() error {
	tmpl := template.New("github-actions").Delims("[[", "]]")

	_, err := tmpl.Parse(templates.GitHubActionsTemplate)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrInvalidTemplate,
			"GitHub Actions template validation failed",
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// GeneratePreview generates a preview without writing to file.
func (g *GitHubActionsGenerator) GeneratePreview(ctx context.Context) (string, error) {
	g.logger.Debug("Generating GitHub Actions workflow preview")

	// Check if should generate
	if !g.shouldGenerate() {
		return "// GitHub Actions generation is disabled", nil
	}

	// Check context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Create template with custom delimiters
	tmpl := template.New("github-actions").Delims("[[", "]]")

	tmpl, err := tmpl.Parse(templates.GitHubActionsTemplate)
	if err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse GitHub Actions template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute GitHub Actions template preview",
			err.Error(),
		).WithCause(err)
	}

	return output.String(), nil
}

// Rollback removes generated GitHub Actions workflow.
func (g *GitHubActionsGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back GitHub Actions workflow generation")

	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Remove generated workflow
	workflowPath := filepath.Join(".github", "workflows", "release.yml")

	err := removeGeneratedFile(
		g.logger,
		workflowPath,
		"Failed to remove generated workflow",
		"Removed generated workflow",
	)
	if err != nil {
		return err
	}

	removeEmptyDirectories(g.logger, []string{
		filepath.Join(".github", "workflows"),
		".github",
	})

	return nil
}

// UpdateConfig updates the generator configuration.
func (g *GitHubActionsGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = types.NewGitHubActionsTemplateData(config)
}

// GetWorkflowPath returns the path where workflow will be generated.
func (g *GitHubActionsGenerator) GetWorkflowPath() string {
	return filepath.Join(".github", "workflows", "release.yml")
}
