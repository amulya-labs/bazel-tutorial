# Makefile for Economic Dashboard
# Convenience wrapper around Bazel commands

.PHONY: help setup build test clean fetch api ui all

# Default target
help:
	@echo "Economic Dashboard - Makefile Commands"
	@echo "========================================"
	@echo ""
	@echo "Setup:"
	@echo "  make setup       - Install all dependencies"
	@echo ""
	@echo "Development:"
	@echo "  make fetch       - Fetch latest economic data"
	@echo "  make api         - Start Python API server"
	@echo "  make ui          - Start React dev server"
	@echo "  make all         - Run all services (fetch + api)"
	@echo ""
	@echo "Build & Test:"
	@echo "  make build       - Build all services"
	@echo "  make test        - Run all tests"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean       - Clean build artifacts"
	@echo ""

# Setup dependencies
setup:
	@echo "📦 Installing dependencies..."
	@cd web_ui && npm install
	@echo "✅ Setup complete!"

# Build all targets
build:
	@echo "🔨 Building all services..."
	@bazel build //...

# Run all tests
test:
	@echo "🧪 Running all tests..."
	@bazel test //...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@bazel clean
	@rm -rf web_ui/node_modules web_ui/dist

# Fetch economic data
fetch:
	@echo "📊 Fetching economic data..."
	@if [ -z "$$FRED_API_KEY" ]; then \
		echo "⚠️  Warning: FRED_API_KEY not set"; \
		echo "Get your free key at: https://fred.stlouisfed.org/docs/api/api_key.html"; \
	fi
	@bazel run //go_fetch:refresh

# Start API server
api:
	@echo "🚀 Starting Python API server..."
	@echo "API: http://localhost:8000"
	@echo "Docs: http://localhost:8000/docs"
	@bazel run //py_api:server

# Start UI dev server
ui:
	@echo "🎨 Starting React dev server..."
	@echo "UI: http://localhost:3000"
	@cd web_ui && npm run dev

# Run all services
all:
	@./run_all.sh
