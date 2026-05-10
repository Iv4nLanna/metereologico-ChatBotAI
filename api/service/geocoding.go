package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var OpenMeteoGeocodingURL = "https://geocoding-api.open-meteo.com"

var geocodingClient = &http.Client{Timeout: 10 * time.Second}

type GeoResult struct {
	Latitude  float64
	Longitude float64
	Name      string
}

func Geocode(city string) (*GeoResult, error) {
	u := fmt.Sprintf(
		"%s/v1/search?name=%s&count=1&language=en&format=json",
		OpenMeteoGeocodingURL, url.QueryEscape(city),
	)
	resp, err := geocodingClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var result struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("geocoding decode failed: %w", err)
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}
	r := result.Results[0]
	return &GeoResult{Latitude: r.Latitude, Longitude: r.Longitude, Name: r.Name}, nil
}
