package main

import (
	"bytes"
	"io"
	"net/http"
)

// mockHTTPClient is a mock HTTP client for testing
type mockHTTPClient struct {
	responses map[string]*http.Response
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if resp, ok := m.responses[url]; ok {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString("Not found")),
	}, nil
}

// newMockResponse creates a mock HTTP response
func newMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
