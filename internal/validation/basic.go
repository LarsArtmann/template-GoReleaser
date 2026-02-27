package validation

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/errors"
)

// Validation patterns.
var (
	// Project name validation: alphanumeric, hyphens, underscores, 1-50 chars.
	projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,49}$`)

	// Binary name validation: alphanumeric, hyphens, 1-30 chars.
	binaryNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,29}$`)

	// Main path validation: relative path, no invalid characters.
	mainPathPattern = regexp.MustCompile(`^(\.|\./|\.\./)?[^\\:*?"<>|]+$`)

	// Project description validation: 1-500 chars, printable.
	projectDescriptionPattern = regexp.MustCompile(`^[[:print:]]{1,500}$`)

	// Docker image name validation: lowercase, alphanum, separators, max 255.
	dockerImagePattern = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*$`)

	// Docker registry URL validation.
	dockerRegistryPattern = regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateProjectName validates project name.
func ValidateProjectName(name string) error {
	if name == "" {
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Project name is required",
			"Project name cannot be empty",
		).WithField("project_name").WithSuggestion("Choose a valid project name")
	}

	if len(name) > 50 {
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Project name too long",
			"Project name must be 50 characters or less",
		).WithField("project_name").WithSuggestion("Use a shorter project name")
	}

	if !projectNamePattern.MatchString(name) {
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Invalid project name format",
			"Project name must contain only letters, numbers, hyphens, and underscores",
		).WithField("project_name").WithSuggestion("Use alphanumeric characters with hyphens and underscores only")
	}

	// Check for reserved names
	reservedNames := []string{
		"test", "temp", "tmp", "default", "admin", "root", "system",
		"null", "undefined", "none", "false", "true",
	}

	lowerName := strings.ToLower(name)
	if slices.Contains(reservedNames, lowerName) {
		return errors.NewValidationError(
			errors.ErrInvalidProject,
			"Reserved project name",
			fmt.Sprintf("'%s' is a reserved name", name),
		).WithField("project_name").WithSuggestion("Choose a different project name")
	}

	return nil
}

// ValidateBinaryName validates binary name.
func ValidateBinaryName(name string) error {
	if name == "" {
		return errors.NewValidationError(
			errors.ErrInvalidBinary,
			"Binary name is required",
			"Binary name cannot be empty",
		).WithField("binary_name").WithSuggestion("Choose a valid binary name")
	}

	if len(name) > 30 {
		return errors.NewValidationError(
			errors.ErrInvalidBinary,
			"Binary name too long",
			"Binary name must be 30 characters or less",
		).WithField("binary_name").WithSuggestion("Use a shorter binary name")
	}

	if !binaryNamePattern.MatchString(name) {
		return errors.NewValidationError(
			errors.ErrInvalidBinary,
			"Invalid binary name format",
			"Binary name must contain only letters, numbers, hyphens, and underscores",
		).WithField("binary_name").WithSuggestion("Use alphanumeric characters with hyphens and underscores only")
	}

	// Check for system reserved names
	reservedNames := []string{
		"test", "temp", "tmp", "debug", "release", "build", "install",
		"setup", "config", "init", "run", "start", "stop", "restart",
	}

	lowerName := strings.ToLower(name)
	if slices.Contains(reservedNames, lowerName) {
		return errors.NewValidationError(
			errors.ErrInvalidBinary,
			"Reserved binary name",
			fmt.Sprintf("'%s' is a reserved binary name", name),
		).WithField("binary_name").WithSuggestion("Choose a different binary name")
	}

	return nil
}

// ValidateMainPath validates main path.
func ValidateMainPath(path string) error {
	if path == "" {
		return errors.NewValidationError(
			errors.ErrInvalidMainPath,
			"Main path is required",
			"Main path cannot be empty",
		).WithField("main_path").WithSuggestion("Specify the main package path (e.g., '.', './cmd/app')")
	}

	// Normalize path
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	if len(normalizedPath) > 200 {
		return errors.NewValidationError(
			errors.ErrInvalidMainPath,
			"Main path too long",
			"Main path must be 200 characters or less",
		).WithField("main_path").WithSuggestion("Use a shorter path")
	}

	if !mainPathPattern.MatchString(normalizedPath) {
		return errors.NewValidationError(
			errors.ErrInvalidMainPath,
			"Invalid main path format",
			"Main path must be a valid relative path",
		).WithField("main_path").WithSuggestion("Use relative paths like '.', './cmd/app', or './main'")
	}

	// Check for invalid path components
	invalidComponents := []string{
		"..", "con", "prn", "aux", "nul", "com1", "com2", "com3",
		"com4", "com5", "com6", "com7", "com8", "com9", "lpt1",
		"lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
	}

	components := strings.SplitSeq(normalizedPath, "/")
	for component := range components {
		lowerComponent := strings.ToLower(component)
		if slices.Contains(invalidComponents, lowerComponent) {
			return errors.NewValidationError(
				errors.ErrInvalidMainPath,
				"Invalid path component",
				fmt.Sprintf("'%s' is not allowed in path", component),
			).WithField("main_path").WithSuggestion("Remove invalid components from path")
		}
	}

	return nil
}

// ValidateProjectDescription validates project description.
func ValidateProjectDescription(description string) error {
	if description == "" {
		return nil // Description is optional
	}

	trimmed := strings.TrimSpace(description)

	if len(trimmed) == 0 {
		return errors.NewValidationError(
			errors.ErrInvalidProjectDescription,
			"Project description cannot be empty",
			"Project description must contain text",
		).WithField("project_description").WithSuggestion("Provide a meaningful description")
	}

	if len(trimmed) > 500 {
		return errors.NewValidationError(
			errors.ErrInvalidProjectDescription,
			"Project description too long",
			"Project description must be 500 characters or less",
		).WithField("project_description").WithSuggestion("Use a shorter description")
	}

	if !projectDescriptionPattern.MatchString(trimmed) {
		return errors.NewValidationError(
			errors.ErrInvalidProjectDescription,
			"Invalid characters in project description",
			"Project description contains invalid characters",
		).WithField("project_description").WithSuggestion("Use only printable characters")
	}

	// Check for common description issues
	if strings.Contains(trimmed, "TODO") || strings.Contains(trimmed, "FIXME") {
		return errors.NewValidationError(
			errors.ErrInvalidProjectDescription,
			"Incomplete project description",
			"Project description contains TODO or FIXME markers",
		).WithField("project_description").WithSuggestion("Complete the project description")
	}

	return nil
}

// ValidateDockerImageName validates Docker image name.
func ValidateDockerImageName(name string) error {
	if name == "" {
		return errors.NewValidationError(
			errors.ErrInvalidDockerImage,
			"Docker image name is required",
			"Docker image name cannot be empty",
		).WithField("docker_image").WithSuggestion("Specify a Docker image name")
	}

	// Split name and tag
	parts := strings.SplitN(name, ":", 2)
	imageName := parts[0]

	if len(imageName) > 255 {
		return errors.NewValidationError(
			errors.ErrInvalidDockerImage,
			"Docker image name too long",
			"Docker image name must be 255 characters or less",
		).WithField("docker_image").WithSuggestion("Use a shorter image name")
	}

	if !dockerImagePattern.MatchString(imageName) {
		return errors.NewValidationError(
			errors.ErrInvalidDockerImage,
			"Invalid Docker image name format",
			"Docker image name must contain only lowercase letters, numbers, dots, hyphens, and underscores",
		).WithField("docker_image").WithSuggestion("Use lowercase alphanumeric characters with valid separators")
	}

	// Validate tag if present
	if len(parts) == 2 {
		tag := parts[1]
		if len(tag) > 128 {
			return errors.NewValidationError(
				errors.ErrInvalidDockerImage,
				"Docker image tag too long",
				"Docker image tag must be 128 characters or less",
			).WithField("docker_image").WithSuggestion("Use a shorter tag")
		}

		if !strings.HasPrefix(tag, "v") &&
			!regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(tag) {
			return errors.NewValidationError(
				errors.ErrInvalidDockerImage,
				"Invalid Docker image tag format",
				"Docker image tag must start with 'v' for version tags or contain only valid characters",
			).WithField("docker_image").WithSuggestion("Use semantic versioning (e.g., v1.0.0)")
		}
	}

	// Check for reserved image names
	reservedNames := []string{
		"latest", "stable", "production", "master", "main", "develop",
		"test", "temp", "debug", "dev", "staging",
	}

	imageNameLower := strings.ToLower(imageName)
	for _, reserved := range reservedNames {
		if strings.Contains(imageNameLower, reserved) && imageNameLower != reserved {
			continue // Allow containing but not equal to reserved names
		}

		if imageNameLower == reserved {
			return errors.NewValidationError(
				errors.ErrInvalidDockerImage,
				"Reserved Docker image name",
				fmt.Sprintf("'%s' is a reserved Docker image name", name),
			).WithField("docker_image").WithSuggestion("Choose a different image name")
		}
	}

	return nil
}

// ValidateDockerRegistry validates Docker registry URL.
func ValidateDockerRegistry(registry string) error {
	if registry == "" {
		return nil // Registry can be empty (use default)
	}

	if len(registry) > 253 {
		return errors.NewValidationError(
			errors.ErrInvalidDockerRegistry,
			"Docker registry URL too long",
			"Docker registry URL must be 253 characters or less",
		).WithField("docker_registry").WithSuggestion("Use a shorter registry URL")
	}

	if !dockerRegistryPattern.MatchString(registry) {
		return errors.NewValidationError(
			errors.ErrInvalidDockerRegistry,
			"Invalid Docker registry URL format",
			"Docker registry URL must be a valid domain name",
		).WithField("docker_registry").WithSuggestion("Use a valid domain name (e.g., registry.example.com)")
	}

	// Check for common registry domains
	knownRegistries := []string{
		"docker.io", "ghcr.io", "registry.gitlab.com", "gcr.io",
		"azurecr.io", "quay.io", "mcr.microsoft.com",
	}

	for _, known := range knownRegistries {
		if strings.EqualFold(registry, known) {
			return nil // Known registry is valid
		}
	}

	// For unknown registries, provide a warning
	// In production, you might want to validate the registry exists
	return nil
}

// ValidateVersion validates version string.
func ValidateVersion(version string) error {
	if version == "" {
		return errors.NewValidationError(
			errors.ErrInvalidVersion,
			"Version is required",
			"Version cannot be empty",
		).WithField("version").WithSuggestion("Specify a version (e.g., v1.0.0)")
	}

	// Accept semantic versioning or git describe format
	semverPattern := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$`)
	gitDescribePattern := regexp.MustCompile(`^[a-f0-9]{7,40}(-dirty)?$`)

	if !semverPattern.MatchString(version) && !gitDescribePattern.MatchString(version) {
		return errors.NewValidationError(
			errors.ErrInvalidVersion,
			"Invalid version format",
			"Version must follow semantic versioning (e.g., v1.0.0) or be a git commit hash",
		).WithField("version").WithSuggestion("Use semantic versioning (e.g., v1.0.0, v1.0.0-alpha.1)")
	}

	return nil
}

// ValidateGitBranch validates Git branch name.
func ValidateGitBranch(branch string) error {
	if branch == "" {
		return errors.NewValidationError(
			errors.ErrInvalidGitBranch,
			"Git branch name is required",
			"Git branch name cannot be empty",
		).WithField("git_branch").WithSuggestion("Specify a Git branch name")
	}

	branchPattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9/_-]{0,252}$`)
	if !branchPattern.MatchString(branch) {
		return errors.NewValidationError(
			errors.ErrInvalidGitBranch,
			"Invalid Git branch name format",
			"Git branch name must contain only letters, numbers, hyphens, underscores, and forward slashes",
		).WithField("git_branch").WithSuggestion("Use valid Git branch naming conventions")
	}

	// Check for invalid branch names
	invalidNames := []string{
		"HEAD", "master", "main", "develop", "feature", "release", "hotfix",
		"MERGE_HEAD", "ORIG_HEAD", "FETCH_HEAD",
	}

	lowerBranch := strings.ToLower(branch)
	if slices.Contains(invalidNames, lowerBranch) {
		return errors.NewValidationError(
			errors.ErrInvalidGitBranch,
			"Reserved Git branch name",
			fmt.Sprintf("'%s' is a reserved branch name", branch),
		).WithField("git_branch").WithSuggestion("Use a descriptive branch name")
	}

	return nil
}

// ValidateGitTag validates Git tag name.
func ValidateGitTag(tag string) error {
	if tag == "" {
		return errors.NewValidationError(
			errors.ErrInvalidGitTag,
			"Git tag name is required",
			"Git tag name cannot be empty",
		).WithField("git_tag").WithSuggestion("Specify a Git tag name")
	}

	tagPattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,254}$`)
	if !tagPattern.MatchString(tag) {
		return errors.NewValidationError(
			errors.ErrInvalidGitTag,
			"Invalid Git tag format",
			"Git tag name must contain only letters, numbers, dots, hyphens, and underscores",
		).WithField("git_tag").WithSuggestion("Use valid Git tag naming conventions")
	}

	return nil
}

// ValidateBuildTags validates build tags.
func ValidateBuildTags(tags []string) error {
	if len(tags) == 0 {
		return nil // No build tags is valid
	}

	for i, tag := range tags {
		if tag == "" {
			return errors.NewValidationError(
				errors.ErrInvalidBuildTag,
				"Build tag cannot be empty",
				fmt.Sprintf("Build tag at index %d is empty", i),
			).WithField("build_tags").WithSuggestion("Remove empty build tags")
		}

		// Validate build tag format
		tagPattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
		if !tagPattern.MatchString(tag) {
			return errors.NewValidationError(
				errors.ErrInvalidBuildTag,
				"Invalid build tag format",
				fmt.Sprintf("Build tag '%s' is invalid", tag),
			).WithField("build_tags").WithSuggestion("Use alphanumeric characters with hyphens and underscores only")
		}

		if len(tag) > 50 {
			return errors.NewValidationError(
				errors.ErrInvalidBuildTag,
				"Build tag too long",
				fmt.Sprintf("Build tag '%s' exceeds 50 characters", tag),
			).WithField("build_tags").WithSuggestion("Use shorter build tags")
		}
	}

	// Check for duplicate build tags
	seen := make(map[string]bool)
	for _, tag := range tags {
		if seen[tag] {
			return errors.NewValidationError(
				errors.ErrInvalidBuildTag,
				"Duplicate build tag",
				fmt.Sprintf("Build tag '%s' is specified multiple times", tag),
			).WithField("build_tags").WithSuggestion("Remove duplicate build tags")
		}

		seen[tag] = true
	}

	return nil
}

// ValidatePort validates port number.
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return errors.NewValidationError(
			errors.ErrInvalidPort,
			"Invalid port number",
			"Port number must be between 1 and 65535",
		).WithField("port").WithSuggestion("Use a valid port number (1-65535)")
	}

	// Check for well-known ports that shouldn't be used by applications
	restrictedPorts := []int{
		1, 7, 9, 11, 13, 17, 19, 20, 21, 22, 23, 25,
		53, 67, 68, 69, 80, 110, 123, 135, 137, 138, 139,
		143, 161, 162, 179, 389, 443, 445, 512, 513, 514,
		515, 520, 521, 522, 523, 524, 525, 526, 527, 528,
		529, 530, 531, 532, 533, 534, 535, 536, 537, 538,
		539, 540, 541, 542, 543, 544, 545, 546, 547, 548,
		549, 550, 551, 552, 553, 554, 555, 556, 557, 558,
		559, 560, 561, 562, 563, 564, 565, 566, 567, 568,
		569, 570, 571, 572, 573, 574, 575, 576, 577, 578,
		579, 580, 581, 582, 583, 584, 585, 586, 587, 588,
		589, 590, 591, 592, 593, 594, 595, 596, 597, 598,
		599, 600, 601, 602, 603, 604, 605, 606, 607, 608,
		609, 610, 611, 612, 613, 614, 615, 616, 617, 618,
		619, 620, 621, 622, 623, 624, 625, 626, 627, 628,
		629, 630, 631, 632, 633, 634, 635, 636, 637, 638,
		639, 640, 641, 642, 643, 644, 645, 646, 647, 648,
		649, 650, 651, 652, 653, 654, 655,
	}

	if slices.Contains(restrictedPorts, port) {
		return errors.NewValidationError(
			errors.ErrInvalidPort,
			"Restricted port number",
			fmt.Sprintf("Port %d is restricted and shouldn't be used by applications", port),
		).WithField("port").WithSuggestion("Use a non-privileged port (1024-65535)")
	}

	return nil
}
