// Package benchmarks provides benchmark tests for the GoReleaser CLI template
package benchmarks

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/template-GoReleaser/internal/container"
	"github.com/LarsArtmann/template-GoReleaser/internal/services"
	"github.com/LarsArtmann/template-GoReleaser/internal/types"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

// BenchmarkCLIExecution benchmarks the execution of various CLI commands
func BenchmarkCLIExecution(b *testing.B) {
	tests := []struct {
		name string
		args []string
	}{
		{"Version", []string{"version"}},
		{"Help", []string{"--help"}},
		{"Validate", []string{"validate", "--help"}},
		{"Verify", []string{"verify", "--help"}},
		{"License", []string{"license", "--help"}},
		{"Config", []string{"config", "--help"}},
		{"Server", []string{"server", "--help"}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Create a new command instance
				cmd := createTestCommand()
				cmd.SetArgs(tt.args)

				// Capture output
				buf := new(bytes.Buffer)
				cmd.SetOut(buf)
				cmd.SetErr(buf)

				// Execute command
				_ = cmd.Execute()
			}
		})
	}
}

// BenchmarkContainerInitialization benchmarks DI container initialization
func BenchmarkContainerInitialization(b *testing.B) {
	b.Run("Container_Creation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			container := container.NewContainer()
			injector := container.GetInjector()

			// Perform some basic operations
			_ = injector != nil

			// Cleanup
			_ = container.Shutdown()
		}
	})

	b.Run("Container_ServiceResolution", func(b *testing.B) {
		// Setup container once
		c := container.NewContainer()
		defer c.Shutdown()
		injector := c.GetInjector()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Benchmark service resolution
			_, err := do.Invoke[services.ValidationService](injector)
			if err != nil {
				b.Fatalf("Failed to resolve ValidationService: %v", err)
			}
		}
	})
}

// BenchmarkValidationService benchmarks validation operations
func BenchmarkValidationService(b *testing.B) {
	// Setup
	c := container.NewContainer()
	defer c.Shutdown()
	injector := c.GetInjector()

	validationService, err := do.Invoke[services.ValidationService](injector)
	if err != nil {
		b.Fatalf("Failed to resolve ValidationService: %v", err)
	}

	b.Run("ValidateEnvironment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = validationService.ValidateEnvironment()
		}
	})

	b.Run("ValidateTools", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = validationService.ValidateTools()
		}
	})

	b.Run("ValidateProject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = validationService.ValidateProject()
		}
	})
}

// BenchmarkConfigService benchmarks configuration operations
func BenchmarkConfigService(b *testing.B) {
	// Setup
	c := container.NewContainer()
	defer c.Shutdown()
	injector := c.GetInjector()

	configService, err := do.Invoke[services.ConfigService](injector)
	if err != nil {
		b.Fatalf("Failed to resolve ConfigService: %v", err)
	}

	// Create a temporary directory for test configs
	tempDir := b.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	b.Run("InitConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			config, err := configService.InitConfig()
			if err != nil {
				b.Fatalf("Failed to init config: %v", err)
			}
			_ = config
		}
	})

	b.Run("LoadConfig", func(b *testing.B) {
		// Create a test config file first
		config, _ := configService.InitConfig()
		_ = configService.SaveConfig(config)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = configService.LoadConfig()
		}
	})

	b.Run("ValidateConfig", func(b *testing.B) {
		config, _ := configService.InitConfig()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = configService.ValidateConfig(config)
		}
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("Config_Creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			config := &types.Config{
				Project: types.ProjectConfig{
					Name:        "test-project",
					Description: "Test description",
					Repository:  "github.com/test/repo",
				},
				Author: types.AuthorConfig{
					Name: "Test Author",
				},
				License: types.LicenseConfig{
					Type: "MIT",
				},
			}
			_ = config
		}
	})

	b.Run("ValidationResult_Creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := &services.ValidationResult{
				Success:  true,
				Errors:   make([]string, 0),
				Warnings: make([]string, 0),
				Checks:   10,
			}
			_ = result
		}
	})

	b.Run("Large_Config_Marshal", func(b *testing.B) {
		config := createLargeTestConfig()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Simulate marshaling operations that might happen in real usage
			_ = fmt.Sprintf("%+v", config)
		}
	})
}

// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	c := container.NewContainer()
	defer c.Shutdown()
	injector := c.GetInjector()

	validationService, err := do.Invoke[services.ValidationService](injector)
	if err != nil {
		b.Fatalf("Failed to resolve ValidationService: %v", err)
	}

	b.Run("Concurrent_Validation", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = validationService.ValidateEnvironment()
			}
		})
	})

	b.Run("Concurrent_Container_Resolution", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = do.Invoke[services.ValidationService](injector)
			}
		})
	})
}

// BenchmarkFileOperations benchmarks file system operations
func BenchmarkFileOperations(b *testing.B) {
	tempDir := b.TempDir()

	b.Run("Config_File_Write", func(b *testing.B) {
		config := createTestConfig()
		c := container.NewContainer()
		defer c.Shutdown()
		injector := c.GetInjector()

		configService, _ := do.Invoke[services.ConfigService](injector)

		// Change to temp directory
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			configPath := filepath.Join(tempDir, fmt.Sprintf("config_%d.yaml", i))
			os.Setenv("GORELEASER_CLI_CONFIG", configPath)
			_ = configService.SaveConfig(config)
		}
	})

	b.Run("Config_File_Read", func(b *testing.B) {
		// Setup: Create a config file
		config := createTestConfig()
		c := container.NewContainer()
		defer c.Shutdown()
		injector := c.GetInjector()

		configService, _ := do.Invoke[services.ConfigService](injector)
		configPath := filepath.Join(tempDir, "bench_config.yaml")

		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		os.Setenv("GORELEASER_CLI_CONFIG", configPath)
		_ = configService.SaveConfig(config)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = configService.LoadConfig()
		}
	})
}

// Helper functions for benchmarks

func createTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test-cli",
		Short: "Test CLI command",
		Run: func(cmd *cobra.Command, args []string) {
			// Basic command execution
		},
	}
}

func createTestConfig() *types.Config {
	return &types.Config{
		Project: types.ProjectConfig{
			Name:        "benchmark-test",
			Description: "Benchmark test configuration",
			Repository:  "github.com/test/benchmark",
		},
		Author: types.AuthorConfig{
			Name: "Benchmark Author",
		},
		License: types.LicenseConfig{
			Type: "MIT",
		},
	}
}

func createLargeTestConfig() *types.Config {
	config := createTestConfig()
	// Add additional fields that might exist in a large configuration
	config.Project.Description = "This is a very long description that simulates a real-world configuration with extensive documentation and detailed explanations of what this project does and how it should be configured for optimal performance and maintainability."
	return config
}

// TestMain provides package-level setup and teardown for benchmarks
func TestMain(m *testing.M) {
	// Setup any global benchmark configuration
	ctx := context.Background()
	_ = ctx // Use context if needed for setup

	// Run tests and benchmarks
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

// Example benchmark runner - can be called from a separate main function
func RunBenchmarks() {
	// This could be extended to run specific benchmark suites
	fmt.Println("To run benchmarks, use: go test -bench=. ./benchmarks/")
	fmt.Println("For memory profiling: go test -bench=. -benchmem ./benchmarks/")
	fmt.Println("For CPU profiling: go test -bench=. -cpuprofile=cpu.prof ./benchmarks/")
	fmt.Println("For memory profiling: go test -bench=. -memprofile=mem.prof ./benchmarks/")
}
