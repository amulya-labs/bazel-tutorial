#!/bin/bash
# validate.sh - Script to validate the Bazel tutorial setup

set -e

echo "🧪 Validating Bazel Tutorial Setup..."
echo ""

# Check if bazel is installed
echo "1️⃣  Checking Bazel installation..."
if command -v bazel &> /dev/null; then
    echo "✅ Bazel is installed: $(bazel version 2>&1 | head -1 || echo "Version check failed")"
else
    echo "❌ Bazel is not installed. Please install Bazel: https://bazel.build/install"
    exit 1
fi

# Check if we're in the right directory
echo ""
echo "2️⃣  Checking repository structure..."
if [ ! -f "MODULE.bazel" ]; then
    echo "❌ MODULE.bazel not found. Please run this script from the repository root."
    exit 1
fi

if [ ! -d "go_service" ] || [ ! -d "py_service" ]; then
    echo "❌ Service directories not found."
    exit 1
fi

echo "✅ Repository structure looks good"

# List all targets
echo ""
echo "3️⃣  Listing all targets..."
if bazel query //... 2>&1 | grep -q "//go_service:server"; then
    echo "✅ Targets are properly defined"
else
    echo "⚠️  Could not query targets (this might be okay on first run)"
fi

# Try to build all targets
echo ""
echo "4️⃣  Building all targets..."
if bazel build //... 2>&1 | tee /tmp/bazel_build.log; then
    echo "✅ All targets built successfully!"
else
    echo "❌ Build failed. Check /tmp/bazel_build.log for details"
    exit 1
fi

# Try to run tests
echo ""
echo "5️⃣  Running all tests..."
if bazel test //... --test_output=errors 2>&1 | tee /tmp/bazel_test.log; then
    echo "✅ All tests passed!"
else
    echo "❌ Tests failed. Check /tmp/bazel_test.log for details"
    exit 1
fi

echo ""
echo "🎉 Success! Your Bazel tutorial is fully functional."
echo ""
echo "Try running the services:"
echo "  bazel run //go_service:server"
echo "  bazel run //py_service:server"
echo ""
echo "Explore the documentation:"
echo "  docs/QUICKSTART.md"
echo "  docs/CONCEPTS.md"
echo "  docs/DEPENDENCIES.md"
