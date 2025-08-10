package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print detailed version information including build metadata.

This command displays the version, commit hash, build date, 
Go version, and platform information.`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("Version:      %s\n", version)
	fmt.Printf("Commit:       %s\n", commit)
	fmt.Printf("Built:        %s\n", date)
	fmt.Printf("Built by:     %s\n", builtBy)
	fmt.Printf("Go version:   %s\n", runtime.Version())
	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)

	if gitDescription != "" {
		fmt.Printf("Git describe: %s\n", gitDescription)
	}

	if gitState != "" {
		fmt.Printf("Git state:    %s\n", gitState)
	}
}
