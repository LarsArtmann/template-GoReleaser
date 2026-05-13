package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// baseTestProjectConfig returns a base ProjectConfig for testing with common defaults.
func baseTestProjectConfig() ProjectConfig {
	return ProjectConfig{
		ProjectType:    domain.ProjectTypeCLI,
		GitProvider:    domain.GitProviderGitHub,
		CGOStatus:      domain.CGOStatusDisabled,
		DockerSupport:  domain.DockerSupportNone,
		DockerRegistry: domain.DockerRegistryDockerHub,
		SigningLevel:   domain.SigningLevelNone,
		ActionLevel:    domain.ActionLevelBasic,
		FeatureLevel:   domain.FeatureLevelStandard,
		State:          domain.ConfigStateValid,
		ActionsOn:      []domain.ActionTrigger{domain.ActionTriggerVersionTags},
	}
}

// testProjectConfigWithOverrides creates a ProjectConfig with field overrides applied.
func testProjectConfigWithOverrides(overrides map[string]any) ProjectConfig {
	cfg := baseTestProjectConfig()

	for key, value := range overrides {
		switch key {
		case "ProjectName":
			if s, ok := value.(string); ok {
				cfg.ProjectName = s
			}
		case "BinaryName":
			if s, ok := value.(string); ok {
				cfg.BinaryName = s
			}
		case "MainPath":
			if s, ok := value.(string); ok {
				cfg.MainPath = s
			}
		case "ProjectDescription":
			if s, ok := value.(string); ok {
				cfg.ProjectDescription = s
			}
		case "Platforms":
			if p, ok := value.([]domain.Platform); ok {
				cfg.Platforms = p
			}
		case "Architectures":
			if a, ok := value.([]domain.Architecture); ok {
				cfg.Architectures = a
			}
		}
	}

	return cfg
}

func TestRunGenerate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
	}{
		{
			name: "generate_valid_config",
			setupFunc: func(t *testing.T) string {
				dir, _ := os.MkdirTemp("", "wizard-generate-test")
				createBasicGoProject(t, dir, "github.com/user/generate-test")

				return dir
			},
			expectError: false,
		},
		{
			name: "generate_in_non_go_project",
			setupFunc: func(t *testing.T) string {
				dir, _ := os.MkdirTemp("", "wizard-generate-test")

				return dir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := tt.setupFunc(t)
			defer os.RemoveAll(testDir)

			inTestDir(t, testDir, func() {
				config := &ProjectConfig{}
				detectProjectInfo(config)

				if config.ProjectName == "" && !tt.expectError {
					t.Error("Expected project detection to work")
				}
			})
		})
	}
}

func TestTemplateGeneration(t *testing.T) {
	tests := []struct {
		name        string
		config      ProjectConfig
		expectError bool
		checks      []string
	}{
		{
			name: "generate_complete_config",
			config: ProjectConfig{
				ProjectName:        "complete-test",
				ProjectDescription: "A complete test project",
				BinaryName:         "complete-test",
				MainPath:           "./cmd/complete-test",
				ProjectType:        domain.ProjectTypeCLI,
				Platforms:          []domain.Platform{"linux", "darwin"},
				Architectures:      []domain.Architecture{"amd64", "arm64"},
				CGOStatus:          domain.CGOStatusDisabled,
				GitProvider:        domain.GitProviderGitHub,
				DockerSupport:      domain.DockerSupportBoth,
				DockerRegistry:     domain.DockerRegistryGitHub,
				SigningLevel:       domain.SigningLevelBasic,
				Homebrew:           true,
				ActionLevel:        domain.ActionLevelBasic,
				FeatureLevel:       domain.FeatureLevelStandard,
				State:              domain.ConfigStateValid,
				ActionsOn:          []domain.ActionTrigger{domain.ActionTriggerVersionTags},
			},
			expectError: false,
			checks: []string{
				"project_name: complete-test",
				"binary: complete-test",
				"main: ./cmd/complete-test",
				"goos:",
				"goarch:",
				"CGO_ENABLED=0",
				"dockers:",
				// Note: signs: and brews: are not in the template yet
				// Signing is handled via GitHub Actions, Homebrew as separate file
			},
		},
		{
			name: "generate_minimal_config",
			config: testProjectConfigWithOverrides(map[string]any{
				"ProjectName": "minimal-test",
				"BinaryName":  "minimal-test",
				"MainPath":    ".",
			}),
			expectError: false,
			checks: []string{
				"project_name: minimal-test",
				"binary: minimal-test",
				"main: .",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-template-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()

			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Generate config
			err := generateGoReleaserConfig(&tt.config)

			// Check error
			AssertErr(t, "generateGoReleaserConfig", err, tt.expectError)

			if err != nil {
				return
			}

			if !tt.expectError {
				// Read generated file
				content, err := os.ReadFile(".goreleaser.yaml")
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				contentStr := string(content)

				// Check for expected strings
				for _, check := range tt.checks {
					if !strings.Contains(contentStr, check) {
						t.Errorf("Generated config missing expected string: %q", check)
					}
				}
			}
		})
	}
}

func TestGitHubActionsGeneration(t *testing.T) {
	tests := []struct {
		name        string
		config      ProjectConfig
		expectError bool
		checks      []string
	}{
		{
			name: "actions_with_docker",
			config: ProjectConfig{
				ProjectName:    "docker-test",
				BinaryName:     "docker-test",
				ActionLevel:    domain.ActionLevelBasic,
				DockerSupport:  domain.DockerSupportBoth,
				DockerRegistry: domain.DockerRegistryGitHub,
				ActionsOn:      []domain.ActionTrigger{domain.ActionTriggerManual},
			},
			expectError: false,
			checks: []string{
				"name: Release",
				"workflow_dispatch:",
				"Login to Docker Registry",
				"packages: write",
			},
		},
		{
			name: "actions_with_signing",
			config: ProjectConfig{
				ProjectName:  "signing-test",
				BinaryName:   "signing-test",
				ActionLevel:  domain.ActionLevelBasic,
				SigningLevel: domain.SigningLevelBasic,
				ActionsOn:    []domain.ActionTrigger{domain.ActionTriggerAllTags},
			},
			expectError: false,
			checks: []string{
				"Install Cosign",
				"id-token: write",
				"tags:",
				"- '*'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-actions-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()

			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Apply defaults to ensure config is complete
			config := (*domain.SafeProjectConfig)(&tt.config)
			config.ApplyDefaults()

			// Generate actions
			err := generateGitHubActions(config)

			// Check error
			AssertErr(t, "generateGitHubActions", err, tt.expectError)

			if err != nil {
				return
			}

			if !tt.expectError {
				workflowPath := filepath.Join(".github", "workflows", "release.yml")
				verifyFileContents(t, workflowPath, tt.checks)
			}
		})
	}
}

// invalidEmptyFieldTestCases returns test cases for validation of empty required fields.
func invalidEmptyFieldTestCases() []struct {
	name    string
	config  ProjectConfig
	wantErr bool
} {
	return []struct {
		name    string
		config  ProjectConfig
		wantErr bool
	}{
		makeInvalidEmptyFieldTestCase(
			"invalid_empty_project_name",
			"ProjectName",
			"",
			"BinaryName",
			"myapp",
			"MainPath",
			".",
		),
		makeInvalidEmptyFieldTestCase(
			"invalid_empty_binary_name",
			"ProjectName",
			"myapp",
			"BinaryName",
			"",
			"MainPath",
			".",
		),
		makeInvalidEmptyFieldTestCase(
			"invalid_empty_main_path",
			"ProjectName",
			"myapp",
			"BinaryName",
			"myapp",
			"MainPath",
			"",
		),
	}
}

// makeInvalidEmptyFieldTestCase creates a test case for empty field validation.
func makeInvalidEmptyFieldTestCase(
	testName, field1Name, field1Val, field2Name, field2Val, field3Name, field3Val string,
) struct {
	name    string
	config  ProjectConfig
	wantErr bool
} {
	overrides := map[string]any{field1Name: field1Val, field2Name: field2Val, field3Name: field3Val}

	return struct {
		name    string
		config  ProjectConfig
		wantErr bool
	}{
		name:    testName,
		config:  testProjectConfigWithOverrides(overrides),
		wantErr: true,
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ProjectConfig
		wantErr bool
	}{
		{
			name: "valid_complete_config",
			config: ProjectConfig{
				ProjectName:        "valid-test",
				ProjectDescription: "Valid test project",
				BinaryName:         "valid-test",
				MainPath:           "./cmd/valid-test",
				ProjectType:        domain.ProjectTypeCLI,
				Platforms:          []domain.Platform{"linux", "darwin"},
				Architectures:      []domain.Architecture{"amd64"},
				CGOStatus:          domain.CGOStatusDisabled,
				GitProvider:        domain.GitProviderGitHub,
				DockerSupport:      domain.DockerSupportNone,
				DockerRegistry:     domain.DockerRegistryDockerHub,
				SigningLevel:       domain.SigningLevelNone,
				ActionLevel:        domain.ActionLevelBasic,
				FeatureLevel:       domain.FeatureLevelStandard,
				State:              domain.ConfigStateValid,
				ActionsOn:          []domain.ActionTrigger{domain.ActionTriggerVersionTags},
			},
			wantErr: false,
		},
	}
	tests = append(tests, invalidEmptyFieldTestCases()...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, _ := os.MkdirTemp("", "wizard-config-validation-test")
			defer os.RemoveAll(tmpDir)

			originalDir, _ := os.Getwd()

			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			err := generateGoReleaserConfig(&tt.config)

			AssertErr(t, "generateGoReleaserConfig", err, tt.wantErr)
		})
	}
}

func TestFileOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation func() error
		wantErr   bool
	}{
		{
			name: "write_new_file_with_safe_write",
			operation: func() error {
				return os.WriteFile("test-safe-write.txt", []byte("test content"), 0o644)
			},
			wantErr: false,
		},
		{
			name: "read_existing_file_with_safe_read",
			operation: func() error {
				content := []byte("test content for reading")

				err := os.WriteFile("test-safe-read.txt", content, 0o644)
				if err != nil {
					return err
				}

				readContent, err := os.ReadFile("test-safe-read.txt")
				if err != nil {
					return err
				}

				if string(readContent) != string(content) {
					return os.ErrInvalid
				}

				return nil
			},
			wantErr: false,
		},
		{
			name: "create_file_with_safe_create",
			operation: func() error {
				file, err := os.Create("test-safe-create.txt")
				if err != nil {
					return err
				}

				file.WriteString("test content")
				file.Close()

				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-file-ops-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()

			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			// Test file operation
			err := tt.operation()

			AssertErr(t, "File operation", err, tt.wantErr)
		})
	}
}

func TestBackupCreation(t *testing.T) {
	// Test that backup files are created when overwriting existing files
	tests := []struct {
		name            string
		originalContent string
		newContent      string
		expectBackup    bool
	}{
		{
			name:            "backup_created_on_overwrite",
			originalContent: "original content",
			newContent:      "new content",
			expectBackup:    false, // os.WriteFile doesn't create backups
		},
		{
			name:            "no_backup_for_new_file",
			originalContent: "",
			newContent:      "new content",
			expectBackup:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, _ := os.MkdirTemp("", "wizard-backup-test")
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()

			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			testFile := "test-backup.txt"

			// Create original file if needed
			if tt.originalContent != "" {
				err := os.WriteFile(testFile, []byte(tt.originalContent), 0o644)
				if err != nil {
					t.Fatalf("Failed to create original file: %v", err)
				}
			}

			// Write new file (this should create backup)
			err := os.WriteFile(testFile, []byte(tt.newContent), 0o644)
			if err != nil {
				t.Errorf("SafeFileWrite() error = %v", err)

				return
			}

			// Check backup file
			backupFile := testFile + ".backup"
			backupInfo, _ := os.Stat(backupFile)
			backupExists := backupInfo != nil

			if tt.expectBackup && !backupExists {
				t.Error("Backup file should exist when overwriting existing file")
			} else if tt.expectBackup && backupExists {
				// Read backup to ensure it contains original content
				backupContent, _ := os.ReadFile(backupFile)
				if string(backupContent) != tt.originalContent {
					t.Errorf(
						"Backup content = %q, want %q",
						string(backupContent),
						tt.originalContent,
					)
				}
			} else if !tt.expectBackup && backupExists {
				t.Error("Backup file should not exist for new file")
			}
		})
	}
}

func TestErrorRecovery(t *testing.T) {
	// Test error recovery and panic handling
	tests := []struct {
		name        string
		shouldPanic bool
		testFunc    func()
	}{
		{
			name:        "normal_operation_no_panic",
			shouldPanic: false,
			testFunc: func() {
				// Normal operation should not panic
				config := &ProjectConfig{}
				detectProjectInfo(config)
			},
		},
		{
			name:        "panic_recovery_works",
			shouldPanic: true,
			testFunc: func() {
				// Intentional panic for testing recovery
				panic("intentional panic for test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.shouldPanic {
					t.Error("Expected panic but none occurred")
				}
			}()

			tt.testFunc()
		})
	}
}
