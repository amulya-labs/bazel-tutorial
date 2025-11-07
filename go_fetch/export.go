package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// IndicatorSummary represents a summary of an economic indicator
type IndicatorSummary struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Unit        string   `json:"unit"`
	Source      string   `json:"source"`
	Category    string   `json:"category,omitempty"`
	ProxyFor    string   `json:"proxy_for,omitempty"`
	LastUpdated string   `json:"last_updated"`
	Value       float64  `json:"value"`
	DeltaMoM    *float64 `json:"delta_mom"`
	DeltaYoY    *float64 `json:"delta_yoy"`
}

// SummaryResponse represents the /api/econ/summary endpoint response
type SummaryResponse struct {
	Indicators  []IndicatorSummary `json:"indicators"`
	Count       int                `json:"count"`
	LastRefresh *string            `json:"last_refresh"`
}

// ObservationData represents a single observation
type ObservationData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// SeriesResponse represents the /api/econ/series endpoint response
type SeriesResponse struct {
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	Unit         string            `json:"unit"`
	Source       string            `json:"source"`
	LastUpdated  *string           `json:"last_updated"`
	LatestValue  *float64          `json:"latest_value"`
	DeltaMoM     *float64          `json:"delta_mom"`
	DeltaYoY     *float64          `json:"delta_yoy"`
	Observations []ObservationData `json:"observations"`
	Count        int               `json:"count"`
	Range        string            `json:"range"`
}

// ExportJSON exports database data to JSON files for static hosting
func ExportJSON(store *SQLiteStore, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export metadata
	if err := exportMetadata(outputDir); err != nil {
		return fmt.Errorf("failed to export metadata: %w", err)
	}

	// Export summary
	if err := exportSummary(store, outputDir); err != nil {
		return fmt.Errorf("failed to export summary: %w", err)
	}

	// Export individual series
	series, err := store.GetAllSeries()
	if err != nil {
		return fmt.Errorf("failed to get series list: %w", err)
	}

	for _, s := range series {
		if err := exportSeries(store, s.ID, outputDir); err != nil {
			return fmt.Errorf("failed to export series %s: %w", s.ID, err)
		}
	}

	fmt.Printf("Successfully exported data to %s\n", outputDir)
	return nil
}

// exportSummary generates summary.json
func exportSummary(store *SQLiteStore, outputDir string) error {
	series, err := store.GetAllSeries()
	if err != nil {
		return err
	}

	indicators := []IndicatorSummary{}
	metadata := GetIndicatorMetadata()

	for _, s := range series {
		// Get latest observation
		date, value, err := store.GetLatestObservation(s.ID)
		if err != nil {
			continue // Skip if no data
		}

		// Calculate deltas
		deltaMoM := calculateDeltaMoM(store, s.ID)
		deltaYoY := calculateDeltaYoY(store, s.ID)

		// Get metadata if available
		meta, hasMetadata := metadata[s.ID]
		category := ""
		proxyFor := ""
		if hasMetadata {
			category = meta.Category
			proxyFor = meta.ProxyFor
		}

		indicator := IndicatorSummary{
			Code:        s.ID,
			Name:        s.Title,
			Unit:        s.Units,
			Source:      "FRED",
			Category:    category,
			ProxyFor:    proxyFor,
			LastUpdated: date,
			Value:       value,
			DeltaMoM:    roundFloat(deltaMoM, 2),
			DeltaYoY:    roundFloat(deltaYoY, 2),
		}

		indicators = append(indicators, indicator)
	}

	// Get last refresh info
	lastRefresh, err := store.GetLastRefresh()
	if err != nil {
		lastRefresh = nil
	}

	response := SummaryResponse{
		Indicators:  indicators,
		Count:       len(indicators),
		LastRefresh: lastRefresh,
	}

	// Write to file
	filePath := filepath.Join(outputDir, "summary.json")
	return writeJSON(filePath, response)
}

// exportSeries generates series-{CODE}.json for a specific series
func exportSeries(store *SQLiteStore, seriesID string, outputDir string) error {
	// Get series metadata
	series, err := store.GetSeriesMetadata(seriesID)
	if err != nil {
		return err
	}

	// Get all observations
	observations, err := store.GetAllObservations(seriesID)
	if err != nil {
		return err
	}

	// Convert to response format
	obsData := make([]ObservationData, len(observations))
	for i, obs := range observations {
		obsData[i] = ObservationData{
			Date:  obs.Date,
			Value: obs.Value,
		}
	}

	// Get latest observation
	date, value, err := store.GetLatestObservation(seriesID)
	var lastUpdated *string
	var latestValue *float64
	if err == nil {
		lastUpdated = &date
		latestValue = &value
	}

	// Calculate deltas
	deltaMoM := calculateDeltaMoM(store, seriesID)
	deltaYoY := calculateDeltaYoY(store, seriesID)

	response := SeriesResponse{
		Code:         series.ID,
		Name:         series.Title,
		Unit:         series.Units,
		Source:       "FRED",
		LastUpdated:  lastUpdated,
		LatestValue:  latestValue,
		DeltaMoM:     roundFloat(deltaMoM, 2),
		DeltaYoY:     roundFloat(deltaYoY, 2),
		Observations: obsData,
		Count:        len(obsData),
		Range:        "max",
	}

	// Write to file
	filePath := filepath.Join(outputDir, fmt.Sprintf("series-%s.json", seriesID))
	return writeJSON(filePath, response)
}

// calculateDeltaMoM calculates month-over-month percentage change
func calculateDeltaMoM(store *SQLiteStore, seriesID string) *float64 {
	rows, err := store.db.Query(`
		SELECT value
		FROM observations
		WHERE series_id = ?
		ORDER BY date DESC
		LIMIT 2
	`, seriesID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var values []float64
	for rows.Next() {
		var value float64
		if err := rows.Scan(&value); err != nil {
			return nil
		}
		values = append(values, value)
	}

	if len(values) < 2 {
		return nil
	}

	current := values[0]
	previous := values[1]

	if previous == 0 {
		return nil
	}

	delta := ((current - previous) / previous) * 100
	return &delta
}

// calculateDeltaYoY calculates year-over-year percentage change
func calculateDeltaYoY(store *SQLiteStore, seriesID string) *float64 {
	// Get latest observation
	latestDate, latestValue, err := store.GetLatestObservation(seriesID)
	if err != nil {
		return nil
	}

	// Parse latest date
	t, err := time.Parse("2006-01-02", latestDate)
	if err != nil {
		return nil
	}

	// Calculate date one year ago
	yearAgo := t.AddDate(-1, 0, 0).Format("2006-01-02")

	// Get observation at or before year ago date
	var yearAgoValue float64
	err = store.db.QueryRow(`
		SELECT value
		FROM observations
		WHERE series_id = ? AND date <= ?
		ORDER BY date DESC
		LIMIT 1
	`, seriesID, yearAgo).Scan(&yearAgoValue)

	if err != nil {
		return nil
	}

	if yearAgoValue == 0 {
		return nil
	}

	delta := ((latestValue - yearAgoValue) / yearAgoValue) * 100
	return &delta
}

// roundFloat rounds a float pointer to specified decimal places
func roundFloat(f *float64, decimals int) *float64 {
	if f == nil {
		return nil
	}

	multiplier := float64(1)
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}

	rounded := float64(int(*f*multiplier+0.5)) / multiplier
	return &rounded
}

// writeJSON writes data to a JSON file
func writeJSON(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Printf("✓ Exported %s\n", filepath.Base(filePath))
	return nil
}

// ObservationRecord represents a database observation record
type ObservationRecord struct {
	Date  string
	Value float64
}

// GetAllObservations retrieves all observations for a series
func (s *SQLiteStore) GetAllObservations(seriesID string) ([]ObservationRecord, error) {
	rows, err := s.db.Query(`
		SELECT date, value
		FROM observations
		WHERE series_id = ?
		ORDER BY date
	`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []ObservationRecord
	for rows.Next() {
		var obs ObservationRecord
		if err := rows.Scan(&obs.Date, &obs.Value); err != nil {
			return nil, err
		}
		observations = append(observations, obs)
	}

	return observations, rows.Err()
}

// GetAllSeries retrieves all series metadata
func (s *SQLiteStore) GetAllSeries() ([]Series, error) {
	rows, err := s.db.Query(`
		SELECT id, name, unit, source
		FROM series
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var series Series
		var source string
		if err := rows.Scan(&series.ID, &series.Title, &series.Units, &source); err != nil {
			return nil, err
		}
		seriesList = append(seriesList, series)
	}

	return seriesList, rows.Err()
}

// GetSeriesMetadata retrieves metadata for a specific series
func (s *SQLiteStore) GetSeriesMetadata(seriesID string) (*Series, error) {
	var series Series
	err := s.db.QueryRow(`
		SELECT id, name, unit
		FROM series
		WHERE id = ?
	`, seriesID).Scan(&series.ID, &series.Title, &series.Units)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("series %s not found", seriesID)
	}

	if err != nil {
		return nil, err
	}

	return &series, nil
}

// GetLastRefresh retrieves the timestamp of the last successful refresh
func (s *SQLiteStore) GetLastRefresh() (*string, error) {
	var finishedAt string
	err := s.db.QueryRow(`
		SELECT finished_at
		FROM refresh_log
		WHERE ok = 1
		ORDER BY started_at DESC
		LIMIT 1
	`).Scan(&finishedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &finishedAt, nil
}

// exportMetadata generates metadata.json with rich indicator information
func exportMetadata(outputDir string) error {
	metadata := GetIndicatorMetadata()
	categories := GetIndicatorCategories()

	// Create metadata export structure
	type MetadataExport struct {
		Indicators map[string]IndicatorMetadata `json:"indicators"`
		Categories map[string][]string          `json:"categories"`
		CoreSeries []string                     `json:"core_series"`
		Version    string                       `json:"version"`
		Generated  string                       `json:"generated"`
	}

	export := MetadataExport{
		Indicators: metadata,
		Categories: categories,
		CoreSeries: GetCoreIndicators(),
		Version:    "1.0.0",
		Generated:  time.Now().Format(time.RFC3339),
	}

	filePath := filepath.Join(outputDir, "metadata.json")
	return writeJSON(filePath, export)
}
