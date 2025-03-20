package models

// CountryResponse represents the API response for country data
type CountryResponse struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int    `json:"population"`
}

// RestCountriesResponse represents the response from REST Countries API
type RestCountriesResponse []struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Capital    []string `json:"capital"`
	Population int      `json:"population"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

// Error represents an API error response
type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
