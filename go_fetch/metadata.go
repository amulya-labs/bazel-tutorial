// Package main - Indicator configuration for high-ROI macro indicators
package main

// defaultSeries lists the High-ROI Global Economic Indicators to fetch (30-Year Baseline)
// These indicators provide comprehensive coverage of growth, inflation, liquidity,
// sentiment, and risk since ~1990.
//
// Supports multiple data sources:
// - FRED series: use code directly (e.g., "GDPC1")
// - World Bank: prefix with "WB:" (e.g., "WB:NY.GDP.MKTP.KD:USA")
// - OECD: prefix with "OECD:" (e.g., "OECD:MEI.LRUNTTTT.USA.M")
var defaultSeries = []string{
	// ===========================================
	// U.S. INDICATORS (FRED)
	// ===========================================

	// 1️⃣ Real GDP (constant prices) - Growth / cycle
	"GDPC1", // US Real GDP (Billions of Chained 2012 Dollars), Quarterly, 1947→ [FRED]

	// 2️⃣ CPI (Consumer Price Index) - Inflation
	"CPIAUCSL", // US CPI All Urban Consumers (Index 1982-84=100), Monthly, 1947→ [FRED]
	"CPILFESL", // US Core CPI Less Food & Energy (Index 1982-84=100), Monthly, 1957→ [FRED]

	// 3️⃣ Unemployment Rate - Labor market
	"UNRATE", // US Unemployment Rate (Percent), Monthly, 1948→ [FRED]

	// 4️⃣ Fed Funds Rate - Monetary policy stance
	"FEDFUNDS", // Federal Funds Effective Rate (Percent), Monthly, 1954→ [FRED]

	// 5️⃣ 10Y – 2Y Treasury Spread - Recession signal / yield curve
	"T10Y2Y", // 10-Year Treasury Minus 2-Year Treasury (Percent), Daily, 1976→ [FRED]

	// 6️⃣ M2 Money Supply - Liquidity / credit conditions
	"M2SL", // US M2 Money Stock (Billions of Dollars), Monthly, 1959→ [FRED]

	// 7️⃣ Brent Crude Oil Price - Inflation driver & demand proxy
	"POILBREUSDM", // Global Price of Brent Crude (Dollars per Barrel), Monthly, 1987→ [FRED]

	// 8️⃣ Industrial Production: Manufacturing - Manufacturing activity
	"IPMAN", // Industrial Production: Manufacturing NAICS (Index 2017=100), Monthly, 1972→ [FRED]

	// 9️⃣ Consumer Sentiment Index - Household confidence
	"UMCSENT", // University of Michigan Consumer Sentiment (Index 1966:Q1=100), Monthly, 1978→ [FRED]

	// 🔟 Inflation Expectations - Market-based inflation forecast
	"T5YIE", // 5-Year Breakeven Inflation Rate (Percent), Daily, 2003→ [FRED]

	// 1️⃣1️⃣ Financial Conditions Index - Financial stress/tightness
	"NFCI", // Chicago Fed National Financial Conditions Index, Weekly, 1971→ [FRED]

	// ===========================================
	// GLOBAL INDICATORS (World Bank)
	// ===========================================

	// 🌍 Global GDP from World Bank
	"WB:NY.GDP.MKTP.KD:WLD", // World Real GDP (constant 2015 USD), Annual, 1960→
	"WB:NY.GDP.MKTP.KD:USA", // US Real GDP from World Bank, Annual, 1960→
	"WB:NY.GDP.MKTP.KD:CHN", // China Real GDP, Annual, 1960→
	"WB:NY.GDP.MKTP.KD:EMU", // Euro Area Real GDP, Annual, 1960→

	// 🌍 Global Inflation from World Bank
	"WB:FP.CPI.TOTL.ZG:WLD", // Global CPI (inflation, annual %), Annual, 1960→
	"WB:FP.CPI.TOTL.ZG:CHN", // China CPI inflation, Annual, 1960→
	"WB:FP.CPI.TOTL.ZG:EMU", // Euro Area CPI inflation (proxy for HICP), Annual, 1960→

	// 🌍 Global Money Supply & Credit
	"WB:FM.LBL.BMNY.CN:CHN",    // China Broad Money (M2), Annual, 1960→
	"WB:FS.AST.PRVT.GD.ZS:WLD", // Global Credit to Private Sector (% of GDP), Annual, 1960→

	// 🌍 Global Trade Volume Index - Global demand flow
	"WB:NE.EXP.GNFS.ZS:WLD", // Global Exports (% of GDP), Annual, 1960→
	"WB:NE.IMP.GNFS.ZS:WLD", // Global Imports (% of GDP), Annual, 1960→

	// 🌍 Global Energy Consumption - Industrial activity proxy
	"WB:EG.USE.PCAP.KG.OE:WLD", // Global Energy Use per capita (kg of oil equivalent), Annual, 1960→
}

// IndicatorMetadata provides rich metadata about economic indicators
type IndicatorMetadata struct {
	Code        string
	Name        string
	Unit        string
	Frequency   string
	Category    string
	Description string
	StartYear   int
	ProxyFor    string
	FREDURL     string
}

// GetIndicatorMetadata returns metadata for all configured indicators
func GetIndicatorMetadata() map[string]IndicatorMetadata {
	return map[string]IndicatorMetadata{
		// 1️⃣ Real GDP - US (FRED)
		"GDPC1": {
			Code:        "GDPC1",
			Name:        "Real Gross Domestic Product",
			Unit:        "Billions of Chained 2017 Dollars",
			Frequency:   "Quarterly",
			Category:    "Growth",
			Description: "Real GDP measures economic output adjusted for inflation",
			StartYear:   1947,
			ProxyFor:    "Economic Growth / Business Cycle",
			FREDURL:     "https://fred.stlouisfed.org/series/GDPC1",
		},

		// 1️⃣ Real GDP - Global (World Bank)
		"WB:NY.GDP.MKTP.KD:USA": {
			Code:        "WB:NY.GDP.MKTP.KD:USA",
			Name:        "Real GDP - United States (World Bank)",
			Unit:        "Constant 2015 US$",
			Frequency:   "Annual",
			Category:    "Growth",
			Description: "Real GDP in constant 2015 prices from World Bank",
			StartYear:   1960,
			ProxyFor:    "Economic Growth / Business Cycle",
			FREDURL:     "https://data.worldbank.org/indicator/NY.GDP.MKTP.KD?locations=US",
		},

		"WB:NY.GDP.MKTP.KD:CHN": {
			Code:        "WB:NY.GDP.MKTP.KD:CHN",
			Name:        "Real GDP - China (World Bank)",
			Unit:        "Constant 2015 US$",
			Frequency:   "Annual",
			Category:    "Growth",
			Description: "Real GDP in constant 2015 prices from World Bank",
			StartYear:   1960,
			ProxyFor:    "Economic Growth / Business Cycle - China",
			FREDURL:     "https://data.worldbank.org/indicator/NY.GDP.MKTP.KD?locations=CN",
		},

		"WB:NY.GDP.MKTP.KD:EMU": {
			Code:        "WB:NY.GDP.MKTP.KD:EMU",
			Name:        "Real GDP - Euro Area (World Bank)",
			Unit:        "Constant 2015 US$",
			Frequency:   "Annual",
			Category:    "Growth",
			Description: "Real GDP in constant 2015 prices from World Bank",
			StartYear:   1960,
			ProxyFor:    "Economic Growth / Business Cycle - Euro Area",
			FREDURL:     "https://data.worldbank.org/indicator/NY.GDP.MKTP.KD?locations=EU",
		},

		"WB:NY.GDP.MKTP.KD:WLD": {
			Code:        "WB:NY.GDP.MKTP.KD:WLD",
			Name:        "Real GDP - World (World Bank)",
			Unit:        "Constant 2015 US$",
			Frequency:   "Annual",
			Category:    "Growth",
			Description: "Global real GDP in constant 2015 prices from World Bank",
			StartYear:   1960,
			ProxyFor:    "Global Economic Growth",
			FREDURL:     "https://data.worldbank.org/indicator/NY.GDP.MKTP.KD?locations=1W",
		},

		// 2️⃣ Global Inflation (World Bank)
		"WB:FP.CPI.TOTL.ZG:WLD": {
			Code:        "WB:FP.CPI.TOTL.ZG:WLD",
			Name:        "Global CPI Inflation (World Bank)",
			Unit:        "Annual %",
			Frequency:   "Annual",
			Category:    "Inflation",
			Description: "Global consumer price index inflation, measures worldwide inflationary pressure",
			StartYear:   1960,
			ProxyFor:    "Global Inflation / Monetary Tightness",
			FREDURL:     "https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG?locations=1W",
		},

		"WB:FP.CPI.TOTL.ZG:CHN": {
			Code:        "WB:FP.CPI.TOTL.ZG:CHN",
			Name:        "China CPI Inflation (World Bank)",
			Unit:        "Annual %",
			Frequency:   "Annual",
			Category:    "Inflation",
			Description: "China consumer price index inflation",
			StartYear:   1960,
			ProxyFor:    "China Inflation",
			FREDURL:     "https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG?locations=CN",
		},

		"WB:FP.CPI.TOTL.ZG:EMU": {
			Code:        "WB:FP.CPI.TOTL.ZG:EMU",
			Name:        "Euro Area CPI Inflation (World Bank)",
			Unit:        "Annual %",
			Frequency:   "Annual",
			Category:    "Inflation",
			Description: "Euro Area consumer price index inflation (proxy for HICP)",
			StartYear:   1960,
			ProxyFor:    "Euro Area Inflation / ECB Policy Driver",
			FREDURL:     "https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG?locations=EU",
		},

		// 5️⃣ China M2 Money Supply (World Bank)
		"WB:FM.LBL.BMNY.CN:CHN": {
			Code:        "WB:FM.LBL.BMNY.CN:CHN",
			Name:        "China Broad Money (M2) - World Bank",
			Unit:        "Current LCU",
			Frequency:   "Annual",
			Category:    "Liquidity",
			Description: "China M2 money supply, major global liquidity driver",
			StartYear:   1960,
			ProxyFor:    "China Credit Expansion / Global Liquidity",
			FREDURL:     "https://data.worldbank.org/indicator/FM.LBL.BMNY.CN?locations=CN",
		},

		// 🔟 Global Credit to Private Sector
		"WB:FS.AST.PRVT.GD.ZS:WLD": {
			Code:        "WB:FS.AST.PRVT.GD.ZS:WLD",
			Name:        "Global Credit to Private Sector (% of GDP)",
			Unit:        "% of GDP",
			Frequency:   "Annual",
			Category:    "Liquidity",
			Description: "Global domestic credit to private sector, indicates leverage buildup and financial cycle risk",
			StartYear:   1960,
			ProxyFor:    "Global Leverage / Financial Cycle Risk",
			FREDURL:     "https://data.worldbank.org/indicator/FS.AST.PRVT.GD.ZS?locations=1W",
		},

		// 3️⃣ 9️⃣ Global Trade (World Bank)
		"WB:NE.EXP.GNFS.ZS:WLD": {
			Code:        "WB:NE.EXP.GNFS.ZS:WLD",
			Name:        "Global Exports of Goods and Services (% of GDP)",
			Unit:        "% of GDP",
			Frequency:   "Annual",
			Category:    "Trade",
			Description: "Global exports as percentage of GDP, proxy for world trade volume and demand",
			StartYear:   1960,
			ProxyFor:    "Global Trade Volume / Cross-Border Demand",
			FREDURL:     "https://data.worldbank.org/indicator/NE.EXP.GNFS.ZS?locations=1W",
		},

		"WB:NE.IMP.GNFS.ZS:WLD": {
			Code:        "WB:NE.IMP.GNFS.ZS:WLD",
			Name:        "Global Imports of Goods and Services (% of GDP)",
			Unit:        "% of GDP",
			Frequency:   "Annual",
			Category:    "Trade",
			Description: "Global imports as percentage of GDP, measures cross-border demand",
			StartYear:   1960,
			ProxyFor:    "Global Trade Volume / Import Demand",
			FREDURL:     "https://data.worldbank.org/indicator/NE.IMP.GNFS.ZS?locations=1W",
		},

		// 11️⃣ Global Energy & Emissions (World Bank)
		"WB:EG.USE.PCAP.KG.OE:WLD": {
			Code:        "WB:EG.USE.PCAP.KG.OE:WLD",
			Name:        "Global Energy Use per Capita",
			Unit:        "kg of oil equivalent per capita",
			Frequency:   "Annual",
			Category:    "Commodities",
			Description: "Global energy consumption per capita, proxy for industrial activity and growth",
			StartYear:   1960,
			ProxyFor:    "Industrial Activity / Long-term Growth",
			FREDURL:     "https://data.worldbank.org/indicator/EG.USE.PCAP.KG.OE?locations=1W",
		},

		// 2️⃣ CPI - Headline
		"CPIAUCSL": {
			Code:        "CPIAUCSL",
			Name:        "Consumer Price Index for All Urban Consumers: All Items",
			Unit:        "Index 1982-1984=100",
			Frequency:   "Monthly",
			Category:    "Inflation",
			Description: "Headline CPI measuring price changes for consumer goods and services",
			StartYear:   1947,
			ProxyFor:    "Inflation",
			FREDURL:     "https://fred.stlouisfed.org/series/CPIAUCSL",
		},

		// 2️⃣ CPI - Core
		"CPILFESL": {
			Code:        "CPILFESL",
			Name:        "Consumer Price Index for All Urban Consumers: All Items Less Food and Energy",
			Unit:        "Index 1982-1984=100",
			Frequency:   "Monthly",
			Category:    "Inflation",
			Description: "Core CPI excluding volatile food and energy components",
			StartYear:   1957,
			ProxyFor:    "Core Inflation",
			FREDURL:     "https://fred.stlouisfed.org/series/CPILFESL",
		},

		// 3️⃣ Unemployment Rate
		"UNRATE": {
			Code:        "UNRATE",
			Name:        "Unemployment Rate",
			Unit:        "Percent",
			Frequency:   "Monthly",
			Category:    "Labor Market",
			Description: "Percentage of labor force that is unemployed and actively seeking work",
			StartYear:   1948,
			ProxyFor:    "Labor Market Health",
			FREDURL:     "https://fred.stlouisfed.org/series/UNRATE",
		},

		// 4️⃣ Fed Funds Rate
		"FEDFUNDS": {
			Code:        "FEDFUNDS",
			Name:        "Federal Funds Effective Rate",
			Unit:        "Percent",
			Frequency:   "Monthly",
			Category:    "Monetary Policy",
			Description: "Overnight lending rate between banks, primary tool of U.S. monetary policy",
			StartYear:   1954,
			ProxyFor:    "Monetary Policy Stance",
			FREDURL:     "https://fred.stlouisfed.org/series/FEDFUNDS",
		},

		// 5️⃣ Treasury Yield Curve - 10Y-2Y Spread
		"T10Y2Y": {
			Code:        "T10Y2Y",
			Name:        "10-Year Treasury Constant Maturity Minus 2-Year Treasury Constant Maturity",
			Unit:        "Percent",
			Frequency:   "Daily",
			Category:    "Yield Curve",
			Description: "Yield curve slope; negative values historically precede recessions",
			StartYear:   1976,
			ProxyFor:    "Recession Signal / Yield Curve Shape",
			FREDURL:     "https://fred.stlouisfed.org/series/T10Y2Y",
		},

		// 6️⃣ M2 Money Supply
		"M2SL": {
			Code:        "M2SL",
			Name:        "M2 Money Stock",
			Unit:        "Billions of Dollars",
			Frequency:   "Monthly",
			Category:    "Liquidity",
			Description: "Broad money supply including cash, deposits, and near money",
			StartYear:   1959,
			ProxyFor:    "Liquidity / Credit Conditions",
			FREDURL:     "https://fred.stlouisfed.org/series/M2SL",
		},

		// 7️⃣ Brent Crude Oil Price
		"POILBREUSDM": {
			Code:        "POILBREUSDM",
			Name:        "Global Price of Brent Crude",
			Unit:        "Dollars per Barrel",
			Frequency:   "Monthly",
			Category:    "Commodities",
			Description: "Global oil price benchmark, key inflation driver",
			StartYear:   1987,
			ProxyFor:    "Inflation Driver & Global Demand",
			FREDURL:     "https://fred.stlouisfed.org/series/POILBREUSDM",
		},

		// 8️⃣ Industrial Production: Manufacturing
		"IPMAN": {
			Code:        "IPMAN",
			Name:        "Industrial Production: Manufacturing (NAICS)",
			Unit:        "Index 2017=100, Seasonally Adjusted",
			Frequency:   "Monthly",
			Category:    "Manufacturing",
			Description: "Manufacturing production index - measures real output of the manufacturing sector",
			StartYear:   1972,
			ProxyFor:    "Manufacturing Activity / Industrial Output",
			FREDURL:     "https://fred.stlouisfed.org/series/IPMAN",
		},

		// 9️⃣ Consumer Sentiment
		"UMCSENT": {
			Code:        "UMCSENT",
			Name:        "University of Michigan: Consumer Sentiment",
			Unit:        "Index 1966:Q1=100",
			Frequency:   "Monthly",
			Category:    "Sentiment",
			Description: "Consumer confidence and economic expectations survey",
			StartYear:   1978,
			ProxyFor:    "Household Confidence",
			FREDURL:     "https://fred.stlouisfed.org/series/UMCSENT",
		},

		// 🔟 Inflation Expectations
		"T5YIE": {
			Code:        "T5YIE",
			Name:        "5-Year Breakeven Inflation Rate",
			Unit:        "Percent",
			Frequency:   "Daily",
			Category:    "Inflation",
			Description: "Market-based measure of expected inflation derived from 5-Year Treasury securities and Treasury Inflation-Protected Securities (TIPS)",
			StartYear:   2003,
			ProxyFor:    "Market Inflation Expectations / Forward Inflation Outlook",
			FREDURL:     "https://fred.stlouisfed.org/series/T5YIE",
		},

		// 1️⃣1️⃣ Financial Conditions Index
		"NFCI": {
			Code:        "NFCI",
			Name:        "Chicago Fed National Financial Conditions Index",
			Unit:        "Index (0 = average conditions)",
			Frequency:   "Weekly",
			Category:    "Financial Conditions",
			Description: "Index of financial stress - positive values indicate tighter than average financial conditions, negative values indicate looser conditions",
			StartYear:   1971,
			ProxyFor:    "Financial Market Stress / Credit Availability",
			FREDURL:     "https://fred.stlouisfed.org/series/NFCI",
		},
	}
}

// GetIndicatorCategories returns a map of categories to indicator codes
func GetIndicatorCategories() map[string][]string {
	return map[string][]string{
		"Growth":               {"GDPC1", "WB:NY.GDP.MKTP.KD:USA", "WB:NY.GDP.MKTP.KD:CHN", "WB:NY.GDP.MKTP.KD:EMU", "WB:NY.GDP.MKTP.KD:WLD"},
		"Inflation":            {"CPIAUCSL", "CPILFESL", "T5YIE", "WB:FP.CPI.TOTL.ZG:WLD", "WB:FP.CPI.TOTL.ZG:CHN", "WB:FP.CPI.TOTL.ZG:EMU"},
		"Labor Market":         {"UNRATE"},
		"Monetary Policy":      {"FEDFUNDS"},
		"Yield Curve":          {"T10Y2Y"},
		"Liquidity":            {"M2SL", "WB:FM.LBL.BMNY.CN:CHN", "WB:FS.AST.PRVT.GD.ZS:WLD"},
		"Commodities":          {"POILBREUSDM", "WB:EG.USE.PCAP.KG.OE:WLD"},
		"Manufacturing":        {"IPMAN"},
		"Sentiment":            {"UMCSENT"},
		"Trade":                {"WB:NE.EXP.GNFS.ZS:WLD", "WB:NE.IMP.GNFS.ZS:WLD"},
		"Financial Conditions": {"NFCI"},
	}
}

// GetCoreIndicators returns the primary U.S. indicators (one per category)
func GetCoreIndicators() []string {
	return []string{
		"GDPC1",       // Growth
		"CPIAUCSL",    // Inflation
		"UNRATE",      // Labor Market
		"FEDFUNDS",    // Monetary Policy
		"T10Y2Y",      // Yield Curve
		"M2SL",        // Liquidity
		"POILBREUSDM", // Commodities
		"IPMAN",       // Manufacturing
		"UMCSENT",     // Sentiment
		"T5YIE",       // Inflation Expectations
		"NFCI",        // Financial Conditions
	}
}
