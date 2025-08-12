package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-GoReleaser/internal/container"
	"github.com/charmbracelet/fang"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Build-time variables set by GoReleaser
	version        = "dev"
	commit         = "none"
	date           = "unknown"
	builtBy        = "unknown"
	gitDescription = ""
	gitState       = ""
)

var (
	cfgFile     string
	diContainer *container.Container
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goreleaser-cli",
	Short: "A powerful CLI tool built with GoReleaser template",
	Long: `GoReleaser CLI is a batteries-included command line tool that provides
validation, verification, and license management capabilities.

Built with Cobra for CLI structure, Fang for enhanced features, 
and Viper for configuration management.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Use fang.Execute for enhanced CLI experience with styling
	ctx := context.Background()
	if err := fang.Execute(ctx, rootCmd,
		fang.WithVersion(version),
		fang.WithCommit(commit),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initContainer)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goreleaser-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Enable completion commands
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".goreleaser-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goreleaser-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// initContainer initializes the dependency injection container
func initContainer() {
	diContainer = container.NewContainer()
}

// GetContainer returns the DI container for use in commands
func GetContainer() *do.Injector {
	if diContainer == nil {
		initContainer()
	}
	return diContainer.GetInjector()
}

func main() {
	// Check for health check flag
	if len(os.Args) > 1 && os.Args[1] == "--health" {
		// Simple health check - return 0 for healthy
		fmt.Println("healthy")
		os.Exit(0)
	}

	// Ensure graceful shutdown of DI container
	defer func() {
		if diContainer != nil {
			_ = diContainer.Shutdown()
		}
	}()

	Execute()
}
