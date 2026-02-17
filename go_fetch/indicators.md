# 📊 High-ROI Global Economic Indicators - Data Dictionary

## Overview

This document defines the **10 core macroeconomic indicators** fetched and versioned for the Economic Dashboard. These indicators provide comprehensive coverage of **growth, inflation, liquidity, sentiment, and risk** with consistent data since ~1990.

## Design Principles

✅ **Economically interpretable and model-relevant**  
✅ **Consistent global coverage** since ~1990  
✅ **Small enough to store in Git** (<50 MB total)  
✅ **High signal-to-noise ratio**  

---

## Indicator Specifications

### 1️⃣ Real GDP (Constant Prices)
**Proxy For:** Growth / Economic Cycle

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `GDPC1` |
| **Name** | Real Gross Domestic Product |
| **Unit** | Billions of Chained 2012 Dollars |
| **Frequency** | Quarterly |
| **Coverage** | 1947 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Economic Analysis |
| **FRED URL** | https://fred.stlouisfed.org/series/GDPC1 |

**Description:** Measures the value of all goods and services produced in the U.S. economy, adjusted for inflation. Key indicator of economic growth and business cycles.

---

### 2️⃣ Consumer Price Index (CPI)
**Proxy For:** Inflation

#### Headline CPI
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `CPIAUCSL` |
| **Name** | Consumer Price Index for All Urban Consumers: All Items |
| **Unit** | Index 1982-1984=100 |
| **Frequency** | Monthly |
| **Coverage** | 1947 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Labor Statistics |
| **FRED URL** | https://fred.stlouisfed.org/series/CPIAUCSL |

**Description:** Measures average change in prices paid by urban consumers for a market basket of goods and services.

#### Core CPI
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `CPILFESL` |
| **Name** | Consumer Price Index for All Urban Consumers: All Items Less Food and Energy |
| **Unit** | Index 1982-1984=100 |
| **Frequency** | Monthly |
| **Coverage** | 1957 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Labor Statistics |
| **FRED URL** | https://fred.stlouisfed.org/series/CPILFESL |

**Description:** Core inflation measure excluding volatile food and energy components. Preferred for assessing underlying inflation trends.

---

### 3️⃣ Unemployment Rate
**Proxy For:** Labor Market Health

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `UNRATE` |
| **Name** | Unemployment Rate |
| **Unit** | Percent |
| **Frequency** | Monthly |
| **Coverage** | 1948 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Labor Statistics |
| **FRED URL** | https://fred.stlouisfed.org/series/UNRATE |

**Description:** Percentage of labor force that is jobless and actively seeking employment. Key indicator of labor market slack.

---

### 4️⃣ Federal Funds Rate
**Proxy For:** Monetary Policy Stance

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `FEDFUNDS` |
| **Name** | Federal Funds Effective Rate |
| **Unit** | Percent |
| **Frequency** | Monthly (daily data available) |
| **Coverage** | 1954 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | Board of Governors of the Federal Reserve System |
| **FRED URL** | https://fred.stlouisfed.org/series/FEDFUNDS |

**Description:** Interest rate at which depository institutions lend reserve balances to other institutions overnight. Primary tool of U.S. monetary policy.

---

### 5️⃣ Treasury Yield Curve (10Y - 2Y Spread)
**Proxy For:** Recession Signal / Yield Curve Shape

#### 10Y-2Y Spread (Direct)
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `T10Y2Y` |
| **Name** | 10-Year Treasury Constant Maturity Minus 2-Year Treasury Constant Maturity |
| **Unit** | Percent |
| **Frequency** | Daily |
| **Coverage** | 1976 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | Federal Reserve Board |
| **FRED URL** | https://fred.stlouisfed.org/series/T10Y2Y |

**Description:** Spread between 10-year and 2-year Treasury yields. Negative values (inverted curve) historically precede recessions.

#### Supporting Series: 10-Year Treasury
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `DGS10` |
| **Name** | Market Yield on U.S. Treasury Securities at 10-Year Constant Maturity |
| **Unit** | Percent |
| **Frequency** | Daily |
| **Coverage** | 1962 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |

#### Supporting Series: 2-Year Treasury
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `DGS2` |
| **Name** | Market Yield on U.S. Treasury Securities at 2-Year Constant Maturity |
| **Unit** | Percent |
| **Frequency** | Daily |
| **Coverage** | 1976 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |

---

### 6️⃣ M2 Money Supply
**Proxy For:** Liquidity / Credit Conditions

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `M2SL` |
| **Name** | M2 Money Stock |
| **Unit** | Billions of Dollars |
| **Frequency** | Monthly |
| **Coverage** | 1959 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | Board of Governors of the Federal Reserve System |
| **FRED URL** | https://fred.stlouisfed.org/series/M2SL |

**Description:** Measure of money supply including cash, checking deposits, and easily convertible near money. Key indicator of liquidity and credit conditions.

---

### 7️⃣ Brent Crude Oil Price
**Proxy For:** Inflation Driver & Demand Proxy

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `POILBREUSDM` |
| **Name** | Global Price of Brent Crude |
| **Unit** | Dollars per Barrel |
| **Frequency** | Monthly |
| **Coverage** | 1987 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | International Monetary Fund via FRED |
| **FRED URL** | https://fred.stlouisfed.org/series/POILBREUSDM |

**Description:** Global benchmark oil price. Key driver of inflation and indicator of global economic demand.

---

### 8️⃣ Manufacturing Sector Activity
**Proxy For:** Business Confidence / Early Cycle Indicator

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `MANEMP` |
| **Name** | All Employees, Manufacturing |
| **Unit** | Thousands of Persons |
| **Frequency** | Monthly |
| **Coverage** | 1939 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Labor Statistics |
| **FRED URL** | https://fred.stlouisfed.org/series/MANEMP |

**Description:** Total manufacturing employment. Proxy for manufacturing activity and early-cycle business conditions.

**Note:** ISM Manufacturing PMI would be ideal but is not available in FRED. Manufacturing employment serves as a reliable alternative indicator.

---

### 9️⃣ Consumer Sentiment Index
**Proxy For:** Household Confidence

| Attribute | Value |
|-----------|-------|
| **FRED Code** | `UMCSENT` |
| **Name** | University of Michigan: Consumer Sentiment |
| **Unit** | Index 1966:Q1=100 |
| **Frequency** | Monthly |
| **Coverage** | 1978 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | University of Michigan via FRED |
| **FRED URL** | https://fred.stlouisfed.org/series/UMCSENT |

**Description:** Measures consumer confidence and expectations for the economy. Leading indicator of consumer spending.

---

### 🔟 Global Trade Volume
**Proxy For:** Global Demand Flow

#### U.S. Real Imports
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `IMPCH` |
| **Name** | Real Imports of Goods and Services |
| **Unit** | Billions of Chained 2012 Dollars |
| **Frequency** | Quarterly |
| **Coverage** | 1947 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Economic Analysis |
| **FRED URL** | https://fred.stlouisfed.org/series/IMPCH |

#### U.S. Real Exports
| Attribute | Value |
|-----------|-------|
| **FRED Code** | `EXPCH` |
| **Name** | Real Exports of Goods and Services |
| **Unit** | Billions of Chained 2012 Dollars |
| **Frequency** | Quarterly |
| **Coverage** | 1947 → Present |
| **30-Year Baseline** | ✅ Yes (1990→) |
| **Source** | U.S. Bureau of Economic Analysis |
| **FRED URL** | https://fred.stlouisfed.org/series/EXPCH |

**Description:** Real trade flows adjusted for inflation. Proxy for global trade volume and international demand.

**Note:** CPB Netherlands World Trade Monitor would be ideal for truly global trade but is not available in FRED. U.S. trade data serves as a reliable proxy for global trade patterns.

---

## Data Coverage Summary

| # | Indicator | FRED Code | Frequency | Start Year | 30Y Coverage |
|---|-----------|-----------|-----------|------------|--------------|
| 1 | Real GDP | GDPC1 | Quarterly | 1947 | ✅ |
| 2a | CPI (Headline) | CPIAUCSL | Monthly | 1947 | ✅ |
| 2b | Core CPI | CPILFESL | Monthly | 1957 | ✅ |
| 3 | Unemployment Rate | UNRATE | Monthly | 1948 | ✅ |
| 4 | Fed Funds Rate | FEDFUNDS | Monthly | 1954 | ✅ |
| 5a | 10Y-2Y Spread | T10Y2Y | Daily | 1976 | ✅ |
| 5b | 10Y Treasury | DGS10 | Daily | 1962 | ✅ |
| 5c | 2Y Treasury | DGS2 | Daily | 1976 | ✅ |
| 6 | M2 Money Supply | M2SL | Monthly | 1959 | ✅ |
| 7 | Brent Crude Oil | POILBREUSDM | Monthly | 1987 | ✅ |
| 8 | Manufacturing Emp. | MANEMP | Monthly | 1939 | ✅ |
| 9 | Consumer Sentiment | UMCSENT | Monthly | 1978 | ✅ |
| 10a | Real Imports | IMPCH | Quarterly | 1947 | ✅ |
| 10b | Real Exports | EXPCH | Quarterly | 1947 | ✅ |

**Total Series:** 14 FRED series covering 10 core economic indicators  
**Estimated Size:** ~5-10 MB total (well under 50 MB limit)

---

## Future Enhancements

### Additional Data Sources to Consider

1. **World Bank WDI** (World Development Indicators)
   - Global GDP data for multiple countries
   - API: https://datahelpdesk.worldbank.org/knowledgebase/articles/889392

2. **OECD API** (Organisation for Economic Co-operation and Development)
   - International unemployment, CPI data
   - API: https://data.oecd.org/api/

3. **Trading Economics API**
   - Manufacturing PMI for US, China, Eurozone
   - China M2 Money Supply
   - API: https://tradingeconomics.com/api

4. **CPB Netherlands**
   - World Trade Monitor (true global trade volume index)
   - Data: https://www.cpb.nl/en/worldtrademonitor

5. **S&P Global**
   - ISM Manufacturing PMI (official PMI data)
   - Requires license/subscription

### Implementation Notes

- Current implementation focuses on **FRED-only data** to simplify initial rollout
- All selected series have excellent coverage back to 1990 or earlier
- Additional sources can be added incrementally without disrupting existing data
- Consider implementing a multi-source adapter pattern for future expansions

---

## Usage

### Fetching Data

```bash
# Set your FRED API key
export FRED_API_KEY=your_key_here

# Fetch all indicators (default series)
bazel run //go_fetch:refresh

# Data is stored in ./data/econ.db
```

### Exporting to JSON

```bash
# Export for static hosting / Git versioning
bazel run //go_fetch:refresh -- --export-json=./docs/app/data

# Creates:
# - summary.json (latest values + deltas)
# - series-{CODE}.json (full time series per indicator)
```

### Querying Database

```bash
# Open SQLite database
sqlite3 ./data/econ.db

# List all series
SELECT * FROM series;

# Get latest observation for a series
SELECT * FROM observations 
WHERE series_id = 'GDPC1' 
ORDER BY date DESC 
LIMIT 1;

# Get 30-year history for CPI
SELECT * FROM observations 
WHERE series_id = 'CPIAUCSL' 
  AND date >= date('now', '-30 years')
ORDER BY date;
```

---

## References

- **FRED API Documentation:** https://fred.stlouisfed.org/docs/api/
- **FRED Series Search:** https://fred.stlouisfed.org/
- **Bazel Tutorial Repo:** https://github.com/amulya-labs/bazel-tutorial
- **Economic Dashboard Documentation:** /docs/economic-dashboard.md

---

**Last Updated:** 2025-11-06  
**Version:** 1.0.0  
**Maintained By:** Economic Dashboard Team
