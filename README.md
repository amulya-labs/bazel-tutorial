🧱 README.md — Multi-Language Bazel Tutorial
🚀 Overview

Welcome to the Bazel Multi-Language Tutorial 👋
This repository is a hands-on introduction to Bazel — Google’s open-source build system for fast, reproducible, and scalable builds across multiple languages.

You’ll learn how to:

Build and test apps in Python and Go

Use Bazel’s WORKSPACE and BUILD.bazel files

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
├── WORKSPACE               # Defines external dependencies and Bazel rules
├── go_service/             # Go service + BUILD.bazel
│   ├── main.go
│   └── BUILD.bazel
├── py_service/             # Python service + BUILD.bazel
│   ├── main.py
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
🧰 WORKSPACE

Defines external dependencies and language rules.
For example:

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Go rules
http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/v0.47.0/rules_go-v0.47.0.tar.gz"],
    sha256 = "...",
)

# Python rules
http_archive(
    name = "rules_python",
    urls = ["https://github.com/bazelbuild/rules_python/releases/download/0.31.0/rules_python-0.31.0.tar.gz"],
    sha256 = "...",
)

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

Bazel provides deterministic builds — same inputs → same outputs.

It understands dependencies deeply, so it rebuilds only what changed.

It can build different languages together, reproducibly.

Ideal for large repos, CI optimization, and polyglot systems.

📚 Further Reading

Bazel Basics

rules_python

rules_go

Bazel Query Guide

Bazel Remote Caching
