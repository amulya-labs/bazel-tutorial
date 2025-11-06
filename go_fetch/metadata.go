// Package main - Indicator configuration for high-ROI macro indicators
package main

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
		// 1️⃣ Real GDP
		"GDPC1": {
			Code:        "GDPC1",
			Name:        "Real Gross Domestic Product",
			Unit:        "Billions of Chained 2012 Dollars",
			Frequency:   "Quarterly",
			Category:    "Growth",
			Description: "Real GDP measures economic output adjusted for inflation",
			StartYear:   1947,
			ProxyFor:    "Economic Growth / Business Cycle",
			FREDURL:     "https://fred.stlouisfed.org/series/GDPC1",
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
		
		"DGS10": {
			Code:        "DGS10",
			Name:        "Market Yield on U.S. Treasury Securities at 10-Year Constant Maturity",
			Unit:        "Percent",
			Frequency:   "Daily",
			Category:    "Yield Curve",
			Description: "10-year Treasury yield, long-term interest rate benchmark",
			StartYear:   1962,
			ProxyFor:    "Long-term Interest Rates",
			FREDURL:     "https://fred.stlouisfed.org/series/DGS10",
		},
		
		"DGS2": {
			Code:        "DGS2",
			Name:        "Market Yield on U.S. Treasury Securities at 2-Year Constant Maturity",
			Unit:        "Percent",
			Frequency:   "Daily",
			Category:    "Yield Curve",
			Description: "2-year Treasury yield, short-term interest rate benchmark",
			StartYear:   1976,
			ProxyFor:    "Short-term Interest Rates",
			FREDURL:     "https://fred.stlouisfed.org/series/DGS2",
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
		
		// 8️⃣ Manufacturing Activity
		"MANEMP": {
			Code:        "MANEMP",
			Name:        "All Employees, Manufacturing",
			Unit:        "Thousands of Persons",
			Frequency:   "Monthly",
			Category:    "Manufacturing",
			Description: "Total manufacturing employment, proxy for industrial activity",
			StartYear:   1939,
			ProxyFor:    "Business Confidence / Early Cycle Indicator",
			FREDURL:     "https://fred.stlouisfed.org/series/MANEMP",
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
		
		// 🔟 Global Trade - Imports
		"IMPCH": {
			Code:        "IMPCH",
			Name:        "Real Imports of Goods and Services",
			Unit:        "Billions of Chained 2012 Dollars",
			Frequency:   "Quarterly",
			Category:    "Trade",
			Description: "Real imports adjusted for inflation, proxy for global trade",
			StartYear:   1947,
			ProxyFor:    "Global Demand Flow / International Trade",
			FREDURL:     "https://fred.stlouisfed.org/series/IMPCH",
		},
		
		// 🔟 Global Trade - Exports
		"EXPCH": {
			Code:        "EXPCH",
			Name:        "Real Exports of Goods and Services",
			Unit:        "Billions of Chained 2012 Dollars",
			Frequency:   "Quarterly",
			Category:    "Trade",
			Description: "Real exports adjusted for inflation, proxy for global trade",
			StartYear:   1947,
			ProxyFor:    "Global Demand Flow / International Trade",
			FREDURL:     "https://fred.stlouisfed.org/series/EXPCH",
		},
	}
}

// GetIndicatorCategories returns a map of categories to indicator codes
func GetIndicatorCategories() map[string][]string {
	return map[string][]string{
		"Growth":          {"GDPC1"},
		"Inflation":       {"CPIAUCSL", "CPILFESL", "POILBREUSDM"},
		"Labor Market":    {"UNRATE"},
		"Monetary Policy": {"FEDFUNDS"},
		"Yield Curve":     {"T10Y2Y", "DGS10", "DGS2"},
		"Liquidity":       {"M2SL"},
		"Commodities":     {"POILBREUSDM"},
		"Manufacturing":   {"MANEMP"},
		"Sentiment":       {"UMCSENT"},
		"Trade":           {"IMPCH", "EXPCH"},
	}
}

// GetCoreIndicators returns the primary indicators (one per category)
func GetCoreIndicators() []string {
	return []string{
		"GDPC1",       // Growth
		"CPIAUCSL",    // Inflation
		"UNRATE",      // Labor Market
		"FEDFUNDS",    // Monetary Policy
		"T10Y2Y",      // Yield Curve
		"M2SL",        // Liquidity
		"POILBREUSDM", // Commodities
		"MANEMP",      // Manufacturing
		"UMCSENT",     // Sentiment
		"IMPCH",       // Trade
	}
}
