package services

// ConfigService handles configuration management operations
type ConfigService interface {
	// LoadConfig loads configuration from file or creates default
	LoadConfig() (*Config, error)
	
	// SaveConfig saves configuration to file
	SaveConfig(config *Config) error
	
	// ValidateConfig validates configuration structure and values
	ValidateConfig(config *Config) error
	
	// InitConfig creates a new configuration with defaults
	InitConfig() (*Config, error)
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

// Config represents the application configuration
type Config struct {
	License LicenseConfig `yaml:"license" json:"license"`
	Author  AuthorConfig  `yaml:"author" json:"author"`
	Project ProjectConfig `yaml:"project" json:"project"`
	CLI     CLIConfig     `yaml:"cli" json:"cli"`
}

// LicenseConfig holds license-related configuration
type LicenseConfig struct {
	Type string `yaml:"type" json:"type"`
}

// AuthorConfig holds author information
type AuthorConfig struct {
	Name  string `yaml:"name" json:"name"`
	Email string `yaml:"email" json:"email"`
}

// ProjectConfig holds project-specific settings
type ProjectConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

// CLIConfig holds CLI behavior settings
type CLIConfig struct {
	Verbose bool `yaml:"verbose" json:"verbose"`
	Colors  bool `yaml:"colors" json:"colors"`
}

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