package main

import (
	"fmt"
	"os"

	"github.com/samber/do"
	"github.com/spf13/cobra"
	"github.com/LarsArtmann/template-GoReleaser/internal/services"
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

	// Get validation service from DI container
	injector := GetContainer()
	validationService := do.MustInvoke[services.ValidationService](injector)

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

	if allFlag || structureFlag {
		fmt.Println("\n📁 Validating project structure...")
		if result, err := validationService.ValidateProject(); err != nil {
			fmt.Printf("❌ Error during project validation: %v\n", err)
			validationsFailed++
		} else if result.Success {
			fmt.Printf("✅ Project structure validation passed (%d checks)\n", result.Checks)
			validationsPassed++
		} else {
			fmt.Printf("❌ Project structure validation failed (%d checks)\n", result.Checks)
			for _, error := range result.Errors {
				fmt.Printf("   • %s\n", error)
			}
			for _, warning := range result.Warnings {
				fmt.Printf("   ⚠️ %s\n", warning)
			}
			validationsFailed++
		}
	}

	if allFlag || envFlag {
		fmt.Println("\n🌍 Validating environment...")
		if result, err := validationService.ValidateEnvironment(); err != nil {
			fmt.Printf("❌ Error during environment validation: %v\n", err)
			validationsFailed++
		} else if result.Success {
			fmt.Printf("✅ Environment validation passed (%d checks)\n", result.Checks)
			validationsPassed++
		} else {
			fmt.Printf("❌ Environment validation failed (%d checks)\n", result.Checks)
			for _, error := range result.Errors {
				fmt.Printf("   • %s\n", error)
			}
			for _, warning := range result.Warnings {
				fmt.Printf("   ⚠️ %s\n", warning)
			}
			validationsFailed++
		}
	}

	if allFlag || configFlag {
		fmt.Println("\n🛠️ Validating tools...")
		if result, err := validationService.ValidateTools(); err != nil {
			fmt.Printf("❌ Error during tools validation: %v\n", err)
			validationsFailed++
		} else if result.Success {
			fmt.Printf("✅ Tools validation passed (%d checks)\n", result.Checks)
			validationsPassed++
		} else {
			fmt.Printf("❌ Tools validation failed (%d checks)\n", result.Checks)
			for _, error := range result.Errors {
				fmt.Printf("   • %s\n", error)
			}
			for _, warning := range result.Warnings {
				fmt.Printf("   ⚠️ %s\n", warning)
			}
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

