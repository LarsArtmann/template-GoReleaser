package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify project build and test status",
	Long: `Verify various aspects of your project build and testing:
- Code compilation
- Unit tests execution
- Code formatting
- Linting status
- Dependencies integrity

This command ensures your project is ready for release by
running comprehensive verification checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		runVerify(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Add flags for verify command
	verifyCmd.Flags().BoolP("build", "b", false, "Verify build compilation")
	verifyCmd.Flags().BoolP("test", "t", false, "Run tests")
	verifyCmd.Flags().BoolP("format", "f", false, "Check code formatting")
	verifyCmd.Flags().BoolP("lint", "l", false, "Run linting")
	verifyCmd.Flags().BoolP("deps", "d", false, "Verify dependencies")
	verifyCmd.Flags().BoolP("all", "a", true, "Run all verifications (default)")
}

func runVerify(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 Running verification...")

	buildFlag, _ := cmd.Flags().GetBool("build")
	testFlag, _ := cmd.Flags().GetBool("test")
	formatFlag, _ := cmd.Flags().GetBool("format")
	lintFlag, _ := cmd.Flags().GetBool("lint")
	depsFlag, _ := cmd.Flags().GetBool("deps")
	allFlag, _ := cmd.Flags().GetBool("all")

	// If no specific flags are set, default to all
	if !buildFlag && !testFlag && !formatFlag && !lintFlag && !depsFlag {
		allFlag = true
	}

	verificationsPassed := 0
	verificationsFailed := 0

	if allFlag || depsFlag {
		fmt.Println("\n📦 Verifying dependencies...")
		if verifyDependencies() {
			fmt.Println("✅ Dependencies verification passed")
			verificationsPassed++
		} else {
			fmt.Println("❌ Dependencies verification failed")
			verificationsFailed++
		}
	}

	if allFlag || formatFlag {
		fmt.Println("\n🎨 Verifying code formatting...")
		if verifyFormatting() {
			fmt.Println("✅ Code formatting verification passed")
			verificationsPassed++
		} else {
			fmt.Println("❌ Code formatting verification failed")
			verificationsFailed++
		}
	}

	if allFlag || buildFlag {
		fmt.Println("\n🔨 Verifying build...")
		if verifyBuild() {
			fmt.Println("✅ Build verification passed")
			verificationsPassed++
		} else {
			fmt.Println("❌ Build verification failed")
			verificationsFailed++
		}
	}

	if allFlag || testFlag {
		fmt.Println("\n🧪 Verifying tests...")
		if verifyTests() {
			fmt.Println("✅ Tests verification passed")
			verificationsPassed++
		} else {
			fmt.Println("❌ Tests verification failed")
			verificationsFailed++
		}
	}

	if allFlag || lintFlag {
		fmt.Println("\n🔍 Verifying linting...")
		if verifyLinting() {
			fmt.Println("✅ Linting verification passed")
			verificationsPassed++
		} else {
			fmt.Println("❌ Linting verification failed")
			verificationsFailed++
		}
	}

	fmt.Printf("\n📊 Verification Summary:\n")
	fmt.Printf("   ✅ Passed: %d\n", verificationsPassed)
	fmt.Printf("   ❌ Failed: %d\n", verificationsFailed)

	if verificationsFailed > 0 {
		fmt.Println("\n❌ Some verifications failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All verifications passed successfully!")
	}
}

func verifyDependencies() bool {
	fmt.Println("   • Running go mod verify...")

	cmd := exec.Command("go", "mod", "verify")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("   ❌ go mod verify failed: %s\n", strings.TrimSpace(string(output)))
		return false
	}

	fmt.Println("   ✅ All dependencies verified")
	return true
}

func verifyFormatting() bool {
	fmt.Println("   • Checking code formatting with gofmt...")

	cmd := exec.Command("gofmt", "-l", ".")
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("   ❌ gofmt failed: %v\n", err)
		return false
	}

	if len(strings.TrimSpace(string(output))) > 0 {
		fmt.Printf("   ❌ Unformatted files found:\n%s\n", string(output))
		fmt.Println("   💡 Run 'go fmt ./...' to fix formatting issues")
		return false
	}

	fmt.Println("   ✅ All files are properly formatted")
	return true
}

func verifyBuild() bool {
	fmt.Println("   • Building project...")

	cmd := exec.Command("go", "build", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("   ❌ Build failed: %s\n", strings.TrimSpace(string(output)))
		return false
	}

	fmt.Println("   ✅ Project builds successfully")
	return true
}

func verifyTests() bool {
	fmt.Println("   • Running tests...")

	cmd := exec.Command("go", "test", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("   ❌ Tests failed: %s\n", strings.TrimSpace(string(output)))
		return false
	}

	fmt.Printf("   ✅ All tests passed\n")
	return true
}

func verifyLinting() bool {
	fmt.Println("   • Checking for golint...")

	// Check if golint is available
	if _, err := exec.LookPath("golint"); err != nil {
		fmt.Println("   ⚠️  golint not found, skipping lint check")
		fmt.Println("   💡 Install with: go install golang.org/x/lint/golint@latest")
		return true // Don't fail if golint is not available
	}

	cmd := exec.Command("golint", "./...")
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("   ❌ Linting failed: %v\n", err)
		return false
	}

	lintOutput := strings.TrimSpace(string(output))
	if lintOutput != "" {
		fmt.Printf("   ⚠️  Lint warnings found:\n%s\n", lintOutput)
		// Don't fail on warnings, just report them
	}

	fmt.Println("   ✅ Linting completed")
	return true
}
