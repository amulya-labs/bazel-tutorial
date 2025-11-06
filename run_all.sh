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
echo "Step 1: Checking economic data..."
echo ""

# Use absolute path for database
DB_PATH="$(pwd)/data/econ.db"
export ECON_DB_PATH="$DB_PATH"

# Check if database exists and is less than 1 hour old
SKIP_REFRESH=false
if [ -f "$DB_PATH" ]; then
    # Get current time and file modification time in seconds since epoch
    CURRENT_TIME=$(date +%s)
    FILE_TIME=$(stat -c %Y "$DB_PATH" 2>/dev/null || stat -f %m "$DB_PATH" 2>/dev/null)
    AGE_SECONDS=$((CURRENT_TIME - FILE_TIME))
    AGE_MINUTES=$((AGE_SECONDS / 60))

    if [ $AGE_SECONDS -lt 3600 ]; then
        echo "✅ Using cached data (last updated $AGE_MINUTES minutes ago)"
        SKIP_REFRESH=true
    else
        echo "⏰ Data is older than 1 hour, refreshing..."
    fi
else
    echo "📥 No cached data found, fetching from FRED..."
fi

if [ "$SKIP_REFRESH" = false ]; then
    echo "Running: bazel run //go_fetch:refresh"
    if bazel run //go_fetch:refresh -- --db="$DB_PATH"; then
        echo "✅ Data refresh complete"
    else
        echo "⚠️  Data refresh failed or skipped"
    fi
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
