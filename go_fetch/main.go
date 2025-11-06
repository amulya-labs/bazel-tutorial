// Package main implements a CLI tool that fetches economic indicators from FRED API
// and stores them in a SQLite database for the Economic Dashboard.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// API key for FRED (Federal Reserve Economic Data)
	apiKey     = flag.String("api-key", os.Getenv("FRED_API_KEY"), "FRED API key")
	dbPath     = flag.String("db", "./data/econ.db", "Path to SQLite database")
	exportJSON = flag.String("export-json", "", "Export data to JSON files in specified directory")

	// Series to fetch - High-ROI Global Economic Indicators (30-Year Baseline)
	// These indicators provide comprehensive coverage of growth, inflation, liquidity,
	// sentiment, and risk since ~1990.
	//
	// Supports multiple data sources:
	// - FRED series: use code directly (e.g., "GDPC1")
	// - World Bank: prefix with "WB:" (e.g., "WB:NY.GDP.MKTP.KD:USA")
	// - OECD: prefix with "OECD:" (e.g., "OECD:MEI.LRUNTTTT.USA.M")
	defaultSeries = []string{
		// 1️⃣ Real GDP (constant prices) - Growth / cycle
		"GDPC1", // US Real GDP (Billions of Chained 2012 Dollars), Quarterly, 1947→ [FRED]

		// 🌍 Global GDP from World Bank
		"WB:NY.GDP.MKTP.KD:WLD", // 1️⃣ World Real GDP (constant 2015 USD), Annual, 1960→
		"WB:NY.GDP.MKTP.KD:USA", // US Real GDP from World Bank, Annual, 1960→
		"WB:NY.GDP.MKTP.KD:CHN", // 4️⃣ China Real GDP, Annual, 1960→
		"WB:NY.GDP.MKTP.KD:EMU", // 6️⃣ Euro Area Real GDP, Annual, 1960→

		// 2️⃣ CPI (Consumer Price Index) - Inflation
		"CPIAUCSL", // US CPI All Urban Consumers (Index 1982-84=100), Monthly, 1947→ [FRED]
		"CPILFESL", // US Core CPI Less Food & Energy (Index 1982-84=100), Monthly, 1957→ [FRED]

		// 🌍 Global Inflation from World Bank
		"WB:FP.CPI.TOTL.ZG:WLD", // 2️⃣ Global CPI (inflation, annual %), Annual, 1960→
		"WB:FP.CPI.TOTL.ZG:CHN", // China CPI inflation, Annual, 1960→
		"WB:FP.CPI.TOTL.ZG:EMU", // 7️⃣ Euro Area CPI inflation (proxy for HICP), Annual, 1960→

		// 3️⃣ Unemployment Rate - Labor market
		"UNRATE", // US Unemployment Rate (Percent), Monthly, 1948→ [FRED]

		// 4️⃣ Fed Funds Rate - Monetary policy stance
		"FEDFUNDS", // Federal Funds Effective Rate (Percent), Monthly, 1954→ [FRED]

		// 5️⃣ 10Y – 2Y Treasury Spread - Recession signal / yield curve
		"T10Y2Y", // 10-Year Treasury Minus 2-Year Treasury (Percent), Daily, 1976→ [FRED]
		"DGS10",  // 10-Year Treasury Constant Maturity Rate (Percent), Daily, 1962→ [FRED]
		"DGS2",   // 2-Year Treasury Constant Maturity Rate (Percent), Daily, 1976→ [FRED]

		// 6️⃣ M2 Money Supply - Liquidity / credit conditions
		"M2SL",                     // US M2 Money Stock (Billions of Dollars), Monthly, 1959→ [FRED]
		"WB:FM.LBL.BMNY.CN:CHN",    // 5️⃣ China Broad Money (M2), Annual, 1960→
		"WB:FS.AST.PRVT.GD.ZS:WLD", // 🔟 Global Credit to Private Sector (% of GDP), Annual, 1960→

		// 7️⃣ Brent Crude Oil Price - Inflation driver & demand proxy
		"POILBREUSDM", // Global Price of Brent Crude (Dollars per Barrel), Monthly, 1987→ [FRED]

		// 8️⃣ Manufacturing PMI - Business confidence / early cycle
		"MANEMP", // Manufacturing Employment (Thousands), Monthly, 1939→ [FRED]

		// 9️⃣ Consumer Sentiment Index - Household confidence
		"UMCSENT", // University of Michigan Consumer Sentiment (Index 1966:Q1=100), Monthly, 1978→ [FRED]

		// 🔟 Global Trade Volume Index - Global demand flow
		"IMPCH",                 // Real Imports of Goods & Services (Billions of Chained 2012 Dollars), Quarterly, 1947→ [FRED]
		"EXPCH",                 // Real Exports of Goods & Services (Billions of Chained 2012 Dollars), Quarterly, 1947→ [FRED]
		"WB:NE.EXP.GNFS.ZS:WLD", // 3️⃣ 9️⃣ Global Exports (% of GDP), Annual, 1960→
		"WB:NE.IMP.GNFS.ZS:WLD", // Global Imports (% of GDP), Annual, 1960→

		// 11️⃣ Global Energy Consumption - Industrial activity proxy
		"WB:EG.USE.PCAP.KG.OE:WLD", // 11️⃣ Global Energy Use per capita (kg of oil equivalent), Annual, 1960→
		"WB:EN.ATM.CO2E.KT:WLD",    // Global CO2 emissions (kt), Annual, 1960→
	}
)

func main() {
	flag.Parse()

	// Initialize database
	store, err := NewSQLiteStore(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()

	// Check if export-only mode
	exportOnly := *exportJSON != "" && *apiKey == ""

	if exportOnly {
		// Export-only mode: skip data refresh
		log.Printf("Export-only mode: exporting existing data from %s", *dbPath)
		if err := ExportJSON(store, *exportJSON); err != nil {
			log.Fatalf("Failed to export JSON: %v", err)
		}
		log.Printf("JSON export complete")
		return
	}

	// Normal mode: require API key for data refresh
	if *apiKey == "" {
		log.Fatal("FRED_API_KEY environment variable or --api-key flag is required.\n" +
			"Get your free API key at: https://fred.stlouisfed.org/docs/api/api_key.html\n\n" +
			"To export existing data without fetching, use --export-json without setting FRED_API_KEY")
	}

	log.Printf("Starting economic data refresh at %s", time.Now().Format(time.RFC3339))
	log.Printf("Database: %s", *dbPath)

	// Log refresh start
	refreshID, err := store.LogRefreshStart("Multi-Source")
	if err != nil {
		log.Printf("Warning: Failed to log refresh start: %v", err)
	}

	// Fetch each series using appropriate data source
	successCount := 0
	errorCount := 0
	for _, seriesID := range defaultSeries {
		log.Printf("Fetching series: %s", seriesID)

		// Create appropriate client based on series ID
		client, actualSeriesID, err := DataSourceFactory(seriesID, *apiKey)
		if err != nil {
			log.Printf("Error creating client for %s: %v", seriesID, err)
			errorCount++
			continue
		}

		series, observations, err := client.FetchSeries(actualSeriesID)
		if err != nil {
			log.Printf("Error fetching %s from %s: %v", seriesID, client.GetSourceName(), err)
			errorCount++
			continue
		}

		// Store series metadata with source information
		if err := store.StoreSeries(series); err != nil {
			log.Printf("Error storing series metadata for %s: %v", seriesID, err)
			errorCount++
			continue
		}

		// Store observations
		if err := store.StoreObservations(seriesID, observations); err != nil {
			log.Printf("Error storing observations for %s: %v", seriesID, err)
			errorCount++
			continue
		}

		log.Printf("Successfully stored %d observations for %s", len(observations), seriesID)
		successCount++
	}

	// Log refresh completion
	message := fmt.Sprintf("Fetched %d/%d series successfully", successCount, len(defaultSeries))
	if err := store.LogRefreshEnd(refreshID, errorCount == 0, message); err != nil {
		log.Printf("Warning: Failed to log refresh end: %v", err)
	}

	log.Printf("Refresh complete: %s", message)
	if errorCount > 0 {
		os.Exit(1)
	}

	// Export JSON if requested
	if *exportJSON != "" {
		log.Printf("Exporting data to JSON files in %s", *exportJSON)
		if err := ExportJSON(store, *exportJSON); err != nil {
			log.Fatalf("Failed to export JSON: %v", err)
		}
		log.Printf("JSON export complete")
	}
}
