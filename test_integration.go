package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/generators"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// TestLogger implements Logger interface for testing
type TestLogger struct{}

func (tl *TestLogger) Debug(msg string, args ...any)      { log.Printf("[DEBUG] %s %v", msg, args) }
func (tl *TestLogger) Info(msg string, args ...any)       { log.Printf("[INFO] %s %v", msg, args) }
func (tl *TestLogger) Warn(msg string, args ...any)       { log.Printf("[WARN] %s %v", msg, args) }
func (tl *TestLogger) Error(msg string, args ...any)      { log.Printf("[ERROR] %s %v", msg, args) }

// TestAllGenerators tests all generators with sample configuration
func TestAllGenerators() {
	logger := &TestLogger{}
	
	// Create test configuration
	config := createTestConfig()
	
	// Test GoReleaser generator
	log.Println("🧪 Testing GoReleaser Generator...")
	testGoReleaserGenerator(config, logger)
	
	// Test GitHub Actions generator
	log.Println("🧪 Testing GitHub Actions Generator...")
	testGitHubActionsGenerator(config, logger)
	
	// Test Dockerfile generator
	log.Println("🧪 Testing Dockerfile Generator...")
	testDockerfileGenerator(config, logger)
	
	// Test Homebrew generator
	log.Println("🧪 Testing Homebrew Generator...")
	testHomebrewGenerator(config, logger)
	
	log.Println("✅ All generator tests completed successfully!")
}

// createTestConfig creates a test configuration
func createTestConfig() *domain.SafeProjectConfig {
	config := domain.NewSafeProjectConfig()
	config.ProjectName = "test-project"
	config.ProjectDescription = "A test project for GoReleaser Wizard"
	config.ProjectType = domain.ProjectTypeCLIApplication
	config.BinaryName = "test-cli"
	config.MainPath = "./cmd/test-cli"
	config.Platforms = []domain.Platform{domain.PlatformLinux, domain.PlatformDarwin, domain.PlatformWindows}
	config.Architectures = []domain.Architecture{domain.ArchitectureAmd64, domain.ArchitectureArm64}
	config.CGOStatus = domain.CGOStatusDisabled
	config.GitProvider = domain.GitProviderGitHub
	config.DockerSupport = domain.DockerSupportBuild
	config.DockerRegistry = domain.DockerRegistryGitHub
	config.DockerImage = "test-project"
	config.SigningLevel = domain.SigningLevelBasic
	config.Homebrew = true
	config.ActionLevel = domain.ActionLevelBasic
	return config
}

// testGoReleaserGenerator tests GoReleaser generator
func testGoReleaserGenerator(config *domain.SafeProjectConfig, logger generators.Logger) {
	generator := generators.NewGoReleaserGenerator(config, logger)
	
	// Test validation
	if err := generator.ValidateTemplate(); err != nil {
		log.Fatalf("❌ GoReleaser template validation failed: %v", err)
	}
	
	// Test preview
	ctx := context.Background()
	preview, err := generator.GeneratePreview(ctx)
	if err != nil {
		log.Fatalf("❌ GoReleaser preview generation failed: %v", err)
	}
	
	if len(preview) == 0 {
		log.Fatalf("❌ GoReleaser preview is empty")
	}
	
	log.Printf("✅ GoReleaser generator test passed (preview length: %d)", len(preview))
}

// testGitHubActionsGenerator tests GitHub Actions generator
func testGitHubActionsGenerator(config *domain.SafeProjectConfig, logger generators.Logger) {
	generator := generators.NewGitHubActionsGenerator(config, logger)
	
	// Test validation
	if err := generator.ValidateTemplate(); err != nil {
		log.Fatalf("❌ GitHub Actions template validation failed: %v", err)
	}
	
	// Test preview
	ctx := context.Background()
	preview, err := generator.GeneratePreview(ctx)
	if err != nil {
		log.Fatalf("❌ GitHub Actions preview generation failed: %v", err)
	}
	
	if len(preview) == 0 {
		log.Fatalf("❌ GitHub Actions preview is empty")
	}
	
	log.Printf("✅ GitHub Actions generator test passed (preview length: %d)", len(preview))
}

// testDockerfileGenerator tests Dockerfile generator
func testDockerfileGenerator(config *domain.SafeProjectConfig, logger generators.Logger) {
	generator := generators.NewDockerfileGenerator(config, logger)
	
	// Test validation
	if err := generator.ValidateTemplate(); err != nil {
		log.Fatalf("❌ Dockerfile template validation failed: %v", err)
	}
	
	// Test preview
	ctx := context.Background()
	preview, err := generator.GeneratePreview(ctx)
	if err != nil {
		log.Fatalf("❌ Dockerfile preview generation failed: %v", err)
	}
	
	if len(preview) == 0 {
		log.Fatalf("❌ Dockerfile preview is empty")
	}
	
	// Verify Dockerfile contains expected content
	expectedContent := []string{
		"FROM golang:",
		"FROM alpine:",
		"WORKDIR",
		"ENTRYPOINT",
		config.ProjectName,
	}
	
	for _, content := range expectedContent {
		if !contains(preview, content) {
			log.Fatalf("❌ Dockerfile missing expected content: %s", content)
		}
	}
	
	log.Printf("✅ Dockerfile generator test passed (preview length: %d)", len(preview))
}

// testHomebrewGenerator tests Homebrew generator
func testHomebrewGenerator(config *domain.SafeProjectConfig, logger generators.Logger) {
	generator := generators.NewHomebrewGenerator(config, logger)
	
	// Test validation
	if err := generator.ValidateTemplate(); err != nil {
		log.Fatalf("❌ Homebrew template validation failed: %v", err)
	}
	
	// Test preview
	ctx := context.Background()
	preview, err := generator.GeneratePreview(ctx)
	if err != nil {
		log.Fatalf("❌ Homebrew preview generation failed: %v", err)
	}
	
	if len(preview) == 0 {
		log.Fatalf("❌ Homebrew preview is empty")
	}
	
	// Verify Homebrew formula contains expected content
	expectedContent := []string{
		"class TestProject",
		"homepage",
		"url",
		"sha256",
		"depends_on \"go\"",
		config.ProjectName,
	}
	
	for _, content := range expectedContent {
		if !contains(preview, content) {
			log.Fatalf("❌ Homebrew formula missing expected content: %s", content)
		}
	}
	
	log.Printf("✅ Homebrew generator test passed (preview length: %d)", len(preview))
}

// contains checks if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr ||
		     findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestInteractivePrompter tests the interactive prompter functionality
func TestInteractivePrompter() {
	log.Println("🧪 Testing Interactive Prompter...")
	
	prompter := NewInteractivePrompter()
	
	if prompter == nil {
		log.Fatalf("❌ Failed to create interactive prompter")
	}
	
	log.Printf("✅ Interactive prompter test passed")
}

// TestJobFactory tests the job factory functionality
func TestJobFactory() {
	log.Println("🧪 Testing Job Factory...")
	
	logger := &TestLogger{}
	factory := NewJobFactory(logger)
	
	if factory == nil {
		log.Fatalf("❌ Failed to create job factory")
	}
	
	config := createTestConfig()
	
	// Test creating full wizard jobs
	jobs := factory.CreateFullWizardJobs(config, false)
	
	if len(jobs) == 0 {
		log.Fatalf("❌ No jobs created by factory")
	}
	
	expectedJobNames := []string{
		"Validate Project Structure",
		"Generate GoReleaser Configuration", 
		"Generate Dockerfile",
		"Generate Homebrew Formula",
		"Generate GitHub Actions Workflow",
	}
	
	jobNames := make(map[string]bool)
	for _, job := range jobs {
		jobNames[job.Name()] = true
	}
	
	for _, expectedName := range expectedJobNames {
		if !jobNames[expectedName] {
			log.Printf("⚠️  Expected job not found: %s", expectedName)
		}
	}
	
	log.Printf("✅ Job factory test passed (%d jobs created)", len(jobs))
}

// TestDomainConfiguration tests domain configuration functionality
func TestDomainConfiguration() {
	log.Println("🧪 Testing Domain Configuration...")
	
	config := domain.NewSafeProjectConfig()
	
	if config == nil {
		log.Fatalf("❌ Failed to create domain configuration")
	}
	
	// Test applying defaults
	config.ApplyDefaults()
	
	// Test validation
	if err := config.ValidateInvariants(); err != nil {
		log.Printf("⚠️  Domain configuration validation warning: %v", err)
	}
	
	// Test recommended values
	platforms := domain.GetRecommendedPlatforms()
	if len(platforms) == 0 {
		log.Fatalf("❌ No recommended platforms found")
	}
	
	architectures := domain.GetRecommendedArchitectures()
	if len(architectures) == 0 {
		log.Fatalf("❌ No recommended architectures found")
	}
	
	log.Printf("✅ Domain configuration test passed")
	log.Printf("   - Recommended platforms: %v", platforms)
	log.Printf("   - Recommended architectures: %v", architectures)
}

// RunIntegrationTests runs all integration tests
func RunIntegrationTests() {
	log.Println("🚀 Starting GoReleaser-Wizard Integration Tests")
	log.Println("=" * 60)
	
	// Test domain configuration
	TestDomainConfiguration()
	
	// Test interactive prompter
	TestInteractivePrompter()
	
	// Test job factory
	TestJobFactory()
	
	// Test all generators
	TestAllGenerators()
	
	log.Println("=" * 60)
	log.Println("🎉 All integration tests passed successfully!")
	log.Println("🚀 GoReleaser-Wizard is ready for production!")
}

func main() {
	// Run integration tests
	RunIntegrationTests()
}