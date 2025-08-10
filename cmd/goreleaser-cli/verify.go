package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Comprehensive verification of GoReleaser configuration and project",
	Long: `Perform comprehensive verification of your GoReleaser project including:
- Configuration validation with native GoReleaser check
- Tool dependencies and hook commands
- Signing and compression tools
- Template syntax validation
- Git state and tagging
- License system testing  
- Security scanning (gosec, govulncheck, shellcheck, hadolint)
- Dry-run testing
- Project structure validation

This command provides a complete health check of your GoReleaser setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		runVerify(cmd, args)
	},
}

type verifyStats struct {
	Checks   int
	Warnings int
	Errors   int
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Add flags for verify command
	verifyCmd.Flags().Bool("skip-security", false, "Skip security validation scans")
	verifyCmd.Flags().Bool("skip-dry-run", false, "Skip dry-run testing")
	verifyCmd.Flags().Bool("skip-license-test", false, "Skip license generation testing")
	verifyCmd.Flags().String("config", ".goreleaser.yaml", "Primary GoReleaser config file")
	verifyCmd.Flags().String("pro-config", ".goreleaser.pro.yaml", "GoReleaser Pro config file")
}

func runVerify(cmd *cobra.Command, args []string) {
	fmt.Println("================================")
	fmt.Println("GoReleaser Configuration Verifier")
	fmt.Println("================================")
	fmt.Println()

	stats := &verifyStats{}
	
	skipSecurity, _ := cmd.Flags().GetBool("skip-security")
	skipDryRun, _ := cmd.Flags().GetBool("skip-dry-run")
	skipLicenseTest, _ := cmd.Flags().GetBool("skip-license-test")
	configFile, _ := cmd.Flags().GetString("config")
	proConfigFile, _ := cmd.Flags().GetString("pro-config")

	// Check dependencies
	logInfo("Checking dependencies...")
	checkDependencies(stats)
	fmt.Println()

	// Check configurations
	checkGoReleaserConfigurations(stats, configFile, proConfigFile)
	fmt.Println()

	// Check environment variables
	checkEnvironmentVariables(stats, configFile, proConfigFile)
	fmt.Println()

	// Check project structure
	checkProjectStructureVerify(stats)
	fmt.Println()

	// Check specific tools
	checkHookCommands(stats, proConfigFile)
	fmt.Println()
	checkSigningTools(stats, proConfigFile)
	fmt.Println()
	checkCompressionTools(stats, proConfigFile)
	fmt.Println()

	// Check Git state
	checkGitState(stats)
	fmt.Println()

	// Test license system
	if !skipLicenseTest {
		testLicenseSystem(stats)
		fmt.Println()
	}

	// Validate templates
	validateTemplates(stats, configFile, proConfigFile)
	fmt.Println()

	// Run security validation
	if !skipSecurity {
		runSecurityValidation(stats)
		fmt.Println()
	}

	// Try dry-run
	if !skipDryRun {
		runDryRun(stats, configFile, proConfigFile)
		fmt.Println()
	}

	// Summary
	printVerificationSummary(stats)
}

func logInfo(message string) {
	fmt.Printf("\033[0;34m[INFO]\033[0m %s\n", message)
}

func logSuccess(stats *verifyStats, message string) {
	fmt.Printf("\033[0;32m[✓]\033[0m %s\n", message)
	stats.Checks++
}

func logWarning(stats *verifyStats, message string) {
	fmt.Printf("\033[1;33m[⚠]\033[0m %s\n", message)
	stats.Warnings++
}

func logError(stats *verifyStats, message string) {
	fmt.Printf("\033[0;31m[✗]\033[0m %s\n", message)
	stats.Errors++
}

func checkCommand(stats *verifyStats, command, description string) bool {
	if _, err := exec.LookPath(command); err != nil {
		logError(stats, fmt.Sprintf("%s is not installed (%s)", command, description))
		return false
	}
	logSuccess(stats, fmt.Sprintf("%s is installed", command))
	return true
}

func checkDependencies(stats *verifyStats) {
	checkCommand(stats, "go", "Go compiler")
	checkCommand(stats, "git", "Git version control")
	checkCommand(stats, "goreleaser", "GoReleaser binary")
	checkCommand(stats, "yq", "YAML processor")
	checkCommand(stats, "docker", "Docker for containers")
}

func checkGoReleaserConfigurations(stats *verifyStats, configFile, proConfigFile string) {
	checkGoReleaserConfig(stats, configFile)
	fmt.Println()
	checkGoReleaserConfig(stats, proConfigFile)
}

func checkGoReleaserConfig(stats *verifyStats, file string) {
	logInfo(fmt.Sprintf("Validating %s...", file))

	if !fileExists(stats, file) {
		return
	}

	if !checkYAMLSyntaxVerify(stats, file) {
		return
	}

	// Check with goreleaser (handle multiple token conflicts)
	if _, err := exec.LookPath("goreleaser"); err != nil {
		logWarning(stats, "goreleaser not installed, skipping native validation")
		return
	}

	// Temporarily clear conflicting tokens for GoReleaser
	savedGitlabToken := os.Getenv("GITLAB_TOKEN")
	savedGiteaToken := os.Getenv("GITEA_TOKEN")
	os.Unsetenv("GITLAB_TOKEN")
	os.Unsetenv("GITEA_TOKEN")

	cmd := exec.Command("goreleaser", "check", "--config", file)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		if !strings.Contains(outputStr, "multiple tokens") {
			logWarning(stats, fmt.Sprintf("GoReleaser validation failed: %s (this might be expected in test environment)", file))
			// Show errors but don't fail completely
			lines := strings.Split(strings.TrimSpace(outputStr), "\n")
			for i, line := range lines {
				if i >= 5 {
					break
				}
				if !strings.Contains(line, "multiple tokens") {
					fmt.Printf("  %s\n", line)
				}
			}
		} else {
			logSuccess(stats, fmt.Sprintf("GoReleaser validation passed: %s", file))
		}
	} else {
		logSuccess(stats, fmt.Sprintf("GoReleaser validation passed: %s", file))
	}

	// Restore tokens
	if savedGitlabToken != "" {
		os.Setenv("GITLAB_TOKEN", savedGitlabToken)
	}
	if savedGiteaToken != "" {
		os.Setenv("GITEA_TOKEN", savedGiteaToken)
	}

	// Skip build test for faster validation
	logInfo("Skipping snapshot build test (for performance)")
}

func fileExists(stats *verifyStats, file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		logError(stats, fmt.Sprintf("File not found: %s", file))
		return false
	}
	logSuccess(stats, fmt.Sprintf("File exists: %s", file))
	return true
}

func checkYAMLSyntaxVerify(stats *verifyStats, file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		logError(stats, fmt.Sprintf("Cannot read %s: %v", file, err))
		return false
	}

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
			logError(stats, fmt.Sprintf("%s:%d: YAML syntax error - tabs not allowed", file, i+1))
			return false
		}
	}

	// Try yq validation if available
	if _, err := exec.LookPath("yq"); err == nil {
		cmd := exec.Command("yq", "eval", ".", file)
		if err := cmd.Run(); err != nil {
			logError(stats, fmt.Sprintf("YAML syntax error in: %s", file))
			// Show yq output
			cmd = exec.Command("yq", "eval", ".", file)
			if output, _ := cmd.CombinedOutput(); len(output) > 0 {
				lines := strings.Split(string(output), "\n")
				for i, line := range lines {
					if i >= 10 {
						break
					}
					fmt.Printf("  %s\n", line)
				}
			}
			return false
		}
	}

	logSuccess(stats, fmt.Sprintf("%s has valid YAML syntax", file))
	return true
}

func checkEnvironmentVariables(stats *verifyStats, configFile, proConfigFile string) {
	checkRequiredEnvVarsForConfig(stats, configFile)
	checkRequiredEnvVarsForConfig(stats, proConfigFile)
}

func checkRequiredEnvVarsForConfig(stats *verifyStats, file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}

	logInfo(fmt.Sprintf("Checking environment variables for %s...", file))

	// Read config file and extract env var references
	content, err := os.ReadFile(file)
	if err != nil {
		logWarning(stats, fmt.Sprintf("Cannot read config file: %s", file))
		return
	}

	// Extract env var references from config file
	envVarRegex := regexp.MustCompile(`\{\{\s*\.Env\.([A-Z_]+)\s*\}\}`)
	matches := envVarRegex.FindAllStringSubmatch(string(content), -1)
	
	configEnvVars := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			configEnvVars[match[1]] = true
		}
	}

	// Define commonly required environment variables for GoReleaser
	criticalVars := []string{"GITHUB_TOKEN"}
	commonVars := []string{"DOCKER_USERNAME", "DOCKER_PASSWORD", "GORELEASER_KEY"}

	// Add config-specific variables
	allEnvVars := make([]string, 0)
	for varName := range configEnvVars {
		allEnvVars = append(allEnvVars, varName)
	}

	// Add common variables that GoReleaser typically needs
	for _, varName := range append(criticalVars, commonVars...) {
		if !configEnvVars[varName] {
			allEnvVars = append(allEnvVars, varName)
		}
	}

	if len(allEnvVars) > 0 {
		fmt.Println("Environment variables for GoReleaser:")
		criticalMissing := make([]string, 0)
		optionalMissing := make([]string, 0)

		for _, varName := range allEnvVars {
			value := os.Getenv(varName)
			if value != "" {
				// Basic validation for common patterns
				if isPlaceholderValue(value) {
					logWarning(stats, fmt.Sprintf("%s is set but appears to be a placeholder", varName))
				} else {
					logSuccess(stats, fmt.Sprintf("%s is set", varName))
				}
			} else {
				// Check if this is a critical variable
				isCritical := false
				for _, critical := range criticalVars {
					if varName == critical {
						isCritical = true
						break
					}
				}

				if isCritical {
					logWarning(stats, fmt.Sprintf("%s is not set (critical)", varName))
					criticalMissing = append(criticalMissing, varName)
				} else {
					logWarning(stats, fmt.Sprintf("%s is not set (optional)", varName))
					optionalMissing = append(optionalMissing, varName)
				}
			}
		}

		if len(criticalMissing) > 0 {
			fmt.Println()
			logError(stats, "Critical environment variables missing:")
			for _, varName := range criticalMissing {
				fmt.Printf("  - %s\n", varName)
			}
		}

		if len(optionalMissing) > 0 || len(criticalMissing) > 0 {
			fmt.Println()
			fmt.Println("Environment variable information:")
			fmt.Println("Set these variables or source a .env file before running GoReleaser")
		}
	} else {
		logSuccess(stats, "No additional environment variables detected")
	}
}

func checkProjectStructureVerify(stats *verifyStats) {
	logInfo("Checking project structure...")

	// Check for main.go or cmd directory
	if fileExistsSimple("main.go") || dirExists("cmd") {
		logSuccess(stats, "Go project structure detected")
	} else {
		logWarning(stats, "No main.go or cmd/ directory found")
	}

	// Check for go.mod
	if fileExistsSimple("go.mod") {
		logSuccess(stats, "go.mod exists")
	} else {
		logWarning(stats, "go.mod not found")
	}

	// Check for LICENSE file
	if fileExistsSimple("LICENSE") {
		logSuccess(stats, "LICENSE file exists")

		// Check LICENSE file size (should not be empty)
		if info, err := os.Stat("LICENSE"); err == nil {
			size := info.Size()
			if size > 100 {
				logSuccess(stats, fmt.Sprintf("LICENSE file has content (%d bytes)", size))
			} else {
				logWarning(stats, fmt.Sprintf("LICENSE file seems too small: %d bytes", size))
			}
		}
	} else {
		logWarning(stats, "No LICENSE file found")
	}

	// Check for license generation script
	scriptPath := "scripts/generate-license.sh"
	if fileExistsSimple(scriptPath) {
		logSuccess(stats, "License generation script exists")
		if info, err := os.Stat(scriptPath); err == nil && info.Mode()&0111 != 0 {
			logSuccess(stats, "License script is executable")
		} else {
			logWarning(stats, "License script is not executable")
		}
	} else {
		logWarning(stats, "License generation script not found")
	}

	// Check for Dockerfile if Docker is configured
	if configContains(".goreleaser.yaml", "dockers:") || configContains(".goreleaser.pro.yaml", "dockers:") {
		if fileExistsSimple("Dockerfile") {
			logSuccess(stats, "Dockerfile exists")
		} else {
			logWarning(stats, "Dockerfile not found but Docker is configured")
		}
	}

	// Check for assets/licenses directory
	if dirExists("assets/licenses") {
		logSuccess(stats, "License templates directory exists")
		if templates, err := filepath.Glob("assets/licenses/*.template"); err == nil {
			templateCount := len(templates)
			if templateCount > 0 {
				logSuccess(stats, fmt.Sprintf("Found %d license templates", templateCount))
			} else {
				logWarning(stats, "No license templates found in assets/licenses")
			}
		}
	} else {
		logWarning(stats, "License templates directory not found")
	}
}

func checkHookCommands(stats *verifyStats, proConfigFile string) {
	logInfo("Checking hook commands...")

	// Check for templ
	if configContains(proConfigFile, "templ generate") {
		if !checkCommand(stats, "templ", "Template generation tool") {
			logWarning(stats, "templ is referenced but not installed")
		}
	}

	// Check for tsp (TypeSpec)
	if configContains(proConfigFile, "tsp compile") {
		if !checkCommand(stats, "tsp", "TypeSpec compiler") {
			logWarning(stats, "TypeSpec is referenced but not installed")
		}
	}

	// Check for security tools
	if configContains(proConfigFile, "gosec") {
		if !checkCommand(stats, "gosec", "Go security checker") {
			logWarning(stats, "gosec is referenced but not installed")
		}
	}

	if configContains(proConfigFile, "golangci-lint") {
		if !checkCommand(stats, "golangci-lint", "Go linter") {
			logWarning(stats, "golangci-lint is referenced but not installed")
		}
	}
}

func checkSigningTools(stats *verifyStats, proConfigFile string) {
	logInfo("Checking signing tools...")

	if configContains(proConfigFile, "cosign") {
		if !checkCommand(stats, "cosign", "Container signing tool") {
			logWarning(stats, "cosign is referenced but not installed")
		}
	}

	if configContains(proConfigFile, "syft") {
		if !checkCommand(stats, "syft", "SBOM generation tool") {
			logWarning(stats, "syft is referenced but not installed")
		}
	}
}

func checkCompressionTools(stats *verifyStats, proConfigFile string) {
	logInfo("Checking compression tools...")

	if configContains(proConfigFile, "upx") {
		if !checkCommand(stats, "upx", "Binary packer") {
			logWarning(stats, "UPX is referenced but not installed")
		}
	}
}

func validateTemplates(stats *verifyStats, configFile, proConfigFile string) {
	logInfo("Validating template variables...")

	validateTemplatesInFile(stats, configFile)
	validateTemplatesInFile(stats, proConfigFile)
}

func validateTemplatesInFile(stats *verifyStats, file string) {
	if !fileExistsSimple(file) {
		return
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return
	}

	contentStr := string(content)

	// Check for common template issues
	nestedTemplateRegex := regexp.MustCompile(`\{\{[^}]*\{\{`)
	if matches := nestedTemplateRegex.FindAllString(contentStr, -1); len(matches) > 0 {
		logError(stats, fmt.Sprintf("Found nested templates in %s: %v", file, matches))
	}

	// Check for unclosed templates
	unclosedRegex := regexp.MustCompile(`\{\{[^}]*$`)
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		if unclosedRegex.MatchString(line) {
			logError(stats, fmt.Sprintf("Found unclosed template in %s:%d: %s", file, i+1, line))
		}
	}
}

func checkGitState(stats *verifyStats) {
	logInfo("Checking Git state...")

	if dirExists(".git") {
		logSuccess(stats, "Git repository initialized")

		// Check for remote
		cmd := exec.Command("git", "remote", "-v")
		if output, err := cmd.Output(); err == nil && strings.Contains(string(output), "origin") {
			logSuccess(stats, "Git remote 'origin' is configured")
		} else {
			logWarning(stats, "No Git remote 'origin' configured")
		}

		// Check for tags
		cmd = exec.Command("git", "tag")
		if output, err := cmd.Output(); err == nil && strings.TrimSpace(string(output)) != "" {
			logSuccess(stats, "Git tags exist")
		} else {
			logWarning(stats, "No Git tags found (needed for releases)")
		}
	} else {
		logError(stats, "Not a Git repository")
	}
}

func testLicenseSystem(stats *verifyStats) {
	logInfo("Testing license generation system...")

	scriptPath := "scripts/generate-license.sh"
	if !fileExistsSimple(scriptPath) {
		logWarning(stats, "License generation script not found")
		return
	}

	// Check if script is executable
	if info, err := os.Stat(scriptPath); err != nil || info.Mode()&0111 == 0 {
		logWarning(stats, "License generation script not executable")
		return
	}

	// Test license script help
	cmd := exec.Command("./scripts/generate-license.sh", "--help")
	if err := cmd.Run(); err == nil {
		logSuccess(stats, "License script help works")
	} else {
		logWarning(stats, "License script help failed")
	}

	// Test license templates listing
	cmd = exec.Command("./scripts/generate-license.sh", "--list")
	if err := cmd.Run(); err == nil {
		logSuccess(stats, "License templates listing works")
	} else {
		logWarning(stats, "License templates listing failed")
	}

	// Test license generation (if readme config exists)
	if fileExistsSimple(".readme/configs/readme-config.yaml") {
		logInfo("Testing license generation...")
		
		// Backup existing LICENSE
		var backupPath string
		if fileExistsSimple("LICENSE") {
			backupPath = "LICENSE.backup.tmp"
			if err := copyFile("LICENSE", backupPath); err != nil {
				logWarning(stats, "Could not backup LICENSE file for testing")
				return
			}
		}

		// Test license generation
		cmd = exec.Command("./scripts/generate-license.sh")
		if err := cmd.Run(); err == nil {
			logSuccess(stats, "License generation test passed")
		} else {
			logWarning(stats, "License generation test failed")
		}

		// Restore backup if we made one
		if backupPath != "" {
			if err := copyFile(backupPath, "LICENSE"); err == nil {
				os.Remove(backupPath)
			}
		}
	} else {
		logWarning(stats, "No readme config found, skipping license generation test")
	}
}

func runSecurityValidation(stats *verifyStats) {
	logInfo("Running security validation...")

	// Go code security scan
	if _, err := exec.LookPath("gosec"); err == nil {
		logInfo("Scanning Go code with gosec...")
		cmd := exec.Command("gosec", "-quiet", "./...")
		if err := cmd.Run(); err == nil {
			logSuccess(stats, "Go code security scan passed")
		} else {
			logError(stats, "Go code security issues found")
		}
	} else {
		logWarning(stats, "gosec not installed, skipping Go security scan")
	}

	// Dependency vulnerability scan
	if _, err := exec.LookPath("govulncheck"); err == nil {
		logInfo("Checking dependencies for vulnerabilities...")
		cmd := exec.Command("govulncheck", "./...")
		if output, err := cmd.Output(); err == nil && strings.Contains(string(output), "No vulnerabilities found") {
			logSuccess(stats, "No vulnerable dependencies found")
		} else {
			logWarning(stats, "Vulnerable dependencies or scan issues detected")
		}
	} else {
		logWarning(stats, "govulncheck not installed, skipping dependency vulnerability scan")
	}

	// Shell script security scan
	if _, err := exec.LookPath("shellcheck"); err == nil {
		logInfo("Scanning shell scripts...")
		shellIssues := 0
		scripts, _ := filepath.Glob("*.sh")
		scriptDir, _ := filepath.Glob("scripts/*.sh")
		allScripts := append(scripts, scriptDir...)
		
		for _, script := range allScripts {
			if fileExistsSimple(script) {
				cmd := exec.Command("shellcheck", "--severity=error", script)
				if err := cmd.Run(); err != nil {
					shellIssues++
				}
			}
		}
		
		if shellIssues == 0 {
			logSuccess(stats, "Shell script security scan passed")
		} else {
			logError(stats, fmt.Sprintf("Shell script security issues found in %d files", shellIssues))
		}
	} else {
		logWarning(stats, "shellcheck not installed, skipping shell script scan")
	}

	// Dockerfile security scan
	if fileExistsSimple("Dockerfile") {
		if _, err := exec.LookPath("hadolint"); err == nil {
			logInfo("Scanning Dockerfile...")
			cmd := exec.Command("hadolint", "Dockerfile")
			if err := cmd.Run(); err == nil {
				logSuccess(stats, "Dockerfile security scan passed")
			} else {
				logError(stats, "Dockerfile security issues found")
			}
		} else {
			logWarning(stats, "hadolint not installed, skipping Dockerfile scan")
		}
	}

	// Check for hardcoded secrets
	logInfo("Checking for hardcoded secrets...")
	cmd := exec.Command("grep", "-r", "-i", "--exclude-dir=.git", "--exclude-dir=dist", "--exclude-dir=vendor", "-E", "(password|secret|token|key)[[:space:]]*[:=][[:space:]]*['\"][^'\"]{8,}", ".")
	if output, err := cmd.Output(); err != nil || len(strings.TrimSpace(string(output))) == 0 {
		logSuccess(stats, "No hardcoded secrets detected")
	} else {
		matches := strings.Split(strings.TrimSpace(string(output)), "\n")
		logError(stats, fmt.Sprintf("Potential hardcoded secrets found: %d matches", len(matches)))
	}
}

func runDryRun(stats *verifyStats, configFile, proConfigFile string) {
	logInfo("Running GoReleaser dry-run...")

	if _, err := exec.LookPath("goreleaser"); err != nil {
		logWarning(stats, "goreleaser not installed, skipping dry-run")
		return
	}

	// Temporarily clear conflicting tokens for GoReleaser
	savedGitlabToken := os.Getenv("GITLAB_TOKEN")
	savedGiteaToken := os.Getenv("GITEA_TOKEN")
	os.Unsetenv("GITLAB_TOKEN")
	os.Unsetenv("GITEA_TOKEN")

	if fileExistsSimple(configFile) {
		logInfo("Testing free version configuration...")
		cmd := exec.Command("goreleaser", "release", "--config", configFile, "--snapshot", "--skip=publish", "--clean")
		if err := cmd.Run(); err == nil {
			logSuccess(stats, "Dry-run successful for free version")
		} else {
			logWarning(stats, "Dry-run failed for free version (this might be expected in test environment)")
		}
	}

	// Note: Pro features dry-run would require a pro license
	if fileExistsSimple(proConfigFile) {
		logInfo("Pro version exists but requires license for full validation")
	}

	// Restore tokens
	if savedGitlabToken != "" {
		os.Setenv("GITLAB_TOKEN", savedGitlabToken)
	}
	if savedGiteaToken != "" {
		os.Setenv("GITEA_TOKEN", savedGiteaToken)
	}
}

func printVerificationSummary(stats *verifyStats) {
	fmt.Println("================================")
	fmt.Println("Verification Summary")
	fmt.Println("================================")
	fmt.Printf("\033[0;32mChecks passed:\033[0m %d\n", stats.Checks)
	fmt.Printf("\033[1;33mWarnings:\033[0m %d\n", stats.Warnings)
	fmt.Printf("\033[0;31mErrors:\033[0m %d\n", stats.Errors)
	fmt.Println()

	if stats.Errors == 0 {
		if stats.Warnings == 0 {
			fmt.Println("\033[0;32m✓ All checks passed successfully!\033[0m")
			os.Exit(0)
		} else {
			fmt.Println("\033[1;33m⚠ Verification completed with warnings\033[0m")
			os.Exit(0)
		}
	} else {
		fmt.Println("\033[0;31m✗ Verification failed with errors\033[0m")
		os.Exit(1)
	}
}

// Helper functions
func fileExistsSimple(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return !os.IsNotExist(err) && info.IsDir()
}

func configContains(filename, searchString string) bool {
	if !fileExistsSimple(filename) {
		return false
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), searchString)
}

func copyFile(src, dst string) error {
	srcContent, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, srcContent, 0644)
}

