# Adding Dependencies Guide

Learn how to add new dependencies to your Bazel project.

## Table of Contents
- [Python Dependencies](#python-dependencies)
- [Go Dependencies](#go-dependencies)
- [Adding New Services](#adding-new-services)
- [Common Issues](#common-issues)

## Python Dependencies

### 1. Add to requirements.txt

Edit `py_service/requirements.txt`:

```txt
fastapi==0.109.0
uvicorn==0.27.0
httpx==0.26.0
requests==2.31.0  # ← New dependency
```

### 2. Update BUILD.bazel

Edit `py_service/BUILD.bazel` to include the new dependency:

```python
py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = [
        requirement("fastapi"),
        requirement("uvicorn"),
        requirement("requests"),  # ← Add here
    ],
)
```

### 3. Rebuild

```bash
bazel build //py_service:server
```

Bazel will automatically fetch the new dependency.

### 4. Use in Code

Now you can import it in `main.py`:

```python
import requests

@app.get("/fetch")
def fetch_data():
    response = requests.get("https://api.example.com/data")
    return response.json()
```

### Example: Adding pandas

**1. Update requirements.txt:**
```txt
pandas==2.1.4
```

**2. Update BUILD.bazel:**
```python
py_binary(
    name = "analytics",
    srcs = ["analytics.py"],
    deps = [
        requirement("pandas"),
    ],
)
```

**3. Use in code:**
```python
import pandas as pd

df = pd.DataFrame({"col1": [1, 2], "col2": [3, 4]})
print(df.head())
```

## Go Dependencies

Go dependencies are managed automatically through `go.mod`. Bazel's Gazelle tool can generate and update BUILD files.

### Method 1: Manual (Simple)

#### 1. Import in your code

Edit `go_service/main.go`:

```go
package main

import (
    "github.com/gin-gonic/gin"  // ← New import
)

func main() {
    r := gin.Default()
    r.GET("/", handler)
    r.Run(":8080")
}
```

#### 2. Initialize go.mod (if not exists)

```bash
cd go_service
go mod init github.com/yourusername/bazel-tutorial/go_service
go mod tidy
```

#### 3. Update MODULE.bazel

Add Go dependency configuration:

```python
go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//go_service:go.mod")
```

#### 4. Update BUILD.bazel

```python
go_binary(
    name = "server",
    srcs = ["main.go"],
    deps = ["@com_github_gin_gonic_gin//:gin"],
)
```

### Method 2: Using Gazelle (Recommended)

Gazelle automatically generates BUILD files for Go projects.

#### 1. Add Gazelle target to root BUILD.bazel

Create `/BUILD.bazel`:

```python
load("@gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")
```

#### 2. Run Gazelle

```bash
bazel run //:gazelle
```

This automatically:
- Detects Go imports
- Updates BUILD.bazel files
- Adds dependencies

#### 3. Update dependencies

```bash
bazel run //:gazelle -- update-repos -from_file=go_service/go.mod
```

### Example: Adding UUID library

**1. Import in code:**
```go
import "github.com/google/uuid"

func handler(w http.ResponseWriter, r *http.Request) {
    id := uuid.New()
    fmt.Fprintf(w, "Request ID: %s\n", id.String())
}
```

**2. Update go.mod:**
```bash
cd go_service
go get github.com/google/uuid
go mod tidy
```

**3. Run Gazelle:**
```bash
bazel run //:gazelle
```

Done! Gazelle updates BUILD files automatically.

## Adding New Services

### Adding a Rust Service

#### 1. Create directory structure

```bash
mkdir -p rust_service
cd rust_service
```

#### 2. Add rules_rust to MODULE.bazel

```python
bazel_dep(name = "rules_rust", version = "0.40.0")

rust = use_extension("@rules_rust//rust:extensions.bzl", "rust")
rust.toolchain(
    edition = "2021",
    versions = ["1.75.0"],
)
use_repo(rust, "rust_toolchains")
```

#### 3. Create Cargo.toml

```toml
[package]
name = "rust_service"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4.4"
```

#### 4. Create BUILD.bazel

```python
load("@rules_rust//rust:defs.bzl", "rust_binary")

rust_binary(
    name = "server",
    srcs = ["src/main.rs"],
)
```

#### 5. Add to CI

Update `.github/workflows/ci.yaml`:

```yaml
- name: Build Rust service
  run: bazel build //rust_service:server
```

### Adding a Node.js Service

#### 1. Add rules_nodejs to MODULE.bazel

```python
bazel_dep(name = "aspect_rules_js", version = "1.35.0")

npm = use_extension("@aspect_rules_js//npm:extensions.bzl", "npm")
npm.npm_translate_lock(
    name = "npm",
    pnpm_lock = "//nodejs_service:pnpm-lock.yaml",
)
use_repo(npm, "npm")
```

#### 2. Create package.json

```json
{
  "name": "nodejs_service",
  "dependencies": {
    "express": "^4.18.0"
  }
}
```

#### 3. Create BUILD.bazel

```python
load("@aspect_rules_js//js:defs.bzl", "js_binary")
load("@npm//:defs.bzl", "npm_link_all_packages")

npm_link_all_packages(name = "node_modules")

js_binary(
    name = "server",
    entry_point = "index.js",
    data = [":node_modules"],
)
```

## Version Pinning

### Python: requirements.txt

Always pin exact versions for reproducibility:

**Good:**
```txt
fastapi==0.109.0
uvicorn==0.27.0
```

**Avoid:**
```txt
fastapi>=0.100.0  # Version can change!
uvicorn           # No version specified!
```

### Go: go.mod

Go modules already pin versions:

```go
module github.com/example/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.6.0
)
```

### MODULE.bazel

Bazel dependencies are pinned in MODULE.bazel:

```python
bazel_dep(name = "rules_python", version = "0.31.0")
```

## Common Issues

### "Cannot find module"

**Python:**
- Check `requirements.txt` has the package
- Verify `requirement("package")` in BUILD.bazel
- Run `bazel clean` and rebuild

**Go:**
- Ensure `go.mod` exists and has the module
- Run `go mod tidy`
- Run `bazel run //:gazelle` if using Gazelle

### "Version conflict"

Bazel's Bzlmod resolves version conflicts automatically. Check:

```bash
bazel mod graph
```

To see the dependency graph and selected versions.

### "External repository not found"

Network issue or typo. Try:

```bash
bazel sync --configure
```

### Stale cache

After changing dependencies:

```bash
bazel clean
bazel build //...
```

## Best Practices

1. **Pin exact versions**: Ensures reproducible builds
2. **Use lock files**: `requirements.txt`, `go.mod`, `pnpm-lock.yaml`
3. **Minimize dependencies**: Fewer deps = faster builds
4. **Use monorepo structure**: Share dependencies across services
5. **Document dependencies**: Add comments explaining why each is needed
6. **Periodic updates**: Keep dependencies up to date for security

## Resources

- [rules_python pip integration](https://github.com/bazelbuild/rules_python/blob/main/docs/pip.md)
- [rules_go dependencies](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/dependencies.md)
- [Gazelle documentation](https://github.com/bazelbuild/bazel-gazelle)
- [Bazel Registry](https://registry.bazel.build/) - Browse available modules
