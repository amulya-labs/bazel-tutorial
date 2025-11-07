package main

import (
	"fmt"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// DefaultRetryConfig returns sensible defaults for API retries
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
	}
}

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(config RetryConfig, operation func() error) error {
	var lastErr error
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			time.Sleep(backoff)

			// Increase backoff exponentially
			backoff = time.Duration(float64(backoff) * config.BackoffFactor)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}

		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on certain errors (e.g., invalid series ID)
		if !isRetryableError(err) {
			return err
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", config.MaxRetries, lastErr)
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Don't retry on client errors (4xx)
	if contains(errStr, "400") || contains(errStr, "404") || contains(errStr, "invalid") {
		return false
	}

	// Retry on network errors and server errors (5xx)
	return contains(errStr, "timeout") ||
		contains(errStr, "connection") ||
		contains(errStr, "500") ||
		contains(errStr, "502") ||
		contains(errStr, "503")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
