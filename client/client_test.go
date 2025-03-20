package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchCountry(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/name/india" {
			t.Errorf("Expected path /name/india, got %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[
			{
				"name": {
					"common": "India",
					"official": "Republic of India"
				},
				"capital": ["New Delhi"],
				"population": 1380004385,
				"currencies": {
					"INR": {
						"name": "Indian rupee",
						"symbol": "₹"
					}
				}
			}
		]`)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewCountryClient(server.URL, 5*time.Second)

	// Test search
	ctx := context.Background()
	resp, err := client.SearchCountry(ctx, "india")

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check response
	if len(resp) != 1 {
		t.Errorf("Expected 1 country in response, got %d", len(resp))
	}

	// Check country data
	country := resp[0]
	if country.Name.Common != "India" {
		t.Errorf("Expected name 'India', got '%s'", country.Name.Common)
	}

	if len(country.Capital) != 1 || country.Capital[0] != "New Delhi" {
		t.Errorf("Expected capital 'New Delhi', got '%v'", country.Capital)
	}

	if country.Population != 1380004385 {
		t.Errorf("Expected population 1380004385, got %d", country.Population)
	}

	// Check currency
	currency, ok := country.Currencies["INR"]
	if !ok {
		t.Errorf("Expected currency INR, but it was not found")
	}

	if currency.Symbol != "₹" {
		t.Errorf("Expected currency symbol '₹', got '%s'", currency.Symbol)
	}
}

func TestSearchCountryError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"message": "Country not found"}`)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewCountryClient(server.URL, 5*time.Second)

	// Test search
	ctx := context.Background()
	_, err := client.SearchCountry(ctx, "nonexistent")

	// Should have an error
	if err == nil {
		t.Errorf("Expected an error for non-existent country, but got nil")
	}
}

func TestSearchCountryTimeout(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"name":{"common":"Test"}}]`)
	}))
	defer server.Close()

	// Create client with very short timeout
	client := NewCountryClient(server.URL, 10*time.Millisecond)

	// Test search with timeout
	ctx := context.Background()
	_, err := client.SearchCountry(ctx, "test")

	// Should have a timeout error
	if err == nil {
		t.Errorf("Expected timeout error, but got nil")
	}
}
