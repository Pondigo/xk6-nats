# Local GitHub Actions Testing

This project supports **Nix-native GitHub Actions testing** - a fast, Docker-free approach for running workflows locally.

## 🚀 Quick Start

```bash
# Enter development environment
nix develop

# Run test workflow natively
just test

# Run Nix workflows natively
just nix

# Run all workflows natively
just all
```

## 📋 Available Commands

| Command | Description |
|---------|-------------|
| `just test` | Run test workflow natively |
| `just build` | Run build workflow natively (requires tag) |
| `just nix` | Run Nix-specific workflows natively |
| `just update-hash` | Run vendor hash update natively |
| `just all` | Run all workflows natively |
| `just dry <workflow>` | Dry run workflow natively |
| `just verbose <workflow>` | Run workflow natively with verbose output |
| `just test-setup` | Test Nix-native setup |

### Direct Script Usage

```bash
# Run specific workflow
./scripts/run-nix-actions.sh test

# Dry run
./scripts/run-nix-actions.sh test --dry-run

# Verbose output
./scripts/run-nix-actions.sh test --verbose

# Run all workflows
./scripts/run-nix-actions.sh all

# Show help
./scripts/run-nix-actions.sh --help
```

## 🔧 Configuration

The Nix-native runner uses your existing Nix environment:

- **No Docker required** - Uses Nix store directly
- **Native performance** - No container overhead
- **Nix integration** - Uses flake.nix and devShell
- **Environment variables** - Inherits from Nix shell
- **Timeout protection** - Prevents hanging builds
- **Error handling** - Graceful failure handling

## 🔧 Workflow-Specific Notes

### Test Workflow

```bash
just test
# or
./scripts/run-nix-actions.sh test
```

- Runs Go tests, linting, and builds
- Starts NATS server for integration tests
- Optimized for local execution

### Nix Workflows

```bash
just nix
# or
./scripts/run-nix-actions.sh nix
```

- Runs `nix flake check` and builds
- Tests cross-platform compilation
- Validates Nix configuration

### Build Workflow

```bash
just build
# or
./scripts/run-nix-actions.sh build
```

- Requires a git tag to trigger
- Creates release artifacts
- Generates SBOM for security

### Vendor Hash Update

```bash
just update-hash
# or
./scripts/run-nix-actions.sh update-hash
```

- Updates vendor hash when Go dependencies change
- Automated dependency management
- Lightweight and fast

## 🐛 Troubleshooting

### Nix Issues

**Error**: `nix command not found`
```bash
# Install Nix
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install

# Or use existing installation
nix --version
```

**Error**: `flake.nix validation failed`
```bash
# Check flake syntax
nix flake check

# Show detailed errors
nix flake check --verbose
```

### Go Issues

**Error**: `go command not found`
```bash
# Enter development environment
nix develop

# Check Go availability
nix develop --command go version
```

**Error**: `Go tests failed`
```bash
# Run tests with verbose output
just verbose test

# Run specific test
nix develop --command go test -v ./path/to/test
```

### Build Issues

**Error**: `Build timeout`
```bash
# The script includes timeout protection
# Try building with longer timeout:
./scripts/run-nix-actions.sh build

# Or build manually:
nix build .#k6-pondigo-nats
```

### Permission Issues

**Error**: `permission denied`
```bash
# Check file permissions
ls -la scripts/

# Make scripts executable
chmod +x scripts/*.sh
```

## 🎯 Best Practices

1. **Test Before Pushing**: Always run workflows locally before pushing
2. **Use Dry Run**: Check what will run with `--dry-run`
3. **Start Small**: Test individual workflows before running all
4. **Monitor Resources**: Watch system resources during execution
5. **Clean Up**: Remove build artifacts periodically

## 📚 Additional Resources

- [Nix Documentation](https://nixos.org/manual/nix/stable/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Documentation](https://golang.org/doc/)

## 🆘 Getting Help

If you encounter issues:

1. Check the [Troubleshooting](#-troubleshooting) section
2. Run with `--verbose` flag for detailed output
3. Check GitHub Actions workflow files in `.github/workflows/`
4. Test your setup with `just test-setup`
5. Open an issue on the project repository