package http

import "errors"

// HTTPClient represents an HTTP client
type HTTPClient struct {
	config interface{}
}

// NewHTTPClient creates a new HTTP client
// TODO: Implement the logic to create a new HTTP client
func NewHTTPClient(config interface{}) (*HTTPClient, error) {
	if config == nil {
		return nil, errors.New("http client: config cannot be nil")
	}
	return &HTTPClient{config: config}, nil
}
