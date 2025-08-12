package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goreleaser-cli",
	Short: "A simple CLI tool built with GoReleaser template",
	Long:  `A simple CLI tool demonstrating GoReleaser capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoReleaser CLI %s\n", version)
		fmt.Printf("Built: %s\n", date)
		fmt.Printf("Commit: %s\n", commit)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add version command
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		fmt.Printf("Built by: %s\n", builtBy)
	},
}

func main() {
	Execute()
}
