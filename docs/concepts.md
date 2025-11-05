# Understanding Bazel Concepts

This guide provides a deeper explanation of key Bazel concepts used in this tutorial.

## Table of Contents
- [MODULE.bazel and Bzlmod](#modulebazel-and-bzlmod)
- [BUILD.bazel Files](#buildbazel-files)
- [Targets and Labels](#targets-and-labels)
- [Dependencies](#dependencies)
- [Caching and Incrementality](#caching-and-incrementality)

## MODULE.bazel and Bzlmod

**Bzlmod** is Bazel's modern external dependency management system, replacing the legacy WORKSPACE approach.

### Key Benefits:
- **Version Resolution**: Automatic conflict resolution for transitive dependencies
- **Reproducibility**: Lock files ensure consistent builds across machines
- **Modularity**: Each module declares its own dependencies
- **Simpler Syntax**: Cleaner than WORKSPACE's `http_archive` approach

### Structure:
```python
module(
    name = "my_project",  # Your project name
    version = "1.0.0",     # Semantic version
)

# Declare dependencies
bazel_dep(name = "rules_python", version = "0.31.0")
bazel_dep(name = "rules_go", version = "0.46.0")
```

### Extensions:
Extensions allow you to configure rules:
```python
python = use_extension("@rules_python//python/extensions:python.bzl", "python")
python.toolchain(python_version = "3.11")
```

## BUILD.bazel Files

BUILD.bazel files define **targets** — the buildable units in your project.

### Common Target Types:

#### `py_binary` - Python Executable
```python
py_binary(
    name = "my_app",
    srcs = ["main.py"],
    deps = [
        requirement("fastapi"),
    ],
)
```

#### `go_binary` - Go Executable
```python
go_binary(
    name = "server",
    srcs = ["main.go"],
)
```

#### `py_test` - Python Test
```python
py_test(
    name = "my_test",
    srcs = ["test_main.py"],
    deps = [":my_app"],
)
```

#### `go_test` - Go Test
```python
go_test(
    name = "my_test",
    srcs = ["main_test.go"],
)
```

## Targets and Labels

A **label** uniquely identifies a target in your workspace.

### Label Syntax:
```
//package/path:target_name
```

Examples:
- `//go_service:server` - The "server" target in the "go_service" package
- `//py_service:py_service_test` - The test target in "py_service"
- `//:all` - All targets in the root package

### Special Labels:
- `//...` - All targets in all packages (recursive)
- `//py_service/...` - All targets under py_service (recursive)

## Dependencies

### Direct Dependencies
Declared in the `deps` attribute:
```python
py_binary(
    name = "server",
    deps = [
        requirement("fastapi"),  # External pip dependency
        "//common:utils",         # Internal dependency
    ],
)
```

### Transitive Dependencies
Bazel automatically tracks transitive dependencies. If A depends on B, and B depends on C, Bazel knows A transitively depends on C.

### Data Dependencies
Use `data` for runtime files:
```python
py_test(
    name = "my_test",
    data = [
        "//config:test_config.yaml",
        "//go_service:server",
    ],
)
```

## Caching and Incrementality

### How Bazel Caches Work:

1. **Content Hashing**: Bazel hashes all inputs (source files, dependencies, build commands)
2. **Action Graph**: Creates a graph of all build actions
3. **Cache Lookup**: Before building, checks if outputs exist for the same inputs
4. **Incremental Builds**: Only rebuilds what changed

### Example:
```bash
# First build - builds everything
$ bazel build //...

# Change only Python code
$ echo "# comment" >> py_service/main.py

# Second build - only rebuilds Python targets
$ bazel build //...
# Go targets are cached and not rebuilt!
```

### Cache Types:

#### Local Cache
Stored on disk in `~/.cache/bazel`:
```bash
build --disk_cache=~/.cache/bazel
```

#### Remote Cache (Advanced)
Shared cache across teams/CI:
```bash
build --remote_cache=grpc://cache-server:9092
```

### Viewing Cache Stats:
```bash
bazel info
bazel clean --expunge  # Clear all cache
```

## Hermetic Builds

Bazel ensures **hermetic** (isolated) builds:
- Uses declared dependencies only
- Ignores system-installed packages
- Same inputs → same outputs (reproducible)

This is why we use:
- `rules_python` for Python toolchains (not system Python)
- `rules_go` for Go toolchains (not system Go)

## Query Commands

Explore your build graph:

```bash
# List all targets
bazel query //...

# Find dependencies of a target
bazel query 'deps(//go_service:server)'

# Find what depends on a target
bazel query 'rdeps(//..., //py_service:server)'

# Visualize build graph
bazel query --output=graph //... > graph.dot
dot -Tpng graph.dot -o graph.png
```

## Best Practices

1. **Keep BUILD files small**: One package per directory
2. **Use visibility**: Control who can depend on your targets
3. **Minimize dependencies**: Faster builds, clearer design
4. **Use tests extensively**: `bazel test //...` is fast
5. **Enable remote cache**: Share build artifacts across team
6. **Pin versions**: Use `.bazelversion` for reproducibility

## Further Reading

- [Bazel Concepts](https://bazel.build/concepts/build-ref)
- [Bzlmod Guide](https://bazel.build/external/overview#bzlmod)
- [rules_python Documentation](https://github.com/bazelbuild/rules_python)
- [rules_go Documentation](https://github.com/bazelbuild/rules_go)
