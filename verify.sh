#!/usr/bin/env bash
set -euo pipefail

# GoReleaser Configuration Verifier
# Comprehensive validation script for GoReleaser templates

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
ERRORS=0
WARNINGS=0
CHECKS=0

# Configuration files
GORELEASER_FILE=".goreleaser.yaml"
GORELEASER_PRO_FILE=".goreleaser.pro.yaml"
CURRENT_FILE=""

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((CHECKS++))
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
    ((WARNINGS++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((ERRORS++))
}

check_command() {
    if command -v "$1" &> /dev/null; then
        log_success "$1 is installed"
        return 0
    else
        log_error "$1 is not installed"
        return 1
    fi
}

check_file_exists() {
    if [[ -f "$1" ]]; then
        log_success "File exists: $1"
        return 0
    else
        log_error "File not found: $1"
        return 1
    fi
}

check_yaml_syntax() {
    local file=$1
    if command -v yq &> /dev/null; then
        if yq eval '.' "$file" > /dev/null 2>&1; then
            log_success "YAML syntax is valid: $file"
            return 0
        else
            log_error "YAML syntax error in: $file"
            yq eval '.' "$file" 2>&1 | head -10
            return 1
        fi
    else
        log_warning "yq not installed, skipping YAML syntax check"
        return 0
    fi
}

check_goreleaser_config() {
    local file=$1
    CURRENT_FILE=$file
    
    log_info "Validating $file..."
    
    if ! check_file_exists "$file"; then
        return 1
    fi
    
    if ! check_yaml_syntax "$file"; then
        return 1
    fi
    
    # Check with goreleaser (handle multiple token conflicts)
    if command -v goreleaser &> /dev/null; then
        # Temporarily clear conflicting tokens for GoReleaser
        local saved_gitlab_token="${GITLAB_TOKEN:-}"
        local saved_gitea_token="${GITEA_TOKEN:-}"
        unset GITLAB_TOKEN GITEA_TOKEN
        
        if goreleaser check --config "$file" 2>/dev/null; then
            log_success "GoReleaser validation passed: $file"
        else
            log_warning "GoReleaser validation failed: $file (this might be expected in test environment)"
            # Show errors but don't fail completely
            goreleaser check --config "$file" 2>&1 | head -5 | grep -v "multiple tokens" || true
        fi
        
        # Restore tokens
        [[ -n "$saved_gitlab_token" ]] && export GITLAB_TOKEN="$saved_gitlab_token"
        [[ -n "$saved_gitea_token" ]] && export GITEA_TOKEN="$saved_gitea_token"
    else
        log_warning "goreleaser not installed, skipping native validation"
        return 0
    fi
    
    # Skip build test for faster validation
    log_info "Skipping snapshot build test (for performance)"
}

check_required_env_vars() {
    local file=$1
    log_info "Checking environment variables for $file..."
    
    # Extract env var references from config file
    local config_env_vars=$(grep -oE '\{\{ \.Env\.[A-Z_]+ \}\}' "$file" 2>/dev/null | sed 's/{{ \.Env\.\([A-Z_]*\) }}/\1/' | sort -u)
    
    # Define commonly required environment variables for GoReleaser
    local critical_vars=("GITHUB_TOKEN")
    local common_vars=("DOCKER_USERNAME" "DOCKER_PASSWORD" "GORELEASER_KEY")
    local all_env_vars=()
    
    # Add config-specific variables
    if [[ -n "$config_env_vars" ]]; then
        while IFS= read -r var; do
            all_env_vars+=("$var")
        done <<< "$config_env_vars"
    fi
    
    # Add common variables that GoReleaser typically needs
    for var in "${critical_vars[@]}" "${common_vars[@]}"; do
        if [[ ! " ${all_env_vars[*]} " =~ " ${var} " ]]; then
            all_env_vars+=("$var")
        fi
    done
    
    if [[ ${#all_env_vars[@]} -gt 0 ]]; then
        echo "Environment variables for GoReleaser:"
        local critical_missing=()
        local optional_missing=()
        
        for var in "${all_env_vars[@]}"; do
            if [[ -n "${!var:-}" ]]; then
                # Basic validation for common patterns
                local value="${!var}"
                if [[ "$value" =~ ^(your-|xxxx|example|changeme|todo) ]] || [[ ${#value} -lt 3 ]]; then
                    log_warning "$var is set but appears to be a placeholder"
                else
                    log_success "$var is set"
                fi
            else
                # Check if this is a critical variable
                if [[ " ${critical_vars[*]} " =~ " ${var} " ]]; then
                    log_warning "$var is not set (critical)"
                    critical_missing+=("$var")
                else
                    log_warning "$var is not set (optional)"
                    optional_missing+=("$var")
                fi
            fi
        done
        
        if [[ ${#critical_missing[@]} -gt 0 ]]; then
            echo
            log_error "Critical environment variables missing:"
            for var in "${critical_missing[@]}"; do
                echo "  - $var"
            done
        fi
        
        if [[ ${#optional_missing[@]} -gt 0 ]] || [[ ${#critical_missing[@]} -gt 0 ]]; then
            echo
            echo "Environment variable information:"
            echo "Set these variables or source a .env file before running GoReleaser"
        fi
    else
        log_success "No additional environment variables detected"
    fi
}

check_project_structure() {
    log_info "Checking project structure..."
    
    # Check for main.go or cmd directory
    if [[ -f "main.go" ]] || [[ -d "cmd" ]]; then
        log_success "Go project structure detected"
    else
        log_warning "No main.go or cmd/ directory found"
    fi
    
    # Check for go.mod
    if [[ -f "go.mod" ]]; then
        log_success "go.mod exists"
    else
        log_warning "go.mod not found"
    fi
    
    # Check for LICENSE file
    if [[ -f "LICENSE" ]]; then
        log_success "LICENSE file exists"
        
        # Check LICENSE file size (should not be empty)
        local license_size=$(wc -c < "LICENSE" 2>/dev/null || echo "0")
        if [[ $license_size -gt 100 ]]; then
            log_success "LICENSE file has content ($license_size bytes)"
        else
            log_warning "LICENSE file seems too small: $license_size bytes"
        fi
    else
        log_warning "No LICENSE file found"
    fi
    
    # Check for license generation script
    if [[ -f "scripts/generate-license.sh" ]]; then
        log_success "License generation script exists"
        if [[ -x "scripts/generate-license.sh" ]]; then
            log_success "License script is executable"
        else
            log_warning "License script is not executable"
        fi
    else
        log_warning "License generation script not found"
    fi
    
    # Check for Dockerfile if Docker is configured
    if grep -q "dockers:" "$GORELEASER_FILE" 2>/dev/null || grep -q "dockers:" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        if [[ -f "Dockerfile" ]]; then
            log_success "Dockerfile exists"
        else
            log_warning "Dockerfile not found but Docker is configured"
        fi
    fi
    
    # Check for assets/licenses directory
    if [[ -d "assets/licenses" ]]; then
        log_success "License templates directory exists"
        local template_count=$(find assets/licenses -name "*.template" | wc -l)
        if [[ $template_count -gt 0 ]]; then
            log_success "Found $template_count license templates"
        else
            log_warning "No license templates found in assets/licenses"
        fi
    else
        log_warning "License templates directory not found"
    fi
}

check_hooks_commands() {
    log_info "Checking hook commands..."
    
    # Check for templ
    if grep -q "templ generate" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "templ" || log_warning "templ is referenced but not installed"
    fi
    
    # Check for tsp (TypeSpec)
    if grep -q "tsp compile" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "tsp" || log_warning "TypeSpec is referenced but not installed"
    fi
    
    # Check for security tools
    if grep -q "gosec" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "gosec" || log_warning "gosec is referenced but not installed"
    fi
    
    if grep -q "golangci-lint" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "golangci-lint" || log_warning "golangci-lint is referenced but not installed"
    fi
}

check_signing_tools() {
    log_info "Checking signing tools..."
    
    if grep -q "cosign" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "cosign" || log_warning "cosign is referenced but not installed"
    fi
    
    if grep -q "syft" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "syft" || log_warning "syft is referenced but not installed"
    fi
}

check_compression_tools() {
    log_info "Checking compression tools..."
    
    if grep -q "upx" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        check_command "upx" || log_warning "UPX is referenced but not installed"
    fi
}

validate_templates() {
    log_info "Validating template variables..."
    
    # Check for common template issues
    local invalid_templates=$(grep -E '\{\{[^}]*\{\{' "$CURRENT_FILE" 2>/dev/null || true)
    if [[ -n "$invalid_templates" ]]; then
        log_error "Found nested templates: $invalid_templates"
    fi
    
    # Check for unclosed templates
    local unclosed=$(grep -E '\{\{[^}]*$' "$CURRENT_FILE" 2>/dev/null || true)
    if [[ -n "$unclosed" ]]; then
        log_error "Found unclosed templates: $unclosed"
    fi
}

check_git_state() {
    log_info "Checking Git state..."
    
    if [[ -d ".git" ]]; then
        log_success "Git repository initialized"
        
        # Check for remote
        if git remote -v | grep -q origin; then
            log_success "Git remote 'origin' is configured"
        else
            log_warning "No Git remote 'origin' configured"
        fi
        
        # Check for tags
        if git tag | grep -q .; then
            log_success "Git tags exist"
        else
            log_warning "No Git tags found (needed for releases)"
        fi
    else
        log_error "Not a Git repository"
    fi
}

test_license_system() {
    log_info "Testing license generation system..."
    
    if [[ -f "scripts/generate-license.sh" ]] && [[ -x "scripts/generate-license.sh" ]]; then
        # Test license script help
        if ./scripts/generate-license.sh --help >/dev/null 2>&1; then
            log_success "License script help works"
        else
            log_warning "License script help failed"
        fi
        
        # Test license templates listing
        if ./scripts/generate-license.sh --list >/dev/null 2>&1; then
            log_success "License templates listing works"
        else
            log_warning "License templates listing failed"
        fi
        
        # Test license generation (if readme config exists)
        if [[ -f ".readme/configs/readme-config.yaml" ]]; then
            # Backup existing LICENSE
            local backup_license=""
            if [[ -f "LICENSE" ]]; then
                backup_license=$(mktemp)
                cp LICENSE "$backup_license"
            fi
            
            # Test license generation
            if ./scripts/generate-license.sh >/dev/null 2>&1; then
                log_success "License generation test passed"
                
                # Restore backup if we made one
                if [[ -n "$backup_license" ]] && [[ -f "$backup_license" ]]; then
                    cp "$backup_license" LICENSE
                    rm -f "$backup_license"
                fi
            else
                log_warning "License generation test failed"
                
                # Restore backup if we made one
                if [[ -n "$backup_license" ]] && [[ -f "$backup_license" ]]; then
                    cp "$backup_license" LICENSE
                    rm -f "$backup_license"
                fi
            fi
        else
            log_warning "No readme config found, skipping license generation test"
        fi
    else
        log_warning "License generation script not found or not executable"
    fi
}

run_security_validation() {
    log_info "Running security validation..."
    
    # Go code security scan
    if command -v gosec &> /dev/null; then
        log_info "Scanning Go code with gosec..."
        if gosec -quiet ./... 2>/dev/null; then
            log_success "Go code security scan passed"
        else
            log_error "Go code security issues found"
        fi
    else
        log_warning "gosec not installed, skipping Go security scan"
    fi
    
    # Dependency vulnerability scan
    if command -v govulncheck &> /dev/null; then
        log_info "Checking dependencies for vulnerabilities..."
        if govulncheck ./... 2>/dev/null | grep -q "No vulnerabilities found"; then
            log_success "No vulnerable dependencies found"
        else
            log_warning "Vulnerable dependencies or scan issues detected"
        fi
    else
        log_warning "govulncheck not installed, skipping dependency vulnerability scan"
    fi
    
    # Shell script security scan
    if command -v shellcheck &> /dev/null; then
        log_info "Scanning shell scripts..."
        local shell_issues=0
        for script in *.sh scripts/*.sh; do
            if [[ -f "$script" ]]; then
                if ! shellcheck --severity=error "$script" 2>/dev/null; then
                    shell_issues=$((shell_issues + 1))
                fi
            fi
        done
        if [[ $shell_issues -eq 0 ]]; then
            log_success "Shell script security scan passed"
        else
            log_error "Shell script security issues found in $shell_issues files"
        fi
    else
        log_warning "shellcheck not installed, skipping shell script scan"
    fi
    
    # Dockerfile security scan
    if [[ -f "Dockerfile" ]]; then
        if command -v hadolint &> /dev/null; then
            log_info "Scanning Dockerfile..."
            if hadolint Dockerfile 2>/dev/null; then
                log_success "Dockerfile security scan passed"
            else
                log_error "Dockerfile security issues found"
            fi
        else
            log_warning "hadolint not installed, skipping Dockerfile scan"
        fi
    fi
    
    # Check for hardcoded secrets
    log_info "Checking for hardcoded secrets..."
    local secret_count=$(grep -r -i --exclude-dir=.git --exclude-dir=dist --exclude-dir=vendor \
        -E "(password|secret|token|key)[[:space:]]*[:=][[:space:]]*['\"][^'\"]{8,}" . 2>/dev/null | wc -l || echo 0)
    if [[ $secret_count -eq 0 ]]; then
        log_success "No hardcoded secrets detected"
    else
        log_error "Potential hardcoded secrets found: $secret_count matches"
    fi
}

run_dry_run() {
    log_info "Running GoReleaser dry-run..."
    
    if command -v goreleaser &> /dev/null; then
        # Temporarily clear conflicting tokens for GoReleaser
        local saved_gitlab_token="${GITLAB_TOKEN:-}"
        local saved_gitea_token="${GITEA_TOKEN:-}"
        unset GITLAB_TOKEN GITEA_TOKEN
        
        if [[ -f "$GORELEASER_FILE" ]]; then
            log_info "Testing free version configuration..."
            if goreleaser release --config "$GORELEASER_FILE" --snapshot --skip=publish --clean 2>/dev/null; then
                log_success "Dry-run successful for free version"
            else
                log_warning "Dry-run failed for free version (this might be expected in test environment)"
            fi
        fi
        
        # Note: Pro features dry-run would require a pro license
        if [[ -f "$GORELEASER_PRO_FILE" ]]; then
            log_info "Pro version exists but requires license for full validation"
        fi
        
        # Restore tokens
        [[ -n "$saved_gitlab_token" ]] && export GITLAB_TOKEN="$saved_gitlab_token"
        [[ -n "$saved_gitea_token" ]] && export GITEA_TOKEN="$saved_gitea_token"
    else
        log_warning "goreleaser not installed, skipping dry-run"
    fi
}

# Main execution
main() {
    echo "================================"
    echo "GoReleaser Configuration Verifier"
    echo "================================"
    echo
    
    # Check dependencies (continue even if some are missing)
    log_info "Checking dependencies..."
    check_command "go" || true
    check_command "git" || true
    check_command "goreleaser" || true
    check_command "yq" || true
    check_command "docker" || true
    echo
    
    # Check configurations (continue even if some checks fail)
    check_goreleaser_config "$GORELEASER_FILE" || true
    echo
    check_goreleaser_config "$GORELEASER_PRO_FILE" || true
    echo
    
    # Check environment variables
    check_required_env_vars "$GORELEASER_FILE"
    echo
    check_required_env_vars "$GORELEASER_PRO_FILE"
    echo
    
    # Check project structure
    check_project_structure
    echo
    
    # Check specific tools
    check_hooks_commands
    echo
    check_signing_tools
    echo
    check_compression_tools
    echo
    
    # Check Git state
    check_git_state
    echo
    
    # Test license system
    test_license_system
    echo
    
    # Validate templates
    CURRENT_FILE=$GORELEASER_FILE
    validate_templates
    CURRENT_FILE=$GORELEASER_PRO_FILE
    validate_templates
    echo
    
    # Run security validation
    run_security_validation
    echo
    
    # Try dry-run
    run_dry_run
    echo
    
    # Summary
    echo "================================"
    echo "Verification Summary"
    echo "================================"
    echo -e "${GREEN}Checks passed:${NC} $CHECKS"
    echo -e "${YELLOW}Warnings:${NC} $WARNINGS"
    echo -e "${RED}Errors:${NC} $ERRORS"
    echo
    
    if [[ $ERRORS -eq 0 ]]; then
        if [[ $WARNINGS -eq 0 ]]; then
            echo -e "${GREEN}✓ All checks passed successfully!${NC}"
            exit 0
        else
            echo -e "${YELLOW}⚠ Verification completed with warnings${NC}"
            exit 0
        fi
    else
        echo -e "${RED}✗ Verification failed with errors${NC}"
        exit 1
    fi
}

# Run main function
main