package domain

// RequiresMainPath returns true if project type requires main path
func (pt ProjectType) RequiresMainPath() bool {
	switch pt {
	case ProjectTypeCLI, ProjectTypeDesktop, ProjectTypeDaemon, ProjectTypeTool:
		return true
	case ProjectTypeLibrary, ProjectTypeWebAPI, ProjectTypeGRPCService, ProjectTypeMicroservice, ProjectTypePlugin, ProjectTypeMobile:
		return false
	default:
		return false
	}
}

// GetGenerateActions returns whether actions should be generated
func (spc *SafeProjectConfig) GetGenerateActions() bool {
	// Simple logic - generate actions for non-library types
	return !spc.ProjectType.IsLibrary()
}
