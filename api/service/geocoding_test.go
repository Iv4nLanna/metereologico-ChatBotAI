package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocode_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"name": "Curitiba", "latitude": -25.4297, "longitude": -49.2711},
			},
		})
	}))
	defer srv.Close()
	OpenMeteoGeocodingURL = srv.URL

	got, err := Geocode("Curitiba")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Curitiba" {
		t.Fatalf("expected Curitiba, got %s", got.Name)
	}
}

func TestGeocode_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer srv.Close()
	OpenMeteoGeocodingURL = srv.URL

	_, err := Geocode("CidadeInexistente999")
	if err == nil {
		t.Fatal("expected error for city not found")
	}
}
