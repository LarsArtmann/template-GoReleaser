package types

import (
	"os"
	"testing"
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
