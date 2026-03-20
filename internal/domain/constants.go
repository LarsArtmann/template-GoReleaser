package domain

const (
	maxBuildTags         = 50
	maxDockerNameLength  = 255
	maxDescriptionLength = 255

	complexityCGOEnabled    = 5
	complexityBuildTag      = 2
	complexityDockerEnabled = 10
	complexitySigningTool   = 3

	baseBuildTimeSeconds       = 30
	cgoBuildTimeAdditional     = 10
	dockerBuildTimeAdditional  = 60
	signingBuildTimeAdditional = 5
	baseDependencies           = 2
	dockerDependency           = 1
)
