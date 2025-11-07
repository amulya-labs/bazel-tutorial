package main

import (
	"net/http"
	"testing"
)

// Note: mockHTTPClient and newMockResponse are now in test_helpers.go

func TestFREDClient_FetchSeriesMetadata(t *testing.T) {
	mockClient := &mockHTTPClient{
		responses: map[string]*http.Response{
			"https://api.stlouisfed.org/fred/series?series_id=CPIAUCSL&api_key=test&file_type=json": newMockResponse(200, `{
				"seriess": [{
					"id": "CPIAUCSL",
					"title": "Consumer Price Index for All Urban Consumers",
					"units": "Index 1982-1984=100"
				}]
			}`),
		},
	}

	client := NewFREDClientWithHTTP("test", mockClient)
	series, err := client.fetchSeriesMetadata("CPIAUCSL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if series.ID != "CPIAUCSL" {
		t.Errorf("Expected ID CPIAUCSL, got %s", series.ID)
	}

	if series.Title != "Consumer Price Index for All Urban Consumers" {
		t.Errorf("Expected correct title, got %s", series.Title)
	}
}

func TestFREDClient_FetchObservations(t *testing.T) {
	mockClient := &mockHTTPClient{
		responses: map[string]*http.Response{
			"https://api.stlouisfed.org/fred/series/observations?series_id=CPIAUCSL&api_key=test&file_type=json": newMockResponse(200, `{
				"observations": [
					{"date": "2023-01-01", "value": "299.170"},
					{"date": "2023-02-01", "value": "300.840"},
					{"date": "2023-03-01", "value": "."}
				]
			}`),
		},
	}

	client := NewFREDClientWithHTTP("test", mockClient)
	observations, err := client.fetchObservations("CPIAUCSL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 2 observations (third one has "." and should be skipped)
	if len(observations) != 2 {
		t.Errorf("Expected 2 observations, got %d", len(observations))
	}

	if observations[0].Date != "2023-01-01" {
		t.Errorf("Expected first date 2023-01-01, got %s", observations[0].Date)
	}

	if observations[0].Value != "299.170" {
		t.Errorf("Expected first value 299.170, got %s", observations[0].Value)
	}
}

func TestFREDClient_FetchSeries(t *testing.T) {
	mockClient := &mockHTTPClient{
		responses: map[string]*http.Response{
			"https://api.stlouisfed.org/fred/series?series_id=UNRATE&api_key=test&file_type=json": newMockResponse(200, `{
				"seriess": [{
					"id": "UNRATE",
					"title": "Unemployment Rate",
					"units": "Percent"
				}]
			}`),
			"https://api.stlouisfed.org/fred/series/observations?series_id=UNRATE&api_key=test&file_type=json": newMockResponse(200, `{
				"observations": [
					{"date": "2023-01-01", "value": "3.5"},
					{"date": "2023-02-01", "value": "3.6"}
				]
			}`),
		},
	}

	client := NewFREDClientWithHTTP("test", mockClient)
	series, observations, err := client.FetchSeries("UNRATE")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if series.ID != "UNRATE" {
		t.Errorf("Expected series ID UNRATE, got %s", series.ID)
	}

	if len(observations) != 2 {
		t.Errorf("Expected 2 observations, got %d", len(observations))
	}
}

func TestFREDClient_ErrorHandling(t *testing.T) {
	mockClient := &mockHTTPClient{
		responses: map[string]*http.Response{
			"https://api.stlouisfed.org/fred/series?series_id=INVALID&api_key=test&file_type=json": newMockResponse(400, `{
				"error_code": 400,
				"error_message": "Bad Request"
			}`),
		},
	}

	client := NewFREDClientWithHTTP("test", mockClient)
	_, err := client.fetchSeriesMetadata("INVALID")

	if err == nil {
		t.Error("Expected error for invalid series, got nil")
	}
}
