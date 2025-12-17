package domain

// This enum replaces string-based registry types for type safety
type DockerRegistry string

const (
	DockerRegistryDockerHub  DockerRegistry = "dockerhub"
	DockerRegistryGitHub     DockerRegistry = "github"
	DockerRegistryGitLab     DockerRegistry = "gitlab"
	DockerRegistryAzure      DockerRegistry = "azure"
	DockerRegistryBitbucket  DockerRegistry = "bitbucket"
	DockerRegistryGoogle     DockerRegistry = "google"
	DockerRegistryAWS        DockerRegistry = "aws"
	DockerRegistryPrivate    DockerRegistry = "private"
	DockerRegistrySelfHosted DockerRegistry = "self-hosted"
)

// IsValid returns true if DockerRegistry is valid
func (dr DockerRegistry) IsValid() bool {
	switch dr {
	case DockerRegistryDockerHub, DockerRegistryGitHub, DockerRegistryGitLab,
		DockerRegistryAzure, DockerRegistryBitbucket, DockerRegistryGoogle,
		DockerRegistryAWS, DockerRegistryPrivate, DockerRegistrySelfHosted:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (dr DockerRegistry) String() string {
	switch dr {
	case DockerRegistryDockerHub:
		return "Docker Hub"
	case DockerRegistryGitHub:
		return "GitHub Container Registry"
	case DockerRegistryGitLab:
		return "GitLab Container Registry"
	case DockerRegistryAzure:
		return "Azure Container Registry"
	case DockerRegistryBitbucket:
		return "Bitbucket Container Registry"
	case DockerRegistryGoogle:
		return "Google Container Registry"
	case DockerRegistryAWS:
		return "Amazon ECR"
	case DockerRegistryPrivate:
		return "Private Registry"
	case DockerRegistrySelfHosted:
		return "Self-Hosted Registry"
	default:
		return "Unknown"
	}
}

// GetURL returns the registry URL
func (dr DockerRegistry) GetURL() string {
	switch dr {
	case DockerRegistryDockerHub:
		return "docker.io"
	case DockerRegistryGitHub:
		return "ghcr.io"
	case DockerRegistryGitLab:
		return "registry.gitlab.com"
	case DockerRegistryAzure:
		return ".azurecr.io"
	case DockerRegistryBitbucket:
		return "bitbucketpipelines.com"
	case DockerRegistryGoogle:
		return "gcr.io"
	case DockerRegistryAWS:
		return ".amazonaws.com"
	case DockerRegistryPrivate, DockerRegistrySelfHosted:
		return "" // Configurable
	default:
		return ""
	}
}

// RequiresAuth returns true if registry requires authentication
func (dr DockerRegistry) RequiresAuth() bool {
	switch dr {
	case DockerRegistryDockerHub, DockerRegistryGitHub, DockerRegistryGitLab,
		DockerRegistryAzure, DockerRegistryGoogle, DockerRegistryAWS:
		return true
	case DockerRegistryPrivate, DockerRegistrySelfHosted:
		return true
	default:
		return false
	}
}

// SupportsScanning returns true if registry supports vulnerability scanning
func (dr DockerRegistry) SupportsScanning() bool {
	switch dr {
	case DockerRegistryGitHub, DockerRegistryGitLab, DockerRegistryAzure,
		DockerRegistryGoogle, DockerRegistryAWS:
		return true
	case DockerRegistryDockerHub, DockerRegistryBitbucket:
		return false
	default:
		return false
	}
}

// IsCloudRegistry returns true for cloud-based registries
func (dr DockerRegistry) IsCloudRegistry() bool {
	switch dr {
	case DockerRegistryGitHub, DockerRegistryGitLab, DockerRegistryAzure,
		DockerRegistryGoogle, DockerRegistryAWS:
		return true
	case DockerRegistryDockerHub, DockerRegistryPrivate, DockerRegistrySelfHosted:
		return false
	default:
		return false
	}
}