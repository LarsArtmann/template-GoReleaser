package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Manage project licenses",
	Long: `Manage project licenses with various subcommands:
- Generate license files from templates
- List available license templates  
- Show current license information
- Validate existing license files

This integrates with the existing license management system
in the GoReleaser template.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLicense(cmd, args)
	},
}

var licenseGenerateCmd = &cobra.Command{
	Use:   "generate [license-type] [copyright-holder]",
	Short: "Generate a license file",
	Long: `Generate a license file from available templates.

Examples:
  goreleaser-cli license generate MIT "John Doe"
  goreleaser-cli license generate Apache-2.0
  goreleaser-cli license generate --interactive`,
	Run: func(cmd *cobra.Command, args []string) {
		runLicenseGenerate(cmd, args)
	},
}

var licenseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available license templates",
	Long: `List all available license templates that can be used
to generate license files.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLicenseList(cmd, args)
	},
}

var licenseShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current license information",
	Long: `Display information about the current project license,
including type, copyright holder, and file contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLicenseShow(cmd, args)
	},
}

var licenseValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate existing license file",
	Long: `Validate that the existing license file is properly formatted
and contains all required information.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLicenseValidate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(licenseGenerateCmd)
	licenseCmd.AddCommand(licenseListCmd)
	licenseCmd.AddCommand(licenseShowCmd)
	licenseCmd.AddCommand(licenseValidateCmd)

	// Add flags for license generate command
	licenseGenerateCmd.Flags().BoolP("interactive", "i", false, "Interactive license generation")
	licenseGenerateCmd.Flags().StringP("type", "t", "", "License type")
	licenseGenerateCmd.Flags().StringP("holder", "c", "", "Copyright holder")
	licenseGenerateCmd.Flags().StringP("output", "o", "LICENSE", "Output file path")

	// Add flags for license show command
	licenseShowCmd.Flags().BoolP("preview", "p", false, "Show preview of license content")
}

func runLicense(cmd *cobra.Command, args []string) {
	fmt.Println("📄 License Management")
	fmt.Println("\nUse one of the following subcommands:")
	fmt.Println("  generate  - Generate a new license file")
	fmt.Println("  list      - List available license templates")
	fmt.Println("  show      - Show current license information")
	fmt.Println("  validate  - Validate existing license file")
	fmt.Println("\nFor more help: goreleaser-cli license --help")
}

func runLicenseGenerate(cmd *cobra.Command, args []string) {
	fmt.Println("📝 Generating license...")

	interactive, _ := cmd.Flags().GetBool("interactive")
	licenseType, _ := cmd.Flags().GetString("type")
	copyrightHolder, _ := cmd.Flags().GetString("holder")
	outputFile, _ := cmd.Flags().GetString("output")

	// Get values from args if provided
	if len(args) > 0 && licenseType == "" {
		licenseType = args[0]
	}
	if len(args) > 1 && copyrightHolder == "" {
		copyrightHolder = args[1]
	}

	// Interactive mode
	if interactive {
		fmt.Println("🤖 Interactive license generation mode")

		if licenseType == "" {
			fmt.Print("Enter license type (e.g., MIT, Apache-2.0): ")
			if _, err := fmt.Scanln(&licenseType); err != nil {
				fmt.Printf("❌ Failed to read license type: %v\n", err)
				os.Exit(1)
			}
		}

		if copyrightHolder == "" {
			fmt.Print("Enter copyright holder name: ")
			if _, err := fmt.Scanln(&copyrightHolder); err != nil {
				fmt.Printf("❌ Failed to read copyright holder: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Use environment variables or config as fallback
	if licenseType == "" {
		licenseType = os.Getenv("LICENSE_TYPE")
		if licenseType == "" {
			licenseType = viper.GetString("license.type")
		}
	}

	if copyrightHolder == "" {
		copyrightHolder = os.Getenv("COPYRIGHT_HOLDER")
		if copyrightHolder == "" {
			copyrightHolder = viper.GetString("author.name")
		}
	}

	// Validate inputs
	if licenseType == "" {
		fmt.Println("❌ License type is required")
		fmt.Println("💡 Use: goreleaser-cli license generate MIT \"Your Name\"")
		fmt.Println("💡 Or: goreleaser-cli license generate --interactive")
		os.Exit(1)
	}

	// Use the existing license generation script
	if err := generateLicenseWithScript(licenseType, copyrightHolder); err != nil {
		fmt.Printf("❌ License generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ License generated successfully: %s\n", outputFile)
}

func runLicenseList(cmd *cobra.Command, args []string) {
	fmt.Println("📋 Available license templates:")

	templatesDir := "assets/licenses"
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		fmt.Printf("❌ Cannot read templates directory: %v\n", err)
		os.Exit(1)
	}

	found := false
	for _, entry := range files {
		if strings.HasSuffix(entry.Name(), ".template") {
			licenseName := strings.TrimSuffix(entry.Name(), ".template")
			fmt.Printf("  • %s\n", licenseName)
			found = true
		}
	}

	if !found {
		fmt.Println("⚠️  No license templates found")
	} else {
		fmt.Printf("\n💡 Generate a license with: goreleaser-cli license generate <type> \"Your Name\"\n")
	}
}

func runLicenseShow(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Current license information:")

	licenseFile := "LICENSE"
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		fmt.Println("❌ No LICENSE file found in current directory")
		fmt.Println("💡 Generate one with: goreleaser-cli license generate <type>")
		os.Exit(1)
	}

	// Show file info
	info, err := os.Stat(licenseFile)
	if err != nil {
		fmt.Printf("❌ Cannot read LICENSE file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  📄 File: %s\n", licenseFile)
	fmt.Printf("  📏 Size: %d bytes\n", info.Size())
	fmt.Printf("  📅 Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))

	preview, _ := cmd.Flags().GetBool("preview")
	if preview {
		fmt.Println("\n📖 License content preview (first 15 lines):")
		content, err := os.ReadFile(licenseFile)
		if err != nil {
			fmt.Printf("❌ Cannot read file content: %v\n", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		maxLines := 15
		if len(lines) < maxLines {
			maxLines = len(lines)
		}

		for i := 0; i < maxLines; i++ {
			fmt.Printf("  %s\n", lines[i])
		}

		if len(lines) > maxLines {
			fmt.Printf("  ... (%d more lines)\n", len(lines)-maxLines)
		}
	} else {
		fmt.Println("💡 Use --preview to see license content")
	}
}

func runLicenseValidate(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Validating license file...")

	licenseFile := "LICENSE"
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		fmt.Println("❌ No LICENSE file found")
		os.Exit(1)
	}

	// Read and validate content
	content, err := os.ReadFile(licenseFile)
	if err != nil {
		fmt.Printf("❌ Cannot read LICENSE file: %v\n", err)
		os.Exit(1)
	}

	contentStr := string(content)
	validationsPassed := 0
	validationsFailed := 0

	// Check file size
	fmt.Println("   • Checking file size...")
	if len(content) < 100 {
		fmt.Println("   ❌ License file seems too small")
		validationsFailed++
	} else {
		fmt.Println("   ✅ File size is reasonable")
		validationsPassed++
	}

	// Check for template placeholders
	fmt.Println("   • Checking for unsubstituted placeholders...")
	if strings.Contains(contentStr, "{{") || strings.Contains(contentStr, "}}") {
		fmt.Println("   ❌ Found unsubstituted template variables")
		validationsFailed++
	} else {
		fmt.Println("   ✅ No template placeholders found")
		validationsPassed++
	}

	// Check for common license patterns
	fmt.Println("   • Checking license format...")
	hasLicenseKeywords := strings.Contains(strings.ToLower(contentStr), "license") ||
		strings.Contains(strings.ToLower(contentStr), "copyright") ||
		strings.Contains(strings.ToLower(contentStr), "permission")

	if hasLicenseKeywords {
		fmt.Println("   ✅ License keywords found")
		validationsPassed++
	} else {
		fmt.Println("   ❌ No common license keywords found")
		validationsFailed++
	}

	fmt.Printf("\n📊 Validation Summary:\n")
	fmt.Printf("   ✅ Passed: %d\n", validationsPassed)
	fmt.Printf("   ❌ Failed: %d\n", validationsFailed)

	if validationsFailed > 0 {
		fmt.Println("\n❌ License validation failed")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 License validation passed!")
	}
}

func generateLicenseWithScript(licenseType, copyrightHolder string) error {
	scriptPath := "scripts/generate-license.sh"

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("license generation script not found: %s", scriptPath)
	}

	// Prepare command arguments
	var args []string
	if licenseType != "" {
		args = append(args, licenseType)
	}
	if copyrightHolder != "" {
		args = append(args, copyrightHolder)
	}

	// Execute the license generation script
	cmd := exec.Command("bash", append([]string{scriptPath}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
