package domain

import "fmt"

// BuildLevel represents build complexity levels
// This enum replaces bool flags for build configuration
type BuildLevel string

const (
	// BuildLevelMinimal builds with minimal optimizations and features
	BuildLevelMinimal BuildLevel = "minimal"
	// BuildLevelBasic includes standard optimizations and basic features
	BuildLevelBasic BuildLevel = "basic"
	// BuildLevelStandard includes full optimizations and standard features
	BuildLevelStandard BuildLevel = "standard"
	// BuildLevelAdvanced includes aggressive optimizations and advanced features
	BuildLevelAdvanced BuildLevel = "advanced"
	// BuildLevelEnterprise includes enterprise-grade optimizations and compliance features
	BuildLevelEnterprise BuildLevel = "enterprise"
)

// IsValid returns true if BuildLevel is valid
func (bl BuildLevel) IsValid() bool {
	switch bl {
	case BuildLevelMinimal, BuildLevelBasic, BuildLevelStandard,
		BuildLevelAdvanced, BuildLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (bl BuildLevel) String() string {
	switch bl {
	case BuildLevelMinimal:
		return "Minimal"
	case BuildLevelBasic:
		return "Basic"
	case BuildLevelStandard:
		return "Standard"
	case BuildLevelAdvanced:
		return "Advanced"
	case BuildLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsOptimized returns true if level includes optimizations
func (bl BuildLevel) IsOptimized() bool {
	return bl == BuildLevelStandard || bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}

// IsProductionReady returns true if level is production-ready
func (bl BuildLevel) IsProductionReady() bool {
	return bl == BuildLevelStandard || bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}

// IsAdvanced returns true if level includes advanced features
func (bl BuildLevel) IsAdvanced() bool {
	return bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}

// IsEnterprise returns true if level includes enterprise features
func (bl BuildLevel) IsEnterprise() bool {
	return bl == BuildLevelEnterprise
}

// ToBool converts to legacy boolean for compatibility
func (bl BuildLevel) ToBool() bool {
	return bl.IsOptimized()
}

// ValidateBuildLevel validates a build level
func ValidateBuildLevel(level BuildLevel) error {
	if !level.IsValid() {
		return fmt.Errorf("invalid build level: %s", level)
	}
	return nil
}

// GetOptimizationFlags returns optimization flags for level
func (bl BuildLevel) GetOptimizationFlags() []string {
	switch bl {
	case BuildLevelMinimal:
		return []string{"-gcflags=all=-l"}
	case BuildLevelBasic:
		return []string{"-gcflags=all=-l"}
	case BuildLevelStandard:
		return []string{"-gcflags=all=-l", "-ldflags=-s -w"}
	case BuildLevelAdvanced:
		return []string{"-gcflags=all=-l", "-ldflags=-s -w"}
	case BuildLevelEnterprise:
		return []string{"-gcflags=all=-l", "-ldflags=-s -w"}
	default:
		return []string{}
	}
}

// GetBuildTags returns recommended build tags for level
func (bl BuildLevel) GetBuildTags() []BuildTag {
	switch bl {
	case BuildLevelMinimal:
		return []BuildTag{CreateBuildTag("pure", "Pure Go compilation")}
	case BuildLevelBasic:
		return []BuildTag{CreateBuildTag("pure", "Pure Go compilation")}
	case BuildLevelStandard:
		return []BuildTag{CreateBuildTag("pure", "Pure Go compilation"), CreateBuildTag("static", "Static linking")}
	case BuildLevelAdvanced:
		return []BuildTag{CreateBuildTag("pure", "Pure Go compilation"), CreateBuildTag("static", "Static linking")}
	case BuildLevelEnterprise:
		return []BuildTag{CreateBuildTag("pure", "Pure Go compilation"), CreateBuildTag("static", "Static linking")}
	default:
		return []BuildTag{}
	}
}

// RequiresCGO returns true if level requires CGO
func (bl BuildLevel) RequiresCGO() bool {
	return false // All build levels support pure Go by default
}

// RequiresCache returns true if level requires build cache
func (bl BuildLevel) RequiresCache() bool {
	return bl == BuildLevelStandard || bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}

// RequiresParallel returns true if level requires parallel builds
func (bl BuildLevel) RequiresParallel() bool {
	return bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}
