package domain

import (
	"fmt"
	"strings"
)

// GitProvider represents git hosting providers
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY.
type GitProvider string

const (
	GitProviderGitHub     GitProvider = "github"      // GitHub
	GitProviderGitLab     GitProvider = "gitlab"      // GitLab
	GitProviderBitbucket  GitProvider = "bitbucket"   // Bitbucket
	GitProviderGitea      GitProvider = "gitea"       // Gitea
	GitProviderSelfHosted GitProvider = "self-hosted" // Self-hosted
)

// GitProvider metadata - generated from TypeSpec invariants.
type gitProviderMeta struct {
	defaultRegistry             DockerRegistry
	actionsSupported            bool
	apiURL                      string
	webURL                      string
	requiresPersonalAccessToken bool
}

var gitProviderMetaMap = map[GitProvider]gitProviderMeta{
	GitProviderGitHub: {
		defaultRegistry:             DockerRegistryGitHub,
		actionsSupported:            true,
		apiURL:                      "https://api.github.com",
		webURL:                      "https://github.com",
		requiresPersonalAccessToken: false,
	},
	GitProviderGitLab: {
		defaultRegistry:             DockerRegistryGitLab,
		actionsSupported:            true,
		apiURL:                      "https://gitlab.com/api/v4",
		webURL:                      "https://gitlab.com",
		requiresPersonalAccessToken: true,
	},
	GitProviderBitbucket: {
		defaultRegistry:             DockerRegistryCustom,
		actionsSupported:            true,
		apiURL:                      "https://api.bitbucket.org/2.0",
		webURL:                      "https://bitbucket.org",
		requiresPersonalAccessToken: true,
	},
	GitProviderGitea: {
		defaultRegistry:             DockerRegistryCustom,
		actionsSupported:            false,
		apiURL:                      "", // Self-hosted
		webURL:                      "", // Self-hosted
		requiresPersonalAccessToken: true,
	},
	GitProviderSelfHosted: {
		defaultRegistry:             DockerRegistryCustom,
		actionsSupported:            false,
		apiURL:                      "", // User-defined
		webURL:                      "", // User-defined
		requiresPersonalAccessToken: true,
	},
}

// IsValid returns true if GitProvider is valid.
func (gp GitProvider) IsValid() bool {
	_, exists := gitProviderMetaMap[gp]

	return exists
}

// String returns human-readable display name.
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
	case GitProviderSelfHosted:
		return "Self-hosted"
	default:
		return string(gp)
	}
}

// DefaultRegistry returns the default Docker registry for this provider.
func (gp GitProvider) DefaultRegistry() DockerRegistry {
	if meta, exists := gitProviderMetaMap[gp]; exists {
		return meta.defaultRegistry
	}

	return DockerRegistryCustom
}

// ActionsSupported returns true if GitHub Actions are supported.
func (gp GitProvider) ActionsSupported() bool {
	if meta, exists := gitProviderMetaMap[gp]; exists {
		return meta.actionsSupported
	}

	return false
}

// APIURL returns the API URL for this provider.
func (gp GitProvider) APIURL() string {
	if meta, exists := gitProviderMetaMap[gp]; exists {
		return meta.apiURL
	}

	return ""
}

// WebURL returns the web URL for this provider.
func (gp GitProvider) WebURL() string {
	if meta, exists := gitProviderMetaMap[gp]; exists {
		return meta.webURL
	}

	return ""
}

// RequiresPersonalAccessToken returns true if provider requires personal access token.
func (gp GitProvider) RequiresPersonalAccessToken() bool {
	if meta, exists := gitProviderMetaMap[gp]; exists {
		return meta.requiresPersonalAccessToken
	}

	return true
}

// ValidateGitProvider validates a git provider.
func ValidateGitProvider(provider GitProvider) error {
	if !provider.IsValid() {
		return NewValidationError(
			ErrInvalidGitProvider,
			"Invalid git provider",
			fmt.Sprintf("'%s' is not a valid git provider", provider),
		)
	}

	return nil
}

// GetAllGitProviders returns all available git providers.
func GetAllGitProviders() []GitProvider {
	return []GitProvider{
		GitProviderGitHub, GitProviderGitLab, GitProviderBitbucket,
		GitProviderGitea, GitProviderSelfHosted,
	}
}

// GetRecommendedGitProvider returns recommended git provider (GitHub).
func GetRecommendedGitProvider() GitProvider {
	return GitProviderGitHub
}

// ConvertToGitProvider converts string display name to GitProvider.
func ConvertToGitProvider(displayName string) GitProvider {
	switch strings.ToLower(strings.TrimSpace(displayName)) {
	case "github", "🐙  github":
		return GitProviderGitHub
	case "gitlab", "🦊  gitlab":
		return GitProviderGitLab
	case "bitbucket", "🪣  bitbucket":
		return GitProviderBitbucket
	case "gitea", "🕊️  gitea":
		return GitProviderGitea
	case "self-hosted", "🏠  self-hosted":
		return GitProviderSelfHosted
	default:
		return GitProviderGitHub
	}
}
