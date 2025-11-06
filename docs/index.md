# Bazel Multi-Language Tutorial

!!! success "🎯 View Live Dashboard"
    The Economic Indicators Dashboard is deployed live on GitHub Pages:

    **[🚀 Open Live Dashboard →](/bazel-tutorial/app/)**

    The dashboard automatically updates daily with fresh economic data from FRED API.

!!! tip "Run Locally"
    Want to run it yourself?

    1. Clone the repo and set your FRED API key
    2. Run `./run_all.sh`
    3. Start the UI: `cd web_ui && npm run dev`
    4. Open [http://localhost:3000](http://localhost:3000)

    **[📖 Full Setup Instructions →](economic-dashboard.md#quick-start)**

A comprehensive, hands-on guide to building polyglot applications with Bazel — Google's fast, reproducible, and scalable build system.

## 🎯 What You'll Learn

This repository demonstrates real-world Bazel usage across multiple languages:

- ✅ **Build Go, Python, and React services** in one unified repository
- ✅ **Incremental compilation and caching** for lightning-fast builds
- ✅ **Hermetic testing** without network dependencies
- ✅ **Modern toolchain integration** (rules_go, rules_python, npm/vite)
- ✅ **CI/CD with GitHub Actions** and Bazel caching
- ✅ **Real-world architecture patterns** for data fetching, APIs, and UIs

## 🚀 Quick Start

```bash
# Install Bazel
brew install bazel  # macOS
# or download from https://bazel.build/install

# Clone the repository
git clone https://github.com/rrl-personal-projects/bazel-tutorial.git
cd bazel-tutorial

# Build everything
bazel build //...

# Run tests
bazel test //...
```

**[Full Quick Start Guide →](quickstart.md)**

## 📚 What's Inside

### Basic Tutorial Services

Simple services to learn Bazel fundamentals:

| Service | Language | Description | Command |
|---------|----------|-------------|---------|
| `py_service` | Python (FastAPI) | "Hello from Python" web service | `bazel run //py_service:server` |
| `go_service` | Go | "Hello from Go" web service | `bazel run //go_service:server` |

**[Learn More →](basic-tutorial.md)**

### 📊 Economic Indicators Dashboard

A full-stack polyglot application demonstrating real-world Bazel usage:

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Go Fetch  │  writes │   SQLite DB  │  reads  │  Python API │
│   Service   │────────→│  (econ.db)   │←────────│  (FastAPI)  │
└─────────────┘         └──────────────┘         └─────────────┘
                                                         │
                                                         │ HTTP
                                                         ↓
                                                  ┌─────────────┐
                                                  │  React UI   │
                                                  │   (Vite)    │
                                                  └─────────────┘
```

**Features:**
- Go service fetches economic data from FRED API
- Python FastAPI serves REST endpoints with computed metrics
- React frontend with interactive charts and responsive design
- Tracks 6 core economic indicators (CPI, Unemployment, Fed Funds, etc.)

**Quick Start:**
```bash
# Get a free FRED API key from https://fred.stlouisfed.org
export FRED_API_KEY=your_key_here

# Run everything
./run_all.sh

# In another terminal, start the UI
cd web_ui && npm install && npm run dev
```

**[Full Dashboard Documentation →](economic-dashboard.md)**

## 🏗️ Repository Structure

```
bazel-tutorial/
├── MODULE.bazel              # Bazel dependencies (Bzlmod)
├── .bazelrc                  # Bazel configuration
├── .bazelversion             # Pin Bazel version
│
├── go_service/              # Basic Go service tutorial
├── py_service/              # Basic Python service tutorial
│
├── go_fetch/                # Economic data fetcher (Go)
├── py_api/                  # REST API server (Python)
├── web_ui/                  # Dashboard UI (React + TypeScript)
│
├── docs/                    # Documentation
├── tests/                   # Integration tests
│
└── .github/workflows/       # CI/CD pipelines
```

## 🧱 Key Bazel Concepts

### MODULE.bazel (Bzlmod)

Modern external dependency management:

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

**[Learn more about Bazel concepts →](concepts.md)**

### BUILD.bazel Files

Each component has a `BUILD.bazel` file defining buildable targets:

```python
# Python binary
py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = ["//py_service:requirements"],
)

# Go binary
go_binary(
    name = "server",
    srcs = ["main.go"],
)
```

### Targets and Labels

- `//go_service:server` - The "server" target in the "go_service" package
- `//py_api:server` - Python API server binary
- `//...` - All targets in all packages (recursive)

**[Deep dive into dependencies →](dependencies.md)**

## 🧪 Testing

All tests use fixtures and mocks - no network access required!

```bash
# Run all tests
bazel test //...

# Run specific tests
bazel test //go_fetch:fred_test
bazel test //py_api:api_test

# View test output
bazel test //... --test_output=all
```

**[Complete testing guide →](testing.md)**

## 🔄 Incremental Builds

Bazel only rebuilds what changed:

```bash
# First build (downloads deps, compiles everything)
bazel build //...

# Edit Python code - only Python targets rebuild!
echo "# comment" >> py_service/main.py
bazel build //...

# Edit Go code - only Go targets rebuild!
echo "// comment" >> go_service/main.go
bazel build //...
```

## ⚙️ CI/CD Integration

The repository includes GitHub Actions workflows demonstrating:

- Bazel build and test automation
- Dependency caching for fast CI runs
- Multi-language build verification

```yaml
- name: Build and Test
  run: |
    bazel build //...
    bazel test //...
```

Bazel's incremental caching allows the same artifacts to be reused between CI runs.

## 🧠 Why Bazel for Polyglot Repos?

### Benefits Demonstrated in This Project

1. **Single Build System**
   - One tool for Go, Python, and JavaScript
   - Consistent commands across languages

2. **Incremental Builds**
   - Only rebuild what changed
   - Shared cache across projects

3. **Hermetic Testing**
   - Tests don't depend on network or system state
   - Use mocks and fixtures

4. **Explicit Dependencies**
   - Clear dependency graph
   - No hidden dependencies

5. **Scalability**
   - Works for small projects (this) and massive monorepos (Google)
   - Parallel builds and tests

## 📖 Documentation

Explore comprehensive guides:

- **[Quick Start Guide](quickstart.md)** - Get up and running quickly
- **[Basic Tutorial](basic-tutorial.md)** - Learn with simple services
- **[Economic Dashboard](economic-dashboard.md)** - Full-stack application example
- **[Bazel Concepts](concepts.md)** - Deep dive into Bazel fundamentals
- **[Dependencies](dependencies.md)** - Managing Python, Go, and npm dependencies
- **[Query Guide](query.md)** - Explore the build graph with `bazel query`
- **[Testing Guide](testing.md)** - Comprehensive testing strategies
- **[Contributing](contributing.md)** - Development guidelines

## 🔧 Troubleshooting

### Common Issues

**"command not found: bazel"**
- Install Bazel: https://bazel.build/install
- Or use Bazelisk (recommended): https://github.com/bazelbuild/bazelisk

**"Failed to fetch external repository"**
- Check your internet connection
- Run `bazel sync` to re-fetch dependencies
- Check MODULE.bazel for typos

**"Builds are slow"**
- First build downloads all dependencies (slow)
- Subsequent builds use cache (fast!)
- Enable remote caching for team collaboration

**"Database not found" (Economic Dashboard)**
- Run the data fetcher first: `bazel run //go_fetch:refresh`
- Make sure you set `FRED_API_KEY` environment variable

### Validation Script

To verify your setup is working correctly:

```bash
./validate.sh
```

This script will:
- Check if Bazel is installed
- Validate repository structure
- Build all targets
- Run all tests

## 🤝 Contributing

Contributions are welcome! See **[CONTRIBUTING.md](contributing.md)** for guidelines on:
- Adding new services
- Improving documentation
- Adding tests
- Submitting pull requests

## 📚 Further Reading

### Official Documentation

- [Bazel Basics](https://bazel.build/basics)
- [Bzlmod Guide](https://bazel.build/external/overview#bzlmod) - The new MODULE.bazel system
- [rules_python](https://github.com/bazelbuild/rules_python)
- [rules_go](https://github.com/bazelbuild/rules_go)
- [Bazel Query Guide](https://bazel.build/query/guide)
- [Bazel Remote Caching](https://bazel.build/remote/caching)

### Economic Data

- [FRED API Documentation](https://fred.stlouisfed.org/docs/api/)

## 📄 License

MIT License - See LICENSE file for details

---

**Happy Building! 🚀**

For questions or issues, please open a [GitHub issue](https://github.com/rrl-personal-projects/bazel-tutorial/issues).
