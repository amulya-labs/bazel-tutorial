package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DataSource represents a generic data source interface
type DataSource interface {
	FetchSeries(seriesID string) (*Series, []Observation, error)
	GetSourceName() string
}

// WorldBankClient handles requests to the World Bank API
type WorldBankClient struct {
	httpClient HTTPClient
}

const (
	worldBankBaseURL = "https://api.worldbank.org/v2"
)

// World Bank API response structures
type wbIndicatorResponse struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type wbDataPoint struct {
	Indicator struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"indicator"`
	Country struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"country"`
	Countryiso3code string   `json:"countryiso3code"`
	Date            string   `json:"date"`
	Value           *float64 `json:"value"` // Pointer to distinguish null from 0
	Unit            string   `json:"unit"`
	ObsStatus       string   `json:"obs_status"`
	Decimal         int      `json:"decimal"`
}

// NewWorldBankClient creates a new World Bank API client
func NewWorldBankClient() *WorldBankClient {
	return &WorldBankClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewWorldBankClientWithHTTP creates a World Bank client with a custom HTTP client (for testing)
func NewWorldBankClientWithHTTP(httpClient HTTPClient) *WorldBankClient {
	return &WorldBankClient{
		httpClient: httpClient,
	}
}

// FetchSeries fetches data from World Bank API
// For World Bank, seriesID format is: "WB:INDICATOR:COUNTRY" e.g., "WB:NY.GDP.MKTP.KD:USA"
func (c *WorldBankClient) FetchSeries(seriesID string) (*Series, []Observation, error) {
	// Parse the series ID to extract indicator and country
	// Format: WB:INDICATOR:COUNTRY
	// Example: WB:NY.GDP.MKTP.KD:USA

	parts := strings.Split(seriesID, ":")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid World Bank series ID format: %s (expected WB:INDICATOR:COUNTRY)", seriesID)
	}

	indicator := parts[1] // e.g., NY.GDP.MKTP.KD
	country := parts[2]   // e.g., USA, CHN, EMU, WLD

	// Fetch metadata
	series := &Series{
		ID:    seriesID,
		Title: fmt.Sprintf("Real GDP (%s)", country),
		Units: "Constant 2015 US$",
	}

	// Fetch data from World Bank API
	// Example: https://api.worldbank.org/v2/country/USA/indicator/NY.GDP.MKTP.KD?format=json&date=1990:2024
	url := fmt.Sprintf("%s/country/%s/indicator/%s?format=json&date=1990:2024&per_page=1000",
		worldBankBaseURL, country, indicator)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// World Bank returns an array with 2 elements: [metadata, data]
	var result []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result) < 2 {
		return nil, nil, fmt.Errorf("unexpected response format from World Bank API")
	}

	// Parse the data points
	dataBytes, err := json.Marshal(result[1])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var dataPoints []wbDataPoint
	if err := json.Unmarshal(dataBytes, &dataPoints); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal data points: %w", err)
	}

	// Convert to observations
	observations := make([]Observation, 0, len(dataPoints))
	for _, dp := range dataPoints {
		if dp.Value != nil { // Skip null values (missing data)
			// World Bank returns annual data by year only, convert to ISO date
			dateStr := fmt.Sprintf("%s-01-01", dp.Date)
			observations = append(observations, Observation{
				Date:  dateStr,
				Value: fmt.Sprintf("%f", *dp.Value),
			})
		}
	}

	return series, observations, nil
}

// GetSourceName returns the name of this data source
func (c *WorldBankClient) GetSourceName() string {
	return "World Bank"
}

// OECDClient handles requests to the OECD API
type OECDClient struct {
	httpClient HTTPClient
}

const (
	oecdBaseURL = "https://stats.oecd.org/sdmx-json/data"
)

// NewOECDClient creates a new OECD API client
func NewOECDClient() *OECDClient {
	return &OECDClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchSeries fetches data from OECD API
// For OECD, seriesID format is: "OECD:DATASET.SERIES" e.g., "OECD:MEI.LRUNTTTT.USA.M"
func (c *OECDClient) FetchSeries(seriesID string) (*Series, []Observation, error) {
	// Placeholder implementation
	// OECD API is more complex and requires specific dataset and dimension parsing
	return nil, nil, fmt.Errorf("OECD data source not yet implemented")
}

// GetSourceName returns the name of this data source
func (c *OECDClient) GetSourceName() string {
	return "OECD"
}

// DataSourceFactory creates the appropriate data source client based on series ID prefix
func DataSourceFactory(seriesID string, fredAPIKey string) (DataSource, string, error) {
	// Determine source from series ID prefix using idiomatic string matching
	if strings.HasPrefix(seriesID, "WB:") {
		// World Bank series
		return NewWorldBankClient(), seriesID, nil
	} else if strings.HasPrefix(seriesID, "OECD:") {
		// OECD series not yet implemented
		return nil, seriesID, fmt.Errorf("OECD data source not yet implemented")
	} else {
		// Default to FRED (no prefix)
		return NewFREDClient(fredAPIKey), seriesID, nil
	}
}
