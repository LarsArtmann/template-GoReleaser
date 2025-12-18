package domain

// GetRecommendedProjectType returns recommended project type (CLI as default)
func GetRecommendedProjectType() ProjectType {
	return ProjectTypeCLI
}