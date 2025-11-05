# 📊 Economic Indicators Dashboard - Bazel Polyglot Tutorial

A comprehensive **Bazel-based polyglot application** that demonstrates building and orchestrating **Go**, **Python**, and **React** services in a single repository. This project fetches real economic data from the FRED API, serves it through a REST API, and visualizes it in a modern web dashboard.

## 🎯 What This Project Demonstrates

This is a **teaching repository** that shows:

- ✅ **Polyglot builds** with Bazel (Go + Python + React)
- ✅ **Incremental compilation** and caching
- ✅ **Hermetic testing** without network dependencies
- ✅ **Modern toolchain integration** (rules_go, rules_python, npm/vite)
- ✅ **CI/CD** with GitHub Actions and Bazel caching
- ✅ **Real-world architecture** patterns (data fetcher → API → UI)

---

## 🏗️ Architecture

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

### Components

1. **`go_fetch/`** - Go service that fetches economic indicators from FRED API
   - Fetches 6 core indicators (CPI, Unemployment, Fed Funds, etc.)
   - Stores data in SQLite database
   - Includes mock HTTP tests

2. **`py_api/`** - Python FastAPI REST service
   - Reads from SQLite database
   - Computes MoM/YoY deltas
   - Exposes normalized JSON endpoints
   - Includes fixture-based tests

3. **`web_ui/`** - React + TypeScript + Vite frontend
   - Dashboard view with indicator cards
   - Detailed series view with charts
   - Responsive design

---

## 🚀 Quick Start

### Prerequisites

1. **Install Bazel** (6.4.0 or later)
   ```bash
   # macOS
   brew install bazel
   
   # Linux
   # Download from https://github.com/bazelbuild/bazel/releases
   
   # Verify installation
   bazel --version
   ```

2. **Get a FRED API Key** (free)
   - Sign up at: https://fred.stlouisfed.org/docs/api/api_key.html
   - Set environment variable:
     ```bash
     export FRED_API_KEY=your_api_key_here
     ```

### Running the Dashboard

#### Option 1: Run All Components (Recommended)

```bash
# Start everything with one script
./run_all.sh
```

This will:
1. Fetch latest economic data from FRED
2. Start the Python API server on http://localhost:8000
3. Show instructions for starting the web UI

Then in another terminal:
```bash
cd web_ui
npm install  # First time only
npm run dev
```

Open http://localhost:3000 in your browser 🎉

#### Option 2: Run Components Individually

```bash
# 1. Fetch data
bazel run //go_fetch:refresh

# 2. Start API (in one terminal)
bazel run //py_api:server

# 3. Start UI (in another terminal)
cd web_ui && npm install && npm run dev
```

---

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

---

## 🏗️ Building

```bash
# Build everything
bazel build //...

# Build specific targets
bazel build //go_fetch:refresh
bazel build //py_api:server
bazel build //web_ui:build

# Query the build graph
bazel query //...
bazel query 'deps(//py_api:server)'
```

---

## 📊 Available Indicators

The dashboard tracks 6 core economic indicators from FRED:

| Code | Name | Unit |
|------|------|------|
| `CPIAUCSL` | Consumer Price Index (Headline) | Index 1982-1984=100 |
| `CPILFESL` | Core CPI | Index 1982-1984=100 |
| `UNRATE` | Unemployment Rate | Percent |
| `FEDFUNDS` | Federal Funds Rate | Percent |
| `DGS10` | 10-Year Treasury Yield | Percent |
| `M2SL` | M2 Money Stock | Billions of Dollars |

---

## 🔌 API Endpoints

### `GET /api/health`
Health check and service status

### `GET /api/econ/summary`
Summary of all indicators with latest values and deltas

### `GET /api/econ/series?code=<CODE>&range=<RANGE>`
Time series data for a specific indicator
- `code`: Series code (e.g., `CPIAUCSL`)
- `range`: Time range - `1y`, `5y`, or `max` (default)

### `GET /api/econ/meta`
Metadata about all available series

API documentation: http://localhost:8000/docs

---

## 🧱 Bazel Concepts

### MODULE.bazel

This file defines external dependencies using **Bzlmod** (Bazel's modern dependency system):

- **rules_go** - Go toolchain and build rules
- **rules_python** - Python toolchain and pip integration
- **aspect_rules_js** - Modern JavaScript/Node.js support

### BUILD.bazel Files

Each component has a `BUILD.bazel` file that defines:
- **Binaries** (`go_binary`, `py_binary`) - Executable programs
- **Libraries** (`go_library`, `py_library`) - Reusable code
- **Tests** (`go_test`, `py_test`) - Test targets

Example from `go_fetch/BUILD.bazel`:
```python
go_binary(
    name = "refresh",
    srcs = ["main.go"],
    deps = [":go_fetch_lib"],
)
```

### Targets and Labels

Targets are referenced with labels: `//package:target`
- `//go_fetch:refresh` - Go binary in go_fetch package
- `//py_api:server` - Python server binary
- `//...` - All targets in the repository

---

## 🔄 Incremental Builds

Bazel only rebuilds what changed. Try this:

```bash
# First build (downloads deps, compiles everything)
bazel build //...

# Edit go_fetch/main.go
# Only Go targets rebuild!
bazel build //...

# Edit py_api/handlers.py
# Only Python targets rebuild!
bazel build //...
```

---

## 🎓 How to Extend

### Adding a New Economic Indicator

1. **Update `go_fetch/main.go`** - Add series code to `defaultSeries`
   ```go
   defaultSeries = []string{
       "CPIAUCSL",
       "UNRATE",
       "GDPC1",  // ← Add new indicator
   }
   ```

2. **Refresh data**
   ```bash
   bazel run //go_fetch:refresh
   ```

3. **Indicator automatically appears in UI** - No frontend changes needed!

### Adding Go Dependencies

1. Add to `go_fetch/go.mod`:
   ```go
   require github.com/some/package v1.0.0
   ```

2. Run:
   ```bash
   cd go_fetch && go mod tidy
   bazel run //:gazelle
   ```

### Adding Python Dependencies

1. Add to `py_api/requirements.txt`:
   ```
   requests==2.31.0
   ```

2. Update MODULE.bazel pip.parse section

3. Rebuild:
   ```bash
   bazel build //py_api:server
   ```

---

## 🔍 Troubleshooting

### "Database not found"
**Solution:** Run the data fetcher first:
```bash
bazel run //go_fetch:refresh
```

### "FRED_API_KEY not set"
**Solution:** Get a free API key and set it:
```bash
export FRED_API_KEY=your_key_here
```

### "Build is slow"
**Solution:** First build downloads all dependencies. Subsequent builds are fast due to caching.

### "Tests failing"
**Solution:** Tests are hermetic and don't need network/API keys. If failing, check:
```bash
bazel test //... --test_output=errors
```

### "Web UI not loading"
**Solution:** Ensure API server is running:
```bash
bazel run //py_api:server  # In one terminal
cd web_ui && npm run dev   # In another terminal
```

---

## 📁 Repository Structure

```
bazel-tutorial/
├── MODULE.bazel              # Bazel dependencies (Bzlmod)
├── .bazelrc                  # Bazel configuration
├── .bazelversion             # Pin Bazel version
├── run_all.sh               # Convenience script to start all services
│
├── go_fetch/                # Go data fetcher service
│   ├── main.go              # CLI entry point
│   ├── fred.go              # FRED API client
│   ├── store_sqlite.go      # SQLite storage
│   ├── fred_test.go         # Tests with mocks
│   ├── go.mod               # Go dependencies
│   └── BUILD.bazel          # Bazel build config
│
├── py_api/                  # Python FastAPI service
│   ├── main.py              # FastAPI server
│   ├── handlers.py          # API endpoints
│   ├── db.py                # Database layer
│   ├── test_api.py          # Tests with fixtures
│   ├── requirements.txt     # Python dependencies
│   └── BUILD.bazel          # Bazel build config
│
├── web_ui/                  # React frontend
│   ├── src/
│   │   ├── App.tsx          # Main app component
│   │   ├── components/
│   │   │   ├── IndicatorCard.tsx
│   │   │   └── LineChart.tsx
│   │   └── main.tsx
│   ├── package.json         # npm dependencies
│   ├── vite.config.ts       # Vite configuration
│   └── BUILD.bazel          # Bazel build config
│
├── data/                    # SQLite database (gitignored)
│   └── econ.db              # Created by go_fetch
│
└── .github/workflows/
    └── ci.yaml              # CI/CD pipeline
```

---

## 🔐 Configuration

### Environment Variables

- `FRED_API_KEY` - Required for data fetching
- `ECON_DB_PATH` - Database path (default: `./data/econ.db`)
- `PORT` - API server port (default: `8000`)

### Files

- `.env` - Local environment variables (create yourself, not committed)
- `.bazelrc` - Bazel configuration (caching, output options)

---

## 🚢 CI/CD

The project includes GitHub Actions workflow (`.github/workflows/ci.yaml`) that:

1. ✅ Sets up Bazel and caches dependencies
2. ✅ Builds all targets (`bazel build //...`)
3. ✅ Runs all tests (`bazel test //...`)
4. ✅ Uses Bazel cache for fast builds

Bazel's hermetic builds ensure CI runs are deterministic and cacheable.

---

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

---

## 📚 Learning Resources

### Bazel
- [Official Bazel Docs](https://bazel.build/)
- [Bzlmod Guide](https://bazel.build/external/overview#bzlmod)
- [Bazel Query](https://bazel.build/query/guide)

### Rules
- [rules_go](https://github.com/bazelbuild/rules_go)
- [rules_python](https://github.com/bazelbuild/rules_python)
- [aspect_rules_js](https://github.com/aspect-build/rules_js)

### Economic Data
- [FRED API Documentation](https://fred.stlouisfed.org/docs/api/)

---

## 🤝 Contributing

This is a tutorial project. Contributions welcome!

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `bazel test //...` passes
5. Submit a pull request

---

## 📄 License

MIT License - See LICENSE file for details

---

## 🎉 Acknowledgments

- **FRED** - Federal Reserve Economic Data
- **Bazel** - Google's build system
- **FastAPI** - Modern Python web framework
- **React** - UI library
- **Vite** - Fast build tool

---

**Happy Building! 🚀**

For questions or issues, please open a GitHub issue.
