# Bazel Query Examples

This guide demonstrates how to use `bazel query` and `bazel aquery` to explore the build graph.

## Basic Queries

### List all targets in the workspace
```bash
bazel query //...
```

Output:
```
//go_service:go_service_test
//go_service:server
//py_service:py_service_test
//py_service:server
//tests:test_go_service
//tests:test_py_service
```

### List targets in a specific package
```bash
bazel query //go_service:all
```

### Find all test targets
```bash
bazel query 'kind(test, //...)'
```

### Find all Python targets
```bash
bazel query 'kind(py_.*, //...)'
```

### Find all Go targets
```bash
bazel query 'kind(go_.*, //...)'
```

## Dependency Queries

### Show all dependencies of a target
```bash
bazel query 'deps(//go_service:server)'
```

This shows everything the Go server depends on, including:
- Source files
- Toolchains
- External dependencies

### Show only direct dependencies
```bash
bazel query 'deps(//go_service:server, 1)'
```

The `1` limits the depth to immediate dependencies only.

### Find what depends on a target (reverse dependencies)
```bash
bazel query 'rdeps(//..., //go_service:server)'
```

This shows all targets that depend on the Go server.

### Show path between two targets
```bash
bazel query 'somepath(//py_service:server, @pip_deps//:fastapi)'
```

## Visualization

### Generate a dependency graph
```bash
bazel query 'deps(//go_service:server)' --output=graph > go_deps.dot
```

Then convert to image (requires Graphviz):
```bash
dot -Tpng go_deps.dot -o go_deps.png
```

### Generate graph for entire workspace
```bash
bazel query '//...' --output=graph > workspace.dot
dot -Tpng workspace.dot -o workspace.png
```

## Action Queries (Advanced)

`aquery` shows the actual actions (compile, link, etc.) that Bazel executes.

### Show all actions for building a target
```bash
bazel aquery //go_service:server
```

### Show only compilation actions
```bash
bazel aquery 'mnemonic("GoCompile", //go_service:server)'
```

### Show linking actions
```bash
bazel aquery 'mnemonic("GoLink", //go_service:server)'
```

### Show actions with inputs and outputs
```bash
bazel aquery //go_service:server --output=text
```

## Practical Examples

### Find all binary targets (runnable services)
```bash
bazel query 'kind(".*_binary", //...)'
```

Output:
```
//go_service:server
//py_service:server
```

### Find targets that changed since last commit
```bash
bazel query --universe_scope=//... \
  --order_output=no \
  "set(//...)" --output=label
```

### Check if target exists
```bash
bazel query //py_service:server &>/dev/null && echo "Target exists" || echo "Target not found"
```

### Find unused dependencies (requires analysis)
```bash
# This is more complex and requires additional tooling
bazel query "rdeps(//..., //some:target, 1)" --output=label_kind
```

## Filtering and Formatting

### Output as a list
```bash
bazel query '//...' --output=label
```

### Output as JSON
```bash
bazel query '//...' --output=jsonproto
```

### Output with kind (target type)
```bash
bazel query '//...' --output=label_kind
```

Example output:
```
go_binary rule //go_service:server
py_binary rule //py_service:server
go_test rule //go_service:go_service_test
py_test rule //py_service:py_service_test
```

## Query Functions Cheatsheet

| Function | Description | Example |
|----------|-------------|---------|
| `deps(target)` | All dependencies | `deps(//go_service:server)` |
| `rdeps(universe, target)` | Reverse dependencies | `rdeps(//..., //py_service:server)` |
| `allpaths(from, to)` | All paths between targets | `allpaths(//go_service:server, @rules_go//:go)` |
| `somepath(from, to)` | One path between targets | `somepath(//go_service:server, @rules_go//:go)` |
| `kind(pattern, input)` | Filter by target type | `kind(test, //...)` |
| `attr(name, pattern, input)` | Filter by attribute | `attr(srcs, ".*\.py", //...)` |
| `labels(attr, input)` | Extract label values | `labels(srcs, //go_service:server)` |
| `filter(pattern, input)` | Filter labels by regex | `filter("test", //...)` |

## Performance Tips

### Use `--universe_scope` for large queries
```bash
bazel query 'allpaths(//service:a, //service:b)' --universe_scope=//service/...
```

This limits the search space, making queries faster.

### Use `--nohost_deps` to skip host dependencies
```bash
bazel query 'deps(//go_service:server)' --nohost_deps
```

### Use `--notool_deps` to skip tool dependencies
```bash
bazel query 'deps(//go_service:server)' --notool_deps
```

## Debugging Build Issues

### Find why a target is built
```bash
bazel query 'rdeps(//..., //some:target)' --output=graph | dot -Tpng > why_built.png
```

### Check for circular dependencies
```bash
bazel query 'somepath(//pkg:target, //pkg:target)'
# If output is found, there's a cycle
```

### List all external dependencies
```bash
bazel query 'kind("http_archive", @*//...)'
```

## Further Reading

- [Bazel Query Reference](https://bazel.build/query/language)
- [Bazel Query How-To](https://bazel.build/query/guide)
- [Action Query (aquery)](https://bazel.build/query/aquery)
- [Configuration Query (cquery)](https://bazel.build/query/cquery)
