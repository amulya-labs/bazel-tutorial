# 📊 Go Fetch - Economic Data Fetcher

## Overview

The `go_fetch` service fetches high-ROI global economic indicators from **multiple data sources** (FRED, World Bank, OECD, and more) and stores them in a SQLite database. It supports 30-year baselines for comprehensive macroeconomic analysis.

## Features

✅ **Multi-source support** - FRED, World Bank WDI, OECD (coming soon)
✅ **24 series** covering 10 core economic indicators (11 FRED + 13 World Bank)
✅ **Global coverage** - US, China, Euro Area, and World GDP  
✅ **30-year baseline** coverage (≈1990 → present)  
✅ **SQLite storage** for fast local queries  
✅ **JSON export** for static hosting and Git versioning  
✅ **Rich metadata** with categories, descriptions, and proxy mappings  
✅ **Automated deltas** (MoM and YoY calculations)  

## Supported Data Sources

### 🇺🇸 FRED (Federal Reserve Economic Data)
- **Status:** ✅ Fully Implemented
- **Coverage:** United States
- **API Key:** Required (free)
- **Series:** Use code directly (e.g., `GDPC1`)

### 🌍 World Bank WDI (World Development Indicators)
- **Status:** ✅ Implemented
- **Coverage:** 200+ countries, global aggregates
- **API Key:** Not required
- **Series:** Prefix with `WB:` (e.g., `WB:NY.GDP.MKTP.KD:USA`)

### 🏛️ OECD API
- **Status:** 🚧 Coming Soon
- **Coverage:** OECD member countries
- **Series:** Will use `OECD:` prefix

**See [global-sources.md](./global-sources.md) for complete documentation**

## Quick Start

### Prerequisites

1. **FRED API Key** (free) - Get one at: https://fred.stlouisfed.org/docs/api/api_key.html
2. **Bazel** - For building the service

### Fetch Data

```bash
# Set your FRED API key (World Bank doesn't need one)
export FRED_API_KEY=your_key_here

# Fetch all indicators from all sources
bazel run //go_fetch:refresh

# Data stored in ./data/econ.db
```

### Export to JSON

```bash
# Export for static hosting / version control
bazel run //go_fetch:refresh -- --export-json=./docs/app/data

# Creates:
# - metadata.json (indicator specifications)
# - summary.json (latest values + deltas)
# - series-{CODE}.json (full time series per indicator)
```

### Export Only (Without Fetching)

```bash
# Export existing data without requiring FRED_API_KEY
bazel run //go_fetch:refresh -- --export-json=./docs/app/data
```

## Indicators

### 10 Core Economic Indicators

The service fetches 11 FRED series covering 10 high-ROI macro indicators:

| # | Indicator | FRED Codes | Frequency | Coverage |
|---|-----------|------------|-----------|----------|
| 1️⃣ | Real GDP | `GDPC1` | Quarterly | 1947→ |
| 2️⃣ | CPI (Inflation) | `CPIAUCSL`, `CPILFESL` | Monthly | 1947→ |
| 3️⃣ | Unemployment Rate | `UNRATE` | Monthly | 1948→ |
| 4️⃣ | Fed Funds Rate | `FEDFUNDS` | Monthly | 1954→ |
| 5️⃣ | Treasury Yield Curve | `T10Y2Y` | Daily | 1976→ |
| 6️⃣ | M2 Money Supply | `M2SL` | Monthly | 1959→ |
| 7️⃣ | Brent Crude Oil Price | `POILBREUSDM` | Monthly | 1987→ |
| 8️⃣ | Industrial Production | `IPMAN` | Monthly | 1972→ |
| 9️⃣ | Consumer Sentiment | `UMCSENT` | Monthly | 1978→ |
| 🔟 | Inflation Expectations | `T5YIE` | Daily | 2003→ |
| 1️⃣1️⃣ | Financial Conditions | `NFCI` | Weekly | 1971→ |

**Full specifications:** See [indicators.md](./indicators.md)

### Categories

Indicators are organized into economic categories:

- **Growth:** Real GDP
- **Inflation:** CPI (Headline, Core), Oil Prices
- **Labor Market:** Unemployment Rate
- **Monetary Policy:** Fed Funds Rate
- **Yield Curve:** 10Y-2Y Spread, Treasury Yields
- **Liquidity:** M2 Money Supply
- **Manufacturing:** Employment
- **Sentiment:** Consumer Confidence
- **Trade:** Imports/Exports

## Architecture

```
┌─────────────────┐
│   FRED API      │
│  (14 series)    │
└────────┬────────┘
         │ HTTP/JSON
         ↓
┌─────────────────┐
│   fred.go       │
│ (API client)    │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│ store_sqlite.go │
│  (SQLite DB)    │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│   export.go     │
│ (JSON generator)│
└─────────────────┘
```

### Files

- **`main.go`** - CLI entry point and orchestration
- **`fred.go`** - FRED API client with HTTP interface for testing
- **`store_sqlite.go`** - SQLite storage layer
- **`export.go`** - JSON export with metadata and deltas
- **`metadata.go`** - Indicator metadata and categorization
- **`indicators.md`** - Complete data dictionary (12KB)
- **`fred_test.go`** - Tests with mock HTTP client

## Database Schema

### Tables

#### `series`
| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT PRIMARY KEY | FRED series code (e.g., "GDPC1") |
| `name` | TEXT NOT NULL | Full indicator name |
| `unit` | TEXT | Unit of measurement |
| `source` | TEXT | Data source (default: "FRED") |

#### `observations`
| Column | Type | Description |
|--------|------|-------------|
| `series_id` | TEXT NOT NULL | Foreign key to series.id |
| `date` | TEXT NOT NULL | ISO date (YYYY-MM-DD) |
| `value` | REAL NOT NULL | Observation value |

**Primary Key:** `(series_id, date)`

#### `refresh_log`
| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PRIMARY KEY | Auto-increment |
| `source` | TEXT NOT NULL | Data source name |
| `started_at` | TEXT NOT NULL | Refresh start timestamp |
| `finished_at` | TEXT | Refresh completion timestamp |
| `ok` | INTEGER | Success flag (0/1) |
| `message` | TEXT | Status message |

## JSON Exports

### `metadata.json`

Rich metadata about all indicators:

```json
{
  "indicators": {
    "GDPC1": {
      "code": "GDPC1",
      "name": "Real Gross Domestic Product",
      "unit": "Billions of Chained 2012 Dollars",
      "frequency": "Quarterly",
      "category": "Growth",
      "description": "Real GDP measures economic output adjusted for inflation",
      "start_year": 1947,
      "proxy_for": "Economic Growth / Business Cycle",
      "fredurl": "https://fred.stlouisfed.org/series/GDPC1"
    }
  },
  "categories": {
    "Growth": ["GDPC1"],
    "Inflation": ["CPIAUCSL", "CPILFESL", "POILBREUSDM"]
  },
  "core_series": ["GDPC1", "CPIAUCSL", ...],
  "version": "1.0.0",
  "generated": "2025-11-06T22:30:00Z"
}
```

### `summary.json`

Latest values with deltas for all indicators:

```json
{
  "indicators": [
    {
      "code": "GDPC1",
      "name": "Real Gross Domestic Product",
      "unit": "Billions of Chained 2012 Dollars",
      "source": "FRED",
      "category": "Growth",
      "proxy_for": "Economic Growth / Business Cycle",
      "last_updated": "2024-07-01",
      "value": 22914.517,
      "delta_mom": 0.7,
      "delta_yoy": 2.5
    }
  ],
  "count": 14,
  "last_refresh": "2025-11-06T22:30:00Z"
}
```

### `series-{CODE}.json`

Complete time series for a specific indicator:

```json
{
  "code": "GDPC1",
  "name": "Real Gross Domestic Product",
  "unit": "Billions of Chained 2012 Dollars",
  "source": "FRED",
  "last_updated": "2024-07-01",
  "latest_value": 22914.517,
  "delta_mom": 0.7,
  "delta_yoy": 2.5,
  "observations": [
    {"date": "1947-01-01", "value": 2033.061},
    {"date": "1947-04-01", "value": 2027.639},
    ...
  ],
  "count": 311,
  "range": "max"
}
```

## CLI Usage

### Options

```
-api-key string
    FRED API key (default: $FRED_API_KEY)
    
-db string
    Path to SQLite database (default: "./data/econ.db")
    
-export-json string
    Export data to JSON files in specified directory
```

### Examples

```bash
# Fetch data with custom database path
bazel run //go_fetch:refresh -- -db=./custom/path/econ.db

# Fetch and export in one command
bazel run //go_fetch:refresh -- -export-json=./output

# Export only (no fetching)
bazel run //go_fetch:refresh -- -export-json=./output
```

## Testing

The service includes comprehensive tests with mock HTTP clients:

```bash
# Run tests
bazel test //go_fetch:fred_test

# Run tests with verbose output
bazel test //go_fetch:fred_test --test_output=all
```

**Test Coverage:**
- ✅ FRED API client (with mock HTTP responses)
- ✅ Series metadata fetching
- ✅ Observations parsing
- ✅ Error handling

## Development

### Adding New Indicators

1. **Find FRED series code** at https://fred.stlouisfed.org/

2. **Add to `main.go`:**
   ```go
   defaultSeries = []string{
       // ... existing series ...
       "NEWSERIES",  // Description
   }
   ```

3. **Add metadata to `metadata.go`:**
   ```go
   "NEWSERIES": {
       Code:        "NEWSERIES",
       Name:        "New Economic Indicator",
       Unit:        "Units",
       Frequency:   "Monthly",
       Category:    "Category Name",
       Description: "What this measures",
       StartYear:   1990,
       ProxyFor:    "What it represents",
       FREDURL:     "https://fred.stlouisfed.org/series/NEWSERIES",
   }
   ```

4. **Fetch data:**
   ```bash
   bazel run //go_fetch:refresh
   ```

### Code Style

- **Go formatting:** `gofmt`
- **Comments:** Document all exported types and functions
- **Error handling:** Always wrap errors with context
- **Testing:** Add tests for new functionality

## Data Size Estimation

Current setup:
- **14 series** × ~300-600 observations each
- **~4,200-8,400 total observations**
- **SQLite database:** ~1-2 MB
- **JSON exports:** ~5-10 MB total

Well under the 50 MB Git storage target! 🎉

## Future Enhancements

### Additional Data Sources

- **World Bank WDI** - Global GDP, multi-country data
- **OECD API** - International unemployment, CPI
- **Trading Economics** - Manufacturing PMI, China M2
- **CPB Netherlands** - World Trade Monitor
- **S&P Global** - Official ISM PMI data

### Multi-Region Support

Expand to cover:
- **United States** (current)
- **European Union**
- **China**
- **Global aggregates**

See [indicators.md](./indicators.md) for detailed roadmap.

## Troubleshooting

### "FRED_API_KEY not set"

**Solution:** Get a free API key and set environment variable:
```bash
export FRED_API_KEY=your_key_here
```

### "Database not found"

**Solution:** The database is created automatically. Ensure `./data/` directory is writable.

### "API rate limit exceeded"

**Solution:** FRED allows 120 requests/minute. Space out requests or contact FRED for higher limits.

### "Series not found"

**Solution:** Verify the series code exists at https://fred.stlouisfed.org/

## Resources

- **FRED API Documentation:** https://fred.stlouisfed.org/docs/api/
- **FRED Series Search:** https://fred.stlouisfed.org/
- **Indicator Specifications:** [indicators.md](./indicators.md)
- **Project Documentation:** [../docs/economic-dashboard.md](../docs/economic-dashboard.md)

## License

MIT License - See ../LICENSE for details

---

**Maintained By:** Economic Dashboard Team  
**Last Updated:** 2025-11-06  
**Version:** 1.0.0
