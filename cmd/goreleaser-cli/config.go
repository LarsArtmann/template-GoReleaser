package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("⚠️  No configuration file in use")
		fmt.Println("💡 Run 'goreleaser-cli config init' to create one")
		return
	}

	fmt.Printf("📄 Config file: %s\n", configFile)
	fmt.Printf("📏 File format: %s\n", filepath.Ext(configFile))

	// Load typed configuration
	config, err := LoadTypedConfig()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return
	}

	fmt.Println("\n🔧 Configuration values:")
	fmt.Printf("  📄 License: %s (year: %d)\n", config.License.Type, config.License.Year)
	fmt.Printf("  👤 Author: %s <%s>\n", config.Author.Name, config.Author.Email)
	fmt.Printf("  📦 Project: %s\n", config.Project.Name)
	if config.Project.Description != "" {
		fmt.Printf("  📋 Description: %s\n", config.Project.Description)
	}
	fmt.Printf("  🛠️  CLI: verbose=%t, colors=%t, timeout=%ds\n", 
		config.CLI.Verbose, config.CLI.Colors, config.CLI.Timeout)

	// Validate configuration
	if errors := config.Validate(); len(errors) > 0 {
		fmt.Println("\n⚠️  Configuration validation warnings:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
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

	// Set the value in viper
	viper.Set(key, value)

	// Write the config file
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// No config file exists, create one
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Cannot determine home directory: %v\n", err)
			os.Exit(1)
		}

		configFile = filepath.Join(home, ".goreleaser-cli.yaml")
		viper.SetConfigFile(configFile)
	}

	if err := viper.WriteConfig(); err != nil {
		// If WriteConfig fails, try WriteConfigAs (in case file doesn't exist)
		if err := viper.WriteConfigAs(configFile); err != nil {
			fmt.Printf("❌ Cannot write config file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("✅ Configuration saved to: %s\n", configFile)
}

func runConfigInit(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Initializing configuration file...")

	force, _ := cmd.Flags().GetBool("force")

	// Determine config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Cannot determine home directory: %v\n", err)
		os.Exit(1)
	}

	configFile := filepath.Join(home, ".goreleaser-cli.yaml")

	// Check if file already exists
	if _, err := os.Stat(configFile); err == nil && !force {
		fmt.Printf("❌ Configuration file already exists: %s\n", configFile)
		fmt.Println("💡 Use --force to overwrite, or 'goreleaser-cli config show' to view current config")
		os.Exit(1)
	}

	// Create default configuration
	defaultConfig := `# GoReleaser CLI Configuration File
# This file stores default values for the CLI tool

# License settings
license:
  type: "MIT"
  
# Author/Copyright settings  
author:
  name: ""
  email: ""

# Project settings
project:
  name: ""
  description: ""

# CLI behavior settings
cli:
  verbose: false
  colors: true
`

	// Write the config file
	if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
		fmt.Printf("❌ Cannot create config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuration file created: %s\n", configFile)
	fmt.Println("\n📝 Edit the file to customize your settings:")
	fmt.Printf("   %s\n", configFile)
	fmt.Println("\n💡 Use 'goreleaser-cli config show' to view current settings")
}

// LoadTypedConfig loads the configuration into typed structs
func LoadTypedConfig() (*Config, error) {
	config := DefaultConfig()
	
	// If no config file is set, try to find one
	if viper.ConfigFileUsed() == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configFile := filepath.Join(home, ".goreleaser-cli.yaml")
			if _, err := os.Stat(configFile); err == nil {
				viper.SetConfigFile(configFile)
				viper.ReadInConfig()
			}
		}
	}
	
	// Unmarshal into typed config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Apply environment variable overrides
	if licenseType := os.Getenv("LICENSE_TYPE"); licenseType != "" {
		config.License.Type = licenseType
	}
	if authorName := os.Getenv("COPYRIGHT_HOLDER"); authorName != "" {
		config.Author.Name = authorName
	}
	if authorName := os.Getenv("AUTHOR_NAME"); authorName != "" && config.Author.Name == "" {
		config.Author.Name = authorName
	}
	if projectAuthor := os.Getenv("PROJECT_AUTHOR"); projectAuthor != "" && config.Author.Name == "" {
		config.Author.Name = projectAuthor
	}
	
	// Set current year if not specified
	if config.License.Year == 0 {
		config.License.Year = time.Now().Year()
	}
	
	return config, nil
}

// SaveTypedConfig saves the typed configuration to file
func SaveTypedConfig(config *Config) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		configFile = filepath.Join(home, ".goreleaser-cli.yaml")
		viper.SetConfigFile(configFile)
	}
	
	// Convert config to viper settings
	viper.Set("license", config.License)
	viper.Set("author", config.Author)
	viper.Set("project", config.Project)
	viper.Set("cli", config.CLI)
	
	if err := viper.WriteConfig(); err != nil {
		// If WriteConfig fails, try WriteConfigAs
		if err := viper.WriteConfigAs(configFile); err != nil {
			return fmt.Errorf("cannot write config file: %w", err)
		}
	}
	
	return nil
}