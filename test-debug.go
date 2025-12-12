package main

import (
	"fmt"
	"log"
	"os"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

func main() {
	config := domain.NewSafeProjectConfig()
	config.ProjectName = "test-project"
	config.BinaryName = "test-project" 
	config.MainPath = "./main.go"
	config.ProjectType = domain.ProjectTypeCLI
	
	fmt.Printf("Before ApplyDefaults:\n")
	fmt.Printf("  ProjectName: %s\n", config.ProjectName)
	fmt.Printf("  ProjectType: %s\n", config.ProjectType)
	fmt.Printf("  IsValid: %t\n", config.ProjectType.IsValid())
	fmt.Printf("  Platforms: %v\n", config.Platforms)
	
	config.ApplyDefaults()
	
	fmt.Printf("\nAfter ApplyDefaults:\n")
	fmt.Printf("  ProjectName: %s\n", config.ProjectName)
	fmt.Printf("  ProjectType: %s\n", config.ProjectType)
	fmt.Printf("  IsValid: %t\n", config.ProjectType.IsValid())
	fmt.Printf("  Platforms: %v\n", config.Platforms)
	
	fmt.Printf("\nIsReadyForGeneration: %t\n", config.IsReadyForGeneration())
}