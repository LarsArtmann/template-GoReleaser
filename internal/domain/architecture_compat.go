package domain

// IsCompatibleWith returns true if architecture is compatible with platform.
func (a Architecture) IsCompatibleWith(platform Platform) bool {
	// Basic compatibility checks
	switch platform {
	case PlatformLinux:
		return a.IsValid() // Most architectures work on Linux
	case PlatformDarwin:
		return a == ArchitectureAMD64 || a == ArchitectureARM64
	case PlatformWindows:
		return a == ArchitectureAMD64 || a == Architecture386
	case PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD:
		return a == ArchitectureAMD64 || a == ArchitectureARM64
	case PlatformAndroid:
		return a == ArchitectureARM64
	case PlatformIOS:
		return a == ArchitectureARM64
	default:
		return false
	}
}
