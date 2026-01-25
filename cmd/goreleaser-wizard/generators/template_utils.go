package generators

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// ValidateTemplate validates a template by parsing it.
func ValidateTemplate(templateName, templateContent string) error {
	tmpl := template.New(templateName)

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
func GeneratePreview(ctx context.Context, logger Logger, templateName, templateContent, logPrefix string, templateData any) (string, error) {
	logger.Debug(logPrefix)

	// Check context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Create and parse template
	tmpl := template.New(templateName)
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
func GenerateTemplate(ctx context.Context, logger Logger, templateName, templateContent, logPrefix string, templateData any) ([]byte, error) {
	logger.Info(logPrefix)

	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
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
		if err := os.Remove(filePath); err != nil {
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
