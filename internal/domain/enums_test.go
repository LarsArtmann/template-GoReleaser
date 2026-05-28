package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DockerRegistry tests.

func TestDockerRegistry_Validation(t *testing.T) {
	tests := []struct {
		name     string
		registry DockerRegistry
		valid    bool
	}{
		{"DockerHub", DockerRegistryDockerHub, true},
		{"GitHub", DockerRegistryGitHub, true},
		{"GitLab", DockerRegistryGitLab, true},
		{"Quay", DockerRegistryQuay, true},
		{"Custom", DockerRegistryCustom, true},
		{"Invalid", DockerRegistry("invalid"), false},
		{"Empty", DockerRegistry(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.registry.IsValid())
		})
	}
}

func TestDockerRegistry_String(t *testing.T) {
	tests := []struct {
		registry DockerRegistry
		expected string
	}{
		{DockerRegistryDockerHub, "Docker Hub"},
		{DockerRegistryGitHub, "GitHub Container Registry"},
		{DockerRegistryGitLab, "GitLab Registry"},
		{DockerRegistryQuay, "Quay.io"},
		{DockerRegistryCustom, "Custom Registry"},
		{DockerRegistry("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.registry.String())
		})
	}
}

func TestDockerRegistry_Metadata(t *testing.T) {
	t.Run("DockerHub metadata", func(t *testing.T) {
		assert.Equal(t, "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$", DockerRegistryDockerHub.URLPattern())
		assert.False(t, DockerRegistryDockerHub.SupportsHTTPSOnly())
		assert.False(t, DockerRegistryDockerHub.RequiresAuthentication())
		assert.Equal(t, "library", DockerRegistryDockerHub.DefaultNamespace())
	})

	t.Run("GitHub metadata", func(t *testing.T) {
		assert.Equal(t, "^ghcr\\.io/[a-z0-9-]+/[a-z0-9-]+$", DockerRegistryGitHub.URLPattern())
		assert.True(t, DockerRegistryGitHub.SupportsHTTPSOnly())
		assert.True(t, DockerRegistryGitHub.RequiresAuthentication())
		assert.Empty(t, DockerRegistryGitHub.DefaultNamespace())
	})

	t.Run("invalid registry returns safe defaults", func(t *testing.T) {
		invalid := DockerRegistry("invalid")
		assert.Empty(t, invalid.URLPattern())
		assert.True(t, invalid.SupportsHTTPSOnly())
		assert.True(t, invalid.RequiresAuthentication())
		assert.Empty(t, invalid.DefaultNamespace())
	})
}

func TestValidateDockerRegistry(t *testing.T) {
	t.Run("valid registry", func(t *testing.T) {
		assert.NoError(t, ValidateDockerRegistry(DockerRegistryDockerHub))
	})

	t.Run("invalid registry", func(t *testing.T) {
		err := ValidateDockerRegistry(DockerRegistry("bad"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid Docker registry")
	})
}

func TestValidateDockerRegistryURL(t *testing.T) {
	tests := []struct {
		name     string
		registry DockerRegistry
		url      string
		wantErr  bool
	}{
		{"GitHub valid", DockerRegistryGitHub, "ghcr.io/user/repo", false},
		{"GitHub missing domain", DockerRegistryGitHub, "user/repo", true},
		{"GitLab valid", DockerRegistryGitLab, "registry.gitlab.com/group/project", false},
		{"Quay valid", DockerRegistryQuay, "quay.io/org/image", false},
		{"Custom any URL", DockerRegistryCustom, "anything", false},
		{"Empty URL", DockerRegistryGitHub, "", true},
		{"Invalid registry", DockerRegistry("bad"), "url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDockerRegistryURL(tt.registry, tt.url)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ConfigState tests.

func TestConfigState_Validation(t *testing.T) {
	tests := []struct {
		name  string
		state ConfigState
		valid bool
	}{
		{"Draft", ConfigStateDraft, true},
		{"Valid", ConfigStateValid, true},
		{"Invalid", ConfigStateInvalid, true},
		{"Processing", ConfigStateProcessing, true},
		{"Generated", ConfigStateGenerated, true},
		{"Unknown", ConfigState("unknown"), false},
		{"Empty", ConfigState(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.state.IsValid())
		})
	}
}

func TestConfigState_String(t *testing.T) {
	tests := []struct {
		state    ConfigState
		expected string
	}{
		{ConfigStateDraft, "Draft"},
		{ConfigStateValid, "Valid"},
		{ConfigStateInvalid, "Invalid"},
		{ConfigStateProcessing, "Processing"},
		{ConfigStateGenerated, "Generated"},
		{ConfigState("other"), "other"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}

func TestConfigState_Metadata(t *testing.T) {
	t.Run("Draft allows validation but not generation", func(t *testing.T) {
		assert.Equal(t, "Configuration is being created or modified", ConfigStateDraft.Description())
		assert.False(t, ConfigStateDraft.IsFinal())
		assert.True(t, ConfigStateDraft.AllowsValidation())
		assert.False(t, ConfigStateDraft.AllowsGeneration())
	})

	t.Run("Generated is final", func(t *testing.T) {
		assert.Equal(t, "Configuration has been generated successfully", ConfigStateGenerated.Description())
		assert.True(t, ConfigStateGenerated.IsFinal())
		assert.True(t, ConfigStateGenerated.AllowsValidation())
		assert.False(t, ConfigStateGenerated.AllowsGeneration())
	})

	t.Run("Valid allows both validation and generation", func(t *testing.T) {
		assert.True(t, ConfigStateValid.AllowsValidation())
		assert.True(t, ConfigStateValid.AllowsGeneration())
	})

	t.Run("Processing allows neither", func(t *testing.T) {
		assert.False(t, ConfigStateProcessing.AllowsValidation())
		assert.False(t, ConfigStateProcessing.AllowsGeneration())
	})

	t.Run("invalid state returns safe defaults", func(t *testing.T) {
		invalid := ConfigState("bad")
		assert.Empty(t, invalid.Description())
		assert.False(t, invalid.IsFinal())
		assert.False(t, invalid.AllowsValidation())
		assert.False(t, invalid.AllowsGeneration())
	})
}

func TestValidateConfigState(t *testing.T) {
	t.Run("valid state", func(t *testing.T) {
		assert.NoError(t, ValidateConfigState(ConfigStateDraft))
	})

	t.Run("invalid state", func(t *testing.T) {
		err := ValidateConfigState(ConfigState("bad"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid configuration state")
	})
}

func TestGetInitialConfigState(t *testing.T) {
	assert.Equal(t, ConfigStateDraft, GetInitialConfigState())
}

func TestGetAllConfigStates(t *testing.T) {
	states := GetAllConfigStates()
	assert.Len(t, states, 5)
	assert.Contains(t, states, ConfigStateDraft)
	assert.Contains(t, states, ConfigStateValid)
	assert.Contains(t, states, ConfigStateInvalid)
	assert.Contains(t, states, ConfigStateProcessing)
	assert.Contains(t, states, ConfigStateGenerated)
}
