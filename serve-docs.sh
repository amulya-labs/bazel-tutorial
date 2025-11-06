#!/bin/bash
# Build and serve the documentation with the dashboard

set -e

echo "Building React dashboard..."
cd web_ui
npm run build
cd ..

echo "Copying dashboard to docs..."
rm -rf docs/app
cp -r web_ui/dist docs/app

echo "Building MkDocs site..."
source .venv/bin/activate
mkdocs build --clean

echo "Copying dashboard to site..."
cp -r docs/app site/

echo "Starting MkDocs server..."
echo "Dashboard will be available at: http://127.0.0.1:8000/bazel-tutorial/app/"
echo "Docs will be available at: http://127.0.0.1:8000/bazel-tutorial/"
mkdocs serve
