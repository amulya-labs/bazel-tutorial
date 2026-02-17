# Contributing to Bazel Tutorial

Thank you for your interest in improving this Bazel tutorial! 🎉

## Ways to Contribute

- **Report Issues**: Found a bug or unclear documentation? [Open an issue](https://github.com/amulya-labs/bazel-tutorial/issues)
- **Improve Documentation**: Fix typos, add examples, or clarify concepts
- **Add Examples**: Create examples for additional languages or use cases
- **Enhance Tests**: Add more comprehensive tests

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/bazel-tutorial.git
   cd bazel-tutorial
   ```
3. Create a branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Making Changes

### Adding a New Service

To add a new service (e.g., Java, Rust, Node.js):

1. **Create directory**:
   ```bash
   mkdir new_service
   ```

2. **Add MODULE.bazel dependencies** (if needed):
   ```python
   bazel_dep(name = "rules_new_lang", version = "X.Y.Z")
   ```

3. **Create BUILD.bazel**:
   ```python
   # Define your targets here
   ```

4. **Update README.md**: Add the new service to the table

5. **Update CI**: Add build step in `.github/workflows/ci.yaml`

### Improving Documentation

Documentation lives in:
- `README.md` - Main overview
- `docs/quickstart.md` - Getting started guide
- `docs/concepts.md` - Deep dive into concepts
- `docs/dependencies.md` - Managing dependencies

Guidelines:
- Keep explanations beginner-friendly
- Include code examples
- Add links to official documentation
- Test all commands before documenting

### Testing Your Changes

Before submitting:

1. **Validate the build**:
   ```bash
   bazel clean --expunge
   bazel build //...
   ```

2. **Run tests**:
   ```bash
   bazel test //...
   ```

3. **Run validation script**:
   ```bash
   ./validate.sh
   ```

4. **Check formatting**:
   ```bash
   # For Bazel files (if buildifier is installed)
   bazel run //:buildifier
   
   # For Go files
   cd go_service && gofmt -w *.go
   
   # For Python files
   cd py_service && black *.py
   ```

## Commit Guidelines

- Use clear, descriptive commit messages
- Reference issues if applicable: `Fix #123: Add Rust service example`
- Keep commits focused and atomic

Example:
```
Add Rust service example with BUILD.bazel

- Create rust_service directory
- Add rust_binary target
- Include basic HTTP server example
- Update README with Rust instructions

Fixes #42
```

## Pull Request Process

1. **Update documentation** if you changed functionality
2. **Add tests** for new features
3. **Ensure all tests pass**: Run `bazel test //...`
4. **Update README** if you added a new service or major feature
5. Submit your PR with a clear description

### PR Description Template

```markdown
## What does this PR do?

Brief description of changes.

## Why is this needed?

Explain the motivation.

## How to test?

Steps to verify the changes work.

## Checklist

- [ ] Code builds successfully (`bazel build //...`)
- [ ] All tests pass (`bazel test //...`)
- [ ] Documentation updated
- [ ] Commit messages are clear
```

## Code Style

### Bazel Files (BUILD.bazel, MODULE.bazel)

- Use 4 spaces for indentation
- Add comments explaining each target
- Group related targets together
- Keep dependencies sorted alphabetically

Example:
```python
# Good
py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = [
        requirement("fastapi"),
        requirement("uvicorn"),
    ],
)

# Not ideal (missing comments, unsorted)
py_binary(
    name = "server",
    srcs = ["main.py"],
    deps = [requirement("uvicorn"), requirement("fastapi")],
)
```

### Python Code

- Follow PEP 8
- Use type hints where appropriate
- Add docstrings to functions
- Keep functions small and focused

### Go Code

- Follow Go conventions
- Use `gofmt` for formatting
- Add comments to exported functions
- Write idiomatic Go

## Questions?

Feel free to:
- Open an issue for discussion
- Ask questions in your PR
- Check existing issues for similar topics

## License

By contributing, you agree that your contributions will be licensed under the same license as this project (MIT License).

## Code of Conduct

Be respectful and constructive in all interactions. We're all here to learn and improve together! 🤝
