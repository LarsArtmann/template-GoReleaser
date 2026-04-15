package generators

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// newTemplateWithFuncs creates a template with custom functions.
func newTemplateWithFuncs(name string, funcs template.FuncMap) *template.Template {
	return template.New(name).Funcs(funcs)
}

// newTemplateWithDelims creates a template with custom delimiters.
func newTemplateWithDelims(name, left, right string) *template.Template {
	return template.New(name).Delims(left, right)
}

// parseAndExecute parses a template and executes it with the given data.
func parseAndExecute(
	tmpl *template.Template,
	content string,
	data any,
) ([]byte, error) {
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// parseTemplateWithError parses a template and wraps errors with context.
func parseTemplateWithError(
	tmpl *template.Template,
	content, templateType string,
) (*template.Template, error) {
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return nil, errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse "+templateType+" template",
			err.Error(),
		).WithCause(err)
	}

	return tmpl, nil
}

// executeTemplateWithError executes a template and wraps errors with context.
func executeTemplateWithError(
	tmpl *template.Template,
	data any,
	templateType string,
) ([]byte, error) {
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		return nil, errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute "+templateType+" template",
			err.Error(),
		).WithCause(err)
	}

	return output.Bytes(), nil
}

// ValidateTemplate validates a template by parsing it.
func ValidateTemplate(templateName, templateContent string) error {
	return ValidateTemplateWithDelims(templateName, "", "", templateContent)
}

// ValidateTemplateWithDelims validates a template with custom delimiters.
func ValidateTemplateWithDelims(templateName, left, right, templateContent string) error {
	tmpl := template.New(templateName).Delims(left, right)

	_, err := tmpl.Parse(templateContent)
	if err != nil {
		return errors.NewConfigError(
			errors.ErrInvalidTemplate,
			templateName+" template validation failed",
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// GeneratePreview generates a preview of a template without writing to file.
func GeneratePreview(
	ctx context.Context,
	logger Logger,
	templateName, templateContent, logPrefix string,
	templateData any,
) (string, error) {
	return GeneratePreviewWithDelims(
		ctx, logger, templateName, "", "", templateContent, logPrefix, templateData,
	)
}

// GeneratePreviewWithDelims generates a preview with custom delimiters.
func GeneratePreviewWithDelims(
	ctx context.Context,
	logger Logger,
	templateName, left, right, templateContent, logPrefix string,
	templateData any,
) (string, error) {
	logger.Debug(logPrefix)

	// Check context cancellation
	if ctx.Err() != nil {
		return "", fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	// Create and parse template with custom delimiters
	tmpl := template.New(templateName).Delims(left, right)

	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse "+templateName+" template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, templateData); err != nil {
		return "", errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute "+templateName+" template preview",
			err.Error(),
		).WithCause(err)
	}

	return output.String(), nil
}

// GenerateTemplate generates and executes a template, returning the output.
func GenerateTemplate(
	ctx context.Context,
	logger Logger,
	templateName, templateContent, logPrefix string,
	templateData any,
) ([]byte, error) {
	logger.Info(logPrefix)

	// Check context cancellation
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	// Create and parse template
	tmpl := template.New(templateName)

	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return nil, errors.NewConfigError(
			errors.ErrTemplateParsing,
			"Failed to parse "+templateName+" template",
			err.Error(),
		).WithCause(err)
	}

	// Execute template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, templateData); err != nil {
		return nil, errors.NewConfigError(
			errors.ErrTemplateRendering,
			"Failed to execute "+templateName+" template",
			err.Error(),
		).WithCause(err)
	}

	return output.Bytes(), nil
}

// removeGeneratedFile removes a generated file if it exists.
func removeGeneratedFile(logger Logger, filePath, errorMsg, successMsg string) error {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return errors.NewFileError(
				errors.ErrFileOperation,
				errorMsg,
				err.Error(),
			).WithCause(err)
		}

		logger.Info(successMsg, "path", filePath)
	}

	return nil
}

// removeEmptyDirectory removes a directory if it's empty.
func removeEmptyDirectory(logger Logger, dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err == nil && len(files) == 0 {
		err := os.Remove(dirPath)
		if err == nil {
			logger.Info("Removed empty directory", "path", dirPath)
		}
	}
}

// removeEmptyDirectories removes empty directories in a parent chain.
func removeEmptyDirectories(logger Logger, parentDirs []string) {
	for _, dir := range parentDirs {
		removeEmptyDirectory(logger, dir)
	}
}
