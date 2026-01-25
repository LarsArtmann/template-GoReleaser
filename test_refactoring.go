package main

import (
	"fmt"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/validation"
)

func main() {
	fmt.Println("Testing refactored validation code...")

	// Test 1: Docker support validation with incompatible project type
	fmt.Println("\n=== Test 1: Docker support with incompatible project type ===")
	config1 := &domain.SafeProjectConfig{
		ProjectType:  domain.ProjectTypeCLI,
		DockerSupport: domain.DockerSupportBoth,
	}
	result1, _ := validation.ValidateConfiguration(config1)
	hasDockerError := false
	for _, err := range result1.Errors {
		fmt.Printf("Error: Field=%s, Message=%s\n", err.Field, err.Message)
		if err.Field == "docker_support" {
			hasDockerError = true
			fmt.Printf("✓ Docker error found: %s\n", err.Message)
		}
	}
	if !hasDockerError {
		fmt.Println("✗ Docker error expected but not found")
	}
	fmt.Printf("Total errors: %d\n", len(result1.Errors))

	// Test 2: Actions with incompatible Git provider
	fmt.Println("\n=== Test 2: Actions with incompatible Git provider ===")
	config2 := &domain.SafeProjectConfig{
		GitProvider: domain.GitProviderBitbucket,
		ActionLevel: domain.ActionLevelStandard,
	}
	result2, _ := validation.ValidateConfiguration(config2)
	hasActionsError := false
	for _, err := range result2.Errors {
		if err.Field == "action_level" {
			hasActionsError = true
			fmt.Printf("✓ Actions error found: %s\n", err.Message)
		}
	}
	if !hasActionsError {
		fmt.Println("✗ Actions error expected but not found")
	}
	fmt.Printf("Total errors: %d\n", len(result2.Errors))

	fmt.Println("\n=== All tests completed ===")
}
