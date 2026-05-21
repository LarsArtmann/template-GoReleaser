package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/config"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// Type alias for backward compatibility during migration
// TODO: Remove after migration complete.
type ProjectConfig = domain.SafeProjectConfig

// LoggerAdapter adapts charmbracelet/log to domain.Logger interface.
type LoggerAdapter struct {
	logger *log.Logger
}

func (la *LoggerAdapter) Debug(msg string, args ...any) {
	la.logger.Debug(msg, args...)
}

func (la *LoggerAdapter) Info(msg string, args ...any) {
	la.logger.Info(msg, args...)
}

func (la *LoggerAdapter) Warn(msg string, args ...any) {
	la.logger.Warn(msg, args...)
}

func (la *LoggerAdapter) Error(msg string, args ...any) {
	la.logger.Error(msg, args...)
}

func (la *LoggerAdapter) Fatal(msg string, args ...any) {
	la.logger.Fatal(msg, args...)
}

func (la *LoggerAdapter) DebugContext(ctx context.Context, msg string, args ...any) {
	la.logger.Debug(msg, args...)
}

func (la *LoggerAdapter) InfoContext(ctx context.Context, msg string, args ...any) {
	la.logger.Info(msg, args...)
}

func (la *LoggerAdapter) WarnContext(ctx context.Context, msg string, args ...any) {
	la.logger.Warn(msg, args...)
}

func (la *LoggerAdapter) ErrorContext(ctx context.Context, msg string, args ...any) {
	la.logger.Error(msg, args...)
}

func (la *LoggerAdapter) WithField(key string, value any) domain.Logger {
	return la // Simplified - doesn't add field
}

func (la *LoggerAdapter) WithFields(fields map[string]any) domain.Logger {
	return la // Simplified - doesn't add fields
}

func (la *LoggerAdapter) WithError(err error) domain.Logger {
	return la // Simplified - doesn't add error
}

var (
	// Build-time variables set by GoReleaser.
	version        = "dev"
	commit         = "none"
	date           = "unknown"
	builtBy        = "unknown"
	gitDescription = ""
	gitState       = ""

	cfgFile string
)

// Domain logger for dependency injection.
var appLogger domain.Logger

// Style definitions.
var titleStyle, successStyle, errorStyle, infoStyle lipgloss.Style

func init() {
	// Create a logger adapter to satisfy domain.Logger interface
	appLogger = &LoggerAdapter{logger: log.New(os.Stderr)}

	// Initialize styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)
	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	if config.GetManager().IsDebug() {
		if la, ok := appLogger.(*LoggerAdapter); ok {
			la.logger.SetLevel(log.DebugLevel)
		}
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "goreleaser-wizard",
	Short: "Interactive setup wizard for GoReleaser",
	Long: `GoReleaser Wizard is an interactive CLI tool that helps you create
perfect GoReleaser configurations for your Go projects.

It guides you through the configuration process with smart defaults
and best practices, generating both .goreleaser.yaml and GitHub Actions
workflows tailored to your project's needs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Set up panic recovery using domain error handling
	defer recoverFromPanic("command execution")

	err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version),
		fang.WithCommit(commit),
	)
	if err != nil {
		displayError(err)
	}
}

// recoverFromPanic provides graceful panic recovery using domain types.
func recoverFromPanic(context string) {
	if r := recover(); r != nil {
		log.Error("Panic recovered", "context", context, "panic", r)

		err := domain.NewSystemError(
			domain.ErrTemplateExecutionFailed,
			"Unexpected error occurred",
			fmt.Sprintf("The wizard encountered an unexpected problem: %v", r),
			fmt.Errorf("panic: %v", r),
		).WithContext(context)

		displayError(err)

		os.Exit(1)
	}
}

// displayError displays errors using domain error handling.
func displayError(err error) {
	if err == nil {
		return
	}

	// Convert to domain error if not already
	var domainErr *domain.DomainError
	if !errors.As(err, &domainErr) {
		domainErr = domain.NewSystemError(
			domain.ErrFileWriteFailed,
			"Unexpected error",
			err.Error(),
			err,
		)
	}

	// Display structured error information
	fmt.Println()
	fmt.Println(errorStyle.Render("❌ Error: " + domainErr.Message))

	if domainErr.Details != "" {
		fmt.Println(infoStyle.Render("Details: " + domainErr.Details))
	}

	if domainErr.Context != "" {
		fmt.Println(infoStyle.Render("Context: " + domainErr.Context))
	}

	suggestion := domainErr.GetRecoverySuggestion()
	if suggestion != "" {
		suggestStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true)
		fmt.Println(suggestStyle.Render("💡 Suggestion: " + suggestion))
	}

	// Log the full error for debugging
	log.Error(
		"Domain error",
		"code", domainErr.Code,
		"message", domainErr.Message,
		"details", domainErr.Details,
		"context", domainErr.Context,
	)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goreleaser-wizard.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

	// Bind flags to config manager
	config.GetManager().RegisterFlags(rootCmd.PersistentFlags())

	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set up panic recovery for config initialization
	defer recoverFromPanic("config initialization")

	err := config.GetManager().Load(cfgFile)
	if err != nil {
		// Only log if it's not a "file not found" error for optional config
		if cfgFile != "" {
			displayError(domain.NewSystemError(
				domain.ErrFileReadFailed,
				"Failed to load configuration",
				err.Error(),
				err,
			).WithContext(cfgFile))
		}
	}

	if config.GetManager().IsDebug() {
		log.Info("Configuration loaded")
	}
}

// validateFileExists validates file existence using domain error types.
func validateFileExists(path string, requireDir bool) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.NewSystemError(
				domain.ErrFileNotFound,
				"File not found",
				fmt.Sprintf("File %s does not exist, require_dir=%t", path, requireDir),
				err,
			).WithContext(path)
		}

		return domain.NewSystemError(
			domain.ErrFileReadFailed,
			"File access error",
			fmt.Sprintf("Cannot access %s, require_dir=%t", path, requireDir),
			err,
		).WithContext(path)
	}

	if requireDir && !info.IsDir() {
		return domain.NewValidationError(
			domain.ErrInvalidCharacters,
			"Expected directory",
			fmt.Sprintf("%s is not a directory, require_dir=%t", path, requireDir),
		).WithContext(path)
	}

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoReleaser Wizard %s\n", version)
		fmt.Printf("  Build Date: %s\n", date)
		fmt.Printf("  Git Commit: %s\n", commit)
		fmt.Printf("  Built By: %s\n", builtBy)

		if gitState != "" {
			fmt.Printf("  Git State: %s\n", gitState)
		}

		if gitDescription != "" {
			fmt.Printf("  Git Summary: %s\n", gitDescription)
		}
	},
}

// Placeholder command definitions - TODO: Implement actual functionality.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize GoReleaser configuration",
	Run:   runInitWizard,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate GoReleaser configuration",
	Run:   runGenerate,
}

// HandleError provides a simple error handling function for tests and CLI usage.
func HandleError(err error) {
	if err != nil {
		displayError(err)
	}
}

// addCommonFlags adds common flags used by both init and generate commands.
func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("force", false, "force overwrite existing configuration")
	cmd.Flags().String("project-name", "", "override project name")
	cmd.Flags().String("main-path", "", "override main.go path")
	cmd.Flags().String("binary-name", "", "override binary name")
	cmd.Flags().String("project-type", "", "override project type (CLI Application, Library)")
}

// configureInitCommand adds flags specific to the init command.
func configureInitCommand(cmd *cobra.Command) {
	addCommonFlags(cmd)
	cmd.Flags().Bool("interactive", true, "run in interactive mode (default true)")
}

// configureGenerateCommand adds flags specific to the generate command.
func configureGenerateCommand(cmd *cobra.Command) {
	addCommonFlags(cmd)
	cmd.Flags().
		Bool("config-only", false, "generate only GoReleaser configuration (no GitHub Actions)")
}

func main() {
	// Set up global panic recovery
	defer recoverFromPanic("main")

	Execute()
}
