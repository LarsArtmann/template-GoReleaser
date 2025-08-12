# API Documentation

This document provides comprehensive API documentation for the GoReleaser CLI template project, including CLI commands, service interfaces, and web APIs.

## Table of Contents

- [CLI Commands](#cli-commands)
- [Service Interfaces](#service-interfaces)
- [Web API Endpoints](#web-api-endpoints)
- [Configuration API](#configuration-api)
- [Data Types](#data-types)
- [Error Handling](#error-handling)

## CLI Commands

### Root Command

```
goreleaser-cli - A powerful CLI tool built with GoReleaser template

USAGE:
  goreleaser-cli [command]

FLAGS:
  -h, --help              help for goreleaser-cli
  -t, --toggle            help message for toggle
      --config string     config file (default is $HOME/.goreleaser-cli.yaml)

AVAILABLE COMMANDS:
  config      Manage configuration settings
  help        Help about any command
  license     Manage project licenses
  server      Start the web server for GoReleaser configuration management
  validate    Validate configuration and project structure
  verify      Verify project setup and perform comprehensive checks
  version     Print version information
```

### Version Command

**Usage:** `goreleaser-cli version`

**Description:** Print detailed version information including build metadata.

**Output:**
```
Version:      v1.0.0
Commit:       abc123def
Built:        2025-08-12T10:30:00Z
Built by:     goreleaser
Go version:   go1.21.0
OS/Arch:      linux/amd64
Git describe: v1.0.0-1-gabc123d
Git state:    clean
```

### Validate Command

**Usage:** `goreleaser-cli validate [flags]`

**Description:** Validate various aspects of your project including configuration files, project structure, environment variables, and dependencies.

**Flags:**
- `-c, --config`: Validate configuration files
- `-s, --structure`: Validate project structure
- `-e, --env`: Validate environment variables
- `-a, --all`: Validate all aspects (default: true)

**Examples:**
```bash
# Validate all aspects (default)
goreleaser-cli validate

# Validate only environment variables
goreleaser-cli validate --env

# Validate configuration and structure
goreleaser-cli validate --config --structure
```

**Output:**
```
🔍 Running validation...

📁 Validating project structure...
✅ Project structure validation passed (15 checks)

🌍 Validating environment...
✅ Environment validation passed (8 checks)

🛠️ Validating tools...
✅ Tools validation passed (5 checks)

📊 Validation Summary:
   ✅ Passed: 3
   ❌ Failed: 0

🎉 All validations passed successfully!
```

### Configuration Command

**Usage:** `goreleaser-cli config <subcommand> [flags]`

**Description:** Manage configuration settings for the CLI tool.

#### config show

**Usage:** `goreleaser-cli config show`

**Description:** Display the current configuration values and their sources.

**Output:**
```
📋 Current Configuration:

🔧 Configuration values:
  📄 License: MIT
  👤 Author: John Doe <john@example.com>
  📦 Project: my-project
  📋 Description: A sample project description
  🛠️ CLI: verbose=false, colors=true

✅ Configuration is valid

🌍 Environment variable overrides:
  (no relevant environment variables set)
```

#### config set

**Usage:** `goreleaser-cli config set <key> <value>`

**Description:** Set a configuration value in the config file.

**Supported Keys:**
- `license.type`: License type (e.g., MIT, Apache-2.0)
- `author.name`: Author's full name
- `author.email`: Author's email address
- `project.name`: Project name
- `project.description`: Project description
- `cli.verbose`: Enable verbose output (true/false)
- `cli.colors`: Enable colored output (true/false)

**Examples:**
```bash
# Set license type
goreleaser-cli config set license.type MIT

# Set author information
goreleaser-cli config set author.name "John Doe"
goreleaser-cli config set author.email "john@example.com"

# Update project details
goreleaser-cli config set project.name "my-awesome-project"
```

#### config init

**Usage:** `goreleaser-cli config init [flags]`

**Description:** Create a default configuration file with common settings.

**Flags:**
- `-f, --force`: Overwrite existing config file

**Examples:**
```bash
# Initialize new configuration
goreleaser-cli config init

# Force overwrite existing configuration
goreleaser-cli config init --force
```

### License Command

**Usage:** `goreleaser-cli license <subcommand> [flags]`

**Description:** Manage project licenses with various operations.

#### license generate

**Usage:** `goreleaser-cli license generate [license-type] [copyright-holder] [flags]`

**Description:** Generate a license file from available templates.

**Flags:**
- `-i, --interactive`: Interactive license generation
- `-t, --type string`: License type
- `-c, --holder string`: Copyright holder
- `-o, --output string`: Output file path (default: "LICENSE")

**Examples:**
```bash
# Generate MIT license
goreleaser-cli license generate MIT "John Doe"

# Interactive generation
goreleaser-cli license generate --interactive

# Generate to custom file
goreleaser-cli license generate Apache-2.0 "Acme Corp" --output LICENSE.txt
```

#### license list

**Usage:** `goreleaser-cli license list`

**Description:** List all available license templates.

**Output:**
```
📋 Available license templates:
  • MIT
  • Apache-2.0
  • BSD-3-Clause
  • EUPL-1.2

💡 Generate a license with: goreleaser-cli license generate <type> "Your Name"
```

#### license show

**Usage:** `goreleaser-cli license show [flags]`

**Description:** Display information about the current project license.

**Flags:**
- `-p, --preview`: Show preview of license content

**Output:**
```
🔍 Current license information:
  📄 File: LICENSE
  📏 Size: 1065 bytes
  📅 Modified: 2025-08-12 10:30:00

💡 Use --preview to see license content
```

#### license validate

**Usage:** `goreleaser-cli license validate`

**Description:** Validate that the existing license file is properly formatted.

**Output:**
```
🔍 Validating license file...
   • Checking file size...
   ✅ File size is reasonable
   • Checking for unsubstituted placeholders...
   ✅ No template placeholders found
   • Checking license format...
   ✅ License keywords found

📊 Validation Summary:
   ✅ Passed: 3
   ❌ Failed: 0

🎉 License validation passed!
```

### Server Command

**Usage:** `goreleaser-cli server [flags]`

**Description:** Start a web server that provides a modern web interface for managing and validating GoReleaser configurations.

**Flags:**
- `-p, --port int`: Port to run the server on (default: 8080)
- `--host string`: Host to bind the server to (default: "localhost")
- `--debug`: Enable debug mode with verbose logging
- `--static-dir string`: Directory for static files (default: "./web/static")

**Examples:**
```bash
# Start server on default port 8080
goreleaser-cli server

# Start on custom port with debug enabled
goreleaser-cli server --port 3000 --debug

# Bind to all interfaces
goreleaser-cli server --host 0.0.0.0 --port 8080
```

**Output:**
```
🚀 Starting GoReleaser Configuration Server
   Host: localhost
   Port: 8080
   Debug: false
   Static Dir: ./web/static

🌐 Server starting at http://localhost:8080
📋 Configuration Editor: http://localhost:8080/config
🏥 Health Check: http://localhost:8080/health
📊 Metrics: http://localhost:8080/metrics

Press Ctrl+C to stop the server
```

### Verify Command

**Usage:** `goreleaser-cli verify [flags]`

**Description:** Perform comprehensive project verification including security scans and dry runs.

**Flags:**
- `--skip-security`: Skip security scan
- `--skip-dry-run`: Skip GoReleaser dry run
- `--config string`: GoReleaser config file path

## Service Interfaces

The application uses dependency injection with well-defined service interfaces.

### ConfigService

Handles configuration management operations.

```go
type ConfigService interface {
    // LoadConfig loads configuration from file or creates default
    LoadConfig() (*types.Config, error)
    
    // SaveConfig saves configuration to file
    SaveConfig(config *types.Config) error
    
    // ValidateConfig validates configuration structure and values
    ValidateConfig(config *types.Config) error
    
    // InitConfig creates a new configuration with defaults
    InitConfig() (*types.Config, error)
}
```

**Methods:**

#### LoadConfig()
- **Returns:** `(*types.Config, error)`
- **Description:** Loads configuration from the default location or creates a new one if it doesn't exist.

#### SaveConfig(config *types.Config)
- **Parameters:** `config` - Configuration object to save
- **Returns:** `error`
- **Description:** Saves the provided configuration to the default file location.

#### ValidateConfig(config *types.Config)
- **Parameters:** `config` - Configuration object to validate
- **Returns:** `error`
- **Description:** Validates the configuration structure and values.

#### InitConfig()
- **Returns:** `(*types.Config, error)`
- **Description:** Creates a new configuration with default values.

### ValidationService

Handles all validation operations for the project.

```go
type ValidationService interface {
    // ValidateProject validates the entire project structure
    ValidateProject() (*ValidationResult, error)
    
    // ValidateEnvironment validates environment variables
    ValidateEnvironment() (*ValidationResult, error)
    
    // ValidateGoReleaser validates GoReleaser configuration files
    ValidateGoReleaser(configPath string) (*ValidationResult, error)
    
    // ValidateTools validates required tools are installed
    ValidateTools() (*ValidationResult, error)
}
```

**Methods:**

#### ValidateProject()
- **Returns:** `(*ValidationResult, error)`
- **Description:** Validates the entire project structure including required files, directories, and configurations.

#### ValidateEnvironment()
- **Returns:** `(*ValidationResult, error)`
- **Description:** Validates that all required environment variables are set and have valid values.

#### ValidateGoReleaser(configPath string)
- **Parameters:** `configPath` - Path to GoReleaser configuration file
- **Returns:** `(*ValidationResult, error)`
- **Description:** Validates GoReleaser configuration syntax and completeness.

#### ValidateTools()
- **Returns:** `(*ValidationResult, error)`
- **Description:** Validates that all required development tools are installed and accessible.

### LicenseService

Manages license generation and validation.

```go
type LicenseService interface {
    // GenerateLicense generates a license file based on configuration
    GenerateLicense(licenseType string, author string) error
    
    // ListAvailableLicenses returns available license templates
    ListAvailableLicenses() ([]LicenseTemplate, error)
    
    // ValidateLicense validates existing license file
    ValidateLicense() (*ValidationResult, error)
}
```

### VerificationService

Provides comprehensive project verification.

```go
type VerificationService interface {
    // RunFullVerification runs all verification checks
    RunFullVerification(opts *VerificationOptions) (*VerificationResult, error)
    
    // RunSecurityScan performs security validation
    RunSecurityScan() (*SecurityScanResult, error)
    
    // RunDryRun performs GoReleaser dry run
    RunDryRun(configPath string) (*DryRunResult, error)
}
```

## Web API Endpoints

When running in server mode, the application exposes several HTTP endpoints.

### Health Check

**Endpoint:** `GET /health`

**Description:** Returns the health status of the application.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-08-12T10:30:00Z",
  "version": "1.0.0",
  "uptime": "1h30m15s"
}
```

### Configuration Endpoints

#### Get Configuration

**Endpoint:** `GET /api/config`

**Description:** Retrieve current configuration.

**Response:**
```json
{
  "project": {
    "name": "my-project",
    "description": "Project description"
  },
  "author": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "license": {
    "type": "MIT"
  },
  "cli": {
    "verbose": false,
    "colors": true
  }
}
```

#### Update Configuration

**Endpoint:** `PUT /api/config`

**Description:** Update configuration settings.

**Request Body:**
```json
{
  "project": {
    "name": "updated-project",
    "description": "Updated description"
  },
  "author": {
    "name": "Jane Doe",
    "email": "jane@example.com"
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Configuration updated successfully"
}
```

### Validation Endpoints

#### Validate Project

**Endpoint:** `POST /api/validate`

**Description:** Run project validation checks.

**Request Body:**
```json
{
  "checks": ["project", "environment", "tools"],
  "options": {
    "verbose": true
  }
}
```

**Response:**
```json
{
  "success": true,
  "results": [
    {
      "name": "project",
      "success": true,
      "checks": 15,
      "errors": [],
      "warnings": []
    },
    {
      "name": "environment", 
      "success": true,
      "checks": 8,
      "errors": [],
      "warnings": ["OPTIONAL_VAR not set"]
    }
  ]
}
```

### License Endpoints

#### List Licenses

**Endpoint:** `GET /api/licenses`

**Description:** Get available license templates.

**Response:**
```json
{
  "licenses": [
    {
      "name": "MIT",
      "description": "MIT License",
      "path": "assets/licenses/MIT.template"
    },
    {
      "name": "Apache-2.0",
      "description": "Apache License 2.0",
      "path": "assets/licenses/Apache-2.0.template"
    }
  ]
}
```

#### Generate License

**Endpoint:** `POST /api/licenses/generate`

**Description:** Generate a license file.

**Request Body:**
```json
{
  "type": "MIT",
  "holder": "John Doe",
  "output": "LICENSE"
}
```

**Response:**
```json
{
  "success": true,
  "message": "License generated successfully",
  "file": "LICENSE"
}
```

### Metrics Endpoint

**Endpoint:** `GET /metrics`

**Description:** Prometheus-compatible metrics endpoint.

**Response:**
```
# HELP goreleaser_cli_requests_total Total number of requests
# TYPE goreleaser_cli_requests_total counter
goreleaser_cli_requests_total{method="GET",endpoint="/health"} 42

# HELP goreleaser_cli_request_duration_seconds Request duration in seconds
# TYPE goreleaser_cli_request_duration_seconds histogram
goreleaser_cli_request_duration_seconds_bucket{le="0.1"} 35
goreleaser_cli_request_duration_seconds_bucket{le="0.5"} 40
```

## Configuration API

### Configuration Structure

The application uses a hierarchical configuration structure defined in YAML:

```yaml
# Project information
project:
  name: "my-project"
  description: "Project description"
  version: "1.0.0"
  
# Author information  
author:
  name: "John Doe"
  email: "john@example.com"
  
# License configuration
license:
  type: "MIT"
  
# CLI behavior
cli:
  verbose: false
  colors: true
  
# Server configuration
server:
  host: "localhost"
  port: 8080
  debug: false
  
# Development settings
development:
  hot_reload: true
  debug_templates: false
```

### Environment Variable Overrides

Configuration values can be overridden using environment variables:

| Environment Variable | Configuration Path | Example |
|---------------------|-------------------|---------|
| `PROJECT_NAME` | `project.name` | `"my-app"` |
| `AUTHOR_NAME` | `author.name` | `"John Doe"` |
| `AUTHOR_EMAIL` | `author.email` | `"john@example.com"` |
| `LICENSE_TYPE` | `license.type` | `"Apache-2.0"` |
| `CLI_VERBOSE` | `cli.verbose` | `"true"` |
| `CLI_COLORS` | `cli.colors` | `"false"` |
| `SERVER_PORT` | `server.port` | `"3000"` |
| `SERVER_HOST` | `server.host` | `"0.0.0.0"` |

## Data Types

### Core Configuration Types

```go
// Config represents the main configuration structure
type Config struct {
    Project Project `yaml:"project" json:"project"`
    Author  Author  `yaml:"author" json:"author"`
    License License `yaml:"license" json:"license"`
    CLI     CLI     `yaml:"cli" json:"cli"`
    Server  Server  `yaml:"server" json:"server"`
}

// Project contains project-specific information
type Project struct {
    Name        string `yaml:"name" json:"name"`
    Description string `yaml:"description" json:"description"`
    Version     string `yaml:"version" json:"version"`
}

// Author contains author information
type Author struct {
    Name  string `yaml:"name" json:"name"`
    Email string `yaml:"email" json:"email"`
}

// License contains license configuration
type License struct {
    Type string `yaml:"type" json:"type"`
}

// CLI contains CLI behavior settings
type CLI struct {
    Verbose bool `yaml:"verbose" json:"verbose"`
    Colors  bool `yaml:"colors" json:"colors"`
}
```

### Validation Types

```go
// ValidationResult represents the result of a validation operation
type ValidationResult struct {
    Success  bool     `json:"success"`
    Errors   []string `json:"errors,omitempty"`
    Warnings []string `json:"warnings,omitempty"`
    Checks   int      `json:"checks"`
}

// VerificationOptions configures verification behavior
type VerificationOptions struct {
    SkipSecurity    bool   `json:"skip_security"`
    SkipDryRun      bool   `json:"skip_dry_run"`
    SkipLicenseTest bool   `json:"skip_license_test"`
    ConfigFile      string `json:"config_file"`
    ProConfigFile   string `json:"pro_config_file"`
}

// VerificationResult contains comprehensive verification results
type VerificationResult struct {
    Success  bool                `json:"success"`
    Checks   int                 `json:"checks"`
    Warnings int                 `json:"warnings"`
    Errors   int                 `json:"errors"`
    Details  map[string][]string `json:"details"`
}
```

### License Types

```go
// LicenseTemplate represents an available license template
type LicenseTemplate struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Path        string `json:"path"`
}
```

## Error Handling

### Error Response Format

All API endpoints follow a consistent error response format:

```json
{
  "error": true,
  "message": "Human-readable error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "specific error details",
    "validation_errors": ["list", "of", "errors"]
  },
  "timestamp": "2025-08-12T10:30:00Z"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `CONFIG_NOT_FOUND` | 404 | Configuration file not found |
| `CONFIG_INVALID` | 400 | Invalid configuration format |
| `VALIDATION_FAILED` | 422 | Validation checks failed |
| `LICENSE_NOT_FOUND` | 404 | License template not found |
| `PERMISSION_DENIED` | 403 | Insufficient permissions |
| `INTERNAL_ERROR` | 500 | Internal server error |

### CLI Error Handling

CLI commands use consistent exit codes:

- `0`: Success
- `1`: General error
- `2`: Misuse of command (invalid arguments)
- `3`: Configuration error
- `4`: Validation failure
- `5`: Network/connectivity error

### Logging

The application provides structured logging with configurable levels:

- `ERROR`: Error conditions
- `WARN`: Warning conditions  
- `INFO`: Informational messages
- `DEBUG`: Debug-level messages

**Log Format:**
```
2025-08-12T10:30:00Z [INFO] goreleaser-cli: Configuration loaded successfully file=.goreleaser-cli.yaml
2025-08-12T10:30:01Z [DEBUG] goreleaser-cli: Validating project structure checks=15
2025-08-12T10:30:02Z [ERROR] goreleaser-cli: Validation failed error="missing LICENSE file"
```

---

For more detailed information about specific endpoints or data structures, refer to the source code or generate API documentation using tools like `swagger` or `godoc`.