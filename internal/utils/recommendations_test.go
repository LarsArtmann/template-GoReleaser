package utils

import (
	"os"
	"runtime"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetRecommendedPlatforms(t *testing.T) {
	platforms := GetRecommendedPlatforms()
	assert.Contains(t, platforms, domain.PlatformLinux, "should always include Linux")

	switch runtime.GOOS {
	case "darwin":
		assert.Contains(t, platforms, domain.PlatformDarwin)
	case "windows":
		assert.Contains(t, platforms, domain.PlatformWindows)
	}
}

func TestGetRecommendedArchitectures(t *testing.T) {
	archs := GetRecommendedArchitectures()
	assert.Contains(t, archs, domain.ArchitectureAMD64, "should always include AMD64")

	switch runtime.GOARCH {
	case "arm64":
		assert.Contains(t, archs, domain.ArchitectureARM64)
	case "386":
		assert.Contains(t, archs, domain.Architecture386)
	}
}

func TestGetRecommendedDockerRegistry(t *testing.T) {
	assert.Equal(t, domain.DockerRegistryGitHub, GetRecommendedDockerRegistry(domain.GitProviderGitHub))
	assert.Equal(t, domain.DockerRegistryGitLab, GetRecommendedDockerRegistry(domain.GitProviderGitLab))
	assert.Equal(t, domain.DockerRegistryCustom, GetRecommendedDockerRegistry(domain.GitProviderGitea))
}

func TestIsDevelopmentEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "development")
	assert.True(t, IsDevelopmentEnvironment())

	os.Setenv("GO_ENV", "production")
	assert.False(t, IsDevelopmentEnvironment())
}

func TestIsProductionEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "production")
	assert.True(t, IsProductionEnvironment())

	os.Setenv("GO_ENV", "development")
	assert.False(t, IsProductionEnvironment())
}

func TestGetEnvironment(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("ENV")
		os.Unsetenv("NODE_ENV")
	})

	os.Setenv("GO_ENV", "staging")
	assert.Equal(t, "staging", GetEnvironment())

	os.Unsetenv("GO_ENV")
	os.Setenv("ENV", "test")
	assert.Equal(t, "test", GetEnvironment())

	os.Unsetenv("ENV")
	os.Setenv("NODE_ENV", "ci")
	assert.Equal(t, "ci", GetEnvironment())

	os.Unsetenv("NODE_ENV")
	assert.Equal(t, "development", GetEnvironment())
}
