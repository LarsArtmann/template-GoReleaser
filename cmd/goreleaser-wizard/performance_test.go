package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// BenchmarkProjectDetection benchmarks project detection performance.
func BenchmarkProjectDetection(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create a moderately complex project structure
	goMod := `module github.com/user/benchmark-test
go 1.21
require github.com/charmbracelet/huh v0.7.0
require charm.land/lipgloss/v2 v2.0.2
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644)

	// Create main and cmd structure
	os.MkdirAll(filepath.Join(tmpDir, "cmd", "benchmark-test"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0o644)
	os.WriteFile(
		filepath.Join(tmpDir, "cmd", "benchmark-test", "main.go"),
		[]byte("package main\n\nfunc main() {}"),
		0o644,
	)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tmpDir)

	for b.Loop() {
		config := &ProjectConfig{}
		detectProjectInfo(config)
	}
}

// BenchmarkConfigGeneration benchmarks GoReleaser config generation.
func BenchmarkConfigGeneration(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-config-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/config-benchmark
go 1.21
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0o644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tmpDir)

	config := &ProjectConfig{
		ProjectName:        "config-benchmark",
		ProjectDescription: "A benchmark test project",
		ProjectType:        domain.ProjectTypeCLI,
		BinaryName:         "config-benchmark",
		MainPath:           ".",
		Platforms:          []domain.Platform{"linux", "darwin", "windows"},
		Architectures:      []domain.Architecture{"amd64", "arm64"},
		CGOStatus:          domain.CGOStatusDisabled,
		GitProvider:        domain.GitProviderGitHub,
		DockerSupport:      domain.DockerSupportBoth,
		DockerRegistry:     domain.DockerRegistryGitHub,
		SigningLevel:       domain.SigningLevelBasic,
		Homebrew:           true,
		ActionLevel:        domain.ActionLevelBasic,
		ActionsOn:          []domain.ActionTrigger{domain.ActionTriggerVersionTags},
	}

	for b.Loop() {
		err := generateGoReleaserConfig(config)
		if err != nil {
			b.Fatalf("Config generation failed: %v", err)
		}

		os.Remove(".goreleaser.yaml") // Clean up for next iteration
	}
}

// BenchmarkGitHubActionsGeneration benchmarks GitHub Actions workflow generation.
func BenchmarkGitHubActionsGeneration(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-actions-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/actions-benchmark
go 1.21
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0o644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tmpDir)

	config := &ProjectConfig{
		ProjectName:   "actions-benchmark",
		BinaryName:    "actions-benchmark",
		ActionLevel:   domain.ActionLevelBasic,
		DockerSupport: domain.DockerSupportBoth,
		SigningLevel:  domain.SigningLevelBasic,
		ActionsOn:     []domain.ActionTrigger{domain.ActionTriggerAllTags},
	}

	for b.Loop() {
		err := generateGitHubActions(config)
		if err != nil {
			b.Fatalf("GitHub Actions generation failed: %v", err)
		}

		os.RemoveAll(".github") // Clean up for next iteration
	}
}

// BenchmarkFileOperations benchmarks file operation performance.
func BenchmarkFileOperations(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "wizard-fileops-benchmark")
	defer os.RemoveAll(tmpDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tmpDir)

	testContent := []string{
		"Small content for testing file operations",
		strings.Repeat("Larger content for testing file operations with more data. ", 50),
		strings.Repeat("Very large content for testing file operations with much more data. ", 200),
	}

	for i := 0; b.Loop(); i++ {
		content := testContent[i%len(testContent)]
		filename := fmt.Sprintf("benchmark-file-%d.txt", i)

		// Test write operation
		err := os.WriteFile(filename, []byte(content), 0o644)
		if err != nil {
			b.Fatalf("SafeFileWrite failed: %v", err)
		}

		// Test read operation
		readContent, err := os.ReadFile(filename)
		if err != nil {
			b.Fatalf("SafeReadFile failed: %v", err)
		}

		if string(readContent) != content {
			b.Fatalf("Content mismatch")
		}

		// Clean up
		os.Remove(filename)
	}
}

// TestPerformanceCharacteristics tests performance characteristics under different conditions.
func TestPerformanceCharacteristics(t *testing.T) {
	tests := []struct {
		name          string
		complexity    int
		expectedMaxMs int64
	}{
		{"simple_project", 1, 500},    // Simple project should complete in <500ms
		{"medium_project", 5, 1000},   // Medium project should complete in <1s
		{"complex_project", 10, 3000}, // Complex project should complete in <3s
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			// Create project with specified complexity
			tmpDir, _ := os.MkdirTemp("", "wizard-perf-"+tt.name)
			defer os.RemoveAll(tmpDir)

			createBenchmarkProject(t, tmpDir, tt.complexity)

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			os.Chdir(tmpDir)

			// Run full wizard workflow
			config := &ProjectConfig{
				GitProvider:    domain.GitProviderGitHub,
				ProjectType:    domain.ProjectTypeCLI,
				CGOStatus:      domain.CGOStatusDisabled,
				DockerSupport:  domain.DockerSupportNone,
				DockerRegistry: domain.DockerRegistryDockerHub,
				SigningLevel:   domain.SigningLevelNone,
				ActionLevel:    domain.ActionLevelBasic,
				FeatureLevel:   domain.FeatureLevelStandard,
				State:          domain.ConfigStateValid,
				ActionsOn:      []domain.ActionTrigger{domain.ActionTriggerVersionTags},
			}
			detectProjectInfo(config)

			err := generateGoReleaserConfig(config)
			if err != nil {
				t.Errorf("Config generation failed: %v", err)
			}

			err = generateGitHubActions(config)
			if err != nil {
				t.Errorf("GitHub Actions generation failed: %v", err)
			}

			duration := time.Since(start)

			// Check performance requirements
			if duration.Milliseconds() > tt.expectedMaxMs {
				t.Errorf("Performance exceeded threshold: %v > %dms", duration, tt.expectedMaxMs)
			}

			t.Logf("Performance: %v for %s (threshold: %dms)", duration, tt.name, tt.expectedMaxMs)
		})
	}
}

// TestMemoryUsage tests memory usage patterns.
func TestMemoryUsage(t *testing.T) {
	// Get initial memory stats
	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Create and configure multiple projects
	tmpDir, _ := os.MkdirTemp("", "wizard-memory-test")
	defer os.RemoveAll(tmpDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	for i := range 10 {
		projectDir := filepath.Join(tmpDir, fmt.Sprintf("project-%d", i))
		os.MkdirAll(projectDir, 0o755)
		os.Chdir(projectDir)

		// Create basic project
		goMod := fmt.Sprintf("module github.com/user/memory-test-%d\ngo 1.21\n", i)
		os.WriteFile("go.mod", []byte(goMod), 0o644)
		os.WriteFile("main.go", []byte("package main\n\nfunc main() {}"), 0o644)

		// Run wizard operations
		config := &ProjectConfig{}
		detectProjectInfo(config)
		generateGoReleaserConfig(config)
	}

	// Get final memory stats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Calculate memory usage
	allocDiff := m2.Alloc - m1.Alloc
	totalAllocDiff := m2.TotalAlloc - m1.TotalAlloc

	t.Logf(
		"Memory usage: Alloc diff = %d bytes, TotalAlloc diff = %d bytes",
		allocDiff,
		totalAllocDiff,
	)

	// Memory usage should be reasonable (less than 50MB for 10 projects)
	if totalAllocDiff > 50*1024*1024 {
		t.Errorf("Memory usage too high: %d bytes (> 50MB)", totalAllocDiff)
	}
}

// TestConcurrentOperations tests sequential wizard operations.
// Note: True concurrency is not tested because os.Chdir is process-wide
// and cannot be used safely across goroutines. This test verifies that
// multiple wizard operations can run in sequence without side effects.
func TestConcurrentOperations(t *testing.T) {
	// Test that wizard can handle multiple operations safely
	operations := 5
	errors := make([]error, 0)

	for i := range operations {
		// Create temporary project
		tmpDir, _ := os.MkdirTemp("", fmt.Sprintf("wizard-sequential-%d", i))

		// Create project
		goMod := fmt.Sprintf("module github.com/user/sequential-test-%d\ngo 1.21\n", i)
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644)
		os.WriteFile(
			filepath.Join(tmpDir, "main.go"),
			[]byte("package main\n\nfunc main() {}"),
			0o644,
		)

		originalDir, _ := os.Getwd()

		os.Chdir(tmpDir)

		// Run wizard operations
		config := &ProjectConfig{
			GitProvider:    domain.GitProviderGitHub,
			ProjectType:    domain.ProjectTypeCLI,
			CGOStatus:      domain.CGOStatusDisabled,
			DockerSupport:  domain.DockerSupportNone,
			DockerRegistry: domain.DockerRegistryDockerHub,
			SigningLevel:   domain.SigningLevelNone,
			ActionLevel:    domain.ActionLevelBasic,
			FeatureLevel:   domain.FeatureLevelStandard,
			State:          domain.ConfigStateValid,
			ActionsOn:      []domain.ActionTrigger{domain.ActionTriggerVersionTags},
		}
		detectProjectInfo(config)

		err := generateGoReleaserConfig(config)
		if err != nil {
			errors = append(errors, fmt.Errorf("project %d: %w", i, err))
		}

		// Cleanup
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	}

	// Check for errors
	for _, err := range errors {
		t.Errorf("Operation error: %v", err)
	}

	t.Logf("Successfully completed %d sequential operations", operations)
}

// createBenchmarkProject creates a project with specified complexity.
func createBenchmarkProject(t *testing.T, dir string, complexity int) {
	// Create basic structure
	goMod := `module github.com/user/benchmark-project
go 1.21
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0o644)

	// Add complexity based on level
	if complexity >= 1 {
		addCmdStructure(dir)
	}

	if complexity >= 3 {
		addInternalStructure(dir)
	}

	if complexity >= 5 {
		addAPIStructure(dir)
	}

	if complexity >= 7 {
		addPkgStructure(dir)
	}

	if complexity >= 10 {
		addExtensiveStructure(dir, 5)
	}
}

// addCmdStructure adds cmd directory structure.
func addCmdStructure(dir string) {
	createDirAndFiles(dir, "cmd/benchmark-project", map[string]string{
		"main.go": "package main\n\nfunc main() {}",
	})
}

// addInternalStructure adds internal/app directory structure.
func addInternalStructure(dir string) {
	createDirAndFiles(dir, "internal/app", map[string]string{
		"app.go":    "package app\n\nfunc Run() {}",
		"config.go": "package app\n\ntype Config struct {}",
	})
}

// addAPIStructure adds api/v1 directory structure.
func addAPIStructure(dir string) {
	createDirAndFiles(dir, "api/v1", map[string]string{
		"handler.go":    "package v1\n\nfunc Handle() {}",
		"middleware.go": "package v1\n\nfunc Middleware() {}",
	})
}

// addPkgStructure adds pkg/utils directory structure.
func addPkgStructure(dir string) {
	createDirAndFiles(dir, "pkg/utils", map[string]string{
		"helper.go":    "package utils\n\nfunc Helper() {}",
		"validator.go": "package utils\n\nfunc Validate() {}",
	})
}

// createDirAndFiles creates a directory and writes multiple files.
func createDirAndFiles(baseDir, dirPath string, files map[string]string) {
	fullDir := filepath.Join(baseDir, filepath.FromSlash(dirPath))
	os.MkdirAll(fullDir, 0o755)

	for filename, content := range files {
		os.WriteFile(filepath.Join(fullDir, filename), []byte(content), 0o644)
	}
}

// addExtensiveStructure adds multiple package directories.
func addExtensiveStructure(dir string, count int) {
	for i := range count {
		pkgName := fmt.Sprintf("pkg%02d", i)
		os.MkdirAll(filepath.Join(dir, pkgName), 0o755)
		os.WriteFile(
			filepath.Join(dir, pkgName, pkgName+".go"),
			fmt.Appendf(nil, "package %s\n\nfunc Func() {}", pkgName),
			0o644,
		)
	}
}
