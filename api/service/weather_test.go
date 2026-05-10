package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func weatherMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"current": map[string]any{
				"temperature_2m":       22.5,
				"apparent_temperature": 21.0,
				"relative_humidity_2m": 68,
				"wind_speed_10m":       12.3,
				"weather_code":         1,
				"uv_index":             4.0,
			},
		})
	}))
}

func forecastMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"daily": map[string]any{
				"time":                          []string{"2026-05-10", "2026-05-11"},
				"weather_code":                  []int{0, 61},
				"temperature_2m_max":            []float64{25.0, 18.0},
				"temperature_2m_min":            []float64{14.0, 12.0},
				"precipitation_sum":             []float64{0.0, 8.4},
				"precipitation_probability_max": []int{5, 85},
			},
		})
	}))
}

func TestGetCurrentWeather(t *testing.T) {
	srv := weatherMockServer()
	defer srv.Close()
	OpenMeteoWeatherURL = srv.URL

	geo := &GeoResult{Name: "São Paulo", Latitude: -23.5, Longitude: -46.6}
	w, err := GetCurrentWeather(geo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.City != "São Paulo" {
		t.Fatalf("expected São Paulo, got %s", w.City)
	}
	if w.TemperatureC != 22.5 {
		t.Fatalf("expected 22.5, got %f", w.TemperatureC)
	}
	if w.Condition != "Mainly clear" {
		t.Fatalf("expected 'Mainly clear', got %s", w.Condition)
	}
}

func TestGetForecast(t *testing.T) {
	srv := forecastMockServer()
	defer srv.Close()
	OpenMeteoWeatherURL = srv.URL

	geo := &GeoResult{Name: "Curitiba", Latitude: -25.4, Longitude: -49.3}
	f, err := GetForecast(geo, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Days) != 2 {
		t.Fatalf("expected 2 days, got %d", len(f.Days))
	}
	if f.Days[1].RainProbabilityPct != 85 {
		t.Fatalf("expected 85%%, got %d", f.Days[1].RainProbabilityPct)
	}
}
