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

	// Series to fetch from FRED - High-ROI Global Economic Indicators (30-Year Baseline)
	// These indicators provide comprehensive coverage of growth, inflation, liquidity,
	// sentiment, and risk since ~1990.
	defaultSeries = []string{
		// 1️⃣ Real GDP (constant prices) - Growth / cycle
		"GDPC1", // US Real GDP (Billions of Chained 2012 Dollars), Quarterly, 1947→

		// 2️⃣ CPI (Consumer Price Index) - Inflation
		"CPIAUCSL", // US CPI All Urban Consumers (Index 1982-84=100), Monthly, 1947→
		"CPILFESL", // US Core CPI Less Food & Energy (Index 1982-84=100), Monthly, 1957→

		// 3️⃣ Unemployment Rate - Labor market
		"UNRATE", // US Unemployment Rate (Percent), Monthly, 1948→

		// 4️⃣ Fed Funds Rate - Monetary policy stance
		"FEDFUNDS", // Federal Funds Effective Rate (Percent), Monthly, 1954→

		// 5️⃣ 10Y – 2Y Treasury Spread - Recession signal / yield curve
		"T10Y2Y", // 10-Year Treasury Minus 2-Year Treasury (Percent), Daily, 1976→
		"DGS10",  // 10-Year Treasury Constant Maturity Rate (Percent), Daily, 1962→
		"DGS2",   // 2-Year Treasury Constant Maturity Rate (Percent), Daily, 1976→

		// 6️⃣ M2 Money Supply - Liquidity / credit conditions
		"M2SL", // M2 Money Stock (Billions of Dollars), Monthly, 1959→

		// 7️⃣ Brent Crude Oil Price - Inflation driver & demand proxy
		"POILBREUSDM", // Global Price of Brent Crude (Dollars per Barrel), Monthly, 1987→

		// 8️⃣ Manufacturing PMI - Business confidence / early cycle
		// Note: ISM Manufacturing PMI is available from FRED
		"MANEMP", // Manufacturing Employment (Thousands), Monthly, 1939→
		// ISM PMI data would be ideal but requires separate source (not in FRED)

		// 9️⃣ Consumer Sentiment Index - Household confidence
		"UMCSENT", // University of Michigan Consumer Sentiment (Index 1966:Q1=100), Monthly, 1978→

		// 🔟 Global Trade Volume Index - Global demand flow
		// Note: CPB Netherlands World Trade Monitor not available in FRED
		// Using US-specific trade proxy as alternative
		"IMPCH", // Real Imports of Goods & Services (Billions of Chained 2012 Dollars), Quarterly, 1947→
		"EXPCH", // Real Exports of Goods & Services (Billions of Chained 2012 Dollars), Quarterly, 1947→
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

	// Create FRED client
	client := NewFREDClient(*apiKey)

	// Log refresh start
	refreshID, err := store.LogRefreshStart("FRED")
	if err != nil {
		log.Printf("Warning: Failed to log refresh start: %v", err)
	}

	// Fetch each series
	successCount := 0
	errorCount := 0
	for _, seriesID := range defaultSeries {
		log.Printf("Fetching series: %s", seriesID)

		series, observations, err := client.FetchSeries(seriesID)
		if err != nil {
			log.Printf("Error fetching %s: %v", seriesID, err)
			errorCount++
			continue
		}

		// Store series metadata
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
