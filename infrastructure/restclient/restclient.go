package restclient

import (
	"errors"
	"net/http"
	"time"
)

// Config represents the configuration for the HTTP client
type Config struct {
	BaseURL string        // Base URL for the HTTP client
	Timeout time.Duration // Timeout for the HTTP client
	// TODO: Add configuration fields
}

// NewRestClient creates a new HTTP client
// TODO: Implement the logic to create a new HTTP client
func NewRestClient(config *Config) (*http.Client, error) {
	if config == nil {
		return nil, errors.New("http client: config cannot be nil")
	}

	// TODO: Implement the logic to create a new HTTP client

	return &http.Client{}, nil
}
