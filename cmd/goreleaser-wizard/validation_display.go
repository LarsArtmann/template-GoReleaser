package main

import (
	"fmt"
)

// displayValidationResults displays validation results to the user.
func displayValidationResults(results *ValidationResults, verbose bool) {
	fmt.Println("📋 Validation Summary:")
	fmt.Println()

	// Configuration status
	if results.ConfigExists {
		if results.ConfigValid {
			fmt.Println(successStyle.Render("✅ GoReleaser configuration: Valid"))
		} else {
			fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Invalid"))
		}
	} else {
		fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Not found"))
	}

	// GitHub Actions status
	if results.ActionsExists {
		if results.ActionsValid {
			fmt.Println(successStyle.Render("✅ GitHub Actions workflow: Valid"))
		} else {
			fmt.Println(errorStyle.Render("❌ GitHub Actions workflow: Invalid"))
		}
	} else {
		fmt.Println(infoStyle.Render("ℹ️  GitHub Actions workflow: Not found"))
	}

	// Project structure status
	if results.ProjectValid {
		fmt.Println(successStyle.Render("✅ Project structure: Valid"))
	} else {
		fmt.Println(errorStyle.Render("❌ Project structure: Invalid"))
	}

	// Display errors
	if len(results.Errors) > 0 {
		fmt.Println(errorStyle.Render("❌ Errors:"))

		for _, err := range results.Errors {
			fmt.Printf("  • %s\n", err.Message)

			if verbose {
				fmt.Printf("    Details: %s\n", err.Details)

				if err.Context != "" {
					fmt.Printf("    Context: %s\n", err.Context)
				}
			}
		}

		fmt.Println()
	}

	// Display warnings
	if len(results.Warnings) > 0 {
		fmt.Println(infoStyle.Render("⚠️  Warnings:"))

		for _, warning := range results.Warnings {
			fmt.Printf("  • %s\n", warning.Message)

			if verbose {
				fmt.Printf("    Details: %s\n", warning.Details)
			}
		}

		fmt.Println()
	}

	// Display recommendations
	if len(results.Recommendations) > 0 {
		fmt.Println(infoStyle.Render("💡 Recommendations:"))

		for _, rec := range results.Recommendations {
			fmt.Printf("  • %s\n", rec)
		}

		fmt.Println()
	}
}
