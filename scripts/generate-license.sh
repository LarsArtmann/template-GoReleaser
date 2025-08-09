#!/usr/bin/env bash
set -euo pipefail

# License Generation Script
# Automatically generates LICENSE file from templates based on readme-config.yaml

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
CONFIG_FILE=".readme/configs/readme-config.yaml"
TEMPLATE_DIR="assets/licenses"
OUTPUT_FILE="LICENSE"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are available
check_dependencies() {
    local missing_deps=()
    
    if ! command -v yq &> /dev/null; then
        missing_deps+=("yq")
    fi
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        log_info "Install with: brew install yq"
        exit 1
    fi
}

# Extract license type from config
get_license_type() {
    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_error "Configuration file not found: $CONFIG_FILE"
        exit 1
    fi
    
    local license_type
    license_type=$(yq eval '.project.license.type' "$CONFIG_FILE" 2>/dev/null || echo "null")
    
    if [[ "$license_type" == "null" || "$license_type" == "" ]]; then
        log_error "License type not found in configuration"
        exit 1
    fi
    
    echo "$license_type"
}

# Get copyright holder from config
get_copyright_holder() {
    local copyright_holder
    copyright_holder=$(yq eval '.project.contact.name' "$CONFIG_FILE" 2>/dev/null || echo "null")
    
    if [[ "$copyright_holder" == "null" || "$copyright_holder" == "" ]]; then
        log_warning "Copyright holder not found in config, using default"
        echo "Project Maintainer"
    else
        echo "$copyright_holder"
    fi
}

# Get current year
get_current_year() {
    date +%Y
}

# Find license template file
find_license_template() {
    local license_type="$1"
    local template_file="$TEMPLATE_DIR/${license_type}.template"
    
    if [[ -f "$template_file" ]]; then
        echo "$template_file"
        return 0
    fi
    
    # Try common variations
    local variations=(
        "$TEMPLATE_DIR/${license_type,,}.template"  # lowercase
        "$TEMPLATE_DIR/${license_type^^}.template"  # uppercase
    )
    
    for variation in "${variations[@]}"; do
        if [[ -f "$variation" ]]; then
            echo "$variation"
            return 0
        fi
    done
    
    return 1
}

# Process template and generate license
generate_license() {
    local license_type="$1"
    local copyright_holder="$2"
    local year="$3"
    local template_file
    
    log_info "Generating license: $license_type"
    log_info "Copyright holder: $copyright_holder"
    log_info "Year: $year"
    
    if ! template_file=$(find_license_template "$license_type"); then
        log_error "License template not found for: $license_type"
        log_info "Available templates:"
        ls -1 "$TEMPLATE_DIR"/*.template 2>/dev/null | sed 's/.*\//  - /' | sed 's/\.template$//' || log_warning "No templates found"
        exit 1
    fi
    
    log_info "Using template: $template_file"
    
    # Process template with variable substitution
    sed -e "s/{{YEAR}}/$year/g" \
        -e "s/{{COPYRIGHT_HOLDER}}/$copyright_holder/g" \
        "$template_file" > "$OUTPUT_FILE"
    
    log_success "License generated: $OUTPUT_FILE"
    
    # Show preview
    echo
    echo "License preview (first 10 lines):"
    head -10 "$OUTPUT_FILE" | sed 's/^/  /'
    echo
}

# Validate generated license
validate_license() {
    if [[ ! -f "$OUTPUT_FILE" ]]; then
        log_error "Generated license file not found: $OUTPUT_FILE"
        return 1
    fi
    
    local file_size
    file_size=$(wc -c < "$OUTPUT_FILE")
    
    if [[ $file_size -lt 100 ]]; then
        log_error "Generated license file seems too small: $file_size bytes"
        return 1
    fi
    
    # Check for template variables that weren't substituted
    if grep -q "{{" "$OUTPUT_FILE"; then
        log_error "Unsubstituted template variables found in license:"
        grep -n "{{" "$OUTPUT_FILE" | sed 's/^/  /'
        return 1
    fi
    
    log_success "License validation passed"
    return 0
}

# List available license templates
list_templates() {
    log_info "Available license templates:"
    if ls "$TEMPLATE_DIR"/*.template &>/dev/null; then
        ls -1 "$TEMPLATE_DIR"/*.template | sed 's/.*\//  - /' | sed 's/\.template$//'
    else
        log_warning "No license templates found in $TEMPLATE_DIR"
    fi
}

# Main function
main() {
    local license_type copyright_holder year
    
    log_info "License Generator Starting..."
    
    # Check dependencies
    check_dependencies
    
    # Handle command line options
    case "${1:-}" in
        "--list"|"-l")
            list_templates
            exit 0
            ;;
        "--help"|"-h")
            echo "Usage: $0 [--list|--help]"
            echo "  --list, -l    List available license templates"
            echo "  --help, -h    Show this help message"
            echo
            echo "Generates LICENSE file based on configuration in $CONFIG_FILE"
            exit 0
            ;;
    esac
    
    # Get configuration values
    license_type=$(get_license_type)
    copyright_holder=$(get_copyright_holder)
    year=$(get_current_year)
    
    # Generate license
    generate_license "$license_type" "$copyright_holder" "$year"
    
    # Validate result
    if validate_license; then
        log_success "License generation completed successfully!"
        
        # Show file info
        echo
        echo "Generated license details:"
        echo "  File: $OUTPUT_FILE"
        echo "  Size: $(wc -c < "$OUTPUT_FILE") bytes"
        echo "  Lines: $(wc -l < "$OUTPUT_FILE") lines"
        echo "  Type: $license_type"
    else
        log_error "License generation failed validation"
        exit 1
    fi
}

# Run main function
main "$@"