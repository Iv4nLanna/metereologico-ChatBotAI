package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"weather-api/cache"
)

func TestWeather_MissingCity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := cache.New(time.Minute)
	r.GET("/weather", Weather(c))

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWeather_CacheHit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := cache.New(time.Minute)

	c.Set("weather:TestCity", map[string]any{"city": "TestCity", "temperature_c": 25.0})

	r.GET("/weather", Weather(c))

	req := httptest.NewRequest(http.MethodGet, "/weather?city=TestCity", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body["city"] != "TestCity" {
		t.Fatalf("expected city 'TestCity', got %v", body["city"])
	}
}
