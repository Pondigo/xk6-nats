#!/usr/bin/env bash
set -euo pipefail

# Script to update vendor hash in flake.nix
# Usage: ./scripts/update-vendor-hash.sh

echo "🔍 Updating vendor hash for xk6-nats..."

# Get current vendor hash from flake.nix
current_hash=$(grep 'vendorHash = ' flake.nix | sed 's/.*"\(.*\)".*/\1/')
echo "Current vendor hash: $current_hash"

# Set a dummy hash to force Nix to calculate the real one
dummy_hash="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
sed -i "s|vendorHash = \".*\"|vendorHash = \"$dummy_hash\"|g" flake.nix

# Try to build and capture the actual hash
echo "🧮 Calculating new vendor hash..."
build_output=$(nix build .#k6-pondigo-nats 2>&1 || true)

# Extract the correct hash from the error message
new_hash=$(echo "$build_output" | grep 'got:' | sed 's/.*got:[[:space:]]*//' | sed 's/[[:space:]]*$//' || echo "")

if [ -z "$new_hash" ]; then
    echo "❌ Could not extract new vendor hash. Build output:"
    echo "$build_output"
    # Restore the original hash
    sed -i "s|vendorHash = \".*\"|vendorHash = \"$current_hash\"|g" flake.nix
    exit 1
fi

echo "✅ New vendor hash: $new_hash"

# Update flake.nix with the correct hash
sed -i "s|vendorHash = \".*\"|vendorHash = \"$new_hash\"|g" flake.nix

# Quick verification that the hash is correct (just check that it starts with sha256-)
if [[ $new_hash == sha256-* ]]; then
    echo "✅ Hash format looks correct!"
    echo "📝 Updated vendorHash from $current_hash to $new_hash"
else
    echo "❌ Invalid hash format. Restoring original hash."
    sed -i "s|vendorHash = \".*\"|vendorHash = \"$current_hash\"|g" flake.nix
    exit 1
fi

echo "🎉 Vendor hash updated successfully!"