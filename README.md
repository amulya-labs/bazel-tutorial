# 🧱 Bazel Multi-Language Tutorial (with Economic Dashboard)

A comprehensive, hands-on guide to building polyglot applications with Bazel — Google's fast, reproducible, and scalable build system.

**📚 [Full Documentation →](https://rrl-personal-projects.github.io/bazel-tutorial/)**

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

**[Full Quick Start Guide →](https://rrl-personal-projects.github.io/bazel-tutorial/quickstart/)**

## 📚 What's Inside

### 1️⃣ Basic Tutorial Services

Simple services to learn Bazel fundamentals:

| Service | Language | Description | Command |
|---------|----------|-------------|---------|
| `py_service` | Python (FastAPI) | "Hello from Python" web service | `bazel run //py_service:server` |
| `go_service` | Go | "Hello from Go" web service | `bazel run //go_service:server` |

**[Learn More →](https://rrl-personal-projects.github.io/bazel-tutorial/basic-tutorial/)**

### 2️⃣ Economic Indicators Dashboard

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
- 📊 Go service fetches economic data from FRED API
- 🐍 Python FastAPI serves REST endpoints with computed metrics
- ⚛️ React frontend with interactive charts and responsive design
- 📈 Tracks 6 core economic indicators (CPI, Unemployment, Fed Funds, etc.)

**Quick Start:**
```bash
# Get a free FRED API key from https://fred.stlouisfed.org
export FRED_API_KEY=your_key_here

# Run everything
./run_all.sh

# In another terminal, start the UI
cd web_ui && npm install && npm run dev
```

Open http://localhost:3000 in your browser! 🎉

**[Full Dashboard Documentation →](https://rrl-personal-projects.github.io/bazel-tutorial/economic-dashboard/)**

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
├── docs/                    # Documentation (MkDocs)
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

**[Learn more about Bazel concepts →](https://rrl-personal-projects.github.io/bazel-tutorial/concepts/)**

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

**[Complete testing guide →](https://rrl-personal-projects.github.io/bazel-tutorial/testing/)**

## 📖 Documentation

Comprehensive guides available at **[https://rrl-personal-projects.github.io/bazel-tutorial/](https://rrl-personal-projects.github.io/bazel-tutorial/)**

- **[Quick Start Guide](https://rrl-personal-projects.github.io/bazel-tutorial/quickstart/)** - Get up and running quickly
- **[Installation Guide](https://rrl-personal-projects.github.io/bazel-tutorial/installation/)** - Install all prerequisites
- **[Basic Tutorial](https://rrl-personal-projects.github.io/bazel-tutorial/basic-tutorial/)** - Learn with simple services
- **[Economic Dashboard](https://rrl-personal-projects.github.io/bazel-tutorial/economic-dashboard/)** - Full-stack application
- **[Bazel Concepts](https://rrl-personal-projects.github.io/bazel-tutorial/concepts/)** - Deep dive into fundamentals
- **[Dependencies](https://rrl-personal-projects.github.io/bazel-tutorial/dependencies/)** - Managing dependencies
- **[Query Guide](https://rrl-personal-projects.github.io/bazel-tutorial/query/)** - Explore the build graph
- **[Testing Guide](https://rrl-personal-projects.github.io/bazel-tutorial/testing/)** - Testing strategies
- **[Contributing](https://rrl-personal-projects.github.io/bazel-tutorial/contributing/)** - Development guidelines

## 🧠 Why Bazel for Polyglot Repos?

### Benefits Demonstrated

1. **Single Build System** - One tool for Go, Python, and JavaScript
2. **Incremental Builds** - Only rebuild what changed
3. **Hermetic Testing** - Tests don't depend on network or system state
4. **Explicit Dependencies** - Clear dependency graph
5. **Scalability** - Works from small projects to massive monorepos

## 🔧 Troubleshooting

### Common Issues

**"command not found: bazel"**
- Install Bazel: https://bazel.build/install
- Or use Bazelisk: https://github.com/bazelbuild/bazelisk

**"Failed to fetch external repository"**
- Check internet connection
- Run `bazel sync` to re-fetch dependencies

**"Database not found" (Economic Dashboard)**
- Run: `bazel run //go_fetch:refresh`
- Set: `export FRED_API_KEY=your_key`

**See more:** [Full Troubleshooting Guide](https://rrl-personal-projects.github.io/bazel-tutorial/installation/#common-installation-issues)

### Validation Script

Verify your setup is working:

```bash
./validate.sh
```

## 🤝 Contributing

Contributions are welcome! See **[Contributing Guide](https://rrl-personal-projects.github.io/bazel-tutorial/contributing/)** for guidelines on:
- Adding new services
- Improving documentation
- Adding tests
- Submitting pull requests

## 📚 Further Reading

- [Bazel Documentation](https://bazel.build/)
- [Bzlmod Guide](https://bazel.build/external/overview#bzlmod)
- [rules_python](https://github.com/bazelbuild/rules_python)
- [rules_go](https://github.com/bazelbuild/rules_go)
- [FRED API Documentation](https://fred.stlouisfed.org/docs/api/)

## 📄 License

MIT License - See LICENSE file for details

---

**Happy Building! 🚀**

For questions or issues, open a [GitHub issue](https://github.com/rrl-personal-projects/bazel-tutorial/issues).
