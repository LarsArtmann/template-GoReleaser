package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

const defaultGitCommandTimeout = 30 * time.Second

// Command represents a git command with context.
type Command struct {
	ctx     context.Context
	dir     string
	env     []string
	timeout time.Duration
}

// NewCommand creates a new git command.
func NewCommand(ctx context.Context) *Command {
	return &Command{
		ctx:     ctx,
		timeout: defaultGitCommandTimeout,
	}
}

// WithDir sets the working directory for git commands.
func (c *Command) WithDir(dir string) *Command {
	c.dir = dir

	return c
}

// WithEnv sets environment variables for git commands.
func (c *Command) WithEnv(env []string) *Command {
	c.env = append(c.env, env...)

	return c
}

// WithTimeout sets the timeout for git commands.
func (c *Command) WithTimeout(timeout time.Duration) *Command {
	c.timeout = timeout

	return c
}

// execute executes a git command and returns the output.
func (c *Command) execute(args ...string) (string, error) {
	cmd := exec.CommandContext(c.ctx, "git", args...)
	if c.dir != "" {
		cmd.Dir = c.dir
	}

	if len(c.env) > 0 {
		cmd.Env = append(os.Environ(), c.env...)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewGitError(
			errors.ErrGitCommand,
			"Git command failed: git "+strings.Join(args, " "),
			err.Error(),
		).WithCause(err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetVersion returns the git version.
func (c *Command) GetVersion() (string, error) {
	return c.execute("version")
}

// GetDescribe gets the git describe output.
func (c *Command) GetDescribe() (string, error) {
	return c.execute("describe", "--tags", "--always", "--dirty")
}

// GetCommitHash gets the current commit hash.
func (c *Command) GetCommitHash() (string, error) {
	return c.execute("rev-parse", "HEAD")
}

// GetRemoteURL gets the URL for a remote.
func (c *Command) GetRemoteURL(remote string) (string, error) {
	return c.execute("remote", "get-url", remote)
}

// GetBranches gets all branches.
func (c *Command) GetBranches() ([]string, error) {
	output, err := c.execute("branch", "-a")
	if err != nil {
		return nil, err
	}

	branches := strings.Split(output, "\n")

	var result []string

	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" && !strings.HasPrefix(branch, "remotes/origin/HEAD") {
			result = append(result, branch)
		}
	}

	return result, nil
}

// GetTags gets all tags.
func (c *Command) GetTags() ([]string, error) {
	output, err := c.execute("tag", "-l")
	if err != nil {
		return nil, err
	}

	tags := strings.Split(output, "\n")

	var result []string

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result, nil
}

// IsRepository checks if the current directory is a git repository.
func (c *Command) IsRepository() bool {
	_, err := c.execute("rev-parse", "--git-dir")

	return err == nil
}

// HasRemote checks if the repository has a remote.
func (c *Command) HasRemote(remote string) bool {
	_, err := c.execute("remote", "get-url", remote)

	return err == nil
}

// GetRepositoryInfo gets information about the repository.
func (c *Command) GetRepositoryInfo() (*RepositoryInfo, error) {
	if !c.IsRepository() {
		return nil, errors.NewGitError(
			errors.ErrGitNotFound,
			"Not a git repository",
			"Current directory is not a git repository",
		)
	}

	info := &RepositoryInfo{}

	// Get commit hash
	if commitHash, err := c.GetCommitHash(); err == nil {
		info.CommitHash = commitHash
	}

	// Get description
	if describe, err := c.GetDescribe(); err == nil {
		info.Description = describe
	}

	// Get remote info
	if c.HasRemote("origin") {
		if remoteURL, err := c.GetRemoteURL("origin"); err == nil {
			info.RemoteURL = remoteURL
			info.Owner, info.Repo = ParseGitHubURL(remoteURL)
		}
	}

	// Get branches
	if branches, err := c.GetBranches(); err == nil {
		info.Branches = branches
	}

	// Get tags
	if tags, err := c.GetTags(); err == nil {
		info.Tags = tags
	}

	return info, nil
}

// RepositoryInfo contains information about a git repository.
type RepositoryInfo struct {
	CommitHash  string   `json:"commit_hash"`
	Description string   `json:"description"`
	RemoteURL   string   `json:"remote_url,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Repo        string   `json:"repo,omitempty"`
	Branches    []string `json:"branches,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ParseGitHubURL parses a GitHub URL and extracts owner and repo.
// Supports formats: git@github.com:owner/repo.git or https://github.com/owner/repo.git
func ParseGitHubURL(url string) (owner, repo string) {
	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "github.com")
		if len(parts) > 1 {
			repoPath := strings.TrimPrefix(parts[1], ":")
			repoPath = strings.TrimPrefix(repoPath, "/")
			repoPath = strings.TrimSuffix(repoPath, ".git")

			pathParts := strings.Split(repoPath, "/")
			if len(pathParts) > 1 {
				return pathParts[0], pathParts[1]
			}
		}
	}

	return "owner", "repo"
}

// GetVersionInfo gets version information from git.
func GetVersionInfo(ctx context.Context) (*VersionInfo, error) {
	cmd := NewCommand(ctx)

	info := &VersionInfo{
		Timestamp: time.Now(),
	}

	// Get version from git describe
	if version, err := cmd.GetDescribe(); err == nil {
		info.Version = version
	}

	// Get commit hash
	if commitHash, err := cmd.GetCommitHash(); err == nil {
		info.CommitHash = commitHash
	}

	// Get repository info
	if repoInfo, err := cmd.GetRepositoryInfo(); err == nil {
		info.Owner = repoInfo.Owner
		info.Repo = repoInfo.Repo
	}

	return info, nil
}

// VersionInfo contains version information from git.
type VersionInfo struct {
	Version    string    `json:"version"`
	CommitHash string    `json:"commit_hash"`
	Owner      string    `json:"owner,omitempty"`
	Repo       string    `json:"repo,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// GetMajorVersion extracts major version from version string.
func GetMajorVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		return parts[0]
	}

	return "0"
}

// GetCurrentDate returns current date in YYYY-MM-DD format.
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetGitHubOwner tries to get GitHub owner from git remote.
func GetGitHubOwner() string {
	ctx := context.Background()
	cmd := NewCommand(ctx)

	info, err := cmd.GetRepositoryInfo()
	if err != nil {
		return "owner"
	}

	return info.Owner
}

// GetGitHubRepo tries to get GitHub repo from git remote.
func GetGitHubRepo() string {
	ctx := context.Background()
	cmd := NewCommand(ctx)

	info, err := cmd.GetRepositoryInfo()
	if err != nil {
		return "repo"
	}

	return info.Repo
}

// IncPatchVersion increments the patch version for snapshots.
// It removes the 'v' prefix if present, splits the version string,
// and returns the version with the patch number incremented and "-next" suffix.
func IncPatchVersion(v string) string {
	if strings.HasPrefix(v, "v") {
		v = v[1:]
	}

	parts := strings.Split(v, ".")
	if len(parts) == 3 {
		patch := 0
		if len(parts) > 2 {
			if p, err := fmt.Sscanf(parts[2], "%d", &patch); err == nil && p == 1 {
				return fmt.Sprintf("v%s.%s.%d-next", parts[0], parts[1], patch+1)
			}
		}
	}

	return v + "-next"
}

// GetGitHubRepoDescription fetches the repository description from GitHub using gh CLI.
// Returns an empty string if gh is not available or the description cannot be fetched.
func GetGitHubRepoDescription() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "repo", "view", "--json", "description", "-q", ".description")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
