package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/do/v2"
	"gopkg.in/yaml.v3"
)

type ConfigService struct{}

type GoReleaserConfig struct {
	ProjectName string                   `yaml:"project_name"`
	Before      *BeforeHooks            `yaml:"before,omitempty"`
	Builds      []BuildConfig           `yaml:"builds,omitempty"`
	Archives    []ArchiveConfig         `yaml:"archives,omitempty"`
	Checksum    *ChecksumConfig         `yaml:"checksum,omitempty"`
	Snapshot    *SnapshotConfig         `yaml:"snapshot,omitempty"`
	Changelog   *ChangelogConfig        `yaml:"changelog,omitempty"`
	Dockers     []DockerConfig          `yaml:"dockers,omitempty"`
	Brews       []BrewConfig           `yaml:"brews,omitempty"`
}

type BeforeHooks struct {
	Hooks []string `yaml:"hooks,omitempty"`
}

type BuildConfig struct {
	Env    []string `yaml:"env,omitempty"`
	Goos   []string `yaml:"goos,omitempty"`
	Goarch []string `yaml:"goarch,omitempty"`
}

type ArchiveConfig struct {
	Format          string            `yaml:"format,omitempty"`
	NameTemplate    string            `yaml:"name_template,omitempty"`
	FormatOverrides []FormatOverride  `yaml:"format_overrides,omitempty"`
}

type FormatOverride struct {
	Goos   string `yaml:"goos"`
	Format string `yaml:"format"`
}

type ChecksumConfig struct {
	NameTemplate string `yaml:"name_template,omitempty"`
}

type SnapshotConfig struct {
	NameTemplate string `yaml:"name_template,omitempty"`
}

type ChangelogConfig struct {
	Sort    string           `yaml:"sort,omitempty"`
	Filters *ChangelogFilter `yaml:"filters,omitempty"`
}

type ChangelogFilter struct {
	Exclude []string `yaml:"exclude,omitempty"`
}

type DockerConfig struct {
	ImageTemplates    []string `yaml:"image_templates,omitempty"`
	Dockerfile        string   `yaml:"dockerfile,omitempty"`
	BuildFlagTemplates []string `yaml:"build_flag_templates,omitempty"`
}

type BrewConfig struct {
	Name        string    `yaml:"name,omitempty"`
	Tap         *TapConfig `yaml:"tap,omitempty"`
	Homepage    string    `yaml:"homepage,omitempty"`
	Description string    `yaml:"description,omitempty"`
	License     string    `yaml:"license,omitempty"`
}

type TapConfig struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

func NewConfigService(injector do.Injector) (*ConfigService, error) {
	return &ConfigService{}, nil
}

func (s *ConfigService) LoadConfig(path string) (*GoReleaserConfig, error) {
	if path == "" {
		path = s.findConfigFile()
	}
	
	if path == "" {
		return nil, fmt.Errorf("no GoReleaser config file found")
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config GoReleaserConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

func (s *ConfigService) SaveConfig(path string, config *GoReleaserConfig) error {
	if path == "" {
		path = ".goreleaser.yaml"
	}
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (s *ConfigService) SaveConfigString(path string, content string) error {
	if path == "" {
		path = ".goreleaser.yaml"
	}
	
	// Validate YAML syntax
	var temp interface{}
	if err := yaml.Unmarshal([]byte(content), &temp); err != nil {
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}
	
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (s *ConfigService) findConfigFile() string {
	candidates := []string{
		".goreleaser.yaml",
		".goreleaser.yml",
		".goreleaser.pro.yaml",
		".goreleaser.pro.yml",
		"goreleaser.yaml",
		"goreleaser.yml",
	}
	
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			abs, _ := filepath.Abs(candidate)
			return abs
		}
	}
	
	return ""
}

func (s *ConfigService) ValidateSyntax(content string) error {
	var temp interface{}
	return yaml.Unmarshal([]byte(content), &temp)
}