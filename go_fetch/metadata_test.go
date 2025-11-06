package main

import (
	"fmt"
	"testing"
)

func TestIndicatorMetadata(t *testing.T) {
	metadata := GetIndicatorMetadata()

	// Verify we have metadata for all 14 series
	expectedSeries := []string{
		"GDPC1", "CPIAUCSL", "CPILFESL", "UNRATE", "FEDFUNDS",
		"T10Y2Y", "DGS10", "DGS2", "M2SL", "POILBREUSDM",
		"MANEMP", "UMCSENT", "IMPCH", "EXPCH",
	}

	for _, code := range expectedSeries {
		meta, exists := metadata[code]
		if !exists {
			t.Errorf("Missing metadata for series: %s", code)
			continue
		}

		// Verify all required fields are populated
		if meta.Code == "" {
			t.Errorf("Series %s missing Code", code)
		}
		if meta.Name == "" {
			t.Errorf("Series %s missing Name", code)
		}
		if meta.Unit == "" {
			t.Errorf("Series %s missing Unit", code)
		}
		if meta.Frequency == "" {
			t.Errorf("Series %s missing Frequency", code)
		}
		if meta.Category == "" {
			t.Errorf("Series %s missing Category", code)
		}
		if meta.ProxyFor == "" {
			t.Errorf("Series %s missing ProxyFor", code)
		}
		if meta.StartYear == 0 {
			t.Errorf("Series %s missing StartYear", code)
		}

		// Verify 30-year baseline (data since 1990 or earlier)
		if meta.StartYear > 1990 {
			t.Errorf("Series %s does not meet 30-year baseline: starts in %d", code, meta.StartYear)
		}
	}

	t.Logf("✅ All %d indicators have complete metadata", len(expectedSeries))
}

func TestIndicatorCategories(t *testing.T) {
	categories := GetIndicatorCategories()

	expectedCategories := []string{
		"Growth", "Inflation", "Labor Market", "Monetary Policy",
		"Yield Curve", "Liquidity", "Commodities", "Manufacturing",
		"Sentiment", "Trade",
	}

	for _, cat := range expectedCategories {
		indicators, exists := categories[cat]
		if !exists {
			t.Errorf("Missing category: %s", cat)
			continue
		}
		if len(indicators) == 0 {
			t.Errorf("Category %s has no indicators", cat)
		}
	}

	t.Logf("✅ All %d categories defined", len(expectedCategories))
}

func TestCoreIndicators(t *testing.T) {
	coreIndicators := GetCoreIndicators()

	// Should have exactly 10 core indicators (one per category)
	if len(coreIndicators) != 10 {
		t.Errorf("Expected 10 core indicators, got %d", len(coreIndicators))
	}

	// Verify all core indicators exist in metadata
	metadata := GetIndicatorMetadata()
	for _, code := range coreIndicators {
		if _, exists := metadata[code]; !exists {
			t.Errorf("Core indicator %s missing from metadata", code)
		}
	}

	t.Logf("✅ %d core indicators defined", len(coreIndicators))
}

func TestDefaultSeries(t *testing.T) {
	// Verify defaultSeries matches what we expect
	if len(defaultSeries) != 14 {
		t.Errorf("Expected 14 series in defaultSeries, got %d", len(defaultSeries))
	}

	// Verify all series in defaultSeries have metadata
	metadata := GetIndicatorMetadata()
	missingMetadata := []string{}

	for _, code := range defaultSeries {
		if _, exists := metadata[code]; !exists {
			missingMetadata = append(missingMetadata, code)
		}
	}

	if len(missingMetadata) > 0 {
		t.Errorf("Series in defaultSeries missing metadata: %v", missingMetadata)
	}

	t.Logf("✅ All %d default series have metadata", len(defaultSeries))
}

func TestMetadataConsistency(t *testing.T) {
	metadata := GetIndicatorMetadata()
	categories := GetIndicatorCategories()

	// Verify all indicators in categories exist in metadata
	for category, indicators := range categories {
		for _, code := range indicators {
			if _, exists := metadata[code]; !exists {
				t.Errorf("Indicator %s in category %s missing from metadata", code, category)
			}
		}
	}

	// Verify all indicators in metadata belong to a category
	allCategorized := make(map[string]bool)
	for _, indicators := range categories {
		for _, code := range indicators {
			allCategorized[code] = true
		}
	}

	for code, meta := range metadata {
		if !allCategorized[code] {
			t.Errorf("Indicator %s (%s) not assigned to any category", code, meta.Name)
		}
	}

	t.Log("✅ Metadata and categories are consistent")
}

// Example test output helper
func ExampleGetIndicatorMetadata() {
	metadata := GetIndicatorMetadata()
	gdp := metadata["GDPC1"]
	fmt.Printf("%s: %s\n", gdp.Code, gdp.Name)
	fmt.Printf("Category: %s\n", gdp.Category)
	fmt.Printf("Proxy For: %s\n", gdp.ProxyFor)
	// Output:
	// GDPC1: Real Gross Domestic Product
	// Category: Growth
	// Proxy For: Economic Growth / Business Cycle
}
