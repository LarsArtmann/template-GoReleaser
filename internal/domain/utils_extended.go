package domain

// GetURL returns the registry URL for this DockerRegistry
func (dr DockerRegistry) GetURL() string {
	switch dr {
	case DockerRegistryDockerHub:
		return "docker.io"
	case DockerRegistryGitHub:
		return "ghcr.io"
	case DockerRegistryGitLab:
		return "registry.gitlab.com"
	case DockerRegistryQuay:
		return "quay.io"
	case DockerRegistryCustom:
		return "" // User-defined
	default:
		return ""
	}
}

// SupportsDocker returns true if git provider supports Docker
func (gp GitProvider) SupportsDocker() bool {
	switch gp {
	case GitProviderGitHub, GitProviderGitLab, GitProviderGitea, GitProviderSelfHosted:
		return true
	case GitProviderBitbucket:
		return false
	default:
		return false
	}
}

// GetGitHubOwner returns GitHub owner from project config
func GetGitHubOwner(projectName string) string {
	// Simple extraction - in real implementation would parse from git remote
	return projectName
}
