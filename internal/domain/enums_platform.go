package domain

import "fmt"

// Platform represents supported operating systems
// This enum replaces string-based platform types for type safety
type Platform string

const (
	// PlatformLinux represents Linux distributions
	PlatformLinux Platform = "linux"
	// PlatformDarwin represents macOS (Apple)
	PlatformDarwin Platform = "darwin"
	// PlatformWindows represents Microsoft Windows
	PlatformWindows Platform = "windows"
	// PlatformFreeBSD represents FreeBSD
	PlatformFreeBSD Platform = "freebsd"
	// PlatformOpenBSD represents OpenBSD
	PlatformOpenBSD Platform = "openbsd"
	// PlatformNetBSD represents NetBSD
	PlatformNetBSD Platform = "netbsd"
	// PlatformAndroid represents Android OS
	PlatformAndroid Platform = "android"
	// PlatformIOS represents iOS (Apple mobile)
	PlatformIOS Platform = "ios"
)

// IsValid returns true if Platform is valid
func (p Platform) IsValid() bool {
	switch p {
	case PlatformLinux, PlatformDarwin, PlatformWindows,
		PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD,
		PlatformAndroid, PlatformIOS:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (p Platform) String() string {
	switch p {
	case PlatformLinux:
		return "Linux"
	case PlatformDarwin:
		return "macOS"
	case PlatformWindows:
		return "Windows"
	case PlatformFreeBSD:
		return "FreeBSD"
	case PlatformOpenBSD:
		return "OpenBSD"
	case PlatformNetBSD:
		return "NetBSD"
	case PlatformAndroid:
		return "Android"
	case PlatformIOS:
		return "iOS"
	default:
		return "Unknown"
	}
}

// IsUnix returns true for Unix-like platforms
func (p Platform) IsUnix() bool {
	switch p {
	case PlatformLinux, PlatformDarwin, PlatformFreeBSD,
		PlatformOpenBSD, PlatformNetBSD:
		return true
	default:
		return false
	}
}

// IsMobile returns true for mobile platforms
func (p Platform) IsMobile() bool {
	switch p {
	case PlatformAndroid, PlatformIOS:
		return true
	default:
		return false
	}
}

// IsDesktop returns true for desktop platforms
func (p Platform) IsDesktop() bool {
	return !p.IsMobile()
}

// IsWindows returns true for Windows platform
func (p Platform) IsWindows() bool {
	return p == PlatformWindows
}

// IsApple returns true for Apple platforms
func (p Platform) IsApple() bool {
	switch p {
	case PlatformDarwin, PlatformIOS:
		return true
	default:
		return false
	}
}

// SupportsCGO returns true if platform supports CGO
func (p Platform) SupportsCGO() bool {
	switch p {
	case PlatformLinux, PlatformDarwin, PlatformWindows, PlatformFreeBSD:
		return true
	case PlatformOpenBSD, PlatformNetBSD:
		return false // Limited CGO support
	case PlatformAndroid, PlatformIOS:
		return true // Via NDK/XCode
	default:
		return false
	}
}

// GetDefaultArchitectures returns default architectures for platform
func (p Platform) GetDefaultArchitectures() []Architecture {
	switch p {
	case PlatformLinux:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case PlatformDarwin:
		return []Architecture{ArchitectureARM64, ArchitectureAMD64}
	case PlatformWindows:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD:
		return []Architecture{ArchitectureAMD64}
	case PlatformAndroid:
		return []Architecture{ArchitectureARM64}
	case PlatformIOS:
		return []Architecture{ArchitectureARM64}
	default:
		return []Architecture{}
	}
}

// GetPackageFormat returns preferred package format for platform
func (p Platform) GetPackageFormat() string {
	switch p {
	case PlatformLinux:
		return "tar.gz"
	case PlatformDarwin:
		return "tar.gz"
	case PlatformWindows:
		return "zip"
	case PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD:
		return "tar.gz"
	default:
		return "tar.gz"
	}
}

// GetBinaryExtension returns binary extension for platform
func (p Platform) GetBinaryExtension() string {
	switch p {
	case PlatformWindows:
		return ".exe"
	default:
		return ""
	}
}

// IsProductionReady returns true if platform is production-ready
func (p Platform) IsProductionReady() bool {
	switch p {
	case PlatformLinux, PlatformDarwin, PlatformWindows:
		return true
	case PlatformFreeBSD:
		return true // Stable but limited usage
	case PlatformOpenBSD, PlatformNetBSD:
		return false // Security-focused, limited Go support
	case PlatformAndroid, PlatformIOS:
		return true // Via cross-compilation
	default:
		return false
	}
}

// ValidatePlatform validates a platform
func ValidatePlatform(platform Platform) error {
	if !platform.IsValid() {
		return NewValidationError(
			ErrInvalidPlatform,
			"Invalid platform",
			fmt.Sprintf("'%s' is not a valid platform", platform),
		)
	}
	return nil
}

// GetRecommendedPlatforms returns recommended platforms for projects
func GetRecommendedPlatforms() []Platform {
	return []Platform{
		PlatformLinux,   // Primary server platform
		PlatformDarwin,  // Primary development platform
		PlatformWindows, // Windows development platform
	}
}
