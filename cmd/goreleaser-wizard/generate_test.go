package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

func TestGenerateGoReleaserConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  ProjectConfig
		wantErr bool
		checks  []string // strings that should be in the output
	}{
		{
			name: "basic_config",
			config: ProjectConfig{
				ProjectName:        "test-app",
				ProjectDescription: "A test application",
				BinaryName:         "test-app",
				MainPath:           ".",
				Platforms:          []domain.Platform{"linux", "darwin"},
				Architectures:      []domain.Architecture{"amd64", "arm64"},
				CGOStatus:          domain.CGOStatusDisabled,
				GitProvider:        domain.GitProviderGitHub,
			},
			wantErr: false,
			checks: []string{
				"project_name: test-app",
				"binary: test-app",
				"- linux",
				"- darwin",
				"- amd64",
				"- arm64",
				"CGO_ENABLED=0",
				`owner: "owner"`, // Template renders to actual value
				`name: "repo"`,   // Template renders to actual value
			},
		},
		{
			name: "docker_enabled",
			config: ProjectConfig{
				ProjectName:    "docker-app",
				BinaryName:     "docker-app",
				MainPath:       "./cmd/app",
				DockerSupport:  domain.DockerSupportBoth,
				DockerRegistry: domain.DockerRegistryGitHub,
				GitProvider:    domain.GitProviderGitHub,
			},
			wantErr: false,
			checks: []string{
				"dockers:",
				"image_templates:",
				"ghcr.io/testuser/docker-app:{{.Tag}}",
				"dockerfile: Dockerfile",
			},
		},
		{
			name: "signing_enabled",
			config: ProjectConfig{
				ProjectName:  "signed-app",
				BinaryName:   "signed-app",
				MainPath:     ".",
				SigningLevel: domain.SigningLevelBasic,
				GitProvider:  domain.GitProviderGitHub,
			},
			wantErr: false,
			checks: []string{
				"signs:",
				"cmd: cosign",
				"certificate:",
			},
		},
		{
			name: "homebrew_enabled",
			config: ProjectConfig{
				ProjectName:        "brew-app",
				ProjectDescription: "App with Homebrew support",
				BinaryName:         "brew-app",
				MainPath:           ".",
				Homebrew:           true,
				GitProvider:        "GitHub",
			},
			wantErr: false,
			checks: []string{
				"brews:",
				"repository:",
				"folder: Formula",
				"App with Homebrew support",
			},
		},
		{
			name: "missing_project_name",
			config: ProjectConfig{
				BinaryName: "test",
				MainPath:   ".",
			},
			wantErr: true,
		},
		{
			name: "missing_binary_name",
			config: ProjectConfig{
				ProjectName: "test",
				MainPath:    ".",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, err := os.MkdirTemp("", "goreleaser-wizard-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalDir)

			// Apply defaults to ensure config is complete
			config := (*domain.SafeProjectConfig)(&tt.config)
			config.ApplyDefaults()

			// Generate config
			err = generateGoReleaserConfig(config)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("generateGoReleaserConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Read generated file
				content, err := os.ReadFile(".goreleaser.yaml")
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				contentStr := string(content)

				// Debug: optionally print content for debugging
				// fmt.Printf("Generated content (full):\n%s\n", contentStr)

				// Check for expected strings
				for _, check := range tt.checks {
					if !strings.Contains(contentStr, check) {
						t.Errorf("Generated config missing expected string: %q", check)
					}
				}

				// Basic YAML structure checks
				if !strings.HasPrefix(contentStr, "# GoReleaser configuration") {
					t.Error("Config should start with comment header")
				}
				if !strings.Contains(contentStr, "version: 2") {
					t.Error("Config should specify version 2")
				}
			}
		})
	}
}

func TestGenerateGitHubActions(t *testing.T) {
	tests := []struct {
		name    string
		config  ProjectConfig
		wantErr bool
		checks  []string
	}{
		{
			name: "basic_actions",
			config: ProjectConfig{
				ProjectName: "test-app",
				BinaryName:  "test-app",
				ActionLevel: domain.ActionLevelBasic,
				ActionsOn:   []domain.ActionTrigger{domain.ActionTriggerVersionTags},
				GitProvider: domain.GitProviderGitHub,
			},
			wantErr: false,
			checks: []string{
				"name: Release",
				"tags:",
				"- 'v*'",
				"uses: goreleaser/goreleaser-action@v6",
				"GITHUB_TOKEN:",
				"GITHUB_OWNER:",
				"GITHUB_REPO:",
			},
		},
		{
			name: "docker_support",
			config: ProjectConfig{
				ProjectName:    "docker-app",
				DockerSupport:  domain.DockerSupportBoth,
				DockerRegistry: domain.DockerRegistryGitHub,
				ActionLevel:    domain.ActionLevelBasic,
				ActionsOn:      []domain.ActionTrigger{domain.ActionTriggerManual},
				GitProvider:    domain.GitProviderGitHub,
			},
			wantErr: false,
			checks: []string{
				"workflow_dispatch:",
				"Login to Docker Registry",
				"packages: write",
			},
		},
		{
			name: "signing_support",
			config: ProjectConfig{
				ProjectName:  "signed-app",
				SigningLevel: domain.SigningLevelBasic,
				ActionLevel:  domain.ActionLevelBasic,
				ActionsOn:    []domain.ActionTrigger{domain.ActionTriggerAllTags},
				GitProvider:  domain.GitProviderGitHub,
			},
			wantErr: false,
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
			tmpDir, err := os.MkdirTemp("", "goreleaser-wizard-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalDir)

			// Apply defaults to ensure config is complete
			config := (*domain.SafeProjectConfig)(&tt.config)
			config.ApplyDefaults()

			// Generate actions
			err = generateGitHubActions(config)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("generateGitHubActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Read generated file
				workflowPath := filepath.Join(".github", "workflows", "release.yml")
				content, err := os.ReadFile(workflowPath)
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				contentStr := string(content)

				// Check for expected strings
				for _, check := range tt.checks {
					if !strings.Contains(contentStr, check) {
						t.Errorf("Generated workflow missing expected string: %q", check)
					}
				}
			}
		})
	}
}

func TestDetectProjectInfo(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected ProjectConfig
	}{
		{
			name: "simple_project",
			setup: func(dir string) error {
				// Create go.mod
				goMod := `module github.com/user/myapp
go 1.21`
				if err := os.WriteFile("go.mod", []byte(goMod), 0o644); err != nil {
					return err
				}
				// Create main.go
				return os.WriteFile("main.go", []byte("package main"), 0o644)
			},
			expected: ProjectConfig{
				ProjectName: "myapp",
				MainPath:    ".",
				BinaryName:  "myapp",
				ProjectType: "CLI Application",
			},
		},
		{
			name: "cmd_structure",
			setup: func(dir string) error {
				// Create go.mod
				goMod := `module github.com/user/complexapp
go 1.21`
				if err := os.WriteFile("go.mod", []byte(goMod), 0o644); err != nil {
					return err
				}
				// Create cmd/complexapp/main.go
				if err := os.MkdirAll("cmd/complexapp", 0o755); err != nil {
					return err
				}
				return os.WriteFile("cmd/complexapp/main.go", []byte("package main"), 0o644)
			},
			expected: ProjectConfig{
				ProjectName: "complexapp",
				MainPath:    "./cmd/complexapp",
				BinaryName:  "complexapp",
				ProjectType: "CLI Application",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, err := os.MkdirTemp("", "goreleaser-wizard-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Change to temp directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalDir)

			// Setup test environment
			if err := tt.setup(tmpDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Test detection
			config := &ProjectConfig{}
			detectProjectInfo(config)

			// Check results
			if config.ProjectName != tt.expected.ProjectName {
				t.Errorf("ProjectName = %q, want %q", config.ProjectName, tt.expected.ProjectName)
			}
			if config.MainPath != tt.expected.MainPath {
				t.Errorf("MainPath = %q, want %q", config.MainPath, tt.expected.MainPath)
			}
			if config.BinaryName != tt.expected.BinaryName {
				t.Errorf("BinaryName = %q, want %q", config.BinaryName, tt.expected.BinaryName)
			}
			if config.ProjectType != tt.expected.ProjectType {
				t.Errorf("ProjectType = %q, want %q", config.ProjectType, tt.expected.ProjectType)
			}
		})
	}
}
