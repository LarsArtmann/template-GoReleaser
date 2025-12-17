package domain

// BuildTag represents a build tag for conditional compilation
// This type-safe enum replaces string-based build tags
type BuildTag string

const (
	BuildTagPureGo      BuildTag = "purego"
	BuildTagCGoFree     BuildTag = "cgo_free"
	BuildTagNoASM       BuildTag = "noasm"
	BuildTagNetgo       BuildTag = "netgo"
	BuildTagOsusergo    BuildTag = "osusergo"
	BuildTagSQLiteFts5  BuildTag = "sqlite_fts5"
	BuildTagSQLiteJSON1 BuildTag = "sqlite_json1"
	BuildTagLMDB        BuildTag = "lmdb"
	BuildTagBtrfs       BuildTag = "btrfs"
	BuildTagNoCgo       BuildTag = "nocgo"
	BuildTagNoPie       BuildTag = "nopie"
	BuildTagStatic      BuildTag = "static"
	BuildTagStaticAll   BuildTag = "static_all"
	BuildTagRace        BuildTag = "race"
	BuildTagCoverage    BuildTag = "coverage"
	BuildTagDebug       BuildTag = "debug"
	BuildTagRelease     BuildTag = "release"
	BuildTagProduction  BuildTag = "production"
)

// IsValid returns true if BuildTag is valid
func (bt BuildTag) IsValid() bool {
	switch bt {
	case BuildTagPureGo, BuildTagCGoFree, BuildTagNoASM, BuildTagNetgo, BuildTagOsusergo,
		BuildTagSQLiteFts5, BuildTagSQLiteJSON1, BuildTagLMDB, BuildTagBtrfs,
		BuildTagNoCgo, BuildTagNoPie, BuildTagStatic, BuildTagStaticAll,
		BuildTagRace, BuildTagCoverage, BuildTagDebug, BuildTagRelease, BuildTagProduction:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (bt BuildTag) String() string {
	switch bt {
	case BuildTagPureGo:
		return "Pure Go (No CGO)"
	case BuildTagCGoFree:
		return "CGO-Free"
	case BuildTagNoASM:
		return "No Assembly"
	case BuildTagNetgo:
		return "Go Net Implementation"
	case BuildTagOsusergo:
		return "Go OS/User Implementation"
	case BuildTagSQLiteFts5:
		return "SQLite FTS5 Support"
	case BuildTagSQLiteJSON1:
		return "SQLite JSON1 Support"
	case BuildTagLMDB:
		return "LMDB Support"
	case BuildTagBtrfs:
		return "Btrfs Support"
	case BuildTagNoCgo:
		return "No CGO"
	case BuildTagNoPie:
		return "No PIE"
	case BuildTagStatic:
		return "Static Linking"
	case BuildTagStaticAll:
		return "Static All"
	case BuildTagRace:
		return "Race Detection"
	case BuildTagCoverage:
		return "Code Coverage"
	case BuildTagDebug:
		return "Debug Build"
	case BuildTagRelease:
		return "Release Build"
	case BuildTagProduction:
		return "Production Build"
	default:
		return "Unknown"
	}
}

// IsPerformanceRelated returns true for performance-related tags
func (bt BuildTag) IsPerformanceRelated() bool {
	switch bt {
	case BuildTagRace, BuildTagCoverage, BuildTagDebug, BuildTagRelease, BuildTagProduction:
		return true
	default:
		return false
	}
}

// IsCompilationRelated returns true for compilation-related tags
func (bt BuildTag) IsCompilationRelated() bool {
	switch bt {
	case BuildTagPureGo, BuildTagCGoFree, BuildTagNoASM, BuildTagNetgo, BuildTagOsusergo,
		BuildTagNoCgo, BuildTagNoPie, BuildTagStatic, BuildTagStaticAll:
		return true
	default:
		return false
	}
}

// IsLibraryRelated returns true for library-related tags
func (bt BuildTag) IsLibraryRelated() bool {
	switch bt {
	case BuildTagSQLiteFts5, BuildTagSQLiteJSON1, BuildTagLMDB, BuildTagBtrfs:
		return true
	default:
		return false
	}
}

// IsSafe returns true for tags safe for production use
func (bt BuildTag) IsSafe() bool {
	switch bt {
	case BuildTagRace, BuildTagCoverage, BuildTagDebug:
		return false // Not safe for production
	default:
		return true
	}
}

// IsExperimental returns true for experimental tags
func (bt BuildTag) IsExperimental() bool {
	switch bt {
	case BuildTagSQLiteFts5, BuildTagSQLiteJSON1, BuildTagLMDB:
		return true
	default:
		return false
	}
}

// CGOStatus represents CGO compilation status
// This enum replaces bool CGOEnabled for better type safety and semantic clarity
type CGOStatus string

const (
	// CGOStatusDisabled disables CGO compilation completely
	CGOStatusDisabled CGOStatus = "disabled"
	// CGOStatusEnabled enables CGO compilation when available
	CGOStatusEnabled CGOStatus = "enabled"
	// CGOStatusRequired requires CGO compilation and will fail if not available
	CGOStatusRequired CGOStatus = "required"
)

// IsValid returns true if CGOStatus is valid
func (cs CGOStatus) IsValid() bool {
	switch cs {
	case CGOStatusDisabled, CGOStatusEnabled, CGOStatusRequired:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (cs CGOStatus) String() string {
	switch cs {
	case CGOStatusDisabled:
		return "Disabled"
	case CGOStatusEnabled:
		return "Enabled"
	case CGOStatusRequired:
		return "Required"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if CGO is enabled (enabled or required)
func (cs CGOStatus) IsEnabled() bool {
	return cs == CGOStatusEnabled || cs == CGOStatusRequired
}

// IsRequired returns true if CGO is required
func (cs CGOStatus) IsRequired() bool {
	return cs == CGOStatusRequired
}

// IsOptional returns true if CGO is optional (enabled but not required)
func (cs CGOStatus) IsOptional() bool {
	return cs == CGOStatusEnabled
}

// IsDisabled returns true if CGO is disabled
func (cs CGOStatus) IsDisabled() bool {
	return cs == CGOStatusDisabled
}

// RequiresCompiler returns true if requires specific compiler
func (cs CGOStatus) RequiresCompiler() bool {
	return cs.IsEnabled()
}

// GetGCCFlags returns GCC flags needed
func (cs CGOStatus) GetGCCFlags() []string {
	switch cs {
	case CGOStatusDisabled:
		return []string{"-tags=nocgo"}
	case CGOStatusEnabled:
		return []string{}
	case CGOStatusRequired:
		return []string{}
	default:
		return []string{}
	}
}

// GetLdFlags returns linker flags needed
func (cs CGOStatus) GetLdFlags() []string {
	switch cs {
	case CGOStatusDisabled:
		return []string{"-ldflags=-s -w"}
	case CGOStatusEnabled:
		return []string{"-ldflags=-s -w"}
	case CGOStatusRequired:
		return []string{"-ldflags=-s -w"}
	default:
		return []string{}
	}
}

// BuildLevel represents build complexity levels
type BuildLevel string

const (
	BuildLevelMinimal    BuildLevel = "minimal"
	BuildLevelBasic      BuildLevel = "basic"
	BuildLevelStandard   BuildLevel = "standard"
	BuildLevelAdvanced   BuildLevel = "advanced"
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

// GetRecommendedBuildTags returns recommended build tags for level
func (bl BuildLevel) GetRecommendedBuildTags() []BuildTag {
	switch bl {
	case BuildLevelMinimal:
		return []BuildTag{BuildTagPureGo, BuildTagStatic}
	case BuildLevelBasic:
		return []BuildTag{BuildTagPureGo, BuildTagStatic, BuildTagNetgo}
	case BuildLevelStandard:
		return []BuildTag{BuildTagPureGo, BuildTagStatic, BuildTagNetgo, BuildTagOsusergo}
	case BuildLevelAdvanced:
		return []BuildTag{BuildTagStatic, BuildTagNetgo, BuildTagOsusergo}
	case BuildLevelEnterprise:
		return []BuildTag{BuildTagStatic, BuildTagNetgo, BuildTagOsusergo}
	default:
		return []BuildTag{BuildTagPureGo}
	}
}

// GetRecommendedCGOStatus returns recommended CGO status for level
func (bl BuildLevel) GetRecommendedCGOStatus() CGOStatus {
	switch bl {
	case BuildLevelMinimal, BuildLevelBasic:
		return CGOStatusDisabled
	case BuildLevelStandard, BuildLevelAdvanced, BuildLevelEnterprise:
		return CGOStatusEnabled
	default:
		return CGOStatusDisabled
	}
}

// IncludesAdvancedBuilds returns true if level includes advanced builds
func (bl BuildLevel) IncludesAdvancedBuilds() bool {
	return bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}

// IncludesStaticLinking returns true if level includes static linking
func (bl BuildLevel) IncludesStaticLinking() bool {
	switch bl {
	case BuildLevelMinimal, BuildLevelBasic, BuildLevelStandard,
		BuildLevelAdvanced, BuildLevelEnterprise:
		return true
	default:
		return false
	}
}

// RequiresCrossCompilation returns true if level requires cross-compilation
func (bl BuildLevel) RequiresCrossCompilation() bool {
	return bl == BuildLevelAdvanced || bl == BuildLevelEnterprise
}
