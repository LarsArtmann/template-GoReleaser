package generators

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// CheckContext returns an error if the context is cancelled.
func CheckContext(ctx context.Context) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	return nil
}

// WriteFile writes data to a file with proper error handling.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	err := os.WriteFile(path, data, perm)
	if err != nil {
		return WrapFileError(err, "failed to write file")
	}

	return nil
}

// WrapFileError wraps a file operation error with context.
func WrapFileError(err error, message string) error {
	return domain.NewFileError(
		domain.ErrFileOperation,
		message,
		err.Error(),
	).WithCause(err)
}

// newTemplateWithFuncs creates a template with custom functions.
func newTemplateWithFuncs(name string, funcs template.FuncMap) *template.Template {
	return template.New(name).Funcs(funcs)
}

// parseTemplateWithError parses a template and wraps errors with context.
func parseTemplateWithError(
	tmpl *template.Template,
	content, templateType string,
) (*template.Template, error) {
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse %s template (content length %d): %w",
			templateType,
			len(content),
			err,
		)
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
		return nil, domain.NewConfigError(
			domain.ErrTemplateRendering,
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
		return fmt.Errorf(
			"templateName=%s validation failed (left=%q, right=%q, content length %d): %w",
			templateName,
			left,
			right,
			len(templateContent),
			err,
		)
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

// parseAndExecuteTemplate parses and executes a template, wrapping errors with context.
func parseAndExecuteTemplate(
	tmpl *template.Template,
	templateName, templateContent string,
	templateData any,
) ([]byte, error) {
	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse %s template (content length %d): %w",
			templateName,
			len(templateContent),
			err,
		)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, templateData); err != nil {
		return nil, fmt.Errorf(
			"failed to execute %s template (content length %d): %w",
			templateName,
			len(templateContent),
			err,
		)
	}

	return output.Bytes(), nil
}

// GeneratePreviewWithDelims generates a preview with custom delimiters.
func GeneratePreviewWithDelims(
	ctx context.Context,
	logger Logger,
	templateName, left, right, templateContent, logPrefix string,
	templateData any,
) (string, error) {
	logger.Debug(logPrefix)

	if ctx.Err() != nil {
		return "", fmt.Errorf(
			"context cancelled for templateName=%s (left=%q, right=%q, logPrefix=%q, content length %d): %w",
			templateName,
			left,
			right,
			logPrefix,
			len(templateContent),
			ctx.Err(),
		)
	}

	tmpl := template.New(templateName).Delims(left, right)

	result, err := parseAndExecuteTemplate(tmpl, templateName, templateContent, templateData)
	if err != nil {
		return "", fmt.Errorf(
			"failed for templateName=%s (left=%q, right=%q, logPrefix=%q, content length %d): %w",
			templateName,
			left,
			right,
			logPrefix,
			len(templateContent),
			err,
		)
	}

	return string(result), nil
}

// GenerateTemplate generates and executes a template, returning the output.
func GenerateTemplate(
	ctx context.Context,
	logger Logger,
	templateName, templateContent, logPrefix string,
	templateData any,
) ([]byte, error) {
	logger.Info(logPrefix)

	if ctx.Err() != nil {
		return nil, fmt.Errorf(
			"context cancelled for templateName=%s (logPrefix=%q, content length %d): %w",
			templateName,
			logPrefix,
			len(templateContent),
			ctx.Err(),
		)
	}

	tmpl := template.New(templateName)

	return parseAndExecuteTemplate(tmpl, templateName, templateContent, templateData)
}

// removeGeneratedFile removes a generated file if it exists.
func removeGeneratedFile(logger Logger, filePath, errorMsg, successMsg string) error {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return domain.NewFileError(
				domain.ErrFileOperation,
				fmt.Sprintf("%s (successMsg=%q)", errorMsg, successMsg),
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
