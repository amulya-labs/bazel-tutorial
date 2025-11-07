package main

import (
	"net/http"
	"testing"
)

func TestWorldBankClient_ParseSeriesID(t *testing.T) {
	tests := []struct {
		name          string
		seriesID      string
		wantIndicator string
		wantCountry   string
		wantErr       bool
		mockResponse  string
	}{
		{
			name:          "Valid World GDP",
			seriesID:      "WB:NY.GDP.MKTP.KD:WLD",
			wantIndicator: "NY.GDP.MKTP.KD",
			wantCountry:   "WLD",
			wantErr:       false,
			mockResponse: `[
				{"page": 1, "pages": 1},
				[{"indicator": {"id": "NY.GDP.MKTP.KD"}, "country": {"id": "WLD"}, "date": "2023", "value": 100000}]
			]`,
		},
		{
			name:          "Valid China GDP",
			seriesID:      "WB:NY.GDP.MKTP.KD:CHN",
			wantIndicator: "NY.GDP.MKTP.KD",
			wantCountry:   "CHN",
			wantErr:       false,
			mockResponse: `[
				{"page": 1, "pages": 1},
				[{"indicator": {"id": "NY.GDP.MKTP.KD"}, "country": {"id": "CHN"}, "date": "2023", "value": 100000}]
			]`,
		},
		{
			name:     "Invalid format - missing parts",
			seriesID: "WB:NY.GDP.MKTP.KD",
			wantErr:  true,
		},
		{
			name:     "Invalid format - too many parts",
			seriesID: "WB:NY.GDP.MKTP.KD:USA:EXTRA",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				responses: map[string]*http.Response{
					"https://api.worldbank.org/v2/country/WLD/indicator/NY.GDP.MKTP.KD?format=json&date=1990:2024&per_page=1000": newMockResponse(200, tt.mockResponse),
					"https://api.worldbank.org/v2/country/CHN/indicator/NY.GDP.MKTP.KD?format=json&date=1990:2024&per_page=1000": newMockResponse(200, tt.mockResponse),
				},
			}
			client := NewWorldBankClientWithHTTP(mockClient)

			// We'll test parsing by trying to fetch
			_, _, err := client.FetchSeries(tt.seriesID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchSeries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorldBankClient_FetchSeries(t *testing.T) {
	mockResponse := `[
		{"page": 1, "pages": 1, "per_page": 50, "total": 3},
		[
			{
				"indicator": {"id": "NY.GDP.MKTP.KD", "value": "GDP (constant 2015 US$)"},
				"country": {"id": "WLD", "value": "World"},
				"countryiso3code": "WLD",
				"date": "2023",
				"value": 100000000000000,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			},
			{
				"indicator": {"id": "NY.GDP.MKTP.KD", "value": "GDP (constant 2015 US$)"},
				"country": {"id": "WLD", "value": "World"},
				"countryiso3code": "WLD",
				"date": "2022",
				"value": 95000000000000,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			},
			{
				"indicator": {"id": "NY.GDP.MKTP.KD", "value": "GDP (constant 2015 US$)"},
				"country": {"id": "WLD", "value": "World"},
				"countryiso3code": "WLD",
				"date": "2021",
				"value": null,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			}
		]
	]`

	mockClient := &mockHTTPClient{
		responses: map[string]*http.Response{
			"https://api.worldbank.org/v2/country/WLD/indicator/NY.GDP.MKTP.KD?format=json&date=1990:2024&per_page=1000": newMockResponse(200, mockResponse),
		},
	}

	client := NewWorldBankClientWithHTTP(mockClient)
	series, observations, err := client.FetchSeries("WB:NY.GDP.MKTP.KD:WLD")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if series.ID != "WB:NY.GDP.MKTP.KD:WLD" {
		t.Errorf("Expected series ID WB:NY.GDP.MKTP.KD:WLD, got %s", series.ID)
	}

	if series.Title != "Real GDP (WLD)" {
		t.Errorf("Expected title 'Real GDP (WLD)', got %s", series.Title)
	}

	// Should have 2 observations (2023, 2022), null value should be skipped
	if len(observations) != 2 {
		t.Errorf("Expected 2 observations, got %d", len(observations))
	}

	// Check first observation
	if observations[0].Date != "2023-01-01" {
		t.Errorf("Expected date 2023-01-01, got %s", observations[0].Date)
	}

	// Check value format (should be a number as string)
	if observations[0].Value == "" {
		t.Error("Expected non-empty value")
	}
}

func TestWorldBankClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		seriesID   string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "API Error - 404",
			seriesID:   "WB:INVALID.INDICATOR:USA",
			statusCode: 404,
			response:   `{"message": "Not found"}`,
			wantErr:    true,
		},
		{
			name:       "Invalid JSON Response",
			seriesID:   "WB:NY.GDP.MKTP.KD:USA",
			statusCode: 200,
			response:   `invalid json`,
			wantErr:    true,
		},
		{
			name:       "Invalid Series ID Format",
			seriesID:   "INVALID",
			statusCode: 200,
			response:   `[]`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				responses: map[string]*http.Response{
					"https://api.worldbank.org/v2/country/USA/indicator/INVALID.INDICATOR?format=json&date=1990:2024&per_page=1000": newMockResponse(tt.statusCode, tt.response),
					"https://api.worldbank.org/v2/country/USA/indicator/NY.GDP.MKTP.KD?format=json&date=1990:2024&per_page=1000":    newMockResponse(tt.statusCode, tt.response),
				},
			}

			client := NewWorldBankClientWithHTTP(mockClient)
			_, _, err := client.FetchSeries(tt.seriesID)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestWorldBankClient_GetSourceName(t *testing.T) {
	client := NewWorldBankClient()
	if client.GetSourceName() != "World Bank" {
		t.Errorf("Expected source name 'World Bank', got %s", client.GetSourceName())
	}
}

func TestDataSourceFactory(t *testing.T) {
	tests := []struct {
		name       string
		seriesID   string
		fredAPIKey string
		wantSource string
		wantErr    bool
	}{
		{
			name:       "FRED series",
			seriesID:   "GDPC1",
			fredAPIKey: "test_key",
			wantSource: "FRED",
			wantErr:    false,
		},
		{
			name:       "World Bank series",
			seriesID:   "WB:NY.GDP.MKTP.KD:USA",
			fredAPIKey: "test_key",
			wantSource: "World Bank",
			wantErr:    false,
		},
		{
			name:       "OECD series - not implemented",
			seriesID:   "OECD:MEI.LRUNTTTT.USA.M",
			fredAPIKey: "test_key",
			wantSource: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, _, err := DataSourceFactory(tt.seriesID, tt.fredAPIKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("DataSourceFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && source.GetSourceName() != tt.wantSource {
				t.Errorf("DataSourceFactory() source = %s, want %s", source.GetSourceName(), tt.wantSource)
			}
		})
	}
}
