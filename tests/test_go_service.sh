#!/bin/bash
# Simple integration test to verify Go service builds correctly
# This is a shell test that can be run with: bazel test //tests:test_go_service

set -e

echo "Testing Go service build..."

# The binary is in the runfiles because it's declared as data dependency
# Find the binary in the runfiles
SERVER_PATH="go_service/server_/server"

if [ -f "$SERVER_PATH" ]; then
    echo "✓ Go service binary exists at $SERVER_PATH"
    exit 0
else
    echo "✗ Go service binary not found"
    echo "Looking for: $SERVER_PATH"
    echo "Available files:"
    find . -name "server" 2>/dev/null || true
    exit 1
fi
