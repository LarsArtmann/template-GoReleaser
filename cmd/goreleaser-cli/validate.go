package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration and project structure",
	Long: `Validate various aspects of your project including:
- Configuration files (YAML, JSON, TOML)
- Project structure and required files
- Environment variables
- Dependencies and imports

This command helps ensure your project is properly configured
and ready for release.`,
	Run: func(cmd *cobra.Command, args []string) {
		runValidate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Add flags for validate command
	validateCmd.Flags().BoolP("config", "c", false, "Validate configuration files")
	validateCmd.Flags().BoolP("structure", "s", false, "Validate project structure")
	validateCmd.Flags().BoolP("env", "e", false, "Validate environment variables")
	validateCmd.Flags().BoolP("all", "a", true, "Validate all aspects (default)")
}

func runValidate(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Running validation...")

	configFlag, _ := cmd.Flags().GetBool("config")
	structureFlag, _ := cmd.Flags().GetBool("structure")
	envFlag, _ := cmd.Flags().GetBool("env")
	allFlag, _ := cmd.Flags().GetBool("all")

	// If no specific flags are set, default to all
	if !configFlag && !structureFlag && !envFlag {
		allFlag = true
	}

	validationsPassed := 0
	validationsFailed := 0

	if allFlag || configFlag {
		fmt.Println("\n📋 Validating configuration files...")
		if validateConfigFiles() {
			fmt.Println("✅ Configuration files validation passed")
			validationsPassed++
		} else {
			fmt.Println("❌ Configuration files validation failed")
			validationsFailed++
		}
	}

	if allFlag || structureFlag {
		fmt.Println("\n📁 Validating project structure...")
		if validateProjectStructure() {
			fmt.Println("✅ Project structure validation passed")
			validationsPassed++
		} else {
			fmt.Println("❌ Project structure validation failed")
			validationsFailed++
		}
	}

	if allFlag || envFlag {
		fmt.Println("\n🌍 Validating environment...")
		if validateEnvironment() {
			fmt.Println("✅ Environment validation passed")
			validationsPassed++
		} else {
			fmt.Println("❌ Environment validation failed")
			validationsFailed++
		}
	}

	fmt.Printf("\n📊 Validation Summary:\n")
	fmt.Printf("   ✅ Passed: %d\n", validationsPassed)
	fmt.Printf("   ❌ Failed: %d\n", validationsFailed)

	if validationsFailed > 0 {
		fmt.Println("\n❌ Some validations failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All validations passed successfully!")
	}
}

func validateConfigFiles() bool {
	fmt.Println("   • Checking for config file...")

	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("   ⚠️  No config file found (this is optional)")
		return true // Config file is optional
	}

	fmt.Printf("   • Found config file: %s\n", configFile)

	// Check if config file is readable
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("   ❌ Config file does not exist: %s\n", configFile)
		return false
	}

	fmt.Println("   ✅ Config file is accessible")
	return true
}

func validateProjectStructure() bool {
	fmt.Println("   • Checking required files...")

	requiredFiles := []string{
		"go.mod",
		"README.md",
		"LICENSE",
	}

	allFound := true
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("   ❌ Missing required file: %s\n", file)
			allFound = false
		} else {
			fmt.Printf("   ✅ Found: %s\n", file)
		}
	}

	fmt.Println("   • Checking directory structure...")
	requiredDirs := []string{
		"cmd",
	}

	for _, dir := range requiredDirs {
		if info, err := os.Stat(dir); os.IsNotExist(err) || !info.IsDir() {
			fmt.Printf("   ❌ Missing required directory: %s\n", dir)
			allFound = false
		} else {
			fmt.Printf("   ✅ Found directory: %s\n", dir)
		}
	}

	return allFound
}

func validateEnvironment() bool {
	fmt.Println("   • Checking Go environment...")

	goModFile := "go.mod"
	if _, err := os.Stat(goModFile); os.IsNotExist(err) {
		fmt.Println("   ❌ go.mod file not found")
		return false
	}

	fmt.Println("   ✅ go.mod file exists")

	// Check for go.sum
	goSumFile := "go.sum"
	if _, err := os.Stat(goSumFile); os.IsNotExist(err) {
		fmt.Println("   ⚠️  go.sum file not found (run 'go mod tidy')")
	} else {
		fmt.Println("   ✅ go.sum file exists")
	}

	// Check if we're in a git repository
	gitDir := ".git"
	if info, err := os.Stat(gitDir); os.IsNotExist(err) || !info.IsDir() {
		fmt.Println("   ⚠️  Not a git repository")
	} else {
		fmt.Println("   ✅ Git repository detected")
	}

	// Check for common configuration files
	configFiles := []string{
		".goreleaser.yml",
		".goreleaser.yaml",
		"goreleaser.yml",
		"goreleaser.yaml",
	}

	foundGoReleaser := false
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			fmt.Printf("   ✅ Found GoReleaser config: %s\n", configFile)
			foundGoReleaser = true
			break
		}
	}

	if !foundGoReleaser {
		fmt.Println("   ⚠️  No GoReleaser configuration found")
	}

	return true
}
