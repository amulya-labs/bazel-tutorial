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
	// Note: defaultSeries is now defined in metadata.go as part of the library
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

	// Initialize retry configuration
	retryConfig := DefaultRetryConfig()

	// Fetch each series using appropriate data source
	startTime := time.Now()
	successCount := 0
	errorCount := 0

	for _, seriesID := range defaultSeries {
		seriesStartTime := time.Now()
		log.Printf("[%s] Fetching series: %s", seriesID, seriesID)

		// Create appropriate client based on series ID
		client, actualSeriesID, err := DataSourceFactory(seriesID, *apiKey)
		if err != nil {
			log.Printf("[%s] ERROR: Failed to create client: %v", seriesID, err)
			errorCount++
			continue
		}

		// Fetch with retry logic
		var series *Series
		var observations []Observation

		err = RetryWithBackoff(retryConfig, func() error {
			var fetchErr error
			series, observations, fetchErr = client.FetchSeries(actualSeriesID)
			return fetchErr
		})

		if err != nil {
			log.Printf("[%s] ERROR: Failed to fetch from %s: %v", seriesID, client.GetSourceName(), err)
			errorCount++
			continue
		}

		duration := time.Since(seriesStartTime)
		log.Printf("[%s] SUCCESS: Fetched %d observations from %s in %v",
			seriesID, len(observations), client.GetSourceName(), duration)

		// Store series metadata with source information
		if err := store.StoreSeries(series); err != nil {
			log.Printf("[%s] ERROR: Failed to store series metadata: %v", seriesID, err)
			errorCount++
			continue
		}

		// Store observations
		if err := store.StoreObservations(seriesID, observations); err != nil {
			log.Printf("[%s] ERROR: Failed to store observations: %v", seriesID, err)
			errorCount++
			continue
		}

		successCount++
	}

	// Log refresh completion with summary statistics
	totalDuration := time.Since(startTime)
	message := fmt.Sprintf("Fetched %d/%d series successfully", successCount, len(defaultSeries))

	log.Printf("========================================")
	log.Printf("REFRESH SUMMARY:")
	log.Printf("  Total series: %d", len(defaultSeries))
	log.Printf("  Successful: %d", successCount)
	log.Printf("  Failed: %d", errorCount)
	log.Printf("  Total duration: %v", totalDuration)
	log.Printf("  Avg time per series: %v", totalDuration/time.Duration(len(defaultSeries)))
	log.Printf("========================================")

	if err := store.LogRefreshEnd(refreshID, errorCount == 0, message); err != nil {
		log.Printf("Warning: Failed to log refresh end: %v", err)
	}

	if errorCount > 0 {
		log.Printf("Refresh completed with errors")
		os.Exit(1)
	}

	log.Printf("Refresh completed successfully")

	// Export JSON if requested
	if *exportJSON != "" {
		log.Printf("Exporting data to JSON files in %s", *exportJSON)
		if err := ExportJSON(store, *exportJSON); err != nil {
			log.Fatalf("Failed to export JSON: %v", err)
		}
		log.Printf("JSON export complete")
	}
}
