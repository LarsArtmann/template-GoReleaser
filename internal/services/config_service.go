package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"github.com/LarsArtmann/template-GoReleaser/internal/types"
)

// ConfigServiceImpl implements ConfigService interface
type ConfigServiceImpl struct {
	viper *viper.Viper
}

// NewConfigService creates a new configuration service
func NewConfigService(v *viper.Viper) ConfigService {
	return &ConfigServiceImpl{
		viper: v,
	}
}

// LoadConfig loads configuration from file or creates default
func (s *ConfigServiceImpl) LoadConfig() (*types.Config, error) {
	config := &types.Config{}
	
	// Try to load from viper first
	if err := s.viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Apply defaults if values are empty
	s.applyDefaults(config)
	
	return config, nil
}

// SaveConfig saves configuration to file
func (s *ConfigServiceImpl) SaveConfig(config *types.Config) error {
	// Determine config file path
	configPath := s.viper.ConfigFileUsed()
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".goreleaser-cli.yaml")
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ValidateConfig validates configuration structure and values
func (s *ConfigServiceImpl) ValidateConfig(config *types.Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	// Validate license type
	if config.License.Type == "" {
		return fmt.Errorf("license type cannot be empty")
	}
	
	validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "MPL-2.0"}
	isValid := false
	for _, valid := range validLicenses {
		if config.License.Type == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid license type: %s", config.License.Type)
	}
	
	// Validate author information for licenses that require it
	needsAuthor := []string{"MIT", "BSD-3-Clause", "ISC"}
	for _, license := range needsAuthor {
		if config.License.Type == license && config.Author.Name == "" {
			return fmt.Errorf("license %s requires author name", config.License.Type)
		}
	}
	
	return nil
}

// InitConfig creates a new configuration with defaults
func (s *ConfigServiceImpl) InitConfig() (*types.Config, error) {
	config := types.DefaultConfig()
	return config, nil
}

// applyDefaults applies default values to configuration
func (s *ConfigServiceImpl) applyDefaults(config *types.Config) {
	if config.License.Type == "" {
		config.License.Type = "MIT"
	}
	
	if config.CLI.Colors == false && config.CLI.Verbose == false {
		// Only set defaults if both are false (meaning uninitialized)
		config.CLI.Colors = true
		config.CLI.Verbose = false
	}
}