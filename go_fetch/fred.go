package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	fredBaseURL = "https://api.stlouisfed.org/fred"
)

// FREDClient handles requests to the FRED API
type FREDClient struct {
	apiKey     string
	httpClient HTTPClient
}

// HTTPClient interface allows for mocking HTTP requests in tests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// Series represents metadata about an economic indicator series
type Series struct {
	ID    string
	Title string
	Units string
}

// Observation represents a single data point in a time series
type Observation struct {
	Date  string
	Value string
}

// FRED API response structures
type fredSeriesResponse struct {
	Seriess []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Units string `json:"units"`
	} `json:"seriess"`
}

type fredObservationsResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

// NewFREDClient creates a new FRED API client
func NewFREDClient(apiKey string) *FREDClient {
	return &FREDClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewFREDClientWithHTTP creates a FRED client with a custom HTTP client (for testing)
func NewFREDClientWithHTTP(apiKey string, httpClient HTTPClient) *FREDClient {
	return &FREDClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// FetchSeries fetches metadata and observations for a given series ID
func (c *FREDClient) FetchSeries(seriesID string) (*Series, []Observation, error) {
	// Fetch series metadata
	series, err := c.fetchSeriesMetadata(seriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch series metadata: %w", err)
	}

	// Fetch observations
	observations, err := c.fetchObservations(seriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch observations: %w", err)
	}

	return series, observations, nil
}

// fetchSeriesMetadata retrieves metadata about a series
func (c *FREDClient) fetchSeriesMetadata(seriesID string) (*Series, error) {
	url := fmt.Sprintf("%s/series?series_id=%s&api_key=%s&file_type=json",
		fredBaseURL, seriesID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp fredSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Seriess) == 0 {
		return nil, fmt.Errorf("no series found with ID %s", seriesID)
	}

	s := apiResp.Seriess[0]
	return &Series{
		ID:    s.ID,
		Title: s.Title,
		Units: s.Units,
	}, nil
}

// fetchObservations retrieves all observations for a series
func (c *FREDClient) fetchObservations(seriesID string) ([]Observation, error) {
	url := fmt.Sprintf("%s/series/observations?series_id=%s&api_key=%s&file_type=json",
		fredBaseURL, seriesID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp fredObservationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	observations := make([]Observation, 0, len(apiResp.Observations))
	for _, obs := range apiResp.Observations {
		// Skip observations with "." as value (missing data)
		if obs.Value != "." {
			observations = append(observations, Observation{
				Date:  obs.Date,
				Value: obs.Value,
			})
		}
	}

	return observations, nil
}

// GetSourceName returns the name of this data source
func (c *FREDClient) GetSourceName() string {
	return "FRED"
}
