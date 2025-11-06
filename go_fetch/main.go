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

	// Series to fetch from FRED
	defaultSeries = []string{
		"CPIAUCSL",  // CPI (Headline)
		"CPILFESL",  // Core CPI
		"UNRATE",    // Unemployment Rate
		"FEDFUNDS",  // Federal Funds Rate
		"DGS10",     // 10-Year Treasury Yield
		"M2SL",      // M2 Money Stock
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
