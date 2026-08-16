package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"charm.land/log/v2"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// quietLogger returns a logger that discards output so tests stay readable.
func quietLogger() *log.Logger {
	return log.New(io.Discard)
}

func dockerTestConfig(t *testing.T, docker domain.DockerSupport, projectType domain.ProjectType) *ProjectConfig {
	t.Helper()

	cfg := baseTestProjectConfig()
	cfg.ProjectName = "docker-test"
	cfg.BinaryName = "docker-test"
	cfg.MainPath = "."
	cfg.ProjectType = projectType
	cfg.DockerSupport = docker
	cfg.DockerRegistry = domain.DockerRegistryGitHub
	cfg.Platforms = []domain.Platform{"linux"}
	cfg.Architectures = []domain.Architecture{"amd64"}

	return &cfg
}

func TestGenerationTargets(t *testing.T) {
	tests := []struct {
		name           string
		docker         domain.DockerSupport
		includeActions bool
		want           []string
	}{
		{
			name:           "config_only_no_docker",
			docker:         domain.DockerSupportNone,
			includeActions: false,
			want:           []string{".goreleaser.yaml"},
		},
		{
			name:           "full_wizard_no_docker",
			docker:         domain.DockerSupportNone,
			includeActions: true,
			want:           []string{".goreleaser.yaml", ".github/workflows/release.yml"},
		},
		{
			name:           "full_wizard_with_docker",
			docker:         domain.DockerSupportBoth,
			includeActions: true,
			want:           []string{".goreleaser.yaml", ".github/workflows/release.yml", "Dockerfile"},
		},
		{
			name:           "config_only_with_docker",
			docker:         domain.DockerSupportBoth,
			includeActions: false,
			want:           []string{".goreleaser.yaml", "Dockerfile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := dockerTestConfig(t, tt.docker, domain.ProjectTypeCLI)

			got := generationTargets(cfg, tt.includeActions)
			if len(got) != len(tt.want) {
				t.Fatalf("generationTargets() = %v, want %v", got, tt.want)
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("generationTargets()[%d] = %q, want %q", i, got[i], want)
				}
			}
		})
	}
}

func TestGenerationPreflightJobExecute(t *testing.T) {
	tests := []struct {
		name       string
		force      bool
		existing   []string
		targets    []string
		wantErr    bool
		wantInText []string
	}{
		{
			name:    "no_existing_files",
			targets: []string{".goreleaser.yaml", "Dockerfile"},
		},
		{
			name:       "existing_target_without_force",
			targets:    []string{".goreleaser.yaml", "Dockerfile"},
			existing:   []string{".goreleaser.yaml"},
			wantErr:    true,
			wantInText: []string{".goreleaser.yaml", "--force"},
		},
		{
			name:       "existing_target_with_force",
			targets:    []string{".goreleaser.yaml", "Dockerfile"},
			existing:   []string{".goreleaser.yaml"},
			force:      true,
			wantErr:    false,
			wantInText: nil,
		},
		{
			name:       "all_existing_targets_listed",
			targets:    []string{".goreleaser.yaml", "Dockerfile"},
			existing:   []string{".goreleaser.yaml", "Dockerfile"},
			wantErr:    true,
			wantInText: []string{".goreleaser.yaml", "Dockerfile"},
		},
		{
			name:     "file_outside_targets_ignored",
			targets:  []string{".goreleaser.yaml"},
			existing: []string{"README.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Chdir(tmpDir)

			for _, file := range tt.existing {
				if err := os.WriteFile(file, []byte("existing"), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			job := NewGenerationPreflightJob(tt.force, tt.targets, quietLogger())
			err := job.Execute(context.Background())

			AssertErr(t, "GenerationPreflightJob.Execute", err, tt.wantErr)

			if err != nil {
				for _, want := range tt.wantInText {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("error %q does not contain %q", err.Error(), want)
					}
				}
			}
		})
	}
}

func TestDockerfileGenerationJobExecute(t *testing.T) {
	t.Run("docker_disabled_skips_generation", func(t *testing.T) {
		t.Chdir(t.TempDir())

		job := NewDockerfileGenerationJob(
			dockerTestConfig(t, domain.DockerSupportNone, domain.ProjectTypeCLI),
			false,
			quietLogger(),
		)
		if err := job.Execute(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat("Dockerfile"); !os.IsNotExist(err) {
			t.Error("Dockerfile should not be generated when Docker is disabled")
		}
	})

	t.Run("cli_scratch_copies_prebuilt_binary", func(t *testing.T) {
		t.Chdir(t.TempDir())

		job := NewDockerfileGenerationJob(
			dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI), false, quietLogger())
		if err := job.Execute(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		verifyFileContents(t, "Dockerfile", []string{
			"FROM scratch",
			"ARG TARGETPLATFORM",
			"COPY $TARGETPLATFORM/docker-test /usr/bin/docker-test",
			`ENTRYPOINT ["/usr/bin/docker-test"]`,
		})

		content, _ := os.ReadFile("Dockerfile")
		if strings.Contains(string(content), "EXPOSE") {
			t.Error("CLI Dockerfile must not EXPOSE a port")
		}
	})

	t.Run("cgo_uses_alpine_runtime", func(t *testing.T) {
		t.Chdir(t.TempDir())

		cfg := dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI)
		cfg.CGOStatus = domain.CGOStatusEnabled

		job := NewDockerfileGenerationJob(cfg, false, quietLogger())
		if err := job.Execute(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		verifyFileContents(t, "Dockerfile", []string{"FROM alpine:latest"})
	})

	t.Run("webapi_exposes_8080", func(t *testing.T) {
		t.Chdir(t.TempDir())

		job := NewDockerfileGenerationJob(
			dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeWebAPI), false, quietLogger())
		if err := job.Execute(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		verifyFileContents(t, "Dockerfile", []string{"EXPOSE 8080"})
	})

	t.Run("existing_dockerfile_without_force_errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		if err := os.WriteFile("Dockerfile", []byte("FROM scratch"), 0o644); err != nil {
			t.Fatal(err)
		}

		job := NewDockerfileGenerationJob(
			dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI), false, quietLogger())

		err := job.Execute(context.Background())
		if err == nil || !strings.Contains(err.Error(), "--force") {
			t.Errorf("expected --force error, got: %v", err)
		}
	})

	t.Run("rollback_removes_generated_dockerfile", func(t *testing.T) {
		t.Chdir(t.TempDir())

		job := NewDockerfileGenerationJob(
			dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI), false, quietLogger())
		if err := job.Execute(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := job.Rollback(context.Background()); err != nil {
			t.Fatalf("rollback failed: %v", err)
		}

		if _, err := os.Stat("Dockerfile"); !os.IsNotExist(err) {
			t.Error("rollback must remove the generated Dockerfile")
		}
	})
}

func jobIDs(jobs []Job) []string {
	ids := make([]string, 0, len(jobs))
	for _, job := range jobs {
		ids = append(ids, job.ID())
	}

	return ids
}

func containsID(ids []string, id string) bool {
	return slices.Contains(ids, id)
}

func TestCreateConfigOnlyJobsIncludesDockerfileJob(t *testing.T) {
	factory := NewJobFactory(quietLogger())

	withDocker := factory.CreateConfigOnlyJobs(
		dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI), false)
	ids := jobIDs(withDocker)

	if !containsID(ids, "dockerfile-generation") {
		t.Errorf("config-only workflow with Docker must include the Dockerfile job, got %v", ids)
	}

	if !containsID(ids, "generation-preflight") {
		t.Errorf("config-only workflow must include the preflight job, got %v", ids)
	}

	withoutDocker := factory.CreateConfigOnlyJobs(
		dockerTestConfig(t, domain.DockerSupportNone, domain.ProjectTypeCLI), false)
	if containsID(jobIDs(withoutDocker), "dockerfile-generation") {
		t.Error("config-only workflow without Docker must not include the Dockerfile job")
	}
}

func TestCreateFullWizardJobsOrdering(t *testing.T) {
	factory := NewJobFactory(quietLogger())

	jobs := factory.CreateFullWizardJobs(
		dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI), false)
	ids := jobIDs(jobs)

	preflightIndex := -1

	for i, id := range ids {
		if id == "generation-preflight" {
			preflightIndex = i
		}
	}

	if preflightIndex == -1 {
		t.Fatalf("full wizard workflow must include the preflight job, got %v", ids)
	}

	for _, mustPrecede := range []string{"config-generation", "github-actions-generation", "dockerfile-generation"} {
		if !containsID(ids, mustPrecede) {
			continue // conditionally included; only check ordering when present
		}

		for i, id := range ids {
			if id == mustPrecede && i < preflightIndex {
				t.Errorf("%s must run after generation-preflight, got order %v", mustPrecede, ids)
			}
		}
	}
}

func TestConfigOnlyWorkflowGeneratesDockerfileEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	createBasicGoProject(t, tmpDir, "github.com/user/config-only-test")
	t.Chdir(tmpDir)

	cfg := dockerTestConfig(t, domain.DockerSupportBoth, domain.ProjectTypeCLI)

	workflow, err := NewWorkflowBuilder(quietLogger()).BuildWorkflow(WorkflowTypeConfigOnly, cfg, false)
	if err != nil {
		t.Fatalf("failed to build workflow: %v", err)
	}

	if err := workflow.Execute(context.Background()); err != nil {
		t.Fatalf("config-only workflow failed: %v", err)
	}

	AssertFileExists(t, ".goreleaser.yaml", "config-only workflow must generate .goreleaser.yaml")
	AssertFileExists(t, "Dockerfile", "config-only workflow with Docker must generate a Dockerfile")

	if _, err := os.Stat(filepath.Join(".github", "workflows", "release.yml")); !os.IsNotExist(err) {
		t.Error("config-only workflow must not generate the GitHub Actions workflow")
	}
}
