package main

import "github.com/spf13/cobra"

// getBoolFlag reads a bool flag from the command.
// Flags are registered during init(), so errors here indicate a programming bug.
func getBoolFlag(cmd *cobra.Command, name string) bool {
	value, _ := cmd.Flags().GetBool(name)

	return value
}

// getStringFlag reads a string flag from the command.
// Flags are registered during init(), so errors here indicate a programming bug.
func getStringFlag(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name)

	return value
}
