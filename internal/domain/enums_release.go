package domain

// GitProvider represents supported Git hosting providers
// This enum replaces string-based Git provider types for type safety
type GitProvider string

const (
	GitProviderGitHub      GitProvider = "github"
	GitProviderGitLab      GitProvider = "gitlab"
	GitProviderBitbucket   GitProvider = "bitbucket"
	GitProviderGitea       GitProvider = "gitea"
	GitProviderGogs        GitProvider = "gogs"
	GitProviderAzureDevOps GitProvider = "azure-devops"
	GitProviderSelfHosted  GitProvider = "self-hosted"
)

// IsValid returns true if GitProvider is valid
func (gp GitProvider) IsValid() bool {
	switch gp {
	case GitProviderGitHub, GitProviderGitLab, GitProviderBitbucket,
		GitProviderGitea, GitProviderGogs, GitProviderAzureDevOps, GitProviderSelfHosted:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (gp GitProvider) String() string {
	switch gp {
	case GitProviderGitHub:
		return "GitHub"
	case GitProviderGitLab:
		return "GitLab"
	case GitProviderBitbucket:
		return "Bitbucket"
	case GitProviderGitea:
		return "Gitea"
	case GitProviderGogs:
		return "Gogs"
	case GitProviderAzureDevOps:
		return "Azure DevOps"
	case GitProviderSelfHosted:
		return "Self-Hosted"
	default:
		return "Unknown"
	}
}

// IsGitBased returns true for Git-based providers
func (gp GitProvider) IsGitBased() bool {
	return gp != GitProviderAzureDevOps
}

// IsOpenSource returns true for open-source platforms
func (gp GitProvider) IsOpenSource() bool {
	switch gp {
	case GitProviderGitHub, GitProviderGitLab, GitProviderGitea, GitProviderGogs:
		return true
	case GitProviderBitbucket, GitProviderAzureDevOps, GitProviderSelfHosted:
		return false
	default:
		return false
	}
}

// IsCommercial returns true for commercial platforms
func (gp GitProvider) IsCommercial() bool {
	switch gp {
	case GitProviderGitHub, GitProviderBitbucket, GitProviderAzureDevOps:
		return true
	case GitProviderGitLab, GitProviderGitea, GitProviderGogs, GitProviderSelfHosted:
		return false
	default:
		return false
	}
}

// SupportsSelfHosting returns true if provider supports self-hosting
func (gp GitProvider) SupportsSelfHosting() bool {
	switch gp {
	case GitProviderGitLab, GitProviderGitea, GitProviderGogs:
		return true
	case GitProviderGitHub, GitProviderBitbucket, GitProviderAzureDevOps:
		return false
	default:
		return false
	}
}

// ActionsSupported returns true if provider supports CI/CD actions
func (gp GitProvider) ActionsSupported() bool {
	switch gp {
	case GitProviderGitHub, GitProviderGitLab, GitProviderAzureDevOps:
		return true
	case GitProviderBitbucket, GitProviderGitea, GitProviderGogs:
		return false
	default:
		return false
	}
}

// DefaultRegistry returns the default container registry for provider
func (gp GitProvider) DefaultRegistry() DockerRegistry {
	switch gp {
	case GitProviderGitHub:
		return DockerRegistryGitHub
	case GitProviderGitLab:
		return DockerRegistryGitLab
	case GitProviderAzureDevOps:
		return DockerRegistryAzure
	case GitProviderBitbucket:
		return DockerRegistryBitbucket
	default:
		return DockerRegistryDockerHub
	}
}

// GetAPIURL returns the API URL for provider
func (gp GitProvider) GetAPIURL() string {
	switch gp {
	case GitProviderGitHub:
		return "https://api.github.com"
	case GitProviderGitLab:
		return "https://gitlab.com/api/v4"
	case GitProviderBitbucket:
		return "https://api.bitbucket.org/2.0"
	case GitProviderGitea:
		return "" // Configurable per instance
	case GitProviderGogs:
		return "" // Configurable per instance
	case GitProviderAzureDevOps:
		return "https://dev.azure.com"
	default:
		return ""
	}
}

// GetWebURL returns the web URL for provider
func (gp GitProvider) GetWebURL() string {
	switch gp {
	case GitProviderGitHub:
		return "https://github.com"
	case GitProviderGitLab:
		return "https://gitlab.com"
	case GitProviderBitbucket:
		return "https://bitbucket.org"
	case GitProviderGitea:
		return "" // Configurable per instance
	case GitProviderGogs:
		return "" // Configurable per instance
	case GitProviderAzureDevOps:
		return "https://dev.azure.com"
	default:
		return ""
	}
}

// DockerRegistry represents container registries
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

// SigningLevel represents code signing maturity levels
// This enum replaces boolean/level combinations for type safety
type SigningLevel string

const (
	SigningLevelNone       SigningLevel = "none"
	SigningLevelBasic      SigningLevel = "basic"
	SigningLevelStandard   SigningLevel = "standard"
	SigningLevelAdvanced   SigningLevel = "advanced"
	SigningLevelEnterprise SigningLevel = "enterprise"
)

// IsValid returns true if SigningLevel is valid
func (sl SigningLevel) IsValid() bool {
	switch sl {
	case SigningLevelNone, SigningLevelBasic, SigningLevelStandard,
		SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (sl SigningLevel) String() string {
	switch sl {
	case SigningLevelNone:
		return "None"
	case SigningLevelBasic:
		return "Basic"
	case SigningLevelStandard:
		return "Standard"
	case SigningLevelAdvanced:
		return "Advanced"
	case SigningLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if signing is enabled
func (sl SigningLevel) IsEnabled() bool {
	return sl != SigningLevelNone
}

// RequiresCosign returns true if level requires cosign
func (sl SigningLevel) RequiresCosign() bool {
	switch sl {
	case SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// RequiresKeyManagement returns true if level requires key management
func (sl SigningLevel) RequiresKeyManagement() bool {
	switch sl {
	case SigningLevelStandard, SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// RequiresCertificate returns true if level requires certificate management
func (sl SigningLevel) RequiresCertificate() bool {
	return sl == SigningLevelEnterprise
}

// GetRequiredTools returns the tools required for this signing level
func (sl SigningLevel) GetRequiredTools() []string {
	var tools []string

	if sl.RequiresCosign() {
		tools = append(tools, "cosign")
	}

	if sl.RequiresKeyManagement() {
		tools = append(tools, "gpg")
	}

	if sl.RequiresCertificate() {
		tools = append(tools, "openssl")
	}

	return tools
}

// GetSigningAlgorithms returns the signing algorithms supported
func (sl SigningLevel) GetSigningAlgorithms() []string {
	switch sl {
	case SigningLevelNone:
		return []string{}
	case SigningLevelBasic:
		return []string{"rsa-sha256"}
	case SigningLevelStandard:
		return []string{"rsa-sha256", "ecdsa-sha256"}
	case SigningLevelAdvanced:
		return []string{"rsa-sha256", "ecdsa-sha256", "ed25519"}
	case SigningLevelEnterprise:
		return []string{"rsa-sha256", "rsa-sha512", "ecdsa-sha256", "ecdsa-sha384", "ecdsa-sha512", "ed25519"}
	default:
		return []string{}
	}
}
