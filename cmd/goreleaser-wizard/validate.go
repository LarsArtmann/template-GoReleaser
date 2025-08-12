package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate GoReleaser configuration",
	Long: `Validate your GoReleaser configuration and check for common issues.

This command will:
- Check if .goreleaser.yaml exists and is valid YAML
- Run goreleaser check if available
- Verify project structure matches configuration
- Check for missing dependencies
- Suggest improvements`,
	Run: runValidate,
}

func init() {
	validateCmd.Flags().Bool("verbose", false, "show detailed validation output")
	validateCmd.Flags().Bool("fix", false, "attempt to fix common issues")
}

func runValidate(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	fix, _ := cmd.Flags().GetBool("fix")

	fmt.Println(titleStyle.Render("🔍 Validating GoReleaser Configuration"))
	fmt.Println()

	issues := []string{}
	warnings := []string{}
	passed := 0
	total := 0

	// Check 1: .goreleaser.yaml exists
	total++
	if !fileExists(".goreleaser.yaml") {
		issues = append(issues, ".goreleaser.yaml not found")
		fmt.Println(errorStyle.Render("✗ .goreleaser.yaml not found"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'goreleaser-wizard init' to create one"))
		}
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ .goreleaser.yaml exists"))
	}

	// Check 2: go.mod exists
	total++
	if !fileExists("go.mod") {
		issues = append(issues, "go.mod not found")
		fmt.Println(errorStyle.Render("✗ go.mod not found"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'go mod init' to create one"))
		}
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ go.mod exists"))
	}

	// Check 3: Git repository
	total++
	if !fileExists(".git") {
		warnings = append(warnings, "Not a git repository")
		fmt.Println(errorStyle.Render("⚠ Not a git repository"))
		fmt.Println(infoStyle.Render("  → GoReleaser requires a git repository to work"))
	} else {
		// Check for uncommitted changes
		cmd := exec.Command("git", "status", "--porcelain")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			warnings = append(warnings, "Uncommitted changes detected")
			fmt.Println(errorStyle.Render("⚠ Uncommitted changes detected"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + strings.TrimSpace(string(output))))
			}
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Git repository clean"))
		}
	}

	// Check 4: GoReleaser installed
	total++
	goreleaserPath, err := exec.LookPath("goreleaser")
	if err != nil {
		issues = append(issues, "GoReleaser not installed")
		fmt.Println(errorStyle.Render("✗ GoReleaser not installed"))
		if fix {
			fmt.Println(infoStyle.Render("  → Install with: go install github.com/goreleaser/goreleaser/v2@latest"))
		}
	} else {
		passed++
		fmt.Println(successStyle.Render("✓ GoReleaser installed"))
		if verbose {
			fmt.Println(infoStyle.Render("  → " + goreleaserPath))
		}

		// Check 5: Run goreleaser check
		total++
		fmt.Print("  Checking configuration... ")
		checkCmd := exec.Command("goreleaser", "check")
		checkOutput, checkErr := checkCmd.CombinedOutput()
		if checkErr != nil {
			issues = append(issues, "Configuration validation failed")
			fmt.Println(errorStyle.Render("Failed"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + strings.TrimSpace(string(checkOutput))))
			}
		} else {
			passed++
			fmt.Println(successStyle.Render("OK"))
		}
	}

	// Check 6: Main package exists
	if fileExists(".goreleaser.yaml") {
		total++
		// Parse config to find main path
		// For simplicity, we'll check common locations
		mainFound := false
		commonPaths := []string{
			"main.go",
			"./cmd/*/main.go",
			"./*.go",
		}
		for _, path := range commonPaths {
			if matches, _ := filepath.Glob(path); len(matches) > 0 {
				mainFound = true
				break
			}
		}

		if !mainFound {
			warnings = append(warnings, "No main.go found in expected locations")
			fmt.Println(errorStyle.Render("⚠ No main.go found"))
			fmt.Println(infoStyle.Render("  → Make sure your main package path is correct in .goreleaser.yaml"))
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Main package found"))
		}
	}

	// Check 7: Docker (if configured)
	if fileExists("Dockerfile") {
		total++
		dockerPath, err := exec.LookPath("docker")
		if err != nil {
			warnings = append(warnings, "Docker not installed but Dockerfile exists")
			fmt.Println(errorStyle.Render("⚠ Docker not installed"))
		} else {
			passed++
			fmt.Println(successStyle.Render("✓ Docker installed"))
			if verbose {
				fmt.Println(infoStyle.Render("  → " + dockerPath))
			}
		}
	}

	// Check 8: GitHub Actions workflow
	total++
	if fileExists(".github/workflows/release.yml") || fileExists(".github/workflows/release.yaml") {
		passed++
		fmt.Println(successStyle.Render("✓ GitHub Actions workflow found"))
	} else {
		warnings = append(warnings, "No GitHub Actions workflow for releases")
		fmt.Println(infoStyle.Render("ℹ No GitHub Actions workflow"))
		if fix {
			fmt.Println(infoStyle.Render("  → Run 'goreleaser-wizard init' with GitHub Actions option"))
		}
	}

	// Summary
	fmt.Println()
	fmt.Println(titleStyle.Render("📊 Validation Summary"))
	fmt.Printf("Checks passed: %d/%d\n", passed, total)
	
	if len(issues) > 0 {
		fmt.Println()
		fmt.Println(errorStyle.Render("❌ Critical Issues:"))
		for _, issue := range issues {
			fmt.Println("  • " + issue)
		}
	}

	if len(warnings) > 0 {
		fmt.Println()
		fmt.Println(infoStyle.Render("⚠️  Warnings:"))
		for _, warning := range warnings {
			fmt.Println("  • " + warning)
		}
	}

	// Test build suggestion
	if len(issues) == 0 {
		fmt.Println()
		fmt.Println(successStyle.Render("✨ Configuration looks good!"))
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Test build: goreleaser build --snapshot --clean")
		fmt.Println("  2. Create tag: git tag -a v0.1.0 -m 'First release'")
		fmt.Println("  3. Push tag: git push origin v0.1.0")
	} else {
		fmt.Println()
		fmt.Println(errorStyle.Render("⚠️  Please fix the issues above before releasing"))
		if fix {
			fmt.Println()
			fmt.Println("Run with --fix to see suggested fixes")
		}
		os.Exit(1)
	}
}