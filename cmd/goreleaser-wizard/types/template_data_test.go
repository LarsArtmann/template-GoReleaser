package types

import (
	"os"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// resetGitHubOverrides clears flag overrides so tests are independent.
func resetGitHubOverrides(t *testing.T) {
	t.Helper()

	previousOwner, previousRepo := githubOwnerOverride, githubRepoOverride

	t.Cleanup(func() {
		githubOwnerOverride = previousOwner
		githubRepoOverride = previousRepo
	})

	githubOwnerOverride = ""
	githubRepoOverride = ""
	cachedGitHubRemote = nil
}

func TestGitHubOverridesTakePrecedence(t *testing.T) {
	resetGitHubOverrides(t)

	SetGitHubOwnerOverride("explicit-owner")
	SetGitHubRepoOverride("explicit-repo")

	if got := GetGitHubOwner(); got != "explicit-owner" {
		t.Errorf("GetGitHubOwner() = %q, want explicit override %q", got, "explicit-owner")
	}

	if got := GetGitHubRepo(); got != "explicit-repo" {
		t.Errorf("GetGitHubRepo() = %q, want explicit override %q", got, "explicit-repo")
	}

	if HasPlaceholderGitHubTarget() {
		t.Error("explicit overrides must not count as placeholders")
	}
}

func TestPlaceholderDetectionWithoutRemote(t *testing.T) {
	resetGitHubOverrides(t)

	// A directory outside any git repository: `git remote get-url origin`
	// fails and detection must fall back to the documented placeholders.
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	if got := GetGitHubOwner(); got != PlaceholderGitHubOwner {
		t.Errorf("GetGitHubOwner() = %q, want placeholder %q", got, PlaceholderGitHubOwner)
	}

	if got := GetGitHubRepo(); got != PlaceholderGitHubRepo {
		t.Errorf("GetGitHubRepo() = %q, want placeholder %q", got, PlaceholderGitHubRepo)
	}

	if !HasPlaceholderGitHubTarget() {
		t.Error("HasPlaceholderGitHubTarget() = false, want true when detection falls back")
	}
}

func TestPlaceholderDetectionWithRemote(t *testing.T) {
	resetGitHubOverrides(t)

	// The wizard's own repository always has an origin remote in CI and local
	// clones; if it somehow does not, detection falling back to placeholders
	// is acceptable, but it must never claim a remote exists when it doesn't.
	if _, err := os.Stat(".git"); err != nil {
		t.Skip("not running inside a git clone; remote-detection case not verifiable")
	}

	if HasPlaceholderGitHubTarget() {
		t.Error("HasPlaceholderGitHubTarget() = true, want false inside a clone with an origin remote")
	}
}

func TestDockerPlatforms(t *testing.T) {
	tests := []struct {
		name          string
		architectures []domain.Architecture
		wantPlatforms []string
	}{
		{
			name:          "amd64_only",
			architectures: []domain.Architecture{"amd64"},
			wantPlatforms: []string{"linux/amd64"},
		},
		{
			name:          "amd64_and_arm64",
			architectures: []domain.Architecture{"amd64", "arm64"},
			wantPlatforms: []string{"linux/amd64", "linux/arm64"},
		},
		{
			name:          "arm64_alone_still_gets_amd64",
			architectures: []domain.Architecture{"arm64"},
			wantPlatforms: []string{"linux/amd64", "linux/arm64"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.SafeProjectConfig{Architectures: tt.architectures}

			got := dockerPlatforms(config)
			if len(got) != len(tt.wantPlatforms) {
				t.Fatalf("dockerPlatforms() = %v, want %v", got, tt.wantPlatforms)
			}

			for i, want := range tt.wantPlatforms {
				if got[i] != want {
					t.Errorf("dockerPlatforms()[%d] = %q, want %q", i, got[i], want)
				}
			}
		})
	}
}

func TestNewGoReleaserTemplateData(t *testing.T) {
	resetGitHubOverrides(t)
	SetGitHubOwnerOverride("owner-a")
	SetGitHubRepoOverride("repo-b")

	config := &domain.SafeProjectConfig{
		ProjectName:   "typed-app",
		BinaryName:    "typed-bin",
		MainPath:      "./cmd/app",
		Platforms:     []domain.Platform{"linux", "darwin"},
		Architectures: []domain.Architecture{"amd64", "arm64"},
		CGOStatus:     domain.CGOStatusEnabled,
		DockerSupport: domain.DockerSupportBoth,
	}

	data := NewGoReleaserTemplateData(config)

	if data.ProjectName != "typed-app" || data.BinaryName != "typed-bin" || data.MainPath != "./cmd/app" {
		t.Errorf("project fields not mapped: %+v", data)
	}

	if data.CGOEnabled != "1" {
		t.Errorf("CGOEnabled = %q, want %q for enabled CGO", data.CGOEnabled, "1")
	}

	if data.GitHubOwner != "owner-a" || data.GitHubRepo != "repo-b" {
		t.Errorf("GitHub target = %q/%q, want overrides owner-a/repo-b", data.GitHubOwner, data.GitHubRepo)
	}

	if !data.DockerEnabled {
		t.Error("DockerEnabled must be true for DockerSupportBoth")
	}

	if len(data.Platforms) != 2 || data.Platforms[0] != "linux" {
		t.Errorf("Platforms = %v, want [linux darwin]", data.Platforms)
	}

	if len(data.Architectures) != 2 || data.Architectures[1] != "arm64" {
		t.Errorf("Architectures = %v, want [amd64 arm64]", data.Architectures)
	}

	if len(data.DockerPlatforms) != 2 || data.DockerPlatforms[1] != "linux/arm64" {
		t.Errorf("DockerPlatforms = %v, want [linux/amd64 linux/arm64]", data.DockerPlatforms)
	}

	if len(data.IgnoreCombinations) == 0 {
		t.Error("IgnoreCombinations must default to the common combinations")
	}
}
