# Examples

This directory contains example configurations and usage patterns for the GoReleaser template project.

## Basic CLI Usage

The included `goreleaser-cli` is a simple example CLI application that demonstrates version information injection by GoReleaser.

```bash
# Build the CLI
go build -o goreleaser-cli ./cmd/goreleaser-cli

# Show version information
./goreleaser-cli version

# Show help
./goreleaser-cli --help
```

## Using This Template

1. **Fork or clone this repository**
2. **Rename the project** - Update module name in `go.mod`
3. **Customize the CLI** - Add your own commands to `cmd/goreleaser-cli/main.go`
4. **Update GoReleaser config** - Modify `.goreleaser.yaml` for your needs
5. **Create a release** - Tag your code and push to trigger the GitHub Actions workflow

## GoReleaser Examples

### Building a snapshot (without releasing)
```bash
goreleaser build --snapshot --clean
```

### Creating a release locally
```bash
# Requires a git tag
git tag -a v1.0.0 -m "Release v1.0.0"
goreleaser release --clean
```

### Validating configuration
```bash
goreleaser check
```

## GitHub Actions Integration

The template includes a complete GitHub Actions workflow that:
- Triggers on version tags (v*)
- Builds binaries for multiple platforms
- Creates GitHub releases with artifacts
- Supports both free and pro GoReleaser versions

## Customization Tips

### Adding New Commands

Edit `cmd/goreleaser-cli/main.go` to add your own commands using Cobra:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description of my command",
    Run: func(cmd *cobra.Command, args []string) {
        // Your command logic here
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

### Modifying Build Targets

Edit `.goreleaser.yaml` to customize which platforms to build for:

```yaml
builds:
  - goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
```

## License

This template is provided under the MIT License. See LICENSE file for details.