package main

import (
	"fmt"
	"os"

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

	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goreleaser-wizard.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

	// Bind flags to viper
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".goreleaser-wizard" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goreleaser-wizard")
	}

	viper.SetEnvPrefix("GORELEASER_WIZARD")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetBool("debug") {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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

func main() {
	Execute()
}