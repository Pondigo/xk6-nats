# Vendor Hash Management

This project uses Nix's `buildGoModule` which requires a `vendorHash` for reproducible builds. The vendor hash is a SHA256 hash of all Go module dependencies.

## Automatic Updates

### GitHub Actions

The project includes automated vendor hash updates:

1. **Update Vendor Hash Workflow** (`.github/workflows/update-vendor-hash.yaml`)
   - Triggers on changes to `go.mod` or `go.sum`
   - Automatically calculates the new hash
   - Creates a pull request with the updated hash

2. **Test Workflow Integration** (`.github/workflows/test.yaml`)
   - Detects hash mismatches during CI
   - Automatically updates the hash if needed
   - Continues with the updated hash

### Local Updates

#### Using the Script

```bash
# Update vendor hash automatically
./scripts/update-vendor-hash.sh
```

#### Using Just

```bash
# Update vendor hash using just
just vendor-hash
```

#### Manual Process

If you need to update the hash manually:

1. Set a dummy hash in `flake.nix`:
   ```nix
   vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
   ```

2. Try to build:
   ```bash
   nix build .#k6-pondigo-nats
   ```

3. Copy the actual hash from the error message (the `got:` value)

4. Update `flake.nix` with the correct hash

## How It Works

The vendor hash system works by:

1. **Hash Calculation**: Nix calculates a hash of all vendored dependencies
2. **Verification**: The hash is compared against the expected hash in `flake.nix`
3. **Reproducibility**: Ensures the same dependencies are used across builds
4. **Security**: Prevents tampering with dependencies

## Troubleshooting

### Hash Mismatch Error

If you see an error like:
```
error: hash mismatch in fixed-output derivation:
         specified: sha256-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
            got:    sha256-YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
```

Run the update script:
```bash
./scripts/update-vendor-hash.sh
```

### Build Fails After Hash Update

If the build fails after updating the hash:

1. Check if `go.mod` or `go.sum` have uncommitted changes
2. Run `go mod tidy` to clean up dependencies
3. Try updating the hash again

### CI Failures

If CI fails with hash errors:

1. The automated system should handle this automatically
2. If it doesn't, check the workflow logs
3. Manually trigger the "Update Vendor Hash" workflow

## Best Practices

1. **Commit `go.mod` and `go.sum` together** - Always commit both files when updating dependencies
2. **Update hash after dependency changes** - Run the update script after modifying Go dependencies
3. **Use the automated script** - Prefer the script over manual hash updates
4. **Check CI logs** - Monitor CI for hash-related issues