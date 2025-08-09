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
    
    # Check with goreleaser
    if command -v goreleaser &> /dev/null; then
        if goreleaser check --config "$file" 2>/dev/null; then
            log_success "GoReleaser validation passed: $file"
        else
            log_error "GoReleaser validation failed: $file"
            goreleaser check --config "$file" 2>&1 | head -20
        fi
    else
        log_warning "goreleaser not installed, skipping native validation"
    fi
}

check_required_env_vars() {
    local file=$1
    log_info "Checking environment variables in $file..."
    
    # Extract env var references
    local env_vars=$(grep -oE '\{\{\.Env\.[A-Z_]+\}\}' "$file" 2>/dev/null | sed 's/{{\.Env\.\([^}]*\)}}/\1/' | sort -u)
    
    if [[ -n "$env_vars" ]]; then
        echo "Required environment variables:"
        for var in $env_vars; do
            if [[ -n "${!var:-}" ]]; then
                log_success "$var is set"
            else
                log_warning "$var is not set"
            fi
        done
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
    
    # Check for Dockerfile if Docker is configured
    if grep -q "dockers:" "$GORELEASER_FILE" 2>/dev/null || grep -q "dockers:" "$GORELEASER_PRO_FILE" 2>/dev/null; then
        if [[ -f "Dockerfile" ]]; then
            log_success "Dockerfile exists"
        else
            log_warning "Dockerfile not found but Docker is configured"
        fi
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

run_dry_run() {
    log_info "Running GoReleaser dry-run..."
    
    if command -v goreleaser &> /dev/null; then
        if [[ -f "$GORELEASER_FILE" ]]; then
            log_info "Testing free version configuration..."
            if goreleaser release --config "$GORELEASER_FILE" --snapshot --skip=publish --clean 2>/dev/null; then
                log_success "Dry-run successful for free version"
            else
                log_warning "Dry-run failed for free version (this might be expected without a Go project)"
            fi
        fi
        
        # Note: Pro features dry-run would require a pro license
        if [[ -f "$GORELEASER_PRO_FILE" ]]; then
            log_info "Pro version exists but requires license for full validation"
        fi
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
    
    # Check dependencies
    log_info "Checking dependencies..."
    check_command "go"
    check_command "git"
    check_command "goreleaser"
    check_command "yq"
    check_command "docker"
    echo
    
    # Check configurations
    check_goreleaser_config "$GORELEASER_FILE"
    echo
    check_goreleaser_config "$GORELEASER_PRO_FILE"
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
    
    # Validate templates
    CURRENT_FILE=$GORELEASER_FILE
    validate_templates
    CURRENT_FILE=$GORELEASER_PRO_FILE
    validate_templates
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