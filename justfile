# Default recipe
default: go-build

# Build the k6 extension
go-build:
    go build ./...

# Run tests
go-test:
    go test ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Run linter
lint:
    golangci-lint run

# Format code
fmt:
    go fmt ./...

# Tidy dependencies
tidy:
    go mod tidy

# Start NATS server for testing
nats-server:
    nats-server

# Clean build artifacts
clean:
    go clean -cache
    rm -f k6-pondigo-nats

# Build and install locally
install: go-build
    go install ./...

# Run all checks (lint, test, build)
check: lint test go-build

# Generate vendor hash for nix
vendor-hash:
    ./scripts/update-vendor-hash.sh

# Generate API documentation
docs:
    cd api-docs && yarn install && yarn generate-docs

# Watch for changes and rebuild
watch:
    @echo "Watching for changes..."
    @while true; do \
        find . -name "*.go" -type f | entr -rd just go-build; \
    done

# Run GitHub Actions locally using Nix-native runner
workflow name="":
    ./scripts/run-nix-actions.sh {{name}}

# Run test workflow locally
test:
    ./scripts/run-nix-actions.sh test

# Run build workflow locally (requires tag)
workflow-build:
    ./scripts/run-nix-actions.sh build

# Run nix workflow locally
nix:
    ./scripts/run-nix-actions.sh nix

# Run vendor hash update workflow locally
update-hash:
    ./scripts/run-nix-actions.sh update-hash

# Run all workflows locally
all:
    ./scripts/run-nix-actions.sh all

# Dry run workflows (show what would run)
dry workflow name="":
    ./scripts/run-nix-actions.sh {{name}} --dry-run

# Verbose workflow execution
verbose workflow name="":
    ./scripts/run-nix-actions.sh {{name}} --verbose

# Test GitHub Actions setup
test-setup:
    ./scripts/test-nix-actions.sh

# Show help
help:
    @just --list