#!/bin/bash
# Run all components of the Economic Dashboard
#
# This script starts all services needed for the dashboard:
# 1. Refreshes economic data from FRED
# 2. Starts the Python API server
# 3. Provides instructions for starting the web UI

set -e

echo "========================================"
echo "Economic Indicators Dashboard - Startup"
echo "========================================"
echo ""

# Check for FRED API key
if [ -z "$FRED_API_KEY" ]; then
    echo "⚠️  Warning: FRED_API_KEY environment variable not set"
    echo "Get your free API key at: https://fred.stlouisfed.org/docs/api/api_key.html"
    echo ""
    echo "Set it with: export FRED_API_KEY=your_key_here"
    echo ""
    read -p "Press Enter to continue anyway (using cached data) or Ctrl+C to exit..."
fi

# Step 1: Refresh data
echo "Step 1: Refreshing economic data..."
echo "Running: bazel run //go_fetch:refresh"
echo ""

if bazel run //go_fetch:refresh; then
    echo "✅ Data refresh complete"
else
    echo "⚠️  Data refresh failed or skipped"
fi

echo ""
echo "========================================"
echo ""

# Step 2: Start API server (in background)
echo "Step 2: Starting Python API server..."
echo "Running: bazel run //py_api:server"
echo ""

echo "The API server will start on http://localhost:8000"
echo "API docs available at: http://localhost:8000/docs"
echo ""
echo "To start the web UI (in another terminal):"
echo "  cd web_ui"
echo "  npm install  # (first time only)"
echo "  npm run dev"
echo ""
echo "Then open http://localhost:3000 in your browser"
echo ""
echo "========================================"
echo ""

# Run the API server (this will block)
exec bazel run //py_api:server
