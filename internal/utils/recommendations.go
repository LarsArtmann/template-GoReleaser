package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/git"
)

// GetRecommendedProjectType analyzes current directory and returns recommended project type
func GetRecommendedProjectType() domain.ProjectType {
	// Check for common project indicators
	if hasFile("go.mod") {
		if hasFile("main.go") || hasFile("cmd/main.go") {
			return domain.ProjectTypeCLI
		}
		if hasFile("api/") || hasFile("server/") {
			return domain.ProjectTypeWebAPI
		}
		if hasFile("grpc/") || hasFile("pb/") {
			return domain.ProjectTypeGRPCService
		}
		return domain.ProjectTypeLibrary
	}

	// Default to CLI for Go projects
	return domain.ProjectTypeCLI
}

// GetRecommendedPlatforms returns recommended platforms based on current OS and common patterns
func GetRecommendedPlatforms() []domain.Platform {
	platforms := []domain.Platform{domain.PlatformLinux} // Always include Linux

	// Add current OS platform
	switch runtime.GOOS {
	case "darwin":
		platforms = append(platforms, domain.PlatformDarwin)
	case "windows":
		platforms = append(platforms, domain.PlatformWindows)
	}

	return platforms
}

// GetRecommendedArchitectures returns recommended architectures based on current arch
func GetRecommendedArchitectures() []domain.Architecture {
	architectures := []domain.Architecture{domain.ArchitectureAMD64} // Always include AMD64

	// Add current architecture if supported
	switch runtime.GOARCH {
	case "arm64":
		architectures = append(architectures, domain.ArchitectureARM64)
	case "386":
		architectures = append(architectures, domain.Architecture386)
	}

	return architectures
}

// GetRecommendedGitProvider returns recommended Git provider based on analysis
func GetRecommendedGitProvider() domain.GitProvider {
	// Try to detect from git remote
	if remote, err := exec.Command("git", "remote", "get-url", "origin").Output(); err == nil {
		remoteStr := strings.TrimSpace(string(remote))
		if strings.Contains(remoteStr, "github.com") {
			return domain.GitProviderGitHub
		}
		if strings.Contains(remoteStr, "gitlab.com") {
			return domain.GitProviderGitLab
		}
		if strings.Contains(remoteStr, "bitbucket.org") {
			return domain.GitProviderBitbucket
		}
	}

	// Default to GitHub
	return domain.GitProviderGitHub
}

// GetRecommendedDockerRegistry returns recommended Docker registry based on Git provider
func GetRecommendedDockerRegistry(gitProvider domain.GitProvider) domain.DockerRegistry {
	return gitProvider.DefaultRegistry()
}

// hasFile checks if a file exists in current directory or common subdirectories
func hasFile(filename string) bool {
	// Check current directory
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	// Check common subdirectories
	dirs := []string{"cmd", "pkg", "internal", "src", "app"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir + "/" + filename); err == nil {
			return true
		}
	}

	return false
}

// GetGitHubOwner tries to extract GitHub owner from git remote
func GetGitHubOwner() string {
	if owner := git.GetGitHubOwner(); owner != "owner" {
		return owner
	}
	return "owner"
}

// GetGitHubRepo tries to extract GitHub repo name from git remote
func GetGitHubRepo() string {
	if repo := git.GetGitHubRepo(); repo != "repo" {
		return repo
	}
	return "repo"
}

// IsDevelopmentEnvironment returns true if running in development environment
func IsDevelopmentEnvironment() bool {
	return os.Getenv("GO_ENV") == "development" ||
		os.Getenv("ENV") == "development" ||
		os.Getenv("NODE_ENV") == "development"
}

// IsProductionEnvironment returns true if running in production environment
func IsProductionEnvironment() bool {
	return os.Getenv("GO_ENV") == "production" ||
		os.Getenv("ENV") == "production" ||
		os.Getenv("NODE_ENV") == "production"
}

// GetEnvironment returns current environment name
func GetEnvironment() string {
	if env := os.Getenv("GO_ENV"); env != "" {
		return env
	}
	if env := os.Getenv("ENV"); env != "" {
		return env
	}
	if env := os.Getenv("NODE_ENV"); env != "" {
		return env
	}
	return "development"
}
