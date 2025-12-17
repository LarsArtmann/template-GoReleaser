package domain

// Platform represents supported operating systems
// This enum replaces string-based platform types for type safety
type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformDarwin  Platform = "darwin"
	PlatformWindows Platform = "windows"
	PlatformFreeBSD Platform = "freebsd"
	PlatformOpenBSD Platform = "openbsd"
	PlatformNetBSD  Platform = "netbsd"
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
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
	return p == PlatformLinux || p == PlatformDarwin ||
		p == PlatformFreeBSD || p == PlatformOpenBSD || p == PlatformNetBSD
}

// IsMobile returns true for mobile platforms
func (p Platform) IsMobile() bool {
	return p == PlatformAndroid || p == PlatformIOS
}

// IsDesktop returns true for desktop platforms
func (p Platform) IsDesktop() bool {
	return p == PlatformLinux || p == PlatformDarwin || p == PlatformWindows
}

// IsServer returns true for server platforms
func (p Platform) IsServer() bool {
	return p == PlatformLinux || p == PlatformFreeBSD ||
		p == PlatformOpenBSD || p == PlatformNetBSD
}

// GetFamily returns the OS family
func (p Platform) GetFamily() string {
	switch p {
	case PlatformLinux, PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD:
		return "unix"
	case PlatformDarwin:
		return "darwin"
	case PlatformWindows:
		return "windows"
	case PlatformAndroid:
		return "android"
	case PlatformIOS:
		return "ios"
	default:
		return "unknown"
	}
}

// GetExecutableExtension returns the executable file extension
func (p Platform) GetExecutableExtension() string {
	switch p {
	case PlatformWindows:
		return ".exe"
	default:
		return ""
	}
}

// GetArchiveFormat returns preferred archive format
func (p Platform) GetArchiveFormat() string {
	switch p {
	case PlatformWindows:
		return "zip"
	default:
		return "tar.gz"
	}
}

// Architecture represents supported CPU architectures
// This enum replaces string-based architecture types for type safety
type Architecture string

const (
	Architecture386      Architecture = "386"
	ArchitectureAMD64    Architecture = "amd64"
	ArchitectureARM      Architecture = "arm"
	ArchitectureARM64    Architecture = "arm64"
	ArchitecturePPC64    Architecture = "ppc64"
	ArchitecturePPC64LE  Architecture = "ppc64le"
	ArchitectureS390X    Architecture = "s390x"
	ArchitectureMIPS     Architecture = "mips"
	ArchitectureMIPSLE   Architecture = "mipsle"
	ArchitectureMIPS64   Architecture = "mips64"
	ArchitectureMIPS64LE Architecture = "mips64le"
	ArchitectureRISCV64  Architecture = "riscv64"
)

// IsValid returns true if Architecture is valid
func (a Architecture) IsValid() bool {
	switch a {
	case Architecture386, ArchitectureAMD64, ArchitectureARM, ArchitectureARM64,
		ArchitecturePPC64, ArchitecturePPC64LE, ArchitectureS390X,
		ArchitectureMIPS, ArchitectureMIPSLE, ArchitectureMIPS64, ArchitectureMIPS64LE,
		ArchitectureRISCV64:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (a Architecture) String() string {
	switch a {
	case Architecture386:
		return "x86 (32-bit)"
	case ArchitectureAMD64:
		return "x86_64"
	case ArchitectureARM:
		return "ARM (32-bit)"
	case ArchitectureARM64:
		return "ARM64"
	case ArchitecturePPC64:
		return "PowerPC64"
	case ArchitecturePPC64LE:
		return "PowerPC64LE"
	case ArchitectureS390X:
		return "s390x"
	case ArchitectureMIPS:
		return "MIPS (32-bit)"
	case ArchitectureMIPSLE:
		return "MIPSLE (32-bit)"
	case ArchitectureMIPS64:
		return "MIPS64"
	case ArchitectureMIPS64LE:
		return "MIPS64LE"
	case ArchitectureRISCV64:
		return "RISC-V 64"
	default:
		return "Unknown"
	}
}

// Is64Bit returns true for 64-bit architectures
func (a Architecture) Is64Bit() bool {
	switch a {
	case ArchitectureAMD64, ArchitectureARM64, ArchitecturePPC64, ArchitecturePPC64LE,
		ArchitectureS390X, ArchitectureMIPS64, ArchitectureMIPS64LE, ArchitectureRISCV64:
		return true
	default:
		return false
	}
}

// Is32Bit returns true for 32-bit architectures
func (a Architecture) Is32Bit() bool {
	return !a.Is64Bit()
}

// IsBigEndian returns true for big-endian architectures
func (a Architecture) IsBigEndian() bool {
	switch a {
	case ArchitecturePPC64, ArchitectureS390X, ArchitectureMIPS64:
		return true
	default:
		return false
	}
}

// IsLittleEndian returns true for little-endian architectures
func (a Architecture) IsLittleEndian() bool {
	return !a.IsBigEndian()
}

// GetFamily returns the CPU architecture family
func (a Architecture) GetFamily() string {
	switch a {
	case Architecture386, ArchitectureAMD64:
		return "x86"
	case ArchitectureARM, ArchitectureARM64:
		return "arm"
	case ArchitecturePPC64, ArchitecturePPC64LE:
		return "ppc"
	case ArchitectureS390X:
		return "s390"
	case ArchitectureMIPS, ArchitectureMIPSLE, ArchitectureMIPS64, ArchitectureMIPS64LE:
		return "mips"
	case ArchitectureRISCV64:
		return "riscv"
	default:
		return "unknown"
	}
}

// IsCompatibleWith checks if architecture is compatible with platform
func (a Architecture) IsCompatibleWith(platform Platform) bool {
	// Basic compatibility rules
	switch platform {
	case PlatformLinux:
		// Linux supports most architectures
		return a.IsValid()
	case PlatformDarwin:
		// macOS supports x86_64 and ARM64
		return a == ArchitectureAMD64 || a == ArchitectureARM64
	case PlatformWindows:
		// Windows supports x86 architectures
		return a == Architecture386 || a == ArchitectureAMD64
	case PlatformFreeBSD:
		// FreeBSD supports x86 and ARM
		return a == Architecture386 || a == ArchitectureAMD64 || a == ArchitectureARM64
	case PlatformOpenBSD, PlatformNetBSD:
		// OpenBSD/NetBSD support x86 and ARM64
		return a == Architecture386 || a == ArchitectureAMD64 || a == ArchitectureARM64
	case PlatformAndroid:
		// Android supports ARM
		return a == ArchitectureARM || a == ArchitectureARM64
	case PlatformIOS:
		// iOS supports ARM64
		return a == ArchitectureARM64
	default:
		return false
	}
}

// GetRecommendedArchitectures returns recommended architectures for platform
func GetRecommendedArchitectures(platform Platform) []Architecture {
	switch platform {
	case PlatformLinux:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386}
	case PlatformDarwin:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case PlatformWindows:
		return []Architecture{ArchitectureAMD64, Architecture386}
	case PlatformFreeBSD:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386}
	case PlatformOpenBSD, PlatformNetBSD:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386}
	case PlatformAndroid:
		return []Architecture{ArchitectureARM64, ArchitectureARM}
	case PlatformIOS:
		return []Architecture{ArchitectureARM64}
	default:
		return []Architecture{ArchitectureAMD64}
	}
}

// GetServerRecommendedArchitectures returns recommended server architectures
func GetServerRecommendedArchitectures() []Architecture {
	return []Architecture{ArchitectureAMD64, ArchitectureARM64}
}

// GetDesktopRecommendedArchitectures returns recommended desktop architectures
func GetDesktopRecommendedArchitectures() []Architecture {
	return []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386}
}

// GetMobileRecommendedArchitectures returns recommended mobile architectures
func GetMobileRecommendedArchitectures() []Architecture {
	return []Architecture{ArchitectureARM64, ArchitectureARM}
}
