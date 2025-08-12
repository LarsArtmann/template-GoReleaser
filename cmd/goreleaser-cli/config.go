package main

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/template-GoReleaser/internal/services"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage configuration settings for the CLI tool.
	
This command allows you to view, set, and manage configuration
values that are used across different commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfig(cmd, args)
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration values and their sources.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfigShow(cmd, args)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a configuration value in the config file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runConfigSet(cmd, args)
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Create a default configuration file with common settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfigInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)

	// Add flags
	configInitCmd.Flags().BoolP("force", "f", false, "Overwrite existing config file")
}

func runConfig(cmd *cobra.Command, args []string) {
	fmt.Println("⚙️  Configuration Management")
	fmt.Println("\nUse one of the following subcommands:")
	fmt.Println("  show  - Show current configuration")
	fmt.Println("  set   - Set a configuration value")
	fmt.Println("  init  - Initialize configuration file")
	fmt.Println("\nFor more help: goreleaser-cli config --help")
}

func runConfigShow(cmd *cobra.Command, args []string) {
	fmt.Println("📋 Current Configuration:")

	// Get config service from DI container
	injector := GetContainer()
	configService := do.MustInvoke[services.ConfigService](injector)

	// Load configuration using the service
	config, err := configService.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		fmt.Println("💡 Run 'goreleaser-cli config init' to create one")
		return
	}

	fmt.Println("\n🔧 Configuration values:")
	fmt.Printf("  📄 License: %s\n", config.License.Type)
	fmt.Printf("  👤 Author: %s <%s>\n", config.Author.Name, config.Author.Email)
	fmt.Printf("  📦 Project: %s\n", config.Project.Name)
	if config.Project.Description != "" {
		fmt.Printf("  📋 Description: %s\n", config.Project.Description)
	}
	fmt.Printf("  🛠️  CLI: verbose=%t, colors=%t\n",
		config.CLI.Verbose, config.CLI.Colors)

	// Validate configuration using the service
	if err := configService.ValidateConfig(config); err != nil {
		fmt.Printf("\n⚠️  Configuration validation error: %v\n", err)
	} else {
		fmt.Println("\n✅ Configuration is valid")
	}

	// Show environment variables that would override
	fmt.Println("\n🌍 Environment variable overrides:")
	envVars := []string{
		"LICENSE_TYPE",
		"COPYRIGHT_HOLDER",
		"AUTHOR_NAME",
		"PROJECT_AUTHOR",
	}

	foundEnvVars := false
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			fmt.Printf("  %s = %s\n", envVar, value)
			foundEnvVars = true
		}
	}

	if !foundEnvVars {
		fmt.Println("  (no relevant environment variables set)")
	}
}

func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	fmt.Printf("🔧 Setting configuration: %s = %s\n", key, value)

	// Get config service from DI container
	injector := GetContainer()
	configService := do.MustInvoke[services.ConfigService](injector)

	// Load current configuration
	config, err := configService.LoadConfig()
	if err != nil {
		// If config doesn't exist, create a new one
		config, err = configService.InitConfig()
		if err != nil {
			fmt.Printf("❌ Cannot initialize configuration: %v\n", err)
			os.Exit(1)
		}
	}

	// Set the value based on key path
	// This is a simplified implementation - could be enhanced with path traversal
	switch key {
	case "license.type":
		config.License.Type = value
	case "author.name":
		config.Author.Name = value
	case "author.email":
		config.Author.Email = value
	case "project.name":
		config.Project.Name = value
	case "project.description":
		config.Project.Description = value
	case "cli.verbose":
		config.CLI.Verbose = (value == "true")
	case "cli.colors":
		config.CLI.Colors = (value == "true")
	default:
		fmt.Printf("❌ Unknown configuration key: %s\n", key)
		fmt.Println("💡 Supported keys: license.type, author.name, author.email, project.name, project.description, cli.verbose, cli.colors")
		os.Exit(1)
	}

	// Save the updated configuration
	if err := configService.SaveConfig(config); err != nil {
		fmt.Printf("❌ Cannot save configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration updated successfully")
}

func runConfigInit(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Initializing configuration file...")

	force, _ := cmd.Flags().GetBool("force")

	// Get config service from DI container
	injector := GetContainer()
	configService := do.MustInvoke[services.ConfigService](injector)

	// Check if file already exists (simplified check)
	if !force {
		if _, err := configService.LoadConfig(); err == nil {
			fmt.Println("❌ Configuration already exists")
			fmt.Println("💡 Use --force to overwrite, or 'goreleaser-cli config show' to view current config")
			os.Exit(1)
		}
	}

	// Create default configuration using the service
	config, err := configService.InitConfig()
	if err != nil {
		fmt.Printf("❌ Cannot create default configuration: %v\n", err)
		os.Exit(1)
	}

	// Save the configuration using the service
	if err := configService.SaveConfig(config); err != nil {
		fmt.Printf("❌ Cannot save configuration file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration file created successfully")
	fmt.Println("\n📝 Edit the file to customize your settings")
	fmt.Println("\n💡 Use 'goreleaser-cli config show' to view current settings")
}
