# 🌍 Global Data Sources - Multi-Source Support

## Overview

The Economic Dashboard now supports fetching data from multiple global sources beyond FRED, enabling truly international macroeconomic analysis.

## Supported Data Sources

### 1. FRED (Federal Reserve Economic Data) 🇺🇸

**Status:** ✅ Fully Implemented

**Coverage:** United States economic indicators

**API:** https://fred.stlouisfed.org/docs/api/

**Usage:**
```bash
# Use series code directly (no prefix)
GDPC1      # US Real GDP
CPIAUCSL   # US CPI
UNRATE     # US Unemployment Rate
```

**API Key:** Required (free)
```bash
export FRED_API_KEY=your_key_here
```

**Advantages:**
- ✅ High-frequency data (daily, monthly, quarterly)
- ✅ Long historical coverage (often back to 1940s)
- ✅ Excellent data quality and documentation
- ✅ Free tier generous (120 requests/minute)

---

### 2. World Bank WDI (World Development Indicators) 🌍

**Status:** ✅ Implemented

**Coverage:** 200+ countries and territories, 1960-present

**API:** https://datahelpdesk.worldbank.org/knowledgebase/articles/889392

**Usage:**
```bash
# Format: WB:INDICATOR:COUNTRY
WB:NY.GDP.MKTP.KD:USA  # US Real GDP
WB:NY.GDP.MKTP.KD:CHN  # China Real GDP
WB:NY.GDP.MKTP.KD:EMU  # Euro Area Real GDP
WB:NY.GDP.MKTP.KD:WLD  # World Real GDP
```

**API Key:** Not required (public API)

**Advantages:**
- ✅ True global coverage
- ✅ Standardized cross-country data
- ✅ No API key required
- ✅ Authoritative source for development indicators

**Limitations:**
- ⚠️ Annual frequency only
- ⚠️ Data published with 1-2 year lag
- ⚠️ Limited to annual aggregates

**Common Indicators:**

| Code | Name | Coverage |
|------|------|----------|
| `NY.GDP.MKTP.KD` | GDP (constant 2015 US$) | 1960→ |
| `NY.GDP.PCAP.KD` | GDP per capita (constant 2015 US$) | 1960→ |
| `FP.CPI.TOTL` | Consumer price index (2010=100) | 1960→ |
| `SL.UEM.TOTL.NE.ZS` | Unemployment, total (% of labor force) | 1991→ |

**Country Codes:**

| Code | Country/Region |
|------|----------------|
| `USA` | United States |
| `CHN` | China |
| `JPN` | Japan |
| `DEU` | Germany |
| `GBR` | United Kingdom |
| `FRA` | France |
| `IND` | India |
| `BRA` | Brazil |
| `CAN` | Canada |
| `AUS` | Australia |
| `EMU` | Euro Area |
| `WLD` | World |
| `EAS` | East Asia & Pacific |
| `ECS` | Europe & Central Asia |

---

### 3. OECD API 🏛️

**Status:** 🚧 Planned (not yet implemented)

**Coverage:** 38 OECD member countries + partners

**API:** https://data.oecd.org/api/

**Usage (planned):**
```bash
# Format: OECD:DATASET.DIMENSIONS
OECD:MEI.LRUNTTTT.USA.M     # US Unemployment Rate
OECD:MEI.CPALTT01.USA.M     # US CPI
OECD:KEI.CLI.OECD.M          # OECD Composite Leading Indicator
```

**Advantages:**
- ✅ High-quality harmonized data
- ✅ Monthly and quarterly frequencies
- ✅ Advanced economic indicators (CLI, PMI)
- ✅ No API key required

**Coming Soon:** Implementation planned for next phase

---

### 4. Trading Economics 📊

**Status:** 🚧 Planned (not yet implemented)

**Coverage:** 196 countries, 300,000+ indicators

**API:** https://tradingeconomics.com/api

**Usage (planned):**
```bash
# Format: TE:COUNTRY:INDICATOR
TE:USA:PMI              # US Manufacturing PMI
TE:CHN:M2               # China M2 Money Supply
TE:EMU:PMI              # Eurozone Manufacturing PMI
```

**API Key:** Required (paid, starts at $50/month)

**Advantages:**
- ✅ Comprehensive global coverage
- ✅ Real-time and historical data
- ✅ Includes PMI, sentiment, forecasts
- ✅ High-frequency updates

**Limitations:**
- ⚠️ Requires paid subscription
- ⚠️ Higher cost for commercial use

---

### 5. CPB Netherlands World Trade Monitor 🚢

**Status:** 🚧 Planned (not yet implemented)

**Coverage:** Global trade volumes and prices

**Data:** https://www.cpb.nl/en/worldtrademonitor

**Usage (planned):**
```bash
# Format: CPB:INDICATOR
CPB:WORLD_TRADE_VOLUME     # World Trade Volume Index
CPB:WORLD_TRADE_PRICE      # World Trade Price Index
```

**API Key:** Not required (public data download)

**Advantages:**
- ✅ Authoritative global trade data
- ✅ Monthly frequency
- ✅ Free access
- ✅ Used by IMF, World Bank, central banks

**Limitations:**
- ⚠️ No formal API (CSV/Excel downloads)
- ⚠️ Requires custom parsing

---

## Implementation Architecture

### Multi-Source Factory Pattern

```go
// DataSource interface - all sources implement this
type DataSource interface {
    FetchSeries(seriesID string) (*Series, []Observation, error)
    GetSourceName() string
}

// Factory creates appropriate client based on series ID prefix
func DataSourceFactory(seriesID string, fredAPIKey string) (DataSource, string, error) {
    if seriesID[:3] == "WB:" {
        return NewWorldBankClient(), seriesID, nil
    } else if seriesID[:5] == "OECD:" {
        return NewOECDClient(), seriesID, nil
    } else {
        // Default to FRED
        return NewFREDClient(fredAPIKey), seriesID, nil
    }
}
```

### Automatic Source Detection

The system automatically detects the data source from the series ID prefix:

- **No prefix** → FRED
- **`WB:`** → World Bank
- **`OECD:`** → OECD
- **`TE:`** → Trading Economics (planned)
- **`CPB:`** → CPB Netherlands (planned)

---

## Configuration

### Series Configuration in `main.go`

```go
defaultSeries = []string{
    // FRED (no prefix)
    "GDPC1",                    // US Real GDP
    "CPIAUCSL",                 // US CPI
    
    // World Bank (WB: prefix)
    "WB:NY.GDP.MKTP.KD:USA",    // US Real GDP (World Bank)
    "WB:NY.GDP.MKTP.KD:CHN",    // China Real GDP
    "WB:NY.GDP.MKTP.KD:EMU",    // Euro Area Real GDP
    "WB:NY.GDP.MKTP.KD:WLD",    // World Real GDP
    
    // OECD (planned - OECD: prefix)
    // "OECD:MEI.LRUNTTTT.USA.M", // US Unemployment
    
    // Trading Economics (planned - TE: prefix)
    // "TE:USA:PMI",               // US Manufacturing PMI
}
```

---

## Usage Examples

### Fetching Mixed-Source Data

```bash
# Set FRED API key (World Bank doesn't need one)
export FRED_API_KEY=your_key_here

# Fetch data from all sources
bazel run //go_fetch:refresh

# The system automatically:
# 1. Detects source from series ID prefix
# 2. Creates appropriate client
# 3. Fetches and stores data
# 4. Logs source in database
```

### Querying Multi-Source Data

```bash
# Open database
sqlite3 ./data/econ.db

# See all sources
SELECT DISTINCT source FROM series;
# Returns: FRED, World Bank

# Get all World Bank series
SELECT * FROM series WHERE source = 'World Bank';

# Get US GDP from different sources
SELECT s.source, s.name, o.date, o.value 
FROM series s 
JOIN observations o ON s.id = o.series_id
WHERE s.name LIKE '%GDP%' AND s.name LIKE '%United States%'
ORDER BY s.source, o.date DESC
LIMIT 20;
```

---

## Data Quality Comparison

| Source | Frequency | Timeliness | Coverage | Quality | Cost |
|--------|-----------|------------|----------|---------|------|
| **FRED** | Daily-Quarterly | Real-time | US-focused | ⭐⭐⭐⭐⭐ | Free |
| **World Bank** | Annual | 1-2 year lag | Global | ⭐⭐⭐⭐ | Free |
| **OECD** | Monthly-Quarterly | 1-3 month lag | OECD countries | ⭐⭐⭐⭐⭐ | Free |
| **Trading Economics** | Real-time | Real-time | Global | ⭐⭐⭐⭐ | Paid |
| **CPB** | Monthly | 1-2 month lag | Global trade | ⭐⭐⭐⭐ | Free |

---

## Best Practices

### 1. Source Selection

**Use FRED for:**
- ✅ US economic data
- ✅ High-frequency analysis (daily/monthly)
- ✅ Real-time monitoring
- ✅ Financial market data

**Use World Bank for:**
- ✅ Cross-country comparisons
- ✅ Long-term trends (annual)
- ✅ Development indicators
- ✅ Global aggregates

**Use OECD for (when implemented):**
- ✅ Harmonized international data
- ✅ Advanced indicators (CLI, PMI)
- ✅ OECD member countries
- ✅ Monthly economic monitoring

### 2. Frequency Alignment

When mixing sources:
- Annual data (World Bank) → Use end-of-year for alignment
- Quarterly data (FRED GDPC1) → Aggregate to annual for comparisons
- Monthly data (FRED, OECD) → Can be aggregated up

### 3. Data Validation

Always check:
1. **Date ranges** - ensure overlap for analysis
2. **Units** - convert if needed (e.g., constant prices year)
3. **Coverage** - verify no missing periods
4. **Source documentation** - understand methodology

---

## Future Roadmap

### Phase 1: ✅ Complete
- ✅ FRED integration
- ✅ World Bank WDI integration
- ✅ Multi-source factory pattern
- ✅ Automatic source detection

### Phase 2: 🚧 In Progress
- 🚧 OECD API integration
- 🚧 Enhanced metadata with source attribution
- 🚧 Cross-source validation

### Phase 3: 📋 Planned
- 📋 Trading Economics integration
- 📋 CPB Netherlands integration
- 📋 IMF IFS integration
- 📋 Eurostat integration

### Phase 4: 💡 Future
- 💡 Real-time data streaming
- 💡 Automatic unit conversion
- 💡 Data quality scoring
- 💡 Multi-source consensus indicators

---

## Troubleshooting

### World Bank API Issues

**Error: "No data available"**
- Check country code (use 3-letter ISO codes)
- Verify indicator code at https://data.worldbank.org/
- Some indicators have limited country coverage

**Error: "API request timeout"**
- World Bank API can be slow for large date ranges
- Try limiting date range: `?date=2000:2024`
- Check World Bank API status

### Source Detection Issues

**Error: "Unknown source prefix"**
- Ensure series ID format is correct
- Check prefix spelling (case-sensitive)
- Valid prefixes: `WB:`, `OECD:`, `TE:`, `CPB:`

---

## Resources

### API Documentation
- **FRED:** https://fred.stlouisfed.org/docs/api/
- **World Bank:** https://datahelpdesk.worldbank.org/knowledgebase/articles/889392
- **OECD:** https://data.oecd.org/api/sdmx-json-documentation/
- **Trading Economics:** https://docs.tradingeconomics.com/

### Data Catalogs
- **World Bank Indicators:** https://data.worldbank.org/indicator
- **OECD Data:** https://data.oecd.org/
- **CPB World Trade Monitor:** https://www.cpb.nl/en/worldtrademonitor

---

**Last Updated:** 2025-11-06  
**Version:** 2.0.0  
**Status:** 🌍 Multi-Source Enabled
