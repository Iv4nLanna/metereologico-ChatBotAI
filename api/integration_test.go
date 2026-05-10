//go:build integration

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestIntegration_Health(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/health", baseURL))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_Weather(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/weather?city=São Paulo", baseURL))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal("invalid JSON")
	}
	if body["city"] == nil || body["temperature_c"] == nil {
		t.Fatal("missing required fields")
	}
}

func TestIntegration_Forecast(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/forecast?city=Curitiba&days=3", baseURL))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		City string `json:"city"`
		Days []any  `json:"days"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal("invalid JSON")
	}
	if len(body.Days) != 3 {
		t.Fatalf("expected 3 days, got %d", len(body.Days))
	}
}

func TestIntegration_CityNotFound(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/weather?city=CidadeQueNaoExiste999", baseURL))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
