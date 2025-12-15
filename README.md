# GoReleaser Wizard 🚀

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![GoReleaser](https://img.shields.io/badge/powered%20by-GoReleaser-blue.svg)](https://goreleaser.com)

**The interactive setup wizard that creates perfect GoReleaser configurations in seconds.**

Stop copy-pasting configs. Stop guessing at YAML. Get a production-ready GoReleaser setup with one command.

## ✨ Features

- 🎯 **Interactive wizard** - Guides you through every option
- 🧠 **Smart defaults** - Detects your project structure automatically
- 🚀 **GitHub Actions included** - Complete CI/CD pipeline ready to go
- 📦 **Multi-platform builds** - Linux, macOS, Windows, ARM, and more
- 🐳 **Docker support** - Multi-arch container images
- 🔒 **Security built-in** - Code signing, SBOM generation
- ✅ **Validation** - Check your config before releasing

## 🎬 Quick Start

```bash
# Install
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest

# Run the wizard
goreleaser-wizard init

# That's it! Your .goreleaser.yaml is ready
```

## 📸 Demo

```bash
$ goreleaser-wizard init
🚀 GoReleaser Configuration Wizard
Let's create the perfect GoReleaser config for your project!

? Project Name › my-awesome-cli
? Project Description › A fantastic CLI tool
? Project Type › CLI Application
? Binary Name › my-awesome-cli
? Main Package Path › ./cmd/my-awesome-cli

? Target Platforms › ✓ linux ✓ darwin ✓ windows
? Target Architectures › ✓ amd64 ✓ arm64
? Enable CGO? › No (recommended)
? Embed Version Info? › Yes (recommended)

? Git Provider › GitHub
? Docker Images? › Yes
? Code Signing? › Yes
? Generate SBOM? › Yes

✓ Created .goreleaser.yaml
✓ Created .github/workflows/release.yml

✨ Setup Complete!
```

## 🛠️ Installation

### Using Go

```bash
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest
```

### From Source

```bash
git clone https://github.com/LarsArtmann/GoReleaser-Wizard.git
cd GoReleaser-Wizard
go build -o goreleaser-wizard ./cmd/goreleaser-wizard
```

### Download Binary

Download the latest release from the [releases page](https://github.com/LarsArtmann/GoReleaser-Wizard/releases).

## 📖 Usage

### Interactive Mode (Recommended)

The wizard will guide you through creating a perfect configuration:

```bash
goreleaser-wizard init
```

Options:

- `--force` - Overwrite existing configuration
- `--minimal` - Create minimal configuration
- `--pro` - Include GoReleaser Pro features

### Non-Interactive Mode

Perfect for CI/CD pipelines:

```bash
goreleaser-wizard generate \
  --name my-project \
  --binary my-app \
  --platforms linux,darwin,windows \
  --docker \
  --github-action
```

### Validate Configuration

Check your existing GoReleaser configuration:

```bash
goreleaser-wizard validate

# With fixes
goreleaser-wizard validate --fix

# Verbose output
goreleaser-wizard validate --verbose
```

## 🎯 What It Creates

### `.goreleaser.yaml`

- Optimized build configuration
- Multi-platform support
- Archive generation
- Checksums and signatures
- Changelog generation
- Release configuration

### `.github/workflows/release.yml`

- Automated releases on tags
- Docker image building
- Code signing with cosign
- SBOM generation
- Multi-platform builds

## 🏗️ Project Types

The wizard adapts to your project:

- **CLI Application** - Single binary with version info
- **Web Service** - Includes Docker configuration
- **Library with CLI** - Focuses on the CLI component
- **Multiple Binaries** - Configures multiple build targets

## 🔧 Advanced Features

### GoReleaser Pro Support

Enable Pro features during setup:

```bash
goreleaser-wizard init --pro
```

Adds support for:

- Custom publishers
- Advanced templating
- Nightlies
- Docker manifests
- And more!

### Docker Integration

When Docker is enabled, the wizard:

- Detects your registry (ghcr.io, Docker Hub, etc.)
- Configures multi-platform images
- Sets up proper labels
- Handles authentication in CI/CD

### Package Managers

Optional support for:

- **Homebrew** - macOS/Linux formula
- **Snap** - Linux snap packages
- **Scoop** - Windows package manager
- **AUR** - Arch Linux (Pro)

## 🧪 Testing Your Configuration

After generating your configuration:

```bash
# 1. Validate the configuration
goreleaser-wizard validate

# 2. Test build locally
goreleaser build --snapshot --clean

# 3. Create a tag
git tag -a v0.1.0 -m 'First release'

# 4. Push to trigger release
git push origin v0.1.0
```

## 📚 Examples

### Minimal CLI Tool

```bash
goreleaser-wizard generate \
  --name simple-cli \
  --binary simple \
  --platforms linux,darwin
```

### Full-Featured Web Service

```bash
goreleaser-wizard generate \
  --name api-server \
  --binary server \
  --docker \
  --signing \
  --github-action \
  --platforms linux,darwin,windows \
  --architectures amd64,arm64
```

### Library with CLI

```bash
goreleaser-wizard init --minimal
# Then select "Library with CLI" in the wizard
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [GoReleaser](https://goreleaser.com) - The amazing release automation tool
- [Charm](https://charm.sh) - Beautiful terminal UI components
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## 🔗 Links

- [GoReleaser Documentation](https://goreleaser.com)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Report Issues](https://github.com/LarsArtmann/GoReleaser-Wizard/issues)

---

**Made with ❤️ to simplify Go releases**
