package services

import "github.com/LarsArtmann/template-GoReleaser/internal/types"

// ConfigService handles configuration management operations
type ConfigService interface {
	// LoadConfig loads configuration from file or creates default
	LoadConfig() (*types.Config, error)
	
	// SaveConfig saves configuration to file
	SaveConfig(config *types.Config) error
	
	// ValidateConfig validates configuration structure and values
	ValidateConfig(config *types.Config) error
	
	// InitConfig creates a new configuration with defaults
	InitConfig() (*types.Config, error)
}

// ValidationService handles all validation operations
type ValidationService interface {
	// ValidateProject validates the entire project structure
	ValidateProject() (*ValidationResult, error)
	
	// ValidateEnvironment validates environment variables
	ValidateEnvironment() (*ValidationResult, error)
	
	// ValidateGoReleaser validates GoReleaser configuration files
	ValidateGoReleaser(configPath string) (*ValidationResult, error)
	
	// ValidateTools validates required tools are installed
	ValidateTools() (*ValidationResult, error)
}

// LicenseService handles license management operations
type LicenseService interface {
	// GenerateLicense generates a license file based on configuration
	GenerateLicense(licenseType string, author string) error
	
	// ListAvailableLicenses returns available license templates
	ListAvailableLicenses() ([]LicenseTemplate, error)
	
	// ValidateLicense validates existing license file
	ValidateLicense() (*ValidationResult, error)
}

// VerificationService provides comprehensive project verification
type VerificationService interface {
	// RunFullVerification runs all verification checks
	RunFullVerification(opts *VerificationOptions) (*VerificationResult, error)
	
	// RunSecurityScan performs security validation
	RunSecurityScan() (*SecurityScanResult, error)
	
	// RunDryRun performs GoReleaser dry run
	RunDryRun(configPath string) (*DryRunResult, error)
}

// Import Config types from main types file to avoid duplication
// These types are defined in cmd/goreleaser-cli/types.go

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Checks   int      `json:"checks"`
}

// VerificationOptions configures verification behavior
type VerificationOptions struct {
	SkipSecurity    bool   `json:"skip_security"`
	SkipDryRun      bool   `json:"skip_dry_run"`
	SkipLicenseTest bool   `json:"skip_license_test"`
	ConfigFile      string `json:"config_file"`
	ProConfigFile   string `json:"pro_config_file"`
}

// VerificationResult contains comprehensive verification results
type VerificationResult struct {
	Success  bool                `json:"success"`
	Checks   int                 `json:"checks"`
	Warnings int                 `json:"warnings"`
	Errors   int                 `json:"errors"`
	Details  map[string][]string `json:"details"`
}

// SecurityScanResult contains security scan findings
type SecurityScanResult struct {
	Success     bool     `json:"success"`
	Issues      []string `json:"issues,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// DryRunResult contains GoReleaser dry run results  
type DryRunResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}

// LicenseTemplate represents an available license template
type LicenseTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}