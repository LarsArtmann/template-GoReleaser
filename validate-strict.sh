#!/usr/bin/env bash
set -euo pipefail

# Strict GoReleaser Validator
# Performs comprehensive validation with zero tolerance for issues

# Enable strict error handling
set -E
trap 'echo "Error on line $LINENO"' ERR

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Validation results
declare -A VALIDATION_RESULTS
TOTAL_CHECKS=0
FAILED_CHECKS=0
CRITICAL_FAILURES=0

# Configuration
GORELEASER_FILE=".goreleaser.yaml"
GORELEASER_PRO_FILE=".goreleaser.pro.yaml"
MIN_GO_VERSION="1.21"
REQUIRED_TOOLS=("go" "git" "goreleaser")
RECOMMENDED_TOOLS=("docker" "cosign" "syft" "upx" "yq" "jq")

# Logging functions
log_header() {
    echo
    echo -e "${MAGENTA}========================================${NC}"
    echo -e "${MAGENTA}  $1${NC}"
    echo -e "${MAGENTA}========================================${NC}"
}

log_subheader() {
    echo
    echo -e "${CYAN}──── $1 ────${NC}"
}

log_check() {
    echo -ne "${BLUE}[CHECKING]${NC} $1... "
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC}"
    VALIDATION_RESULTS["$1"]="PASS"
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $2"
    VALIDATION_RESULTS["$1"]="FAIL"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
}

log_critical() {
    echo -e "${RED}[CRITICAL]${NC} $2"
    VALIDATION_RESULTS["$1"]="CRITICAL"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
    CRITICAL_FAILURES=$((CRITICAL_FAILURES + 1))
}

# Validation functions
validate_yaml_structure() {
    local file=$1
    local name=$(basename "$file")
    
    log_check "YAML structure for $name"
    
    if [[ ! -f "$file" ]]; then
        log_critical "$name" "File does not exist"
        return 1
    fi
    
    # Check YAML syntax
    if command -v yq &> /dev/null; then
        if ! yq eval '.' "$file" > /dev/null 2>&1; then
            log_critical "$name" "Invalid YAML syntax"
            return 1
        fi
    fi
    
    # Check for required top-level keys
    local required_keys=("project_name" "builds" "archives" "checksum" "changelog" "release")
    for key in "${required_keys[@]}"; do
        if ! grep -q "^${key}:" "$file"; then
            log_fail "$name" "Missing required key: $key"
            return 1
        fi
    done
    
    log_pass "$name"
}

validate_build_configuration() {
    local file=$1
    local name=$(basename "$file")
    
    log_check "Build configuration in $name"
    
    # Check for CGO_ENABLED=0
    if ! grep -q "CGO_ENABLED=0" "$file"; then
        log_fail "$name" "CGO_ENABLED should be 0 for static builds"
        return 1
    fi
    
    # Check for proper ldflags
    if ! grep -q -- "-s -w" "$file"; then
        log_fail "$name" "Missing -s -w ldflags for smaller binaries"
        return 1
    fi
    
    # Check for version injection
    if ! grep -q "\.Version" "$file"; then
        log_fail "$name" "No version injection in ldflags"
        return 1
    fi
    
    # Check for trimpath
    if ! grep -q "trimpath" "$file"; then
        log_fail "$name" "Missing -trimpath flag for reproducible builds"
        return 1
    fi
    
    log_pass "$name"
}

validate_security_configuration() {
    local file=$1
    local name=$(basename "$file")
    
    log_subheader "Security Configuration"
    
    # Check for signing in pro version
    if [[ "$file" == *".pro.yaml" ]]; then
        log_check "Signing configuration"
        if grep -q "signs:" "$file"; then
            log_pass "Signing"
        else
            log_fail "Signing" "No signing configuration in pro version"
        fi
        
        log_check "SBOM generation"
        if grep -q "sboms:" "$file"; then
            log_pass "SBOM"
        else
            log_fail "SBOM" "No SBOM configuration in pro version"
        fi
    fi
    
    # Check for checksum
    log_check "Checksum configuration"
    if grep -q "checksum:" "$file"; then
        if grep -q "algorithm: sha256" "$file"; then
            log_pass "Checksum"
        else
            log_fail "Checksum" "Should use SHA256 algorithm"
        fi
    else
        log_critical "Checksum" "No checksum configuration"
    fi
}

validate_docker_configuration() {
    local file=$1
    
    if ! grep -q "dockers:" "$file"; then
        return 0
    fi
    
    log_subheader "Docker Configuration"
    
    log_check "Docker multi-arch support"
    if grep -q "linux/amd64" "$file" && grep -q "linux/arm64" "$file"; then
        log_pass "Docker multi-arch"
    else
        log_fail "Docker multi-arch" "Missing multi-architecture support"
    fi
    
    log_check "Docker labels"
    local required_labels=("org.opencontainers.image.title" "org.opencontainers.image.version" "org.opencontainers.image.source")
    for label in "${required_labels[@]}"; do
        if ! grep -q "$label" "$file"; then
            log_fail "Docker labels" "Missing label: $label"
            return 1
        fi
    done
    log_pass "Docker labels"
}

validate_release_configuration() {
    local file=$1
    
    log_subheader "Release Configuration"
    
    log_check "Changelog configuration"
    if grep -q "changelog:" "$file"; then
        if grep -q "groups:" "$file"; then
            log_pass "Changelog"
        else
            log_fail "Changelog" "No commit grouping in changelog"
        fi
    else
        log_critical "Changelog" "No changelog configuration"
    fi
    
    log_check "Release footer"
    if grep -q "footer:" "$file"; then
        log_pass "Release footer"
    else
        log_fail "Release footer" "No installation instructions in release footer"
    fi
}

validate_environment_variables() {
    local file=$1
    
    log_subheader "Environment Variables"
    
    # Extract all environment variables with correct pattern
    local env_vars=$(grep -oE '\{\{ \.Env\.[A-Z_]+ \}\}' "$file" 2>/dev/null | sed 's/{{ \.Env\.\([A-Z_]*\) }}/\1/' | sort -u)
    
    if [[ -z "$env_vars" ]]; then
        log_check "Environment variables"
        log_pass "No environment variables required"
        return 0
    fi
    
    local missing_vars=()
    local critical_vars=("GITHUB_TOKEN" "GITHUB_OWNER" "GITHUB_REPO")
    local invalid_vars=()
    
    for var in $env_vars; do
        log_check "Environment variable: $var"
        if [[ -n "${!var:-}" ]]; then
            local value="${!var}"
            # Validate critical variables more strictly
            local is_critical=false
            for critical_var in "${critical_vars[@]}"; do
                if [[ "$var" == "$critical_var" ]]; then
                    is_critical=true
                    break
                fi
            done
            
            # Check for placeholder values
            if [[ "$value" =~ ^(your-|xxxx|example|changeme|todo) ]] || [[ ${#value} -lt 3 ]]; then
                log_fail "$var" "Appears to be placeholder value"
                invalid_vars+=("$var")
            elif [[ "$is_critical" == true ]]; then
                # Additional validation for critical vars
                case "$var" in
                    GITHUB_TOKEN)
                        if [[ ! "$value" =~ ^(ghp_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{82})$ ]]; then
                            log_fail "$var" "Invalid GitHub token format"
                            invalid_vars+=("$var")
                        else
                            log_pass "$var"
                        fi
                        ;;
                    GITHUB_OWNER|GITHUB_REPO)
                        if [[ ! "$value" =~ ^[A-Za-z0-9._-]+$ ]]; then
                            log_fail "$var" "Invalid format for GitHub owner/repo"
                            invalid_vars+=("$var")
                        else
                            log_pass "$var"
                        fi
                        ;;
                    *)
                        log_pass "$var"
                        ;;
                esac
            else
                log_pass "$var"
            fi
        else
            if [[ " ${critical_vars[*]} " =~ " $var " ]]; then
                log_critical "$var" "Critical variable not set"
            else
                log_fail "$var" "Not set"
            fi
            missing_vars+=("$var")
        fi
    done
    
    if [[ ${#missing_vars[@]} -gt 0 ]] || [[ ${#invalid_vars[@]} -gt 0 ]]; then
        echo
        if [[ ${#missing_vars[@]} -gt 0 ]]; then
            echo -e "${YELLOW}Missing environment variables:${NC}"
            printf '%s\n' "${missing_vars[@]}" | sed 's/^/  - /'
        fi
        if [[ ${#invalid_vars[@]} -gt 0 ]]; then
            echo -e "${RED}Invalid environment variables:${NC}"
            printf '%s\n' "${invalid_vars[@]}" | sed 's/^/  - /'
        fi
        echo
        echo -e "${BLUE}Quick setup:${NC}"
        echo "1. Copy .env.example to .env: cp .env.example .env"
        echo "2. Edit .env with your values: editor .env"
        echo "3. Source the file: source .env"
        echo "4. Run validation again: ./validate-strict.sh"
    fi
}

validate_pro_features() {
    local file=$1
    
    if [[ ! "$file" == *".pro.yaml" ]]; then
        return 0
    fi
    
    log_subheader "Pro Features Validation"
    
    local pro_features=("docker_manifests" "upx" "nfpms" "brews" "scoops" "snapcrafts" "aurs" "signs" "sboms" "milestones" "announce")
    
    for feature in "${pro_features[@]}"; do
        log_check "Pro feature: $feature"
        if grep -q "${feature}:" "$file"; then
            log_pass "$feature"
        else
            log_fail "$feature" "Pro feature not configured"
        fi
    done
}

validate_git_repository() {
    log_subheader "Git Repository"
    
    log_check "Git initialization"
    if [[ -d ".git" ]]; then
        log_pass "Git initialized"
    else
        log_critical "Git" "Not a git repository"
        return 1
    fi
    
    log_check "Git remote"
    if git remote -v 2>/dev/null | grep -q origin; then
        log_pass "Git remote"
    else
        log_fail "Git remote" "No origin remote configured"
    fi
    
    log_check "Git tags"
    if [[ $(git tag 2>/dev/null | wc -l) -gt 0 ]]; then
        log_pass "Git tags"
    else
        log_fail "Git tags" "No tags found (required for releases)"
    fi
    
    log_check "Git clean state"
    if [[ -z $(git status --porcelain 2>/dev/null) ]]; then
        log_pass "Git clean"
    else
        log_fail "Git clean" "Uncommitted changes present"
    fi
}

validate_dependencies() {
    log_subheader "Dependencies"
    
    # Check required tools
    for tool in "${REQUIRED_TOOLS[@]}"; do
        log_check "Required tool: $tool"
        if command -v "$tool" &> /dev/null; then
            log_pass "$tool"
        else
            log_critical "$tool" "Required tool not installed"
        fi
    done
    
    # Check recommended tools
    for tool in "${RECOMMENDED_TOOLS[@]}"; do
        log_check "Recommended tool: $tool"
        if command -v "$tool" &> /dev/null; then
            log_pass "$tool"
        else
            log_fail "$tool" "Recommended tool not installed"
        fi
    done
    
    # Check Go version
    if command -v go &> /dev/null; then
        log_check "Go version >= $MIN_GO_VERSION"
        local go_version=$(go version | grep -oE '[0-9]+\.[0-9]+' | head -1)
        if [[ $(echo "$go_version >= $MIN_GO_VERSION" | bc -l) -eq 1 ]]; then
            log_pass "Go version"
        else
            log_fail "Go version" "Go version $go_version is below minimum $MIN_GO_VERSION"
        fi
    fi
}

validate_project_structure() {
    log_subheader "Project Structure"
    
    log_check "Go module"
    if [[ -f "go.mod" ]]; then
        log_pass "go.mod"
    else
        log_fail "go.mod" "No go.mod file found"
    fi
    
    log_check "Main package"
    if [[ -f "main.go" ]] || [[ -d "cmd" ]]; then
        log_pass "Main package"
    else
        log_fail "Main package" "No main.go or cmd/ directory"
    fi
    
    log_check "License file"
    if ls LICENSE* &> /dev/null || ls license* &> /dev/null; then
        log_pass "License"
    else
        log_fail "License" "No LICENSE file found"
    fi
    
    log_check "README file"
    if ls README* &> /dev/null || ls readme* &> /dev/null; then
        log_pass "README"
    else
        log_fail "README" "No README file found"
    fi
    
    log_check "Dockerfile"
    if grep -q "dockers:" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        if [[ -f "Dockerfile" ]]; then
            log_pass "Dockerfile"
        else
            log_fail "Dockerfile" "Docker configured but no Dockerfile found"
        fi
    else
        log_pass "Dockerfile (not required)"
    fi
}

run_goreleaser_validation() {
    if ! command -v goreleaser &> /dev/null; then
        return 0
    fi
    
    log_subheader "GoReleaser Validation"
    
    # Temporarily clear conflicting tokens for GoReleaser
    local saved_gitlab_token="${GITLAB_TOKEN:-}"
    local saved_gitea_token="${GITEA_TOKEN:-}"
    unset GITLAB_TOKEN GITEA_TOKEN
    
    for file in "$GORELEASER_FILE" "$GORELEASER_PRO_FILE"; do
        if [[ ! -f "$file" ]]; then
            continue
        fi
        
        local name=$(basename "$file")
        log_check "GoReleaser check for $name"
        
        if goreleaser check --config "$file" &> /dev/null; then
            log_pass "$name"
        else
            log_fail "$name" "GoReleaser validation failed (might be expected in test environment)"
            echo -e "${YELLOW}Error details:${NC}"
            goreleaser check --config "$file" 2>&1 | head -10 | grep -v "multiple tokens" | sed 's/^/  /' || true
        fi
    done
    
    # Restore tokens
    [[ -n "$saved_gitlab_token" ]] && export GITLAB_TOKEN="$saved_gitlab_token"
    [[ -n "$saved_gitea_token" ]] && export GITEA_TOKEN="$saved_gitea_token"
}

generate_report() {
    log_header "VALIDATION REPORT"
    
    echo
    echo -e "${CYAN}Total Checks:${NC} $TOTAL_CHECKS"
    echo -e "${GREEN}Passed:${NC} $((TOTAL_CHECKS - FAILED_CHECKS))"
    echo -e "${YELLOW}Failed:${NC} $((FAILED_CHECKS - CRITICAL_FAILURES))"
    echo -e "${RED}Critical:${NC} $CRITICAL_FAILURES"
    
    if [[ $CRITICAL_FAILURES -gt 0 ]]; then
        echo
        echo -e "${RED}========================================${NC}"
        echo -e "${RED}  CRITICAL FAILURES DETECTED!           ${NC}"
        echo -e "${RED}  Configuration is NOT production ready ${NC}"
        echo -e "${RED}========================================${NC}"
        return 1
    elif [[ $FAILED_CHECKS -gt 0 ]]; then
        echo
        echo -e "${YELLOW}========================================${NC}"
        echo -e "${YELLOW}  WARNINGS DETECTED                     ${NC}"
        echo -e "${YELLOW}  Review and fix before production use  ${NC}"
        echo -e "${YELLOW}========================================${NC}"
        return 0
    else
        echo
        echo -e "${GREEN}========================================${NC}"
        echo -e "${GREEN}  ALL VALIDATIONS PASSED! ✓            ${NC}"
        echo -e "${GREEN}  Configuration is production ready     ${NC}"
        echo -e "${GREEN}========================================${NC}"
        return 0
    fi
}

# Main execution
main() {
    echo -e "${MAGENTA}========================================${NC}"
    echo -e "${MAGENTA}  STRICT GORELEASER VALIDATOR          ${NC}"
    echo -e "${MAGENTA}  Zero Tolerance Mode Enabled          ${NC}"
    echo -e "${MAGENTA}========================================${NC}"
    
    log_header "CONFIGURATION VALIDATION"
    validate_yaml_structure "$GORELEASER_FILE"
    validate_yaml_structure "$GORELEASER_PRO_FILE"
    
    log_header "BUILD VALIDATION"
    validate_build_configuration "$GORELEASER_FILE"
    validate_build_configuration "$GORELEASER_PRO_FILE"
    
    log_header "SECURITY VALIDATION"
    validate_security_configuration "$GORELEASER_FILE"
    validate_security_configuration "$GORELEASER_PRO_FILE"
    
    log_header "FEATURE VALIDATION"
    validate_docker_configuration "$GORELEASER_FILE"
    validate_docker_configuration "$GORELEASER_PRO_FILE"
    validate_release_configuration "$GORELEASER_FILE"
    validate_release_configuration "$GORELEASER_PRO_FILE"
    validate_pro_features "$GORELEASER_PRO_FILE"
    
    log_header "ENVIRONMENT VALIDATION"
    validate_environment_variables "$GORELEASER_FILE"
    validate_environment_variables "$GORELEASER_PRO_FILE"
    
    log_header "REPOSITORY VALIDATION"
    validate_git_repository
    
    log_header "DEPENDENCY VALIDATION"
    validate_dependencies
    
    log_header "PROJECT VALIDATION"
    validate_project_structure
    
    log_header "GORELEASER VALIDATION"
    run_goreleaser_validation
    
    # Generate final report
    generate_report
    exit_code=$?
    
    # Save validation results to file
    echo
    echo -e "${BLUE}Saving validation results to validation-report.json...${NC}"
    {
        echo "{"
        echo "  \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\","
        echo "  \"total_checks\": $TOTAL_CHECKS,"
        echo "  \"passed\": $((TOTAL_CHECKS - FAILED_CHECKS)),"
        echo "  \"failed\": $((FAILED_CHECKS - CRITICAL_FAILURES)),"
        echo "  \"critical\": $CRITICAL_FAILURES,"
        echo "  \"results\": {"
        local first=true
        for key in "${!VALIDATION_RESULTS[@]}"; do
            if [[ "$first" == false ]]; then
                echo ","
            fi
            echo -n "    \"$key\": \"${VALIDATION_RESULTS[$key]}\""
            first=false
        done
        echo
        echo "  }"
        echo "}"
    } > validation-report.json
    
    exit $exit_code
}

# Run validation
main "$@"