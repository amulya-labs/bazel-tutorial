🧱 README.md — Multi-Language Bazel Tutorial
🚀 Overview

Welcome to the Bazel Multi-Language Tutorial 👋
This repository is a hands-on introduction to Bazel — Google’s open-source build system for fast, reproducible, and scalable builds across multiple languages.

You’ll learn how to:

Build and test apps in Python and Go

Use Bazel’s MODULE.bazel and BUILD.bazel files

Leverage caching, parallelism, and hermetic builds

Run Bazel inside GitHub Actions CI

🧭 What You’ll Build

This project contains two simple services:

Service	Language	Description	Run Command
py_service	Python (FastAPI)	“Hello from Python” web service	bazel run //py_service:server
go_service	Go	“Hello from Go” web service	bazel run //go_service:server

You’ll build, test, and run them independently — and learn how Bazel manages both under one unified build system.

🏗️ Repository Structure
bazel-multilang-tutorial/
├── MODULE.bazel            # Defines external dependencies using Bzlmod (new standard)
├── .bazelversion           # Specifies the Bazel version to use
├── .bazelrc                # Bazel configuration options
├── go_service/             # Go service + BUILD.bazel
│   ├── main.go
│   ├── main_test.go
│   └── BUILD.bazel
├── py_service/             # Python service + BUILD.bazel
│   ├── main.py
│   ├── test_main.py
│   ├── requirements.txt
│   └── BUILD.bazel
├── tests/                  # Cross-service tests
│   ├── test_go_service.sh
│   ├── test_py_service.py
│   └── BUILD.bazel
└── .github/workflows/ci.yaml  # Bazel build/test on GitHub Actions

🪄 Quick Start
1. Install Bazel

Follow instructions for your OS:
👉 https://bazel.build/install

Verify:

bazel --version

2. Build Everything
bazel build //...


This builds all targets defined in every BUILD.bazel.

3. Run Each Service
bazel run //py_service:server
bazel run //go_service:server

4. Run Tests
bazel test //...

📦 Understanding Bazel Files
🧰 MODULE.bazel

The new standard for defining external dependencies using Bzlmod (Bazel Module system).
For example:

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

🧱 BUILD.bazel

Each sub-project has its own build file.
Example (Go):

go_binary(
    name = "server",
    srcs = ["main.go"],
)


Example (Python):

py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = ["//py_service:requirements"],
)

🧩 Labels & Targets

Each target (e.g. //py_service:server) uniquely identifies a buildable or testable entity.

⚙️ CI Integration

The repo includes a .github/workflows/ci.yaml file to demonstrate Bazel in CI.

Example job:

- name: Build and Test
  run: |
    bazel build //...
    bazel test //...


Bazel’s incremental caching allows the same artifacts to be reused between CI runs.

🧪 Learn by Changing Things

Try these exercises:

Modify main.py and rerun bazel build //py_service:server.
→ Notice how Bazel only rebuilds that target.

Add a new dependency in requirements.txt.
→ Observe how Bazel fetches it via rules_python.

Edit Go code and run bazel test //go_service:go_default_test.
→ Only Go targets rebuild.

Enable remote cache (optional, advanced).

🧠 Key Takeaways

✅ Bazel provides deterministic builds — same inputs → same outputs.

✅ It understands dependencies deeply, so it rebuilds only what changed.

✅ It can build different languages together, reproducibly.

✅ **Bzlmod (MODULE.bazel)** is the new standard, replacing WORKSPACE.

✅ Ideal for large repos, CI optimization, and polyglot systems.

📖 Detailed Documentation

This repository includes comprehensive guides:

- **[Quick Start Guide](docs/QUICKSTART.md)** - Get up and running quickly
- **[Bazel Concepts](docs/CONCEPTS.md)** - Deep dive into Bazel fundamentals
- **[Adding Dependencies](docs/DEPENDENCIES.md)** - How to add Python, Go, and other dependencies
- **[Query Guide](docs/QUERY.md)** - Explore the build graph with `bazel query`

## 🔧 Troubleshooting

### Common Issues

**"command not found: bazel"**
- Install Bazel: https://bazel.build/install
- Or use Bazelisk (recommended): https://github.com/bazelbuild/bazelisk

**"Failed to fetch external repository"**
- Check your internet connection
- Run `bazel sync` to re-fetch dependencies
- Check MODULE.bazel for typos

**Builds are slow**
- First build downloads all dependencies (slow)
- Subsequent builds use cache (fast!)
- Enable remote caching for team collaboration

**Need help?**
- Check [docs/QUICKSTART.md](docs/QUICKSTART.md) for detailed setup
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
- Open an issue for bugs or questions

📚 Further Reading

📘 Official Documentation

[Bazel Basics](https://bazel.build/basics)

[Bzlmod Guide](https://bazel.build/external/overview#bzlmod) - The new MODULE.bazel system

[rules_python](https://github.com/bazelbuild/rules_python)

[rules_go](https://github.com/bazelbuild/rules_go)

[Bazel Query Guide](https://bazel.build/query/guide)

[Bazel Remote Caching](https://bazel.build/remote/caching)

🎯 Tutorial Guides

Check the [docs/](docs/) folder for detailed guides:
- [QUICKSTART.md](docs/QUICKSTART.md) - Installation and basic commands
- [CONCEPTS.md](docs/CONCEPTS.md) - Understanding Bazel fundamentals
- [DEPENDENCIES.md](docs/DEPENDENCIES.md) - Managing dependencies
- [QUERY.md](docs/QUERY.md) - Exploring the build graph

## ✅ Validation

To verify your setup is working correctly, run:

```bash
./validate.sh
```

This script will:
- Check if Bazel is installed
- Validate repository structure
- Build all targets
- Run all tests

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:
- Adding new services
- Improving documentation
- Adding tests
- Submitting pull requests

---

## 🆕 Economic Indicators Dashboard

**New: A comprehensive full-stack polyglot application!**

This advanced tutorial demonstrates real-world Bazel usage with:
- **Go** service that fetches economic data from FRED API
- **Python** FastAPI serving REST endpoints with computed metrics
- **React** frontend with interactive charts and responsive design

### Quick Start

```bash
# Get a free FRED API key from https://fred.stlouisfed.org
export FRED_API_KEY=your_key_here

# Run everything
./run_all.sh

# In another terminal, start the UI
cd web_ui && npm install && npm run dev
```

📖 **[Full Documentation](ECONOMIC_DASHBOARD.md)**

This teaching project demonstrates:
- ✅ Multi-language builds (Go + Python + React)
- ✅ Incremental compilation and caching
- ✅ Hermetic testing without network dependencies
- ✅ Modern CI/CD with GitHub Actions

---

