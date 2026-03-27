package main

import (
	"fmt"
)

// displayValidationResults displays validation results to the user.
func displayValidationResults(results *ValidationResults, verbose bool) {
	fmt.Println("📋 Validation Summary:")
	fmt.Println()

	displayConfigStatus(results)
	displayActionsStatus(results)
	displayProjectStatus(results)
	displayErrors(results.Errors, verbose)
	displayWarnings(results.Warnings, verbose)
	displayRecommendations(results.Recommendations)
}

func displayConfigStatus(results *ValidationResults) {
	switch {
	case !results.ConfigExists:
		fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Not found"))
	case results.ConfigValid:
		fmt.Println(successStyle.Render("✅ GoReleaser configuration: Valid"))
	default:
		fmt.Println(errorStyle.Render("❌ GoReleaser configuration: Invalid"))
	}
}

func displayActionsStatus(results *ValidationResults) {
	switch {
	case !results.ActionsExists:
		fmt.Println(infoStyle.Render("ℹ️  GitHub Actions workflow: Not found"))
	case results.ActionsValid:
		fmt.Println(successStyle.Render("✅ GitHub Actions workflow: Valid"))
	default:
		fmt.Println(errorStyle.Render("❌ GitHub Actions workflow: Invalid"))
	}
}

func displayProjectStatus(results *ValidationResults) {
	if results.ProjectValid {
		fmt.Println(successStyle.Render("✅ Project structure: Valid"))
	} else {
		fmt.Println(errorStyle.Render("❌ Project structure: Invalid"))
	}
}

func displayErrors(errors []*DomainError, verbose bool) {
	if len(errors) == 0 {
		return
	}

	fmt.Println(errorStyle.Render("❌ Errors:"))

	for _, err := range errors {
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

func displayWarnings(warnings []*DomainError, verbose bool) {
	if len(warnings) == 0 {
		return
	}

	fmt.Println(infoStyle.Render("⚠️  Warnings:"))

	for _, warning := range warnings {
		fmt.Printf("  • %s\n", warning.Message)

		if verbose {
			fmt.Printf("    Details: %s\n", warning.Details)
		}
	}

	fmt.Println()
}

func displayRecommendations(recommendations []string) {
	if len(recommendations) == 0 {
		return
	}

	fmt.Println(infoStyle.Render("💡 Recommendations:"))

	for _, rec := range recommendations {
		fmt.Printf("  • %s\n", rec)
	}

	fmt.Println()
}
