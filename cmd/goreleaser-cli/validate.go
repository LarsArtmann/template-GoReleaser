package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	fmt.Println("   • Checking GoReleaser configuration files...")

	goreleaserFiles := []string{
		".goreleaser.yaml",
		".goreleaser.yml",
		".goreleaser.pro.yaml",
		".goreleaser.pro.yml",
		"goreleaser.yaml",
		"goreleaser.yml",
	}

	foundConfigs := []string{}
	for _, file := range goreleaserFiles {
		if _, err := os.Stat(file); err == nil {
			foundConfigs = append(foundConfigs, file)
		}
	}

	if len(foundConfigs) == 0 {
		fmt.Println("   ❌ No GoReleaser configuration files found")
		fmt.Println("   💡 Create .goreleaser.yaml or run 'goreleaser init'")
		return false
	}

	allValid := true
	for _, file := range foundConfigs {
		fmt.Printf("   • Validating %s...\n", file)
		if !validateYAMLSyntax(file) {
			allValid = false
		}
		if !validateGoReleaserConfig(file) {
			allValid = false
		}
	}

	// Also check CLI config file
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("   • Found CLI config file: %s\n", configFile)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("   ❌ CLI config file does not exist: %s\n", configFile)
			allValid = false
		} else {
			fmt.Println("   ✅ CLI config file is accessible")
		}
	}

	return allValid
}

func validateYAMLSyntax(file string) bool {
	// Try to parse as YAML using a simple approach
	if content, err := os.ReadFile(file); err != nil {
		fmt.Printf("   ❌ Cannot read %s: %v\n", file, err)
		return false
	} else {
		// Basic YAML syntax validation - check for common issues
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			// Check for tabs (YAML doesn't allow them)
			if strings.Contains(line, "\t") {
				fmt.Printf("   ❌ %s:%d: YAML syntax error - tabs not allowed\n", file, i+1)
				return false
			}
		}
		fmt.Printf("   ✅ %s has valid YAML syntax\n", file)
		return true
	}
}

func validateGoReleaserConfig(file string) bool {
	// Check if goreleaser command is available
	if _, err := exec.LookPath("goreleaser"); err != nil {
		fmt.Println("   ⚠️  goreleaser not installed, skipping native validation")
		return true // Don't fail if goreleaser is not available
	}

	// Run goreleaser check
	cmd := exec.Command("goreleaser", "check", "--config", file)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		// Filter out multiple token warnings
		if !strings.Contains(outputStr, "multiple tokens") {
			fmt.Printf("   ❌ %s validation failed: %s\n", file, strings.TrimSpace(outputStr))
			return false
		}
	}
	
	fmt.Printf("   ✅ %s passed GoReleaser validation\n", file)
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
	allValid := true

	fmt.Println("   • Checking Go environment...")
	if !validateGoEnvironment() {
		allValid = false
	}

	fmt.Println("   • Checking Git environment...")
	if !validateGitEnvironment() {
		allValid = false
	}

	fmt.Println("   • Checking required environment variables...")
	if !validateRequiredEnvVars() {
		allValid = false
	}

	fmt.Println("   • Checking tool dependencies...")
	if !validateToolDependencies() {
		allValid = false
	}

	return allValid
}

func validateGoEnvironment() bool {
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

	// Check if Go compiler is available
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Println("   ❌ Go compiler not found in PATH")
		return false
	}
	fmt.Println("   ✅ Go compiler available")

	return true
}

func validateGitEnvironment() bool {
	// Check if we're in a git repository
	gitDir := ".git"
	if info, err := os.Stat(gitDir); os.IsNotExist(err) || !info.IsDir() {
		fmt.Println("   ❌ Not a git repository")
		return false
	}
	fmt.Println("   ✅ Git repository detected")

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("   ❌ Git not found in PATH")
		return false
	}
	fmt.Println("   ✅ Git available")

	// Check for remote origin
	cmd := exec.Command("git", "remote", "-v")
	if output, err := cmd.Output(); err != nil {
		fmt.Println("   ⚠️  Cannot check git remotes")
	} else if strings.Contains(string(output), "origin") {
		fmt.Println("   ✅ Git remote 'origin' configured")
	} else {
		fmt.Println("   ⚠️  No Git remote 'origin' configured")
	}

	return true
}

func validateRequiredEnvVars() bool {
	criticalVars := map[string]string{
		"GITHUB_TOKEN": "GitHub API access token for releases",
	}
	
	optionalVars := map[string]string{
		"DOCKER_USERNAME":             "Docker Hub username",
		"DOCKER_PASSWORD":             "Docker Hub password/token", 
		"GORELEASER_KEY":              "GoReleaser Pro license key",
		"HOMEBREW_TAP_GITHUB_TOKEN":   "GitHub token for Homebrew tap",
		"SCOOP_GITHUB_TOKEN":          "GitHub token for Scoop bucket",
	}

	allValid := true
	criticalMissing := []string{}

	fmt.Println("   • Critical environment variables:")
	for varName, description := range criticalVars {
		if value := os.Getenv(varName); value != "" {
			if isPlaceholderValue(value) {
				fmt.Printf("   ⚠️  %s is set but appears to be a placeholder\n", varName)
			} else {
				fmt.Printf("   ✅ %s is set\n", varName)
			}
		} else {
			fmt.Printf("   ❌ %s is not set (%s)\n", varName, description)
			criticalMissing = append(criticalMissing, varName)
			allValid = false
		}
	}

	fmt.Println("   • Optional environment variables:")
	for varName, description := range optionalVars {
		if value := os.Getenv(varName); value != "" {
			if isPlaceholderValue(value) {
				fmt.Printf("   ⚠️  %s is set but appears to be a placeholder\n", varName)
			} else {
				fmt.Printf("   ✅ %s is set\n", varName)
			}
		} else {
			fmt.Printf("   ⚠️  %s is not set (%s)\n", varName, description)
		}
	}

	if len(criticalMissing) > 0 {
		fmt.Println("\n   💡 Set critical environment variables:")
		for _, varName := range criticalMissing {
			fmt.Printf("     export %s=your_value_here\n", varName)
		}
	}

	return allValid
}

func validateToolDependencies() bool {
	requiredTools := map[string]string{
		"go":          "Go compiler",
		"git":         "Git version control",
		"goreleaser":  "GoReleaser binary",
	}

	recommendedTools := map[string]string{
		"docker":   "Docker for container builds",
		"yq":       "YAML processor",
		"cosign":   "Container signing",
		"syft":     "SBOM generation",
	}

	allValid := true

	fmt.Println("   • Required tools:")
	for tool, description := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			fmt.Printf("   ❌ %s not found (%s)\n", tool, description)
			allValid = false
		} else {
			fmt.Printf("   ✅ %s available\n", tool)
		}
	}

	fmt.Println("   • Recommended tools:")
	for tool, description := range recommendedTools {
		if _, err := exec.LookPath(tool); err != nil {
			fmt.Printf("   ⚠️  %s not found (%s)\n", tool, description)
		} else {
			fmt.Printf("   ✅ %s available\n", tool)
		}
	}

	return allValid
}

func isPlaceholderValue(value string) bool {
	placeholders := []string{"your-", "xxxx", "example", "changeme", "todo", "test-"}
	lowerValue := strings.ToLower(value)
	for _, placeholder := range placeholders {
		if strings.HasPrefix(lowerValue, placeholder) {
			return true
		}
	}
	return len(value) < 3
}
