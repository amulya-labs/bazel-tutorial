# 🚀 Economic Dashboard Quick Start

Get the Economic Indicators Dashboard running in 5 minutes!

## Prerequisites

- **Bazel 6.4.0+** - [Install Guide](https://bazel.build/install)
- **Node.js 18+** - For the React UI (download from [nodejs.org](https://nodejs.org))
- **FRED API Key** - Free from [FRED](https://fred.stlouisfed.org/docs/api/api_key.html)

## 1️⃣ Clone & Setup

```bash
# Clone the repository
git clone https://github.com/rrl-personal-projects/bazel-tutorial.git
cd bazel-tutorial

# Install UI dependencies
cd web_ui
npm install
cd ..
```

## 2️⃣ Get Your FRED API Key

1. Visit https://fred.stlouisfed.org/
2. Create a free account (takes 1 minute)
3. Go to **My Account → API Keys**
4. Click **Request API Key**
5. Copy your API key

Set it in your environment:
```bash
export FRED_API_KEY=your_actual_api_key_here
```

## 3️⃣ Fetch Economic Data

```bash
# This downloads the latest economic indicators from FRED
bazel run //go_fetch:refresh
```

You should see output like:
```
Starting economic data refresh...
Fetching series: CPIAUCSL
Successfully stored 893 observations for CPIAUCSL
Fetching series: UNRATE
Successfully stored 947 observations for UNRATE
...
Refresh complete: Fetched 6/6 series successfully
```

## 4️⃣ Start the API Server

In terminal 1:
```bash
bazel run //py_api:server
```

The server starts on http://localhost:8000

Test it:
```bash
curl http://localhost:8000/api/health
```

## 5️⃣ Start the Web UI

In terminal 2:
```bash
cd web_ui
npm run dev
```

The UI starts on http://localhost:3000

**Open http://localhost:3000 in your browser!** 🎉

## 🎨 What You'll See

### Dashboard View
- 6 indicator cards showing latest values
- Month-over-Month (MoM) changes
- Year-over-Year (YoY) changes
- Click any card for details

### Series Detail View
- Interactive line chart
- Historical data
- Computed metrics

## 🔄 Updating Data

To refresh with the latest data:
```bash
bazel run //go_fetch:refresh
```

The API automatically picks up the new data!

## 🧪 Running Tests

```bash
# Run all tests
bazel test //...

# Run specific tests
bazel test //go_fetch:fred_test
bazel test //py_api:api_test
```

## 🛠️ Development Helpers

We provide a `dev.sh` script for common tasks:

```bash
# Install dependencies
./dev.sh setup

# Fetch data
./dev.sh fetch

# Start API
./dev.sh api

# Start UI
./dev.sh ui

# Build all
./dev.sh build

# Run tests
./dev.sh test
```

## 📖 Next Steps

- Read [ECONOMIC_DASHBOARD.md](ECONOMIC_DASHBOARD.md) for full documentation
- Learn about [adding new indicators](ECONOMIC_DASHBOARD.md#adding-a-new-economic-indicator)
- Explore [Bazel concepts](ECONOMIC_DASHBOARD.md#-bazel-concepts)
- Check out the [CI/CD setup](ECONOMIC_DASHBOARD.md#-cicd)

## ❓ Troubleshooting

### "Database not found"
**Solution:** Run the data fetcher first:
```bash
bazel run //go_fetch:refresh
```

### "FRED_API_KEY not set"
**Solution:** Set your API key:
```bash
export FRED_API_KEY=your_key_here
```

### "Port already in use"
**Solution:** Kill the existing process or use a different port:
```bash
# For API
PORT=8001 bazel run //py_api:server

# For UI (edit web_ui/vite.config.ts)
```

### "npm install fails"
**Solution:** Make sure you have Node.js 18+ installed:
```bash
node --version  # Should be v18.0.0 or higher
```

## 🆘 Need Help?

- Check [ECONOMIC_DASHBOARD.md](ECONOMIC_DASHBOARD.md#-troubleshooting) for more troubleshooting
- Open an issue on GitHub
- Review the inline comments in BUILD.bazel files

---

**Happy Building! 📊🚀**
