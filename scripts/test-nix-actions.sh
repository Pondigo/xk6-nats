#!/usr/bin/env bash
set -euo pipefail

# Test script for Nix-native GitHub Actions runner
# Usage: ./scripts/test-nix-actions.sh

echo "🧪 Testing Nix-native GitHub Actions setup..."

# Test 1: Check if script exists and is executable
echo "1. Testing script availability..."
if [[ -x "scripts/run-nix-actions.sh" ]]; then
    echo "✅ run-nix-actions.sh script exists and is executable"
else
    echo "❌ run-nix-actions.sh script missing or not executable"
    exit 1
fi

# Test 2: Check if script shows help
echo "2. Testing help functionality..."
if ./scripts/run-nix-actions.sh --help > /dev/null 2>&1; then
    echo "✅ Help functionality works"
else
    echo "❌ Help functionality failed"
    exit 1
fi

# Test 3: Check if dry run works
echo "3. Testing dry run functionality..."
if ./scripts/run-nix-actions.sh test --dry-run > /dev/null 2>&1; then
    echo "✅ Dry run functionality works"
else
    echo "❌ Dry run functionality failed"
    exit 1
fi

# Test 4: Check if Nix environment is ready
echo "4. Testing Nix environment..."
if nix --version > /dev/null 2>&1; then
    echo "✅ Nix is available"
else
    echo "❌ Nix is not available"
    exit 1
fi

# Test 5: Check if flake.nix is valid
echo "5. Testing flake.nix validity..."
if nix flake check --quiet > /dev/null 2>&1; then
    echo "✅ flake.nix is valid"
else
    echo "❌ flake.nix is invalid"
    exit 1
fi

# Test 6: Check if Go is available in dev shell
echo "6. Testing Go availability..."
if nix develop --command go version > /dev/null 2>&1; then
    echo "✅ Go is available in dev shell"
else
    echo "❌ Go is not available in dev shell"
    exit 1
fi

# Test 7: Check if just commands are available
echo "7. Testing just commands..."
if grep -q "nix-" justfile; then
    nix_count=$(grep -c "nix-" justfile)
    echo "✅ Nix commands are available in justfile ($nix_count commands found)"
else
    echo "❌ Nix commands not found in justfile"
    exit 1
fi

# Test 8: Try a minimal dry run of test workflow
echo "8. Testing minimal workflow execution..."
if timeout 10s ./scripts/run-nix-actions.sh test --dry-run > /dev/null 2>&1; then
    echo "✅ Workflow execution test passed"
else
    echo "❌ Workflow execution test failed"
    exit 1
fi

echo ""
echo "🎉 All tests passed! Your Nix-native GitHub Actions setup is ready to use."
echo ""
echo "📋 Quick start commands:"
echo "   just test                      # Run test workflow natively"
echo "   just nix                       # Run Nix workflow natively"
echo "   just dry test                  # Dry run test workflow"
echo "   just verbose test               # Run test workflow with verbose output"
echo "   just all                       # Run all workflows natively"
echo ""
echo "✨ Benefits:"
echo "   • No Docker required - Uses Nix store directly"
echo "   • Fast execution - No container overhead"
echo "   • Native integration - Works with your Nix environment"
echo ""
echo "📚 For more info, see: docs/LOCAL_ACTIONS.md"