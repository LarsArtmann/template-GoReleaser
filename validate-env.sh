#!/usr/bin/env bash
set -euo pipefail

# Environment Variable Validation Script
# Validates all environment variables used in GoReleaser configurations

# Change to script directory
cd "$(dirname "${BASH_SOURCE[0]}")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Validation state
ERRORS=0
WARNINGS=0
CHECKS=0
CRITICAL_MISSING=()
OPTIONAL_MISSING=()

# Configuration files
GORELEASER_FILE=".goreleaser.yaml"
GORELEASER_PRO_FILE=".goreleaser.pro.yaml"
ENV_EXAMPLE_FILE=".env.example"

# Environment variable categories
declare -A CRITICAL_VARS
declare -A OPTIONAL_VARS
declare -A VAR_DESCRIPTIONS
declare -A VAR_VALIDATORS

# Critical environment variables (required for basic functionality)
CRITICAL_VARS=(
    ["GITHUB_TOKEN"]="GitHub API access token for releases"
    ["GITHUB_OWNER"]="GitHub repository owner/organization"
    ["GITHUB_REPO"]="GitHub repository name"
)

# Optional environment variables (enhance functionality)
OPTIONAL_VARS=(
    ["DOCKER_USERNAME"]="Docker Hub username for container publishing"
    ["DOCKER_TOKEN"]="Docker Hub access token"
    ["PROJECT_DESCRIPTION"]="Project description for packages"
    ["VENDOR_NAME"]="Vendor/company name for packages"
    ["MAINTAINER_NAME"]="Package maintainer name"
    ["MAINTAINER_EMAIL"]="Package maintainer email"
    ["HOMEBREW_TAP_GITHUB_TOKEN"]="GitHub token for Homebrew tap"
    ["SCOOP_GITHUB_TOKEN"]="GitHub token for Scoop bucket"
    ["AUR_KEY"]="AUR SSH private key path"
    ["FURY_TOKEN"]="Fury.io API token"
    ["FURY_ACCOUNT"]="Fury.io account name"
    ["S3_BUCKET"]="AWS S3 bucket for releases"
    ["AWS_ACCESS_KEY_ID"]="AWS access key ID"
    ["AWS_SECRET_ACCESS_KEY"]="AWS secret access key"
    ["AZURE_STORAGE_ACCOUNT"]="Azure storage account name"
    ["AZURE_STORAGE_CONTAINER"]="Azure storage container name"
    ["AZURE_STORAGE_KEY"]="Azure storage access key"
    ["GCS_BUCKET"]="Google Cloud Storage bucket"
    ["GOOGLE_APPLICATION_CREDENTIALS"]="GCS service account credentials path"
    ["ARTIFACTORY_HOST"]="Artifactory server hostname"
    ["ARTIFACTORY_REPO"]="Artifactory repository name"
    ["ARTIFACTORY_USERNAME"]="Artifactory username"
    ["ARTIFACTORY_PASSWORD"]="Artifactory password/token"
    ["DISCORD_WEBHOOK_ID"]="Discord webhook ID for notifications"
    ["DISCORD_WEBHOOK_TOKEN"]="Discord webhook token"
    ["SLACK_WEBHOOK_URL"]="Slack webhook URL for notifications"
    ["TEAMS_WEBHOOK_URL"]="Microsoft Teams webhook URL"
    ["SMTP_FROM"]="SMTP sender email address"
    ["SMTP_TO"]="SMTP recipient email address"
    ["SMTP_USERNAME"]="SMTP server username"
    ["SMTP_PASSWORD"]="SMTP server password"
    ["WEBHOOK_URL"]="Custom webhook URL"
    ["WEBHOOK_TOKEN"]="Custom webhook authentication token"
    ["COSIGN_PRIVATE_KEY"]="Cosign private key for signing"
    ["COSIGN_PASSWORD"]="Cosign private key password"
    ["GITLAB_TOKEN"]="GitLab API token (alternative to GitHub)"
    ["GITEA_TOKEN"]="Gitea API token (alternative to GitHub)"
)

# Validators for specific environment variables
VAR_VALIDATORS=(
    ["GITHUB_TOKEN"]="validate_github_token"
    ["DOCKER_TOKEN"]="validate_docker_token"
    ["MAINTAINER_EMAIL"]="validate_email"
    ["SMTP_FROM"]="validate_email"
    ["SMTP_TO"]="validate_email"
    ["SLACK_WEBHOOK_URL"]="validate_url"
    ["TEAMS_WEBHOOK_URL"]="validate_url"
    ["WEBHOOK_URL"]="validate_url"
    ["S3_BUCKET"]="validate_aws_bucket_name"
    ["GCS_BUCKET"]="validate_gcs_bucket_name"
    ["AZURE_STORAGE_ACCOUNT"]="validate_azure_storage_name"
    ["ARTIFACTORY_HOST"]="validate_hostname"
    ["AUR_KEY"]="validate_file_path"
    ["COSIGN_PRIVATE_KEY"]="validate_file_path"
    ["GOOGLE_APPLICATION_CREDENTIALS"]="validate_file_path"
)

# Logging functions
log_header() {
    echo
    echo -e "${MAGENTA}════════════════════════════════════════${NC}"
    echo -e "${MAGENTA}  $1${NC}"
    echo -e "${MAGENTA}════════════════════════════════════════${NC}"
}

log_subheader() {
    echo
    echo -e "${CYAN}──── $1 ────${NC}"
}

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

log_critical() {
    echo -e "${RED}[CRITICAL]${NC} $1"
    ((ERRORS++))
}

# Validation functions for specific variable types
validate_github_token() {
    local token="$1"
    if [[ "$token" =~ ^ghp_[A-Za-z0-9]{36}$ ]] || [[ "$token" =~ ^github_pat_[A-Za-z0-9_]{82}$ ]]; then
        return 0
    else
        echo "Invalid GitHub token format"
        return 1
    fi
}

validate_docker_token() {
    local token="$1"
    if [[ "$token" =~ ^dckr_pat_[A-Za-z0-9_-]{30,}$ ]]; then
        return 0
    else
        echo "Invalid Docker token format"
        return 1
    fi
}

validate_email() {
    local email="$1"
    if [[ "$email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]]; then
        return 0
    else
        echo "Invalid email format"
        return 1
    fi
}

validate_url() {
    local url="$1"
    if [[ "$url" =~ ^https?://[A-Za-z0-9.-]+.*$ ]]; then
        return 0
    else
        echo "Invalid URL format (must start with http:// or https://)"
        return 1
    fi
}

validate_aws_bucket_name() {
    local bucket="$1"
    if [[ "$bucket" =~ ^[a-z0-9][a-z0-9.-]*[a-z0-9]$ ]] && [[ ${#bucket} -ge 3 ]] && [[ ${#bucket} -le 63 ]]; then
        return 0
    else
        echo "Invalid S3 bucket name format"
        return 1
    fi
}

validate_gcs_bucket_name() {
    local bucket="$1"
    if [[ "$bucket" =~ ^[a-z0-9][a-z0-9._-]*[a-z0-9]$ ]] && [[ ${#bucket} -ge 3 ]] && [[ ${#bucket} -le 63 ]]; then
        return 0
    else
        echo "Invalid GCS bucket name format"
        return 1
    fi
}

validate_azure_storage_name() {
    local name="$1"
    if [[ "$name" =~ ^[a-z0-9]{3,24}$ ]]; then
        return 0
    else
        echo "Invalid Azure storage account name (must be 3-24 lowercase alphanumeric)"
        return 1
    fi
}

validate_hostname() {
    local hostname="$1"
    if [[ "$hostname" =~ ^[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]]; then
        return 0
    else
        echo "Invalid hostname format"
        return 1
    fi
}

validate_file_path() {
    local path="$1"
    # Expand tilde if present
    path="${path/#\~/$HOME}"
    
    if [[ -f "$path" ]]; then
        return 0
    else
        echo "File does not exist: $path"
        return 1
    fi
}

# Extract environment variables from GoReleaser configs
extract_env_vars_from_configs() {
    log_subheader "Analyzing GoReleaser configurations"
    
    local all_vars=()
    
    # Extract from all config files
    for config_file in "$GORELEASER_FILE" "$GORELEASER_PRO_FILE"; do
        if [[ -f "$config_file" ]]; then
            log_info "Scanning $config_file for environment variables..."
            local vars_in_file
            vars_in_file=$(grep -oE '\{\{ \.Env\.[A-Z_]+ \}\}' "$config_file" 2>/dev/null | sed 's/{{ \.Env\.\([A-Z_]*\) }}/\1/' | sort -u)
            
            if [[ -n "$vars_in_file" ]]; then
                log_info "Found variables in $config_file:"
                while IFS= read -r var; do
                    echo "  - $var"
                    all_vars+=("$var")
                done <<< "$vars_in_file"
            else
                log_info "No environment variables found in $config_file"
            fi
        fi
    done
    
    # Remove duplicates
    IFS=$'\n' all_vars=($(sort -u <<< "${all_vars[*]}"))
    unset IFS
    
    echo
    log_info "Total unique environment variables found: ${#all_vars[@]}"
}

# Check if .env.example is in sync with actual usage
validate_env_example_sync() {
    log_subheader "Validating .env.example synchronization"
    
    if [[ ! -f "$ENV_EXAMPLE_FILE" ]]; then
        log_error ".env.example file not found"
        return 1
    fi
    
    # Extract variables from .env.example
    local env_example_vars
    env_example_vars=$(grep -E '^[A-Z_]+=' "$ENV_EXAMPLE_FILE" 2>/dev/null | cut -d'=' -f1 | sort)
    
    # Extract variables from GoReleaser configs
    local config_vars
    config_vars=$(for f in "$GORELEASER_FILE" "$GORELEASER_PRO_FILE"; do
        if [[ -f "$f" ]]; then
            grep -oE '\{\{ \.Env\.[A-Z_]+ \}\}' "$f" 2>/dev/null | sed 's/{{ \.Env\.\([A-Z_]*\) }}/\1/'
        fi
    done | sort -u)
    
    # Check for variables in configs but not in .env.example
    log_info "Checking for missing variables in .env.example..."
    local missing_in_example=()
    while IFS= read -r var; do
        if ! echo "$env_example_vars" | grep -q "^$var$"; then
            missing_in_example+=("$var")
        fi
    done <<< "$config_vars"
    
    if [[ ${#missing_in_example[@]} -gt 0 ]]; then
        log_error "Variables used in configs but missing from .env.example:"
        printf '  - %s\n' "${missing_in_example[@]}"
    else
        log_success ".env.example contains all variables used in configs"
    fi
    
    # Check for variables in .env.example but not used in configs
    log_info "Checking for unused variables in .env.example..."
    local unused_in_example=()
    while IFS= read -r var; do
        if [[ -n "$var" ]] && ! echo "$config_vars" | grep -q "^$var$"; then
            unused_in_example+=("$var")
        fi
    done <<< "$env_example_vars"
    
    if [[ ${#unused_in_example[@]} -gt 0 ]]; then
        log_warning "Variables in .env.example but not used in configs:"
        printf '  - %s\n' "${unused_in_example[@]}"
    else
        log_success "All variables in .env.example are used in configs"
    fi
}

# Validate individual environment variable
validate_env_var() {
    local var_name="$1"
    local var_value="${!var_name:-}"
    local is_critical="$2"
    
    if [[ -z "$var_value" ]]; then
        if [[ "$is_critical" == "true" ]]; then
            log_error "$var_name: MISSING (CRITICAL)"
            CRITICAL_MISSING+=("$var_name")
        else
            log_warning "$var_name: not set (optional)"
            OPTIONAL_MISSING+=("$var_name")
        fi
        return 1
    fi
    
    # Check if variable has a custom validator
    if [[ -n "${VAR_VALIDATORS[$var_name]:-}" ]]; then
        local validator="${VAR_VALIDATORS[$var_name]}"
        local validation_error
        if validation_error=$($validator "$var_value" 2>&1); then
            log_success "$var_name: set and valid"
        else
            log_error "$var_name: set but invalid - $validation_error"
            return 1
        fi
    else
        # Basic validation - just check if it's not empty and not placeholder
        if [[ "$var_value" =~ ^(your-|xxxx|example|changeme|todo) ]] || [[ ${#var_value} -lt 3 ]]; then
            log_warning "$var_name: appears to be a placeholder value"
            return 1
        else
            log_success "$var_name: set"
        fi
    fi
    
    return 0
}

# Validate all environment variables
validate_all_env_vars() {
    log_header "ENVIRONMENT VARIABLE VALIDATION"
    
    log_subheader "Critical Environment Variables"
    log_info "These variables are required for basic GoReleaser functionality"
    
    for var_name in "${!CRITICAL_VARS[@]}"; do
        echo -e "${CYAN}$var_name${NC}: ${CRITICAL_VARS[$var_name]}"
        validate_env_var "$var_name" "true"
    done
    
    log_subheader "Optional Environment Variables"
    log_info "These variables enable additional features"
    
    for var_name in "${!OPTIONAL_VARS[@]}"; do
        echo -e "${CYAN}$var_name${NC}: ${OPTIONAL_VARS[$var_name]}"
        validate_env_var "$var_name" "false"
    done
}

# Test environment variable loading by GoReleaser
test_goreleaser_env_loading() {
    log_subheader "Testing GoReleaser Environment Variable Loading"
    
    if ! command -v goreleaser &> /dev/null; then
        log_warning "goreleaser not installed, skipping env loading test"
        return 0
    fi
    
    # Temporarily clear conflicting tokens for GoReleaser
    local saved_gitlab_token="${GITLAB_TOKEN:-}"
    local saved_gitea_token="${GITEA_TOKEN:-}"
    unset GITLAB_TOKEN GITEA_TOKEN
    
    for config_file in "$GORELEASER_FILE" "$GORELEASER_PRO_FILE"; do
        if [[ -f "$config_file" ]]; then
            log_info "Testing environment variable loading for $config_file..."
            
            # Use goreleaser's check command which will fail if required env vars are missing
            if goreleaser check --config "$config_file" &> /dev/null; then
                log_success "Environment variables properly loaded by $config_file"
            else
                log_warning "GoReleaser check failed for $config_file (this might be expected in test environment)"
                # Show specific error but filter out multiple token warnings
                echo -e "${YELLOW}Error details:${NC}"
                goreleaser check --config "$config_file" 2>&1 | head -5 | grep -v "multiple tokens" | sed 's/^/  /' || true
            fi
        fi
    done
    
    # Restore tokens
    [[ -n "$saved_gitlab_token" ]] && export GITLAB_TOKEN="$saved_gitlab_token"
    [[ -n "$saved_gitea_token" ]] && export GITEA_TOKEN="$saved_gitea_token"
}

# Generate environment variable documentation
generate_env_documentation() {
    log_subheader "Generating environment variable documentation"
    
    local doc_file="ENV_VARS.md"
    
    cat > "$doc_file" << 'EOF'
# Environment Variables

This document lists all environment variables used by the GoReleaser configurations.

## Critical Environment Variables

These variables are required for basic functionality:

EOF
    
    for var_name in "${!CRITICAL_VARS[@]}"; do
        echo "### $var_name" >> "$doc_file"
        echo "${CRITICAL_VARS[$var_name]}" >> "$doc_file"
        echo >> "$doc_file"
        
        # Add validator info if available
        if [[ -n "${VAR_VALIDATORS[$var_name]:-}" ]]; then
            case "${VAR_VALIDATORS[$var_name]}" in
                "validate_github_token")
                    echo "**Format:** GitHub Personal Access Token (ghp_...) or Fine-grained token (github_pat_...)" >> "$doc_file"
                    ;;
                "validate_docker_token")
                    echo "**Format:** Docker Personal Access Token (dckr_pat_...)" >> "$doc_file"
                    ;;
                "validate_email")
                    echo "**Format:** Valid email address" >> "$doc_file"
                    ;;
                "validate_url")
                    echo "**Format:** Valid HTTP/HTTPS URL" >> "$doc_file"
                    ;;
                "validate_file_path")
                    echo "**Format:** Path to existing file" >> "$doc_file"
                    ;;
            esac
            echo >> "$doc_file"
        fi
    done
    
    echo "## Optional Environment Variables" >> "$doc_file"
    echo >> "$doc_file"
    echo "These variables enable additional features:" >> "$doc_file"
    echo >> "$doc_file"
    
    for var_name in "${!OPTIONAL_VARS[@]}"; do
        echo "### $var_name" >> "$doc_file"
        echo "${OPTIONAL_VARS[$var_name]}" >> "$doc_file"
        echo >> "$doc_file"
    done
    
    log_success "Environment variable documentation generated: $doc_file"
}

# Show helpful setup instructions
show_setup_instructions() {
    if [[ ${#CRITICAL_MISSING[@]} -gt 0 ]] || [[ ${#OPTIONAL_MISSING[@]} -gt 0 ]]; then
        log_header "SETUP INSTRUCTIONS"
        
        if [[ ${#CRITICAL_MISSING[@]} -gt 0 ]]; then
            echo -e "${RED}Critical variables that must be set:${NC}"
            for var in "${CRITICAL_MISSING[@]}"; do
                echo -e "  ${RED}✗${NC} $var: ${CRITICAL_VARS[$var]}"
            done
            echo
        fi
        
        if [[ ${#OPTIONAL_MISSING[@]} -gt 0 ]]; then
            echo -e "${YELLOW}Optional variables for additional features:${NC}"
            for var in "${OPTIONAL_MISSING[@]}"; do
                echo -e "  ${YELLOW}⚠${NC} $var: ${OPTIONAL_VARS[$var]}"
            done
            echo
        fi
        
        echo -e "${BLUE}To set up your environment:${NC}"
        echo "1. Copy .env.example to .env"
        echo "   cp .env.example .env"
        echo
        echo "2. Edit .env with your actual values"
        echo "   editor .env"
        echo
        echo "3. Source the environment file before running GoReleaser"
        echo "   source .env"
        echo "   # or export the variables directly"
        echo
        echo "4. For GitHub Actions, add secrets to your repository:"
        echo "   GitHub Settings → Secrets and variables → Actions"
    fi
}

# Generate validation report
generate_validation_report() {
    local report_file="env-validation-report.json"
    
    log_info "Generating validation report: $report_file"
    
    cat > "$report_file" << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "total_checks": $CHECKS,
  "errors": $ERRORS,
  "warnings": $WARNINGS,
  "critical_missing": [
EOF
    
    if [[ ${#CRITICAL_MISSING[@]} -gt 0 ]]; then
        local first=true
        for var in "${CRITICAL_MISSING[@]}"; do
            if [[ "$first" == false ]]; then
                echo "," >> "$report_file"
            fi
            echo -n "    \"$var\"" >> "$report_file"
            first=false
        done
        echo >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF
  ],
  "optional_missing": [
EOF
    
    if [[ ${#OPTIONAL_MISSING[@]} -gt 0 ]]; then
        local first=true
        for var in "${OPTIONAL_MISSING[@]}"; do
            if [[ "$first" == false ]]; then
                echo "," >> "$report_file"
            fi
            echo -n "    \"$var\"" >> "$report_file"
            first=false
        done
        echo >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF
  ],
  "validation_status": "$(if [[ ${#CRITICAL_MISSING[@]} -eq 0 ]]; then echo "READY"; else echo "NEEDS_SETUP"; fi)"
}
EOF
    
    log_success "Validation report saved to $report_file"
}

# Main execution
main() {
    echo -e "${MAGENTA}╔════════════════════════════════════════╗${NC}"
    echo -e "${MAGENTA}║     ENVIRONMENT VARIABLE VALIDATOR    ║${NC}"
    echo -e "${MAGENTA}║   GoReleaser Configuration Analyzer   ║${NC}"
    echo -e "${MAGENTA}╚════════════════════════════════════════╝${NC}"
    
    # Extract environment variables from configs
    extract_env_vars_from_configs
    
    # Validate .env.example sync
    validate_env_example_sync
    
    # Validate all environment variables
    validate_all_env_vars
    
    # Test GoReleaser environment loading
    test_goreleaser_env_loading
    
    # Generate documentation
    generate_env_documentation
    
    # Show setup instructions
    show_setup_instructions
    
    # Generate validation report
    generate_validation_report
    
    # Final summary
    log_header "VALIDATION SUMMARY"
    echo -e "${GREEN}Checks completed:${NC} $CHECKS"
    echo -e "${YELLOW}Warnings:${NC} $WARNINGS"
    echo -e "${RED}Errors:${NC} $ERRORS"
    echo -e "${RED}Critical missing:${NC} ${#CRITICAL_MISSING[@]}"
    echo -e "${YELLOW}Optional missing:${NC} ${#OPTIONAL_MISSING[@]}"
    
    if [[ ${#CRITICAL_MISSING[@]} -eq 0 ]]; then
        echo
        echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║  ✓ READY FOR GORELEASER EXECUTION     ║${NC}"
        echo -e "${GREEN}║  All critical variables are set       ║${NC}"
        echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
        exit 0
    else
        echo
        echo -e "${RED}╔════════════════════════════════════════╗${NC}"
        echo -e "${RED}║  ✗ SETUP REQUIRED                     ║${NC}"
        echo -e "${RED}║  Missing critical environment vars    ║${NC}"
        echo -e "${RED}╚════════════════════════════════════════╝${NC}"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Environment Variable Validator for GoReleaser"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h    Show this help message"
        echo "  --list, -l    List all environment variables"
        echo "  --docs, -d    Generate documentation only"
        echo "  --test, -t    Test GoReleaser loading only"
        echo
        echo "This script validates all environment variables used in GoReleaser"
        echo "configurations and provides setup instructions for missing ones."
        exit 0
        ;;
    --list|-l)
        echo "Critical environment variables:"
        for var in "${!CRITICAL_VARS[@]}"; do
            echo "  $var: ${CRITICAL_VARS[$var]}"
        done
        echo
        echo "Optional environment variables:"
        for var in "${!OPTIONAL_VARS[@]}"; do
            echo "  $var: ${OPTIONAL_VARS[$var]}"
        done
        exit 0
        ;;
    --docs|-d)
        generate_env_documentation
        exit 0
        ;;
    --test|-t)
        test_goreleaser_env_loading
        exit $?
        ;;
esac

# Run full validation
main "$@"