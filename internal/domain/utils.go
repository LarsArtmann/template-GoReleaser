package domain

// WithField adds field context to error.
func (de *DomainError) WithField(field string) *DomainError {
	return de.WithContext(field)
}

// GetRecommendedProjectType returns recommended project type (CLI as default).
func GetRecommendedProjectType() ProjectType {
	return ProjectTypeCLI
}

// GetAllPlatforms returns all available platforms.
func GetAllPlatforms() []Platform {
	return []Platform{
		PlatformLinux,
		PlatformDarwin,
		PlatformWindows,
		PlatformFreeBSD,
		PlatformOpenBSD,
		PlatformNetBSD,
		PlatformAndroid,
		PlatformIOS,
	}
}
