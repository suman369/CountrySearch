package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"countrysearch/models"
)

// CountryClient is a custom HTTP client for REST Countries API
type CountryClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCountryClient creates a new HTTP client for REST Countries API
func NewCountryClient(baseURL string, timeout time.Duration) *CountryClient {
	return &CountryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SearchCountry searches for a country by name
func (c *CountryClient) SearchCountry(ctx context.Context, name string) (models.RestCountriesResponse, error) {
	// Build URL
	endpoint := fmt.Sprintf("%s/name/%s", c.baseURL, url.QueryEscape(name))

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	// Parse response
	var countries models.RestCountriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return countries, nil
}
