#!/bin/bash
# Development helper script for the Economic Dashboard

set -e

echo "🔧 Economic Dashboard Development Helper"
echo "========================================"
echo ""

show_help() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  setup       - Install all dependencies"
    echo "  fetch       - Fetch latest economic data"
    echo "  api         - Start Python API server"
    echo "  ui          - Start React dev server"
    echo "  build       - Build all services"
    echo "  test        - Run all tests"
    echo "  clean       - Clean build artifacts"
    echo "  help        - Show this help message"
    echo ""
}

setup() {
    echo "📦 Installing dependencies..."
    echo ""
    
    # Install Node dependencies for UI
    if [ -d "web_ui" ]; then
        echo "Installing Node.js dependencies..."
        cd web_ui
        npm install
        cd ..
        echo "✅ Node.js dependencies installed"
    fi
    
    echo ""
    echo "✅ Setup complete!"
    echo ""
    echo "Next steps:"
    echo "1. Set your FRED API key: export FRED_API_KEY=your_key"
    echo "2. Fetch data: $0 fetch"
    echo "3. Start API: $0 api (in one terminal)"
    echo "4. Start UI: $0 ui (in another terminal)"
}

fetch() {
    echo "📊 Fetching economic data from FRED..."
    
    if [ -z "$FRED_API_KEY" ]; then
        echo "⚠️  Warning: FRED_API_KEY not set"
        echo "Get your free key at: https://fred.stlouisfed.org/docs/api/api_key.html"
        echo ""
    fi
    
    bazel run //go_fetch:refresh
}

api() {
    echo "🚀 Starting Python API server..."
    echo "API will be available at http://localhost:8000"
    echo "API docs at http://localhost:8000/docs"
    echo ""
    bazel run //py_api:server
}

ui() {
    echo "🎨 Starting React dev server..."
    echo "UI will be available at http://localhost:3000"
    echo ""
    
    if [ ! -d "web_ui/node_modules" ]; then
        echo "Node modules not found. Running setup first..."
        setup
    fi
    
    cd web_ui
    npm run dev
}

build() {
    echo "🔨 Building all services..."
    bazel build //...
    echo "✅ Build complete!"
}

test() {
    echo "🧪 Running all tests..."
    bazel test //...
    echo "✅ Tests complete!"
}

clean() {
    echo "🧹 Cleaning build artifacts..."
    bazel clean
    rm -rf web_ui/node_modules web_ui/dist
    echo "✅ Clean complete!"
}

# Main command router
case "${1:-help}" in
    setup)
        setup
        ;;
    fetch)
        fetch
        ;;
    api)
        api
        ;;
    ui)
        ui
        ;;
    build)
        build
        ;;
    test)
        test
        ;;
    clean)
        clean
        ;;
    help|*)
        show_help
        ;;
esac
