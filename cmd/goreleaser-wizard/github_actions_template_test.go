package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	yaml "github.com/go-faster/yaml"
)

// renderReleaseWorkflow renders the GitHub Actions workflow for the given
// combo into a temp directory and returns its content.
func renderReleaseWorkflow(t *testing.T, cfg *ProjectConfig) string {
	t.Helper()

	t.Chdir(t.TempDir())

	if err := generateGitHubActions(cfg, quietLogger()); err != nil {
		t.Fatalf("generateGitHubActions() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("failed to read generated workflow: %v", err)
	}

	return string(content)
}

// assertWorkflowYAMLParses proves the rendered workflow is valid YAML.
func assertWorkflowYAMLParses(t *testing.T, content string) {
	t.Helper()

	var parsed any
	if err := yaml.Unmarshal([]byte(content), &parsed); err != nil {
		t.Fatalf("rendered workflow is not valid YAML: %v\n---\n%s", err, content)
	}
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(content, want) {
		t.Errorf("rendered workflow missing %q\n---\n%s", want, content)
	}
}

func assertNotContains(t *testing.T, content, unwanted string) {
	t.Helper()

	if strings.Contains(content, unwanted) {
		t.Errorf("rendered workflow must not contain %q\n---\n%s", unwanted, content)
	}
}

// TestGitHubActionsTemplateDockerSigningMatrix renders release.yml for all
// four docker×signing combinations; each must be valid YAML with the
// conditional steps present exactly when enabled (F31).
func TestGitHubActionsTemplateDockerSigningMatrix(t *testing.T) {
	combinations := []struct {
		name       string
		docker     domain.DockerSupport
		signing    domain.SigningLevel
		wantBuildx bool
		wantLogin  bool
		wantCosign bool
	}{
		{
			name:       "docker_and_signing",
			docker:     domain.DockerSupportBoth,
			signing:    domain.SigningLevelBasic,
			wantBuildx: true,
			wantLogin:  true,
			wantCosign: true,
		},
		{
			name:       "docker_only",
			docker:     domain.DockerSupportBoth,
			signing:    domain.SigningLevelNone,
			wantBuildx: true,
			wantLogin:  true,
		},
		{
			name:       "signing_only",
			docker:     domain.DockerSupportNone,
			signing:    domain.SigningLevelBasic,
			wantCosign: true,
		},
		{
			name:    "neither_docker_nor_signing",
			docker:  domain.DockerSupportNone,
			signing: domain.SigningLevelNone,
		},
	}

	for _, combo := range combinations {
		t.Run(combo.name, func(t *testing.T) {
			cfg := baseTestProjectConfig()
			cfg.ProjectName = "combo-app"
			cfg.BinaryName = "combo-app"
			cfg.MainPath = "."
			cfg.DockerSupport = combo.docker
			cfg.DockerRegistry = domain.DockerRegistryGitHub
			cfg.SigningLevel = combo.signing

			content := renderReleaseWorkflow(t, &cfg)

			assertWorkflowYAMLParses(t, content)
			assertContains(t, content, "runs-on: ubuntu-latest")
			assertContains(t, content, "goreleaser/goreleaser-action@v6")

			if combo.wantBuildx {
				assertContains(t, content, "docker/setup-buildx-action@v3")
			} else {
				assertNotContains(t, content, "docker/setup-buildx-action@v3")
			}

			if combo.wantLogin {
				assertContains(t, content, "docker/login-action@v3")
			} else {
				assertNotContains(t, content, "docker/login-action@v3")
			}

			if combo.wantCosign {
				assertContains(t, content, "sigstore/cosign-installer@v3")
			} else {
				assertNotContains(t, content, "sigstore/cosign-installer@v3")
			}
		})
	}
}
