#!/bin/bash
# Simple integration test to verify Go service builds correctly
# This is a shell test that can be run with: bazel test //tests:test_go_service

set -e

echo "Testing Go service build..."

# Check if the Go binary was built successfully
if bazel build //go_service:server; then
    echo "✓ Go service builds successfully"
    exit 0
else
    echo "✗ Go service failed to build"
    exit 1
fi
