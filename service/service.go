package service

import (
	"context"
	"fmt"
	"time"

	"countrysearch/cache"
	"countrysearch/client"
	"countrysearch/models"
	"countrysearch/utils"
)

// CountryService handles business logic for country data
type CountryService struct {
	client *client.CountryClient
	cache  cache.Cache
	logger utils.Logger
}

// NewCountryService creates a new country service
func NewCountryService(client *client.CountryClient, cache cache.Cache, logger utils.Logger) *CountryService {
	return &CountryService{
		client: client,
		cache:  cache,
		logger: logger,
	}
}

// SearchCountry searches for a country by name
func (s *CountryService) SearchCountry(ctx context.Context, name string) (*models.CountryResponse, error) {
	s.logger.Info("Searching for country", "name", name)

	// Check cache first
	cacheKey := fmt.Sprintf("country:%s", name)
	if cachedData, found := s.cache.Get(cacheKey); found {
		s.logger.Info("Country found in cache", "name", name)
		return cachedData.(*models.CountryResponse), nil
	}

	// If not in cache, call API
	s.logger.Info("Country not found in cache, calling API", "name", name)
	countries, err := s.client.SearchCountry(ctx, name)
	if err != nil {
		s.logger.Error("Error fetching country from API", "error", err)
		return nil, err
	}

	// No results found
	if len(countries) == 0 {
		s.logger.Info("No country found with name", "name", name)
		return nil, fmt.Errorf("no country found with name: %s", name)
	}

	// Process first result (most relevant match)
	country := countries[0]

	// Extract currency information
	var currencySymbol string
	for _, currency := range country.Currencies {
		currencySymbol = currency.Symbol
		break // Just use the first currency
	}

	// Extract capital
	var capital string
	if len(country.Capital) > 0 {
		capital = country.Capital[0]
	}

	// Create response model
	response := &models.CountryResponse{
		Name:       country.Name.Common,
		Capital:    capital,
		Currency:   currencySymbol,
		Population: country.Population,
	}

	// Store in cache for 1 hour
	s.cache.(*cache.InMemoryCache).SetWithExpiration(cacheKey, response, 1*time.Hour)
	s.logger.Info("Country data stored in cache", "name", name)

	return response, nil
}
