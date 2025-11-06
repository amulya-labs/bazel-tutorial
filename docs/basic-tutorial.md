# Basic Services Tutorial

This tutorial introduces Bazel fundamentals through two simple services: a Python FastAPI service and a Go HTTP service.

## Overview

You'll learn how to:

- Build and test apps in Python and Go
- Use Bazel's MODULE.bazel and BUILD.bazel files
- Leverage caching, parallelism, and hermetic builds
- Run Bazel inside GitHub Actions CI

## Services

This project contains two simple services:

| Service | Language | Description | Run Command |
|---------|----------|-------------|-------------|
| py_service | Python (FastAPI) | "Hello from Python" web service | `bazel run //py_service:server` |
| go_service | Go | "Hello from Go" web service | `bazel run //go_service:server` |

## Repository Structure

```
bazel-tutorial/
├── MODULE.bazel            # Defines external dependencies using Bzlmod
├── .bazelversion           # Specifies the Bazel version to use
├── .bazelrc                # Bazel configuration options
├── go_service/             # Go service + BUILD.bazel
│   ├── main.go            → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/go_service/main.go
│   ├── main_test.go       → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/go_service/main_test.go
│   └── BUILD.bazel        → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/go_service/BUILD.bazel
├── py_service/             # Python service + BUILD.bazel
│   ├── main.py            → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/py_service/main.py
│   ├── test_main.py       → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/py_service/test_main.py
│   ├── requirements.txt   → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/py_service/requirements.txt
│   └── BUILD.bazel        → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/py_service/BUILD.bazel
├── tests/                  # Cross-service tests
│   ├── test_go_service.sh
│   ├── test_py_service.py
│   └── BUILD.bazel        → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/tests/BUILD.bazel
└── .github/workflows/ci.yaml  → View: https://github.com/rrl-personal-projects/bazel-tutorial/blob/main/.github/workflows/ci.yaml
```

## Quick Start

### 1. Build Everything

```bash
bazel build //...
```

This builds all targets defined in every BUILD.bazel.

**What happens:**
- Bazel reads `MODULE.bazel` to resolve dependencies
- Downloads Python and Go toolchains
- Compiles all services
- Caches results for future builds

### 2. Run the Python Service

```bash
bazel run //py_service:server
```

Visit `http://localhost:8081` - you should see:
```json
{"message": "Hello from Python! 🐍"}
```

Press `Ctrl+C` to stop.

### 3. Run the Go Service

```bash
bazel run //go_service:server
```

Visit `http://localhost:8080` - you should see:
```
Hello from Go! 🚀
```

Press `Ctrl+C` to stop.

### 4. Run Tests

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

## Understanding Bazel Files

### MODULE.bazel

The new standard for defining external dependencies using Bzlmod (Bazel Module system).

```python
module(
    name = "bazel_multilang_tutorial",
    version = "1.0.0",
)

# Python rules
bazel_dep(name = "rules_python", version = "0.31.0")

# Go rules
bazel_dep(name = "rules_go", version = "0.46.0")
```

**Why MODULE.bazel instead of WORKSPACE?**
- Cleaner dependency management
- Better version resolution
- More maintainable for large projects
- The new Bazel standard (Bzlmod)

### BUILD.bazel Files

Each sub-project has its own build file.

**Example (Go):**
```python
go_binary(
    name = "server",
    srcs = ["main.go"],
)

go_test(
    name = "go_service_test",
    srcs = ["main_test.go"],
)
```

**Example (Python):**
```python
py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = ["//py_service:requirements"],
)

py_test(
    name = "py_service_test",
    srcs = ["test_main.py"],
    deps = [":server"],
)
```

### Labels & Targets

Each target (e.g. `//py_service:server`) uniquely identifies a buildable or testable entity.

**Label syntax:**
```
//package/path:target_name
```

**Examples:**
- `//go_service:server` - The "server" binary in go_service package
- `//py_service:py_service_test` - The test in py_service package
- `//...` - All targets in all packages (recursive)

## Incremental Builds Demo

Try this to see Bazel's incremental build in action:

```bash
# 1. Build everything
bazel build //...

# 2. Change only Python code
echo "# New comment" >> py_service/main.py

# 3. Build again - only Python targets rebuild!
bazel build //...

# 4. Change Go code
echo "// New comment" >> go_service/main.go

# 5. Build again - only Go targets rebuild!
bazel build //...
```

Notice how Bazel only rebuilds what changed. This is one of Bazel's superpowers! 🚀

## Exploring the Services

### Python Service (py_service/)

**Files:**
- `main.py` - FastAPI application
- `test_main.py` - Unit tests
- `requirements.txt` - Python dependencies
- `BUILD.bazel` - Bazel build configuration

**Key code (main.py):**
```python
from fastapi import FastAPI
import uvicorn

app = FastAPI()

@app.get("/")
def read_root():
    return {"message": "Hello from Python! 🐍"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8081)
```

### Go Service (go_service/)

**Files:**
- `main.go` - HTTP server
- `main_test.go` - Unit tests
- `BUILD.bazel` - Bazel build configuration

**Key code (main.go):**
```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from Go! 🚀")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

## Testing

### Run All Tests
```bash
bazel test //...
```

### Run Specific Tests
```bash
# Python tests
bazel test //py_service:py_service_test

# Go tests
bazel test //go_service:go_service_test

# Integration tests
bazel test //tests:test_go_service
bazel test //tests:test_py_service
```

### Test with Verbose Output
```bash
bazel test //... --test_output=all
```

### Test and Show Only Failures
```bash
bazel test //... --test_output=errors
```

## CI Integration

The repository includes `.github/workflows/ci.yaml` to demonstrate Bazel in CI.

**Example job:**
```yaml
- name: Build and Test
  run: |
    bazel build //...
    bazel test //...
```

Bazel's incremental caching allows the same artifacts to be reused between CI runs.

## Learning Exercises

Try these exercises to deepen your understanding:

### Exercise 1: Modify and Rebuild
1. Modify `py_service/main.py` - change the message
2. Run `bazel build //py_service:server`
3. Notice how only Python targets rebuild

### Exercise 2: Add a Dependency
1. Add `requests` to `py_service/requirements.txt`
2. Update `py_service/main.py` to use requests
3. Run `bazel build //py_service:server`
4. Observe how Bazel fetches the dependency

### Exercise 3: Add a Test
1. Create a new test in `py_service/test_main.py`
2. Run `bazel test //py_service:py_service_test`
3. See the test pass!

### Exercise 4: Explore Dependencies
```bash
# Show what py_service depends on
bazel query 'deps(//py_service:server)'

# Show what depends on py_service
bazel query 'rdeps(//..., //py_service:server)'

# Visualize the dependency graph
bazel query --output=graph //... > graph.dot
```

## Key Takeaways

✅ Bazel provides **deterministic builds** — same inputs → same outputs

✅ It understands dependencies deeply, so it **rebuilds only what changed**

✅ It can build **different languages together**, reproducibly

✅ **Bzlmod (MODULE.bazel)** is the new standard, replacing WORKSPACE

✅ Ideal for large repos, CI optimization, and polyglot systems

## Next Steps

- **[Economic Dashboard](economic-dashboard.md)** - Build a full-stack application
- **[Bazel Concepts](concepts.md)** - Deep dive into Bazel fundamentals
- **[Dependencies](dependencies.md)** - Learn how to add dependencies
- **[Query Guide](query.md)** - Explore the build graph

## Troubleshooting

### "command not found: bazel"
- Install Bazel: https://bazel.build/install
- Or use Bazelisk (recommended): https://github.com/bazelbuild/bazelisk

### "Failed to fetch external repository"
- Check your internet connection
- Run `bazel sync` to re-fetch dependencies
- Check MODULE.bazel for typos

### Builds are slow
- First build downloads all dependencies (slow)
- Subsequent builds use cache (fast!)
- Enable remote caching for team collaboration

---

**Continue to:** [Economic Dashboard Tutorial →](economic-dashboard.md)
