# Bazel Quick Start Guide

This is a hands-on guide to get you started with this Bazel tutorial.

## Prerequisites

- **Bazel** installed (see [Installation](#installation))
- **Git** for cloning the repository

## Installation

### Option 1: Using Bazelisk (Recommended)

Bazelisk automatically downloads and uses the correct Bazel version:

**macOS:**
```bash
brew install bazelisk
```

**Linux:**
```bash
npm install -g @bazel/bazelisk
# or
wget https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 -O /usr/local/bin/bazel
chmod +x /usr/local/bin/bazel
```

**Windows:**
```bash
choco install bazelisk
```

### Option 2: Direct Bazel Installation

Follow the official guide: https://bazel.build/install

### Verify Installation

```bash
bazel --version
# Should show: bazel 6.4.0 or later
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/rrl-personal-projects/bazel-tutorial.git
cd bazel-tutorial
```

### 2. Build Everything

The simplest command - builds all targets:

```bash
bazel build //...
```

**What happens:**
- Bazel reads `MODULE.bazel` to resolve dependencies
- Downloads Python and Go toolchains
- Compiles all services
- Caches results for future builds

### 3. Run the Go Service

```bash
bazel run //go_service:server
```

Visit `http://localhost:8080` - you should see:
```
Hello from Go! 🚀
```

Press `Ctrl+C` to stop.

### 4. Run the Python Service

```bash
bazel run //py_service:server
```

Visit `http://localhost:8081` - you should see:
```json
{"message": "Hello from Python! 🐍"}
```

Press `Ctrl+C` to stop.

### 5. Run All Tests

```bash
bazel test //...
```

**Output:**
```
//go_service:go_service_test        PASSED in 0.5s
//py_service:py_service_test        PASSED in 1.2s
//tests:test_go_service             PASSED in 0.3s
//tests:test_py_service             PASSED in 0.4s
```

## Common Commands

### Build Commands

```bash
# Build all targets
bazel build //...

# Build specific target
bazel build //go_service:server

# Build multiple packages
bazel build //go_service/... //py_service/...
```

### Run Commands

```bash
# Run a binary target
bazel run //go_service:server

# Run with arguments
bazel run //go_service:server -- --port=9090
```

### Test Commands

```bash
# Test all targets
bazel test //...

# Test specific package
bazel test //go_service:go_service_test

# Test with verbose output
bazel test //... --test_output=all

# Test and show only failures
bazel test //... --test_output=errors
```

### Query Commands

```bash
# List all targets
bazel query //...

# List targets in a package
bazel query //py_service:all

# Show dependencies
bazel query 'deps(//go_service:server)'

# Find reverse dependencies
bazel query 'rdeps(//..., //py_service:server)'
```

### Clean Commands

```bash
# Clean build outputs (keeps cache)
bazel clean

# Clean everything including cache
bazel clean --expunge
```

## Understanding the Output

### Build Output

```
INFO: Analyzed 2 targets (15 packages loaded, 89 targets configured).
INFO: Found 2 targets...
INFO: Elapsed time: 3.456s, Critical Path: 2.123s
```

- **Analyzed**: How many targets Bazel processed
- **Packages loaded**: How many BUILD files were read
- **Critical Path**: Longest chain of dependent actions

### Cached Builds

After the first build:
```
INFO: Analyzed 2 targets (0 packages loaded, 0 targets configured).
INFO: Found 2 targets...
INFO: Elapsed time: 0.234s, Critical Path: 0.001s
```

Everything is cached! ⚡

## Incremental Builds Demo

Try this to see Bazel's incremental build in action:

```bash
# 1. Build everything
bazel build //...

# 2. Change only Python code
echo "# New comment" >> py_service/main.py

# 3. Build again
bazel build //...
# Only Python targets rebuild - Go is cached!

# 4. Change Go code
echo "// New comment" >> go_service/main.go

# 5. Build again
bazel build //...
# Only Go targets rebuild - Python is cached!
```

## Troubleshooting

### "Cannot find MODULE.bazel"

Make sure you're in the repository root:
```bash
cd bazel-tutorial
ls MODULE.bazel  # Should exist
```

### "command not found: bazel"

Bazel isn't installed or not in PATH. See [Installation](#installation).

### "Failed to fetch external repository"

Network issue or typo in MODULE.bazel. Check:
```bash
bazel sync
```

### Build is slow

First build is always slow (downloads dependencies). Subsequent builds are fast due to caching.

### Clean slate

Start fresh:
```bash
bazel clean --expunge
bazel build //...
```

## Next Steps

1. **Read the concepts**: Check out [Bazel Concepts](concepts.md)
2. **Modify code**: Change `main.py` and see incremental builds
3. **Add a dependency**: Update `requirements.txt` and rebuild
4. **Create a test**: Add a new test file and `py_test` target
5. **Explore queries**: Use `bazel query` to understand the build graph

## Resources

- [Bazel Documentation](https://bazel.build/docs)
- [rules_python](https://github.com/bazelbuild/rules_python)
- [rules_go](https://github.com/bazelbuild/rules_go)
- [Bazel Glossary](https://bazel.build/reference/glossary)
- [Bazel Best Practices](https://bazel.build/basics/best-practices)
